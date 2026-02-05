# AI CLI Integration

## Job to be Done
Execute OODA loop prompts through a configurable AI CLI tool that can read files, modify code, run commands, and interact with the repository autonomously.

## Activities
1. Resolve AI CLI command from configuration (flag > config > default)
2. Pipe assembled OODA prompt to AI CLI via stdin
3. Pass flags to enable autonomous operation (no interactive prompts)
4. Trust all tool invocations without permission prompts
5. Allow AI to read/write files, execute commands, and commit changes
6. Capture AI CLI exit status for error handling

## Configuration

### Precedence System

AI CLI command resolution follows this precedence order (highest to lowest):

1. `--ai-cli` flag - Direct command override
2. `--ai-tool` preset - Resolves to command via hardcoded presets or config
3. `$ROODA_AI_CLI` environment variable
4. Default: `kiro-cli chat --no-interactive --trust-all-tools`

### --ai-cli Flag

Override the AI CLI command with a direct command string:

```bash
./rooda.sh build --ai-cli "custom-cli --flags"
```

**Properties:**
- **Type:** String (full command with flags)
- **Precedence:** Highest (overrides all other settings)
- **Use case:** One-off custom commands, testing new tools

### --ai-tool Preset

Specify an AI tool by preset name:

```bash
./rooda.sh build --ai-tool claude
```

**Hardcoded presets:**
- `kiro-cli` → `kiro-cli chat --no-interactive --trust-all-tools`
- `claude` → `claude-cli --no-interactive`
- `aider` → `aider --yes`

**Custom presets** can be defined in `rooda-config.yml`:

```yaml
ai_tools:
  my-tool: "my-cli --autonomous --trust-tools"

procedures:
  bootstrap:
    # ... procedure config
```

**Properties:**
- **Type:** String (preset name)
- **Precedence:** Second (after --ai-cli flag)
- **Resolution:** Hardcoded presets checked first, then config `ai_tools` section
- **Error handling:** Unknown presets show helpful error with available options
- **Use case:** Team-wide tool standardization, convenient shortcuts

### $ROODA_AI_CLI Environment Variable

Set a default AI CLI command via environment variable:

```bash
export ROODA_AI_CLI="claude-cli --no-interactive"
./rooda.sh build
```

**Properties:**
- **Type:** String (full command with flags)
- **Precedence:** Third (after --ai-cli and --ai-tool)
- **Use case:** User-specific defaults, CI/CD environments

### Default

If no configuration is provided, defaults to:

```bash
kiro-cli chat --no-interactive --trust-all-tools
```

**Properties:**
- **Precedence:** Lowest (used when nothing else specified)
- **Backward compatibility:** Existing installations work without changes

## Acceptance Criteria
- [x] Prompt piped to AI CLI via stdin
- [x] AI CLI command configurable via rooda-config.yml
- [x] AI CLI command overridable via --ai-cli flag
- [x] Precedence: flag > config > default
- [x] Default remains kiro-cli for backward compatibility
- [x] --no-interactive flag (or equivalent) disables interactive prompts
- [x] --trust-all-tools flag (or equivalent) bypasses permission prompts
- [x] AI can read files from repository
- [x] AI can write/modify files in repository
- [x] AI can execute bash commands
- [x] AI can commit changes to git
- [x] Script continues to next iteration regardless of AI CLI exit status

## Data Structures

### AI CLI Command Resolution
```bash
# Precedence: --ai-cli > --ai-tool > $ROODA_AI_CLI > default
if [ -n "$AI_CLI_FLAG" ]; then
    AI_CLI_COMMAND="$AI_CLI_FLAG"
elif [ -n "$AI_TOOL_PRESET" ]; then
    AI_CLI_COMMAND=$(resolve_ai_tool_preset "$AI_TOOL_PRESET" "$CONFIG_FILE")
elif [ -n "$ROODA_AI_CLI" ]; then
    AI_CLI_COMMAND="$ROODA_AI_CLI"
else
    AI_CLI_COMMAND="kiro-cli chat --no-interactive --trust-all-tools"
fi
```

