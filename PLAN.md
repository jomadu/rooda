# Draft Plan: Spec to Implementation Gap Analysis

## Critical Gaps (Missing Core Features)

### 1. AI CLI Integration - Missing ai_cli_command config field
**Gap:** Spec `ai-cli-integration.md` describes `ai_cli_command` field in config for setting default AI CLI, but `rooda-config.yml` doesn't have this field and implementation doesn't query it.

**Current behavior:** Uses hardcoded default `kiro-cli chat --no-interactive --trust-all-tools`, overridable only via `$ROODA_AI_CMD` env var or `--ai-cmd` flag.

**Spec requirement:** Four-tier precedence: `--ai-cmd` flag > `--ai-cmd-preset` > `$ROODA_AI_CMD` > config `ai_cli_command` > default

**Implementation needed:**
- Add `ai_cli_command` field to config schema documentation
- Query config for `ai_cli_command` when resolving AI CLI command
- Insert into precedence chain between `$ROODA_AI_CMD` and default

**Acceptance:** Config with `ai_cli_command: "custom-cli"` uses that command when no flag/preset/env var specified

---

### 2. External Dependencies - Missing bd/beads dependency check
**Gap:** Spec `external-dependencies.md` lists bd (beads) as project-specific dependency with version check, but implementation doesn't check for bd at startup.

**Current behavior:** No bd availability check. Commands fail at runtime if bd not installed.

**Spec requirement:** Check for bd at startup, provide installation instructions if missing.

**Implementation needed:**
- Add bd availability check in dependency section (after yq, before kiro-cli)
- Check bd version >= 0.1.0 (requires --json flag support)
- Provide installation instructions: `cargo install beads-cli`
- Note: This is project-specific, not framework-required

**Acceptance:** Running `./rooda.sh bootstrap` without bd installed shows error with installation instructions

---

### 3. CLI Interface - Missing --version output format
**Gap:** Spec `cli-interface.md` Example 10 shows version output as `rooda.sh version 0.1.0`, but implementation outputs `rooda.sh version $VERSION` (variable not expanded in example verification).

**Current behavior:** Implementation correctly outputs version.

**Spec accuracy:** Spec example is correct, no implementation gap. This is a false positive from initial analysis.

**Action:** No implementation change needed. Verify spec example matches implementation.

---

### 4. Iteration Loop - Missing iteration timing display
**Gap:** Spec `iteration-loop.md` "Areas for Improvement" mentions displaying elapsed time per iteration, but not implemented.

**Current behavior:** No timing information displayed.

**Spec requirement:** This is listed as "Areas for Improvement", not acceptance criteria.

**Action:** Not a gap - this is a future enhancement, not a missing feature.

---

## High-Impact Gaps (Frequently Used Functionality)

### 5. Configuration Schema - ai_tools section not documented in config comments
**Gap:** Spec `configuration-schema.md` describes `ai_tools` section for custom presets, but `rooda-config.yml` has this section commented out with example presets.

**Current behavior:** Config has commented-out `ai_tools` section with examples.

**Spec requirement:** Config should document the `ai_tools` section structure and usage.

**Implementation status:** Config has examples but they're commented out. Users must uncomment to use.

**Action:** Config is correct - examples are intentionally commented out. No gap.

---

### 6. CLI Interface - Missing fuzzy procedure name matching
**Gap:** Implementation has fuzzy matching for unknown procedures (lines 230-260), but spec `cli-interface.md` doesn't document this feature.

**Current behavior:** When procedure not found, suggests closest match based on character overlap.

**Spec requirement:** Spec only says "Error 'Procedure X not found in config', exit 1"

**Action:** Spec is incomplete - implementation has better UX than spec describes. Update spec to document fuzzy matching.

---

### 7. User Documentation - Missing docs/beads.md content verification
**Gap:** Spec `user-documentation.md` lists `docs/beads.md` as existing file, but need to verify it exists and has correct content.

**Current behavior:** File exists (confirmed by glob).

**Action:** Verify `docs/beads.md` content matches spec requirements.

---

## Low-Priority Gaps (Nice-to-Have, Edge Cases)

