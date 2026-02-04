# Specifications

This directory contains specifications for the ralph-wiggum-ooda framework.

See [TEMPLATE.md](TEMPLATE.md) for the specification structure.
See [specification-system.md](specification-system.md) for the spec system design.

## Specifications

### [agents-md-format.md](agents-md-format.md)
AGENTS.md is the interface between agents and the repository. It defines how agents interact with project-specific workflows, tools, and conventions.

### [ai-cli-integration.md](ai-cli-integration.md)
Execute OODA loop prompts through an AI CLI tool that can read files, modify code, run commands, and interact with the repository autonomously.

### [cli-interface.md](cli-interface.md)
Enable users to invoke OODA loop procedures through a command-line interface, supporting both named procedures from configuration and explicit OODA phase file specification.

### [component-authoring.md](component-authoring.md)
Enable developers to create and modify OODA component prompt files that can be composed into executable procedures.

### [component-system.md](component-system.md) (deprecated)
Enable agents to execute procedures through composable, reusable prompt components that maintain consistency across iterations while allowing flexible procedure definitions.

### [configuration-schema.md](configuration-schema.md)
Enable users to define custom OODA loop procedures by mapping procedure names to composable prompt component files, supporting both predefined framework procedures and user-defined custom procedures.

### [external-dependencies.md](external-dependencies.md)
Enable users to install and verify all required external tools before running ralph-wiggum-ooda procedures, preventing runtime failures due to missing dependencies.

### [iteration-loop.md](iteration-loop.md)
Execute OODA loop procedures through controlled iteration cycles that clear context between runs, preventing LLM degradation while maintaining file-based state continuity.

### [prompt-composition.md](prompt-composition.md) (deprecated)
Assemble four OODA phase prompt files into a single executable prompt that can be piped to the AI CLI.

