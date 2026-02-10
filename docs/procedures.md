# Procedures

rooda includes 21 built-in procedures organized by job type. Each procedure is an OODA loop workflow composed from reusable prompt fragments.

## Bootstrap & Sync

### bootstrap

**Job**: Analyze repository and create/update AGENTS.md operational guide.

**When to use**: First time setting up rooda in a repository, or when project structure changes significantly.

**Example**:
```bash
rooda bootstrap --ai-cmd-alias kiro-cli
```

**Output**: Creates or updates `AGENTS.md` with detected build/test/lint commands, work tracking system, spec/impl patterns.

### agents-sync

**Job**: Detect drift between AGENTS.md and actual repository state, then update AGENTS.md.

**When to use**: When commands in AGENTS.md stop working or paths change.

**Example**:
```bash
rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1
```

**Output**: Updates AGENTS.md to fix drifted commands and paths.

## Audit Procedures

### audit-spec

**Job**: Review spec files against quality criteria and generate audit report.

**When to use**: Before planning work, to identify spec quality issues.

**Example**:
```bash
rooda audit-spec --ai-cmd-alias kiro-cli
```

**Output**: Audit report identifying missing sections, broken cross-references, incomplete acceptance criteria.

### audit-impl

**Job**: Review implementation files, run tests and lints, generate audit report.

**When to use**: Before planning work, to identify code quality issues.

**Example**:
```bash
rooda audit-impl --ai-cmd-alias kiro-cli
```

**Output**: Audit report with test results, lint errors, code quality issues.

### audit-agents

**Job**: Verify AGENTS.md matches repository state and commands work correctly.

**When to use**: When suspecting AGENTS.md is out of sync.

**Example**:
```bash
rooda audit-agents --ai-cmd-alias kiro-cli
```

**Output**: Audit report listing drifted commands and paths.

### audit-spec-to-impl

**Job**: Identify features specified but not yet implemented.

**When to use**: Gap analysis - what's in specs but missing from code.

**Example**:
```bash
rooda audit-spec-to-impl --ai-cmd-alias kiro-cli
```

**Output**: Gap analysis report listing unimplemented features.

### audit-impl-to-spec

**Job**: Identify code that exists but is not documented in specifications.

**When to use**: Reverse gap analysis - what's in code but missing from specs.

**Example**:
```bash
rooda audit-impl-to-spec --ai-cmd-alias kiro-cli
```

**Output**: Gap analysis report listing undocumented features.

## Planning Procedures

All planning procedures write to `PLAN.md` at project root. Use `publish-plan` to import the plan into your work tracking system.

### draft-plan-spec-feat

**Job**: Analyze feature requirements and create implementation plan focused on specifications.

**When to use**: Planning a new feature that requires spec changes.

**Example**:
```bash
rooda draft-plan-spec-feat --ai-cmd-alias kiro-cli --context "Add user authentication"
```

**Output**: `PLAN.md` with prioritized tasks for spec changes.

### draft-plan-impl-feat

**Job**: Analyze feature requirements and create implementation plan focused on code.

**When to use**: Planning a new feature that requires code changes.

**Example**:
```bash
rooda draft-plan-impl-feat --ai-cmd-alias kiro-cli --context "Add user authentication"
```

**Output**: `PLAN.md` with prioritized tasks for code changes.

### draft-plan-spec-fix

**Job**: Analyze bug root cause and create fix plan focused on specifications.

**When to use**: Bug requires spec clarification or correction.

**Example**:
```bash
rooda draft-plan-spec-fix --ai-cmd-alias kiro-cli --context "Authentication fails for OAuth users"
```

**Output**: `PLAN.md` with tasks to fix spec issues.

### draft-plan-impl-fix

**Job**: Analyze bug root cause and create fix plan focused on code.

**When to use**: Bug requires code fix.

**Example**:
```bash
rooda draft-plan-impl-fix --ai-cmd-alias kiro-cli --context "Authentication fails for OAuth users"
```

