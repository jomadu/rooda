package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantOutputContains []string
	}{
		{
			name:    "help flag shows usage",
			args:    []string{"--help"},
			wantErr: false,
			wantOutputContains: []string{
				"rooda orchestrates",
				"Usage:",
				"Flags:",
			},
		},
		{
			name:    "version flag shows version",
			args:    []string{"--version"},
			wantErr: false,
			wantOutputContains: []string{
				"rooda",
			},
		},
		{
			name:    "no args shows help",
			args:    []string{},
			wantErr: false,
			wantOutputContains: []string{
				"Usage:",
			},
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
			cmd.SetArgs(tt.args)

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check output contains expected strings
			output := buf.String()
			for _, want := range tt.wantOutputContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing expected string %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestPersistentFlags(t *testing.T) {
	cmd := newRootCommand()

	// Check persistent flags exist
	persistentFlags := []string{
		"config",
		"verbose",
		"quiet",
		"log-level",
	}

	for _, flagName := range persistentFlags {
		flag := cmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected persistent flag %q to exist", flagName)
		}
	}
}

func TestFlagMutualExclusion(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "verbose and quiet are mutually exclusive",
			args:    []string{"--verbose", "--quiet"},
			wantErr: true,
			errMsg:  "verbose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
