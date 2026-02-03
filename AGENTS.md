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
shellcheck rooda.sh  # Lint bash script
```

**Verification:**
```bash
./rooda.sh --help  # Verify script runs
bd ready --json    # Verify beads integration works
```

## Specification Definition

**Location:** `docs/*.md` and `prompts/*.md`

**Format:** Markdown documentation

**Patterns:**
- `docs/agents-md-specification.md` - AGENTS.md format specification
- `docs/ooda-loop.md` - OODA framework explanation
- `docs/ralph-loop.md` - Original methodology
- `docs/specs.md` - Specification system design
- `docs/spec-template.md` - Template for specs
- `prompts/*.md` - OODA phase prompt components

**Rationale:** This is a framework/methodology repository. The "specs" are the documentation that defines how the system works. Implementation is the bash script that executes the framework.

## Implementation Definition

**Location:** `rooda.sh` (bash script)

**Patterns:** `rooda.sh`

**Exclude:** 
- `.beads/*` (work tracking database)
- `prompts/*` (prompt templates, not implementation)
- `docs/*` (documentation, not implementation)
- `*.md` (documentation files)

**Rationale:** The only implementation is the bash loop script. Everything else is documentation, configuration, or prompt templates that define the framework's behavior.

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

