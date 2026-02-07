# Configuration

## Job to be Done

Define custom OODA procedures, AI CLI presets, and project-specific settings through a three-tier configuration system — workspace (`./`), global (`~/.config/rooda/`), and environment variables — with sensible built-in defaults for zero-config startup. Tiers merge with clear precedence (CLI flags > env vars > workspace > global > built-in defaults) and provenance tracking so users know where each setting comes from.

The developer wants to start using rooda on a new project with zero configuration, customize behavior as needs evolve, and share team-wide settings across repositories — all without modifying the binary or prompt files.

## Activities

1. Load built-in default configuration embedded in the Go binary
2. Discover and parse global config file (`~/.config/rooda/config.yml`)
3. Discover and parse workspace config file (`./rooda-config.yml`)
4. Resolve environment variables (`ROODA_*` prefix)
5. Merge configuration tiers with precedence, tracking provenance of each resolved value
6. Apply CLI flag overrides (highest precedence)
7. Validate merged configuration (required fields, type constraints, file path existence)
8. Expose resolved configuration to iteration loop, prompt composition, and AI CLI integration

## Acceptance Criteria

- [ ] Zero-config startup works — `rooda build` executes using only built-in defaults (embedded procedures, default AI CLI, default iteration limits)
- [ ] Global config at `~/.config/rooda/config.yml` is loaded if present, ignored if absent
- [ ] Workspace config at `./rooda-config.yml` is loaded if present, ignored if absent
- [ ] Environment variables with `ROODA_` prefix override config file values
- [ ] CLI flags override all other sources
- [ ] Precedence order is: CLI flags > env vars > workspace config > global config > built-in defaults
- [ ] Workspace config overrides global config for all overlapping fields
- [ ] Procedure definitions in workspace config merge with (not replace) built-in defaults — workspace procedures add to or override individual built-in procedures, but don't remove other built-in procedures
- [ ] AI tool presets in workspace config merge with built-in presets — workspace presets add to or override individual built-in presets
- [ ] Loop settings (`default_max_iterations`, `failure_threshold`) follow the same precedence chain
- [ ] Provenance tracked for each resolved setting — can report which tier provided each value
- [ ] `--verbose` or `--dry-run` displays provenance for resolved configuration
- [ ] Config file validated at load time — invalid YAML produces clear error with file path and line number
- [ ] Unknown top-level keys in config files produce warnings (not errors) for forward compatibility
- [ ] Missing config files are silently skipped (not errors)
- [ ] Procedure `default_max_iterations` overrides `loop.default_max_iterations` for that procedure
- [ ] `ROODA_AI_CLI` environment variable sets the AI CLI command
- [ ] `ROODA_MAX_ITERATIONS` environment variable sets the default max iterations
- [ ] Prompt file paths in procedures resolved relative to config file location (workspace config) or as embedded resources (built-in defaults)
- [ ] Built-in default procedures include all 16 v2 procedures with embedded prompt files
- [ ] Custom procedures can reference prompt files on the filesystem

## Data Structures

### Config

The fully resolved configuration after merging all tiers.

```go
type Config struct {
    Loop       LoopConfig              // Global loop settings
    Procedures map[string]Procedure    // Named procedure definitions
    AITools    map[string]string       // AI CLI preset name -> command
    Provenance map[string]ConfigSource // Setting path -> source that provided it
}
```

**Fields:**
- `Loop` — Global defaults for iteration behavior
- `Procedures` — Map of procedure names to their definitions; includes both built-in and user-defined
- `AITools` — Map of preset names to AI CLI command strings; includes built-in presets (`kiro-cli`, `claude`, `copilot`, `cursor-agent`) and user-defined
- `Provenance` — Records which tier provided each resolved value, keyed by dot-path (e.g., `"loop.default_max_iterations"`, `"procedures.build.default_max_iterations"`)

### LoopConfig

```go
type LoopConfig struct {
    DefaultMaxIterations int // Global default (built-in default: 5)
    FailureThreshold     int // Consecutive failures before abort (built-in default: 3)
}
```

### Procedure

```go
type Procedure struct {
    Display              string // Human-readable name (optional)
    Summary              string // One-line description (optional)
    Description          string // Detailed description (optional)
    Observe              string // Path to observe phase prompt file, or embedded resource name
    Orient               string // Path to orient phase prompt file, or embedded resource name
    Decide               string // Path to decide phase prompt file, or embedded resource name
    Act                  string // Path to act phase prompt file, or embedded resource name
    DefaultMaxIterations *int   // Override loop.default_max_iterations for this procedure (nil = use global)
}
```

