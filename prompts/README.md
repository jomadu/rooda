# Minimal Composable Prompt Set

## Procedure Breakdown by Phase

### Procedure 0: Bootstrap (create AGENTS.md)

**OBSERVE**
- Study repository structure (file tree, languages, build files)
- Study existing documentation (README, specs if present)
- Study implementation patterns

**ORIENT**
- Identify project type and tech stack
- Determine what constitutes "specification" vs "implementation"
- Identify build/test/run commands empirically
- Synthesize operational understanding

**DECIDE**
- Determine AGENTS.md structure
- Define specification and implementation locations
- Identify quality criteria for this project

**ACT**
- Create AGENTS.md with operational guide
- Commit changes

---

### Procedure 1: Building from plan

**OBSERVE**
- Study AGENTS.md (how to build/test, what is specification/implementation)
- Study plan file per AGENTS.md (find most important task)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Understand task requirements
- Search codebase (don't assume not implemented)
- Identify what needs to be built/modified
- Determine test strategy

**DECIDE**
- Pick most important task from plan file
- Determine implementation approach using parallel subagents
- Identify which files to modify

**ACT**
- Implement using parallel subagents (only 1 subagent for build/tests)
- Run tests per AGENTS.md (backpressure)
- Update plan file (mark complete/update)
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit when tests pass

---

### Procedure 2: Plan spec-to-impl

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study plan file per AGENTS.md (current state)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Gap analysis: compare specs against existing code
- Search codebase (don't assume not implemented)
- Identify what's in specifications but missing from implementation
- Assess implementation completeness and accuracy

**DECIDE**
- Structure plan by priority (most important tasks first)
- Break gaps into tight, implementable tasks
- Determine task dependencies

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 3: Plan impl-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study plan file per AGENTS.md (current state)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Gap analysis: compare existing code against specs
- Search codebase thoroughly
- Identify what's in implementation but missing from specifications
- Assess specification completeness and accuracy

**DECIDE**
- Structure plan by priority (most important tasks first)
- Break gaps into tight documentation tasks
- Determine which specs need updates

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 4: Plan spec refactoring

**OBSERVE**
- Study AGENTS.md (what is specification, quality criteria definitions)
- Study specifications (per AGENTS.md definition)

**ORIENT**
- Apply boolean criteria: clarity, completeness, consistency, testability
- Identify human markers (TODOs, "REFACTORME", unclear language)
- Score each criterion PASS/FAIL
- Resolve issues or document them

**DECIDE**
- If criteria fail threshold: propose spec refactoring
- Structure spec refactoring plan by priority
- Prioritize by impact (most important first)

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 5: Plan impl refactoring

**OBSERVE**
- Study AGENTS.md (what is implementation, quality criteria definitions)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Apply boolean criteria: cohesion, coupling, complexity, maintainability
- Identify human markers (TODOs, long functions, code smells)
- Score each criterion PASS/FAIL
- Resolve issues or document them

**DECIDE**
- If criteria fail threshold: propose implementation refactoring
- Structure implementation refactoring plan by priority
- Prioritize by impact (most important first)

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 2: Plan spec-to-impl

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study plan file per AGENTS.md (current state)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Gap analysis: compare specs against existing code
- Search codebase (don't assume not implemented)
- Identify what's in specifications but missing from implementation
- Assess implementation completeness and accuracy

**DECIDE**
- Structure plan by priority (most important tasks first)
- Break gaps into tight, implementable tasks
- Determine task dependencies

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 3: Plan impl-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study plan file per AGENTS.md (current state)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Gap analysis: compare existing code against specs
- Search codebase thoroughly
- Identify what's in implementation but missing from specifications
- Assess specification completeness and accuracy

**DECIDE**
- Structure plan by priority (most important tasks first)
- Break gaps into tight documentation tasks
- Determine which specs need updates

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 4: Plan spec refactoring

**OBSERVE**
- Study AGENTS.md (what is specification, quality criteria definitions)
- Study specifications (per AGENTS.md definition)

**ORIENT**
- Apply boolean criteria: clarity, completeness, consistency, testability
- Identify human markers (TODOs, "REFACTORME", unclear language)
- Score each criterion PASS/FAIL
- Resolve issues or document them

**DECIDE**
- If criteria fail threshold: propose spec refactoring
- Structure spec refactoring plan by priority
- Prioritize by impact (most important first)

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 5: Plan impl refactoring

**OBSERVE**
- Study AGENTS.md (what is implementation, quality criteria definitions)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Apply boolean criteria: cohesion, coupling, complexity, maintainability
- Identify human markers (TODOs, long functions, code smells)
- Score each criterion PASS/FAIL
- Resolve issues or document them

**DECIDE**
- If criteria fail threshold: propose implementation refactoring
- Structure implementation refactoring plan by priority
- Prioritize by impact (most important first)

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 6: Plan story-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study task file per AGENTS.md (the feature/story to be incorporated)
- Study plan file in task directory (current plan state - may not exist on first iteration)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Analyze the story from task file per AGENTS.md (scope, requirements, integration points)
- Understand existing spec structure and patterns
- Determine how the story should be incorporated (create new specs, update existing, refactor)
- If plan file exists: critique it (completeness, accuracy, priorities, clarity)
- Identify what tasks are needed for proper incorporation

**DECIDE**
- Generate a complete plan for incorporating the story into specs
- Structure the plan by priority (most important tasks first)
- Break incorporation into tight, actionable tasks
- Determine task dependencies and which specs need creation/updates

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Procedure 7: Plan bug-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study task file per AGENTS.md (the bug to be addressed)
- Study plan file in task directory (current plan state - may not exist on first iteration)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Analyze the bug from task file per AGENTS.md (symptoms, root cause, affected functionality)
- Understand existing spec structure and patterns
- Determine how the spec should be adjusted to drive the fix (acceptance criteria, edge cases, clarifications)
- If plan file exists: critique it (completeness, accuracy, priorities, clarity)
- Identify what spec changes are needed to prevent this class of bug

**DECIDE**
- Generate a complete plan for adjusting specs to drive the bug fix
- Structure the plan by priority (most important spec changes first)
- Break spec adjustments into tight, actionable tasks
- Determine which specs need updates and what acceptance criteria to add

**ACT**
- Write plan file with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

## Prompt Component Set

**OBSERVE (6 variants)**
1. `observe_bootstrap.md` - Repository structure, docs, implementation patterns
2. `observe_plan_specs_impl.md` - AGENTS.md + plan file + specifications + implementation
3. `observe_specs.md` - AGENTS.md + specifications
4. `observe_impl.md` - AGENTS.md + implementation
5. `observe_story_task_specs_impl.md` - AGENTS.md + task file + plan file + specifications + implementation
6. `observe_bug_task_specs_impl.md` - AGENTS.md + task file + plan file + specifications + implementation

**ORIENT (6 variants)**
1. `orient_bootstrap.md` - Identify project type, determine definitions, synthesize understanding
2. `orient_build.md` - Understand task, identify what to build
3. `orient_gap.md` - Compare sources, identify gaps, assess completeness/accuracy
4. `orient_quality.md` - Apply criteria, identify markers, score PASS/FAIL
5. `orient_story_task_incorporation.md` - Analyze story/task, determine incorporation strategy, critique existing plan
6. `orient_bug_task_incorporation.md` - Analyze bug/task, determine spec adjustments to drive fix, critique existing plan

**DECIDE (6 variants)**
1. `decide_bootstrap.md` - Determine AGENTS.md structure and content
2. `decide_build.md` - Pick task, determine approach, identify files
3. `decide_gap_plan.md` - Structure plan, break gaps into tasks, determine dependencies/updates
4. `decide_refactor_plan.md` - If threshold fails: propose refactoring, structure plan, prioritize
5. `decide_story_task_plan.md` - Generate complete plan for story/task incorporation, structure by priority
6. `decide_bug_task_plan.md` - Generate complete plan for spec adjustments to drive bug fix, structure by priority

**ACT (3 variants)**
1. `act_bootstrap.md` - Create AGENTS.md, commit
2. `act_build.md` - Implement, test, update plan file/AGENTS.md, commit if passing
3. `act_plan.md` - Write plan file, update AGENTS.md, commit

**Total: 21 prompt files** (6+6+6+3)

---

## Procedure Compositions

| Procedure | Observe | Orient | Decide | Act |
|------|---------|--------|--------|-----|
| 0. Bootstrap | bootstrap | bootstrap | bootstrap | bootstrap |
| 1. Building from plan | plan_specs_impl | build | build | build |
| 2. Plan spec-to-impl | plan_specs_impl | gap | gap_plan | plan |
| 3. Plan impl-to-spec | plan_specs_impl | gap | gap_plan | plan |
| 4. Plan spec refactoring | specs | quality | refactor_plan | plan |
| 5. Plan impl refactoring | impl | quality | refactor_plan | plan |
| 6. Plan story-to-spec | story_task_specs_impl | story_task_incorporation | story_task_plan | plan |
| 7. Plan bug-to-spec | bug_task_specs_impl | bug_task_incorporation | bug_task_plan | plan |

---

## Principles

- **AGENTS.md always studied first** - Defines what constitutes "specification" and "implementation", plus quality criteria
- **Definitions defer to AGENTS.md** - What files/locations constitute specs and implementation varies by project
- **Don't assume not implemented** - Always search codebase before implementing (critical - the Achilles' heel)
- **Use parallel subagents** - Main agent as scheduler, spawn subagents for work (only 1 for build/tests)
- **Backpressure forces correctness** - Tests must pass before commit
- **Capture the why** - Update AGENTS.md with learnings, keep it up to date
- **Tight tasks** - 1 task per loop = 100% smart zone utilization
- **Commit after updates** - All file modifications complete before commit
- **Composability** - 21 files generate 8 procedure types through different combinations
