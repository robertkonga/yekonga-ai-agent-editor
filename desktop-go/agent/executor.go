package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ToolExecutor handles all tool executions with proper error handling and security
type ToolExecutor struct {
	Agent       *Agent
	MaxFileSize int64
	Timeout     time.Duration
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor(agent *Agent) *ToolExecutor {
	return &ToolExecutor{
		Agent:       agent,
		MaxFileSize: 10 * 1024 * 1024, // 10MB default
		Timeout:     30 * time.Second,
	}
}

// Execute executes a tool with the given name and raw input
func (e *ToolExecutor) Execute(name string, rawInput json.RawMessage) (string, error) {
	// Resolve paths relative to workspace
	resolvePath := func(p string) string {
		if p == "" {
			return ""
		}
		if filepath.IsAbs(p) {
			return p
		}
		if e.Agent.ActivePath != "" {
			return filepath.Join(e.Agent.ActivePath, p)
		}
		return p
	}

	// Validate path is within workspace (security)
	validatePath := func(p string) error {
		if p == "" {
			return fmt.Errorf("path cannot be empty")
		}
		if !filepath.IsAbs(p) {
			return fmt.Errorf("path must be absolute: %s", p)
		}
		if e.Agent.ActivePath != "" {
			cleanPath := filepath.Clean(p)
			cleanWorkspace := filepath.Clean(e.Agent.ActivePath)
			if !strings.HasPrefix(cleanPath, cleanWorkspace) {
				return fmt.Errorf("path escapes workspace: %s", p)
			}
		}
		return nil
	}

	switch name {
	case "read_file":
		return e.handleReadFile(rawInput, resolvePath, validatePath)
	case "read_file_range":
		return e.handleReadFileRange(rawInput, resolvePath, validatePath)
	case "list_directory":
		return e.handleListDirectory(rawInput, resolvePath, validatePath)
	case "list_directory_recursive":
		return e.handleListDirectoryRecursive(rawInput, resolvePath, validatePath)
	case "write_file":
		return e.handleWriteFile(rawInput, resolvePath, validatePath)
	case "patch_file":
		return e.handlePatchFile(rawInput, resolvePath, validatePath)
	case "delete_file":
		return e.handleDeleteFile(rawInput, resolvePath, validatePath)
	case "rename_file":
		return e.handleRenameFile(rawInput, resolvePath, validatePath)
	case "file_exists":
		return e.handleFileExists(rawInput, resolvePath, validatePath)
	case "summarize_file":
		return e.handleSummarizeFile(rawInput, resolvePath, validatePath)
	case "format_file":
		return e.handleFormatFile(rawInput, resolvePath, validatePath)
	case "create_directory":
		return e.handleCreateDirectory(rawInput, resolvePath, validatePath)
	case "search_in_files":
		return e.handleSearchInFiles(rawInput, resolvePath, validatePath)
	case "find_definition":
		return e.handleFindDefinition(rawInput, resolvePath, validatePath)
	case "find_usages":
		return e.handleFindUsages(rawInput, resolvePath, validatePath)
	case "get_file_tree":
		return e.handleGetFileTree(rawInput, resolvePath, validatePath)
	case "run_shell":
		return e.handleRunShell(rawInput, resolvePath)
	case "run_tests":
		return e.handleRunTests(rawInput, resolvePath)
	case "lint_file":
		return e.handleLintFile(rawInput, resolvePath, validatePath)
	case "git_status":
		return e.handleGitStatus(rawInput, resolvePath)
	case "git_diff":
		return e.handleGitDiff(rawInput, resolvePath)
	case "git_log":
		return e.handleGitLog(rawInput, resolvePath)
	case "git_commit":
		return e.handleGitCommit(rawInput, resolvePath)
	case "read_dependencies":
		return e.handleReadDependencies(rawInput, resolvePath, validatePath)
	case "install_dependency":
		return e.handleInstallDependency(rawInput, resolvePath)
	case "initialize_project":
		return e.handleInitializeProject(rawInput, resolvePath)
	case "generate_dataschema":
		return e.handleGenerateDataSchema(rawInput)
	case "generate_crud":
		return e.handleGenerateCRUD(rawInput, resolvePath, validatePath)
	case "read_template":
		return e.handleReadTemplate(rawInput)
	case "list_templates":
		return e.handleListTemplates()
	case "get_active_workspace_path":
		return e.handleGetActiveWorkspacePath()
	case "get_active_file":
		return e.handleGetActiveFile()
	case "get_workspace_config":
		return e.handleGetWorkspaceConfig(rawInput, resolvePath, validatePath)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// ── Filesystem handlers ──────────────────────────────────────────────────────

func (e *ToolExecutor) handleReadFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.Size() > e.MaxFileSize {
		return "", fmt.Errorf("file exceeds max size limit (%d bytes)", e.MaxFileSize)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *ToolExecutor) handleReadFileRange(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path      string `json:"path"`
		StartLine int    `json:"start_line"`
		EndLine   int    `json:"end_line"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if in.StartLine < 1 {
		return "", fmt.Errorf("start_line must be >= 1")
	}
	if in.EndLine < in.StartLine {
		return "", fmt.Errorf("end_line must be >= start_line")
	}

	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum >= in.StartLine && lineNum <= in.EndLine {
			lines = append(lines, scanner.Text())
		}
		if lineNum > in.EndLine {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

func (e *ToolExecutor) handleListDirectory(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if path == "" {
		path = "."
	}
	if err := validatePath(path); err != nil {
		return "", err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			lines = append(lines, "📁 "+e.Name()+"/")
		} else {
			info, _ := e.Info()
			size := ""
			if info != nil {
				size = fmt.Sprintf(" (%d bytes)", info.Size())
			}
			lines = append(lines, "  "+e.Name()+size)
		}
	}
	return strings.Join(lines, "\n"), nil
}

func (e *ToolExecutor) handleListDirectoryRecursive(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path    string `json:"path"`
		Exclude string `json:"exclude"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	var excludePatterns []string
	if in.Exclude != "" {
		excludePatterns = strings.Split(in.Exclude, ",")
		for i := range excludePatterns {
			excludePatterns[i] = strings.TrimSpace(excludePatterns[i])
		}
	}

	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			for _, pattern := range excludePatterns {
				if matched, _ := filepath.Match(pattern, info.Name()); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}
		relPath, _ := filepath.Rel(path, p)
		files = append(files, relPath)
		return nil
	})
	if err != nil {
		return "", err
	}
	return strings.Join(files, "\n"), nil
}

func (e *ToolExecutor) handleWriteFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(in.Content), 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Successfully wrote file: %s", path), nil
}

