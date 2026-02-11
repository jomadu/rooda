# Configuration

rooda uses a three-tier configuration system with clear precedence: CLI flags > environment variables > workspace config > global config > built-in defaults.

## Configuration Tiers

### 1. Built-in Defaults (Lowest Precedence)

Embedded in the binary. Includes:
- 21 procedures with embedded prompt fragments
- Built-in AI command aliases (`kiro-cli`, `claude`, `copilot`, `cursor-agent`)
- Default loop settings (5 max iterations, 3 failure threshold, info log level)

No configuration files required to start using rooda.

### 2. Global Config

**Location**: `~/.config/rooda/rooda-config.yml` (Unix/macOS) or `%APPDATA%\rooda\rooda-config.yml` (Windows)

**Override**: Set `ROODA_CONFIG_HOME` environment variable to use a different directory.

**Purpose**: Team-wide settings shared across all repositories.

**Example**:
```yaml
# ~/.config/rooda/rooda-config.yml
loop:
  ai_cmd_alias: kiro-cli
  default_max_iterations: 5
  log_level: info

ai_cmd_aliases:
  my-ai: "my-custom-ai-tool --flag"
```

### 3. Workspace Config

**Location**: `./rooda-config.yml` (project root)

**Override**: Use `--config <path>` flag to specify alternate location.

**Purpose**: Project-specific settings.

**Example**:
```yaml
# ./rooda-config.yml
loop:
  default_max_iterations: 3
  log_level: debug

procedures:
  build:
    default_max_iterations: 10
    ai_cmd_alias: claude

  custom-procedure:
    display: "My Custom Procedure"
    summary: "Does something custom"
    observe:
      - path: "prompts/observe_custom.md"
    orient:
      - path: "prompts/orient_custom.md"
    decide:
      - path: "prompts/decide_custom.md"
    act:
      - path: "prompts/act_custom.md"
```

### 4. Environment Variables

**Prefix**: `ROODA_`

**Common variables**:
- `ROODA_LOOP_AI_CMD` - Direct AI command string
- `ROODA_LOOP_AI_CMD_ALIAS` - AI command alias name
- `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` - Default max iterations (must be >= 1)
- `ROODA_LOOP_ITERATION_MODE` - `max-iterations` or `unlimited`
- `ROODA_LOOP_LOG_LEVEL` - `debug`, `info`, `warn`, `error`
- `ROODA_LOOP_LOG_TIMESTAMP_FORMAT` - `time`, `relative`, `iso`, `none`
- `ROODA_CONFIG_HOME` - Override global config directory

**Example**:
```bash
export ROODA_LOOP_AI_CMD_ALIAS=kiro-cli
export ROODA_LOOP_LOG_LEVEL=debug
rooda build
```

### 5. CLI Flags (Highest Precedence)

Override everything. See [CLI Reference](cli-reference.md) for all flags.

**Example**:
```bash
rooda run build --ai-cmd-alias kiro-cli --max-iterations 10 --verbose
```

## Configuration Schema

### Loop Settings

```yaml
loop:
  iteration_mode: max-iterations  # or "unlimited"
  default_max_iterations: 5       # Must be >= 1
  iteration_timeout: 3600          # Seconds, nil = no timeout
  max_output_buffer: 10485760      # Bytes (10MB default)
  failure_threshold: 3             # Consecutive failures before abort
  log_level: info                  # debug, info, warn, error
  log_timestamp_format: time       # time, relative, iso, none
  show_ai_output: false            # Stream AI output to terminal
  ai_cmd: ""                       # Direct command string (optional)
  ai_cmd_alias: ""                 # Alias name (optional)
```

### AI Command Aliases

```yaml
ai_cmd_aliases:
  kiro-cli: "kiro-cli chat"
  claude: "claude --no-cache"
  my-custom: "my-ai-tool --flag value"
```

Built-in aliases: `kiro-cli`, `claude`, `copilot`, `cursor-agent`.

### Procedures

```yaml
procedures:
  my-procedure:
    display: "Human-readable name"
    summary: "One-line description"
    description: "Detailed description"
    
    # OODA phase fragments (arrays)
    observe:
      - path: "prompts/observe_something.md"
      - content: "Inline prompt content"
        parameters:
          key: value
    
    orient:
      - path: "prompts/orient_something.md"
    
    decide:
      - path: "prompts/decide_something.md"
    
    act:
      - path: "prompts/act_something.md"
    
    # Override loop settings for this procedure
    iteration_mode: max-iterations
    default_max_iterations: 10
    iteration_timeout: 1800
    max_output_buffer: 5242880
    ai_cmd_alias: claude
```

