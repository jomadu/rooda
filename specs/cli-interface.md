# CLI Interface

## Job to be Done

Expose all framework capabilities through a CLI that supports named procedures, explicit OODA phase flags, global options, and helpful error messages.

The developer wants to invoke rooda procedures with minimal typing, override configuration at runtime, get clear feedback when something goes wrong, and discover available procedures and flags through help text — all without reading documentation.

## Activities

1. Parse command-line arguments into procedure name, flags, and context
2. Validate procedure name exists in merged configuration
3. Resolve flag values and merge with configuration (flags have highest precedence)
4. Display help text when requested or when invoked incorrectly
5. Execute the specified procedure with resolved configuration
6. Report errors with actionable messages and exit codes

## Acceptance Criteria

- [ ] `rooda <procedure>` invokes the named procedure with default configuration
- [ ] `rooda --help` displays usage summary, global flags, and available procedures
- [ ] `rooda <procedure> --help` displays procedure-specific help (description, OODA phases, iteration limits)
- [ ] `rooda --list-procedures` lists all available procedures (built-in and custom) with one-line descriptions
- [ ] `rooda --version` displays version number and build information
- [ ] Unknown procedure name produces error: "Unknown procedure '<name>'. Run 'rooda --list-procedures' to see available procedures."
- [ ] `--max-iterations <n>` overrides default max iterations for the procedure (must be >= 1)
- [ ] `--unlimited` sets iteration mode to unlimited (overrides `--max-iterations`)
- [ ] `--dry-run` displays assembled prompt without executing AI CLI
- [ ] `--context <text>` injects user-provided context into prompt composition
- [ ] `--context-file <path>` reads context from file and injects into prompt composition
- [ ] `--ai-cmd <command>` overrides AI command for this execution (direct command string)
- [ ] `--ai-cmd-alias <alias>` overrides AI command using a named alias
- [ ] `--ai-cmd` takes precedence over `--ai-cmd-alias` when both provided
- [ ] `--verbose` enables verbose output (streams AI CLI output, shows provenance)
- [ ] `--quiet` suppresses all non-error output
- [ ] `--log-level <level>` sets log level (debug, info, warn, error)
- [ ] `--config <path>` specifies alternate workspace config file path
- [ ] `--observe <file>`, `--orient <file>`, `--decide <file>`, `--act <file>` override individual OODA phase prompt files
- [ ] Multiple `--context` flags accumulate (all contexts injected)
- [ ] `--verbose` and `--quiet` are mutually exclusive (error if both provided)
- [ ] `--max-iterations` and `--unlimited` are mutually exclusive (error if both provided)
- [ ] Invalid flag values produce clear error messages with expected format
- [ ] Missing required AI command produces error listing all ways to configure it
- [ ] Exit code 0 on success, non-zero on failure
- [ ] Exit code 1 for user errors (invalid flags, unknown procedure)
- [ ] Exit code 2 for configuration errors (invalid config file, missing AI command)
- [ ] Exit code 3 for execution errors (AI CLI failure, iteration timeout)
- [ ] Flag parsing follows POSIX conventions (supports `--flag=value` and `--flag value`)
- [ ] Short flags supported for common options: `-v` (verbose), `-q` (quiet), `-n` (max-iterations), `-u` (unlimited), `-d` (dry-run), `-c` (context)
- [ ] Help text includes examples for common use cases
- [ ] Help text groups flags by category (Loop Control, AI Command, Prompt Overrides, Output Control, Configuration)

## Data Structures

### CLIArgs

Parsed command-line arguments before merging with configuration.

```go
type CLIArgs struct {
    ProcedureName    string            // Procedure to execute
    MaxIterations    *int              // --max-iterations <n>
    Unlimited        bool              // --unlimited
    DryRun           bool              // --dry-run
    Contexts         []string          // --context <text> (multiple allowed)
    ContextFiles     []string          // --context-file <path> (multiple allowed)
    AICmd            string            // --ai-cmd <command>
    AICmdAlias       string            // --ai-cmd-alias <alias>
    Verbose          bool              // --verbose
    Quiet            bool              // --quiet
    LogLevel         string            // --log-level <level>
    ConfigPath       string            // --config <path>
    ObserveFile      string            // --observe <file>
    OrientFile       string            // --orient <file>
    DecideFile       string            // --decide <file>
    ActFile          string            // --act <file>
    ShowHelp         bool              // --help
    ShowVersion      bool              // --version
    ListProcedures   bool              // --list-procedures
}
```

