# Ralph Wiggum OODA Loop

Autonomous AI coding that maintains fresh context across iterations—no degradation, just eventual consistency through composable prompts and empirical feedback.

## Installation

```bash
# Clone the repository
git clone https://github.com/jomadu/ralph-wiggum-ooda.git

# Copy necessary files to your project
cp ralph-wiggum-ooda/rooda.sh .
cp ralph-wiggum-ooda/ooda-config.yml .
cp -r ralph-wiggum-ooda/prompts .
chmod +x rooda.sh
```

## Basic Workflow

```bash
# 1. Bootstrap: create operational guide
./rooda.sh bootstrap                              # Create AGENTS.md

# 2. Plan: incorporate story into specs
./rooda.sh plan-story-to-spec --max-iterations 5  # Incorporate story into specs

# 3. Build: write the specs
./rooda.sh build --max-iterations 5               # Write specs

# 4. Plan: gap analysis from specs to implementation
./rooda.sh plan-spec-to-impl                      # Gap analysis: specs to code

# 5. Build: implement the specs
./rooda.sh build --max-iterations 5               # Implement code

# 6. Refactor implementation
./rooda.sh plan-impl-refactor                     # Quality assessment of code
./rooda.sh build --max-iterations 5               # Refactor code

# 7. Refactor specs
./rooda.sh plan-spec-refactor                     # Quality assessment of specs
./rooda.sh build --max-iterations 5               # Refine specs
```

## How It Works

The system runs as a bash loop: each iteration loads prompt files, executes them through your AI CLI, updates files on disk, then exits—clearing context completely. This fresh-context-per-iteration approach prevents LLM degradation, keeping the AI in its "smart zone" (40-60% utilization) indefinitely. File-based state (AGENTS.md, work tracking, specs, code) persists across iterations, providing memory without context pollution.

Each iteration follows the OODA framework: observe, orient, decide, act. These four phases are implemented as composable prompt files that combine into different procedures via configuration. Quality control happens through backpressure—tests and lints reject invalid work downstream, while boolean criteria trigger refactoring upstream. The result is eventual consistency: iteration converges to solution through empirical feedback.

### The Loop

A single iteration starts by loading four prompt files (one for each OODA phase), combining them into a single prompt, and piping it to your AI CLI tool. The agent reads files, analyzes the situation, makes decisions, and executes changes—updating code, work tracking, or AGENTS.md as needed. When the iteration completes, the script exits completely, clearing all context from the AI's memory.

This exit-and-restart pattern is critical. LLMs advertise 200K token windows but degrade in quality as context fills—usable capacity is closer to 176K, and performance drops significantly beyond 60% utilization. By clearing context each iteration, the AI stays perpetually in its "smart zone" (40-60% utilization) where output quality remains high. File-based state provides continuity: AGENTS.md, work tracking, specs, and code persist on disk, giving the next iteration everything it needs without carrying forward conversational baggage.

You control iteration count with `--max-iterations N` or stop manually with Ctrl+C. Each procedure defines a sensible default (bootstrap runs once, build runs five times), but you can override as needed.

### OODA Phases

The OODA loop (Observe, Orient, Decide, Act) is a decision-making framework developed by military strategist John Boyd. Breaking the monolithic prompt into these four phases creates clear separation of concerns and enables prompt reuse across different procedures.

**Observe** - Gather information from the environment. The agent reads relevant files and data sources needed for the specific procedure—this might include AGENTS.md, specs, implementation, work tracking, or task descriptions.

**Orient** - Analyze and synthesize the observations. The agent processes what it observed, identifying patterns, gaps, or issues. This phase varies widely by procedure: gap analysis, quality assessment, task understanding, or incorporation strategy.

**Decide** - Determine the course of action. Based on the orientation, the agent makes decisions about what to do next: which task to tackle, how to structure a plan, what approach to take, or which files to modify.

**Act** - Execute the decision. The agent carries out the chosen action: implementing code, writing plans, updating AGENTS.md, running tests, or committing changes. The specific actions depend on the procedure type.

### AGENTS.md: The Agent-Project Interface

AGENTS.md is the interface between agents and your repository. It defines how agents interact with project-specific workflows, tools, and conventions. Every procedure starts by reading AGENTS.md—it's the source of truth for operational details.

