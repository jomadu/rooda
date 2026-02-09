# Gap Analysis Plan: v2 Specs vs v0.1.0 Implementation

**Context:** v2 specifications describe a Go rewrite with 16 procedures and 13 critical features. Current v0.1.0 bash implementation has 9 procedures and is missing most v2 features.

**Root Cause:** v2 specs describe the target Go implementation; v0.1.0 bash is the current working prototype.

**Resolution Strategy:** Implement v2 Go rewrite following the prioritized task list below.

---

## Priority 0: Foundation (Go Implementation Bootstrap)

### Task 1: Initialize Go Module and Project Structure
- Create `go.mod` with module name `github.com/maxdunn/ralph-wiggum-ooda`
- Create directory structure: `cmd/rooda/`, `internal/{config,loop,ai,cli,prompt}/`
- Add `.gitignore` for Go artifacts (`*.exe`, `*.test`, `vendor/`, `dist/`)
- **Acceptance:** `go mod init` succeeds, directory structure matches spec

### Task 2: Implement Core Data Structures
- Define `Config`, `LoopConfig`, `Procedure` structs per configuration.md
- Define `IterationState`, `LoopStatus`, `IterationStats` per iteration-loop.md
- Define `LogLevel`, `TimestampFormat`, `LogEvent` per observability.md
- **Acceptance:** All structs compile, match spec field definitions

### Task 3: Implement Built-in Defaults
- Embed 25 prompt files using `go:embed` directive
- Define built-in procedure definitions (16 procedures)
- Define built-in AI command aliases (kiro-cli, claude, copilot, cursor-agent)
- **Acceptance:** `BuiltInDefaults` variable compiles, all 16 procedures defined

---

## Priority 1: Core Loop (Iteration Engine)

### Task 4: Implement Prompt Composition
- Implement `AssemblePrompt()` function per prompt-composition.md
- Support four OODA phase files + optional user context
- Support embedded prompts (builtin: prefix) and filesystem prompts
- **Acceptance:** Assembled prompt matches spec format, context injection works

### Task 5: Implement AI CLI Execution
- Implement `ExecuteAICLI()` function per ai-cli-integration.md
- Pipe assembled prompt to AI CLI process
- Capture stdout/stderr with configurable buffer size
- **Acceptance:** AI CLI executes, output captured, exit code returned

### Task 6: Implement Promise Signal Scanning
- Scan output for `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>`
- Implement outcome matrix per iteration-loop.md (FAILURE wins if both present)
- **Acceptance:** Signal detection works, outcome matrix logic correct

### Task 7: Implement Consecutive Failure Tracking
- Track `ConsecutiveFailures` counter across iterations
- Reset counter on success, increment on failure
- Abort when `ConsecutiveFailures >= FailureThreshold`
- **Acceptance:** Failure tracking works, abort at threshold

### Task 8: Implement Iteration Loop
- Implement `RunLoop()` function per iteration-loop.md
- Termination conditions: max iterations, failure threshold, SUCCESS signal, SIGINT
- Iteration counter (0-indexed internal, 1-indexed display)
- **Acceptance:** Loop executes, terminates correctly, exit codes match spec

### Task 9: Implement Iteration Timeouts
- Add per-iteration timeout with SIGTERM → SIGKILL escalation
- Timeout always counts as failure (promise signals ignored)
- **Acceptance:** Timeout kills process, increments failure counter

### Task 10: Implement Signal Handling
- Handle SIGINT/SIGTERM: kill AI CLI, wait for termination, exit 130
- **Acceptance:** Ctrl+C terminates cleanly, no zombie processes

---

## Priority 2: Observability (Logging and Diagnostics)

### Task 11: Implement Structured Logging
- Implement log event emission at four levels (debug, info, warn, error)
- Implement logfmt formatting (timestamp, level, message, fields)
- Support five timestamp formats (time, time-ms, relative, iso, none)
- **Acceptance:** Log output matches spec format, levels filter correctly

