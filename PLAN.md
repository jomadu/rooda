# Draft Plan: Spec to Implementation Gap Analysis

## Critical Gaps (Missing Core Features)

### 1. AI CLI Configuration Support (ai_cli_command)
**Specified:** `ai-cli-integration.md` and `configuration-schema.md` define `ai_cli_command` field in rooda-config.yml with three-tier precedence (--ai-cli flag > ai_cli_command config > default)

**Implementation Status:** NOT IMPLEMENTED
- rooda-config.yml has no `ai_cli_command` field
- rooda.sh has no --ai-cli flag parsing
- rooda.sh hardcodes `kiro-cli chat --no-interactive --trust-all-tools` at line 436
- No precedence resolution logic exists

**Acceptance Criteria:**
- Add `ai_cli_command` field to rooda-config.yml root level
- Add --ai-cli flag parsing in argument loop
- Implement three-tier precedence: flag > config > default
- Query config for ai_cli_command when procedure specified
- Use resolved AI_CLI_COMMAND variable in iteration loop

**Dependencies:** None

---

### 2. Dependency Philosophy Documentation
**Specified:** `external-dependencies.md` defines dependency philosophy: only yq required, kiro-cli configurable, bd project-specific

**Implementation Status:** CONTRADICTS SPEC
- rooda.sh lines 60-88 check for kiro-cli and bd as REQUIRED dependencies
- Script exits if kiro-cli or bd not installed
- No documentation that these are optional/configurable
- Contradicts spec's "Minimal Required, Maximum Flexibility" philosophy

