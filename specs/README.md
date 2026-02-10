# Specifications

This directory contains specifications for **rooda v2** — a Go rewrite of the ralph-wiggum-ooda framework. Each specification follows the Jobs-to-be-Done (JTBD) structure: one spec per topic of concern, organized by the job it serves.

## Customer Types

### Job Executor
The developer running rooda to orchestrate AI coding agents. They invoke procedures, monitor iteration progress, and steer the loop when it drifts. They care about: reliable execution, clear feedback, fast iteration cycles, and minimal setup friction.

### Product Lifecycle Support Team
People who install, configure, upgrade, and maintain rooda across projects and teams. This includes DevOps engineers setting up CI/CD pipelines, team leads standardizing tooling, and contributors to the framework itself. They care about: easy distribution, cross-platform support, testability, and clean extensibility.

### Buyer
Engineering managers and tech leads deciding whether to adopt rooda for their teams. They care about: clear value proposition, low adoption risk, observable outcomes, and alignment with existing workflows.

## Core Functional Job

**Orchestrate AI coding agents through structured OODA iteration loops to autonomously build, plan, and maintain software from specifications.**

The developer wants to define what should be built (specs), point an AI agent at the work, and have it iterate toward a solution — with fresh context each cycle and file-based state providing continuity. The loop should run unattended, self-correct through empirical feedback, and produce working, tested software. Quality gates (tests, lints) and git operations are driven by the prompts, not the loop orchestrator.

Prompts are composed from reusable fragments — small, focused markdown files organized by OODA phase. Each procedure defines which fragments to use for each phase, enabling high reusability across the 16 built-in procedures.

## Related Jobs

- **Sync the agent-project interface** — analyze a repository and create or update the operational guide (AGENTS.md) so agents can interact with the project effectively. Works whether AGENTS.md exists or not — creates from scratch on first run, reconciles with actual repo state on subsequent runs. This is a direct-action procedure (not planning) because AGENTS.md must exist before any planning procedure can run (chicken-and-egg: plans need AGENTS.md to know where plans live).
- **Audit what exists** — assess specs, implementation, AGENTS.md, or the gap between specs and implementation, producing an audit report. Audits don't modify anything — they identify what needs attention. Audit output feeds as context into planning procedures.
- **Plan work by type and target** — given a work item classified by conventional commit type (feat, fix, refactor, chore) and a target (spec or impl), produce a prioritized task list that agents can execute. Procedures are cheap because OODA components are reused across types — a `draft-plan-spec-feat` and `draft-plan-spec-fix` may share most of their observe/act components, differing only in orient/decide.
- **Publish plans to work tracking** — import a converged draft plan into the project's work tracking system (beads, GitHub Issues, file-based, etc.).
- **Build from plan** — implement tasks from the work tracking system. The only procedure that modifies specs and implementation files based on work tracking tasks (other procedures modify specific files: AGENTS.md, PLAN.md).
- **Provide context to guide a procedure** — pass runtime hints to any procedure execution (e.g., "focus on the auth module", "the new feature should integrate with the payment service") that steer the agent's focus without changing procedure definitions or prompt files.
- **Configure procedures for a team** — define custom OODA procedures using fragment arrays, AI command aliases, and project-specific settings without modifying framework code. Procedures compose prompts from reusable fragments, with support for inline content and template parameters.
- **Distribute and install the tool** — get rooda running on a new machine or in a CI/CD pipeline with minimal friction and no external dependencies.

## Emotional Jobs

- **Feel confident the loop is working** — clear progress indicators, iteration counts, and error messages so the developer isn't anxious about what's happening.
- **Feel in control** — ability to stop, dry-run, and override at every level. The loop is a tool, not a black box.
- **Feel safe** — sandboxing guidance, blast radius awareness, and no silent failures that corrupt the codebase.

## Desired Outcomes (Success Metrics)

When executing the core functional job, the developer measures success by:

- Minimize the time to go from specs to working implementation
- Minimize the number of iterations that produce no useful progress
- Minimize the setup effort to start using rooda on a new project
- Minimize the time to diagnose why an iteration failed
- Minimize the risk of AI-generated code breaking existing functionality
- Maximize the percentage of iterations where the AI stays in its "smart zone" (40-60% context utilization)
- Maximize the reusability of prompt fragments across different procedures
- Maximize the observability of what the loop is doing at any moment

## What Changed from v1 (Archive)

