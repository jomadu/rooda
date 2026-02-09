# Procedures

## Job to be Done

Define the 16 built-in procedures that ship as defaults — their OODA phase compositions, iteration limits, and use cases.

The developer wants to invoke procedures by name without understanding their internal OODA composition, trust that iteration limits are sensible for each procedure's purpose, and extend the framework with custom procedures by following the same composition pattern.

## Activities

1. Define procedure metadata (name, description, use case)
2. Specify OODA phase composition (which prompt files for observe, orient, decide, act)
3. Set default iteration limits (max iterations or unlimited mode)
4. Map prompt files to embedded resources or filesystem paths
5. Validate procedure definitions at framework startup
6. Expose procedures through CLI and configuration system

## Acceptance Criteria

- [ ] 16 built-in procedures defined with complete metadata
- [ ] Each procedure specifies all four OODA phase prompt files
- [ ] Each procedure has a default iteration limit appropriate for its purpose
- [ ] Direct-action procedures (agents-sync, build, publish-plan) have iteration limits
- [ ] Audit procedures have iteration limits (typically 1-3 iterations)
- [ ] Planning procedures have iteration limits (typically 3-10 iterations)
- [ ] Procedure names follow kebab-case convention
- [ ] Procedure descriptions are concise (one sentence)
- [ ] Prompt files for built-in procedures are embedded in the Go binary
- [ ] Custom procedures can reference filesystem prompt files
- [ ] Procedure definitions validate at startup (prompt files exist, iteration limits valid)
- [ ] Invalid procedure definitions produce clear error messages
- [ ] `rooda --list-procedures` displays all procedures with descriptions
- [ ] Procedures are grouped by category (direct-action, audit, planning)
- [ ] Procedure prerequisites documented (AGENTS.md, TASK.md, PLAN.md requirements)
- [ ] Framework validates prerequisites before procedure execution

## Data Structures

### Procedure

Defined in configuration.md, repeated here for reference:

```go
type Procedure struct {
    Name                 string        // Procedure name (kebab-case)
    Description          string        // One-sentence description
    ObserveFile          string        // Path to observe prompt file
    OrientFile           string        // Path to orient prompt file
    DecideFile           string        // Path to decide prompt file
    ActFile              string        // Path to act prompt file
    DefaultMaxIterations *int          // Default max iterations (nil = inherit from loop)
    IterationMode        IterationMode // Iteration mode (nil = inherit from loop)
    AICmd                string        // AI command override (optional)
    AICmdAlias           string        // AI command alias override (optional)
    Category             string        // Category (direct-action, audit, planning)
}
```

### ProcedureCategory

```go
const (
    CategoryDirectAction = "direct-action"
    CategoryAudit        = "audit"
    CategoryPlanning     = "planning"
)
```

## Algorithm

### Procedure Definition

Each built-in procedure is defined with:

```yaml
procedures:
  <procedure-name>:
    description: "<one-sentence description>"
    observe: "prompts/<observe-file>.md"
    orient: "prompts/<orient-file>.md"
    decide: "prompts/<decide-file>.md"
    act: "prompts/<act-file>.md"
    default_max_iterations: <n>
    category: "<category>"
```

### Procedure Validation

At framework startup:

```
For each procedure:
  1. Validate name is non-empty and kebab-case
  2. Validate description is non-empty
  3. Validate all four OODA phase files are specified
  4. Validate prompt files exist (embedded or filesystem)
  5. Validate default_max_iterations >= 1 (if set)
  6. Validate category is valid (direct-action, audit, planning)
  7. If validation fails: report error and exit
```

### Procedure Invocation

When user runs `rooda <procedure-name>`:

```
1. Look up procedure by name in merged configuration
2. If not found: error "Unknown procedure"
3. Load OODA phase prompt files
4. Resolve iteration limit (CLI flag > procedure default > loop default)
5. Execute iteration loop with assembled prompts
```

## Built-in Procedures

### Procedure Prerequisites

All procedures require AGENTS.md to exist and be valid. The framework validates AGENTS.md sections before execution:

- **All procedures** require AGENTS.md with all 10 required sections
- **Planning procedures** (`draft-plan-*`) require AGENTS.md Task Input section to specify TASK.md location
- **publish-plan** requires AGENTS.md Planning System section to specify PLAN.md location
- **Audit procedures** require AGENTS.md Audit Output section to specify report location pattern

If AGENTS.md doesn't exist, run `rooda agents-sync` first to bootstrap it.

### Direct-Action Procedures

