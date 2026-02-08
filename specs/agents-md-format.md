# AGENTS.md Format

## Job to be Done

Specify the AGENTS.md format — required sections, field definitions, and structural conventions — that serves as the contract between AI agents and the repository. This file is the source of truth for project-specific behavior: build commands, file paths, work tracking, quality criteria. Agents read it first, verify it empirically, and update it when drift is detected.

## Activities

1. Define required sections and their purposes
2. Specify field formats and constraints for each section
3. Document structural conventions (markdown format, code blocks, lists)
4. Provide examples of well-formed AGENTS.md content
5. Define validation rules for each section
6. Specify how agents should interpret ambiguous or missing content

## Acceptance Criteria

- [ ] All required sections are defined with clear purposes
- [ ] Each section has field definitions with types and constraints
- [ ] Structural conventions are documented (markdown syntax, code blocks, lists)
- [ ] Examples provided for each required section
- [ ] Validation rules specified for detecting malformed content
- [ ] Bootstrap case documented (what to do when AGENTS.md doesn't exist)
- [ ] Update semantics defined (when to modify vs append vs replace)
- [ ] Cross-references to operational-knowledge.md for lifecycle behavior

## Data Structures

### AGENTS.md Schema

```markdown
# Agent Instructions

## Issue Tracking
[Brief description of work tracking system and onboarding steps]

## Quick Reference
[Code block with essential commands for common operations]

## Landing the Plane (Session Completion)
[Mandatory workflow steps for completing a work session]
[Critical rules about pushing to remote]

## Work Tracking System
**System:** [name of work tracking system]
[Commands for querying, updating, closing, creating work items]

## Story/Bug Input
[Where and how stories/bugs are documented]
[Prerequisites for planning procedures]

## Planning System
[Where draft plans are stored]
[How plans are published to work tracking]

## Build/Test/Lint Commands
**Dependencies:** [list of required tools with versions]
**Test:** [command to run tests, or "Manual verification"]
**Build:** [command to build, or "Not required"]
**Lint:** [command to run linters]
[Additional verification commands]

## Specification Definition
**Location:** [file paths or glob patterns]
**Format:** [description of spec format]
**Exclude:** [paths to exclude from specs]
**Current state:** [brief description of spec status]

## Implementation Definition
**Location:** [file paths or glob patterns]
**Patterns:** [list of implementation file patterns with descriptions]
**Exclude:** [paths to exclude from implementation]
**Note:** [any important context about implementation structure]

## Quality Criteria
**For specifications:** [list of PASS/FAIL criteria]
**For implementation:** [list of PASS/FAIL criteria]
**Refactoring triggers:** [conditions that require refactoring]

## Operational Learnings
**Last Bootstrap Verification:** [YYYY-MM-DD]
**Verified Working:** [list of commands/behaviors confirmed working]
**Verified Not Working / Missing:** [list of known issues]
**Why These Definitions:** [rationale for key decisions]
```

### Section Definitions

#### Issue Tracking (Required)
**Purpose:** Identify the work tracking system and provide onboarding instructions.

**Fields:**
- System name (e.g., "bd (beads)", "GitHub Issues", "file-based")
- Onboarding command or instructions

**Format:** Brief paragraph with system name in bold, followed by onboarding steps.

**Example:**
```markdown
## Issue Tracking

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.
```

#### Quick Reference (Required)
**Purpose:** Provide essential commands for common operations without explanation.

**Fields:**
- Command with inline comment explaining purpose

**Format:** Code block (bash) with one command per line, inline comments with `#`.

**Example:**
```markdown
## Quick Reference

\```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
\```
```

#### Landing the Plane (Required)
**Purpose:** Define mandatory workflow for completing a work session, emphasizing git push.

**Fields:**
- Numbered workflow steps
- Critical rules (bullet list)

**Format:** Bold heading, numbered list of steps (may include sub-bullets), followed by "CRITICAL RULES:" with bullet list.

**Example:**
```markdown
## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work**
2. **Run quality gates** (if code changed)
3. **PUSH TO REMOTE** - This is MANDATORY:
   \```bash
   git push
   \```

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing
```

#### Work Tracking System (Required)
**Purpose:** Document commands for interacting with the work tracking system.

**Fields:**
- System name (bold, prefixed with "System:")
- Query command with description
- Update command with description
- Close command with description
- Create command with description

**Format:** Bold "System:" line, followed by subsections with bold labels and code blocks.

**Example:**
```markdown
## Work Tracking System

**System:** beads (bd CLI)

**Query ready work:**
\```bash
bd ready --json
\```

**Update status:**
\```bash
bd update <id> --status in_progress
\```
```

#### Story/Bug Input (Optional)
**Purpose:** Document where stories/bugs are documented and prerequisites for planning.

**Fields:**
- Location of story/bug documentation
- Prerequisites for planning procedures

**Format:** Paragraph with file path and instructions.

#### Planning System (Optional)
**Purpose:** Document where draft plans are stored and how they're published.

**Fields:**
- Draft plan location
- Publishing mechanism

**Format:** Bold labels with descriptions.

#### Build/Test/Lint Commands (Required)
**Purpose:** Define commands for running tests, builds, and linters.

**Fields:**
- Dependencies (list with versions and descriptions)
- Test command or "Manual verification"
- Build command or "Not required"
- Lint command(s)
- Additional verification commands

**Format:** Bold labels followed by code blocks or descriptions. Dependencies as bullet list.

**Constraints:**
- If no automated tests, must say "Manual verification"
- If no build step, must say "Not required"
- Each command must be executable (agents will verify empirically)

**Example:**
```markdown
## Build/Test/Lint Commands

**Dependencies:**
- yq >= 4.0.0 (required) — YAML parsing
- shellcheck (optional) — bash linting

**Test:** Manual verification (no automated tests)

**Build:** Not required (bash scripts are interpreted)

**Lint:**
\```bash
shellcheck rooda.sh
\```
```

#### Specification Definition (Required)
**Purpose:** Define what constitutes a specification in this repository.

**Fields:**
- Location (file paths or glob patterns)
- Format (description of spec structure)
- Exclude (paths to exclude)
- Current state (brief status description)

**Format:** Bold labels with descriptions.

**Constraints:**
- Location must be a valid glob pattern or file path
- Exclude must list files that match Location but aren't specs

**Example:**
```markdown
## Specification Definition

**Location:** `specs/*.md`

**Format:** Markdown specifications following JTBD structure

**Exclude:** `specs/README.md` (index file, not a spec)

**Current state:** v2 Go rewrite specs in progress
```

#### Implementation Definition (Required)
**Purpose:** Define what constitutes implementation in this repository.

**Fields:**
- Location (file paths or glob patterns)
- Patterns (list of implementation file patterns with descriptions)
- Exclude (paths to exclude)
- Note (optional context about implementation structure)

**Format:** Bold labels with descriptions. Patterns as bullet list.

**Constraints:**
- Location must be valid file paths or glob patterns
- Exclude must list paths that should not be modified by agents

**Example:**
```markdown
## Implementation Definition

**Location:** `rooda.sh`, `rooda-config.yml`, `prompts/*.md`

**Patterns:**
- `rooda.sh` — Main OODA loop script
- `rooda-config.yml` — Procedure definitions
- `prompts/*.md` — OODA prompt components

**Exclude:**
- `archive/` — archived v1 implementation
- `.beads/` — work tracking database
```

#### Quality Criteria (Required)
**Purpose:** Define PASS/FAIL criteria for specifications and implementation.

**Fields:**
- For specifications (bullet list of PASS/FAIL criteria)
- For implementation (bullet list of PASS/FAIL criteria)
- Refactoring triggers (bullet list of conditions)

**Format:** Bold subsection labels with bullet lists.

**Constraints:**
- Each criterion must be boolean (PASS/FAIL)
- Each criterion must be verifiable (executable command or file check)
- Must include "(PASS/FAIL)" suffix on each criterion

**Example:**
```markdown
## Quality Criteria

**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)
- All specs have "Acceptance Criteria" section (PASS/FAIL)

**For implementation:**
- `./rooda.sh --version` executes without errors (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails
- Documentation contradicts actual behavior
```

#### Operational Learnings (Required)
**Purpose:** Document empirically verified behavior and known issues.

**Fields:**
- Last Bootstrap Verification (YYYY-MM-DD date)
- Verified Working (bullet list)
- Verified Not Working / Missing (bullet list)
- Why These Definitions (bullet list of rationale)

**Format:** Bold labels with bullet lists or date.

**Constraints:**
- Last Bootstrap Verification must be ISO date (YYYY-MM-DD)
- Verified Working must list commands/behaviors confirmed working
- Verified Not Working must list known issues or missing features
- Why These Definitions must explain key decisions

**Update semantics:**
- Agents update this section in-place when drift detected
- Add inline rationale for changes (brief comment)
- Don't append dated diary entries

**Example:**
```markdown
## Operational Learnings

**Last Bootstrap Verification:** 2026-02-06

**Verified Working:**
- `./rooda.sh --version` returns v0.1.0
- `bd ready --json` returns valid JSON

**Verified Not Working / Missing:**
- shellcheck not installed on this machine
- No automated test suite exists

**Why These Definitions:**
- Implementation is at root level because `goify` branch restructured
- Quality criteria are boolean PASS/FAIL for automated verification
```

## Algorithm

### Validation Algorithm

```
function ValidateAGENTSMD(content):
    sections = ParseMarkdownSections(content)
    
    required = [
        "Issue Tracking",
        "Quick Reference",
        "Landing the Plane",
        "Work Tracking System",
        "Build/Test/Lint Commands",
        "Specification Definition",
        "Implementation Definition",
        "Quality Criteria",
        "Operational Learnings"
    ]
    
    for section in required:
        if section not in sections:
            return FAIL, "Missing required section: " + section
    
    # Validate Build/Test/Lint Commands
    btl = sections["Build/Test/Lint Commands"]
    if not (contains(btl, "**Test:**") or contains(btl, "**Build:**") or contains(btl, "**Lint:**")):
        return FAIL, "Build/Test/Lint Commands missing required fields"
    
    # Validate Quality Criteria has PASS/FAIL markers
    qc = sections["Quality Criteria"]
    if not contains(qc, "(PASS/FAIL)"):
        return FAIL, "Quality Criteria must have (PASS/FAIL) markers"
    
    # Validate Operational Learnings has date
    ol = sections["Operational Learnings"]
    if not matches(ol, /\d{4}-\d{2}-\d{2}/):
        return FAIL, "Operational Learnings missing Last Bootstrap Verification date"
    
    return PASS, "Valid AGENTS.md"
```

### Bootstrap Algorithm (AGENTS.md doesn't exist)

```
function BootstrapAGENTSMD(repo_path):
    # Detect work tracking system
    if exists(repo_path + "/.beads"):
        work_tracking = "beads (bd CLI)"
        query_cmd = "bd ready --json"
    elif exists(repo_path + "/.github"):
        work_tracking = "GitHub Issues"
        query_cmd = "gh issue list --json number,title,state"
    else:
        work_tracking = "file-based (TASK.md)"
        query_cmd = "cat TASK.md"
    
    # Detect build system
    if exists(repo_path + "/Makefile"):
        build_cmd = "make"
        test_cmd = "make test"
    elif exists(repo_path + "/package.json"):
        build_cmd = "npm run build"
        test_cmd = "npm test"
    elif exists(repo_path + "/go.mod"):
        build_cmd = "go build ./..."
        test_cmd = "go test ./..."
    else:
        build_cmd = "Not required"
        test_cmd = "Manual verification"
    
    # Detect spec location
    if exists(repo_path + "/specs"):
        spec_location = "specs/*.md"
    elif exists(repo_path + "/docs/specs"):
        spec_location = "docs/specs/*.md"
    else:
        spec_location = "Not defined"
    
    # Generate AGENTS.md from template with detected values
    return RenderTemplate(AGENTS_MD_TEMPLATE, {
        work_tracking: work_tracking,
        query_cmd: query_cmd,
        build_cmd: build_cmd,
        test_cmd: test_cmd,
        spec_location: spec_location
    })
```

## Edge Cases

### Missing AGENTS.md
**Scenario:** Repository has no AGENTS.md file.

**Behavior:** Bootstrap algorithm creates AGENTS.md from template with detected values. Agent verifies empirically (runs commands, checks paths) and updates if detection was wrong.

### Malformed AGENTS.md
**Scenario:** AGENTS.md exists but missing required sections or has invalid format.

**Behavior:** Validation fails with clear error message. Agent can either fix in-place or regenerate from bootstrap.

### Conflicting Information
**Scenario:** AGENTS.md says "Test: `npm test`" but `npm test` fails with "script not found".

**Behavior:** Agent detects drift through empirical verification, updates AGENTS.md to reflect actual state (e.g., "Test: Manual verification").

### Multiple Work Tracking Systems
**Scenario:** Repository has both `.beads/` and `.github/` directories.

**Behavior:** Precedence order: beads > GitHub Issues > file-based. Agent documents the active system in AGENTS.md.

### Ambiguous Spec/Impl Definitions
**Scenario:** AGENTS.md says "Location: `src/*.go`" but `src/` contains both specs and implementation.

**Behavior:** Agent refines definition by examining file content (looks for "Job to be Done" section for specs) and updates AGENTS.md with more precise patterns.

### Empty Operational Learnings
**Scenario:** AGENTS.md has Operational Learnings section but no content under "Verified Working".

**Behavior:** Valid (section exists). Agent populates on first verification pass.

### Outdated Last Bootstrap Verification
**Scenario:** Last Bootstrap Verification is 6 months old.

**Behavior:** Not an error. Date indicates when last full verification occurred. Agent updates date when running bootstrap or full verification.

### Quality Criteria Without (PASS/FAIL)
**Scenario:** Quality Criteria lists criteria but omits "(PASS/FAIL)" markers.

**Behavior:** Validation warning (not error). Agent can still interpret as boolean criteria but should add markers for clarity.

### Build Command Fails
**Scenario:** AGENTS.md says "Build: `go build ./...`" but command exits with error.

**Behavior:** Agent distinguishes between "command doesn't exist" (drift, update AGENTS.md) and "command failed due to code issue" (don't update AGENTS.md, this is a code problem).

## Dependencies

- **operational-knowledge.md** — Defines the read-verify-update lifecycle that uses this schema
- **procedures.md** — All procedures read AGENTS.md first per this schema

## Implementation Mapping

**Related specs:**
- `operational-knowledge.md` — Runtime behavior for reading/verifying/updating AGENTS.md
- `procedures.md` — Built-in procedures that consume AGENTS.md

**Implementation files (v2 Go):**
- `internal/agents/schema.go` — AGENTS.md schema definition and validation
- `internal/agents/parser.go` — Markdown parsing for AGENTS.md sections
- `internal/agents/bootstrap.go` — Bootstrap algorithm for creating AGENTS.md
- `internal/agents/verify.go` — Empirical verification (run commands, check paths)
- `cmd/rooda/agents_sync.go` — CLI command for `rooda agents-sync` procedure

## Examples

### Example 1: Minimal Valid AGENTS.md

**Input:**
```markdown
# Agent Instructions

## Issue Tracking
This project uses **file-based** tracking. Create `TASK.md` to document work.

## Quick Reference
\```bash
cat TASK.md  # View current tasks
\```

## Landing the Plane (Session Completion)
**MANDATORY WORKFLOW:**
1. **PUSH TO REMOTE**
   \```bash
   git push
   \```

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds

## Work Tracking System
**System:** file-based (TASK.md)

**Query ready work:**
\```bash
cat TASK.md
\```

**Update status:** Edit TASK.md manually

**Close issue:** Remove from TASK.md

**Create issue:** Add to TASK.md

## Build/Test/Lint Commands
**Test:** Manual verification
**Build:** Not required
**Lint:** Not configured

## Specification Definition
**Location:** Not defined
**Format:** Not defined
**Exclude:** None
**Current state:** No specs

## Implementation Definition
**Location:** `*.sh`
**Patterns:**
- `*.sh` — Shell scripts
**Exclude:** None

## Quality Criteria
**For specifications:** Not applicable (no specs)
**For implementation:**
- Scripts execute without syntax errors (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails

## Operational Learnings
**Last Bootstrap Verification:** 2026-02-07
**Verified Working:** None yet
**Verified Not Working / Missing:** No automated tests
**Why These Definitions:** Minimal project, no formal specs
```

**Verification:** PASS (all required sections present, valid format)

### Example 2: Full-Featured AGENTS.md (Current Project)

**Input:** See `/Users/maxdunn/Dev/ralph-wiggum-ooda/AGENTS.md`

**Verification:** PASS (all required sections present, valid format, comprehensive content)

### Example 3: Invalid AGENTS.md (Missing Required Section)

**Input:**
```markdown
# Agent Instructions

## Issue Tracking
This project uses **bd** (beads).

## Build/Test/Lint Commands
**Test:** `go test ./...`
```

**Verification:** FAIL ("Missing required section: Quick Reference")

### Example 4: Invalid AGENTS.md (Quality Criteria Without PASS/FAIL)

**Input:**
```markdown
# Agent Instructions
[... other sections ...]

## Quality Criteria
**For specifications:**
- All specs have "Job to be Done" section
- All specs have "Acceptance Criteria" section
```

**Verification:** WARNING ("Quality Criteria should have (PASS/FAIL) markers")

### Example 5: Bootstrap Detection (Go Project)

**Input:** Repository with `go.mod`, `.beads/`, `specs/` directory

**Output:**
```markdown
# Agent Instructions

## Issue Tracking
This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference
\```bash
bd ready              # Find available work
bd show <id>          # View issue details
\```

## Landing the Plane (Session Completion)
**MANDATORY WORKFLOW:**
1. **PUSH TO REMOTE**
   \```bash
   git push
   \```

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds

## Work Tracking System
**System:** beads (bd CLI)

**Query ready work:**
\```bash
bd ready --json
\```

## Build/Test/Lint Commands
**Dependencies:**
- Go >= 1.21 (required)

**Test:**
\```bash
go test ./...
\```

**Build:**
\```bash
go build ./...
\```

**Lint:**
\```bash
golangci-lint run
\```

## Specification Definition
**Location:** `specs/*.md`
**Format:** Markdown specifications
**Exclude:** `specs/README.md`
**Current state:** Detected specs directory

## Implementation Definition
**Location:** `*.go`, `cmd/`, `internal/`
**Patterns:**
- `*.go` — Go source files
- `cmd/` — CLI entry points
- `internal/` — Internal packages
**Exclude:** `vendor/`, `*_test.go`

## Quality Criteria
**For specifications:**
- All specs have "Job to be Done" section (PASS/FAIL)

**For implementation:**
- `go build ./...` executes without errors (PASS/FAIL)
- `go test ./...` passes (PASS/FAIL)

**Refactoring triggers:**
- Any quality criterion fails

## Operational Learnings
**Last Bootstrap Verification:** 2026-02-07
**Verified Working:** Bootstrap detection completed
**Verified Not Working / Missing:** Not yet verified
**Why These Definitions:** Auto-detected from repository structure
```

**Verification:** PASS (generated from bootstrap algorithm)

## Notes

### Design Rationale

**Why required sections?**
- Agents need consistent structure to parse and interpret AGENTS.md
- Missing sections cause agents to make incorrect assumptions
- Required sections cover the minimum information needed for any repository

**Why PASS/FAIL markers in Quality Criteria?**
- Forces boolean criteria that can be automated
- Prevents ambiguous criteria like "code should be clean"
- Enables automated quality gate verification

**Why Operational Learnings section?**
- Captures empirical verification results
- Documents known issues and working behavior
- Provides rationale for decisions (the "why")
- Prevents agents from repeating failed approaches

**Why bootstrap algorithm?**
- Solves chicken-and-egg problem: agents need AGENTS.md to work, but AGENTS.md doesn't exist yet
- Auto-detection reduces manual setup burden
- Empirical verification corrects wrong detections

**Why in-place updates instead of append-only?**
- Keeps AGENTS.md concise and current
- Prevents accumulation of outdated information
- Inline rationale provides context without bloat

**Why markdown format?**
- Human-readable and editable
- Git-friendly (diffs, merges)
- Widely supported by tools and editors
- Structured enough for parsing, flexible enough for prose

### Alternative Approaches Considered

**YAML/JSON format:**
- Pros: Easier to parse, strict schema validation
- Cons: Less human-readable, harder to edit, poor for prose explanations
- Decision: Markdown chosen for human-first design

**Separate files per section:**
- Pros: Easier to version control individual sections
- Cons: Harder to get holistic view, more files to manage
- Decision: Single file chosen for simplicity

**Append-only Operational Learnings:**
- Pros: Full history preserved
- Cons: File grows unbounded, hard to find current state
- Decision: In-place updates with inline rationale chosen for conciseness

**Optional vs Required sections:**
- Pros: Flexibility for minimal projects
- Cons: Agents can't rely on sections existing
- Decision: Required sections cover minimum viable information; optional sections for advanced features

### Future Considerations

**Schema versioning:**
- If AGENTS.md format changes significantly, may need version marker
- Could use frontmatter: `---\nversion: 2\n---`

**Multi-language support:**
- Current schema assumes English
- Could add optional `language: en` field for i18n

**Machine-readable metadata:**
- Could add YAML frontmatter for structured metadata
- Would preserve human-readable body while enabling strict parsing

**Validation tooling:**
- Could provide `rooda validate-agents` command
- Would run validation algorithm and report errors

**Template library:**
- Could provide templates for common project types (Go, Python, Node.js)
- Would speed up bootstrap process
