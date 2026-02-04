# Specification Refactoring Plan

## Priority 1: Fix Quality Criteria (Critical)

**Issue:** AGENTS.md quality criteria are all subjective and non-boolean, violating the framework's core principle that quality assessment must use PASS/FAIL criteria with clear thresholds.

**Tasks:**
- Replace subjective spec criteria with boolean checks:
  - "All specs have Job to be Done section" (PASS/FAIL)
  - "All specs have Acceptance Criteria section" (PASS/FAIL)
  - "All specs have Examples section" (PASS/FAIL)
  - "All command examples in specs are verified working" (PASS/FAIL)
  - "No specs marked as DEPRECATED without replacement" (PASS/FAIL)

- Replace subjective implementation criteria with boolean checks:
  - "shellcheck passes with no errors" (PASS/FAIL)
  - "All procedures in config have corresponding component files" (PASS/FAIL)
  - "Script executes bootstrap procedure successfully" (PASS/FAIL)
  - "Script executes on macOS without errors" (PASS/FAIL)
  - "Script executes on Linux without errors" (PASS/FAIL)

**Impact:** High - Current criteria cannot be used for quality assessment procedures
**Effort:** Low - Update AGENTS.md with boolean criteria

## Priority 2: Remove Deprecated Specs (High Impact)

**Issue:** Two specs are marked DEPRECATED (prompt-composition.md, component-system.md) but are still referenced by other specs and may confuse users.

**Tasks:**
- Remove prompt-composition.md (superseded by component-authoring.md)
- Remove component-system.md (superseded by component-authoring.md)
- Update any references in other specs to point to component-authoring.md
- Update specs/README.md index if it references deprecated specs

**Impact:** High - Deprecated specs create confusion about which documentation is authoritative
**Effort:** Low - Delete files and update references

## Priority 3: Complete Acceptance Criteria (Medium Impact)

**Issue:** iteration-loop.md has 5 of 7 acceptance criteria unchecked, indicating incomplete specification.

**Tasks:**
- Review iteration-loop.md acceptance criteria
- Verify which criteria are actually met by implementation
- Check or document why criteria are not met
- Add implementation notes if criteria are met but not documented

**Impact:** Medium - Incomplete acceptance criteria make it unclear if spec matches implementation
**Effort:** Low - Review and update checkboxes with verification

## Priority 4: Address Known Issues (Low Priority)

**Issue:** Multiple specs document known issues but don't propose solutions or track remediation.

**Tasks:**
- Review known issues in:
  - external-dependencies.md (4 issues)
  - cli-interface.md (2 issues)
  - iteration-loop.md (3 issues)
  - ai-cli-integration.md (4 issues)
- Determine which issues should become work tracking items
- File issues for problems that need fixing
- Document issues that are acceptable trade-offs

**Impact:** Low - Known issues are documented, not blocking current functionality
**Effort:** Medium - Requires analysis and decision-making per issue
