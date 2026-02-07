# Configuration

## Job to be Done

Define custom OODA procedures, AI command aliases, and project-specific settings through a three-tier configuration system — workspace (`./`), global (`<config_dir>/`), and environment variables — with sensible built-in defaults for zero-config startup. Tiers merge with clear precedence (CLI flags > env vars > workspace > global > built-in defaults) and provenance tracking so users know where each setting comes from.

The global config directory resolves as: `ROODA_CONFIG_HOME` env var (if set), else `$XDG_CONFIG_HOME/rooda/` (if `XDG_CONFIG_HOME` is set), else `~/.config/rooda/` (cross-platform default).

The developer wants to start using rooda on a new project with zero configuration, customize behavior as needs evolve, and share team-wide settings across repositories — all without modifying the binary or prompt files.

## Activities

1. Load built-in default configuration embedded in the Go binary
2. Resolve global config directory (`ROODA_CONFIG_HOME` > `$XDG_CONFIG_HOME/rooda/` > `~/.config/rooda/`) and parse global config file (`<config_dir>/rooda-config.yml`)
3. Discover and parse workspace config file (`./rooda-config.yml`)
4. Resolve environment variables (`ROODA_*` prefix)
5. Merge configuration tiers with precedence, tracking provenance of each resolved value
6. Apply CLI flag overrides (highest precedence)
7. Validate merged configuration (required fields, type constraints, file path existence)
8. Expose resolved configuration to iteration loop, prompt composition, and AI CLI integration

## Acceptance Criteria

- [ ] Zero-config startup works for procedures and loop settings — `rooda build` uses built-in defaults (embedded procedures, default iteration limits) but requires the user to configure an AI command
- [ ] If no AI command is configured via any source, rooda exits with a clear error listing all ways to set one (`--ai-cmd`, `--ai-cmd-alias`, `ROODA_LOOP_AI_CMD`, `ROODA_LOOP_AI_CMD_ALIAS`, `loop.ai_cmd`, `loop.ai_cmd_alias`, procedure-level `ai_cmd`/`ai_cmd_alias`)
- [ ] Global config at `<config_dir>/rooda-config.yml` is loaded if present, ignored if absent
- [ ] Global config directory resolved as: `ROODA_CONFIG_HOME` env var > `$XDG_CONFIG_HOME/rooda/` > `~/.config/rooda/`
- [ ] Workspace config at `./rooda-config.yml` is loaded if present, ignored if absent
- [ ] Environment variables with `ROODA_` prefix override config file values
- [ ] CLI flags override all other sources
- [ ] Precedence order is: CLI flags > env vars > workspace config > global config > built-in defaults
- [ ] Workspace config overrides global config for all overlapping fields
- [ ] Procedure definitions in workspace config merge with (not replace) built-in defaults — workspace procedures add to or override individual built-in procedures, but don't remove other built-in procedures
- [ ] AI command aliases in workspace config merge with built-in aliases — workspace aliases add to or override individual built-in aliases
- [ ] Loop settings (`iteration_mode`, `default_max_iterations`, `failure_threshold`, `ai_cmd`, `ai_cmd_alias`) follow the same precedence chain
- [ ] `iteration_mode` field accepted at loop and procedure levels, resolved through tier precedence
- [ ] When `iteration_mode` is `unlimited`, `default_max_iterations` is ignored
- [ ] Procedure `ai_cmd` or `ai_cmd_alias` overrides `loop.ai_cmd` / `loop.ai_cmd_alias` for that procedure
- [ ] Within any level, `ai_cmd` (direct command) takes precedence over `ai_cmd_alias`
- [ ] Provenance tracked for each resolved setting — can report which tier provided each value
- [ ] `--verbose` displays provenance for resolved configuration
- [ ] Config file validated at load time — invalid YAML produces clear error with file path and line number
- [ ] Unknown top-level keys in config files produce warnings (not errors) for forward compatibility
- [ ] Missing config files are silently skipped (not errors)
- [ ] Procedure `iteration_mode` overrides `loop.iteration_mode` for that procedure
- [ ] Procedure `default_max_iterations` overrides `loop.default_max_iterations` for that procedure
- [ ] `ROODA_LOOP_AI_CMD` environment variable sets `loop.ai_cmd` (overrides config file, but procedure-level overrides it)
- [ ] `ROODA_LOOP_AI_CMD_ALIAS` environment variable sets `loop.ai_cmd_alias` (overrides config file, but procedure-level overrides it)
- [ ] `ROODA_LOOP_ITERATION_MODE` environment variable sets `loop.iteration_mode` (`max-iterations` or `unlimited`)
- [ ] `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` environment variable sets the default max iterations (must be >= 1)
- [ ] Prompt file paths in procedures resolved relative to the config file that defines them (workspace or global) or as embedded resources (built-in defaults)
- [ ] Built-in default procedures include all 16 v2 procedures with embedded prompt files
- [ ] Custom procedures can reference prompt files on the filesystem

## Data Structures

### Config

The fully resolved configuration after merging all tiers.

```go
type Config struct {
    Loop          LoopConfig              // Global loop settings
    Procedures    map[string]Procedure    // Named procedure definitions
    AICmdAliases  map[string]string       // AI command alias name -> command string
    Provenance    map[string]ConfigSource // Setting path -> source that provided it
}
```

**Fields:**
- `Loop` — Global defaults for iteration behavior
- `Procedures` — Map of procedure names to their definitions; includes both built-in and user-defined
- `AICmdAliases` — Map of alias names to AI command strings; includes built-in aliases (`kiro-cli`, `claude`, `copilot`, `cursor-agent`) and user-defined
- `Provenance` — Records which tier provided each resolved value, keyed by dot-path (e.g., `"loop.default_max_iterations"`, `"procedures.build.default_max_iterations"`)

### LoopConfig

