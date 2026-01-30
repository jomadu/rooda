# Ralph Wiggum OODA Loop

An autonomous AI coding methodology using composable [OODA-based](ooda-loop.md) prompts to maintain fresh context across iterations.

## Core Concept

Each loop iteration: **observe > orient > decide > act > clear context > repeat**

Fresh context each iteration keeps the AI in its "smart zone" (40-60% context utilization). File-based memory (AGENTS.md, PLAN.md) persists learnings between iterations.

This methodology evolved from the [Ralph Loop](ralph-loop.md) by Geoff Huntley. It applies the OODA framework to break the monolithic prompt into discrete phases, creating composable prompt components.

## The Loop Mechanism

```bash
./ooda.sh <procedure> --task TASK-123 \
          [--max-iterations N]
```

Or with explicit prompt files:

```bash
./ooda.sh --observe prompts/observe_X.md \
          --orient prompts/orient_Y.md \
          --decide prompts/decide_Z.md \
          --act prompts/act_W.md \
          --task TASK-123 \
          [--max-iterations N]
```

The script interpolates the 4 prompt components into a template and feeds it to an LLM agent. Each iteration:
1. Loads prompt template with 4 OODA phase components
2. Injects task-specific context (task ID, file paths)
3. Agent executes observe > orient > decide > act
4. Updates task files on disk
5. Exits (context cleared)
6. Loop restarts with fresh context

Exits when max iterations reached.

## OODA Phase Responsibilities

### Planning Procedures (Procedures 2-7)

1. **Observe** - Gather information from specs, implementation, plan file, story/bug file (if applicable), and AGENTS.md
2. **Orient** - Analyze observations using procedure-specific criteria and synthesize understanding
3. **Decide** - Determine plan structure, priorities, tasks for plan file, and necessary AGENTS.md updates
4. **Act** - Write the plan to plan file and update AGENTS.md

### Building Procedures (Procedure 1)

