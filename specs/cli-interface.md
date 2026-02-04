# CLI Interface

## Job to be Done
Enable users to invoke OODA loop procedures through a command-line interface, supporting both named procedures from configuration and explicit OODA phase file specification.

## Activities
1. Parse command-line arguments to determine invocation mode
2. Resolve procedure configuration or explicit OODA phase files
3. Validate all required arguments and file paths
4. Display execution parameters before starting loop
5. Execute iteration loop with configured parameters

## Acceptance Criteria
- [x] Procedure-based invocation loads OODA files from config
- [x] Explicit flag invocation accepts four OODA phase files directly
- [x] Explicit flags override config-based procedure settings
- [x] Config file resolves relative to script location
- [x] Missing files produce clear error messages
- [x] Invalid arguments produce usage help
- [x] Max iterations can be specified or defaults to procedure config

## Data Structures

### Command-Line Arguments
```bash
# Procedure-based invocation
./rooda.sh <procedure> [--config <file>] [--max-iterations N]

# Explicit flag invocation
./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]
```

**Arguments:**
- `<procedure>` - Named procedure from config (optional, positional)
- `--config <file>` - Path to config file (default: rooda-config.yml in script directory)
- `--observe <file>` - Path to observe phase prompt
- `--orient <file>` - Path to orient phase prompt
- `--decide <file>` - Path to decide phase prompt
- `--act <file>` - Path to act phase prompt
- `--max-iterations N` - Maximum iterations (default: from config or 0 for unlimited)

## Algorithm

1. Check for yq dependency, exit if missing
2. Initialize variables (OBSERVE, ORIENT, DECIDE, ACT, MAX_ITERATIONS, PROCEDURE, CONFIG_FILE)
3. Resolve CONFIG_FILE relative to script location
4. Parse first positional argument as PROCEDURE if not flag
5. Parse remaining arguments in while loop
6. If PROCEDURE specified:
   - Validate config file exists
   - Query config for procedure OODA files using yq
   - Exit if procedure not found
   - Load default_iterations if MAX_ITERATIONS not specified
7. Validate all four OODA phase files specified
8. Validate all four files exist on filesystem
9. Display execution parameters
10. Enter iteration loop

**Pseudocode:**
```bash
if ! command -v yq; then
    error "yq required"
    exit 1

# Parse arguments
if first_arg not starts_with "--"; then
    PROCEDURE = first_arg
    shift

while has_args:
    case arg:
        --config: CONFIG_FILE = next_arg
        --observe: OBSERVE = next_arg
        --orient: ORIENT = next_arg
        --decide: DECIDE = next_arg
        --act: ACT = next_arg
        --max-iterations: MAX_ITERATIONS = next_arg
        default: error "Unknown option"

# Resolve from config if procedure specified
if PROCEDURE:
    if not exists(CONFIG_FILE):
        error "Config not found"
    OBSERVE = yq(".procedures.$PROCEDURE.observe", CONFIG_FILE)
    ORIENT = yq(".procedures.$PROCEDURE.orient", CONFIG_FILE)
    DECIDE = yq(".procedures.$PROCEDURE.decide", CONFIG_FILE)
    ACT = yq(".procedures.$PROCEDURE.act", CONFIG_FILE)
    if any_null:
        error "Procedure not found"
    if MAX_ITERATIONS == 0:
        MAX_ITERATIONS = yq(".procedures.$PROCEDURE.default_iterations", CONFIG_FILE)

# Validate
if not all(OBSERVE, ORIENT, DECIDE, ACT):
    error "All four OODA phases required"
for file in [OBSERVE, ORIENT, DECIDE, ACT]:
    if not exists(file):
        error "File not found: $file"
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| yq not installed | Error message with installation instructions, exit 1 |
| No arguments provided | Error with usage help, exit 1 |
| Unknown flag | Error with flag name and usage help, exit 1 |
| Procedure not in config | Error "Procedure 'X' not found in config", exit 1 |
| Config file missing | Error "Config not found", exit 1 |
| OODA phase file missing | Error "File not found: path", exit 1 |
| Only some OODA phases specified | Error "All four OODA phases required", exit 1 |
| Explicit flags with procedure | Explicit flags take precedence, procedure ignored |
| --max-iterations 0 | Unlimited iterations (loop until Ctrl+C) |
| --max-iterations not specified | Use default_iterations from config, or 0 if not in config |

## Dependencies

- yq - YAML query tool for parsing rooda-config.yml
- bash - Shell interpreter
- rooda-config.yml - Configuration file with procedure definitions

## Implementation Mapping

**Source files:**
- `src/rooda.sh` - Lines 1-141 implement argument parsing and validation

**Related specs:**
- `configuration-schema.md` - Defines rooda-config.yml structure (to be created)
- `iteration-loop.md` - Defines loop execution behavior (to be created)

## Examples

### Example 1: Procedure-Based Invocation

**Input:**
```bash
./rooda.sh bootstrap
```

**Expected Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Procedure: bootstrap
Observe:   src/prompts/observe_bootstrap.md
Orient:    src/prompts/orient_bootstrap.md
Decide:    src/prompts/decide_bootstrap.md
Act:       src/prompts/act_bootstrap.md
Branch:    main
Max:       1 iterations
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Verification:**
- OODA files loaded from config
- default_iterations from config used (1 for bootstrap)
- Iteration loop starts

### Example 2: Explicit Flag Invocation

**Input:**
```bash
./rooda.sh \
  --observe src/prompts/observe_specs.md \
  --orient src/prompts/orient_gap.md \
  --decide src/prompts/decide_gap_plan.md \
  --act src/prompts/act_plan.md \
  --max-iterations 1
