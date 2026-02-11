# AI CLI Integration

## Job to be Done

Pipe assembled prompts to a configurable AI CLI tool with support for command aliases, environment variables, and direct command override. Built-in support for kiro-cli, claude, github copilot, and cursor agent, with extensibility for custom tools.

The developer wants to use their preferred AI CLI tool without modifying framework code, switch between tools for different procedures (fast model for audits, thorough model for builds), and have the loop capture output for signal scanning while optionally streaming it to the terminal for visibility.

## Activities

1. Resolve AI command from precedence chain (CLI flags > procedure config > loop config > error)
2. Validate AI command binary exists and is executable (dry-run mode)
3. Spawn AI CLI process with assembled prompt as stdin
4. Capture stdout/stderr to buffer (with configurable max size)
5. Optionally stream output to terminal in real-time (if `--verbose`)
6. Wait for process termination or timeout
7. Scan captured output for `<promise>` signals
8. Return exit code, captured output, and any errors

## Acceptance Criteria

- [ ] AI command resolved from precedence chain: `--ai-cmd` > `--ai-cmd-alias` > procedure `ai_cmd` > procedure `ai_cmd_alias` > `loop.ai_cmd` > `loop.ai_cmd_alias` > error
- [ ] If no AI command configured via any source, error with clear guidance listing all ways to set one
- [ ] Built-in aliases available: `kiro-cli`, `claude`, `copilot`, `cursor-agent`
- [ ] Built-in alias `kiro-cli` maps to: `kiro-cli chat --no-interactive --trust-all-tools`
- [ ] Built-in alias `claude` maps to: `claude -p --dangerously-skip-permissions`
- [ ] Built-in alias `copilot` maps to: `copilot --yolo`
- [ ] Built-in alias `cursor-agent` maps to: `cursor-wrapper.sh` (wrapper script that parses JSON output)
- [ ] User-defined aliases in config merge with built-in aliases (user aliases can override built-in)
- [ ] Direct command (`ai_cmd`) takes precedence over alias (`ai_cmd_alias`) at same config level
- [ ] Dry-run mode displays resolved config and assembled prompt (binary validation already done at config load)
- [ ] Assembled prompt piped to AI CLI stdin
- [ ] AI CLI stdout and stderr captured to buffer
- [ ] Output buffer size configurable via `loop.max_output_buffer` (default: 10485760 bytes = 10MB)
- [ ] Output buffer size overridable per-procedure via `procedure.max_output_buffer`
- [ ] AI command binary validated at config load time (fail fast before any iteration)
- [ ] If output exceeds buffer size, buffer truncated from beginning (keeps most recent output), warning logged
- [ ] Output always captured regardless of `--verbose` flag (needed for signal scanning)
- [ ] When `--verbose` flag set, output streamed to terminal in real-time while also being captured
- [ ] When `--verbose` not set, output captured silently (only loop-level progress displayed)
- [ ] Process exit code captured and returned
- [ ] If process exceeds AI execution timeout, process killed and error returned
- [ ] If process crashes (SIGSEGV, SIGKILL, OOM), partial output captured and returned
- [ ] Captured output scanned for `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>` signals using substring matching
- [ ] Signal scanning happens after process exits (not during streaming)
- [ ] Signals detected anywhere in output (substring match), but should be on own line for clarity
- [ ] Prompts should instruct AI to emit `<promise>` signals at end of output (after all work complete)
- [ ] If output is truncated, signals at end are preserved (signals at beginning may be lost)
- [ ] Environment variables from parent process inherited by AI CLI process
- [ ] Working directory for AI CLI process is current working directory (where rooda was invoked)
- [ ] SIGINT/SIGTERM to rooda kills AI CLI process and waits for termination (with 5s timeout)
- [ ] If AI CLI doesn't terminate within timeout after kill signal, log warning and exit anyway

## Data Structures

### AICommand

```go
type AICommand struct {
    Command string   // Resolved command string (e.g., "kiro-cli chat --no-interactive --trust-all-tools")
    Source  string   // Provenance: where this command came from (e.g., "--ai-cmd flag", "procedure.ai_cmd", "loop.ai_cmd_alias=kiro-cli")
}
```

**Fields:**
- `Command` — Full command string to execute, including arguments
- `Source` — Provenance for debugging and dry-run display

### AIExecutionResult

