# Draft Plan: Spec to Implementation Gap Analysis

## Priority 1: Core Framework Specifications

### Create rooda-loop-execution.md spec
Document the main loop execution mechanism - how rooda.sh orchestrates OODA iterations.

**What's specified:**
- README.md describes "bash loop: each iteration loads prompt files, executes through AI CLI, updates files, exits"
- AGENTS.md defines build/test/lint commands
- specs/specification-system.md defines JTBD methodology

**What's implemented:**
- src/rooda.sh implements argument parsing, config loading, prompt composition, iteration loop
- src/rooda-config.yml defines 9 procedures with OODA phase mappings

**Gap:** No JTBD-based spec documenting the loop execution mechanism itself.

**Acceptance criteria:**
- Documents JTBD: "Execute autonomous coding iterations with fresh context"
- Activities: parse arguments, load config, compose prompt, execute AI CLI, iterate
- Algorithm: iteration loop with exit-and-restart pattern
- Data structures: config YAML format, prompt composition structure
- Implementation mapping: references src/rooda.sh functions

**Dependencies:** None

---

### Create procedure-system.md spec
Document the 9 procedures and how they compose OODA phases.

**What's specified:**
- README.md lists 9 procedures in table format
- src/README.md describes procedures and component composition
- specs/specification-system.md defines JTBD methodology

