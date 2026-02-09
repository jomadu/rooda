# Error Handling

## Job to be Done

Detect, report, and recover from failures — AI CLI crashes, network issues, test failures, invalid configs — with configurable retry logic, timeouts, and graceful degradation. The loop must distinguish transient failures (retry) from permanent failures (abort), provide clear diagnostics, and prevent silent corruption.

## Activities

1. Validate configuration at load time (fail fast before any iteration)
2. Detect AI CLI process failures (non-zero exit, crash, timeout)
3. Scan AI CLI output for `<promise>FAILURE</promise>` signals
4. Track consecutive failures across iterations
5. Abort loop when failure threshold exceeded
6. Handle partial output from crashed processes
7. Kill and cleanup AI CLI processes on timeout
8. Handle SIGINT/SIGTERM gracefully (cleanup and exit)
9. Log errors with context (iteration number, command, exit code)
10. Provide actionable error messages to user

## Acceptance Criteria

### Configuration Validation

- [ ] Config file syntax errors detected at load time (before any iteration)
- [ ] Missing required fields detected at load time
- [ ] Invalid field values detected at load time (e.g., negative timeouts)
- [ ] Missing prompt files detected at load time (fail fast)
- [ ] AI command binary validated at load time (exists and executable)
- [ ] Dry-run mode validates all above without executing AI CLI
- [ ] Dry-run mode exits with code 0 if all validations pass
- [ ] Dry-run mode exits with code 1 if any validation fails (user error, config error, missing files)
- [ ] Validation errors include file path, line number (if applicable), and clear fix guidance
- [ ] If config validation fails, exit with code 1 and error message

### AI CLI Failure Detection

- [ ] Promise signals scanned first, before checking exit code
- [ ] `<promise>FAILURE</promise>` in output counts as failure (overrides exit code 0)
- [ ] If both `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>` present, FAILURE takes precedence
- [ ] `<promise>SUCCESS</promise>` counts as success (overrides non-zero exit code)
- [ ] If no promise signals, non-zero exit code counts as failure
- [ ] If no promise signals and exit code 0, counts as success
- [ ] Process crash (SIGSEGV, SIGKILL, OOM) counts as failure
- [ ] Process timeout always counts as failure (promise signals in timed-out output ignored)
- [ ] Promise signals are case-sensitive: exact match for `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>`
- [ ] Invalid formats (lowercase, extra spaces, unclosed tags) are not recognized
- [ ] Failure increments `ConsecutiveFailures` counter
- [ ] Success resets `ConsecutiveFailures` to 0

### Consecutive Failure Tracking

- [ ] `ConsecutiveFailures` counter starts at 0
- [ ] Each failure increments counter
- [ ] Each success resets counter to 0
- [ ] Loop aborts when `ConsecutiveFailures >= FailureThreshold`
- [ ] Default `FailureThreshold` is 3
- [ ] `FailureThreshold` configurable via `loop.failure_threshold`
- [ ] `FailureThreshold` overridable per-procedure via `procedure.failure_threshold`
- [ ] Abort logs final failure count and threshold
- [ ] Abort exits with code 1 and status `aborted`

### Timeout Handling

- [ ] Iteration timeout configurable via `loop.iteration_timeout` (seconds, nil = no timeout)
- [ ] Iteration timeout overridable per-procedure via `procedure.iteration_timeout`
- [ ] `ROODA_LOOP_ITERATION_TIMEOUT` environment variable sets `loop.iteration_timeout`
- [ ] If AI CLI exceeds timeout, process killed with SIGTERM
- [ ] After SIGTERM, wait up to 5 seconds for graceful termination
- [ ] If process doesn't terminate within 5 seconds, send SIGKILL
- [ ] If process still doesn't terminate after SIGKILL, log warning and continue
- [ ] Timeout always counts as failure (increments `ConsecutiveFailures`)
- [ ] Timeout always fails even if output contains `<promise>SUCCESS</promise>`
- [ ] Promise signals in timed-out output logged for diagnostics but don't affect outcome
- [ ] Timeout logs iteration number, command, and timeout duration
- [ ] Partial output from timed-out process captured for diagnostic logging

### Crash Recovery

- [ ] Partial output from crashed AI CLI processes captured
- [ ] Partial output scanned for `<promise>` signals
- [ ] Crash logs process exit signal (SIGSEGV, SIGKILL, etc.)
- [ ] Crash counts as failure
- [ ] Crash doesn't corrupt loop state (next iteration starts clean)

