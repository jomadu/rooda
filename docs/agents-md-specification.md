# AGENTS.md Specification

## Purpose

AGENTS.md is the interface between agents and the repository. It defines how agents interact with project-specific workflows, tools, and conventions.

## What's in AGENTS.md?

When an agent is prompted to update the work tracking system, AGENTS.md tells the agent how to update the work tracking system for that particular repository. Same for build/test/lint commands, and all the rest of the operational details.

## Required Sections

### Task/Story/Bug Descriptions
- Where to find task descriptions for story/bug incorporation procedures (if applicable)
- Where to find or write plans during incorporation (if applicable)

### Work Tracking System
- What system tracks work (beads, GitHub issues, linear, etc.)
- How to query ready work
- How to update status
- How to mark complete

### Build/Test/Lint Commands
- Specific commands to run tests
- Specific commands to run builds
- Specific commands to run linters
- Any setup required before running these commands

### Specification Definition
- What constitutes "specification" in this repository
- File paths or patterns (e.g., `specs/*.md`, `README.md sections`, inline docs)

### Implementation Definition
- What constitutes "implementation" in this repository
- File paths or patterns (e.g., `src/`, `lib/`, specific file extensions)

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

## Anti-Patterns

### Don't Make It a Changelog
AGENTS.md is not a log of what happened. It's a guide for what to do now.

### Don't Make It a Progress Diary
Status and progress belong in the work tracking system, not AGENTS.md.

### Don't Assume It's Correct
Always verify empirically. If AGENTS.md says "run `npm test`" but that fails, update it.

### Don't Let It Get Stale
Update immediately when you discover it's wrong. Stale AGENTS.md causes repeated failures.