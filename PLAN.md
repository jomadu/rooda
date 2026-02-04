# Draft Plan: Spec to Implementation Gap Analysis

## Priority 1: AI CLI Configuration Support (CRITICAL)

**Gap:** Specs define `ai_cli_command` field in rooda-config.yml and `--ai-cli` flag for configurable AI CLI tools, but implementation hardcodes `kiro-cli chat --no-interactive --trust-all-tools` with no configuration support.

**Tasks:**
1. Add `ai_cli_command` field parsing from rooda-config.yml (root level, string type)
2. Add `--ai-cli` flag to argument parser with precedence: flag > config > default
3. Update AI CLI invocation to use resolved command variable instead of hardcoded string
4. Update execution banner to display AI CLI command when non-default
5. Validate AI CLI command is non-empty string when loaded from config

**Acceptance Criteria:**
- `ai_cli_command` field in config overrides default kiro-cli command
- `--ai-cli` flag overrides both config and default
- Default remains `kiro-cli chat --no-interactive --trust-all-tools` for backward compatibility
- Execution banner shows AI CLI command when using non-default

**Dependencies:** None

---

## Priority 2: Remove Hardcoded Dependency Checks (HIGH)

**Gap:** Specs state kiro-cli and bd are configurable/project-specific dependencies, but implementation has hardcoded startup checks that exit if they're missing. This prevents users from substituting alternative AI CLIs or work tracking systems.

**Tasks:**
1. Remove hardcoded kiro-cli dependency check (lines 76-80)
2. Remove hardcoded bd dependency check (lines 82-91)
3. Remove kiro-cli version check (lines 93-102)
4. Remove bd version check (lines 104-112)
5. Keep yq dependency check (required for config parsing)
6. Update external-dependencies.md to reflect that dependency checks are removed for configurable tools

**Acceptance Criteria:**
- Script starts successfully without kiro-cli installed (if using alternative AI CLI)
- Script starts successfully without bd installed (if using alternative work tracking)
- yq dependency check remains (required)
- Script fails at runtime with clear error if configured AI CLI is not installed

**Dependencies:** Should be done after Priority 1 (AI CLI configuration) to ensure alternative AI CLIs can be configured

---

## Priority 3: Documentation Cross-Reference Link Validation (MEDIUM)

**Gap:** Quality criteria states "All cross-document links work correctly (PASS/FAIL)" but specs note this is not verified. README.md and specs/ contain numerous cross-references that may be broken.

**Tasks:**
1. Run `scripts/audit-links.sh` to identify broken links
2. Fix broken links in README.md
3. Fix broken links in specs/*.md
4. Fix broken links in docs/*.md
5. Update quality criteria verification process to include link checking

**Acceptance Criteria:**
- All markdown links resolve to existing files or valid URLs
- No broken cross-references between documentation files
- Link audit script runs successfully with no errors

**Dependencies:** None

---

## Priority 4: Verbose and Quiet Mode Implementation (LOW)

**Gap:** CLI interface spec documents `--verbose` and `--quiet` flags with expected behavior, but implementation only partially supports them. Verbose mode shows full prompt but doesn't show detailed execution. Quiet mode suppresses some output but not all.

**Tasks:**
1. Review current verbose/quiet implementation (VERBOSE variable usage)
2. Ensure `--verbose` shows full prompt before execution (already implemented)
3. Ensure `--quiet` suppresses all non-error output (check all echo statements)
4. Add verbose logging for AI CLI invocation details
5. Add verbose logging for git push operations
6. Test both modes to ensure consistent behavior

**Acceptance Criteria:**
- `--verbose` shows full prompt, AI CLI command, git operations, iteration progress
- `--quiet` suppresses all output except errors
- Default mode shows execution banner and iteration progress (current behavior)

**Dependencies:** None

---

## Priority 5: Help Text Generation from Config Metadata (LOW)

**Gap:** Configuration schema spec defines optional `display`, `summary`, and `description` fields for procedures, noting they're "for future help text generation." These fields exist in rooda-config.yml but are not used by the implementation.

**Tasks:**
1. Add `--list-procedures` flag to show available procedures with display names and summaries
2. Update `show_help()` to include procedure listing
3. Consider adding `--describe <procedure>` to show detailed procedure information
4. Use display/summary/description fields from config when generating help text

**Acceptance Criteria:**
- `--list-procedures` shows all available procedures with display names and summaries
- Help text includes procedure listing or reference to `--list-procedures`
- Procedure metadata from config is displayed to users

**Dependencies:** None

---

## Priority 6: Error Handling for AI CLI Failures (LOW)

**Gap:** Specs note "No error handling: Script continues to git push even if AI CLI fails" as a known issue. Implementation intentionally ignores AI CLI exit status per design, but could provide better user feedback.

**Tasks:**
1. Capture AI CLI exit status (without changing loop behavior)
2. Log warning if AI CLI exits with non-zero status
3. Consider adding `--strict` flag to exit loop on AI CLI failure (optional enhancement)
4. Update ai-cli-integration.md to document error handling behavior

**Acceptance Criteria:**
- AI CLI exit status is captured and logged
- Loop continues regardless of exit status (preserve current design)
- User is informed when AI CLI fails but loop continues
- Optional: `--strict` flag exits loop on failure

**Dependencies:** None

---

## Notes

**Gaps Not Included:**
- Version validation improvements (yq v4 check) - mentioned in specs as "areas for improvement" but not critical
- Automated installer - mentioned as improvement, not required functionality
- Docker image - mentioned as improvement, not required functionality
- Resume capability - mentioned as improvement, not required functionality
- Timeout mechanism - mentioned as improvement, not required functionality

**Implementation Verification:**
- All 25 prompt files exist in src/prompts/ (verified via glob)
- All 9 procedures defined in config have corresponding prompt files (verified via config review)
- create_prompt() function exists and matches spec (verified lines 397-416)
- Iteration loop exists and matches spec (verified lines 418-452)
- Argument parsing exists and matches spec (verified lines 219-295)

**Quality Criteria Status:**
- Specifications: All specs have JTBD, Acceptance Criteria, Examples sections ✓
- Implementation: shellcheck not run yet (should be run before commit)
- Documentation: Cross-reference links not verified (Priority 3 task)
