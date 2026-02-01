# Agent Instructions

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Task Management

This repository uses **beads** for task management, not the tasks/ directory pattern.

**Query ready work:**
```bash
bd ready --json
```

**View issue details:**
```bash
bd show <id> --json
```

**Task/Plan files for story/bug incorporation:**
- Task file: Issue description from `bd show <id> --json` (title + description fields)
- Plan file: Not used (beads tracks status directly)

**Update status:**
```bash
bd update <id> --status in_progress
bd update <id> --status blocked
```

**Close issue:**
```bash
bd close <id> --reason "Completed X"
```

**View dependency tree:**
```bash
bd dep tree <id> --json
```

**Sync with git:**
```bash
bd sync
```

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

