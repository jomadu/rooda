package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantErrContain string
	}{
		{
			name:           "no procedure name",
			args:           []string{"run"},
			wantErr:        true,
			wantErrContain: "accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"run", "agents-sync", "extra"},
			wantErr:        true,
			wantErrContain: "accepts 1 arg(s), received 2",
		},
		{
			name:           "unknown procedure",
			args:           []string{"run", "nonexistent-procedure"},
			wantErr:        true,
			wantErrContain: "unknown procedure",
		},
		{
			name:    "valid procedure with help",
			args:    []string{"run", "agents-sync", "--help"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs(tt.args)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErrContain != "" && !strings.Contains(buf.String()+err.Error(), tt.wantErrContain) {
				t.Errorf("expected error containing %q, got: %v\nOutput: %s", tt.wantErrContain, err, buf.String())
			}
		})
	}
}

func TestRunCommandFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "max-iterations flag",
			args:    []string{"run", "agents-sync", "--max-iterations", "5", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "max-iterations short flag",
			args:    []string{"run", "agents-sync", "-n", "3", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "unlimited flag",
			args:    []string{"run", "agents-sync", "--unlimited", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "unlimited short flag",
			args:    []string{"run", "agents-sync", "-u", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "dry-run flag",
			args:    []string{"run", "agents-sync", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "dry-run short flag",
			args:    []string{"run", "agents-sync", "-d"},
			wantErr: false,
		},
		{
			name:    "ai-cmd flag",
			args:    []string{"run", "agents-sync", "--ai-cmd", "kiro-cli chat", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "ai-cmd-alias flag",
			args:    []string{"run", "agents-sync", "--ai-cmd-alias", "kiro-cli", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "context flag",
			args:    []string{"run", "agents-sync", "--context", "test context", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "context short flag",
			args:    []string{"run", "agents-sync", "-c", "test context", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "multiple context flags",
			args:    []string{"run", "agents-sync", "-c", "ctx1", "-c", "ctx2", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "observe flag",
			args:    []string{"run", "agents-sync", "--observe", "custom.md", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "mutually exclusive max-iterations and unlimited",
			args:    []string{"run", "agents-sync", "--max-iterations", "5", "--unlimited"},
			wantErr: true,
		},
		{
			name:    "invalid max-iterations",
			args:    []string{"run", "agents-sync", "--max-iterations", "0"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs(tt.args)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v\nOutput: %s", err, buf.String())
			}
		})
	}
}

func TestRunCommandCompletion(t *testing.T) {
	cmd := newRootCommand()
	runCmd := findCommand(cmd, "run")
	if runCmd == nil {
		t.Fatal("run command not found")
	}

	// Test that ValidArgsFunction is set
	if runCmd.ValidArgsFunction == nil {
		t.Error("ValidArgsFunction not set on run command")
	}
}

// Helper to find a command by name
func findCommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
