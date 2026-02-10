# Gap Analysis Plan: v2 Go Implementation

**Status:** Draft  
**Created:** 2026-02-09  
**Source:** Gap analysis between v2 specifications and current implementation

## Summary

All 11 v2 specifications are complete with JTBD structure, acceptance criteria, and examples. However, **no Go implementation exists** — the project currently runs on v0.1.0 bash (`rooda.sh`). This plan covers building the entire v2 Go implementation from scratch.

## Priority: Critical (P0) — Foundation

### T1: Initialize Go Project Structure
**Description:** Create `go.mod`, directory structure (`cmd/rooda/`, `internal/`), and basic `main.go` entry point.

**Acceptance Criteria:**
- `go.mod` exists with module name `github.com/maxdunn/ralph-wiggum-ooda`
- `cmd/rooda/main.go` exists with version flag working
- `go build ./...` succeeds
- `./rooda --version` prints version number

**Dependencies:** None

---

### T2: Implement Configuration System (Three-Tier Merge)
**Description:** Build config loading with built-in defaults, global config, workspace config, and environment variable support per `configuration.md`.

**Acceptance Criteria:**
- `internal/config/config.go` defines `Config`, `LoopConfig`, `Procedure` structs
- `internal/config/loader.go` implements three-tier merge (built-in > global > workspace)
- `internal/config/env.go` handles `ROODA_*` environment variables
- `internal/config/validate.go` validates merged config
- Global config directory resolves via `ROODA_CONFIG_HOME` > `XDG_CONFIG_HOME/rooda` > `~/.config/rooda`
- Config validation fails fast with clear error messages
- Unit tests cover merge precedence and validation

**Dependencies:** T1

---

### T3: Implement CLI Argument Parsing
**Description:** Parse command-line arguments per `cli-interface.md` with support for named procedures, OODA phase overrides, and global flags.

**Acceptance Criteria:**
- `internal/cli/parser.go` parses all flags from spec
- `rooda <procedure>` invokes named procedure
- `rooda --help` displays usage
- `rooda --list-procedures` lists available procedures
- `--max-iterations`, `--unlimited`, `--dry-run`, `--context`, `--ai-cmd`, `--ai-cmd-alias` flags work
- `--observe`, `--orient`, `--decide`, `--act` flags accumulate into arrays
- Flag validation with clear error messages
- Exit codes: 0 (success), 1 (user error), 2 (config error), 3 (execution error)

**Dependencies:** T2

---

## Priority: High (P1) — Core Loop

### T4: Implement AI CLI Integration
**Description:** Execute AI CLI tools with prompt piping, output capture, signal scanning per `ai-cli-integration.md`.

**Acceptance Criteria:**
- `internal/ai/executor.go` implements `ExecuteAICLI` function
- Prompt piped to AI CLI stdin
- Stdout/stderr captured to buffer (configurable max size, default 10MB)
- Output streamed to terminal when `--verbose` flag set
- Exit code captured
- `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>` signals detected
- Timeout handling (kill process after configured duration)
- Built-in aliases: `kiro-cli`, `claude`, `copilot`, `cursor-agent`
- AI command resolution precedence: CLI flag > procedure > loop config > error
- Unit tests for signal scanning, timeout, output truncation

**Dependencies:** T2, T3

---

### T5: Implement Prompt Composition (Fragment System)
**Description:** Assemble prompts from fragment arrays per `prompt-composition.md` and `procedures.md`.

**Acceptance Criteria:**
- `internal/procedures/fragments.go` loads fragments from embedded resources or filesystem
- Fragment paths with `builtin:` prefix load from embedded resources
- Fragment paths without prefix resolve relative to config directory
- Fragments support inline content via `content` field
- Fragments support Go text/template with parameters
- Fragment arrays concatenate with double newlines
- Validation: exactly one of `content` or `path` required per fragment
- Unit tests for fragment loading, template processing, path resolution

**Dependencies:** T2

---