```go
type LoopConfig struct {
    IterationMode        IterationMode // Iteration mode (built-in default: ModeMaxIterations)
    DefaultMaxIterations int           // Global default (built-in default: 5). Must be >= 1 when set. 0 = not set.
    FailureThreshold     int           // Consecutive failures before abort (built-in default: 3)
    AICmd                string        // Default AI command (direct command string, optional)
    AICmdAlias           string        // Default AI command alias name (resolved from AICmdAliases, optional)
}
```

**Fields:**
- `IterationMode` — Controls whether iterations are limited or unlimited. Built-in default: `ModeMaxIterations`. Empty string means not set (for merging — after built-in defaults are applied, always has a value).
- `DefaultMaxIterations` — Global default iteration limit. Built-in default: 5. Must be >= 1 when set. 0 means not set (for merging). Ignored when `IterationMode` is `ModeUnlimited`.
- `FailureThreshold` — Consecutive failures before the loop aborts. Built-in default: 3. Must be >= 1.
- `AICmd` — Direct AI command string. If set, takes precedence over `AICmdAlias`. Empty means not set.
- `AICmdAlias` — Name of an alias from `AICmdAliases` to use as the default AI command. Empty means not set. If both `AICmd` and `AICmdAlias` are set, `AICmd` wins.

### Procedure

```go
type Procedure struct {
    Display              string        // Human-readable name (optional)
    Summary              string        // One-line description (optional)
    Description          string        // Detailed description (optional)
    Observe              string        // Path to observe phase prompt file, or embedded resource name
    Orient               string        // Path to orient phase prompt file, or embedded resource name
    Decide               string        // Path to decide phase prompt file, or embedded resource name
    Act                  string        // Path to act phase prompt file, or embedded resource name
    IterationMode        IterationMode // Override loop iteration mode ("" = inherit from loop)
    DefaultMaxIterations int           // Override loop.default_max_iterations (0 = inherit from loop). Must be >= 1 when set.
    AICmd                string        // Override AI command for this procedure (optional)
    AICmdAlias           string        // Override AI command alias for this procedure (optional)
}
```

**Fields:**
- `Observe`, `Orient`, `Decide`, `Act` — Paths to prompt markdown files. For built-in procedures, these reference embedded resources. For user-defined procedures, these are filesystem paths resolved relative to the config file location.
- `IterationMode` — Optional per-procedure override. Empty string means inherit from loop settings. When set to `ModeUnlimited`, `DefaultMaxIterations` is ignored.
- `DefaultMaxIterations` — Optional per-procedure override. 0 means inherit from `loop.default_max_iterations`. Must be >= 1 when set.
- `AICmd` — Optional per-procedure AI command override. If set, takes precedence over `AICmdAlias`. Empty means inherit from loop settings.
- `AICmdAlias` — Optional per-procedure AI command alias. Empty means inherit from loop settings. If both `AICmd` and `AICmdAlias` are set, `AICmd` wins.

### ConfigSource

```go
type ConfigSource struct {
    Tier  ConfigTier // Which tier provided this value
    File  string     // File path (for workspace/global tiers) or "" for built-in/env/cli
    Value any        // The resolved value
}

type ConfigTier string

const (
    TierBuiltIn   ConfigTier = "built-in"
    TierGlobal    ConfigTier = "global"    // <config_dir>/rooda-config.yml
    TierWorkspace ConfigTier = "workspace" // ./rooda-config.yml
    TierEnvVar    ConfigTier = "env"       // ROODA_* environment variables
    TierCLIFlag   ConfigTier = "cli"       // --flag values
)

type IterationMode string

const (
    ModeMaxIterations IterationMode = "max-iterations" // Run up to DefaultMaxIterations
    ModeUnlimited     IterationMode = "unlimited"      // Run until SUCCESS signal, failure threshold, or Ctrl+C
)
```

### YAML Config File Schema

Both workspace and global config files share the same schema:

```yaml
# Loop settings (all optional — built-in defaults apply)
loop:
  iteration_mode: max-iterations  # "max-iterations" or "unlimited" (built-in default: max-iterations)
  default_max_iterations: 5       # Default max iterations (must be >= 1). Ignored when mode is unlimited.
  failure_threshold: 3            # Consecutive failures before abort
  ai_cmd_alias: claude            # Default AI command alias for all procedures
  # ai_cmd: "custom-cli --flags"  # Or set a direct command (overrides ai_cmd_alias)

# AI command aliases (optional — merges with built-in aliases)
ai_cmd_aliases:
  fast: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-5-haiku-20241022"
  thorough: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-7-sonnet-20250219"

# Procedure definitions (optional — merges with built-in procedures)
procedures:
  build:
    display: "Build from Plan"
    summary: "Implements tasks from plan"
    description: "Reads from work tracking, picks task, implements, tests, commits"
    observe: prompts/observe_plan_specs_impl.md
    orient: prompts/orient_build.md
    decide: prompts/decide_build.md
    act: prompts/act_build.md
    default_max_iterations: 10
    ai_cmd_alias: thorough     # This procedure uses a beefier model

  my-custom-procedure:
    observe: my-prompts/observe.md
    orient: my-prompts/orient.md
    decide: my-prompts/decide.md
    act: my-prompts/act.md
    default_max_iterations: 3
    # ai_cmd_alias inherits from loop.ai_cmd_alias if not set

  my-long-running:
    observe: my-prompts/observe.md
    orient: my-prompts/orient.md
    decide: my-prompts/decide.md
    act: my-prompts/act.md
    iteration_mode: unlimited    # This procedure runs until SUCCESS or failure threshold
    # default_max_iterations ignored when mode is unlimited
```

### Built-in Defaults

These are compiled into the Go binary and always available:

