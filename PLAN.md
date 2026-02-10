# Plan: Make Composed Prompts More Actionable for AI Agents

## Problem Statement

AI agents occasionally fail to recognize that the composed prompt is an executable procedure rather than a template or reference document. The current prompt structure uses passive language ("Load and parse...", "Extract:", "Identify:") which reads like documentation rather than instructions.

## Root Cause Analysis

Current prompt structure issues:
1. **Passive voice** - "Load specification files" vs "You must load specification files"
2. **No explicit execution directive** - Missing clear "Execute this procedure now" framing
3. **Template-like formatting** - Bullet lists and section headers look like documentation
4. **No iteration context** - Agent doesn't know this is part of an OODA loop cycle
5. **Missing success criteria** - No clear "when you're done" signal
6. **No role definition** - Agent doesn't know it's an autonomous executor

## Proposed Solution

Add a **procedure execution preamble** that wraps the composed prompt with explicit instructions and context.

### Changes Required

#### 1. Update `prompt-composition.md` Specification

**Location:** `specs/prompt-composition.md`

**Change:** Add new section "Procedure Execution Preamble" that defines a header injected before OODA phases.

**Content to add:**
```markdown
### Procedure Execution Preamble

Before the OODA phase content, inject an execution preamble that:
- Identifies this as an executable procedure
- Defines the agent's role and autonomy level
- Explains the OODA loop iteration context
- Specifies success/failure signaling requirements
- Sets expectations for tool usage and file modifications

Format:
```
# ROODA PROCEDURE EXECUTION

You are an autonomous AI coding agent executing the **{procedure_name}** procedure.

## Your Role
- Execute all four OODA phases (Observe, Orient, Decide, Act) in sequence
- Use available tools to read files, run commands, and modify code
- Make decisions independently based on the information you gather
- Complete the work described in this iteration

## Iteration Context
- This is iteration {N} of up to {M} iterations (or "unlimited" if no max)
- Each iteration is a fresh process with clean context
- File-based state (AGENTS.md, specs, code) persists between iterations
- Previous iterations may have made progress - check current state first

## Success Signaling
When you have completed the work for this iteration:
- Output exactly: `<promise>SUCCESS</promise>`
- Include a summary of what was accomplished

If you are blocked and cannot make further progress:
- Output exactly: `<promise>FAILURE</promise>`
- Explain what is blocking you

If more work remains but you made progress:
- Do NOT output any promise signal
- The loop will continue to the next iteration

## Execution Instructions
The following sections define what you must do in each OODA phase.
Execute them in order. Use all available tools. Modify files as needed.

---
```
```

#### 2. Update `iteration-loop.md` Specification

**Location:** `specs/iteration-loop.md`

**Change:** Document that assembled prompts include the execution preamble with iteration context.

**Section to update:** "Algorithm" - step 4 "Assemble prompt from OODA phase files"

**Add:**
```markdown
Prompt assembly includes:
1. Procedure execution preamble (with iteration number, max iterations, success signaling instructions)
2. User context (if --context provided)
3. OODA phase content (observe, orient, decide, act)
```

#### 3. Update Fragment Language to Imperative Voice

**Location:** All fragment files in `internal/prompt/fragments/`

**Change:** Convert passive documentation language to imperative commands.

**Examples:**

**Before (passive):**
```markdown
# Read AGENTS.md

Load and parse the AGENTS.md configuration file from the repository root.

Extract:
- Work tracking system configuration
- Build/test/lint commands
```

**After (imperative):**
```markdown
# Read AGENTS.md

**You must load and parse AGENTS.md from the repository root.**

Use the file reading tool to load `./AGENTS.md`, then extract:
- Work tracking system configuration (under "Work Tracking System" section)
- Build/test/lint commands (under "Build/Test/Lint Commands" section)
- Specification paths (under "Specification Definition" section)
- Implementation paths (under "Implementation Definition" section)
- Quality criteria (under "Quality Criteria" section)

Store this information - you will use it throughout this procedure.
```

**Pattern to apply across all fragments:**
- Start with "You must..." or "Your task is to..."
- Use imperative verbs: "Load", "Execute", "Identify", "Modify"
- Add tool usage hints: "Use the file reading tool", "Run the command"
- Include explicit storage/usage instructions: "Store this for later", "You will use this in the Act phase"

#### 4. Add Iteration State to Prompt Context

**Location:** `internal/prompt/composer.go` (implementation)

**Change:** Pass iteration state to AssemblePrompt so preamble can include current iteration number.

**Function signature change:**
```go
// Before
func AssemblePrompt(proc Procedure, userContext string, configDir string) (string, error)

// After
func AssemblePrompt(proc Procedure, userContext string, configDir string, iterState *IterationState) (string, error)
```

**Preamble template:**
```
This is iteration {iterState.Iteration + 1} of {iterState.MaxIterations} (or "unlimited")
```

#### 5. Update Section Markers for Clarity

**Location:** `specs/prompt-composition.md`

**Change:** Make section markers more directive.

**Before:**
```
=== OBSERVE ===
```

**After:**
```
═══════════════════════════════════════════════════════════════
PHASE 1: OBSERVE
Execute these observation tasks to gather information.
═══════════════════════════════════════════════════════════════
```

Apply to all four phases with appropriate descriptions:
- OBSERVE: "Execute these observation tasks to gather information."
- ORIENT: "Analyze the information you gathered and form your understanding."
- DECIDE: "Make decisions about what actions to take."
- ACT: "Execute the actions you decided on. Modify files, run commands, commit changes."

## Implementation Priority

### Phase 1: Minimal Changes (Immediate Impact)
1. Add procedure execution preamble to prompt composition
2. Update section markers to be more directive
3. Test with existing fragments

### Phase 2: Fragment Updates (Incremental)
4. Update top 5 most-used fragments to imperative voice
5. Test and validate improvements
6. Update remaining fragments

### Phase 3: Iteration Context (Enhancement)
7. Add iteration state to prompt assembly
8. Include iteration number in preamble

## Success Criteria

- [ ] Dry-run output shows clear "ROODA PROCEDURE EXECUTION" header
- [ ] Prompt explicitly states agent role and autonomy
- [ ] Success/failure signaling instructions are prominent
- [ ] Section markers clearly indicate execution phases
- [ ] Fragment language uses imperative voice
- [ ] Iteration context included in preamble
- [ ] AI agents consistently recognize prompts as executable procedures

## Testing Strategy

1. Run `./bin/rooda build --dry-run` and verify preamble appears
2. Test with actual AI CLI execution - monitor for improved recognition
3. Compare iteration success rates before/after changes
4. Validate that agents emit promise signals correctly

## Risks & Mitigations

**Risk:** Longer prompts may exceed context windows
**Mitigation:** Preamble is ~500 chars, well within budget. Monitor total prompt size.

**Risk:** Imperative language may be too prescriptive
**Mitigation:** Keep instructions clear but allow agent autonomy in implementation details.

**Risk:** Breaking changes to existing procedures
**Mitigation:** Changes are additive (preamble) and stylistic (fragment voice). No structural changes.

## Related Specifications

- `specs/prompt-composition.md` - Prompt assembly algorithm
- `specs/iteration-loop.md` - Iteration context and state
- `specs/procedures.md` - Fragment structure and organization
- `specs/ai-cli-integration.md` - How prompts are piped to AI CLI

## Notes

The key insight is that AI agents need **explicit role definition** and **clear execution directives**. The current prompt assumes the agent will infer its role from context, but making it explicit eliminates ambiguity.

The preamble serves as a "system prompt" that frames the entire OODA procedure as an executable task rather than a reference document.
