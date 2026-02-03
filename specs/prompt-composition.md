# Prompt Composition

## Job to be Done
Assemble four OODA phase prompt files into a single executable prompt that can be piped to the AI CLI.

## Activities
1. Read observe phase prompt file
2. Read orient phase prompt file
3. Read decide phase prompt file
4. Read act phase prompt file
5. Concatenate files with OODA section headers
6. Output combined prompt to stdout

## Acceptance Criteria
- [ ] Prompt assembly algorithm documented
- [ ] File validation behavior specified
- [ ] Error handling for missing files defined
- [ ] Output format produces valid markdown
- [ ] Section headers clearly delineate OODA phases

## Data Structures

### Input Variables
```bash
OBSERVE="path/to/observe.md"
ORIENT="path/to/orient.md"
DECIDE="path/to/decide.md"
ACT="path/to/act.md"
```

**Fields:**
- `OBSERVE` - Path to observation phase prompt file
- `ORIENT` - Path to orientation phase prompt file
- `DECIDE` - Path to decision phase prompt file
- `ACT` - Path to action phase prompt file

### Output Format
```markdown
# OODA Loop Iteration

## OBSERVE
[contents of observe file]

## ORIENT
[contents of orient file]

## DECIDE
[contents of decide file]

## ACT
[contents of act file]
```

## Algorithm

1. Validate all four prompt file paths are set
2. Validate all four files exist on disk
3. Create heredoc with OODA section structure
4. Embed file contents using command substitution
5. Output to stdout

**Pseudocode:**
```bash
function create_prompt():
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
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Missing prompt file path | Script exits with error before create_prompt() is called |
| File does not exist | Script exits with error before create_prompt() is called |
| Empty prompt file | File contents are empty, section header still present |
| File read error | Command substitution fails, bash reports error |

## Dependencies

- Bash shell with heredoc support
- File system access to read prompt files
- All four OODA phase files must exist before calling create_prompt()

## Implementation Mapping

**Source files:**
- `src/rooda.sh` (lines 143-159) - create_prompt() function implementation
- `src/rooda.sh` (lines 107-111) - File validation before prompt creation

**Related specs:**
- `cli-interface.md` - Defines how prompt is passed to AI CLI
- `configuration-schema.md` - Defines how prompt file paths are resolved from config

## Examples

### Example 1: Standard Prompt Assembly

**Input:**
```bash
OBSERVE="src/components/observe_plan_specs_impl.md"
ORIENT="src/components/orient_build.md"
DECIDE="src/components/decide_build.md"
ACT="src/components/act_build.md"
```

**Expected Output:**
```markdown
# OODA Loop Iteration

## OBSERVE
# Observe: Plan, Specs, Implementation

## O1: Study AGENTS.md as a Whole
[... full contents of observe file ...]

## ORIENT
# Orient: Build

## R5: Understand Task Requirements
[... full contents of orient file ...]

## DECIDE
# Decide: Build

## D4: Pick the Most Important Task
[... full contents of decide file ...]

## ACT
# Act: Build

## A3: Implement Using Parallel Subagents
[... full contents of act file ...]
```

**Verification:**
- Output contains all four OODA section headers
- Each section contains full file contents
- Markdown structure is valid
- Output can be piped to kiro-cli

### Example 2: File Validation Prevents Bad Prompt

**Input:**
```bash
OBSERVE="src/components/observe_plan_specs_impl.md"
ORIENT="src/components/nonexistent.md"
DECIDE="src/components/decide_build.md"
ACT="src/components/act_build.md"
```

**Expected Output:**
```
Error: File not found: src/components/nonexistent.md
```

**Verification:**
- Script exits with error code 1
- create_prompt() is never called
- Error message identifies missing file

## Notes

The create_prompt() function uses bash heredoc syntax to create a template with embedded command substitution. Each `$(cat "$VAR")` is evaluated when the heredoc is executed, inserting the file contents at that position.

File validation happens before create_prompt() is called (lines 107-111 in rooda.sh), ensuring all files exist before attempting to read them. This prevents partial prompt assembly.

The output is piped directly to kiro-cli, so the prompt must be valid markdown that the AI can parse and execute.

## Known Issues

None.

## Areas for Improvement

- Could add validation that prompt files contain expected markdown structure
- Could add size limits to prevent extremely large prompts that exceed AI context windows
- Could add caching to avoid re-reading unchanged prompt files across iterations