### AI CLI Invocation
```bash
create_prompt | $AI_CLI_COMMAND
```

**Components:**
- `create_prompt` - Function that assembles OODA prompt from four phase files
- `$AI_CLI_COMMAND` - Resolved AI CLI command (configurable)
- Default: `kiro-cli chat --no-interactive --trust-all-tools`

**Common AI CLI tools:**
- `kiro-cli chat --no-interactive --trust-all-tools` (default)
- `claude-cli --no-interactive`
- `aider --yes`
- Custom wrapper scripts

**Hardcoded presets:**
- `kiro-cli` - Kiro CLI with autonomous flags
- `claude` - Claude CLI with non-interactive mode
- `aider` - Aider with auto-yes mode

**Custom presets** defined in `rooda-config.yml`:
```yaml
ai_tools:
  my-tool: "my-cli --autonomous"
```

### Prompt Format
```markdown
# OODA Loop Iteration

## OBSERVE
[Content from observe phase file]

## ORIENT
[Content from orient phase file]

## DECIDE
[Content from decide phase file]

## ACT
[Content from act phase file]
```

## Algorithm

1. Resolve AI CLI command from configuration
   - Check for --ai-cli flag (highest priority - direct command)
   - Check for --ai-tool preset (resolve via hardcoded or config)
   - Check for $ROODA_AI_CLI environment variable
   - Fall back to default: `kiro-cli chat --no-interactive --trust-all-tools`
2. Assemble OODA prompt using `create_prompt` function
3. Pipe prompt to AI CLI via stdin
4. AI CLI reads prompt and executes OODA phases
5. AI reads files, analyzes situation, makes decisions
6. AI executes actions (modify files, run commands, commit changes)
7. AI CLI exits (status ignored by script)
8. Script continues to git push and next iteration

