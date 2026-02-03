# Draft Plan: Spec to Implementation Gap Analysis

## Priority 1 (High Impact)

### Task 1: Automate Spec Creation from Template
**Description:** Create helper script or procedure to generate new spec files from TEMPLATE.md

**Rationale:** Currently agents must manually copy TEMPLATE.md when creating new specs. Automation ensures consistency and reduces friction.

**Acceptance Criteria:**
- [ ] Script/command creates new spec file from TEMPLATE.md
- [ ] Prompts for spec name and basic metadata
- [ ] Places file in specs/ directory with correct naming convention
- [ ] Preserves all template sections and structure

**Implementation Approach:**
- Add `create-spec.sh` helper script to src/
- Or add `create-spec` procedure to rooda-config.yml
- Use sed/awk to replace placeholder text in template

**Dependencies:** None

---

### Task 2: Validate Spec Structure Against Template
**Description:** Add validation that spec files follow TEMPLATE.md structure

**Rationale:** Specs should consistently follow the template structure (JTBD, Activities, Acceptance Criteria, etc.). Validation catches deviations early.

**Acceptance Criteria:**
- [ ] Script checks for required template sections
- [ ] Reports missing or out-of-order sections
- [ ] Can be run manually or as part of quality assessment
- [ ] Returns exit code 0 for valid, 1 for invalid

**Implementation Approach:**
- Add `validate-spec.sh` script to src/
- Parse markdown headers to verify structure
- Compare against TEMPLATE.md section order

**Dependencies:** Task 1 (needs template structure defined)

---

## Priority 2 (Medium Impact)

### Task 3: Auto-Generate specs/README.md Index
**Description:** Generate specs/README.md index from existing spec files

**Rationale:** specs/README.md should list all JTBDs, topics, and specs. Manual maintenance leads to drift. Auto-generation ensures accuracy.

**Acceptance Criteria:**
- [ ] Script reads all specs/*.md files
- [ ] Extracts JTBD and topic information from each spec
- [ ] Generates specs/README.md with organized index
- [ ] Groups specs by category or JTBD
- [ ] Includes links to individual spec files

**Implementation Approach:**
- Add `generate-spec-index.sh` script to src/
- Parse "Job to be Done" section from each spec
- Build categorized list with links
- Overwrite specs/README.md

**Dependencies:** None

---

### Task 4: Validate JTBD/Topic Pattern in Specs
**Description:** Enforce that spec files follow JTBD → topic of concern pattern

**Rationale:** Specification system defines JTBD as organizing principle. Validation ensures specs follow this methodology.

**Acceptance Criteria:**
- [ ] Script checks each spec has "Job to be Done" section
- [ ] Verifies JTBD is outcome-focused (not mechanism-focused)
- [ ] Can suggest which JTBD a spec belongs to
- [ ] Reports specs that don't fit JTBD pattern

**Implementation Approach:**
- Extend validate-spec.sh or create separate script
- Parse JTBD section and check for keywords
- Cross-reference with specs/README.md JTBDs

**Dependencies:** Task 3 (needs README structure to validate against)

---

## Priority 3 (Low Priority)

### Task 5: Automated Spec Completeness Metrics
**Description:** Check that all template sections are filled with meaningful content

**Rationale:** Specs should be complete, not just structurally valid. Metrics identify incomplete specs.

**Acceptance Criteria:**
- [ ] Script checks each template section has content
- [ ] Identifies placeholder text or empty sections
- [ ] Reports completeness percentage per spec
- [ ] Can be used in quality assessment procedures

**Implementation Approach:**
- Extend validate-spec.sh with content checks
- Look for common placeholder patterns
- Count filled vs empty sections

**Dependencies:** Task 2 (needs validation logic)

---

## Summary

**Total Tasks:** 5
**P1 (High Impact):** 2 tasks
**P2 (Medium Impact):** 2 tasks  
**P3 (Low Priority):** 1 task

**Recommended Sequence:**
1. Task 1 (spec creation automation) - immediate value, no dependencies
2. Task 3 (README generation) - immediate value, no dependencies
3. Task 2 (structure validation) - depends on Task 1
4. Task 4 (JTBD validation) - depends on Task 3
5. Task 5 (completeness metrics) - depends on Task 2

**Note:** All tasks are tooling/automation improvements. The core framework (rooda.sh, OODA loop, AGENTS.md) is fully operational. These tasks reduce manual work and enforce spec quality.
