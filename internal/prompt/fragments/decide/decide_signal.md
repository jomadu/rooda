# Decide Signal

Decide which signal to emit based on blockers, completion, and remaining work.

## Emit FAILURE if Blocked

- Missing information, tools, dependencies, or permissions
- Work tracking unavailable
- Conflicting requirements

## Emit SUCCESS if Complete

- **build:** No ready work remains
- **Planning/auditing:** Output complete
- **sync:** Task complete

## Continue if Work Remains

No blockers and work exists: continue iterating (no signal).
