# External Dependencies

## Job to be Done
Enable users to install and verify all required external tools before running ralph-wiggum-ooda procedures, preventing runtime failures due to missing dependencies.

## Dependency Philosophy

**Minimal Required, Maximum Flexibility**

Only `yq` is truly required for ralph-wiggum-ooda to function. Everything else is configurable or project-specific:

- **yq (required)** - Core dependency for parsing rooda-config.yml. No substitution possible.
- **kiro-cli (default, configurable)** - Default AI CLI tool, but can be substituted with any AI CLI that supports piping prompts and tool access (claude-cli, aider, cursor, etc.). The framework pipes prompts to stdin and expects the AI to have filesystem and command execution capabilities.
- **bd/beads (project-specific, optional)** - Default work tracking system used in examples and AGENTS.md templates, but entirely project-specific. Projects can use GitHub Issues, file-based tracking, Jira, Linear, or any other system. AGENTS.md documents whatever system the project uses.

This philosophy enables ralph-wiggum-ooda to adapt to existing project workflows rather than forcing specific tools.

## Activities
1. Check for required dependencies at script startup
2. Report missing dependencies with installation instructions
3. Verify dependency versions meet minimum requirements
4. Document all external tools and their purposes

## Acceptance Criteria
- [x] Core dependency (yq) documented with installation instructions
- [x] Default configurable dependencies (kiro-cli) documented with substitution guidance
- [x] Project-specific optional dependencies (bd) documented with alternatives
- [x] Dependency philosophy clearly stated
- [x] Version requirements specified where applicable
- [x] Dependency checking implemented in rooda.sh for critical tools

## Data Structures

### Dependency Information
```json
{
  "name": "string",
  "purpose": "string",
  "required": "boolean",
  "minimum_version": "string|null",
  "check_command": "string",
  "install_macos": "string",
  "install_linux": "string"
}
```

**Fields:**
- `name` - Tool name (e.g., "yq", "kiro-cli")
- `purpose` - What the tool is used for in the framework
- `required` - Whether the tool is mandatory for operation
- `minimum_version` - Minimum version required (null if any version works)
- `check_command` - Command to verify installation (e.g., "yq --version")
- `install_macos` - Installation command for macOS
- `install_linux` - Installation command for Linux

## Dependencies

### Required Dependencies

#### yq (YAML processor)
**Purpose:** Parse rooda-config.yml to extract procedure configurations

**Minimum Version:** v4.0.0 (uses v4 syntax)

**Check:**
```bash
yq --version
```

**Installation:**
- macOS: `brew install yq`
- Linux: `wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/local/bin/yq && chmod +x /usr/local/bin/yq`

**Verification:** rooda.sh checks for yq at startup (lines 15-19)

### Default Dependencies (Configurable)

#### kiro-cli (AI CLI tool)
**Purpose:** Execute OODA prompts through AI chat interface with tool access

**Default:** kiro-cli is the default AI CLI used in documentation and examples

**Substitution:** Can be replaced with any AI CLI that supports:
- Reading prompts from stdin (piped input)
- Filesystem access (read/write files)
- Command execution (run bash commands)
- Non-interactive mode

**Alternatives:** claude-cli, aider, cursor, or custom AI CLI wrappers

**Minimum Version:** 1.0.0 (requires --no-interactive and --trust-all-tools flags)

**Check:**
```bash
kiro-cli --version
```

**Installation:**
- macOS/Linux: Follow AWS Kiro CLI installation instructions at https://docs.aws.amazon.com/kiro/

**Verification:** Script pipes prompt to `kiro-cli chat --no-interactive --trust-all-tools`

**Configuration:** To use a different AI CLI, modify the pipe command in rooda.sh (line 161)

### Project-Specific Dependencies (Optional)

#### bd (beads work tracking)
**Purpose:** Query and update work tracking system for task management

**Default:** bd (beads) is used in examples and AGENTS.md templates

**Project-Specific:** Work tracking is entirely project-specific. Projects can use:
- GitHub Issues
- File-based tracking (PLAN.md, TODO.md)
- Jira, Linear, Asana
- Any system that can be queried and updated via CLI or API

**AGENTS.md Role:** The "Work Tracking System" section in AGENTS.md documents whatever system the project uses, including query/update commands.

**Minimum Version:** 0.1.0 (requires --json flag support)

**Check:**
```bash
bd --version
```

**Installation:**
- macOS/Linux: `cargo install beads-cli` (requires Rust toolchain)
- Alternative: Download binary from https://github.com/jomadu/beads/releases

**Verification:** AGENTS.md documents `bd ready --json` command

**Note:** If not using bd, update AGENTS.md with your project's work tracking commands

### Optional Dependencies

#### shellcheck (bash linter)
**Purpose:** Lint rooda.sh for bash script quality

**Minimum Version:** None specified

**Check:**
```bash
shellcheck --version
```

**Installation:**
- macOS: `brew install shellcheck`
- Linux: `apt-get install shellcheck` or `yum install shellcheck`

**Verification:** AGENTS.md notes this is optional but recommended

#### git (version control)
**Purpose:** Commit changes after successful iterations

