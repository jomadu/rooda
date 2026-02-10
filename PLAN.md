# Plan: Update Phase Marker Format from `#` Headers to `===` Separators

## Objective

Change OODA phase markers in assembled prompts from markdown headers (`# OBSERVE`) to visual separators (`=== OBSERVE ===`) to prevent collision with fragment content headers.

## Issues Identified

1. **Phase marker format collision**: `#` headers used for both OODA phases and fragment content
2. **Context file not interpolated**: `--context TASK.md` passes literal string "TASK.md" instead of showing source path and file content

## Design Decisions

### Phase Marker Format
Use `=== PHASE ===` separators instead of `#` headers to avoid collision with fragment content headers.

### Context Display Format
When context is from a file:
```
=== CONTEXT ===
Source: ./TASK.md

[file content]
```

When context is inline:
```
=== CONTEXT ===
[inline text]
```

Rationale: Provides both provenance (source path) and immediate content access for agents, with minimal overhead.

## Affected Specifications

### Primary Changes Required

1. **prompt-composition.md**
   - Update "Assembled Prompt Structure" section to show `===` separators
   - Update `AssemblePrompt` algorithm:
     - Use `prompt += "=== " + phase.toUpperCase() + " ===\n"` for phase markers
     - Add context source path when content is from file: `"Source: " + filePath + "\n\n" + content`
   - Update all examples showing assembled prompts
   - Update "Context Injection" examples to show both file-based and inline context formats
   - Update algorithm to show file existence heuristic and content interpolation

2. **iteration-loop.md**
   - Update all example outputs showing assembled prompts
   - Examples 4, 5, 7, 8 show dry-run output with prompt structure
   - Update "Dry-Run Mode" example sections
   - Update Example 5 to show context file with "Source:" line

3. **observability.md**
   - Update "Dry-Run Mode (Validation Only)" example
   - Shows assembled prompt with section markers

4. **cli-interface.md**
   - Update "Dry Run" examples showing assembled prompt output
   - Update "Dry-Run Validation (Success)" example
   - Update context injection examples to show file vs inline handling
   - Clarify file existence heuristic behavior in examples

### Secondary Review (May Not Need Changes)

5. **ai-cli-integration.md**
   - Review: mentions prompts are piped to AI CLI but doesn't show format
   - Likely no changes needed

6. **error-handling.md**
   - Review: focuses on failure detection, not prompt format
   - Likely no changes needed

7. **procedures.md**
   - Review: defines fragment arrays but not final assembly format
   - Likely no changes needed

8. **configuration.md**
   - Review: defines config schema, not prompt format
   - Likely no changes needed

9. **agents-md-format.md**, **operational-knowledge.md**, **distribution.md**
   - Review: unrelated to prompt assembly
   - No changes needed

## Implementation Tasks

1. ✅ Update prompt-composition.md:
   - ✅ Algorithm for phase markers (`===`)
   - ✅ Algorithm for context source path injection
   - ✅ All examples showing assembled prompts
   - ✅ Context injection examples (file vs inline)
2. ✅ Update iteration-loop.md dry-run examples
3. ✅ Update observability.md dry-run examples
4. ✅ Update cli-interface.md:
   - ✅ No changes needed (no prompt examples found)
5. ✅ Verify no other specs reference the `#` header format

## Acceptance Criteria

- [x] All specs show `=== PHASE ===` format for OODA phase markers
- [x] Fragment content headers remain as `#` markdown headers
- [x] CONTEXT section uses `=== CONTEXT ===` format
- [x] All dry-run examples updated consistently
- [x] No references to old `# OBSERVE` format remain in specs
- [x] Context from file shows: `Source: <path>\n\n<content>`
- [x] Context inline shows: `<content>` (no Source line)
- [x] File existence heuristic behavior clearly documented
- [x] Algorithm shows context file reading and interpolation
