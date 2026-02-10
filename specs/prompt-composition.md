# Prompt Composition

## Job to be Done

Assemble prompts from fragment arrays for each OODA phase (observe, orient, decide, act) and optional user-provided context into a single prompt that can be piped to an AI CLI tool, supporting embedded defaults, user-provided custom fragments, inline content, and Go text/template processing with parameters.

## Activities

1. **Assemble preamble** — Create procedure execution preamble with agent role, iteration context, and success signaling instructions
2. **Resolve fragment paths** — For each fragment in OODA phase arrays, determine whether to use embedded resources (builtin: prefix), filesystem paths (relative to config directory), or inline content
3. **Load fragment content** — Load content from embedded resources, filesystem, or use inline content directly
4. **Process templates** — When parameters are provided, execute Go text/template processing on fragment content
5. **Concatenate fragments** — Join fragments within each phase with double newlines
6. **Inject user context** — If --context flag provided, insert user-supplied text after preamble
7. **Format with section markers** — Wrap each phase with clear delimiters for readability and debugging
8. **Validate completeness** — Ensure all fragments are loadable and templates are valid at config load time (fail fast)

## Acceptance Criteria

- [ ] Assembles preamble with procedure name, iteration context, and success signaling instructions
- [ ] Preamble includes current iteration number and max iterations (or "unlimited")
- [ ] Preamble instructs agent to emit `<promise>SUCCESS</promise>` when job complete
- [ ] Preamble instructs agent to emit `<promise>FAILURE</promise>` when blocked
- [ ] Assembles prompts from fragment arrays for each OODA phase (observe, orient, decide, act)
- [ ] Supports embedded fragments via builtin: prefix (e.g., builtin:fragments/observe/read_agents_md.md)
- [ ] Supports filesystem fragments via relative paths (e.g., fragments/observe/custom.md)
- [ ] Supports inline content via content field in fragment actions
- [ ] Filesystem paths resolved relative to config file directory, not current working directory
- [ ] Processes Go text/template syntax when parameters are provided
- [ ] Concatenates fragments within each phase with double newlines (\n\n)
- [ ] Injects user-provided context when --context flag is supplied (after preamble, before OODA phases)
- [ ] Wraps user context with "=== CONTEXT ===" section marker
- [ ] Wraps each phase with section markers (e.g., "=== OBSERVE ===", "=== ORIENT ===")
- [ ] When context is from file, shows "Source: <path>" followed by file content
- [ ] When context is inline, shows content only (no Source line)
- [ ] Validates all fragment paths at config load time (fail fast)
- [ ] Validates template syntax at config load time
- [ ] Returns error if fragment has both content and path specified
- [ ] Returns error if fragment has neither content nor path specified
- [ ] Returns error if any fragment file is missing
- [ ] Preserves markdown formatting from source fragments
- [ ] Normalizes trailing newlines (trim then add consistent spacing)
- [ ] Handles multi-line context injection without breaking prompt structure
- [ ] Empty fragment arrays result in empty phase content (not an error)

## Data Structures

### Procedure Definition (from configuration.md and procedures.md)

```yaml
procedures:
  build:
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/query_work_tracking.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
      - content: "Focus on maintainability and test coverage."
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
      - path: "builtin:fragments/act/run_tests.md"
    default_max_iterations: 5
```

### Fragment Action Structure

```go
type FragmentAction struct {
    Content    string                 // Inline prompt content (optional)
    Path       string                 // Path to fragment file (optional)
    Parameters map[string]interface{} // Template parameters (optional)
}
```

Each fragment action must specify exactly one of:
- `content`: Inline prompt text
- `path`: Path to a fragment file (with optional `parameters` for templates)

Specifying both or neither is an error.

### Assembled Prompt Structure

