# Operational Knowledge

## Job to be Done

Every procedure reads AGENTS.md first as the source of truth for project-specific behavior — build commands, file paths, work tracking, quality criteria. Agents defer to it, verify it empirically (run commands, check paths), and update it in-place when something is wrong or a new learning occurs.

The developer wants rooda to adapt to each repository's unique conventions without manual configuration, maintain accuracy as the project evolves, and capture operational learnings so future iterations don't repeat mistakes — all without maintaining a separate configuration file.

## Activities

1. **Read** — Load AGENTS.md at procedure start, parse sections into structured data
2. **Verify** — Empirically validate claims (run commands, check paths, query work tracking)
3. **Execute** — Use AGENTS.md data to guide procedure behavior (run tests, update issues, commit)
4. **Detect Drift** — Compare expected vs actual outcomes during execution
5. **Update** — Modify AGENTS.md in-place when drift detected or new learning occurs
6. **Bootstrap** — Create AGENTS.md from scratch if it doesn't exist (first-run case)

## Acceptance Criteria

- [ ] All procedures read AGENTS.md before executing any work
- [ ] If AGENTS.md doesn't exist, procedure creates it by analyzing the repository
- [ ] AGENTS.md parsed into structured data (build commands, test commands, spec paths, impl paths, work tracking config)
- [ ] Build/test/lint commands from AGENTS.md are executed verbatim (no interpretation)
- [ ] If a command from AGENTS.md fails, procedure attempts to fix the command and updates AGENTS.md
- [ ] File path patterns from AGENTS.md are validated (glob patterns resolve to actual files)
- [ ] If path patterns don't match expected files, procedure updates AGENTS.md with correct patterns
- [ ] Work tracking system commands from AGENTS.md are executed to query/update issues
- [ ] If work tracking commands fail, procedure detects the failure and updates AGENTS.md
- [ ] Quality criteria from AGENTS.md are evaluated (PASS/FAIL checks)
- [ ] When a quality criterion fails, procedure logs the failure and may update AGENTS.md if the criterion is incorrect
- [ ] Operational learnings are incorporated into existing AGENTS.md sections (not appended as diary entries)
- [ ] Updates to AGENTS.md include inline rationale comments explaining why the change was made
- [ ] AGENTS.md updates are committed with descriptive messages referencing the drift or learning
- [ ] Bootstrap case (no AGENTS.md) analyzes repository structure, detects build system, identifies spec/impl patterns, and creates initial AGENTS.md
- [ ] Bootstrap detection heuristics: package.json → npm/yarn, Cargo.toml → cargo, go.mod → go, Makefile → make, etc.
- [ ] Bootstrap spec detection: looks for specs/, docs/, documentation/, README.md patterns
- [ ] Bootstrap impl detection: looks for src/, lib/, cmd/, main.*, index.*, etc.
- [ ] Bootstrap work tracking detection: .beads/ → beads, .github/issues/ → GitHub Issues, TODO.md → file-based
- [ ] Empirical verification runs before trusting AGENTS.md data (e.g., run `--version` flag to verify command exists)
- [ ] Verification failures trigger update workflow (fix + commit)
- [ ] AGENTS.md is treated as living documentation, not static configuration
- [ ] Procedures never assume AGENTS.md is correct — always verify empirically

## Data Structures

### AgentsMD

Parsed representation of AGENTS.md content.

```go
type AgentsMD struct {
    BuildCommand      string              // Command to build the project
    TestCommand       string              // Command to run tests
    LintCommands      []string            // Commands to run linters
    SpecPaths         []string            // Glob patterns for specification files
    SpecExcludes      []string            // Glob patterns to exclude from specs
    ImplPaths         []string            // Glob patterns for implementation files
    ImplExcludes      []string            // Glob patterns to exclude from impl
    WorkTracking      WorkTrackingConfig  // Work tracking system configuration
    QualityCriteria   []QualityCriterion  // Quality checks (PASS/FAIL)
    RawContent        string              // Full AGENTS.md content for updates
    FilePath          string              // Path to AGENTS.md (usually ./AGENTS.md)
}
```

### WorkTrackingConfig

```go
type WorkTrackingConfig struct {
    System        string   // System name (beads, github-issues, file-based)
    QueryCommand  string   // Command to query ready work (e.g., "bd ready --json")
    UpdateCommand string   // Command template to update status (e.g., "bd update <id> --status <status>")
    CloseCommand  string   // Command template to close issue (e.g., "bd close <id> --reason <reason>")
    CreateCommand string   // Command template to create issue (e.g., "bd create --title <title> --description <desc>")
}
```

### QualityCriterion

