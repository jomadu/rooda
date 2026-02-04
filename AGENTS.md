# Agent Instructions

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

**Create issue with dependencies:**
```bash
bd create --title "Title" --description "Desc" --deps blocks:issue-id --priority 2
```

## Story/Bug Input

Stories and bugs are documented in `TASK.md` at project root.

## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues

## Build/Test/Lint Commands

**Test:** Manual verification (no automated tests)

**Build:** Not required (bash scripts are interpreted)

**Lint:**
```bash
shellcheck src/rooda.sh
```

**Verification:**
```bash
./src/rooda.sh bootstrap --max-iterations 1
bd ready --json
```

## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure

**Exclude:** `specs/README.md`, `specs/TEMPLATE.md`, `specs/specification-system.md`

## Implementation Definition

**Location:** `src/rooda.sh`, `src/components/*.md`, and `docs/*.md`

**Patterns:**
- `src/rooda.sh` - Main loop script
- `src/rooda-config.yml` - Procedure configuration
- `src/components/*.md` - OODA prompt components
- `docs/*.md` - User-facing documentation

**Exclude:**
- `.beads/*` (work tracking database)
- `specs/*` (specifications)
- `README.md`, `AGENTS.md`, `PLAN.md`, `TASK.md`, `LICENSE.md` (root files)

## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)
- All specs have "Examples" section (PASS/FAIL)
- All command examples in specs are verified working (PASS/FAIL)
- No specs marked as DEPRECATED without replacement (PASS/FAIL)

**For implementation:**
- shellcheck passes with no errors (PASS/FAIL)
- All procedures in config have corresponding component files (PASS/FAIL)
- Script executes bootstrap procedure successfully (PASS/FAIL)
- Script executes on macOS without errors (PASS/FAIL)
- Script executes on Linux without errors (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts script behavior
- Script fails on documented use cases

## Operational Learnings

**2025-02-03:** Updated quality criteria from subjective assessments to boolean PASS/FAIL checks. Previous criteria ("Clarity: Can a new user understand?", "Maintainability: Is bash code readable?") were not actionable for automated quality assessment. New criteria provide clear thresholds that can be verified empirically.

