# Draft Plan: Bug Fix - VALIDATION-*.md Files in Project Root

## Priority 1: Create Component System Specification

### Task 1: Create specs/component-system.md
**Description:** Create formal specification for OODA component system following TEMPLATE.md format. Incorporate common steps content from `src/README.md`.

**Content Requirements:**
- **Job to be Done:** Define how OODA components are structured and how they guide agent behavior
- **Activities:** Component composition, common step reference, AGENTS.md update guidance
- **Acceptance Criteria:**
  - Components follow four-phase OODA structure (observe/orient/decide/act)
  - Common steps defined and referenceable by code (O1-O15, R1-R22, D1-D15, A1-A9) - source from `src/README.md`
  - Key principles for writing components documented (from `src/README.md`)
  - A6 step explicitly defines what qualifies as "operational learning" for AGENTS.md:
    - **Operational:** Commands that failed/succeeded, file paths discovered, quality criteria refined, workflow patterns learned
    - **Non-operational:** Test artifacts, validation patterns, historical notes, temporary debugging files
  - Clear examples distinguish operational vs non-operational updates
  - References `specs/agents-md-format.md` for AGENTS.md content boundaries
- **Data Structures:** Component file format (markdown with step codes and descriptions)
- **Algorithm:** How components reference common steps, how agents interpret step codes, how procedures compose components
- **Edge Cases:** Missing step codes, undefined steps, conflicting guidance between components
- **Dependencies:** Requires `specs/agents-md-format.md` (already exists)
- **Implementation Mapping:** `src/components/*.md` files implement this spec
- **Examples:** 
  - Proper A6 update: "Command `npm test` failed, updated AGENTS.md with correct command `npm run test:unit`"
  - Improper A6 update: "Created VALIDATION-issue-123.md with test cases" (should NOT update AGENTS.md)

**Why This Matters:** Without this spec, components lack clear guidance on AGENTS.md updates, leading to test artifacts being added inappropriately.

---

## Priority 2: Clean Up Validation Files

### Task 2: Delete all VALIDATION-*.md files from project root
**Description:** Remove 18 existing VALIDATION-*.md files that clutter the repository.

**Acceptance Criteria:**
- All files matching pattern `VALIDATION-*.md` deleted from project root
- No validation files remain in project root
- Commit message documents cleanup rationale

**Why This Matters:** Immediate clutter removal, prevents confusion about what these files are for.

---

## Priority 3: Update Documentation

### Task 3: Update src/README.md to reference component spec
**Description:** Modify `src/README.md` to reference the new `specs/component-system.md` specification.

**Acceptance Criteria:**
- Add reference to `specs/component-system.md` at top of src/README.md
- Clarify that src/README.md is a quick reference, spec is authoritative
- Maintain existing content (procedure list, component list, common steps)
- Add note about A6 guidance being defined in the spec

**Why This Matters:** Maintains consistency between documentation and specifications, makes spec discoverable.

---

## Summary

**Total Tasks:** 3

**Dependencies:**
- Task 1 must complete first (defines the spec)
- Tasks 2 and 3 can run in parallel after Task 1

**Root Cause Addressed:** Missing specification for component system led to unclear A6 guidance, which caused validation files to be created in project root.

**Expected Outcome:** 
- Components have formal specification defining behavior
- Clear guidance on what belongs in AGENTS.md vs what doesn't
- Project root cleaned of validation file clutter
- Future build procedures won't create validation files inappropriately
