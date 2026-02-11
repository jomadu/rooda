# Emit Success

You must output a SUCCESS promise when the procedure's goal is fully achieved.

For work tracking procedures (build): Signal SUCCESS only when no ready work remains.
For planning/auditing procedures: Signal SUCCESS only when output is minimal, complete, and accurate.
For single-task procedures (sync): Signal SUCCESS when the task completes.

Actions:
- Output the exact signal: <promise>SUCCESS</promise>
- After the signal, include summary of what was accomplished
- List files modified or created
- Note any follow-up actions needed

Example (build procedure, no work remaining):
<promise>SUCCESS</promise>

All ready work completed:
- Implemented task #42: Add user authentication
- All tests passing
- No ready tasks remaining in work tracking

Example (audit procedure, validated output):
<promise>SUCCESS</promise>

Audit completed and validated:
- Reviewed 12 specification files
- Generated audit report at docs/audit-2024-01-15.md
- Report is minimal yet complete and accurate
- Found 3 issues requiring attention

Example (draft plan procedure, validated output):
<promise>SUCCESS</promise>

Draft plan completed and validated:
- Created plan at docs/draft-plan-auth-feature.md
- Plan is minimal yet complete and accurate
- Broken down into 8 actionable tasks
- Ready for import to work tracking