```go
var BuiltInDefaults = Config{
    Loop: LoopConfig{
        IterationMode:        ModeMaxIterations,
        DefaultMaxIterations: 5,
        FailureThreshold:     3,
    },
    AICmdAliases: map[string]string{
        "kiro-cli":     "kiro-cli chat --no-interactive --trust-all-tools",
        "claude":       "claude-cli --no-interactive",
        "copilot":      "github-copilot-cli",
        "cursor-agent": "cursor-agent -p -f --stream-partial-output --output-format stream-json",
    },
    Procedures: builtInProcedures, // All 16 procedures with embedded prompts
}
```

### Environment Variable Mapping

| Environment Variable | Config Path | Type |
|---|---|---|
| `ROODA_CONFIG_HOME` | Global config directory (overrides XDG resolution) | string (path) |
| `ROODA_LOOP_AI_CMD` | AI command (direct, overrides alias selection) | string |
| `ROODA_LOOP_AI_CMD_ALIAS` | AI command alias name (resolved from merged config) | string |
| `ROODA_LOOP_ITERATION_MODE` | `loop.iteration_mode` (`max-iterations` or `unlimited`) | string |
| `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` | `loop.default_max_iterations` (must be >= 1) | int |
| `ROODA_LOOP_FAILURE_THRESHOLD` | `loop.failure_threshold` | int |

Environment variables use the `ROODA_` prefix. They override config file values but are overridden by CLI flags.

## Algorithm

### Config Loading

```
function LoadConfig(cliFlags CLIFlags) -> (Config, error):
    // 1. Start with built-in defaults
    config = BuiltInDefaults.DeepCopy()
    provenance = initProvenance(config, TierBuiltIn)

    // 2. Resolve global config directory and load config
    globalDir = resolveGlobalConfigDir()  // ROODA_CONFIG_HOME > XDG_CONFIG_HOME/rooda > ~/.config/rooda
    globalPath = filepath.Join(globalDir, "rooda-config.yml")
    if fileExists(globalPath):
        globalConfig, err = parseYAML(globalPath)
        if err:
            return error("global config %s: %v", globalPath, err)
        config = mergeConfig(config, globalConfig, provenance, TierGlobal, globalPath)

    // 3. Load workspace config (./rooda-config.yml)
    workspacePath = "./rooda-config.yml"
    if fileExists(workspacePath):
        workspaceConfig, err = parseYAML(workspacePath)
        if err:
            return error("workspace config %s: %v", workspacePath, err)
        config = mergeConfig(config, workspaceConfig, provenance, TierWorkspace, workspacePath)

    // 4. Apply environment variables (override config file values at loop level)
    if env("ROODA_LOOP_AI_CMD") != "":
        config.Loop.AICmd = env("ROODA_LOOP_AI_CMD")
        provenance["loop.ai_cmd"] = ConfigSource{TierEnvVar, "", config.Loop.AICmd}
    if env("ROODA_LOOP_AI_CMD_ALIAS") != "":
        config.Loop.AICmdAlias = env("ROODA_LOOP_AI_CMD_ALIAS")
        provenance["loop.ai_cmd_alias"] = ConfigSource{TierEnvVar, "", config.Loop.AICmdAlias}
    if env("ROODA_LOOP_ITERATION_MODE") != "":
        config.Loop.IterationMode = IterationMode(env("ROODA_LOOP_ITERATION_MODE"))
        provenance["loop.iteration_mode"] = ConfigSource{TierEnvVar, "", config.Loop.IterationMode}
    if env("ROODA_LOOP_DEFAULT_MAX_ITERATIONS") != "":
        config.Loop.DefaultMaxIterations = parseInt(env("ROODA_LOOP_DEFAULT_MAX_ITERATIONS"))
        provenance["loop.default_max_iterations"] = ConfigSource{TierEnvVar, "", config.Loop.DefaultMaxIterations}
    if env("ROODA_LOOP_FAILURE_THRESHOLD") != "":
        config.Loop.FailureThreshold = parseInt(env("ROODA_LOOP_FAILURE_THRESHOLD"))
        provenance["loop.failure_threshold"] = ConfigSource{TierEnvVar, "", config.Loop.FailureThreshold}

    // 5. Apply CLI flags (highest precedence)
    if cliFlags.MaxIterations != nil:
        config.Loop.DefaultMaxIterations = *cliFlags.MaxIterations
        provenance["loop.default_max_iterations"] = ConfigSource{TierCLIFlag, "", *cliFlags.MaxIterations}
    // Additional CLI flag overrides applied by caller (--ai-cmd, --ai-cmd-alias, etc.)

    // 6. Validate
    err = validateConfig(config)
    if err:
        return error("config validation: %v", err)

    config.Provenance = provenance
    return config, nil
```

### Global Config Directory Resolution

```
function resolveGlobalConfigDir() -> string:
    // 1. ROODA_CONFIG_HOME env var (highest precedence, explicit override)
    if env("ROODA_CONFIG_HOME") != "":
        return env("ROODA_CONFIG_HOME")

    // 2. XDG_CONFIG_HOME (cross-platform standard)
    if env("XDG_CONFIG_HOME") != "":
        return filepath.Join(env("XDG_CONFIG_HOME"), "rooda")

    // 3. Default: ~/.config/rooda
    return filepath.Join(homeDir(), ".config", "rooda")
```

### Config Merging