```
═══════════════════════════════════════════════════════════════
ROODA PROCEDURE EXECUTION
═══════════════════════════════════════════════════════════════

Procedure: <procedure-name>

Your Role:
You are an AI coding agent executing a structured OODA loop procedure.
This is NOT a template or example - this is an EXECUTABLE PROCEDURE.
You must complete all phases and produce concrete outputs.

Success Signaling:
- When you complete all tasks successfully, output: <promise>SUCCESS</promise>
- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>
- Explanations should come AFTER the signal, not embedded in the tag
- The loop orchestrator uses these signals to determine iteration outcome.

═══════════════════════════════════════════════════════════════
CONTEXT
═══════════════════════════════════════════════════════════════
Source: ./path/to/file.md

[File content - if context from file]

═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
[Content from observe phase fragments]

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from orient phase fragments]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from decide phase fragments]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from act phase fragments]
```

Note: When context is inline (not from file), the "Source:" line is omitted. When no context is provided, the entire CONTEXT section is omitted.

## Procedure Execution Preamble

The preamble wraps the composed prompt with explicit execution context that frames the agent's role and responsibilities. It appears at the very beginning of every assembled prompt, before any user context or OODA phases.

### Purpose

The preamble serves three critical functions:

1. **Agent Role Definition** — Establishes that the agent is executing a structured OODA loop procedure, not having a freeform conversation
2. **Iteration Context** — Provides awareness of progress (iteration N of M) so the agent can gauge how much work remains
3. **Success Signaling** — Instructs the agent how to communicate completion or blockage through `<promise>` output markers

### Format

```
═══════════════════════════════════════════════════════════════
ROODA PROCEDURE EXECUTION
═══════════════════════════════════════════════════════════════

Procedure: <procedure-name>

Your Role:
You are an AI coding agent executing a structured OODA loop procedure.
This is NOT a template or example - this is an EXECUTABLE PROCEDURE.
You must complete all phases and produce concrete outputs.

Success Signaling:
- When you complete all tasks successfully, output: <promise>SUCCESS</promise>
- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>
- Explanations should come AFTER the signal, not embedded in the tag
- The loop orchestrator uses these signals to determine iteration outcome.
```

### Template Variables

- `<procedure-name>` — The name of the procedure being executed (e.g., "build", "agents-sync")

### Success Signal Instructions

The preamble explicitly instructs the agent to emit promise markers:

- `<promise>SUCCESS</promise>` — Job is complete, loop should terminate with success status
- `<promise>FAILURE</promise>` — Agent is blocked and cannot make further progress, loop should count this as a failure

These markers are scanned by the iteration loop to determine outcome (see iteration-loop.md Iteration Outcome Matrix).

### Design Rationale

**Why a preamble?**
- Agents often treat prompts as templates or documentation rather than executable procedures
- Explicit framing as "EXECUTABLE PROCEDURE" improves agent recognition of the task
- Clear role definition helps agents understand they must produce concrete outputs

**Why enhanced section markers with double lines?**
- Visual prominence: Double-line borders make phase transitions unmissable
- Directive language: Phase descriptions tell agents what to do, not just what the phase is
- Reduced ambiguity: "Execute these observation tasks" is more directive than "OBSERVE"
- Better agent recognition: Enhanced markers help agents understand prompts as procedures, not templates

**Why include phase descriptions in markers?**
- Reinforces the OODA structure at each phase transition
- Provides a mental model for how to process the subsequent content
- Reduces confusion about what each phase marker means
- Makes each phase self-documenting

**Why explicit success signaling instructions?**
- Agents don't inherently know to emit `<promise>` markers
- Clear instructions increase signal emission rate
- Reduces iterations where agent completes work but doesn't signal completion

**Why emphasize "NOT a template"?**
- Agents frequently misinterpret structured prompts as examples to follow rather than tasks to execute
- Explicit negation ("This is NOT a template") helps overcome this pattern
- Improves agent recognition that they must produce actual outputs, not discuss what they would do

## Algorithm

