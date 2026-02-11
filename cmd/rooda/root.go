package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Persistent flags (available to all commands)
	cfgFile  string
	verbose  bool
	quiet    bool
	logLevel string
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rooda",
		Short: "AI-powered OODA loop orchestrator",
		Long: `rooda orchestrates AI coding agents through structured OODA (Observe-Orient-Decide-Act) 
iteration loops to autonomously build, plan, and maintain software from specifications.`,
		Version:       fmt.Sprintf("%s (commit: %s, built: %s)", Version, CommitSHA, BuildDate),
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand, show help
			cmd.Help()
		},
	}

	// Persistent flags available to all commands
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./rooda-config.yml)")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (sets show_ai_output=true and log_level=debug)")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all non-error output")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level (debug, info, warn, error)")

	// Mark verbose and quiet as mutually exclusive
	cmd.MarkFlagsMutuallyExclusive("verbose", "quiet")

	// Add subcommands
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newInfoCommand())
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newRunCommand())

	return cmd
}

func Execute() error {
	return newRootCommand().Execute()
}