```

**Expected Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Observe:   src/prompts/observe_specs.md
Orient:    src/prompts/orient_gap.md
Decide:    src/prompts/decide_gap_plan.md
Act:       src/prompts/act_plan.md
Branch:    main
Max:       1 iterations
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Verification:**
- No procedure name shown
- Explicit files used
- Max iterations from command line

### Example 3: Override Default Iterations

**Input:**
```bash
./rooda.sh build --max-iterations 10
```

**Expected Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Procedure: build
Observe:   src/prompts/observe_plan_specs_impl.md
Orient:    src/prompts/orient_build.md
Decide:    src/prompts/decide_build.md
Act:       src/prompts/act_build.md
Branch:    main
Max:       10 iterations
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Verification:**
- Procedure loaded from config
- Command-line --max-iterations overrides config default (5)

### Example 4: Missing yq Dependency

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
- Clear installation instructions provided

### Example 5: Unknown Procedure

**Input:**
```bash
./rooda.sh nonexistent
```

**Expected Output:**
```
Error: Procedure 'nonexistent' not found in /path/to/rooda-config.yml
```

**Verification:**
- Script exits with status 1
- Error message includes procedure name and config path

### Example 6: Missing OODA Phase File

**Input:**
```bash
./rooda.sh --observe missing.md --orient o.md --decide d.md --act a.md
```

**Expected Output:**
```
Error: File not found: missing.md
```

**Verification:**
- Script exits with status 1
- Error message includes missing file path

## Notes

**Design Rationale:**

The CLI supports two invocation modes to balance convenience and flexibility:

1. **Procedure-based** - Named procedures in config provide convenient shortcuts for common workflows (bootstrap, build, draft-plan-*). This reduces typing and ensures consistent OODA phase combinations.

2. **Explicit flags** - Direct OODA file specification enables custom procedures without modifying config. This supports experimentation and one-off workflows.

**Precedence Rules:**

Explicit flags always override config-based procedure settings. This allows users to customize a procedure's behavior (e.g., swap one OODA phase) without creating a new config entry.

**Config Resolution:**

CONFIG_FILE resolves relative to script location (`SCRIPT_DIR`), not current working directory. This ensures the script finds its config regardless of where it's invoked from. Users can override with `--config` if needed.

**Iteration Defaults:**

The three-tier default system provides flexibility:
1. Command-line `--max-iterations` (highest priority)
2. Config `default_iterations` per procedure
3. 0 (unlimited) if neither specified

This allows procedures to define sensible defaults (bootstrap=1, build=5) while letting users override when needed.

## Known Issues

**Duplicate validation blocks:** Lines 95-103 and 117-125 contain identical validation logic. This duplication exists in the current implementation but should be refactored.

**Inconsistent usage messages:** Error messages show different usage patterns (some include `<task-id>`, some don't). The actual implementation doesn't use task-id as a positional argument.

## Areas for Improvement

**Help flag:** No `--help` or `-h` flag support. Users must trigger an error to see usage.

**Version flag:** No `--version` flag to show script version.

**Verbose/quiet modes:** No control over output verbosity.

**Config validation:** Script doesn't validate config file structure before querying with yq.

**Short flags:** No short flag alternatives (e.g., `-o` for `--observe`, `-m` for `--max-iterations`).