**Fields:**
- `Observe`, `Orient`, `Decide`, `Act` — Paths to prompt markdown files. For built-in procedures, these reference embedded resources. For user-defined procedures, these are filesystem paths resolved relative to the config file location.
- `DefaultMaxIterations` — Optional per-procedure override. Nil means inherit from `loop.default_max_iterations`.

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
    TierGlobal    ConfigTier = "global"    // ~/.config/rooda/config.yml
    TierWorkspace ConfigTier = "workspace" // ./rooda-config.yml
    TierEnvVar    ConfigTier = "env"       // ROODA_* environment variables
    TierCLIFlag   ConfigTier = "cli"       // --flag values
)
```

### YAML Config File Schema

Both workspace and global config files share the same schema:

```yaml
# Loop settings (all optional — built-in defaults apply)
loop:
  default_max_iterations: 5    # Default max iterations for all procedures
  failure_threshold: 3         # Consecutive failures before abort

# AI tool presets (optional — merges with built-in presets)
ai_tools:
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

  my-custom-procedure:
    observe: my-prompts/observe.md
    orient: my-prompts/orient.md
    decide: my-prompts/decide.md
    act: my-prompts/act.md
    default_max_iterations: 3
```

### Built-in Defaults

These are compiled into the Go binary and always available:

```go
var BuiltInDefaults = Config{
    Loop: LoopConfig{
        DefaultMaxIterations: 5,
        FailureThreshold:     3,
    },
    AITools: map[string]string{
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
| `ROODA_AI_CLI` | AI CLI command (direct, not a preset) | string |
| `ROODA_MAX_ITERATIONS` | `loop.default_max_iterations` | int |
| `ROODA_FAILURE_THRESHOLD` | `loop.failure_threshold` | int |

Environment variables use the `ROODA_` prefix. They override config file values but are overridden by CLI flags.

## Algorithm

### Config Loading

```
function LoadConfig(cliFlags CLIFlags) -> (Config, error):
    // 1. Start with built-in defaults
    config = BuiltInDefaults.DeepCopy()
    provenance = initProvenance(config, TierBuiltIn)

    // 2. Load global config (~/.config/rooda/config.yml)
    globalPath = expandHome("~/.config/rooda/config.yml")
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

    // 4. Apply environment variables
    if env("ROODA_AI_CLI") != "":
        // Stored separately — not a preset, it's a direct command
        provenance["ai_cli"] = ConfigSource{TierEnvVar, "", env("ROODA_AI_CLI")}
    if env("ROODA_MAX_ITERATIONS") != "":
        config.Loop.DefaultMaxIterations = parseInt(env("ROODA_MAX_ITERATIONS"))
        provenance["loop.default_max_iterations"] = ConfigSource{TierEnvVar, "", config.Loop.DefaultMaxIterations}
    if env("ROODA_FAILURE_THRESHOLD") != "":
        config.Loop.FailureThreshold = parseInt(env("ROODA_FAILURE_THRESHOLD"))
        provenance["loop.failure_threshold"] = ConfigSource{TierEnvVar, "", config.Loop.FailureThreshold}

    // 5. Apply CLI flags (highest precedence)
    if cliFlags.MaxIterations != nil:
        config.Loop.DefaultMaxIterations = *cliFlags.MaxIterations
        provenance["loop.default_max_iterations"] = ConfigSource{TierCLIFlag, "", *cliFlags.MaxIterations}
    // Additional CLI flag overrides applied by caller (--ai-cli, --ai-tool, etc.)

    // 6. Validate
    err = validateConfig(config)
    if err:
        return error("config validation: %v", err)

    config.Provenance = provenance
    return config, nil
```

### Config Merging

```
function mergeConfig(base Config, overlay ConfigFile, provenance map, tier ConfigTier, filePath string) -> Config:
    // Merge loop settings (overlay wins if non-zero)
    if overlay.Loop.DefaultMaxIterations != 0:
        base.Loop.DefaultMaxIterations = overlay.Loop.DefaultMaxIterations
        provenance["loop.default_max_iterations"] = ConfigSource{tier, filePath, overlay.Loop.DefaultMaxIterations}
    if overlay.Loop.FailureThreshold != 0:
        base.Loop.FailureThreshold = overlay.Loop.FailureThreshold
        provenance["loop.failure_threshold"] = ConfigSource{tier, filePath, overlay.Loop.FailureThreshold}

    // Merge AI tool presets (overlay adds to or overrides individual presets)
    for name, command in overlay.AITools:
        base.AITools[name] = command
        provenance["ai_tools." + name] = ConfigSource{tier, filePath, command}

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
    if config.Loop.DefaultMaxIterations < 0:
        return error("loop.default_max_iterations must be >= 0, got %d", config.Loop.DefaultMaxIterations)
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
        if proc.DefaultMaxIterations != nil && *proc.DefaultMaxIterations < 0:
            return error("procedure %s: default_max_iterations must be >= 0", name)

    // Validate AI tools (values must be non-empty strings)
    for name, command in config.AITools:
        if command == "":
            return error("ai_tools.%s: command must be non-empty", name)

    return nil
```

### AI CLI Command Resolution

The AI CLI command is resolved from multiple sources with precedence. This happens after config loading, using both the merged config and CLI flags:

```
function ResolveAICLICommand(config Config, cliFlags CLIFlags) -> (string, ConfigSource):
    // 1. --ai-cli flag (direct command, highest precedence)
    if cliFlags.AICLI != "":
        return cliFlags.AICLI, ConfigSource{TierCLIFlag, "", cliFlags.AICLI}

    // 2. --ai-tool flag (preset name, resolved from merged config)
    if cliFlags.AITool != "":
        command, exists = config.AITools[cliFlags.AITool]
        if !exists:
            return error("unknown AI tool preset: %s\nAvailable: %v", cliFlags.AITool, keys(config.AITools))
        return command, ConfigSource{TierCLIFlag, "", command}

    // 3. ROODA_AI_CLI environment variable
    if source, exists = config.Provenance["ai_cli"]; exists && source.Tier == TierEnvVar:
        return source.Value.(string), source

    // 4. Built-in default
    return "kiro-cli chat --no-interactive --trust-all-tools", ConfigSource{TierBuiltIn, "", "kiro-cli chat --no-interactive --trust-all-tools"}
```

### Max Iterations Resolution

```
function ResolveMaxIterations(config Config, procedureName string, cliFlags CLIFlags) -> (*int, ConfigSource):
    // 1. --max-iterations CLI flag (highest)
    if cliFlags.MaxIterations != nil:
        return cliFlags.MaxIterations, ConfigSource{TierCLIFlag, "", *cliFlags.MaxIterations}

    // 2. --unlimited CLI flag
    if cliFlags.Unlimited:
        return nil, ConfigSource{TierCLIFlag, "", "unlimited"}

    // 3. Procedure default_max_iterations
    proc, exists = config.Procedures[procedureName]
    if exists && proc.DefaultMaxIterations != nil:
        return proc.DefaultMaxIterations, config.Provenance["procedures." + procedureName]

    // 4. loop.default_max_iterations (already resolved through tier precedence)
    val = config.Loop.DefaultMaxIterations
    return &val, config.Provenance["loop.default_max_iterations"]
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No config files exist, no env vars | Works using built-in defaults only (zero-config) |
| Global config exists, no workspace config | Global config merges with built-in defaults |
| Both global and workspace config exist | Workspace overrides global for overlapping values; non-overlapping values from both are kept |
| Workspace config defines procedure that overrides a built-in | Workspace procedure replaces the built-in procedure entirely (all fields) |
| Workspace config defines new procedure not in built-ins | New procedure added alongside built-in procedures |
| Global config defines procedure, workspace does not | Global procedure is available alongside built-in procedures |
| Environment variable `ROODA_MAX_ITERATIONS` set to non-integer | Error: "ROODA_MAX_ITERATIONS must be an integer, got 'abc'" |
| Environment variable `ROODA_MAX_ITERATIONS=0` | Sets default max iterations to 0 (unlimited) |
| Config file has invalid YAML | Error with file path and parse error details |
| Config file has unknown top-level key | Warning logged, key ignored (forward compatibility) |
| Config file has unknown key inside `procedures` | Warning logged, key ignored |
| Procedure references prompt file that doesn't exist | Error at validation time: "procedure 'X': observe file not found: path/to/file.md" |
| Built-in procedure prompt file (embedded) | Loaded from Go binary via `go:embed`, not filesystem |
| Workspace procedure prompt file path | Resolved relative to workspace config file directory |
| Global procedure prompt file path | Resolved relative to global config file directory |
| `~/.config/rooda/` directory doesn't exist | Silently skipped (no global config) |
| Config file is empty | Valid — all built-in defaults apply |
| Config file has only `ai_tools` section | Valid — procedures and loop use built-in defaults |
| Two config files define same AI tool preset | Workspace wins over global, both override built-in |
| `loop.failure_threshold: 0` in config | Error: "loop.failure_threshold must be >= 1" |
| `loop.default_max_iterations: 0` in config | Valid — means unlimited iterations |
| Procedure `default_max_iterations: 0` in config | Valid — means unlimited iterations for that procedure |
| CLI `--max-iterations 0` | Not valid — use `--unlimited` flag instead for clarity |
| Workspace config and `--config` flag both present | `--config` flag specifies which config file to load as workspace config |

## Dependencies

- **Go standard library** — `os`, `path/filepath`, `encoding/json` for file operations and path resolution
- **Go YAML library** — `gopkg.in/yaml.v3` or `github.com/goccy/go-yaml` for YAML parsing (replaces yq dependency)
- **go:embed** — For embedding built-in default prompts and config in the binary
- **cli-interface** — Provides CLI flags that override config values
- **prompt-composition** — Consumes procedure definitions to assemble prompts
- **ai-cli-integration** — Consumes AI tool presets and resolved AI CLI command
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
- `iteration-loop.md` — Consumes `LoopConfig` and max iterations resolution
- `prompt-composition.md` — Consumes `Procedure` definitions to assemble prompts
- `ai-cli-integration.md` — Consumes `AITools` presets and resolved AI CLI command
- `cli-interface.md` — Defines CLI flags that override config values
- `procedures.md` — Defines the 16 built-in procedures that ship as defaults

## Examples

### Example 1: Zero-Config Startup

**Scenario:** Developer installs rooda and runs it immediately with no config files.

**Input:**
```bash
rooda build
```

**Config resolution:**
```
loop.default_max_iterations: 5        (built-in)
loop.failure_threshold: 3             (built-in)
procedure: build                      (built-in, embedded prompts)
ai_cli: kiro-cli chat --no-interactive --trust-all-tools (built-in)
```

**Verification:**
- No config files needed
- Built-in `build` procedure uses embedded prompt files
- Default 5 iterations, 3 failure threshold
- kiro-cli used as default AI CLI

### Example 2: Workspace Config Override

**Scenario:** Team wants to use claude instead of kiro-cli and increase build iterations.

**Workspace config (`./rooda-config.yml`):**
```yaml
ai_tools:
  default: "claude-cli --no-interactive"

procedures:
  build:
    default_max_iterations: 10
```

**Input:**
```bash
rooda build --ai-tool default
```

**Config resolution:**
```
loop.default_max_iterations: 5        (built-in)
loop.failure_threshold: 3             (built-in)
procedure: build
  observe/orient/decide/act:          (built-in, embedded — not overridden)
  default_max_iterations: 10          (workspace: ./rooda-config.yml)
ai_cli: claude-cli --no-interactive   (cli: --ai-tool resolved from workspace preset)
```

**Verification:**
- Build procedure prompt files still use built-in embedded defaults (workspace only overrode `default_max_iterations`)
- Iteration limit increased to 10
- Claude CLI used via custom preset
- Provenance shows `default_max_iterations` came from workspace config

### Example 3: Three-Tier Merge

**Scenario:** Developer has global preferences and project-specific overrides.

**Global config (`~/.config/rooda/config.yml`):**
```yaml
loop:
  default_max_iterations: 8

ai_tools:
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
```

**Input:**
```bash
rooda my-lint
```

**Config resolution:**
```
loop.default_max_iterations: 8        (global: ~/.config/rooda/config.yml)
loop.failure_threshold: 5             (workspace: ./rooda-config.yml)
ai_tools.fast: kiro-cli ... haiku     (global: ~/.config/rooda/config.yml)
ai_tools.kiro-cli: kiro-cli ...       (built-in)
ai_tools.claude: claude-cli ...       (built-in)
ai_tools.copilot: github-copilot-cli  (built-in)
ai_tools.cursor-agent: cursor-agent ...(built-in)
procedure: my-lint                    (workspace: ./rooda-config.yml)
  default_max_iterations: 1           (workspace: ./rooda-config.yml)
ai_cli: kiro-cli chat --no-interactive --trust-all-tools (built-in)
```

**Verification:**
- `default_max_iterations` from global (8), but overridden by procedure-level (1) for `my-lint`
- `failure_threshold` from workspace (5) overrides global which used built-in default
- `fast` AI tool preset from global is available
- All built-in AI tool presets still available
- All 16 built-in procedures still available alongside custom `my-lint`
- Custom procedure prompt paths resolved relative to workspace config directory

### Example 4: Environment Variable Override

**Input:**
```bash
export ROODA_AI_CLI="aider --yes"
export ROODA_MAX_ITERATIONS=3
rooda build
```

**Config resolution:**
```
loop.default_max_iterations: 3        (env: ROODA_MAX_ITERATIONS)
loop.failure_threshold: 3             (built-in)
procedure: build                      (built-in, embedded prompts)
ai_cli: aider --yes                   (env: ROODA_AI_CLI)
```

**Verification:**
- Environment variables override config file values
- `ROODA_AI_CLI` sets AI CLI command directly (not a preset name)
- `ROODA_MAX_ITERATIONS` overrides `loop.default_max_iterations`
- Built-in procedure still used

### Example 5: CLI Flag Takes Precedence Over Everything

**Input:**
```bash
export ROODA_MAX_ITERATIONS=3
rooda build --max-iterations 1 --ai-cli "claude-cli --no-interactive"
```

**Config resolution:**
```
max_iterations: 1                     (cli: --max-iterations)
ai_cli: claude-cli --no-interactive   (cli: --ai-cli)
```

**Verification:**
- CLI `--max-iterations 1` overrides env var `ROODA_MAX_ITERATIONS=3`
- CLI `--ai-cli` overrides env var `ROODA_AI_CLI`
- Highest precedence wins

### Example 6: Provenance Display (Dry-Run)

**Input:**
```bash
rooda build --dry-run
```

**Expected Output (provenance section):**
```
[DRY RUN] Configuration provenance:
  loop.default_max_iterations: 5 (built-in)
  loop.failure_threshold: 3 (built-in)
  procedure: build (built-in)
  ai_cli: kiro-cli chat --no-interactive --trust-all-tools (built-in)

[DRY RUN] Procedure: build
[DRY RUN] Would execute with: kiro-cli chat --no-interactive --trust-all-tools
```

**Verification:**
- Every resolved setting shows its source tier
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
# Global config at ~/.config/rooda/config.yml loaded if present
# /path/to/team-config.yml loaded as workspace config (instead of ./rooda-config.yml)
```

**Verification:**
- `--config` flag replaces the workspace config file path
- Global config still loads from `~/.config/rooda/config.yml`
- Prompt paths in team config resolved relative to `/path/to/`

## Notes

**Design Rationale — Three-Tier System:**

The three-tier system (built-in > global > workspace) mirrors established tooling conventions (Git, npm, EditorConfig). Each tier serves a distinct use case:
- **Built-in defaults** — Zero-config startup, always available, ships with the binary
- **Global config** — Personal preferences that follow the developer across projects (preferred AI tool, iteration limits)
- **Workspace config** — Project-specific settings committed to the repo (custom procedures, team presets)

This separation means a developer can set their preferred AI tool globally, while each project can define custom procedures — without conflict.

**Why Merge Instead of Replace:**

Config merging (additive) instead of replacement (destructive) is critical for a good experience. If a workspace config defining one custom procedure replaced all built-in procedures, the developer would lose access to `build`, `bootstrap`, etc. Merging means workspace configs only need to specify what's different.

**Provenance Tracking:**

Provenance answers the question "where did this value come from?" — essential for debugging configuration issues. When a developer runs `rooda build --dry-run` and sees unexpected settings, provenance immediately shows which file or environment variable set each value.

**Environment Variable Convention:**

The `ROODA_` prefix follows Go CLI conventions and avoids namespace collisions. Only a small set of environment variables is supported — complex configuration belongs in YAML files.

**Path Resolution:**

Prompt file paths in config files are resolved relative to the config file's directory, not the current working directory. This ensures a workspace config committed to a repo works regardless of where `rooda` is invoked from within the project. Built-in procedure prompt files are embedded in the binary via `go:embed` and don't need filesystem resolution.

**Forward Compatibility:**

Unknown keys produce warnings, not errors. This allows newer config files to work with older rooda versions — new keys are simply ignored. This is important for teams where not everyone upgrades simultaneously.

**Migration from v1:**

v1 used a single `rooda-config.yml` file with yq for YAML parsing. The v2 workspace config file uses the same filename and a compatible schema, so existing v1 config files work as v2 workspace configs with minimal changes (rename `default_iterations` to `default_max_iterations`).
