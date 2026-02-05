# Draft Plan: Spec to Implementation Gaps

## Priority 1: Fix User Documentation Link Validation

**Gap:** Spec `user-documentation.md` has unchecked acceptance criterion "Links between documents work correctly" but `./scripts/audit-links.sh` reports "All links valid".

**Task:** Update `specs/user-documentation.md` acceptance criteria line 21 from `- [ ]` to `- [x]` to reflect that link validation is implemented and passing.

**Acceptance:**
- Line 21 in user-documentation.md shows `- [x] Links between documents work correctly`
- Matches empirical evidence from audit-links.sh passing

**Dependencies:** None

---

## Priority 2: Remove Obsolete "to be created" References

**Gap:** Specs contain references to specs "to be created" that already exist:
- `cli-interface.md` line 170: references `configuration-schema.md` as "to be created" (exists)
- `cli-interface.md` line 171: references `iteration-loop.md` as "to be created" (exists)
- `iteration-loop.md` line 96: references `ai-cli-integration.md` as "to be created" (exists)

**Task:** Remove "(to be created)" annotations from Related specs sections in these three files.

**Acceptance:**
- cli-interface.md lines 170-171 no longer say "to be created"
- iteration-loop.md line 96 no longer says "to be created"
- References remain but without obsolete annotations

**Dependencies:** None

---

## Priority 3: Document Known Issues in CLI Interface

**Gap:** Spec `cli-interface.md` documents three known issues (duplicate validation blocks, inconsistent usage messages, config validation) but implementation may have addressed some of these.

**Task:** Verify each known issue empirically:
1. Search rooda.sh for duplicate validation blocks (lines 95-103 and 117-125 mentioned)
2. Check if usage messages are consistent
3. Check if config validation exists before yq queries
4. Update "Known Issues" section to reflect current state
5. Move resolved issues to implementation notes or remove if no longer relevant

**Acceptance:**
- Each known issue verified against current implementation
- Known Issues section accurately reflects current state
- No false positives (claiming issues that don't exist)

**Dependencies:** None

---

## Priority 4: Verify Quality Criteria Implementation

**Gap:** AGENTS.md documents quality criteria with verification commands but need to confirm all are implemented:
- Prompt validation: `./scripts/validate-prompts.sh` (verified working)
- Link audit: `./scripts/audit-links.sh` (verified working)
- shellcheck: `shellcheck src/rooda.sh` (verified working)
- Bootstrap execution: `./src/rooda.sh bootstrap --max-iterations 1` (needs verification)
- Cross-platform: macOS verified, Linux needs verification

**Task:** 
1. Run bootstrap procedure to verify it executes successfully
2. Document Linux verification status in AGENTS.md
3. If Linux verification missing, note as limitation or add verification

**Acceptance:**
- Bootstrap procedure verified working
- AGENTS.md accurately reflects cross-platform verification status
- Quality criteria section matches empirical reality

**Dependencies:** None

---

## Priority 5: Validate Spec Completeness

**Gap:** All specs have "Job to be Done", "Acceptance Criteria", and "Examples" sections (quality criterion PASS), but need to verify command examples are working vs pseudocode.

**Task:** Audit all command examples in specs to distinguish:
- Executable commands (should be verified working)
- Pseudocode/algorithm descriptions (not meant to execute)
- Mark which examples have been verified
- Update quality criteria verification process if needed

**Acceptance:**
- Clear distinction between executable vs pseudocode examples
- Executable examples verified or marked for verification
- Quality criteria accurately reflects verification status

**Dependencies:** None

---

## Notes

**Why these priorities:**
1. Quick wins first (checkbox updates, annotation removal)
2. Known issues verification ensures specs match reality
3. Quality criteria validation ensures framework integrity
4. Spec completeness ensures documentation accuracy

**Gaps NOT included:**
- TASK.md unchecked items: These are for a different feature (AI CLI configuration refactoring), not spec-to-impl gaps
- Beads issues: Work tracking system contains separate backlog, not spec gaps
- "to be created" in grep results: Most are in prompts/docs, not specs claiming missing implementation

**Empirical findings:**
- All 25 prompt files validated successfully
- shellcheck passes with no errors
- All links validated successfully
- 9 procedures defined in config, all have corresponding prompt files
- No missing OODA phase files
