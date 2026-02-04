# Draft Plan: Spec to Implementation Gap Analysis

## Priority 1: Implement --ai-cli Flag and ai_cli_command Config Support

**Gap:** Specs specify configurable AI CLI command via --ai-cli flag and ai_cli_command config field, but implementation hardcodes kiro-cli invocation.

**Specified in:**
- `cli-interface.md` - Documents --ai-cli flag with precedence rules
- `ai-cli-integration.md` - Documents ai_cli_command config field and resolution algorithm
- `configuration-schema.md` - Documents ai_cli_command as optional root-level field

**Current State:**
- Line 443: Hardcoded `kiro-cli chat --no-interactive --trust-all-tools`
- No --ai-cli flag in argument parser
- No AI_CLI_COMMAND variable
- No query for .ai_cli_command from config
- Lines 72-86: Hardcoded kiro-cli dependency check

**Implementation:**
- Add AI_CLI_COMMAND variable initialization with default
- Add --ai-cli flag to argument parser
- Query .ai_cli_command from config when procedure specified
- Implement precedence: --ai-cli flag > config ai_cli_command > default
- Replace hardcoded invocation at line 443 with $AI_CLI_COMMAND
- Remove hardcoded kiro-cli dependency check (configurable tool)

**Acceptance:**
- `./rooda.sh build --ai-cli "claude-cli"` uses claude-cli
- Config with `ai_cli_command: "aider"` uses aider when no flag
- No flag + no config = default kiro-cli
- Script runs without kiro-cli installed if using different AI CLI

---

## Priority 2: Remove Hardcoded bd Dependency Check

**Gap:** `external-dependencies.md` specifies bd is project-specific and optional, but implementation requires it at startup.

**Specified in:**
- `external-dependencies.md` - Documents bd as project-specific optional dependency

**Current State:**
- Lines 88-97: Hardcoded bd dependency check with exit on failure
- Lines 104-111: bd version validation

**Implementation:**
- Remove bd dependency check (project-specific work tracking)
- Remove bd version validation
- Keep yq dependency check (required)
- Update AGENTS.md to clarify only yq is framework-required

**Acceptance:**
- Script runs without bd installed
- yq check remains and exits with clear error if missing
- AGENTS.md documents dependency philosophy

---

## Priority 3: Document --verbose and --quiet Flags

**Gap:** Implementation has --verbose and --quiet flags but specs don't document them.

**Current State:**
- Lines 283-289: --verbose and --quiet flag parsing
- VERBOSE variable (0=default, 1=verbose, -1=quiet)
- Line 428-434: Verbose mode shows full prompt
- Lines 383-385, 451: Conditional output based on VERBOSE

**Implementation:**
- Update `cli-interface.md` to document --verbose and --quiet
- Add examples showing verbose mode output
- Add examples showing quiet mode suppression
- Update README.md with usage examples

**Acceptance:**
- cli-interface.md documents both flags with examples
- README.md shows common usage patterns
- Help text already includes these (verified in show_help)

---

## Priority 4: Update cli-interface.md for --version Flag

**Gap:** Implementation has --version flag but cli-interface.md lists it as "Areas for Improvement".

**Current State:**
- Line 2: VERSION="0.1.0" defined
- Lines 237-240, 251-254: --version flag implemented
- cli-interface.md "Areas for Improvement" incorrectly states "No --version flag"

**Implementation:**
- Update cli-interface.md acceptance criteria to mark --version as implemented
- Add example showing version output
- Remove from "Areas for Improvement" section
- Document in data structures and algorithm sections

**Acceptance:**
- cli-interface.md accurately reflects implementation
- Example shows `rooda.sh version 0.1.0` output
- "Areas for Improvement" section updated

---

## Priority 5: Add Short Flag Support

**Gap:** Help text shows short flags (-o, -r, -d, -a, -m, -c) but argument parser only accepts long flags.

**Current State:**
- show_help documents short flags (lines 15-49)
- Argument parser only handles long flags (lines 250-289)
- No short flag cases in switch statement

**Implementation:**
- Add short flag cases to argument parser
- Map -o to --observe, -r to --orient, -d to --decide, -a to --act
- Map -m to --max-iterations, -c to --config, -h to --help
- Test all short flags work identically to long flags