**Acceptance Criteria:**
- Remove or make optional the kiro-cli dependency check (since it's configurable via ai_cli_command)
- Remove or make optional the bd dependency check (since work tracking is project-specific)
- Update dependency checks to reflect philosophy: only yq is truly required
- Add comments explaining why kiro-cli/bd checks are optional

**Dependencies:** Should be implemented after task #1 (AI CLI configuration)

---

### 3. Version Validation for yq
**Specified:** `external-dependencies.md` requires yq v4.0.0+ with validation

**Implementation Status:** IMPLEMENTED
- Lines 90-96 validate yq version >= 4.0.0
- Provides clear error message with upgrade instructions

**Gap:** None - this is correctly implemented

---

### 4. Version Validation for kiro-cli
**Specified:** `external-dependencies.md` requires kiro-cli v1.0.0+ with validation

**Implementation Status:** IMPLEMENTED
- Lines 98-103 validate kiro-cli version >= 1.0.0
- Provides clear error message with upgrade instructions

**Gap:** None - this is correctly implemented (but should be made optional per task #2)

---

### 5. Version Validation for bd
**Specified:** `external-dependencies.md` requires bd v0.1.0+ with validation

**Implementation Status:** IMPLEMENTED
- Lines 105-111 validate bd version >= 0.1.0
- Provides clear error message with upgrade instructions

**Gap:** None - this is correctly implemented (but should be made optional per task #2)

---

## High-Impact Gaps (Frequently Used Functionality)

### 6. Help Flag Support
**Specified:** `cli-interface.md` notes "No --help or -h flag support" as area for improvement

**Implementation Status:** IMPLEMENTED
- Lines 233-237 check for --help/-h flags
- Lines 251-254 handle --help/-h in argument loop
- show_help() function exists at lines 15-49

**Gap:** None - this was implemented despite spec noting it as missing

---

### 7. Version Flag Support
**Specified:** `cli-interface.md` notes "No --version flag" as area for improvement

**Implementation Status:** IMPLEMENTED
- Lines 239-242 check for --version flag
- Lines 247-250 handle --version in argument loop
- VERSION variable defined (need to verify)

**Gap:** None - this was implemented despite spec noting it as missing

---

### 8. Short Flag Alternatives
**Specified:** `cli-interface.md` notes "No short flag alternatives" as area for improvement

**Implementation Status:** PARTIALLY IMPLEMENTED
- Lines 255-277 show short flags: -c, -o, -r, -d, -a, -m, -h
- Short flags exist for all major options

**Gap:** None - this was implemented despite spec noting it as missing

---

### 9. Verbose/Quiet Modes
**Specified:** `cli-interface.md` notes "No control over output verbosity" as area for improvement

**Implementation Status:** IMPLEMENTED
- Line 228 initializes VERBOSE variable
- Lines 278-284 handle --verbose and --quiet flags
- Lines 383-393 use VERBOSE for conditional output
- Lines 426-432 show full prompt in verbose mode

**Gap:** None - this was implemented despite spec noting it as missing

---

## Documentation Gaps

### 10. Update CLI Interface Spec
**Specified:** `cli-interface.md` lists several features as "Areas for Improvement" that are actually implemented

**Implementation Status:** SPEC OUTDATED
- Help flag is implemented but spec says it's missing
- Version flag is implemented but spec says it's missing
- Short flags are implemented but spec says they're missing
- Verbose/quiet modes are implemented but spec says they're missing

**Acceptance Criteria:**
- Update cli-interface.md to reflect implemented features
- Move implemented features from "Areas for Improvement" to acceptance criteria
- Mark acceptance criteria as [x] completed
- Add examples for new flags (--help, --version, --verbose, --quiet)

**Dependencies:** None

---

### 11. Update AI CLI Integration Spec Examples
**Specified:** `ai-cli-integration.md` shows examples with ai_cli_command configuration

**Implementation Status:** SPEC AHEAD OF IMPLEMENTATION
- Examples show ai_cli_command in config (not yet implemented)
- Examples show --ai-cli flag usage (not yet implemented)
- Examples show precedence resolution (not yet implemented)

**Acceptance Criteria:**
- Mark examples as "Not Yet Implemented" until task #1 is complete
- Or update examples to reflect current hardcoded behavior

**Dependencies:** Should be updated after task #1 (AI CLI configuration)

---

### 12. Update External Dependencies Spec
**Specified:** `external-dependencies.md` describes dependency philosophy but implementation contradicts it

**Implementation Status:** SPEC AHEAD OF IMPLEMENTATION
- Spec says kiro-cli is "default, configurable" but implementation requires it
- Spec says bd is "project-specific, optional" but implementation requires it
- Spec philosophy contradicts implementation behavior

**Acceptance Criteria:**
- Update spec to reflect current implementation (kiro-cli and bd are required)
- Or update implementation to match spec (make them optional)
- Ensure consistency between spec and implementation

**Dependencies:** Should be updated after task #2 (dependency philosophy)

---

## Low-Priority Gaps (Nice-to-Have)

### 13. Config Validation Function
**Specified:** `cli-interface.md` notes "Script doesn't validate config file structure before querying"

**Implementation Status:** IMPLEMENTED
- Lines 114-217 implement validate_config() function
- Validates YAML parseability
- Validates procedure exists
- Validates required OODA fields present
- Validates OODA files exist

**Gap:** None - this was implemented despite spec noting it as missing

---

### 14. Error Handling for AI CLI Failures
**Specified:** `iteration-loop.md` notes "No kiro-cli error handling" as known issue

**Implementation Status:** INTENTIONALLY NOT IMPLEMENTED
- Line 436 comment explains exit status intentionally ignored
- Design decision: allow loop to self-correct through empirical feedback

**Gap:** None - this is intentional per design

---

### 15. Git Push Error Handling
**Specified:** `iteration-loop.md` notes "Git push failures" as known issue

**Implementation Status:** IMPLEMENTED
- Lines 438-448 handle git push failures
- Attempts to create remote branch if push fails
- Provides clear error messages
- Allows user to continue or stop

**Gap:** None - this was implemented despite spec noting it as issue

---

## Summary

**Critical gaps requiring implementation:**
1. AI CLI configuration support (ai_cli_command field and --ai-cli flag)
2. Dependency philosophy alignment (make kiro-cli and bd optional)

**Documentation gaps requiring updates:**
3. Update cli-interface.md to reflect implemented features
4. Update ai-cli-integration.md examples after task #1
5. Update external-dependencies.md after task #2

**Total tasks:** 5 (2 implementation, 3 documentation)

**Observation:** Many features listed as "missing" or "areas for improvement" in specs have actually been implemented. The specs are outdated relative to the implementation. Priority should be: implement missing features (#1, #2), then update specs to match reality (#3, #4, #5).
