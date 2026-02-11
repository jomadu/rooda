package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jomadu/rooda/internal/config"
	"github.com/jomadu/rooda/internal/loop"
	"github.com/jomadu/rooda/internal/observability"
	"github.com/jomadu/rooda/internal/procedures"
	"github.com/jomadu/rooda/internal/prompt"
)

func init() {
	// Register built-in procedures with config package
	config.BuiltInProceduresFunc = procedures.BuiltInProcedures
}

var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

const (
	ExitSuccess        = 0
	ExitUserError      = 1
	ExitConfigError    = 2
	ExitExecutionError = 3
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitUserError)
	}
}

func parseFlags() config.CLIFlags {
	var flags config.CLIFlags
	var maxIterLong, maxIterShort int

	// Define flags
	fs := flag.NewFlagSet("rooda", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	// Info flags
	fs.BoolVar(&flags.ShowVersion, "version", false, "Print version information")
	fs.BoolVar(&flags.ShowHelp, "help", false, "Show help")
	fs.BoolVar(&flags.ListProcedures, "list-procedures", false, "List all available procedures")

	// Loop control flags
	fs.IntVar(&maxIterLong, "max-iterations", -1, "Maximum iterations (>= 1)")
	fs.IntVar(&maxIterShort, "n", -1, "Maximum iterations (short form)")
	fs.BoolVar(&flags.Unlimited, "unlimited", false, "Run unlimited iterations")
	fs.BoolVar(&flags.Unlimited, "u", false, "Run unlimited iterations (short form)")
	fs.BoolVar(&flags.DryRun, "dry-run", false, "Validate without executing")
	fs.BoolVar(&flags.DryRun, "d", false, "Validate without executing (short form)")

	// AI command flags
	fs.StringVar(&flags.AICmd, "ai-cmd", "", "Override AI command")
	fs.StringVar(&flags.AICmdAlias, "ai-cmd-alias", "", "Override AI command using alias")

	// Output control flags
	fs.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose output")
	fs.BoolVar(&flags.Verbose, "v", false, "Enable verbose output (short form)")
	fs.BoolVar(&flags.Quiet, "quiet", false, "Suppress non-error output")
	fs.BoolVar(&flags.Quiet, "q", false, "Suppress non-error output (short form)")
	fs.StringVar(&flags.LogLevel, "log-level", "", "Set log level (debug, info, warn, error)")

	// Configuration flags
	fs.StringVar(&flags.ConfigPath, "config", "", "Alternate workspace config file path")

	// Custom flag parsing for repeatable flags
	fs.Func("context", "Context file path or inline text (repeatable)", func(s string) error {
		flags.Contexts = append(flags.Contexts, s)
		return nil
	})
	fs.Func("c", "Context file path or inline text (short form, repeatable)", func(s string) error {
		flags.Contexts = append(flags.Contexts, s)
		return nil
	})

	fs.Func("observe", "Observe phase fragment (repeatable)", func(s string) error {
		flags.ObserveFragments = append(flags.ObserveFragments, s)
		return nil
	})
	fs.Func("orient", "Orient phase fragment (repeatable)", func(s string) error {
		flags.OrientFragments = append(flags.OrientFragments, s)
		return nil
	})
	fs.Func("decide", "Decide phase fragment (repeatable)", func(s string) error {
		flags.DecideFragments = append(flags.DecideFragments, s)
		return nil
	})
	fs.Func("act", "Act phase fragment (repeatable)", func(s string) error {
		flags.ActFragments = append(flags.ActFragments, s)
		return nil
	})

	// Parse flags - we need to handle procedure name specially
	// First, extract procedure name if it's the first non-flag arg
	args := os.Args[1:]
	procedureIdx := -1
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") && procedureIdx == -1 {
			flags.ProcedureName = arg
			procedureIdx = i
			break
		}
	}

	// Remove procedure name from args for flag parsing
	flagArgs := args
	if procedureIdx >= 0 {
		flagArgs = append(args[:procedureIdx], args[procedureIdx+1:]...)
	}

	if err := fs.Parse(flagArgs); err != nil {
		if err == flag.ErrHelp {
			flags.ShowHelp = true
			return flags
		}
		os.Exit(ExitUserError)
	}

	// Handle max-iterations from both long and short forms
	if maxIterLong >= 0 {
		flags.MaxIterations = &maxIterLong
	} else if maxIterShort >= 0 {
		flags.MaxIterations = &maxIterShort
	}

	return flags
}