These procedures modify the repository directly (not planning).

#### agents-sync

**Description:** Create or update AGENTS.md by analyzing the repository.

**OODA Composition:**
- Observe: `prompts/observe_bootstrap.md` — Study repository structure, documentation, implementation patterns
- Orient: `prompts/orient_bootstrap.md` — Identify build system, test system, spec/impl paths, work tracking
- Decide: `prompts/decide_bootstrap.md` — Determine what to write or update in AGENTS.md
- Act: `prompts/act_bootstrap.md` — Write or update AGENTS.md, commit changes

**Default Max Iterations:** 3

**Use Case:** First-run bootstrap or periodic reconciliation of AGENTS.md with actual repository state.

**Rationale:** 3 iterations allows for initial creation, verification, and one round of corrections if detections were wrong.

#### build

**Description:** Implement tasks from work tracking (only procedure that modifies specs and implementation).

**OODA Composition:**
- Observe: `prompts/observe_plan_specs_impl.md` — Study AGENTS.md, work tracking, specs, implementation
- Orient: `prompts/orient_build.md` — Understand task requirements, search codebase, identify what to build
- Decide: `prompts/decide_build.md` — Pick task, determine implementation approach, identify files to modify
- Act: `prompts/act_build.md` — Implement, run tests, update work tracking, commit

**Default Max Iterations:** 5

**Use Case:** Autonomous implementation of tasks from the work tracking system. This is the only procedure that modifies specs and implementation files based on work tracking tasks. Other procedures modify specific files (AGENTS.md via agents-sync, PLAN.md via planning procedures) but don't implement features or fix bugs in the project codebase.

