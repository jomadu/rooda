# Documentation Structure

## Job to be Done

Define what documentation must exist, how it's organized, and how AI agents verify it matches reality.

Developers want documentation that stays synchronized with specs and code through automated verification, uses natural human writing (not AI patterns), and provides clear examples that work as documented — all without manual maintenance.

## Activities

1. Define required documentation files and their purposes
2. Specify section structure and content requirements per doc type
3. Define quality criteria agents use to audit docs
4. Specify humanizer skill requirement for removing AI writing patterns
5. Define verification procedures for checking docs against source truth
6. Specify cross-reference conventions and validation rules
7. Define detection logic for missing or outdated documentation

## Acceptance Criteria

- [ ] List of required documentation files with purpose and required sections
- [ ] Quality criteria for each doc type (testable by AI agents)
- [ ] Humanizer skill requirement: all user-facing docs must pass humanizer review before commit
- [ ] Verification procedures: how agents check docs against specs/code
- [ ] Cross-reference syntax and resolution rules
- [ ] Missing/outdated doc detection logic
- [ ] Writing style guidelines that agents must follow
- [ ] Example verification workflow showing doc → code comparison

## Data Structures

### DocumentationFile

```go
type DocumentationFile struct {
    Path             string              // File path relative to project root
    Purpose          string              // What this doc explains
    RequiredSections []string            // Section headings that must exist
    SourceTruth      []string            // Where to verify content (specs, code files)
    QualityCriteria  []QualityCriterion  // Testable quality checks
}
```

### QualityCriterion

```go
type QualityCriterion struct {
    Description     string   // Human-readable criterion
    VerificationCmd string   // Command to verify (empty if manual)
    PassCondition   string   // What indicates pass
}
```

### CrossReference

```go
type CrossReference struct {
    SourceFile string   // Doc file containing reference
    TargetFile string   // Referenced file
    LineNumber int      // Line in source file
    RefType    string   // Type: spec-ref, code-ref, doc-ref
}
```

## Algorithm

### Documentation Verification Workflow

```
function VerifyDocumentation(docFile DocumentationFile) -> (bool, []string):
    errors = []
    
    // 1. Check file exists
    if !fileExists(docFile.Path):
        errors.append("Missing: " + docFile.Path)
        return false, errors
    
    // 2. Check required sections exist
    content = readFile(docFile.Path)
    for section in docFile.RequiredSections:
        if !containsSection(content, section):
            errors.append("Missing section: " + section + " in " + docFile.Path)
    
    // 3. Verify content against source truth
    for sourcePath in docFile.SourceTruth:
        sourceContent = readFile(sourcePath)
        drift = detectDrift(content, sourceContent, docFile.Path)
        if drift:
            errors.append(drift)
    
    // 4. Run quality criteria checks
    for criterion in docFile.QualityCriteria:
        if criterion.VerificationCmd != "":
            result = executeCommand(criterion.VerificationCmd)
            if !matchesPassCondition(result, criterion.PassCondition):
                errors.append("Quality check failed: " + criterion.Description)
    
    // 5. Check for AI writing patterns (humanizer skill)
    aiPatterns = detectAIPatterns(content)
    if len(aiPatterns) > 0:
        errors.append("AI writing patterns detected: " + join(aiPatterns, ", "))
    
    return len(errors) == 0, errors

function DetectDrift(docContent, sourceContent, docPath) -> string:
    // Extract documented behavior from doc
    documentedBehavior = extractExamples(docContent)
    
    // Extract actual behavior from source
    actualBehavior = extractBehavior(sourceContent)
    
    // Compare
    for docExample in documentedBehavior:
        if !matchesActual(docExample, actualBehavior):
            return "Drift in " + docPath + ": documented '" + docExample + "' doesn't match source"
    
    return ""

function DetectAIPatterns(content) -> []string:
    patterns = []
    
    // Common AI writing patterns to avoid
    aiPhrases = [
        "It's important to note that",
        "It's worth noting",
        "Keep in mind that",
        "Please note that",
        "Simply",
        "Just",
        "Easily",
        "Seamlessly",
        "Robust",
        "Leverage",
        "Utilize",
        "In order to",
        "At the end of the day",
    ]
    
    for phrase in aiPhrases:
        if contains(content, phrase):
            patterns.append(phrase)
    
    return patterns
```