```
function mergeConfig(base Config, overlay ConfigFile, provenance map, tier ConfigTier, filePath string) -> Config:
    // Merge loop settings (overlay wins if non-zero/non-empty)
    if overlay.Loop.IterationMode != "":
        base.Loop.IterationMode = overlay.Loop.IterationMode
        provenance["loop.iteration_mode"] = ConfigSource{tier, filePath, overlay.Loop.IterationMode}
    if overlay.Loop.DefaultMaxIterations != 0:
        base.Loop.DefaultMaxIterations = overlay.Loop.DefaultMaxIterations
        provenance["loop.default_max_iterations"] = ConfigSource{tier, filePath, overlay.Loop.DefaultMaxIterations}
    if overlay.Loop.FailureThreshold != 0:
        base.Loop.FailureThreshold = overlay.Loop.FailureThreshold
        provenance["loop.failure_threshold"] = ConfigSource{tier, filePath, overlay.Loop.FailureThreshold}
    if overlay.Loop.AICmd != "":
        base.Loop.AICmd = overlay.Loop.AICmd
        provenance["loop.ai_cmd"] = ConfigSource{tier, filePath, overlay.Loop.AICmd}
    if overlay.Loop.AICmdAlias != "":
        base.Loop.AICmdAlias = overlay.Loop.AICmdAlias
        provenance["loop.ai_cmd_alias"] = ConfigSource{tier, filePath, overlay.Loop.AICmdAlias}

    // Merge AI command aliases (overlay adds to or overrides individual aliases)
    for name, command in overlay.AICmdAliases:
        base.AICmdAliases[name] = command
        provenance["ai_cmd_aliases." + name] = ConfigSource{tier, filePath, command}

    // Merge procedures (overlay adds to or overrides individual procedures)
    for name, proc in overlay.Procedures:
        // Resolve prompt file paths relative to config file directory
        configDir = dirname(filePath)
        proc.Observe = resolvePath(configDir, proc.Observe)
        proc.Orient = resolvePath(configDir, proc.Orient)
        proc.Decide = resolvePath(configDir, proc.Decide)
        proc.Act = resolvePath(configDir, proc.Act)

        base.Procedures[name] = proc
        provenance["procedures." + name] = ConfigSource{tier, filePath, proc}

    return base
```

### Config Validation

```
function validateConfig(config Config) -> error:
    // Validate loop settings
    if config.Loop.IterationMode != ModeMaxIterations && config.Loop.IterationMode != ModeUnlimited:
        return error("loop.iteration_mode must be 'max-iterations' or 'unlimited', got '%s'", config.Loop.IterationMode)
    if config.Loop.DefaultMaxIterations < 0:
        return error("loop.default_max_iterations must be >= 1 when set, got %d", config.Loop.DefaultMaxIterations)
    if config.Loop.IterationMode == ModeMaxIterations && config.Loop.DefaultMaxIterations < 1:
        return error("loop.default_max_iterations must be >= 1 when iteration_mode is 'max-iterations', got %d", config.Loop.DefaultMaxIterations)
    if config.Loop.FailureThreshold < 1:
        return error("loop.failure_threshold must be >= 1, got %d", config.Loop.FailureThreshold)

    // Validate procedures
    for name, proc in config.Procedures:
        if proc.Observe == "":
            return error("procedure %s: observe is required", name)
        if proc.Orient == "":
            return error("procedure %s: orient is required", name)
        if proc.Decide == "":
            return error("procedure %s: decide is required", name)
        if proc.Act == "":
            return error("procedure %s: act is required", name)
        if proc.IterationMode != "" && proc.IterationMode != ModeMaxIterations && proc.IterationMode != ModeUnlimited:
            return error("procedure %s: iteration_mode must be '', 'max-iterations', or 'unlimited', got '%s'", name, proc.IterationMode)
        if proc.DefaultMaxIterations != 0 && proc.DefaultMaxIterations < 1:
            return error("procedure %s: default_max_iterations must be >= 1 when set, got %d", name, proc.DefaultMaxIterations)

    // Validate AI command aliases (values must be non-empty strings)
    for name, command in config.AICmdAliases:
        if command == "":
            return error("ai_cmd_aliases.%s: command must be non-empty", name)

    // Validate alias references exist in merged alias map
    if config.Loop.AICmdAlias != "":
        if _, exists = config.AICmdAliases[config.Loop.AICmdAlias]; !exists:
            return error("loop.ai_cmd_alias references unknown alias: %s\nAvailable: %v",
                config.Loop.AICmdAlias, keys(config.AICmdAliases))
    for name, proc in config.Procedures:
        if proc.AICmdAlias != "":
            if _, exists = config.AICmdAliases[proc.AICmdAlias]; !exists:
                return error("procedure %s: ai_cmd_alias references unknown alias: %s\nAvailable: %v",
                    name, proc.AICmdAlias, keys(config.AICmdAliases))

    return nil
```

### AI Command Resolution

The AI command is resolved from multiple sources with precedence. This happens after config loading, using the merged config, procedure name, and CLI flags:

```
function ResolveAICommand(config Config, procedureName string, cliFlags CLIFlags) -> (string, ConfigSource):
    // 1. --ai-cmd flag (direct command, highest precedence)
    if cliFlags.AICmd != "":
        return cliFlags.AICmd, ConfigSource{TierCLIFlag, "", cliFlags.AICmd}

    // 2. --ai-cmd-alias flag (alias name, resolved from merged config)
    if cliFlags.AICmdAlias != "":
        return resolveAlias(config, cliFlags.AICmdAlias, "flag --ai-cmd-alias")

    // 3. Procedure-level ai_cmd (direct command)
    proc, exists = config.Procedures[procedureName]
    if exists && proc.AICmd != "":
        return proc.AICmd, config.Provenance["procedures." + procedureName]

    // 4. Procedure-level ai_cmd_alias
    if exists && proc.AICmdAlias != "":
        return resolveAlias(config, proc.AICmdAlias, "procedure " + procedureName)

    // 5. loop.ai_cmd (direct command, already merged from config tiers + env vars)
    if config.Loop.AICmd != "":
        return config.Loop.AICmd, config.Provenance["loop.ai_cmd"]

    // 6. loop.ai_cmd_alias (already merged from config tiers + env vars)
    if config.Loop.AICmdAlias != "":
        return resolveAlias(config, config.Loop.AICmdAlias, "loop.ai_cmd_alias")

    // 7. No AI command configured — error with guidance
    return error("no AI command configured\n\nSet one via:\n" +
        "  --ai-cmd \"your-command\"           CLI flag (direct command)\n" +
        "  --ai-cmd-alias <name>             CLI flag (alias from config)\n" +
        "  ROODA_LOOP_AI_CMD=your-command    Environment variable\n" +
        "  ROODA_LOOP_AI_CMD_ALIAS=<name>    Environment variable\n" +
        "  loop.ai_cmd or loop.ai_cmd_alias  rooda-config.yml\n" +
        "\nAvailable aliases: %v", keys(config.AICmdAliases))

function resolveAlias(config Config, aliasName string, source string) -> (string, ConfigSource):
    command, exists = config.AICmdAliases[aliasName]
    if !exists:
        return error("unknown AI command alias: %s (from %s)\nAvailable: %v", aliasName, source, keys(config.AICmdAliases))
    return command, ConfigSource{source, "", command}
```