**Fragment Actions**:
- `path` - Path to fragment file (relative to config file directory, or embedded resource)
- `content` - Inline prompt content (alternative to `path`)
- `parameters` - Template parameters for substitution

**Procedure-level overrides**:
- `iteration_mode` - Override loop iteration mode
- `default_max_iterations` - Override loop default
- `iteration_timeout` - Override loop timeout
- `max_output_buffer` - Override loop buffer size
- `ai_cmd` - Direct command string (overrides loop.ai_cmd)
- `ai_cmd_alias` - Alias name (overrides loop.ai_cmd_alias)

## Precedence Rules

### AI Command Resolution

1. CLI `--ai-cmd` (direct command)
2. CLI `--ai-cmd-alias` (alias name)
3. Procedure `ai_cmd` (direct command)
4. Procedure `ai_cmd_alias` (alias name)
5. Loop `ai_cmd` (direct command)
6. Loop `ai_cmd_alias` (alias name)
7. Error if none configured

Within any level, `ai_cmd` (direct) takes precedence over `ai_cmd_alias`.

### Iteration Settings

1. CLI `--max-iterations` or `--unlimited`
2. Environment `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` or `ROODA_LOOP_ITERATION_MODE`
3. Procedure `default_max_iterations` or `iteration_mode`
4. Loop `default_max_iterations` or `iteration_mode`
5. Built-in default (5 iterations, max-iterations mode)

### Procedure Merging

- Workspace procedures **merge with** (not replace) built-in procedures
- Workspace procedure fields override only the specified fields
- Unspecified fields inherit from built-in procedure
- This allows customizing individual procedures without redefining everything

**Example**:
```yaml
# Workspace config - only override max_iterations for build procedure
procedures:
  build:
    default_max_iterations: 10
    # All other fields (observe, orient, decide, act) inherit from built-in
```

## Provenance Tracking

Use `--verbose` to see where each setting came from:

```bash
rooda run build --ai-cmd-alias kiro-cli --verbose
```

Output includes provenance for each resolved setting:
```
[INFO] Configuration loaded:
  loop.ai_cmd_alias: kiro-cli (CLI flag)
  loop.default_max_iterations: 5 (built-in default)
  procedures.build.default_max_iterations: 10 (workspace config)
```

## Validation

Configuration is validated at load time:
- Invalid YAML produces error with file path and line number
- Unknown top-level keys produce warnings (not errors)
- Missing config files are silently skipped
- Invalid values (negative numbers, unknown enums) produce errors

## Examples

### Minimal Setup (Zero Config)

```bash
# No config files needed - uses built-in defaults
# But AI command must be specified
rooda run build --ai-cmd-alias kiro-cli
```

### Global Config for Team

```yaml
# ~/.config/rooda/rooda-config.yml
loop:
  ai_cmd_alias: kiro-cli
  default_max_iterations: 5
  log_level: info
  log_timestamp_format: relative
```

All team members use the same AI command and log settings.

### Project-Specific Overrides

```yaml
# ./rooda-config.yml
loop:
  default_max_iterations: 3  # Override global default

procedures:
  build:
    default_max_iterations: 10  # Build needs more iterations
    ai_cmd_alias: claude        # Use different AI for build
```

### Custom Procedure

```yaml
# ./rooda-config.yml
procedures:
  my-audit:
    display: "Custom Audit"
    summary: "Audits custom aspects of the codebase"
    observe:
      - path: "prompts/observe_specs.md"
      - path: "prompts/observe_impl.md"
    orient:
      - content: |
          # Orient: Custom Audit
          
          Analyze the codebase for custom quality criteria:
          - Check for TODO comments
          - Verify all functions have docstrings
          - Identify unused imports
    decide:
      - path: "prompts/decide_gap_plan.md"
    act:
      - path: "prompts/act_plan.md"
    default_max_iterations: 1
```

Run with:
```bash
rooda my-audit --ai-cmd-alias kiro-cli
```

## See Also

- [CLI Reference](cli-reference.md) - All CLI flags
- [Procedures](procedures.md) - Built-in procedures
- [Troubleshooting](troubleshooting.md) - Common configuration errors
