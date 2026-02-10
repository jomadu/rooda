# Plan: Implement Fragment-Based Procedures System

## Authority

**Source of Truth**: `./procedures.md` — This file defines the complete fragment-based procedures system and serves as the authoritative specification for all implementation work.

**Important**: The `./procedures.md` file will be removed after plan execution. All content must be copied into specs, not referenced.

## Overview

The `./procedures.md` file defines a fragment-based composition system that replaces the single-file prompt approach described in current specs. This plan ensures all content from `./procedures.md` is copied into the specs directory before `./procedures.md` is removed.

## Acceptance Criteria

- [ ] **Complete Content Copy**: All content from `./procedures.md` must be copied into specs (not referenced)
- [ ] **Fragment Directory Structure**: The complete fragments directory structure from `./procedures.md` must be copied into specs
- [ ] **Schema Compliance**: All specs must contain the exact schema defined in `./procedures.md`
- [ ] **Template System**: Go text/template system from `./procedures.md` must be copied into specs
- [ ] **Built-in Procedures**: All 16 built-in procedures from `./procedures.md` must be copied into specs
- [ ] **Self-Contained Specs**: Specs must be complete without `./procedures.md` dependency

## Tasks

### 1. Create Missing Procedures Spec
**File**: `specs/procedures.md`
**Status**: ✅ COMPLETE
**Authority**: `./procedures.md` (complete content migration required)

**Requirements from `./procedures.md`**:
- Copy complete fragments directory structure from `./procedures.md` fragments section
- Copy exact schema from `./procedures.md` "Procedures Configuration Schema"
- Copy complete template example from `./procedures.md`
- Copy all 16 built-in procedures from `./procedures.md` "Built-in Procedures Config"
- Make spec self-contained (no references to `./procedures.md`)

### 2. Update Prompt Composition Spec
**File**: `specs/prompt-composition.md`
**Status**: ✅ COMPLETE
**Authority**: `./procedures.md` schema and template system

**Current Issue**: Assumes single file per OODA phase, but `./procedures.md` defines fragment arrays
**Required Changes Based on `./procedures.md`**:
- Replace single file logic with fragment array processing per copied schema
- Implement Go text/template processing as shown in copied template example
- Update data structures to match copied fragment format
- Add inline content support as defined in copied content
- Include fragment path resolution rules from `./procedures.md`

### 3. Update Configuration Spec
**File**: `specs/configuration.md`
**Status**: ✅ COMPLETE
**Authority**: `./procedures.md` "Procedures Configuration Schema"

**Current Issue**: Procedure data structure doesn't match `./procedures.md` schema
**Required Changes**:
- Update Procedure struct to match exact schema copied from `./procedures.md`
- Add fragment array validation per copied requirements
- Implement template parameter validation as implied by copied content
- Update path resolution to handle `builtin:` prefix from copied content
- Include schema details copied from `./procedures.md`

### 4. Update README Spec
**File**: `specs/README.md`
**Status**: ✅ COMPLETE
**Authority**: `./procedures.md` for procedure count and system description

**Required Changes**:
- Add `procedures.md` to spec list (missing from current README)
- Update procedure count to 16 (from copied built-in procedures)
- Update system description to reflect fragment-based architecture from copied content
- Update procedure library section with copied information

### 5. Update CLI Interface Spec
**File**: `specs/cli-interface.md`
**Status**: Minor updates for fragment system
**Authority**: `./procedures.md` for OODA phase structure

**Current Issue**: 
1. OODA phase override flags assume single file per phase, but `./procedures.md` defines fragment arrays
2. Inconsistent flag design: `--context` vs `--context-file` creates unnecessary complexity

**Required Changes**:

1. **Unify Context Flags**:
   - Remove `--context-file` flag entirely
   - Change `--context` to accept both file paths and inline content
   - Use file existence heuristic: if value exists as file, read it; otherwise treat as inline content
   - Multiple `--context` flags accumulate in order
   - Update acceptance criteria to remove all `--context-file` references
   - Update examples to show `--context file.txt` instead of `--context-file file.txt`

