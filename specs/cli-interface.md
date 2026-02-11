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

- [ ] `rooda run <procedure>` invokes the named procedure with default configuration
- [ ] `rooda --help` displays usage summary, available commands, and global flags
- [ ] `rooda run <procedure> --help` displays procedure-specific help (description, OODA phases, iteration limits)
- [ ] `rooda list` lists all available procedures (built-in and custom) with one-line descriptions
- [ ] `rooda info <procedure>` displays detailed procedure information (metadata, description, OODA phases, configuration)
- [ ] `rooda version` displays version number, commit SHA, and build date
- [ ] `rooda --version` also works (cobra convention)
- [ ] Unknown procedure name produces error: "Unknown procedure '<name>'. Run 'rooda list' to see available procedures."
- [ ] `--max-iterations <n>` overrides default max iterations for the procedure (must be >= 1)
- [ ] `--unlimited` sets iteration mode to unlimited (overrides `--max-iterations`)
- [ ] `--dry-run` displays assembled prompt without executing AI CLI
- [ ] `--dry-run` exits with code 0 if validation passes (config valid, prompts exist, AI command found)
- [ ] `--dry-run` exits with code 1 if validation fails (invalid flags, unknown procedure, missing AI command, invalid config, missing prompt files)
- [ ] `--context <value>` accepts file path or inline text (file existence check determines interpretation)
- [ ] Multiple `--context` flags accumulate (all values processed in order)
- [ ] Context file existence heuristic: if value exists as file, read it; otherwise treat as inline content
- [ ] `--ai-cmd <command>` overrides AI command for this execution (direct command string)
- [ ] `--ai-cmd-alias <alias>` overrides AI command using a named alias
- [ ] `--ai-cmd` takes precedence over `--ai-cmd-alias` when both provided
- [ ] `--verbose` enables verbose output (sets show_ai_output=true and log_level=debug)
- [ ] `--quiet` suppresses all non-error output
- [ ] `--log-level <level>` sets log level (debug, info, warn, error)
- [ ] `--config <path>` specifies alternate workspace config file path
- [ ] Fragment values in CLI overrides (--observe, --orient, --decide, --act) resolve relative to --config file's directory if --config provided, else relative to ./rooda-config.yml directory
- [ ] Multiple `--observe` flags accumulate into fragment array
- [ ] Multiple `--orient` flags accumulate into fragment array
- [ ] Multiple `--decide` flags accumulate into fragment array
- [ ] Multiple `--act` flags accumulate into fragment array
- [ ] OODA phase flags use file existence heuristic (file path vs inline content)
- [ ] Providing any OODA phase flag replaces entire phase array (not appended to config)
- [ ] Fragment order preserved: processed left-to-right as specified on CLI
- [ ] OODA phase override fragments validated at config load time (fail fast before execution)
- [ ] OODA phase validation skipped for `rooda list` and `rooda version` (info commands only)
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
    Command            string            // Command: "run", "list", "info", "version"
    ProcedureName      string            // Procedure to execute (for "run" and "info" commands)
    MaxIterations      *int              // --max-iterations <n>
    Unlimited          bool              // --unlimited
    DryRun             bool              // --dry-run
    Contexts           []string          // --context <value> (file path or inline, multiple allowed)
    AICmd              string            // --ai-cmd <command>
    AICmdAlias         string            // --ai-cmd-alias <alias>
    Verbose            bool              // --verbose
    Quiet              bool              // --quiet
    LogLevel           string            // --log-level <level>
    ConfigPath         string            // --config <path>
    ObserveFragments   []string          // --observe <value> (file path or inline, multiple allowed)
    OrientFragments    []string          // --orient <value> (file path or inline, multiple allowed)
    DecideFragments    []string          // --decide <value> (file path or inline, multiple allowed)
    ActFragments       []string          // --act <value> (file path or inline, multiple allowed)
    ShowHelp           bool              // --help
}
```

Note: `--version` is handled by cobra automatically and shows version info.

### ExitCode

```go
const (
    ExitSuccess         = 0  // Successful execution
    ExitUserError       = 1  // Invalid flags, unknown procedure, validation failures
    ExitConfigError     = 2  // Invalid config file, missing AI command (runtime only, not dry-run)
    ExitExecutionError  = 3  // AI CLI failure, iteration timeout
)
```

## Algorithm

### Main CLI Flow

```
1. Parse command-line arguments into command and flags
2. If command is "version":
   - Display version, commit SHA, and build date
   - Exit 0
3. If command is "list":
   - Load configuration (built-in + global + workspace)
   - List all procedures with descriptions
   - Exit 0
4. If command is "info <procedure>":
   - Load configuration
   - Validate procedure exists
   - Display procedure metadata, description, OODA phases, configuration
   - Exit 0
