# Procedures Specification

## Job to be Done

Define the 16 built-in procedures that ship as defaults — their fragment-based OODA phase compositions, iteration limits, and use cases. Each procedure is composed of reusable prompt fragments organized by OODA loop phase, allowing flexible procedure definitions through configuration rather than hardcoded prompts.

## Activities

1. Load procedure definition from configuration (built-in or custom)
2. Resolve fragment paths (builtin: prefix or relative to config directory)
3. Load fragment content from embedded resources or filesystem
4. Process templates when parameters are provided
5. Concatenate fragments in order for each OODA phase
6. Provide composed prompts to iteration loop

## Acceptance Criteria

- [ ] All 16 built-in procedures defined with fragment arrays for each OODA phase
- [ ] Fragment paths support builtin: prefix for embedded resources
- [ ] Fragment paths support relative paths from config file directory
- [ ] Fragments support inline content via content field
- [ ] Fragments support Go text/template syntax with parameters
- [ ] Fragment arrays concatenate in order to form complete phase prompts
- [ ] Procedure definitions support all configuration fields (display, summary, description, iteration settings, AI command overrides)
- [ ] Built-in fragments organized by OODA phase (observe/, orient/, decide/, act/)
- [ ] Fragment directory structure includes all 55 fragments (13 observe, 20 orient, 10 decide, 12 act)
- [ ] Template processing executes with provided parameters
- [ ] Path resolution validates fragment existence at config load time (fail fast)
- [ ] Schema validation ensures required fields present
- [ ] Fragment validation rejects fragments with both content and path specified

## Data Structures

See [configuration.md](configuration.md) for the Procedure struct definition. The procedures system extends this with fragment-based composition:

```go
type FragmentAction struct {
    Content    string                 // Inline prompt content (optional)
    Path       string                 // Path to fragment file (optional)
    Parameters map[string]interface{} // Template parameters (optional)
}

type ProcedureV2 struct {
    Display              string            // Human-readable name
    Summary              string            // Brief description
    Description          string            // Detailed explanation
    Observe              []FragmentAction  // Array of observe phase fragments
    Orient               []FragmentAction  // Array of orient phase fragments
    Decide               []FragmentAction  // Array of decide phase fragments
    Act                  []FragmentAction  // Array of act phase fragments
    IterationMode        IterationMode     // Override loop iteration mode
    DefaultMaxIterations *int              // Override loop max iterations
    IterationTimeout     *int              // Override loop timeout (seconds)
    MaxOutputBuffer      *int              // Override output buffer (bytes)
    AICmd                string            // Override AI command
    AICmdAlias           string            // Override AI command alias
}
```

## Algorithm

### Fragment Loading and Composition