**What's implemented:**
- src/rooda-config.yml defines all 9 procedures with phase mappings
- src/components/*.md implement 28 OODA phase components

**Gap:** No JTBD-based spec documenting the procedure system design.

**Acceptance criteria:**
- Documents JTBD: "Compose reusable OODA phases into specialized procedures"
- Activities: define procedure, map phases, set default iterations
- Lists all 9 procedures with purpose and phase composition
- Data structures: procedure config format
- Implementation mapping: references src/rooda-config.yml and src/components/

**Dependencies:** None

---

### Create component-system.md spec
Document the composable OODA component architecture.

**What's specified:**
- src/README.md describes component principles and common steps
- README.md mentions "composable prompt files"
- specs/specification-system.md defines JTBD methodology

**What's implemented:**
- src/components/*.md implement 28 components across 4 OODA phases
- Components use step codes (O1-O15, R1-R22, D1-D15, A1-A9)
- src/README.md documents all common steps

**Gap:** No JTBD-based spec documenting the component system architecture.

**Acceptance criteria:**
- Documents JTBD: "Reuse prompt logic across multiple procedures"
- Activities: define common steps, reference by code, compose into phases
- Lists all common steps with codes
- Edge cases: component reuse patterns, step code consistency
- Implementation mapping: references src/components/ and src/README.md

**Dependencies:** None

---

## Priority 2: Workflow Specifications

### Create bootstrap-procedure.md spec
Document the bootstrap procedure that creates AGENTS.md.

**What's specified:**
- README.md describes bootstrap as "create operational guide"
- src/README.md lists bootstrap components
- specs/agents-md-format.md defines AGENTS.md structure

**What's implemented:**
- src/rooda-config.yml defines bootstrap procedure
- src/components/observe_bootstrap.md, orient_bootstrap.md, decide_bootstrap.md, act_bootstrap.md

**Gap:** No JTBD-based spec for the bootstrap procedure workflow.

**Acceptance criteria:**
- Documents JTBD: "Initialize agent-project interface"
- Activities: analyze repository, determine definitions, create AGENTS.md
- Algorithm: observe repo → orient analysis → decide structure → act create
- Implementation mapping: references bootstrap components

**Dependencies:** component-system.md (references component architecture)

---

### Create build-procedure.md spec
Document the build procedure that implements code from work tracking.

**What's specified:**
- README.md describes build as "implements tasks from work tracking"
- src/README.md lists build components
- AGENTS.md defines work tracking system

**What's implemented:**
- src/rooda-config.yml defines build procedure
- src/components/observe_plan_specs_impl.md, orient_build.md, decide_build.md, act_build.md

**Gap:** No JTBD-based spec for the build procedure workflow.

**Acceptance criteria:**
- Documents JTBD: "Implement code from prioritized work tracking"
- Activities: query work, pick task, implement, test, commit
- Algorithm: observe work → orient understand → decide approach → act implement
- Edge cases: test failures trigger backpressure, parallel subagents (only 1 for build/test)
- Implementation mapping: references build components

**Dependencies:** component-system.md (references component architecture)

---

### Create planning-procedures.md spec
Document the 6 draft planning procedures and publish-plan procedure.

**What's specified:**
- README.md lists 6 draft planning procedures + publish-plan
- src/README.md describes planning components
- AGENTS.md defines planning system (draft location, publishing mechanism)

**What's implemented:**
- src/rooda-config.yml defines 7 planning procedures
- src/components/act_plan.md, act_publish.md, orient_gap.md, orient_quality.md, etc.

**Gap:** No JTBD-based spec for the planning procedure workflows.

**Acceptance criteria:**
- Documents JTBD: "Generate prioritized implementation plans"
- Activities: analyze gaps/quality, structure plan, publish to work tracking
- Lists all 7 planning procedures with purpose
- Algorithm: draft iterations converge plan → publish imports to work tracking
- Implementation mapping: references planning components

**Dependencies:** component-system.md (references component architecture)

---

## Priority 3: Quality and Integration Specifications

### Create backpressure-system.md spec
Document the quality control through backpressure (tests, lints, quality criteria).

**What's specified:**
- README.md mentions "backpressure for quality control"
- AGENTS.md defines quality criteria
- src/README.md mentions "backpressure is mandatory"

**What's implemented:**
- src/components/act_build.md enforces test passing before commit
- src/components/orient_quality.md applies boolean criteria
- Quality assessment triggers refactoring procedures

**Gap:** No JTBD-based spec documenting the backpressure system design.

**Acceptance criteria:**
- Documents JTBD: "Enforce quality through empirical feedback"
- Activities: run tests, apply criteria, trigger refactoring
- Algorithm: downstream backpressure (tests) + upstream backpressure (quality criteria)
- Edge cases: test failures block commit, criteria failures trigger planning
- Implementation mapping: references act_build.md, orient_quality.md

**Dependencies:** None

---

### Create agents-md-lifecycle.md spec
Document how AGENTS.md is created, maintained, and updated across procedures.

**What's specified:**
- specs/agents-md-format.md defines AGENTS.md structure
- README.md mentions "assumed inaccurate until verified empirically"
- src/README.md mentions "capture the why, keep it up to date"

**What's implemented:**
- src/components/act_bootstrap.md creates AGENTS.md
- src/components/act_build.md, act_plan.md update AGENTS.md when learned
- All observe components read AGENTS.md first

**Gap:** No JTBD-based spec documenting AGENTS.md lifecycle management.

**Acceptance criteria:**
- Documents JTBD: "Maintain accurate operational guide through empirical learning"
- Activities: create, verify, update, capture rationale
- Algorithm: bootstrap creates → procedures verify → update when wrong
- Edge cases: command failures trigger updates, new patterns captured
- Implementation mapping: references act_bootstrap.md, act_build.md, act_plan.md

**Dependencies:** bootstrap-procedure.md (references bootstrap creation)

---

## Priority 4: Advanced Topics

### Create context-management.md spec
Document the fresh-context-per-iteration approach and smart zone utilization.

**What's specified:**
- README.md describes "exit-and-restart pattern" and "smart zone (40-60% utilization)"
- README.md explains "LLMs degrade beyond 60% context"

**What's implemented:**
- src/rooda.sh exits completely after each iteration
- Iteration loop clears context by restarting process
- File-based state (AGENTS.md, work tracking, specs, code) persists

**Gap:** No JTBD-based spec documenting context management strategy.

**Acceptance criteria:**
- Documents JTBD: "Prevent LLM degradation through context clearing"
- Activities: execute iteration, exit process, restart fresh
- Algorithm: load → execute → exit → repeat
- Edge cases: file-based state provides continuity, no conversational baggage
- Implementation mapping: references src/rooda.sh iteration loop

**Dependencies:** rooda-loop-execution.md (references loop mechanism)

---

### Create subagent-system.md spec
Document parallel subagent spawning for auxiliary work.

**What's specified:**
- README.md mentions "parallel subagents" and "only one subagent for build/tests"
- src/README.md describes "using parallel subagents"

**What's implemented:**
- src/components/act_build.md spawns subagents for implementation
- src/components/decide_build.md determines parallel work breakdown
- Constraint: only 1 subagent for build/test to avoid conflicts

**Gap:** No JTBD-based spec documenting subagent system design.

**Acceptance criteria:**
- Documents JTBD: "Scale work through parallel execution"
- Activities: identify parallel work, spawn subagents, coordinate results
- Algorithm: decide phase identifies parallelism → act phase spawns subagents
- Edge cases: build/test serialization, context isolation
- Implementation mapping: references act_build.md, decide_build.md

**Dependencies:** build-procedure.md (references build workflow)

---

## Summary

**Total tasks:** 10 specifications to create

**Priority breakdown:**
- Priority 1 (Core Framework): 3 specs - rooda-loop-execution, procedure-system, component-system
- Priority 2 (Workflows): 3 specs - bootstrap-procedure, build-procedure, planning-procedures
- Priority 3 (Quality/Integration): 2 specs - backpressure-system, agents-md-lifecycle
- Priority 4 (Advanced): 2 specs - context-management, subagent-system

**Rationale:** Core framework specs enable understanding the system architecture. Workflow specs document how to use the system. Quality/integration specs explain the feedback mechanisms. Advanced specs cover optimization strategies.

**Implementation approach:** Each spec follows TEMPLATE.md structure with JTBD, Activities, Acceptance Criteria, Algorithm, Data Structures, Edge Cases, Implementation Mapping, and Examples sections.