5. If command is "run <procedure>":
   - Load configuration (built-in + global + workspace + env vars)
   - Validate procedure exists
     - If not: display error with suggestion to run 'rooda list'
     - Exit 1
   - Merge CLI flags with configuration (flags override config)
   - Validate merged configuration
     - Check mutually exclusive flags
     - Validate flag value constraints
     - Resolve AI command (see configuration.md AI Command Resolution)
     - If no AI command configured: error with guidance, exit 2
     - Validate OODA phase override files exist (if provided)
     - If errors: display clear messages, exit 1 or 2
   - Execute procedure with merged configuration
   - Exit with appropriate code based on outcome
6. If no command or --help:
   - Display global help
   - Exit 0
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

For context values (Contexts array):
  For each value in array:
    If file exists at path:
      Read file content and use as context
    Else:
      Use value directly as inline content

For OODA phase fragments (ObserveFragments, OrientFragments, DecideFragments, ActFragments):
  If CLI provides any fragments for a phase:
    Replace entire phase array (do not merge with config)
  For each fragment value:
    If file exists at path:
      Create FragmentAction with {path: <file>, content: "", parameters: nil}
    Else:
      Create FragmentAction with {path: "", content: <value>, parameters: nil}
  Preserve order from CLI arguments (left to right)
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
Error: Unknown procedure 'unknown-proc'. Run 'rooda list' to see available procedures.
```

Exit code: 1

### Mutually Exclusive Flags

```bash
$ rooda run build --verbose --quiet
Error: --verbose and --quiet are mutually exclusive.
```

Exit code: 1

```bash
$ rooda run build --max-iterations 10 --unlimited
Error: --max-iterations and --unlimited are mutually exclusive.
```

Exit code: 1

### Invalid Flag Value

```bash
$ rooda run build --max-iterations 0
Error: --max-iterations must be >= 1.
```

Exit code: 1

```bash
$ rooda run build --log-level invalid
Error: Invalid log level 'invalid'. Valid levels: debug, info, warn, error.
```

Exit code: 1

### Missing AI Command

```bash
$ rooda run build
Error: No AI command configured. Set one of:
  - CLI flag: --ai-cmd <command> or --ai-cmd-alias <alias>
  - Environment: ROODA_LOOP_AI_CMD or ROODA_LOOP_AI_CMD_ALIAS
  - Config file: loop.ai_cmd or loop.ai_cmd_alias
  - Procedure-level: procedures.<name>.ai_cmd or procedures.<name>.ai_cmd_alias
```

Exit code: 2

### Empty Inline Content

```bash
$ rooda run build --observe ""
Error: Empty inline content not allowed for --observe flag.
```

Exit code: 1

```bash
$ rooda run build --context ""
Error: Empty inline content not allowed for --context flag.
```

Exit code: 1

### Ambiguous Filename

If user wants inline content "file.md" but a file named "file.md" exists, the file wins (file existence heuristic).

```bash
$ rooda run build --observe "file.md"
# If file.md exists: uses file content
# If file.md doesn't exist: uses "file.md" as inline content
```

To force inline content that looks like a filename, ensure the file doesn't exist or use a non-existent absolute path.

### Dry-Run Validation Success

```bash
$ rooda run build --dry-run
[21:00:15.200] INFO Dry-run mode enabled procedure=build
[21:00:15.201] INFO Configuration validated procedure=build
[21:00:15.202] INFO Prompt files validated procedure=build
[21:00:15.203] INFO AI command found path=/usr/local/bin/kiro-cli procedure=build
[21:00:15.204] INFO Assembled prompt (1234 chars) procedure=build
--- OBSERVE ---
...prompt content...
--- END ---
```

Exit code: 0

### Dry-Run Validation Failure

```bash
$ rooda run build --dry-run --ai-cmd nonexistent
[21:00:15.200] INFO Dry-run mode enabled procedure=build
[21:00:15.201] ERROR AI command binary not found path=nonexistent procedure=build
Error: AI command binary not found: nonexistent
```

Exit code: 1

```bash
$ rooda run build --dry-run --observe missing.md
[21:00:15.200] INFO Dry-run mode enabled procedure=build
[21:00:15.201] ERROR Observe file not found path=missing.md procedure=build
Error: Observe file not found: missing.md
```

Exit code: 1

## Dependencies

- **configuration.md** — Defines Config, LoopConfig, Procedure structures and merge logic
- **iteration-loop.md** — Defines iteration modes, max iterations, timeout behavior
- **prompt-composition.md** — Defines fragment processing and how --context and OODA phase overrides are used
- **procedures.md** — Defines fragment array structure and template system
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
$ rooda run build
# Executes 'build' procedure with default configuration
```

### Override Max Iterations

```bash
$ rooda run build --max-iterations 10
# Runs build procedure with max 10 iterations
```

### Unlimited Iterations

```bash
$ rooda run build --unlimited
# Runs build procedure until convergence (no iteration limit)
```

### Dry Run

```bash
$ rooda run build --dry-run
# Displays assembled prompt without executing AI CLI
```

### Inject Context

```bash
$ rooda run build --context "Focus on the auth module"
# Injects inline context into prompt composition
```

