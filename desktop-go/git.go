package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitCommit represents a single commit entry
type GitCommit struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

// GitFileChange represents a file changed in a commit
type GitFileChange struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Patch     string `json:"patch"`
}

// GitCommitDetail holds full details of a commit
type GitCommitDetail struct {
	Hash    string          `json:"hash"`
	Author  string          `json:"author"`
	Date    string          `json:"date"`
	Message string          `json:"message"`
	Files   []GitFileChange `json:"files"`
}

// SearchResult represents a single search match
type SearchResult struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

// ReplaceResult represents the result of a replace operation on one file
type ReplaceResult struct {
	File    string `json:"file"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ── helpers ──────────────────────────────────────────────────────────────────

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git error: %s", string(exitErr.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func findGitRoot(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	dir := abs
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .git directory found")
		}
		dir = parent
	}
}

var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
	".svg": true, ".woff": true, ".woff2": true, ".eot": true, ".ttf": true,
	".otf": true, ".pdf": true, ".zip": true, ".tar": true, ".gz": true,
	".exe": true, ".dll": true, ".so": true, ".dylib": true, ".bin": true,
	".o": true, ".a": true, ".lib": true, ".obj": true, ".class": true,
	".jar": true, ".war": true, ".pyc": true, ".pyo": true,
	".icns": true, ".webp": true, ".avif": true,
	".mp3": true, ".mp4": true, ".avi": true, ".mov": true, ".wav": true,
	".flac": true, ".ogg": true, ".wasm": true,
}

func shouldSkipDir(name string) bool {
	skipped := map[string]bool{
		"node_modules": true, ".git": true, "dist": true, "build": true,
		".build": true, "__pycache__": true, "vendor": true, ".next": true,
		".nuxt": true, ".output": true, "target": true, ".cache": true,
	}
	return skipped[name] || (strings.HasPrefix(name, ".") && name != ".")
}

// ── Exported methods ─────────────────────────────────────────────────────────

// GetGitLog returns the commit history for the workspace
func (a *App) GetGitLog(path string) ([]GitCommit, error) {
	if path == "" {
		path = a.agent.ActivePath
	}
	gitRoot, err := findGitRoot(path)
	if err != nil {
		return nil, fmt.Errorf("not a git repository")
	}

	const sep = "\x1F"
	format := fmt.Sprintf("%%H%s%%an%s%%ai%s%%s", sep, sep, sep)
	out, err := runGit(gitRoot, "log", "--pretty=format:"+format, "--max-count=50")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return []GitCommit{}, nil
	}

	lines := strings.Split(out, "\n")
	commits := make([]GitCommit, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, sep, 4)
		if len(parts) < 4 {
			continue
		}
		commits = append(commits, GitCommit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		})
	}
	return commits, nil
}

// GetGitCommitDetail returns full details of a specific commit
func (a *App) GetGitCommitDetail(path string, hash string) (*GitCommitDetail, error) {
	if path == "" {
		path = a.agent.ActivePath
	}
	gitRoot, err := findGitRoot(path)
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %v", err)
	}

	const sep = "\x1F"
	metaFormat := fmt.Sprintf("%%H%s%%an%%s%%ai%s%%s", sep, sep)
	meta, err := runGit(gitRoot, "show", "--quiet", "--pretty=format:"+metaFormat, hash)
	if err != nil {
		return nil, err
	}
	metaParts := strings.SplitN(meta, sep, 3)
	if len(metaParts) < 3 {
		return nil, fmt.Errorf("unable to parse commit metadata")
	}

	// Commit message is everything after the date on the same line – re-parse properly
	const fullSep = "\x1E"
	fullFormat := fmt.Sprintf("%%H%s%%an%s%%ai%s%%s", fullSep, fullSep, fullSep)
	fullMeta, err := runGit(gitRoot, "show", "--quiet", "--pretty=format:"+fullFormat, hash)
	if err != nil {
		return nil, err
	}
	fullParts := strings.SplitN(fullMeta, fullSep, 4)
	hashVal, authorVal, dateVal, messageVal := "", "", "", ""
	if len(fullParts) >= 4 {
		hashVal = fullParts[0]
		authorVal = fullParts[1]
		dateVal = fullParts[2]
		messageVal = fullParts[3]
	} else if len(fullParts) >= 3 {
		hashVal = fullParts[0]
		authorVal = fullParts[1]
		dateVal = fullParts[2]
	}

	detail := &GitCommitDetail{
		Hash:    hashVal,
		Author:  authorVal,
		Date:    dateVal,
		Message: messageVal,
	}

	// Get changed files
	filesOut, err := runGit(gitRoot, "diff-tree", "--no-commit-id", "--name-status", "-r", "-M", hash)
	if err != nil {
		return detail, nil
	}

	var files []GitFileChange
	if filesOut != "" {
		for _, line := range strings.Split(filesOut, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) < 2 {
				continue
			}
			statusChar := parts[0]
			filePath := parts[1]
			status := "modified"
			if len(statusChar) > 0 {
				switch statusChar[0] {
				case 'A':
					status = "added"
				case 'M':
					status = "modified"
				case 'D':
					status = "deleted"
				case 'R':
					status = "renamed"
				case 'C':
					status = "copied"
				}
			}
			files = append(files, GitFileChange{
				Path:   filePath,
				Status: status,
			})
		}
	}

	// Get diff
	diffOut, err := runGit(gitRoot, "diff", hash+"^!", "--unified=3")
	if err == nil && diffOut != "" {
		fileMap := make(map[string]*GitFileChange)
		for i := range files {
			fileMap[files[i].Path] = &files[i]
		}

		diffFiles := strings.Split(diffOut, "\ndiff --git ")
		for _, df := range diffFiles {
			if df == "" {
				continue
			}
			if !strings.HasPrefix(df, "a/") && !strings.HasPrefix(df, "---") {
				df = "diff --git " + df
			}

			scanner := bufio.NewScanner(strings.NewReader(df))
			var filePath string
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "+++ b/") {
					filePath = strings.TrimPrefix(line, "+++ b/")
					break
				}
			}
			if filePath == "" {
				continue
			}

			if strings.Contains(df, "Binary files") {
				if fc, ok := fileMap[filePath]; ok {
					fc.Patch = "[Binary file]"
				}
				continue
			}

			if fc, ok := fileMap[filePath]; ok {
				fc.Patch = df
			} else {
				files = append(files, GitFileChange{
					Path:   filePath,
					Status: "modified",
					Patch:  df,
				})
			}
		}
	}

	detail.Files = files
	return detail, nil
}

// SearchInFiles searches for text in all files within the workspace directory
func (a *App) SearchInFiles(path string, query string) ([]SearchResult, error) {
	if path == "" {
		path = a.agent.ActivePath
	}
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	var results []SearchResult
	queryLower := strings.ToLower(query)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		if binaryExts[ext] {
			return nil
		}
		if info.Size() > 1_000_000 {
			return nil
		}

		f, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*64), 1024*64)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), queryLower) {
				relPath, _ := filepath.Rel(path, filePath)
				results = append(results, SearchResult{
					File:    relPath,
					Line:    lineNum,
					Content: strings.TrimSpace(line),
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return results, nil
}

// ReplaceInFiles replaces text in files matching the search query
func (a *App) ReplaceInFiles(path string, searchText string, replaceText string) ([]ReplaceResult, error) {
	if path == "" {
		path = a.agent.ActivePath
	}
	if searchText == "" {
		return nil, fmt.Errorf("search text cannot be empty")
	}

	var results []ReplaceResult

	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		if binaryExts[ext] {
			return nil
		}
		if info.Size() > 1_000_000 {
			return nil
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		text := string(content)
		if !strings.Contains(text, searchText) {
			return nil
		}

		newText := strings.ReplaceAll(text, searchText, replaceText)
		if newText == text {
			return nil
		}

		err = os.WriteFile(filePath, []byte(newText), info.Mode())
		relPath, _ := filepath.Rel(path, filePath)
		if err != nil {
			results = append(results, ReplaceResult{
				File:    relPath,
				Success: false,
				Error:   err.Error(),
			})
		} else {
			results = append(results, ReplaceResult{
				File:    relPath,
				Success: true,
			})
		}
		return nil
	})

	return results, nil
}
