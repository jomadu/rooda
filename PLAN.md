# Specification Quality Refactoring Plan

## Quality Assessment Results

### Criterion 1: All specs have "Job to be Done" section
**Status:** FAIL

**Findings:**
- 8 of 9 specs have "## Job to be Done" section
- agents-md-format.md uses "## Purpose" instead of "## Job to be Done"
- DEPRECATED specs (component-system.md, prompt-composition.md) have "## Job to be Done" sections

**Failed Specs:**
- specs/agents-md-format.md

### Criterion 2: All specs have "Acceptance Criteria" section
**Status:** FAIL

**Findings:**
- 8 of 9 specs have "## Acceptance Criteria" section
- agents-md-format.md missing "## Acceptance Criteria" section
- DEPRECATED specs have "## Acceptance Criteria" sections

**Failed Specs:**
- specs/agents-md-format.md

### Criterion 3: All specs have "Examples" section
**Status:** FAIL

**Findings:**
- 8 of 9 specs have "## Examples" section
- agents-md-format.md missing "## Examples" section
- DEPRECATED specs have "## Examples" sections

**Failed Specs:**
- specs/agents-md-format.md

### Criterion 4: All command examples in specs are verified working
**Status:** FAIL

**Findings:**
- No verification process defined in AGENTS.md
- No empirical testing performed on command examples
- No distinction between executable commands and pseudocode/illustrative examples
- Command examples exist across all specs but verification status unknown

**Impact:**
- Specs may contain outdated or incorrect commands
- Users may encounter failures when following documentation
- No systematic way to ensure specs remain accurate as implementation evolves

### Criterion 5: No specs marked as DEPRECATED without replacement
**Status:** PASS

**Findings:**
- 2 specs marked DEPRECATED: component-system.md, prompt-composition.md
- Both reference replacement: component-authoring.md
- Replacement spec exists and is complete

## Refactoring Tasks

### Task 1: Fix agents-md-format.md Structure (CRITICAL)
**Priority:** 1 (High Impact, Blocks Criteria 1-3)

**Description:**
Restructure agents-md-format.md to follow TEMPLATE.md format:
- Replace "## Purpose" with "## Job to be Done"
- Add "## Acceptance Criteria" section with boolean criteria
- Add "## Examples" section with concrete AGENTS.md examples

**Rationale:**
agents-md-format.md is a specification like all others and should follow the same structure. Current structure causes criteria 1, 2, and 3 to fail. This is the critical path blocker.

**Acceptance Criteria:**
- agents-md-format.md has "## Job to be Done" section
- agents-md-format.md has "## Acceptance Criteria" section
- agents-md-format.md has "## Examples" section
- Content from "## Purpose" migrated to "## Job to be Done"
- All existing content preserved and reorganized

### Task 2: Define Command Verification Process
**Priority:** 2 (High Impact, Enables Task 3)

**Description:**
Define systematic process for verifying command examples in specs:
- Identify categories: executable commands vs pseudocode/illustrative examples
- Define verification approach per category
- Document verification process in AGENTS.md quality criteria section
- Create verification tracking mechanism

**Rationale:**
Criterion 4 requires verification but no process exists. Must define process before executing verification pass.

**Acceptance Criteria:**
- Verification process documented in AGENTS.md
- Clear distinction between executable and non-executable examples
- Verification approach defined for each category
- Tracking mechanism established

### Task 3: Execute Verification Pass on All Specs
**Priority:** 3 (High Impact, Depends on Task 2)

**Description:**
Systematically verify all command examples across all specs:
- Identify all bash code blocks in specs/*.md
- Categorize each as executable or illustrative
- Execute all executable commands and validate output
- Document verification results
- Fix or mark any failing commands

**Rationale:**
Empirical testing ensures specs remain accurate. Depends on Task 2 defining the verification process.

**Acceptance Criteria:**
- All bash code blocks identified and categorized
- All executable commands tested
- Verification results documented
- Failing commands fixed or marked as non-executable

### Task 4: Mark Non-Executable Examples Clearly
**Priority:** 4 (Medium Impact, Improves Clarity)

**Description:**
Add clear markers to distinguish pseudocode/illustrative examples from executable commands:
- Add comments or labels to non-executable code blocks
- Use consistent notation (e.g., "# Pseudocode", "# Example pattern")
- Update TEMPLATE.md with guidance on marking examples

**Rationale:**
Prevents confusion between commands that should work as-is vs patterns that need adaptation.

**Acceptance Criteria:**
- All non-executable examples clearly marked
- Consistent notation used across all specs
- TEMPLATE.md updated with marking guidance

### Task 5: Automate Verification Where Possible
**Priority:** 5 (Low Impact, Long-term Improvement)

**Description:**
Create automation to verify command examples remain accurate:
- Script to extract and execute testable commands
- CI integration to run verification on spec changes
- Automated reporting of verification failures

**Rationale:**
Automation ensures ongoing accuracy without manual effort. Lower priority as manual verification (Task 3) provides immediate value.

**Acceptance Criteria:**
- Verification script created
- Script can extract and test commands from specs
- CI integration documented (implementation optional)

## Summary

**Critical Path:**
1. Fix agents-md-format.md structure (Task 1) → Achieves criteria 1-3 compliance
2. Define verification process (Task 2) → Enables criterion 4 compliance
3. Execute verification pass (Task 3) → Achieves criterion 4 compliance

**Current Status:**
- Criteria 1, 2, 3: FAIL (single spec non-compliant)
- Criterion 4: FAIL (no verification process)
- Criterion 5: PASS

**Expected Outcome:**
All 5 quality criteria pass after completing Tasks 1-3.
