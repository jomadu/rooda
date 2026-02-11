package main

import (
	"github.com/spf13/cobra"
)

// ExecutionFlags holds all execution-related flags
type ExecutionFlags struct {
	// Loop control
	MaxIterations int
	Unlimited     bool
	DryRun        bool

	// AI command
	AICmd      string
	AICmdAlias string

	// Context
	Contexts []string

	// OODA phase overrides
	ObserveFragments []string
	OrientFragments  []string
	DecideFragments  []string
	ActFragments     []string
}

// AddExecutionFlags adds all execution flags to a command
func AddExecutionFlags(cmd *cobra.Command, flags *ExecutionFlags) {
	// Loop control flags
	cmd.Flags().IntVarP(&flags.MaxIterations, "max-iterations", "n", 0, "maximum number of iterations (must be >= 1)")
	cmd.Flags().BoolVarP(&flags.Unlimited, "unlimited", "u", false, "run until SUCCESS signal or failure threshold")
	cmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "d", false, "display assembled prompt without executing")

	// AI command flags
	cmd.Flags().StringVar(&flags.AICmd, "ai-cmd", "", "AI command to use (direct command string)")
	cmd.Flags().StringVar(&flags.AICmdAlias, "ai-cmd-alias", "", "AI command alias name")

	// Context flags (repeatable)
	cmd.Flags().StringArrayVarP(&flags.Contexts, "context", "c", nil, "inject context (file path or inline text, repeatable)")

	// OODA phase override flags (repeatable)
	cmd.Flags().StringArrayVar(&flags.ObserveFragments, "observe", nil, "observe phase fragment (file path or inline, repeatable)")
	cmd.Flags().StringArrayVar(&flags.OrientFragments, "orient", nil, "orient phase fragment (file path or inline, repeatable)")
	cmd.Flags().StringArrayVar(&flags.DecideFragments, "decide", nil, "decide phase fragment (file path or inline, repeatable)")
	cmd.Flags().StringArrayVar(&flags.ActFragments, "act", nil, "act phase fragment (file path or inline, repeatable)")

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("max-iterations", "unlimited")
}

// ValidateExecutionFlags validates execution flags
func ValidateExecutionFlags(flags *ExecutionFlags) error {
	// Validate max-iterations
	if flags.MaxIterations != 0 && flags.MaxIterations < 1 {
		return ErrInvalidMaxIterations
	}

	// Validate contexts are not empty
	for _, ctx := range flags.Contexts {
		if ctx == "" {
			return ErrEmptyContext
		}
	}

	// Validate OODA fragments are not empty
	for _, frag := range flags.ObserveFragments {
		if frag == "" {
			return ErrEmptyFragment
		}
	}
	for _, frag := range flags.OrientFragments {
		if frag == "" {
			return ErrEmptyFragment
		}
	}
	for _, frag := range flags.DecideFragments {
		if frag == "" {
			return ErrEmptyFragment
		}
	}
	for _, frag := range flags.ActFragments {
		if frag == "" {
			return ErrEmptyFragment
		}
	}

	return nil
}