### Signal Handling

- [ ] SIGINT (Ctrl+C) kills AI CLI process and exits with status `interrupted`
- [ ] SIGTERM kills AI CLI process and exits with status `interrupted`
- [ ] After SIGINT/SIGTERM, wait up to 5 seconds for AI CLI to terminate
- [ ] If AI CLI doesn't terminate within 5 seconds, send SIGKILL
- [ ] If AI CLI still doesn't terminate after SIGKILL, log warning and exit anyway
- [ ] Interrupted loops exit with code 130 (standard SIGINT exit code)
- [ ] Signal handling doesn't leave zombie processes

### Error Logging

- [ ] All errors logged with structured context (iteration, command, exit code, duration)
- [ ] Config validation errors include file path and line number
- [ ] AI CLI failures include command, exit code, and first/last 500 chars of output
- [ ] Timeout errors include timeout duration and command
- [ ] Crash errors include signal name and partial output
- [ ] Consecutive failure abort includes failure count and threshold
- [ ] Error messages actionable (suggest fixes where possible)

### Output Buffer Overflow

- [ ] Output buffer size configurable via `loop.max_output_buffer` (default: 10485760 bytes = 10MB)
- [ ] Output buffer size overridable per-procedure via `procedure.max_output_buffer`
- [ ] If output exceeds buffer size, truncate from beginning (keep most recent output)
- [ ] Truncation logs warning with actual size and buffer limit
- [ ] Truncation sets `AIExecutionResult.Truncated = true`
- [ ] `<promise>` signals at end of output preserved (signals at beginning may be lost)
- [ ] If promise signal is split across truncation boundary, it may not be detected
- [ ] Prompts must instruct AI to emit promise signals at the very end of output
- [ ] If truncation is frequent, increase `max_output_buffer` size

## Edge Cases

### Promise Signal Format

**Valid signals** (case-sensitive, exact match):
- `<promise>SUCCESS</promise>`
- `<promise>FAILURE</promise>`

**Invalid signals** (not recognized):
- `<promise>success</promise>` (lowercase)
- `<promise> SUCCESS </promise>` (extra spaces)
- `<PROMISE>SUCCESS</PROMISE>` (uppercase tags)
- `<promise>SUCCESS` (unclosed tag)
- `<Promise>Success</Promise>` (mixed case)

**Rationale:** Strict format forces AI to follow exact specification, prevents ambiguity, and enables fast string matching.

### Timeout with Promise Signals

**Scenario:** AI CLI times out but output contains `<promise>SUCCESS</promise>`

**Behavior:** Timeout always counts as failure. Promise signals in timed-out output are logged for diagnostics but do not affect outcome.

**Rationale:** If agent can't complete within timeout, it's a failure regardless of partial output. User should increase timeout if legitimate work takes longer.

### Truncated Output with Split Promise Signal

**Scenario:** Promise signal is split across truncation boundary

**Example:**
```
Output: [9.9MB of text]...<prom[TRUNCATE]ise>SUCCESS</promise>
Result: Buffer contains "ise>SUCCESS</promise>" (incomplete tag)
```

**Behavior:** Signal scanning will fail because `<promise>SUCCESS</promise>` is incomplete. Iteration outcome determined by exit code.

**Mitigation:** Prompts must instruct AI to emit promise signals at the very end of output. 10MB buffer is sufficient for most iterations.

### All Failures Treated Equally

**Current behavior:** All failures (timeout, crash, exit code, explicit FAILURE signal) increment the same `ConsecutiveFailures` counter.

**Limitation:** Cannot distinguish transient failures (network blip, rate limit) from permanent failures (invalid config, AI agent blocked).

**Rationale:** Simple implementation for v2. Future versions may add failure type classification with separate thresholds.

## Data Structures

### FailureContext

```go
type FailureContext struct {
    Iteration       int           // Which iteration failed
    Command         string        // AI command that was executed
    ExitCode        int           // Process exit code (if available)
    Signal          string        // Signal name if process crashed (e.g., "SIGSEGV")
    Duration        time.Duration // How long execution took before failure
    OutputPreview   string        // First 500 chars of output
    OutputSuffix    string        // Last 500 chars of output
    Truncated       bool          // True if output was truncated
    TimeoutDuration *int          // Timeout value if failure was due to timeout
}
```

