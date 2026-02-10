# Agent Instructions

## Issue Tracking

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

## Work Tracking System

**System:** beads (bd CLI)

**Query ready work:**
```bash
bd ready --json
```

**Update status:**
```bash
bd update <id> --status in_progress
```

**Close issue:**
```bash
bd close <id> --reason "Completed X"
```

**Create issue:**
```bash
bd create --title "Title" --description "Desc" --priority 2  # Priority range: 0-4 or P0-P4
```

## Story/Bug Input

Stories and bugs are documented in `TASK.md` at project root. Create this file before running `draft-plan-story-to-spec` or `draft-plan-bug-to-spec` procedures.

## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues

## Build/Test/Lint Commands

**Dependencies:**
- yq >= 4.0.0 (required) — YAML parsing for rooda-config.yml
- AI CLI tool (configurable) — default: kiro-cli, can substitute with claude, aider, cursor-agent
- bd (beads CLI) — issue tracking
- Go >= 1.24.5 (required for v2 Go implementation)

**Test:**
```bash
# v0.1.0 bash: Manual verification (no automated tests)
./rooda.sh --version
./rooda.sh --list-procedures

# v2 Go: Automated test suite
go test ./...                    # Run all tests
go test -v ./internal/...        # Run with verbose output
```

**Build:**
```bash
# v0.1.0 bash: Not required (interpreted)

# v2 Go: Build binary
go build -o bin/rooda ./cmd/rooda
# or use build script:
./scripts/build.sh
```

**Lint:**
```bash
# v0.1.0 bash
shellcheck rooda.sh  # Requires shellcheck installed; brew install shellcheck

# v2 Go
go vet ./...         # Built-in Go linter
# golangci-lint run  # Comprehensive linting (not yet configured)
```

## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure (see `specs/README.md` for index and template)

**Exclude:** `specs/README.md` (index file, not a spec)

**Current state:** v2 Go rewrite specs complete — all 11 specs written with JTBD structure, acceptance criteria, examples. Go implementation started: fragment embedding (P0.3) complete, config/AI/loop not yet implemented.

**v2 Rewrite Status:** The specifications in `specs/` describe a planned v2 Go rewrite with fundamentally different architecture (fragment-based composition, three-tier config, embedded resources, 16 procedures). This is a separate planned initiative, not missing features in v0.1.0. The v0.1.0 bash implementation (`rooda.sh`) is the current working version. When gap analysis identifies v2 features as "not implemented," this is expected—specs were written first per spec-driven development, and v2 implementation has not started.

## Implementation Definition

**Location:** `rooda.sh`, `rooda-config.yml`, `prompts/*.md` (v0.1.0 bash), `cmd/`, `internal/` (v2 Go - in progress)

**Patterns:**
- `rooda.sh` — Main OODA loop script (bash, v0.1.0)
- `rooda-config.yml` — Procedure definitions and AI tool presets
- `prompts/*.md` — 25 OODA prompt component files (observe_*, orient_*, decide_*, act_*)
- `cmd/rooda/main.go` — Go binary entry point (v2, stub only)
- `internal/prompt/` — Prompt composition and fragment loading (v2, partial)
- `internal/prompt/fragments/` — 25 embedded fragment files organized by OODA phase

**Exclude:**
- `archive/` — archived v1 implementation (preserved for reference)
- `.beads/` — work tracking database
- `specs/` — specifications
- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files

**v2 Rewrite Status:** Go implementation in progress. Completed: fragment embedding (P0.3), config loader with backward compatibility, AI executor, loop state management, error handling, observability/logging. Test suite exists with 10 test files covering all implemented packages. Binary builds successfully (`./bin/rooda --list-procedures` works). Not yet implemented: full loop integration, all 16 planned procedures. The v0.1.0 bash implementation (`rooda.sh`) remains the current working version.

## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)
- All specs have "Examples" section (PASS/FAIL)
- No broken cross-references between specs (PASS/FAIL)

**For implementation:**
- All procedures in rooda-config.yml have corresponding prompt files that exist (PASS/FAIL)
- `./rooda.sh --version` executes without errors (PASS/FAIL)
- `./rooda.sh --list-procedures` executes without errors (PASS/FAIL)
- shellcheck passes on rooda.sh with no errors (PASS/FAIL) — requires shellcheck installed
- `go test ./...` passes all tests (PASS/FAIL) — for v2 Go implementation
- `go build -o bin/rooda ./cmd/rooda` succeeds (PASS/FAIL) — for v2 Go implementation

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts actual behavior
- Referenced files or paths don't exist

## Operational Learnings

**Last Bootstrap Verification:** 2026-02-09

**Verified Working:**
- `./rooda.sh --version` returns v0.1.0
- `./rooda.sh --list-procedures` lists 9 procedures
- `bd ready --json` returns valid JSON
- All 25 prompt files in `prompts/` exist and are referenced by rooda-config.yml
- rooda-config.yml parses without errors
- All 11 v2 specs complete with JTBD structure, acceptance criteria, examples
- v2 Go binary `./bin/rooda --list-procedures` works with v0.1.0 config format
- Config loader supports both v0.1.0 string format and v2 array format (backward compatible)
- Go test suite exists with 10 test files across internal/ packages
- `go build -o bin/rooda ./cmd/rooda` succeeds

**Verified Not Working / Missing:**
- shellcheck not installed on this machine — `shellcheck rooda.sh` cannot run
- No CI/CD pipeline configured
- v2 Go implementation partial — fragment embedding complete (P0.3), config/AI/loop not implemented

**Why These Definitions:**
- Implementation is at root level (not `src/`) because `goify` branch restructured the project
- Specs use JTBD format per `specs/README.md` — v2 Go rewrite follows jobs-to-be-done methodology
- Archive preserved for reference but excluded from active implementation — prevents agents from modifying deprecated code
- Quality criteria are boolean PASS/FAIL for clear automated verification
- No build step because current implementation is bash (v2 Go will need `go build`, `go test`, `golangci-lint`)
