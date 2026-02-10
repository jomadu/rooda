# rooda v2 Implementation Plan

**Status:** Draft  
**Created:** 2026-02-09  
**Target:** Minimal viable v2 Go implementation with single working procedure

## Priority 0: Foundation (Blocking Everything)

### P0.1: Project Bootstrap
- [ ] Create `go.mod` with module `github.com/jomadu/rooda`
- [ ] Create directory structure: `cmd/rooda/`, `internal/{config,loop,prompt,ai,agents}/`
- [ ] Add dependencies: `gopkg.in/yaml.v3`, `github.com/kballard/go-shellquote`
- [ ] Create `.gitignore` for Go artifacts

**Acceptance:** `go mod tidy` succeeds, directory structure exists

### P0.2: Version and Build Metadata
- [ ] Create `cmd/rooda/main.go` with version variables (Version, CommitSHA, BuildDate)
- [ ] Implement `--version` flag that prints version info
- [ ] Create `scripts/build.sh` with ldflags injection

**Acceptance:** `go run cmd/rooda/main.go --version` prints version

### P0.3: Embed Default Fragments
- [ ] Copy `prompts/*.md` to `fragments/` with OODA phase subdirectories
- [ ] Reorganize into `fragments/{observe,orient,decide,act}/` structure per procedures.md
- [ ] Add `//go:embed fragments` directive in `internal/prompt/embed.go`
- [ ] Implement `LoadFragment(path)` that handles `builtin:` prefix

**Acceptance:** Embedded fragments accessible at runtime, no external files needed

## Priority 1: Configuration System (Blocks Execution)

### P1.1: Core Config Types
- [ ] Define `Config`, `LoopConfig`, `Procedure`, `FragmentAction` structs in `internal/config/types.go`
- [ ] Define `LogLevel`, `TimestampFormat`, `IterationMode` enums
- [ ] Define built-in defaults as constants

**Acceptance:** Types compile, built-in defaults defined

### P1.2: YAML Parsing and Merging
- [ ] Implement `LoadConfig(cliFlags)` in `internal/config/loader.go`
- [ ] Implement three-tier merge: built-in → global → workspace → env vars
- [ ] Implement global config directory resolution (ROODA_CONFIG_HOME > XDG_CONFIG_HOME > ~/.config/rooda)
- [ ] Implement provenance tracking

**Acceptance:** Config loads from all tiers, precedence correct, provenance tracked

### P1.3: Config Validation
- [ ] Implement `ValidateConfig(config)` in `internal/config/validate.go`
- [ ] Validate required fields, type constraints, file existence
- [ ] Validate AI command binary exists and is executable
- [ ] Return clear error messages with file paths and suggestions

**Acceptance:** Invalid configs rejected with actionable errors

### P1.4: AI Command Resolution
- [ ] Implement `ResolveAICommand(config, procedure, flags)` in `internal/config/resolver.go`
- [ ] Implement built-in aliases map (kiro-cli, claude, copilot, cursor-agent)
- [ ] Implement precedence chain per ai-cli-integration.md
- [ ] Return error with guidance if no AI command configured

**Acceptance:** AI command resolved correctly, aliases work, clear error when missing

## Priority 2: Prompt Composition (Blocks Execution)

### P2.1: Fragment Loading
- [ ] Implement `LoadFragment(path, configDir, embeddedFS)` in `internal/prompt/loader.go`
- [ ] Handle `builtin:` prefix for embedded resources
- [ ] Handle relative paths from config directory
- [ ] Handle inline content from FragmentAction.Content

**Acceptance:** Fragments load from embedded FS and filesystem

### P2.2: Template Processing
- [ ] Implement `ProcessTemplate(content, params)` in `internal/prompt/template.go`
- [ ] Use Go text/template for parameterized fragments
- [ ] Return clear errors for template syntax issues

**Acceptance:** Templates execute with parameters, errors clear

