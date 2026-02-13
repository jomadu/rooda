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

Stories and bugs are documented in `TASK.md` at project root. Create this file before running planning procedures like `draft-plan-spec-feat`, `draft-plan-spec-fix`, `draft-plan-impl-feat`, or `draft-plan-impl-fix`.

## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues

## Build/Test/Lint Commands

**Dependencies:**
- AI CLI tool (configurable) — default: kiro-cli, can substitute with claude, aider, cursor-agent
- bd (beads CLI) — issue tracking
- Go >= 1.24.5 (required)
- make (optional but recommended) — unified build interface

**Unified Interface (via Makefile):**
```bash
make test    # Run all tests
make build   # Build Go binary
make lint    # Run all linters (go vet)
make all     # Run lint, test, and build
make clean   # Remove build artifacts
```

**Direct Commands (alternative):**
```bash
# Test
go test ./...                    # Run all tests
go test -v ./internal/...        # Run tests (verbose)

# Build
go build -o bin/rooda ./cmd/rooda  # Build binary

# Lint
go vet ./...                     # Run linter
```

## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure (see `specs/README.md` for index and template)

**Exclude:** `specs/README.md` (index file, not a spec)

**Current state:** All 11 specs complete with JTBD structure, acceptance criteria, examples. Go implementation in progress.

## Implementation Definition

**Location:** `cmd/`, `internal/`, `docs/`

**Patterns:**
- `cmd/rooda/main.go` — Go binary entry point
- `internal/prompt/` — Prompt composition and fragment loading
- `internal/prompt/fragments/` — 25 embedded fragment files organized by OODA phase
- `docs/*.md` — User documentation (installation, procedures, configuration, CLI reference, troubleshooting, AGENTS.md format)
- `README.md` — Project overview, quick start, installation

**Exclude:**
- `archive/` — archived v1 implementation (preserved for reference)
- `.beads/` — work tracking database
- `specs/` — specifications
- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files

**Implementation Status:** Core infrastructure complete (fragment embedding, config loader, AI executor, loop state management, error handling, observability/logging, full loop integration). All 16 procedure definitions exist. Test suite with 10 test files. Binary builds successfully and executes procedures end-to-end.

## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)
- All specs have "Examples" section (PASS/FAIL)
- No broken cross-references between specs (PASS/FAIL)

**For implementation:**
- `go test ./...` passes all tests (PASS/FAIL)
- `go build -o bin/rooda ./cmd/rooda` succeeds (PASS/FAIL)
- `./bin/rooda list` shows all 16 procedures (PASS/FAIL)
- Documentation examples execute successfully (PASS/FAIL) — verify commands in docs/ work as documented

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts actual behavior
- Referenced files or paths don't exist

## Operational Learnings

**Last Bootstrap Verification:** 2026-02-12

**Verified Working:**
- `bd ready --json` returns valid JSON
- All 11 specs complete with JTBD structure, acceptance criteria, examples
- Config loader supports both string format and array format (backward compatible)
- Go test suite exists with 10 test files across internal/ packages
- `go build -o bin/rooda ./cmd/rooda` succeeds
- `./bin/rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1` executes end-to-end successfully
- `go test ./...` passes all tests
- `./bin/rooda list` shows all 16 procedures available

**Verified Not Working / Missing:**
- Fragment files not yet verified — all 16 procedures defined but fragment file existence not confirmed

**Why These Definitions:**
- Implementation is at root level (not `src/`) because `goify` branch restructured the project
- Specs use JTBD format per `specs/README.md` — follows jobs-to-be-done methodology
- Archive preserved for reference but excluded from active implementation — prevents agents from modifying deprecated code
- Quality criteria are boolean PASS/FAIL for clear automated verification