```go
type QualityCriterion struct {
    Description string   // Human-readable description
    Command     string   // Command to run (empty if manual check)
    PassPattern string   // Regex pattern indicating pass (empty if exit code 0 = pass)
    Category    string   // Category (specs, implementation, etc.)
}
```

### DriftDetection

```go
type DriftDetection struct {
    Field         string   // Which AGENTS.md field drifted
    Expected      string   // What AGENTS.md claimed
    Actual        string   // What was actually observed
    FixApplied    string   // What fix was applied (if any)
    Rationale     string   // Why the drift occurred
}
```

## Algorithm

### Read-Verify-Update Lifecycle

```
1. Check if AGENTS.md exists
   - If not: run bootstrap workflow
   - If yes: proceed to read

2. Read AGENTS.md
   - Parse file into AgentsMD struct
   - Extract build/test/lint commands
   - Extract spec/impl path patterns
   - Extract work tracking configuration
   - Extract quality criteria

3. Verify AGENTS.md Claims
   - For each command: run with --help or --version to verify it exists
   - For each path pattern: glob to verify files match
   - For work tracking: run query command to verify system is accessible
   - Record verification results

4. Execute Procedure Using AGENTS.md Data
   - Run build/test/lint commands as specified
   - Query work tracking system as specified
   - Evaluate quality criteria as specified
   - Monitor for drift (expected vs actual outcomes)

5. Detect Drift
   - Command failed but AGENTS.md claimed it should work
   - Path pattern matched no files but AGENTS.md claimed it should
   - Work tracking command failed but AGENTS.md claimed system was accessible
   - Quality criterion failed but AGENTS.md claimed it should pass

6. Update AGENTS.md (if drift detected)
   - Identify which section needs update
   - Modify section in-place (don't append)
   - Add inline rationale comment
   - Commit with descriptive message

7. Return Execution Results
   - Success/failure status
   - Drift detections (if any)
   - Updated AGENTS.md content (if modified)
```

### Bootstrap Workflow (No AGENTS.md)

```
1. Detect Build System
   - Check for package.json → npm/yarn
   - Check for Cargo.toml → cargo
   - Check for go.mod → go build
   - Check for Makefile → make
   - Check for build.sh → ./build.sh
   - Default: "Not required (interpreted language)"

2. Detect Test System
   - npm: "npm test"
   - cargo: "cargo test"
   - go: "go test ./..."
   - pytest: "pytest"
   - Default: "Manual verification (no automated tests)"

3. Detect Lint System
   - npm + eslint: "npm run lint"
   - cargo: "cargo clippy"
   - go: "golangci-lint run"
   - python + flake8: "flake8"
   - shellcheck for .sh files: "shellcheck <file>"
   - Default: None

4. Detect Spec Paths
   - Check for specs/ directory
   - Check for docs/ directory
   - Check for documentation/ directory
   - Check for README.md with spec sections
   - Default: "Not defined"

5. Detect Impl Paths
   - Check for src/ directory
   - Check for lib/ directory
   - Check for cmd/ directory
   - Check for main.* or index.* files
   - Default: "." (project root)

6. Detect Work Tracking
   - Check for .beads/ directory → beads
   - Check for .github/ directory → GitHub Issues
   - Check for TODO.md or TASKS.md → file-based
   - Default: "Not configured"

7. Generate AGENTS.md
   - Create file with detected values
   - Include rationale comments for each detection
   - Mark uncertain detections with "# Verify this"
   - Commit with message "Bootstrap AGENTS.md"

8. Verify Generated AGENTS.md
   - Run verification workflow on generated file
   - Update if any detections were incorrect
```

## Edge Cases

### AGENTS.md Doesn't Exist (Bootstrap)

First run on a new repository:

```
$ rooda build
INFO: AGENTS.md not found. Bootstrapping...
INFO: Detected build system: go build
INFO: Detected test system: go test ./...
INFO: Detected spec paths: specs/*.md
INFO: Detected impl paths: *.go, internal/**/*.go, cmd/**/*.go
INFO: Detected work tracking: beads (bd CLI)
INFO: Created AGENTS.md
INFO: Verifying generated AGENTS.md...
INFO: Verification passed. Proceeding with build procedure.
```

### Command in AGENTS.md Fails

AGENTS.md claims:
```markdown
**Test:** `npm test`
```

Actual behavior:
```
$ npm test
bash: npm: command not found
```

Drift detection:
```
WARN: Test command failed: npm test
WARN: Attempting to detect correct test command...
INFO: Detected: go test ./...
INFO: Updating AGENTS.md: Test command changed from 'npm test' to 'go test ./...'
INFO: Rationale: npm not installed, project uses Go
```