```
function ComposePhasePrompt(fragments []FragmentAction, configDir string, embeddedFS embed.FS) -> (string, error):
    var parts []string
    
    for fragment in fragments:
        var content string
        
        // 1. Validate content vs path (exactly one required)
        hasContent = fragment.Content != ""
        hasPath = fragment.Path != ""
        
        if hasContent && hasPath:
            return error("fragment cannot specify both content and path")
        if !hasContent && !hasPath:
            return error("fragment must specify either content or path")
        
        // 2. Determine content source
        if hasContent:
            content = fragment.Content
        else:
            content, err = LoadFragment(fragment.Path, configDir, embeddedFS)
            if err:
                return error("failed to load fragment %s: %v", fragment.Path, err)
        
        
        // 3. Process template if parameters provided
        if len(fragment.Parameters) > 0:
            content, err = ProcessTemplate(content, fragment.Parameters)
            if err:
                return error("failed to process template: %v", err)
        
        // 4. Append to parts
        parts = append(parts, content)
    
    // 5. Concatenate with newlines
    return strings.Join(parts, "\n\n"), nil

function LoadFragment(path string, configDir string, embeddedFS embed.FS) -> (string, error):
    // Check for builtin: prefix
    if strings.HasPrefix(path, "builtin:"):
        embeddedPath = strings.TrimPrefix(path, "builtin:")
        content, err = embeddedFS.ReadFile(embeddedPath)
        if err:
            return error("embedded fragment not found: %s", path)
        return string(content), nil
    
    // Filesystem path, resolved relative to config directory
    fsPath = filepath.Join(configDir, path)
    content, err = os.ReadFile(fsPath)
    if err:
        return error("fragment file not found: %s (resolved to %s)", path, fsPath)
    return string(content), nil

function ProcessTemplate(content string, parameters map[string]interface{}) -> (string, error):
    tmpl, err = template.New("fragment").Parse(content)
    if err:
        return error("template parse error: %v", err)
    
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, parameters)
    if err:
        return error("template execution error: %v", err)
    
    return buf.String(), nil
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Fragment has neither content nor path | Error at config load: "fragment must specify either content or path" |
| Fragment has both content and path | Error at config load: "fragment cannot specify both content and path" |
| Fragment path with builtin: prefix not found | Error at config load: "embedded fragment not found: [path]" |
| Fragment path without prefix not found | Error at config load: "fragment file not found: [path] (resolved to [absolute])" |
| Template syntax error in fragment | Error: "template parse error: [details]" |
| Template parameter missing | Template executes with zero value for missing parameter (Go default) |
| Template execution error (e.g., range over nil) | Error: "template execution error: [details]" |
| Empty fragment array for OODA phase | Empty string returned for that phase |
| Fragment parameters provided but no template syntax | Parameters ignored, content returned as-is |
| Relative path in workspace config | Resolved relative to workspace config directory (./) |
| Relative path in global config | Resolved relative to global config directory |
| Fragment concatenation | Joined with double newline (\n\n) separator |

## Dependencies

- **configuration.md** — Defines Procedure struct and config loading
- **prompt-composition.md** — Consumes composed phase prompts
- **iteration-loop.md** — Executes procedures with composed prompts
- **Go text/template** — Template processing for parameterized fragments
- **go:embed** — Embedding built-in fragment files in binary

## Implementation Mapping

**Source files:**
- `internal/procedures/procedures.go` — Procedure loading and fragment composition
- `internal/procedures/fragments.go` — Fragment loading and template processing
- `internal/procedures/builtin.go` — Built-in procedure definitions
- `fragments/` — Built-in fragment files (embedded via go:embed)
  - `fragments/observe/` — 13 observe phase fragments
  - `fragments/orient/` — 20 orient phase fragments
  - `fragments/decide/` — 10 decide phase fragments
  - `fragments/act/` — 12 act phase fragments

**Related specs:**
- `configuration.md` — Procedure configuration schema
- `prompt-composition.md` — Prompt assembly from composed phases
- `iteration-loop.md` — Procedure execution

## Examples

### Example 1: Simple Built-in Procedure

**Procedure definition:**
```yaml
procedures:
  agents-sync:
    display: "Agents Sync"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/scan_repo_structure.md"
```

**Fragment composition:**
```
1. Load builtin:fragments/observe/read_agents_md.md from embedded FS
2. Load builtin:fragments/observe/scan_repo_structure.md from embedded FS
3. Concatenate with \n\n separator
4. Return composed observe prompt
```

**Verification:**
- Both fragments loaded from embedded resources
- No template processing (no parameters)
- Concatenated in order

### Example 2: Custom Procedure with Template

**Procedure definition:**
```yaml
procedures:
  custom-audit:
    observe:
      - path: "fragments/observe/read_files.md"
        parameters:
          file_type: "specification"
          paths:
            - "specs/api.md"
            - "specs/database.md"
          include_tests: true
```

**Fragment file** (`fragments/observe/read_files.md`):
```markdown
Read the following {{.file_type}} files from the repository:

{{range .paths}}
- {{.}}
{{end}}

{{if .include_tests}}
Also include any associated test files.
{{end}}
```

**Rendered output:**
```markdown
Read the following specification files from the repository:

- specs/api.md
- specs/database.md

