# Minimal Composable Prompt Set

## Task Breakdown by Phase

### Task 0: Bootstrap (create AGENTS.md)

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

### Task 1: Building from plan

**OBSERVE**
- Study AGENTS.md (how to build/test, what is specification/implementation)
- Study PLAN.md (find most important task)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Understand task requirements
- Search codebase (don't assume not implemented)
- Identify what needs to be built/modified
- Determine test strategy

**DECIDE**
- Pick most important task from PLAN.md
- Determine implementation approach using parallel subagents
- Identify which files to modify

**ACT**
- Implement using parallel subagents (only 1 subagent for build/tests)
- Run tests per AGENTS.md (backpressure)
- Update PLAN.md (mark complete/update)
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit when tests pass

---

### Task 2: Plan spec-to-impl

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study PLAN.md (current state)
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
- Write PLAN.md with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Task 3: Plan impl-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study PLAN.md (current state)
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
- Write PLAN.md with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Task 4: Plan spec refactoring

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
- Write PLAN.md with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Task 5: Plan impl refactoring

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
- Write PLAN.md with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

### Task 6: Plan feature-to-spec

**OBSERVE**
- Study AGENTS.md (what is specification/implementation)
- Study FEATURE.md (the feature to be incorporated)
- Study PLAN.md (current plan state - may not exist on first iteration)
- Study specifications (per AGENTS.md definition)
- Study implementation (per AGENTS.md definition, file tree, symbols)

**ORIENT**
- Analyze the feature from FEATURE.md (scope, requirements, integration points)
- Understand existing spec structure and patterns
- Determine how the feature should be incorporated (create new specs, update existing, refactor)
- If PLAN.md exists: critique it (completeness, accuracy, priorities, clarity)
- Identify what tasks are needed for proper incorporation

**DECIDE**
- Generate a complete plan for incorporating the feature into specs
- Structure the plan by priority (most important tasks first)
- Break incorporation into tight, actionable tasks
- Determine task dependencies and which specs need creation/updates

**ACT**
- Write PLAN.md with prioritized bullet-point task list
- Update AGENTS.md if learned something new (capture the why, keep it up to date)
- Commit changes

---

## Prompt Component Set

**OBSERVE (5 variants)**
1. `observe_bootstrap.md` - Repository structure, docs, implementation patterns
2. `observe_plan_specs_impl.md` - AGENTS.md + PLAN.md + specifications + implementation
3. `observe_specs.md` - AGENTS.md + specifications
4. `observe_impl.md` - AGENTS.md + implementation
5. `observe_feature_specs_impl.md` - AGENTS.md + FEATURE.md + PLAN.md + specifications + implementation

**ORIENT (5 variants)**
1. `orient_bootstrap.md` - Identify project type, determine definitions, synthesize understanding
2. `orient_build.md` - Understand task, identify what to build
3. `orient_gap.md` - Compare sources, identify gaps, assess completeness/accuracy
4. `orient_quality.md` - Apply criteria, identify markers, score PASS/FAIL
5. `orient_feature_incorporation.md` - Analyze feature, determine incorporation strategy, critique existing plan

**DECIDE (5 variants)**
1. `decide_bootstrap.md` - Determine AGENTS.md structure and content
2. `decide_build.md` - Pick task, determine approach, identify files
3. `decide_gap_plan.md` - Structure plan, break gaps into tasks, determine dependencies/updates
4. `decide_refactor_plan.md` - If threshold fails: propose refactoring, structure plan, prioritize
5. `decide_feature_plan.md` - Generate complete plan for feature incorporation, structure by priority

**ACT (3 variants)**
1. `act_bootstrap.md` - Create AGENTS.md, commit
2. `act_build.md` - Implement, test, update PLAN.md/AGENTS.md, commit if passing
3. `act_plan.md` - Write PLAN.md, update AGENTS.md, commit

**Total: 18 prompt files** (5+5+5+3)

---

## Task Compositions

| Task | Observe | Orient | Decide | Act |
|------|---------|--------|--------|-----|
| 0. Bootstrap | bootstrap | bootstrap | bootstrap | bootstrap |
| 1. Building from plan | plan_specs_impl | build | build | build |
| 2. Plan spec-to-impl | plan_specs_impl | gap | gap_plan | plan |
| 3. Plan impl-to-spec | plan_specs_impl | gap | gap_plan | plan |
| 4. Plan spec refactoring | specs | quality | refactor_plan | plan |
| 5. Plan impl refactoring | impl | quality | refactor_plan | plan |
| 6. Plan feature-to-spec | feature_specs_impl | feature_incorporation | feature_plan | plan |

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
- **Composability** - 18 files generate 7 task types through different combinations