**Rationale:** 5 iterations balances progress (implement multiple tasks) with context freshness (don't let AI drift too far).

#### publish-plan

**Description:** Publish converged draft plan to work tracking system.

**OODA Composition:**
- Observe: `prompts/observe_draft_plan.md` — Study PLAN.md and work tracking system
- Orient: `prompts/orient_publish.md` — Parse PLAN.md into task list
- Decide: `prompts/decide_publish.md` — Determine how to map tasks to work tracking system
- Act: `prompts/act_publish.md` — Run work tracking commands to create issues, update PLAN.md status

**Default Max Iterations:** 1

**Use Case:** Import a converged draft plan into the work tracking system.

**Rationale:** 1 iteration is sufficient — this is a straightforward import operation with no iteration needed.

### Audit Procedures

These procedures produce reports without modifying anything.

#### audit-spec

**Description:** Quality assessment of specs.

**OODA Composition:**
- Observe: `prompts/observe_specs.md` — Study all spec files per AGENTS.md definition
- Orient: `prompts/orient_quality.md` — Evaluate against quality criteria (JTBD sections, acceptance criteria, examples)
- Decide: `prompts/decide_refactor_plan.md` — Identify what needs attention
- Act: `prompts/act_plan.md` — Write audit report

**Default Max Iterations:** 1

**Use Case:** Assess spec quality before planning or implementation work.

**Rationale:** 1 iteration is sufficient — audits are read-only assessments.

#### audit-impl

**Description:** Quality assessment of implementation.

**OODA Composition:**
- Observe: `prompts/observe_impl.md` — Study implementation files per AGENTS.md definition
- Orient: `prompts/orient_quality.md` — Evaluate against quality criteria (tests pass, lints pass, build succeeds)
- Decide: `prompts/decide_refactor_plan.md` — Identify what needs attention
- Act: `prompts/act_plan.md` — Write audit report

**Default Max Iterations:** 1

**Use Case:** Assess implementation quality before planning or deployment.

**Rationale:** 1 iteration is sufficient — audits are read-only assessments.

#### audit-agents

**Description:** Accuracy assessment of AGENTS.md against actual repo state.

**OODA Composition:**
- Observe: `prompts/observe_bootstrap.md` — Study repository structure, documentation, implementation
- Orient: `prompts/orient_bootstrap.md` — Detect actual build system, test system, spec/impl paths, work tracking
- Decide: `prompts/decide_bootstrap.md` — Compare AGENTS.md claims vs actual state, identify drift
- Act: `prompts/act_plan.md` — Write audit report listing drift

**Default Max Iterations:** 1

**Use Case:** Verify AGENTS.md accuracy without modifying it (dry-run for agents-sync).

**Rationale:** 1 iteration is sufficient — audits are read-only assessments.

#### audit-spec-to-impl

**Description:** Gap analysis: what's in specs but not in code (execute SDD).

**OODA Composition:**
- Observe: `prompts/observe_plan_specs_impl.md` — Study specs and implementation per AGENTS.md
- Orient: `prompts/orient_gap.md` — Identify features specified but not implemented
- Decide: `prompts/decide_gap_plan.md` — Prioritize gaps by impact
- Act: `prompts/act_plan.md` — Write gap analysis report

**Default Max Iterations:** 1

**Use Case:** Identify what needs to be built to achieve spec compliance.

**Rationale:** 1 iteration is sufficient — audits are read-only assessments.

#### audit-impl-to-spec

**Description:** Gap analysis: what's in code but not in specs (align brownfield with SDD).

**OODA Composition:**
- Observe: `prompts/observe_plan_specs_impl.md` — Study implementation and specs per AGENTS.md
- Orient: `prompts/orient_gap.md` — Identify implemented features not documented in specs
- Decide: `prompts/decide_gap_plan.md` — Prioritize gaps by impact
- Act: `prompts/act_plan.md` — Write gap analysis report

**Default Max Iterations:** 1

**Use Case:** Identify what needs to be documented to achieve spec compliance (brownfield projects).

**Rationale:** 1 iteration is sufficient — audits are read-only assessments.

### Planning Procedures (Spec-Targeted)

These procedures produce draft plans for spec modifications.

#### draft-plan-spec-feat

**Description:** Plan new capability incorporation into specs.

**OODA Composition:**
- Observe: `prompts/observe_story_task_specs_impl.md` — Study TASK.md, specs, implementation
- Orient: `prompts/orient_story_task_incorporation.md` — Understand feature requirements, identify affected specs
- Decide: `prompts/decide_story_task_plan.md` — Break down into spec modification tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to incorporate a new feature into specifications.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-spec-fix

**Description:** Plan spec adjustment to drive a correction.

**OODA Composition:**
- Observe: `prompts/observe_bug_task_specs_impl.md` — Study TASK.md (bug description), specs, implementation
- Orient: `prompts/orient_bug_task_incorporation.md` — Understand bug root cause, identify spec deficiencies
- Decide: `prompts/decide_bug_task_plan.md` — Break down into spec correction tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to fix specs to prevent or document a bug.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-spec-refactor

**Description:** Plan spec restructuring from task file.

**OODA Composition:**
- Observe: `prompts/observe_draft_plan.md` — Study TASK.md (refactor description), specs
- Orient: `prompts/orient_quality.md` — Identify structural issues, duplication, clarity problems
- Decide: `prompts/decide_refactor_plan.md` — Break down into spec restructuring tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to restructure specs for clarity or maintainability.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-spec-chore

**Description:** Plan spec maintenance from task file.

**OODA Composition:**
- Observe: `prompts/observe_draft_plan.md` — Study TASK.md (chore description), specs
- Orient: `prompts/orient_quality.md` — Identify maintenance needs (typos, outdated examples, broken links)
- Decide: `prompts/decide_refactor_plan.md` — Break down into spec maintenance tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 3

**Use Case:** Plan spec maintenance work (typos, formatting, link updates).

**Rationale:** 3 iterations is sufficient for straightforward maintenance planning.

### Planning Procedures (Impl-Targeted)

These procedures produce draft plans for implementation modifications.

#### draft-plan-impl-feat

**Description:** Plan new capability implementation.

**OODA Composition:**
- Observe: `prompts/observe_story_task_specs_impl.md` — Study TASK.md, specs, implementation
- Orient: `prompts/orient_story_task_incorporation.md` — Understand feature requirements, identify affected code
- Decide: `prompts/decide_story_task_plan.md` — Break down into implementation tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to implement a new feature.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-impl-fix

**Description:** Plan implementation correction.

**OODA Composition:**
- Observe: `prompts/observe_bug_task_specs_impl.md` — Study TASK.md (bug description), specs, implementation
- Orient: `prompts/orient_bug_task_incorporation.md` — Understand bug root cause, identify affected code
- Decide: `prompts/decide_bug_task_plan.md` — Break down into bug fix tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to fix a bug in the implementation.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-impl-refactor

**Description:** Plan implementation restructuring from task file.

**OODA Composition:**
- Observe: `prompts/observe_plan_specs_impl.md` — Study TASK.md (refactor description), implementation
- Orient: `prompts/orient_quality.md` — Identify code smells, duplication, complexity issues
- Decide: `prompts/decide_refactor_plan.md` — Break down into refactoring tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 5

**Use Case:** Plan how to refactor implementation for maintainability.

**Rationale:** 5 iterations allows for iterative refinement of the plan as understanding deepens.

#### draft-plan-impl-chore

**Description:** Plan implementation maintenance from task file.

**OODA Composition:**
- Observe: `prompts/observe_plan_specs_impl.md` — Study TASK.md (chore description), implementation
- Orient: `prompts/orient_quality.md` — Identify maintenance needs (dependency updates, lint fixes, comment updates)
- Decide: `prompts/decide_refactor_plan.md` — Break down into maintenance tasks
- Act: `prompts/act_plan.md` — Write PLAN.md with prioritized task list

**Default Max Iterations:** 3

**Use Case:** Plan implementation maintenance work (dependency updates, lint fixes).

**Rationale:** 3 iterations is sufficient for straightforward maintenance planning.

## Edge Cases

### Unknown Procedure Name

```bash
$ rooda unknown-proc
Error: Unknown procedure 'unknown-proc'. Run 'rooda --list-procedures' to see available procedures.
```

### Missing Prompt File

Built-in procedure references non-existent embedded prompt file:

```
FATAL: Procedure 'build' references missing prompt file: prompts/observe_build.md
This is a framework bug. Please report it.
```

Custom procedure references non-existent filesystem prompt file:

```
ERROR: Procedure 'custom-proc' references missing prompt file: custom-prompts/observe.md
Check your rooda-config.yml and ensure the file exists.
```

### Invalid Iteration Limit

Procedure defined with invalid iteration limit:

```yaml
procedures:
  custom-proc:
    default_max_iterations: 0  # Invalid
```

Error:
```
ERROR: Procedure 'custom-proc' has invalid default_max_iterations: 0 (must be >= 1)
```

### Procedure Override in Workspace Config

Built-in procedure overridden in workspace config:

```yaml
# ./rooda-config.yml
procedures:
  build:
    default_max_iterations: 10  # Override built-in default of 5
```

Outcome: `rooda build` uses 10 iterations instead of 5.

### Partial Procedure Override

Built-in procedure partially overridden in workspace config:

```yaml
# ./rooda-config.yml
procedures:
  build:
    default_max_iterations: 10  # Override this field
    # observe, orient, decide, act inherit from built-in
```

Outcome: Only `default_max_iterations` is overridden; OODA phase files remain as built-in defaults.

## Dependencies

- **configuration.md** — Defines Procedure struct and configuration merging
- **prompt-composition.md** — Defines how OODA phase files are assembled
- **iteration-loop.md** — Defines iteration execution and termination
- **operational-knowledge.md** — All procedures read AGENTS.md first

## Implementation Mapping

**Source files:**
- `internal/config/defaults.go` — Built-in procedure definitions
- `internal/config/procedures.go` — Procedure validation and lookup
- `prompts/*.md` — Embedded prompt files (25 files)

**Embedded prompt files:**
- `prompts/observe_bootstrap.md`
- `prompts/observe_draft_plan.md`
- `prompts/observe_impl.md`
- `prompts/observe_plan_specs_impl.md`
- `prompts/observe_specs.md`
- `prompts/observe_story_task_specs_impl.md`
- `prompts/observe_bug_task_specs_impl.md`
- `prompts/orient_bootstrap.md`
- `prompts/orient_build.md`
- `prompts/orient_gap.md`
- `prompts/orient_publish.md`
- `prompts/orient_quality.md`
- `prompts/orient_story_task_incorporation.md`
- `prompts/orient_bug_task_incorporation.md`
- `prompts/decide_bootstrap.md`
- `prompts/decide_build.md`
- `prompts/decide_gap_plan.md`
- `prompts/decide_publish.md`
- `prompts/decide_refactor_plan.md`
- `prompts/decide_story_task_plan.md`
- `prompts/decide_bug_task_plan.md`
- `prompts/act_bootstrap.md`
- `prompts/act_build.md`
- `prompts/act_plan.md`
- `prompts/act_publish.md`

**Related specs:**
- `configuration.md` — Procedure definitions and merging
- `prompt-composition.md` — OODA phase assembly
- `iteration-loop.md` — Iteration execution
- `operational-knowledge.md` — AGENTS.md lifecycle

## Examples

### List All Procedures

```bash
$ rooda --list-procedures

Direct-Action Procedures:
  agents-sync    Create or update AGENTS.md by analyzing the repository
  build          Implement tasks from work tracking (only procedure that modifies code)
  publish-plan   Publish converged draft plan to work tracking system

Audit Procedures:
  audit-spec           Quality assessment of specs
  audit-impl           Quality assessment of implementation
  audit-agents         Accuracy assessment of AGENTS.md against actual repo state
  audit-spec-to-impl   Gap analysis: what's in specs but not in code
  audit-impl-to-spec   Gap analysis: what's in code but not in specs

Planning Procedures (Spec-Targeted):
  draft-plan-spec-feat      Plan new capability incorporation into specs
  draft-plan-spec-fix       Plan spec adjustment to drive a correction
  draft-plan-spec-refactor  Plan spec restructuring from task file
  draft-plan-spec-chore     Plan spec maintenance from task file

Planning Procedures (Impl-Targeted):
  draft-plan-impl-feat      Plan new capability implementation
  draft-plan-impl-fix       Plan implementation correction
  draft-plan-impl-refactor  Plan implementation restructuring from task file
  draft-plan-impl-chore     Plan implementation maintenance from task file
```

### Invoke Procedure

```bash
$ rooda build
INFO: Executing procedure: build
INFO: OODA composition:
  Observe: prompts/observe_plan_specs_impl.md
  Orient: prompts/orient_build.md
  Decide: prompts/decide_build.md
  Act: prompts/act_build.md
INFO: Default max iterations: 5
INFO: Starting iteration loop...
```

### Override Iteration Limit

```bash
$ rooda build --max-iterations 10
INFO: Executing procedure: build
INFO: Max iterations overridden: 10 (default: 5)
INFO: Starting iteration loop...
```

### Procedure Help

```bash
$ rooda build --help

Procedure: build
Description: Implement tasks from work tracking (only procedure that modifies code)
Category: direct-action

OODA Composition:
  Observe: prompts/observe_plan_specs_impl.md
  Orient: prompts/orient_build.md
  Decide: prompts/decide_build.md
  Act: prompts/act_build.md

Default Max Iterations: 5

Use Case:
  Autonomous implementation of tasks from the work tracking system.

Usage:
  rooda build [flags]

Flags:
  --max-iterations <n>   Override default max iterations
  --unlimited            Run until convergence (no iteration limit)
  --dry-run              Display assembled prompt without executing
  --context <text>       Inject user-provided context
  --verbose              Stream AI CLI output
  (see 'rooda --help' for all flags)
```

### Custom Procedure in Workspace Config

```yaml
# ./rooda-config.yml
procedures:
  custom-audit:
    description: "Custom audit procedure for security checks"
    observe: "custom-prompts/observe-security.md"
    orient: "custom-prompts/orient-security.md"
    decide: "custom-prompts/decide-security.md"
    act: "custom-prompts/act-security.md"
    default_max_iterations: 1
    category: "audit"
```

```bash
$ rooda custom-audit
INFO: Executing procedure: custom-audit
INFO: OODA composition:
  Observe: custom-prompts/observe-security.md
  Orient: custom-prompts/orient-security.md
  Decide: custom-prompts/decide-security.md
  Act: custom-prompts/act-security.md
INFO: Default max iterations: 1
INFO: Starting iteration loop...
```

## Notes

### Design Rationale

**Why 16 procedures?**
Covers the core use cases: direct actions (3), audits (5), and planning for both specs and impl across 4 conventional commit types (8). This is the minimal complete set.

**Why separate spec-targeted and impl-targeted planning procedures?**
Specs and implementation have different concerns. Spec planning focuses on requirements and acceptance criteria; impl planning focuses on code structure and integration points. Separate procedures keep prompts focused.

**Why 4 conventional commit types (feat, fix, refactor, chore)?**
These are the most common types in conventional commits. They map cleanly to different planning concerns: feat = new capability, fix = correction, refactor = restructuring, chore = maintenance.

**Why iteration limits vary by procedure?**
Direct-action procedures need more iterations (build = 5) because they do real work. Audits need only 1 iteration (read-only assessment). Planning procedures need moderate iterations (3-5) for iterative refinement.

**Why embed prompt files in the binary?**
Zero-config startup — users shouldn't need to copy prompt files to use built-in procedures. Embedding makes the binary self-contained.

**Why allow custom procedures to reference filesystem prompt files?**
Extensibility — teams can define custom procedures without modifying the framework. Filesystem paths provide flexibility.

**Why validate procedures at startup?**
Fail fast — if a procedure is misconfigured, the user should know immediately, not after invoking it.

**Why categorize procedures?**
Helps users understand procedure purpose and discover related procedures. Categories map to different use cases.

**Why kebab-case for procedure names?**
Consistent with CLI conventions (e.g., `git commit`, `docker build`). Easier to type than snake_case or camelCase.

### Prompt File Reuse

Many procedures share prompt files:

**Shared observe files:**
- `observe_bootstrap.md` — Used by agents-sync, audit-agents
- `observe_draft_plan.md` — Used by publish-plan, draft-plan-spec-refactor, draft-plan-spec-chore
- `observe_plan_specs_impl.md` — Used by build, audit-spec-to-impl, audit-impl-to-spec, draft-plan-impl-refactor, draft-plan-impl-chore
- `observe_specs.md` — Used by audit-spec
- `observe_impl.md` — Used by audit-impl
- `observe_story_task_specs_impl.md` — Used by draft-plan-spec-feat, draft-plan-impl-feat
- `observe_bug_task_specs_impl.md` — Used by draft-plan-spec-fix, draft-plan-impl-fix

**Shared orient files:**
- `orient_bootstrap.md` — Used by agents-sync, audit-agents
- `orient_quality.md` — Used by audit-spec, audit-impl, draft-plan-spec-refactor, draft-plan-spec-chore, draft-plan-impl-refactor, draft-plan-impl-chore
- `orient_gap.md` — Used by audit-spec-to-impl, audit-impl-to-spec
- `orient_story_task_incorporation.md` — Used by draft-plan-spec-feat, draft-plan-impl-feat
- `orient_bug_task_incorporation.md` — Used by draft-plan-spec-fix, draft-plan-impl-fix

**Shared decide files:**
- `decide_bootstrap.md` — Used by agents-sync, audit-agents (audit-agents uses it to compare, not to decide changes)
- `decide_gap_plan.md` — Used by audit-spec-to-impl, audit-impl-to-spec
- `decide_refactor_plan.md` — Used by audit-spec, audit-impl, draft-plan-spec-refactor, draft-plan-spec-chore, draft-plan-impl-refactor, draft-plan-impl-chore
- `decide_story_task_plan.md` — Used by draft-plan-spec-feat, draft-plan-impl-feat
- `decide_bug_task_plan.md` — Used by draft-plan-spec-fix, draft-plan-impl-fix

**Shared act files:**
- `act_plan.md` — Used by all audit and planning procedures (write report or plan)
- `act_bootstrap.md` — Used by agents-sync (write/update AGENTS.md)
- `act_build.md` — Used by build (implement, test, commit)
- `act_publish.md` — Used by publish-plan (import to work tracking)

This reuse reduces the total number of prompt files from 64 (16 procedures × 4 phases) to 25 unique files.

### Iteration Limit Rationale

**Direct-action:**
- agents-sync: 3 iterations (create, verify, correct)
- build: 5 iterations (implement multiple tasks, test, fix)
- publish-plan: 1 iteration (straightforward import)

**Audit:**
- All audits: 1 iteration (read-only assessment, no iteration needed)

**Planning:**
- feat/fix procedures: 5 iterations (complex planning, iterative refinement)
- refactor procedures: 5 iterations (complex planning, iterative refinement)
- chore procedures: 3 iterations (simpler planning, less refinement needed)

These are defaults — users can override with `--max-iterations` or `--unlimited`.

## Iteration Limit Tuning

Default limits are conservative starting points. Adjust based on observed behavior:

**Signs limit is too low:**
- Procedure frequently hits max iterations without `<promise>SUCCESS</promise>` signal
- Work is incomplete when loop terminates
- **Action:** Increase by 2-3 iterations, or use `--unlimited` to find natural convergence point

**Signs limit is too high:**
- Procedure consistently succeeds in fewer iterations than the limit
- Wasted API calls on unnecessary iterations
- **Action:** Reduce to observed average + 1 buffer iteration

**Signs of context degradation:**
- AI output quality degrades in later iterations
- Agent makes mistakes it didn't make earlier
- **Action:** Reduce limit to keep AI in "smart zone" (40-60% context utilization)

**How to tune:**
1. Run procedure with `--verbose` to observe iteration behavior
2. Note which iteration produces `<promise>SUCCESS</promise>`
3. Set limit to that number + 1-2 buffer iterations
4. Override per-procedure in `rooda-config.yml` or per-invocation with `--max-iterations`

**Example tuning:**
```yaml
# rooda-config.yml
procedures:
  build:
    default_max_iterations: 8  # Increased from 5 after observing consistent 6-7 iteration completions
  draft-plan-impl-feat:
    default_max_iterations: 3  # Reduced from 5 after observing consistent 2-iteration completions
```