```bash
$ rooda run build --context task.md
# Reads context from task.md file and injects into prompt
```

```bash
$ rooda run build --context task.md --context "Additional notes"
# Mixed: file content + inline content (processed in order)
```

### Override AI Command

```bash
$ rooda run build --ai-cmd "kiro-cli chat"
# Uses direct command string
```

```bash
$ rooda run build --ai-cmd-alias claude
# Uses 'claude' alias from configuration
```

### Verbose Output

```bash
$ rooda run build --verbose
# Sets show_ai_output=true and log_level=debug
# Streams AI CLI output and shows debug-level logs
```

### Override OODA Phase

```bash
$ rooda run build --observe custom.md
# Single file fragment (replaces entire observe phase)
```

```bash
$ rooda run build --observe file1.md --observe file2.md
# Multiple file fragments (replaces entire observe phase)
```

```bash
$ rooda run build --observe "Focus on auth module"
# Inline content fragment (replaces entire observe phase)
```

```bash
$ rooda run build --observe custom.md --observe "Additional instructions"
# Mixed: file + inline content (replaces entire observe phase)
```

### Multiple Contexts (Unified Flag)

```bash
$ rooda run build --context "Focus on auth" --context "Prioritize security"
# Both inline contexts injected into prompt (in order)
```

```bash
$ rooda run build --context task.md --context notes.md
# Both file contents injected into prompt (in order)
```

```bash
$ rooda run build --context task.md --context "Focus on auth" --context notes.md
# Mixed: file + inline + file (all processed in order)
```

### Procedure Help

```bash
$ rooda run build --help
# Displays build procedure details
```

### List Procedures

```bash
$ rooda list
# Lists all available procedures with descriptions
```

### Version

```bash
$ rooda version
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

**Why unify `--context` and `--context-file` into a single flag?**
Consistency with OODA phase flags and reduced cognitive load. The file existence heuristic is intuitive: if it looks like a file and exists, it's a file; otherwise it's inline content. This eliminates the need to remember two separate flags.

**Why use file existence heuristic for context and OODA phases?**
Intuitive and requires no special syntax. Users naturally specify file paths for files and text for inline content. The heuristic makes the right choice 99% of the time without requiring quotes or prefixes.

**Why do OODA phase overrides replace the entire phase array?**
Predictability. Users specify complete phase composition, not partial modifications. Element-by-element merging would create ambiguity about which fragments come from config vs CLI.

**Why support repeatable flags for OODA phases?**
Standard CLI convention (matches `--context` pattern). Enables composing phases from multiple fragments without requiring array syntax in shell.

**Why preserve fragment order from CLI arguments?**
Predictable fragment composition. Left-to-right order matches user intent and makes debugging easier.

**Why doesn't CLI support template parameters for fragments?**
Template parameters are a config-only feature. CLI focuses on simple overrides (file paths or inline content). Complex parameterization belongs in configuration files where it can be documented and versioned.

**What about ambiguous filenames?**
File existence wins. This is a design tradeoff for simplicity. If a user wants inline content "file.md" but a file exists with that name, the file will be used. To force inline content, ensure the file doesn't exist or use a non-existent absolute path.

**Why mutually exclusive `--verbose` and `--quiet`?**
These represent opposite intents. Allowing both would create ambiguity about which takes precedence.

**What does `--verbose` actually do?**
Sets two configuration values: `show_ai_output=true` (streams AI CLI output to console) and `log_level=debug` (shows debug-level logs). This provides maximum visibility into loop execution.

**Why does dry-run use exit code 1 for all validation failures?**
Dry-run is a pre-flight check tool. All validation failures (missing files, invalid config, missing AI command) are actionable by the user before execution. Using a single exit code (1) simplifies scripting: 0 means "ready to run", 1 means "fix something first". Exit code 2 is reserved for runtime configuration errors that occur during actual execution.

**Why different exit codes for different error types?**
Enables scripting and CI/CD integration to distinguish user errors (fix the command) from configuration errors (fix the config) from execution errors (retry or investigate).

**Why support both `--flag=value` and `--flag value`?**
POSIX convention — users expect both forms to work.

**Why short flags for common options?**
Reduces typing for frequently used flags. Only the most common flags get short versions to avoid namespace pollution.

**Short Flag Policy:**
Short flags are reserved for the most frequently used operations only. Future flags will use long form only to avoid namespace pollution. Current short flags (`-v`, `-q`, `-n`, `-u`, `-d`, `-c`, `-h`) are considered stable and will not change.

### Flag Categories

**Loop Control:**
- `--max-iterations`, `-n`
- `--unlimited`, `-u`
- `--dry-run`, `-d`

**AI Command:**
- `--ai-cmd`
- `--ai-cmd-alias`

**Prompt Overrides:**
- `--observe` (repeatable, file or inline)
- `--orient` (repeatable, file or inline)
- `--decide` (repeatable, file or inline)
- `--act` (repeatable, file or inline)
- `--context`, `-c` (repeatable, file or inline)

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
