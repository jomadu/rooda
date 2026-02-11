package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantErr            bool
		wantOutputContains []string
	}{
		{
			name:    "version shows version info",
			args:    []string{"version"},
			wantErr: false,
			wantOutputContains: []string{
				"Version:",
				"Commit:",
				"Built:",
			},
		},
		{
			name:    "version help shows usage",
			args:    []string{"version", "--help"},
			wantErr: false,
			wantOutputContains: []string{
				"Display version, commit SHA, and build date",
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

			for _, want := range tt.wantOutputContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
