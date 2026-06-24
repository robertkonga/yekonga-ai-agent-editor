package agent

import (
	"yekonga-builder/agent/provider"
)

var agentTools = []provider.Tool{

	// ── Filesystem ────────────────────────────────────────────────────────────

	{
		Name:        "read_file",
		Description: "Read the full content of a file by its absolute path",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the file"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "read_file_range",
		Description: "Read a specific line range from a file. Use this instead of read_file when you only need part of a large file to save tokens.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":       {Type: "string", Description: "Absolute path to the file"},
				"start_line": {Type: "number", Description: "First line to read (1-indexed)"},
				"end_line":   {Type: "number", Description: "Last line to read (inclusive)"},
			},
			Required: []string{"path", "start_line", "end_line"},
		},
	},
	{
		Name:        "list_directory",
		Description: "List all files and folders inside a directory",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the directory"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "list_directory_recursive",
		Description: "Recursively list all files in a directory tree. Returns file paths only, no directories. Use to understand full project structure.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":    {Type: "string", Description: "Absolute path to the root directory"},
				"exclude": {Type: "string", Description: "Comma-separated glob patterns to exclude, e.g. 'node_modules,dist,.git'"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "write_file",
		Description: "Write or overwrite a file with new content",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":    {Type: "string"},
				"content": {Type: "string"},
			},
			Required: []string{"path", "content"},
		},
	},
	{
		Name:        "patch_file",
		Description: "Apply a targeted find-and-replace patch to a file without rewriting the whole thing. Faster and uses fewer tokens than write_file for small edits.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":        {Type: "string", Description: "Absolute path to the file"},
				"old_content": {Type: "string", Description: "Exact string to find (must be unique in the file)"},
				"new_content": {Type: "string", Description: "String to replace it with"},
			},
			Required: []string{"path", "old_content", "new_content"},
		},
	},
	{
		Name:        "delete_file",
		Description: "Delete a file at the given path.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the file to delete"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "rename_file",
		Description: "Rename or move a file from one path to another.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"from": {Type: "string", Description: "Current absolute path"},
				"to":   {Type: "string", Description: "New absolute path"},
			},
			Required: []string{"from", "to"},
		},
	},
	{
		Name:        "file_exists",
		Description: "Check whether a file or directory exists at a given path. Returns true or false.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to check"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "summarize_file",
		Description: "Return a token-efficient summary of a file: exported functions/types/components and their signatures only, no implementation. Use this instead of read_file when you only need to understand the API surface.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the file"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "format_file",
		Description: "Run the appropriate code formatter on a file (prettier, gofmt, black, etc.) based on its extension. Returns formatted content.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the file to format"},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "create_directory",
		Description: "Create a directory (and any missing parent directories) at the given path.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path of the directory to create"},
			},
			Required: []string{"path"},
		},
	},

	// ── Search & navigation ───────────────────────────────────────────────────

	{
		Name:        "search_in_files",
		Description: "Search for a string or pattern across all files in the workspace",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"query":   {Type: "string"},
				"regex":   {Type: "string", Description: "Enable regex search"},
				"include": {Type: "string", Description: "File pattern like *.go, *.ts"},
				"root":    {Type: "string", Description: "Root directory to search from"},
			},
			Required: []string{"query", "root"},
		},
	},
	{
		Name:        "find_definition",
		Description: "Find where a function, type, variable, or component is defined in the codebase. Returns file path and line number.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"symbol": {Type: "string", Description: "Name of the function, type, or variable to find"},
				"root":   {Type: "string", Description: "Root directory to search from"},
			},
			Required: []string{"symbol", "root"},
		},
	},
	{
		Name:        "find_usages",
		Description: "Find all files and lines that reference a given symbol, function, or import. Useful before refactoring.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"symbol": {Type: "string", Description: "Symbol name to find usages of"},
				"root":   {Type: "string", Description: "Root directory to search from"},
			},
			Required: []string{"symbol", "root"},
		},
	},
	{
		Name:        "get_file_tree",
		Description: "Get a compact tree representation of a directory (like the `tree` command). Good for understanding project layout in one shot.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":  {Type: "string", Description: "Root directory"},
				"depth": {Type: "number", Description: "Maximum depth to traverse (default 3)"},
			},
			Required: []string{"path"},
		},
	},

	// ── Shell & process ───────────────────────────────────────────────────────

	{
		Name:        "run_shell",
		Description: "Run a bash shell command. Use for installing packages, running build tools, linters, formatters, or test runners. Output returned as text.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"command": {Type: "string", Description: "The bash command to run"},
				"cwd":     {Type: "string", Description: "Working directory to run the command in (optional, defaults to workspace root)"},
			},
			Required: []string{"command"},
		},
	},
	{
		Name:        "run_tests",
		Description: "Run the project test suite or a specific test file/pattern. Returns pass/fail output.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":    {Type: "string", Description: "Path to test file or directory (optional, runs all tests if omitted)"},
				"pattern": {Type: "string", Description: "Test name pattern to filter (optional)"},
				"cwd":     {Type: "string", Description: "Working directory"},
			},
			Required: []string{"cwd"},
		},
	},
	{
		Name:        "lint_file",
		Description: "Run the linter on a specific file and return diagnostics (errors and warnings).",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "Absolute path to the file to lint"},
			},
			Required: []string{"path"},
		},
	},

	// ── Git ───────────────────────────────────────────────────────────────────

	{
		Name:        "git_status",
		Description: "Get the current git status: staged, unstaged, and untracked files.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"cwd": {Type: "string", Description: "Repository root path"},
			},
			Required: []string{"cwd"},
		},
	},
	{
		Name:        "git_diff",
		Description: "Get the git diff for a file or the whole repository. Use before making changes to understand what has already been modified.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path": {Type: "string", Description: "File path to diff (optional, diffs everything if omitted)"},
				"cwd":  {Type: "string", Description: "Repository root path"},
			},
			Required: []string{"cwd"},
		},
	},
	{
		Name:        "git_log",
		Description: "Get recent git commit history for a file or the whole repo.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"path":  {Type: "string", Description: "File to get history for (optional)"},
				"limit": {Type: "number", Description: "Number of commits to return (default 10)"},
				"cwd":   {Type: "string", Description: "Repository root path"},
			},
			Required: []string{"cwd"},
		},
	},
	{
		Name:        "git_commit",
		Description: "Stage all changes and create a git commit with the given message.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"message": {Type: "string", Description: "Commit message"},
				"cwd":     {Type: "string", Description: "Repository root path"},
			},
			Required: []string{"message", "cwd"},
		},
	},

	// ── Dependencies ──────────────────────────────────────────────────────────

	{
		Name:        "read_dependencies",
		Description: "Read the project dependency file (package.json, go.mod, requirements.txt, Cargo.toml) and return its contents.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"root": {Type: "string", Description: "Project root path"},
			},
			Required: []string{"root"},
		},
	},
	{
		Name:        "install_dependency",
		Description: "Add and install a new dependency to the project using the appropriate package manager (npm, go get, pip, cargo add).",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"package": {Type: "string", Description: "Package name (and optional version, e.g. 'axios@1.6.0')"},
				"dev":     {Type: "string", Description: "Pass 'true' to install as a dev dependency"},
				"cwd":     {Type: "string", Description: "Project root path"},
			},
			Required: []string{"package", "cwd"},
		},
	},

	// ── Code generation ───────────────────────────────────────────────────────

	{
		Name:        "create_project",
		Description: "Create a new project with requirements gathering, module grouping, and navigation planning. Takes project specifications and generates a comprehensive plan with module groupings and navigation structure before scaffolding.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"project_name":  {Type: "string", Description: "Name of the project"},
				"description":   {Type: "string", Description: "Detailed description of what the project does"},
				"size":          {Type: "string", Description: "Project size: 'small', 'middle', or 'large'"},
				"complexity":    {Type: "string", Description: "Project complexity: 'low', 'medium', or 'high'"},
				"modules":       {Type: "string", Description: "Comma-separated list of modules/features to include"},
				"features":      {Type: "string", Description: "Detailed feature descriptions and requirements"},
				"db_tables_min": {Type: "number", Description: "Minimum number of database tables/collections expected"},
				"db_tables_max": {Type: "number", Description: "Maximum number of database tables/collections expected"},
				"framework":     {Type: "string", Description: "Preferred framework/language (optional, will be auto-recommended if omitted)"},
				"root_path":     {Type: "string", Description: "The absolute path where the project should be created"},
			},
			Required: []string{"project_name", "description", "size", "complexity", "modules", "db_tables_min", "db_tables_max", "root_path"},
		},
	},
	{
		Name:        "initialize_project",
		Description: "Initialize a new project with a standard directory structure and configuration files based on a description.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"prompt":    {Type: "string", Description: "Detailed description of the project requirements, framework, and language."},
				"root_path": {Type: "string", Description: "The absolute path where the project should be initialized."},
			},
			Required: []string{"prompt", "root_path"},
		},
	},
	{
		Name:        "generate_dataschema",
		Description: "Generate a comprehensive database schema (JSON) based on a description of the data requirements.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"description": {Type: "string", Description: "Detailed description of the tables, fields, and relationships needed."},
			},
			Required: []string{"description"},
		},
	},
	{
		Name:        "generate_crud",
		Description: "Generate CRUD (Create, Read, Update, Delete) components and backend logic for a specific entity or the whole schema.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"schema_json": {Type: "string", Description: "The database schema in JSON format."},
				"target_path": {Type: "string", Description: "The absolute path where the CRUD files should be generated."},
			},
			Required: []string{"schema_json", "target_path"},
		},
	},
	{
		Name:        "read_template",
		Description: "Read a code generation template by name (e.g., 'form', 'list', 'report').",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"name": {Type: "string", Description: "Name of the template to read (without extension)."},
			},
			Required: []string{"name"},
		},
	},
	{
		Name:        "list_templates",
		Description: "List all available code generation templates.",
		InputSchema: provider.Schema{
			Type:       "object",
			Properties: map[string]provider.Property{},
		},
	},

	// ── Workspace ─────────────────────────────────────────────────────────────
	{
		Name:        "get_active_workspace_path",
		Description: "Get the absolute path of the currently open workspace/project.",
		InputSchema: provider.Schema{
			Type:       "object",
			Properties: map[string]provider.Property{},
		},
	},
	{
		Name:        "get_active_file",
		Description: "Get the absolute path and content of the file currently open in the editor.",
		InputSchema: provider.Schema{
			Type:       "object",
			Properties: map[string]provider.Property{},
		},
	},
	{
		Name:        "get_workspace_config",
		Description: "Read project-level config files: tsconfig.json, .eslintrc, vite.config, go.mod, etc. Returns all found config files as a map.",
		InputSchema: provider.Schema{
			Type: "object",
			Properties: map[string]provider.Property{
				"root": {Type: "string", Description: "Project root path"},
			},
			Required: []string{"root"},
		},
	},
}