**Precedence summary (consistent with max iterations):**
1. `--ai-cmd` (CLI flag, direct command)
2. `--ai-cmd-alias` (CLI flag, alias)
3. Procedure `ai_cmd` (config, direct command)
4. Procedure `ai_cmd_alias` (config, alias)
5. `loop.ai_cmd` (merged: env var > workspace > global > built-in)
6. `loop.ai_cmd_alias` (merged: env var > workspace > global > built-in)
7. Error — no AI command configured

Note: `ROODA_LOOP_AI_CMD` and `ROODA_LOOP_AI_CMD_ALIAS` environment variables set `loop.ai_cmd` and `loop.ai_cmd_alias` respectively (same as how `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` sets `loop.default_max_iterations`). Procedure-level settings still override them.

### Max Iterations Resolution

```
function ResolveMaxIterations(config Config, procedureName string, cliFlags CLIFlags) -> (*int, ConfigSource):
    // 1. --max-iterations N CLI flag (highest precedence, must be >= 1)
    if cliFlags.MaxIterations != nil:
        return cliFlags.MaxIterations, ConfigSource{TierCLIFlag, "", *cliFlags.MaxIterations}

    // 2. --unlimited CLI flag
    if cliFlags.Unlimited:
        return nil, ConfigSource{TierCLIFlag, "", "unlimited"}

    // 3. Procedure-level iteration_mode and default_max_iterations
    proc, exists = config.Procedures[procedureName]
    if exists:
        if proc.IterationMode == ModeUnlimited:
            return nil, config.Provenance["procedures." + procedureName]
        if proc.IterationMode == ModeMaxIterations && proc.DefaultMaxIterations > 0:
            return &proc.DefaultMaxIterations, config.Provenance["procedures." + procedureName]
        // proc.IterationMode == ModeMaxIterations && proc.DefaultMaxIterations == 0:
        //   mode set but count inherits from loop — fall through
        // proc.IterationMode == "": inherit everything from loop — fall through
        if proc.IterationMode == "" && proc.DefaultMaxIterations > 0:
            // No mode override but explicit count — use count with loop's mode
            if config.Loop.IterationMode == ModeUnlimited:
                return nil, config.Provenance["loop.iteration_mode"]  // mode governs
            return &proc.DefaultMaxIterations, config.Provenance["procedures." + procedureName]

    // 4. Loop-level iteration_mode and default_max_iterations (already merged from tiers)
    if config.Loop.IterationMode == ModeUnlimited:
        return nil, config.Provenance["loop.iteration_mode"]
    return &config.Loop.DefaultMaxIterations, config.Provenance["loop.default_max_iterations"]
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No config files exist, no env vars, no AI cmd flag | Error: "no AI command configured" with guidance on all ways to set one |
| No config files exist, AI command provided via flag | Works using built-in defaults for procedures and loop settings |
| Global config exists, no workspace config | Global config merges with built-in defaults |
| Both global and workspace config exist | Workspace overrides global for overlapping values; non-overlapping values from both are kept |
| Workspace config defines procedure that overrides a built-in | Workspace procedure replaces the built-in procedure entirely (all fields) |
| Workspace config defines new procedure not in built-ins | New procedure added alongside built-in procedures |
| Global config defines procedure, workspace does not | Global procedure is available alongside built-in procedures |
| Environment variable `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` set to non-integer | Error: "ROODA_LOOP_DEFAULT_MAX_ITERATIONS must be an integer, got 'abc'" |
| Environment variable `ROODA_LOOP_DEFAULT_MAX_ITERATIONS=0` | Error: "ROODA_LOOP_DEFAULT_MAX_ITERATIONS must be >= 1, got 0" |
| `ROODA_LOOP_ITERATION_MODE=unlimited` | Sets loop iteration mode to unlimited |
| `ROODA_LOOP_ITERATION_MODE=invalid-value` | Error at config validation: "loop.iteration_mode must be 'max-iterations' or 'unlimited'" |
| Config file has invalid YAML | Error with file path and parse error details |
| Config file has unknown top-level key | Warning logged, key ignored (forward compatibility) |
| Config file has unknown key inside `procedures` | Warning logged, key ignored |
| Procedure references prompt file that doesn't exist | Error at validation time: "procedure 'X': observe file not found: path/to/file.md" |
| Built-in procedure prompt file (embedded) | Loaded from Go binary via `go:embed`, not filesystem |
| Workspace procedure prompt file path | Resolved relative to workspace config file directory |
| Global procedure prompt file path | Resolved relative to global config file directory |
| Global config directory doesn't exist | Silently skipped (no global config) |
| `ROODA_CONFIG_HOME` set to non-existent directory | Silently skipped (same as missing directory) |
| `XDG_CONFIG_HOME` set, `ROODA_CONFIG_HOME` not set | Global config at `$XDG_CONFIG_HOME/rooda/rooda-config.yml` |
| Both `ROODA_CONFIG_HOME` and `XDG_CONFIG_HOME` set | `ROODA_CONFIG_HOME` wins (more specific override) |
| Config file is empty | Valid — all built-in defaults apply |
| Config file has only `ai_cmd_aliases` section | Valid — procedures and loop use built-in defaults |
| Two config files define same AI command alias | Workspace wins over global, both override built-in |
| `ROODA_LOOP_AI_CMD_ALIAS` references non-existent alias | Error at config validation: "loop.ai_cmd_alias references unknown alias: X" |
| Both `ROODA_LOOP_AI_CMD` and `ROODA_LOOP_AI_CMD_ALIAS` set | Both set on `config.Loop`; `loop.ai_cmd` wins over `loop.ai_cmd_alias` during resolution |
| `loop.ai_cmd` and `loop.ai_cmd_alias` both set | `loop.ai_cmd` wins (direct command overrides alias at same level) |
| Procedure `ai_cmd_alias` set, `loop.ai_cmd_alias` also set | Procedure-level wins for that procedure; other procedures use loop-level |
| Procedure `ai_cmd_alias` references non-existent alias | Error at config validation: "procedure X: ai_cmd_alias references unknown alias: Y" |
| Global config sets `loop.ai_cmd_alias`, workspace does not | Global loop-level alias is used (merged through config tiers) |
| Workspace sets `loop.ai_cmd_alias`, overriding global | Workspace wins per standard tier precedence |
| `ROODA_LOOP_AI_CMD` set, procedure has `ai_cmd_alias` | Procedure-level wins — env vars set loop-level defaults, procedure overrides loop (same as max iterations) |
| `loop.failure_threshold: 0` in config | Error: "loop.failure_threshold must be >= 1" |
| `loop.default_max_iterations: 0` in config | Valid — means not set (inherits built-in default). Note: 0 no longer means unlimited; use `iteration_mode: unlimited` instead. |
| Procedure `default_max_iterations: 0` in config | Valid — means inherit from loop. Note: 0 no longer means unlimited; use `iteration_mode: unlimited` instead. |
| `iteration_mode: unlimited` at global, `default_max_iterations: 10` at workspace without mode | Unlimited — mode governs. Workspace only set count, didn't override mode, so global mode applies. |
| `iteration_mode: invalid-value` in config | Error at validation: "loop.iteration_mode must be 'max-iterations' or 'unlimited'" |
| Procedure `iteration_mode: unlimited` with `default_max_iterations: 5` | Unlimited for that procedure — mode governs, count ignored |
| Procedure `iteration_mode: max-iterations` without `default_max_iterations` | Valid — count inherits from `loop.default_max_iterations` |
| Workspace config and `--config` flag both present | `--config` flag specifies which config file to load as workspace config |

## Dependencies

- **Go standard library** — `os`, `path/filepath`, `encoding/json` for file operations and path resolution
- **Go YAML library** — `gopkg.in/yaml.v3` or `github.com/goccy/go-yaml` for YAML parsing (replaces yq dependency)
- **go:embed** — For embedding built-in default prompts and config in the binary
- **cli-interface** — Provides CLI flags that override config values
- **prompt-composition** — Consumes procedure definitions to assemble prompts
- **ai-cli-integration** — Consumes AI command aliases and resolved AI command
- **iteration-loop** — Consumes loop settings and resolved max iterations

## Implementation Mapping

**Source files:**
- `internal/config/config.go` — Core config types, loading, and merging logic
- `internal/config/defaults.go` — Built-in default configuration and embedded prompts
- `internal/config/validate.go` — Config validation
- `internal/config/provenance.go` — Provenance tracking
- `internal/config/env.go` — Environment variable resolution
- `cmd/rooda/main.go` — CLI flag parsing, calls `LoadConfig`

**Related specs:**
- `iteration-loop.md` — Consumes `LoopConfig`, max iterations resolution, and AI command resolution
- `prompt-composition.md` — Consumes `Procedure` definitions to assemble prompts
- `ai-cli-integration.md` — Consumes `AICmdAliases` and resolved AI command
- `cli-interface.md` — Defines CLI flags that override config values
- `procedures.md` — Defines the 16 built-in procedures that ship as defaults

## Examples

### Example 1: Minimal Startup

**Scenario:** Developer installs rooda and runs it with only an AI command specified.

**Input:**
```bash
rooda build --ai-cmd-alias claude
```

**Config resolution:**
```
loop.iteration_mode: max-iterations   (built-in)
loop.default_max_iterations: 5        (built-in)
loop.failure_threshold: 3             (built-in)
procedure: build                      (built-in, embedded prompts)
ai_cmd: claude-cli --no-interactive (cli: --ai-cmd-alias "claude" → built-in alias)
```

**Verification:**
- No config files needed
- Built-in `build` procedure uses embedded prompt files
- Default mode `max-iterations` with 5 iterations, 3 failure threshold
- AI command must be explicitly chosen — no built-in default command

### Example 1b: No AI Command Configured

**Scenario:** Developer runs rooda without configuring an AI command.

**Input:**
```bash
rooda build
```

**Expected Output:**
```
Error: no AI command configured

