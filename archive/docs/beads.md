# Beads: External Memory for AI Agents

A distributed, git-backed issue tracker that provides persistent, structured memory for AI coding agents across sessions and context resets.

## Core Concept

**Beads solves the cross-session memory problem for AI agents.**

Markdown TODOs are write-only memory. After 200 messages or a new session, the agent either:
1. Hopes the TODO list is still in context
2. Asks the human to re-read it
3. Guesses

Beads provides **queryable, persistent memory** that survives context window resets. Dependencies are first-class data, not prose. Work state persists across sessions via git.

## The Three-Layer Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Layer                            │
│  bd create, ready, update, close, dep, sync              │
│  - Every command has --json output                       │
│  - Agent-optimized semantics                             │
└──────────────────────┬──────────────────────────────────┘
                       │
                       v
┌─────────────────────────────────────────────────────────┐
│                  SQLite Database                         │
│               (.beads/beads.db)                          │
│  - Local working copy (gitignored)                       │
│  - Fast queries, indexes, dependency graphs              │
│  - Blocked issues cache (25x speedup)                    │
└──────────────────────┬──────────────────────────────────┘
                       │
                  auto-sync
                  (5s debounce)
                       │
                       v
┌─────────────────────────────────────────────────────────┐
│                    JSONL File                            │
│                (.beads/issues.jsonl)                     │
│  - Git-tracked source of truth                           │
│  - One JSON line per issue                               │
│  - Merge-friendly (concurrent appends rarely conflict)   │
└──────────────────────┬──────────────────────────────────┘
                       │
                  git push/pull
                       │
                       v
┌─────────────────────────────────────────────────────────┐
│                  Remote Repository                       │
│  - Issues travel with code                               │
│  - No special sync server needed                         │
│  - Offline work just works                               │
└─────────────────────────────────────────────────────────┘
```

**Why this design?**
- SQLite for speed (queries in milliseconds)
- JSONL for git-friendliness (readable diffs, mergeable)
- Git for distribution (no coordination server)

## Hash-Based IDs: Zero-Collision Coordination

Sequential IDs (bd-1, bd-2, bd-3) cause collisions when multiple agents/branches create issues concurrently.

**The problem:**
```bash
Branch A: bd create "Add OAuth"   → bd-10
Branch B: bd create "Add Stripe"  → bd-10 (collision!)
```

**The solution:**
```bash
Branch A: bd create "Add OAuth"   → bd-a1b2 (from random UUID)
Branch B: bd create "Add Stripe"  → bd-f14c (different UUID, no collision)
```

Hash IDs scale progressively:
- 4 chars (0-500 issues): `bd-a1b2`
- 5 chars (500-1,500 issues): `bd-f14c3`
- 6 chars (1,500+ issues): `bd-3e7a5b`

**Hierarchical IDs** for epics: `bd-a3f8e9.1`, `bd-a3f8e9.2` provide human-readable structure while maintaining unique parent namespace.

## Dependency-Aware Execution

Dependencies are first-class data, not prose annotations.

### Dependency Types

| Type | Semantics | Affects `bd ready`? |
|------|-----------|---------------------|
| `blocks` | Issue X must close before Y starts | Yes |
| `parent-child` | Hierarchical (epic/subtask) | Yes (children blocked if parent blocked) |
| `related` | Soft link for reference | No |
| `discovered-from` | Found during work on parent | No |

### Ready Work Detection

```bash
bd ready --json
```

Returns issues with **no open blockers**. This is deterministic, offline, and completes in ~29ms (even on 10K issue databases).

**The cognitive load difference:** Instead of scanning markdown and mentally parsing "blocked by X", agents query structured data. No interpretation needed.

### Example Workflow

```bash
# Create issues with dependencies
bd create "Set up database" -p 1 -t task
# Returns: bd-a1b2

bd create "Create API" -p 2 -t feature
# Returns: bd-f14c

bd create "Add authentication" -p 2 -t feature
# Returns: bd-g25d

# Add dependencies (API depends on database)
bd dep add bd-f14c bd-a1b2

# Auth depends on API
bd dep add bd-g25d bd-f14c

