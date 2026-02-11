package main

import (
	"fmt"
	"sort"

	"github.com/jomadu/rooda/internal/config"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available procedures",
		Long:  `List all available procedures with their names and summaries.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd)
		},
	}

	return cmd
}

func runList(cmd *cobra.Command) error {
	// Build CLIFlags from persistent flags
	flags := config.CLIFlags{
		ConfigPath: cfgFile,
		Verbose:    verbose,
		Quiet:      quiet,
		LogLevel:   logLevel,
	}

	// Load merged config (built-in + global + workspace)
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(cfg.Procedures) == 0 {
		cmd.Println("No procedures defined.")
		return nil
	}

	// Sort procedure names for consistent output
	names := make([]string, 0, len(cfg.Procedures))
	for name := range cfg.Procedures {
		names = append(names, name)
	}
	sort.Strings(names)

	// Display procedures
	cmd.Println("Available procedures:")
	for _, name := range names {
		proc := cfg.Procedures[name]
		summary := proc.Summary
		if summary == "" {
			summary = "(no summary)"
		}
		// Truncate long summaries
		if len(summary) > 80 {
			summary = summary[:77] + "..."
		}
		cmd.Printf("  %-20s %s\n", name, summary)
	}

	return nil
}
