package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jomadu/rooda/internal/config"
	"github.com/jomadu/rooda/internal/loop"
	"github.com/jomadu/rooda/internal/observability"
	"github.com/jomadu/rooda/internal/prompt"
	"github.com/spf13/cobra"
)

func newRunCommand() *cobra.Command {
	var execFlags ExecutionFlags

	cmd := &cobra.Command{
		Use:   "run <procedure>",
		Short: "Execute a procedure",
		Long:  `Execute a named OODA loop procedure with the specified configuration.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate execution flags
			return ValidateExecutionFlags(&execFlags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			procedureName := args[0]
			return runProcedure(cmd, procedureName, &execFlags)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// Load config to get procedure names
			flags := buildCLIFlags(&execFlags, "")
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

	// Add all execution flags
	AddExecutionFlags(cmd, &execFlags)

	return cmd
}

func runProcedure(cmd *cobra.Command, procedureName string, execFlags *ExecutionFlags) error {
	// Build CLIFlags from cobra flags
	flags := buildCLIFlags(execFlags, procedureName)

	// Handle dry-run mode
	if execFlags.DryRun {
		return runDryRunMode(cmd, flags)
	}

	// Load configuration
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate procedure exists
	proc, exists := cfg.Procedures[procedureName]
	if !exists {
		return fmt.Errorf("unknown procedure '%s'\n\nRun 'rooda list' to see available procedures", procedureName)
	}

	// Resolve AI command
	aiCmd, err := config.ResolveAICommand(*cfg, procedureName, flags)
	if err != nil {
		return fmt.Errorf("failed to resolve AI command: %w", err)
	}

	// Determine max iterations
	var maxIterations *int
	if execFlags.Unlimited {
		maxIterations = nil
	} else if execFlags.MaxIterations > 0 {
		maxIterations = &execFlags.MaxIterations
	} else if proc.DefaultMaxIterations != nil {
		maxIterations = proc.DefaultMaxIterations
	} else if cfg.Loop.DefaultMaxIterations != nil {
		maxIterations = cfg.Loop.DefaultMaxIterations
	} else {
		defaultMax := config.DefaultMaxIterations
		maxIterations = &defaultMax
	}

	// Determine iteration timeout
	var iterationTimeout *int
	if proc.IterationTimeout != nil {
		iterationTimeout = proc.IterationTimeout
	} else {
		iterationTimeout = cfg.Loop.IterationTimeout
	}

	// Determine max output buffer
	maxOutputBuffer := cfg.Loop.MaxOutputBuffer
	if proc.MaxOutputBuffer != nil {
		maxOutputBuffer = *proc.MaxOutputBuffer
	}

	// Determine log level
	logLevel := cfg.Loop.LogLevel
	if verbose {
		logLevel = config.LogLevelDebug
	} else if quiet {
		logLevel = config.LogLevelError
	} else if logLevel != "" {
		logLevel = logLevel
	}

	// Determine show AI output
	showAIOutput := cfg.Loop.ShowAIOutput
	if verbose {
		showAIOutput = true
	}

	// Create logger
	logger := observability.NewLogger(logLevel, cfg.Loop.LogTimestampFormat, time.Now())

	// Create iteration state
	state := &loop.IterationState{
		Iteration:           0,
		MaxIterations:       maxIterations,
		IterationTimeout:    iterationTimeout,
		MaxOutputBuffer:     maxOutputBuffer,
		ConsecutiveFailures: 0,
		FailureThreshold:    cfg.Loop.FailureThreshold,
		StartedAt:           time.Now(),
		Status:              loop.StatusRunning,
		ProcedureName:       procedureName,
		Stats:               loop.IterationStats{},
	}

	// Join user contexts
	userContext := strings.Join(execFlags.Contexts, "\n\n")

	// Run loop
	status := loop.RunLoop(state, *cfg, aiCmd, userContext, showAIOutput, logger)

	// Map status to exit code (return nil for success, error for failure)
	switch status {
	case loop.StatusSuccess, loop.StatusMaxIters, loop.StatusInterrupted:
		return nil
	case loop.StatusAborted:
		return fmt.Errorf("procedure aborted")
	default:
		return fmt.Errorf("procedure failed with status: %s", status)
	}
}

func runDryRunMode(cmd *cobra.Command, flags config.CLIFlags) error {
	cmd.Println("=== DRY RUN MODE ===")
	cmd.Println()

	// Load and validate configuration
	cmd.Println("Loading configuration...")
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	cmd.Println("✓ Configuration valid")
	cmd.Println()

	// Validate procedure exists
	proc, exists := cfg.Procedures[flags.ProcedureName]
	if !exists {
		return fmt.Errorf("unknown procedure '%s'\n\nRun 'rooda list' to see available procedures", flags.ProcedureName)
	}
	cmd.Printf("✓ Procedure '%s' found\n\n", flags.ProcedureName)

	// Assemble prompt
	cmd.Println("Assembling prompt...")
	userContext := strings.Join(flags.Contexts, "\n\n")

	// Determine config directory for fragment resolution
	configPath := "./rooda-config.yml"
	if flags.ConfigPath != "" {
		configPath = flags.ConfigPath
	}
	configDir := "."
	if dir := filepath.Dir(configPath); dir != "" {
		configDir = dir
	}

	// For dry-run, we don't have iteration state, so pass nil
	assembledPrompt, err := prompt.AssemblePrompt(proc, userContext, configDir, nil)
	if err != nil {
		return fmt.Errorf("prompt assembly failed: %w", err)
	}
	cmd.Printf("✓ Prompt assembled (%d characters)\n", len(assembledPrompt))
	cmd.Println()

	// Display assembled prompt
	cmd.Println("--- Assembled Prompt ---")
	cmd.Println()
	cmd.Print(assembledPrompt)
	cmd.Println("--- End Prompt ---")
	cmd.Println()

	cmd.Println("✓ Dry-run validation passed")
	return nil
}

func buildCLIFlags(execFlags *ExecutionFlags, procedureName string) config.CLIFlags {
	flags := config.CLIFlags{
		ProcedureName: procedureName,
		ConfigPath:    cfgFile,
		Verbose:       verbose,
		Quiet:         quiet,
		LogLevel:      logLevel,
		DryRun:        execFlags.DryRun,
		Unlimited:     execFlags.Unlimited,
		AICmd:         execFlags.AICmd,
		AICmdAlias:    execFlags.AICmdAlias,
		Contexts:      execFlags.Contexts,
	}

	if execFlags.MaxIterations > 0 {
		flags.MaxIterations = &execFlags.MaxIterations
	}

	// Convert OODA fragment overrides
	if len(execFlags.ObserveFragments) > 0 {
		flags.ObserveFragments = execFlags.ObserveFragments
	}
	if len(execFlags.OrientFragments) > 0 {
		flags.OrientFragments = execFlags.OrientFragments
	}
	if len(execFlags.DecideFragments) > 0 {
		flags.DecideFragments = execFlags.DecideFragments
	}
	if len(execFlags.ActFragments) > 0 {
		flags.ActFragments = execFlags.ActFragments
	}

	return flags
}