Also include any associated test files.
```

**Verification:**
- Fragment loaded from filesystem relative to config directory
- Template processed with provided parameters
- Go text/template syntax executed correctly

### Example 3: Inline Content Fragment

**Procedure definition:**
```yaml
procedures:
  quick-check:
    observe:
      - content: "Read the current git status and list uncommitted changes."
      - path: "builtin:fragments/observe/read_specs.md"
```

**Fragment composition:**
```
1. Use inline content directly (no file loading)
2. Load builtin:fragments/observe/read_specs.md from embedded FS
3. Concatenate with \n\n separator
```

**Verification:**
- Inline content used without file loading
- Mixed with file-based fragment
- Both concatenated in order

## Notes

**Design Rationale — Fragment-Based Composition:**

The fragment-based system replaces single-file prompts per OODA phase with composable arrays of fragments. This provides:

1. **Reusability** — Common fragments (read_agents_md.md, read_specs.md) used across multiple procedures
2. **Flexibility** — Procedures compose different combinations of fragments without duplication
3. **Maintainability** — Update one fragment, all procedures using it benefit
4. **Extensibility** — Add new fragments without modifying existing procedures

**Why Arrays Instead of Single Files:**

v1 used single prompt files per OODA phase (observe.md, orient.md, decide.md, act.md). This led to:
- Duplication across procedures (many procedures read AGENTS.md first)
- Large monolithic prompt files
- Difficulty reusing prompt components

Fragment arrays solve this by allowing procedures to compose from shared building blocks.

**Template System Choice:**

Go's text/template provides:
- Built-in to Go standard library (no external dependency)
- Familiar syntax for Go developers
- Safe execution (no arbitrary code execution)
- Good error messages for template syntax issues

**Path Resolution:**

The builtin: prefix distinguishes embedded resources from filesystem paths without ambiguity:
- No naming conflicts (builtin: never appears in filesystem paths)
- Clear intent (builtin: = shipped with binary, no prefix = user-provided)
- Simple implementation (string prefix check)

**Fragment Organization:**

Fragments organized by OODA phase (observe/, orient/, decide/, act/) because:
- Clear mental model (fragments map to phases)
- Easy to find relevant fragments
- Matches procedure structure
- Prevents cross-phase confusion

The procedures system uses a fragment-based composition approach where each procedure is composed of reusable prompt fragments organized by OODA loop phase (Observe, Orient, Decide, Act).

## Fragment Directory Structure

```
fragments/
├── observe/
│   ├── read_agents_md.md              # Load and parse AGENTS.md configuration
│   ├── scan_repo_structure.md          # Examine directory structure and files
│   ├── detect_build_system.md          # Identify build tools (go.mod, package.json, etc.)
│   ├── detect_work_tracking.md         # Identify work tracking system (.beads/, .github/)
│   ├── verify_commands.md              # Test that commands from AGENTS.md work
│   ├── query_work_tracking.md          # Fetch ready work items
│   ├── read_specs.md                   # Load specification files
│   ├── read_impl.md                    # Load implementation files
│   ├── read_task_input.md              # Load task description
│   ├── read_draft_plan.md              # Load draft plan file
│   ├── read_task_details.md            # Load specific task from work tracking
│   ├── run_tests.md                    # Execute test commands
│   └── run_lints.md                    # Execute lint commands
├── orient/
│   ├── compare_detected_vs_documented.md        # Find drift between actual and documented
│   ├── compare_documented_vs_actual.md          # Find drift in AGENTS.md
│   ├── identify_drift.md                        # Categorize inconsistencies
│   ├── evaluate_against_quality_criteria.md     # Check PASS/FAIL criteria
│   ├── understand_task_requirements.md          # Parse task into requirements
│   ├── understand_feature_requirements.md       # Parse feature requirements
│   ├── understand_bug_root_cause.md             # Analyze bug cause
│   ├── search_codebase.md                       # Find relevant code sections
│   ├── identify_affected_files.md               # Determine what needs to change
│   ├── identify_affected_specs.md               # Determine which specs need changes
│   ├── identify_affected_code.md                # Determine which code needs changes
│   ├── identify_spec_deficiencies.md            # Find gaps in specs
│   ├── parse_plan_tasks.md                      # Extract tasks from draft plan
│   ├── map_to_work_tracking_format.md           # Convert plan to work tracking format
│   ├── identify_specified_but_not_implemented.md # Gap analysis (specs → impl)
│   ├── identify_implemented_but_not_specified.md # Gap analysis (impl → specs)
│   ├── identify_structural_issues.md            # Find spec/code structure problems
│   ├── identify_duplication.md                  # Find duplicated content
│   ├── identify_code_smells.md                  # Find code quality issues
│   ├── identify_complexity_issues.md            # Find overly complex code
│   └── identify_maintenance_needs.md            # Find maintenance work
├── decide/
│   ├── determine_sections_to_update.md          # What to change in AGENTS.md
│   ├── check_if_blocked.md                      # Can we proceed? (emit FAILURE if not)
│   ├── pick_task.md                             # Select work item from work tracking
│   ├── plan_implementation_approach.md          # How to implement
│   ├── break_down_into_tasks.md                 # Decompose work into tasks
│   ├── prioritize_tasks.md                      # Order by impact/dependency
│   ├── prioritize_findings.md                   # Order audit findings
│   ├── prioritize_gaps_by_impact.md             # Order gap analysis findings
│   ├── identify_issues.md                       # List problems found
│   ├── categorize_drift_severity.md             # Rank drift items
│   └── determine_import_strategy.md             # How to import plan to work tracking
└── act/
    ├── write_agents_md.md                       # Update AGENTS.md file
    ├── write_audit_report.md                    # Create audit report
    ├── write_gap_report.md                      # Create gap analysis report
    ├── write_draft_plan.md                      # Create draft plan
    ├── modify_files.md                          # Edit specs or implementation files
    ├── commit_changes.md                        # Git commit with message
    ├── update_work_tracking.md                  # Mark tasks complete
    ├── update_draft_plan_status.md              # Update draft plan status
    ├── create_work_items.md                     # Import plan to work tracking
    ├── run_tests.md                             # Execute tests (verification)
    ├── emit_success.md                          # Output <promise>SUCCESS</promise> signal
    └── emit_failure.md                          # Output <promise>FAILURE</promise> signal