```
function AssemblePrompt(procedure, contextValues []string, configDir):
    prompt = ""
    
    // 1. Assemble preamble
    prompt += AssemblePreamble(procedure.Name)
    prompt += "\n\n"
    
    // 2. Inject user context if provided
    if len(contextValues) > 0:
        prompt += "═══════════════════════════════════════════════════════════════\n"
        prompt += "CONTEXT\n"
        prompt += "═══════════════════════════════════════════════════════════════\n"
        
        for _, contextValue in contextValues:
            // Check if context is a file path (file existence heuristic)
            if fileExists(contextValue):
                // Read file content
                content, err = os.ReadFile(contextValue)
                if err:
                    return error("failed to read context file %s: %v", contextValue, err)
                
                // Add source path and content
                prompt += "Source: " + contextValue + "\n\n"
                prompt += string(content) + "\n\n"
            else:
                // Inline content - no source line
                prompt += contextValue + "\n\n"
    
    // 3. Process each OODA phase in order
    phaseDescriptions = map[string]string{
        "observe": "Execute these observation tasks to gather information.",
        "orient": "Analyze the information you gathered and form your understanding.",
        "decide": "Make decisions about what actions to take.",
        "act": "Execute the actions you decided on. Modify files, run commands, commit changes.",
    }
    
    phaseNumbers = map[string]int{
        "observe": 1,
        "orient": 2,
        "decide": 3,
        "act": 4,
    }
    
    for phase in [observe, orient, decide, act]:
        fragmentActions = procedure[phase]  // Array of FragmentAction
        
        // Compose phase prompt from fragment array
        phaseContent, err = ComposePhasePrompt(fragmentActions, configDir)
        if err:
            return error("failed to compose %s phase: %v", phase, err)
        
        // Add section marker and content (normalize trailing newlines)
        if strings.TrimSpace(phaseContent) != "":
            prompt += "═══════════════════════════════════════════════════════════════\n"
            prompt += fmt.Sprintf("PHASE %d: %s\n", phaseNumbers[phase], strings.ToUpper(phase))
            prompt += phaseDescriptions[phase] + "\n"
            prompt += "═══════════════════════════════════════════════════════════════\n"
            prompt += strings.TrimRight(phaseContent, "\n") + "\n\n"
    
    return prompt

function AssemblePreamble(procedureName string) -> string:
    preamble = "═══════════════════════════════════════════════════════════════\n"
    preamble += "ROODA PROCEDURE EXECUTION\n"
    preamble += "═══════════════════════════════════════════════════════════════\n\n"
    preamble += "Procedure: " + procedureName + "\n\n"
    
    // Role definition
    preamble += "Your Role:\n"
    preamble += "You are an AI coding agent executing a structured OODA loop procedure.\n"
    preamble += "This is NOT a template or example - this is an EXECUTABLE PROCEDURE.\n"
    preamble += "You must complete all phases and produce concrete outputs.\n\n"
    
    // Success signaling instructions
    preamble += "Success Signaling:\n"
    preamble += "- When you complete all tasks successfully, output: <promise>SUCCESS</promise>\n"
    preamble += "- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>\n"
    preamble += "- Explanations should come AFTER the signal, not embedded in the tag\n"
    preamble += "- The loop orchestrator uses these signals to determine iteration outcome.\n"
    
    return preamble

function ComposePhasePrompt(fragments []FragmentAction, configDir string) -> (string, error):
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
            content, err = LoadFragment(fragment.Path, configDir)
            if err:
                return error("failed to load fragment %s: %v", fragment.Path, err)
        
        // 3. Process template if parameters provided
        if len(fragment.Parameters) > 0:
            content, err = ProcessTemplate(content, fragment.Parameters)
            if err:
                return error("failed to process template: %v", err)
        
        // 4. Append to parts
        parts = append(parts, strings.TrimSpace(content))
    
    // 5. Concatenate with double newlines
    return strings.Join(parts, "\n\n"), nil

function LoadFragment(path string, configDir string) -> (string, error):
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

### Validation at Config Load Time

```
function ValidateProcedure(procedure, configDir):
    // Validate all four phases have valid fragment arrays
    for phase in [observe, orient, decide, act]:
        fragments = procedure[phase]  // Array of FragmentAction
        
        for i, fragment in enumerate(fragments):
            // Validate content vs path
            hasContent = fragment.Content != ""
            hasPath = fragment.Path != ""
            
            if hasContent && hasPath:
                return error("procedure %s: %s phase fragment %d: cannot specify both content and path",
                    procedure.name, phase, i)
            if !hasContent && !hasPath:
                return error("procedure %s: %s phase fragment %d: must specify either content or path",
                    procedure.name, phase, i)
            
            // Validate path exists if specified
            if hasPath:
                if strings.HasPrefix(fragment.Path, "builtin:"):
                    embeddedPath = strings.TrimPrefix(fragment.Path, "builtin:")
                    if !embeddedFragmentExists(embeddedPath):
                        return error("procedure %s: %s phase fragment %d: embedded fragment not found: %s\n" +
                                   "Available builtin fragments for %s phase:\n%s",
                                   procedure.name, phase, i, fragment.Path, phase,
                                   listBuiltinFragmentsForPhase(phase))
                else:
                    absolutePath = resolvePath(configDir, fragment.Path)
                    if !fileExists(absolutePath):
                        return error("procedure %s: %s phase fragment %d: fragment file not found: %s\n" +
                                   "Resolved to: %s\n\n" +
                                   "Tip: Check that the file exists and the path is correct.\n" +
                                   "     Paths are resolved relative to the config file directory.",
                                   procedure.name, phase, i, fragment.Path, absolutePath)
            
            // Validate template syntax if parameters provided
            if len(fragment.Parameters) > 0:
                content = ""
                if hasContent:
                    content = fragment.Content
                else:
                    content, err = loadFragmentContent(fragment.Path, configDir)
                    if err:
                        return err
                
                _, err = template.New("test").Parse(content)
                if err:
                    return error("procedure %s: %s phase fragment %d: template parse error: %v",
                        procedure.name, phase, i, err)
    
    return nil