### Task 12: Implement Iteration Statistics
- Implement Welford's online algorithm for constant-memory stats
- Calculate count, min, max, mean, stddev
- Display at loop completion (omit stddev when count < 2)
- **Acceptance:** Statistics correct, memory usage O(1)

### Task 13: Implement Dry-Run Mode
- Validate config, prompt files, AI command without executing
- Display assembled prompt and resolved config with provenance
- Exit 0 if valid, 1 if invalid
- **Acceptance:** Dry-run validates, displays prompt, doesn't execute AI CLI

### Task 14: Implement Verbose Mode
- Stream AI CLI output to terminal in real-time
- Set `show_ai_output=true` and `log_level=debug`
- **Acceptance:** Verbose mode streams output, shows debug logs

### Task 15: Implement Output Buffering
- Configurable max buffer size (default 10MB)
- Truncate from beginning if exceeded, keep most recent output
- Log warning when truncated
- **Acceptance:** Buffer truncation works, signals at end preserved

---

## Priority 3: Configuration (Three-Tier System)

### Task 16: Implement Config Loading
- Load built-in defaults → global config → workspace config → env vars → CLI flags
- Resolve global config directory (ROODA_CONFIG_HOME > XDG_CONFIG_HOME/rooda > ~/.config/rooda)
- Parse YAML config files using `gopkg.in/yaml.v3`
- **Acceptance:** Config loads from all tiers, merges correctly

