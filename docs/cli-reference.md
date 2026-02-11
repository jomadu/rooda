# CLI Reference

## Synopsis

```bash
rooda <command> [flags]
rooda run <procedure> [flags]
rooda list
rooda info <procedure>
rooda version
rooda --help
```

## Commands

### `rooda run <procedure>`

Execute a named OODA loop procedure.

```bash
rooda run build --ai-cmd-alias kiro-cli --max-iterations 5
rooda run agents-sync --ai-cmd-alias claude
```

### `rooda list`

List all available procedures (built-in and custom) with one-line descriptions.

```bash
rooda list
```

### `rooda info <procedure>`

Show detailed information about a specific procedure including metadata, description, OODA phases, and configuration.

```bash
rooda info build
rooda info agents-sync
```

### `rooda version`

Display version number, commit SHA, and build date.

```bash
rooda version
```

Note: `rooda version` also works (cobra convention).

### `rooda --help`

Display usage summary, available commands, and global flags.

```bash
rooda --help
rooda run --help  # Command-specific help
```

## Global Flags

These flags are available for all commands:

**`--config <path>`**  
Specify alternate config file path (default: `./rooda-config.yml`).

```bash
rooda run build --config /path/to/config.yml
```

**`--verbose` / `-v`**  
Enable verbose output (sets `show_ai_output=true` and `log_level=debug`).

```bash
rooda run build -v
```

**`--quiet` / `-q`**  
Suppress all non-error output.

```bash
rooda run build -q
```

**`--log-level <level>`**  
Set log level: `debug`, `info`, `warn`, `error`.

```bash
rooda run build --log-level debug
```

## Run Command Flags

These flags are specific to `rooda run <procedure>`:

### Loop Control

**`--max-iterations <n>` / `-n <n>`**  
Override default max iterations for the procedure. Must be >= 1.

```bash
rooda run build --max-iterations 10
rooda run build -n 10
```

**`--unlimited` / `-u`**  
Set iteration mode to unlimited (runs until success or failure threshold). Overrides `--max-iterations`.

```bash
rooda run build --unlimited
rooda run build -u
```

**`--dry-run` / `-d`**  
Display assembled prompt without executing AI CLI. Validates configuration, prompts, and AI command.

Exit codes:
- 0 if validation passes
- 1 if validation fails

```bash
rooda run build --dry-run
rooda run build -d
```

### AI Command

**`--ai-cmd <command>`**  
Override AI command with direct command string. Takes precedence over `--ai-cmd-alias`.

```bash
rooda run build --ai-cmd "kiro-cli chat"
rooda run build --ai-cmd "claude --no-cache"
```

**`--ai-cmd-alias <alias>`**  
Override AI command using a named alias. Built-in aliases: `kiro-cli`, `claude`, `copilot`, `cursor-agent`.

```bash
rooda run build --ai-cmd-alias kiro-cli
rooda run build --ai-cmd-alias claude
```

### Context

**`--context <value>` / `-c <value>`**  
Pass runtime context to the procedure. Can be file path or inline text. Multiple `--context` flags accumulate.

File existence heuristic: if value exists as file, read it; otherwise treat as inline content.

```bash
# Inline context
rooda run draft-plan-impl-feat --context "Add user authentication"

# File context
rooda run draft-plan-impl-feat --context feature-requirements.md

# Multiple contexts
rooda run build --context "Focus on auth module" --context notes.md
```

### Output Control

**`--verbose` / `-v`**  
Enable verbose output. Sets `show_ai_output=true` and `log_level=debug`.

```bash
rooda run build --verbose
rooda run build -v
```

**`--quiet` / `-q`**  
Suppress all non-error output.

```bash
rooda run build --quiet
rooda run build -q
```

**`--log-level <level>`**  
Set log level. Valid values: `debug`, `info`, `warn`, `error`.

```bash
rooda run build --log-level debug
```

### Configuration

**`--config <path>`**  
Specify alternate workspace config file path. Overrides `./rooda-config.yml`.

Fragment paths in CLI overrides resolve relative to this config file's directory.

```bash
rooda run build --config /path/to/custom-config.yml
```

### Prompt Overrides

Override OODA phase fragments for this execution. Multiple flags accumulate into fragment array. Replaces entire phase array (not appended to config).

**`--observe <value>`**  
Override observe phase fragments.

```bash
rooda run build --observe prompts/observe_custom.md
rooda run build --observe "# Observe\nCustom inline content"
```

**`--orient <value>`**  
Override orient phase fragments.

```bash
rooda run build --orient prompts/orient_custom.md
```

