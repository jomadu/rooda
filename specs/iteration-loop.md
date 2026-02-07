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
- [ ] If AI CLI output contains `<promise>SUCCESS</promise>`, loop terminates with status `completed` regardless of exit code
- [ ] If AI CLI output contains `<promise>FAILURE</promise>`, iteration counts as a failure (increments `ConsecutiveFailures`) even if exit code is 0
- [ ] Each iteration invokes the AI CLI as a separate process, ensuring fresh context
- [ ] Iteration counter increments correctly (0-indexed internal, 1-indexed display)
- [ ] Max iterations resolved with precedence: `--max-iterations` CLI flag > `--unlimited` CLI flag > procedure `default_max_iterations` > `loop.default_max_iterations` > built-in default (5)
- [ ] `--unlimited` flag overrides max iterations to nil (no limit)
- [ ] Procedure-level `default_max_iterations` overrides global `loop.default_max_iterations`
- [ ] Progress displayed between iterations (iteration number, elapsed time)
- [ ] AI CLI exit status is captured and checked after each iteration
- [ ] Iteration outcome determined by output signals and exit code per the outcome matrix
- [ ] Consecutive failure count tracked; loop aborts after configurable threshold (`loop.failure_threshold`, default: 3)
- [ ] Single failure resets consecutive counter on next success
- [ ] Dry-run mode assembles and displays the prompt once without executing AI CLI, then exits
- [ ] AI CLI output always captured and scanned for `<promise>` signals, regardless of `--verbose`
- [ ] AI CLI output streamed to the terminal in real-time when `--verbose` flag is set
- [ ] Without `--verbose`, only loop-level progress (iteration start/complete, timing, outcome) is displayed
- [ ] Loop logs iteration number, elapsed time, and outcome (success/failure/completed) after each iteration
- [ ] SIGINT/SIGTERM kills the AI CLI process and exits immediately

## Data Structures

### IterationState

```go
type IterationState struct {
    Iteration           int           // Current iteration number (0-indexed)
    MaxIterations       *int          // Termination threshold (nil = unlimited)
    ConsecutiveFailures int           // Consecutive AI CLI failures
    FailureThreshold    int           // Max consecutive failures before abort (default: 3)
    StartedAt           time.Time     // When the loop started
    Status              LoopStatus    // running, completed, aborted
    ProcedureName       string        // Name of the procedure being executed
    ElapsedPerIteration []time.Duration // Timing for each completed iteration
}
```

**Fields:**
- `Iteration` — Current iteration count, starts at 0, increments after each completed cycle
- `MaxIterations` — Termination threshold from `--max-iterations` > `--unlimited` > procedure `default_max_iterations` > `loop.default_max_iterations` > built-in default (5). Nil means unlimited.
- `ConsecutiveFailures` — Resets to 0 on any successful iteration
- `FailureThreshold` — From `loop.failure_threshold`; default 3. Loop aborts when `ConsecutiveFailures >= FailureThreshold`
- `StartedAt` — Wall clock time when loop began (for total elapsed calculation)
- `Status` — State machine: running → completed | aborted
- `ProcedureName` — Which procedure is executing (for logging)
- `ElapsedPerIteration` — Duration of each iteration for performance visibility

### Loop Configuration

The `loop` section in `rooda-config.yml` defines global defaults for loop execution. Procedures can override `default_max_iterations`.

```yaml
loop:
  default_max_iterations: 5    # Global default (built-in default: 5)
  failure_threshold: 3         # Consecutive failures before abort (built-in default: 3)

procedures:
  bootstrap:
    default_max_iterations: 1  # Overrides loop.default_max_iterations for this procedure
    observe: prompts/observe_bootstrap.md
    orient: prompts/orient_bootstrap.md
    decide: prompts/decide_bootstrap.md
    act: prompts/act_bootstrap.md
  build:
    default_max_iterations: 10
    observe: prompts/observe_plan_specs_impl.md
    orient: prompts/orient_build.md
    decide: prompts/decide_build.md
    act: prompts/act_build.md
```

**Precedence for max iterations:**
1. `--max-iterations N` CLI flag (highest)
2. `--unlimited` CLI flag (sets to nil)
3. Procedure `default_max_iterations`
4. `loop.default_max_iterations`
5. Built-in default: 5

### LoopStatus

```go
type LoopStatus string

const (
    StatusRunning     LoopStatus = "running"
    StatusCompleted   LoopStatus = "completed"   // Max iterations reached or AI signaled SUCCESS
    StatusAborted     LoopStatus = "aborted"     // Failure threshold exceeded
)
```

### Iteration Outcome Matrix

Each iteration's outcome is determined by two independent signals: the AI CLI process exit code and output markers emitted by the AI agent.

**Output signals:**
- `<promise>SUCCESS</promise>` — AI agent declares the job is done
- `<promise>FAILURE</promise>` — AI agent declares it is blocked and cannot make further progress (not a single test failure — the agent has exhausted what it can do)