func printGlobalHelp() {
	fmt.Println(`rooda - OODA Loop Framework

USAGE:
  rooda <procedure> [flags]
  rooda --help
  rooda --version
  rooda --list-procedures

LOOP CONTROL FLAGS:
  -n, --max-iterations <n>    Maximum iterations (>= 1)
  -u, --unlimited              Run unlimited iterations
  -d, --dry-run                Validate without executing

AI COMMAND FLAGS:
  --ai-cmd <command>           Override AI command
  --ai-cmd-alias <alias>       Override AI command using alias

PROMPT OVERRIDE FLAGS:
  --observe <value>            Observe phase fragment (repeatable)
  --orient <value>             Orient phase fragment (repeatable)
  --decide <value>             Decide phase fragment (repeatable)
  --act <value>                Act phase fragment (repeatable)

OUTPUT CONTROL FLAGS:
  -v, --verbose                Enable verbose output
  -q, --quiet                  Suppress non-error output
  --log-level <level>          Set log level (debug, info, warn, error)

CONFIGURATION FLAGS:
  --config <path>              Alternate workspace config file path
  -c, --context <value>        Context file path or inline text (repeatable)

INFO FLAGS:
  --help                       Show this help
  --version                    Print version information
  --list-procedures            List all available procedures

EXAMPLES:
  rooda build                           # Run build procedure
  rooda build --max-iterations 5        # Run with iteration limit
  rooda build --dry-run                 # Validate without executing
  rooda build --context task.md         # Add context from file
  rooda build --ai-cmd-alias claude     # Use Claude AI
  rooda --list-procedures               # List available procedures

For procedure-specific help:
  rooda <procedure> --help`)
}

func printProcedureHelp(procedureName string, flags config.CLIFlags) {
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(ExitConfigError)
	}

	proc, exists := cfg.Procedures[procedureName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Unknown procedure '%s'\n", procedureName)
		fmt.Fprintln(os.Stderr, "\nRun 'rooda --list-procedures' to see available procedures.")
		os.Exit(ExitUserError)
	}

	// Display name or fallback to procedure name
	display := proc.Display
	if display == "" {
		display = procedureName
	}

	fmt.Printf("rooda %s - %s\n\n", procedureName, display)

	// Description
	if proc.Description != "" {
		fmt.Println("DESCRIPTION:")
		fmt.Printf("  %s\n\n", proc.Description)
	} else if proc.Summary != "" {
		fmt.Println("DESCRIPTION:")
		fmt.Printf("  %s\n\n", proc.Summary)
	}

	// OODA phases
	fmt.Println("OODA PHASES:")
	fmt.Printf("  Observe:  %d fragment(s)\n", len(proc.Observe))
	fmt.Printf("  Orient:   %d fragment(s)\n", len(proc.Orient))
	fmt.Printf("  Decide:   %d fragment(s)\n", len(proc.Decide))
	fmt.Printf("  Act:      %d fragment(s)\n\n", len(proc.Act))

	// Configuration overrides
	hasOverrides := false
	fmt.Println("CONFIGURATION:")
	if proc.IterationMode != "" {
		fmt.Printf("  Iteration mode: %s\n", proc.IterationMode)
		hasOverrides = true
	}
	if proc.DefaultMaxIterations != nil {
		fmt.Printf("  Max iterations: %d\n", *proc.DefaultMaxIterations)
		hasOverrides = true
	}
	if proc.IterationTimeout != nil {
		fmt.Printf("  Timeout: %ds\n", *proc.IterationTimeout)
		hasOverrides = true
	}
	if proc.AICmd != "" {
		fmt.Printf("  AI command: %s\n", proc.AICmd)
		hasOverrides = true
	}
	if proc.AICmdAlias != "" {
		fmt.Printf("  AI alias: %s\n", proc.AICmdAlias)
		hasOverrides = true
	}
	if !hasOverrides {
		fmt.Println("  (uses global defaults)")
	}
	fmt.Println()

	// Usage examples
	fmt.Println("USAGE:")
	fmt.Printf("  rooda %s\n", procedureName)
	fmt.Printf("  rooda %s --max-iterations 10\n", procedureName)
	fmt.Printf("  rooda %s --dry-run\n", procedureName)
	fmt.Printf("  rooda %s --context task.md\n\n", procedureName)

	fmt.Println("Run 'rooda --help' for all available flags.")
}

