package main

import (
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display version, commit SHA, and build date for rooda.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\n", Version)
			cmd.Printf("Commit:  %s\n", CommitSHA)
			cmd.Printf("Built:   %s\n", BuildDate)
		},
	}

	return cmd
}
