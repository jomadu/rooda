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

**2026-02-03:** Beads dependency format uses `--deps` flag with format `type:id` or just `id`:
- Use `--deps blocks:issue-id` to indicate this issue is blocked by another issue
- Priority format is numeric (0-4) or P0-P4, not words like high/medium/low
- Multiple dependencies can be comma-separated: `--deps blocks:id1,blocks:id2`
- Rationale: Discovered during publish-plan procedure when attempting to create issues with dependencies

**2026-02-03:** Bootstrap validation added to act_bootstrap.md:
- Step A2 now validates AGENTS.md structure against required sections
- Missing sections trigger warnings in Operational Learnings
- Guidance provided for incomplete sections
- Required sections: Work Tracking, Story/Bug Input, Planning System, Build/Test/Lint, Spec Definition, Impl Definition, Quality Criteria
- Rationale: Ensures AGENTS.md completeness per agents-md-format.md specification
- yq dependency installed and working (required for YAML parsing in rooda.sh)
- shellcheck installed and linting passes cleanly on src/rooda.sh
- All 25 prompt components present in src/components/
- Beads work tracking responding correctly (empty queue is valid state)
- Repository structure matches documented patterns (src/, specs/, docs/ separation)
- Rationale: Empirical verification during bootstrap iteration confirms AGENTS.md accuracy

**2026-02-03:** Spec index generation integrated into build procedure:
- Builder reads specification-system.md to understand README structure
- Automatically regenerates specs/README.md when specs are modified
- Extracts JTBD from each spec file
- Excludes README.md, TEMPLATE.md, and specification-system.md from listing
- Rationale: Index is derived artifact from specs, not separate script. Dogfoods framework methodology (spec → implementation)

**2026-02-03 21:54:** Bootstrap validation confirms all systems operational:
- yq installed and working (v4.52.2) - required for YAML parsing in rooda.sh
- shellcheck installed and working - linting passes cleanly on src/rooda.sh
- All 25 prompt components present in src/components/
- Beads work tracking responding correctly (10 ready tasks in queue)
- Repository structure matches documented patterns (src/, specs/, docs/ separation)
- All required AGENTS.md sections present and complete
- Rationale: Empirical verification during bootstrap iteration confirms AGENTS.md accuracy and system readiness

**2026-02-03 22:06:** Bootstrap validation reconfirms operational status:
- yq installed and working - YAML parsing functional
- shellcheck installed and working - linting passes with no errors
- All 25 prompt components present in src/components/
- Beads work tracking operational (2 ready tasks: ralph-wiggum-ooda-abj, ralph-wiggum-ooda-4qd)
- Repository structure validated (src/, specs/, docs/ separation maintained)
- All required AGENTS.md sections present per agents-md-format.md specification
- Rationale: Periodic bootstrap validation ensures AGENTS.md remains accurate as work progresses

**2026-02-03:** Manual validation approach for bash script testing:
- Since this is a bash script repository with no automated test framework, validation uses manual test cases
- Create VALIDATION-<issue-id>.md files documenting test commands, expected behavior, and actual results
- Test cases verify acceptance criteria through empirical execution
- Validation documents serve as regression test documentation
- Rationale: Discovered during ralph-wiggum-ooda-1w0 (CLI procedure validation) - manual verification is the appropriate testing methodology for this framework

**2026-02-03:** Bug discovered in OODA file path resolution (ralph-wiggum-ooda-cuy):
- Config file resolution works correctly (resolves relative to script location)
- OODA phase file paths from config are NOT resolved relative to script directory
- File validation at line 107 checks paths relative to current working directory
- Impact: Script only works when invoked from project root, fails from other directories
- Root cause: Paths loaded from config (lines 77-80) are used as-is without prepending SCRIPT_DIR
- Fix needed: Resolve OODA paths relative to SCRIPT_DIR after loading from config
- Rationale: Empirical testing during ralph-wiggum-ooda-cuy validation revealed this gap between config resolution (works) and OODA file resolution (broken)

**2026-02-03 14:07:** Bootstrap validation confirms continued operational status:
- yq installed and working (v4.52.2) - YAML parsing functional
- shellcheck installed and working - linting passes cleanly
- All 25 prompt components present in src/components/
- Beads work tracking operational (2 ready tasks: ralph-wiggum-ooda-abj, ralph-wiggum-ooda-4qd)
- Repository structure validated (src/, specs/, docs/ separation maintained)
- All required AGENTS.md sections present and complete per agents-md-format.md specification
- Rationale: Periodic bootstrap validation ensures AGENTS.md remains accurate as work progresses

**2026-02-03 14:10:** Version validation implemented (ralph-wiggum-ooda-5ib):
- yq version >= 4.0.0 now validated at startup (lines 64-93 in src/rooda.sh)
- kiro-cli version >= 1.0.0 validated
- bd version >= 0.1.0 validated
- Clear error messages for incompatible versions with upgrade instructions
- Rationale: Prevents cryptic YAML parsing errors from yq v3, ensures all tools meet minimum requirements per external-dependencies.md specification

