# Component Authoring

## Job to be Done
Enable developers to create and modify OODA component prompt files that can be composed into executable procedures.

## Activities
- **Write Prompt Files** - Create markdown files with phase headers, step codes, and prose instructions
- **Reference Common Steps** - Use step codes (O1-O15, R1-R22, D1-D15, A1-A9) to structure components
- **Compose Procedures** - Combine four components (observe, orient, decide, act) into complete procedures
- **Maintain Consistency** - Follow naming conventions and structural patterns across all components

## Acceptance Criteria
- [ ] Prompt file structure documented (phase header, step headers, prose instructions)
- [ ] Step code patterns explained (O1-O15, R1-R22, D1-D15, A1-A9)
- [ ] Complete common steps reference provided
- [ ] Prompt assembly algorithm documented
- [ ] Authoring guidelines included
- [ ] Real examples from actual prompt files shown
- [ ] Dual purpose of step codes clarified (structure + cross-reference)

## Data Structures

### Prompt File Format
```markdown
# [Phase Name]: [Purpose]

## [Step Code]: [Step Name]

[Detailed prose instructions for this step]

## [Step Code]: [Step Name]

[Detailed prose instructions for this step]
```

**Fields:**
- `Phase Name` - One of: Observe, Orient, Decide, Act
- `Purpose` - Brief description of what this component does
- `Step Code` - Unique identifier (O1-O15, R1-R22, D1-D15, A1-A9)
- `Step Name` - Human-readable step description
- `Detailed prose instructions` - Full instructions for executing this step

