# Gap Analysis Plan: v2 Specs vs v0.1.0 Bash Implementation

**Generated:** 2026-02-08T22:26:05-08:00  
**Source:** Gap analysis between v2 specifications (11 spec files) and v0.1.0 bash implementation  
**Recommendation:** Proceed with Go rewrite per v2 specs. Bash served its purpose validating the OODA loop pattern.

## Summary

- **v2 specs:** 11 complete spec files defining 16 procedures, promise signals, failure tracking, timeouts, dry-run, context injection, signal handling, iteration statistics, output buffering, iteration modes, exit code semantics, provenance tracking, config tiers, Go implementation
- **v0.1.0 bash:** 9 procedures (56% coverage), basic iteration loop, OODA composition, AI CLI integration, config file, verbose/quiet modes
- **Gap:** 13 critical features missing (0% coverage), 7 procedures missing (44% coverage)
- **Root cause:** v2 specs describe Go rewrite; v0.1.0 is bash proof-of-concept

## Priority 0: Go Implementation Foundation

**Goal:** Establish Go project structure and core abstractions

1. **Initialize Go module and project structure**
   - Create `go.mod` with module name `github.com/maxdunn/ralph-wiggum-ooda`
   - Create directory structure: `cmd/rooda/`, `internal/{cli,config,loop,prompt,ai}/`, `prompts/`
   - Set up `.gitignore` for Go artifacts (`*.exe`, `*.test`, `vendor/`)
   - **Acceptance:** `go mod init` succeeds, directory structure matches spec

2. **Implement configuration system (configuration.md)**
   - Define `Config`, `LoopConfig`, `Procedure` structs
   - Implement three-tier config loading (workspace > global > built-in defaults)
   - Implement YAML parsing with validation
   - Implement provenance tracking
   - **Acceptance:** Can load rooda-config.yml, merge with built-in defaults, track provenance

3. **Implement CLI argument parsing (cli-interface.md)**
   - Define `CLIArgs` struct
   - Implement flag parsing with POSIX conventions
   - Implement help text generation
   - Implement `--list-procedures`, `--version`
   - **Acceptance:** `rooda --help`, `rooda --version`, `rooda --list-procedures` work

4. **Embed default prompts and procedures**
   - Use `//go:embed` to embed 25 prompt files from `prompts/` directory
   - Define 16 built-in procedures in `internal/config/defaults.go`
   - **Acceptance:** Binary contains embedded prompts, `rooda --list-procedures` shows all 16

## Priority 1: Core Loop Features

**Goal:** Implement iteration loop with promise signals, failure tracking, timeouts

5. **Implement basic iteration loop (iteration-loop.md)**
   - Define `IterationState`, `LoopStatus` types
   - Implement loop termination logic (max iterations, Ctrl+C)
   - Implement iteration counter and timing
   - **Acceptance:** `rooda build --max-iterations 3` runs 3 iterations and exits

6. **Implement promise signal scanning (iteration-loop.md, error-handling.md)**
   - Scan AI CLI output for `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>`
   - Implement outcome matrix (exit code + output signals → iteration outcome)
   - Terminate loop on SUCCESS signal
   - **Acceptance:** Loop terminates when AI outputs `<promise>SUCCESS</promise>`

7. **Implement failure tracking (error-handling.md)**
   - Track `ConsecutiveFailures` counter
   - Reset counter on success
   - Abort loop when threshold exceeded
   - **Acceptance:** Loop aborts after 3 consecutive failures (default threshold)

8. **Implement iteration timeouts (iteration-loop.md, error-handling.md)**
   - Add `iteration_timeout` config field
   - Kill AI CLI process if timeout exceeded
   - Count timeout as failure
   - **Acceptance:** Loop kills AI CLI after configured timeout, increments failure counter

9. **Implement signal handling (iteration-loop.md, error-handling.md)**
   - Register SIGINT/SIGTERM handlers
   - Kill AI CLI process on interrupt
   - Wait for cleanup (5s timeout)
   - Exit with code 130
   - **Acceptance:** Ctrl+C kills AI CLI cleanly, exits with code 130

## Priority 2: Observability

**Goal:** Provide visibility and control over loop execution

10. **Implement dry-run mode (cli-interface.md, iteration-loop.md)**
    - Add `--dry-run` flag
    - Validate config, prompt files, AI command
    - Display assembled prompt and resolved config with provenance
    - Exit with code 0 (success) or 1 (validation failed)
    - **Acceptance:** `rooda build --dry-run` validates and displays prompt without executing

11. **Implement iteration statistics (iteration-loop.md)**
    - Track min/max/mean/stddev using Welford's online algorithm
    - Display statistics at loop completion
    - Use constant memory (O(1))
    - **Acceptance:** Loop displays "Iteration timing: count=N min=Xs max=Xs mean=Xs stddev=Xs"

12. **Implement context injection (cli-interface.md, prompt-composition.md)**
    - Add `--context <text>` and `--context-file <path>` flags
    - Accumulate multiple contexts
    - Inject as dedicated section before OODA phases
    - **Acceptance:** `rooda build --context "focus on auth"` injects context into prompt

13. **Implement provenance display (configuration.md, observability.md)**
    - Track where each config value came from
    - Display provenance in dry-run mode
    - Display provenance with `--verbose`
    - **Acceptance:** Dry-run shows "max_iterations: 10 (from: workspace config)"

