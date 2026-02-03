#!/bin/bash
# Usage: ./rooda.sh <procedure> [--config <file>] [--max-iterations N]
#    OR: ./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]
# Examples:
#   ./rooda.sh build
#   ./rooda.sh build --max-iterations 5
#   ./rooda.sh build --config custom-config.yml
#   ./rooda.sh --observe prompts/observe_specs.md \
#             --orient prompts/orient_gap.md \
#             --decide prompts/decide_gap_plan.md \
#             --act prompts/act_plan.md \
#             --max-iterations 1

# Parse YAML (minimal parser for our simple structure)
parse_yaml() {
    local file=$1
    local section=$2
    local field=$3
    local in_section=false
    
    while IFS= read -r line; do
        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// }" ]] && continue
        
        # Check if we're entering the target section
        if [[ "$line" =~ ^${section}: ]]; then
            in_section=true
            continue
        fi
        
        # Check if we've entered a different top-level section
        if [[ "$line" =~ ^[a-z_-]+: ]] && [[ ! "$line" =~ ^${section}: ]]; then
            in_section=false
            continue
        fi
        
        # If we're in the target section, look for the field
        if [ "$in_section" = true ]; then
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
CONFIG_FILE="rooda-config.yml"

# First positional argument is procedure name (optional)
if [[ $# -gt 0 ]] && [[ ! "$1" =~ ^-- ]]; then
    PROCEDURE="$1"
    shift
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --config)
            CONFIG_FILE="$2"
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
            echo "Usage: ./rooda.sh <procedure> [--config <file>] [--max-iterations N]"
            echo "   OR: ./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
            exit 1
            ;;
    esac
done

# If procedure specified, load from config
if [ -n "$PROCEDURE" ]; then
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Error: $CONFIG_FILE not found"
        exit 1
    fi
    
    OBSERVE=$(parse_yaml "$CONFIG_FILE" "procedures" "observe")
    ORIENT=$(parse_yaml "$CONFIG_FILE" "procedures" "orient")
    DECIDE=$(parse_yaml "$CONFIG_FILE" "procedures" "decide")
    ACT=$(parse_yaml "$CONFIG_FILE" "procedures" "act")
    
    # Need to parse within the specific procedure
    OBSERVE=$(parse_yaml "$CONFIG_FILE" "$PROCEDURE" "observe")
    ORIENT=$(parse_yaml "$CONFIG_FILE" "$PROCEDURE" "orient")
    DECIDE=$(parse_yaml "$CONFIG_FILE" "$PROCEDURE" "decide")
    ACT=$(parse_yaml "$CONFIG_FILE" "$PROCEDURE" "act")
    
    if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
        echo "Error: Procedure '$PROCEDURE' not found in $CONFIG_FILE"
        exit 1
    fi
    
    # Use default iterations if not specified
    if [ $MAX_ITERATIONS -eq 0 ]; then
        DEFAULT_ITER=$(parse_yaml "$CONFIG_FILE" "$PROCEDURE" "default_iterations")
        [ -n "$DEFAULT_ITER" ] && MAX_ITERATIONS=$DEFAULT_ITER
    fi
fi

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    echo "Usage: ./rooda.sh <task-id> <procedure> [--config <file>] [--max-iterations N] [--task-file <file>] [--plan-file <file>]"
    echo "   OR: ./rooda.sh <task-id> --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
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

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    echo "Usage: ./rooda.sh <procedure> [--config <file>] [--max-iterations N]"
    echo "   OR: ./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]"
    exit 1
fi

# Validate files exist
for file in "$OBSERVE" "$ORIENT" "$DECIDE" "$ACT"; do
    if [ ! -f "$file" ]; then
        echo "Error: File not found: $file"
        exit 1
    fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
[ -n "$PROCEDURE" ] && echo "Procedure: $PROCEDURE"
echo "Observe:   $OBSERVE"
echo "Orient:    $ORIENT"
echo "Decide:    $DECIDE"
echo "Act:       $ACT"
echo "Branch:    $CURRENT_BRANCH"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:       $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create prompt template
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
