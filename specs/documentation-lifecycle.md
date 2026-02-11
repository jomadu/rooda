# Documentation Lifecycle

## Job to be Done

Specify how existing procedures interact with docs through read-verify-update cycles.

Developers want documentation to stay synchronized with specs and code through the same automated workflow that maintains AGENTS.md — no special doc procedures, just audit → plan → build cycles that treat docs as implementation artifacts.

## Activities

1. Bootstrap documents docs/ patterns in AGENTS.md
2. Audit-impl-to-spec finds undocumented features (code/docs exist, specs missing)
3. Audit-spec-to-impl finds specified features not in docs
4. Planning procedures convert gaps into tasks
5. Build synthesizes docs by reading specs/code when implementing doc tasks
6. Build applies humanizer skill before committing
7. Build commits doc changes with descriptive messages

## Acceptance Criteria

- [ ] Bootstrap procedure creates/updates AGENTS.md with docs/ patterns and quality criteria
- [ ] Bootstrap documents humanizer skill requirement for user-facing docs
- [ ] Audit-impl-to-spec finds docs that exist but aren't referenced in specs
- [ ] Audit-spec-to-impl finds features specified but not documented
- [ ] Draft-plan-* procedures convert doc gaps into actionable tasks
- [ ] Build procedure reads specs/code to synthesize docs when implementing doc tasks
- [ ] Build applies humanizer skill to remove AI patterns before committing docs
- [ ] Build commits doc changes with descriptive messages (e.g., "docs: add CLI reference for --context flag")
- [ ] No special doc procedures needed (reuse existing audit/plan/build)
- [ ] Docs follow same read-verify-update lifecycle as AGENTS.md

## Data Structures

### AGENTS.md Documentation Section

