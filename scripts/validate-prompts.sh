#!/usr/bin/env bash
set -euo pipefail

# Validate prompt file structure per component-authoring.md

PROMPTS_DIR="src/prompts"
EXIT_CODE=0

echo "Validating prompt files in $PROMPTS_DIR..."

for file in "$PROMPTS_DIR"/*.md; do
    filename=$(basename "$file")
    echo "Checking $filename..."
    
    # Check phase header: # [Phase]: [Purpose]
    if ! grep -q '^# \(Observe\|Orient\|Decide\|Act\):' "$file"; then
        echo "  ERROR: Missing or invalid phase header (must be: # Observe:|Orient:|Decide:|Act:)"
        EXIT_CODE=1
        continue
    fi
    
    # Determine expected phase code prefix
    phase=$(grep '^# \(Observe\|Orient\|Decide\|Act\):' "$file" | head -1 | sed 's/^# \([^:]*\):.*/\1/')
    case "$phase" in
        Observe) prefix="O" ;;
        Orient) prefix="R" ;;
        Decide) prefix="D" ;;
        Act) prefix="A" ;;
    esac
    
    # Check step headers match phase
    while IFS= read -r line; do
        # Extract step code (e.g., O1, R5, A3.5)
        code=$(echo "$line" | sed 's/^## \([^:]*\):.*/\1/')
        
        # Check if code starts with correct prefix
        if [[ ! "$code" =~ ^${prefix}[0-9]+(\.[0-9]+)?$ ]]; then
            echo "  ERROR: Step code '$code' doesn't match phase '$phase' (expected ${prefix}1-${prefix}99)"
            EXIT_CODE=1
        fi
    done < <(grep '^## [A-Z][0-9]' "$file")
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "All prompt files valid"
else
    echo "Validation failed"
fi

exit $EXIT_CODE
