# Task: Consolidate src/README.md into Specs

## Goal
Remove `src/README.md` by folding its content into existing specs, eliminating the developer/user documentation split.

## Background
`src/README.md` currently documents the prompt composition system (`create_prompt()` function, prompt file structure, assembly process). This content belongs in specs, not as a separate developer reference. Multiple specs reference `src/README.md`:

- `user-documentation.md` (3 references)
- `iteration-loop.md` (1 reference)
- `configuration-schema.md` (1 reference)
- `ai-cli-integration.md` (1 reference)
- `component-authoring.md` (1 reference)

## Approach

### 1. Merge Content into component-authoring.md
`specs/component-authoring.md` already covers prompt authoring. Add sections for:
- The `create_prompt()` function implementation
- How prompt files are assembled into final prompt
- The heredoc mechanism and command substitution
- Assembled prompt structure

### 2. Update Cross-References
Replace all `src/README.md` references with `component-authoring.md`:
- `user-documentation.md` line 113, 176, 201
- `iteration-loop.md` line 95
- `configuration-schema.md` line 100
- `ai-cli-integration.md` line 114
- `component-authoring.md` line 129

### 3. Update AGENTS.md
Remove `src/README.md` from implementation definition since it will no longer exist.

### 4. Delete src/README.md
Once content is merged and references updated, delete the file.

### 5. Move audit-links.sh to scripts/
Create `scripts/` directory and move `audit-links.sh` into it. Update any references or documentation that mentions the script location.

## Acceptance Criteria
- [ ] All content from `src/README.md` incorporated into `component-authoring.md`
- [ ] All cross-references updated to point to `component-authoring.md`
- [ ] `src/README.md` deleted
- [ ] AGENTS.md implementation definition updated
- [ ] No broken links in specs or docs
- [ ] `component-authoring.md` maintains coherent flow with new content
- [ ] `scripts/` directory created
- [ ] `audit-links.sh` moved to `scripts/audit-links.sh`
- [ ] Any references to `audit-links.sh` updated with new path
