# AGENTS.md Format

AGENTS.md is the operational guide that tells AI agents how to interact with your repository. It documents build/test/lint commands, work tracking system, file patterns, and quality criteria.

## Purpose

AGENTS.md serves as the interface between AI agents and your project. It answers:
- How do I run tests?
- How do I build the project?
- Where are the specs?
- Where is the implementation?
- What work tracking system is used?
- What are the quality criteria?

## Lifecycle

### Bootstrap

Create AGENTS.md for a new repository:

```bash
rooda bootstrap --ai-cmd-alias kiro-cli
```

This analyzes the repository and generates AGENTS.md with detected:
- Build system (Makefile, package.json, go.mod, etc.)
- Test system (go test, npm test, pytest, etc.)
- Lint system (golangci-lint, eslint, pylint, etc.)
- Spec patterns (specs/, docs/, README.md)
- Implementation patterns (src/, lib/, cmd/, internal/)
- Work tracking system (beads, GitHub Issues, TASKS.md)

### Sync

Update AGENTS.md when repository structure changes:

```bash
rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1
```

This detects drift between AGENTS.md and actual repository state, then updates AGENTS.md to fix:
- Commands that no longer work
- Paths that no longer exist
- Patterns that no longer match

### Audit

Verify AGENTS.md is accurate:

```bash
rooda audit-agents --ai-cmd-alias kiro-cli
```

This runs all commands and checks all paths, producing an audit report with drift detections.

## Format

AGENTS.md is structured markdown with specific sections that agents parse.

### Required Sections

#### Issue Tracking

Documents the work tracking system.

```markdown
## Issue Tracking

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference

\`\`\`bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
\`\`\`
```

#### Work Tracking System

Detailed commands for querying and updating work.

```markdown
## Work Tracking System

**System:** beads (bd CLI)

**Query ready work:**
\`\`\`bash
bd ready --json
\`\`\`

**Update status:**
\`\`\`bash
bd update <id> --status in_progress
\`\`\`

**Close issue:**
\`\`\`bash
bd close <id> --reason "Completed X"
\`\`\`

**Create issue:**
\`\`\`bash
bd create --title "Title" --description "Desc" --priority 2
\`\`\`
```

#### Build/Test/Lint Commands

How to run quality gates.

```markdown
## Build/Test/Lint Commands

**Dependencies:**
- Go >= 1.24.5 (required for v2 Go implementation)
- yq >= 4.0.0 (required for YAML parsing)

**Test:**
\`\`\`bash
go test ./...                    # Run all tests
go test -v ./internal/...        # Run with verbose output
\`\`\`

**Build:**
\`\`\`bash
go build -o bin/rooda ./cmd/rooda
# or use build script:
./scripts/build.sh
\`\`\`

**Lint:**
\`\`\`bash
go vet ./...         # Built-in Go linter
\`\`\`
```

#### Specification Definition

Where specs live and what format they use.

```markdown
## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure (see `specs/README.md` for index and template)

**Exclude:** `specs/README.md` (index file, not a spec)

**Current state:** v2 Go rewrite specs complete — all 11 specs written with JTBD structure, acceptance criteria, examples.
```

#### Implementation Definition

Where code lives and what patterns to follow.

