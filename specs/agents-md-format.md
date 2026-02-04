# AGENTS.md Specification

## Job to be Done

AGENTS.md is the interface between agents and the repository. It defines how agents interact with project-specific workflows, tools, and conventions.

## Acceptance Criteria

- [ ] AGENTS.md contains Work Tracking System section with query/update/complete commands
- [ ] AGENTS.md contains Story/Bug Input section defining where to read task descriptions
- [ ] AGENTS.md contains Planning System section with draft location and publishing mechanism
- [ ] AGENTS.md contains Build/Test/Lint Commands section with specific commands
- [ ] AGENTS.md contains Specification Definition section with file paths/patterns
- [ ] AGENTS.md contains Implementation Definition section with file paths/patterns
- [ ] AGENTS.md contains Quality Criteria section with boolean PASS/FAIL criteria
- [ ] All commands in AGENTS.md are empirically verified to work
- [ ] AGENTS.md is updated when operational learnings occur
- [ ] AGENTS.md includes rationale (the "why") for key decisions

## What's in AGENTS.md?

## Required Sections

### Work Tracking System
- What system tracks work
- How to query ready work
- How to update status
- How to mark complete

**Examples:**
```
# Beads
Query: bd ready --json
Update: bd update <id> --status in_progress
Complete: bd close <id> --reason "Done"

# GitHub Issues
Query: gh issue list --label ready --json
Update: gh issue edit <id> --add-label in-progress
Complete: gh issue close <id>

# File-based
Query: ls tasks/ready/*.md
Update: mv tasks/ready/<id>.md tasks/in-progress/
Complete: mv tasks/in-progress/<id>.md tasks/done/
```

### Story/Bug Input
Where agents read the story/bug description for `draft-plan-story-to-spec` and `draft-plan-bug-to-spec` procedures.

**Fixed file (simplest):**
```
TASK.md at project root
```

**Environment variable pointing to file:**
```
stories/$TASK_ID.md
```

**Environment variable with command:**
```
bd show $TASK_ID --json  # Use title + description fields
gh issue view $TASK_ID --json
```

### Planning System
- **Draft plan location** - Where agents write/read plans during convergence iterations
- **Publishing mechanism** - How converged plans get imported into work tracking
- Used by draft planning procedures to iterate toward a complete plan, then publish to work tracking

**Workflow:**
1. Draft procedures (`draft-plan-*`) iterate to converge on a plan in the draft location
2. Publish procedure (`publish-plan`) imports the converged plan into work tracking system
3. Build procedure (`build`) implements from work tracking system

**Examples:**
```
# Simple PLAN.md (simplest)
Draft plan: PLAN.md
Publishing: Agent runs work tracking commands to create issues from PLAN.md

# Beads with draft plans
Draft plan: plans/draft-<topic>.md
Publishing: Agent runs `bd create` commands to file epics/issues with dependencies

# GitHub Issues with draft plans
Draft plan: plans/draft-<topic>.md
Publishing: Agent runs `gh issue create` commands with labels and milestones

# File-based work tracking
Draft plan: tasks/<id>/plan.md
Publishing: Plan file becomes the work tracking (no import needed)
```

### Build/Test/Lint Commands
- Specific commands to run tests
- Specific commands to run builds
- Specific commands to run linters
- Any setup required before running these commands

**Examples:**
```
# Node.js
Test: npm test
Build: npm run build
Lint: npm run lint

# Go
Test: go test ./...
Build: go build ./...
Lint: golangci-lint run

# Python
Test: pytest
Build: python -m build
Lint: ruff check .
```

### Specification Definition
- What constitutes "specification" in this repository
- File paths or patterns

**Examples:**
```
# Dedicated specs directory
specs/*.md

# README sections
README.md (## Specification sections only)

# Inline documentation
src/**/*.md (excluding README.md files)

# API documentation
docs/api/*.md
```

### Implementation Definition
- What constitutes "implementation" in this repository
- File paths or patterns

**Examples:**
```
# Source directory
src/**/*.{js,ts,py,go,rs}

# Library directory
lib/**/*.rb

# Multiple source roots
src/**/*.java, test/**/*.java

# Exclude patterns
src/**/*.ts (excluding *.test.ts, *.spec.ts)
```