**Required sections:**
- **Work tracking system** - How to query ready work, update status, mark complete
- **Build/test/lint commands** - Specific commands to run
- **Specification definition** - What files/patterns constitute specs
- **Implementation definition** - What files/patterns constitute code
- **Quality criteria** - Boolean triggers for refactoring
- **Task/story/bug locations** - Where to find descriptions for incorporation procedures

AGENTS.md is a living document, assumed inaccurate until verified empirically. When agents discover that commands fail, file paths are wrong, or quality criteria don't match project needs, they update AGENTS.md immediately—capturing not just what changed, but why. The bootstrap procedure creates it initially through repository analysis; all subsequent procedures maintain it. Think of it as operational memory that accumulates learnings across iterations, keeping agents aligned with how your project actually works.

### Example Iteration

**Build Cycle:**

1. **Observe** - Reads AGENTS.md to understand work tracking and build commands, queries work tracking system which shows "implement user authentication" as ready, checks specs for authentication requirements, examines existing implementation
2. **Orient** - Searches codebase and finds auth/ directory exists but password reset functionality is missing, understands the gap between spec requirements and current implementation
3. **Decide** - Picks password reset as the most important task, determines approach (modify auth/reset.go, add tests), identifies files that need changes
4. **Act** - Spawns subagent to implement password reset, runs tests per AGENTS.md commands, updates work tracking to mark task complete, commits changes

**Planning Cycle:**

1. **Observe** - Reads AGENTS.md to understand spec and implementation definitions, studies all files in specs/ directory, examines implementation file tree and symbols
2. **Orient** - Performs gap analysis comparing specs to code, finds specs describe email verification feature but no corresponding implementation exists
3. **Decide** - Structures plan with email verification as priority task, breaks it into implementable subtasks (email service integration, verification token generation, verification endpoint)
4. **Act** - Writes prioritized task list to work tracking system per AGENTS.md, updates AGENTS.md if discovered new patterns, commits plan

## Composable Architecture

The system uses composable prompt files to create different procedures. Each OODA phase (observe, orient, decide, act) is implemented as a separate markdown file in the `prompts/` directory. Procedures are defined in `ooda-config.yml` by specifying which four prompt files to combine—one for each phase.

This composition enables significant reuse. For example, `observe_plan_specs_impl.md` (which reads AGENTS.md, work tracking, specs, and implementation) is shared by the `build` procedure, `plan-spec-to-impl` procedure, and `plan-impl-to-spec` procedure. They differ only in their orient, decide, and act components. The build procedure orients around understanding tasks and implements code, while the planning procedures orient around gap analysis and write plans.

The configuration-driven approach means you can create custom procedures without writing new code—just specify which existing prompt components to combine. Want a procedure that observes only specs, orients around quality assessment, decides on refactoring, and acts by writing a plan? Map those four files in the config. The separation of concerns across OODA phases makes different combinations naturally express different task types.

## The 8 Procedures

| Procedure                       | ID                   | Description                                                        | Modifies Code | Iterations |
| ------------------------------- | -------------------- | ------------------------------------------------------------------ | ------------- | ---------- |
| Bootstrap Repository            | `bootstrap`          | Creates or updates AGENTS.md operational guide for the repository  | No            | 1          |
| Build from Plan                 | `build`              | Implements tasks from plan (only procedure that modifies code)     | Yes           | 5          |
| Plan Spec to Implementation     | `plan-spec-to-impl`  | Creates plan from gap analysis - what's in specs but not in code   | No            | 1          |
| Plan Implementation to Spec     | `plan-impl-to-spec`  | Creates plan from gap analysis - what's in code but not in specs   | No            | 1          |
| Plan Spec Refactoring           | `plan-spec-refactor` | Creates refactoring plan from quality assessment of specs          | No            | 1          |
| Plan Implementation Refactoring | `plan-impl-refactor` | Creates refactoring plan from quality assessment of implementation | No            | 1          |
| Plan Story to Spec              | `plan-story-to-spec` | Creates plan for incorporating story into specs                    | No            | 5          |
| Plan Bug to Spec                | `plan-bug-to-spec`   | Creates plan for spec adjustments needed to drive bug fix          | No            | 3          |

