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
bd create --title "Title" --description "Desc" --deps blocks:issue-id --priority 2  # Priority range: 0-4 or P0-P4
```

## Story/Bug Input

Stories and bugs are documented in `TASK.md` at project root.

## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues

## Build/Test/Lint Commands

**Framework Dependencies:**
- yq (required) - YAML parsing for rooda-config.yml
- kiro-cli (default, configurable) - AI CLI tool, can substitute with claude-cli, aider, etc.

**Project Dependencies:**
- bd (optional) - Work tracking system used in this project, not framework-required

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

**Location:** `src/rooda.sh`, `src/prompts/*.md`, `docs/*.md`, and `scripts/*.sh`

**Patterns:**
- `src/rooda.sh` - Main loop script (root `rooda.sh` is wrapper for framework development)
- `src/rooda-config.yml` - Procedure configuration
- `src/prompts/*.md` - OODA prompt components (25 files: observe, orient, decide, act variants)
- `docs/*.md` - User-facing documentation (4 files)
- `scripts/*.sh` - Utility scripts (audit-links.sh, validate-prompts.sh)

**Exclude:**
- `.beads/*` (work tracking database)
- `specs/*` (specifications)
- `README.md`, `AGENTS.md`, `PLAN.md`, `TASK.md`, `LICENSE.md` (root files)

## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)
- All specs have "Examples" section (PASS/FAIL)
- All command examples in specs are verified working (PASS/FAIL)  # Verification process: execute commands, validate output; distinguish executable vs pseudocode
- No specs marked as DEPRECATED without replacement (PASS/FAIL)

**For implementation:**
- shellcheck passes with no errors (PASS/FAIL)
- All procedures in config have corresponding component files (PASS/FAIL)
- All prompt files follow structure per component-authoring.md (PASS/FAIL)  # Verify with: ./scripts/validate-prompts.sh
- Script executes bootstrap procedure successfully (PASS/FAIL)
- Script executes on macOS without errors (PASS/FAIL)
- Script executes on Linux without errors (PASS/FAIL)

**For documentation:**
- All code examples in docs/ are verified working (PASS/FAIL)
- Documentation matches script behavior (PASS/FAIL)
- All cross-document links work correctly (PASS/FAIL)  # Verify with: ./scripts/audit-links.sh (checks internal relative paths and external URLs with 10s timeout)
- Each procedure has usage examples (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts script behavior
- Script fails on documented use cases

**Note:** Quality criteria evolved from subjective assessments to boolean PASS/FAIL checks for automated verification.

## Operational Learnings

**Last Bootstrap Verification:** 2026-02-04T20:23:34-08:00

**Verified Working:**
- shellcheck src/rooda.sh executes without errors (clean pass)
- bd ready --json returns valid JSON with issue list
- ./scripts/audit-links.sh validates all cross-document links
- ./scripts/validate-prompts.sh confirms all 25 prompt files valid
- All commands in AGENTS.md tested and functional
- All short flags work correctly: -o, -r, -d, -a, -m, -c, -h  # Fixed -m 0 to support unlimited iterations (was being overridden by config default)
- Repository structure matches documented patterns

**Why These Definitions:**
- Specs location chosen because project uses JTBD-based markdown specifications in dedicated directory
- Implementation includes prompts/ because OODA prompt components are core framework logic
- Quality criteria are boolean PASS/FAIL to enable automated verification via scripts
- Work tracking uses beads because it provides JSON output for programmatic access
- Planning system uses PLAN.md for draft convergence before publishing to work tracking
- MAX_ITERATIONS uses empty string for "not set" to distinguish from explicit 0 (unlimited)  # Enables three-tier default: CLI flag > config default > 0