**Fields:**
- `Iteration` — 0-indexed iteration number where failure occurred
- `Command` — Full AI command string that was executed
- `ExitCode` — Process exit code (0 if crashed before exit)
- `Signal` — Signal name if process was killed (SIGSEGV, SIGKILL, SIGTERM, etc.)
- `Duration` — Wall clock time from process start to failure
- `OutputPreview` — First 500 chars of captured output (for diagnostics)
- `OutputSuffix` — Last 500 chars of captured output (where `<promise>` signals likely are)
- `Truncated` — True if output exceeded buffer size
- `TimeoutDuration` — Timeout value in seconds if failure was due to timeout

### ValidationError

```go
type ValidationError struct {
    FilePath   string // Config or prompt file path
    LineNumber int    // Line number where error occurred (0 if not applicable)
    Field      string // Config field name (if applicable)
    Message    string // Human-readable error message
    Suggestion string // Suggested fix (if available)
}
```

**Fields:**
- `FilePath` — Path to file that failed validation
- `LineNumber` — Line number where error occurred (0 if not line-specific)
- `Field` — Config field name that failed validation (empty if not field-specific)
- `Message` — Clear description of what's wrong
- `Suggestion` — Actionable fix suggestion (e.g., "Set loop.ai_cmd or use --ai-cmd flag")

## Future Enhancements (Out of Scope for v2)

### Failure Type Classification

Distinguish transient failures (network, timeout, rate limit) from permanent failures (config error, explicit FAILURE signal) with separate thresholds:

```go
type FailureType string

const (
    FailureTransient FailureType = "transient" // Network, timeout, rate limit
    FailurePermanent FailureType = "permanent" // Config, explicit FAILURE signal
)

type IterationState struct {
    // ...
    TransientFailures int // Consecutive transient failures
    PermanentFailures int // Consecutive permanent failures
}

// Abort if: PermanentFailures >= 1 OR TransientFailures >= 5
```

**Tradeoffs:**
- ✅ More intelligent retry logic
- ✅ Don't abort on transient issues
- ⚠️ Complex: how to classify each failure?
- ⚠️ More configuration needed

### Exponential Backoff

Wait between iterations after transient failures:

```go
On timeout/network failure:
  - Wait before next iteration (1s, 2s, 4s, 8s...)
  - Don't increment counter if backoff succeeds
```

**Tradeoffs:**
- ✅ Handles transient issues gracefully
- ✅ Simple classification (timeout = transient)
- ⚠️ Adds delay to loop
- ⚠️ More complex state management

### Failure Pattern Detection

Abort immediately on repeated identical failures:

```go
If same error message 3 times in a row:
  - Abort immediately (don't wait for threshold)
  - Log: "Repeated identical failure detected"
```

**Tradeoffs:**
- ✅ Faster abort on permanent issues
- ✅ Prevents wasted iterations
- ⚠️ Requires error message comparison
- ⚠️ May abort too early on legitimate retries

## Algorithm

### Detect Iteration Failure

```
Input: AIExecutionResult
Output: bool (is failure)

1. Scan result.Output for <promise> signals (case-sensitive exact match):
   - If contains <promise>FAILURE</promise> → failure (overrides everything)
   - If contains <promise>SUCCESS</promise> (and no FAILURE) → success (overrides exit code)
2. If no promise signals:
   - If result.Error != nil → failure
   - If result.ExitCode != 0 → failure
   - Otherwise → success
```

**Precedence:**
1. `<promise>FAILURE</promise>` overrides everything (even exit code 0)
2. `<promise>SUCCESS</promise>` overrides non-zero exit code
3. Process errors (timeout, crash) count as failure if no promise signals
4. Exit code used if no promise signals and no process errors

**Implementation:**
```go
func DetectIterationFailure(result AIExecutionResult) bool {
    hasSuccess := strings.Contains(result.Output, "<promise>SUCCESS</promise>")
    hasFailure := strings.Contains(result.Output, "<promise>FAILURE</promise>")
    
    // Promise signals take precedence
    if hasFailure {
        return true  // Failure wins if both present
    }
    if hasSuccess {
        return false  // Success overrides exit code
    }
    
    // No promise signals - fall back to process status
    if result.Error != nil {
        return true
    }
    return result.ExitCode != 0
}
```

### Track Consecutive Failures

