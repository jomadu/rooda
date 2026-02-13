# rooda

Run AI agents through OODA loops to build software autonomously.

## What is rooda?

rooda runs AI coding agents through OODA (Observe-Orient-Decide-Act) loops. Each iteration gives the agent fresh context to work with. Point it at a task, set a max iteration count, and it works toward a solution.

The agent observes your codebase, orients itself to the problem, decides what to do, and acts. Then it does it again with updated context. No massive prompts that lose focus halfway through.

## Quick Start

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/jomadu/rooda/main/scripts/install.sh | bash

# Bootstrap a repository (creates/updates AGENTS.md)
rooda run bootstrap --ai-cmd-alias kiro-cli

# List available procedures
rooda list

# Run a procedure
rooda run build --ai-cmd-alias kiro-cli --max-iterations 3
```

## Installation

See [docs/installation.md](docs/installation.md) for all installation methods (direct download, build from source).

## Core Concepts

**Procedures** - Named OODA workflows like `build`, `audit-spec`, or `draft-plan-impl-feat`. Each procedure defines which prompt fragments to use for each phase.

**OODA Loop** - Observe (gather context) → Orient (analyze) → Decide (plan) → Act (execute). Each iteration runs all four phases with fresh AI context.

**AGENTS.md** - A file that tells agents how to interact with your project: build commands, test commands, work tracking system, file patterns.

**Work Tracking** - rooda works with beads, GitHub Issues, or file-based systems. The `build` procedure reads tasks and implements them.

## Configuration

rooda uses a three-tier configuration system:

1. **Built-in defaults**: 21 procedures embedded in the binary
2. **Global config**: `~/.config/rooda/rooda-config.yml` (team-wide settings)
3. **Workspace config**: `./rooda-config.yml` (project-specific settings)

CLI flags override everything. See [docs/configuration.md](docs/configuration.md) for details.

## Common Workflows

**Start a new project:**
```bash
rooda run bootstrap --ai-cmd-alias kiro-cli
rooda run audit-spec --ai-cmd-alias kiro-cli
```

**Implement from work tracking:**
```bash
rooda run build --ai-cmd-alias kiro-cli --max-iterations 5
```

**Plan a feature:**
```bash
rooda run draft-plan-impl-feat --ai-cmd-alias kiro-cli --context "Add user authentication"
```

**Audit implementation:**
```bash
rooda run audit-impl --ai-cmd-alias kiro-cli
```

## Documentation

- [Installation](docs/installation.md) - All install methods
- [Procedures](docs/procedures.md) - All 21 procedures with examples
- [Configuration](docs/configuration.md) - Three-tier config system
- [CLI Reference](docs/cli-reference.md) - All flags and exit codes
- [Troubleshooting](docs/troubleshooting.md) - Common errors and solutions
- [AGENTS.md Format](docs/agents-md.md) - Repository operational guide format

## Requirements

- Go >= 1.24.5 (if building from source)
- AI CLI tool (kiro-cli, claude, cursor-agent, etc.)
- Work tracking system (beads recommended, GitHub Issues supported)

## License

MIT