### 8. External Dependencies - No automated installer
**Gap:** Spec `external-dependencies.md` "Areas for Improvement" mentions automated installer, but not implemented.

**Current behavior:** Users manually install dependencies.

**Spec requirement:** This is "Areas for Improvement", not acceptance criteria.

**Action:** Not a gap - this is a future enhancement.

---

### 9. AI CLI Integration - No timeout mechanism
**Gap:** Spec `ai-cli-integration.md` "Known Issues" mentions no timeout for AI CLI invocation, but not implemented.

**Current behavior:** Script waits indefinitely if AI CLI hangs.

**Spec requirement:** This is "Known Issues" / "Areas for Improvement", not acceptance criteria.

**Action:** Not a gap - this is a known limitation, not a missing feature.

---

### 10. Iteration Loop - No dry-run mode
**Gap:** Spec `iteration-loop.md` "Areas for Improvement" mentions --dry-run flag, but not implemented.

**Current behavior:** No dry-run mode.

**Spec requirement:** This is "Areas for Improvement", not acceptance criteria.

**Action:** Not a gap - this is a future enhancement.

---

## Spec Drift (Specified Differently Than Implemented)

### 11. AI CLI Integration - Precedence documentation inconsistency
**Gap:** Spec `ai-cli-integration.md` describes four-tier precedence including config `ai_cli_command`, but implementation only has three tiers (no config field).

**Current behavior:** Three-tier precedence: `--ai-cmd` > `--ai-cmd-preset` > `$ROODA_AI_CMD` > default

**Spec requirement:** Four-tier precedence: `--ai-cmd` > `--ai-cmd-preset` > `$ROODA_AI_CMD` > config `ai_cli_command` > default

**Action:** This is the same as Gap #1. Implement config `ai_cli_command` field.

---

### 12. Configuration Schema - ai_tools validation not implemented
**Gap:** Spec `configuration-schema.md` mentions "ai_tools presets validated as string type" and "Unknown presets return helpful error messages", but implementation doesn't validate string type.

**Current behavior:** Implementation checks if preset is null or empty, returns helpful error for unknown presets.

**Spec requirement:** Validate that custom presets are strings.

**Implementation needed:**
- Add type validation in `resolve_ai_tool_preset` function
- Check that custom preset value is a string (not array, object, number, boolean)
- Return error if wrong type

**Acceptance:** Config with `ai_tools: { bad: 123 }` returns error "Preset 'bad' must be a string"

---

## Undocumented Implementation (Implemented But Not Specified)

### 13. CLI Interface - Fuzzy procedure name matching
**Gap:** Implementation has fuzzy matching (lines 230-260) but spec doesn't document it.

**Action:** Update spec `cli-interface.md` to document fuzzy matching behavior.

---

### 14. CLI Interface - Detailed git push error handling
**Gap:** Implementation has detailed git push error parsing (lines 580-600) but spec doesn't document it.

**Current behavior:** Parses git push errors and provides specific guidance (authentication, network, merge conflicts).

**Spec requirement:** Spec `iteration-loop.md` only says "Attempts to create remote branch, continues loop"

**Action:** Update spec `iteration-loop.md` to document detailed error handling.

---

### 15. Configuration Schema - Config validation function
**Gap:** Implementation has comprehensive `validate_config` function (lines 195-298) but spec doesn't document validation behavior.

**Current behavior:** Validates YAML parseability, procedures key exists, procedure exists, required OODA fields present and non-empty.

**Spec requirement:** Spec only mentions "Missing procedures return clear error messages"

**Action:** Update spec `configuration-schema.md` to document validation behavior.

---

## Summary

**Critical gaps requiring implementation:**
1. AI CLI Integration - Add config `ai_cli_command` field support
2. Configuration Schema - Add ai_tools preset type validation

**Spec updates required (implementation exceeds spec):**
3. CLI Interface - Document fuzzy procedure name matching
4. Iteration Loop - Document detailed git push error handling  
5. Configuration Schema - Document config validation behavior

**False positives (no action needed):**
- External Dependencies bd check (project-specific, not framework-required)
- Version output format (spec is correct)
- Future enhancements listed in "Areas for Improvement" (not gaps)
- Commented-out ai_tools examples (intentional)