```

## Procedures Configuration Schema

```yaml
procedures:
  <procedure-name>:                    # string: Unique identifier for the procedure
    display: string                    # Human-readable name for the procedure
    summary: string                    # Brief description of what the procedure does
    description: string                # Detailed explanation of the procedure's purpose
    observe:                           # Array of actions (concatenated to form full prompt)
      - content: string                # Inline prompt content (optional)
        path: string                   # Path to prompt file (optional)
        parameters:                    # Template parameters (optional)
          <param-name>: <param-value>
    orient:                            # Array of actions (concatenated to form full prompt)
      - content: string                # Inline prompt content (optional)
        path: string                   # Path to prompt file (optional)
        parameters:                    # Template parameters (optional)
          <param-name>: <param-value>
    decide:                            # Array of actions (concatenated to form full prompt)
      - content: string                # Inline prompt content (optional)
        path: string                   # Path to prompt file (optional)
        parameters:                    # Template parameters (optional)
          <param-name>: <param-value>
    act:                               # Array of actions (concatenated to form full prompt)
      - content: string                # Inline prompt content (optional)
        path: string                   # Path to prompt file (optional)
        parameters:                    # Template parameters (optional)
          <param-name>: <param-value>
    iteration_mode: string             # "max-iterations" or "unlimited" (optional)
    default_max_iterations: integer    # Maximum number of OODA loop iterations (optional)
    iteration_timeout: integer         # Timeout per iteration in seconds (optional)
    max_output_buffer: integer         # Maximum output buffer size in bytes (optional)
    ai_cmd: string                     # AI command to use (optional)
    ai_cmd_alias: string               # Model configuration alias (optional)
