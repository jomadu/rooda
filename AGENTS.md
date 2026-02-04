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

**2026-02-03:** Quality criterion "All command examples in specs are verified working" requires verification process definition. Command examples should be empirically tested by executing them and validating output. Distinguish between executable commands (./rooda.sh, yq, bd, kiro-cli) that must work as documented, and pseudocode/illustrative examples that should be clearly marked as non-executable. Verification ensures specs remain accurate as implementation evolves.

**2026-02-03:** Quality criterion "All command examples in specs are verified working" requires verification process definition. Command examples should be empirically tested by executing them and validating output. Distinguish between executable commands (./rooda.sh, yq, bd, kiro-cli) that must work as documented, and pseudocode/illustrative examples that should be clearly marked as non-executable. Verification ensures specs remain accurate as implementation evolves.

**2026-02-03:** Quality assessment of specifications completed. Criteria 1-3 and 5 pass. Criterion 4 fails: command examples exist but verification process not defined and examples not empirically tested. Refactoring plan created in PLAN.md to address: (1) define verification process, (2) execute verification on all specs, (3) mark non-executable examples clearly, (4) create verification tracking system, (5) automate verification where possible. Priority is defining verification process and executing initial verification pass.

**2026-02-03:** Quality assessment iteration completed. Systematic evaluation of all 9 specifications against boolean criteria. Results: Criteria 1 (Job to be Done sections), 2 (Acceptance Criteria sections), 3 (Examples sections), and 5 (DEPRECATED specs have replacements) all pass. Criterion 4 (command examples verified working) fails: 44 bash code blocks identified across 8 specs, no verification process defined, no empirical testing performed, no distinction between executable commands and pseudocode. Refactoring plan written to PLAN.md with 5 prioritized tasks. Critical path: define verification process (Task 1) before executing verification pass (Task 2).

**2026-02-03:** Quality assessment re-executed. Same results as previous iteration: Criteria 1, 2, 3, and 5 pass. Criterion 4 fails (command examples not verified). Refactoring plan regenerated in PLAN.md with identical 5-task structure. No operational learnings - quality criteria remain accurate, verification process still undefined, no changes to AGENTS.md needed.

**2026-02-03:** Quality assessment iteration identified new failure: agents-md-format.md does not follow spec template structure. Missing "Job to be Done", "Acceptance Criteria", and "Examples" sections (uses "Purpose" instead). This causes criteria 1, 2, and 3 to fail. Criterion 4 (command examples verified) continues to fail - no verification process defined. Criterion 5 passes. Refactoring plan updated in PLAN.md with 5 prioritized tasks: (1) fix agents-md-format.md structure, (2) define verification process, (3) execute verification pass, (4) mark non-executable examples, (5) automate verification. Critical path: fix agents-md-format.md structure (Task 1) to achieve criteria 1-3 compliance.

**2026-02-03:** Quality assessment re-executed. Same structural failure in agents-md-format.md (criteria 1, 2, 3 fail). Criterion 4 continues to fail (no verification process). Criterion 5 passes. DEPRECATED specs (component-system.md, prompt-composition.md) correctly have replacement (component-authoring.md) and follow template structure. Refactoring plan regenerated in PLAN.md with identical 5-task structure and detailed findings per criterion. No changes to AGENTS.md needed - quality criteria remain accurate and boolean.