```go
type AIExecutionResult struct {
    Output       string        // Captured stdout/stderr
    ExitCode     int           // Process exit code
    Duration     time.Duration // How long execution took
    Truncated    bool          // True if output was truncated due to buffer size
    Error        error         // Non-nil if execution failed (timeout, spawn failure, etc.)
}
```

**Fields:**
- `Output` — Combined stdout/stderr, truncated from beginning if exceeds max buffer size
- `ExitCode` — Process exit code (0 = success, non-zero = failure). Only valid if Error is nil.
- `Duration` — Wall clock time from process start to termination
- `Truncated` — True if output exceeded max buffer size and was truncated
- `Error` — Non-nil if process couldn't be spawned, timed out, or other execution failure. If Error is non-nil, ExitCode may not be meaningful.

### Built-in Aliases

```go
var BuiltinAliases = map[string]string{
    "kiro-cli":     "kiro-cli chat --no-interactive --trust-all-tools",
    "claude":       "claude -p --dangerously-skip-permissions",
    "copilot":      "copilot --yolo",
    "cursor-agent": "cursor-wrapper.sh",
}
```

## Algorithm

### Resolve AI Command

This is defined in `configuration.md` under "AI Command Resolution". Summary:

```
function ResolveAICommand(config Config, procedureName string, cliFlags CLIFlags) -> (AICommand, error):
    // 1. --ai-cmd flag (direct command, highest precedence)
    if cliFlags.AICmd != "":
        return AICommand{cliFlags.AICmd, "--ai-cmd flag"}, nil
    
    // 2. --ai-cmd-alias flag (alias from merged config)
    if cliFlags.AICmdAlias != "":
        return resolveAlias(config, cliFlags.AICmdAlias, "--ai-cmd-alias flag")
    
    // 3. procedure.ai_cmd (direct command)
    proc = config.Procedures[procedureName]
    if proc.AICmd != "":
        return AICommand{proc.AICmd, "procedure." + procedureName + ".ai_cmd"}, nil
    
    // 4. procedure.ai_cmd_alias (alias from merged config)
    if proc.AICmdAlias != "":
        return resolveAlias(config, proc.AICmdAlias, "procedure." + procedureName + ".ai_cmd_alias")
    
    // 5. loop.ai_cmd (already merged from config tiers + env vars)
    if config.Loop.AICmd != "":
        return AICommand{config.Loop.AICmd, "loop.ai_cmd"}, nil
    
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
        "  procedure.ai_cmd or ai_cmd_alias  rooda-config.yml\n\n" +
        "Available aliases: %v", keys(config.AICmdAliases))

function resolveAlias(config Config, aliasName string, source string) -> (AICommand, error):
    command, exists = config.AICmdAliases[aliasName]
    if !exists:
        return error("unknown AI command alias: %s (from %s)\nAvailable: %v", 
            aliasName, source, keys(config.AICmdAliases))
    return AICommand{command, source + "=" + aliasName}, nil
```

### Execute AI CLI

```go
function ExecuteAICLI(aiCmd AICommand, prompt string, verbose bool, aiExecutionTimeout *int, maxBuffer int) -> AIExecutionResult:
    startTime = time.Now()
    
    // Parse command string into binary and args using shell-style quoting
    parts, err = shellquote.Split(aiCmd.Command)
    if len(parts) == 0:
        return AIExecutionResult{Error: error("empty AI command")}
    
    binary = parts[0]
    args = parts[1:]
    
    // Create command
    cmd = exec.Command(binary, args...)
    cmd.Dir = os.Getwd()  // Inherit current working directory
    cmd.Env = os.Environ()  // Inherit environment variables
    
    // Set up stdin with prompt
    cmd.Stdin = strings.NewReader(prompt)
    
    // Set up output capture
    var outputBuffer bytes.Buffer
    var outputWriter io.Writer = &outputBuffer
    
    // If verbose, also stream to terminal
    if verbose:
        outputWriter = io.MultiWriter(&outputBuffer, os.Stdout)
    
    cmd.Stdout = outputWriter
    cmd.Stderr = outputWriter
    
    // Start process
    if err = cmd.Start(); err != nil:
        return AIExecutionResult{
            Error: error("failed to start AI CLI: %w", err),
            Duration: time.Since(startTime),
        }
    
    // Wait for completion or timeout
    done = make(chan error)
    go func():
        done <- cmd.Wait()
    
    var waitErr error
    if aiExecutionTimeout != nil:
        select:
            case waitErr = <-done:
                // Process completed
            case <-time.After(time.Duration(*aiExecutionTimeout) * time.Second):
                // Timeout — kill process
                cmd.Process.Kill()
                <-done  // Wait for process to actually terminate
                return AIExecutionResult{
                    Output: outputBuffer.String(),
                    Duration: time.Since(startTime),
                    Error: ErrTimeout,
                }
    else:
        waitErr = <-done
    
    duration = time.Since(startTime)
    output = outputBuffer.String()
    truncated = false
    
    // Truncate output if exceeds max buffer
    if len(output) > maxBuffer:
        truncated = true
        output = output[len(output)-maxBuffer:]  // Keep most recent output
    
    // Get exit code
    exitCode = 0
    if waitErr != nil:
        if exitError, ok = waitErr.(*exec.ExitError); ok:
            exitCode = exitError.ExitCode()
        else:
            // Process failed to start or other error
            return AIExecutionResult{
                Output: output,
                Duration: duration,
                Truncated: truncated,
                Error: error("process execution failed: %w", waitErr),
            }
    
    return AIExecutionResult{
        Output: output,
        ExitCode: exitCode,
        Duration: duration,
        Truncated: truncated,
        Error: nil,
    }
```

