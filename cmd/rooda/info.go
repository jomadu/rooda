package main

import (
	"fmt"

	"github.com/jomadu/rooda/internal/config"
	"github.com/spf13/cobra"
)

func newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <procedure>",
		Short: "Show detailed information about a procedure",
		Long:  `Display metadata, description, and configuration for a specific procedure.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			procedureName := args[0]
			return runInfo(cmd, procedureName)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// Load config to get procedure names
			flags := config.CLIFlags{
				ConfigPath: cfgFile,
				Verbose:    verbose,
				Quiet:      quiet,
				LogLevel:   logLevel,
			}
			cfg, err := config.LoadConfig(flags)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			// Return procedure names
			names := make([]string, 0, len(cfg.Procedures))
			for name := range cfg.Procedures {
				names = append(names, name)
			}
			return names, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

func runInfo(cmd *cobra.Command, procedureName string) error {
	// Build CLIFlags from persistent flags
	flags := config.CLIFlags{
		ConfigPath: cfgFile,
		Verbose:    verbose,
		Quiet:      quiet,
		LogLevel:   logLevel,
	}

	// Load merged config
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate procedure exists
	proc, exists := cfg.Procedures[procedureName]
	if !exists {
		return fmt.Errorf("unknown procedure '%s'\n\nRun 'rooda list' to see available procedures", procedureName)
	}

	// Display procedure information
	display := proc.Display
	if display == "" {
		display = procedureName
	}

	cmd.Printf("Procedure: %s\n", procedureName)
	cmd.Printf("Display Name: %s\n", display)
	cmd.Println()

	// Summary
	if proc.Summary != "" {
		cmd.Println("Summary:")
		cmd.Printf("  %s\n", proc.Summary)
		cmd.Println()
	}

	// Description
	if proc.Description != "" {
		cmd.Println("Description:")
		cmd.Printf("  %s\n", proc.Description)
		cmd.Println()
	}

	// OODA phases
	cmd.Println("OODA Phases:")
	cmd.Printf("  Observe:  %d fragment(s)\n", len(proc.Observe))
	cmd.Printf("  Orient:   %d fragment(s)\n", len(proc.Orient))
	cmd.Printf("  Decide:   %d fragment(s)\n", len(proc.Decide))
	cmd.Printf("  Act:      %d fragment(s)\n", len(proc.Act))
	cmd.Println()

	// Configuration overrides
	hasOverrides := false
	cmd.Println("Configuration:")
	if proc.IterationMode != "" {
		cmd.Printf("  Iteration mode: %s\n", proc.IterationMode)
		hasOverrides = true
	}
	if proc.DefaultMaxIterations != nil {
		cmd.Printf("  Max iterations: %d\n", *proc.DefaultMaxIterations)
		hasOverrides = true
	}
	if proc.IterationTimeout != nil {
		cmd.Printf("  Timeout: %ds\n", *proc.IterationTimeout)
		hasOverrides = true
	}
	if proc.AICmd != "" {
		cmd.Printf("  AI command: %s\n", proc.AICmd)
		hasOverrides = true
	}
	if proc.AICmdAlias != "" {
		cmd.Printf("  AI alias: %s\n", proc.AICmdAlias)
		hasOverrides = true
	}
	if !hasOverrides {
		cmd.Println("  (uses global defaults)")
	}

	return nil
}
