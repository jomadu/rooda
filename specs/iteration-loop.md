# Iteration Loop

## Job to be Done

Execute AI coding agents through controlled OODA iteration cycles that clear AI context between runs, preventing LLM degradation while maintaining file-based state continuity. Each iteration invokes the AI CLI as a fresh process — the agent starts clean, processes the assembled prompt, executes tools, then exits. The loop orchestrator persists across iterations, managing termination, state, and quality gates.

This is the core loop — everything else feeds into or out of it.

## Activities

1. Initialize iteration state (counter, config, termination criteria)
2. Check termination conditions before each iteration
3. Assemble the OODA prompt from four phase files and optional context
4. Validate prompt is within token budget
5. Pipe assembled prompt to AI CLI tool
6. Capture AI CLI exit status and detect failures
7. Run quality gates (tests, lints) if configured
8. Persist iteration state to enable resume
9. Increment iteration counter and record timing
10. Display iteration progress
11. Repeat until termination condition met

## Acceptance Criteria

- [ ] Loop executes until max iterations reached, Ctrl+C pressed, or consecutive failure threshold hit
- [ ] Each iteration invokes the AI CLI as a separate process, ensuring fresh context
- [ ] Iteration counter increments correctly (0-indexed internal, 1-indexed display)
- [ ] Max iterations of 0 means unlimited (loop until Ctrl+C or failure threshold)
- [ ] Max iterations defaults to procedure config, then global config, then 0
- [ ] Progress displayed between iterations (iteration number, elapsed time)
- [ ] AI CLI exit status is captured and checked after each iteration
- [ ] Consecutive failure count tracked; loop aborts after configurable threshold (default: 3)
- [ ] Single failure resets consecutive counter on next success
- [ ] Iteration state persisted to disk after each iteration (enables resume)
- [ ] Resume from persisted state after interruption (Ctrl+C, crash, reboot)
- [ ] Dry-run mode assembles and displays prompt without executing AI CLI
- [ ] Token budget validated before piping prompt to AI CLI; warn if over budget
- [ ] Structured logging emits timing, phase, iteration number, and outcome per iteration
- [ ] Graceful shutdown on SIGINT/SIGTERM: finish current phase, persist state, exit cleanly
- [ ] Quality gates (if configured) run after AI CLI completes; failures count as iteration failures

## Data Structures

### IterationState

```go
type IterationState struct {
    Iteration           int           // Current iteration number (0-indexed)
    MaxIterations       int           // Termination threshold (0 = unlimited)
    ConsecutiveFailures int           // Consecutive AI CLI failures
    FailureThreshold    int           // Max consecutive failures before abort (default: 3)
    StartedAt           time.Time     // When the loop started
    LastIterationAt     time.Time     // When the last iteration completed
    Status              LoopStatus    // running, paused, completed, aborted, interrupted
    ProcedureName       string        // Name of the procedure being executed
    ElapsedPerIteration []time.Duration // Timing for each completed iteration
}
```

**Fields:**
- `Iteration` — Current iteration count, starts at 0, increments after each completed cycle
- `MaxIterations` — Termination threshold from CLI flag > procedure config > global config > 0
- `ConsecutiveFailures` — Resets to 0 on any successful iteration
- `FailureThreshold` — Configurable via config; default 3. Loop aborts when `ConsecutiveFailures >= FailureThreshold`
- `StartedAt` — Wall clock time when loop began (for total elapsed calculation)
- `LastIterationAt` — Timestamp of most recent iteration completion (for resume context)
- `Status` — State machine: running → completed | aborted | interrupted
- `ProcedureName` — Which procedure is executing (for state file naming and logging)
- `ElapsedPerIteration` — Duration of each iteration for performance visibility

### LoopStatus

```go
type LoopStatus string

const (
    StatusRunning     LoopStatus = "running"
    StatusCompleted   LoopStatus = "completed"   // Max iterations reached
    StatusAborted     LoopStatus = "aborted"     // Failure threshold exceeded
    StatusInterrupted LoopStatus = "interrupted" // Signal received (Ctrl+C)
)
```

### StateFile

```json
{
  "iteration": 3,
  "max_iterations": 10,
  "consecutive_failures": 0,
  "failure_threshold": 3,
  "started_at": "2026-02-06T10:00:00Z",
  "last_iteration_at": "2026-02-06T10:05:32Z",
  "status": "interrupted",
  "procedure_name": "build",
  "elapsed_per_iteration": ["45.2s", "38.7s", "52.1s"]
}
```

**Location:** `.rooda/state/<procedure-name>.json` in the workspace directory.

**Lifecycle:** Created on first iteration, updated after each iteration, deleted on clean completion. Presence of a state file with status `interrupted` signals that resume is available.

## Algorithm

1. Load configuration (CLI flags > env vars > workspace config > global config > built-in defaults)
2. Check for existing state file for this procedure
   - If found with status `interrupted`: offer resume, restore `IterationState`
   - If found with status `running`: warn about possible concurrent execution, abort
   - Otherwise: initialize fresh `IterationState`
