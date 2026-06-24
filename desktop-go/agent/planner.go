package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yekonga-builder/agent/provider"
	agentTypes "yekonga-builder/agent/types"
	"yekonga-builder/types"
)

// ── Data structures for project planning ──────────────────────────────────────

// ProjectRequirements holds all the collected project specifications.
type ProjectRequirements struct {
	Name        string `json:"project_name"`
	Description string `json:"description"`
	Size        string `json:"size"`       // "small", "middle", "large"
	Complexity  string `json:"complexity"` // "low", "medium", "high"
	Framework   string `json:"framework,omitempty"`
	DBTablesMin int    `json:"db_tables_min"`
	DBTablesMax int    `json:"db_tables_max"`
	Modules     string `json:"modules"`  // comma-separated
	Features    string `json:"features"` // detailed feature description
	RootPath    string `json:"root_path"`
}

// ModuleGroup represents a group of related modules/features.
type ModuleGroup struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Modules     []string `json:"modules"`
}

// NavigationItem represents a navigation entry in the UI.
type NavigationItem struct {
	Label    string           `json:"label"`
	Icon     string           `json:"icon"`
	Route    string           `json:"route"`
	Children []NavigationItem `json:"children,omitempty"`
}

// ProjectPlan is the complete output from the planner.
type ProjectPlan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Complexity  string `json:"complexity"`
	Framework   string `json:"framework"`
	DBTablesMin int    `json:"db_tables_min"`
	DBTablesMax int    `json:"db_tables_max"`

	// Module groups: features grouped into logical domains
	ModuleGroups []ModuleGroup `json:"module_groups"`

	// Navigation: suggested UI navigation structure
	Navigation []NavigationItem `json:"navigation"`

	// Database collections/tables suggested
	DatabaseCollections []string `json:"database_collections"`

	// Directory structure suggestion
	DirectoryStructure []string `json:"directory_structure"`

	// Technologies stack
	TechStack map[string]string `json:"tech_stack"`

	// Any additional notes
	Notes []string `json:"notes"`
}

const plannerSystemPrompt = `You are a professional system architect and project planner.
Given a set of project requirements, you must produce a comprehensive project plan as a JSON object.

Respond ONLY with a valid JSON object — no markdown fences, no explanation, no preamble.

The JSON must match this exact schema:
{
  "name": "project-name",
  "description": "Brief description of the project",
  "size": "small|middle|large",
  "complexity": "low|medium|high",
  "framework": "Recommended framework/language",
  "db_tables_min": 5,
  "db_tables_max": 20,
  "module_groups": [
    {
      "name": "Module Group Name",
      "description": "What this group does",
      "modules": ["module1", "module2"]
    }
  ],
  "navigation": [
    {
      "label": "Dashboard",
      "icon": "dashboard",
      "route": "/dashboard",
      "children": []
    }
  ],
  "database_collections": ["collection1", "collection2"],
  "directory_structure": ["src/modules/auth/", "src/modules/users/"],
  "tech_stack": {
    "frontend": "framework",
    "backend": "framework",
    "database": "database_type"
  },
  "notes": ["Important architectural note"]
}

Rules:
- Group modules logically by feature domain (e.g., "User Management" group contains login, register, profile, roles)
- Navigation should reflect the module groups and be hierarchical where appropriate
- Directory structure should follow standard conventions for the recommended framework
- Database collections should cover all entities needed
- The plan must be comprehensive and ready for implementation
- Never use ".." in any path
- All paths must use forward slashes`

// GenerateProjectPlan uses the LLM to create a comprehensive project plan from requirements.
func (a *Agent) GenerateProjectPlan(req ProjectRequirements) (*ProjectPlan, error) {
	// Default to Anthropic for planning as it handles structured output well
	selectedProvider := provider.NewAnthropicProvider(a.AnthropicApiKey, "")

	a.Emit(types.ScaffoldProgress{File: "Analyzing project requirements…"})

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	// Build the user prompt from requirements
	userPrompt := fmt.Sprintf(`I need a comprehensive project plan for the following requirements:

Project Name: %s
Description: %s
Project Size: %s
Complexity: %s
Framework (if specified): %s
Database Tables/Collections: %d to %d
Modules to consider: %s
Features: %s

Please provide a detailed project plan with:
1. Module groupings (group related features into logical domains)
2. Navigation structure (what should appear in the sidebar/menu)
3. Database collections/tables needed
4. Directory structure following best practices
5. Recommended technology stack
6. Any architectural notes

Make sure module groups are meaningful and related to features.
Navigation should reflect the module groups and provide a clear user flow.`,
		req.Name,
		req.Description,
		req.Size,
		req.Complexity,
		req.Framework,
		req.DBTablesMin,
		req.DBTablesMax,
		req.Modules,
		req.Features,
	)

	a.Emit(types.ScaffoldProgress{File: "Designing module groups and navigation…"})

	resp, err := selectedProvider.Complete(ctx, plannerSystemPrompt, []agentTypes.ChatMessage{{Role: "user", Content: userPrompt}}, nil)
	if err != nil {
		return nil, fmt.Errorf("planning LLM call failed: %w", err)
	}

	a.Emit(types.ScaffoldProgress{File: "Structuring project plan…"})

	plan, err := parsePlanJSON(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("parse plan failed: %w", err)
	}

	// Fill in fields from requirements
	plan.Name = req.Name
	plan.Description = req.Description
	plan.Size = req.Size
	plan.Complexity = req.Complexity
	if req.Framework != "" {
		plan.Framework = req.Framework
	}
	plan.DBTablesMin = req.DBTablesMin
	plan.DBTablesMax = req.DBTablesMax

	return plan, nil
}

