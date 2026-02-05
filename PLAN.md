# Draft Plan: Flexible AI CLI Configuration

## Priority 1: Add resolve_ai_tool_preset() function to src/rooda.sh
- Add function after show_help() and before argument parsing
- Implement hardcoded presets: kiro-cli, claude, aider
- Query custom presets from config: `.ai_tools.$preset`
- Return error with helpful message for unknown presets
- List available hardcoded presets and suggest checking config

**Dependencies:** None

**Acceptance Criteria:**
- Function returns correct command for hardcoded presets
- Function queries config for custom presets
- Function returns error for unknown presets with helpful message

## Priority 2: Add $ROODA_AI_CLI environment variable support
- After default AI_CLI_COMMAND initialization, check for $ROODA_AI_CLI
- If set, override AI_CLI_COMMAND with environment variable value
- Occurs before argument parsing (so flags can still override)

**Dependencies:** None

**Acceptance Criteria:**
- $ROODA_AI_CLI overrides default when set
- Environment variable checked before argument parsing

## Priority 3: Add --ai-tool flag parsing
- Add case for --ai-tool in argument parsing while loop
- Store preset name in AI_TOOL_PRESET variable
- Shift 2 arguments after parsing

**Dependencies:** Priority 1 (resolve_ai_tool_preset function)

**Acceptance Criteria:**
- --ai-tool flag parsed correctly
- Preset name stored in variable

## Priority 4: Implement precedence resolution for AI CLI command
- After argument parsing, before procedure config loading
- If AI_TOOL_PRESET set, resolve via resolve_ai_tool_preset()
- Resolved preset overrides $ROODA_AI_CLI and default
- --ai-cli flag already has highest priority (set during parsing)

**Dependencies:** Priority 1, Priority 2, Priority 3

**Acceptance Criteria:**
- Precedence order: --ai-cli > --ai-tool > $ROODA_AI_CLI > default
- Preset resolution happens at correct point in script flow

## Priority 5: Remove per-procedure ai_cli_command from config loading
- Remove yq query for `.procedures.$PROCEDURE.ai_cli_command`
- Remove conditional that sets AI_CLI_COMMAND from procedure config
- Root-level ai_cli_command already removed per task requirements

**Dependencies:** None (breaking change, safe to do independently)

**Acceptance Criteria:**
- Per-procedure ai_cli_command no longer queried
- Root-level ai_cli_command remains unsupported

## Priority 6: Update show_help() with new flags
- Add --ai-tool <preset> to usage examples
- Document precedence order
- List hardcoded presets (kiro-cli, claude, aider)
- Mention custom presets in config

**Dependencies:** Priority 1, Priority 3

**Acceptance Criteria:**
- Help text includes --ai-tool flag
- Precedence documented clearly
- Examples show usage patterns

## Priority 7: Update src/rooda-config.yml with ai_tools section
- Add commented example ai_tools section at root level
- Include example presets: fast, thorough
- Add comments explaining custom preset usage

**Dependencies:** None

**Acceptance Criteria:**
- Config has ai_tools section (commented)
- Examples are clear and actionable

## Priority 8: Update specs/ai-cli-integration.md
- Replace "Configuration" section with new precedence system
- Document --ai-tool flag and presets
- Document $ROODA_AI_CLI environment variable
- Update examples to show all configuration methods
- Remove per-procedure ai_cli_command references
- Add edge cases for preset resolution errors

**Dependencies:** Priority 1-7 (implementation complete)

**Acceptance Criteria:**
- Spec accurately reflects implemented behavior
- All configuration methods documented
- Examples match implementation
- Edge cases cover preset resolution

## Priority 9: Update specs/configuration-schema.md
- Add ai_tools section documentation
- Document structure: map of preset names to commands
- Add examples of custom presets
- Remove per-procedure ai_cli_command field docs
- Update root-level schema to show ai_tools is optional

**Dependencies:** Priority 7 (config updated)

**Acceptance Criteria:**
- ai_tools section fully documented
- Schema examples match config examples
- Optional fields clearly marked

## Priority 10: Update specs/cli-interface.md
- Add --ai-tool flag to arguments table
- Document precedence order in algorithm section
- Add examples showing --ai-tool usage
- Update edge cases for unknown presets
- Add example showing preset resolution error message

**Dependencies:** Priority 1-7 (implementation complete)

**Acceptance Criteria:**
- --ai-tool flag documented in arguments
- Precedence clear in algorithm
- Examples show preset usage
- Edge cases match implementation

## Priority 11: Update README.md with configuration section
- Add "Configuring Your AI CLI" section after "Installation"
- Document all four configuration methods (--ai-cli, --ai-tool, $ROODA_AI_CLI, default)
- Show team workflow examples
- Include troubleshooting for unknown presets

**Dependencies:** Priority 1-7 (implementation complete)

**Acceptance Criteria:**
- Configuration section added to README
- User-friendly examples for common scenarios
- Clear guidance for team usage