# Query ready work
bd ready --json
# Returns: [{"id": "bd-a1b2", ...}]  (only database is ready)

# Work on it
bd update bd-a1b2 --status in_progress

# Complete it
bd close bd-a1b2 --reason "Database setup complete"

# Query again
bd ready --json
# Returns: [{"id": "bd-f14c", ...}]  (API is now unblocked)
```

## Molecules: Workflow Graphs

**Molecules = epics with execution semantics.**

Any epic with children is a molecule. Dependencies control execution flow:
- No dependency = parallel execution
- `blocks` dependency = sequential execution
- `parent-child` = hierarchical structure

### Three Phases (Optional Templates)

| Phase | Name | Storage | Synced | Purpose |
|-------|------|---------|--------|---------|
| **Solid** | Proto | `.beads/` | Yes | Frozen template |
| **Liquid** | Mol | `.beads/` | Yes | Active persistent work |
| **Vapor** | Wisp | `.beads/` | No | Ephemeral operations |

**Wisps are local-only:** Never exported to JSONL, never synced via git. Perfect for routine operations with no audit value.

### Bonding: Connecting Work Graphs

```bash
bd mol bond A B  # B depends on A (sequential)
```

Bonding creates dependencies between work graphs, enabling agents to traverse compound workflows across multiple sessions.

## Auto-Sync: Invisible Infrastructure

### Write Path

```
CLI Command → SQLite Write → Mark Dirty → 5s Debounce → Export to JSONL → Git Commit
```

**FlushManager** (event-driven, single-owner pattern):
- Batches multiple operations within 5-second window
- Incremental export (only changed issues)
- Full export after ID changes (e.g., prefix rename)
- Race-free via channel-based communication

### Read Path

```
git pull → Auto-Import Detection → Import to SQLite → Merge via Content Hashes → Query
```

**Auto-import** runs on first command after `git pull`. Hash-based comparison prevents false positives.

### Daemon Architecture

Each workspace runs its own background daemon (LSP-style):
- Socket at `.beads/bd.sock` (Windows: named pipes)
- Auto-starts on first command
- Handles RPC, auto-sync, background tasks
- Complete database isolation per workspace

**When to disable daemon:**
- Git worktrees (required: `bd --no-daemon`)
- CI/CD pipelines
- Resource-constrained environments

## Agent-Native Design

### Every Command Has --json

```bash
bd ready --json
bd create "Task" -p 1 --json
bd show bd-a1b2 --json
bd dep tree bd-a1b2 --json
```

Structured output, no text parsing needed.

### Discovery During Execution

```bash
# Agent discovers bug while implementing feature
bd create "Fix validation bug" -t bug -p 1 --deps discovered-from:bd-a1b2
```

The `discovered-from` dependency type maps to how agents actually work - discovering issues during implementation, not just planned work.

### Multi-Agent Coordination

```bash
# Agent 1
bd ready --assignee agent-1 --json
bd update bd-a1b2 --status in_progress --assignee agent-1

# Agent 2
bd ready --assignee agent-2 --json
bd update bd-f14c --status in_progress --assignee agent-2
```

Both agents query the same logical database via git. No coordination server needed.

### Session Persistence

**The test:** After using Beads, going back to markdown TODOs feels like trying to remember a phone number without writing it down.

Between sessions:
- Markdown TODO: Human must copy-paste list back to agent
- Beads: Agent runs `bd ready --json` and is immediately back in context

## Key Commands for Agents

### Core Workflow

```bash
# Query ready work
bd ready --json

# Create issue
bd create "Title" -p 0 -t task --json

# Add dependency
bd dep add <child-id> <parent-id>

# Update status
bd update <id> --status in_progress --assignee agent-name

# Close issue
bd close <id> --reason "Completed"

# Force immediate sync (bypasses 5s debounce)
bd sync
```

### Discovery and Exploration

```bash
# View dependency tree
bd dep tree <id> --json

# Show issue details and audit trail
bd show <id> --json

# List blocked issues
bd blocked --json

# View statistics
bd stats --json
```

### Multi-Agent Coordination

```bash
# Query by assignee
bd ready --assignee agent-name --json

