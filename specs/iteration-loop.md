# Iteration Loop

## Job to be Done

Execute AI coding agents through controlled OODA iteration cycles that clear AI context between runs, preventing LLM degradation while maintaining file-based state continuity. Each iteration invokes the AI CLI as a fresh process — the agent starts clean, processes the assembled prompt, executes tools, then exits. The loop orchestrator persists across iterations, managing termination and state.

This is the core loop — everything else feeds into or out of it.

## Activities

1. Initialize iteration state (counter, config, termination criteria)
2. Check termination conditions before each iteration
3. Assemble the OODA prompt from four phase files and optional context
4. Pipe assembled prompt to AI CLI tool, capturing output (stream to terminal if `--verbose`)
5. Scan captured output for `<promise>` signals
6. Capture AI CLI exit status and detect failures
7. Increment iteration counter and record timing
8. Display iteration progress
9. Repeat until termination condition met

Note: Output is always captured and scanned for `<promise>` signals regardless of `--verbose`. The flag only controls whether output is also displayed to the terminal.

Note: Quality gates (tests, lints) and git operations (commit, push) are the AI agent's responsibility within each iteration, driven by the prompts — not orchestrated by the loop.

## Acceptance Criteria

- [ ] Loop executes until max iterations reached, Ctrl+C pressed, consecutive failure threshold hit, or AI signals success
- [ ] If AI CLI output contains `<promise>SUCCESS</promise>`, loop terminates with status `success` regardless of exit code
- [ ] If AI CLI output contains `<promise>FAILURE</promise>`, iteration counts as a failure (increments `ConsecutiveFailures`) even if exit code is 0
- [ ] If both `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>` are present, FAILURE takes precedence (conservative choice)
- [ ] Each iteration invokes the AI CLI as a separate process, ensuring fresh context
- [ ] Iteration counter increments correctly (0-indexed internal, 1-indexed display)
- [ ] Iteration counting example: `--max-iterations 5` runs iterations 0-4 (displayed as 1-5), termination check `Iteration >= 5` prevents iteration 5 from starting
- [ ] Max iterations resolved with precedence: `--max-iterations` CLI flag > `--unlimited` CLI flag > procedure `iteration_mode`/`default_max_iterations` > `loop.iteration_mode`/`loop.default_max_iterations` > built-in default (mode=max-iterations, count=5)
- [ ] `--unlimited` flag overrides max iterations to nil (no limit)
- [ ] Procedure-level `iteration_mode` overrides global `loop.iteration_mode`
- [ ] Procedure-level `default_max_iterations` overrides global `loop.default_max_iterations`
- [ ] Progress displayed between iterations (iteration number, elapsed time)
- [ ] AI CLI exit status is captured and checked after each iteration
- [ ] Iteration outcome determined by output signals and exit code per the outcome matrix
- [ ] Consecutive failure count tracked; loop aborts after configurable threshold (`loop.failure_threshold`, default: 3)
- [ ] Single failure resets consecutive counter on next success
- [ ] Dry-run mode validates all prompt files exist and are readable
- [ ] Dry-run mode validates AI command binary exists and is executable
- [ ] Dry-run mode displays assembled prompt with clear section markers
- [ ] Dry-run mode displays resolved configuration with provenance
- [ ] Dry-run mode exits with code 0 if all validations pass
- [ ] Dry-run mode exits with code 1 if any validation fails (with clear error message)
- [ ] Dry-run mode does not execute AI CLI
- [ ] AI CLI output always captured and scanned for `<promise>` signals, regardless of `--verbose`
- [ ] AI CLI output buffered with configurable max size (`loop.max_output_buffer`, default: 10485760 bytes = 10MB)
- [ ] If output exceeds buffer size, buffer truncated from beginning (keeps most recent output), warning logged
- [ ] Output buffer size can be overridden per-procedure via procedure `max_output_buffer`
- [ ] Promise signal scanning happens after AI CLI exits (not during streaming)
- [ ] AI CLI output streamed to the terminal in real-time when `--verbose` flag is set
- [ ] Without `--verbose`, only loop-level progress (iteration start/complete, timing, outcome) is displayed
- [ ] Loop displays iteration statistics (count, min, max, mean, stddev) at completion
- [ ] Iteration statistics use constant memory (O(1)) regardless of iteration count
- [ ] Partial output from crashed AI CLI processes scanned for `<promise>` signals
- [ ] Loop terminates with status `success` when AI signals `<promise>SUCCESS</promise>`
- [ ] Loop terminates with status `max-iters` when max iterations reached
- [ ] Loop terminates with status `aborted` when failure threshold exceeded
- [ ] Loop terminates with status `interrupted` when SIGINT/SIGTERM received
- [ ] Exit code 0 for `success`, 1 for `aborted`, 2 for `max-iters`, 130 for `interrupted`
- [ ] Loop log level configurable via `loop.log_level` (debug, info, warn, error, built-in default: info)
- [ ] AI output streaming configurable via `loop.show_ai_output` (true, false, built-in default: false)
- [ ] `ROODA_LOG_LEVEL` environment variable sets `loop.log_level`
- [ ] `ROODA_SHOW_AI_OUTPUT` environment variable sets `loop.show_ai_output`
- [ ] `--verbose` flag overrides `loop.show_ai_output` to true
- [ ] `--quiet` flag overrides `loop.log_level` to warn
- [ ] `--log-level=<level>` flag overrides `loop.log_level`
- [ ] SIGINT/SIGTERM kills the AI CLI process, waits for termination (with timeout), and exits with status `interrupted`
- [ ] If AI CLI exceeds iteration timeout, process is killed and iteration counts as failure (increments `ConsecutiveFailures`)
- [ ] Iteration timeout configurable via `loop.iteration_timeout` (seconds, nil = no timeout, built-in default: nil)
- [ ] Timeout can be overridden per-procedure via procedure `iteration_timeout`
- [ ] `ROODA_LOOP_ITERATION_TIMEOUT` environment variable sets `loop.iteration_timeout`
- [ ] If AI CLI doesn't terminate within timeout, log warning and exit anyway
- [ ] Interrupted loops exit with code 130 (standard SIGINT exit code)

