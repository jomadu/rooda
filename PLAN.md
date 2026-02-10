# Gap Analysis: v2 Go Implementation

**Generated:** 2026-02-09
**Status:** Draft

## Summary

All 11 v2 specifications are complete with JTBD structure, acceptance criteria, and examples. The current implementation is bash-based (v0.1.0) with 25 prompt files and 9 procedures. No Go implementation exists yet (no go.mod, no *.go files, no cmd/ or internal/ directories).

**Gap:** Complete v2 Go implementation required per specifications.

## Specifications Status

✅ **Complete (11/11):**
- cli-interface.md
- iteration-loop.md
- procedures.md
- ai-cli-integration.md
- configuration.md
- agents-md-format.md
- observability.md
- error-handling.md
- distribution.md
- operational-knowledge.md
- prompt-composition.md

## Implementation Status

**Current (v0.1.0 bash):**
- rooda.sh — Main OODA loop script
- rooda-config.yml — 9 procedure definitions
- prompts/*.md — 25 OODA prompt component files
- No automated tests
- No Go code

**Missing (v2 Go):**
- No go.mod
- No cmd/ directory
- No internal/ directory
- No *.go files
- No embedded prompts
- No built-in procedures
- No fragment-based composition system

## Ready Work Items (from bd)

10 tasks queued, priority 0-2:

**P0 (Critical):**
1. ralph-wiggum-ooda-hrz0: Embed default prompts and procedures

**P1 (High):**
2. ralph-wiggum-ooda-gm29: Implement basic iteration loop
3. ralph-wiggum-ooda-sont: Implement promise signal scanning
4. ralph-wiggum-ooda-pn2x: Implement failure tracking
5. ralph-wiggum-ooda-ac1c: Implement iteration timeouts
6. ralph-wiggum-ooda-07xa: Implement signal handling

**P2 (Medium):**
7. ralph-wiggum-ooda-xdf7: Implement dry-run mode
8. ralph-wiggum-ooda-jobd: Implement iteration statistics
9. ralph-wiggum-ooda-hur3: Implement context injection
10. ralph-wiggum-ooda-rayx: Implement provenance display

## Gap Analysis: Specified vs Implemented

### Phase 1: Foundation (P0)

**Task: ralph-wiggum-ooda-hrz0 (Embed default prompts and procedures)**
- Spec: procedures.md, prompt-composition.md
- Status: NOT STARTED
- Gap: No Go project structure, no go:embed, no built-in procedures
- Blocks: All other tasks (need foundation first)
- Acceptance: Binary contains embedded prompts, rooda --list-procedures shows all 16

### Phase 2: Core Loop (P1)

**Task: ralph-wiggum-ooda-gm29 (Basic iteration loop)**
- Spec: iteration-loop.md
- Status: NOT STARTED
- Gap: No loop implementation, no IterationState, no termination logic
- Depends: hrz0 (need embedded prompts)
- Acceptance: rooda build --max-iterations 3 runs 3 iterations and exits

**Task: ralph-wiggum-ooda-sont (Promise signal scanning)**
- Spec: iteration-loop.md, error-handling.md
- Status: NOT STARTED
- Gap: No output scanning, no outcome matrix
- Depends: gm29 (need loop first)
- Acceptance: Loop terminates when AI outputs <promise>SUCCESS</promise>

**Task: ralph-wiggum-ooda-pn2x (Failure tracking)**
- Spec: error-handling.md
- Status: NOT STARTED
- Gap: No ConsecutiveFailures counter, no abort logic
- Depends: sont (need outcome detection)
- Acceptance: Loop aborts after 3 consecutive failures

**Task: ralph-wiggum-ooda-ac1c (Iteration timeouts)**
- Spec: iteration-loop.md, error-handling.md
- Status: NOT STARTED
- Gap: No timeout handling, no process killing
- Depends: gm29 (need loop first)
- Acceptance: Loop kills AI CLI after configured timeout

**Task: ralph-wiggum-ooda-07xa (Signal handling)**
- Spec: iteration-loop.md, error-handling.md
- Status: NOT STARTED
- Gap: No SIGINT/SIGTERM handlers
- Depends: gm29 (need loop first)
- Acceptance: Ctrl+C kills AI CLI cleanly, exits with code 130

### Phase 3: Enhanced Features (P2)

**Task: ralph-wiggum-ooda-xdf7 (Dry-run mode)**
- Spec: cli-interface.md, iteration-loop.md
- Status: NOT STARTED
- Gap: No validation, no prompt display
- Depends: hrz0, gm29 (need prompts and loop)
- Acceptance: rooda build --dry-run validates and displays prompt without executing

**Task: ralph-wiggum-ooda-jobd (Iteration statistics)**
- Spec: iteration-loop.md
- Status: NOT STARTED
- Gap: No Welford's algorithm, no stats display
- Depends: gm29 (need loop first)
- Acceptance: Loop displays 'Iteration timing: count=N min=Xs max=Xs mean=Xs stddev=Xs'

**Task: ralph-wiggum-ooda-hur3 (Context injection)**
- Spec: cli-interface.md, prompt-composition.md
- Status: NOT STARTED
- Gap: No --context flag, no context section in prompts
- Depends: hrz0 (need prompt composition)
- Acceptance: rooda build --context 'focus on auth' injects context into prompt

**Task: ralph-wiggum-ooda-rayx (Provenance display)**
- Spec: configuration.md, observability.md
- Status: NOT STARTED
- Gap: No provenance tracking, no display
- Depends: hrz0 (need config system)
- Acceptance: Dry-run shows 'max_iterations: 10 (from: workspace config)'

## Missing from Work Tracking

The following spec features have no corresponding work items:

### CLI Interface (cli-interface.md)
- Flag parsing (--help, --version, --list-procedures)
- OODA phase overrides (--observe, --orient, --decide, --act)
- AI command resolution (--ai-cmd, --ai-cmd-alias)
- Exit code handling (0, 1, 2, 3, 130)
- Short flags (-v, -q, -n, -u, -d, -c)

### Configuration (configuration.md)
- Three-tier config loading (built-in > global > workspace)
- Global config directory resolution (ROODA_CONFIG_HOME, XDG_CONFIG_HOME)
- Config merging with provenance
- Environment variable mapping (ROODA_*)
- Config validation
- AI command alias system

### AI CLI Integration (ai-cli-integration.md)
- AI command execution
- Output capture and buffering
- Built-in aliases (kiro-cli, claude, copilot, cursor-agent)
- Shell-style command parsing
- Binary validation

### Prompt Composition (prompt-composition.md)
- Fragment loading (builtin: prefix, filesystem paths)
- Template processing (Go text/template)
- Fragment concatenation
- OODA phase assembly

### Observability (observability.md)
- Structured logging (debug, info, warn, error)
- Log timestamp formats (time, relative, iso, none)
- Progress display
- Verbose mode (AI output streaming)
- Quiet mode

### Agents.md Format (agents-md-format.md)
- Schema definition
- Section parsing
- Validation rules
- Bootstrap algorithm

### Operational Knowledge (operational-knowledge.md)
- Read-verify-update lifecycle
- Bootstrap detection
- Drift detection
- Empirical verification

### Distribution (distribution.md)
- Binary packaging
- Installation methods
- Version management
- Platform support

## Implementation Mapping

Per specifications, the Go implementation should have this structure:

```
ralph-wiggum-ooda/
├── cmd/
│   └── rooda/
│       └── main.go                    # CLI entry point
├── internal/
│   ├── ai/
│   │   ├── executor.go                # AI CLI execution
│   │   ├── resolver.go                # AI command resolution
│   │   └── aliases.go                 # Built-in aliases
│   ├── cli/
│   │   ├── parser.go                  # Flag parsing
│   │   ├── help.go                    # Help text
│   │   └── validator.go               # Flag validation
│   ├── config/
│   │   ├── config.go                  # Config types and loading
│   │   ├── defaults.go                # Built-in defaults
│   │   ├── validate.go                # Config validation
│   │   ├── provenance.go              # Provenance tracking
│   │   ├── env.go                     # Environment variables
│   │   └── merge.go                   # Config merging
│   ├── loop/
│   │   ├── loop.go                    # Core iteration loop
│   │   ├── signals.go                 # Signal handling
│   │   └── errors.go                  # Failure detection
│   ├── procedures/
│   │   ├── procedures.go              # Procedure loading
│   │   ├── fragments.go               # Fragment loading
│   │   └── builtin.go                 # Built-in procedures
│   ├── agents/
│   │   ├── schema.go                  # AGENTS.md schema
│   │   └── parser.go                  # Markdown parsing
│   └── observability/
│       ├── logger.go                  # Structured logging
│       └── stats.go                   # Statistics calculation
├── fragments/                         # Built-in fragments (embedded)
│   ├── observe/                       # 13 observe fragments
│   ├── orient/                        # 20 orient fragments
│   ├── decide/                        # 10 decide fragments
│   └── act/                           # 12 act fragments
├── go.mod
├── go.sum
└── README.md
```

## Recommended Task Breakdown

### Phase 1: Foundation (1 task, P0)
1. ✅ ralph-wiggum-ooda-hrz0: Embed default prompts and procedures

### Phase 2: Core Loop (5 tasks, P1)
2. ✅ ralph-wiggum-ooda-gm29: Implement basic iteration loop
3. ✅ ralph-wiggum-ooda-sont: Implement promise signal scanning
4. ✅ ralph-wiggum-ooda-pn2x: Implement failure tracking
5. ✅ ralph-wiggum-ooda-ac1c: Implement iteration timeouts
6. ✅ ralph-wiggum-ooda-07xa: Implement signal handling

### Phase 3: Enhanced Features (4 tasks, P2)
7. ✅ ralph-wiggum-ooda-xdf7: Implement dry-run mode
8. ✅ ralph-wiggum-ooda-jobd: Implement iteration statistics
9. ✅ ralph-wiggum-ooda-hur3: Implement context injection
10. ✅ ralph-wiggum-ooda-rayx: Implement provenance display

### Phase 4: Missing Features (10 new tasks needed)

**CLI Interface:**
11. NEW: Implement CLI flag parsing and help system
    - Spec: cli-interface.md
    - Acceptance: rooda --help displays usage, rooda --version shows version
    - Priority: P1 (blocks all CLI usage)

12. NEW: Implement OODA phase overrides
    - Spec: cli-interface.md, prompt-composition.md
    - Acceptance: rooda build --observe custom.md replaces observe phase
    - Priority: P2 (enhancement)

**Configuration:**
13. NEW: Implement three-tier config system
    - Spec: configuration.md
    - Acceptance: Global and workspace configs merge correctly
    - Priority: P0 (foundation)

14. NEW: Implement environment variable support
    - Spec: configuration.md
    - Acceptance: ROODA_LOOP_AI_CMD sets AI command
    - Priority: P1 (common use case)

**AI CLI Integration:**
15. NEW: Implement AI command execution and output capture
    - Spec: ai-cli-integration.md
    - Acceptance: AI CLI output captured and scanned for signals
    - Priority: P0 (core functionality)

**Prompt Composition:**
16. NEW: Implement fragment-based prompt composition
    - Spec: prompt-composition.md, procedures.md
    - Acceptance: Fragments load and concatenate correctly
    - Priority: P0 (core functionality)

**Observability:**
17. NEW: Implement structured logging system
    - Spec: observability.md
    - Acceptance: Logs formatted with timestamp, level, fields
    - Priority: P1 (debugging essential)

**Agents.md:**
18. NEW: Implement AGENTS.md parser and validator
    - Spec: agents-md-format.md, operational-knowledge.md
    - Acceptance: AGENTS.md sections parsed correctly
    - Priority: P2 (operational feature)

**Distribution:**
19. NEW: Implement build and packaging
    - Spec: distribution.md
    - Acceptance: Binary builds for Linux, macOS, Windows
    - Priority: P2 (distribution)

**Testing:**
20. NEW: Implement test suite
    - Spec: All specs have acceptance criteria
    - Acceptance: go test ./... passes, coverage >80%
    - Priority: P1 (quality gate)

## Critical Path

The critical path to a working v2 implementation:

1. **Foundation** (P0): hrz0, task 13 (config), task 15 (AI exec), task 16 (prompts)
2. **Core Loop** (P1): gm29, sont, pn2x, ac1c, 07xa
3. **CLI** (P1): task 11 (flag parsing), task 14 (env vars)
4. **Observability** (P1): task 17 (logging), jobd (stats)
5. **Testing** (P1): task 20 (test suite)
6. **Enhanced** (P2): xdf7, hur3, rayx, task 12 (OODA overrides)
7. **Operational** (P2): task 18 (AGENTS.md)
8. **Distribution** (P2): task 19 (packaging)

## Estimated Effort

**Phase 1 (Foundation):** 4 tasks, ~8-12 hours
- Complex: Config system, prompt composition, AI execution

**Phase 2 (Core Loop):** 5 tasks, ~6-8 hours
- Moderate: Loop logic, signal scanning, failure tracking

**Phase 3 (CLI & Observability):** 4 tasks, ~4-6 hours
- Straightforward: Flag parsing, logging, env vars

**Phase 4 (Enhanced & Operational):** 7 tasks, ~6-8 hours
- Mixed: Some simple (stats), some complex (AGENTS.md parser)

**Total:** 20 tasks, ~24-34 hours

## Next Steps

1. File 10 new tasks for Phase 4 (missing features)
2. Start with Phase 1 foundation tasks (P0)
3. Implement in dependency order per critical path
4. Run quality gates after each phase
5. Update AGENTS.md with Go-specific commands as implementation progresses

## Notes

**Why this order:**
- Foundation tasks (config, AI exec, prompts) block everything else
- Core loop is the heart of the system — must work before enhancements
- CLI and observability enable debugging and usage
- Enhanced features add value but aren't blocking
- Operational features (AGENTS.md) are nice-to-have for v2.0

**Why 20 tasks:**
- Specs are comprehensive — each major feature needs implementation
- Current 10 tasks cover ~50% of spec surface area
- Missing tasks are equally important (config, CLI, AI exec)
- Breaking into small tasks enables incremental progress

**Why no bash migration:**
- v2 is a clean rewrite, not a port
- Bash implementation archived for reference
- Go provides better error handling, testing, cross-platform support
- Fragment-based composition is fundamentally different from v1 monolithic prompts
