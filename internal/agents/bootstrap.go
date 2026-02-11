package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BootstrapAgentsMD analyzes a repository and generates an AgentsMD structure
func BootstrapAgentsMD(repoPath string) (*AgentsMD, error) {
	agentsMD := &AgentsMD{
		FilePath: filepath.Join(repoPath, "AGENTS.md"),
	}

	agentsMD.BuildCommand = detectBuildSystem(repoPath)
	agentsMD.TestCommand = detectTestSystem(repoPath)
	agentsMD.LintCommands = detectLintSystem(repoPath)

	specPaths, specExcludes := detectSpecPaths(repoPath)
	agentsMD.SpecPaths = specPaths
	agentsMD.SpecExcludes = specExcludes

	implPaths, implExcludes := detectImplPaths(repoPath)
	agentsMD.ImplPaths = implPaths
	agentsMD.ImplExcludes = implExcludes

	agentsMD.DocsPaths = detectDocsPaths(repoPath)

	workTracking, err := detectWorkTracking(repoPath)
	if err != nil {
		return nil, fmt.Errorf("detect work tracking: %w", err)
	}
	agentsMD.WorkTracking = *workTracking

	agentsMD.TaskInput = TaskInputConfig{
		Location: "TASK.md at project root",
		Format:   "Markdown file with task description, requirements, and acceptance criteria",
	}

	agentsMD.PlanningSystem = PlanningConfig{
		DraftPlanLocation: "PLAN.md at project root",
		PublishingMethod:  "Agent reads PLAN.md and runs work tracking commands to create issues",
	}

	agentsMD.AuditOutput = AuditOutputConfig{
		LocationPattern: "AUDIT-{procedure}.md at project root",
		Format:          "Markdown",
	}

	return agentsMD, nil
}

func detectBuildSystem(repoPath string) string {
	checks := []struct {
		file    string
		command string
	}{
		{"go.mod", "go build -o bin/rooda ./cmd/rooda"},
		{"Cargo.toml", "cargo build"},
		{"package.json", "npm run build"},
		{"Makefile", "make build"},
		{"build.sh", "./build.sh"},
	}

	for _, check := range checks {
		if fileExists(filepath.Join(repoPath, check.file)) {
			return check.command
		}
	}

	return "Not required (interpreted language)"
}

func detectTestSystem(repoPath string) string {
	if fileExists(filepath.Join(repoPath, "go.mod")) {
		if hasFilesWithSuffix(repoPath, "_test.go") {
			return "go test ./..."
		}
	}

	if fileExists(filepath.Join(repoPath, "Cargo.toml")) {
		return "cargo test"
	}

	if fileExists(filepath.Join(repoPath, "package.json")) {
		return "npm test"
	}

	if fileExists(filepath.Join(repoPath, "pytest.ini")) || fileExists(filepath.Join(repoPath, "conftest.py")) {
		return "pytest"
	}

	if hasFilesWithSuffix(repoPath, "_test.py") {
		return "pytest"
	}

	return "Manual verification (no automated tests)"
}

func detectLintSystem(repoPath string) []string {
	var linters []string

	if fileExists(filepath.Join(repoPath, "go.mod")) {
		linters = append(linters, "go vet ./...")
	}

	if fileExists(filepath.Join(repoPath, "Cargo.toml")) {
		linters = append(linters, "cargo clippy")
	}

	if fileExists(filepath.Join(repoPath, ".eslintrc.json")) || fileExists(filepath.Join(repoPath, ".eslintrc.js")) {
		linters = append(linters, "npm run lint")
	}

	if hasFilesWithSuffix(repoPath, ".sh") {
		linters = append(linters, "shellcheck *.sh")
	}

	return linters
}

func detectSpecPaths(repoPath string) ([]string, []string) {
	var paths []string
	var excludes []string

	if dirExists(filepath.Join(repoPath, "specs")) {
		paths = append(paths, "specs/*.md")
		excludes = append(excludes, "specs/README.md")
	} else if dirExists(filepath.Join(repoPath, "docs")) {
		paths = append(paths, "docs/**/*.md")
	} else if dirExists(filepath.Join(repoPath, "documentation")) {
		paths = append(paths, "documentation/**/*.md")
	}

	return paths, excludes
}

func detectImplPaths(repoPath string) ([]string, []string) {
	var paths []string
	var excludes []string

	if fileExists(filepath.Join(repoPath, "go.mod")) {
		paths = append(paths, "*.go", "cmd/**/*.go", "internal/**/*.go")
		excludes = append(excludes, "*_test.go")
	} else if fileExists(filepath.Join(repoPath, "Cargo.toml")) {
		paths = append(paths, "src/**/*.rs")
	} else if fileExists(filepath.Join(repoPath, "package.json")) {
		paths = append(paths, "src/**/*.{ts,tsx,js,jsx}")
	} else if dirExists(filepath.Join(repoPath, "src")) {
		paths = append(paths, "src/**/*")
	} else if hasFilesWithSuffix(repoPath, ".sh") {
		paths = append(paths, "*.sh")
	}

	excludes = append(excludes, "node_modules", "vendor", ".git")

	return paths, excludes
}

func detectWorkTracking(repoPath string) (*WorkTrackingConfig, error) {
	if dirExists(filepath.Join(repoPath, ".beads")) {
		return &WorkTrackingConfig{
			System:        "beads",
			QueryCommand:  "bd ready --json",
			UpdateCommand: "bd update <id> --status <status>",
			CloseCommand:  "bd close <id> --reason <reason>",
			CreateCommand: "bd create --title <title> --description <desc> --priority <priority>",
		}, nil
	}

	if dirExists(filepath.Join(repoPath, ".github")) {
		return &WorkTrackingConfig{
			System:        "github-issues",
			QueryCommand:  "gh issue list --json number,title,labels --label ready",
			UpdateCommand: "gh issue edit <id> --add-label in-progress",
			CloseCommand:  "gh issue close <id> --comment <reason>",
			CreateCommand: "gh issue create --title <title> --body <desc>",
		}, nil
	}

	if fileExists(filepath.Join(repoPath, "TODO.md")) || fileExists(filepath.Join(repoPath, "TASKS.md")) {
		return &WorkTrackingConfig{
			System:        "file-based",
			QueryCommand:  "cat TODO.md",
			UpdateCommand: "Edit TODO.md manually",
			CloseCommand:  "Remove from TODO.md",
			CreateCommand: "Add to TODO.md",
		}, nil
	}

	return &WorkTrackingConfig{
		System:        "not-configured",
		QueryCommand:  "Not configured",
		UpdateCommand: "Not configured",
		CloseCommand:  "Not configured",
		CreateCommand: "Not configured",
	}, nil
}

func detectDocsPaths(repoPath string) []string {
	if dirExists(filepath.Join(repoPath, "docs")) {
		return []string{"docs/*.md"}
	}
	return []string{}
}

// Helper functions

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func hasFilesWithSuffix(dir, suffix string) bool {
	found := false
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
