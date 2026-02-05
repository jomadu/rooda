# Task: Flexible AI CLI Configuration

## Problem

Multi-user teams can't have individual AI CLI preferences because config is committed to git. User A commits their preferred AI CLI, User B pulls and their workflow breaks.

## Solution

Environment variable + presets system. No RC files needed.

### Precedence (Highest to Lowest)

1. `--ai-cli "command"` - Explicit CLI flag override
2. `--ai-tool preset` - Named preset (hardcoded or custom)
3. `$ROODA_AI_CLI` - Environment variable (primary user preference)
4. `ai_tools` presets in rooda-config.yml - Project-defined presets
5. Default: `kiro-cli chat --no-interactive --trust-all-tools`

## Implementation

### 1. Hardcoded Presets

Built into script for convenience:

```bash
resolve_ai_tool_preset() {
    local preset="$1"
    
    case "$preset" in
        kiro-cli)
            echo "kiro-cli chat --no-interactive --trust-all-tools"
            return 0
            ;;
        claude)
            echo "claude-cli --autonomous --trust-tools"
            return 0
            ;;
        aider)
            echo "aider --yes --auto-commits"
            return 0
            ;;
    esac
    
    # Check custom presets in config
    local custom_preset=$(yq eval ".ai_tools.$preset" "$CONFIG_FILE" 2>&1)
    if [ "$custom_preset" != "null" ] && [ -n "$custom_preset" ]; then
        echo "$custom_preset"
        return 0
    fi
    
    echo "Error: Unknown AI tool preset: $preset" >&2
    echo "" >&2
    echo "Available hardcoded presets: kiro-cli, claude, aider" >&2
    echo "Check rooda-config.yml for custom presets in 'ai_tools' section" >&2
    return 1
}
```

### 2. Custom Presets in Config

Optional `ai_tools` section in rooda-config.yml:

```yaml
# Optional: Define custom AI tool presets
ai_tools:
  fast: "kiro-cli chat --model haiku --no-interactive --trust-all-tools"
  thorough: "claude-cli --model opus --autonomous"

procedures:
  bootstrap:
    observe: src/prompts/observe_bootstrap.md
    orient: src/prompts/orient_bootstrap.md
    decide: src/prompts/decide_bootstrap.md
    act: src/prompts/act_bootstrap.md
    default_iterations: 1
```

### 3. Precedence Resolution

```bash
# Initialize with default
AI_CLI_COMMAND="kiro-cli chat --no-interactive --trust-all-tools"

# Apply environment variable (user preference)
[ -n "$ROODA_AI_CLI" ] && AI_CLI_COMMAND="$ROODA_AI_CLI"

# Apply --ai-tool preset (if specified)
if [ -n "$AI_TOOL_PRESET" ]; then
    AI_CLI_COMMAND=$(resolve_ai_tool_preset "$AI_TOOL_PRESET") || exit 1
fi

# Apply --ai-cli flag (highest priority, already set during parsing)
```

### 4. Argument Parsing

```bash
while [[ $# -gt 0 ]]; do
    case $1 in
        --ai-cli)
            AI_CLI_COMMAND="$2"
            shift 2
            ;;
        --ai-tool)
            AI_TOOL_PRESET="$2"
            shift 2
            ;;
        # ... other flags
    esac
done
```

## Breaking Changes

**Remove:** Per-procedure `ai_cli_command` field
- Current location: `.procedures.$PROCEDURE.ai_cli_command` in rooda.sh
- Reason: Locks users to single AI CLI, causes multi-user conflicts
- Migration: None needed (only maintainer uses this)

## Files to Modify

**src/rooda.sh:**
- Remove per-procedure `ai_cli_command` query logic
- Add `resolve_ai_tool_preset()` function
- Add `--ai-tool` flag parsing
- Add precedence resolution logic
- Add `$ROODA_AI_CLI` environment variable support

**src/rooda-config.yml:**
- Add `ai_tools` section with example presets (commented)
- Remove any per-procedure `ai_cli_command` fields

**Documentation:**
- README.md - Add "Configuring Your AI CLI" section
- specs/ai-cli-integration.md - Update configuration section
- specs/configuration-schema.md - Add `ai_tools` section docs
- specs/cli-interface.md - Add `--ai-tool` flag docs

## User Experience

### Team member using Claude

```bash
# Add to ~/.zshrc once
echo 'export ROODA_AI_CLI="claude-cli --autonomous --trust-tools"' >> ~/.zshrc

# Works everywhere
./rooda.sh build
```

### One-off override

```bash
./rooda.sh build --ai-tool aider
```

### Project-specific preset

```yaml
# rooda-config.yml (committed)
ai_tools:
  fast: "kiro-cli chat --model haiku --no-interactive"
  thorough: "claude-cli --model opus --autonomous"
```

```bash
./rooda.sh build --ai-tool fast
```

### Explicit command

```bash
./rooda.sh build --ai-cli "my-custom-cli --flags"
```

## Acceptance Criteria

- [ ] `--ai-cli "command"` flag overrides all other sources
- [ ] `--ai-tool preset` resolves hardcoded presets (kiro-cli, claude, aider)
- [ ] `--ai-tool preset` resolves custom presets from `ai_tools` section
- [ ] `$ROODA_AI_CLI` environment variable works
- [ ] Unknown preset shows helpful error with available options
- [ ] Per-procedure `ai_cli_command` removed from config and code
- [ ] Precedence verified: `--ai-cli` > `--ai-tool` > `$ROODA_AI_CLI` > default
- [ ] Default remains `kiro-cli chat --no-interactive --trust-all-tools`
- [ ] Backward compatible: existing users without config continue to work

## Testing Checklist

- [ ] Test `--ai-cli` flag override
- [ ] Test `--ai-tool` with hardcoded presets (kiro-cli, claude, aider)
- [ ] Test `--ai-tool` with custom preset from config
- [ ] Test `$ROODA_AI_CLI` environment variable
- [ ] Test unknown preset shows helpful error
- [ ] Test precedence order
- [ ] Test backward compatibility (no config, uses default)
- [ ] Test cross-platform: macOS, Linux, Windows (Git Bash/WSL)
