package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfoCommand(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantErr            bool
		wantOutputContains []string
		wantErrContains    string
	}{
		{
			name:    "info shows agents-sync metadata",
			args:    []string{"info", "agents-sync"},
			wantErr: false,
			wantOutputContains: []string{
				"agents-sync",
				"Agents Sync",
				"Synchronize AGENTS.md with actual repository state",
			},
		},
		{
			name:    "info shows build metadata",
			args:    []string{"info", "build"},
			wantErr: false,
			wantOutputContains: []string{
				"build",
				"Build",
				"Implement a task from work tracking",
			},
		},
		{
			name:            "info with unknown procedure shows error",
			args:            []string{"info", "nonexistent"},
			wantErr:         true,
			wantErrContains: "unknown procedure 'nonexistent'",
		},
		{
			name:            "info without procedure name shows error",
			args:            []string{"info"},
			wantErr:         true,
			wantErrContains: "accepts 1 arg(s), received 0",
		},
		{
			name:    "info help shows usage",
			args:    []string{"info", "--help"},
			wantErr: false,
			wantOutputContains: []string{
				"Display metadata, description, and configuration",
				"Usage:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			buf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(errBuf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String() + errBuf.String()

			if tt.wantErr && tt.wantErrContains != "" {
				if !strings.Contains(output, tt.wantErrContains) && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Expected error to contain %q, got: %s (err: %v)", tt.wantErrContains, output, err)
				}
			}

			for _, want := range tt.wantOutputContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