**Acceptance:**
- `-o file.md` works identically to `--observe file.md`
- `-m 5` works identically to `--max-iterations 5`
- `-h` works identically to `--help`
- All documented short flags functional

---

## Priority 6: Implement --list-procedures Flag

**Gap:** Config has display/summary/description fields but no command to list available procedures.

**Specified in:**
- `configuration-schema.md` - Documents display/summary/description as "reserved for future help text generation"

**Current State:**
- Config has display/summary/description for all 9 procedures
- show_help doesn't use these fields
- No --list-procedures flag

**Implementation:**
- Add --list-procedures flag to argument parser
- Query config for all procedure names using yq
- Display procedure name, display field, and summary
- Format as table or list
- Consider --help <procedure> for detailed description

**Acceptance:**
- `./rooda.sh --list-procedures` shows all 9 procedures
- Output includes name and summary from config
- Graceful handling if display/summary missing
- Clear, readable format

---

## Priority 7: Add Cross-Document Link Validation Script

**Gap:** `user-documentation.md` acceptance criteria requires working links but no validation exists.

**Specified in:**
- `user-documentation.md` - Acceptance criteria: "All cross-document links work correctly"
- AGENTS.md quality criteria - "All cross-document links work correctly (PASS/FAIL)"

**Current State:**
- No automated link checking
- scripts/audit-links.sh exists but may not be comprehensive
- Quality criteria includes link validation but no tool

**Implementation:**
- Verify scripts/audit-links.sh functionality
- Ensure it checks internal links (relative paths)
- Ensure it checks external links (with timeout)
- Add to quality criteria verification process
- Document usage in AGENTS.md

**Acceptance:**
- Script validates all markdown links in docs/ and specs/
- Broken internal links reported with file and line number
- Broken external links reported
- Can be run as part of quality assessment

---

## Priority 8: Add Prompt File Structure Validation Script

**Gap:** `component-authoring.md` documents prompt structure but no validation ensures compliance.

**Specified in:**
- `component-authoring.md` - Documents phase headers, step codes, prose structure
- AGENTS.md quality criteria - "All procedures in config have corresponding component files (PASS/FAIL)"

**Current State:**
- No automated validation of prompt file structure
- 25 prompt files in src/prompts/
- No linting for step code consistency

**Implementation:**
- Create script to validate prompt file structure
- Check for phase header: `# [Phase]: [Purpose]`
- Check for step headers: `## [Code]: [Name]`
- Validate step codes match phase (O1-O15, R1-R22, D1-D15, A1-A9)
- Add to quality criteria for implementation

**Acceptance:**
- Script validates all 25 prompt files
- Reports missing phase headers
- Reports invalid step codes
- Reports step codes that don't match phase
- Can be run as part of quality assessment

---

## Priority 9: Verify act_build.md Substep Numbering

**Gap:** `component-authoring.md` documents substep numbering pattern but need to verify act_build.md follows it.

**Specified in:**
- `component-authoring.md` - Documents substep numbering (A3, A3.5, A3.6, A4) with example from act_build.md

**Current State:**
- component-authoring.md shows act_build.md example with substeps
- Need to verify actual file matches documented pattern

**Implementation:**
- Read src/prompts/act_build.md
- Verify substep numbering follows pattern
- Ensure conditional steps clearly marked
- Confirm backpressure steps properly labeled
- Update if needed to match documented pattern

**Acceptance:**
- act_build.md uses substep numbering correctly
- Conditional steps have "If X Modified" markers
- Critical warnings use **bold** emphasis
- Matches example in component-authoring.md

---

## Priority 10: Update external-dependencies.md Implementation Mapping

**Gap:** Spec documents dependency checks but implementation has evolved (version validation added).

**Specified in:**
- `external-dependencies.md` - Documents dependency checking algorithm

**Current State:**
- Lines 99-113: Version validation for yq, kiro-cli, bd (not documented in spec)
- Spec "Implementation Mapping" only mentions lines 15-19 for yq check
- Spec "Known Issues" mentions "No version validation" but it's implemented

**Implementation:**
- Update external-dependencies.md "Implementation Mapping" section
- Document version validation logic (lines 99-113)
- Remove "No version validation" from "Known Issues"
- Update algorithm section to include version checks

**Acceptance:**
- external-dependencies.md accurately reflects implementation
- Version validation documented in algorithm
- Implementation mapping includes all relevant line numbers
- Known issues section updated
