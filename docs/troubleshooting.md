# Troubleshooting

Common errors and solutions when using rooda.

## Installation Issues

### "rooda: command not found"

**Cause**: Binary not in PATH.

**Solution**:
```bash
# Check if rooda is installed
which rooda

# If not found, reinstall
curl -fsSL https://raw.githubusercontent.com/jomadu/rooda/main/scripts/install.sh | bash

# Or add to PATH manually
export PATH="$PATH:/usr/local/bin"
```

### "permission denied" when running rooda

**Cause**: Binary not executable.

**Solution**:
```bash
chmod +x /usr/local/bin/rooda
```

## Configuration Errors

### "No AI command configured"

**Error**:
```
Error: No AI command configured. Set one of:
  --ai-cmd <command>
  --ai-cmd-alias <alias>
  ROODA_LOOP_AI_CMD
  ROODA_LOOP_AI_CMD_ALIAS
  loop.ai_cmd in config file
  loop.ai_cmd_alias in config file
```

**Cause**: No AI command specified in any configuration source.

**Solution**:
```bash
# Use CLI flag
rooda run build --ai-cmd-alias kiro-cli

# Or set environment variable
export ROODA_LOOP_AI_CMD_ALIAS=kiro-cli
rooda build

# Or add to global config
mkdir -p ~/.config/rooda
cat > ~/.config/rooda/rooda-config.yml <<EOF
loop:
  ai_cmd_alias: kiro-cli
EOF
```

### "Unknown procedure '<name>'"

**Error**:
```
Error: Unknown procedure 'my-procedure'. Run 'rooda list' to see available procedures.
```

**Cause**: Procedure name doesn't exist in configuration.

**Solution**:
```bash
# List available procedures
rooda list

# Check spelling
rooda run build  # not "rooda Build" or "rooda BUILD"

# If custom procedure, verify it's defined in config
cat rooda-config.yml
```

### "Invalid YAML in config file"

**Error**:
```
Error: Failed to parse config file ./rooda-config.yml: yaml: line 5: mapping values are not allowed in this context
```

**Cause**: Syntax error in YAML file.

**Solution**:
```bash
# Validate YAML syntax
yq eval . rooda-config.yml

# Common issues:
# - Missing colon after key
# - Incorrect indentation (use spaces, not tabs)
# - Unquoted strings with special characters
```

### "Fragment file not found"

**Error**:
```
Error: Fragment file not found: prompts/observe_custom.md
```

**Cause**: Referenced prompt file doesn't exist.

**Solution**:
```bash
# Check file exists
ls -la prompts/observe_custom.md

# Verify path is relative to config file directory
# If config is at ./rooda-config.yml, fragments resolve from ./
# If config is at /path/to/config.yml, fragments resolve from /path/to/

# Use absolute path if needed
rooda run build --observe /absolute/path/to/prompts/observe_custom.md
```

## Execution Errors

### "AI CLI command failed with exit code 1"

**Error**:
```
Error: AI CLI execution failed with exit code 1
```

**Cause**: AI command returned non-zero exit code.

**Solution**:
```bash
# Test AI command manually
kiro-cli chat

# Check AI command is in PATH
which kiro-cli

# Verify AI command alias is correct
rooda run build --ai-cmd "kiro-cli chat" --verbose

# Check AI command output for errors
rooda run build --ai-cmd-alias kiro-cli --verbose 2>&1 | less
```

### "Iteration timeout exceeded"

**Error**:
```
Error: Iteration timeout exceeded (3600s)
```

**Cause**: AI CLI execution took longer than configured timeout.

**Solution**:
```bash
# Increase timeout (in seconds)
rooda run build --ai-cmd-alias kiro-cli  # No timeout by default

# Or set in config
cat > rooda-config.yml <<EOF
loop:
  iteration_timeout: 7200  # 2 hours
EOF

# Or disable timeout
cat > rooda-config.yml <<EOF
loop:
  iteration_timeout: null  # No timeout
EOF
```

### "Failure threshold exceeded"

**Error**:
```
Error: Failure threshold exceeded (3 consecutive failures)
```

**Cause**: AI CLI failed 3 times in a row.

**Solution**:
```bash
# Check AI output for errors
rooda run build --ai-cmd-alias kiro-cli --verbose

# Increase failure threshold
cat > rooda-config.yml <<EOF
loop:
  failure_threshold: 5
EOF

# Or fix underlying issue causing failures
# - Check AGENTS.md is accurate
# - Verify test/build commands work
# - Review AI output for error patterns
```

### "Max iterations reached without success"

**Error**:
```
Error: Max iterations reached (5) without success signal
```

**Cause**: Loop completed all iterations but didn't detect success.

