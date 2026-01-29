#!/bin/bash
# Usage: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]
# Example:
#   ./ooda.sh --observe prompts/observe_specs.md \
#             --orient prompts/orient_gap_analysis.md \
#             --decide prompts/decide_planning.md \
#             --act prompts/act_plan.md \
#             --max-iterations 10

# Parse arguments
OBSERVE=""
ORIENT=""
DECIDE=""
ACT=""
MAX_ITERATIONS=0

while [[ $# -gt 0 ]]; do
    case $1 in
        --observe)
            OBSERVE="$2"
            shift 2
            ;;
        --orient)
            ORIENT="$2"
            shift 2
            ;;
        --decide)
            DECIDE="$2"
            shift 2
            ;;
        --act)
            ACT="$2"
            shift 2
            ;;
        --max-iterations)
            MAX_ITERATIONS="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
            exit 1
            ;;
    esac
done

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    echo "Usage: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
    exit 1
fi

# Validate files exist
for file in "$OBSERVE" "$ORIENT" "$DECIDE" "$ACT"; do
    if [ ! -f "$file" ]; then
        echo "Error: File not found: $file"
        exit 1
    fi
done

ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Observe: $OBSERVE"
echo "Orient:  $ORIENT"
echo "Decide:  $DECIDE"
echo "Act:     $ACT"
echo "Branch:  $CURRENT_BRANCH"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:     $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create interpolated prompt template
create_prompt() {
    cat <<EOF
# OODA Loop Iteration

## OBSERVE
$(cat "$OBSERVE")

## ORIENT
$(cat "$ORIENT")

## DECIDE
$(cat "$DECIDE")

## ACT
$(cat "$ACT")
EOF
}

while true; do
    if [ $MAX_ITERATIONS -gt 0 ] && [ $ITERATION -ge $MAX_ITERATIONS ]; then
        echo "Reached max iterations: $MAX_ITERATIONS"
        break
    fi

    create_prompt | kiro-cli chat --no-interactive --trust-all-tools

    git push origin "$CURRENT_BRANCH" || {
        echo "Failed to push. Creating remote branch..."
        git push -u origin "$CURRENT_BRANCH"
    }

    ITERATION=$((ITERATION + 1))
    echo -e "\n\n======================== LOOP $ITERATION ========================\n"
done
