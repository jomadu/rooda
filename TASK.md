# Task: Refactor CLI to Use Cobra Framework

## Type
Refactor (implementation) + CLI improvement

## Priority
P2 (Medium) - Improves maintainability, extensibility, and user experience

## Context

The current v2 Go implementation uses manual flag parsing in `cmd/rooda/main.go`. As the CLI grows to support 16 procedures with multiple flags per procedure, manual parsing becomes error-prone and difficult to maintain.

The Cobra CLI framework (https://github.com/spf13/cobra) is the de facto standard for Go CLI applications, used by kubectl, Hugo, GitHub CLI, and many others. It provides:
- Automatic help generation
- Subcommand structure
- Flag inheritance (global flags + command-specific flags)
- Shell completion
- Better error messages

## Current State

**File:** `cmd/rooda/main.go`
- Manual flag parsing using `flag` package
- Flat command structure: `rooda <procedure> [flags]`
- Flags defined inline with parsing logic
- Help text manually constructed
- `--version` and `--list-procedures` as flags, not commands

**Current CLI:**
```bash
rooda <procedure-name> [flags]
rooda version
rooda list
```

**Behavior to preserve:**
- All existing flag semantics
- Exit codes (0=success, 1=aborted, 2=max-iters, 130=interrupted)
- Flag precedence (CLI > env > workspace > global > built-in)
- Dynamic procedure loading from config

## Desired State

**New CLI Structure:**
```bash
rooda run <procedure-name> [flags]      # Execute a procedure
rooda list [flags]                      # List available procedures
rooda info <procedure-name>             # Show procedure details
rooda version                           # Show version information
rooda help [command]                    # Show help
```

**Implementation:**
- Root command with global flags (--config, --verbose, --quiet, --log-level)
- `run` subcommand: executes procedures with all execution flags
- `list` subcommand: shows available procedures (replaces --list-procedures)
- `info` subcommand: shows procedure metadata (new feature)
- `version` subcommand: shows version (replaces --version flag)
- Procedures loaded dynamically from configuration (built-in + custom)
- Cobra handles flag parsing, validation, and help generation

## Requirements

### Must Have
1. Install Cobra dependency: `go get -u github.com/spf13/cobra@latest`
2. Create `run` subcommand that accepts procedure name as argument
3. Create `list` subcommand to show available procedures
4. Create `info` subcommand to show procedure details
5. Create `version` subcommand
6. Load procedures dynamically from merged config (built-in + global + workspace)
7. Validate procedure exists before execution
8. Preserve all existing flag behavior and semantics
9. All tests pass after refactor

### Should Have
1. Flag validation using Cobra's built-in validators
2. Consistent flag descriptions across commands
3. Shell completion support (bash, zsh, fish) with dynamic procedure names

### Nice to Have
1. `info` command shows OODA phase composition
2. `list --json` for machine-readable output
3. Custom usage templates for better formatting

## Acceptance Criteria

### Build & Test
- [ ] `go build -o bin/rooda ./cmd/rooda` succeeds
- [ ] `go test ./...` passes all tests

### New Commands
- [ ] `./bin/rooda run agents-sync --ai-cmd-alias kiro-cli --max-iterations 1` executes successfully
- [ ] `./bin/rooda list` shows all available procedures (built-in + custom)
- [ ] `./bin/rooda info agents-sync` shows procedure metadata (display, summary, description)
- [ ] `./bin/rooda version` displays version information
- [ ] `./bin/rooda help` displays root help with subcommands
- [ ] `./bin/rooda run --help` displays run command help with all flags
- [ ] `./bin/rooda help run` works (alias for `rooda run --help`)

### Dynamic Behavior
- [ ] Custom procedures from workspace config appear in `list` output
- [ ] Unknown procedure name produces clear error with suggestion to run `rooda list`
- [ ] Exit codes unchanged (0, 1, 2, 130)
- [ ] Error messages are clear and actionable

### Breaking Changes
- [ ] Old syntax `rooda <procedure>` no longer works (shows helpful error)
- [ ] Old flags `--version` and `--list-procedures` no longer work (shows helpful error)

## Implementation Notes

### File Structure
```
cmd/rooda/
├── main.go              # Entry point, calls root command
├── root.go              # Root command definition, global flags
├── run.go               # Run subcommand: execute procedures
├── list.go              # List subcommand: show available procedures
├── info.go              # Info subcommand: show procedure details
├── version.go           # Version subcommand
└── flags.go             # Shared flag definitions and helpers
```

### Key Cobra Concepts
- Root `cobra.Command` with subcommands
- `run` command with `Args: cobra.ExactArgs(1)` for procedure name
- `PersistentFlags()` on root for global flags (--config, --verbose, --quiet)
- `Flags()` on run command for execution flags (--max-iterations, --ai-cmd, etc.)
- `PreRunE` on run command for config loading and procedure validation
- `RunE` for command execution
- `ValidArgsFunction` for shell completion of procedure names

### Migration Strategy
1. Add Cobra dependency to go.mod
2. Create root command with global flags
3. Create `run` subcommand with procedure name as argument and all execution flags
4. Create `list` subcommand that loads config and displays procedures
5. Create `info` subcommand that loads config and displays procedure metadata
6. Create `version` subcommand
7. Implement PreRunE on `run`: load config, validate procedure exists
8. Implement RunE on `run`: execute procedure (existing loop logic)
9. Test with built-in and custom procedures
10. Remove old manual parsing code
11. Update all documentation to reflect new CLI structure
12. Update specs/cli-interface.md with new command structure

### Breaking Changes

**This is a breaking change for v0.x users.**

**Old syntax (no longer works):**
```bash
rooda agents-sync --ai-cmd-alias kiro-cli
rooda version
rooda list
```

**New syntax:**
```bash
rooda run agents-sync --ai-cmd-alias kiro-cli
rooda version
rooda list
```

**Error messages:**
- If user runs `rooda agents-sync`, show: "Unknown command 'agents-sync'. Did you mean 'rooda run agents-sync'? Run 'rooda help' for usage."
- If user runs `rooda version`, show: "Flag --version has been replaced by 'rooda version' command."
- If user runs `rooda list`, show: "Flag --list-procedures has been replaced by 'rooda list' command."

## Testing Strategy

1. **Unit tests:** Test flag parsing and validation
2. **Integration tests:** Test full command execution
3. **Regression tests:** Verify all existing commands work identically
4. **Manual testing:** Run through common workflows

## References

- Cobra documentation: https://github.com/spf13/cobra
- CLI spec: `specs/cli-interface.md`
- Current implementation: `cmd/rooda/main.go`
- Flag definitions in specs: `specs/cli-interface.md` (Acceptance Criteria section)

## Non-Goals

- Changing flag semantics (--max-iterations behavior unchanged)
- Adding new procedures (separate task)
- Changing configuration system (separate concern)
- Adding interactive prompts (rooda is non-interactive by design)
- Hardcoding procedure names in CLI code (must remain dynamic)
- Changing exit codes or error handling behavior

## Risks

1. **Breaking changes:** New CLI structure breaks existing scripts/workflows
   - Mitigation: Clear error messages with hints to new syntax
   - Acceptable: We're in v0.x, breaking changes expected
2. **Dependency bloat:** Cobra adds ~500KB to binary
   - Mitigation: Acceptable tradeoff for maintainability
3. **Learning curve:** Team needs to learn Cobra patterns
   - Mitigation: Well-documented, widely-used framework
4. **Dynamic completion:** Shell completion for procedure names requires runtime config loading
   - Mitigation: Use ValidArgsFunction to load config and return procedure names
5. **Documentation updates:** All docs need updating for new CLI structure
   - Mitigation: Update docs as part of this task (README, specs, examples)

## Success Metrics

- New CLI structure is more intuitive and discoverable
- Help text is clearer and more consistent
- Adding new procedures requires no CLI code changes
- Flag validation errors are more helpful
- Shell completion works out of the box
- `list` and `info` commands improve procedure discoverability
