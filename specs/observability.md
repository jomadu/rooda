# Observability

## Job to be Done

Provide visibility into what the loop is doing (timing, iteration progress, outcome) and controls to stop, dry-run, and override behavior. The developer wants to understand loop execution state, diagnose failures, and control output verbosity without modifying configuration files.

This spec defines the logging, progress display, and runtime control mechanisms that make the loop observable and controllable.

## Activities

1. Emit structured log events at defined levels (debug, info, warn, error)
2. Display iteration progress (start, completion, timing, outcome)
3. Calculate and display iteration statistics (count, min, max, mean, stddev)
4. Validate configuration and prompts in dry-run mode without executing AI CLI
5. Stream AI CLI output to terminal when verbose mode enabled
6. Format log output with timestamp, level, and structured fields
7. Route log output to appropriate destination (stdout, stderr, file)
8. Handle log level configuration from multiple sources (config, env, flags)
9. Display resolved configuration with provenance in dry-run mode

## Acceptance Criteria

### Structured Logging

- [ ] Log events emitted at four levels: debug, info, warn, error
- [ ] Log level configurable via `loop.log_level` (debug, info, warn, error, built-in default: info)
- [ ] `ROODA_LOG_LEVEL` environment variable sets `loop.log_level`
- [ ] `--log-level=<level>` flag overrides `loop.log_level`
- [ ] `--quiet` flag overrides `loop.log_level` to warn
- [ ] Log format includes timestamp, level, message, and structured fields
- [ ] Timestamp format: ISO 8601 with timezone (e.g., `2026-02-08T20:59:35.877-08:00`)
- [ ] Structured fields include: procedure, iteration, elapsed_time, exit_code, status
- [ ] Log output routed to stderr by default
- [ ] Log output can be redirected to file via `loop.log_file` (optional, default: stderr)
- [ ] Invalid log level produces error: "Invalid log level '<level>'. Valid levels: debug, info, warn, error."

### Progress Display

- [ ] Iteration start logged at info level: "Starting iteration N/M (procedure: <name>)"
- [ ] Iteration completion logged at info level: "Completed iteration N/M (elapsed: Xs, status: <status>)"
- [ ] Loop start logged at info level: "Starting loop (procedure: <name>, max_iterations: <n>)"
- [ ] Loop completion logged at info level: "Loop completed (status: <status>, iterations: N, total_elapsed: Xs)"
- [ ] Progress messages suppressed when log level > info
- [ ] Progress messages include iteration number (1-indexed for display)
- [ ] Progress messages include elapsed time in human-readable format (e.g., "1.23s", "2m 15s")
- [ ] Loop status displayed: success, max-iters, aborted, interrupted

### Iteration Statistics

- [ ] Statistics calculated using constant memory (O(1)) regardless of iteration count
- [ ] Statistics include: count, min, max, mean, stddev
- [ ] Statistics displayed at loop completion (info level)
- [ ] Statistics format: "Iteration timing: count=N, min=Xs, max=Xs, mean=Xs, stddev=Xs"
- [ ] Statistics omitted if count < 2 (stddev undefined for single iteration)
- [ ] Statistics use Welford's online algorithm for numerical stability

### Dry-Run Mode

- [ ] `--dry-run` flag enables dry-run mode
- [ ] Dry-run validates all prompt files exist and are readable
- [ ] Dry-run validates AI command binary exists and is executable
- [ ] Dry-run displays assembled prompt with clear section markers
- [ ] Dry-run displays resolved configuration with provenance
- [ ] Dry-run exits with code 0 if all validations pass
- [ ] Dry-run exits with code 1 if user error (invalid flags, unknown procedure)
- [ ] Dry-run exits with code 2 if config error (missing AI command, invalid config, missing prompt files)
- [ ] Dry-run does not execute AI CLI
- [ ] Dry-run does not modify any files
- [ ] Dry-run displays configuration provenance: "max_iterations: 5 (from: procedure default)"
- [ ] Provenance sources: built-in default, global config, procedure config, environment variable, CLI flag

### Verbose Mode

