# Emit Failure

You must output the exact FAILURE signal to indicate the procedure cannot proceed.

Output the exact signal format, then explain why work is blocked.

Actions:
- Output exactly: `<promise>FAILURE</promise>`
- After the signal, explain why work is blocked
- List missing prerequisites or dependencies
- Suggest what needs to happen to unblock
- Provide actionable next steps

Example:
```
<promise>FAILURE</promise>

Cannot proceed: Missing authentication module specification. The OAuth2 integration requires a detailed spec defining token refresh behavior and error handling patterns.

Next steps:
1. Create specs/auth-oauth2.md with token lifecycle specification
2. Define error handling patterns for expired tokens
3. Document refresh token rotation policy
```

Note: The signal `<promise>FAILURE</promise>` must be exact. The loop orchestrator scans for this exact format to determine iteration outcome.
