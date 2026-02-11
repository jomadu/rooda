package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestCLIStructure verifies the new CLI command structure
func TestCLIStructure(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectContains string
	}{
		{
			name:           "rooda version shows version info",
			args:           []string{"version"},
			expectError:    false,
			expectContains: "Version:",
		},
		{
			name:           "rooda list shows procedures",
			args:           []string{"list"},
			expectError:    false,
			expectContains: "Available procedures:",
		},
		{
			name:           "rooda info <procedure> shows procedure details",
			args:           []string{"info", "agents-sync"},
			expectError:    false,
			expectContains: "Procedure: agents-sync",
		},
		{
			name:        "rooda run requires procedure argument",
			args:        []string{"run"},
			expectError: true,
			// Error message goes to stderr, which cobra captures
			expectContains: "",
		},
		{
			name:           "rooda with no args shows help",
			args:           []string{},
			expectError:    false,
			expectContains: "rooda orchestrates AI coding agents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create root command
			cmd := newRootCommand()
			
			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			
			// Set args
			cmd.SetArgs(tt.args)
			
			// Execute
			err := cmd.Execute()
			
			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			// Check output contains expected string
			if tt.expectContains != "" {
				output := buf.String()
				if !strings.Contains(output, tt.expectContains) {
					t.Errorf("expected output to contain %q, got:\n%s", tt.expectContains, output)
				}
			}
		})
	}
}

// TestDeprecatedFlagsNotSupported verifies old --list-procedures flag is not supported
func TestDeprecatedFlagsNotSupported(t *testing.T) {
	// Note: --version is supported by cobra automatically and shows version info
	// This is acceptable as it's a common convention
	
	cmd := newRootCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--list-procedures"})
	
	err := cmd.Execute()
	output := buf.String()
	
	// Should not list procedures for --list-procedures
	if strings.Contains(output, "Available procedures:") {
		t.Errorf("--list-procedures flag should not be supported, but it listed procedures")
	}
	
	// Either error or help is acceptable
	if err == nil && !strings.Contains(output, "Usage:") {
		t.Errorf("expected error or help output, got: %s", output)
	}
}

// TestRunCommandStructure verifies rooda run <procedure> works
func TestRunCommandStructure(t *testing.T) {
	// This test verifies the command structure exists
	// Note: unknown-procedure actually exists as a built-in, so we use a truly unknown name
	cmd := newRootCommand()
	cmd.SetArgs([]string{"run", "this-procedure-definitely-does-not-exist-xyz", "--dry-run"})
	
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	
	err := cmd.Execute()
	
	// Should error about unknown procedure
	if err == nil {
		t.Error("expected error for unknown procedure")
	}
	
	if !strings.Contains(err.Error(), "unknown procedure") && !strings.Contains(err.Error(), "Unknown procedure") {
		t.Errorf("expected unknown procedure error, got: %v", err)
	}
}
