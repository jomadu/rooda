# Ralph Wiggum OODA

An autonomous AI coding methodology using composable [OODA-based](ooda-loop.md) prompts to maintain fresh context across iterations.

## Core Concept

Each loop iteration: **observe → orient → decide → act → clear context → repeat**

Fresh context each iteration keeps the AI in its "smart zone" (40-60% context utilization). File-based memory (AGENTS.md, PLAN.md) persists learnings between iterations.

This methodology evolved from the [Ralph Loop](ralph-loop.md) by Geoff Huntley, applying the OODA framework to create composable prompt components.

## The Loop Mechanism

```bash
./ooda.sh --observe prompts/observe_X.md \
          --orient prompts/orient_Y.md \
          --decide prompts/decide_Z.md \
          --act prompts/act_W.md \
          [--max-iterations N]
```

The script interpolates the 4 prompt components into a template and feeds it to an LLM agent. Each iteration:
1. Loads prompt template with 4 OODA phase components
2. Agent executes observe → orient → decide → act
3. Updates PLAN.md on disk
4. Exits (context cleared)
5. Loop restarts with fresh context

Exits when max iterations reached.

## OODA Phase Responsibilities

### Planning Tasks (Tasks 2-5)

1. **Observe** - Gather information from specs, implementation, PLAN.md, and AGENTS.md
2. **Orient** - Analyze observations using task-specific criteria and synthesize understanding
3. **Decide** - Determine plan structure, priorities, and tasks for PLAN.md
4. **Act** - Write the plan to PLAN.md

### Building Tasks (Task 1)

1. **Observe** - Gather information from PLAN.md, AGENTS.md, specs, and implementation
2. **Orient** - Understand task requirements and identify what needs to be built
3. **Decide** - Pick highest priority task and determine implementation approach
4. **Act** - Implement task, run quality checks per AGENTS.md, commit when passing, update PLAN.md

## Task Types

The methodology supports multiple task types through prompt composition:

1. **Building from plan** - Implement tasks from PLAN.md
   - Only task type that modifies implementation code
   - Backpressure from tests ensures correctness

2. **Plan spec→impl** - Create plan to make implementation match specifications
   - Gap analysis: what's in specs but not in code

3. **Plan impl→spec** - Create plan to make specifications match implementation
   - Gap analysis: what's in code but not in specs

4. **Plan spec refactoring** - Create plan to refactor specs out of local optimums
   - Orient applies boolean criteria (e.g. clarity, completeness, consistency, testability, human markers)
   - Triggers on threshold failures
   - Proposes refactoring in PLAN.md, doesn't execute

5. **Plan impl refactoring** - Create plan to refactor implementation out of local optimums
   - Orient applies boolean criteria (e.g. cohesion, coupling, complexity, maintainability, human markers)
   - Triggers on threshold failures
   - Proposes refactoring in PLAN.md, doesn't execute

## Key Principles

### Composable Prompts
- Minimal yet complete set of prompt variants per phase
- Most variants in observe (different data sources) and orient (different analysis types)
- Decide/act more stable across task types
- Same orient variant can be reused with different observe inputs

### File-Based State

**AGENTS.md** - Operational guide for the repository
- How to build/run/test the project
- Definition of what constitutes "specification" vs "implementation"
- Created by orient phase if missing
- Assumed inaccurate/incomplete until verified empirically
- Updated in the case of discovered errors

**PLAN.md** - Prioritized task list and progress tracking
- Generated and updated by act phase
- Can contain refactoring proposals with criteria scores
- Assumed inaccurate until verified empirically
- Updated in the case of discovered errors

**specs/** - Specification documents (optional, see [specs.md](specs.md))
- One spec per topic of concern using [spec-template.md](spec-template.md)
- Source of truth for requirements
- Acceptance criteria define backpressure for act phase
- Implementation Mapping bridges specs ↔ code for gap analysis

**prompts/** - OODA phase component library
- `prompts/observe_*.md` - Different observation sources
- `prompts/orient_*.md` - Different analysis types
- `prompts/decide_*.md` - Different decision strategies
- `prompts/act_*.md` - Different execution modes

## Sample Repository Structure

```
project-root/
├── ooda.sh                    # Loop script
├── AGENTS.md                  # Operational guide (generated/verified by orient)
├── PLAN.md                    # Task list and progress (generated/updated by act)
├── prompts/                   # OODA phase components
│   ├── observe_specs.md
│   ├── observe_implementation.md
│   ├── observe_gaps.md
│   ├── orient_gap_analysis.md
│   ├── orient_refactoring.md
│   ├── decide_planning.md
│   ├── decide_building.md
│   ├── act_plan.md
│   └── act_build.md
├── specs/                     # Requirements (if using spec-driven approach)
│   ├── feature-a.md
│   └── feature-b.md
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

When criteria fail threshold, decide/act write refactoring proposal to PLAN.md. Future iteration with building task executes it.

## Why It Works

1. **Fresh context** - No degradation from context pollution
2. **Composable prompts** - Reusable components for different task types
3. **OODA framework** - Clear separation of concerns across phases
4. **File-based state** - PLAN.md and AGENTS.md persist learnings
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
- Regenerate PLAN.md if trajectory goes wrong

---

*Methodology evolved from [Ralph Loop](https://ghuntley.com/ralph/) by Geoff Huntley*