**Pseudocode:**
```bash
# Resolve AI CLI command
if [ -n "$AI_CLI_FLAG" ]; then
    AI_CLI_COMMAND="$AI_CLI_FLAG"
elif [ -n "$AI_TOOL_PRESET" ]; then
    AI_CLI_COMMAND=$(resolve_ai_tool_preset "$AI_TOOL_PRESET" "$CONFIG_FILE")
elif [ -n "$ROODA_AI_CLI" ]; then
    AI_CLI_COMMAND="$ROODA_AI_CLI"
else
    AI_CLI_COMMAND="kiro-cli chat --no-interactive --trust-all-tools"
fi

resolve_ai_tool_preset() {
    case "$preset" in
        kiro-cli) echo "kiro-cli chat --no-interactive --trust-all-tools" ;;
        claude) echo "claude-cli --no-interactive" ;;
        aider) echo "aider --yes" ;;
        *) 
            # Query custom preset from config ai_tools section
            custom=$(yq eval ".ai_tools.$preset" "$config_file")
            if [ "$custom" != "null" ]; then
                echo "$custom"
            else
                echo "Error: Unknown AI tool preset: $preset" >&2
                return 1
            fi
            ;;
    esac
}

create_prompt() {
    cat <<EOF
# OODA Loop Iteration

## OBSERVE
$(cat "$OBSERVE")

## ORIENT
$(cat "$ORIENT")

## DECIDE
$(cat "$DECIDE")

## ACT
$(cat "$ACT")
EOF
}

# Execute AI CLI
create_prompt | $AI_CLI_COMMAND
# Exit status not checked - script continues regardless
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| AI CLI not installed | Command fails, script exits with error |
| AI CLI exits with error | Script continues to git push (no error handling) |
| AI refuses to execute action | Iteration completes, next iteration may retry |
| AI modifies unexpected files | Changes committed and pushed (no validation) |
| AI executes dangerous command | Command runs (sandboxing required for safety) |
| Prompt exceeds token limit | AI CLI may truncate or fail (no size validation) |
| Network failure during AI call | AI CLI fails, script continues (no retry logic) |
| Invalid AI CLI command via --ai-cli | Command fails at runtime, script exits |
| Unknown preset via --ai-tool | Error with available presets listed, script exits |
| Custom preset not in config | Error with instructions to add to config, script exits |
| $ROODA_AI_CLI with invalid command | Command fails at runtime, script exits |
| AI CLI doesn't support stdin | Script fails, no fallback mechanism |
| Multiple configuration methods | Precedence: --ai-cli > --ai-tool > $ROODA_AI_CLI > default |

## Dependencies

- AI CLI tool (configurable, defaults to kiro-cli)
- AI CLI must support:
  - Reading prompts from stdin
  - Non-interactive operation mode
  - Tool invocation without permission prompts
  - File read/write capabilities
  - Command execution capabilities

## Implementation Mapping

**Source files:**
- `src/rooda.sh` - Lines 143-159 implement `create_prompt` function
- `src/rooda.sh` - Line 169 implements AI CLI invocation

**Related specs:**
- `component-authoring.md` - Defines how OODA phases are assembled
- `iteration-loop.md` - Defines loop execution behavior
- `cli-interface.md` - Defines command-line argument parsing

## Examples

### Example 1: Default AI CLI (kiro-cli)

**Input:**
```bash
./rooda.sh build
```

**Expected Output:**
```
[Prompt piped to kiro-cli chat --no-interactive --trust-all-tools]
[AI reads files, analyzes, makes decisions, executes actions]
[AI commits changes]
[AI CLI exits with status 0]
```

**Verification:**
- Files modified by AI exist on disk
- Git commits created by AI
- Script continues to next iteration

### Example 2: Preset via --ai-tool Flag

**Input:**
```bash
./rooda.sh build --ai-tool claude
```

**Expected Output:**
```
[Prompt piped to claude-cli --no-interactive]
[AI executes OODA loop]
```

**Verification:**
- claude-cli invoked instead of kiro-cli
- Iteration completes successfully

### Example 3: Custom Preset from Config

**Config (rooda-config.yml):**
```yaml
ai_tools:
  my-tool: "my-cli --autonomous --trust-tools"

procedures:
  build:
    # ... procedure config
```

**Input:**
```bash
./rooda.sh build --ai-tool my-tool
```

**Expected Output:**
```
[Prompt piped to my-cli --autonomous --trust-tools]
[AI executes OODA loop]
```

**Verification:**
- Custom preset resolved from config
- my-cli invoked
- Iteration completes successfully

### Example 4: Environment Variable

**Input:**
```bash
export ROODA_AI_CLI="aider --yes"
./rooda.sh build
```

**Expected Output:**
```
[Prompt piped to aider --yes]
[AI executes OODA loop]
```

**Verification:**
- Environment variable used (no flag or preset specified)
- aider invoked
- Iteration completes successfully

### Example 5: Override via --ai-cli Flag

**Input:**
```bash
./rooda.sh build --ai-cli "custom-cli --flags"
```

**Expected Output:**
```
[Prompt piped to custom-cli --flags]
[AI executes OODA loop]
```

**Verification:**
- --ai-cli flag overrides all other settings
- custom-cli invoked
- Iteration completes successfully

### Example 6: Unknown Preset

**Input:**
```bash
./rooda.sh build --ai-tool nonexistent
```

**Expected Output:**
```
Error: Unknown AI tool preset: nonexistent

Available hardcoded presets:
  - kiro-cli
  - claude
  - aider

To define custom presets, add to rooda-config.yml:
  ai_tools:
    nonexistent: "your-command-here"
