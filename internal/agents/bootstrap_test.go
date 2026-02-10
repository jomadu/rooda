package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectBuildSystem(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected string
	}{
		{
			name:     "go.mod detected",
			files:    map[string]string{"go.mod": "module test"},
			expected: "go build -o bin/rooda ./cmd/rooda",
		},
		{
			name:     "Cargo.toml detected",
			files:    map[string]string{"Cargo.toml": "[package]"},
			expected: "cargo build",
		},
		{
			name:     "package.json detected",
			files:    map[string]string{"package.json": `{"scripts":{"build":"tsc"}}`},
			expected: "npm run build",
		},
		{
			name:     "Makefile detected",
			files:    map[string]string{"Makefile": "build:\n\tgo build"},
			expected: "make build",
		},
		{
			name:     "build.sh detected",
			files:    map[string]string{"build.sh": "#!/bin/bash"},
			expected: "./build.sh",
		},
		{
			name:     "no build system",
			files:    map[string]string{},
			expected: "Not required (interpreted language)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			result := detectBuildSystem(tmpDir)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectTestSystem(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected string
	}{
		{
			name:     "go test files",
			files:    map[string]string{"go.mod": "module test", "main_test.go": "package main"},
			expected: "go test ./...",
		},
		{
			name:     "cargo test",
			files:    map[string]string{"Cargo.toml": "[package]"},
			expected: "cargo test",
		},
		{
			name:     "npm test",
			files:    map[string]string{"package.json": `{"scripts":{"test":"jest"}}`},
			expected: "npm test",
		},
		{
			name:     "pytest",
			files:    map[string]string{"pytest.ini": "[pytest]"},
			expected: "pytest",
		},
		{
			name:     "no test system",
			files:    map[string]string{},
			expected: "Manual verification (no automated tests)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			result := detectTestSystem(tmpDir)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectWorkTracking(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected string
	}{
		{
			name:     "beads detected",
			files:    map[string]string{".beads/config.json": "{}"},
			expected: "beads",
		},
		{
			name:     "github detected",
			files:    map[string]string{".github/workflows/ci.yml": "name: CI"},
			expected: "github-issues",
		},
		{
			name:     "file-based detected",
			files:    map[string]string{"TODO.md": "# Tasks"},
			expected: "file-based",
		},
		{
			name:     "no work tracking",
			files:    map[string]string{},
			expected: "not-configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tt.files {
				dir := filepath.Dir(filepath.Join(tmpDir, name))
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			result, err := detectWorkTracking(tmpDir)
			if err != nil {
				t.Fatal(err)
			}
			if result.System != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.System)
			}
		})
	}
}

func TestBootstrapAgentsMD(t *testing.T) {
	t.Run("go project with beads", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		// Create go.mod
		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
			t.Fatal(err)
		}
		
		// Create test file
		if err := os.WriteFile(filepath.Join(tmpDir, "main_test.go"), []byte("package main"), 0644); err != nil {
			t.Fatal(err)
		}
		
		// Create .beads directory
		if err := os.MkdirAll(filepath.Join(tmpDir, ".beads"), 0755); err != nil {
			t.Fatal(err)
		}
		
		// Create specs directory
		if err := os.MkdirAll(filepath.Join(tmpDir, "specs"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "specs", "test.md"), []byte("# Test"), 0644); err != nil {
			t.Fatal(err)
		}
		
		result, err := BootstrapAgentsMD(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		
		if result.BuildCommand != "go build -o bin/rooda ./cmd/rooda" {
			t.Errorf("expected go build command, got %q", result.BuildCommand)
		}
		if result.TestCommand != "go test ./..." {
			t.Errorf("expected go test command, got %q", result.TestCommand)
		}
		if result.WorkTracking.System != "beads" {
			t.Errorf("expected beads, got %q", result.WorkTracking.System)
		}
		if len(result.SpecPaths) == 0 || result.SpecPaths[0] != "specs/*.md" {
			t.Errorf("expected specs/*.md, got %v", result.SpecPaths)
		}
	})

	t.Run("empty repository", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		result, err := BootstrapAgentsMD(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		
		if result.BuildCommand != "Not required (interpreted language)" {
			t.Errorf("expected 'Not required', got %q", result.BuildCommand)
		}
		if result.TestCommand != "Manual verification (no automated tests)" {
			t.Errorf("expected 'Manual verification', got %q", result.TestCommand)
		}
		if result.WorkTracking.System != "not-configured" {
			t.Errorf("expected 'not-configured', got %q", result.WorkTracking.System)
		}
	})
}