3. Register signal handlers for SIGINT and SIGTERM
4. Enter iteration loop:

```
function RunLoop(state IterationState, config Config) -> LoopStatus:
    while true:
        // Termination check
        if state.MaxIterations > 0 AND state.Iteration >= state.MaxIterations:
            state.Status = completed
            break

        if state.ConsecutiveFailures >= state.FailureThreshold:
            state.Status = aborted
            log.Error("Aborting: %d consecutive failures", state.ConsecutiveFailures)
            break

        // Assemble prompt (delegates to prompt-composition)
        prompt = ComposePrompt(config.Procedure)

        // Token budget check (delegates to prompt-composition)
        if TokenCount(prompt) > config.TokenBudget:
            log.Warn("Prompt exceeds token budget: %d > %d", TokenCount(prompt), config.TokenBudget)
            // Continue anyway — warn but don't block

        // Dry-run exits here
        if config.DryRun:
            Display(prompt)
            state.Status = completed
            break

        // Execute AI CLI (delegates to ai-cli-integration)
        iterationStart = time.Now()
        exitCode = ExecuteAICLI(config.AICommand, prompt)
        iterationDuration = time.Since(iterationStart)

        // Record timing
        state.ElapsedPerIteration = append(state.ElapsedPerIteration, iterationDuration)

        // Handle AI CLI result
        if exitCode != 0:
            state.ConsecutiveFailures++
            log.Warn("AI CLI failed (exit %d), consecutive failures: %d/%d",
                exitCode, state.ConsecutiveFailures, state.FailureThreshold)
        else:
            state.ConsecutiveFailures = 0

            // Run quality gates if configured
            if config.QualityGates != nil:
                gateResult = RunQualityGates(config.QualityGates)
                if gateResult.Failed:
                    state.ConsecutiveFailures++
                    log.Warn("Quality gates failed: %s", gateResult.Summary)

        // Increment and persist
        state.Iteration++
        state.LastIterationAt = time.Now()
        PersistState(state)

        // Display progress
        log.Info("Iteration %d completed in %s", state.Iteration, iterationDuration)

    // Clean up
    if state.Status == completed:
        RemoveStateFile(state)
    else:
        PersistState(state)  // Preserve for resume

    return state.Status
```