**Solution**:
```bash
# Increase max iterations
rooda run build --ai-cmd-alias kiro-cli --max-iterations 10

# Or use unlimited mode
rooda run build --ai-cmd-alias kiro-cli --unlimited

# Or check if success signal is being emitted
# Success signal: "ROODA_SUCCESS" in AI output
# Verify tests pass and emit success signal
```

## Work Tracking Issues

### "bd: command not found"

**Cause**: beads CLI not installed.

**Solution**:
```bash
# Install beads
# See https://github.com/beadslabs/beads for installation

# Or use different work tracking system
# Update AGENTS.md to use GitHub Issues or file-based tracking
```

### "No ready work found"

**Error**:
```
Error: No ready work found in work tracking system
```

**Cause**: No tasks with status "open" or "ready".

**Solution**:
```bash
# Check work tracking system
bd ready --json

# Create a task
bd create --title "Implement feature X" --priority 2

# Or use planning procedures to generate tasks
rooda draft-plan-impl-feat --ai-cmd-alias kiro-cli --context "Feature X"
rooda publish-plan --ai-cmd-alias kiro-cli
```

## AGENTS.md Issues

### "AGENTS.md not found"

**Error**:
```
Error: AGENTS.md not found. Run 'rooda bootstrap' to create it.
```

**Cause**: AGENTS.md doesn't exist in repository root.

**Solution**:
```bash
# Bootstrap the repository
rooda bootstrap --ai-cmd-alias kiro-cli

# This creates AGENTS.md with detected settings
```

### "Command in AGENTS.md failed"

**Error**:
```
Error: Test command failed: go test ./...
```

**Cause**: Command documented in AGENTS.md doesn't work.

**Solution**:
```bash
# Test command manually
go test ./...

# If command is wrong, update AGENTS.md
rooda agents-sync --ai-cmd-alias kiro-cli --max-iterations 1

# Or manually edit AGENTS.md
vim AGENTS.md
```

## Prompt Issues

### "Empty prompt generated"

**Error**:
```
Error: Empty prompt generated for phase: observe
```

**Cause**: No fragments defined for OODA phase.

**Solution**:
```bash
# Check procedure definition
rooda run build --help

# Verify fragments exist
ls -la prompts/

# Add fragments to procedure
cat > rooda-config.yml <<EOF
procedures:
  build:
    observe:
      - path: "prompts/observe_specs.md"
EOF
```

### "Template parameter not found"

**Error**:
```
Error: Template parameter 'project_name' not found
```

**Cause**: Fragment uses template parameter that wasn't provided.

**Solution**:
```yaml
# Add parameters to fragment action
procedures:
  my-procedure:
    observe:
      - path: "prompts/observe_custom.md"
        parameters:
          project_name: "my-project"
```

## Performance Issues

### "AI CLI output too large"

**Error**:
```
Warning: AI CLI output truncated (exceeded 10MB buffer)
```

**Cause**: AI output exceeded max buffer size.

**Solution**:
```yaml
# Increase buffer size (in bytes)
loop:
  max_output_buffer: 20971520  # 20MB
```

### "Iterations taking too long"

**Symptom**: Each iteration takes several minutes.

**Solution**:
```bash
# Use faster AI model
rooda run build --ai-cmd "claude --model claude-3-haiku"

# Reduce context size
# - Simplify prompts
# - Reduce number of fragments
# - Use more focused procedures

# Set iteration timeout
cat > rooda-config.yml <<EOF
loop:
  iteration_timeout: 300  # 5 minutes
EOF
```

## Debugging

### Enable Verbose Output

```bash
rooda run build --ai-cmd-alias kiro-cli --verbose
```

Shows:
- Configuration provenance
- Assembled prompts
- AI CLI output
- Debug logs

### Dry Run

```bash
rooda run build --ai-cmd-alias kiro-cli --dry-run
```

Validates:
- Configuration is valid
- Prompts exist and can be assembled
- AI command is found
- No execution

### Check Configuration

```bash
# List procedures
rooda list

# Show procedure help
rooda run build --help

# Validate config file
yq eval . rooda-config.yml
```

### Check Logs

```bash
# Run with debug logging
rooda run build --ai-cmd-alias kiro-cli --log-level debug

# Save logs to file
rooda run build --ai-cmd-alias kiro-cli --verbose 2>&1 | tee rooda.log
```

## Getting Help

If you're still stuck:

1. **Check documentation**: [docs/](.)
2. **Search issues**: [GitHub Issues](https://github.com/jomadu/rooda/issues)
3. **File a bug**: [New Issue](https://github.com/jomadu/rooda/issues/new)

Include in bug reports:
- `rooda version` output
- Full error message
- Configuration files (sanitize secrets)
- Steps to reproduce
- Operating system and Go version (if building from source)