**Key Characteristics:**
- Step codes are section headers within the file (## O1, ## R5, ## A3)
- Instructions are written in full prose under each step header
- No external references required - each step is self-contained
- Step codes enable cross-referencing when needed

### Assembled Prompt Format
```markdown
# OODA Loop Iteration

## OBSERVE
[contents of observe prompt file]

## ORIENT
[contents of orient prompt file]

## DECIDE
[contents of decide prompt file]

## ACT
[contents of act prompt file]
```

## Algorithm

**Prompt Creation:**
1. Choose OODA phase (observe, orient, decide, act)
2. Define purpose (what this component accomplishes)
3. Select relevant step codes from common steps reference
4. Write phase header: `# [Phase]: [Purpose]`
5. For each step: write step header `## [Code]: [Name]` and full prose instructions
6. Save as `src/prompts/[phase]_[purpose].md`

**Prompt Assembly:**
1. Load procedure configuration from rooda-config.yml
2. Read four prompt files (observe, orient, decide, act)
3. Validate all files exist
4. Concatenate with OODA section headers
5. Output combined prompt to stdout

**Pseudocode:**
```bash
function create_prompt():
    cat <<EOF
# OODA Loop Iteration

## OBSERVE
$(cat "$OBSERVE")

## ORIENT
$(cat "$ORIENT")

## DECIDE
$(cat "$DECIDE")

## ACT
$(cat "$ACT")
EOF
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt file uses undefined step code | Valid - step codes are headers, not references |
| Multiple prompt files use same step code | Valid - each prompt file is independent |
| Prompt file omits step codes entirely | Valid - step codes are optional structure |
| Step code used but instructions differ from common reference | Valid - prompt file instructions override |
| Missing prompt file | Script exits with error before prompt assembly |
| Empty prompt file | Section header present, no content under it |

## Dependencies

- `src/rooda.sh` - Script that loads and combines prompts
- `src/rooda-config.yml` - Procedure configuration mapping prompts to procedures

## Implementation Mapping

**Source files:**
- `src/prompts/observe_*.md` - Observation phase components
- `src/prompts/orient_*.md` - Analysis phase components
- `src/prompts/decide_*.md` - Decision phase components
- `src/prompts/act_*.md` - Execution phase components
- `src/rooda.sh` (lines 143-159) - create_prompt() function

**Related specs:**
- `specs/configuration-schema.md` - Procedure configuration

## Examples

### Example 1: Act Component Structure

**File:** `src/prompts/act_bootstrap.md`

```markdown
# Act: Bootstrap

## A1: Create AGENTS.md with Operational Guide

Write AGENTS.md with the following sections:

**Work Tracking System:**
- What system is used
- How to query ready work
- How to update status
- How to mark complete
- Specific commands

**Story/Bug Input:**
- Where agents read story/bug descriptions
- File location or command to retrieve

[... additional prose instructions ...]

## A2: Validate AGENTS.md Structure

Check AGENTS.md for required sections per agents-md-format.md:
- Work Tracking System
- Story/Bug Input
- Planning System

[... additional prose instructions ...]

## A3: Commit Changes

Commit AGENTS.md with descriptive message:
- What was created
- Why these definitions were chosen
- What was learned during bootstrap
```

**Verification:**
- Phase header present: `# Act: Bootstrap`
- Step headers use A-codes: `## A1`, `## A2`, `## A3`
- Full prose instructions under each step
- Self-contained - no external references needed

### Example 2: Observe Component Structure

**File:** `src/prompts/observe_bootstrap.md`

```markdown
# Observe: Bootstrap

## O13: Study Repository Structure

Examine the repository to understand:
- File tree structure (directories, organization)
- What programming languages are used?
- What build files exist? (package.json, Cargo.toml, go.mod, etc.)
- What configuration files exist?

[... additional prose instructions ...]

## O14: Study Existing Documentation

Read available documentation:
- README.md (project description, setup, usage)
- CONTRIBUTING.md (development workflow)
- docs/ directory (if present)

[... additional prose instructions ...]

## O15: Study Implementation Patterns

Examine existing code to identify:
- Code organization patterns
- Naming conventions
- Module structure
```

**Verification:**
- Phase header present: `# Observe: Bootstrap`
- Step headers use O-codes: `## O13`, `## O14`, `## O15`
- Full prose instructions under each step
- Instructions are specific and actionable

### Example 3: Act Component with Substeps

**File:** `src/prompts/act_build.md`

```markdown
# Act: Build

## A3: Implement Using Parallel Subagents

Execute the implementation plan:
- Spawn subagents for parallel work as determined in decide phase
- **Critical: Only 1 subagent for build/tests** (avoid parallel test conflicts)
- Each subagent handles its assigned files/tasks

## A3.5: Validate Spec Structure (If Specs Modified)

If any spec files in `specs/*.md` were created or modified:
- Check filename matches naming convention (lowercase-with-hyphens.md)
- Check for required sections from TEMPLATE.md
- Warn if required sections are missing

## A3.6: Regenerate Spec Index (If Specs Modified)

If any spec files in `specs/*.md` were created, modified, or deleted:
- Read `specs/specification-system.md` for README structure requirements
- Scan all `specs/*.md` files
- Extract "## Job to be Done" section from each spec

## A4: Run Tests per AGENTS.md (Backpressure)

Execute test commands from AGENTS.md:
- Run all relevant tests
- Verify tests pass
- If tests fail: fix issues, don't proceed
```

**Verification:**
- Substep numbering: A3, A3.5, A3.6, A4
- Conditional steps clearly marked: "If Specs Modified"
- Critical warnings emphasized: **Critical: Only 1 subagent**
- Backpressure concept explained in step name

## Common Steps Reference

### Observe Steps (O1-O15)

- **O1:** Study AGENTS.md as a whole (operational guide, conventions, learnings)
- **O2:** Study AGENTS.md for build/test commands
- **O3:** Study AGENTS.md for specification/implementation definitions
- **O4:** Study AGENTS.md for quality criteria definitions
- **O5:** Study AGENTS.md for task file location
- **O6:** Study AGENTS.md for draft plan location
- **O7:** Study AGENTS.md for work tracking system
- **O8:** Study work tracking system per AGENTS.md (query ready work)
- **O9:** Study task file per AGENTS.md (story/bug description)
- **O10:** Study draft plan file per AGENTS.md (current plan state, may not exist)
- **O11:** Study specifications per AGENTS.md definition
- **O12:** Study implementation per AGENTS.md definition (file tree, symbols)
- **O13:** Study repository structure (file tree, languages, build files)
- **O14:** Study existing documentation (README, specs if present)
- **O15:** Study implementation patterns

### Orient Steps (R1-R22)

- **R1:** Identify project type and tech stack
- **R2:** Determine what constitutes "specification" vs "implementation"
- **R3:** Identify build/test/run commands empirically
- **R4:** Synthesize operational understanding
- **R5:** Understand task requirements
- **R6:** Search codebase (don't assume not implemented)
- **R7:** Identify what needs to be built/modified
- **R8:** Determine test strategy
- **R9:** Analyze story from task file (scope, requirements, integration points)
- **R10:** Analyze bug from task file (symptoms, root cause, affected functionality)
- **R11:** Understand existing spec structure and patterns
- **R12:** Determine how story should be incorporated (create new specs, update existing, refactor)
- **R13:** Determine how spec should be adjusted to drive bug fix (acceptance criteria, edge cases, clarifications)
- **R14:** If draft plan exists: critique it (completeness, accuracy, priorities, clarity)
- **R15:** Identify tasks needed
- **R16:** Gap analysis: compare specs vs implementation
- **R17:** Assess completeness and accuracy
- **R18:** Apply boolean criteria per AGENTS.md
- **R19:** Identify human markers (TODOs, code smells, unclear language)
- **R20:** Score each criterion PASS/FAIL
- **R21:** Parse draft plan structure
- **R22:** Understand task breakdown and dependencies

### Decide Steps (D1-D15)

- **D1:** Determine AGENTS.md structure
- **D2:** Define specification and implementation locations
- **D3:** Identify quality criteria for this project
- **D4:** Pick the most important task from work tracking
- **D5:** Determine implementation approach using parallel subagents
- **D6:** Identify which files to modify
- **D7:** Generate complete plan for story incorporation into specs
- **D8:** Generate complete plan for spec adjustments to drive bug fix
- **D9:** Structure plan by priority (most important first)
- **D10:** Break into tight, actionable tasks
- **D11:** Determine task dependencies
- **D12:** If criteria fail threshold: propose refactoring
- **D13:** Prioritize by impact
- **D14:** Map plan tasks to work tracking issues (title, description, dependencies)
- **D15:** Identify order of issue creation

### Act Steps (A1-A9)

- **A1:** Create AGENTS.md with operational guide
- **A2:** Validate AGENTS.md structure
- **A3:** Implement using parallel subagents (only 1 subagent for build/tests)
- **A4:** Run tests per AGENTS.md (backpressure)
- **A5:** Update work tracking per AGENTS.md (mark complete/update status)
- **A6:** Incorporate learnings into AGENTS.md (update commands/paths/criteria in place, add inline rationale, no dated diary entries)
- **A7:** Commit when tests pass
- **A8:** Write draft plan file per AGENTS.md with prioritized bullet-point task list
- **A9:** Execute work tracking commands per AGENTS.md to create issues from draft plan

## Authoring Guidelines

**Reference Common Steps by Code** - Use step codes (O1, R5, D3, A2) as section headers with full prose instructions. This provides structure and enables cross-referencing when needed.

**One Component, One Concern** - Each component handles exactly one OODA phase. Observe gathers, Orient analyzes, Decide chooses, Act executes. No overlap.

**AGENTS.md is Always the Source of Truth** - Components must defer to AGENTS.md for all project-specific definitions: what constitutes specs/implementation, where files live, what commands to run, what quality criteria apply.

**Explicit Over Implicit** - State exactly what to read, analyze, decide, or do. "Study specifications per AGENTS.md definition" is better than "look at the specs."

**Use Precise Language** - Follow terminology from the Ralph Loop methodology:
- "study" (not "read" or "look at")
- "don't assume not implemented" (critical - the Achilles' heel)
- "using parallel subagents" / "only 1 subagent for build/tests"
- "capture the why" when updating AGENTS.md
- "keep it up to date" for maintaining accuracy
- "backpressure" for quality gates

**Search Before Assuming** - Orient components must emphasize searching the codebase before concluding something doesn't exist. This is the critical failure mode.

**Backpressure is Mandatory** - Act components that modify code must run tests and only commit when passing. No exceptions.

**Incorporate Learnings, Don't Append** - When updating AGENTS.md, components must instruct agents to incorporate learnings into existing sections (commands, paths, criteria) with inline rationale, not append dated diary entries. Exception: significant architectural decisions may warrant dated entries if historical context matters.

**Parallel Subagents for Scale** - Act components should use parallel subagents for independent work, but only 1 subagent for build/test operations to avoid conflicts.

**Plans are Disposable** - Planning components should generate complete plans each iteration, not incrementally patch. Cheap to regenerate beats expensive to maintain.

**Tight Tasks Win** - Decide components should break work into the smallest implementable units. One task per build iteration maximizes smart zone utilization.

**Commit After Complete** - Act components must complete all file modifications before committing. No partial work commits.

**Boolean Criteria Only** - Quality assessment components use PASS/FAIL criteria, not subjective scores. Clear thresholds trigger refactoring.

## Notes

**Dual Purpose of Step Codes:**

Step codes (O1-O15, R1-R22, D1-D15, A1-A9) serve two purposes:

1. **Structure** - They are section headers within prompt files (## O1, ## R5, ## A3). Each header is followed by full prose instructions that are self-contained and executable.

2. **Cross-Reference** - They enable referencing common step definitions when needed. The common steps reference provides a quick lookup of what each code typically means, but prompt files contain the actual instructions.

This dual purpose means components are both structured (using consistent step codes) and self-contained (containing full instructions). You don't need to look up O1 to understand what to do - the instructions are right there under the ## O1 header.

**Prompt Reuse:**

**Prompt File Reuse:**

The prompt file system enables significant reuse across procedures. For example, `observe_plan_specs_impl.md` is shared by `build`, `draft-plan-spec-to-impl`, and `draft-plan-impl-to-spec` procedures. They differ only in their orient, decide, and act prompts.

**File Naming Convention:**

Prompt files follow the pattern: `[phase]_[purpose].md`
- Phase: observe, orient, decide, act
- Purpose: descriptive name (bootstrap, build, gap, quality, plan, publish)
- Examples: `observe_bootstrap.md`, `act_build.md`, `orient_gap.md`

## Known Issues

None identified during specification creation.

## Areas for Improvement

- Could add validation tooling to check prompt files follow structure
- Could add linting for step code consistency across prompts
- Could add examples of custom prompt creation for project-specific procedures
- Could document how to extend common steps reference with project-specific codes
