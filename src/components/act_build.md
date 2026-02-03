# Act: Build

## A3: Implement Using Parallel Subagents

Execute the implementation plan:
- Spawn subagents for parallel work as determined in decide phase
- **Critical: Only 1 subagent for build/tests** (avoid parallel test conflicts)
- Each subagent handles its assigned files/tasks
- Coordinate file modifications to avoid conflicts
- Complete all file modifications before proceeding

## A4: Run Tests per AGENTS.md (Backpressure)

Execute test commands from AGENTS.md:
- Run all relevant tests
- Verify tests pass
- If tests fail: fix issues, don't proceed
- This is backpressure - quality gate before commit

## A5: Update Work Tracking per AGENTS.md

Update the work tracking system:
- Mark task as complete if fully done
- Update status if partially done or blocked
- Use commands specified in AGENTS.md
- Document what was accomplished

## A6: Update AGENTS.md if Learned Something New

If operational learnings occurred:
- Update AGENTS.md with new information
- Capture the why - document rationale
- Keep it up to date
- Examples: commands that didn't work, better patterns discovered, new conventions

## A7: Commit When Tests Pass

Commit all changes with descriptive message:
- What was implemented
- What task was completed
- Reference task ID from work tracking
- Only commit when tests pass (backpressure enforced)