Set one via:
  --ai-cmd "your-command"           CLI flag (direct command)
  --ai-cmd-alias <name>             CLI flag (alias from config)
  ROODA_LOOP_AI_CMD=your-command    Environment variable
  ROODA_LOOP_AI_CMD_ALIAS=<name>    Environment variable
  loop.ai_cmd or loop.ai_cmd_alias  rooda-config.yml

Available aliases: claude, copilot, cursor-agent, kiro-cli
```

**Verification:**
- Clear error with all configuration options listed
- Available built-in aliases shown to guide the user
- No silent fallback to a default command

### Example 2: Workspace Config Override

**Scenario:** Team wants to use claude by default and increase build iterations.

**Workspace config (`./rooda-config.yml`):**
```yaml
loop:
  ai_cmd_alias: claude

procedures:
  build:
    default_max_iterations: 10
```

**Input:**
```bash
rooda build
```

**Config resolution:**
```
loop.default_max_iterations: 5        (built-in)
loop.failure_threshold: 3             (built-in)
loop.ai_cmd_alias: claude             (workspace: ./rooda-config.yml)
procedure: build
  observe/orient/decide/act:          (built-in, embedded — not overridden)
  default_max_iterations: 10          (workspace: ./rooda-config.yml)
ai_cmd: claude-cli --no-interactive (loop.ai_cmd_alias → built-in alias)
```

**Verification:**
- Build procedure prompt files still use built-in embedded defaults (workspace only overrode `default_max_iterations`)
- Iteration limit increased to 10
- Claude CLI used via `loop.ai_cmd_alias` — no flag needed
- Provenance shows `default_max_iterations` came from workspace config

### Example 3: Three-Tier Merge with Procedure-Level AI Command

**Scenario:** Developer has global preferences, project uses a default AI command, and one procedure uses a different model.

**Global config (`$XDG_CONFIG_HOME/rooda/rooda-config.yml` or `~/.config/rooda/rooda-config.yml`):**
```yaml
loop:
  default_max_iterations: 8
  ai_cmd_alias: claude