2. **Update OODA Phase Flags for Fragment Arrays**:
   - Change `--observe`, `--orient`, `--decide`, `--act` to accept multiple values (repeatable flags)
   - Each flag invocation adds one fragment to the phase array
   - Use same file existence heuristic: if value exists as file, treat as path; otherwise treat as inline content
   - Examples:
     - File: `--observe prompts/custom.md`
     - Inline: `--observe "Focus on auth module"`
     - Multiple files: `--observe file1.md --observe file2.md`
     - Mixed: `--observe file.md --observe "Additional context"`
   - Providing any phase flag replaces the entire phase array (not appended to config)
   - Order preserved: fragments processed left-to-right as specified on CLI

3. **Update CLIArgs Data Structure**:
   - Change `Contexts []string` to remain as-is (already array)
   - Remove `ContextFiles []string` field entirely
   - Change `ObserveFile string` to `ObserveFragments []string` (array)
   - Change `OrientFile string` to `OrientFragments []string` (array)
   - Change `DecideFile string` to `DecideFragments []string` (array)
   - Change `ActFile string` to `ActFragments []string` (array)

4. **Update Flag Precedence Resolution Algorithm**:
   - For each context value in `Contexts` array:
     - Check if file exists at that path
     - If exists: read file content and use as context
     - If not exists: use value directly as inline content
   - For each OODA phase fragment value:
     - Check if file exists at that path
     - If exists: create FragmentAction with `{path: <file>, content: "", parameters: nil}`
     - If not exists: create FragmentAction with `{path: "", content: <value>, parameters: nil}`
   - Preserve order from CLI arguments (left to right)
   - CLI-provided fragment arrays completely replace config-defined arrays (no merge)

5. **Update Acceptance Criteria**:
   - Remove: `--context-file <path>` reads context from file
   - Remove: Multiple `--context-file` flags accumulate
   - Remove: `--context` and `--context-file` can be mixed
   - Change: `--context <value>` accepts file path or inline text (file existence check)
   - Add: Multiple `--observe` flags accumulate into fragment array
   - Add: Multiple `--orient` flags accumulate into fragment array
   - Add: Multiple `--decide` flags accumulate into fragment array
   - Add: Multiple `--act` flags accumulate into fragment array
   - Add: OODA phase flags use file existence heuristic (file path vs inline content)
   - Update: "Prompt file paths in CLI overrides" to "Fragment values in CLI overrides"
   - Update: "OODA phase override files validated" to "OODA phase override fragments validated"

6. **Update Validation Section**:
   - File existence check happens during flag resolution (determines file vs inline)
   - No separate validation needed for "file not found" - non-existent paths become inline content
   - Validate inline content is non-empty (error if empty string provided)
   - Add note: to force inline content that looks like a filename, ensure file doesn't exist or use absolute path that doesn't exist

7. **Update Examples Section**:
   - Change: `rooda build --context-file task.md` to `rooda build --context task.md`
   - Add: `rooda build --context "Focus on auth"` (inline content)
   - Add: `rooda build --context task.md --context "Additional notes"` (mixed)
   - Add: `rooda build --observe custom.md` (single file fragment)
   - Add: `rooda build --observe file1.md --observe file2.md` (multiple file fragments)
   - Add: `rooda build --observe "Focus on auth module"` (inline content fragment)
   - Add: `rooda build --observe custom.md --observe "Additional instructions"` (mixed)
   - Update: "Multiple Contexts" example to show unified flag

8. **Update Edge Cases Section**:
   - Remove: "Context File Not Found" edge case (now handled by heuristic)
   - Remove: "OODA Phase File Not Found" edge case (now handled by heuristic)
   - Add: "Empty inline content" - `rooda build --observe ""` produces error
   - Add: "Ambiguous filename" - if user wants inline content "file.md" but file exists, file wins
   - Add: "Force inline content" - use non-existent absolute path or ensure file doesn't exist

