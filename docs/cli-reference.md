# CLI Reference

## Synopsis

```bash
rooda [flags] <procedure>
rooda --help
rooda --version
rooda --list-procedures
```

## Global Flags

### Information Commands

**`--help`**  
Display usage summary, global flags, and available procedures.

```bash
rooda --help
rooda build --help  # Procedure-specific help
```

**`--version`**  
Display version number and build information.

```bash
rooda --version
```

**`--list-procedures`**  
List all available procedures (built-in and custom) with one-line descriptions.

```bash
rooda --list-procedures
```

### Loop Control

**`--max-iterations <n>` / `-n <n>`**  
Override default max iterations for the procedure. Must be >= 1.

```bash
rooda build --max-iterations 10
rooda build -n 10
```

**`--unlimited` / `-u`**  
Set iteration mode to unlimited (runs until success or failure threshold). Overrides `--max-iterations`.

```bash
rooda build --unlimited
rooda build -u
```

**`--dry-run` / `-d`**  
Display assembled prompt without executing AI CLI. Validates configuration, prompts, and AI command.

Exit codes:
- 0 if validation passes
- 1 if validation fails

```bash
rooda build --dry-run
rooda build -d
```

### AI Command

**`--ai-cmd <command>`**  
Override AI command with direct command string. Takes precedence over `--ai-cmd-alias`.

```bash
rooda build --ai-cmd "kiro-cli chat"
rooda build --ai-cmd "claude --no-cache"
```

**`--ai-cmd-alias <alias>`**  
Override AI command using a named alias. Built-in aliases: `kiro-cli`, `claude`, `copilot`, `cursor-agent`.

```bash
rooda build --ai-cmd-alias kiro-cli
rooda build --ai-cmd-alias claude
```

### Context

**`--context <value>` / `-c <value>`**  
Pass runtime context to the procedure. Can be file path or inline text. Multiple `--context` flags accumulate.

File existence heuristic: if value exists as file, read it; otherwise treat as inline content.

```bash
# Inline context
rooda draft-plan-impl-feat --context "Add user authentication"

# File context
rooda draft-plan-impl-feat --context feature-requirements.md

# Multiple contexts
rooda build --context "Focus on auth module" --context notes.md
```

### Output Control

**`--verbose` / `-v`**  
Enable verbose output. Sets `show_ai_output=true` and `log_level=debug`.

```bash
rooda build --verbose
rooda build -v
```

**`--quiet` / `-q`**  
Suppress all non-error output.

```bash
rooda build --quiet
rooda build -q
```

**`--log-level <level>`**  
Set log level. Valid values: `debug`, `info`, `warn`, `error`.

```bash
rooda build --log-level debug
```

### Configuration

**`--config <path>`**  
Specify alternate workspace config file path. Overrides `./rooda-config.yml`.

Fragment paths in CLI overrides resolve relative to this config file's directory.

```bash
rooda build --config /path/to/custom-config.yml
```

### Prompt Overrides

Override OODA phase fragments for this execution. Multiple flags accumulate into fragment array. Replaces entire phase array (not appended to config).

**`--observe <value>`**  
Override observe phase fragments.

```bash
rooda build --observe prompts/observe_custom.md
rooda build --observe "# Observe\nCustom inline content"
```

**`--orient <value>`**  
Override orient phase fragments.

```bash
rooda build --orient prompts/orient_custom.md
```

**`--decide <value>`**  
Override decide phase fragments.

```bash
rooda build --decide prompts/decide_custom.md
```

**`--act <value>`**  
Override act phase fragments.

```bash
rooda build --act prompts/act_custom.md
```

**Multiple fragments**:
```bash
rooda build \
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
rooda build --ai-cmd-alias kiro-cli

# Run with limited iterations
rooda build --ai-cmd-alias kiro-cli --max-iterations 3

# Run until success or failure threshold
rooda build --ai-cmd-alias kiro-cli --unlimited
```

### Dry Run

```bash
# Validate configuration and prompts without executing
rooda build --ai-cmd-alias kiro-cli --dry-run

# Check exit code
rooda build --dry-run && echo "Valid" || echo "Invalid"
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
rooda build --ai-cmd-alias kiro-cli \
  --context "Focus on authentication module" \
  --context notes/auth-implementation.md
```

### Prompt Overrides

```bash
# Override single phase
rooda build --ai-cmd-alias kiro-cli \
  --observe prompts/observe_custom.md

# Override multiple phases
rooda build --ai-cmd-alias kiro-cli \
  --observe prompts/observe_specs.md \
  --observe prompts/observe_impl.md \
  --orient prompts/orient_custom.md

# Inline content
rooda build --ai-cmd-alias kiro-cli \
  --observe "# Observe\nRead AGENTS.md and specs/"
```

### Verbose Output

```bash
# See all AI output and debug logs
rooda build --ai-cmd-alias kiro-cli --verbose

# See configuration provenance
rooda build --ai-cmd-alias kiro-cli --verbose 2>&1 | grep "Configuration loaded"
```

### Custom Config

```bash
# Use alternate config file
rooda build --config configs/production.yml --ai-cmd-alias kiro-cli

# Fragments resolve relative to config file directory
rooda build --config /path/to/config.yml \
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