- [ ] `--verbose` flag enables verbose mode
- [ ] Verbose mode streams AI CLI output to terminal in real-time
- [ ] Verbose mode overrides `loop.show_ai_output` to true
- [ ] AI CLI output always captured and scanned for `<promise>` signals regardless of verbose mode
- [ ] Without verbose mode, only loop-level progress displayed (iteration start/complete, timing, outcome)
- [ ] Verbose mode displays configuration provenance for all resolved settings
- [ ] `loop.show_ai_output` configurable (true, false, built-in default: false)
- [ ] `ROODA_SHOW_AI_OUTPUT` environment variable sets `loop.show_ai_output`
- [ ] `--verbose` and `--quiet` are mutually exclusive (error if both provided)

### Log Format

- [ ] Log format: `<timestamp> <level> <message> [<fields>]`
- [ ] Example: `2026-02-08T20:59:35.877-08:00 INFO Starting iteration 1/5 (procedure: build) [elapsed=0s]`
- [ ] Structured fields formatted as key=value pairs in brackets
- [ ] Multi-word field values quoted: `status="max-iters"`
- [ ] Numeric field values unquoted: `iteration=3`
- [ ] Log format is human-readable (not JSON) by default
- [ ] Future: JSON log format via `loop.log_format` (text, json, built-in default: text)

## Data Structures

### LogLevel

```go
type LogLevel int

const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
)

func (l LogLevel) String() string {
    return [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[l]
}
```

### LogEvent

```go
type LogEvent struct {
    Timestamp time.Time
    Level     LogLevel
    Message   string
    Fields    map[string]interface{} // Structured fields (procedure, iteration, elapsed_time, etc.)
}
```

### IterationStats

Defined in `iteration-loop.md`, referenced here for statistics display:

```go
type IterationStats struct {
    Count     int           // Total iterations completed
    Min       time.Duration // Fastest iteration
    Max       time.Duration // Slowest iteration
    Mean      time.Duration // Average iteration time
    M2        float64       // Sum of squared differences (for stddev calculation)
}
```

## Examples

### Basic Progress Display (Default)

```
$ rooda build
2026-02-08T21:00:00.000-08:00 INFO Starting loop (procedure: build, max_iterations: 5)
2026-02-08T21:00:00.100-08:00 INFO Starting iteration 1/5 (procedure: build)
2026-02-08T21:00:15.200-08:00 INFO Completed iteration 1/5 (elapsed: 15.1s, status: success)
2026-02-08T21:00:15.300-08:00 INFO Starting iteration 2/5 (procedure: build)
2026-02-08T21:00:28.400-08:00 INFO Completed iteration 2/5 (elapsed: 13.1s, status: success)
2026-02-08T21:00:28.500-08:00 INFO Loop completed (status: success, iterations: 2, total_elapsed: 28.5s)
2026-02-08T21:00:28.500-08:00 INFO Iteration timing: count=2, min=13.1s, max=15.1s, mean=14.1s, stddev=1.0s
```

### Verbose Mode (AI Output Streaming)

```
$ rooda build --verbose
2026-02-08T21:00:00.000-08:00 INFO Starting loop (procedure: build, max_iterations: 5)
2026-02-08T21:00:00.100-08:00 INFO Starting iteration 1/5 (procedure: build)
2026-02-08T21:00:00.200-08:00 DEBUG Executing AI CLI: kiro-cli chat --prompt-file /tmp/rooda-prompt-12345.md
--- AI CLI Output Start ---
I'll execute the OODA loop iteration systematically.

## OBSERVE
...
<promise>SUCCESS</promise>
--- AI CLI Output End ---
2026-02-08T21:00:15.200-08:00 INFO Completed iteration 1/5 (elapsed: 15.1s, status: success)
2026-02-08T21:00:15.300-08:00 INFO Loop completed (status: success, iterations: 1, total_elapsed: 15.2s)
```

### Dry-Run Mode (Validation Only)

