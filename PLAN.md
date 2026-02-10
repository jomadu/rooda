# Implementation Plan: v2 Go Rewrite

**Status:** Draft  
**Created:** 2026-02-09  
**Source:** Gap analysis (specs → implementation)

## Summary

Implement the v2 Go rewrite per the 11 complete specifications. The v0.1.0 bash implementation validates the core OODA loop concept but lacks the fragment-based composition, three-tier configuration, structured logging, and distribution features specified in v2.

**Key Gap:** No Go implementation exists (no go.mod, no *.go files). All v2 specs are complete but unimplemented.

## Priority Breakdown

- **P0 (Critical):** 15 tasks — Core framework, iteration loop, configuration, prompt composition
- **P1 (High):** 5 tasks — CLI interface, distribution, installation
- **P2 (Medium):** 3 tasks — Procedures, fragments, validation
- **P3 (Low):** 5 tasks — Observability, error handling, AGENTS.md automation

**Total:** 28 tasks

## Tasks

### Phase 1: Foundation (P0)

**P0-01: Initialize Go project structure**
- Create go.mod with module name `github.com/jomadu/rooda`
- Create directory structure: cmd/rooda/, internal/config/, internal/loop/, internal/prompt/, internal/ai/, internal/procedures/, internal/agents/
- Create cmd/rooda/main.go with basic CLI entry point
- Verify: `go build ./cmd/rooda` produces binary

**P0-02: Implement configuration loading (built-in defaults)**
- Define Config, LoopConfig, Procedure, FragmentAction structs per configuration.md
- Implement built-in default configuration (embedded in code)
- Implement config validation (required fields, type constraints)
- Verify: Config struct can be instantiated with built-in defaults

**P0-03: Implement configuration loading (workspace config)**
- Implement YAML parsing for ./rooda-config.yml
- Implement config merging (workspace overrides built-in)
- Implement field-level merging for procedures
- Verify: Workspace config overrides built-in defaults correctly

**P0-04: Implement configuration loading (global config)**
- Implement global config directory resolution (ROODA_CONFIG_HOME > XDG_CONFIG_HOME/rooda > ~/.config/rooda)
- Implement YAML parsing for <config_dir>/rooda-config.yml
- Implement three-tier merging (workspace > global > built-in)
- Verify: Global config loads and merges correctly

**P0-05: Implement configuration loading (environment variables)**
- Implement ROODA_* environment variable parsing
- Implement env var overrides for loop settings
- Implement precedence: CLI flags > env vars > workspace > global > built-in
- Verify: Environment variables override config file values

**P0-06: Implement fragment-based prompt composition**
- Implement fragment loading (builtin: prefix vs filesystem paths)
- Implement fragment concatenation with double newlines
- Implement phase assembly (observe, orient, decide, act)
- Verify: Assembled prompt matches expected structure

**P0-07: Implement template processing for fragments**
- Integrate Go text/template for parameterized fragments
- Implement template execution with provided parameters
- Implement template validation at config load time
- Verify: Templates render correctly with parameters

**P0-08: Implement AI command resolution**
- Implement ResolveAICommand with precedence chain
- Implement built-in aliases (kiro-cli, claude, copilot, cursor-agent)
- Implement alias resolution from merged config
- Verify: AI command resolves correctly from all sources

**P0-09: Implement AI CLI execution**
- Implement ExecuteAICLI with stdin piping
- Implement output capture with configurable buffer size
- Implement exit code capture
- Verify: AI CLI executes and output is captured

**P0-10: Implement promise signal scanning**
- Implement ScanOutputForSignals (case-sensitive exact match)
- Implement signal precedence (FAILURE > SUCCESS)
- Implement outcome determination per iteration-loop.md matrix
- Verify: Signals detected correctly in output

**P0-11: Implement basic iteration loop**
- Implement RunLoop with iteration counter
- Implement max iterations termination check
- Implement iteration timing and progress display
- Verify: Loop executes N iterations and terminates

**P0-12: Implement consecutive failure tracking**
- Implement ConsecutiveFailures counter
- Implement failure threshold abort logic
- Implement counter reset on success
- Verify: Loop aborts after threshold consecutive failures

**P0-13: Implement iteration statistics**
- Implement IterationStats with Welford's online algorithm
- Implement statistics display at loop completion
- Implement constant memory calculation (O(1))
- Verify: Statistics (count, min, max, mean, stddev) calculated correctly

**P0-14: Implement structured logging**
- Implement log levels (debug, info, warn, error)
- Implement log format with timestamp, level, message, fields
- Implement log level configuration (config, env, flags)
- Verify: Logs display at correct levels with structured fields

**P0-15: Implement dry-run mode**
- Implement --dry-run flag
- Implement config validation (prompt files, AI command binary)
- Implement assembled prompt display
- Verify: Dry-run validates and displays prompt without executing

### Phase 2: Core Features (P0-P1)

**P1-01: Implement CLI interface (basic flags)**
- Implement flag parsing for --max-iterations, --unlimited, --dry-run, --verbose, --quiet
- Implement flag precedence over config
- Implement --help and --version flags
- Verify: Flags override config correctly

**P1-02: Implement CLI interface (OODA phase overrides)**
- Implement --observe, --orient, --decide, --act flags (repeatable)
- Implement file existence heuristic (file vs inline content)
- Implement phase array replacement (not merge)
- Verify: OODA phase overrides work correctly

**P1-03: Implement CLI interface (context injection)**
- Implement --context flag (repeatable)
- Implement file existence heuristic
- Implement context injection at top of prompt
- Verify: Context appears in assembled prompt