1. **Observe** - Gather information from plan file, AGENTS.md, specs, and implementation
2. **Orient** - Understand task requirements, search codebase (don't assume not implemented), identify what needs to be built
3. **Decide** - Pick most important task from plan file, determine implementation approach using parallel subagents, identify AGENTS.md updates
4. **Act** - Implement using parallel subagents (only 1 for build/tests), run tests (backpressure), update plan file and AGENTS.md, commit when passing

## Procedure Types

The methodology supports multiple procedure types through prompt composition:

0. **Bootstrap** - Create AGENTS.md operational guide
   - Studies repository structure, tech stack, build files
   - Determines what constitutes "specification" vs "implementation"
   - Identifies build/test/run commands empirically
   - Run first if AGENTS.md doesn't exist

1. **Building from plan** - Implement tasks from plan file
   - Only procedure type that modifies implementation code
   - Uses parallel subagents (only 1 for build/tests)
   - Backpressure from tests ensures correctness

2. **Plan spec-to-impl** - Create plan to make implementation match specifications
   - Gap analysis: what's in specs but not in code
   - Searches codebase (don't assume not implemented)

3. **Plan impl-to-spec** - Create plan to make specifications match implementation
   - Gap analysis: what's in code but not in specs
   - Searches codebase thoroughly

4. **Plan spec refactoring** - Create plan to refactor specs out of local optimums
   - Orient applies boolean criteria (clarity, completeness, consistency, testability, human markers)
   - Triggers on threshold failures
   - Proposes refactoring in plan file, doesn't execute

5. **Plan impl refactoring** - Create plan to refactor implementation out of local optimums
   - Orient applies boolean criteria (cohesion, coupling, complexity, maintainability, human markers)
   - Triggers on threshold failures
   - Proposes refactoring in plan file, doesn't execute

6. **Plan story-to-spec** - Create plan to incorporate new story into specifications
   - Iteratively converges on proper incorporation
   - Analyzes story file and existing specs/implementation
   - Each iteration critiques and improves the plan
   - Runs until plan stabilizes or max iterations reached

7. **Plan bug-to-spec** - Create plan to adjust specs to drive bug fix
   - Analyzes bug symptoms and root cause
   - Determines spec adjustments needed (acceptance criteria, edge cases, clarifications)
   - Focuses on preventing bug recurrence through better specs
   - Runs until plan stabilizes or max iterations reached

## Key Principles

### Composable Prompts
- Minimal yet complete set of prompt variants per phase (21 files total)
- Most variants in observe (different data sources) and orient (different analysis types)
- Decide/act more stable across task types
- Same orient variant can be reused with different observe inputs
- See [prompts/README.md](prompts/README.md) for detailed breakdown

### Ralph Loop Language Patterns
- "study" not "read" - Active, intentional engagement
- "don't assume not implemented" - Critical Achilles' heel, always search first
- "using parallel subagents" - Main agent as scheduler (only 1 for build/tests)
- "capture the why, keep it up to date" - AGENTS.md learnings
- "most important task" - Priority-driven execution
- "tight tasks" - 1 task per loop = 100% smart zone utilization

### File-Based State

**AGENTS.md** - Operational guide for the repository
- How to build/run/test the project
- Definition of what constitutes "specification" vs "implementation"
- Created by orient phase if missing
- Assumed inaccurate/incomplete until verified empirically
- Updated in the case of discovered errors

**tasks/{task-id}/** - Task-specific working directory
- `PLAN.md` - Prioritized task list and progress tracking
- `STORY.md` - Story/feature description (for feature work)
- `BUG.md` - Bug description (for bug fixes)
- Generated and updated by act phase
- Assumed inaccurate until verified empirically
- Updated in the case of discovered errors

**specs/** - Specification documents (optional, see [specs.md](specs.md))
- One spec per topic of concern using [spec-template.md](spec-template.md)
- Source of truth for requirements
- Acceptance criteria define backpressure for act phase
- Implementation Mapping bridges specs ↔ code for gap analysis

**prompts/** - OODA phase component library (21 files)
- `prompts/observe_*.md` - Different observation sources (6 variants)
- `prompts/orient_*.md` - Different analysis types (6 variants)
- `prompts/decide_*.md` - Different decision strategies (6 variants)
- `prompts/act_*.md` - Different execution modes (3 variants)
- The loop script injects task-specific file paths into prompts

## Sample Repository Structure

```
project-root/
├── ooda.sh                    # Loop script
├── ooda-procedures.yml        # Procedure compositions
├── AGENTS.md                  # Operational guide
├── tasks/                     # Task-specific working directories
│   └── {task-id}/
│       ├── PLAN.md            # Task list and progress
│       ├── STORY.md           # Story/feature description (optional)
│       └── BUG.md             # Bug description (optional)
├── prompts/                   # OODA phase components
│   ├── observe_*.md           # Observation variants
│   ├── orient_*.md            # Analysis variants
│   ├── decide_*.md            # Decision variants
│   └── act_*.md               # Execution variants
├── specs/                     # Requirements (if using spec-driven approach)
│   ├── README.md              # Index of JTBDs, topics, and specs
│   ├── TEMPLATE.md            # Template for new specs
│   └── topic-name.md          # One spec per topic of concern
└── src/                       # Implementation
    └── ...
```

### Context Management
- 200K tokens advertised ≈ 176K usable
- 40-60% utilization = "smart zone"
- Fresh context each iteration prevents degradation
- Use main agent as scheduler, spawn subagents for parallel work

### Steering via Backpressure
- Tests, lints, type checks reject invalid work (in act phase)
- Refactoring criteria provide quality gates (in orient phase)
- Existing code patterns guide generation
- Eventual consistency through iteration

### Refactoring Triggers
Boolean criteria scored as PASS/FAIL in orient phase:
- Quality metrics (cohesion, coupling, complexity, completeness, etc)
- Human markers (TODOs, comments, "REFACTORME", spec phrases)
- Custom criteria defined in prompt variants

When criteria fail threshold, decide/act write refactoring proposal to plan file. Future iteration with building task executes it.

## Why It Works

1. **Fresh context** - No degradation from context pollution
2. **Composable prompts** - Reusable components for different task types
3. **OODA framework** - Clear separation of concerns across phases
4. **File-based state** - Plan file and AGENTS.md persist learnings
5. **Backpressure** - Tests and criteria force correctness
6. **Eventual consistency** - Iteration converges to solution
7. **Simplicity** - Bash loop, prompt interpolation, file I/O

## Safety

Requires `--dangerously-skip-permissions` to run autonomously. Run in isolated sandbox environments:
- Docker containers (local)
- Fly Sprites / E2B (remote)
- Minimum viable access (only needed API keys)
- No access to private data beyond requirements

Philosophy: "It's not if it gets popped, it's when. And what is the blast radius?"

## Escape Hatches

- Max iterations prevents infinite loops
- Ctrl+C stops the loop
- `git reset --hard` reverts uncommitted changes
- Regenerate plan file if trajectory goes wrong

---

*Methodology evolved from [Ralph Loop](https://ghuntley.com/ralph/) by Geoff Huntley*