```

## Fragment Composition Rules

### Path Resolution

- **Builtin fragments**: Use `builtin:` prefix (e.g., `builtin:fragments/observe/read_agents_md.md`)
- **Custom fragments**: Use relative path from config file directory (e.g., `fragments/observe/custom.md`)

### Content vs Path

Each fragment action must specify exactly one of:
- `content`: Inline prompt text
- `path`: Path to a fragment file (with optional `parameters` for templates)

Specifying both `content` and `path` is an error and will be rejected at config load time.

### Array Concatenation

Fragments within each OODA phase are concatenated in order to form the complete prompt for that phase.

## Template System

Fragments support Go's text/template syntax. When parameters are supplied, the file at path is processed as a template.

### Template Example

**Fragment file** (`fragments/observe/read_files.md`):
```markdown
Read the following {{.file_type}} files from the repository:

{{range .paths}}
- {{.}}
{{end}}

{{if .include_tests}}
Also include any associated test files.
{{end}}
```

**Procedure configuration**:
```yaml
procedures:
  custom-audit:
    observe:
      - path: "fragments/observe/read_files.md"
        parameters:
          file_type: "specification"
          paths:
            - "specs/api.md"
            - "specs/database.md"
          include_tests: true
```

**Rendered output**:
```markdown
Read the following specification files from the repository:

- specs/api.md
- specs/database.md

Also include any associated test files.
```

## Built-in Procedures

The system includes 16 built-in procedures organized into three categories:

### Direct Action Procedures

1. **agents-sync**: Synchronize AGENTS.md with actual repository state
2. **build**: Implement a task from work tracking
3. **publish-plan**: Import draft plan into work tracking system

### Audit Procedures

4. **audit-spec**: Audit specification files for quality issues
5. **audit-impl**: Audit implementation files for quality issues
6. **audit-agents**: Audit AGENTS.md for accuracy and completeness
7. **audit-spec-to-impl**: Find specifications not implemented in code
8. **audit-impl-to-spec**: Find implementation not covered by specifications

### Planning Procedures

9. **draft-plan-spec-feat**: Create plan for new specification feature
10. **draft-plan-spec-fix**: Create plan for specification bug fix
11. **draft-plan-spec-refactor**: Create plan for specification refactoring
12. **draft-plan-spec-chore**: Create plan for specification maintenance tasks
13. **draft-plan-impl-feat**: Create plan for new implementation feature
14. **draft-plan-impl-fix**: Create plan for implementation bug fix
15. **draft-plan-impl-refactor**: Create plan for code refactoring
16. **draft-plan-impl-chore**: Create plan for code maintenance tasks

## SUCCESS Criteria by Procedure Type

Each procedure type has specific criteria for when to emit `<promise>SUCCESS</promise>` versus continuing iteration.

### Direct Action Procedures

**agents-sync:**
- Emit `<promise>SUCCESS</promise>` when: AGENTS.md has been updated to match repository state and changes are committed
- Continue iterating when: Drift detected but not yet resolved, or verification failed
- Emit `<promise>FAILURE</promise>` when: Cannot determine repository state or AGENTS.md is corrupted

**build:**
- Emit `<promise>SUCCESS</promise>` when: No ready work remains in work tracking system (all tasks completed)
- Continue iterating when: Ready work exists and can be implemented
- Emit `<promise>FAILURE</promise>` when: Blocked by missing dependencies, broken tests that can't be fixed, or corrupted work tracking state

**publish-plan:**
- Emit `<promise>SUCCESS</promise>` when: All tasks from draft plan have been created in work tracking system
- Continue iterating when: Draft plan exists but work items not yet created
- Emit `<promise>FAILURE</promise>` when: Draft plan is missing, malformed, or work tracking system is unavailable

### Audit Procedures

**audit-spec, audit-impl, audit-agents:**
- Emit `<promise>SUCCESS</promise>` when: Audit report has been generated and is complete
- Continue iterating when: Report generation in progress or incomplete
- Emit `<promise>FAILURE</promise>` when: Cannot read target files or quality criteria are undefined

**audit-spec-to-impl, audit-impl-to-spec:**
- Emit `<promise>SUCCESS</promise>` when: Gap report has been generated and is complete
- Continue iterating when: Gap analysis in progress or incomplete
- Emit `<promise>FAILURE</promise>` when: Cannot read specs or implementation files

### Planning Procedures

**draft-plan-* (all 8 variants):**
- Emit `<promise>SUCCESS</promise>` when: Draft plan file has been created and is complete
- Continue iterating when: Plan generation in progress or incomplete
- Emit `<promise>FAILURE</promise>` when: Task input is missing, malformed, or requirements are unclear

## Example Outputs

### agents-sync Procedure

**Successful completion:**
```
[Iteration 1]
Detected drift in AGENTS.md:
- Build command documented as "go build" but actual is "go build -o bin/rooda ./cmd/rooda"
- Test command missing from documentation