## Data Structures

### IterationState

```go
type IterationState struct {
    Iteration           int           // Current iteration number (0-indexed)
    MaxIterations       *int          // Termination threshold (nil = unlimited)
    IterationTimeout    *int          // Per-iteration timeout in seconds (nil = no timeout)
    MaxOutputBuffer     int           // Max AI CLI output buffer size in bytes (default: 10485760 = 10MB)
    ConsecutiveFailures int           // Consecutive AI CLI failures
    FailureThreshold    int           // Max consecutive failures before abort (default: 3)
    StartedAt           time.Time     // When the loop started
    Status              LoopStatus    // running, completed, aborted, interrupted
    ProcedureName       string        // Name of the procedure being executed
    Stats               IterationStats // Running statistics for iteration timing
}

type IterationStats struct {
    Count     int           // Total iterations completed
    TotalTime time.Duration // Sum of all iteration durations
    MinTime   time.Duration // Fastest iteration (0 if no iterations)
    MaxTime   time.Duration // Slowest iteration (0 if no iterations)
    M2        float64       // Sum of squared differences from mean (for variance calculation)
}
```

**Fields:**
- `Iteration` — Current iteration count, starts at 0, increments after each completed cycle
- `MaxIterations` — Termination threshold from `--max-iterations` > `--unlimited` > procedure `default_max_iterations` > `loop.default_max_iterations` > built-in default (5). Nil means unlimited.
- `IterationTimeout` — Per-iteration timeout in seconds. From `loop.iteration_timeout`; nil means no timeout (default). If AI CLI execution exceeds this duration, process is killed and iteration counts as failure.
- `MaxOutputBuffer` — Maximum AI CLI output buffer size in bytes. From `loop.max_output_buffer`; default 10485760 (10MB). If output exceeds this, buffer is truncated from beginning (keeping most recent output for signal scanning).
- `ConsecutiveFailures` — Resets to 0 on any successful iteration
- `FailureThreshold` — From `loop.failure_threshold`; default 3. Loop aborts when `ConsecutiveFailures >= FailureThreshold`
- `StartedAt` — Wall clock time when loop began (for total elapsed calculation)
- `Status` — State machine: running → success | max-iters | aborted | interrupted
- `ProcedureName` — Which procedure is executing (for logging)
- `Stats` — Running statistics for iteration timing (constant memory regardless of iteration count)

