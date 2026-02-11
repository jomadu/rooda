package docs

import (
	"os"
	"path/filepath"
	"testing"
)

// getProjectRoot returns the project root directory
func getProjectRoot(t *testing.T) string {
	t.Helper()
	// Start from current directory and walk up until we find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

func TestDocumentationStructure(t *testing.T) {
	root := getProjectRoot(t)

	tests := []struct {
		name     string
		docPath  string
		required []string
	}{
		{
			name:    "installation.md has required sections",
			docPath: "docs/installation.md",
			required: []string{
				"Installation",
				"Quick Install",
			},
		},
		{
			name:    "procedures.md has required sections",
			docPath: "docs/procedures.md",
			required: []string{
				"Procedures",
				"Bootstrap",
			},
		},
		{
			name:    "cli-reference.md has required sections",
			docPath: "docs/cli-reference.md",
			required: []string{
				"Commands",
				"Global Flags",
			},
		},
		{
			name:    "configuration.md has required sections",
			docPath: "docs/configuration.md",
			required: []string{
				"Configuration",
				"Configuration Tiers",
			},
		},
		{
			name:    "troubleshooting.md has required sections",
			docPath: "docs/troubleshooting.md",
			required: []string{
				"Common",
			},
		},
		{
			name:    "agents-md.md has required sections",
			docPath: "docs/agents-md.md",
			required: []string{
				"Format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(root, tt.docPath)
			// Check file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Fatalf("Documentation file missing: %s", tt.docPath)
			}

			// Read content
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", tt.docPath, err)
			}

			// Check required sections exist
			contentStr := string(content)
			for _, section := range tt.required {
				if !containsSection(contentStr, section) {
					t.Errorf("Missing required section '%s' in %s", section, tt.docPath)
				}
			}
		})
	}
}

func TestREADMEStructure(t *testing.T) {
	root := getProjectRoot(t)

	required := []string{
		"What is rooda",
		"Quick Start",
		"Installation",
		"Core Concepts",
	}

	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	contentStr := string(content)
	for _, section := range required {
		if !containsSection(contentStr, section) {
			t.Errorf("Missing required section '%s' in README.md", section)
		}
	}
}

func TestDocumentationFilesExist(t *testing.T) {
	root := getProjectRoot(t)

	requiredDocs := []string{
		"docs/installation.md",
		"docs/procedures.md",
		"docs/cli-reference.md",
		"docs/configuration.md",
		"docs/troubleshooting.md",
		"docs/agents-md.md",
		"README.md",
	}

	for _, doc := range requiredDocs {
		t.Run(doc, func(t *testing.T) {
			fullPath := filepath.Join(root, doc)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("Required documentation file missing: %s", doc)
			}
		})
	}
}

func TestDocumentationNotEmpty(t *testing.T) {
	root := getProjectRoot(t)

	docs := []string{
		"docs/installation.md",
		"docs/procedures.md",
		"docs/cli-reference.md",
		"docs/configuration.md",
		"docs/troubleshooting.md",
		"docs/agents-md.md",
		"README.md",
	}

	for _, doc := range docs {
		t.Run(doc, func(t *testing.T) {
			fullPath := filepath.Join(root, doc)
			info, err := os.Stat(fullPath)
			if err != nil {
				t.Fatalf("Failed to stat %s: %v", doc, err)
			}

			if info.Size() == 0 {
				t.Errorf("Documentation file is empty: %s", doc)
			}
		})
	}
}

func TestCrossReferencesResolve(t *testing.T) {
	root := getProjectRoot(t)

	docs := []string{
		"docs/installation.md",
		"docs/procedures.md",
		"docs/cli-reference.md",
		"docs/configuration.md",
		"docs/troubleshooting.md",
		"docs/agents-md.md",
		"README.md",
	}

	for _, doc := range docs {
		t.Run(doc, func(t *testing.T) {
			fullPath := filepath.Join(root, doc)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", doc, err)
			}

			refs := extractMarkdownLinks(string(content))
			for _, ref := range refs {
				// Skip external URLs
				if isExternalURL(ref) {
					continue
				}

				// Resolve relative to doc directory
				docDir := filepath.Dir(fullPath)
				targetPath := filepath.Join(docDir, ref)

				// Check if target exists
				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					t.Errorf("Broken cross-reference in %s: %s (resolved to %s)", doc, ref, targetPath)
				}
			}
		})
	}
}

// Helper functions

func containsSection(content, section string) bool {
	// Check for markdown headers containing the section text
	patterns := []string{
		"# " + section,
		"## " + section,
		"### " + section,
	}

	for _, pattern := range patterns {
		if contains(content, pattern) {
			return true
		}
	}

	return false
}

func contains(content, substr string) bool {
	return len(content) >= len(substr) && findSubstring(content, substr)
}

func findSubstring(content, substr string) bool {
	for i := 0; i <= len(content)-len(substr); i++ {
		if content[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractMarkdownLinks(content string) []string {
	var links []string
	// Simple extraction: look for [text](link) patterns
	i := 0
	for i < len(content) {
		if i+1 < len(content) && content[i] == ']' && content[i+1] == '(' {
			// Found potential link
			start := i + 2
			end := start
			for end < len(content) && content[end] != ')' {
				end++
			}
			if end < len(content) {
				link := content[start:end]
				// Remove anchor if present
				if anchorIdx := findChar(link, '#'); anchorIdx != -1 {
					link = link[:anchorIdx]
				}
				if link != "" {
					links = append(links, link)
				}
			}
			i = end
		}
		i++
	}
	return links
}

func findChar(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func isExternalURL(link string) bool {
	if len(link) >= 8 && link[:8] == "https://" {
		return true
	}
	if len(link) >= 7 && link[:7] == "http://" {
		return true
	}
	return false
}