Updated AGENTS.md with correct commands.
Verified all commands execute successfully.
Committed changes.

<promise>SUCCESS</promise>

AGENTS.md synchronized with repository state:
- Updated build command
- Added test command
- Verified all documented commands work
```

### build Procedure

**Successful completion (no work remaining):**
```
[Iteration 3]
Implemented task ralph-wiggum-ooda-abc123: Add user authentication
- Created internal/auth/auth.go
- Added tests in internal/auth/auth_test.go
- All tests passing
- Updated work tracking status to complete

Queried work tracking for ready work: no tasks found.

<promise>SUCCESS</promise>

All ready work completed:
- Implemented 1 task
- All tests passing
- No ready tasks remaining in work tracking
```

**Blocked (cannot proceed):**
```
[Iteration 2]
Attempted to implement task ralph-wiggum-ooda-xyz789: Add OAuth2 integration

<promise>FAILURE</promise>

Cannot proceed: Missing authentication module specification. The OAuth2 integration requires a detailed spec defining token refresh behavior and error handling patterns.

Next steps:
1. Create specs/auth-oauth2.md with token lifecycle specification
2. Define error handling patterns for expired tokens
3. Document refresh token rotation policy
```

### audit-spec Procedure

**Successful completion:**
```
[Iteration 1]
Reviewed 12 specification files against quality criteria:
- All specs have "Job to be Done" section: PASS
- All specs have "Acceptance Criteria" section: PASS
- All specs have "Examples" section: PASS
- No broken cross-references: FAIL (3 broken links found)

Generated audit report at docs/audit-2024-01-15.md

<promise>SUCCESS</promise>

Audit completed:
- Reviewed 12 specification files
- Generated audit report at docs/audit-2024-01-15.md
- Found 3 issues requiring attention
```

### draft-plan-impl-feat Procedure

**Successful completion:**
```
[Iteration 2]
Analyzed feature requirements for user authentication.
Identified affected code files:
- internal/auth/ (new package)
- cmd/rooda/main.go (integration point)
- internal/config/config.go (auth config)

Created draft plan at PLAN.md with 8 tasks:
1. Create auth package structure
2. Implement JWT token generation
3. Implement token validation
4. Add auth middleware
5. Integrate with main.go
6. Add configuration options
7. Write unit tests
8. Write integration tests

<promise>SUCCESS</promise>

