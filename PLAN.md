# Draft Plan: Rename AI CLI Flags for Clarity

## Priority-Ordered Tasks

### 1. Update specs/cli-interface.md for flag renaming
- Replace all `--ai-cli` with `--ai-cmd`
- Replace all `--ai-tool` with `--ai-cmd-preset`
- Replace all `$ROODA_AI_CLI` with `$ROODA_AI_CMD`
- Replace all internal variable references (AI_CLI_COMMAND → AI_CMD_COMMAND, AI_CLI_FLAG → AI_CMD_FLAG, AI_TOOL_PRESET → AI_CMD_PRESET)
- Update acceptance criteria
- Update data structures section
- Update algorithm pseudocode
- Update all 13 examples
- Update edge cases table

**Status:** Complete
**Acceptance:** All references to old flag names replaced with new names in cli-interface.md

### 2. Update specs/ai-cli-integration.md for flag renaming
- Replace all `--ai-cli` with `--ai-cmd`
- Replace all `--ai-tool` with `--ai-cmd-preset`
- Replace all `$ROODA_AI_CLI` with `$ROODA_AI_CMD`
- Replace all internal variable references (AI_CLI_COMMAND → AI_CMD_COMMAND, AI_CLI_FLAG → AI_CMD_FLAG, AI_TOOL_PRESET → AI_CMD_PRESET)
- Update configuration section
- Update acceptance criteria
- Update data structures section
- Update algorithm pseudocode
- Update all 9 examples
- Update edge cases table
- Update notes section

**Status:** Complete
**Acceptance:** All references to old flag names replaced with new names in ai-cli-integration.md

### 3. Update specs/configuration-schema.md for flag renaming
- Replace all `--ai-tool` with `--ai-cmd-preset` (5 occurrences)
- Update ai_tools section description to reference new flag name
- Update examples showing preset usage
- Update algorithm pseudocode for preset resolution

**Status:** Ready
**Acceptance:** All references to `--ai-tool` replaced with `--ai-cmd-preset` in configuration-schema.md

---

## Plan Status

All tasks identified. Ready for build iterations to implement spec updates.
