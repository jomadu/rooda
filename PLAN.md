# Draft Plan: Spec to Implementation Gap Analysis

## Critical Gaps (P0)

### 1. Git Commit Strategy Specification Missing
**Gap:** Specs describe git push after each iteration (iteration-loop.md, ai-cli-integration.md) but no spec defines when/how commits are created. Implementation only pushes, never commits.

**Impact:** AI CLI must handle commits autonomously, but this critical requirement is not specified. Current implementation assumes kiro-cli commits changes, but this is not documented.

**Tasks:**
- Create `specs/git-workflow.md` spec defining commit strategy
- Document when commits should occur (after successful iteration, after tests pass, etc.)
- Specify commit message format and conventions
- Define how AI CLI creates commits vs how rooda.sh pushes them
- Clarify responsibility boundary between rooda.sh and kiro-cli

**Acceptance Criteria:**
- Spec clearly defines commit timing and responsibility
- Spec documents commit message conventions
- Spec explains push-only behavior in rooda.sh

---

## High Priority Gaps (P1)

### 2. Error Handling Not Specified or Implemented
**Gap:** ai-cli-integration.md identifies lack of error handling as "Known Issue" but no spec defines error handling strategy. Implementation has no error checking for kiro-cli exit status.

**Impact:** Loop continues after AI CLI failures, potentially pushing broken changes or looping indefinitely on repeated failures.

**Tasks:**
- Create `specs/error-handling.md` spec defining error handling strategy
- Specify kiro-cli exit status checking
- Define failure thresholds (N consecutive failures → abort)
- Document error recovery strategies
- Specify what errors are fatal vs recoverable

**Acceptance Criteria:**
- Spec defines error detection mechanisms
- Spec specifies failure thresholds and abort conditions
- Spec documents error recovery strategies

### 3. Dependency Checking Incomplete
**Gap:** external-dependencies.md specifies all dependencies but implementation only checks yq. kiro-cli and bd failures discovered at runtime, not startup.

**Impact:** Users discover missing tools late in execution, wasting time and causing confusing errors.

**Tasks:**
- Implement startup checks for kiro-cli (verify command exists)
- Implement startup checks for bd (verify command exists)
- Add version validation for yq (ensure v4+, not v3)
- Provide clear error messages with installation instructions
- Update external-dependencies.md to document implemented checks

**Acceptance Criteria:**
- rooda.sh checks for kiro-cli at startup
- rooda.sh checks for bd at startup
- rooda.sh validates yq version >= 4.0.0
- Error messages include installation instructions

### 4. Safety and Sandboxing Not Specified
**Gap:** README.md and ai-cli-integration.md mention sandboxing requirements (Docker, Fly Sprites, E2B) but no spec defines safety requirements or sandboxing strategy.

**Impact:** Users may run framework in unsafe environments without understanding blast radius risks.

**Tasks:**
- Create `specs/safety-sandboxing.md` spec defining safety requirements
- Document required isolation mechanisms
- Specify minimum viable access principles
- Define blast radius containment strategies
- Document unsafe operations and their risks

**Acceptance Criteria:**
- Spec defines required isolation mechanisms
- Spec documents blast radius containment
- Spec provides concrete sandboxing examples

---

## Medium Priority Gaps (P2)

### 5. Prompt Component Structure Not Specified
**Gap:** prompt-composition.md documents how prompts are assembled but no spec defines the structure/format of individual component files.

**Impact:** No guidance for creating custom OODA components. Users don't know what conventions to follow.

