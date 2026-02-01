# Ralph Wiggum OODA Loop

Autonomous AI coding that maintains fresh context across iterations using composable OODA prompts.

Evolved from the [Ralph Loop](https://ghuntley.com/ralph/) by Geoff Huntley, applying the OODA framework to break monolithic prompts into composable components.

## TL;DR

```bash
# 1. Bootstrap: create operational guide
./rooda.sh bootstrap

# 2. Plan: incorporate story into specs
./rooda.sh plan-story-to-spec --max-iterations 5

# 3. Build: implement the specs
./rooda.sh build --max-iterations 5

# 4. Refactor implementation
./rooda.sh plan-impl-refactor
./rooda.sh build --max-iterations 5
./rooda.sh plan-impl-to-spec
./rooda.sh build --max-iterations 5

# 5. Refactor specs
./rooda.sh plan-spec-refactor
./rooda.sh build --max-iterations 5
```

Each iteration clears context. File-based memory persists. AI stays in its "smart zone" (40-60% utilization) indefinitely.

**What you get:** Autonomous coding without context degradation. Tests must pass before commit. Eventual consistency through iteration.

## What You Get

**8 Procedure Types**
- Bootstrap your repository (creates operational guide)
- Build from plan (only procedure that modifies code)
- Gap analysis (spec to implementation)
- Quality assessment (refactoring triggers)
- Story/bug incorporation (iterative convergence)

**Composable Architecture**
- Prompt files combine into 8 procedures
- Same components reused in different combinations
- Customize by editing prompts or creating new compositions

**Built-in Quality Control**
- Tests/lints reject invalid work (backpressure)
- Boolean quality criteria trigger refactoring
- "Don't assume not implemented" - always searches first
- Parallel subagents (main agent as scheduler)

## How It Works

### The Loop

```bash
./rooda.sh <procedure> [--max-iterations N]
```

Each iteration:
1. Loads 4 OODA prompt components (observe, orient, decide, act)
2. Injects context (file paths, configuration)
3. Agent executes: **observe > orient > decide > act**
4. Updates files on disk (PLAN.md, AGENTS.md, code)
5. Exits - **context cleared**
6. Loop restarts with fresh context

**Why fresh context?** LLMs degrade as context fills. 200K advertised ≈ 176K usable. At 40-60% utilization, quality stays high. Fresh context each iteration prevents degradation.

### OODA Phases

**Observe** - Gather information
- Study AGENTS.md (how to build/test, what is spec/implementation)
- Study plan file, specs, implementation
- Study story/bug file (if applicable)

**Orient** - Analyze and synthesize
- Gap analysis (compare specs to code)
- Quality assessment (apply boolean criteria)
- Understand task requirements
- Search codebase (don't assume not implemented)

**Decide** - Determine course of action
- Pick most important task from plan
- Structure plan by priority
- Determine implementation approach
- Identify files to modify

**Act** - Execute
- Implement using parallel subagents (only 1 for build/tests)
- Run tests (backpressure)
- Write/update plan file
- Update AGENTS.md (capture the why)
- Commit when tests pass

### Example Iteration

```bash
./rooda.sh build
```

1. **Observe:** Reads AGENTS.md, PLAN.md, specs/, src/
2. **Orient:** "Most important task: implement user authentication. Search shows auth/ directory exists but missing password reset."
3. **Decide:** "Implement password reset flow. Modify auth/reset.go, add tests."
4. **Act:** Spawns subagent to implement, runs tests, updates PLAN.md, commits
5. **Exit:** Context cleared
6. **Loop:** Restarts, picks next task

## The 8 Procedures

### 0. Bootstrap
**When:** First time in a repository, or AGENTS.md doesn't exist  
**What:** Creates AGENTS.md operational guide  
**Output:** AGENTS.md with build/test commands, spec/implementation definitions, quality criteria

```bash
./rooda.sh bootstrap
```

### 1. Build
**When:** You have a plan and want to implement it  
**What:** Implements tasks from PLAN.md (only procedure that modifies code)  
**Output:** Code changes, test runs, updated PLAN.md

```bash
./rooda.sh build --max-iterations 5
```

### 2. Plan Spec-to-Impl
**When:** You have specs and want a plan to implement them  
**What:** Gap analysis - what's in specs but not in code  
**Output:** PLAN.md with prioritized implementation tasks

```bash
./rooda.sh plan-spec-to-impl
```

### 3. Plan Impl-to-Spec
**When:** You have code and want specs to document it  
**What:** Gap analysis - what's in code but not in specs  
**Output:** PLAN.md with prioritized documentation tasks

```bash
./rooda.sh plan-impl-to-spec
```

### 4. Plan Spec Refactoring
**When:** Your specs feel unclear or incomplete  
**What:** Quality assessment using boolean criteria (clarity, completeness, consistency, testability)  
**Output:** PLAN.md with spec refactoring tasks (if criteria fail threshold)

```bash
./rooda.sh plan-spec-refactor
```

### 5. Plan Impl Refactoring
**When:** Your code feels messy or hard to maintain  
**What:** Quality assessment using boolean criteria (cohesion, coupling, complexity, maintainability)  
**Output:** PLAN.md with implementation refactoring tasks (if criteria fail threshold)

```bash
./rooda.sh plan-impl-refactor
```

### 6. Plan Story-to-Spec
**When:** You have a new feature/story to incorporate  
**What:** Iteratively converges on how to incorporate story into specs  
**Output:** PLAN.md with tasks to update/create specs

```bash
# Create TASK.md first
./rooda.sh plan-story-to-spec --max-iterations 5
```

### 7. Plan Bug-to-Spec
**When:** You have a bug and want specs to drive the fix  
**What:** Determines spec adjustments needed (acceptance criteria, edge cases)  
**Output:** PLAN.md with spec changes to prevent bug recurrence

```bash
# Create TASK.md first
./rooda.sh plan-bug-to-spec --max-iterations 3
```

## File-Based State

### AGENTS.md
Operational guide for the repository. Created by bootstrap, updated by all procedures.

**Contains:**
- How to build/run/test (specific commands)
- What constitutes "specification" (file paths/patterns)
- What constitutes "implementation" (file paths/patterns)
- Quality criteria definitions

**Philosophy:** Assumed inaccurate until verified empirically. Updated when errors discovered.

### tasks/{task-id}/ (optional)

Task-specific working directory. This is one possible organizational pattern - projects may use different structures.

**Files:**
- `PLAN.md` - Prioritized task list and progress tracking
- `TASK.md` - Task description (optional, for story/bug procedures)

**Philosophy:** Generated and updated by act phase. Assumed inaccurate until verified. Disposable - regenerate if trajectory goes wrong.

### specs/

Specification documents (optional, see [specs.md](specs.md)).

**Structure:**
- One spec per topic of concern using [spec-template.md](spec-template.md)
- Source of truth for requirements
- Acceptance criteria define backpressure for act phase
- Implementation Mapping bridges specs to code

### prompts/

OODA phase component library.

**Organization:**
- `observe_*.md` - Different data sources
- `orient_*.md` - Different analysis types
- `decide_*.md` - Different decision strategies
- `act_*.md` - Different execution modes

**Composition:** See [prompts/README.md](prompts/README.md) for how these combine into procedures.

## Composable Architecture

### How Procedures Compose

Each procedure = observe + orient + decide + act prompt files

**Example: Build procedure**
```yaml
build:
  observe: prompts/observe_plan_specs_impl.md
  orient: prompts/orient_build.md
  decide: prompts/decide_build.md
  act: prompts/act_build.md
```

**Reuse:** Same `observe_plan_specs_impl.md` used by build, plan-spec-to-impl, and plan-impl-to-spec procedures. Different orient/decide/act create different behaviors.

### Custom Procedures

Create your own by editing `ooda-config.yml`:

```yaml
# Add custom file path patterns
paths:
  task_dir: "tasks/{task-id}"
  task_file: "{task_dir}/TASK.md"
  plan_file: "{task_dir}/PLAN.md"

# Add custom procedures
procedures:
  my-custom-procedure:
    observe: prompts/observe_specs.md
    orient: prompts/orient_quality.md
    decide: prompts/decide_refactor_plan.md
    act: prompts/act_plan.md
    default_iterations: 1
```

Or specify prompts directly:

```bash
./rooda.sh \
  --observe prompts/observe_specs.md \
  --orient prompts/orient_gap.md \
  --decide prompts/decide_gap_plan.md \
  --act prompts/act_plan.md \
  --max-iterations 1
```

### The Prompt Files

**Observe**
1. `observe_bootstrap.md` - Repository structure, docs, patterns
2. `observe_plan_specs_impl.md` - AGENTS.md + plan + specs + implementation
3. `observe_specs.md` - AGENTS.md + specs only
4. `observe_impl.md` - AGENTS.md + implementation only
5. `observe_story_specs_impl.md` - AGENTS.md + story + plan + specs + implementation
6. `observe_bug_specs_impl.md` - AGENTS.md + bug + plan + specs + implementation

**Orient**
1. `orient_bootstrap.md` - Identify project type, determine definitions
2. `orient_build.md` - Understand task, identify what to build
3. `orient_gap.md` - Compare sources, identify gaps
4. `orient_quality.md` - Apply criteria, score PASS/FAIL
5. `orient_story_incorporation.md` - Analyze story, determine incorporation strategy
6. `orient_bug_incorporation.md` - Analyze bug, determine spec adjustments

**Decide**
1. `decide_bootstrap.md` - Determine AGENTS.md structure
2. `decide_build.md` - Pick task, determine approach
3. `decide_gap_plan.md` - Structure plan, break gaps into tasks
4. `decide_refactor_plan.md` - Propose refactoring if threshold fails
5. `decide_story_plan.md` - Generate plan for story incorporation
6. `decide_bug_plan.md` - Generate plan for spec adjustments

**Act**
1. `act_bootstrap.md` - Create AGENTS.md, commit
2. `act_build.md` - Implement, test, update files, commit if passing
3. `act_plan.md` - Write plan file, update AGENTS.md, commit

## Key Principles

### Ralph Loop Language Patterns

These specific phrases matter (from [Ralph Loop](ralph-loop.md) by Geoff Huntley):

- **"study" not "read"** - Active, intentional engagement with code
- **"don't assume not implemented"** - Critical Achilles' heel. Always search codebase first
- **"using parallel subagents"** - Main agent as scheduler, subagents do work
- **"only 1 subagent for build/tests"** - Prevents parallel test conflicts
- **"capture the why, keep it up to date"** - AGENTS.md stores learnings and rationale
- **"most important task"** - Priority-driven execution, not sequential
- **"tight tasks"** - 1 task per loop = 100% smart zone utilization

### Context Management

- 200K tokens advertised ≈ 176K usable
- 40-60% utilization = "smart zone" (high quality output)
- Fresh context each iteration prevents degradation
- Main agent as scheduler, spawn subagents for parallel work
- Tight tasks (1 per loop) maximize smart zone utilization

### Steering via Backpressure

**Downstream (in act phase):**
- Tests must pass before commit
- Lints, type checks reject invalid work
- Build failures prevent progression

**Upstream (in orient phase):**
- Boolean quality criteria (PASS/FAIL scoring)
- Refactoring triggers when criteria fail threshold
- Existing code patterns guide generation

**Result:** Eventual consistency through iteration. System self-corrects.

### Refactoring Triggers

Orient phase applies boolean criteria:

**For specs:** clarity, completeness, consistency, testability  
**For implementation:** cohesion, coupling, complexity, maintainability

Each criterion scored PASS or FAIL. When threshold fails, decide/act write refactoring proposal to PLAN.md. Future build iteration executes it.

**Human markers also trigger refactoring:**
- TODOs, "REFACTORME" comments
- Unclear language in specs
- Long functions, code smells in implementation

## Sample Repository Structure

```
project-root/
├── rooda.sh                   # Loop script
├── ooda-config.yml            # File paths and procedure compositions
├── AGENTS.md                  # Operational guide
├── tasks/ (optional)          # Task-specific working directories (project-specific)
│   └── {task-id}/
│       ├── PLAN.md            # Task list and progress
│       └── TASK.md            # Task description (optional)
├── prompts/                   # OODA phase components
│   ├── observe_*.md           # Observation variants
│   ├── orient_*.md            # Analysis variants
│   ├── decide_*.md            # Decision variants
│   └── act_*.md               # Execution variants
├── specs/                     # Requirements (optional)
│   ├── README.md              # Index of JTBDs, topics, and specs
│   ├── TEMPLATE.md            # Template for new specs
│   └── topic-name.md          # One spec per topic of concern
└── src/                       # Implementation
    └── ...
```

## Safety

Requires `--dangerously-skip-permissions` to run autonomously (bypasses all permission prompts).

**Run in isolated sandbox environments:**
- Docker containers (local)
- Fly Sprites / E2B (remote)
- Minimum viable access (only needed API keys)
- No access to private data beyond requirements

**Philosophy:** "It's not if it gets popped, it's when. And what is the blast radius?"

Limit blast radius through isolation, not through hoping the AI won't do something bad.

## Troubleshooting

### Escape Hatches

- **Max iterations** - Prevents infinite loops (`--max-iterations N`)
- **Ctrl+C** - Stops the loop immediately
- **`git reset --hard`** - Reverts uncommitted changes
- **Regenerate plan** - If trajectory goes wrong, delete PLAN.md and run planning procedure again

### Common Issues

**"Agent keeps implementing the same thing"**
- Check PLAN.md - is it being updated?
- Check AGENTS.md - does it define implementation locations correctly?
- Run bootstrap again to regenerate AGENTS.md

**"Tests keep failing"**
- Check AGENTS.md - are test commands correct?
- Run tests manually to verify they work
- Update AGENTS.md with correct commands

**"Agent doesn't find existing code"**
- This is the Achilles' heel: "don't assume not implemented"
- Check orient prompts - do they emphasize searching?
- Add explicit search instructions to AGENTS.md

**"Plan goes off track"**
- Delete PLAN.md
- Run planning procedure again (cheap to regenerate)
- Adjust specs if needed

## Why It Works

1. **Fresh context** - No degradation from context pollution
2. **Composable prompts** - Reusable components for different task types
3. **OODA framework** - Clear separation of concerns across phases
4. **File-based state** - PLAN.md and AGENTS.md persist learnings
5. **Backpressure** - Tests and criteria force correctness
6. **Eventual consistency** - Iteration converges to solution
7. **Simplicity** - Bash loop, prompt interpolation, file I/O

## Learn More

- [OODA Loop](ooda-loop.md) - The decision-making framework
- [Ralph Loop](ralph-loop.md) - Original methodology by Geoff Huntley
- [Specs System](specs.md) - How to structure specifications
- [Spec Template](spec-template.md) - Template for new specs
- [Prompts README](prompts/README.md) - Detailed prompt composition breakdown

---

*Methodology evolved from [Ralph Loop](https://ghuntley.com/ralph/) by Geoff Huntley*

<a href="https://buymeacoffee.com/Max.dunn" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
