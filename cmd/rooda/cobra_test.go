package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCobraDependencyAvailable(t *testing.T) {
	// Verify Cobra can be imported and basic functionality works
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	if cmd.Use != "test" {
		t.Errorf("Expected Use='test', got '%s'", cmd.Use)
	}

	if cmd.Short != "Test command" {
		t.Errorf("Expected Short='Test command', got '%s'", cmd.Short)
	}
}
