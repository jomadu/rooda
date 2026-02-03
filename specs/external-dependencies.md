# External Dependencies

## Job to be Done
Enable users to install and verify all required external tools before running ralph-wiggum-ooda procedures, preventing runtime failures due to missing dependencies.

## Activities
1. Check for required dependencies at script startup
2. Report missing dependencies with installation instructions
3. Verify dependency versions meet minimum requirements
4. Document all external tools and their purposes

## Acceptance Criteria
- [x] All external dependencies documented (yq, kiro-cli, bd)
- [x] Version requirements specified where applicable
- [x] Installation instructions provided per platform
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

#### kiro-cli (AI CLI tool)
**Purpose:** Execute OODA prompts through AI chat interface with tool access

**Minimum Version:** 1.0.0 (requires --no-interactive and --trust-all-tools flags)

**Check:**
```bash
kiro-cli --version
```

**Installation:**
- macOS/Linux: Follow AWS Kiro CLI installation instructions at https://docs.aws.amazon.com/kiro/

**Verification:** Script pipes prompt to `kiro-cli chat --no-interactive --trust-all-tools`

#### bd (beads work tracking)
**Purpose:** Query and update work tracking system for task management

**Minimum Version:** 0.1.0 (requires --json flag support)

**Check:**
```bash
bd --version
```

**Installation:**
- macOS/Linux: `cargo install beads-cli` (requires Rust toolchain)
- Alternative: Download binary from https://github.com/jomadu/beads/releases

**Verification:** AGENTS.md documents `bd ready --json` command

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
    
    for each optional_dependency:
        if not command_exists(dependency):
            print warning_message(dependency)
            print installation_instructions(dependency)
    
    continue execution
```

**Current Implementation:**
- rooda.sh checks for yq at startup (required)
- No checks for kiro-cli or bd (assumed present)
- No version validation

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| yq not installed | Script exits with error and installation instructions |
| yq v3 installed (incompatible) | Script may fail with cryptic YAML parsing errors |
| kiro-cli not installed | Script fails when piping prompt (no early detection) |
| bd not installed | AGENTS.md commands fail (no early detection) |
| shellcheck not installed | Lint command fails but documented as optional |
| git not installed | Commit operations fail (not critical for testing) |

## Implementation Mapping

**Source files:**
- `src/rooda.sh` (lines 15-19) - yq dependency check
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

### Example 3: Version Check (Not Implemented)

**Input:**
```bash
yq --version
# yq (https://github.com/mikefarah/yq/) version v3.4.1
./rooda.sh bootstrap
```

**Expected Behavior:**
Script should detect incompatible yq v3 and warn user, but currently does not validate version.

## Notes

**Design Decision:** Only yq is checked at startup because it's required immediately for config parsing. kiro-cli and bd are used later in execution, so their absence causes runtime failures rather than startup failures.

**Rationale for Minimum Versions:**
- yq v4.0.0: Script uses v4 syntax (`.procedures.$PROCEDURE.observe`)
- kiro-cli 1.0.0: Requires --no-interactive and --trust-all-tools flags
- bd 0.1.0: Requires --json flag for structured output

**Consumer vs Framework:**
- Consumers need: yq, kiro-cli, bd (required for all procedures)
- Framework development needs: shellcheck (optional, for bash linting)
- Both need: git (optional, for version control)

## Known Issues

1. **No version validation:** Script checks if yq exists but not if it's v4+. Users with yq v3 will get cryptic YAML parsing errors.

2. **Late failure for kiro-cli/bd:** Script doesn't check for kiro-cli or bd at startup, so users discover missing tools only when procedures execute.

3. **Platform-specific instructions:** Installation commands assume macOS (brew) or Linux package managers. Windows/WSL users need different instructions.

4. **No automated installer:** Users must manually install all dependencies. No bootstrap script to automate setup.

## Areas for Improvement

1. **Add version validation:** Check yq version and warn if < v4.0.0
2. **Early dependency checks:** Validate kiro-cli and bd presence at startup
3. **Automated installer:** Provide setup script that installs all dependencies
4. **Platform detection:** Detect OS and provide appropriate installation commands
5. **Dependency matrix:** Document tested version combinations
6. **Docker image:** Provide pre-configured container with all dependencies