## Custom Procedures

You can create custom procedures by editing `ooda-config.yml` or using command-line flags. Each procedure needs four prompt files (one per OODA phase) and optionally a default iteration count.

**Adding to config:**

```yaml
procedures:
  my-custom-procedure:
    display: "My Custom Procedure"
    summary: "Brief description"
    observe: prompts/observe_specs.md
    orient: prompts/orient_gap.md
    decide: prompts/decide_gap_plan.md
    act: prompts/act_plan.md
    default_iterations: 1
```

Then run with `./rooda.sh my-custom-procedure`.

**Using command-line flags:**

```bash
./rooda.sh \
  --observe prompts/observe_specs.md \
  --orient prompts/orient_gap.md \
  --decide prompts/decide_gap_plan.md \
  --act prompts/act_plan.md \
  --max-iterations 1
```

**Example: Project-specific procedure to migrate plan to beads**

Say you've run a planning procedure that wrote tasks to a plan file, but you want to migrate those tasks into beads for better tracking. Create a custom procedure that observes the plan file and AGENTS.md, orients around understanding the task structure, decides how to map tasks to beads issues, and acts by running `bd create` commands for each task.

## Key Principles

**AGENTS.md as Source of Truth** - AGENTS.md defines how agents interact with the project, but it's assumed inaccurate until verified empirically. When commands fail, file paths are wrong, or quality criteria don't match reality, agents update AGENTS.md immediately—capturing not just what changed, but why. This creates operational memory that accumulates learnings across iterations.

**Don't Assume Not Implemented** - This is the critical Achilles' heel. Before implementing anything, agents must search the codebase thoroughly to verify what exists. Assumptions lead to duplicate work and wasted iterations. The orient phase emphasizes this: search first, then decide.

**Backpressure for Quality Control** - Quality enforcement happens in two directions. Downstream (in the act phase), tests must pass before commit—lints, type checks, and builds reject invalid work. Upstream (in the orient phase), boolean quality criteria trigger refactoring when thresholds fail. The result is eventual consistency: iteration converges to solution through empirical feedback.

**Eventual Consistency Through Iteration** - Trust the loop to self-correct. Plans are disposable—cheap to regenerate when trajectory goes wrong. The system doesn't need to be perfect on the first iteration; it needs to improve each iteration. LLMs can self-identify, self-correct, and self-improve when given clear feedback through backpressure.

**Parallel Subagents** - The main agent acts as scheduler, spawning subagents for parallel work. This prevents context bloat in the main loop while handling auxiliary tasks efficiently. Exception: only one subagent for build/tests to avoid parallel test conflicts.

**Capture the Why** - When updating AGENTS.md or work tracking, include rationale. Why this command instead of another? Why these file paths? What was learned that prompted the update? This context helps future iterations understand the reasoning behind decisions.

**Priority-Driven Execution** - Each iteration picks the most important task from work tracking, not the next sequential task. This ensures the loop always works on what matters most. Tight tasks (one per loop) maximize smart zone utilization—the agent stays focused and effective.

## Workflow Patterns

Common sequences showing when to use which procedure.

**Greenfield Project (Starting from Scratch)**
```bash
./rooda.sh bootstrap                              # Create AGENTS.md
./rooda.sh plan-story-to-spec --max-iterations 5  # Incorporate story into specs
./rooda.sh build --max-iterations 5               # Write specs
./rooda.sh plan-spec-to-impl                      # Gap analysis: specs to code
./rooda.sh build --max-iterations 5               # Implement code
```
When: New project, no existing code or specs.

**Brownfield Project (Existing Code, No Specs)**
```bash
./rooda.sh bootstrap                              # Create AGENTS.md
./rooda.sh plan-impl-to-spec                      # Gap analysis: code to specs
./rooda.sh build --max-iterations 5               # Write specs
./rooda.sh plan-spec-refactor                     # Quality assessment of specs
./rooda.sh build --max-iterations 5               # Refine specs
```
When: Legacy codebase needs documentation.