### Cross-Reference Validation

```
function ValidateCrossReferences(docFiles []DocumentationFile) -> []string:
    errors = []
    
    for docFile in docFiles:
        content = readFile(docFile.Path)
        refs = extractCrossReferences(content)
        
        for ref in refs:
            // Check reference target exists
            if !fileExists(ref.TargetFile):
                errors.append("Broken reference in " + docFile.Path + ":" + ref.LineNumber + 
                            " -> " + ref.TargetFile + " (file not found)")
            
            // Check section exists if section reference
            if ref.Section != "":
                targetContent = readFile(ref.TargetFile)
                if !containsSection(targetContent, ref.Section):
                    errors.append("Broken reference in " + docFile.Path + ":" + ref.LineNumber + 
                                " -> " + ref.TargetFile + "#" + ref.Section + " (section not found)")
    
    return errors

function ExtractCrossReferences(content) -> []CrossReference:
    refs = []
    lines = split(content, "\n")
    
    for i, line in enumerate(lines):
        // Match markdown links: [text](path) or [text](path#section)
        matches = regex.FindAll(line, `\[([^\]]+)\]\(([^)]+)\)`)
        
        for match in matches:
            linkText = match[1]
            linkTarget = match[2]
            
            // Parse target (may include #section)
            parts = split(linkTarget, "#")
            targetFile = parts[0]
            section = ""
            if len(parts) > 1:
                section = parts[1]
            
            refs.append(CrossReference{
                SourceFile: currentFile,
                TargetFile: targetFile,
                LineNumber: i + 1,
                Section:    section,
            })
    
    return refs
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Documentation file missing | Error: "Missing: [path]" |
| Required section missing | Error: "Missing section: [section] in [file]" |
| Example in doc doesn't match code | Error: "Drift in [file]: documented '[example]' doesn't match source" |
| Cross-reference to non-existent file | Error: "Broken reference in [file]:[line] -> [target] (file not found)" |
| Cross-reference to non-existent section | Error: "Broken reference in [file]:[line] -> [target]#[section] (section not found)" |
| AI writing patterns detected | Error: "AI writing patterns detected: [list of phrases]" |
| Command example fails when executed | Error: "Example verification failed: [command] (exit code [N])" |
| Documentation more recent than source | Warning: "Doc may be ahead of implementation" |
| Empty documentation file | Error: "Empty file: [path]" |
| Documentation without examples | Warning: "No examples found in [file]" |

## Dependencies

- **operational-knowledge.md** — Docs follow same read-verify-update lifecycle as AGENTS.md
- **agents-md-format.md** — AGENTS.md documents docs/ location and patterns
- **procedures.md** — Audit/plan/build procedures treat docs as implementation
- **skills/humanizer/SKILL.md** — Defines AI writing patterns to remove

## Implementation Mapping

**Source files:**
- `internal/docs/verifier.go` — Documentation verification logic
- `internal/docs/parser.go` — Parse doc structure and extract examples
- `internal/docs/crossref.go` — Cross-reference validation
- `internal/docs/humanizer.go` — AI pattern detection
- `skills/humanizer/SKILL.md` — Humanizer skill definition

**Related specs:**
- `operational-knowledge.md` — Read-verify-update lifecycle
- `agents-md-format.md` — AGENTS.md schema
- `procedures.md` — How procedures interact with docs

## Examples

### Example 1: Required Documentation Files

**docs/installation.md**
- **Purpose:** Explain how to install rooda on different platforms
- **Required sections:** Prerequisites, Homebrew, Direct Download, Build from Source, Verification
- **Source truth:** `scripts/install.sh`, `Makefile`, `README.md`
- **Quality criteria:**
  - All installation methods documented
  - Installation commands execute successfully
  - Version command works after installation
  - No AI writing patterns

**docs/procedures.md**
- **Purpose:** Document all 16 built-in procedures with examples
- **Required sections:** Overview, Procedure List (one subsection per procedure with description and example)
- **Source truth:** `internal/procedures/builtin.go`, `rooda-config.yml`
- **Quality criteria:**
  - All procedures from builtin.go are documented
  - Each procedure has working example
  - Example commands execute without errors
  - No AI writing patterns

**docs/cli-reference.md**
- **Purpose:** Complete CLI flag reference with descriptions and examples
- **Required sections:** Commands, Global Flags, Procedure Flags, Exit Codes
- **Source truth:** `cmd/rooda/main.go`, `specs/cli-interface.md`
- **Quality criteria:**
  - All flags from main.go are documented
  - Flag descriptions match code comments
  - Examples execute successfully
  - Exit codes match implementation
  - No AI writing patterns

**docs/configuration.md**
- **Purpose:** Explain three-tier config system and YAML schema
- **Required sections:** Overview, Config Tiers, YAML Schema, Examples
- **Source truth:** `internal/config/config.go`, `specs/configuration.md`
- **Quality criteria:**
  - All config fields from config.go are documented
  - Tier precedence matches implementation
  - Example configs parse successfully
  - No AI writing patterns

**docs/troubleshooting.md**
- **Purpose:** Common errors and solutions
- **Required sections:** Installation Issues, Configuration Errors, Execution Errors, Getting Help
- **Source truth:** `internal/errors/`, user reports
- **Quality criteria:**
  - Error messages match actual output
  - Solutions resolve documented problems
  - No AI writing patterns

**docs/agents-md.md**
- **Purpose:** Explain AGENTS.md format and lifecycle
- **Required sections:** Purpose, Format, Required Sections, Bootstrap, Verification, Updates
- **Source truth:** `specs/agents-md-format.md`, `specs/operational-knowledge.md`
- **Quality criteria:**
  - Format matches agents-md-format.md spec
  - Lifecycle matches operational-knowledge.md spec
  - Examples are valid AGENTS.md content
  - No AI writing patterns

**README.md**
- **Purpose:** Project overview, quick start, installation
- **Required sections:** What is rooda, Quick Start, Installation, Core Concepts, Common Workflows, Documentation, Requirements, License
- **Source truth:** `specs/README.md`, `docs/installation.md`
- **Quality criteria:**
  - Quick start commands execute successfully
  - Installation links point to docs/installation.md
  - Core concepts match spec definitions
  - No AI writing patterns

### Example 2: Verification Workflow

**Scenario:** Verify docs/cli-reference.md matches cmd/rooda/main.go

**Step 1: Extract documented flags**
```bash
# Read docs/cli-reference.md
grep -E "^- \`--" docs/cli-reference.md
```

Output:
```
- `--max-iterations <n>` — Override max iterations
- `--unlimited` — Run until convergence
- `--dry-run` — Display prompt without executing
- `--context <value>` — Inject user context
```

**Step 2: Extract actual flags from code**
```bash
# Read cmd/rooda/main.go
grep -E "flags\.(Int|Bool|String)" cmd/rooda/main.go
```

Output:
```go
flags.IntP("max-iterations", "n", 0, "Override max iterations")
flags.BoolP("unlimited", "u", false, "Run until convergence")
flags.BoolP("dry-run", "d", false, "Display prompt without executing")
flags.StringArrayP("context", "c", nil, "Inject user context")
```

**Step 3: Compare**
- `--max-iterations`: ✓ Documented and implemented
- `--unlimited`: ✓ Documented and implemented
- `--dry-run`: ✓ Documented and implemented
- `--context`: ✓ Documented and implemented

**Result:** No drift detected

### Example 3: Drift Detection

**Scenario:** docs/cli-reference.md documents `--verbose` but code has `--debug`

**Documented:**
```markdown
- `--verbose` — Enable verbose output
```

**Actual code:**
```go
flags.BoolP("debug", "d", false, "Enable debug output")
```

**Drift detected:**
```
Error: Drift in docs/cli-reference.md: documented '--verbose' flag not found in cmd/rooda/main.go
Found similar flag: '--debug' (Enable debug output)
```

**Agent action:**
1. Update docs/cli-reference.md to use `--debug`
2. Add rationale comment: `# Changed from --verbose to match implementation`
3. Commit with message: "docs: fix CLI flag name (--debug not --verbose)"