```

## Edge Cases

### Fragment with Neither Content nor Path

**Scenario:** Fragment action has neither content nor path field

**Behavior:** Return error at config load time

**Example:**
```
Error: procedure build: observe phase fragment 0: must specify either content or path
```

### Fragment with Both Content and Path

**Scenario:** Fragment action specifies both content and path

**Behavior:** Return error at config load time

**Example:**
```
Error: procedure build: observe phase fragment 1: cannot specify both content and path
```

### Missing Fragment File

**Scenario:** Fragment path references a file that doesn't exist

**Behavior:** Return error at config load time

**Example:**
```
Error: procedure build: observe phase fragment 2: fragment file not found: fragments/missing.md
Resolved to: /project/root/fragments/missing.md

Tip: Check that the file exists and the path is correct.
     Paths are resolved relative to the config file directory.
```

### Empty Fragment Array

**Scenario:** OODA phase has empty fragment array

**Behavior:** Phase content is empty string (not an error)

**Rationale:** Some procedures may not need all phases; empty arrays are valid

### Template Syntax Error

**Scenario:** Fragment content has invalid Go template syntax

**Behavior:** Return error at config load time (when parameters provided)

**Example:**
```
Error: procedure custom-audit: observe phase fragment 0: template parse error: unclosed action
```

### Template Execution Error

**Scenario:** Template executes but encounters runtime error (e.g., range over nil)

**Behavior:** Return error at prompt composition time

**Example:**
```
Error: failed to compose observe phase: failed to process template: range can't iterate over <nil>
```

### Missing Template Parameter

**Scenario:** Template references parameter not provided in parameters map

**Behavior:** Go template default behavior (zero value for type)

**Example:**
```markdown
# Template: "Hello {{.name}}"
# Parameters: {}
# Result: "Hello "
```

**Rationale:** Fragment authors should use defensive patterns like `{{if .param}}...{{end}}`

### Relative Path Resolution

**Scenario:** User provides custom fragment with relative path `fragments/custom.md`

**Behavior:** Resolve relative to config file directory, not current working directory

**Example:**
```bash
# Config at /project/root/rooda-config.yml references fragments/custom.md
# Resolves to /project/root/fragments/custom.md

# Works regardless of where rooda is invoked from:
cd /project/root
rooda build  # Resolves to /project/root/fragments/custom.md

cd /project/root/src
rooda build  # Still resolves to /project/root/fragments/custom.md
```

**Rationale:** Config-relative paths ensure reproducibility regardless of invocation directory

### Context Injection with Special Characters

**Scenario:** User context contains markdown formatting, code blocks, or special characters

**Behavior:** Inject verbatim, preserve all formatting

**Example:**
```bash
rooda build --context "Focus on the auth module:
\`\`\`go
func Authenticate(token string) error
\`\`\`"
```

