# Orient: Quality Assessment

## R18: Apply Boolean Criteria per AGENTS.md

Apply the quality criteria defined in AGENTS.md:
- What boolean criteria apply? (specs or implementation)
- What are the specific thresholds?
- What metrics or checks are defined?
- Evaluate each criterion systematically
- Document findings for each criterion

## R19: Identify Human Markers

Search for indicators of quality issues:
- TODOs, FIXMEs, HACKs in code or specs
- Code smells (duplication, long functions, complex conditionals)
- Unclear language in specs (ambiguous, vague, contradictory)
- Missing documentation
- Incomplete implementations
- Inconsistent patterns
- Dead code or unused specs

## R20: Score Each Criterion PASS/FAIL

For each quality criterion:
- PASS: criterion is met, no action needed
- FAIL: criterion is not met, refactoring required

Use boolean scoring only - no subjective grades or percentages.

Determine overall assessment:
- If all criteria PASS: no refactoring needed
- If any criteria FAIL: refactoring plan required
- Identify which specific criteria failed
- Prioritize by impact and severity
