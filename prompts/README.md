# Procedures and Components Specification

## Procedures

### Core Procedures

#### bootstrap

Create/update AGENTS.md operational guide.

- **Observe:** `observe_bootstrap.md`
- **Orient:** `orient_bootstrap.md`
- **Decide:** `decide_bootstrap.md`
- **Act:** `act_bootstrap.md`

#### build

Implement tasks from work tracking system.

- **Observe:** `observe_plan_specs_impl.md`
- **Orient:** `orient_build.md`
- **Decide:** `decide_build.md`
- **Act:** `act_build.md`

### Draft Planning Procedures

#### draft-plan-story-to-spec

Converge plan for incorporating story into specs.

- **Observe:** `observe_story_task_specs_impl.md`
- **Orient:** `orient_story_task_incorporation.md`
- **Decide:** `decide_story_task_plan.md`
- **Act:** `act_plan.md`

#### draft-plan-bug-to-spec

Converge plan for spec adjustments to drive bug fix.

- **Observe:** `observe_bug_task_specs_impl.md`
- **Orient:** `orient_bug_task_incorporation.md`
- **Decide:** `decide_bug_task_plan.md`
- **Act:** `act_plan.md`

#### draft-plan-spec-to-impl

Converge plan from gap analysis (specs → code).

- **Observe:** `observe_plan_specs_impl.md`
- **Orient:** `orient_gap.md`
- **Decide:** `decide_gap_plan.md`
- **Act:** `act_plan.md`

#### draft-plan-impl-to-spec

Converge plan from gap analysis (code → specs).

- **Observe:** `observe_plan_specs_impl.md`
- **Orient:** `orient_gap.md`
- **Decide:** `decide_gap_plan.md`
- **Act:** `act_plan.md`

#### draft-plan-spec-refactor

Converge plan from quality assessment of specs.

- **Observe:** `observe_specs.md`
- **Orient:** `orient_quality.md`
- **Decide:** `decide_refactor_plan.md`
- **Act:** `act_plan.md`

#### draft-plan-impl-refactor

Converge plan from quality assessment of code.

- **Observe:** `observe_impl.md`
- **Orient:** `orient_quality.md`
- **Decide:** `decide_refactor_plan.md`
- **Act:** `act_plan.md`

### Publishing Procedure

#### publish-plan

Publish converged draft plan to work tracking system.

- **Observe:** `observe_draft_plan.md` (new)
- **Orient:** `orient_publish.md` (new)
- **Decide:** `decide_publish.md` (new)
- **Act:** `act_publish.md` (new)

## Components

## Key Principles for Writing Components

**Reference Common Steps by Code** - Use step codes (O1, R5, D3, A2) rather than rewriting instructions. This ensures consistency and makes updates propagate automatically.

**One Component, One Concern** - Each component handles exactly one OODA phase. Observe gathers, Orient analyzes, Decide chooses, Act executes. No overlap.

**AGENTS.md is Always the Source of Truth** - Components must defer to AGENTS.md for all project-specific definitions: what constitutes specs/implementation, where files live, what commands to run, what quality criteria apply.

**Explicit Over Implicit** - State exactly what to read, analyze, decide, or do. "Study specifications per AGENTS.md definition" is better than "look at the specs."

