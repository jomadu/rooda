# Gap Analysis: Specifications vs Implementation

**Generated:** 2026-02-09
**Status:** Draft

## Executive Summary

**Current State:** v2 specifications complete (11 specs with JTBD structure, acceptance criteria, examples). v0.1.0 bash implementation exists but is v1 architecture. No Go implementation exists yet.

**Gap:** All v2 specifications are unimplemented. The project has complete architectural documentation but zero implementation of the v2 Go rewrite.

**Priority:** Start with foundational infrastructure (configuration, CLI parsing) before building the iteration loop and procedures system.

---

## Specified But Not Implemented

### Critical Path (P0) - Foundation

1. **Configuration System** (`specs/configuration.md`)
   - Three-tier config loading (built-in > global > workspace)
   - YAML parsing with go-yaml
   - Environment variable resolution
   - CLI flag precedence
   - Provenance tracking
   - AI command resolution
   - Max iterations resolution
   - **Dependencies:** None
   - **Blocks:** Everything else

2. **CLI Interface** (`specs/cli-interface.md`)
   - Argument parsing (procedure name, flags)
   - Help text generation
   - Version display
   - List procedures
   - Flag validation (mutually exclusive, constraints)
   - Exit codes (0=success, 1=user error, 2=config error, 3=execution error)
   - **Dependencies:** Configuration system
   - **Blocks:** All user-facing functionality

3. **Procedures System** (`specs/procedures.md`)
   - Fragment-based composition (55 fragments across 4 phases)
   - Fragment loading (builtin: prefix, filesystem paths, inline content)
   - Template processing (Go text/template)
   - Fragment validation at config load time
   - 16 built-in procedure definitions
   - **Dependencies:** Configuration system
   - **Blocks:** Prompt composition, iteration loop

### High Priority (P1) - Core Loop