func (e *ToolExecutor) handlePatchFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path       string `json:"path"`
		OldContent string `json:"old_content"`
		NewContent string `json:"new_content"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)

	if !strings.Contains(content, in.OldContent) {
		return "", fmt.Errorf("old_content not found in file")
	}
	if strings.Count(content, in.OldContent) > 1 {
		return "", fmt.Errorf("old_content appears multiple times in file, cannot safely patch")
	}

	newContent := strings.Replace(content, in.OldContent, in.NewContent, 1)
	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Successfully patched file: %s", path), nil
}

func (e *ToolExecutor) handleDeleteFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}
	if err := os.Remove(path); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Successfully deleted file: %s", path), nil
}

func (e *ToolExecutor) handleRenameFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	from := resolvePath(in.From)
	to := resolvePath(in.To)
	if err := validatePath(from); err != nil {
		return "", err
	}
	if err := validatePath(to); err != nil {
		return "", err
	}
	if err := os.Rename(from, to); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Successfully moved/renamed: %s -> %s", from, to), nil
}

func (e *ToolExecutor) handleFileExists(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}
	_, err := os.Stat(path)
	if err == nil {
		return "true", nil
	}
	if os.IsNotExist(err) {
		return "false", nil
	}
	return "", err
}

func (e *ToolExecutor) handleSummarizeFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	// Simple summarization: extract functions/types based on file extension
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)
	ext := filepath.Ext(path)

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("📄 Summary of %s:\n\n", filepath.Base(path)))

	switch ext {
	case ".go":
		re := regexp.MustCompile(`(?m)^\s*(?:func|type)\s+(\w+)\s*(?:\([^)]*\))?\s*(?:[\[{])?`)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			summary.WriteString(fmt.Sprintf("  • %s\n", match[0]))
		}
	case ".js", ".ts", ".jsx", ".tsx":
		re := regexp.MustCompile(`(?m)^\s*(?:export\s+)?(?:function|const|let|var|class|interface|type)\s+(\w+)`)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			summary.WriteString(fmt.Sprintf("  • %s\n", match[0]))
		}
	case ".py":
		re := regexp.MustCompile(`(?m)^\s*(?:def|class)\s+(\w+)`)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			summary.WriteString(fmt.Sprintf("  • %s\n", match[0]))
		}
	default:
		summary.WriteString("  (Summarization not available for this file type)\n")
	}
	return summary.String(), nil
}

func (e *ToolExecutor) handleFormatFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	ext := filepath.Ext(path)
	var cmd *exec.Cmd

	switch ext {
	case ".go":
		cmd = exec.Command("gofmt", path)
	case ".js", ".ts", ".jsx", ".tsx", ".json":
		cmd = exec.Command("prettier", "--write", path)
	case ".py":
		cmd = exec.Command("black", path)
	case ".rb":
		cmd = exec.Command("ruby", "-c", path) // Just validate, no formatter
	default:
		return "", fmt.Errorf("no formatter available for extension: %s", ext)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("format error: %s", string(output))
	}
	return fmt.Sprintf("✅ Formatted %s", filepath.Base(path)), nil
}

func (e *ToolExecutor) handleCreateDirectory(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Directory created: %s", path), nil
}

// ── Search & navigation handlers ────────────────────────────────────────────

func (e *ToolExecutor) handleSearchInFiles(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Query   string `json:"query"`
		Regex   string `json:"regex"`
		Include string `json:"include"`
		Root    string `json:"root"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	root := resolvePath(in.Root)
	if err := validatePath(root); err != nil {
		return "", err
	}

	var pattern *regexp.Regexp
	var err error
	if in.Regex != "" {
		pattern, err = regexp.Compile(in.Regex)
		if err != nil {
			return "", fmt.Errorf("invalid regex: %w", err)
		}
	} else {
		pattern, err = regexp.Compile(regexp.QuoteMeta(in.Query))
		if err != nil {
			return "", err
		}
	}

	var results []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			return nil
		}
		if in.Include != "" {
			if matched, _ := filepath.Match(in.Include, filepath.Base(path)); !matched {
				return nil
			}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if pattern.MatchString(line) {
				relPath, _ := filepath.Rel(root, path)
				results = append(results, fmt.Sprintf("%s:%d: %s", relPath, i+1, strings.TrimSpace(line)))
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "No matches found", nil
	}
	return strings.Join(results, "\n"), nil
}

func (e *ToolExecutor) handleFindDefinition(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Symbol string `json:"symbol"`
		Root   string `json:"root"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	root := resolvePath(in.Root)
	if err := validatePath(root); err != nil {
		return "", err
	}

	var result string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(data)
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			// Simple pattern matching for definitions
			if strings.Contains(line, in.Symbol) {
				// Check if it looks like a definition
				if strings.Contains(line, "func") || strings.Contains(line, "function") ||
					strings.Contains(line, "class") || strings.Contains(line, "interface") ||
					strings.Contains(line, "type") || strings.Contains(line, "def") {
					relPath, _ := filepath.Rel(root, path)
					result = fmt.Sprintf("%s:%d: %s", relPath, i+1, strings.TrimSpace(line))
					return filepath.SkipAll
				}
			}
		}
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		return "", err
	}
	if result == "" {
		return fmt.Sprintf("Definition for '%s' not found", in.Symbol), nil
	}
	return result, nil
}

func (e *ToolExecutor) handleFindUsages(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Symbol string `json:"symbol"`
		Root   string `json:"root"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	root := resolvePath(in.Root)
	if err := validatePath(root); err != nil {
		return "", err
	}

	var usages []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.Contains(line, in.Symbol) {
				relPath, _ := filepath.Rel(root, path)
				usages = append(usages, fmt.Sprintf("%s:%d: %s", relPath, i+1, strings.TrimSpace(line)))
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(usages) == 0 {
		return "No usages found", nil
	}
	return strings.Join(usages, "\n"), nil
}

func (e *ToolExecutor) handleGetFileTree(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path  string `json:"path"`
		Depth int    `json:"depth"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if in.Depth == 0 {
		in.Depth = 3
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	var tree strings.Builder
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(path, p)
		if relPath == "." {
			tree.WriteString(filepath.Base(path) + "\n")
			return nil
		}
		depth := strings.Count(relPath, string(filepath.Separator))
		if depth > in.Depth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		indent := strings.Repeat("  ", depth)
		prefix := "├── "
		if info.IsDir() {
			prefix = "├── 📁 "
		}
		tree.WriteString(indent + prefix + info.Name() + "\n")
		return nil
	})
	if err != nil {
		return "", err
	}
	return tree.String(), nil
}

// ── Shell & process handlers ─────────────────────────────────────────────────

func (e *ToolExecutor) handleRunShell(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Command string `json:"command"`
		Cwd     string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	cmd := exec.Command("bash", "-c", in.Command)
	cmd.Dir = cwd

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %s\nstderr: %s", err.Error(), stderr.String())
	}
	if stderr.Len() > 0 {
		return stdout.String() + "\n[stderr]\n" + stderr.String(), nil
	}
	return stdout.String(), nil
}

func (e *ToolExecutor) handleRunTests(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Path    string `json:"path"`
		Pattern string `json:"pattern"`
		Cwd     string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	// Detect test runner
	var cmd *exec.Cmd
	if _, err := exec.LookPath("go"); err == nil && filepath.Base(cwd) == "go" {
		args := []string{"test"}
		if in.Path != "" {
			args = append(args, in.Path)
		}
		if in.Pattern != "" {
			args = append(args, "-run", in.Pattern)
		}
		cmd = exec.Command("go", args...)
	} else if _, err := exec.LookPath("npm"); err == nil {
		cmd = exec.Command("npm", "test", "--")
		if in.Path != "" {
			cmd.Args = append(cmd.Args, in.Path)
		}
		if in.Pattern != "" {
			cmd.Args = append(cmd.Args, "--testNamePattern", in.Pattern)
		}
	} else if _, err := exec.LookPath("pytest"); err == nil {
		args := []string{}
		if in.Path != "" {
			args = append(args, in.Path)
		}
		if in.Pattern != "" {
			args = append(args, "-k", in.Pattern)
		}
		cmd = exec.Command("pytest", args...)
	} else {
		return "", fmt.Errorf("no test runner found")
	}

	cmd.Dir = cwd
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func (e *ToolExecutor) handleLintFile(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.Path)
	if err := validatePath(path); err != nil {
		return "", err
	}

	ext := filepath.Ext(path)
	var cmd *exec.Cmd

	switch ext {
	case ".go":
		cmd = exec.Command("golangci-lint", "run", path)
	case ".js", ".ts", ".jsx", ".tsx":
		cmd = exec.Command("eslint", path)
	case ".py":
		cmd = exec.Command("pylint", path)
	default:
		return "", fmt.Errorf("no linter available for extension: %s", ext)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), nil // return linter errors as output, not error
	}
	if len(output) == 0 {
		return "✅ No linting issues found", nil
	}
	return string(output), nil
}

// ── Git handlers ─────────────────────────────────────────────────────────────

func (e *ToolExecutor) handleGitStatus(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Cwd string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "Working directory clean, no changes", nil
	}
	return string(output), nil
}

func (e *ToolExecutor) handleGitDiff(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Path string `json:"path"`
		Cwd  string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	args := []string{"diff"}
	if in.Path != "" {
		args = append(args, "--", in.Path)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "No differences", nil
	}
	return string(output), nil
}

func (e *ToolExecutor) handleGitLog(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Path  string `json:"path"`
		Limit int    `json:"limit"`
		Cwd   string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if in.Limit == 0 {
		in.Limit = 10
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	args := []string{"log", "--oneline", "-n", fmt.Sprintf("%d", in.Limit)}
	if in.Path != "" {
		args = append(args, "--", in.Path)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (e *ToolExecutor) handleGitCommit(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Message string `json:"message"`
		Cwd     string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	// Add all changes
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = cwd
	if err := addCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to stage changes: %w", err)
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", in.Message)
	commitCmd.Dir = cwd
	output, err := commitCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("commit failed: %s", string(output))
	}
	return fmt.Sprintf("✅ Commit successful: %s", in.Message), nil
}

// ── Dependencies handlers ───────────────────────────────────────────────────

func (e *ToolExecutor) handleReadDependencies(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Root string `json:"root"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	root := resolvePath(in.Root)
	if err := validatePath(root); err != nil {
		return "", err
	}

	dependencyFiles := []string{"package.json", "go.mod", "requirements.txt", "Cargo.toml", "pyproject.toml"}
	var results []string

	for _, file := range dependencyFiles {
		path := filepath.Join(root, file)
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			results = append(results, fmt.Sprintf("📦 %s:\n%s\n", file, string(data)))
		}
	}

	if len(results) == 0 {
		return "No dependency files found", nil
	}
	return strings.Join(results, "\n"), nil
}

func (e *ToolExecutor) handleInstallDependency(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Package string `json:"package"`
		Dev     string `json:"dev"`
		Cwd     string `json:"cwd"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	cwd := resolvePath(in.Cwd)
	if cwd == "" {
		cwd = e.Agent.ActivePath
	}

	var cmd *exec.Cmd
	if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
		args := []string{"install"}
		if in.Dev == "true" {
			args = append(args, "--save-dev")
		} else {
			args = append(args, "--save")
		}
		args = append(args, in.Package)
		cmd = exec.Command("npm", args...)
	} else if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		cmd = exec.Command("go", "get", in.Package)
	} else if _, err := os.Stat(filepath.Join(cwd, "requirements.txt")); err == nil {
		cmd = exec.Command("pip", "install", in.Package)
	} else {
		return "", fmt.Errorf("no recognized package manager found")
	}

	cmd.Dir = cwd
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("installation failed: %s", string(output))
	}
	return fmt.Sprintf("✅ Installed %s\n%s", in.Package, string(output)), nil
}

// ── Code generation handlers ────────────────────────────────────────────────

func (e *ToolExecutor) handleInitializeProject(rawInput json.RawMessage, resolvePath func(string) string) (string, error) {
	var in struct {
		Prompt   string `json:"prompt"`
		RootPath string `json:"root_path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.RootPath)
	if err := e.Agent.GenerateProject(in.Prompt, path, ""); err != nil {
		return "", err
	}
	return fmt.Sprintf("✅ Project initialized at: %s", path), nil
}

func (e *ToolExecutor) handleGenerateDataSchema(rawInput json.RawMessage) (string, error) {
	var in struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	schema, err := e.Agent.GenerateSchema(in.Description)
	if err != nil {
		return "", err
	}
	return schema, nil
}

func (e *ToolExecutor) handleGenerateCRUD(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		SchemaJSON string `json:"schema_json"`
		TargetPath string `json:"target_path"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	path := resolvePath(in.TargetPath)
	if err := validatePath(path); err != nil {
		return "", err
	}

	// Placeholder - implement actual CRUD generation
	return fmt.Sprintf("🔄 CRUD generation started for target: %s", path), nil
}

func (e *ToolExecutor) handleReadTemplate(rawInput json.RawMessage) (string, error) {
	var in struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	content := e.Agent.getSystemInstruction(in.Name)
	if content == "" {
		return "", fmt.Errorf("template not found: %s", in.Name)
	}
	return content, nil
}

func (e *ToolExecutor) handleListTemplates() (string, error) {
	entries, err := assets.ReadDir("templates")
	if err != nil {
		return "", err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return strings.Join(names, ", "), nil
}

// ── Workspace handlers ──────────────────────────────────────────────────────

func (e *ToolExecutor) handleGetActiveWorkspacePath() (string, error) {
	if e.Agent.ActivePath == "" {
		return "No active workspace path set. Please use an absolute path or initialize a project.", nil
	}
	return e.Agent.ActivePath, nil
}

func (e *ToolExecutor) handleGetActiveFile() (string, error) {
	// This would typically come from an editor integration
	// For now, return a placeholder
	return "No active file currently open", nil
}

func (e *ToolExecutor) handleGetWorkspaceConfig(rawInput json.RawMessage, resolvePath func(string) string, validatePath func(string) error) (string, error) {
	var in struct {
		Root string `json:"root"`
	}
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	root := resolvePath(in.Root)
	if err := validatePath(root); err != nil {
		return "", err
	}

	configFiles := []string{
		"tsconfig.json", ".eslintrc", ".eslintrc.json", ".eslintrc.js",
		"vite.config.js", "vite.config.ts", "webpack.config.js",
		"go.mod", "package.json", "Cargo.toml", "requirements.txt",
		"pom.xml", "build.gradle", "CMakeLists.txt",
	}

	var results []string
	for _, file := range configFiles {
		path := filepath.Join(root, file)
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			if len(data) > 1024*10 { // Truncate large files
				data = data[:1024*10]
				results = append(results, fmt.Sprintf("📄 %s (truncated):\n%s\n---\n", file, string(data)))
			} else {
				results = append(results, fmt.Sprintf("📄 %s:\n%s\n---\n", file, string(data)))
			}
		}
	}

	if len(results) == 0 {
		return "No configuration files found", nil
	}
	return strings.Join(results, "\n"), nil
}