```markdown
## Documentation Definition

**Location:** `docs/*.md`, `README.md`

**Patterns:**
- `docs/installation.md` — Installation instructions for all platforms
- `docs/procedures.md` — All 16 built-in procedures with examples
- `docs/cli-reference.md` — Complete CLI flag reference
- `docs/configuration.md` — Three-tier config system and YAML schema
- `docs/troubleshooting.md` — Common errors and solutions
- `docs/agents-md.md` — AGENTS.md format and lifecycle
- `README.md` — Project overview, quick start, installation

**Exclude:**
- `specs/` — specifications (not user-facing docs)
- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files

**Humanizer Skill Requirement:**
All user-facing documentation (README.md, docs/) must pass humanizer review before commit. See `skills/humanizer/SKILL.md` for AI writing patterns to remove.

## Quality Criteria

**For documentation:**
- All required doc files exist (installation, procedures, cli-reference, configuration, troubleshooting, agents-md, README)
- All documented commands execute successfully (PASS/FAIL)
- All documented flags exist in CLI implementation (PASS/FAIL)
- All cross-references resolve to existing files/sections (PASS/FAIL)
- No AI writing patterns detected (PASS/FAIL) — requires humanizer skill
- Examples in docs match actual behavior (PASS/FAIL)
```

### Audit Report Structure (Doc Gaps)

```markdown
# Documentation Gap Report

## Undocumented Features (Code/Docs Exist, Specs Missing)

- `docs/cli-reference.md` exists but not referenced in specs
- `docs/troubleshooting.md` exists but not referenced in specs

## Specified Features Not Documented

- `--context` flag specified in specs/cli-interface.md but not in docs/cli-reference.md
- Template parameters specified in specs/procedures.md but not in docs/configuration.md

## Outdated Documentation

- docs/installation.md documents Homebrew install but specs/distribution.md specifies direct download only
```

### Task Structure (Doc Work)

```markdown
# Task: Document CLI --context flag

**Type:** chore (documentation maintenance)
**Priority:** 2
**Description:**

Add documentation for --context flag to docs/cli-reference.md.

Source truth:
- specs/cli-interface.md (design intent)
- cmd/rooda/main.go (actual implementation)

Acceptance criteria:
- Flag description matches code
- Example command executes successfully
- Cross-references to related flags
- No AI writing patterns
```

## Algorithm

### Bootstrap Procedure (AGENTS.md Documentation Section)

```
function BootstrapDocumentationSection() -> string:
    section = "## Documentation Definition\n\n"
    
    // Detect docs directory
    if directoryExists("docs/"):
        section += "**Location:** `docs/*.md`, `README.md`\n\n"
        
        // List doc files
        section += "**Patterns:**\n"
        docFiles = listFiles("docs/*.md")
        for file in docFiles:
            purpose = inferPurpose(file)
            section += "- `" + file + "` — " + purpose + "\n"
        
        if fileExists("README.md"):
            section += "- `README.md` — Project overview, quick start, installation\n"
        
        section += "\n"
        
        // Exclude patterns
        section += "**Exclude:**\n"
        section += "- `specs/` — specifications (not user-facing docs)\n"
        section += "- `AGENTS.md`, `PLAN.md`, `TASK.md` — operational files\n\n"
        
        // Humanizer requirement
        section += "**Humanizer Skill Requirement:**\n"
        section += "All user-facing documentation (README.md, docs/) must pass humanizer review before commit. "
        section += "See `skills/humanizer/SKILL.md` for AI writing patterns to remove.\n\n"
    else:
        section += "**Location:** No docs/ directory found\n\n"
    
    return section

function InferPurpose(filename) -> string:
    // Map common doc filenames to purposes
    purposes = {
        "installation.md": "Installation instructions for all platforms",
        "procedures.md": "All built-in procedures with examples",
        "cli-reference.md": "Complete CLI flag reference",
        "configuration.md": "Configuration system and YAML schema",
        "troubleshooting.md": "Common errors and solutions",
        "agents-md.md": "AGENTS.md format and lifecycle",
    }
    
    basename = path.Base(filename)
    if purpose, exists = purposes[basename]; exists:
        return purpose
    
    return "Documentation file"
```

### Audit-Impl-To-Spec (Find Undocumented Features)

```
function AuditImplToSpec() -> []Gap:
    gaps = []
    
    // Find docs that exist but aren't in specs
    docFiles = listFiles("docs/*.md")
    specFiles = listFiles("specs/*.md")
    
    for docFile in docFiles:
        referenced = false
        for specFile in specFiles:
            specContent = readFile(specFile)
            if contains(specContent, docFile):
                referenced = true
                break
        
        if !referenced:
            gaps.append(Gap{
                Type: "undocumented-feature",
                Description: docFile + " exists but not referenced in specs",
                SourceTruth: docFile,
                Action: "Add reference to appropriate spec",
            })
    
    return gaps
```

### Audit-Spec-To-Impl (Find Missing Docs)

```
function AuditSpecToImpl() -> []Gap:
    gaps = []
    
    // Find features specified but not documented
    specFiles = listFiles("specs/*.md")
    
    for specFile in specFiles:
        specContent = readFile(specFile)
        
        // Extract documented features (CLI flags, procedures, config fields)
        features = extractFeatures(specContent)
        
        for feature in features:
            // Check if feature is documented in docs/
            documented = isFeatureDocumented(feature, "docs/")
            
            if !documented:
                gaps.append(Gap{
                    Type: "missing-documentation",
                    Description: feature.Name + " specified in " + specFile + " but not in docs",
                    SourceTruth: specFile,
                    Action: "Document " + feature.Name + " in appropriate doc file",
                })
    
    return gaps

function IsFeatureDocumented(feature, docsDir) -> bool:
    docFiles = listFiles(docsDir + "*.md")
    
    for docFile in docFiles:
        docContent = readFile(docFile)
        if contains(docContent, feature.Name):
            return true
    
    return false
```

### Build Procedure (Synthesize Docs)

```
function BuildDocTask(task) -> error:
    // 1. Read source truth
    specs = readSourceTruthFiles(task.SourceTruth.Specs)
    code = readSourceTruthFiles(task.SourceTruth.Code)
    existingDocs = readFile(task.TargetFile)
    
    // 2. Synthesize doc content
    newContent = synthesizeDocContent(specs, code, existingDocs, task.Description)
    
    // 3. Apply humanizer skill
    humanizedContent = applyHumanizerSkill(newContent)
    
    // 4. Verify quality criteria
    errors = verifyDocQuality(humanizedContent, task.TargetFile)
    if len(errors) > 0:
        return error("Quality verification failed: " + join(errors, ", "))
    
    // 5. Write file
    writeFile(task.TargetFile, humanizedContent)
    
    // 6. Commit with descriptive message
    commitMessage = generateDocCommitMessage(task)
    gitCommit(task.TargetFile, commitMessage)
    
    return nil

function SynthesizeDocContent(specs, code, existingDocs, taskDescription) -> string:
    // Extract relevant information from source truth
    specInfo = extractRelevantInfo(specs, taskDescription)
    codeInfo = extractRelevantInfo(code, taskDescription)
    
    // Merge with existing docs (preserve structure)
    if existingDocs != "":
        return updateExistingDoc(existingDocs, specInfo, codeInfo)
    else:
        return createNewDoc(specInfo, codeInfo)

function ApplyHumanizerSkill(content) -> string:
    // Load humanizer skill
    humanizerSkill = readFile("skills/humanizer/SKILL.md")
    
    // Remove AI patterns
    humanized = content
    
    aiPatterns = [
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
    ]
    
    for pattern in aiPatterns:
        humanized = removePattern(humanized, pattern)
    
    return humanized

function GenerateDocCommitMessage(task) -> string:
    // Format: "docs: <action> <target>"
    // Examples:
    // - "docs: add CLI reference for --context flag"
    // - "docs: update installation instructions"
    // - "docs: fix broken cross-reference in configuration.md"
    
    action = inferAction(task.Description)
    target = inferTarget(task.TargetFile, task.Description)
    
    return "docs: " + action + " " + target
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No docs/ directory exists | Bootstrap documents "No docs/ directory found" in AGENTS.md |
| Doc file exists but empty | Audit-impl-to-spec reports as undocumented feature |
| Feature documented in README but not docs/ | Audit-spec-to-impl considers it documented (README counts) |
| Multiple docs reference same feature | Not an error, cross-references are valid |
| Doc task has no source truth specified | Build reads AGENTS.md to infer source truth from file patterns |
| Humanizer skill file missing | Build emits warning but continues (best-effort humanization) |
| Doc synthesis produces AI patterns | Humanizer skill removes them before commit |
| Commit message exceeds 72 chars | Truncate and add "..." |
| Doc file in specs/ directory | Excluded from doc lifecycle (specs are not user-facing docs) |
| AGENTS.md itself needs updating | Bootstrap procedure handles it (not build) |

## Dependencies

- **documentation-structure.md** — Defines required docs, quality criteria, verification procedures
- **operational-knowledge.md** — Docs follow same read-verify-update lifecycle as AGENTS.md
- **agents-md-format.md** — AGENTS.md documents docs/ location and patterns
- **procedures.md** — Bootstrap, audit, plan, build procedures
- **skills/humanizer/SKILL.md** — AI writing patterns to remove

## Implementation Mapping

**Source files:**
- `internal/prompt/fragments/observe/read_docs.md` — Fragment for reading docs/
- `internal/prompt/fragments/orient/identify_doc_gaps.md` — Fragment for finding missing/outdated docs
- `internal/prompt/fragments/decide/plan_doc_updates.md` — Fragment for planning doc changes
- `internal/prompt/fragments/act/write_docs.md` — Fragment for synthesizing docs
- `internal/prompt/fragments/act/apply_humanizer.md` — Fragment for removing AI patterns

**Related specs:**
- `documentation-structure.md` — What docs must exist and how to verify them
- `operational-knowledge.md` — Read-verify-update lifecycle
- `agents-md-format.md` — AGENTS.md schema
- `procedures.md` — Procedure definitions

## Examples

### Example 1: Greenfield Workflow (New Feature with Docs)

**Scenario:** Add user authentication feature with documentation

**Step 1: Bootstrap**
```bash
rooda run bootstrap --ai-cmd-alias kiro-cli
```

**Output:** Creates AGENTS.md with docs/ patterns section

**Step 2: Create task**
```bash
cat > TASK.md << 'EOF'
# Add User Authentication

Implement JWT-based user authentication with OAuth2 integration.

Requirements:
- JWT token generation and validation
- OAuth2 flow integration
- Session management
- Documentation in docs/authentication.md
EOF
```

**Step 3: Draft plan**
```bash
rooda run draft-plan-spec-feat --ai-cmd-alias kiro-cli
```

**Output:** Creates PLAN.md with tasks including:
- Write specs/authentication.md
- Implement internal/auth/ package
- Write docs/authentication.md
- Update docs/cli-reference.md with auth flags

**Step 4: Publish plan**
```bash
rooda run publish-plan --ai-cmd-alias kiro-cli
```

**Output:** Imports tasks to work tracking

**Step 5: Build**
```bash
rooda run build --ai-cmd-alias kiro-cli --max-iterations 10
```

**Output:** Implements specs, code, and docs in order:
1. Writes specs/authentication.md
2. Implements internal/auth/
3. Reads specs/authentication.md and internal/auth/ code
4. Synthesizes docs/authentication.md from specs and code
5. Applies humanizer skill to remove AI patterns
6. Commits: "docs: add authentication guide"

**Verification:**
```bash
# Check doc exists
ls docs/authentication.md

# Check quality criteria
grep -E "(Simply|Seamlessly|Leverage)" docs/authentication.md  # Should be empty

# Check examples work
grep -E "^\`\`\`bash" docs/authentication.md -A 5 | bash  # Should execute successfully
```

### Example 2: Brownfield Workflow (Code Exists, Docs Missing from Specs)

**Scenario:** docs/cli-reference.md exists but not referenced in specs

**Step 1: Bootstrap**
```bash
rooda run bootstrap --ai-cmd-alias kiro-cli
```

**Output:** Updates AGENTS.md with docs/ patterns (includes cli-reference.md)

**Step 2: Audit**
```bash
rooda run audit-impl-to-spec --ai-cmd-alias kiro-cli
```

**Output:** Gap report:
```markdown
# Implementation to Spec Gap Report

## Undocumented Features

- `docs/cli-reference.md` exists but not referenced in specs
- `docs/troubleshooting.md` exists but not referenced in specs
```

**Step 3: Convert to task**
```bash
cat > TASK.md << 'EOF'
# Document CLI Reference in Specs

Add reference to docs/cli-reference.md in specs/cli-interface.md.

The CLI reference doc exists and is accurate, but specs don't mention it.
EOF
```

**Step 4: Draft plan**
```bash
rooda run draft-plan-spec-chore --ai-cmd-alias kiro-cli
```

**Output:** Creates PLAN.md with task:
- Update specs/cli-interface.md to reference docs/cli-reference.md

**Step 5: Publish and build**
```bash
rooda run publish-plan --ai-cmd-alias kiro-cli
rooda run build --ai-cmd-alias kiro-cli
```

**Output:** Updates specs/cli-interface.md with cross-reference to docs/cli-reference.md

**Verification:**
```bash
# Check reference exists
grep "docs/cli-reference.md" specs/cli-interface.md
```

### Example 3: Doc Drift Workflow (Docs Outdated, Specs Specify Feature Not in Docs)

**Scenario:** --context flag specified in specs but not documented

**Step 1: Audit**
```bash
rooda run audit-spec-to-impl --ai-cmd-alias kiro-cli
```

**Output:** Gap report:
```markdown
# Spec to Implementation Gap Report

## Specified Features Not Documented

- `--context` flag specified in specs/cli-interface.md but not in docs/cli-reference.md
```

**Step 2: Convert to task**
```bash
cat > TASK.md << 'EOF'
# Document --context Flag

Add documentation for --context flag to docs/cli-reference.md.

Source truth:
- specs/cli-interface.md (design intent)
- cmd/rooda/main.go (actual implementation)
EOF
```

**Step 3: Draft plan**
```bash
rooda run draft-plan-impl-chore --ai-cmd-alias kiro-cli
```

**Output:** Creates PLAN.md with task:
- Add --context flag documentation to docs/cli-reference.md

**Step 4: Publish and build**
```bash
rooda run publish-plan --ai-cmd-alias kiro-cli
rooda run build --ai-cmd-alias kiro-cli
```

**Build process:**
1. Reads specs/cli-interface.md (design intent for --context)
2. Reads cmd/rooda/main.go (actual flag definition)
3. Reads docs/cli-reference.md (existing structure)
4. Synthesizes new section for --context flag
5. Applies humanizer skill
6. Writes updated docs/cli-reference.md
7. Commits: "docs: add CLI reference for --context flag"

**Verification:**
```bash
# Check flag documented
grep "\-\-context" docs/cli-reference.md

# Check example works
grep -A 5 "\-\-context" docs/cli-reference.md | grep "rooda run" | bash
```

### Example 4: Humanizer Skill Application

**Before humanization:**
```markdown
## Installation

It's important to note that rooda requires Go 1.24.5 or later. Simply run the following command to easily install rooda:

```bash
brew install rooda
```

This will seamlessly install rooda and leverage Homebrew's robust package management to ensure all dependencies are properly configured.
```

**After humanization:**
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
- Removed "easily" (filler word)
- Removed "seamlessly" (marketing language)
- Removed "leverage" (corporate jargon)
- Removed "robust" (marketing language)
- Made sentences direct and concise

### Example 5: Commit Message Generation

**Task:** Add CLI reference for --context flag

**Generated commit message:**
```
docs: add CLI reference for --context flag
```

**Task:** Update installation instructions

**Generated commit message:**
```
docs: update installation instructions
```

**Task:** Fix broken cross-reference in configuration.md

**Generated commit message:**
```
docs: fix broken cross-reference in configuration.md
```

**Format:** `docs: <action> <target>`

**Actions:** add, update, fix, remove

**Target:** Specific feature or file being documented

## Notes

### Design Rationale

**Why reuse existing procedures instead of creating doc-specific procedures?**
Documentation is implementation. The same workflow that maintains code (audit → plan → build) should maintain docs. Special doc procedures would create unnecessary complexity.

**Why bootstrap documents docs/ patterns?**
AGENTS.md is the source of truth for project structure. Documenting docs/ patterns there ensures all procedures know where docs live and what quality criteria apply.

**Why require humanizer skill for docs?**
AI-generated text has recognizable patterns that make docs feel robotic. Humanizer skill removes these patterns, making docs read naturally.

**Why synthesize docs from specs and code?**
Specs define design intent, code defines actual behavior. Docs must explain both. Synthesizing from both sources ensures accuracy.

**Why commit docs separately from code?**
Clear git history. Doc changes are distinct from code changes. Separate commits make it easier to review and revert.

**Why use "docs:" prefix in commit messages?**
Conventional commit format. Makes it easy to filter doc changes in git log.

**Why treat README.md as documentation?**
README is user-facing and must stay synchronized with specs/code. It follows the same lifecycle as docs/.

**Why exclude specs/ from doc lifecycle?**
Specs are design artifacts, not user-facing documentation. They follow a different lifecycle (spec-driven development).

### Source Truth Hierarchy for Docs

**For CLI documentation:**
1. cmd/rooda/main.go (actual flags and behavior)
2. specs/cli-interface.md (design intent)
3. docs/cli-reference.md (existing documentation)

**For procedure documentation:**
1. internal/procedures/builtin.go (actual procedures)
2. specs/procedures.md (design intent)
3. docs/procedures.md (existing documentation)

**For configuration documentation:**
1. internal/config/config.go (actual config schema)
2. specs/configuration.md (design intent)
3. docs/configuration.md (existing documentation)

**Agents resolve conflicts in favor of source truth (code > specs > existing docs).**

### Workflow Integration

**Bootstrap:**
- Creates/updates AGENTS.md with docs/ patterns
- Documents humanizer skill requirement
- No doc synthesis (only documents where docs live)

**Audit-impl-to-spec:**
- Finds docs that exist but aren't in specs
- Reports as "undocumented features"
- Output feeds into planning procedures

**Audit-spec-to-impl:**
- Finds features specified but not documented
- Reports as "missing documentation"
- Output feeds into planning procedures

**Draft-plan-*:**
- Converts doc gaps into tasks
- Specifies source truth files (specs, code)
- Defines acceptance criteria (quality checks)

**Publish-plan:**
- Imports doc tasks to work tracking
- No special handling for doc tasks

**Build:**
- Implements doc tasks like code tasks
- Reads source truth (specs, code, existing docs)
- Synthesizes new/updated doc content
- Applies humanizer skill
- Verifies quality criteria
- Commits with descriptive message

### Quality Criteria for Doc Tasks

**All doc tasks must pass:**
- Required sections exist
- Examples execute successfully
- Cross-references resolve
- No AI writing patterns
- Content matches source truth

**Build procedure verifies these before committing.**