## Priority 3: Missing Procedures

**Goal:** Implement remaining 7 procedures from v2 specs

14. **Implement audit procedures (procedures.md)**
    - `audit-spec` — Quality assessment of specs
    - `audit-impl` — Quality assessment of implementation
    - `audit-agents` — Accuracy assessment of AGENTS.md
    - Create corresponding prompt files if not already present
    - **Acceptance:** `rooda audit-spec`, `rooda audit-impl`, `rooda audit-agents` execute

15. **Implement missing planning procedures (procedures.md)**
    - `draft-plan-spec-chore` — Plan spec maintenance
    - `draft-plan-impl-feat` — Plan new capability implementation
    - `draft-plan-impl-fix` — Plan implementation correction
    - `draft-plan-impl-chore` — Plan implementation maintenance
    - Create corresponding prompt files if not already present
    - **Acceptance:** All 16 procedures listed by `rooda --list-procedures`

## Priority 4: Polish and Documentation

**Goal:** Document undocumented features, improve error messages

16. **Document undocumented bash features in AGENTS.md**
    - Git push automation with fallback
    - Platform detection
    - Fuzzy procedure name matching
    - AI tool preset resolution
    - **Acceptance:** AGENTS.md Operational Learnings section updated

17. **Implement output buffering (iteration-loop.md, ai-cli-integration.md)**
    - Add `max_output_buffer` config field (default: 10MB)
    - Truncate from beginning if exceeded
    - Log warning when truncated
    - **Acceptance:** Large AI output truncated, warning logged

18. **Implement iteration modes (iteration-loop.md, configuration.md)**
    - Add `iteration_mode` field (max-iterations, unlimited)
    - Add `--unlimited` flag
    - Implement per-procedure mode override
    - **Acceptance:** `rooda build --unlimited` runs until SUCCESS signal or failure threshold

19. **Implement exit code semantics (iteration-loop.md, cli-interface.md)**
    - Exit 0 for success (SUCCESS signal)
    - Exit 1 for aborted (failure threshold)
    - Exit 2 for max-iters (incomplete work)
    - Exit 130 for interrupted (Ctrl+C)
    - **Acceptance:** Exit codes match spec for each termination condition

20. **Implement environment variable support (configuration.md)**
    - `ROODA_LOOP_AI_CMD`, `ROODA_LOOP_AI_CMD_ALIAS`
    - `ROODA_LOOP_ITERATION_MODE`, `ROODA_LOOP_DEFAULT_MAX_ITERATIONS`
    - `ROODA_LOOP_ITERATION_TIMEOUT`, `ROODA_LOOP_LOG_LEVEL`
    - `ROODA_LOOP_SHOW_AI_OUTPUT`, `ROODA_CONFIG_HOME`
    - **Acceptance:** Env vars override config file values

21. **Implement global config directory resolution (configuration.md)**
    - Resolve as: `ROODA_CONFIG_HOME` > `$XDG_CONFIG_HOME/rooda/` > `~/.config/rooda/` (Unix/macOS) > `%APPDATA%\rooda\` (Windows)
    - Load `<config_dir>/rooda-config.yml` if present
    - Merge with workspace config (workspace overrides global)
    - **Acceptance:** Global config at `~/.config/rooda/rooda-config.yml` loaded and merged

22. **Implement AI CLI integration (ai-cli-integration.md)**
    - Resolve AI command from precedence chain
    - Validate binary exists and is executable
    - Spawn process with prompt as stdin
    - Capture stdout/stderr to buffer
    - Stream to terminal if `--verbose`
    - **Acceptance:** AI CLI executed, output captured and scanned for signals

23. **Implement prompt composition (prompt-composition.md)**
    - Assemble four OODA phase files
    - Inject user context if provided
    - Support embedded and filesystem prompts
    - **Acceptance:** Prompt assembled from observe/orient/decide/act files

## Dependencies

- **Tasks 1-4** must complete before any other work (foundation)
- **Task 5** (basic loop) must complete before tasks 6-9 (loop features)
- **Task 22** (AI CLI integration) must complete before task 5 (loop needs AI execution)
- **Task 23** (prompt composition) must complete before task 5 (loop needs prompts)
- **Tasks 14-15** (procedures) can run in parallel after task 4 (embedded prompts)
- **Tasks 10-13** (observability) can run in parallel after task 5 (basic loop)
- **Tasks 16-21** (polish) can run in parallel after task 5 (basic loop)

## Estimated Effort

- **P0 (Foundation):** 16 hours (4 tasks × 4 hours)
- **P1 (Core Loop):** 20 hours (5 tasks × 4 hours)
- **P2 (Observability):** 12 hours (4 tasks × 3 hours)
- **P3 (Procedures):** 8 hours (2 tasks × 4 hours)
- **P4 (Polish):** 24 hours (8 tasks × 3 hours)
- **Total:** 80 hours (~2 weeks full-time)

## Notes

- Bash implementation (v0.1.0) validated the OODA loop pattern and should not be extended further
- Go rewrite provides: testability, cross-platform support, structured error handling, embedded prompts
- All 25 prompt files from bash implementation can be reused in Go (copy to `prompts/` directory)
- Config file format (YAML) remains compatible between bash and Go implementations
- Procedures can be implemented incrementally — start with most-used (build, bootstrap, publish-plan)
