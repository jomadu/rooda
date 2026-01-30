#!/bin/bash
# Usage: ./ooda.sh [task] [--tasks <file>] [--max-iterations N]
#    OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]
# Examples:
#   ./ooda.sh build
#   ./ooda.sh build --max-iterations 5
#   ./ooda.sh build --tasks custom-tasks.yml
#   ./ooda.sh --observe prompts/observe_specs.md \
#             --orient prompts/orient_gap.md \
#             --decide prompts/decide_gap_plan.md \
#             --act prompts/act_plan.md \
#             --max-iterations 1

# Parse YAML (minimal parser for our simple structure)
parse_yaml() {
    local file=$1
    local task=$2
    local field=$3
    local in_task=false
    
    while IFS= read -r line; do
        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// }" ]] && continue
        
        # Check if we're entering the target task
        if [[ "$line" =~ ^${task}: ]]; then
            in_task=true
            continue
        fi
        
        # Check if we've entered a different task
        if [[ "$line" =~ ^[a-z-]+: ]] && [[ ! "$line" =~ ^${task}: ]]; then
            in_task=false
            continue
        fi
        
        # If we're in the target task, look for the field
        if [ "$in_task" = true ]; then
            if [[ "$line" =~ ^[[:space:]]+${field}:[[:space:]]*(.+)$ ]]; then
                echo "${BASH_REMATCH[1]}"
                return 0
            fi
        fi
    done < "$file"
    
    return 1
}

# Parse arguments
OBSERVE=""
ORIENT=""
DECIDE=""
ACT=""
MAX_ITERATIONS=0
TASK=""
TASKS_FILE="ooda-tasks.yml"

# Check if first arg is a positional task name (not a flag)
if [[ $# -gt 0 ]] && [[ ! "$1" =~ ^-- ]]; then
    TASK="$1"
    shift
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --tasks)
            TASKS_FILE="$2"
            shift 2
            ;;
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
            echo "Usage: ./ooda.sh [task] [--tasks <file>] [--max-iterations N]"
            echo "   OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
            exit 1
            ;;
    esac
done

# If task specified, load from config
if [ -n "$TASK" ]; then
    if [ ! -f "$TASKS_FILE" ]; then
        echo "Error: $TASKS_FILE not found"
        exit 1
    fi
    
    OBSERVE=$(parse_yaml "$TASKS_FILE" "$TASK" "observe")
    ORIENT=$(parse_yaml "$TASKS_FILE" "$TASK" "orient")
    DECIDE=$(parse_yaml "$TASKS_FILE" "$TASK" "decide")
    ACT=$(parse_yaml "$TASKS_FILE" "$TASK" "act")
    
    if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
        echo "Error: Task '$TASK' not found in $TASKS_FILE"
        exit 1
    fi
    
    # Use default iterations if not specified
    if [ $MAX_ITERATIONS -eq 0 ]; then
        DEFAULT_ITER=$(parse_yaml "$TASKS_FILE" "$TASK" "default_iterations")
        [ -n "$DEFAULT_ITER" ] && MAX_ITERATIONS=$DEFAULT_ITER
    fi
fi

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    echo "Usage: ./ooda.sh [task] [--tasks <file>] [--max-iterations N]"
    echo "   OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
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
PLAN_FILE="PLAN-${CURRENT_BRANCH}.md"
FEATURE_FILE="FEATURE-${CURRENT_BRANCH}.md"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
[ -n "$TASK" ] && echo "Task:    $TASK"
echo "Observe: $OBSERVE"
echo "Orient:  $ORIENT"
echo "Decide:  $DECIDE"
echo "Act:     $ACT"
echo "Branch:  $CURRENT_BRANCH"
echo "Plan:    $PLAN_FILE"
echo "Feature: $FEATURE_FILE"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:     $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create interpolated prompt template
create_prompt() {
    cat <<EOF
# OODA Loop Iteration

## Context
- Current branch: \`$CURRENT_BRANCH\`
- Plan file: \`$PLAN_FILE\`
- Feature file: \`$FEATURE_FILE\`

Use these file names when reading or writing plan and feature documents.

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