ai_cmd_aliases:
  fast: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-5-haiku-20241022"
```

**Workspace config (`./rooda-config.yml`):**
```yaml
loop:
  failure_threshold: 5

procedures:
  my-lint:
    observe: prompts/observe_lint.md
    orient: prompts/orient_lint.md
    decide: prompts/decide_lint.md
    act: prompts/act_lint.md
    default_max_iterations: 1
    ai_cmd_alias: fast               # Lint is cheap — use a faster model
```

**Input:**
```bash
rooda my-lint
```

**Config resolution:**
```
loop.default_max_iterations: 8        (global: ~/.config/rooda/rooda-config.yml)
loop.failure_threshold: 5             (workspace: ./rooda-config.yml)
loop.ai_cmd_alias: claude             (global: ~/.config/rooda/rooda-config.yml)
ai_cmd_aliases.fast: kiro-cli ... haiku (global: ~/.config/rooda/rooda-config.yml)
ai_cmd_aliases.kiro-cli: kiro-cli ...  (built-in)
ai_cmd_aliases.claude: claude-cli ...  (built-in)
ai_cmd_aliases.copilot: github-copilot-cli (built-in)
ai_cmd_aliases.cursor-agent: cursor-agent ... (built-in)
procedure: my-lint                    (workspace: ./rooda-config.yml)
  default_max_iterations: 1           (workspace: ./rooda-config.yml)
  ai_cmd_alias: fast                  (workspace: ./rooda-config.yml)
ai_cmd: kiro-cli ... haiku   (procedure my-lint.ai_cmd_alias "fast" → global alias)
```

**Verification:**
- `default_max_iterations` from global (8), but overridden by procedure-level (1) for `my-lint`
- `failure_threshold` from workspace (5) overrides global which used built-in default
- `loop.ai_cmd_alias: claude` from global provides the default, but `my-lint` overrides with `fast`
- Running `rooda build` would use `claude` (from `loop.ai_cmd_alias`), since `build` has no procedure-level override
- All 16 built-in procedures still available alongside custom `my-lint`
- Custom procedure prompt paths resolved relative to workspace config directory

### Example 4: Environment Variable Override

**Input:**
```bash
export ROODA_LOOP_AI_CMD="aider --yes"
export ROODA_LOOP_DEFAULT_MAX_ITERATIONS=3
rooda build
```

**Config resolution:**
```
loop.default_max_iterations: 3        (env: ROODA_LOOP_DEFAULT_MAX_ITERATIONS)
loop.failure_threshold: 3             (built-in)
loop.ai_cmd: aider --yes             (env: ROODA_LOOP_AI_CMD)
procedure: build                      (built-in, embedded prompts)
ai_cmd: aider --yes          (loop.ai_cmd)
```

**Verification:**
- Environment variables override config file values at the loop level
- `ROODA_LOOP_AI_CMD` sets `loop.ai_cmd` directly (not an alias name)
- `ROODA_LOOP_DEFAULT_MAX_ITERATIONS` overrides `loop.default_max_iterations`
- Procedure-level `ai_cmd`/`ai_cmd_alias` or `default_max_iterations` would still override these env vars
- Built-in procedure still used

### Example 5: CLI Flag Takes Precedence Over Everything

**Input:**
```bash
export ROODA_LOOP_DEFAULT_MAX_ITERATIONS=3
rooda build --max-iterations 1 --ai-cmd "claude-cli --no-interactive"
```

**Config resolution:**
```
max_iterations: 1                     (cli: --max-iterations)
ai_cmd: claude-cli --no-interactive (cli: --ai-cmd)
```

**Verification:**
- CLI `--max-iterations 1` overrides everything (env var, procedure, loop config)
- CLI `--ai-cmd` overrides everything (env var, procedure, loop config)
- CLI flags have highest precedence

### Example 6: Provenance Display (Verbose)

**Input:**
```bash
rooda build --ai-cmd-alias claude --verbose
```

**Expected Output (provenance section at startup):**
```
[VERBOSE] Configuration provenance:
  loop.iteration_mode: max-iterations (built-in)
  loop.default_max_iterations: 5 (built-in)
  loop.failure_threshold: 3 (built-in)
  procedure: build (built-in)
  ai_cmd: claude-cli --no-interactive (cli: --ai-cmd-alias "claude" → built-in alias)
