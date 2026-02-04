# User Documentation

## Job to be Done
Enable users to understand and effectively use the ralph-wiggum-ooda framework through clear, accessible documentation that guides them from installation through advanced usage.

## Activities
1. Understand framework concepts (OODA loop, fresh context, composable prompts)
2. Install and configure the framework in their project
3. Execute procedures and workflows for their use case
4. Troubleshoot common issues and errors
5. Extend the framework with custom procedures

## Acceptance Criteria
- [x] README.md contains installation instructions, basic workflow, and links to detailed docs
- [x] docs/ directory contains detailed guides for concepts, workflows, and troubleshooting
- [x] All code examples in documentation are verified working
- [x] Documentation matches actual script behavior (no contradictions)
- [x] Each procedure has usage examples with expected outcomes
- [x] Common error scenarios have troubleshooting guidance
- [x] Documentation follows progressive disclosure (quick start → detailed guides)
- [ ] Links between documents work correctly

## Data Structures

### Documentation Hierarchy
```yaml
README.md:
  sections:
    - installation
    - basic_workflow
    - procedures_table
    - key_principles
    - workflow_patterns
    - sample_structure
    - safety
    - troubleshooting
    - learn_more_links

docs/:
  files:
    - ooda-loop.md      # OODA framework explanation
    - ralph-loop.md     # Original methodology
    - beads.md          # Work tracking system
    - README.md         # Docs index

specs/:
  purpose: "Requirements (not user-facing)"
  audience: "Agents and developers"

AGENTS.md:
  purpose: "Operational guide"
  audience: "Agents only"
```

**Fields:**
- `README.md` - Main user entry point, progressive disclosure from quick start to detailed docs
- `docs/` - Detailed conceptual guides, deep-dive explanations, extended examples
- `specs/` - Requirements and specifications (not user-facing)
- `AGENTS.md` - Operational guide for agents (not user-facing)

## Algorithm

1. **User discovers project** → README.md provides value proposition and installation
2. **User installs** → README.md installation section with copy-paste commands
3. **User learns basics** → README.md basic workflow section with common procedures
4. **User needs details** → README.md links to docs/ for deep dives
5. **User encounters issue** → README.md troubleshooting section with common solutions
6. **User extends framework** → README.md custom procedures section with examples

**Pseudocode:**
```
function DocumentationFlow(user_need):
    if user_need == "what_is_this":
        return README.md#introduction
    elif user_need == "how_to_install":
        return README.md#installation
    elif user_need == "quick_start":
        return README.md#basic_workflow
    elif user_need == "understand_concepts":
        return docs/ooda-loop.md
    elif user_need == "troubleshoot":
        return README.md#troubleshooting
    elif user_need == "extend":
        return README.md#custom_procedures
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User has existing installation | README.md "For Existing Installations" section guides update process |
| Documentation contradicts behavior | Quality criteria fails, triggers documentation update |
| Code example doesn't work | Quality criteria fails, example must be fixed and verified |
| User needs agent-level details | README.md links to specs/ for developers, but clarifies it's not user-facing |
| Multiple docs cover same topic | Cross-reference with "See also" links, maintain single source of truth |

## Dependencies

- README.md must exist at project root
- docs/ directory must exist with conceptual guides
- specs/ directory exists but is not user-facing
- AGENTS.md exists but is not user-facing
- src/rooda.sh must be functional for examples to work

## Implementation Mapping

**Source files:**
- `README.md` - Main user entry point with installation, workflows, troubleshooting
- `docs/ooda-loop.md` - Explains OODA decision-making framework
- `docs/ralph-loop.md` - Original Ralph Loop methodology by Geoff Huntley
- `docs/beads.md` - Work tracking system documentation
- `docs/README.md` - Index of detailed documentation
- `src/README.md` - Developer reference for prompt composition (not user-facing)

**Related specs:**
- `specification-system.md` - Defines how specs are structured (not user-facing)
- `component-authoring.md` - Defines how to create prompts (developer-facing)

## Examples

### Example 1: New User Wants to Start Project

**Input:**
User visits repository, wants to understand what this is and how to use it.

**Expected Output:**
1. README.md introduction explains value proposition
2. README.md installation section provides copy-paste commands
3. README.md basic workflow shows first procedures to run
4. README.md links to docs/ for deeper understanding

**Verification:**
- User can install without reading anything beyond README.md
- User can run first procedure (bootstrap) successfully
- User knows where to find detailed guides

### Example 2: User Wants to Understand OODA Phases

**Input:**
User reads README.md, sees mention of "OODA phases", wants to understand deeply.

**Expected Output:**
1. README.md "Learn More" section links to docs/ooda-loop.md
2. docs/ooda-loop.md explains Observe, Orient, Decide, Act in detail
3. Examples show how phases map to prompt files

**Verification:**
- Link from README.md to docs/ooda-loop.md works
- docs/ooda-loop.md provides comprehensive explanation
- User understands how to read prompt files after reading

### Example 3: User Encounters "Command Not Found" Error

**Input:**
User runs `./rooda.sh bootstrap` and gets "command not found" error.

**Expected Output:**
1. README.md troubleshooting section lists common errors
2. "Command not found" → check if rooda.sh is executable
3. Solution: `chmod +x rooda.sh`

**Verification:**
- Error is listed in troubleshooting section
- Solution is actionable and correct
- User can resolve issue without external help

### Example 4: User Wants to Create Custom Procedure

**Input:**
User has specific workflow not covered by built-in procedures.

**Expected Output:**
1. README.md "Custom Procedures" section explains two approaches
2. Example shows editing rooda-config.yml
3. Example shows using command-line flags
4. Links to src/README.md for prompt composition details

**Verification:**
- Both approaches are documented with working examples
- User can create custom procedure following examples
- Link to detailed prompt composition guide works

## Notes

**Documentation Philosophy:**
- Progressive disclosure: README.md for quick start, docs/ for deep dives
- Action-oriented: Every section should enable user to do something
- Problem-solution framing: Address user needs, not just feature lists
- Verified examples: All code examples must be tested and working

**Data Flow:**
1. specs/ define requirements (what documentation must accomplish)
2. Implementation creates/updates README.md and docs/
3. Users consume README.md and docs/ to use the framework

**Separation of Concerns:**
- README.md: User-facing, installation, workflows, quick reference
- docs/: User-facing, detailed guides, concepts, extended examples
- specs/: Developer-facing, requirements, specifications
- AGENTS.md: Agent-facing, operational guide
- src/README.md: Developer-facing, prompt composition reference

## Known Issues

- Some documentation may contradict script behavior (quality criteria catches this)
- Code examples may become outdated as script evolves (quality criteria catches this)
- Cross-references between docs may break if files are renamed

## Areas for Improvement

- Add video tutorials for visual learners
- Create interactive examples or playground
- Add FAQ section based on common user questions
- Create migration guides for major version changes
- Add troubleshooting decision tree for complex issues