```

**Verification:**
- Script exits with error
- Helpful message shows available presets
- Instructions for adding custom preset

### Example 7: AI CLI Not Installed

**Input:**
```bash
./rooda.sh build --ai-cli "nonexistent-cli"
```

**Expected Output:**
```
bash: nonexistent-cli: command not found
```

**Verification:**
- Script exits with error
- No iteration executed

### Example 8: AI Refuses Action

**Input:**
```bash
./rooda.sh build
```

**Expected Output:**
```
[AI analyzes situation]
[AI responds: "I cannot complete this action because..."]
[AI CLI exits]
```

**Verification:**
- No files modified
- No commits created
- Script continues to next iteration (may retry)

### Example 9: Precedence - Flag Overrides Environment

**Input:**
```bash
export ROODA_AI_CLI="aider --yes"
./rooda.sh build --ai-tool claude
```

**Expected Output:**
```
[Prompt piped to claude-cli --no-interactive]
[AI executes OODA loop]
```

**Verification:**
- --ai-tool preset overrides $ROODA_AI_CLI
- claude-cli invoked (not aider)
- Precedence order respected

## Notes

**Design Rationale:**

The AI CLI integration is designed for autonomous operation with minimal human intervention. The four-tier precedence system provides flexibility for different use cases:

1. **--ai-cli flag** - Direct command override for one-off experiments
2. **--ai-tool preset** - Convenient shortcuts for common tools (team standardization)
3. **$ROODA_AI_CLI** - User-specific defaults without modifying config
4. **Default (kiro-cli)** - Backward compatibility for existing users

**Preset System:**

Presets simplify AI CLI configuration by providing named shortcuts. Hardcoded presets (kiro-cli, claude, aider) work out-of-the-box. Custom presets in `rooda-config.yml` enable team-specific tools without hardcoding in the script.

**Configuration Flexibility:**

The precedence system enables:
- **Individual developers** - Set $ROODA_AI_CLI for personal preference
- **Teams** - Define custom presets in rooda-config.yml for consistency
- **Experimentation** - Use --ai-cli or --ai-tool flags without modifying config
- **CI/CD** - Set $ROODA_AI_CLI in environment for automated workflows

**Security Implications:**

The AI CLI must support autonomous operation (no interactive prompts, no permission prompts for tool invocations). This is inherently risky and requires sandboxed execution environments (Docker, Fly Sprites, E2B) to limit blast radius.

**Error Handling:**

The script does not check AI CLI exit status. This design choice allows the loop to continue even if the AI encounters errors or refuses actions. The assumption is that subsequent iterations can self-correct through empirical feedback.

**AI CLI Requirements:**

Any AI CLI tool can be used if it supports:
- Reading prompts from stdin
- Non-interactive operation mode
- Tool invocation without permission prompts
- File read/write capabilities
- Command execution capabilities

**Token Limits:**

The script does not validate prompt size before piping to AI CLI. Large OODA phase files or extensive file contents could exceed token limits. The AI CLI is responsible for handling this (truncation, error, or chunking).

## Known Issues

**No error handling:** Script continues to git push even if AI CLI fails. This could result in pushing incomplete or invalid changes.

**No retry logic:** If AI CLI fails due to transient issues (network, rate limits), the iteration is lost. No automatic retry mechanism exists.

**No validation:** Script does not validate that AI CLI command is valid or that the tool supports required capabilities before invocation. Incompatible tools will fail at runtime.

**No timeout:** If AI CLI hangs, the script waits indefinitely. No timeout mechanism exists.

## Areas for Improvement

**Dependency checking:** Add validation that configured AI CLI is installed and accessible before starting loop.

**Capability detection:** Detect if AI CLI supports required features (stdin, non-interactive mode, tool invocation) and provide clear error messages if not.

**Error handling:** Check AI CLI exit status and handle failures gracefully (retry, skip push, abort loop).

**Timeout mechanism:** Add timeout for AI CLI invocation to prevent indefinite hangs.

**Prompt size validation:** Check assembled prompt size before piping to AI CLI, warn if approaching token limits.

**AI CLI profiles:** Support multiple AI CLI configurations (e.g., "fast" vs "thorough" models) selectable per procedure or via flag.

**Version requirements:** Document minimum version requirements for supported AI CLI tools.