```
Input: IterationState, isFailure bool
Output: Updated IterationState, shouldAbort bool

1. If isFailure:
   - Increment state.ConsecutiveFailures
   - If state.ConsecutiveFailures >= state.FailureThreshold:
     - Log abort message with failure count
     - Return shouldAbort = true
2. If not isFailure:
   - Reset state.ConsecutiveFailures = 0
3. Return shouldAbort = false
```

### Handle Process Timeout

```
Input: Process, timeout duration
Output: AIExecutionResult

1. Start process with assembled prompt as stdin
2. Start timer for timeout duration
3. Wait for process exit or timeout:
   - If process exits before timeout → capture output and exit code
   - If timeout expires:
     a. Send SIGTERM to process
     b. Wait up to 5 seconds for graceful termination
     c. If still running, send SIGKILL
     d. Wait up to 1 second for forced termination
     e. If still running, log warning and continue
     f. Capture partial output
     g. Return result with Error = "timeout exceeded"
4. Return result (timeout always counts as failure, promise signals ignored)
```

**Note:** Partial output from timed-out processes is logged for diagnostics but does not affect failure detection. Timeout always increments `ConsecutiveFailures` regardless of output content.ecutionResult
```

### Handle SIGINT/SIGTERM

```
Input: Signal (SIGINT or SIGTERM)
Output: Exit with code 130

1. Log "Received signal, shutting down..."
2. If AI CLI process running:
   a. Send SIGTERM to AI CLI process
   b. Wait up to 5 seconds for graceful termination
   c. If still running, send SIGKILL
   d. Wait up to 1 second for forced termination
   e. If still running, log warning