**Output**: `PLAN.md` with tasks to fix code issues.

### draft-plan-spec-refactor

**Job**: Apply quality criteria to specifications, write refactoring tasks to draft plan.

**When to use**: Improving spec quality (structure, clarity, completeness).

**Example**:
```bash
rooda draft-plan-spec-refactor --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with spec refactoring tasks.

### draft-plan-impl-refactor

**Job**: Apply quality criteria to implementation, write refactoring tasks to draft plan.

**When to use**: Improving code quality (structure, tests, documentation).

**Example**:
```bash
rooda draft-plan-impl-refactor --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with code refactoring tasks.

### draft-plan-spec-chore

**Job**: Identify maintenance needs in specs and create chore plan.

**When to use**: Routine spec maintenance (update examples, fix formatting, etc.).

**Example**:
```bash
rooda draft-plan-spec-chore --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with spec maintenance tasks.

### draft-plan-impl-chore

**Job**: Identify maintenance needs in code and create chore plan.

**When to use**: Routine code maintenance (update dependencies, fix warnings, etc.).

**Example**:
```bash
rooda draft-plan-impl-chore --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with code maintenance tasks.

### draft-plan-spec-to-impl

**Job**: Compare specifications to implementation, identify missing features, write implementation plan.

**When to use**: Implementing features that are already specified.

**Example**:
```bash
rooda draft-plan-spec-to-impl --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with tasks to implement specified features.

### draft-plan-impl-to-spec

**Job**: Compare implementation to specifications, identify undocumented features, write documentation plan.

**When to use**: Documenting features that are already implemented.

**Example**:
```bash
rooda draft-plan-impl-to-spec --ai-cmd-alias kiro-cli
```

**Output**: `PLAN.md` with tasks to document implemented features.

### draft-plan-story-to-spec

**Job**: Analyze new feature requirements, determine how to expand specs to accommodate the story.

**When to use**: Converting user stories into spec changes.

**Example**:
```bash
# Create TASK.md with story details first
rooda draft-plan-story-to-spec --ai-cmd-alias kiro-cli
```

**Input**: Reads `TASK.md` at project root.

**Output**: `PLAN.md` with tasks to update specs for the story.

### draft-plan-bug-to-spec

**Job**: Analyze bug report, determine spec adjustments (acceptance criteria, edge cases, etc.).

**When to use**: Bug reveals spec gaps or ambiguities.

**Example**:
```bash
# Create TASK.md with bug details first
rooda draft-plan-bug-to-spec --ai-cmd-alias kiro-cli
```

**Input**: Reads `TASK.md` at project root.

**Output**: `PLAN.md` with tasks to update specs for the bug.

### publish-plan

**Job**: Read converged draft plan and create issues in work tracking system per AGENTS.md.

**When to use**: After a planning procedure converges, to import tasks into work tracking.

**Example**:
```bash
rooda publish-plan --ai-cmd-alias kiro-cli
```

**Input**: Reads `PLAN.md` at project root.

**Output**: Creates issues in work tracking system (beads, GitHub Issues, etc.).

## Build Procedure

### build

**Job**: Read from work tracking system, pick most important task, implement it using TDD.

**When to use**: Autonomous implementation of planned work.

**Example**:
```bash
rooda build --ai-cmd-alias kiro-cli --max-iterations 5
```

**Input**: Reads ready work from work tracking system per AGENTS.md.

**Output**: Implements task, runs tests, commits changes, updates work tracking.

**Note**: This is the only procedure that modifies specs and implementation files based on work tracking tasks.

## Iteration Control

All procedures support iteration control flags:

```bash
# Limit iterations
rooda build --ai-cmd-alias kiro-cli --max-iterations 3

# Unlimited iterations (runs until success or failure threshold)
rooda build --ai-cmd-alias kiro-cli --unlimited

# Dry run (validate without executing)
rooda build --ai-cmd-alias kiro-cli --dry-run
```

See [CLI Reference](cli-reference.md) for all flags.
