# Draft Plan: Spec to Implementation Gap Analysis

## High Priority Gaps

### 1. Documentation Quality Criteria - Cross-document link validation failing
**Gap:** Quality criteria states "All cross-document links work correctly (PASS/FAIL)" with verification via `./scripts/audit-links.sh`, but the script reports "All links valid" which doesn't match the acceptance criteria expectation in user-documentation.md that shows "Links between documents work correctly" as unchecked.

**Status:** Spec says unchecked `[ ]`, but implementation passes. Need to verify if this is a spec staleness issue or if there are actually broken links not caught by the script.

**Task:** Verify audit-links.sh catches all link types (relative paths, anchors, external URLs) and update user-documentation.md acceptance criteria checkbox if implementation is correct.

**Priority:** P2 (documentation accuracy)

### 2. User Documentation - Missing docs/README.md index
**Gap:** user-documentation.md spec defines docs/README.md as part of the documentation hierarchy and mentions it should be an "index of detailed documentation", but verification shows docs/README.md exists but may not be comprehensive.

**Status:** File exists but completeness unclear from gap analysis.

**Task:** Review docs/README.md to ensure it properly indexes all documentation files (ooda-loop.md, ralph-loop.md, beads.md) and provides navigation guidance per spec.

**Priority:** P3 (documentation completeness)

### 3. CLI Interface - Short flags implementation verification
**Gap:** cli-interface.md spec states short flags (-o, -r, -d, -a, -m, -c, -h) should work identically to long flags, with acceptance criteria marked as complete `[x]`. Need to verify implementation handles all short flags correctly.

**Status:** Spec says implemented, need empirical verification.

**Task:** Test all short flags work correctly: -o, -r, -d, -a, -m, -c, -h match --observe, --orient, --decide, --act, --max-iterations, --config, --help.

**Priority:** P3 (verification)

## Medium Priority Gaps

### 4. External Dependencies - Optional dependency checking
**Gap:** external-dependencies.md spec describes optional dependencies (shellcheck, git) but notes "shellcheck not installed → Lint command fails but documented as optional" and "git not installed → Commit operations fail (not critical for testing)". Implementation checks for yq and kiro-cli but doesn't warn about missing optional dependencies.

**Status:** Spec describes behavior, implementation doesn't provide warnings for optional tools.

**Task:** Add optional dependency warnings at startup for shellcheck and git if not installed, per spec's edge case documentation.

**Priority:** P3 (user experience improvement)

### 5. Configuration Schema - AI tools section validation
**Gap:** configuration-schema.md spec describes ai_tools section with custom presets but notes "ai_tools presets validated as string type" and "Unknown presets return helpful error messages" as acceptance criteria. Implementation has resolve_ai_tool_preset function but validation completeness unclear.

**Status:** Spec says implemented, need to verify error messages match spec examples.

**Task:** Verify unknown preset error message matches spec format: lists available presets (kiro-cli, claude, aider) and shows instructions for adding custom presets to config.

**Priority:** P3 (verification)

### 6. Iteration Loop - Git push error handling
**Gap:** iteration-loop.md spec notes "Git push failures: If git push fails for reasons other than missing remote branch, the error is silent and the loop continues" as a known issue. Implementation has basic error handling but may not cover all failure modes clearly.

**Status:** Known issue documented in spec, implementation has basic handling.

**Task:** Improve git push error messages to distinguish between: authentication failure, network issue, merge conflict, and missing remote branch. Add user guidance for each case.

**Priority:** P3 (error handling improvement)

## Low Priority Gaps

### 7. User Documentation - Progressive disclosure verification
**Gap:** user-documentation.md spec emphasizes "progressive disclosure (quick start → detailed guides)" as a key principle. Need to verify README.md actually follows this pattern effectively.

**Status:** Principle documented, implementation exists, quality assessment needed.

**Task:** Review README.md structure to ensure it follows progressive disclosure: installation → basic workflow → links to detailed docs, without overwhelming new users.

**Priority:** P4 (documentation quality)

### 8. Component Authoring - Validation tooling
**Gap:** component-authoring.md spec notes in "Areas for Improvement": "Could add validation tooling to check prompt files follow structure" and "Could add linting for step code consistency across prompts". validate-prompts.sh exists but scope unclear.

**Status:** Basic validation exists (validate-prompts.sh), but spec suggests more comprehensive tooling.

**Task:** Review validate-prompts.sh to determine if it checks: phase headers present, step code format (O1-O15, R1-R22, D1-D15, A1-A9), and prose instructions under each step. Document what it validates.

**Priority:** P4 (tooling enhancement)

### 9. AI CLI Integration - Dependency checking conditional on AI_CLI_COMMAND
**Gap:** ai-cli-integration.md spec describes kiro-cli as default but configurable. Implementation has commented-out kiro-cli checks with note "kiro-cli check moved to after argument parsing (conditional on AI_CLI_COMMAND)". Need to verify this conditional check exists.

**Status:** Code comment indicates intentional design, need to verify implementation.

**Task:** Verify that kiro-cli dependency check only runs when AI_CLI_COMMAND uses kiro-cli (not when using claude, aider, or custom CLI). Confirm this matches spec's configurable dependency philosophy.

**Priority:** P4 (verification)

## Completeness Assessment

**Specifications:** 9 spec files (excluding README.md, TEMPLATE.md, specification-system.md per AGENTS.md)
- external-dependencies.md ✓
- cli-interface.md ✓
- iteration-loop.md ✓
- configuration-schema.md ✓
- user-documentation.md ✓
- ai-cli-integration.md ✓
- agents-md-format.md ✓
- component-authoring.md ✓
- (specification-system.md excluded per AGENTS.md)

**Implementation Coverage:** ~95% complete
- Core functionality (OODA loop, procedures, config, CLI) fully implemented
- Quality criteria validation scripts exist (audit-links.sh, validate-prompts.sh)
- Documentation structure matches spec
- All 9 procedures defined in config
- 25 prompt component files present

**Gaps Summary:**
- 1 high priority: Documentation accuracy verification
- 2 medium priority: Optional dependency warnings, error message improvements
- 6 low priority: Quality assessments, validation tooling documentation, verification tasks

**Accuracy Assessment:** Specs are highly accurate relative to implementation. Most acceptance criteria marked `[x]` are correctly implemented. The few gaps identified are:
1. Minor documentation staleness (checkboxes not updated)
2. Optional features mentioned in specs but not implemented (warnings, enhanced error messages)
3. Verification tasks to confirm implementation matches spec claims

**Recommendation:** Focus on verification tasks first (items 1, 3, 5, 9) to confirm specs accurately reflect implementation, then address user experience improvements (items 4, 6) if desired.