3. Set loop status = "interrupted"
4. Exit with code 130
```

## Dependencies

- **iteration-loop** — Consumes failure detection and consecutive failure tracking
- **ai-cli-integration** — Provides process execution and output capture
- **configuration** — Provides failure threshold and timeout settings
- **observability** — Logs errors with structured context

## Implementation Mapping

**Source files:**
- `internal/loop/errors.go` — Failure detection, consecutive failure tracking
- `internal/loop/signals.go` — SIGINT/SIGTERM handling
- `internal/ai/executor.go` — Process timeout, crash recovery
- `internal/config/validation.go` — Config validation at load time

**Related specs:**
- `iteration-loop.md` — Defines loop termination conditions and failure threshold
- `ai-cli-integration.md` — Defines process execution and output capture
- `configuration.md` — Defines config schema and validation rules
- `observability.md` — Defines logging format and levels

## Examples

### Example 1: Consecutive Failures Leading to Abort

**Input:**
- `loop.failure_threshold = 3`
- Iteration 0: AI CLI exits with code 1
- Iteration 1: AI CLI exits with code 1
- Iteration 2: AI CLI exits with code 1

**Output:**
```
[21:00:00.000] INFO Starting iteration 1/5 procedure=build
[21:00:05.100] ERROR AI CLI failed exit_code=1 consecutive_failures=1 threshold=3
[21:00:05.200] INFO Starting iteration 2/5 procedure=build
[21:00:10.300] ERROR AI CLI failed exit_code=1 consecutive_failures=2 threshold=3
[21:00:10.400] INFO Starting iteration 3/5 procedure=build
[21:00:15.500] ERROR AI CLI failed exit_code=1 consecutive_failures=3 threshold=3
[21:00:15.600] ERROR Loop aborted consecutive_failures=3 threshold=3
Exit code: 1
```

### Example 2: Failure Followed by Success (Counter Reset)

**Input:**
- `loop.failure_threshold = 3`
- Iteration 0: AI CLI exits with code 1
- Iteration 1: AI CLI exits with code 0, output contains `<promise>SUCCESS</promise>`

**Output:**
```
[21:00:00.000] INFO Starting iteration 1/5 procedure=build
[21:00:05.100] ERROR AI CLI failed exit_code=1 consecutive_failures=1 threshold=3
[21:00:05.200] INFO Starting iteration 2/5 procedure=build
[21:00:20.300] INFO AI CLI succeeded promise_signal=SUCCESS
[21:00:20.400] INFO Loop completed status=success iterations=2 total_elapsed=20.4s
Exit code: 0
```

### Example 3: Promise Signal Overrides Exit Code

**Input:**
- Iteration 0: AI CLI exits with code 1, output contains `<promise>SUCCESS</promise>`

**Output:**
```
[21:00:00.000] INFO Starting iteration 1/5 procedure=build
[21:00:15.100] WARN AI CLI exited with non-zero code but promise signal indicates success exit_code=1 promise_signal=SUCCESS
[21:00:15.200] INFO Loop completed status=success iterations=1 total_elapsed=15.2s
Exit code: 0
```

### Example 4: Timeout Kills Process

**Input:**
- `loop.iteration_timeout = 60` (seconds)
- AI CLI runs for 65 seconds without completing

**Output:**
```
[21:00:00.000] INFO Starting iteration 1/5 procedure=build
[21:01:00.100] WARN AI CLI exceeded timeout timeout=60s action=sending_sigterm
[21:01:05.200] WARN AI CLI did not terminate gracefully action=sending_sigkill
[21:01:05.300] ERROR AI CLI failed reason=timeout timeout=60s consecutive_failures=1 threshold=3
[21:01:05.400] INFO Starting iteration 2/5 procedure=build
...
```

### Example 5: Config Validation Failure

**Input:**
- `rooda-config.yml` contains `loop.iteration_timeout: -10`

**Output:**
```
[21:00:00.000] ERROR Config validation failed file=rooda-config.yml line=12 field=loop.iteration_timeout error="timeout must be positive or nil" suggestion="Set to a positive number of seconds, or remove to disable timeout"
Exit code: 1
```

### Example 6: Missing AI Command

**Input:**
- No `--ai-cmd` flag
- No `--ai-cmd-alias` flag
- No `loop.ai_cmd` in config
- No `loop.ai_cmd_alias` in config

**Output:**
```
[21:00:00.000] ERROR No AI command configured error="must specify AI command via CLI flag (--ai-cmd or --ai-cmd-alias), config file (loop.ai_cmd or loop.ai_cmd_alias), or built-in alias (kiro-cli, claude, copilot, cursor-agent)" suggestion="Example: rooda build --ai-cmd-alias kiro-cli"
Exit code: 1
```

### Example 7: Output Buffer Overflow

**Input:**
- `loop.max_output_buffer = 1048576` (1MB)
- AI CLI produces 5MB of output

**Output:**
```
[21:00:00.000] INFO Starting iteration 1/5 procedure=build
[21:00:15.100] WARN AI CLI output exceeded buffer size actual_size=5242880 buffer_limit=1048576 action=truncating_from_beginning
[21:00:15.200] INFO Completed iteration 1/5 elapsed=15.2s status=success truncated=true
```

**Result:**
- `AIExecutionResult.Output` contains last ~1MB of output
- `AIExecutionResult.Truncated = true`
- `<promise>` signals at end of output preserved

### Example 8: SIGINT During Iteration

**Input:**
- User presses Ctrl+C during iteration 2

**Output:**
```
[21:00:00.000] INFO Starting iteration 2/5 procedure=build
[21:00:05.100] INFO Received signal signal=SIGINT action=shutting_down
[21:00:05.200] INFO Killing AI CLI process pid=12345
[21:00:05.300] INFO AI CLI terminated
[21:00:05.400] INFO Loop interrupted status=interrupted
Exit code: 130
```

## Rationale

### Why Consecutive Failure Threshold?

A single failure might be transient (network blip, rate limit). But 3+ consecutive failures indicate a systemic problem (bad config, broken code, impossible task). Aborting after threshold prevents infinite loops burning API credits.

### Why Promise Signals Override Exit Code?

AI CLI tools have inconsistent exit code semantics. Some return 0 even when the agent failed its task. Promise signals are explicit, agent-controlled success/failure indicators that are more reliable than exit codes.

### Why FAILURE Takes Precedence Over SUCCESS?

Conservative choice. If the agent emits both signals (bug, confusion, partial success), treat as failure to avoid false positives. Better to retry than to silently accept broken output.

### Why Truncate Output from Beginning?

`<promise>` signals are emitted at the end of output (after all work complete). Keeping the end of output ensures signals are preserved. Losing the beginning is acceptable — it's usually verbose tool output, not the final result.

### Why 5 Second Termination Timeout?

AI CLI tools should respond to SIGTERM quickly (flush buffers, cleanup). 5 seconds is generous for graceful shutdown. If they don't respond, SIGKILL forces termination. This prevents hung processes from blocking the loop.

### Why Validate Config at Load Time?

Fail fast. If config is broken, no point starting iterations. Validation at load time provides immediate feedback and prevents wasted API calls.

### Why Track Output Truncation?

Diagnostics. If output is truncated, the user needs to know — it might explain why signal scanning failed or why error messages are incomplete. Logging truncation and setting `Truncated = true` makes this visible.