**Tasks:**
- Create `specs/prompt-component-format.md` spec defining component structure
- Document markdown conventions for OODA phase files
- Specify section header formats (## O1, ## R5, etc.)
- Define naming conventions for component files
- Provide examples of well-structured components

**Acceptance Criteria:**
- Spec defines component file structure
- Spec documents naming conventions
- Spec provides component examples

### 6. Procedure Metadata Not Fully Utilized
**Gap:** configuration-schema.md defines display/summary/description fields but implementation doesn't use them. No help text generation, no procedure listing.

**Impact:** Users can't discover available procedures or understand their purposes without reading config file.

**Tasks:**
- Implement `--help` flag showing available procedures
- Implement `--list-procedures` flag showing all procedures with summaries
- Generate help text from procedure metadata
- Update cli-interface.md spec to document help functionality

**Acceptance Criteria:**
- `./rooda.sh --help` shows usage and available procedures
- `./rooda.sh --list-procedures` shows all procedures with summaries
- Help text generated from config metadata

### 7. Iteration Timing and Progress Not Specified
**Gap:** iteration-loop.md identifies "iteration timing" as area for improvement but no spec defines timing/progress requirements.

**Impact:** Users have no visibility into iteration performance or progress during long-running procedures.

**Tasks:**
- Create `specs/iteration-progress.md` spec defining progress reporting
- Specify elapsed time display per iteration
- Define progress indicators during OODA phases
- Document iteration summary format

**Acceptance Criteria:**
- Spec defines timing display requirements
- Spec specifies progress indicator format
- Spec documents iteration summary structure

---

## Low Priority Gaps (P3)

### 8. Duplicate Validation Code
**Gap:** cli-interface.md identifies duplicate validation blocks (lines 95-103 and 117-125) as "Known Issue" but no refactoring spec exists.

**Impact:** Code duplication makes maintenance harder and increases bug risk.

**Tasks:**
- Refactor duplicate validation into single function
- Update cli-interface.md to remove "Known Issue" note
- Verify all validation paths still work correctly

**Acceptance Criteria:**
- Validation logic exists in single location
- All invocation modes still validate correctly
- cli-interface.md updated to reflect refactoring

### 9. Short Flag Support Missing
**Gap:** cli-interface.md identifies lack of short flags as "Area for Improvement" but no spec defines short flag mappings.

**Impact:** Verbose command lines for explicit flag invocation.

**Tasks:**
- Create `specs/cli-flags.md` spec defining short flag mappings
- Implement short flags (-o, -r, -d, -a, -m, -c)
- Update cli-interface.md with short flag documentation
- Add short flag examples to README.md

**Acceptance Criteria:**
- Spec defines short flag mappings
- Implementation supports short flags
- Documentation includes short flag examples

### 10. Version and Help Flags Missing
**Gap:** cli-interface.md identifies missing --help and --version flags as "Areas for Improvement".

**Impact:** Users must trigger errors to see usage, no way to check framework version.

**Tasks:**
- Implement `--version` flag showing rooda.sh version
- Implement `--help` flag showing usage and procedures
- Define version numbering scheme for framework
- Update cli-interface.md with flag documentation

**Acceptance Criteria:**
- `./rooda.sh --version` shows version number
- `./rooda.sh --help` shows comprehensive usage
- Version scheme documented

---

## Documentation Gaps (P2)

### 11. Spec Index Incomplete
**Gap:** specs/README.md exists but external-dependencies.md is not listed in index.

**Impact:** Spec is discoverable via filesystem but not via index.

**Tasks:**
- Add external-dependencies.md to specs/README.md index
- Verify all other specs are indexed
- Update spec index generation logic if automated

**Acceptance Criteria:**
- All specs listed in specs/README.md
- Index includes JTBD for each spec

---

## Summary

**Total Gaps Identified:** 11

**By Priority:**
- P0 (Critical): 1 gap
- P1 (High): 3 gaps
- P2 (Medium): 4 gaps
- P3 (Low): 3 gaps

**By Type:**
- Missing Specs: 6 gaps (git-workflow, error-handling, safety-sandboxing, prompt-component-format, iteration-progress, cli-flags)
- Incomplete Implementation: 3 gaps (dependency checking, procedure metadata, duplicate validation)
- Documentation: 2 gaps (spec index, help/version flags)

**Recommended Sequence:**
1. Git commit strategy spec (P0) - Critical for understanding current behavior
2. Error handling spec (P1) - Prevents loop failures
3. Dependency checking implementation (P1) - Better user experience
4. Safety/sandboxing spec (P1) - Critical for safe operation
5. Prompt component format spec (P2) - Enables custom procedures
6. Spec index update (P2) - Quick documentation fix
7. Remaining P2/P3 gaps as time permits
