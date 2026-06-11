package types

// ── Public types ─────────────────────────────────────────────────────────────

// ScaffoldFile is a single file in the generated project plan.
type ScaffoldFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// ScaffoldPlan is the full structured response from the LLM.
type ScaffoldPlan struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Files       []ScaffoldFile `json:"files"`
}

// ScaffoldProgress is emitted to the frontend after each file is written.
type ScaffoldProgress struct {
	File  string `json:"file"`
	Index int    `json:"index"`
	Total int    `json:"total"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

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