func listProcedures(flags config.CLIFlags) {
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(ExitConfigError)
	}

	if len(cfg.Procedures) == 0 {
		fmt.Println("No procedures defined.")
		return
	}

	fmt.Println("Available procedures:")
	for name, proc := range cfg.Procedures {
		desc := proc.Description
		if desc == "" {
			desc = "(no description)"
		}
		// Truncate long descriptions
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Printf("  %-20s %s\n", name, desc)
	}
}

func runDryRun(flags config.CLIFlags) int {
	fmt.Println("=== DRY RUN MODE ===")
	fmt.Println()

	// 1. Load and validate configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Configuration validation failed: %v\n", err)
		return ExitConfigError
	}
	fmt.Println("✓ Configuration valid")
	fmt.Println()

	// 2. Validate procedure exists
	proc, exists := cfg.Procedures[flags.ProcedureName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Unknown procedure '%s'\n", flags.ProcedureName)
		fmt.Fprintln(os.Stderr, "Run 'rooda --list-procedures' to see available procedures.")
		return ExitUserError
	}
	fmt.Printf("✓ Procedure '%s' found\n\n", flags.ProcedureName)

	// 3. Assemble prompt
	fmt.Println("Assembling prompt...")
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
		fmt.Fprintf(os.Stderr, "Error: Prompt assembly failed: %v\n", err)
		return ExitUserError
	}
	fmt.Printf("✓ Prompt assembled (%d characters)\n", len(assembledPrompt))
	fmt.Println()

	// 4. Display assembled prompt (already has === PHASE === markers)
	fmt.Println("--- Assembled Prompt ---")
	fmt.Println()
	fmt.Print(assembledPrompt)
	fmt.Println("--- End Prompt ---")
	fmt.Println()

	// 5. Display resolved configuration with provenance
	fmt.Println("=== RESOLVED CONFIGURATION ===")
	displayConfigWithProvenance(cfg, flags.ProcedureName)
	fmt.Println("=== END CONFIGURATION ===")
	fmt.Println()

	fmt.Println("✓ Dry-run validation passed")
	return ExitSuccess
}

func runLoop(flags config.CLIFlags) int {
	// Load configuration
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitConfigError
	}

	// Resolve AI command
	aiCmd, err := config.ResolveAICommand(*cfg, flags.ProcedureName, flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitConfigError
	}

	// Determine max iterations
	var maxIterations *int
	if flags.Unlimited {
		maxIterations = nil
	} else if flags.MaxIterations != nil {
		maxIterations = flags.MaxIterations
	} else {
		// Use procedure override or loop default
		proc, ok := cfg.Procedures[flags.ProcedureName]
		if ok && proc.DefaultMaxIterations != nil {
			maxIterations = proc.DefaultMaxIterations
		} else if cfg.Loop.DefaultMaxIterations != nil {
			maxIterations = cfg.Loop.DefaultMaxIterations
		} else {
			defaultMax := config.DefaultMaxIterations
			maxIterations = &defaultMax
		}
	}

	// Determine iteration timeout
	var iterationTimeout *int
	proc, ok := cfg.Procedures[flags.ProcedureName]
	if ok && proc.IterationTimeout != nil {
		iterationTimeout = proc.IterationTimeout
	} else {
		iterationTimeout = cfg.Loop.IterationTimeout
	}

	// Determine max output buffer
	maxOutputBuffer := cfg.Loop.MaxOutputBuffer
	if ok && proc.MaxOutputBuffer != nil {
		maxOutputBuffer = *proc.MaxOutputBuffer
	}

	// Determine log level
	logLevel := cfg.Loop.LogLevel
	if flags.Verbose {
		logLevel = config.LogLevelDebug
	} else if flags.Quiet {
		logLevel = config.LogLevelError
	}

	// Determine show AI output
	showAIOutput := cfg.Loop.ShowAIOutput
	if flags.Verbose {
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
		ProcedureName:       flags.ProcedureName,
		Stats:               loop.IterationStats{},
	}

	// Join user contexts
	userContext := strings.Join(flags.Contexts, "\n\n")

	// Run loop
	status := loop.RunLoop(state, *cfg, aiCmd, userContext, showAIOutput, logger)

	// Map status to exit code
	switch status {
	case loop.StatusSuccess:
		return ExitSuccess
	case loop.StatusMaxIters:
		return ExitSuccess // Max iterations is not an error
	case loop.StatusInterrupted:
		return ExitSuccess // User interrupt is not an error
	case loop.StatusAborted:
		return ExitExecutionError
	default:
		return ExitExecutionError
	}
}