Updated AGENTS.md:
```markdown
**Test:** `go test ./...`  # Changed from npm test - npm not installed, project uses Go
```

### Path Pattern Matches No Files

AGENTS.md claims:
```markdown
**Specifications:** `specs/*.md`
```

Actual behavior:
```
$ ls specs/*.md
ls: specs/*.md: No such file or directory
```

Drift detection:
```
WARN: Spec path pattern matched no files: specs/*.md
WARN: Searching for spec files...
INFO: Found specs in: documentation/*.md
INFO: Updating AGENTS.md: Spec paths changed from 'specs/*.md' to 'documentation/*.md'
```

Updated AGENTS.md:
```markdown
**Specifications:** `documentation/*.md`  # Changed from specs/*.md - specs moved to documentation/
```

### Work Tracking Command Fails

AGENTS.md claims:
```markdown
**Query ready work:** `bd ready --json`
```

Actual behavior:
```
$ bd ready --json
bash: bd: command not found
```

Drift detection:
```
WARN: Work tracking query failed: bd ready --json
WARN: Checking for alternative work tracking systems...
INFO: Found .github/issues/ directory
INFO: Updating AGENTS.md: Work tracking changed from 'beads' to 'GitHub Issues'
```

Updated AGENTS.md:
```markdown
**System:** GitHub Issues  # Changed from beads - bd CLI not installed
**Query ready work:** `gh issue list --json number,title,labels --label ready`
```

### Quality Criterion Incorrect

AGENTS.md claims:
```markdown
- `./rooda.sh --version` executes without errors (PASS/FAIL)
```

Actual behavior:
```
$ ./rooda.sh --version
bash: ./rooda.sh: No such file or directory
```

Drift detection:
```
WARN: Quality criterion failed: ./rooda.sh --version executes without errors
WARN: File ./rooda.sh does not exist
INFO: Project has been rewritten in Go, rooda.sh is archived
INFO: Updating AGENTS.md: Remove outdated quality criterion
```

Updated AGENTS.md:
```markdown
# Removed: ./rooda.sh --version - v1 bash implementation archived, v2 is Go rewrite (not yet implemented)
```

### Empirical Verification Prevents Bad Execution

AGENTS.md claims:
```markdown
**Build:** `make build`
```

Verification:
```
$ make --version
make: command not found
```

Outcome:
```
WARN: Build command verification failed: make not installed
WARN: Skipping build step, updating AGENTS.md
INFO: Detected alternative: go build -o rooda ./cmd/rooda
INFO: Updated AGENTS.md with correct build command
```

No bad execution occurs — verification catches the issue before attempting the build.

## Dependencies

- **agents-md-format.md** — Defines AGENTS.md structure, required sections, field definitions
- **iteration-loop.md** — Procedures execute within iteration loops, may update AGENTS.md across iterations
- **error-handling.md** — Drift detection and recovery are forms of error handling

## Implementation Mapping

**Source files:**
- `internal/agents/parser.go` — Parse AGENTS.md into AgentsMD struct
- `internal/agents/verifier.go` — Empirical verification of AGENTS.md claims
- `internal/agents/updater.go` — In-place updates to AGENTS.md with rationale
- `internal/agents/bootstrap.go` — Bootstrap workflow for first-run case
- `internal/agents/drift.go` — Drift detection logic

**Related specs:**
- `agents-md-format.md` — AGENTS.md schema
- `procedures.md` — All procedures use this read-verify-update lifecycle
- `iteration-loop.md` — AGENTS.md may be updated across iterations

## Examples

### Successful Read-Verify-Execute

```
$ rooda build
INFO: Reading AGENTS.md...
INFO: Verifying build command: go build -o rooda ./cmd/rooda
INFO: Verification passed.
INFO: Verifying test command: go test ./...
INFO: Verification passed.
INFO: Verifying work tracking: bd ready --json
INFO: Verification passed.
INFO: Executing build procedure...
```

### Bootstrap from Scratch

```
$ rooda build
INFO: AGENTS.md not found. Bootstrapping...
INFO: Analyzing repository structure...
INFO: Detected build system: go build (found go.mod)
INFO: Detected test system: go test ./... (found *_test.go files)
INFO: Detected spec paths: specs/*.md (found specs/ directory)
INFO: Detected impl paths: **/*.go (found .go files)
INFO: Detected work tracking: beads (found .beads/ directory)
INFO: Created AGENTS.md with detected values
INFO: Verifying generated AGENTS.md...
INFO: Verification passed.
INFO: Executing build procedure...
```

### Drift Detection and Update

```
$ rooda build
INFO: Reading AGENTS.md...
INFO: Verifying test command: npm test
WARN: Test command failed: npm: command not found
INFO: Detecting alternative test command...
INFO: Found: go test ./...
INFO: Updating AGENTS.md: Test command changed from 'npm test' to 'go test ./...'
INFO: Rationale: npm not installed, project uses Go
INFO: Committing AGENTS.md update...
INFO: Executing build procedure with corrected test command...
```

### Path Pattern Correction

```
$ rooda audit-spec
INFO: Reading AGENTS.md...
INFO: Verifying spec paths: specs/*.md
WARN: Pattern matched 0 files
INFO: Searching for spec files...
INFO: Found 12 files matching: documentation/**/*.md
INFO: Updating AGENTS.md: Spec paths changed from 'specs/*.md' to 'documentation/**/*.md'
INFO: Rationale: Specs moved to documentation/ directory
INFO: Committing AGENTS.md update...
INFO: Executing audit-spec procedure with corrected paths...
```

### Quality Criterion Evaluation

```
$ rooda build
INFO: Reading AGENTS.md...
INFO: Evaluating quality criteria...
INFO: Running: shellcheck rooda.sh
WARN: Quality criterion failed: shellcheck rooda.sh (file not found)
INFO: rooda.sh is archived in archive/ directory
INFO: Updating AGENTS.md: Remove outdated quality criterion
INFO: Rationale: v1 bash implementation archived, v2 is Go rewrite
INFO: Committing AGENTS.md update...
INFO: Continuing build procedure...
```

## Notes

### Design Rationale

**Why read AGENTS.md at every procedure start?**
AGENTS.md may change between procedure invocations (manual edits, other procedures updating it). Reading fresh ensures procedures always use current data.

**Why empirical verification instead of trusting AGENTS.md?**
AGENTS.md can drift from reality (commands change, files move, tools uninstalled). Verification catches drift before it causes failures.

**Why update in-place instead of appending?**
In-place updates keep AGENTS.md concise and accurate. Appending creates a diary that grows unbounded and makes it hard to find current truth.

**Why inline rationale comments?**
Future readers (human or AI) need to understand why a value is what it is. Rationale prevents re-introducing old mistakes.

**Why bootstrap instead of requiring manual AGENTS.md creation?**
Zero-config startup — rooda should work on first run without requiring the user to write AGENTS.md manually.

**Why commit AGENTS.md updates?**
AGENTS.md is living documentation. Committing updates makes changes visible in git history and shareable across team members.

**Why treat AGENTS.md as source of truth but verify empirically?**
"Trust but verify" — AGENTS.md is the contract, but reality is the ultimate authority. Verification reconciles the two.

### Bootstrap Heuristics

**Build system detection priority:**
1. go.mod → `go build`
2. Cargo.toml → `cargo build`
3. package.json → `npm run build` or `yarn build`
4. Makefile → `make` or `make build`
5. build.sh → `./build.sh`
6. Default → "Not required"

**Test system detection priority:**
1. go.mod + *_test.go → `go test ./...`
2. Cargo.toml → `cargo test`
3. package.json + jest → `npm test`
4. pytest.ini or conftest.py → `pytest`
5. *_test.py files → `pytest`
6. Default → "Manual verification"

**Spec path detection priority:**
1. specs/ directory → `specs/*.md`
2. docs/ directory → `docs/**/*.md`
3. documentation/ directory → `documentation/**/*.md`
4. README.md with "## Specification" section → `README.md`
5. Default → "Not defined"

**Impl path detection priority:**
1. src/ directory → `src/**/*`
2. lib/ directory → `lib/**/*`
3. cmd/ directory → `cmd/**/*`
4. main.* or index.* files → `*.{go,rs,js,ts,py}`
5. Default → `.` (project root)

**Work tracking detection priority:**
1. .beads/ directory → beads (bd CLI)
2. .github/ directory → GitHub Issues (gh CLI)
3. TODO.md or TASKS.md → file-based
4. Default → "Not configured"

### Update Patterns

**Command correction:**
```markdown
**Test:** `go test ./...`  # Changed from npm test - npm not installed, project uses Go
```

**Path correction:**
```markdown
**Specifications:** `documentation/**/*.md`  # Changed from specs/*.md - specs moved to documentation/
```

**Quality criterion removal:**
```markdown
# Removed: ./rooda.sh --version - v1 bash implementation archived
```

**Work tracking correction:**
```markdown
**System:** GitHub Issues  # Changed from beads - bd CLI not installed
```

All updates include inline rationale explaining the change.
