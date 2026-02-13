# Task: Remove Homebrew Distribution Channel

## Problem

Homebrew is currently listed as a distribution channel in the installation documentation, but it represents a breaking change without proper migration support or documentation.

## Required Changes

1. **Remove Homebrew references from documentation**
   - Remove Homebrew installation instructions from `docs/installation.md`
   - Remove Homebrew references from `README.md`
   - Update quick start guide to use alternative installation method

2. **Update installation script**
   - Verify `scripts/install.sh` does not depend on Homebrew
   - Ensure direct download method is the primary installation path

3. **Fix AGENTS.md documentation**
   - Remove references to non-existent procedures `draft-plan-story-to-spec` and `draft-plan-bug-to-spec`
   - Update "Story/Bug Input" section with correct procedure names
   - Verify all procedure references in AGENTS.md match actual available procedures

4. **Breaking change considerations**
   - No migration path exists for users who installed via Homebrew
   - No documentation exists explaining the removal
   - Users will need to uninstall Homebrew version and reinstall via alternative method

## Impact

- **Breaking change**: Existing Homebrew users will not receive updates
- **No migration support**: Users must manually uninstall and reinstall
- **Documentation gap**: No communication plan for affected users

## Priority

High - This is a breaking change that affects existing users and creates confusion in the installation documentation.
