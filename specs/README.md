# Specifications

This directory contains specifications for the ralph-wiggum-ooda framework. Each specification follows the Jobs-to-be-Done (JTBD) structure defined in [specification-system.md](specification-system.md).

## How to Write Specs

Use [TEMPLATE.md](TEMPLATE.md) as a starting point for new specifications. Follow the structure and guidelines in [specification-system.md](specification-system.md).

## Specifications

### [AGENTS.md Specification](agents-md-format.md)
AGENTS.md is the interface between agents and the repository. It defines how agents interact with project-specific workflows, tools, and conventions.

### [AI CLI Integration](ai-cli-integration.md)
Execute OODA loop prompts through an AI CLI tool that can read files, modify code, run commands, and interact with the repository autonomously.

### [CLI Interface](cli-interface.md)
Enable users to invoke OODA loop procedures through a command-line interface, supporting both named procedures from configuration and explicit OODA phase file specification.

### [Component Authoring](component-authoring.md)
Enable developers to create and modify OODA component prompt files that can be composed into executable procedures.

### [Configuration Schema](configuration-schema.md)
Enable users to define custom OODA loop procedures by mapping procedure names to composable prompt component files, supporting both predefined framework procedures and user-defined custom procedures.

### [External Dependencies](external-dependencies.md)
Enable users to install and verify all required external tools before running ralph-wiggum-ooda procedures, preventing runtime failures due to missing dependencies.

### [Iteration Loop Control](iteration-loop.md)
Execute OODA loop procedures through controlled iteration cycles that clear context between runs, preventing LLM degradation while maintaining file-based state continuity.

### [User Documentation](user-documentation.md)
Enable users to understand and effectively use the ralph-wiggum-ooda framework through clear, accessible documentation that guides them from installation through advanced usage.