# Query by priority
bd ready --priority 0 --json

# Query by label
bd ready --label backend --json
```

## File Structure

```
project-root/
├── .beads/
│   ├── beads.db          # SQLite (gitignored)
│   ├── issues.jsonl      # JSONL source of truth (git-tracked)
│   ├── bd.sock           # Daemon socket (gitignored)
│   └── config.yaml       # Project config (optional)
└── src/                  # Implementation
```

## Initialization

### Basic Setup

```bash
cd ~/project
bd init
```

Creates `.beads/` directory, database, and prompts for git hooks.

### Agent Setup (Non-Interactive)

```bash
bd init --quiet
```

Auto-installs hooks, no prompts. Perfect for agent-driven initialization.

### Contributor Mode (Fork Workflow)

```bash
bd init --contributor
```

Routes planning issues to separate repo (e.g., `~/.beads-planning`). Keeps experimental work out of PRs.

### Protected Branches

```bash
bd init --branch beads-sync
```

Commits to separate branch using internal worktree. Main branch stays clean.

### Stealth Mode (Local-Only)

```bash
bd init --stealth
```

Uses Beads locally without committing files to repo. Perfect for personal use on shared projects.

## Safety and Isolation

### Per-Project Isolation

Each project gets its own:
- `.beads/` directory
- SQLite database
- Daemon process
- JSONL file

No cross-project pollution. Issues cannot reference issues in other projects.

### Git Hooks (Recommended)

```bash
bd hooks install
```

Installs:
- **pre-commit** - Flushes pending changes before commit
- **post-merge** - Imports updated JSONL after pull
- **pre-push** - Exports database before push (prevents stale JSONL)
- **post-checkout** - Imports JSONL after branch checkout

**Why hooks matter:** Without pre-push hook, database changes can be committed locally but stale JSONL pushed to remote, causing multi-workspace divergence.

### Offline-First

- All queries run against local SQLite
- No network required for any commands
- Sync happens via git push/pull when online
- Full functionality available without internet

## Performance Characteristics

### Query Speed

- `bd ready`: ~29ms (with blocked issues cache)
- `bd list`: <100ms (thousands of issues)
- `bd show`: <10ms (single issue lookup)
- `bd dep tree`: <50ms (recursive CTE)

### Write Overhead

- SQLite write: <5ms (immediate)
- Export to JSONL: <50ms (incremental)
- Full export: <100ms (10K issues)
- Daemon debounce: 5 seconds (configurable)

### Blocked Issues Cache

Materialized cache table provides 25x speedup:
- Before cache: ~752ms (recursive CTE every query)
- With cache: ~29ms (NOT EXISTS check)
- Cache rebuild: <50ms (on dependency/status changes)

## Key Language Patterns

Beads-specific phrasing that matters:

- **"bd ready"** (not "what should I work on")
- **"query structured data"** (not "parse markdown")
- **"discovered-from dependency"** (not "found this bug")
- **"hash-based IDs prevent collisions"** (not "sequential IDs")
- **"auto-sync via git"** (not "manual export/import")
- **"external memory for agents"** (not "issue tracker")
- **"dependency-aware execution"** (not "task list")

## Why It Works

1. **Structured data** - Dependencies are queryable, not prose
2. **Git-backed** - No special sync server, works offline
3. **Hash-based IDs** - Zero-collision multi-agent coordination
4. **Auto-sync** - Invisible infrastructure, just works
5. **Agent-native** - Designed for how agents actually work
6. **Session persistence** - Memory survives context resets
7. **Offline-first** - No network required for queries

## The Meta-Observation

Beads isn't "issue tracking for agents" - it's **external memory for agents**, with dependency tracking and query capabilities that make it feel like a reliable extension of working memory across sessions.

The test is simple: After using Beads, going back to markdown TODOs feels like trying to remember a phone number without writing it down. Sure, you can do it for a little while, but why would you?

---

*Beads developed by [Steve Yegge](https://github.com/steveyegge/beads)*

*Agent perspective by Claude Sonnet 4.5*