The previous iteration (archived in `archive/`) was a bash script (`rooda.sh`) that shelled out to `yq` for YAML parsing and piped prompts to an AI CLI. It validated the core concept — composable OODA prompts, fresh context per iteration, file-based state — but had significant operational limitations:

| v1 Limitation | v2 Approach |
|---|---|
| Bash script — fragile, hard to test, platform-specific | Go binary — testable, cross-platform, single artifact |
| Required `yq` external dependency | Built-in YAML/config parsing |
| No error handling for AI CLI failures | Structured error handling with configurable failure threshold and `<promise>` output signals |
| No dry-run mode | `--dry-run` flag shows assembled prompt without executing |
| No observability during iteration | Loop-level progress display with timing per iteration; `--verbose` streams AI CLI output |
| Duplicate validation code | Clean architecture with shared validation |
| Config validation at runtime only | Upfront config validation with clear error messages |
| Manual file copying for installation | Single binary distribution, embedded default prompts |

## Jobs to Be Done

### J1: Execute OODA Iterations
Run AI coding agents through controlled OODA iteration cycles with fresh context per run and file-based state continuity. This is the core loop — everything else feeds into or out of it.

### J2: Compose and Assemble Prompts
Combine fragment arrays for each OODA phase (observe, orient, decide, act) and optional user-provided context into a single executable prompt. Each phase uses an array of fragments that are concatenated together, supporting embedded defaults, user-provided custom fragments, inline content, and Go template parameterization.

### J3: Integrate with AI CLI Tools
Pipe assembled prompts to a configurable AI CLI tool with support for command aliases, environment variables, and direct command override. Built-in support for kiro-cli, claude, github copilot, and cursor agent, with extensibility for custom tools.

### J4: Configure Procedures and Settings
Define custom OODA procedures, AI command aliases, and project-specific settings through a three-tier configuration system — workspace (`./`), global (`<config_dir>/`), and environment variables — with sensible built-in defaults for zero-config startup. Tiers merge with clear precedence (CLI flags > env vars > workspace > global > built-in defaults) and provenance tracking so users know where each setting comes from.

### J5: Provide a Command-Line Interface
Expose all framework capabilities through a CLI that supports named procedures, explicit OODA phase flags, global options, and helpful error messages.

### J6: Define the Agent-Project Interface
Specify the AGENTS.md format — required sections, field definitions, and structural conventions — that serves as the contract between AI agents and the repository. Covers work tracking, build commands, spec/impl definitions, and quality criteria. This is the schema; J10 covers the runtime lifecycle.

### J7: Distribute and Install
Enable users to install rooda as a single binary with no external dependencies, supporting macOS, Linux, and CI/CD environments.

### J8: Handle Errors and Build Resilience
Detect, report, and recover from failures — AI CLI crashes, network issues, test failures, invalid configs — with configurable retry logic, timeouts, and graceful degradation.

### J9: Observe and Control Loop Execution
Provide visibility into what the loop is doing (timing, iteration progress, outcome) and controls to stop, dry-run, and override behavior.

### J10: Maintain Project Operational Knowledge
Every procedure reads AGENTS.md first as the source of truth for project-specific behavior — build commands, file paths, work tracking, quality criteria. Agents defer to it, verify it empirically (run commands, check paths), and update it in-place when something is wrong or a new learning occurs. This is the read-verify-update lifecycle that keeps AGENTS.md accurate across iterations.

## Topics of Concern

### Execution Engine
| Topic | Job | Description |
|---|---|---|
| [iteration-loop](iteration-loop.md) | J1 | Execute OODA cycles with fresh context, termination control, and state continuity |
| [prompt-composition](prompt-composition.md) | J2, J5 | Assemble fragment arrays for each OODA phase and optional user-provided context into a single prompt, with embedded defaults and template support |
| [ai-cli-integration](ai-cli-integration.md) | J3 | Pipe prompts to configurable AI CLI tools with alias resolution |
| [error-handling](error-handling.md) | J8 | Retry logic, timeouts, failure detection, and graceful degradation |

### Configuration & Interface
| Topic | Job | Description |
|---|---|---|
| [cli-interface](cli-interface.md) | J5 | Command-line argument parsing, procedure invocation, help text |
| [configuration](configuration.md) | J4 | YAML config schema, procedure definitions, AI command aliases, defaults |
| [agents-md-format](agents-md-format.md) | J6 | AGENTS.md structure, required sections, field definitions |
| [operational-knowledge](operational-knowledge.md) | J10 | Read-verify-update lifecycle for AGENTS.md across all procedures |