**Feature Development**
```bash
./rooda.sh plan-story-to-spec --max-iterations 5  # Incorporate story into specs
./rooda.sh build --max-iterations 5               # Write specs
./rooda.sh plan-spec-to-impl                      # Gap analysis: specs to code
./rooda.sh build --max-iterations 5               # Implement code
```
When: Adding new functionality to existing project.

**Bug Fix**
```bash
./rooda.sh plan-bug-to-spec --max-iterations 3    # Adjust specs to drive fix
./rooda.sh build --max-iterations 5               # Update specs
./rooda.sh plan-spec-to-impl                      # Gap analysis: specs to code
./rooda.sh build --max-iterations 5               # Fix code
```
When: Bug reveals gap in specifications.

**Refactoring Cycle**
```bash
./rooda.sh plan-impl-refactor                     # Quality assessment of code
./rooda.sh build --max-iterations 5               # Refactor code
./rooda.sh plan-impl-to-spec                      # Gap analysis: code to specs
./rooda.sh build --max-iterations 5               # Update specs
```
When: Code quality degrades, needs cleanup.

**Continuous Improvement**
```bash
./rooda.sh plan-spec-refactor                     # Quality assessment of specs
./rooda.sh build --max-iterations 5               # Refine specs
./rooda.sh plan-impl-refactor                     # Quality assessment of code
./rooda.sh build --max-iterations 5               # Refactor code
```
When: Regular maintenance, keeping both specs and implementation quality high.

## Sample Repository Structure

```
project-root/
├── rooda.sh                   # Loop script
├── ooda-config.yml            # Procedure definitions and file paths
├── AGENTS.md                  # Operational guide (created by bootstrap)
├── prompts/                   # OODA phase components
│   ├── observe_*.md           # Observation variants
│   ├── orient_*.md            # Analysis variants
│   ├── decide_*.md            # Decision variants
│   └── act_*.md               # Execution variants
├── specs/                     # Requirements (optional)
│   ├── README.md              # Index of specs
│   ├── TEMPLATE.md            # Template for new specs
│   └── topic-name.md          # One spec per topic of concern
└── src/                       # Implementation
    └── ...
```

## Safety

The loop requires `--dangerously-skip-permissions` flag in your AI CLI to run autonomously, bypassing all permission prompts. This is inherently risky.

**Run in isolated sandbox environments:**
- Docker containers (local isolation)
- Fly Sprites / E2B (remote sandboxes)
- Minimum viable access (only needed API keys)
- No access to private data beyond requirements

**Philosophy:** "It's not if it gets popped, it's when. And what is the blast radius?"

Limit blast radius through isolation, not through hoping the AI won't do something bad. The loop will modify files, run commands, and commit changes—make sure it can't access anything you can't afford to lose.

## Troubleshooting

**"Agent keeps implementing the same thing"**
- Check work tracking system—is it being updated after each iteration?
- Verify AGENTS.md defines implementation locations correctly
- Run bootstrap again to regenerate AGENTS.md from current repository state

**"Tests keep failing"**
- Verify AGENTS.md has correct test commands
- Run tests manually to confirm they work
- Update AGENTS.md with correct commands and any required setup

**"Agent doesn't find existing code"**
- This is the Achilles' heel: agents must search before implementing
- Check if orient prompts emphasize searching the codebase
- Verify AGENTS.md defines implementation locations accurately

**"Plan goes off track"**
- Plans are disposable—run the planning procedure again
- Adjust specs if they're unclear or incomplete
- Check if AGENTS.md quality criteria need refinement

**"Loop runs forever"**
- Use `--max-iterations N` to set a limit
- Press Ctrl+C to stop immediately
- Check if work tracking is being marked complete

## Learn More

- [OODA Loop](docs/ooda-loop.md) - The decision-making framework
- [Ralph Loop](docs/ralph-loop.md) - Original methodology by Geoff Huntley
- [Specs System](docs/specs.md) - How to structure specifications
- [Spec Template](docs/spec-template.md) - Template for new specs
- [AGENTS.md Specification](docs/agents-md-specification.md) - Complete AGENTS.md format
- [Prompts README](prompts/README.md) - Detailed prompt composition breakdown

---

*Methodology evolved from [Ralph Loop](https://ghuntley.com/ralph/) by Geoff Huntley*
