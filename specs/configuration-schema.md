# Configuration Schema

## Job to be Done
Enable users to define custom OODA loop procedures by mapping procedure names to composable prompt component files, supporting both predefined framework procedures and user-defined custom procedures.

## Activities
1. Define procedure in YAML configuration file
2. Map procedure name to four OODA phase prompt files
3. Specify optional metadata (display name, summary, description)
4. Set default iteration count per procedure
5. Query configuration at runtime using yq

## Acceptance Criteria
- [x] YAML structure supports nested procedure definitions
- [x] Required fields (observe, orient, decide, act) validated at runtime
- [x] Optional fields (display, summary, description, default_iterations) supported
- [x] Optional ai_tools section supports custom preset definitions
- [x] yq queries successfully extract procedure configuration
- [x] yq queries successfully extract ai_tools presets
- [x] Missing procedures return clear error messages
- [x] File paths resolved relative to script directory
- [x] ai_tools presets validated as string type
- [x] Unknown presets return helpful error messages

## Data Structures

### Configuration File Structure
```yaml
ai_tools:
  fast: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-5-haiku-20241022"
  thorough: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-7-sonnet-20250219"
  custom: "your-ai-cli-command-here"

procedures:
  procedure-name:
    display: "Human Readable Name"
    summary: "Brief one-line description"
    description: "Detailed multi-line description"
    observe: path/to/observe.md
    orient: path/to/orient.md
    decide: path/to/decide.md
    act: path/to/act.md
    default_iterations: 5
```

**Root-level fields:**
- `ai_tools` - Map of preset names to AI CLI commands (optional)

**Procedure fields:**
- `procedures` - Top-level map of procedure definitions (required)
- `procedure-name` - Unique identifier for procedure, kebab-case (required)
- `display` - Human-readable name for UI/help text (optional)
- `summary` - One-line description of procedure purpose (optional)
- `description` - Detailed explanation of procedure behavior (optional)
- `observe` - File path to observation phase prompt (required)
- `orient` - File path to orientation phase prompt (required)
- `decide` - File path to decision phase prompt (required)
- `act` - File path to action phase prompt (required)
- `default_iterations` - Default max iterations if not specified via CLI (optional, defaults to 0)

### AI Tools Section

The `ai_tools` section defines custom AI CLI tool presets that can be used with the `--ai-tool` flag.

**Structure:**
```yaml
ai_tools:
  preset-name: "ai-cli-command with flags"
```

**Hardcoded presets** (always available, no config needed):
- `kiro-cli` - `kiro-cli chat --no-interactive --trust-all-tools`
- `claude` - `claude-cli --no-interactive`
- `aider` - `aider --yes`

**Custom presets** can be defined in config to:
- Use different AI models (e.g., fast vs thorough)
- Configure team-specific AI CLI tools
- Set project-specific flags or options

**Usage:**
```bash
./rooda.sh build --ai-tool fast
./rooda.sh build --ai-tool thorough
./rooda.sh build --ai-tool kiro-cli  # hardcoded preset
```

## Algorithm

1. Parse command-line arguments to extract procedure name
2. Resolve config file path relative to script directory
3. Query config for procedure OODA files: `.procedures.$PROCEDURE.observe|orient|decide|act`
4. Validate all four OODA phase paths are non-null
5. Extract default_iterations if max-iterations not specified via CLI
6. If --ai-tool flag specified, query config for preset: `.ai_tools.$PRESET`
7. Return error if procedure not found or required fields missing

