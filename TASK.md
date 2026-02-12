# Task: Configure semantic-release for GitHub Releases

## Objective

Set up semantic-release to automate GitHub releases and version tagging using conventional commits.

## Requirements

### Branch Configuration
- **Main distribution branch**: `main`
- **Pre-release branch**: `beta`

### Commit Convention
- Use conventional commit messages
- Use `conventionalcommits` preset for semantic-release

### Dependencies
Install the following npm packages:
- `semantic-release`
- `@semantic-release/commit-analyzer`
- `@semantic-release/release-notes-generator`
- `@semantic-release/github`
- `conventional-changelog-conventionalcommits`
- `husky`
- `@commitlint/cli`
- `@commitlint/config-conventional`

### Configuration File
Create `.releaserc.json` with:
- Branch configuration for `main` and `beta` (prerelease)
- Commit analyzer using `conventionalcommits` preset
- Release notes generator
- GitHub plugin for creating releases

### Commit Message Linting
Set up husky pre-commit hooks:
- Initialize husky with `npx husky init`
- Configure commitlint with `.commitlintrc.json` using `@commitlint/config-conventional`
- Add commit-msg hook to validate conventional commit format

### CI/CD Integration
Configure GitHub Actions workflow to run semantic-release on push to `main` and `beta` branches.

Add commit message validation to CI pipeline:
- Validate conventional commit format in pull requests
- Use commitlint in GitHub Actions workflow

### Cleanup
Remove all Homebrew distribution:
- Delete `.github/workflows/homebrew.yml` (if exists)
- Remove Homebrew tap references from documentation
- Update installation docs to remove Homebrew instructions

## Reference
See `~/Dev/ai-resource-manager` for working implementation example.
