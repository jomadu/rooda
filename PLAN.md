# Draft Plan: Consolidate src/README.md into Specs

## Priority Tasks

1. **Merge prompt composition content into component-authoring.md**
   - Add "Prompt Assembly Mechanism" subsection within existing "Algorithm" section
   - Document the `create_prompt()` function implementation (heredoc + command substitution)
   - Document how prompt files are assembled into final prompt structure
   - Document the assembled prompt format with OODA section headers
   - Keep existing component-authoring.md structure intact
   - Acceptance: All technical content from src/README.md incorporated into component-authoring.md with coherent flow

2. **Update cross-references in user-documentation.md**
   - Line 113: Change `src/README.md` reference to `component-authoring.md`
   - Line 176: Change `src/README.md` reference to `component-authoring.md`
   - Line 201: Change `src/README.md` reference to `component-authoring.md`
   - Acceptance: All three references point to component-authoring.md

3. **Update cross-references in iteration-loop.md**
   - Line 95: Change `../src/README.md` reference to `component-authoring.md`
   - Acceptance: Reference points to component-authoring.md

4. **Update cross-references in configuration-schema.md**
   - Line 100: Change `../src/README.md` reference to `component-authoring.md`
   - Acceptance: Reference points to component-authoring.md

5. **Update cross-references in ai-cli-integration.md**
   - Line 114: Change `../src/README.md` reference to `component-authoring.md`
   - Acceptance: Reference points to component-authoring.md

6. **Update cross-reference in component-authoring.md itself**
   - Line 129: Remove self-reference to `src/README.md` from "Related specs" section
   - Acceptance: No circular reference to deleted file

7. **Create scripts/ directory and move audit-links.sh**
   - Create `scripts/` directory at project root
   - Move `audit-links.sh` to `scripts/audit-links.sh`
   - Preserve executable permissions
   - Acceptance: `scripts/audit-links.sh` exists and is executable, root `audit-links.sh` deleted

8. **Update AGENTS.md implementation definition**
   - Remove `src/README.md` from implementation patterns (no longer exists)
   - Add `scripts/*.sh` to implementation patterns (new location for utility scripts)
   - Acceptance: AGENTS.md accurately reflects file locations

9. **Delete src/README.md**
   - Remove the file after all content merged and references updated
   - Acceptance: `src/README.md` does not exist

10. **Verify no broken links**
    - Run `scripts/audit-links.sh` to check for broken references
    - Manually verify all updated cross-references resolve correctly
    - Acceptance: No broken links in specs/ or docs/

## Dependencies

- Task 1 must complete before task 9 (content must be merged before deletion)
- Tasks 2-6 must complete before task 9 (references must be updated before deletion)
- Task 7 must complete before task 10 (audit script must be in new location)
- Tasks 1-9 must complete before task 10 (verification happens last)