```
$ rooda build --dry-run
2026-02-08T21:00:00.000-08:00 INFO Dry-run mode enabled (no AI CLI execution)
2026-02-08T21:00:00.100-08:00 INFO Validating configuration...
2026-02-08T21:00:00.200-08:00 INFO Resolved configuration:
  procedure: build (from: CLI argument)
  max_iterations: 5 (from: procedure default)
  iteration_timeout: nil (from: built-in default)
  ai_command: kiro-cli chat (from: global config)
  log_level: info (from: built-in default)
  show_ai_output: false (from: built-in default)
2026-02-08T21:00:00.300-08:00 INFO Validating prompt files...
  observe: prompts/observe_plan_specs_impl.md (exists, readable)
  orient: prompts/orient_build.md (exists, readable)
  decide: prompts/decide_build.md (exists, readable)
  act: prompts/act_build.md (exists, readable)
2026-02-08T21:00:00.400-08:00 INFO Validating AI command...
  command: kiro-cli (found at /usr/local/bin/kiro-cli, executable)
2026-02-08T21:00:00.500-08:00 INFO Assembled prompt (1234 bytes):
--- Prompt Start ---
# OODA Loop Iteration

## OBSERVE
...
--- Prompt End ---
2026-02-08T21:00:00.600-08:00 INFO Dry-run validation passed
```

### Quiet Mode (Warnings and Errors Only)

```
$ rooda build --quiet
2026-02-08T21:00:28.500-08:00 WARN Iteration 3 exceeded timeout (30s), killing AI CLI process
2026-02-08T21:00:28.600-08:00 ERROR Loop aborted (consecutive failures: 3, threshold: 3)
```

### Debug Mode (Detailed Logging)

```
$ rooda build --log-level=debug
2026-02-08T21:00:00.000-08:00 DEBUG Loading configuration from rooda-config.yml
2026-02-08T21:00:00.050-08:00 DEBUG Merged configuration: 15 procedures, 3 AI command aliases
2026-02-08T21:00:00.100-08:00 DEBUG Resolving procedure: build
2026-02-08T21:00:00.150-08:00 DEBUG Procedure found: build (observe: observe_plan_specs_impl.md, orient: orient_build.md, decide: decide_build.md, act: act_build.md)
2026-02-08T21:00:00.200-08:00 INFO Starting loop (procedure: build, max_iterations: 5)
2026-02-08T21:00:00.250-08:00 DEBUG Assembling prompt from 4 phase files
2026-02-08T21:00:00.300-08:00 DEBUG Prompt assembled: 1234 bytes
2026-02-08T21:00:00.350-08:00 INFO Starting iteration 1/5 (procedure: build)
2026-02-08T21:00:00.400-08:00 DEBUG Executing AI CLI: kiro-cli chat --prompt-file /tmp/rooda-prompt-12345.md
2026-02-08T21:00:15.200-08:00 DEBUG AI CLI exited with code 0
2026-02-08T21:00:15.250-08:00 DEBUG Scanning output for promise signals (1024 bytes)
2026-02-08T21:00:15.300-08:00 DEBUG Found promise signal: SUCCESS
2026-02-08T21:00:15.350-08:00 INFO Completed iteration 1/5 (elapsed: 15.1s, status: success)
2026-02-08T21:00:15.400-08:00 INFO Loop completed (status: success, iterations: 1, total_elapsed: 15.2s)
```

## Related Specs

- `iteration-loop.md` — Defines iteration statistics calculation and loop termination
- `cli-interface.md` — Defines flags for controlling observability (--verbose, --quiet, --log-level, --dry-run)
- `error-handling.md` — Defines error logging format and levels
- `ai-cli-integration.md` — Defines AI CLI execution logging events
- `configuration.md` — Defines configuration precedence for log_level and show_ai_output

## Topics of Concern

**Logging:**
- Structured log events with timestamp, level, message, fields
- Log level configuration (config, env, flags)
- Log output destination (stderr, file)
- Log format (human-readable text, future: JSON)

**Progress Display:**
- Iteration start/completion messages
- Loop start/completion messages
- Elapsed time formatting
- Status display (success, max-iters, aborted, interrupted)

**Statistics:**
- Constant memory calculation (Welford's algorithm)
- Min, max, mean, stddev
- Display at loop completion

**Dry-Run Mode:**
- Configuration validation
- Prompt file validation
- AI command validation
- Assembled prompt display
- Configuration provenance display
- Exit codes (0=pass, 1=user error, 2=config error)

**Verbose Mode:**
- AI CLI output streaming
- Configuration provenance display
- Overrides show_ai_output setting

**Controls:**
- --verbose / --quiet mutual exclusivity
- Log level precedence (flag > env > config > default)
- Output suppression based on log level