### Distribution & Operations
| Topic | Job | Description |
|---|---|---|
| [distribution](distribution.md) | J7 | Single binary build, embedded prompts, cross-platform support |
| [observability](observability.md) | J9 | Structured logging, progress display, iteration timing, dry-run mode |

### Procedure Library
| Topic | Job | Description |
|---|---|---|
| [procedures](procedures.md) | J1, J2 | Fragment-based procedure system with embedded fragments composing 16 built-in procedures |

## Specification Status

Written specs with extracted JTBDs:

| Spec | Job to be Done |
|------|----------------|
| [agents-md-format](agents-md-format.md) | Specify the AGENTS.md format — required sections, field definitions, and structural conventions — that serves as the contract between AI agents and the repository. This file is the source of truth for project-specific behavior: build commands, file paths, work tracking, quality criteria. Agents read it first, verify it empirically, and update it when drift is detected. |
| [ai-cli-integration](ai-cli-integration.md) | Pipe assembled prompts to a configurable AI CLI tool with support for command aliases, environment variables, and direct command override. Built-in support for kiro-cli, claude, github copilot, and cursor agent, with extensibility for custom tools. |
| [cli-interface](cli-interface.md) | Expose all framework capabilities through a CLI that supports named procedures, explicit OODA phase flags, global options, and helpful error messages. |
| [configuration](configuration.md) | Define custom OODA procedures, AI command aliases, and project-specific settings through a three-tier configuration system — workspace (`./`), global (`<config_dir>/`), and environment variables — with sensible built-in defaults for zero-config startup. Tiers merge with clear precedence (CLI flags > env vars > workspace > global > built-in defaults) and provenance tracking so users know where each setting comes from. |
| [distribution](distribution.md) | Enable users to install rooda as a single binary with no external dependencies, supporting macOS, Linux, and CI/CD environments. |
| [error-handling](error-handling.md) | Detect, report, and recover from failures — AI CLI crashes, network issues, test failures, invalid configs — with configurable retry logic, timeouts, and graceful degradation. The loop must distinguish transient failures (retry) from permanent failures (abort), provide clear diagnostics, and prevent silent corruption. |
| [iteration-loop](iteration-loop.md) | Execute AI coding agents through controlled OODA iteration cycles that clear AI context between runs, preventing LLM degradation while maintaining file-based state continuity. Each iteration invokes the AI CLI as a fresh process — the agent starts clean, processes the assembled prompt, executes tools, then exits. The loop orchestrator persists across iterations, managing termination and state. |
| [observability](observability.md) | Provide visibility into what the loop is doing (timing, iteration progress, outcome) and controls to stop, dry-run, and override behavior. The developer wants to understand loop execution state, diagnose failures, and control output verbosity without modifying configuration files. |
| [operational-knowledge](operational-knowledge.md) | Every procedure reads AGENTS.md first as the source of truth for project-specific behavior — build commands, file paths, work tracking, quality criteria. Agents defer to it, verify it empirically (run commands, check paths), and update it in-place when something is wrong or a new learning occurs. |
| [procedures](procedures.md) | Define the 16 built-in procedures that ship as defaults — their OODA phase compositions, iteration limits, and use cases. |
| [prompt-composition](prompt-composition.md) | Assemble fragment arrays for each OODA phase (observe, orient, decide, act) and optional user-provided context into a single prompt that can be piped to an AI CLI tool. Supports embedded defaults, user-provided custom fragments, inline content, and Go template parameterization with clear path resolution. |

## How to Write Specs

Each spec follows the JTBD template structure:

```markdown
# [Topic Name]

## Job to be Done
[What user outcome does this enable?]

## Activities
[Key steps or operations]

## Acceptance Criteria
- [ ] [Verifiable outcome]

## Data Structures
[Types, schemas, formats]

## Algorithm
[Logic, pseudocode, flow]

## Edge Cases
[Boundary conditions and error scenarios]

## Dependencies
[Prerequisites]

## Implementation Mapping
[Source files and related specs]

## Examples
[Input/output pairs with verification]

## Notes
[Design rationale and decisions]
```

**Principles:**
- One topic per file, named with lowercase hyphens (`iteration-loop.md`)
- Outcome-focused, not mechanism-focused
- Acceptance criteria must be testable
- Include design rationale ("capture the why")
- Reference related specs in Implementation Mapping
