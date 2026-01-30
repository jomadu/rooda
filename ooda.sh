#!/bin/bash
# Usage: ./ooda.sh [procedure] [--procedures <file>] [--task <task-id>] [--max-iterations N]
#    OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> --task <task-id> [--max-iterations N]
# Examples:
#   ./ooda.sh build --task TASK-123
#   ./ooda.sh build --task TASK-123 --max-iterations 5
#   ./ooda.sh build --procedures custom-procedures.yml --task TASK-123
#   ./ooda.sh --observe prompts/observe_specs.md \
#             --orient prompts/orient_gap.md \
#             --decide prompts/decide_gap_plan.md \
#             --act prompts/act_plan.md \
#             --task TASK-123 \
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
PROCEDURE=""
PROCEDURES_FILE="ooda-procedures.yml"
TASK_ID=""

# Check if first arg is a positional procedure name (not a flag)
if [[ $# -gt 0 ]] && [[ ! "$1" =~ ^-- ]]; then
    PROCEDURE="$1"
    shift
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --procedures)
            PROCEDURES_FILE="$2"
            shift 2
            ;;
        --task)
            TASK_ID="$2"
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
            echo "Usage: ./ooda.sh [procedure] [--procedures <file>] [--task <task-id>] [--max-iterations N]"
            echo "   OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> --task <task-id> [--max-iterations N]"
            exit 1
            ;;
    esac
done

# If procedure specified, load from config
if [ -n "$PROCEDURE" ]; then
    if [ ! -f "$PROCEDURES_FILE" ]; then
        echo "Error: $PROCEDURES_FILE not found"
        exit 1
    fi
    
    OBSERVE=$(parse_yaml "$PROCEDURES_FILE" "$PROCEDURE" "observe")
    ORIENT=$(parse_yaml "$PROCEDURES_FILE" "$PROCEDURE" "orient")
    DECIDE=$(parse_yaml "$PROCEDURES_FILE" "$PROCEDURE" "decide")
    ACT=$(parse_yaml "$PROCEDURES_FILE" "$PROCEDURE" "act")
    
    if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
        echo "Error: Procedure '$PROCEDURE' not found in $PROCEDURES_FILE"
        exit 1
    fi
    
    # Use default iterations if not specified
    if [ $MAX_ITERATIONS -eq 0 ]; then
        DEFAULT_ITER=$(parse_yaml "$PROCEDURES_FILE" "$PROCEDURE" "default_iterations")
        [ -n "$DEFAULT_ITER" ] && MAX_ITERATIONS=$DEFAULT_ITER
    fi
fi

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    echo "Usage: ./ooda.sh [procedure] [--procedures <file>] [--task <task-id>] [--max-iterations N]"
    echo "   OR: ./ooda.sh --observe <file> --orient <file> --decide <file> --act <file> --task <task-id> [--max-iterations N]"
    exit 1
fi

# Task ID is required for non-bootstrap procedures
if [ "$PROCEDURE" != "bootstrap" ] && [ -z "$TASK_ID" ]; then
    echo "Error: --task <task-id> is required"
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

# Set up task-specific paths
if [ -n "$TASK_ID" ]; then
    TASK_DIR="tasks/${TASK_ID}"
    PLAN_FILE="${TASK_DIR}/PLAN.md"
    STORY_FILE="${TASK_DIR}/STORY.md"
    BUG_FILE="${TASK_DIR}/BUG.md"
    
    # Create task directory if it doesn't exist
    mkdir -p "$TASK_DIR"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
[ -n "$PROCEDURE" ] && echo "Procedure: $PROCEDURE"
echo "Observe:   $OBSERVE"
echo "Orient:    $ORIENT"
echo "Decide:    $DECIDE"
echo "Act:       $ACT"
echo "Branch:    $CURRENT_BRANCH"
[ -n "$TASK_ID" ] && echo "Task:      $TASK_ID"
[ -n "$TASK_ID" ] && echo "Task Dir:  $TASK_DIR"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:       $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create interpolated prompt template
create_prompt() {
    cat <<EOF
# OODA Loop Iteration

## Context
- Current branch: \`$CURRENT_BRANCH\`
EOF
    
    if [ -n "$TASK_ID" ]; then
        cat <<EOF
- Task ID: \`$TASK_ID\`
- Task directory: \`$TASK_DIR\`
- Plan file: \`$PLAN_FILE\`
- Story file: \`$STORY_FILE\` (if this is a feature/story)
- Bug file: \`$BUG_FILE\` (if this is a bug fix)

Use these file paths when reading or writing task-related documents.
EOF
    fi
    
    cat <<EOF

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