**`--decide <value>`**  
Override decide phase fragments.

```bash
rooda run build --decide prompts/decide_custom.md
```

**`--act <value>`**  
Override act phase fragments.

```bash
rooda run build --act prompts/act_custom.md
```

**Multiple fragments**:
```bash
rooda run build \
  --observe prompts/observe_specs.md \
  --observe prompts/observe_impl.md \
  --orient prompts/orient_custom.md
```

## Exit Codes

| Code | Meaning | Examples |
|------|---------|----------|
| 0 | Success | Procedure completed successfully, dry-run validation passed |
| 1 | User error | Invalid flags, unknown procedure, validation failures |
| 2 | Configuration error | Invalid config file, missing AI command (runtime only, not dry-run) |
| 3 | Execution error | AI CLI failure, iteration timeout |
| 130 | Interrupted | User pressed Ctrl+C (SIGINT) |

## Flag Precedence

CLI flags have highest precedence and override all other configuration sources:

1. CLI flags (highest)
2. Environment variables (`ROODA_*`)
3. Workspace config (`./rooda-config.yml`)
4. Global config (`~/.config/rooda/rooda-config.yml`)
5. Built-in defaults (lowest)

## Mutually Exclusive Flags

- `--verbose` and `--quiet` cannot be used together
- `--max-iterations` and `--unlimited` cannot be used together
- `--ai-cmd` takes precedence over `--ai-cmd-alias` when both provided

## Short Flags

| Short | Long | Description |
|-------|------|-------------|
| `-v` | `--verbose` | Enable verbose output |
| `-q` | `--quiet` | Suppress non-error output |
| `-n` | `--max-iterations` | Set max iterations |
| `-u` | `--unlimited` | Unlimited iterations |
| `-d` | `--dry-run` | Validate without executing |
| `-c` | `--context` | Pass runtime context |

## Examples

### Basic Execution

```bash
# Run build procedure with default settings
rooda run build --ai-cmd-alias kiro-cli

# Run with limited iterations
rooda run build --ai-cmd-alias kiro-cli --max-iterations 3

# Run until success or failure threshold
rooda run build --ai-cmd-alias kiro-cli --unlimited
```

### Dry Run

```bash
# Validate configuration and prompts without executing
rooda run build --ai-cmd-alias kiro-cli --dry-run

# Check exit code
rooda run build --dry-run && echo "Valid" || echo "Invalid"
```

### Context Passing

```bash
# Inline context
rooda draft-plan-impl-feat --ai-cmd-alias kiro-cli \
  --context "Add OAuth2 authentication with Google and GitHub providers"

# File context
rooda draft-plan-impl-feat --ai-cmd-alias kiro-cli \
  --context requirements/auth-feature.md

# Multiple contexts
rooda run build --ai-cmd-alias kiro-cli \
  --context "Focus on authentication module" \
  --context notes/auth-implementation.md
```

### Prompt Overrides

```bash
# Override single phase
rooda run build --ai-cmd-alias kiro-cli \
  --observe prompts/observe_custom.md

# Override multiple phases
rooda run build --ai-cmd-alias kiro-cli \
  --observe prompts/observe_specs.md \
  --observe prompts/observe_impl.md \
  --orient prompts/orient_custom.md

# Inline content
rooda run build --ai-cmd-alias kiro-cli \
  --observe "# Observe\nRead AGENTS.md and specs/"
```

### Verbose Output

```bash
# See all AI output and debug logs
rooda run build --ai-cmd-alias kiro-cli --verbose

# See configuration provenance
rooda run build --ai-cmd-alias kiro-cli --verbose 2>&1 | grep "Configuration loaded"
```

### Custom Config

```bash
# Use alternate config file
rooda run build --config configs/production.yml --ai-cmd-alias kiro-cli

# Fragments resolve relative to config file directory
rooda run build --config /path/to/config.yml \
  --observe custom-prompts/observe.md
```

## Environment Variables

All `ROODA_*` environment variables can be used instead of flags:

```bash
# Set AI command
export ROODA_LOOP_AI_CMD_ALIAS=kiro-cli

# Set log level
export ROODA_LOOP_LOG_LEVEL=debug

# Set default max iterations
export ROODA_LOOP_DEFAULT_MAX_ITERATIONS=10

# Run without flags
rooda build
```

See [Configuration](configuration.md) for all environment variables.

## See Also

- [Procedures](procedures.md) - All built-in procedures
- [Configuration](configuration.md) - Three-tier config system
- [Troubleshooting](troubleshooting.md) - Common errors and solutions
