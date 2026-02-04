# Specification Refactoring Plan

## Quality Assessment Results

**Criterion 1: All specs have "Job to be Done" section** - PASS
- All 9 specs contain "## Job to be Done" section

**Criterion 2: All specs have "Acceptance Criteria" section** - PASS
- All 9 specs contain "## Acceptance Criteria" section

**Criterion 3: All specs have "Examples" section** - PASS
- All 9 specs contain "## Examples" section

**Criterion 4: All command examples in specs are verified working** - FAIL
- 44 bash code blocks identified across 8 specs
- No verification process defined
- No empirical testing performed
- No distinction between executable commands and pseudocode

**Criterion 5: No specs marked as DEPRECATED without replacement** - PASS
- 2 deprecated specs (component-system.md, prompt-composition.md)
- Both reference replacement: component-authoring.md
- Replacement exists and is complete

## Refactoring Tasks

### Task 1: Define Command Example Verification Process
**Priority:** Critical (blocks Task 2)
**Impact:** High - establishes methodology for criterion 4

Define verification process in AGENTS.md or new spec:
- What constitutes "executable" vs "pseudocode" example
- How to test executable commands (run and validate output)
- How to mark non-executable examples (comment, label, or section)
- What validation criteria apply (exit code, output format, side effects)
- How to track verification status per spec

### Task 2: Execute Verification Pass on All Specs
**Priority:** High (depends on Task 1)
**Impact:** High - validates all 44 bash code blocks

Systematically verify each bash code block:
- external-dependencies.md (8 blocks)
- cli-interface.md (8 blocks)
- component-system.md (8 blocks)
- iteration-loop.md (6 blocks)
- ai-cli-integration.md (5 blocks)
- configuration-schema.md (4 blocks)
- prompt-composition.md (4 blocks)
- component-authoring.md (1 block)

For each block:
- Classify as executable or pseudocode
- If executable: run command and validate
- If pseudocode: mark clearly as non-executable
- Document results

### Task 3: Mark Non-Executable Examples Clearly
**Priority:** Medium (can parallel with Task 2)
**Impact:** Medium - prevents confusion

Update specs to distinguish:
- Executable commands: verified working, can be copy-pasted
- Pseudocode: illustrative only, marked with comment or label
- Examples: `# Pseudocode - not executable` or section header

### Task 4: Create Verification Tracking System
**Priority:** Low (after Task 2 completes)
**Impact:** Medium - enables ongoing maintenance

Track verification status:
- Per-spec verification checklist
- Last verified date
- Known issues or limitations
- Re-verification triggers (implementation changes)

### Task 5: Automate Verification Where Possible
**Priority:** Low (after Task 2 completes)
**Impact:** Low - reduces manual effort

Create automation:
- Script to extract bash blocks from specs
- Script to execute and validate commands
- CI integration to re-verify on changes
- Report generation for verification status
