# Write AGENTS.md

You must update the AGENTS.md file with the corrected information. Use the file writing tool.

Update the AGENTS.md file with corrected information.

Actions:
- Apply determined updates to appropriate sections
- Maintain existing structure and formatting
- Preserve operational learnings
- Update timestamps or version information
- Ensure all commands are accurate

Required sections in AGENTS.md:
- Work Tracking System
- Quick Reference
- Task Input (or Story/Bug Input)
- Planning System
- Build/Test/Lint Commands
- Specification Definition
- Implementation Definition
- Documentation Definition (if docs/ directory exists)
- Audit Output
- Quality Criteria
- Operational Learnings

Documentation Definition section format (include only if docs/ directory exists):
```markdown
## Documentation Definition

**Location:** `docs/*.md`

**Patterns:**
- `docs/*.md` — User-facing documentation files

**Required files:**
- `docs/installation.md` — Installation instructions
- `docs/procedures.md` — Procedure documentation
- `docs/configuration.md` — Configuration guide
- `docs/cli-reference.md` — CLI reference
- `docs/troubleshooting.md` — Troubleshooting guide
- `docs/agents-md.md` — AGENTS.md format documentation

**Quality criteria:**
- All documentation examples execute successfully (PASS/FAIL)
- No AI writing patterns detected (PASS/FAIL)
- All cross-references resolve (PASS/FAIL)
```