| Exit Code | Output Signal | Outcome |
|---|---|---|
| 0 | none | Success — reset `ConsecutiveFailures`, continue |
| 0 | `SUCCESS` | Job done — terminate loop as `completed` |
| 0 | `FAILURE` | Agent-reported failure — increment `ConsecutiveFailures`, continue |
| non-zero | none | Process failure — increment `ConsecutiveFailures`, continue |
| non-zero | `SUCCESS` | Job done — terminate loop as `completed` (signal wins) |
| non-zero | `FAILURE` | Both failed — increment `ConsecutiveFailures`, continue |

## Algorithm

1. Load configuration (CLI flags > env vars > workspace config > global config > built-in defaults)
2. Initialize fresh `IterationState`
3. Register signal handlers for SIGINT and SIGTERM
4. Enter iteration loop:

```
function RunLoop(state IterationState, config Config) -> LoopStatus:
    while true:
        // Termination check
        if state.MaxIterations != nil AND state.Iteration >= *state.MaxIterations:
            state.Status = completed
            break

        if state.ConsecutiveFailures >= state.FailureThreshold:
            state.Status = aborted
            log.Error("Aborting: %d consecutive failures", state.ConsecutiveFailures)
            break

        // Assemble prompt (delegates to prompt-composition)
        prompt = ComposePrompt(state, config)

        // Dry-run exits here
        if config.DryRun:
            Display(prompt)
            state.Status = completed
            break

        // Execute AI CLI (delegates to ai-cli-integration)
        // Output is streamed to terminal and scanned for promise signals
        iterationStart = time.Now()
        exitCode, output = ExecuteAICLI(config.AICommand, prompt)
        iterationDuration = time.Since(iterationStart)

        // Record timing
        state.ElapsedPerIteration = append(state.ElapsedPerIteration, iterationDuration)

        // Check output signals (take precedence over exit code)
        if output.Contains("<promise>SUCCESS</promise>"):
            state.Status = completed
            log.Info("AI signaled success — job complete")
            break

        if output.Contains("<promise>FAILURE</promise>"):
            state.ConsecutiveFailures++
            log.Warn("AI signaled failure, consecutive failures: %d/%d",
                state.ConsecutiveFailures, state.FailureThreshold)
        else if exitCode != 0:
            // No output signal — fall back to exit code
            state.ConsecutiveFailures++
            log.Warn("AI CLI failed (exit %d), consecutive failures: %d/%d",
                exitCode, state.ConsecutiveFailures, state.FailureThreshold)
        else:
            state.ConsecutiveFailures = 0

        // Increment
        state.Iteration++

        // Display progress
        log.Info("Iteration %d completed in %s", state.Iteration, iterationDuration)

    return state.Status
```

5. On signal (SIGINT/SIGTERM):
   - Log interruption
   - Exit with appropriate code

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--unlimited` flag passed | Loop runs indefinitely until Ctrl+C, failure threshold, or AI signals `SUCCESS` |
| No max iterations configured anywhere | Uses built-in default of 5 |
| AI output contains `<promise>SUCCESS</promise>` | Loop terminates with status `completed` regardless of exit code |
| AI output contains `<promise>FAILURE</promise>` with exit code 0 | Counts as failure — increment `ConsecutiveFailures` (output signal overrides exit code) |
| AI CLI exits non-zero, no output signal | Process failure — increment `ConsecutiveFailures`, continue loop |
| AI output contains both `SUCCESS` and `FAILURE` | `SUCCESS` takes precedence — loop terminates as `completed` |
| MaxIterations = 1 | Single iteration, then exit with status `completed` (or `aborted` if iteration fails and threshold is 1) |
| Ctrl+C during AI CLI execution | AI CLI process killed, loop exits cleanly |
| Ctrl+C between iterations | Loop exits cleanly |
| Consecutive failures reach configured threshold (default: 3) | Loop aborts with status `aborted` |
| Successful iteration after 2 failures | `ConsecutiveFailures` resets to 0, loop continues |
| Dry-run mode | Assemble and display prompt, exit without executing AI CLI |

## Dependencies

- **prompt-composition** — Assembles four OODA phase files and optional user context into a single prompt, given iteration state and config
- **ai-cli-integration** — Executes AI CLI tool and captures exit status
- **error-handling** — Retry logic, timeout, failure detection patterns
- **observability** — Structured logging, timing, progress display
- **cli-interface** — Provides `--max-iterations`, `--unlimited`, `--dry-run`, `--context`
- **configuration** — Resolves iteration settings from three-tier config system

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
```

**Verification:**
- Three OODA cycles execute
- Only loop-level progress shown (no AI CLI output without `--verbose`)
- Loop terminates after iteration 3 with status `completed`
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
- Exit with status `completed`

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
- Exit with status `completed`

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