**IterationStats Fields:**
- `Count` — Total iterations completed
- `TotalTime` — Sum of all iteration durations (for mean calculation)
- `MinTime` — Fastest iteration duration (0 if no iterations completed)
- `MaxTime` — Slowest iteration duration (0 if no iterations completed)
- `M2` — Sum of squared differences from mean (for variance/stddev calculation using Welford's online algorithm)

### Loop Configuration

The `loop` section in `rooda-config.yml` defines global defaults for loop execution. Procedures can override `default_max_iterations` and `ai_cmd_alias`.

```yaml
loop:
  iteration_mode: max-iterations  # "max-iterations" or "unlimited" (built-in default: max-iterations)
  default_max_iterations: 5       # Global default (built-in default: 5). Must be >= 1. Ignored when mode is unlimited.
  iteration_timeout: 300          # Per-iteration timeout in seconds (nil/omitted = no timeout, built-in default: nil)
  max_output_buffer: 10485760     # Max AI CLI output buffer in bytes (built-in default: 10485760 = 10MB)
  failure_threshold: 3            # Consecutive failures before abort (built-in default: 3)
  log_level: info                 # "debug", "info", "warn", "error" (built-in default: info)
  show_ai_output: false           # Stream AI CLI output to terminal (built-in default: false)
  ai_cmd_alias: claude            # Default AI command alias for all procedures

procedures:
  bootstrap:
    default_max_iterations: 1  # Overrides loop.default_max_iterations for this procedure
    observe: prompts/observe_bootstrap.md
    orient: prompts/orient_bootstrap.md
    decide: prompts/decide_bootstrap.md
    act: prompts/act_bootstrap.md
  build:
    default_max_iterations: 10
    ai_cmd_alias: thorough     # This procedure uses a beefier model
    observe: prompts/observe_plan_specs_impl.md
    orient: prompts/orient_build.md
    decide: prompts/decide_build.md
    act: prompts/act_build.md
  long-running:
    iteration_mode: unlimited  # Overrides loop.iteration_mode for this procedure
    observe: prompts/observe_plan_specs_impl.md
    orient: prompts/orient_build.md
    decide: prompts/decide_build.md
    act: prompts/act_build.md
```

**Precedence for max iterations:**
1. `--max-iterations N` CLI flag (highest)
2. `--unlimited` CLI flag (sets to nil)
3. Procedure `iteration_mode` + `default_max_iterations`
4. `loop.iteration_mode` + `loop.default_max_iterations`
5. Built-in default: mode=max-iterations, count=5

**Precedence for AI command (consistent with max iterations — procedure overrides loop, env vars set loop level):**
1. `--ai-cmd` CLI flag (direct command, highest)
2. `--ai-cmd-alias` CLI flag (alias name)
3. Procedure `ai_cmd` (direct command)
4. Procedure `ai_cmd_alias` (alias name)
5. `loop.ai_cmd` (merged: env var > workspace > global)
6. `loop.ai_cmd_alias` (merged: env var > workspace > global)
7. Error — no AI command configured

### LoopStatus

```go
type LoopStatus string

const (
    StatusRunning     LoopStatus = "running"
    StatusSuccess     LoopStatus = "success"     // AI signaled SUCCESS
    StatusMaxIters    LoopStatus = "max-iters"   // Max iterations reached
    StatusAborted     LoopStatus = "aborted"     // Failure threshold exceeded
    StatusInterrupted LoopStatus = "interrupted" // User pressed Ctrl+C (SIGINT/SIGTERM)
)
```

**Exit Codes:**
- 0: StatusSuccess (AI signaled completion)
- 1: StatusAborted (failure threshold exceeded)
- 2: StatusMaxIters (max iterations reached, work may be incomplete)
- 130: StatusInterrupted (SIGINT/SIGTERM)
```

### Iteration Outcome Matrix

Each iteration's outcome is determined by two independent signals: the AI CLI process exit code and output markers emitted by the AI agent.

**Output signals:**
- `<promise>SUCCESS</promise>` — AI agent declares the job is done
- `<promise>FAILURE</promise>` — AI agent declares it is blocked and cannot make further progress (not a single test failure — the agent has exhausted what it can do)

**Signal precedence:** If both SUCCESS and FAILURE signals are present in output, FAILURE takes precedence.

| Exit Code | Output Signal | Outcome |
|---|---|---|
| 0 | none | Success — reset `ConsecutiveFailures`, continue |
| 0 | `SUCCESS` | Job done — terminate loop as `completed` |
| 0 | `FAILURE` | Agent-reported failure — increment `ConsecutiveFailures`, continue |
| 0 | both | FAILURE wins — increment `ConsecutiveFailures`, continue |
| non-zero | none | Process failure — increment `ConsecutiveFailures`, continue |
| non-zero | `SUCCESS` | Job done — terminate loop as `completed` (signal wins) |
| non-zero | `FAILURE` | Both failed — increment `ConsecutiveFailures`, continue |
| non-zero | both | FAILURE wins — increment `ConsecutiveFailures`, continue |

## Algorithm

1. Load configuration (CLI flags > env vars > workspace config > global config > built-in defaults)
2. Resolve AI command (see configuration.md AI Command Resolution)
3. Initialize fresh `IterationState`
4. Register signal handlers for SIGINT and SIGTERM:

```
function HandleSignal(state *IterationState, aiProcess *os.Process):
    log.Info("Received interrupt signal")
    
    if aiProcess != nil:
        // Kill AI CLI process
        aiProcess.Kill()
        
        // Wait for termination with timeout
        done = make(chan bool)
        go func():
            aiProcess.Wait()
            done <- true
        
        select:
            case <-done:
                log.Info("AI CLI process terminated")
            case <-time.After(5 * time.Second):
                log.Warn("AI CLI process did not terminate within 5s timeout")
    
    state.Status = interrupted
    os.Exit(130)  // Standard SIGINT exit code
```

5. Enter iteration loop:

```
function RunLoop(state IterationState, config Config) -> LoopStatus:
    while true:
        // Termination check
        if state.MaxIterations != nil AND state.Iteration >= *state.MaxIterations:
            state.Status = max-iters
            break

        if state.ConsecutiveFailures >= state.FailureThreshold:
            state.Status = aborted
            log.Error("Aborting: %d consecutive failures", state.ConsecutiveFailures)
            break

        // Start iteration
        iterationStart = time.Now()
        log.Info("Starting iteration %d", state.Iteration+1)

        // Assemble prompt from OODA phase files
        prompt, err = AssemblePrompt(config.Procedure)
        if err:
            log.Error("Prompt assembly failed: %v", err)
            state.Status = aborted
            break

        // Execute AI CLI with assembled prompt
        output, exitCode, err = ExecuteAICLI(config.AICommand, prompt, config.Verbose, state.IterationTimeout)
        if err:
            if err == ErrTimeout:
                log.Warn("Iteration %d: AI CLI exceeded timeout (%ds)", state.Iteration+1, *state.IterationTimeout)
                state.ConsecutiveFailures++
                elapsed = time.Since(iterationStart)
                updateStats(&state.Stats, elapsed)
                state.Iteration++
                continue
            log.Error("AI CLI execution failed: %v", err)
            state.Status = aborted
            break

        // Scan output for promise signals
        hasSuccess = strings.Contains(output, "<promise>SUCCESS</promise>")
        hasFailure = strings.Contains(output, "<promise>FAILURE</promise>")

        // Determine outcome per matrix (FAILURE wins if both present)
        if hasFailure:
            // Agent explicitly blocked - increment failure counter
            state.ConsecutiveFailures++
            log.Warn("Iteration %d: AI signaled FAILURE (consecutive: %d)", 
                state.Iteration+1, state.ConsecutiveFailures)
        else if hasSuccess:
            // Job complete - terminate loop
            log.Info("Iteration %d: AI signaled SUCCESS", state.Iteration+1)
            state.Status = success
            elapsed = time.Since(iterationStart)
            updateStats(&state.Stats, elapsed)
            log.Info("Iteration %d completed in %v (SUCCESS)", state.Iteration+1, elapsed)
            break
        else if exitCode == 0:
            // Success - reset failure counter
            state.ConsecutiveFailures = 0
            log.Info("Iteration %d succeeded", state.Iteration+1)
        else:
            // Process failure - increment counter
            state.ConsecutiveFailures++
            log.Warn("Iteration %d failed with exit code %d (consecutive: %d)", 
                state.Iteration+1, exitCode, state.ConsecutiveFailures)

        // Record timing
        elapsed = time.Since(iterationStart)
        updateStats(&state.Stats, elapsed)
        log.Info("Iteration %d completed in %v", state.Iteration+1, elapsed)

        // Increment iteration counter
        state.Iteration++

    return state.Status
```

### Statistics Update (Welford's Online Algorithm)

```
function updateStats(stats *IterationStats, elapsed time.Duration):
    stats.Count++
    stats.TotalTime += elapsed
    
    // Update min/max
    if stats.Count == 1 OR elapsed < stats.MinTime:
        stats.MinTime = elapsed
    if elapsed > stats.MaxTime:
        stats.MaxTime = elapsed
    
    // Welford's online algorithm for variance
    delta = float64(elapsed) - (float64(stats.TotalTime) / float64(stats.Count))
    stats.M2 += delta * (float64(elapsed) - (float64(stats.TotalTime) / float64(stats.Count)))

function getMean(stats IterationStats) -> time.Duration:
    if stats.Count == 0:
        return 0
    return stats.TotalTime / time.Duration(stats.Count)

function getStdDev(stats IterationStats) -> time.Duration:
    if stats.Count < 2:
        return 0
    variance = stats.M2 / float64(stats.Count)
    return time.Duration(math.Sqrt(variance))
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--unlimited` flag passed | Loop runs indefinitely until Ctrl+C, failure threshold, or AI signals `SUCCESS` |
| No max iterations configured anywhere | Uses built-in default: mode=max-iterations, count=5 |
| AI output contains `<promise>SUCCESS</promise>` | Loop terminates with status `success` regardless of exit code |
| AI output contains `<promise>FAILURE</promise>` with exit code 0 | Counts as failure — increment `ConsecutiveFailures` (output signal overrides exit code) |
| AI CLI exits non-zero, no output signal | Process failure — increment `ConsecutiveFailures`, continue loop |
| AI output contains both `SUCCESS` and `FAILURE` | FAILURE takes precedence — increment `ConsecutiveFailures`, continue |
| MaxIterations = 1 | Single iteration (iteration 0, displayed as 1), then exit with status `max-iters` (or `aborted` if iteration fails and threshold is 1) |
| MaxIterations = 5 | Runs iterations 0-4 (displayed as 1-5), termination check `Iteration >= 5` prevents iteration 5 from starting |
| Ctrl+C during AI CLI execution | Signal handler kills AI CLI process, waits for termination (5s timeout), exits with status `interrupted` and code 130 |
| Ctrl+C between iterations | Signal handler exits immediately with status `interrupted` and code 130 |
| Consecutive failures reach configured threshold (default: 3) | Loop aborts with status `aborted` |
| Successful iteration after 2 failures | `ConsecutiveFailures` resets to 0, loop continues |
| AI CLI output exceeds max buffer size | Buffer truncated from beginning, keeps most recent output, warning logged, signal scanning uses truncated buffer |
| AI CLI execution exceeds iteration timeout | Process killed, iteration counts as failure, loop continues (unless failure threshold reached) |
| No iteration timeout configured | AI CLI can run indefinitely (user must Ctrl+C to interrupt) |
| AI CLI crashes (SIGSEGV, SIGKILL, OOM) | Partial output scanned for signals; outcome determined per matrix (crash exit codes are non-zero) |
| Dry-run mode | Validates prompt files and AI command exist, displays assembled prompt and resolved config with provenance, exits with code 0 (success) or 1 (validation failed) |

## Dependencies

- **prompt-composition** — Assembles four OODA phase files and optional user context into a single prompt, given iteration state and config
- **ai-cli-integration** — Executes AI CLI tool and captures exit status
- **error-handling** — Retry logic, timeout, failure detection patterns
- **observability** — Structured logging, timing, progress display
- **cli-interface** — Provides `--max-iterations`, `--unlimited`, `--dry-run`, `--context`
- **configuration** — Resolves iteration settings and AI command from three-tier config system

## Implementation Mapping

**Source files:**
- `cmd/rooda/main.go` — CLI entry point, parses flags, invokes loop
- `internal/loop/loop.go` — Core iteration loop (`RunLoop` function)
- `internal/loop/signals.go` — Signal handling (SIGINT, SIGTERM)

**Related specs:**
- `prompt-composition.md` — How prompts are assembled before each iteration
- `ai-cli-integration.md` — How prompts are piped to AI CLI tools
- `error-handling.md` — Failure detection and consecutive failure logic
- `observability.md` — Logging, timing, progress display
- `cli-interface.md` — CLI flags that control loop behavior
- `configuration.md` — Where iteration settings come from

## Examples

### Example 1: Limited Iterations (Happy Path)

**Input:**
```bash
rooda build --max-iterations 3
```

**Expected Output:**
```
[10:00:00] Starting procedure: build (max 3 iterations)
[10:00:00] Iteration 1/3 starting...
[10:00:45] Iteration 1/3 completed in 45.2s (success)
[10:00:45] Iteration 2/3 starting...
[10:01:24] Iteration 2/3 completed in 38.7s (success)
[10:01:24] Iteration 3/3 starting...
[10:02:16] Iteration 3/3 completed in 52.1s (success)
[10:02:16] Reached max iterations: 3 (total: 2m16s)
  Iteration timing: min=38.7s, max=52.1s, mean=45.3s, stddev=5.5s
```

**Verification:**
- Three OODA cycles execute
- Only loop-level progress shown (no AI CLI output without `--verbose`)
- Loop terminates after iteration 3 with status `max-iters` and exit code 2
- Each iteration shows timing and outcome

### Example 2: Consecutive Failure Abort

**Input:**
```bash
rooda build --max-iterations 10
# AI CLI fails on iterations 4, 5, 6
```

**Expected Output:**
```
[10:00:00] Starting procedure: build (max 10 iterations)
[10:00:00] Iteration 1/10 starting...
[10:00:42] Iteration 1/10 completed in 42.0s (success)
[10:00:42] Iteration 2/10 starting...
[10:01:18] Iteration 2/10 completed in 36.0s (success)
[10:01:18] Iteration 3/10 starting...
[10:01:55] Iteration 3/10 completed in 37.0s (success)
[10:01:55] Iteration 4/10 starting...
[10:02:10] Iteration 4/10 completed in 15.0s (failure, consecutive: 1/3)
[10:02:10] Iteration 5/10 starting...
[10:02:22] Iteration 5/10 completed in 12.0s (failure, consecutive: 2/3)
[10:02:22] Iteration 6/10 starting...
[10:02:35] Iteration 6/10 completed in 13.0s (failure, consecutive: 3/3)
[10:02:35] ERROR: Aborting after 3 consecutive failures (6 iterations completed, total: 2m35s)
  Iteration timing: min=12.0s, max=42.0s, mean=25.8s, stddev=12.3s
```

**Verification:**
- Loop aborts after configured threshold of consecutive failures
- Status is `aborted`
- Failure count progression visible in iteration summaries

### Example 3: Verbose Mode

**Input:**
```bash
rooda build --max-iterations 2 --verbose
```

**Expected Output:**
```
[10:00:00] Starting procedure: build (max 2 iterations)
[10:00:00] Iteration 1/2 starting...
  ... created internal/loop/loop.go
  ... running go test ./...
  ok  	rooda/internal/loop	0.342s
  ✓ All tests passing
[10:00:45] Iteration 1/2 completed in 45.2s (success)
[10:00:45] Iteration 2/2 starting...
  ... updated internal/loop/signals.go
  ... running go test ./...
  ok  	rooda/internal/loop	0.298s
  ✓ All tests passing
[10:01:24] Iteration 2/2 completed in 38.7s (success)
[10:01:24] Reached max iterations: 2 (total: 1m24s)
```

**Verification:**
- AI CLI output streamed live between iteration markers
- Full visibility into what the agent is doing each iteration

### Example 4: Dry-Run Mode

**Input:**
```bash
rooda build --dry-run
```

**Expected Output:**
```
[DRY RUN] Procedure: build
[DRY RUN] Would execute with: claude-cli --no-interactive

# OODA Loop Iteration

## OBSERVE
[contents of observe prompt file]

## ORIENT
[contents of orient prompt file]

## DECIDE
[contents of decide prompt file]

## ACT
[contents of act prompt file]
```

**Verification:**
- Full assembled prompt displayed
- AI CLI not invoked
- Exit with status `max-iters` (dry-run doesn't execute, so no success signal)

### Example 5: Dry-Run Mode with User Context

**Input:**
```bash
rooda build --dry-run --context "focus on the auth module, the JWT validation is broken"
```

**Expected Output:**
```
[DRY RUN] Procedure: build
[DRY RUN] Would execute with: claude-cli --no-interactive

# OODA Loop Iteration

## CONTEXT
focus on the auth module, the JWT validation is broken

## OBSERVE
[contents of observe prompt file]

## ORIENT
[contents of orient prompt file]

## DECIDE
[contents of decide prompt file]

## ACT
[contents of act prompt file]
```

**Verification:**
- User context appears as a dedicated section before the OODA phases
- Context is passed through verbatim, not interpreted by the loop
- AI CLI not invoked
- Exit with status `max-iters` (dry-run doesn't execute, so no success signal)

### Example 6: Unlimited Iterations with Recovery

**Input:**
```bash
rooda build --unlimited
# AI CLI fails once on iteration 3, then succeeds on iteration 4
```

**Expected Output:**
```
[10:00:00] Starting procedure: build (unlimited)
...
[10:01:55] Iteration 3 starting...
[10:02:10] Iteration 3 completed in 15.0s (failure, consecutive: 1/3)
[10:02:10] Iteration 4 starting...
[10:02:52] Iteration 4 completed in 42.0s (success)
...
```

**Verification:**
- Consecutive failure counter resets to 0 after iteration 4 succeeds
- Loop continues until Ctrl+C, failure threshold, or `SUCCESS` signal

### Example 7: Dry-Run Validation (Success)

**Input:**
```bash
rooda build --dry-run --ai-cmd-alias claude
```

**Expected Output:**
```
=== Dry-Run: build ===

Configuration:
  AI Command: claude-cli --no-interactive (cli: --ai-cmd-alias "claude" → built-in alias)
  Max Iterations: 10 (workspace: ./rooda-config.yml)
  Iteration Timeout: none (built-in)
  Max Output Buffer: 10MB (built-in)
  Failure Threshold: 3 (built-in)

Validation:
  ✓ AI command binary exists: /usr/local/bin/claude-cli
  ✓ Prompt file exists: builtin:prompts/observe_plan_specs_impl.md
  ✓ Prompt file exists: builtin:prompts/orient_build.md
  ✓ Prompt file exists: builtin:prompts/decide_build.md
  ✓ Prompt file exists: builtin:prompts/act_build.md

Assembled Prompt (12,450 bytes):
────────────────────────────────────────
# OBSERVE
[observe phase content...]

# ORIENT
[orient phase content...]

# DECIDE
[decide phase content...]

# ACT
[act phase content...]
────────────────────────────────────────

Dry-run complete. Ready to execute: rooda build --ai-cmd-alias claude
```

**Verification:**
- All validations pass
- Prompt assembled and displayed with section markers
- Configuration shown with provenance
- Exit code 0

### Example 8: Dry-Run Validation (Failure)

**Input:**
```bash
rooda build --dry-run --ai-cmd nonexistent-cli
```

**Expected Output:**
```
=== Dry-Run: build ===

Configuration:
  AI Command: nonexistent-cli (cli: --ai-cmd)
  Max Iterations: 10 (workspace: ./rooda-config.yml)
  Iteration Timeout: none (built-in)
  Max Output Buffer: 10MB (built-in)
  Failure Threshold: 3 (built-in)

Validation:
  ✗ AI command binary not found: nonexistent-cli
    Searched PATH: /usr/local/bin:/usr/bin:/bin
  ✓ Prompt file exists: builtin:prompts/observe_plan_specs_impl.md
  ✓ Prompt file exists: builtin:prompts/orient_build.md
  ✓ Prompt file exists: builtin:prompts/decide_build.md
  ✓ Prompt file exists: builtin:prompts/act_build.md

Error: Dry-run validation failed
```

**Verification:**
- AI command validation fails
- Clear error message with searched paths
- Exit code 1
- Prompt not displayed (validation failed before assembly)

## Notes

**Origins — The Ralph Loop:**

The iteration loop in rooda descends directly from the [Ralph Loop](https://ghuntley.com/ralph/), an autonomous AI coding methodology developed by [Geoffrey Huntley](https://ghuntley.com). Huntley's core innovation was recognizing that a dumb bash loop — `while :; do cat PROMPT.md | claude ; done` — could drive an AI agent to build software autonomously by exploiting two properties: fresh context each iteration (preventing LLM degradation) and file-based state continuity (specs, plan, and AGENTS.md persist on disk between cycles). The Ralph Loop demonstrated that steering happens not through conversation but through deterministic file loading, upstream patterns, and downstream backpressure (tests reject invalid work). rooda evolves this by decomposing the single monolithic prompt into four composable OODA phase files (observe, orient, decide, act) that can be mixed and matched across 16 procedures — turning one loop with two modes (plan/build) into a general-purpose orchestration framework. The iteration mechanism itself remains faithful to Huntley's original insight: invoke the AI CLI as a fresh process, let it work, let it exit, persist state to files, repeat.

**Design Rationale — Fresh Context Per Iteration:**

LLMs advertise large context windows (200K+ tokens) but degrade in quality as context fills. Usable capacity is closer to 60% of the advertised window. By invoking the AI CLI as a fresh process each iteration, the agent stays perpetually in its "smart zone" where output quality remains high. The loop orchestrator (Go binary) persists across iterations, but the AI's memory is cleared each time the CLI process exits.

**File-Based State Continuity:**

While the AI's context clears between iterations, file-based state persists: AGENTS.md, work tracking, specs, and code remain on disk. The next iteration's prompt is assembled fresh from these files, giving the AI current context without conversational baggage. This is the key insight — state lives in files, not in conversation history.

**Why Go Instead of Bash:**

The v1 bash loop (archived in `archive/`) validated the core concept but had significant limitations: no error handling, no timing, no structured logging, platform-specific behavior. Go provides: testable code, cross-platform binary, structured error handling, signal handling, and the ability to embed default prompts.

**Iteration Counter Convention:**

Internal state is 0-indexed (starts at 0) for clean comparison logic (`Iteration >= MaxIterations`). Display is 1-indexed ("Iteration 1/10") to match user expectations. User-facing configuration (`--max-iterations`, `default_max_iterations` in config) is 1-indexed: `--max-iterations 3` means three iterations. Unlimited requires the explicit `--unlimited` flag or is never the implicit default — the built-in default is 5.

**No State Persistence By Design:**

The loop does not persist its own state to disk. File-based state continuity — AGENTS.md, work tracking, specs, and code on disk — is the resume mechanism. If the loop is interrupted, the user simply runs it again; the AI agent picks up from whatever state the files are in. This keeps the loop simple and avoids a class of bugs around stale state files, concurrent execution guards, and corrupted checkpoints.

**No Pause/Resume Support:**

The loop does not support pausing and resuming execution. Ctrl+C terminates the loop immediately (after killing the AI CLI process and waiting for cleanup). If interrupted, the user must restart the procedure from iteration 0.

This is by design: file-based state continuity (AGENTS.md, work tracking, specs, code on disk) provides the resume mechanism. The agent picks up from whatever state the files are in, regardless of which iteration the loop was on when interrupted. Restarting from iteration 0 gives the agent fresh context, which is often desirable.

If preserving iteration count or statistics across restarts becomes important, a future version could add state persistence via `--resume-from=<state-file>`. For now, the added complexity is not justified by the use case.

**Crash Handling:**

When the AI CLI crashes (segfault, OOM kill, kernel termination), Go's `exec.Command` returns whatever output was written to stdout/stderr before termination. This partial output is scanned for `<promise>` signals using the same logic as clean exits. Crashes produce non-zero exit codes, so the outcome matrix applies identically: if partial output contains `<promise>FAILURE</promise>`, the agent-reported failure is logged; otherwise, it's logged as a process failure. Both increment `ConsecutiveFailures`.

**Logging and Verbosity:**

Loop log level and AI output streaming follow the standard precedence chain: CLI flags (`--verbose`, `--quiet`, `--log-level`) > environment variables (`ROODA_LOG_LEVEL`, `ROODA_SHOW_AI_OUTPUT`) > workspace config > global config > built-in defaults (log_level=info, show_ai_output=false).

Log levels:
- **debug**: Prompt assembly, config resolution, signal scanning, internal state
- **info**: Iteration start/complete, timing, outcome, statistics (default)
- **warn**: Failures, timeouts, buffer truncation, signal handling
- **error**: Abort conditions, pre-execution failures

The `--verbose` flag is shorthand for setting `show_ai_output=true`. The `--quiet` flag is shorthand for setting `log_level=warn`. Both override all lower-precedence sources.

**Observability:**

Detailed logging behavior (format, output destination, structured fields), metrics export, and observability platform integrations are specified in `observability.md`. The iteration loop emits log events at defined levels (debug, info, warn, error) that the observability system captures, formats, and exports according to configuration.