### P2.3: Prompt Assembly
- [ ] Implement `AssemblePrompt(procedure, userContext, configDir)` in `internal/prompt/composer.go`
- [ ] Concatenate fragments per phase with double newlines
- [ ] Add section markers (# OBSERVE, # ORIENT, # DECIDE, # ACT)
- [ ] Inject user context at top if provided

**Acceptance:** Assembled prompt matches format in prompt-composition.md examples

## Priority 3: AI CLI Integration (Blocks Execution)

### P3.1: Process Execution
- [ ] Implement `ExecuteAICLI(aiCmd, prompt, verbose, timeout, maxBuffer)` in `internal/ai/executor.go`
- [ ] Spawn process with prompt as stdin
- [ ] Capture stdout/stderr to buffer with size limit
- [ ] Stream to terminal if verbose=true
- [ ] Handle timeout with SIGTERM → SIGKILL escalation
- [ ] Return `AIExecutionResult` with output, exit code, duration, truncated flag

**Acceptance:** AI CLI executes, output captured, timeout works, verbose streams

### P3.2: Signal Scanning
- [ ] Implement `ScanOutputForSignals(output)` in `internal/ai/executor.go`
- [ ] Simple string matching for `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>`
- [ ] Return both flags (caller decides precedence)

**Acceptance:** Signals detected correctly, case-sensitive exact match

## Priority 4: Iteration Loop (Core Functionality)

### P4.1: Loop State Management
- [ ] Define `IterationState`, `IterationStats` structs in `internal/loop/state.go`
- [ ] Implement Welford's algorithm for statistics in `updateStats()`
- [ ] Implement `getMean()`, `getStdDev()` helpers

**Acceptance:** Statistics calculated correctly with constant memory

### P4.2: Failure Detection
- [ ] Implement `DetectIterationFailure(result)` in `internal/loop/errors.go`
- [ ] Implement outcome matrix per iteration-loop.md
- [ ] Promise signals override exit code
- [ ] FAILURE takes precedence over SUCCESS

**Acceptance:** Failure detection matches outcome matrix

### P4.3: Core Loop
- [ ] Implement `RunLoop(state, config)` in `internal/loop/loop.go`
- [ ] Check termination conditions (max iterations, failure threshold)
- [ ] Assemble prompt, execute AI CLI, scan signals
- [ ] Track consecutive failures, reset on success
- [ ] Update statistics after each iteration
- [ ] Return loop status (success, max-iters, aborted, interrupted)

**Acceptance:** Loop executes iterations, terminates correctly, statistics displayed

### P4.4: Signal Handling
- [ ] Implement SIGINT/SIGTERM handlers in `internal/loop/signals.go`
- [ ] Kill AI CLI process on signal
- [ ] Wait for termination with 5s timeout
- [ ] Exit with code 130

**Acceptance:** Ctrl+C kills AI CLI and exits cleanly

## Priority 5: CLI Interface (User-Facing)

### P5.1: Argument Parsing
- [ ] Implement flag parsing in `cmd/rooda/main.go`
- [ ] Support all flags from cli-interface.md
- [ ] Validate mutually exclusive flags (--verbose/--quiet, --max-iterations/--unlimited)
- [ ] Parse repeatable flags (--context, --observe, --orient, --decide, --act)

**Acceptance:** All flags parse correctly, validation works

### P5.2: Help and Info Commands
- [ ] Implement `--help` with usage, flags, examples
- [ ] Implement `--list-procedures` with descriptions
- [ ] Implement procedure-specific help (`rooda <proc> --help`)

**Acceptance:** Help text clear and complete

### P5.3: Dry-Run Mode
- [ ] Implement `--dry-run` validation in `cmd/rooda/main.go`
- [ ] Validate config, prompt files, AI command
- [ ] Display assembled prompt with section markers
- [ ] Display resolved config with provenance
- [ ] Exit with code 0 (pass) or 1 (fail)

**Acceptance:** Dry-run validates without executing, output matches observability.md

## Priority 6: Observability (Essential for Debugging)

### P6.1: Structured Logging
- [ ] Implement `LogEvent` struct and logger in `internal/observability/logger.go`
- [ ] Implement log levels (debug, info, warn, error)
- [ ] Implement timestamp formats (time, relative, iso, none)
- [ ] Format as `[timestamp] LEVEL message key=value`
- [ ] Route to stderr by default

**Acceptance:** Logs formatted correctly, levels filter, timestamps configurable

### P6.2: Progress Display
- [ ] Log iteration start/complete with timing
- [ ] Log loop start/complete with status
- [ ] Display statistics at loop completion
- [ ] Suppress progress when log level > info

**Acceptance:** Progress matches examples in observability.md

## Priority 7: Single Working Procedure (Milestone)

### P7.1: Implement agents-sync Procedure
- [ ] Define `agents-sync` in built-in procedures with fragment arrays
- [ ] Create minimal fragments: `observe/read_agents_md.md`, `observe/scan_repo_structure.md`
- [ ] Create minimal fragments: `orient/compare_detected_vs_documented.md`, `orient/identify_drift.md`
- [ ] Create minimal fragments: `decide/determine_sections_to_update.md`, `decide/check_if_blocked.md`
- [ ] Create minimal fragments: `act/write_agents_md.md`, `act/commit_changes.md`, `act/emit_success.md`

**Acceptance:** `rooda agents-sync --dry-run` shows assembled prompt

### P7.2: End-to-End Test
- [ ] Run `rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1`
- [ ] Verify loop executes, AI CLI invoked, output captured
- [ ] Verify statistics displayed
- [ ] Verify exit code correct

**Acceptance:** Single procedure works end-to-end

## Priority 8: Remaining Built-in Procedures (Completeness)

### P8.1: Create All 55 Fragment Files
- [ ] Write all fragments per procedures.md fragment directory structure
- [ ] 13 observe fragments, 20 orient fragments, 10 decide fragments, 12 act fragments
- [ ] Embed in binary via go:embed

**Acceptance:** All fragments exist and embedded

### P8.2: Define All 16 Built-in Procedures
- [ ] Define all procedures in `internal/procedures/builtin.go`
- [ ] Map fragment arrays per procedures.md built-in procedures configuration

**Acceptance:** `rooda --list-procedures` shows all 16

## Priority 9: Distribution (Deployment)

### P9.1: Cross-Compilation
- [ ] Create `scripts/build.sh` for cross-platform builds
- [ ] Build for darwin/arm64, darwin/amd64, linux/amd64, linux/arm64, windows/amd64
- [ ] Generate checksums.txt

**Acceptance:** Binaries build for all platforms, checksums generated

### P9.2: Installation Script
- [ ] Create `scripts/install.sh` with platform detection
- [ ] Verify checksums before installation
- [ ] Install to /usr/local/bin (Unix) or %USERPROFILE%\bin (Windows)

**Acceptance:** Install script works on macOS and Linux

### P9.3: GitHub Release Workflow
- [ ] Create `.github/workflows/release.yml`
- [ ] Build on tag push, upload binaries to GitHub Releases
- [ ] Include install.sh and checksums.txt in release assets

**Acceptance:** Release workflow publishes binaries on tag

## Priority 10: AGENTS.md Integration (Operational Knowledge)

### P10.1: AGENTS.md Parser
- [ ] Implement `ParseAgentsMD(content)` in `internal/agents/parser.go`
- [ ] Parse sections per agents-md-format.md schema
- [ ] Return `AgentsMD` struct with all fields

**Acceptance:** AGENTS.md parsed correctly

### P10.2: Bootstrap Workflow
- [ ] Implement `BootstrapAgentsMD(repoPath)` in `internal/agents/bootstrap.go`
- [ ] Detect build system, test system, spec paths, impl paths, work tracking
- [ ] Generate AGENTS.md with detected values and rationale comments

**Acceptance:** Bootstrap creates valid AGENTS.md from scratch

### P10.3: Drift Detection and Update
- [ ] Implement `VerifyAgentsMD(agentsMD)` in `internal/agents/verifier.go`
- [ ] Run commands, check paths, detect drift
- [ ] Implement `UpdateAgentsMD(agentsMD, drifts)` in `internal/agents/updater.go`
- [ ] Update in-place with inline rationale

**Acceptance:** Drift detected and AGENTS.md updated correctly

## Priority 11: User Documentation (Part of Implementation)

### P11.1: Write Core Documentation
- [ ] Write README.md (overview, installation, quick start)
- [ ] Write docs/installation.md (all install methods, platform-specific)
- [ ] Write docs/procedures.md (all 16 procedures with examples)
- [ ] Write docs/configuration.md (three-tier system, all settings)
- [ ] Write docs/cli-reference.md (all flags, exit codes)
- [ ] Write docs/troubleshooting.md (common errors, solutions)
- [ ] Write docs/agents-md.md (format, lifecycle, bootstrap)

**Acceptance:** All 7 documentation files exist with working examples

### P11.2: Update AGENTS.md Implementation Definition
- [ ] Add docs/ to Implementation Definition patterns
- [ ] Add quality criterion: "Documentation examples execute successfully (PASS/FAIL)"
- [ ] Verify audit-impl checks documentation (already built-in)

**Acceptance:** AGENTS.md reflects that docs are implementation

## Out of Scope for v2.0.0

- Homebrew formula (manual installation sufficient for initial release)
- JSON log format (text format sufficient)
- Failure type classification (all failures treated equally)
- Exponential backoff (simple consecutive failure threshold sufficient)
- Advanced template validation (basic Go template errors sufficient)
- Fragment caching (performance not critical for v2.0.0)
- External link validation in docs (only internal links checked)
- Interactive examples in docs (asciinema recordings)
- Versioned docs (docs for each release tag)
- Search index for docs site
- API docs with Go package examples
- Diagram generation from data structures

## Success Criteria

**v2.0.0 is complete when:**
1. `rooda --version` prints version
2. `rooda --list-procedures` shows all 16 procedures
3. `rooda agents-sync --dry-run` validates and displays prompt
4. `rooda agents-sync --ai-cmd-alias kiro-cli` executes end-to-end
5. Single binary with embedded prompts (no external dependencies except AI CLI)
6. Cross-platform binaries available (macOS, Linux, Windows)
7. Installation script works
8. All acceptance criteria from specs pass
9. User documentation written (README.md + 6 docs/ files)
10. `rooda audit-impl` validates documentation (docs are implementation)

## Task Breakdown Summary

- **P0:** 3 tasks (foundation)
- **P1:** 4 tasks (configuration)
- **P2:** 3 tasks (prompt composition)
- **P3:** 2 tasks (AI CLI integration)
- **P4:** 4 tasks (iteration loop)
- **P5:** 3 tasks (CLI interface)
- **P6:** 2 tasks (observability)
- **P7:** 2 tasks (single working procedure)
- **P8:** 2 tasks (remaining procedures)
- **P9:** 3 tasks (distribution)
- **P10:** 3 tasks (AGENTS.md integration)
- **P11:** 2 tasks (user documentation as implementation)

**Total:** 33 tasks across 12 priority levels

## Implementation Order

Work sequentially through priorities P0 → P11. Within each priority, tasks can be parallelized if independent. P7 (single working procedure) is the first major milestone — everything before it is foundational, everything after is expansion.

## Estimated Effort

- **P0-P6:** ~55% of effort (core framework)
- **P7:** ~10% of effort (first working procedure validates architecture)
- **P8-P11:** ~35% of effort (completeness, polish, documentation)

## Notes

- This plan follows spec-driven development: all 11 specs are complete, implementation follows specs exactly
- No Go code exists yet — starting from scratch
- v0.1.0 bash implementation remains functional during v2 development
- Fragment content from `prompts/*.md` will be reorganized into OODA phase subdirectories
- Built-in procedures use embedded fragments; custom procedures can reference filesystem fragments