// CreateProjectFromPlan generates the project scaffolding based on the plan.
func (a *Agent) CreateProjectFromPlan(plan *ProjectPlan, rootPath string) error {
	// Build a comprehensive prompt for the scaffold generator
	var moduleGroups []string
	for _, group := range plan.ModuleGroups {
		modules := strings.Join(group.Modules, ", ")
		moduleGroups = append(moduleGroups, fmt.Sprintf("- %s (%s): %s", group.Name, group.Description, modules))
	}

	var navStr []string
	var buildNav func(items []NavigationItem, depth int)
	buildNav = func(items []NavigationItem, depth int) {
		indent := strings.Repeat("  ", depth)
		for _, item := range items {
			navStr = append(navStr, fmt.Sprintf("%s- %s -> %s", indent, item.Label, item.Route))
			if len(item.Children) > 0 {
				buildNav(item.Children, depth+1)
			}
		}
	}
	buildNav(plan.Navigation, 0)

	dirStructure := strings.Join(plan.DirectoryStructure, "\n")
	techStack := ""
	for k, v := range plan.TechStack {
		techStack += fmt.Sprintf("- %s: %s\n", k, v)
	}

	scaffoldPrompt := fmt.Sprintf(`Create a complete project called "%s" based on the following plan:

DESCRIPTION: %s
SIZE: %s
COMPLEXITY: %s
FRAMEWORK: %s
TECH STACK:
%s

MODULE GROUPS:
%s

NAVIGATION STRUCTURE:
%s

DATABASE COLLECTIONS NEEDED:
%s

DIRECTORY STRUCTURE:
%s

ARCHITECTURAL NOTES:
%s

Create all necessary files including package.json/config files, main app structure, module directories, and placeholder components for each module group. Follow best practices for the chosen framework.`,
		plan.Name,
		plan.Description,
		plan.Size,
		plan.Complexity,
		plan.Framework,
		techStack,
		strings.Join(moduleGroups, "\n"),
		strings.Join(navStr, "\n"),
		strings.Join(plan.DatabaseCollections, ", "),
		dirStructure,
		strings.Join(plan.Notes, "\n"),
	)

	extraConventions := fmt.Sprintf(`Project structure must follow the module groups defined above.
Navigation structure must be implemented as defined.
The project name is: %s
Project size: %s`, plan.Name, plan.Size)

	return a.GenerateProject(scaffoldPrompt, rootPath, extraConventions)
}

// parsePlanJSON extracts the ProjectPlan from raw LLM text.
func parsePlanJSON(raw string) (*ProjectPlan, error) {
	cleaned := strings.TrimSpace(raw)

	// Strip optional ```json ... ``` fences
	if strings.HasPrefix(cleaned, "```") {
		start := strings.Index(cleaned, "\n")
		end := strings.LastIndex(cleaned, "```")
		if start != -1 && end > start {
			cleaned = strings.TrimSpace(cleaned[start+1 : end])
		}
	}

	var plan ProjectPlan
	if err := json.Unmarshal([]byte(cleaned), &plan); err != nil {
		return nil, fmt.Errorf("parse plan JSON: %w\nraw output: %.400s", err, cleaned)
	}

	if err := validatePlanContent(&plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

// validatePlanContent checks the plan for required fields.
func validatePlanContent(plan *ProjectPlan) error {
	if strings.TrimSpace(plan.Name) == "" {
		return fmt.Errorf("plan has no project name")
	}
	if len(plan.ModuleGroups) == 0 {
		return fmt.Errorf("plan has no module groups defined")
	}
	if len(plan.Navigation) == 0 {
		return fmt.Errorf("plan has no navigation defined")
	}
	if len(plan.DatabaseCollections) == 0 && plan.DBTablesMin > 0 {
		return fmt.Errorf("plan has no database collections but min tables is %d", plan.DBTablesMin)
	}
	return nil
}

// CreateProjectDirectory creates the project root directory and returns its path.
func CreateProjectDirectory(rootPath string, projectName string) (string, error) {
	projectDir := filepath.Join(rootPath, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", fmt.Errorf("create project directory: %w", err)
	}
	return projectDir, nil
}