### Task 17: Implement Config Merging
- Field-level merge (overlay non-empty/non-zero values override base)
- Procedures merge additively (workspace adds to built-in, doesn't replace)
- AI command aliases merge additively
- **Acceptance:** Merging works, workspace overrides global, both override built-in

### Task 18: Implement Provenance Tracking
- Track which tier provided each resolved value
- Display provenance in dry-run and verbose modes
- **Acceptance:** Provenance correct, displayed in dry-run

### Task 19: Implement Config Validation
- Validate at load time (fail fast before execution)
- Check required fields, type constraints, file existence
- Clear error messages with file path and line number
- **Acceptance:** Validation catches errors, messages actionable

### Task 20: Implement Environment Variable Resolution
- Support `ROODA_LOOP_*` environment variables per configuration.md
- Override config file values at loop level
- **Acceptance:** Env vars override config, CLI flags override env vars

### Task 21: Implement AI Command Resolution
- Resolve with precedence: --ai-cmd > --ai-cmd-alias > procedure ai_cmd > procedure ai_cmd_alias > loop.ai_cmd > loop.ai_cmd_alias
- Error if no AI command configured (list all ways to set one)
- **Acceptance:** Resolution follows precedence, error message helpful

### Task 22: Implement Max Iterations Resolution
- Resolve with precedence: --max-iterations > --unlimited > procedure iteration_mode/default_max_iterations > loop settings
- **Acceptance:** Resolution follows precedence, unlimited mode works

---

## Priority 4: Procedures (Missing 7 of 16)

### Task 23: Add Missing Audit Procedures
- Add `audit-spec`, `audit-impl`, `audit-agents`, `audit-spec-to-impl`, `audit-impl-to-spec`
- All use existing prompt files (observe_specs.md, observe_impl.md, etc.)
- Default max iterations: 1 (read-only assessments)
- **Acceptance:** All 5 audit procedures defined, use correct prompts

### Task 24: Add Missing Planning Procedures
- Add `draft-plan-spec-feat`, `draft-plan-impl-feat`
- Use existing prompt files (observe_story_task_specs_impl.md, etc.)
- Default max iterations: 5
- **Acceptance:** Both planning procedures defined, use correct prompts

### Task 25: Rename Existing Procedures to Match v2
- Rename `draft-plan-story-to-spec` → `draft-plan-spec-feat`
- Rename `draft-plan-bug-to-spec` → `draft-plan-spec-fix`
- Rename `draft-plan-spec-to-impl` → `audit-spec-to-impl` (audit, not planning)
- Rename `draft-plan-impl-to-spec` → `audit-impl-to-spec` (audit, not planning)
- Rename `draft-plan-spec-refactor` → `draft-plan-spec-refactor` (no change)
- Rename `draft-plan-impl-refactor` → `draft-plan-impl-refactor` (no change)
- Add `draft-plan-spec-chore`, `draft-plan-impl-chore`
- **Acceptance:** All 16 procedures match v2 naming, categories correct

---

## Priority 5: CLI Interface

### Task 26: Implement CLI Argument Parsing
- Parse procedure name, flags, and context per cli-interface.md
- Support short flags (-v, -q, -n, -u, -d, -c, -h)
- Validate mutually exclusive flags (--verbose/--quiet, --max-iterations/--unlimited)
- **Acceptance:** All flags parse correctly, validation works

### Task 27: Implement Help Text Generation
- Global help (--help with no procedure)
- Procedure-specific help (rooda <procedure> --help)
- List procedures (--list-procedures)
- **Acceptance:** Help text matches spec format, all procedures listed

### Task 28: Implement Exit Codes
- 0: success, 1: aborted, 2: max-iters, 130: interrupted
- User errors: 1, config errors: 2, execution errors: 3
- **Acceptance:** Exit codes match spec

---

## Priority 6: Distribution

### Task 29: Implement Single Binary Build
- Add `go build` target to produce `rooda` binary
- Embed all 25 prompt files using `go:embed`
- **Acceptance:** `go build` produces single binary, prompts embedded

### Task 30: Add Cross-Platform Support
- Test on macOS, Linux
- Handle platform-specific config directory resolution
- **Acceptance:** Binary runs on macOS and Linux

### Task 31: Add Installation Instructions
- Update README.md with installation steps
- Document `go install` or binary download
- **Acceptance:** Installation instructions clear, tested

---

## Priority 7: Testing

### Task 32: Add Unit Tests for Core Functions
- Test prompt composition (AssemblePrompt)
- Test config loading and merging
- Test AI command resolution
- Test max iterations resolution
- Test promise signal scanning
- Test failure tracking
- **Acceptance:** Tests pass, coverage > 70%

### Task 33: Add Integration Tests
- Test full loop execution with mock AI CLI
- Test dry-run mode
- Test verbose mode
- Test signal handling
- **Acceptance:** Integration tests pass

---

## Out of Scope (Documented but Not Implemented in v0.1.0)

**Undocumented bash features (working but not in specs):**
- Git push automation with fallback (creates remote branch if needed)
- Platform detection (macOS/Linux)
- Fuzzy procedure name matching (suggests closest match)
- AI tool preset resolution (hardcoded + custom from config)

**Action:** Document these in v2 specs or remove from Go implementation if not needed.

---

## Summary

**Total Tasks:** 33
- Priority 0 (Foundation): 3 tasks
- Priority 1 (Core Loop): 7 tasks
- Priority 2 (Observability): 5 tasks
- Priority 3 (Configuration): 7 tasks
- Priority 4 (Procedures): 3 tasks
- Priority 5 (CLI Interface): 3 tasks
- Priority 6 (Distribution): 3 tasks
- Priority 7 (Testing): 2 tasks

**Estimated Effort:** 
- Foundation: 1-2 days
- Core Loop: 3-4 days
- Observability: 2-3 days
- Configuration: 3-4 days
- Procedures: 1 day
- CLI Interface: 2 days
- Distribution: 1 day
- Testing: 2-3 days

**Total:** 15-23 days (3-5 weeks)

**Critical Path:** Foundation → Core Loop → Configuration → CLI Interface → Distribution

**Parallel Work Possible:** Observability, Procedures, Testing can be done in parallel with Configuration/CLI after Core Loop is complete.