```

**Verification:**
- Every resolved setting shows its source tier
- Only shown with `--verbose` flag
- Helps diagnose "where did this value come from?" questions

### Example 7: Invalid Config File

**Workspace config (`./rooda-config.yml`):**
```yaml
loop:
  default_max_iterations: abc  # Not an integer
```

**Input:**
```bash
rooda build
```

**Expected Output:**
```
Error: workspace config ./rooda-config.yml: line 2: cannot unmarshal "abc" into int for field default_max_iterations
```

**Verification:**
- Clear error with file path and line number
- Error identifies the field and type mismatch

### Example 8: Custom Procedure with Filesystem Prompts

**Workspace config (`./rooda-config.yml`):**
```yaml
procedures:
  security-audit:
    display: "Security Audit"
    summary: "Audit codebase for security vulnerabilities"
    observe: .rooda/prompts/observe_security.md
    orient: .rooda/prompts/orient_security.md
    decide: .rooda/prompts/decide_security.md
    act: .rooda/prompts/act_security.md
    default_max_iterations: 1
```

**Input:**
```bash
rooda security-audit
```

**Verification:**
- Custom procedure loads from workspace config
- Prompt files resolved relative to `./` (workspace config directory)
- Runs alongside all 16 built-in procedures
- If prompt files don't exist, clear error at startup

### Example 9: `--config` Flag Override

**Input:**
```bash
rooda build --config /path/to/team-config.yml
```

**Config resolution:**
```
# Built-in defaults loaded first
# Global config at <config_dir>/rooda-config.yml loaded if present
# /path/to/team-config.yml loaded as workspace config (instead of ./rooda-config.yml)
```

**Verification:**
- `--config` flag replaces the workspace config file path
- Global config still loads from resolved global config directory
- Prompt paths in team config resolved relative to `/path/to/`

### Example 10: Unlimited Iterations via Config

**Scenario:** Team wants the build procedure to run without iteration limits, relying on `<promise>SUCCESS</promise>` signals or failure threshold to terminate.

**Workspace config (`./rooda-config.yml`):**
```yaml
loop:
  ai_cmd_alias: claude

procedures:
  build:
    iteration_mode: unlimited
```

**Input:**
```bash
rooda build
```

**Config resolution:**
```
loop.iteration_mode: max-iterations   (built-in)
loop.default_max_iterations: 5        (built-in)
loop.failure_threshold: 3             (built-in)
loop.ai_cmd_alias: claude             (workspace: ./rooda-config.yml)
procedure: build
  iteration_mode: unlimited           (workspace: ./rooda-config.yml)
  observe/orient/decide/act:          (built-in, embedded — not overridden)
max_iterations: nil                   (procedure build.iteration_mode → unlimited)
ai_cmd: claude-cli --no-interactive (loop.ai_cmd_alias → built-in alias)
```

**Verification:**
- Build procedure runs unlimited despite loop-level `max-iterations` mode — procedure `iteration_mode` overrides loop
- `loop.default_max_iterations: 5` is ignored for build (mode is unlimited)
- Other procedures without `iteration_mode` override still use loop defaults (max-iterations, 5)
- Loop terminates on `<promise>SUCCESS</promise>`, failure threshold (3), or Ctrl+C

## Notes

**Design Rationale — Three-Tier System:**

The three-tier system (built-in > global > workspace) mirrors established tooling conventions (Git, npm, EditorConfig). Each tier serves a distinct use case:
- **Built-in defaults** — Zero-config startup, always available, ships with the binary
- **Global config** — Personal preferences that follow the developer across projects (preferred AI command, iteration limits). Located via `ROODA_CONFIG_HOME` > `$XDG_CONFIG_HOME/rooda/` > `~/.config/rooda/`, following the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/)
- **Workspace config** — Project-specific settings committed to the repo (custom procedures, team aliases)

This separation means a developer can set their preferred AI command globally, while each project can define custom procedures — without conflict.

**Why Merge Instead of Replace:**

Config merging (additive) instead of replacement (destructive) is critical for a good experience. If a workspace config defining one custom procedure replaced all built-in procedures, the developer would lose access to `build`, `bootstrap`, etc. Merging means workspace configs only need to specify what's different.

**Provenance Tracking:**

Provenance answers the question "where did this value come from?" — essential for debugging configuration issues. When a developer runs `rooda build --dry-run` and sees unexpected settings, provenance immediately shows which file or environment variable set each value.

**Environment Variable Convention:**

The `ROODA_` prefix follows Go CLI conventions and avoids namespace collisions. Only a small set of environment variables is supported — complex configuration belongs in YAML files. `ROODA_CONFIG_HOME` provides a rooda-specific override for the global config directory, useful for CI/CD environments or when XDG isn't appropriate. When not set, the standard `XDG_CONFIG_HOME` is respected, falling back to `~/.config/` — the conventional default for CLI tools across Linux and macOS.

**Path Resolution:**

Prompt file paths in config files are resolved relative to the config file's directory, not the current working directory. This ensures a workspace config committed to a repo works regardless of where `rooda` is invoked from within the project. Built-in procedure prompt files are embedded in the binary via `go:embed` and don't need filesystem resolution.

**Forward Compatibility:**

Unknown keys produce warnings, not errors. This allows newer config files to work with older rooda versions — new keys are simply ignored. This is important for teams where not everyone upgrades simultaneously.

**Migration from v1:**

v1 used a single `rooda-config.yml` file with yq for YAML parsing. The v2 workspace config file uses the same filename and a compatible schema, so existing v1 config files work as v2 workspace configs with minimal changes:
- Rename `default_iterations` to `default_max_iterations`
- Replace `default_iterations: 0` (unlimited) with `iteration_mode: unlimited` — v2 no longer uses 0 to mean unlimited; `default_max_iterations` must be >= 1 when set