### T6: Implement Iteration Loop
**Description:** Core OODA iteration loop per `iteration-loop.md` with fresh context per iteration, termination control, and state continuity.

**Acceptance Criteria:**
- `internal/loop/loop.go` implements `RunLoop` function
- Each iteration invokes AI CLI as fresh process
- Iteration counter (0-indexed internal, 1-indexed display)
- Max iterations resolved: CLI flag > procedure > loop config > built-in default (5)
- `--unlimited` flag sets max iterations to nil
- Consecutive failure tracking with configurable threshold (default: 3)
- Termination conditions: max iterations, failure threshold, SUCCESS signal, Ctrl+C
- Iteration outcome matrix (exit code + output signals)
- Iteration timing statistics (count, min, max, mean, stddev) using Welford's algorithm
- SIGINT/SIGTERM handling (kill AI CLI, wait for termination, exit 130)
- Unit tests for termination logic, failure tracking, statistics

**Dependencies:** T4, T5

---

## Priority: High (P1) — Procedures

### T7: Embed Built-in Prompt Fragments
**Description:** Create 55 prompt fragments organized by OODA phase and embed in binary via `go:embed`.

**Acceptance Criteria:**
- `fragments/observe/` contains 13 fragments per `procedures.md`
- `fragments/orient/` contains 20 fragments
- `fragments/decide/` contains 10 fragments
- `fragments/act/` contains 12 fragments
- `internal/procedures/builtin.go` embeds fragments via `//go:embed fragments/*`
- Fragment loading resolves `builtin:` prefix to embedded resources
- All 55 fragments exist and are readable

**Dependencies:** T5

---

### T8: Define 16 Built-in Procedures
**Description:** Define all 16 built-in procedures per `procedures.md` with fragment-based composition.

**Acceptance Criteria:**
- `internal/procedures/builtin.go` defines all 16 procedures
- Procedures: `agents-sync`, `build`, `publish-plan`, `audit-spec`, `audit-impl`, `audit-agents`, `audit-spec-to-impl`, `audit-impl-to-spec`, `draft-plan-spec-feat`, `draft-plan-spec-fix`, `draft-plan-spec-refactor`, `draft-plan-spec-chore`, `draft-plan-impl-feat`, `draft-plan-impl-fix`, `draft-plan-impl-refactor`, `draft-plan-impl-chore`
- Each procedure has fragment arrays for all 4 OODA phases
- Each procedure has display name, summary, description
- Procedures use `builtin:` prefix for embedded fragments
- Built-in procedures available without config file

**Dependencies:** T7

---

## Priority: Medium (P2) — Observability & Error Handling

### T9: Implement Structured Logging
**Description:** Structured logging per `observability.md` with configurable log level and timestamp format.

**Acceptance Criteria:**
- `internal/log/logger.go` implements structured logging
- Log levels: debug, info, warn, error
- Timestamp formats: time, time-ms, relative, iso, none
- `loop.log_level` and `loop.log_timestamp_format` config fields
- `ROODA_LOOP_LOG_LEVEL` and `ROODA_LOOP_LOG_TIMESTAMP_FORMAT` env vars
- `--verbose` sets log_level=debug and show_ai_output=true
- `--quiet` sets log_level=warn
- `--log-level` flag overrides config
- Iteration progress logged at info level
- Iteration statistics logged at completion

**Dependencies:** T6

---

### T10: Implement Error Handling & Retry Logic
**Description:** Error detection, reporting, and recovery per `error-handling.md`.

**Acceptance Criteria:**
- Transient failures (network, rate limits) distinguished from permanent failures
- Configurable retry logic with exponential backoff
- Clear error messages with actionable guidance
- Exit codes: 0 (success), 1 (aborted), 2 (max-iters), 130 (interrupted)
- Timeout handling for AI CLI execution
- Graceful degradation when non-critical components fail
- Unit tests for retry logic, timeout, error classification

**Dependencies:** T6, T9

---

### T11: Implement Dry-Run Mode
**Description:** Dry-run validation per `cli-interface.md` and `iteration-loop.md`.