### Example 4: AI Pattern Detection

**Input (docs/installation.md):**
```markdown
## Installation

It's important to note that rooda requires Go 1.24.5 or later. Simply run the following command to install:

```bash
brew install rooda
```

This will seamlessly install rooda and all its dependencies.
```

**AI patterns detected:**
- "It's important to note that"
- "Simply"
- "seamlessly"

**Humanized version:**
```markdown
## Installation

rooda requires Go 1.24.5 or later. Install with Homebrew:

```bash
brew install rooda
```

This installs rooda and its dependencies.
```

**Changes:**
- Removed "It's important to note that" (unnecessary preamble)
- Removed "Simply" (filler word)
- Removed "seamlessly" (marketing language)
- Made sentences direct and concise

### Example 5: Cross-Reference Validation

**Input (docs/configuration.md):**
```markdown
See [CLI Interface](cli-interface.md) for flag precedence.
See [Procedures](procedures.md#custom-procedures) for custom procedure definitions.
See [Missing File](missing.md) for more info.
```

**Validation:**
1. Check `cli-interface.md` exists: ✓
2. Check `procedures.md` exists: ✓
3. Check `procedures.md` has section "custom-procedures": ✓
4. Check `missing.md` exists: ✗

**Error:**
```
Broken reference in docs/configuration.md:3 -> missing.md (file not found)
```