### Validate AI Command (Config Load Time)

```go
function ValidateAICommand(aiCmd AICommand) -> error:
    // Parse command string using shell-style quoting
    parts, err = shellquote.Split(aiCmd.Command)
    if err != nil:
        return error("invalid AI command syntax: %w\nCommand: %s", err, aiCmd.Command)
    if len(parts) == 0:
        return error("empty AI command")
    
    binary = parts[0]
    
    // Check if binary exists and is executable
    path, err = exec.LookPath(binary)
    if err != nil:
        return error("AI command binary not found: %s\nSource: %s\nCommand: %s", 
            binary, aiCmd.Source, aiCmd.Command)
    
    // Check if file is executable
    info, err = os.Stat(path)
    if err != nil:
        return error("cannot stat AI command binary: %s", path)
    
    if info.Mode() & 0111 == 0:
        return error("AI command binary is not executable: %s", path)
    
    return nil
```

### Scan Output for Signals

```go
function ScanOutputForSignals(output string) -> (hasSuccess bool, hasFailure bool):
    hasSuccess = strings.Contains(output, "<promise>SUCCESS</promise>")
    hasFailure = strings.Contains(output, "<promise>FAILURE</promise>")
    return hasSuccess, hasFailure
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| AI command binary not in PATH | Error during validation (dry-run) or execution: "AI command binary not found: X" |
| AI command binary not executable | Error during validation (dry-run): "AI command binary is not executable: X" |
| AI command is empty string | Error: "empty AI command" |
| AI command has only whitespace | Parsed as empty, error: "empty AI command" |
| AI command with quoted arguments | Parsed correctly using shell-style quoting: `claude -p --model "claude-3-5-sonnet"` → binary=`claude`, args=`["-p", "--model", "claude-3-5-sonnet"]` |
| AI command with complex shell syntax (pipes, redirects) | Not supported — no shell execution. User must wrap in a script. Example: create `my-ai.sh` with `#!/bin/bash\nclaude -p | tee log.txt`, then use `--ai-cmd my-ai.sh` |
| Output exceeds max buffer during execution | Buffer truncated from beginning after process exits, warning logged, `Truncated=true` in result |
| Process times out | Process killed, partial output captured, `Error=ErrTimeout` returned |
| Process crashes (SIGSEGV, SIGKILL) | Partial output captured, exit code captured, `Error=nil` (crash is not an execution error, just a non-zero exit) |
| Process writes to stdout and stderr | Both captured to same buffer (combined output) |
| Verbose mode enabled | Output streamed to terminal in real-time AND captured to buffer |
| Verbose mode disabled | Output captured silently, only loop-level progress displayed |
| SIGINT sent to rooda during AI CLI execution | AI CLI process killed, rooda waits up to 5s for termination, then exits with code 130 |
| AI CLI doesn't terminate after kill signal | After 5s timeout, log warning and exit anyway (don't hang forever) |
| Alias references non-existent alias | Error at resolution time: "unknown AI command alias: X (from Y)\nAvailable: [...]" |
| Both `ai_cmd` and `ai_cmd_alias` set at same level | `ai_cmd` wins (direct command overrides alias) |
| User-defined alias overrides built-in alias | User alias wins (config merging allows override) |
| Prompt is empty string | Valid — empty stdin piped to AI CLI (AI CLI may error, but that's its responsibility) |
| AI CLI requires interactive input | Not supported — all AI CLIs must support non-interactive mode and be pre-configured (API keys, etc.). Built-in aliases include non-interactive flags. |
| AI CLI writes binary data to stdout | Captured as-is, signal scanning may fail if binary data corrupts string matching |
| Output contains multiple `<promise>SUCCESS</promise>` | Treated same as one — `hasSuccess=true` |
| Output contains both SUCCESS and FAILURE signals | Both flags returned; iteration loop decides precedence (FAILURE wins) |

## Dependencies

- **Go standard library** — `os/exec` for process spawning, `io` for stream handling, `bytes` for buffering
- **github.com/kballard/go-shellquote** — Shell-style command string parsing (handles quoted arguments)
- **jq** — Required for cursor-agent wrapper script (JSON parsing)
- **configuration** — Provides resolved AI command via `ResolveAICommand`, validates at config load time
- **iteration-loop** — Calls `ExecuteAICLI` each iteration, interprets exit code and signals
- **observability** — Logs execution start, completion, timeout, truncation warnings

## Implementation Mapping

**Source files:**
- `internal/ai/executor.go` — `ExecuteAICLI`, `ValidateAICommand`, `ScanOutputForSignals`
- `internal/ai/resolver.go` — `ResolveAICommand`, `resolveAlias` (or in `internal/config/resolver.go`)
- `internal/ai/aliases.go` — `BuiltinAliases` constant

**Related specs:**
- `configuration.md` — Defines AI command resolution precedence and alias merging
- `iteration-loop.md` — Consumes `ExecuteAICLI`, interprets exit codes and signals
- `prompt-composition.md` — Produces assembled prompt that is piped to AI CLI
- `error-handling.md` — Defines timeout handling, retry logic, failure detection
- `observability.md` — Defines logging for AI CLI execution events

## Examples

### Example 1: Execute with Built-in Alias

**Input:**
```go
config = Config{
    Loop: LoopConfig{
        AICmdAlias: "kiro-cli",
        MaxOutputBuffer: 10485760,
    },
    AICmdAliases: BuiltinAliases,
}
aiCmd, _ = ResolveAICommand(config, "build", CLIFlags{})
prompt = "# Observe\n...\n# Act\n..."
result = ExecuteAICLI(aiCmd, prompt, false, nil, 10485760)
```

**Expected Output:**
```go
aiCmd = AICommand{
    Command: "kiro-cli chat --no-interactive --trust-all-tools",
    Source: "loop.ai_cmd_alias=kiro-cli",
}
result = AIExecutionResult{
    Output: "[AI CLI output here]",
    ExitCode: 0,
    Duration: 45.3 * time.Second,
    Truncated: false,
    Error: nil,
}
```

**Verification:**
- Alias resolved to built-in command
- Process spawned with correct command
- Output captured
- Exit code captured

### Example 2: Execute with Direct Command Override

**Input:**
```go
config = Config{
    Loop: LoopConfig{
        AICmdAlias: "kiro-cli",
    },
    AICmdAliases: BuiltinAliases,
}
cliFlags = CLIFlags{
    AICmd: "custom-ai-tool --flag",
}
aiCmd, _ = ResolveAICommand(config, "build", cliFlags)
result = ExecuteAICLI(aiCmd, prompt, false, nil, 10485760)
```

**Expected Output:**
```go
aiCmd = AICommand{
    Command: "custom-ai-tool --flag",
    Source: "--ai-cmd flag",
}
// Process spawns custom-ai-tool with --flag argument
```

**Verification:**
- CLI flag overrides config alias
- Custom command executed

### Example 3: Verbose Mode Streaming

**Input:**
```go
result = ExecuteAICLI(aiCmd, prompt, true, nil, 10485760)
```

**Expected Behavior:**
- AI CLI output appears on terminal in real-time as it's generated
- Output also captured to buffer for signal scanning
- User sees progress without waiting for iteration to complete

**Verification:**
- Output streamed to stdout during execution
- Output also available in `result.Output` after completion

### Example 4: Timeout Handling

**Input:**
```go
timeout = 60  // 60 seconds
result = ExecuteAICLI(aiCmd, prompt, false, &timeout, 10485760)
// AI CLI runs for 65 seconds
```

**Expected Output:**
```go
result = AIExecutionResult{
    Output: "[partial output before timeout]",
    Duration: ~60 * time.Second,
    Truncated: false,
    Error: ErrTimeout,
}
```

**Verification:**
- Process killed after 60 seconds
- Partial output captured
- Error indicates timeout

### Example 5: Output Buffer Truncation

**Input:**
```go
maxBuffer = 1024  // 1KB buffer
result = ExecuteAICLI(aiCmd, prompt, false, nil, maxBuffer)
// AI CLI produces 10KB of output
```

**Expected Output:**
```go
result = AIExecutionResult{
    Output: "[last 1KB of output]",
    ExitCode: 0,
    Duration: 45 * time.Second,
    Truncated: true,
    Error: nil,
}
```

**Verification:**
- Output truncated from beginning
- Most recent 1KB kept (for signal scanning)
- `Truncated=true` flag set
- Warning logged by caller

### Example 6: No AI Command Configured

**Input:**
```go
config = Config{
    Loop: LoopConfig{},
    AICmdAliases: BuiltinAliases,
}
aiCmd, err = ResolveAICommand(config, "build", CLIFlags{})
```

**Expected Output:**
```go
err = error("no AI command configured\n\nSet one via:\n" +
    "  --ai-cmd \"your-command\"           CLI flag (direct command)\n" +
    "  --ai-cmd-alias <name>             CLI flag (alias from config)\n" +
    "  ROODA_LOOP_AI_CMD=your-command    Environment variable\n" +
    "  ROODA_LOOP_AI_CMD_ALIAS=<name>    Environment variable\n" +
    "  loop.ai_cmd or loop.ai_cmd_alias  rooda-config.yml\n" +
    "  procedure.ai_cmd or ai_cmd_alias  rooda-config.yml\n\n" +
    "Available aliases: [kiro-cli claude copilot cursor-agent]")
```

**Verification:**
- Clear error message
- Lists all ways to configure AI command
- Shows available aliases

### Example 7: Dry-Run Validation

**Input:**
```go
aiCmd = AICommand{
    Command: "nonexistent-binary --flag",
    Source: "--ai-cmd flag",
}
err = ValidateAICommand(aiCmd)
```

**Expected Output:**
```go
err = error("AI command binary not found: nonexistent-binary\n" +
    "Source: --ai-cmd flag\n" +
    "Command: nonexistent-binary --flag")
```

**Verification:**
- Binary existence checked
- Clear error with provenance

### Example 8: Signal Scanning

**Input:**
```go
output1 = "Some output\n<promise>SUCCESS</promise>\nMore output"
output2 = "Some output\n<promise>FAILURE</promise>\nMore output"
output3 = "Some output\n<promise>SUCCESS</promise>\n<promise>FAILURE</promise>\n"
output4 = "No signals here"

hasSuccess1, hasFailure1 = ScanOutputForSignals(output1)
hasSuccess2, hasFailure2 = ScanOutputForSignals(output2)
hasSuccess3, hasFailure3 = ScanOutputForSignals(output3)
hasSuccess4, hasFailure4 = ScanOutputForSignals(output4)
```

**Expected Output:**
```go
hasSuccess1 = true,  hasFailure1 = false
hasSuccess2 = false, hasFailure2 = true
hasSuccess3 = true,  hasFailure3 = true   // Both present
hasSuccess4 = false, hasFailure4 = false
```

**Verification:**
- Simple string matching
- Both signals can be present simultaneously
- Caller (iteration loop) decides precedence

## Notes

**Design Rationale — Command String Parsing:**

AI commands are parsed using shell-style quoting (via `github.com/kballard/go-shellquote`) but NOT executed through a shell. This allows quoted arguments like `--model "claude-3-5-sonnet"` to work correctly while avoiding shell injection vulnerabilities. Complex shell features (pipes, redirects, environment variable expansion) are not supported — users should wrap such commands in a script.

**Design Rationale — No Shell Interpretation:**

The AI command is split on whitespace and executed directly via `exec.Command`, not through a shell. This avoids shell injection vulnerabilities and makes behavior predictable across platforms. If users need complex shell syntax (pipes, redirects, environment variable expansion), they should wrap their command in a script and invoke the script.

**Design Rationale — Combined stdout/stderr:**

Both stdout and stderr are captured to the same buffer. This simplifies signal scanning (signals can appear on either stream) and matches how most AI CLI tools emit output (mixed informational and response content). If separate streams become necessary, the implementation can be extended.

**Design Rationale — Truncation from Beginning:**

When output exceeds the max buffer size, we truncate from the beginning (keeping the most recent output). This ensures `<promise>` signals at the end of execution are preserved for scanning. The trade-off is losing early output, but that's acceptable — the signals are what matter for loop control.

**Design Rationale — Built-in Aliases:**

The four built-in aliases (`kiro-cli`, `claude`, `copilot`, `cursor-agent`) cover the most common AI CLI tools as of 2026. These are hardcoded in the binary for zero-config startup. Users can override these or add custom aliases via config. The alias names are intentionally short and memorable.

**Design Rationale — No Interactive Mode Support:**

All AI CLI tools must support non-interactive mode. The loop pipes a prompt to stdin and expects the AI to process it and exit. Interactive prompts, confirmation dialogs, or TTY requirements break this model. Built-in aliases include `--no-interactive` flags to enforce this.

**Design Rationale — AI Execution Timeout:**

Timeouts are configured as `ai_execution_timeout` (via `loop.iteration_timeout` or `procedure.iteration_timeout` in config) and apply only to AI CLI execution, not the full iteration. Prompt assembly and signal scanning typically take <1s, so actual iteration duration may exceed the configured timeout by 1-2s. This is acceptable — the timeout prevents runaway AI execution, which is the primary concern.

**Design Rationale — Timeout at Iteration Level:**

Timeouts are configured per-iteration (via `loop.iteration_timeout` or `procedure.iteration_timeout`), not per-procedure or per-loop. This gives fine-grained control — a single runaway iteration is killed, but the loop continues. If the timeout is too aggressive, users can increase it or disable it (nil = no timeout).

**Design Rationale — Environment Variable Inheritance:**

The AI CLI process inherits all environment variables from the parent rooda process. This allows users to configure AI CLI tools via environment variables (e.g., `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`) without rooda needing to know about them. The AI CLI tool is responsible for reading its own config from the environment.

**Design Rationale — Working Directory Inheritance:**

The AI CLI process runs in the same working directory where rooda was invoked. This ensures file paths in prompts (e.g., "read AGENTS.md") resolve correctly. The AI CLI tool can use relative paths naturally.

**Design Rationale — Signal Placement:**

Prompts should instruct the AI to emit `<promise>` signals at the END of output, after all work is complete. This ensures signals are preserved even if output is truncated (truncation keeps the most recent output). If signals appear early and output exceeds the buffer, they may be lost. The 10MB default buffer is sufficient for most iterations — if truncation occurs frequently, users should increase `max_output_buffer`.

Signals are detected using simple substring matching (`strings.Contains()`), so they can technically appear anywhere in the output, even with surrounding text. However, for clarity and reliability, signals should appear on their own line:

**Recommended:**
```
<promise>SUCCESS</promise>
```

**Works but discouraged:**
```
Task complete <promise>SUCCESS</promise> - all tests passing
```

The substring matching approach is intentionally simple and permissive, but prompts should guide agents to use the cleaner format.

**Design Rationale — Cursor Agent Wrapper:**

The cursor agent outputs JSON (`--output-format stream-json`) which requires parsing to extract text and emit `<promise>` signals. Rather than adding JSON parsing to the core loop, we provide a wrapper script (`cursor-wrapper.sh`) that handles this. The script requires `jq` as a dependency. Users who don't use cursor agent don't pay the complexity cost.

**Why Not Retry Logic Here:**

Retry logic for transient failures (network issues, rate limits) belongs in `error-handling.md`, not here. This module is responsible for executing the AI CLI once and reporting the result. The iteration loop decides whether to retry based on failure patterns.

**Why Signal Scanning After Exit:**

Signal scanning happens after the process exits, not during streaming. This simplifies implementation (no need for concurrent scanning) and avoids race conditions (signal might appear mid-stream). The trade-off is we can't terminate early when SUCCESS appears, but that's acceptable — most iterations complete quickly anyway.