**Pseudocode:**
```bash
if procedure_specified:
    # Query procedure OODA files
    observe = yq eval ".procedures.$PROCEDURE.observe" config
    orient = yq eval ".procedures.$PROCEDURE.orient" config
    decide = yq eval ".procedures.$PROCEDURE.decide" config
    act = yq eval ".procedures.$PROCEDURE.act" config
    
    if any_field_is_null:
        error "Procedure not found"
    
    if max_iterations == 0:
        default_iter = yq eval ".procedures.$PROCEDURE.default_iterations" config
        if default_iter != null:
            max_iterations = default_iter

if ai_tool_preset_specified:
    # Query custom preset from config
    custom_command = yq eval ".ai_tools.$PRESET" config
    if custom_command != null:
        AI_CLI_COMMAND = custom_command
    else:
        # Check hardcoded presets or error
        error "Unknown AI tool preset: $PRESET"
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Procedure name not in config | Error: "Procedure 'name' not found in config" |
| Config file missing | Error: "config.yml not found" |
| Required OODA field missing | Error: "Procedure 'name' not found" (null check) |
| Invalid YAML syntax | yq parse error propagated to user |
| File path doesn't exist | Error occurs later when attempting to read prompt file |
| default_iterations not specified | Defaults to 0 (run indefinitely until Ctrl+C) |

## Dependencies

- yq - YAML query tool for parsing configuration
- rooda.sh - Script that reads and validates configuration
- Prompt component files - Markdown files referenced by OODA phase paths

## Implementation Mapping

**Source files:**
- `src/rooda-config.yml` - Default configuration with 9 predefined procedures
- `src/rooda.sh` - Lines 70-90 implement config lookup and validation

**Related specs:**
- `cli-interface.md` - Defines how procedure names are parsed from CLI
- `component-authoring.md` - Defines structure of prompt component files
- `iteration-loop.md` - Defines how default_iterations affects loop behavior

## Examples

### Example 1: Standard Procedure Definition

**Input:**
```yaml
procedures:
  build:
    display: "Build from Plan"
    summary: "Implements tasks from plan"
    observe: src/prompts/observe_plan_specs_impl.md
    orient: src/prompts/orient_build.md
    decide: src/prompts/decide_build.md
    act: src/prompts/act_build.md
    default_iterations: 5
```

**Query:**
```bash
yq eval ".procedures.build.observe" rooda-config.yml
# Returns: src/prompts/observe_plan_specs_impl.md
```

**Verification:**
- All four OODA phase paths resolve to existing files
- default_iterations is 5 when not overridden via CLI

### Example 2: Minimal Custom Procedure

**Input:**
```yaml
procedures:
  my-custom:
    observe: prompts/observe_specs.md
    orient: prompts/orient_gap.md
    decide: prompts/decide_gap_plan.md
    act: prompts/act_plan.md
```

**Query:**
```bash
./rooda.sh my-custom
```

**Verification:**
- Procedure executes with specified OODA files
- default_iterations defaults to 0 (infinite loop)
- No display/summary/description required

### Example 3: Custom AI Tool Presets

**Input:**
```yaml
ai_tools:
  fast: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-5-haiku-20241022"
  thorough: "kiro-cli chat --no-interactive --trust-all-tools --model claude-3-7-sonnet-20250219"

procedures:
  build:
    display: "Build from Plan"
    summary: "Implements tasks from plan"
    observe: src/prompts/observe_plan_specs_impl.md
    orient: src/prompts/orient_build.md
    decide: src/prompts/decide_build.md
    act: src/prompts/act_build.md
    default_iterations: 5
```

**Query:**
```bash
yq eval ".ai_tools.fast" rooda-config.yml
# Returns: kiro-cli chat --no-interactive --trust-all-tools --model claude-3-5-haiku-20241022
```

**Usage:**
```bash
./rooda.sh build --ai-tool fast
./rooda.sh build --ai-tool thorough
```

**Verification:**
- Custom presets resolve to configured commands
- Hardcoded presets (kiro-cli, claude, aider) always available
- Unknown presets return error with helpful message

### Example 4: Missing Procedure Error

**Input:**
```bash
./rooda.sh nonexistent-procedure
```

**Expected Output:**
```
Error: Procedure 'nonexistent-procedure' not found in rooda-config.yml
```

**Verification:**
- Script exits with non-zero status
- Clear error message identifies missing procedure

## Notes

**File Path Resolution:** All OODA phase file paths are resolved relative to the script directory, not the current working directory. This allows rooda.sh to be invoked from any location while maintaining consistent path resolution.

**yq Dependency:** The configuration system requires yq for YAML parsing. The script checks for yq availability at startup and provides installation instructions if missing.

**Procedure Naming:** Procedure names use kebab-case convention (lowercase with hyphens) to match CLI argument conventions.

**Optional Metadata:** The display, summary, and description fields are optional and primarily used for documentation and help text generation. They do not affect procedure execution.

## Known Issues

None identified during specification.

## Areas for Improvement

- Schema validation could be more comprehensive (validate file paths exist at config load time)
- Help text generation from config metadata not yet implemented
- No support for procedure aliases or inheritance
- No validation of default_iterations range (negative values not prevented)