**P1-04: Implement iteration timeout**
- Implement timeout configuration (loop.iteration_timeout, procedure.iteration_timeout)
- Implement process killing on timeout (SIGTERM, then SIGKILL)
- Implement timeout as failure (increments ConsecutiveFailures)
- Verify: Process killed after timeout, iteration counts as failure

**P1-05: Implement output buffer management**
- Implement configurable max buffer size (loop.max_output_buffer, procedure.max_output_buffer)
- Implement truncation from beginning when exceeded
- Implement truncation warning logging
- Verify: Output truncated correctly, signals at end preserved

**P1-06: Implement verbose mode**
- Implement --verbose flag (sets show_ai_output=true and log_level=debug)
- Implement AI CLI output streaming to terminal
- Implement output capture alongside streaming
- Verify: Output streams to terminal and is captured

**P1-07: Implement cross-platform builds**
- Create build script for macOS arm64/amd64, Linux amd64/arm64, Windows amd64
- Implement version embedding via -ldflags
- Generate SHA256 checksums
- Verify: Binaries build for all platforms

**P1-08: Create installation methods**
- Create Homebrew formula
- Create install.sh script with checksum verification
- Document go install method
- Verify: Installation works via all three methods

### Phase 3: Procedures & Fragments (P2)

**P2-01: Create fragment files (observe phase)**
- Create 13 observe phase fragments per procedures.md
- Organize in fragments/observe/ directory
- Embed via go:embed
- Verify: All observe fragments loadable

**P2-02: Create fragment files (orient, decide, act phases)**
- Create 20 orient, 10 decide, 12 act phase fragments per procedures.md
- Organize in fragments/orient/, fragments/decide/, fragments/act/ directories
- Embed via go:embed
- Verify: All fragments loadable

**P2-03: Define 16 built-in procedures**
- Define all 16 procedures per procedures.md (agents-sync, build, publish-plan, 4 audit, 8 draft-plan)
- Implement procedure loading from config
- Implement procedure validation
- Verify: All 16 procedures loadable and valid

### Phase 4: Distribution & Polish (P1-P3)

**P3-01: Implement timestamp format configuration**
- Implement TimestampFormat enum (time, time-ms, relative, iso, none)
- Implement timestamp formatting in logs
- Implement configuration via loop.log_timestamp_format
- Verify: All timestamp formats work correctly

**P3-02: Implement signal handling**
- Implement SIGINT/SIGTERM handlers
- Implement AI CLI process killing on interrupt
- Implement graceful termination with timeout
- Verify: Ctrl+C kills AI CLI and exits cleanly

**P3-03: Implement crash recovery**
- Implement partial output capture from crashed processes
- Implement signal scanning on partial output
- Implement crash logging with signal name
- Verify: Crashes handled gracefully, partial output captured

**P3-04: Implement detailed error messages**
- Implement ValidationError struct with file path, line number, message, suggestion
- Implement FailureContext struct with iteration, command, exit code, duration, output preview
- Implement error messages with provenance
- Verify: Errors display clear, actionable messages

**P3-05: Implement AGENTS.md bootstrap workflow**
- Implement repository structure detection (build system, test system, spec paths, impl paths, work tracking)
- Implement AGENTS.md generation from template
- Implement bootstrap procedure
- Verify: Bootstrap creates valid AGENTS.md

### Phase 5: Testing & Documentation

**P3-06: Manual verification of all procedures**
- Test all 16 built-in procedures
- Verify quality criteria pass
- Document any issues found

**P3-07: Update README with installation and usage**
- Document installation methods (Homebrew, curl|sh, go install)
- Document basic usage examples
- Document configuration options

**P3-08: Update AGENTS.md for Go implementation**
- Update Build/Test/Lint Commands section
- Update Implementation Definition section
- Update Quality Criteria section
- Update Operational Learnings section

## Dependencies

- Go 1.21+ (for go:embed and modern stdlib)
- yq >= 4.0.0 (for YAML parsing in bash during transition)
- git (for version metadata)
- shellcheck (optional, for linting bash during transition)

## Success Criteria

- [ ] All 28 tasks completed
- [ ] All quality criteria from AGENTS.md pass
- [ ] All 16 built-in procedures work correctly
- [ ] Installation works via all three methods (Homebrew, curl|sh, go install)
- [ ] Binary runs on all target platforms (macOS arm64/amd64, Linux amd64/arm64, Windows amd64)
- [ ] Documentation updated (README, AGENTS.md)

## Notes

**Transition Strategy:**
- v0.1.0 bash implementation remains functional during Go development
- Go implementation developed in parallel on main branch
- Once Go implementation passes all quality criteria, bash archived to archive/
- No breaking changes to config file format (v0.1.0 configs work with v2)

**Config Compatibility:**
- v2 uses fragment arrays instead of single files per phase
- v0.1.0 configs can be migrated by wrapping single file paths in arrays
- Example: `observe: prompts/observe_bootstrap.md` → `observe: [{path: "prompts/observe_bootstrap.md"}]`

**Fragment Migration:**
- Existing 25 prompt files in prompts/ can be reused as fragments
- Some prompts may need to be split into multiple fragments for reusability
- New fragments created for procedures not in v0.1.0

**Bash Features to Preserve:**
- `--list-procedures` flag
- Fuzzy matching for unknown procedure names
- Dependency version checking
- Git push error parsing with guidance
- Platform detection for install instructions

**Bash Features to Drop:**
- AI tool presets (replaced by ai_cmd_aliases in config)
- Single file per phase (replaced by fragment arrays)
- Unlimited iterations via 0 (replaced by --unlimited flag)