```markdown
## Implementation Definition

**Location:** `cmd/`, `internal/`, `docs/`

**Patterns:**
- `cmd/rooda/main.go` — Go binary entry point
- `internal/prompt/` — Prompt composition and fragment loading
- `internal/config/` — Configuration loading and validation
- `internal/loop/` — OODA loop orchestration
- `docs/*.md` — User documentation
- `README.md` — Project overview

**Exclude:**
- `archive/` — archived v1 implementation
- `.beads/` — work tracking database
- `specs/` — specifications
- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files
```

#### Quality Criteria

Pass/fail criteria for specs and implementation.

```markdown
## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)
- All specs have "Examples" section (PASS/FAIL)
- No broken cross-references between specs (PASS/FAIL)

**For implementation:**
- `go test ./...` passes all tests (PASS/FAIL)
- `go build -o bin/rooda ./cmd/rooda` succeeds (PASS/FAIL)
- Documentation examples execute successfully (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts actual behavior
- Referenced files or paths don't exist
```

### Optional Sections

#### Planning System

Where draft plans are written and how they're published.

```markdown
## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues
```

#### Story/Bug Input

Where task input is documented for planning procedures.

```markdown
## Story/Bug Input

Stories and bugs are documented in `TASK.md` at project root. Create this file before running `draft-plan-story-to-spec` or `draft-plan-bug-to-spec` procedures.
```

#### Operational Learnings

Accumulated knowledge about what works and what doesn't.

```markdown
## Operational Learnings

**Last Bootstrap Verification:** 2026-02-09

**Verified Working:**
- `go test ./...` passes all tests
- `go build -o bin/rooda ./cmd/rooda` succeeds
- `./bin/rooda --list-procedures` works

**Verified Not Working / Missing:**
- shellcheck not installed on this machine
- No CI/CD pipeline configured

**Why These Definitions:**
- Implementation is at root level (not `src/`) because `goify` branch restructured the project
- Specs use JTBD format per `specs/README.md`
```

## Best Practices

### Keep Commands Accurate

AGENTS.md is the source of truth for how to interact with the repository. If commands change, update AGENTS.md immediately.

```bash
# After changing test command
rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1
```

### Use Inline Rationale

When documenting non-obvious choices, add brief inline comments:

```markdown
**Test:**
\`\`\`bash
go test ./...  # Using Go's built-in test runner
\`\`\`
```

### Document Dependencies

List required tools and versions:

```markdown
**Dependencies:**
- Go >= 1.24.5 (required)
- yq >= 4.0.0 (required for YAML parsing)
- shellcheck (optional, for linting bash scripts)
```

### Exclude Non-Implementation Files

Be explicit about what's NOT implementation:

```markdown
**Exclude:**
- `archive/` — archived code
- `.beads/` — work tracking database
- `specs/` — specifications
- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files
```

### Update Operational Learnings

When you discover something important, update the relevant section inline (don't append diary entries):

```markdown
## Build/Test/Lint Commands

**Test:**
\`\`\`bash
go test ./...  # Using -v flag causes output truncation issues
\`\`\`
```

## Examples

### Minimal AGENTS.md

```markdown
# Agent Instructions

## Issue Tracking

This project uses **GitHub Issues** for issue tracking.

## Work Tracking System

**System:** GitHub Issues

**Query ready work:**
\`\`\`bash
gh issue list --label "ready" --json number,title,body
\`\`\`

## Build/Test/Lint Commands

**Test:**
\`\`\`bash
npm test
\`\`\`

**Build:**
\`\`\`bash
npm run build
\`\`\`

**Lint:**
\`\`\`bash
npm run lint
\`\`\`

## Specification Definition

**Location:** `docs/*.md`

**Format:** Markdown documentation

## Implementation Definition

**Location:** `src/`

**Patterns:**
- `src/**/*.ts` — TypeScript source files
- `src/**/*.test.ts` — Test files

## Quality Criteria

**For implementation:**
- `npm test` passes all tests (PASS/FAIL)
- `npm run lint` passes with no errors (PASS/FAIL)
```

### Full AGENTS.md

See the AGENTS.md in this repository for a complete example with all sections.

## Troubleshooting

### "AGENTS.md not found"

Run bootstrap:
```bash
rooda bootstrap --ai-cmd-alias kiro-cli
```

### "Command in AGENTS.md failed"

Run sync to detect and fix drift:
```bash
rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1
```

### "Quality criteria failing"

Update quality criteria to match actual project state:
```bash
vim AGENTS.md  # Edit Quality Criteria section
```

## See Also

- [Bootstrap Procedure](procedures.md#bootstrap) - Create AGENTS.md
- [Agents-Sync Procedure](procedures.md#agents-sync) - Update AGENTS.md
- [Audit-Agents Procedure](procedures.md#audit-agents) - Verify AGENTS.md
- [Configuration](configuration.md) - rooda configuration system