Draft plan completed:
- Created plan at PLAN.md
- Broken down into 8 actionable tasks
- Ready for import to work tracking
```

## Built-in Procedures Configuration

```yml
procedures:
  agents-sync:
    display: "Agents Sync"
    summary: "Synchronize AGENTS.md with actual repository state"
    description: "Detects drift between documented and actual repository configuration, then updates AGENTS.md to match reality"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/scan_repo_structure.md"
      - path: "builtin:fragments/observe/detect_build_system.md"
      - path: "builtin:fragments/observe/detect_work_tracking.md"
    orient:
      - path: "builtin:fragments/orient/compare_detected_vs_documented.md"
      - path: "builtin:fragments/orient/identify_drift.md"
    decide:
      - path: "builtin:fragments/decide/determine_sections_to_update.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_agents_md.md"
      - path: "builtin:fragments/act/commit_changes.md"
      - path: "builtin:fragments/act/emit_success.md"

  build:
    display: "Build"
    summary: "Implement a task from work tracking"
    description: "Picks a ready task, implements it, runs tests, and marks it complete"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/query_work_tracking.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
      - path: "builtin:fragments/observe/read_task_details.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
      - path: "builtin:fragments/orient/search_codebase.md"
      - path: "builtin:fragments/orient/identify_affected_files.md"
    decide:
      - path: "builtin:fragments/decide/pick_task.md"
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
      - path: "builtin:fragments/act/run_tests.md"
      - path: "builtin:fragments/act/update_work_tracking.md"
      - path: "builtin:fragments/act/commit_changes.md"
      - path: "builtin:fragments/act/emit_success.md"

  publish-plan:
    display: "Publish Plan"
    summary: "Import draft plan into work tracking system"
    description: "Takes a draft plan and creates work items in the configured work tracking system"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_draft_plan.md"
      - path: "builtin:fragments/observe/query_work_tracking.md"
    orient:
      - path: "builtin:fragments/orient/parse_plan_tasks.md"
      - path: "builtin:fragments/orient/map_to_work_tracking_format.md"
    decide:
      - path: "builtin:fragments/decide/determine_import_strategy.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/create_work_items.md"
      - path: "builtin:fragments/act/update_draft_plan_status.md"
      - path: "builtin:fragments/act/emit_success.md"

  audit-spec:
    display: "Audit Specifications"
    summary: "Audit specification files for quality issues"
    description: "Reviews spec files against quality criteria and generates audit report"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/evaluate_against_quality_criteria.md"
    decide:
      - path: "builtin:fragments/decide/identify_issues.md"
      - path: "builtin:fragments/decide/prioritize_findings.md"
    act:
      - path: "builtin:fragments/act/write_audit_report.md"
      - path: "builtin:fragments/act/emit_success.md"

  audit-impl:
    display: "Audit Implementation"
    summary: "Audit implementation files for quality issues"
    description: "Reviews implementation files, runs tests and lints, generates audit report"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_impl.md"
      - path: "builtin:fragments/observe/run_tests.md"
      - path: "builtin:fragments/observe/run_lints.md"
    orient:
      - path: "builtin:fragments/orient/evaluate_against_quality_criteria.md"
    decide:
      - path: "builtin:fragments/decide/identify_issues.md"
      - path: "builtin:fragments/decide/prioritize_findings.md"
    act:
      - path: "builtin:fragments/act/write_audit_report.md"
      - path: "builtin:fragments/act/emit_success.md"

  audit-agents:
    display: "Audit Agents Configuration"
    summary: "Audit AGENTS.md for accuracy and completeness"
    description: "Verifies AGENTS.md matches repository state and commands work correctly"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/scan_repo_structure.md"
      - path: "builtin:fragments/observe/detect_build_system.md"
      - path: "builtin:fragments/observe/verify_commands.md"
    orient:
      - path: "builtin:fragments/orient/compare_documented_vs_actual.md"
      - path: "builtin:fragments/orient/identify_drift.md"
    decide:
      - path: "builtin:fragments/decide/categorize_drift_severity.md"
    act:
      - path: "builtin:fragments/act/write_audit_report.md"
      - path: "builtin:fragments/act/emit_success.md"

  audit-spec-to-impl:
    display: "Audit Spec to Implementation Gap"
    summary: "Find specifications not implemented in code"
    description: "Identifies features specified but not yet implemented"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/identify_specified_but_not_implemented.md"
    decide:
      - path: "builtin:fragments/decide/prioritize_gaps_by_impact.md"
    act:
      - path: "builtin:fragments/act/write_gap_report.md"
      - path: "builtin:fragments/act/emit_success.md"

  audit-impl-to-spec:
    display: "Audit Implementation to Spec Gap"
    summary: "Find implementation not covered by specifications"
    description: "Identifies code that exists but is not documented in specifications"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_impl.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/identify_implemented_but_not_specified.md"
    decide:
      - path: "builtin:fragments/decide/prioritize_gaps_by_impact.md"
    act:
      - path: "builtin:fragments/act/write_gap_report.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-spec-feat:
    display: "Draft Plan: Spec Feature"
    summary: "Create plan for new specification feature"
    description: "Analyzes feature requirements and creates implementation plan focused on specifications"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/understand_feature_requirements.md"
      - path: "builtin:fragments/orient/identify_affected_specs.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-spec-fix:
    display: "Draft Plan: Spec Bug Fix"
    summary: "Create plan for specification bug fix"
    description: "Analyzes bug root cause and creates fix plan focused on specifications"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/understand_bug_root_cause.md"
      - path: "builtin:fragments/orient/identify_spec_deficiencies.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-spec-refactor:
    display: "Draft Plan: Spec Refactor"
    summary: "Create plan for specification refactoring"
    description: "Identifies structural issues in specs and creates refactoring plan"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/identify_structural_issues.md"
      - path: "builtin:fragments/orient/identify_duplication.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-spec-chore:
    display: "Draft Plan: Spec Maintenance"
    summary: "Create plan for specification maintenance tasks"
    description: "Identifies maintenance needs in specs and creates chore plan"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/identify_maintenance_needs.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-impl-feat:
    display: "Draft Plan: Implementation Feature"
    summary: "Create plan for new implementation feature"
    description: "Analyzes feature requirements and creates implementation plan focused on code"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/understand_feature_requirements.md"
      - path: "builtin:fragments/orient/identify_affected_code.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-impl-fix:
    display: "Draft Plan: Implementation Bug Fix"
    summary: "Create plan for implementation bug fix"
    description: "Analyzes bug root cause and creates fix plan focused on code"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_specs.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/understand_bug_root_cause.md"
      - path: "builtin:fragments/orient/identify_affected_code.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-impl-refactor:
    display: "Draft Plan: Implementation Refactor"
    summary: "Create plan for code refactoring"
    description: "Identifies code smells and complexity issues, creates refactoring plan"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/identify_code_smells.md"
      - path: "builtin:fragments/orient/identify_complexity_issues.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"

  draft-plan-impl-chore:
    display: "Draft Plan: Implementation Maintenance"
    summary: "Create plan for code maintenance tasks"
    description: "Identifies maintenance needs in code and creates chore plan"
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_task_input.md"
      - path: "builtin:fragments/observe/read_impl.md"
    orient:
      - path: "builtin:fragments/orient/identify_maintenance_needs.md"
    decide:
      - path: "builtin:fragments/decide/break_down_into_tasks.md"
      - path: "builtin:fragments/decide/prioritize_tasks.md"
      - path: "builtin:fragments/decide/check_if_blocked.md"
    act:
      - path: "builtin:fragments/act/write_draft_plan.md"
      - path: "builtin:fragments/act/emit_success.md"
```

## Implementation Requirements

### Fragment Loading

The system must:
1. Resolve `builtin:` prefix to embedded fragment resources
2. Resolve relative paths to repository-local fragments
3. Load fragment content from files
4. Process templates when parameters are provided
5. Concatenate fragments in order for each OODA phase

### Template Processing

When a fragment has parameters:
1. Load the file content
2. Parse as Go text/template
3. Execute template with provided parameters
4. Return rendered output

### Validation

The system must validate at config load time (fail fast):
1. Each fragment specifies exactly one of `content` or `path` (not both, not neither)
2. All referenced fragment paths exist (both builtin: and filesystem paths)
3. OODA phase arrays are properly structured

The system validates at template execution time:
1. Template syntax is valid (parse errors)
2. Template execution succeeds (execution errors)

Note: Template parameter validation is deferred to future versions. For v2, missing parameters use Go template default behavior (zero values). Fragment authors should use defensive template patterns (e.g., `{{if .param}}...{{end}}`) to handle missing parameters gracefully.
