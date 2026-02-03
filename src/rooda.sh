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

show_help() {
    cat <<EOF
Usage: ./rooda.sh <procedure> [--config <file>] [--max-iterations N]
   OR: ./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]

Options:
  <procedure>           Named procedure from config (bootstrap, build, etc.)
  --config <file>       Path to config file (default: rooda-config.yml)
  --observe <file>      Path to observe phase prompt
  --orient <file>       Path to orient phase prompt
  --decide <file>       Path to decide phase prompt
  --act <file>          Path to act phase prompt
  --max-iterations N    Maximum iterations (default: see below)
  --help, -h            Show this help message

Max Iterations Default Behavior (three-tier system):
  1. Command-line --max-iterations takes precedence
  2. Config default_iterations used if CLI not specified
  3. Defaults to 0 (unlimited) if neither specified

Examples:
  ./rooda.sh bootstrap
  ./rooda.sh build --max-iterations 5
  ./rooda.sh --observe prompts/observe_specs.md \\
            --orient prompts/orient_gap.md \\
            --decide prompts/decide_gap_plan.md \\
            --act prompts/act_plan.md \\
            --max-iterations 1
EOF
}

# Detect OS for platform-specific instructions
OS="$(uname -s)"
case "$OS" in
    Darwin*) PLATFORM="macos" ;;
    Linux*)  PLATFORM="linux" ;;
    *)       PLATFORM="unknown" ;;
esac

# Check for required dependencies
if ! command -v yq &> /dev/null; then
    echo "Error: yq is required for YAML parsing"
    if [ "$PLATFORM" = "macos" ]; then
        echo "Install with: brew install yq"
    elif [ "$PLATFORM" = "linux" ]; then
        echo "Install with: wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/local/bin/yq && chmod +x /usr/local/bin/yq"
    else
        echo "See: https://github.com/mikefarah/yq#install"
    fi
    exit 1
fi

if ! command -v kiro-cli &> /dev/null; then
    echo "Error: kiro-cli is required for AI CLI integration"
    echo "Install from: https://docs.aws.amazon.com/kiro/"
    exit 1
fi

if ! command -v bd &> /dev/null; then
    echo "Error: bd (beads) is required for work tracking"
    if [ "$PLATFORM" = "macos" ] || [ "$PLATFORM" = "linux" ]; then
        echo "Install with: cargo install beads-cli"
    fi
    echo "Or download from: https://github.com/jomadu/beads/releases"
    exit 1
fi

# Check versions
YQ_VERSION=$(yq --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
YQ_MAJOR=$(echo "$YQ_VERSION" | cut -d. -f1)
if [ "$YQ_MAJOR" -lt 4 ]; then
    echo "Error: yq version 4.0.0 or higher required (found $YQ_VERSION)"
    echo "Upgrade with: brew upgrade yq"
    exit 1
fi

KIRO_VERSION=$(kiro-cli --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
KIRO_MAJOR=$(echo "$KIRO_VERSION" | cut -d. -f1)
if [ "$KIRO_MAJOR" -lt 1 ]; then
    echo "Error: kiro-cli version 1.0.0 or higher required (found $KIRO_VERSION)"
    echo "Upgrade from: https://docs.aws.amazon.com/kiro/"
    exit 1
fi

BD_VERSION=$(bd --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
BD_MAJOR=$(echo "$BD_VERSION" | cut -d. -f1)
BD_MINOR=$(echo "$BD_VERSION" | cut -d. -f2)
if [ "$BD_MAJOR" -eq 0 ] && [ "$BD_MINOR" -lt 1 ]; then
    echo "Error: bd version 0.1.0 or higher required (found $BD_VERSION)"
    echo "Upgrade with: cargo install beads-cli"
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

# Check for help flag first
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    show_help
    exit 0
fi

# First positional argument is procedure name (optional)
if [[ $# -gt 0 ]] && [[ ! "$1" =~ ^-- ]]; then
    PROCEDURE="$1"
    shift
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            show_help
            exit 0
            ;;
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
            show_help
            exit 1
            ;;
    esac
done

# If procedure specified, load from config (explicit flags override config)
if [ -n "$PROCEDURE" ]; then
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Error: Configuration file not found"
        echo "  Path: $CONFIG_FILE"
        echo "  Specify a different config with --config <file>"
        exit 1
    fi
    
    # Only load from config if not already set via explicit flags
    [ -z "$OBSERVE" ] && OBSERVE=$(yq eval ".procedures.$PROCEDURE.observe" "$CONFIG_FILE")
    [ -z "$ORIENT" ] && ORIENT=$(yq eval ".procedures.$PROCEDURE.orient" "$CONFIG_FILE")
    [ -z "$DECIDE" ] && DECIDE=$(yq eval ".procedures.$PROCEDURE.decide" "$CONFIG_FILE")
    [ -z "$ACT" ] && ACT=$(yq eval ".procedures.$PROCEDURE.act" "$CONFIG_FILE")
    
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
    show_help
    exit 1
fi

# Validate files exist
for file in "$OBSERVE" "$ORIENT" "$DECIDE" "$ACT"; do
    if [ ! -f "$file" ]; then
        echo "Error: OODA phase file not found"
        echo "  Path: $file"
        echo "  Check that all four phase files exist (observe, orient, decide, act)"
        exit 1
    fi
done

ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

# Validate files exist
for file in "$OBSERVE" "$ORIENT" "$DECIDE" "$ACT"; do
    if [ ! -f "$file" ]; then
        echo "Error: OODA phase file not found"
        echo "  Path: $file"
        echo "  Check that all four phase files exist (observe, orient, decide, act)"
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