**Use Precise Language** - Follow terminology from the Ralph Loop methodology:
- "study" (not "read" or "look at")
- "don't assume not implemented" (critical - the Achilles' heel)
- "using parallel subagents" / "only 1 subagent for build/tests"
- "capture the why" when updating AGENTS.md
- "keep it up to date" for maintaining accuracy
- "resolve them or document them" for issues found

**Search Before Assuming** - Orient components must emphasize searching the codebase before concluding something doesn't exist. This is the critical failure mode.

**Backpressure is Mandatory** - Act components that modify code must run tests and only commit when passing. No exceptions.

**Capture the Why** - When updating AGENTS.md, components must instruct agents to document rationale, not just changes. Why this command? Why this location? What was learned?

**Parallel Subagents for Scale** - Act components should use parallel subagents for independent work, but only 1 subagent for build/test operations to avoid conflicts.

**Plans are Disposable** - Planning components should generate complete plans each iteration, not incrementally patch. Cheap to regenerate beats expensive to maintain.

**Tight Tasks Win** - Decide components should break work into the smallest implementable units. One task per build iteration maximizes smart zone utilization.

**Commit After Complete** - Act components must complete all file modifications before committing. No partial work commits.

**Boolean Criteria Only** - Quality assessment components use PASS/FAIL criteria, not subjective scores. Clear thresholds trigger refactoring.

### Observe

**Common Steps:**
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

#### observe_bootstrap.md

Gather information about repository structure and existing documentation to understand the project.

O13, O14, O15

#### observe_plan_specs_impl.md

Gather information about work tracking, specifications, and implementation to understand current state.

O1, O2, O3, O7, O8, O11, O12

#### observe_story_task_specs_impl.md

Gather information about a story task, draft plan, specifications, and implementation for incorporation planning.

O1, O3, O5, O6, O9, O10, O11, O12

#### observe_bug_task_specs_impl.md

Gather information about a bug task, draft plan, specifications, and implementation for fix planning.

O1, O3, O5, O6, O9, O10, O11, O12

#### observe_specs.md

Gather information about specifications and quality criteria for assessment.

O1, O3, O4, O11

#### observe_impl.md

Gather information about implementation and quality criteria for assessment.

O1, O3, O4, O12

#### observe_draft_plan.md

Gather information about the converged draft plan ready for publishing.

O1, O6, O7, O10

### Orient

**Common Steps:**
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

#### orient_bootstrap.md

Analyze repository to determine project type, definitions, and operational commands.

R1, R2, R3, R4

#### orient_build.md

Analyze task requirements and codebase to determine what needs to be built.

R5, R6, R7, R8

#### orient_story_task_incorporation.md

Analyze story requirements and determine how to incorporate into specifications.

R9, R11, R12, R14, R15

#### orient_bug_task_incorporation.md

Analyze bug details and determine spec adjustments needed to drive the fix.

R10, R11, R13, R14, R15

#### orient_gap.md

Compare specifications and implementation to identify gaps in either direction.

R16, R6, R17

#### orient_quality.md

Assess quality using boolean criteria and identify areas needing improvement.

R18, R19, R20

#### orient_publish.md

Parse draft plan structure to prepare for work tracking system import.

R21, R22

### Decide

**Common Steps:**
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

#### decide_bootstrap.md

Determine AGENTS.md structure, definitions, and quality criteria for the project.

D1, D2, D3

#### decide_build.md

Select the most important task and determine implementation approach.

D4, D5, D6

#### decide_story_task_plan.md

Create a complete plan for incorporating the story into specifications.

D7, D9, D10, D11

#### decide_bug_task_plan.md

Create a complete plan for spec adjustments to drive the bug fix.

D8, D9, D10, D11

#### decide_gap_plan.md

Structure a prioritized plan to address identified gaps.

D9, D10, D11

#### decide_refactor_plan.md

Propose refactoring plan if quality criteria fail threshold.

D12, D9, D13

#### decide_publish.md

Determine how to map draft plan tasks to work tracking issues.

D14, D15

### Act

**Common Steps:**
- **A1:** Create AGENTS.md with operational guide
- **A2:** Commit changes
- **A3:** Implement using parallel subagents (only 1 subagent for build/tests)
- **A4:** Run tests per AGENTS.md (backpressure)
- **A5:** Update work tracking per AGENTS.md (mark complete/update status)
- **A6:** Update AGENTS.md if learned something new (capture the why, keep it up to date)
- **A7:** Commit when tests pass
- **A8:** Write draft plan file per AGENTS.md with prioritized bullet-point task list
- **A9:** Execute work tracking commands per AGENTS.md to create issues from draft plan

#### act_bootstrap.md

Create AGENTS.md operational guide and commit.

A1, A2

#### act_build.md

Implement code, run tests, update tracking, and commit when passing.

A3, A4, A5, A6, A7

#### act_plan.md

Write draft plan file, update AGENTS.md if needed, and commit.

A8, A6, A2

#### act_publish.md

Create work tracking issues from draft plan and commit.

A9, A6, A2