Result: Context appears at top of prompt with code block intact

### Builtin Prefix Case Sensitivity

**Scenario:** User writes `Builtin:` or `BUILTIN:` instead of `builtin:`

**Behavior:** Return error, require exact `builtin:` prefix

**Example:**
```
Error: procedure build: observe phase fragment 0: fragment file not found: Builtin:fragments/observe/read_agents_md.md
Resolved to: /project/root/Builtin:fragments/observe/read_agents_md.md

Did you mean: builtin:fragments/observe/read_agents_md.md
Note: The builtin: prefix is case-sensitive and must be lowercase.
```

**Rationale:** YAML is case-sensitive; clear error message guides users to correct syntax

### Fragment Concatenation

**Scenario:** Multiple fragments in a phase array

**Behavior:** Concatenate with double newlines (\n\n) between fragments

**Example:**
```yaml
observe:
  - path: "builtin:fragments/observe/read_agents_md.md"
  - path: "builtin:fragments/observe/read_specs.md"
```

**Result:**
```
[Content from read_agents_md.md]

[Content from read_specs.md]
```

**Rationale:** Double newlines provide clear visual separation between fragments

### Inline Content with Template Parameters

**Scenario:** Fragment uses inline content with template parameters

**Behavior:** Process template on inline content

**Example:**
```yaml
observe:
  - content: "Read {{.file_type}} files: {{range .paths}}{{.}} {{end}}"
    parameters:
      file_type: "spec"
      paths: ["a.md", "b.md"]
```

**Result:**
```
Read spec files: a.md b.md
```

## Dependencies

- **configuration.md** — Procedure definitions specify fragment arrays for each OODA phase
- **procedures.md** — Defines fragment structure, template system, and built-in fragments
- **Embedded fragment files** — Default fragments shipped with rooda binary (55 files in fragments/ directory)
- **Filesystem access** — For reading custom user-provided fragments
- **Go text/template** — Template processing for parameterized fragments

## Implementation Mapping

### Source Files (v2 Go implementation)

- `internal/prompt/composer.go` — Core assembly logic for fragment arrays
- `internal/prompt/resolver.go` — Path resolution (builtin: vs filesystem)
- `internal/prompt/template.go` — Go text/template processing
- `internal/prompt/embed.go` — Embedded fragment access via go:embed
- `fragments/` — 55 embedded default fragment files organized by OODA phase
  - `fragments/observe/` — 13 observe phase fragments
  - `fragments/orient/` — 20 orient phase fragments
  - `fragments/decide/` — 10 decide phase fragments
  - `fragments/act/` — 12 act phase fragments

### Related Specs

- [configuration.md](configuration.md) — Defines procedure structure with OODA phase fragment arrays
- [procedures.md](procedures.md) — Defines fragment system, template syntax, and built-in procedures
- [iteration-loop.md](iteration-loop.md) — Consumes assembled prompts for each iteration
- [cli-interface.md](cli-interface.md) — Defines --context flag for user context injection
- [ai-cli-integration.md](ai-cli-integration.md) — Receives assembled prompt as stdin

## Examples

### Example 1: Basic Assembly with Builtin Fragments

**Input:**
```yaml
# rooda-config.yml
procedures:
  build:
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
      - path: "builtin:fragments/act/run_tests.md"
```

**Command:**
```bash
rooda build --max-iterations 5
# Assuming this is iteration 1
```

**Output (assembled prompt):**
```markdown
═══════════════════════════════════════════════════════════════
ROODA PROCEDURE EXECUTION
═══════════════════════════════════════════════════════════════

Procedure: build

Your Role:
You are an AI coding agent executing a structured OODA loop procedure.
This is NOT a template or example - this is an EXECUTABLE PROCEDURE.
You must complete all phases and produce concrete outputs.

Success Signaling:
- When you complete all tasks successfully, output: <promise>SUCCESS</promise>
- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>
- Explanations should come AFTER the signal, not embedded in the tag
- The loop orchestrator uses these signals to determine iteration outcome.

═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
[Content from read_agents_md.md]

[Content from read_specs.md]

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from understand_task_requirements.md]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from plan_implementation_approach.md]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from modify_files.md]

[Content from run_tests.md]
```