### Quality Criteria
- Boolean criteria for triggering spec refactoring (clarity, completeness, consistency, testability)
- Boolean criteria for triggering implementation refactoring (cohesion, coupling, complexity, maintainability)

## Key Principles

### Assumed Inaccurate Until Verified
AGENTS.md is a working hypothesis about how the project operates. It should be updated when:
- Commands fail or produce unexpected results
- File paths are incorrect
- Quality criteria don't match project needs
- New patterns or conventions are discovered

### Capture the Why
When updating AGENTS.md, include rationale:
- Why this command instead of another
- Why these file paths
- Why these quality criteria
- What was learned that prompted the update

### Keep It Up to Date
- **Agents update** when they discover errors or learn new operational details during any procedure
- **Humans update** when they change project structure, tooling, or conventions
- Update immediately when discovered, don't defer

### Source of Truth
AGENTS.md is the authoritative definition of how agents interact with the project. When there's ambiguity:
1. Check AGENTS.md first
2. Verify empirically (run commands, check paths)
3. Update AGENTS.md with findings
4. Commit the update

### Concise and Operational
- Focus on "how to" not "what is"
- Specific commands, not general descriptions
- Concrete paths, not abstract concepts
- ~60 lines is a good target (not a hard limit)

## When to Update

### During Bootstrap
- Initial creation based on repository analysis
- Empirical discovery of build/test commands
- Identification of spec and implementation patterns

### During Build
- When commands fail or need adjustment
- When file paths are incorrect
- When new patterns are discovered

### During Planning
- When quality criteria need refinement
- When work tracking system changes
- When new operational learnings emerge

### By Humans
- When changing project structure
- When adding/removing tools
- When adopting new conventions
- When onboarding new team members (clarify ambiguities)

## Path Conventions

**Consumer projects** (using ralph-wiggum-ooda):
- Copy components to `prompts/` directory at project root (flat structure)
- Config file references `prompts/*.md`
- Examples in this spec use `prompts/` to match consumer convention

**Framework repository** (ralph-wiggum-ooda itself):
- Components stored in `src/components/` (organized internal structure)
- Config file references `src/components/*.md`
- Consumers copy from `src/components/` to their `prompts/`

This separation allows the framework to maintain organized structure while consumers get simple flat layout.

## Examples

### Example 1: Node.js Project with GitHub Issues

**AGENTS.md:**
```markdown
## Work Tracking System
System: GitHub Issues
Query: gh issue list --label ready --json
Update: gh issue edit <id> --add-label in-progress
Complete: gh issue close <id>

## Build/Test/Lint Commands
Test: npm test
Build: npm run build
Lint: npm run lint

## Specification Definition
Location: specs/*.md

## Implementation Definition
Location: src/**/*.{js,ts}
```

### Example 2: Go Project with Beads

**AGENTS.md:**
```markdown
## Work Tracking System
System: beads (bd CLI)
Query: bd ready --json
Update: bd update <id> --status in_progress
Complete: bd close <id> --reason "Completed X"

## Build/Test/Lint Commands
Test: go test ./...
Build: go build ./...
Lint: golangci-lint run

## Specification Definition
Location: docs/specs/*.md

## Implementation Definition
Location: pkg/**/*.go, cmd/**/*.go
```

### Example 3: Python Project with File-Based Tracking

**AGENTS.md:**
```markdown
## Work Tracking System
System: File-based
Query: ls tasks/ready/*.md
Update: mv tasks/ready/<id>.md tasks/in-progress/
Complete: mv tasks/in-progress/<id>.md tasks/done/

## Story/Bug Input
Location: TASK.md at project root

## Build/Test/Lint Commands
Test: pytest
Build: python -m build
Lint: ruff check .

## Specification Definition
Location: specs/*.md

## Implementation Definition
Location: src/**/*.py
```

## Anti-Patterns

### Don't Make It a Changelog
AGENTS.md is not a log of what happened. It's a guide for what to do now.

### Don't Make It a Progress Diary
Status and progress belong in the work tracking system, not AGENTS.md.

### Don't Assume It's Correct
Always verify empirically. If AGENTS.md says "run `npm test`" but that fails, update it.

### Don't Let It Get Stale
Update immediately when you discover it's wrong. Stale AGENTS.md causes repeated failures.