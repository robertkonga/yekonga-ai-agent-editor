package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"yekonga-builder/console"
	"yekonga-builder/helper"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileNode represents the metadata schema for files and directories
type FileNode struct {
	ID        string      `json:"id"` // Masked, deterministic unique identifier
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Extension string      `json:"extension,omitempty"` // Omitted if directory
	Type      string      `json:"type"`                // "file" or "directory"
	Lang      string      `json:"lang,omitempty"`      // Language key for editor highlighting
	Extended  bool        `json:"extended"`            // UI state default
	Children  []*FileNode `json:"children,omitempty"`  // Omitted if file
}

// OpenWorkspaceDialog triggers a native folder prompt and returns the selected route string
func (a *App) OpenWorkspaceDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title:                "Select Project Workspace Root",
		CanCreateDirectories: true,
		// Filters out individual files so only directory folders are selectable
	}

	// This blocks thread execution until user clicks "Select" or "Cancel"
	selectedDirectory, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	return selectedDirectory, nil
}

func (a *App) SaveFile(data any, target string) error {
	err := helper.CreateFile(data, target)

	return err
}

func (a *App) ReadFile(target string) (string, error) {
	content := helper.ReadFile(target)

	return content, nil
}

func (a *App) MoveFile(source string, destination string) error {
	if err := os.Rename(source, destination); err != nil {
		// os.Rename fails across devices/partitions — fallback to copy+delete
		if err = a.CopyFile(source, destination); err != nil {
			return err
		}

		return os.Remove(source)
	}

	return nil
}

func (a *App) RenameFile(name string, source string) error {
	dir := filepath.Dir(source)
	dest := filepath.Join(dir, name)
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists: %s", dest)
	}

	return os.Rename(source, dest)
}

func (a *App) DeleteFile(target string) error {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", target)
	}

	return os.RemoveAll(target) // handles both files and dirs
}

// CopyFile is the cross-device fallback used by MoveFile
func (a *App) CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// ReadDirectory recursively scans the target path and returns a nested tree structure
func (a *App) ReadDirectory(rootPath string) (*FileNode, error) {

	// Convert to absolute path to guarantee uniqueness across execution contexts
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(absPath)
	_, err = os.Stat(cleanPath)
	if err != nil {
		return nil, err
	}

	// Bootstrap the root directory node
	rootNode := &FileNode{
		ID:       generateID(cleanPath),
		Name:     filepath.Base(cleanPath),
		Path:     cleanPath,
		Type:     "directory",
		Extended: true,
		Children: []*FileNode{},
	}

	err = buildTree(cleanPath, rootNode)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func (a *App) HomeDirectory(name string) string {
	dir, err := os.UserHomeDir()
	if err != nil {
		console.Error("HomeDirectory", err.Error())
	}

	if dir == "/" || helper.IsEmpty(dir) {
		dir = "/root"
	}

	appDir := dir + string(os.PathSeparator) + ".yekonga-server" + string(os.PathSeparator) + name

	if info, err := os.Stat(appDir); err != nil {
		if info != nil && !info.IsDir() {
			os.MkdirAll(appDir, 0755)
		} else {
			err := os.MkdirAll(appDir, 0755)
			if err != nil {
				console.Error("HomeDirectory", appDir, err)
			}
		}
	}

	return appDir
}

// generateID creates a stable, non-plain-text hash from the file path
func generateID(absolutePath string) string {
	// Standardize path separators (slashes) so IDs are identical across Windows/Linux/macOS hosts
	standardizedPath := filepath.ToSlash(absolutePath)

	hash := sha256.Sum256([]byte(standardizedPath))

	// Return a truncated 16-character hex string for clean frontend handling (or remove [:16] for full hash)
	return hex.EncodeToString(hash[:])[:16]
}

// buildTree helper handles the recursive traversal with IDE sorting rules
func buildTree(currentPath string, parentNode *FileNode) error {
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	// SORTING LOGIC: Directories first, then files, both alphabetically
	sort.Slice(entries, func(i, j int) bool {
		// If one is a directory and the other isn't, prioritize the directory
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir() // Returns true if 'i' is a directory, pushing it up
		}
		// If both are directories OR both are files, sort alphabetically by lowercase name
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	// Process the sorted entries
	for _, entry := range entries {
		entryPath := filepath.Join(currentPath, entry.Name())

		node := &FileNode{
			ID:   generateID(entryPath),
			Name: entry.Name(),
			Path: entryPath,
		}

		if entry.IsDir() {
			node.Type = "directory"
			node.Extended = false
			node.Children = []*FileNode{}

			// Recursive dive
			err := buildTree(entryPath, node)
			if err != nil {
				return err
			}
		} else {
			node.Type = "file"
			node.Extension = filepath.Ext(entry.Name())
			node.Lang = detectLanguage(node.Extension)
		}

		parentNode.Children = append(parentNode.Children, node)
	}

	return nil
}

// Language extension dictionary cache
var extensionMap = map[string]string{
	// Web & Frontend
	".js":   "javascript",
	".jsx":  "javascript",
	".ts":   "typescript",
	".tsx":  "typescript",
	".css":  "css",
	".scss": "scss",
	".less": "less",
	".json": "json",
	".html": "html",
	".htm":  "html",
	".vue":  "html", // Monaco fallback for SFC templates
	".svg":  "xml",

	// Backend & Mainstream Systems
	".go":    "go",
	".py":    "python",
	".pyw":   "python",
	".rs":    "rust",
	".java":  "java",
	".class": "java",
	".cpp":   "cpp",
	".cc":    "cpp",
	".cxx":   "cpp",
	".h":     "cpp",
	".hpp":   "cpp",
	".c":     "c",
	".cs":    "csharp",
	".rb":    "ruby",
	".php":   "php",
	".swift": "swift",
	".kt":    "kotlin",
	".kts":   "kotlin",

	// Shell, Scripting & Configs
	".sh":         "shell",
	".bash":       "shell",
	".zsh":        "shell",
	".ps1":        "powershell",
	".bat":        "bat",
	".cmd":        "bat",
	".yaml":       "yaml",
	".yml":        "yaml",
	".toml":       "toml",
	".ini":        "ini",
	".dockerfile": "dockerfile",

	// Data, Query, & Documents
	".sql": "sql",
	".r":   "r",
	".md":  "markdown",
	".xml": "xml",
	".csv": "plaintext",
}

// DetectLanguage matches file extensions to their Monaco code-highlighter target strings
func detectLanguage(ext string) string {
	// Standardize input string casing and dots
	cleanedExt := strings.ToLower(strings.TrimSpace(ext))
	if !strings.HasPrefix(cleanedExt, ".") && cleanedExt != "" {
		cleanedExt = "." + cleanedExt
	}

	if lang, verified := extensionMap[cleanedExt]; verified {
		return lang
	}

	return "plaintext"
}