4. **Prompt Composition** (`specs/prompt-composition.md`)
   - Assemble prompts from fragment arrays
   - Resolve fragment paths (builtin: vs filesystem)
   - Process templates with parameters
   - Concatenate fragments with double newlines
   - Inject user context (--context flag)
   - Format with section markers (# OBSERVE, # ORIENT, etc.)
   - **Dependencies:** Procedures system, configuration system
   - **Blocks:** Iteration loop

5. **AI CLI Integration** (`specs/ai-cli-integration.md`)
   - Resolve AI command from precedence chain
   - Spawn AI CLI process with prompt as stdin
   - Capture stdout/stderr to buffer
   - Stream output to terminal (--verbose)
   - Handle process timeout
   - Scan output for `<promise>` signals
   - Built-in aliases (kiro-cli, claude, copilot, cursor-agent)
   - **Dependencies:** Configuration system
   - **Blocks:** Iteration loop

6. **Error Handling** (`specs/error-handling.md`)
   - Config validation at load time (fail fast)
   - AI CLI failure detection (exit code, signals, timeout)
   - Consecutive failure tracking
   - Timeout handling (SIGTERM → SIGKILL)
   - Signal handling (SIGINT/SIGTERM)
   - Output buffer overflow handling
   - **Dependencies:** AI CLI integration
   - **Blocks:** Iteration loop

7. **Iteration Loop** (`specs/iteration-loop.md`)
   - Execute OODA cycles with fresh AI context per iteration
   - Check termination conditions (max iterations, failure threshold, SUCCESS signal)
   - Assemble prompt from OODA phases
   - Pipe prompt to AI CLI
   - Scan output for `<promise>` signals
   - Track consecutive failures
   - Calculate iteration statistics (Welford's algorithm)
   - Display progress
   - **Dependencies:** Prompt composition, AI CLI integration, error handling
   - **Blocks:** All procedures

### Medium Priority (P2) - Observability & Distribution

8. **Observability** (`specs/observability.md`)
   - Structured logging (debug, info, warn, error)
   - Log level configuration (config, env, flags)
   - Timestamp format configuration (time, relative, iso, none)
   - Progress display (iteration start/complete, timing, outcome)
   - Iteration statistics display (count, min, max, mean, stddev)
   - Dry-run mode (validate without executing)
   - Verbose mode (stream AI output, show provenance)
   - **Dependencies:** Iteration loop, configuration system
   - **Blocks:** User visibility into loop execution

9. **AGENTS.md Format** (`specs/agents-md-format.md`)
   - Schema definition (10 required sections)
   - Parsing logic (markdown sections, bold labels, code blocks)
   - Validation rules (PASS/FAIL criteria, required fields)
   - **Dependencies:** None (pure schema definition)
   - **Blocks:** Operational knowledge system

10. **Operational Knowledge** (`specs/operational-knowledge.md`)
    - Read AGENTS.md at iteration start
    - Parse into structured data (build commands, spec paths, work tracking)
    - Execute commands from AGENTS.md
    - Detect drift (expected vs actual)
    - Update AGENTS.md in-place with rationale
    - Bootstrap workflow (create AGENTS.md if missing)
    - **Dependencies:** AGENTS.md format, iteration loop
    - **Blocks:** Adaptive behavior, drift detection

11. **Distribution** (`specs/distribution.md`)
    - Single binary build with embedded prompts (go:embed)
    - Cross-compilation (macOS arm64/amd64, Linux amd64/arm64, Windows amd64)
    - Version embedding (-ldflags)
    - SHA256 checksums
    - Install script (curl | sh)
    - Homebrew formula
    - go install support
    - **Dependencies:** All implementation complete
    - **Blocks:** User installation

---

## Implementation Not Covered by Specifications

### v0.1.0 Bash Implementation (Archived)

The current `rooda.sh` (v0.1.0) is a bash implementation that:
- Implements v1 architecture (single monolithic prompts per phase)
- Uses yq for YAML parsing
- Has 9 procedures (not 16)
- Uses 25 prompt files (not 55 fragments)
- Lacks fragment-based composition
- Lacks template system
- Lacks three-tier configuration
- Lacks provenance tracking

**Status:** Archived in `archive/` directory. Not part of v2 implementation. Preserved for reference but excluded from active development per AGENTS.md.

### Prompt Files (v1 Architecture)

The `prompts/` directory contains 25 v1 prompt files:
- Single-file prompts per OODA phase (not fragment arrays)
- No template support
- No builtin: prefix system
- Organized by procedure, not by reusable fragments

**Status:** v1 architecture. v2 specs define 55 fragments organized by OODA phase for reusability. Current prompts will need restructuring into fragment-based system.

### Configuration File (v1 Schema)

The `rooda-config.yml` exists but uses v1 schema:
- Single-file prompts per phase (not fragment arrays)
- No fragment-based composition
- No template parameters
- No three-tier system (only workspace config)

**Status:** v1 schema. v2 specs define new schema with fragment arrays, template support, and three-tier merging.

---

## Structural Issues

### No Go Implementation

**Issue:** No `go.mod`, no `*.go` files, no `cmd/` or `internal/` directories exist.

**Impact:** Cannot build, test, or run v2 implementation. All 11 specs are unimplemented.

**Root Cause:** The `goify` branch restructured the project (moved files from `src/` to root, archived v1 in `archive/`) but did not create the Go implementation. Specifications were written but implementation was not started.

**Fix Required:** Bootstrap Go project structure:
- Create `go.mod` with module path
- Create `cmd/rooda/main.go` entry point
- Create `internal/` package structure
- Create `fragments/` directory with 55 embedded fragments
- Implement specs in priority order (P0 → P1 → P2)

### Prompt Organization Mismatch

**Issue:** Current `prompts/` directory has 25 v1 single-file prompts. v2 specs define 55 reusable fragments organized by OODA phase.

**Impact:** Cannot use existing prompts with v2 fragment-based composition system.

**Gap:**
- v1: `prompts/observe_plan_specs_impl.md` (single file for build procedure observe phase)
- v2: `fragments/observe/read_agents_md.md` + `fragments/observe/read_specs.md` + `fragments/observe/read_impl.md` (reusable fragments)

**Fix Required:**
1. Create `fragments/` directory structure (observe/, orient/, decide/, act/)
2. Decompose v1 prompts into reusable v2 fragments
3. Map v1 procedures to v2 fragment arrays
4. Embed fragments in Go binary via go:embed

### Configuration Schema Incompatibility

**Issue:** `rooda-config.yml` uses v1 schema (single-file prompts). v2 specs define new schema (fragment arrays with template support).

**Impact:** Existing config file will not parse with v2 implementation.

**Gap:**
- v1: `observe: prompts/observe_plan_specs_impl.md` (single string)
- v2: `observe: [{path: "builtin:fragments/observe/read_agents_md.md"}, {path: "builtin:fragments/observe/read_specs.md"}]` (array of fragment actions)

**Fix Required:**
1. Implement v2 config parser (supports fragment arrays)
2. Migrate `rooda-config.yml` to v2 schema
3. Document migration path for users

---

## Dependencies and Blockers

### Critical Path

```
Configuration System (P0)
  ↓
CLI Interface (P0) + Procedures System (P0)
  ↓
Prompt Composition (P1) + AI CLI Integration (P1)
  ↓
Error Handling (P1)
  ↓
Iteration Loop (P1)
  ↓
Observability (P2) + Operational Knowledge (P2)
  ↓
Distribution (P2)
```

### Parallel Work Opportunities

After Configuration System is complete:
- CLI Interface and Procedures System can be developed in parallel
- After Procedures System: Prompt Composition can start
- After CLI Interface: AI CLI Integration can start
- After both Prompt Composition and AI CLI Integration: Error Handling can start

### No External Blockers

All dependencies are internal to the project. No external libraries or services are blocking implementation.

---

## Recommended Implementation Order

### Phase 1: Foundation (P0)

1. **Bootstrap Go Project**
   - Create `go.mod` with module path `github.com/jomadu/rooda`
   - Create `cmd/rooda/main.go` entry point
   - Create `internal/` package structure
   - Verify `go build` produces binary

2. **Configuration System**
   - Implement `internal/config/` package
   - YAML parsing with go-yaml
   - Three-tier loading (built-in > global > workspace)
   - Environment variable resolution
   - Provenance tracking
   - Validation

3. **CLI Interface**
   - Implement `cmd/rooda/main.go` argument parsing
   - Help text generation
   - Version display
   - List procedures
   - Flag validation
   - Exit codes

4. **Procedures System**
   - Implement `internal/procedures/` package
   - Fragment loading (builtin: prefix, filesystem, inline)
   - Template processing (Go text/template)
   - Fragment validation
   - Define 16 built-in procedures

### Phase 2: Core Loop (P1)

5. **Create Fragment Files**
   - Create `fragments/` directory structure
   - Decompose v1 prompts into 55 v2 fragments
   - Organize by OODA phase (observe/, orient/, decide/, act/)
   - Embed via go:embed

6. **Prompt Composition**
   - Implement `internal/prompt/` package
   - Assemble prompts from fragment arrays
   - Resolve fragment paths
   - Process templates
   - Inject user context

7. **AI CLI Integration**
   - Implement `internal/ai/` package
   - Resolve AI command
   - Spawn process with prompt as stdin
   - Capture output
   - Stream to terminal (--verbose)
   - Scan for `<promise>` signals

8. **Error Handling**
   - Implement `internal/loop/errors.go`
   - Config validation at load time
   - AI CLI failure detection
   - Consecutive failure tracking
   - Timeout handling
   - Signal handling

9. **Iteration Loop**
   - Implement `internal/loop/loop.go`
   - Execute OODA cycles
   - Check termination conditions
   - Track consecutive failures
   - Calculate statistics
   - Display progress

### Phase 3: Observability & Distribution (P2)

10. **Observability**
    - Implement `internal/observability/` package
    - Structured logging
    - Log level configuration
    - Timestamp format configuration
    - Progress display
    - Statistics display
    - Dry-run mode
    - Verbose mode

11. **Operational Knowledge**
    - Implement `internal/agents/` package
    - Parse AGENTS.md
    - Execute commands from AGENTS.md
    - Detect drift
    - Update AGENTS.md in-place
    - Bootstrap workflow

12. **Distribution**
    - Cross-compilation script
    - Version embedding
    - SHA256 checksums
    - Install script
    - Homebrew formula
    - Documentation

---

## Quality Criteria Gaps

### Specifications

**All specs have required sections:** ✅ PASS
- All 11 specs have "Job to be Done" section
- All 11 specs have "Acceptance Criteria" section
- All 11 specs have "Examples" section
- No broken cross-references between specs

### Implementation

**All quality criteria fail:** ❌ FAIL (no Go implementation exists)
- No `go.mod` file
- No `*.go` files
- No `cmd/` or `internal/` directories
- Cannot run `go build`
- Cannot run `go test`
- No linter configuration

**Refactoring trigger:** All quality criteria fail → requires full implementation from scratch.

---

## Acceptance Criteria

This gap analysis is complete when:
- ✅ All specified features identified (11 specs analyzed)
- ✅ All unspecified implementation identified (v1 bash, v1 prompts, v1 config)
- ✅ All structural issues identified (no Go implementation, prompt mismatch, config incompatibility)
- ✅ Dependencies and blockers documented
- ✅ Implementation order recommended (P0 → P1 → P2)
- ✅ Quality criteria gaps documented

---

## Notes

### Why Start with Configuration System?

Configuration is the foundation. CLI parsing, procedures, prompt composition, and AI CLI integration all depend on resolved configuration. Building it first unblocks parallel work on CLI and procedures.

### Why Decompose v1 Prompts into v2 Fragments?

v2 fragment-based composition enables reusability. Instead of duplicating "read AGENTS.md" instructions across 16 procedures, it's a single fragment used by all. This reduces maintenance burden and ensures consistency.

### Why Three-Tier Configuration?

Enables zero-config startup (built-in defaults), personal preferences (global config), and project-specific customization (workspace config) without conflict. Mirrors established tooling conventions (Git, npm, EditorConfig).

### Why Go Instead of Bash?

Bash implementation (v0.1.0) validated the core concept but has limitations:
- No error handling
- No structured logging
- Platform-specific behavior
- No testability
- No cross-compilation

Go provides: testable code, cross-platform binary, structured error handling, signal handling, and ability to embed default prompts.

### Why 55 Fragments Instead of 25 Prompts?

v1 used 25 single-file prompts (one per procedure per phase). v2 uses 55 reusable fragments that compose into 16 procedures. This increases the fragment count but dramatically reduces duplication and improves maintainability.

**Example:**
- v1: 9 procedures × 4 phases = 36 prompt files (some shared, some duplicated)
- v2: 55 fragments → 16 procedures (fragments reused across procedures)

The fragment count is higher, but each fragment is smaller, more focused, and reusable.
