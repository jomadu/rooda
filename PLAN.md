# Draft Plan: Unify Component System Specifications

## Priority Tasks

### 1. Create specs/component-authoring.md

**Description:** Create unified specification that accurately describes actual component file structure by examining real component files in `src/components/*.md`.

**What to include:**
- **Job to be Done:** Enable developers to create and modify OODA component prompt files that can be composed into executable procedures
- **Component File Structure:** Markdown format with:
  - Phase header: `# [Phase]: [Purpose]` (e.g., `# Act: Bootstrap`)
  - Step headers: `## [Code]: [Name]` (e.g., `## A1: Create AGENTS.md with Operational Guide`)
  - Full prose instructions under each step header
- **Component Pattern:** All examined component files follow the same pattern:
  - Step code headers (O1-O15, R1-R22, D1-D15, A1-A9) serve as section markers
  - Full prose instructions appear under each step header
  - No files contain only step code references without prose
  - Step codes provide structure and cross-reference capability
- **Common Steps Reference:** Complete O1-O15, R1-R22, D1-D15, A1-A9 list with descriptions (from component-system.md and src/README.md)
- **Prompt Assembly Algorithm:** How create_prompt() combines four files with OODA section headers (from prompt-composition.md)
- **Authoring Guidelines:** Key principles for writing components (from component-system.md)
- **Examples:** Show actual excerpts from `act_bootstrap.md`, `observe_bootstrap.md`, `act_build.md` demonstrating structure

**Acceptance criteria:**
- Spec accurately describes structure of files in `src/components/*.md`
- Clarifies that step codes are section headers, not references to external definitions
- Explains dual purpose: structure + cross-reference capability
- Includes prompt assembly algorithm from prompt-composition.md
- Includes common steps reference from component-system.md
- Includes authoring guidelines from component-system.md
- Shows real examples from actual component files

**Dependencies:** None

---

### 2. Deprecate specs/component-system.md

**Description:** Add deprecation notice at top of component-system.md directing readers to new unified spec.

**Content:**
```markdown
> **DEPRECATED:** This specification has been superseded by [component-authoring.md](component-authoring.md), which provides a unified and more accurate description of component file structure, prompt assembly, and authoring guidelines. This file is retained for historical reference only.
```

**Acceptance criteria:**
- Deprecation notice added at top of file
- Notice references component-authoring.md
- Original content preserved

**Dependencies:** Task 1 (component-authoring.md must exist)

---

### 3. Deprecate specs/prompt-composition.md

**Description:** Add deprecation notice at top of prompt-composition.md directing readers to new unified spec.

**Content:**
```markdown
> **DEPRECATED:** This specification has been superseded by [component-authoring.md](component-authoring.md), which provides a unified description of both component file structure and prompt assembly. This file is retained for historical reference only.
```

**Acceptance criteria:**
- Deprecation notice added at top of file
- Notice references component-authoring.md
- Original content preserved

**Dependencies:** Task 1 (component-authoring.md must exist)

---

### 4. Update specs/README.md

**Description:** Update spec index to list new component-authoring.md and mark deprecated specs.

**Changes needed:**
- Add component-authoring.md to the list with its JTBD
- Mark component-system.md as (deprecated)
- Mark prompt-composition.md as (deprecated)
- Maintain alphabetical or logical ordering

**Acceptance criteria:**
- New spec listed with extracted JTBD
- Deprecated specs marked clearly
- Index structure follows specification-system.md requirements

**Dependencies:** Task 1 (component-authoring.md must exist)

---

### 5. Update src/README.md reference

**Description:** Update the note at top of src/README.md to reference component-authoring.md as authoritative source instead of component-system.md.

**Current text:**
```markdown
**Note:** This is a quick reference guide. For the authoritative specification of the component system, including detailed A6 operational learning criteria, see [specs/component-system.md](../specs/component-system.md).
```

**New text:**
```markdown
**Note:** This is a quick reference guide. For the authoritative specification of component authoring, including file structure, prompt assembly, and detailed A6 operational learning criteria, see [specs/component-authoring.md](../specs/component-authoring.md).
```

**Acceptance criteria:**
- Reference updated to component-authoring.md
- Description reflects unified nature of new spec
- Link path is correct

**Dependencies:** Task 1 (component-authoring.md must exist)

---

## Summary

This plan creates a single unified specification that accurately describes:
1. What component files actually contain (full prose, not just step codes)
2. How component files are structured (phase header, step headers, instructions)
3. How prompt assembly works (from prompt-composition.md)
4. What common steps exist (from component-system.md)
5. How to write new components (authoring guidelines)

The old specs are deprecated but preserved for historical reference. All indexes and references are updated to point to the new authoritative source.