9. **Update Design Rationale Notes**:
   - Add note explaining unified flag design (consistency across `--context` and OODA phases)
   - Add note explaining file existence heuristic (intuitive, no special syntax needed)
   - Add note explaining why OODA overrides replace entire phase (predictability)
   - Add note explaining repeatable flag pattern (standard CLI convention, matches `--context`)
   - Add note explaining order preservation (predictable fragment composition)
   - Add note explaining why CLI doesn't support template parameters (config-only feature)
   - Add note about ambiguous filenames (file existence wins, design tradeoff for simplicity)

10. **Update Dependencies Section**:
   - Add reference to `procedures.md` for fragment array structure
   - Update reference to `prompt-composition.md` to mention fragment processing

11. **Update Short Flags**:
   - Keep `-c` for `--context` (already defined)
   - No short flags for OODA phases (less commonly used, avoid namespace pollution)

**Backward Compatibility**:
- No breaking changes - nothing implemented yet, this is the initial design
- Single file usage for OODA phases: `rooda build --observe custom.md`
- Internally converts to single-element fragment array: `[{path: "custom.md"}]`
- Unified `--context` flag handles both files and inline content seamlessly

## Implementation Priority

1. **Critical**: Create `specs/procedures.md` with complete `./procedures.md` content migration
2. **High**: Update `specs/prompt-composition.md` for fragment processing
3. **High**: Update `specs/configuration.md` for schema compliance
4. **Medium**: Update `specs/README.md` for documentation alignment
5. **Low**: Update `specs/cli-interface.md` for fragment system

## Content Copy from `./procedures.md`

### Fragment Directory Structure
The complete fragments directory from `./procedures.md` must be copied into specs:
- 13 observe fragments (read_agents_md.md, scan_repo_structure.md, etc.)
- 20 orient fragments (compare_detected_vs_documented.md, identify_drift.md, etc.)
- 10 decide fragments (determine_sections_to_update.md, check_if_blocked.md, etc.)
- 12 act fragments (write_agents_md.md, write_audit_report.md, etc.)

### Schema Copy
The exact schema from `./procedures.md` "Procedures Configuration Schema" section must be copied:
```yaml
# Copy from ./procedures.md
procedures:
  <procedure-name>:
    display: string
    summary: string
    description: string
    observe: # Array of actions
      - content: string # optional
        path: string # optional
        parameters: # optional
          <param-name>: <param-value>
```

### Template System Copy
The Go text/template system and example from `./procedures.md` must be copied into specs.

### Built-in Procedures Copy
All 16 procedures from `./procedures.md` "Built-in Procedures Config" must be copied:
- agents-sync, build, publish-plan (direct actions)
- audit-spec, audit-impl, audit-agents, audit-spec-to-impl, audit-impl-to-spec (audits)
- 8 planning procedures (draft-plan-{spec,impl}-{feat,fix,refactor,chore})

## Validation Against `./procedures.md`

All spec changes must be validated against `./procedures.md` to ensure:
- No content is lost during copy
- Schema exactly matches the source
- Fragment organization is preserved
- Template examples are identical
- Built-in procedure definitions are complete
- Specs are self-contained after `./procedures.md` removal

---

## Notes

### Task 1 Completion: Create Missing Procedures Spec
**Completed**: specs/procedures.md created with complete content migration from procedures.md

**What was done**:
- Created comprehensive procedures.md specification in specs/ directory
- Copied complete fragment directory structure (55 total fragments across 4 OODA phases)
- Copied exact schema definition for procedures configuration
- Copied complete Go text/template system with example
- Copied all 16 built-in procedure definitions with full configuration
- Added implementation requirements section for fragment loading, template processing, and validation
- Made spec fully self-contained with no references to procedures.md

**Key learnings**:
- The fragment-based system is well-organized with clear separation by OODA phase
- Built-in procedures cover 3 categories: direct actions (3), audits (5), and planning (8)
- Template system uses Go text/template for parameterized fragments
- Path resolution supports both builtin: prefix and relative paths from config file directory
- Each OODA phase uses array concatenation to compose full prompts from fragments
- Spec already follows JTBD template structure from specs/README.md with all required sections

