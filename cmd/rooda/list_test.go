package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantErr            bool
		wantOutputContains []string
	}{
		{
			name:    "list shows built-in procedures",
			args:    []string{"list"},
			wantErr: false,
			wantOutputContains: []string{
				"agents-sync",
				"build",
				"publish-plan",
				"audit-spec",
			},
		},
		{
			name:    "list shows procedure summaries",
			args:    []string{"list"},
			wantErr: false,
			wantOutputContains: []string{
				"Synchronize AGENTS.md",
				"Implement a task",
			},
		},
		{
			name:    "list help shows usage",
			args:    []string{"list", "--help"},
			wantErr: false,
			wantOutputContains: []string{
				"List all available procedures",
				"Usage:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			for _, want := range tt.wantOutputContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing expected string %q\nGot:\n%s", want, output)
				}
			}
		})
	}
}
