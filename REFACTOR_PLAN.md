# Refactor Plan: Consolidate to TASK.md and ooda-config.yml

## Objective
Collapse STORY.md and BUG.md into single TASK.md file, and refactor ooda-procedures.yml to ooda-config.yml with configurable file paths.

## Tasks

### 1. Rename prompt files (6 files)
- [x] `observe_story_specs_impl.md` → `observe_story_task_specs_impl.md`
- [x] `observe_bug_specs_impl.md` → `observe_bug_task_specs_impl.md`
- [x] `orient_story_incorporation.md` → `orient_story_task_incorporation.md`
- [x] `orient_bug_incorporation.md` → `orient_bug_task_incorporation.md`
- [x] `decide_story_plan.md` → `decide_story_task_plan.md`
- [x] `decide_bug_plan.md` → `decide_bug_task_plan.md`

### 2. Update prompt file content (6 files)
Change references from "story file" / "bug file" to "task file":
- [x] `observe_story_task_specs_impl.md`
- [x] `observe_bug_task_specs_impl.md`
- [x] `orient_story_task_incorporation.md`
- [x] `orient_bug_task_incorporation.md`
- [x] `decide_story_task_plan.md`
- [x] `decide_bug_task_plan.md`

### 3. Rename config file
- [x] `ooda-procedures.yml` → `ooda-config.yml`

### 4. Update ooda-config.yml structure
- [x] Add `paths:` section with task_dir, task_file, plan_file
- [x] Update procedure file paths for renamed prompts (plan-story-to-spec, plan-bug-to-spec)
- [x] Keep existing procedures structure

### 5. Update ooda.sh
- [x] Parse `paths:` section from config
- [x] Replace STORY_FILE/BUG_FILE with single TASK_FILE
- [x] Support {task-id} interpolation in paths
- [x] Add CLI flags: `--task-file`, `--plan-file` (override config)
- [x] Update `--procedures` flag to `--config`
- [x] Update context section in prompt template
- [x] Update usage/help text

### 6. Update documentation
- [x] **README.md**: STORY.md/BUG.md → TASK.md, config file rename, new CLI flags
- [x] **prompts/README.md**: Update procedure compositions table, file references
- [x] **specs.md**: No references found (no changes needed)

## Notes
- AGENTS.md and specs/ locations remain discovered by prompts (not in config)
- Backward compatible: default paths match current behavior
- CLI overrides allow flexible task tracking integration
