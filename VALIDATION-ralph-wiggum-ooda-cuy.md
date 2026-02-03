# Validation: Config File Resolution Relative to Script Location

**Issue:** ralph-wiggum-ooda-cuy  
**Date:** 2026-02-03  
**Validator:** Kiro AI Agent

## Acceptance Criteria

- [x] Config resolves correctly when script invoked from different directories
- [x] Test case verifies resolution behavior
- [ ] **Bug found:** OODA file paths from config don't resolve relative to script location

## Test Cases

### Test 1: Invoke from Project Root

**Command:**
```bash
cd /Users/maxdunn/Dev/ralph-wiggum-ooda
./src/rooda.sh bootstrap --max-iterations 0
```

**Expected Behavior:**
- Script finds config at `src/rooda-config.yml` relative to script location
- Display shows resolved OODA phase files
- Script proceeds to iteration loop

**Result:** ✅ PASS

**Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Procedure: bootstrap
Observe:   src/components/observe_bootstrap.md
Orient:    src/components/orient_bootstrap.md
Decide:    src/components/decide_bootstrap.md
Act:       src/components/act_bootstrap.md
Branch:    main
Max:       1 iterations
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Verification:**
- Config file found successfully
- Procedure resolved from config
- Script executed without config-related errors

### Test 2: Invoke from Parent Directory

**Command:**
```bash
cd /Users/maxdunn/Dev
./ralph-wiggum-ooda/src/rooda.sh bootstrap --max-iterations 0
```

**Expected Behavior:**
- Script finds config at `src/rooda-config.yml` relative to script location (not relative to current directory)
- Display shows resolved OODA phase files
- Script proceeds to iteration loop

**Result:** ❌ FAIL

**Output:**
```
Error: File not found: src/components/observe_bootstrap.md
```

**Verification:**
- Config file found successfully (config resolution works)
- Procedure resolved from config (OODA paths loaded)
- **Bug discovered:** OODA phase file paths are relative to config location but validated relative to current directory
- File validation at line 107 checks `[ ! -f "$file" ]` which uses current directory, not script directory

### Test 3: Invoke with Absolute Path

**Command:**
```bash
cd /tmp
/Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda.sh bootstrap --max-iterations 0
```

**Expected Behavior:**
- Script finds config at `src/rooda-config.yml` relative to script location (not relative to /tmp)
- Display shows resolved OODA phase files
- Script proceeds to iteration loop

**Result:** ❌ FAIL

**Output:**
```
Error: File not found: src/components/observe_bootstrap.md
```

**Verification:**
- Config file found successfully (config resolution works)
- Procedure resolved from config (OODA paths loaded)
- **Bug discovered:** OODA phase file paths are relative to config location but validated relative to current directory
- File validation at line 107 checks `[ ! -f "$file" ]` which uses current directory, not script directory

### Test 4: Invoke from src/ Directory

**Command:**
```bash
cd /Users/maxdunn/Dev/ralph-wiggum-ooda/src
./rooda.sh bootstrap --max-iterations 0
```

**Expected Behavior:**
- Script finds config at `rooda-config.yml` in same directory as script
- Display shows resolved OODA phase files
- Script proceeds to iteration loop

**Result:** ❌ FAIL

**Output:**
```
Error: File not found: src/components/observe_bootstrap.md
```

**Verification:**
- Config file found successfully (config resolution works)
- Procedure resolved from config (OODA paths loaded)
- **Bug discovered:** OODA phase file paths are relative to config location but validated relative to current directory
- When invoked from src/, the paths like "src/components/..." don't exist (would need to be "../src/components/..." or just "components/...")

## Implementation Analysis

**Source:** `src/rooda.sh` lines 27-28

```bash
# Resolve config file relative to script location
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/rooda-config.yml"
```

**How it works:**
1. `dirname "$0"` extracts the directory path from the script invocation path
2. `cd "$(dirname "$0")" && pwd` changes to that directory and prints the absolute path
3. `CONFIG_FILE="${SCRIPT_DIR}/rooda-config.yml"` constructs config path relative to script location

**Why this works:**
- `$0` contains the path used to invoke the script (relative or absolute)
- `dirname` extracts just the directory portion
- `cd` and `pwd` resolve to absolute path, handling symlinks and relative paths
- Config path is always relative to script location, not current working directory

**Validation:**
- ✅ Works with relative invocation from project root
- ✅ Works with relative invocation from parent directory
- ✅ Works with absolute path invocation from any directory
- ✅ Works when invoked from script's own directory

## Conclusion

**Acceptance criteria NOT fully met:**

1. **Config resolves correctly when script invoked from different directories** - ✅ PASS - Config file resolution works correctly
2. **Test case verifies resolution behavior** - ❌ FAIL - Testing revealed a bug in OODA file path resolution

**Bug discovered:** While the config file itself resolves correctly relative to script location (lines 27-28), the OODA phase file paths loaded from the config are not resolved relative to the script directory. They are validated and used relative to the current working directory (line 107).

**Root cause:**
- Config paths in `rooda-config.yml` are relative (e.g., `src/components/observe_bootstrap.md`)
- These paths are stored as-is in variables (lines 77-80)
- File validation at line 107 uses `[ ! -f "$file" ]` which checks relative to current directory
- This only works when invoked from project root where `src/components/` exists

**Impact:**
- Script only works when invoked from project root directory
- Fails when invoked from parent directory, absolute path from elsewhere, or from src/ directory
- This contradicts cli-interface.md AC line 3 which specifies config should resolve relative to script location

## Recommendations

**Fix required:** OODA phase file paths from config should be resolved relative to `SCRIPT_DIR`, not current directory.

**Suggested implementation:**
```bash
# After loading from config (around line 80)
if [ -n "$PROCEDURE" ]; then
    # ... existing config loading ...
    
    # Resolve OODA paths relative to script directory
    OBSERVE="${SCRIPT_DIR}/${OBSERVE}"
    ORIENT="${SCRIPT_DIR}/${ORIENT}"
    DECIDE="${SCRIPT_DIR}/${DECIDE}"
    ACT="${SCRIPT_DIR}/${ACT}"
fi
```

**Alternative:** Update config paths to be absolute or use a different path resolution strategy.

**Related issue:** This bug affects all procedure-based invocations. Explicit flag invocations work correctly if absolute paths are provided.
