# Iteration Loop Control

## Job to be Done
Execute OODA loop procedures through controlled iteration cycles that clear context between runs, preventing LLM degradation while maintaining file-based state continuity.

## Activities
1. Initialize iteration counter to 0
2. Check termination conditions before each iteration
3. Execute single OODA cycle (create prompt, pipe to AI CLI, push changes)
4. Increment iteration counter
5. Display iteration progress
6. Repeat until termination condition met

## Acceptance Criteria
- [ ] Loop executes until max iterations reached or Ctrl+C pressed
- [x] Each iteration exits completely, clearing AI context (kiro-cli exits after each invocation; bash script persists by design)
- [ ] Iteration counter increments correctly
- [ ] Max iterations of 0 means unlimited (loop until Ctrl+C)
- [ ] Max iterations defaults to procedure config or 0 if not specified
- [ ] Progress displayed between iterations
- [ ] Git push happens after each iteration

## Data Structures

### Iteration State
```bash
ITERATION=0           # Current iteration number (0-indexed)
MAX_ITERATIONS=N      # Maximum iterations (0 = unlimited)
```

**Variables:**
- `ITERATION` - Current iteration count, starts at 0, increments after each cycle
- `MAX_ITERATIONS` - Termination threshold from CLI flag, config, or default 0

## Algorithm

1. Initialize ITERATION to 0
2. Enter infinite while loop
3. Check if MAX_ITERATIONS > 0 AND ITERATION >= MAX_ITERATIONS
4. If termination condition met, display message and break
5. Create combined OODA prompt from four phase files
6. Pipe prompt to kiro-cli with --no-interactive and --trust-all-tools flags
7. Push changes to git remote (create branch if needed)
8. Increment ITERATION counter
9. Display iteration separator with next iteration number
10. Loop back to step 3

**Pseudocode:**
```bash
ITERATION=0

while true; do
    if [ "$MAX_ITERATIONS" -gt 0 ] && [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
        echo "Reached max iterations: $MAX_ITERATIONS"
        break
    fi

    create_prompt | kiro-cli chat --no-interactive --trust-all-tools

    git push origin "$CURRENT_BRANCH" || {
        echo "Failed to push. Creating remote branch..."
        git push -u origin "$CURRENT_BRANCH"
    }

    ITERATION=$((ITERATION + 1))
    echo "\n\n======================== LOOP $ITERATION ========================\n"
done
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| MAX_ITERATIONS = 0 | Loop runs indefinitely until Ctrl+C |
| MAX_ITERATIONS = 1 | Single iteration, then exit |
| Ctrl+C during iteration | Bash catches signal, exits immediately |
| kiro-cli fails | Script continues to git push (no error handling) |
| git push fails | Attempts to create remote branch, continues loop |
| ITERATION reaches MAX_ITERATIONS | Displays message, breaks loop, script exits |

## Dependencies

- bash - Shell interpreter with arithmetic expansion
- kiro-cli - AI CLI tool for executing OODA prompts
- git - Version control for pushing changes after each iteration
- create_prompt() - Function that combines four OODA phase files

## Implementation Mapping

**Source files:**
- `src/rooda.sh` - Lines 161-179 implement the iteration loop

**Related specs:**
- `cli-interface.md` - Defines how MAX_ITERATIONS is set from CLI or config
- `prompt-composition.md` - Defines create_prompt() function behavior
- `ai-cli-integration.md` - Defines kiro-cli invocation (to be created)

## Examples

### Example 1: Limited Iterations

**Input:**
```bash
./rooda.sh build --max-iterations 3
```

**Expected Behavior:**
```
# Iteration 0 executes
======================== LOOP 1 ========================

# Iteration 1 executes
======================== LOOP 2 ========================

# Iteration 2 executes
======================== LOOP 3 ========================

Reached max iterations: 3
```

**Verification:**
- Three OODA cycles execute
- Loop terminates after ITERATION reaches 3
- Script exits cleanly

### Example 2: Unlimited Iterations

**Input:**
```bash
./rooda.sh bootstrap --max-iterations 0
```

**Expected Behavior:**
```
# Iteration 0 executes
======================== LOOP 1 ========================

# Iteration 1 executes
======================== LOOP 2 ========================

# ... continues until Ctrl+C
```

**Verification:**
- Loop runs indefinitely
- Only Ctrl+C terminates execution
- Each iteration increments counter

### Example 3: Default Iterations from Config

**Input:**
```bash
./rooda.sh build
# (config specifies default_iterations: 5)
```

**Expected Behavior:**
```
# 5 iterations execute
Reached max iterations: 5
```

**Verification:**
- MAX_ITERATIONS loaded from config (5 for build)
- Loop terminates after 5 iterations

### Example 4: Single Iteration

**Input:**
```bash
./rooda.sh bootstrap
# (config specifies default_iterations: 1)
```

**Expected Behavior:**
```
# Iteration 0 executes
======================== LOOP 1 ========================

Reached max iterations: 1
```

**Verification:**
- Single OODA cycle executes
- Loop terminates immediately after first iteration

## Notes

**Design Rationale:**

The iteration loop is the core mechanism that prevents LLM context degradation. Each iteration invokes kiro-cli as a separate process—the AI CLI starts fresh, processes the prompt, executes tools, then exits completely. This exit-and-restart pattern clears all AI context between iterations. The bash script itself persists across iterations (it's a single bash process running a while loop), but the AI's memory is cleared each time kiro-cli exits.

**Why Exit Between Iterations:**

LLMs advertise 200K token windows but degrade in quality as context fills. Usable capacity is closer to 176K, and performance drops significantly beyond 60% utilization. By invoking kiro-cli fresh each iteration (via pipe: `create_prompt | kiro-cli chat`), the AI stays perpetually in its "smart zone" (40-60% utilization) where output quality remains high.

**File-Based State:**

While the AI's context clears (kiro-cli exits), file-based state persists: AGENTS.md, work tracking, specs, and code remain on disk. The bash script continues running, and the next iteration pipes a fresh prompt to a new kiro-cli invocation. This provides continuity without conversational baggage.

**Iteration Counter:**

The counter is 0-indexed (starts at 0) but displays as 1-indexed in progress messages ("LOOP 1"). This matches user expectations (first iteration is "iteration 1") while keeping code logic simple (0-based comparison).

**Git Push Per Iteration:**

Pushing after each iteration creates a commit history that shows incremental progress. If the loop goes off track, you can revert to a previous iteration's state. The fallback branch creation handles first-time pushes gracefully.

**No Error Handling:**

The loop doesn't check if kiro-cli succeeds. If the AI CLI fails, the script continues to git push and loop. This is intentional—the loop trusts that file-based backpressure (tests, lints) will catch issues in subsequent iterations.

## Known Issues

**No kiro-cli error handling:** If kiro-cli exits with non-zero status, the loop continues anyway. This could lead to repeated failures without termination.

**Git push failures:** If git push fails for reasons other than missing remote branch, the error is silent and the loop continues.

**Iteration display off-by-one:** The separator shows "LOOP $ITERATION" after incrementing, so it displays the next iteration number, not the one that just completed. This is confusing but matches the implementation.

## Areas for Improvement

**Graceful error handling:** Check kiro-cli exit status and break loop after N consecutive failures.

**Iteration timing:** Display elapsed time per iteration to help users understand performance.

**Progress indicators:** Show which OODA phase is executing during long-running iterations.

**Dry-run mode:** Support --dry-run flag to show what would execute without actually running the AI CLI.

**Resume capability:** Save iteration state to file, allowing resume after Ctrl+C or failure.