### ExitCode

```go
const (
    ExitSuccess         = 0  // Successful execution
    ExitUserError       = 1  // Invalid flags, unknown procedure
    ExitConfigError     = 2  // Invalid config file, missing AI command
    ExitExecutionError  = 3  // AI CLI failure, iteration timeout
)
```

## Algorithm

### Main CLI Flow

```
1. Parse command-line arguments into CLIArgs
2. If ShowHelp:
   - If ProcedureName empty: display global help
   - Else: display procedure-specific help
   - Exit 0
3. If ShowVersion:
   - Display version and build info
   - Exit 0
4. If ListProcedures:
   - Load configuration (built-in + global + workspace)
   - List all procedures with descriptions
   - Exit 0
5. If ProcedureName empty:
   - Display error: "No procedure specified. Run 'rooda --help' for usage."
   - Exit 1
6. Load configuration (built-in + global + workspace + env vars)
7. Validate ProcedureName exists in configuration
   - If not: display error with suggestion to run --list-procedures
   - Exit 1
8. Merge CLIArgs with configuration (flags override config)
9. Validate merged configuration
   - Check mutually exclusive flags
   - Validate flag value constraints
   - Verify AI command configured
   - If errors: display clear messages, exit 1 or 2
10. Execute procedure with merged configuration
11. Exit with appropriate code based on outcome
```

### Flag Precedence Resolution

When merging CLIArgs with Config:

```
For each setting:
  If CLI flag provided:
    Use CLI flag value (highest precedence)
  Else if environment variable set:
    Use environment variable value
  Else if workspace config defines it:
    Use workspace config value
  Else if global config defines it:
    Use global config value
  Else:
    Use built-in default value
```

## Edge Cases

### No Procedure Specified

```bash
$ rooda
Error: No procedure specified. Run 'rooda --help' for usage.
```

Exit code: 1

### Unknown Procedure

```bash
$ rooda unknown-proc
Error: Unknown procedure 'unknown-proc'. Run 'rooda --list-procedures' to see available procedures.
```

Exit code: 1

### Mutually Exclusive Flags

```bash
$ rooda build --verbose --quiet
Error: --verbose and --quiet are mutually exclusive.
```

Exit code: 1

```bash
$ rooda build --max-iterations 10 --unlimited
Error: --max-iterations and --unlimited are mutually exclusive.
```

Exit code: 1

### Invalid Flag Value

```bash
$ rooda build --max-iterations 0
Error: --max-iterations must be >= 1.
```

Exit code: 1

```bash
$ rooda build --log-level invalid
Error: Invalid log level 'invalid'. Valid levels: debug, info, warn, error.
```

Exit code: 1

### Missing AI Command

```bash
$ rooda build
Error: No AI command configured. Set one of:
  - CLI flag: --ai-cmd <command> or --ai-cmd-alias <alias>
  - Environment: ROODA_LOOP_AI_CMD or ROODA_LOOP_AI_CMD_ALIAS
  - Config file: loop.ai_cmd or loop.ai_cmd_alias
  - Procedure-level: procedures.<name>.ai_cmd or procedures.<name>.ai_cmd_alias
```

Exit code: 2

### Context File Not Found

```bash
$ rooda build --context-file missing.txt
Error: Context file not found: missing.txt
```

Exit code: 1

### OODA Phase File Not Found

```bash
$ rooda build --observe custom-observe.md
Error: Observe file not found: custom-observe.md
```

Exit code: 1

## Dependencies

- **configuration.md** — Defines Config, LoopConfig, Procedure structures and merge logic
- **iteration-loop.md** — Defines iteration modes, max iterations, timeout behavior
- **prompt-composition.md** — Defines how --context and OODA phase overrides are used
- **ai-cli-integration.md** — Defines AI command resolution and execution

