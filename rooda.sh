#!/bin/bash
# Wrapper script for framework developers to invoke src/rooda.sh from project root
# NOTE: This wrapper is NOT part of the consumer installation.
# Consumers copy files from src/ to their project root per installation instructions.
exec "$(dirname "$0")/src/rooda.sh" "$@"