func displayConfigWithProvenance(cfg *config.Config, procedureName string) {
	// Display loop-level configuration
	fmt.Println("Loop Configuration:")
	displaySetting("  iteration_mode", cfg.Loop.IterationMode, cfg.Provenance["loop.iteration_mode"])
	if cfg.Loop.DefaultMaxIterations != nil {
		displaySetting("  default_max_iterations", *cfg.Loop.DefaultMaxIterations, cfg.Provenance["loop.default_max_iterations"])
	}
	displaySetting("  max_output_buffer", cfg.Loop.MaxOutputBuffer, cfg.Provenance["loop.max_output_buffer"])
	displaySetting("  failure_threshold", cfg.Loop.FailureThreshold, cfg.Provenance["loop.failure_threshold"])
	displaySetting("  log_level", cfg.Loop.LogLevel, cfg.Provenance["loop.log_level"])
	displaySetting("  log_timestamp_format", cfg.Loop.LogTimestampFormat, cfg.Provenance["loop.log_timestamp_format"])
	displaySetting("  show_ai_output", cfg.Loop.ShowAIOutput, cfg.Provenance["loop.show_ai_output"])
	if cfg.Loop.AICmd != "" {
		displaySetting("  ai_cmd", cfg.Loop.AICmd, cfg.Provenance["loop.ai_cmd"])
	}

	// Display procedure-specific configuration
	if proc, exists := cfg.Procedures[procedureName]; exists {
		fmt.Printf("\nProcedure '%s' Configuration:\n", procedureName)
		if proc.Display != "" {
			fmt.Printf("  display: %s\n", proc.Display)
		}
		if proc.Summary != "" {
			fmt.Printf("  summary: %s\n", proc.Summary)
		}
		if proc.IterationMode != "" {
			fmt.Printf("  iteration_mode: %s (procedure override)\n", proc.IterationMode)
		}
		if proc.DefaultMaxIterations != nil {
			fmt.Printf("  default_max_iterations: %d (procedure override)\n", *proc.DefaultMaxIterations)
		}
		if proc.IterationTimeout != nil {
			fmt.Printf("  iteration_timeout: %ds (procedure override)\n", *proc.IterationTimeout)
		}
		if proc.MaxOutputBuffer != nil {
			fmt.Printf("  max_output_buffer: %d (procedure override)\n", *proc.MaxOutputBuffer)
		}
		if proc.AICmd != "" {
			fmt.Printf("  ai_cmd: %s (procedure override)\n", proc.AICmd)
		}
		if proc.AICmdAlias != "" {
			fmt.Printf("  ai_cmd_alias: %s (procedure override)\n", proc.AICmdAlias)
		}
		fmt.Printf("  observe_fragments: %d\n", len(proc.Observe))
		fmt.Printf("  orient_fragments: %d\n", len(proc.Orient))
		fmt.Printf("  decide_fragments: %d\n", len(proc.Decide))
		fmt.Printf("  act_fragments: %d\n", len(proc.Act))
	}
}

func displaySetting(name string, value interface{}, source config.ConfigSource) {
	if source.Tier == "" {
		fmt.Printf("%s: %v\n", name, value)
	} else {
		fmt.Printf("%s: %v (from: %s)\n", name, value, formatProvenance(source))
	}
}

func formatProvenance(source config.ConfigSource) string {
	switch source.Tier {
	case config.TierBuiltIn:
		return "built-in default"
	case config.TierGlobal:
		return "global config"
	case config.TierWorkspace:
		return "workspace config"
	case config.TierEnvVar:
		return "environment variable"
	case config.TierCLIFlag:
		return "CLI flag"
	default:
		return string(source.Tier)
	}
}