**Agent action:**
1. Remove broken reference or fix target path
2. Commit with message: "docs: fix broken cross-reference"

## Notes

### Design Rationale

**Why treat docs as implementation artifacts?**
Documentation describes implementation behavior. When implementation changes, docs must change. Treating them the same way (read-verify-update) keeps them synchronized.

**Why require humanizer skill?**
AI-generated text has recognizable patterns (filler words, hedging, marketing language) that make docs feel robotic. Humanizer skill removes these patterns before commit.

**Why verify docs against code, not just specs?**
Specs define intent, code defines reality. Docs must match reality (what actually happens) not just intent (what should happen).

**Why cross-reference validation?**
Broken links frustrate users and erode trust. Automated validation catches broken references before they reach users.

**Why required sections per doc type?**
Consistency helps users find information. Required sections ensure all docs have minimum necessary content.

**Why detect AI patterns automatically?**
Manual review is slow and inconsistent. Automated detection catches common patterns reliably.

### Writing Style Guidelines

**Direct and concise:**
- ✗ "It's important to note that you should run tests"
- ✓ "Run tests"

**Active voice:**
- ✗ "The configuration can be modified by editing the file"
- ✓ "Edit the file to modify configuration"

**No hedging:**
- ✗ "This might help you install rooda"
- ✓ "Install rooda with this command"

**No marketing language:**
- ✗ "Seamlessly integrate with your workflow"
- ✓ "Integrate with your workflow"

**No filler words:**
- ✗ "Simply run the command"
- ✓ "Run the command"

**Concrete examples:**
- ✗ "You can use various flags"
- ✓ "Use `--max-iterations 5` to limit iterations"

### Verification Frequency

**When to verify docs:**
- Before committing doc changes (pre-commit hook)
- During `audit-impl` procedure (finds outdated docs)
- During `audit-spec-to-impl` procedure (finds missing docs)
- During `build` procedure when implementing doc tasks

**What to verify:**
- Required sections exist
- Examples execute successfully
- Cross-references resolve
- No AI writing patterns
- Content matches source truth (specs/code)

### Source Truth Hierarchy

**For behavior documentation:**
1. Code (what actually happens)
2. Specs (what should happen)
3. Existing docs (what was previously documented)

**For design documentation:**
1. Specs (design intent)
2. Code (implementation decisions)
3. Existing docs (previous explanations)

**Agents synthesize docs by reading all three sources and resolving conflicts in favor of source truth.**
