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

# Check for yq dependency
if ! command -v yq &> /dev/null; then
    echo "Error: yq is required for YAML parsing"
    echo "Install with: brew install yq"
    exit 1
fi

# Parse arguments
OBSERVE=""
ORIENT=""
DECIDE=""
ACT=""
MAX_ITERATIONS=0
PROCEDURE=""
# Resolve config file relative to script location
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/rooda-config.yml"

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
    
    OBSERVE=$(yq eval ".procedures.$PROCEDURE.observe" "$CONFIG_FILE")
    ORIENT=$(yq eval ".procedures.$PROCEDURE.orient" "$CONFIG_FILE")
    DECIDE=$(yq eval ".procedures.$PROCEDURE.decide" "$CONFIG_FILE")
    ACT=$(yq eval ".procedures.$PROCEDURE.act" "$CONFIG_FILE")
    
    if [ "$OBSERVE" = "null" ] || [ "$ORIENT" = "null" ] || [ "$DECIDE" = "null" ] || [ "$ACT" = "null" ]; then
        echo "Error: Procedure '$PROCEDURE' not found in $CONFIG_FILE"
        exit 1
    fi
    
    # Use default iterations if not specified
    if [ "$MAX_ITERATIONS" -eq 0 ]; then
        DEFAULT_ITER=$(yq eval ".procedures.$PROCEDURE.default_iterations" "$CONFIG_FILE")
        [ "$DEFAULT_ITER" != "null" ] && MAX_ITERATIONS=$DEFAULT_ITER
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
[ "$MAX_ITERATIONS" -gt 0 ] && echo "Max:       $MAX_ITERATIONS iterations"
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
    if [ "$MAX_ITERATIONS" -gt 0 ] && [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
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
