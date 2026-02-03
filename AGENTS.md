# Agent Instructions

## Work Tracking System

**System:** beads (bd CLI)

**Query ready work:**
```bash
bd ready --json
```

**View issue details:**
```bash
bd show <id> --json
```

**Update status:**
```bash
bd update <id> --status in_progress
bd update <id> --status blocked
```

**Close issue:**
```bash
bd close <id> --reason "Completed X"
```

**Sync with git:**
```bash
bd sync
```

## Story/Bug Input

**For draft-plan-story-to-spec and draft-plan-bug-to-spec procedures:**
- Use `bd show $TASK_ID --json` to read story/bug description
- Extract `title` and `description` fields from JSON output
- Task ID should be provided via environment variable or command-line argument

## Planning System

**Draft plan location:** `PLAN.md` at project root

**Publishing mechanism:** Agent reads `PLAN.md` and runs `bd create` commands to file issues with appropriate dependencies and priorities

**Workflow:**
1. Draft procedures write/iterate on `PLAN.md`
2. Publish procedure reads `PLAN.md` and creates beads issues
3. Build procedure implements from beads work tracking

## Build/Test/Lint Commands

**This is a documentation repository (bash scripts + markdown).**

**Test:** No automated tests (manual verification of script execution)

**Build:** No build step required (bash scripts are interpreted)

**Lint:** 
```bash
shellcheck rooda.sh  # Lint bash script (if installed)
```

**Note:** shellcheck may not be installed on all systems. Install with `brew install shellcheck` on macOS.

**Verification:**
```bash
./rooda.sh bootstrap  # Verify script runs (use explicit flags due to YAML parser issue)
bd ready --json       # Verify beads integration works
```

## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure

**Patterns:**
- `specs/agents-md-format.md` - AGENTS.md format specification
- `specs/specification-system.md` - Spec system design
- `specs/TEMPLATE.md` - Template for new specs

**Rationale:** Specifications define the framework methodology (JTBD, spec system, AGENTS.md format). These are distinct from user-facing documentation.

## Implementation Definition

**Location:** `src/rooda.sh` and `src/components/*.md`

**Patterns:** 
- `src/rooda.sh` - Main loop script
- `src/rooda-config.yml` - Procedure configuration
- `src/components/*.md` - OODA prompt components

**Exclude:** 
- `.beads/*` (work tracking database)
- `docs/*` (user-facing documentation)
- `specs/*` (specifications)
- `README.md`, `AGENTS.md`, `LICENSE.md` (root-level files)

**Rationale:** Implementation is the bash script and composable prompt components that execute the framework. Configuration and components are co-located with the script.

## Quality Criteria

**For specifications (documentation):**
- Clarity: Can a new user understand the framework from README.md?
- Completeness: Are all 9 procedures documented with examples?
- Consistency: Do docs match actual script behavior?
- Accuracy: Do command examples work when executed?

**For implementation (rooda.sh):**
- Correctness: Does the script execute procedures as documented?
- Robustness: Does error handling work for missing files/commands?
- Maintainability: Is the bash code readable and commented?
- Compatibility: Does it work on macOS and Linux?

**Refactoring triggers:**
- Documentation contradicts script behavior
- Script fails on documented use cases
- Error messages are unclear or misleading
- YAML parsing doesn't support documented config structure

## Operational Learnings

**2026-02-03:** Bootstrap procedure identified that this is a framework repository, not a typical application. The "specification" is the documentation that defines the methodology, and the "implementation" is the bash script that executes it. This differs from typical projects where specs describe features and implementation is source code.

**2026-02-03:** Beads integration is working correctly. The `bd ready --json` command returns structured JSON with issue details. Work tracking commands are operational.

**2026-02-03:** There's an open issue (ralph-wiggum-ooda-i2c) about YAML parser in rooda.sh not properly handling procedure lookups. Workaround exists using explicit flags. This should be fixed to support the documented `./rooda.sh bootstrap` syntax.

**2026-02-03:** shellcheck is not installed on this system. The lint command will fail until shellcheck is installed via `brew install shellcheck`. This is optional for framework operation but recommended for bash script quality.

**2026-02-03:** Restructured repository into src/, specs/, docs/ to enable dogfooding:
- Separating implementation (src/) from specifications (specs/) allows running draft-plan-impl-to-spec
- Framework can now use its own methodology to generate specs from implementation
- Clear separation makes it obvious what agents should analyze vs what they should read for guidance
- Consumers copy from src/ to their project root (flat structure), while framework repo has internal organization
- Config file paths remain prompts/*.md for consumer compatibility

