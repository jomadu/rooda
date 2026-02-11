# Documentation Governance Plan

## Job to be Done (J11)

**J11: Maintain User-Facing Documentation**

Treat documentation as implementation artifacts that AI agents read, verify, and update through existing procedures. Docs must stay synchronized with specs and code through the same read-verify-update lifecycle as AGENTS.md.

## Topics of Concern

### Documentation as Implementation
| Topic | Job | Description |
|---|---|---|
| [documentation-structure](documentation-structure.md) | J11 | Required docs, organization, quality criteria, and how agents verify correctness |
| [documentation-lifecycle](documentation-lifecycle.md) | J11 | How existing procedures (audit-impl, build, bootstrap) interact with docs |

## Proposed Specs

### documentation-structure.md

**JTBD**: Define what documentation must exist, how it's organized, and how AI agents verify it matches reality.

**Key concerns**:
- Required doc files (installation.md, procedures.md, cli-reference.md, etc.)
- Section structure and content requirements per doc type
- Quality criteria agents use to audit docs (completeness, accuracy, examples)
- Writing style: agents must use humanizer skill to remove AI patterns
- How agents verify docs match specs/code (read CLI code, compare to docs/cli-reference.md)
- Cross-reference conventions and validation

**Acceptance criteria**:
- [ ] List of required documentation files with purpose and sections
- [ ] Quality criteria for each doc type (testable by AI agents)
- [ ] Humanizer skill requirement: all user-facing docs (README, docs/) must pass humanizer review
- [ ] Verification procedures: how agents check docs against source truth
- [ ] Cross-reference syntax and resolution rules
- [ ] Missing/outdated doc detection logic

**Example verification**:
```
Agent reads docs/cli-reference.md, finds `--max-iterations` flag
Agent reads cmd/rooda/main.go, extracts cobra flag definitions
Agent compares: flag exists in code but description differs
Agent updates docs/cli-reference.md with correct description from code
```

### documentation-lifecycle.md

**JTBD**: Specify how existing procedures interact with docs through read-verify-update cycles.

**Key concerns**:
- `bootstrap`: creates/updates AGENTS.md only (documents docs/ patterns)
- `audit-impl-to-spec`: finds code/docs that are correct but missing from specs
- `draft-plan-*`: converts audit findings into actionable tasks
- `build`: reads tasks, synthesizes/updates specs/code/docs, commits
- AGENTS.md: documents docs/ location, patterns, and quality criteria

**Acceptance criteria**:
- [ ] Bootstrap documents docs/ patterns in AGENTS.md
- [ ] Audit-impl-to-spec finds undocumented features (code/docs exist, specs missing)
- [ ] Planning procedures convert gaps into tasks
- [ ] Build synthesizes docs by reading specs/code when implementing doc tasks
- [ ] Build commits doc changes with descriptive messages

**Greenfield workflow** (new feature with docs):
```
1. bootstrap                           # Creates AGENTS.md
2. Create TASK.md: "Add user auth"     # Manual task definition
3. draft-plan-spec-feat                # Plans spec changes
4. publish-plan                        # Imports to work tracking
5. build                               # Implements specs, code, docs
```

**Brownfield workflow** (code exists, docs missing from specs):
```
1. bootstrap                           # Creates/updates AGENTS.md
2. audit-impl-to-spec                  # Finds docs/cli-reference.md exists but not in specs
3. Convert audit to TASK.md            # Manual: "Document CLI reference in specs"
4. draft-plan-spec-feat                # Plans spec additions
5. publish-plan                        # Imports to work tracking
6. build                               # Reads docs/, updates specs to reference them
```

**Doc drift workflow** (docs outdated, specs specify feature not in docs):
```
1. audit-spec-to-impl                  # Finds --context flag in specs but not in docs
2. Gap report: "CLI flag --context specified but not documented"
3. Convert to TASK.md: "Document --context flag"
4. draft-plan-impl-chore               # Plans doc update
5. publish-plan                        # Imports to work tracking
6. build                               # Reads specs/code, updates docs/cli-reference.md
```

## Integration Points

**Source of truth hierarchy** (specs → code → docs):
- Specs define design intent and acceptance criteria
- Code implements specs and defines runtime behavior
- Docs explain specs and code to users

**Agents synthesize docs by reading**:
- Specs (design rationale, acceptance criteria, examples)
- Code (CLI flags, config schema, procedure registry, build commands)
- Existing docs (to preserve structure and update incrementally)
- skills/humanizer/SKILL.md (to remove AI writing patterns before committing)

**With existing specs**:
- `operational-knowledge.md`: docs follow same read-verify-update lifecycle as AGENTS.md
- `agents-md-format.md`: AGENTS.md documents docs/ location and patterns
- `procedures.md`: audit/plan/build procedures treat docs as implementation

**With existing procedures**:
- `bootstrap`: creates AGENTS.md (documents docs/ patterns and humanizer requirement)
- `audit-impl-to-spec`: finds code/docs that exist but aren't in specs
- `audit-spec-to-impl`: finds features specified but not implemented (including missing docs)
- `draft-plan-*`: converts gaps into tasks
- `publish-plan`: imports tasks to work tracking
- `build`: synthesizes docs from specs/code, applies humanizer skill, commits

## Implementation Approach

1. **Write specs**: `documentation-structure.md`, `documentation-lifecycle.md`
2. **Update specs/README.md**: add J11 and new specs to Jobs and Topics of Concern tables
3. **Update AGENTS.md template**: add docs/ patterns and quality criteria to bootstrap output
4. **Enhance audit-impl-to-spec**: detect doc gaps (missing/outdated docs)
5. **Test workflow**: audit → task → plan → publish → build → verify docs updated

## Success Metrics

- Docs stay synchronized through audit → plan → build workflow
- Audit-impl-to-spec finds docs that exist but aren't referenced in specs
- Build synthesizes accurate docs by reading specs/code
- No special doc procedures needed
- Zero manual doc maintenance