### Task 2 Completion: Update Prompt Composition Spec
**Completed**: specs/prompt-composition.md updated to support fragment-based composition

**What was done**:
- Replaced single-file-per-phase model with fragment array processing
- Added support for inline content via content field in FragmentAction
- Added Go text/template processing for parameterized fragments
- Updated data structures to show FragmentAction with content/path/parameters
- Replaced single-file algorithm with ComposePhasePrompt that processes fragment arrays
- Updated validation to check fragment arrays, content vs path exclusivity, and template syntax
- Added edge cases for fragment validation, template errors, and concatenation
- Updated all examples to show fragment arrays instead of single files
- Added embedded fragments list (55 fragments: 13 observe, 20 orient, 10 decide, 12 act)
- Updated design rationale to explain fragment-based composition benefits

**Key learnings**:
- Fragment arrays enable reusability across procedures (e.g., read_agents_md.md used by multiple procedures)
- Inline content field allows quick customization without creating separate files
- Template processing happens per-fragment, not per-phase, enabling fine-grained parameterization
- Double newline concatenation between fragments provides clear visual separation
- Validation at config load time (fail fast) catches all fragment and template errors before execution
- Empty fragment arrays are valid (not all procedures need all phases)
- Fragment system is backward compatible - can still use single fragment per phase if desired

### Task 3 Completion: Update Configuration Spec
**Completed**: specs/configuration.md updated to support fragment-based procedure schema

**What was done**:
- Updated Procedure struct to use []FragmentAction arrays instead of single string paths for Observe/Orient/Decide/Act
- Added FragmentAction struct definition with Content, Path, and Parameters fields
- Updated YAML schema examples to show fragment arrays with inline content and path options
- Modified config merging algorithm to handle fragment arrays and resolve fragment paths
- Updated resolveFragmentPaths function to handle builtin: prefix and relative path resolution
- Updated validation to check fragment arrays (at least one fragment per phase required)
- Added fragment-level validation for content vs path exclusivity (exactly one required)
- Updated edge cases to cover fragment validation scenarios
- Updated all examples to show fragment arrays (Example 2, 3, 9 modified)
- Updated prompt file loading section to reference procedures.md for fragment loading
- Updated path resolution notes to explain fragment array replacement behavior (not element-by-element merge)

**Key learnings**:
- Fragment arrays replace entire OODA phase when specified in config overlay (not merged element-by-element)
- This is critical for predictability - users specify complete phase composition, not partial modifications
- resolveFragmentPaths only resolves non-builtin paths (builtin: prefix preserved for embedded resources)
- Validation happens at two levels: config load (structure) and fragment load (file existence)
- Content vs path exclusivity enforced at config validation time (fail fast)
- Field-level merge still applies to procedure metadata (display, summary, iteration settings)
- Fragment path resolution happens during config merging, not during fragment loading
- This allows provenance tracking and early validation before procedure execution

### Task 4 Completion: Update README Spec
**Completed**: specs/README.md updated to reflect fragment-based architecture

**What was done**:
- Updated Core Functional Job section to mention fragment-based composition
- Updated "Configure procedures for a team" related job to explain fragment arrays and template parameters
- Updated success metrics to reference 55 built-in fragments enabling 16 procedures
- Updated J2 (Compose and Assemble Prompts) to describe fragment array processing instead of single files
- Updated Execution Engine table entry for prompt-composition to mention fragment arrays and template support
- Updated Procedure Library table entry to show "55 embedded fragments composing 16 built-in procedures"
- Updated prompt-composition spec description in Specification Status to mention fragment arrays, inline content, and Go template parameterization
- All references to "four OODA phase files" replaced with "fragment arrays for each OODA phase"

**Key learnings**:
- README serves as the entry point - must accurately reflect the fragment-based architecture
- Fragment count (55) and procedure count (16) are key metrics that demonstrate reusability
- The phrase "fragment arrays" is more accurate than "phase files" for describing the composition model
- Template parameterization is a key feature that enables fragment reuse across different contexts
- Inline content support is important to mention as it allows quick customization without file creation
- The README now consistently describes the system as fragment-based throughout all sections