**Acceptance Criteria:**
- `--dry-run` flag validates config without executing
- Validates: config file syntax, procedure exists, prompt files exist, AI command binary exists
- Displays: assembled prompt, resolved configuration with provenance
- Exit code 0 if validation passes, 1 if validation fails
- No AI CLI execution in dry-run mode
- Provenance shows source tier for each config value

**Dependencies:** T3, T5, T9

---

## Priority: Medium (P2) — Distribution

### T12: Implement Single Binary Build
**Description:** Build single binary with embedded prompts per `distribution.md`.

**Acceptance Criteria:**
- `go build -o rooda cmd/rooda/main.go` produces single binary
- Binary runs on macOS and Linux
- Embedded fragments accessible via `builtin:` prefix
- No external dependencies (yq, jq, etc.) required for core functionality
- Binary size reasonable (<20MB)
- `./rooda --version` works without config file

**Dependencies:** T1, T7, T8

---

### T13: Create Installation Documentation
**Description:** Document installation and setup per `distribution.md`.

**Acceptance Criteria:**
- `README.md` has installation instructions
- Homebrew formula documented (if applicable)
- Binary download instructions
- Zero-config startup documented
- AI command configuration documented
- Examples for common use cases

**Dependencies:** T12

---

## Priority: Low (P3) — Polish

### T14: Implement AGENTS.md Bootstrap
**Description:** Bootstrap algorithm per `operational-knowledge.md` to create/update AGENTS.md.

**Acceptance Criteria:**
- Detects work tracking system (.beads/, .github/, etc.)
- Detects build system (go.mod, package.json, etc.)
- Runs build/test/lint commands to verify they work
- Creates AGENTS.md from template if missing
- Updates AGENTS.md in-place when drift detected
- Operational Learnings section updated with verification results

**Dependencies:** T6, T8

---

### T15: Implement Provenance Tracking
**Description:** Track configuration provenance per `configuration.md`.

**Acceptance Criteria:**
- `Config.Provenance` map tracks source tier for each setting
- Provenance displayed in dry-run mode
- Provenance displayed with `--verbose` flag
- Shows: setting name, value, source tier (built-in/global/workspace/env/cli), file path (if applicable)

**Dependencies:** T2, T11

---

### T16: Write Integration Tests
**Description:** End-to-end tests for common workflows.

**Acceptance Criteria:**
- Test: `rooda build` with minimal config
- Test: `rooda --dry-run` validation
- Test: Three-tier config merge
- Test: AI CLI execution with mock
- Test: Iteration loop termination conditions
- Test: Fragment loading and composition
- Tests run in CI/CD pipeline

**Dependencies:** T12

---

## Notes

**Implementation Strategy:**
- Build foundation first (T1-T3) — can't do anything without config and CLI
- Core loop next (T4-T6) — this is the heart of the system
- Procedures (T7-T8) — makes the loop useful
- Polish last (T9-T16) — improves UX but not blocking

**Bash Implementation:**
- Keep `rooda.sh` working during Go development
- Archive bash implementation when Go version reaches feature parity
- Use bash version for dogfooding (building the Go version)

**Testing Strategy:**
- Unit tests for each component (config, CLI, loop, fragments)
- Integration tests for end-to-end workflows
- Manual testing with real AI CLI tools

**Estimated Effort:**
- Foundation (T1-T3): 3-5 build iterations
- Core Loop (T4-T6): 5-8 build iterations
- Procedures (T7-T8): 2-3 build iterations
- Observability (T9-T11): 3-4 build iterations
- Distribution (T12-T13): 1-2 build iterations
- Polish (T14-T16): 3-5 build iterations
- **Total: 17-27 build iterations**

**Risk Factors:**
- Go text/template complexity for fragment parameters
- Cross-platform path resolution (Windows vs Unix)
- Embedding 55 fragments without bloating binary
- AI CLI integration edge cases (crashes, timeouts, partial output)