5. On signal (SIGINT/SIGTERM):
   - Set `state.Status = interrupted`
   - Persist state to disk
   - Log interruption with resume instructions
   - Exit with appropriate code

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| MaxIterations = 0 | Loop runs indefinitely until Ctrl+C or failure threshold |
| MaxIterations = 1 | Single iteration, then clean exit with status `completed` |
| Ctrl+C during AI CLI execution | AI CLI process killed, state persisted as `interrupted`, resume available |
| Ctrl+C between iterations | State persisted immediately, clean exit |
| AI CLI exits non-zero | Increment `ConsecutiveFailures`, continue loop |
| 3 consecutive AI CLI failures | Loop aborts with status `aborted`, state persisted |
| AI CLI succeeds after 2 failures | `ConsecutiveFailures` resets to 0, loop continues |
| State file exists with `running` status | Warn about possible concurrent execution, refuse to start |
| State file exists with `interrupted` status | Offer resume from last completed iteration |
| Resume after iteration 5 of 10 | Restore state, continue from iteration 6 |
| Prompt exceeds token budget | Log warning, continue execution (non-blocking) |
| Dry-run mode | Assemble and display prompt, exit without executing AI CLI |
| Quality gate fails | Counts as iteration failure, increments `ConsecutiveFailures` |
| No quality gates configured | Skip quality gate step, only AI CLI exit status matters |
| State file corrupted or unparseable | Log warning, start fresh (don't crash) |
| Disk full — can't persist state | Log error, continue execution (state persistence is best-effort) |

## Dependencies

- **prompt-composition** — Assembles four OODA phase files into a single prompt
- **ai-cli-integration** — Executes AI CLI tool and captures exit status
- **error-handling** — Retry logic, timeout, failure detection patterns
- **observability** — Structured logging, timing, progress display
- **cli-interface** — Provides `--max-iterations`, `--dry-run`, `resume` subcommand
- **configuration** — Resolves iteration settings from three-tier config system

## Implementation Mapping

**Source files:**
- `cmd/rooda/main.go` — CLI entry point, parses flags, invokes loop
- `internal/loop/loop.go` — Core iteration loop (`RunLoop` function)
- `internal/loop/state.go` — State persistence and resume logic
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
[10:00:45] Iteration 1/3 completed in 45.2s
[10:00:45] Iteration 2/3 starting...
[10:01:24] Iteration 2/3 completed in 38.7s
[10:01:24] Iteration 3/3 starting...
[10:02:16] Iteration 3/3 completed in 52.1s
[10:02:16] Reached max iterations: 3 (total: 2m16s)
```

**Verification:**
- Three OODA cycles execute
- Loop terminates after iteration 3 with status `completed`
- State file cleaned up (no resume needed)
- Each iteration shows timing

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
[10:00:42] Iteration 1/10 completed in 42.0s
[10:00:42] Iteration 2/10 starting...
[10:01:18] Iteration 2/10 completed in 36.0s
[10:01:18] Iteration 3/10 starting...
[10:01:55] Iteration 3/10 completed in 37.0s
[10:01:55] Iteration 4/10 starting...
[10:02:10] WARNING: AI CLI failed (exit 1), consecutive failures: 1/3
[10:02:10] Iteration 4/10 completed in 15.0s
[10:02:10] Iteration 5/10 starting...
[10:02:22] WARNING: AI CLI failed (exit 1), consecutive failures: 2/3
[10:02:22] Iteration 5/10 completed in 12.0s
[10:02:22] Iteration 6/10 starting...
[10:02:35] WARNING: AI CLI failed (exit 1), consecutive failures: 3/3
[10:02:35] ERROR: Aborting after 3 consecutive failures (6 iterations completed, total: 2m35s)
```

**Verification:**
- Loop aborts after 3 consecutive failures, not at max iterations
- Status is `aborted`, state file preserved
- Warning messages show failure count progression

### Example 3: Resume After Interruption

**Input:**
```bash
# First run — interrupted at iteration 5
rooda build --max-iterations 10
# User presses Ctrl+C during iteration 5
```

**Output on interruption:**
```
[10:05:00] Iteration 5/10 starting...
^C
[10:05:12] Interrupted. State saved. Resume with: rooda resume build
```

**Resume:**
```bash
rooda resume build
```

**Output on resume:**
```
[10:10:00] Resuming procedure: build from iteration 5 (max 10)
[10:10:00] Previous session: 5 iterations completed in 4m30s
[10:10:00] Iteration 6/10 starting...
...
```

**Verification:**
- State file preserved on Ctrl+C with status `interrupted`
- Resume restores iteration counter to 5, continues from 6
- Previous session timing displayed for context

### Example 4: Dry-Run Mode

**Input:**
```bash
rooda build --dry-run
```

**Expected Output:**
```
[DRY RUN] Procedure: build
[DRY RUN] Would execute with: claude-cli --no-interactive
[DRY RUN] Token count: 4,230 / 100,000 budget

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
- Token count shown relative to budget
- Exit with status `completed`

### Example 5: Unlimited Iterations with Recovery

**Input:**
```bash
rooda build
# (config specifies default_iterations: 0)
# AI CLI fails once on iteration 3, then succeeds on iteration 4
```

**Expected Output:**
```
[10:00:00] Starting procedure: build (unlimited iterations)
...
[10:01:55] Iteration 3 starting...
[10:02:10] WARNING: AI CLI failed (exit 1), consecutive failures: 1/3
[10:02:10] Iteration 3 completed in 15.0s
[10:02:10] Iteration 4 starting...
[10:02:52] Iteration 4 completed in 42.0s
...
```

**Verification:**
- Consecutive failure counter resets to 0 after iteration 4 succeeds
- Loop continues indefinitely until Ctrl+C

## Notes

**Design Rationale — Fresh Context Per Iteration:**

LLMs advertise large context windows (200K+ tokens) but degrade in quality as context fills. Usable capacity is closer to 60% of the advertised window. By invoking the AI CLI as a fresh process each iteration, the agent stays perpetually in its "smart zone" where output quality remains high. The loop orchestrator (Go binary) persists across iterations, but the AI's memory is cleared each time the CLI process exits.

**File-Based State Continuity:**

While the AI's context clears between iterations, file-based state persists: AGENTS.md, work tracking, specs, and code remain on disk. The next iteration's prompt is assembled fresh from these files, giving the AI current context without conversational baggage. This is the key insight — state lives in files, not in conversation history.

**Why Go Instead of Bash:**

The v1 bash loop (archived in `archive/`) validated the core concept but had significant limitations: no error handling, no resume, no timing, no structured logging, platform-specific behavior. Go provides: testable code, cross-platform binary, structured error handling, signal handling, JSON state persistence, and the ability to embed default prompts.

**Iteration Counter Convention:**

Internal state is 0-indexed (starts at 0) for clean comparison logic (`Iteration >= MaxIterations`). Display is 1-indexed ("Iteration 1/10") to match user expectations. This matches the v1 convention.

**Quality Gates Are Optional:**

Quality gates (tests, lints) run after AI CLI completion but are configured per-project via AGENTS.md, not hardcoded. If no quality gates are configured, only AI CLI exit status determines success. This keeps the loop generic — it doesn't assume the project has tests.

**State Persistence Is Best-Effort:**

If state can't be written to disk (permissions, disk full), the loop continues. State persistence enables resume but is not required for basic operation. The loop should never crash because it can't write a state file.

**Concurrent Execution Guard:**

The `running` status in the state file acts as a simple lock to prevent two loop instances from executing the same procedure simultaneously. This is advisory, not enforced — the check happens at startup and logs a warning. A crashed loop that didn't clean up will leave a `running` state file; the user can manually delete it or the resume flow handles it.
