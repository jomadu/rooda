# AI CLI Integration

## Job to be Done
Execute OODA loop prompts through an AI CLI tool that can read files, modify code, run commands, and interact with the repository autonomously.

## Activities
1. Pipe assembled OODA prompt to AI CLI via stdin
2. Pass flags to enable autonomous operation (no interactive prompts)
3. Trust all tool invocations without permission prompts
4. Allow AI to read/write files, execute commands, and commit changes
5. Capture AI CLI exit status for error handling

## Acceptance Criteria
- [x] Prompt piped to kiro-cli via stdin
- [x] --no-interactive flag disables interactive prompts
- [x] --trust-all-tools flag bypasses permission prompts
- [x] AI can read files from repository
- [x] AI can write/modify files in repository
- [x] AI can execute bash commands
- [x] AI can commit changes to git
- [x] Script continues to next iteration regardless of AI CLI exit status

## Data Structures

### AI CLI Invocation
```bash
create_prompt | kiro-cli chat --no-interactive --trust-all-tools
```

**Components:**
- `create_prompt` - Function that assembles OODA prompt from four phase files
- `kiro-cli chat` - AI CLI command for chat-based interaction
- `--no-interactive` - Flag to disable interactive prompts (no user input required)
- `--trust-all-tools` - Flag to bypass permission prompts for tool invocations

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

1. Assemble OODA prompt using `create_prompt` function
2. Pipe prompt to kiro-cli via stdin
3. kiro-cli reads prompt and executes OODA phases
4. AI reads files, analyzes situation, makes decisions
5. AI executes actions (modify files, run commands, commit changes)
6. kiro-cli exits (status ignored by script)
7. Script continues to git push and next iteration

**Pseudocode:**
```bash
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
create_prompt | kiro-cli chat --no-interactive --trust-all-tools
# Exit status not checked - script continues regardless
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| kiro-cli not installed | Command fails, script exits with error |
| kiro-cli exits with error | Script continues to git push (no error handling) |
| AI refuses to execute action | Iteration completes, next iteration may retry |
| AI modifies unexpected files | Changes committed and pushed (no validation) |
| AI executes dangerous command | Command runs (sandboxing required for safety) |
| Prompt exceeds token limit | kiro-cli may truncate or fail (no size validation) |
| Network failure during AI call | kiro-cli fails, script continues (no retry logic) |

## Dependencies

- kiro-cli - AI CLI tool (must be installed and in PATH)
- kiro-cli chat command - Chat-based interaction mode
- --no-interactive flag support - Must be supported by kiro-cli version
- --trust-all-tools flag support - Must be supported by kiro-cli version

## Implementation Mapping

**Source files:**
- `src/rooda.sh` - Lines 143-159 implement `create_prompt` function
- `src/rooda.sh` - Line 169 implements AI CLI invocation

**Related specs:**
- `../src/README.md` - Defines how OODA phases are assembled
- `iteration-loop.md` - Defines loop execution behavior
- `cli-interface.md` - Defines command-line argument parsing

## Examples

### Example 1: Successful Iteration

**Input:**
```bash
create_prompt | kiro-cli chat --no-interactive --trust-all-tools
```

**Expected Output:**
```
[AI reads files, analyzes, makes decisions, executes actions]
[AI commits changes]
[kiro-cli exits with status 0]
```

**Verification:**
- Files modified by AI exist on disk
- Git commits created by AI
- Script continues to next iteration

### Example 2: AI CLI Not Installed

**Input:**
```bash
create_prompt | kiro-cli chat --no-interactive --trust-all-tools
```

**Expected Output:**
```
bash: kiro-cli: command not found
```

**Verification:**
- Script exits with error
- No iteration executed

### Example 3: AI Refuses Action

**Input:**
```bash
create_prompt | kiro-cli chat --no-interactive --trust-all-tools
```

**Expected Output:**
```
[AI analyzes situation]
[AI responds: "I cannot complete this action because..."]
[kiro-cli exits]
```

**Verification:**
- No files modified
- No commits created
- Script continues to next iteration (may retry)

## Notes

**Design Rationale:**

The AI CLI integration is designed for autonomous operation with minimal human intervention. The `--no-interactive` and `--trust-all-tools` flags are critical for enabling the loop to run unattended.

**Security Implications:**

The `--trust-all-tools` flag bypasses all permission prompts, allowing the AI to execute arbitrary commands and modify any files. This is inherently risky and requires sandboxed execution environments (Docker, Fly Sprites, E2B) to limit blast radius.

**Error Handling:**

The script does not check kiro-cli exit status. This design choice allows the loop to continue even if the AI encounters errors or refuses actions. The assumption is that subsequent iterations can self-correct through empirical feedback.

**Alternative AI CLIs:**

While the implementation uses kiro-cli, the specification is designed to be compatible with any AI CLI that supports:
- Reading prompts from stdin
- Non-interactive operation
- Tool invocation without permission prompts
- File read/write capabilities
- Command execution capabilities

**Token Limits:**

The script does not validate prompt size before piping to kiro-cli. Large OODA phase files or extensive file contents could exceed token limits. The AI CLI is responsible for handling this (truncation, error, or chunking).

## Known Issues

**No error handling:** Script continues to git push even if kiro-cli fails. This could result in pushing incomplete or invalid changes.

**No retry logic:** If kiro-cli fails due to transient issues (network, rate limits), the iteration is lost. No automatic retry mechanism exists.

**No validation:** Script does not validate that kiro-cli supports required flags before invocation. Incompatible versions will fail at runtime.

**No timeout:** If kiro-cli hangs, the script waits indefinitely. No timeout mechanism exists.

## Areas for Improvement

**Dependency checking:** Add validation that kiro-cli is installed and supports required flags before starting loop.

**Error handling:** Check kiro-cli exit status and handle failures gracefully (retry, skip push, abort loop).

**Timeout mechanism:** Add timeout for AI CLI invocation to prevent indefinite hangs.

**Prompt size validation:** Check assembled prompt size before piping to kiro-cli, warn if approaching token limits.

**Alternative CLI support:** Document how to use other AI CLIs (Claude CLI, OpenAI CLI, etc.) with appropriate flag mappings.

**Version requirements:** Specify minimum kiro-cli version required for compatibility.
