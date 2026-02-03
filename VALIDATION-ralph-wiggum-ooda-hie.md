# Validation: ralph-wiggum-ooda-hie

**Task:** Implement missing procedure error messages  
**Date:** 2026-02-03  
**Status:** ✅ PASSED

## Acceptance Criteria

### AC1: Error message lists available procedures
**Status:** ✅ PASSED (already implemented)

**Test:**
```bash
./src/rooda.sh nonexistent-procedure
```

**Expected:** Error message includes list of available procedures

**Actual:**
```
Error: Procedure 'nonexistent-procedure' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml

Available procedures:
  - bootstrap
  - build
  - draft-plan-story-to-spec
  - draft-plan-bug-to-spec
  - draft-plan-spec-to-impl
  - draft-plan-impl-to-spec
  - draft-plan-spec-refactor
  - draft-plan-impl-refactor
  - publish-plan
```

### AC2: Suggests closest match (fuzzy matching)
**Status:** ✅ PASSED (newly implemented)

**Test Case 1: Typo in procedure name**
```bash
./src/rooda.sh bootstra
```

**Expected:** Suggests "bootstrap"

**Actual:**
```
Error: Procedure 'bootstra' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml

Did you mean: bootstrap

Available procedures:
  - bootstrap
  ...
```

**Test Case 2: Similar word**
```bash
./src/rooda.sh bild
```

**Expected:** Suggests "build"

**Actual:**
```
Error: Procedure 'bild' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml

Did you mean: build

Available procedures:
  - build
  ...
```

**Test Case 3: Partial match**
```bash
./src/rooda.sh plan
```

**Expected:** Suggests a procedure containing "plan"

**Actual:**
```
Error: Procedure 'plan' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml

Did you mean: draft-plan-story-to-spec

Available procedures:
  ...
```

**Test Case 4: Unrelated input (threshold test)**
```bash
./src/rooda.sh xyz
```

**Expected:** No suggestion (score too low)

**Actual:**
```
Error: Procedure 'xyz' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml

Available procedures:
  - bootstrap
  ...
```

### AC3: Includes config file path in error
**Status:** ✅ PASSED (already implemented)

**Test:**
```bash
./src/rooda.sh nonexistent-procedure
```

**Expected:** Error message includes full path to config file

**Actual:**
```
Error: Procedure 'nonexistent-procedure' not found in /Users/maxdunn/Dev/ralph-wiggum-ooda/src/rooda-config.yml
```

## Implementation Details

**File:** `src/rooda.sh`  
**Lines:** 135-181 (validate_config function)

**Algorithm:**
1. When procedure not found, iterate through all available procedures
2. For each procedure, calculate similarity score:
   - Substring match: score = 100 (high priority)
   - Character overlap: score = count of matching characters
3. Track best match with highest score
4. Only suggest if score >= 3 (threshold to avoid unrelated suggestions)
5. Display suggestion before listing all available procedures

**Bash Compatibility:**
- Uses `tr '[:upper:]' '[:lower:]'` instead of `${var,,}` for older bash versions
- Avoids bash 4.0+ specific features

## Shellcheck Results

```bash
shellcheck src/rooda.sh
```

**Output:**
```
In src/rooda.sh line 179:
            echo "$available_procs" | sed 's/^/  - /'
            ^-- SC2001 (style): See if you can use ${variable//search/replace} instead.
```

**Note:** This is a style suggestion, not an error. The pattern is used consistently throughout the codebase.

## Conclusion

All acceptance criteria met. Fuzzy matching successfully suggests closest procedure name when user makes typos or enters partial names. Threshold prevents suggesting unrelated procedures for completely invalid input.