## Implementation Mapping

**Source files:**
- `cmd/rooda/main.go` — CLI entry point, argument parsing
- `internal/cli/parser.go` — Flag parsing logic
- `internal/cli/help.go` — Help text generation
- `internal/cli/validator.go` — Flag validation and mutual exclusion checks
- `internal/config/merge.go` — Configuration merging with CLI flag precedence

**Related specs:**
- `configuration.md` — Configuration system
- `iteration-loop.md` — Loop execution
- `prompt-composition.md` — Prompt assembly
- `ai-cli-integration.md` — AI command execution

## Examples

### Basic Procedure Invocation

```bash
$ rooda build
# Executes 'build' procedure with default configuration
```

### Override Max Iterations

```bash
$ rooda build --max-iterations 10
# Runs build procedure with max 10 iterations
```

### Unlimited Iterations

```bash
$ rooda build --unlimited
# Runs build procedure until convergence (no iteration limit)
```

### Dry Run

```bash
$ rooda build --dry-run
# Displays assembled prompt without executing AI CLI
```

### Inject Context

```bash
$ rooda build --context "Focus on the auth module"
# Injects context into prompt composition
```

```bash
$ rooda build --context-file task.md
# Reads context from task.md and injects into prompt
```

### Override AI Command

```bash
$ rooda build --ai-cmd "kiro-cli chat"
# Uses direct command string
```

```bash
$ rooda build --ai-cmd-alias claude
# Uses 'claude' alias from configuration
```

### Verbose Output

```bash
$ rooda build --verbose
# Streams AI CLI output, shows configuration provenance
```

### Override OODA Phase

```bash
$ rooda build --observe custom-observe.md
# Uses custom observe prompt file
```

### Multiple Contexts

```bash
$ rooda build --context "Focus on auth" --context "Prioritize security"
# Both contexts injected into prompt
```

### Procedure Help

```bash
$ rooda build --help
# Displays build procedure details
```

### List Procedures

```bash
$ rooda --list-procedures
# Lists all available procedures with descriptions
```

### Version

```bash
$ rooda --version
# Displays version and build info
```

## Notes

### Design Rationale

**Why named procedures instead of explicit OODA flags?**
Procedures are the primary abstraction — they encapsulate OODA phase composition, iteration limits, and use cases. Explicit OODA flags (`--observe`, `--orient`, etc.) are overrides for customization, not the primary interface.

**Why allow both `--ai-cmd` and `--ai-cmd-alias`?**
Direct commands (`--ai-cmd`) are useful for one-off testing or custom tools. Aliases (`--ai-cmd-alias`) are better for shared team configurations. Supporting both provides flexibility.

**Why accumulate multiple `--context` flags?**
Users may want to inject multiple independent contexts (e.g., task description + architectural constraints). Accumulating them is more flexible than requiring a single concatenated string.

**Why mutually exclusive `--verbose` and `--quiet`?**
These represent opposite intents. Allowing both would create ambiguity about which takes precedence.

**Why different exit codes for different error types?**
Enables scripting and CI/CD integration to distinguish user errors (fix the command) from configuration errors (fix the config) from execution errors (retry or investigate).

**Why support both `--flag=value` and `--flag value`?**
POSIX convention — users expect both forms to work.

**Why short flags for common options?**
Reduces typing for frequently used flags. Only the most common flags get short versions to avoid namespace pollution.

### Flag Categories

**Loop Control:**
- `--max-iterations`, `-n`
- `--unlimited`, `-u`
- `--dry-run`, `-d`

**AI Command:**
- `--ai-cmd`
- `--ai-cmd-alias`

**Prompt Overrides:**
- `--observe`
- `--orient`
- `--decide`
- `--act`
- `--context`, `-c`
- `--context-file`

**Output Control:**
- `--verbose`, `-v`
- `--quiet`, `-q`
- `--log-level`

**Configuration:**
- `--config`

**Help & Info:**
- `--help`, `-h`
- `--version`
- `--list-procedures`