**Minimum Version:** 2.0.0 (standard git operations)

**Check:**
```bash
git --version
```

**Installation:**
- macOS: Pre-installed or `brew install git`
- Linux: `apt-get install git` or `yum install git`

**Verification:** Script uses git for commits (not enforced at startup)

## Algorithm

**Dependency Check Flow:**
```
function check_dependencies():
    for each required_dependency:
        if not command_exists(dependency):
            print error_message(dependency)
            print installation_instructions(dependency)
            exit 1
        
        if version_required(dependency):
            version = extract_version(dependency)
            if version < minimum_version:
                print error_message(dependency, version, minimum_version)
                print upgrade_instructions(dependency)
                exit 1
    
    for each optional_dependency:
        if not command_exists(dependency):
            print warning_message(dependency)
            print installation_instructions(dependency)
    
    continue execution
```

**Current Implementation:**
- rooda.sh checks for yq at startup (required)
- rooda.sh validates yq version >= 4.0.0 (lines 112-120)
- rooda.sh checks for kiro-cli at startup (default, but can be modified)
- rooda.sh checks for bd at startup (project-specific, can be removed if using different work tracking)

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| yq not installed | Script exits with error and installation instructions |
| yq v3 installed (incompatible) | Script exits with error: "yq version 4.0.0 or higher required (found X.X.X)" and upgrade instructions |
| kiro-cli not installed | Script exits with error (default behavior, can be modified for other AI CLIs) |
| bd not installed | Script exits with error (can be removed if using different work tracking) |
| shellcheck not installed | Lint command fails but documented as optional |
| git not installed | Commit operations fail (not critical for testing) |

## Implementation Mapping

**Source files:**
- `src/rooda.sh` (lines 15-19) - yq dependency check
- `src/rooda.sh` (lines 112-120) - yq version validation (requires v4.0.0+)
- `src/rooda.sh` (line 161) - kiro-cli invocation
- `AGENTS.md` - Documents bd commands and shellcheck usage

**Related specs:**
- `cli-interface.md` - Defines command-line interface that depends on yq
- `ai-cli-integration.md` - Defines kiro-cli integration
- `iteration-loop.md` - Lists dependencies briefly

## Examples

### Example 1: Successful Dependency Check

**Input:**
```bash
./rooda.sh bootstrap
```

**Expected Output:**
```
(no output if yq is installed)
(script proceeds to execute bootstrap procedure)
```

**Verification:**
- Script does not exit with error
- yq command succeeds in parsing rooda-config.yml

### Example 2: Missing yq Dependency

**Input:**
```bash
./rooda.sh bootstrap
# (with yq not installed)
```

**Expected Output:**
```
Error: yq is required for YAML parsing
Install with: brew install yq
```

**Verification:**
- Script exits with status 1
- Error message provides installation instructions

### Example 3: Version Check (Implemented)

**Input:**
```bash
yq --version
# yq (https://github.com/mikefarah/yq/) version v3.4.1
./rooda.sh bootstrap
```

**Expected Output:**
```
Error: yq version 4.0.0 or higher required (found 3.4.1)
Upgrade with: brew upgrade yq
```

**Verification:**
- Script exits with status 1
- Error message shows detected version and required minimum
- Provides upgrade instructions

## Notes

**Design Decision:** Only yq is truly required. kiro-cli and bd are checked at startup by default but can be modified or removed based on project needs.

**Rationale for Minimum Versions:**
- yq v4.0.0: Script uses v4 syntax (`.procedures.$PROCEDURE.observe`)
- kiro-cli 1.0.0: Requires --no-interactive and --trust-all-tools flags (if using kiro-cli)
- bd 0.1.0: Requires --json flag for structured output (if using bd)

**Consumer vs Framework:**
- Consumers need: yq (required), AI CLI of choice (kiro-cli or alternative), work tracking system of choice (bd or alternative)
- Framework development needs: shellcheck (optional, for bash linting)
- Both need: git (optional, for version control)

**Substitution Examples:**
- Replace kiro-cli with claude-cli: Modify pipe command in rooda.sh
- Replace bd with GitHub Issues: Update AGENTS.md with `gh issue` commands
- Replace bd with file-based: Update AGENTS.md to read/write PLAN.md or TODO.md

## Known Issues

1. **Dependency checks assume defaults:** Script checks for kiro-cli and bd at startup, but these are configurable. Projects using alternatives need to modify rooda.sh dependency checks.

2. **Platform-specific instructions:** Installation commands assume macOS (brew) or Linux package managers. Windows/WSL users need different instructions.

3. **No automated installer:** Users must manually install all dependencies. No bootstrap script to automate setup.

## Areas for Improvement

1. **Configurable dependency checks:** Allow projects to disable kiro-cli/bd checks via config if using alternatives
2. **Automated installer:** Provide setup script that installs all dependencies
3. **Platform detection:** Detect OS and provide appropriate installation commands
4. **Dependency matrix:** Document tested version combinations
5. **Docker image:** Provide pre-configured container with all dependencies
6. **AI CLI abstraction:** Create adapter layer so different AI CLIs can be swapped without modifying rooda.sh