**Verification:** Preamble appears first with procedure name and enhanced markers, followed by all fragments loaded from embedded resources, concatenated with double newlines within each phase

---

### Example 2: Custom Fragments from Filesystem with Inline Content

**Input:**
```yaml
# rooda-config.yml
procedures:
  custom-build:
    observe:
      - path: "fragments/observe/custom_observe.md"
      - content: "Also check for any TODO comments in the code."
    orient:
      - path: "fragments/orient/custom_orient.md"
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
```

**Command:**
```bash
rooda custom-build
```

**Output:**
```markdown
═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
[Content from fragments/observe/custom_observe.md]

Also check for any TODO comments in the code.

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from fragments/orient/custom_orient.md]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from plan_implementation_approach.md]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from modify_files.md]
```

**Verification:** Mix of filesystem fragments, inline content, and builtin fragments; paths resolved correctly

---

### Example 3: User Context Injection

**Input:**
```yaml
# rooda-config.yml
procedures:
  build:
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
      - path: "builtin:fragments/observe/read_specs.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
```

**Command:**
```bash
rooda build --context "Focus on the authentication module. The new feature should integrate with the existing OAuth2 flow." --max-iterations 10
# Assuming this is iteration 3
```

**Output:**
```markdown
═══════════════════════════════════════════════════════════════
ROODA PROCEDURE EXECUTION
═══════════════════════════════════════════════════════════════

Procedure: build

Your Role:
You are an AI coding agent executing a structured OODA loop procedure.
This is NOT a template or example - this is an EXECUTABLE PROCEDURE.
You must complete all phases and produce concrete outputs.

Success Signaling:
- When you complete all tasks successfully, output: <promise>SUCCESS</promise>
- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>
- Explanations should come AFTER the signal, not embedded in the tag
- The loop orchestrator uses these signals to determine iteration outcome.

═══════════════════════════════════════════════════════════════
CONTEXT
═══════════════════════════════════════════════════════════════
Focus on the authentication module. The new feature should integrate with the existing OAuth2 flow.

═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
[Content from read_agents_md.md]

[Content from read_specs.md]

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from understand_task_requirements.md]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from plan_implementation_approach.md]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from modify_files.md]
```

**Verification:** Preamble appears first with enhanced markers, user context appears after preamble and before OODA phases, followed by all phases with enhanced section markers and phase descriptions

---

### Example 4: Unlimited Iterations Mode

**Input:**
```yaml
# rooda-config.yml
procedures:
  build:
    observe:
      - path: "builtin:fragments/observe/read_agents_md.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
```

**Command:**
```bash
rooda build --unlimited
# Assuming this is iteration 7
```

**Output:**
```markdown
═══════════════════════════════════════════════════════════════
ROODA PROCEDURE EXECUTION
═══════════════════════════════════════════════════════════════

Procedure: build

Your Role:
You are an AI coding agent executing a structured OODA loop procedure.
This is NOT a template or example - this is an EXECUTABLE PROCEDURE.
You must complete all phases and produce concrete outputs.

Success Signaling:
- When you complete all tasks successfully, output: <promise>SUCCESS</promise>
- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>
- Explanations should come AFTER the signal, not embedded in the tag
- The loop orchestrator uses these signals to determine iteration outcome.

═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
[Content from read_agents_md.md]

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from understand_task_requirements.md]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from plan_implementation_approach.md]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from modify_files.md]
```

**Verification:** Preamble uses enhanced markers with clear role definition and success signaling instructions

---

### Example 5: Template Processing with Parameters

**Input:**
```yaml
# rooda-config.yml
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
    orient:
      - path: "builtin:fragments/orient/evaluate_against_quality_criteria.md"
    decide:
      - path: "builtin:fragments/decide/identify_issues.md"
    act:
      - path: "builtin:fragments/act/write_audit_report.md"
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

**Command:**
```bash
rooda custom-audit
```

**Output (observe phase):**
```markdown
═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
Read the following specification files from the repository:

- specs/api.md
- specs/database.md

Also include any associated test files.

═══════════════════════════════════════════════════════════════
PHASE 2: ORIENT
Analyze the information you gathered and form your understanding.
═══════════════════════════════════════════════════════════════
[Content from evaluate_against_quality_criteria.md]

═══════════════════════════════════════════════════════════════
PHASE 3: DECIDE
Make decisions about what actions to take.
═══════════════════════════════════════════════════════════════
[Content from identify_issues.md]

═══════════════════════════════════════════════════════════════
PHASE 4: ACT
Execute the actions you decided on. Modify files, run commands, commit changes.
═══════════════════════════════════════════════════════════════
[Content from write_audit_report.md]
```

**Verification:** Template processed with parameters, Go text/template syntax executed correctly

---

### Example 6: Missing Fragment File Error

**Input:**
```yaml
# rooda-config.yml
procedures:
  broken:
    observe:
      - path: "builtin:fragments/observe/missing.md"
    orient:
      - path: "builtin:fragments/orient/understand_task_requirements.md"
    decide:
      - path: "builtin:fragments/decide/plan_implementation_approach.md"
    act:
      - path: "builtin:fragments/act/modify_files.md"
```

**Command:**
```bash
rooda broken
```

**Output (error):**
```
Error: procedure broken: observe phase fragment 0: embedded fragment not found: builtin:fragments/observe/missing.md

Available builtin fragments for observe phase:
  read_agents_md.md
  scan_repo_structure.md
  detect_build_system.md
  detect_work_tracking.md
  verify_commands.md
  query_work_tracking.md
  read_specs.md
  read_impl.md
  read_task_input.md
  read_draft_plan.md
  read_task_details.md
  run_tests.md
  run_lints.md
```

**Verification:** Clear error message at config load time, suggests available alternatives

## Notes

### Design Rationale

**Why fragment arrays instead of single files?**
- Reusability: Common fragments (read_agents_md.md, read_specs.md) used across multiple procedures
- Flexibility: Procedures compose different combinations without duplication
- Maintainability: Update one fragment, all procedures using it benefit
- Extensibility: Add new fragments without modifying existing procedures

**Why section markers?**
- Debugging: Easy to identify which phase produced output
- Readability: Clear visual separation in assembled prompt
- Parsing: AI can reference specific phases in its response

**Why builtin: prefix instead of @builtin/ or similar?**
- Simplicity: Single character prefix, no special escaping needed
- Familiarity: Similar to URL schemes (http:, file:)
- Clarity: Unambiguous distinction from filesystem paths

**Why inject user context at the top?**
- Precedence: User intent should frame the entire OODA cycle
- Visibility: AI sees context before any phase-specific instructions
- Simplicity: No need to parse or merge context into specific phases

**Why validate all fragments at config load time?**
- Fail fast: Catch configuration errors immediately, before any iteration starts
- Clear feedback: User knows about broken config before invoking a procedure
- Atomicity: Either get complete prompt or clear error, no partial states
- Debugging: Error messages point to exact missing file with resolved absolute path

**Why double newlines between fragments?**
- Visual separation: Clear boundaries between fragment content
- Markdown compatibility: Ensures proper paragraph/section separation
- Readability: Easier for both humans and AI to parse

**Why Go text/template?**
- Built-in: No external dependency
- Familiar: Standard Go template syntax
- Safe: No arbitrary code execution
- Good errors: Clear messages for syntax issues

### Embedded Fragments

The complete list of embedded fragments (55 files organized by OODA phase) is maintained in [procedures.md](procedures.md#fragment-directory-structure). See the Fragment Directory Structure section for the authoritative list of all built-in fragments.

### Future Enhancements (Out of Scope for v2)

- **Advanced template validation** — Validate that all referenced parameters are provided
- **Conditional phases** — Skip phases based on runtime conditions
- **Fragment includes** — Fragments that reference other fragments
- **Fragment validation** — Lint fragments for common mistakes or missing instructions
- **Fragment versioning** — Track which fragment versions produced which outputs
- **Fragment caching** — Cache loaded and processed fragments for performance
