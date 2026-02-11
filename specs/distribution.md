# Distribution

## Job to be Done

Enable users to install rooda as a single binary with no external dependencies, supporting macOS, Linux, Windows, and CI/CD environments.

## Activities

1. **Build single binary** — Compile Go source to standalone executable with embedded prompts
2. **Cross-compile for platforms** — Generate binaries for macOS (arm64, amd64), Linux (amd64, arm64), Windows (amd64)
3. **Embed default prompts** — Package `prompts/*.md` files using `go:embed` so binary is self-contained
4. **Install via package manager** — Support Homebrew for macOS users
5. **Install via direct download** — Provide curl-based installation for Linux/CI environments
6. **Install via Go toolchain** — Support `go install` for Go developers
7. **Version the binary** — Embed version string at build time using `-ldflags`

## Acceptance Criteria

- [ ] `go build` produces single binary with no runtime dependencies (no external yq, no separate prompt files)
- [ ] Binary runs on macOS arm64, macOS amd64, Linux amd64, Linux arm64, Windows amd64
- [ ] SHA256 checksums generated for all binaries in checksums.txt
- [ ] Install script verifies checksums before installation
- [ ] Install script hosted in GitHub Releases (not main branch)
- [ ] `rooda version` reports correct version string embedded at build time
- [ ] Default prompts are accessible when no custom prompts provided (embedded via `go:embed`)
- [ ] Homebrew formula installs binary to standard location and adds to PATH
- [ ] `curl | sh` installation script downloads correct binary for detected platform
- [ ] `go install github.com/jomadu/rooda@latest` installs from source
- [ ] Binary size is reasonable (< 20MB uncompressed)
- [ ] Installation instructions documented in README.md

## Data Structures

### Build Metadata
```go
// Embedded at compile time via -ldflags
var (
    Version   string // e.g., "v2.0.0"
    CommitSHA string // e.g., "a1b2c3d"
    BuildDate string // e.g., "2026-02-08T20:00:00Z"
)
```

### Embedded Prompts
```go
//go:embed prompts/*.md
var defaultPrompts embed.FS
```

### Platform Targets
```
GOOS=darwin GOARCH=arm64  → rooda-darwin-arm64
GOOS=darwin GOARCH=amd64  → rooda-darwin-amd64
GOOS=linux  GOARCH=amd64  → rooda-linux-amd64
GOOS=linux  GOARCH=arm64  → rooda-linux-arm64
GOOS=windows GOARCH=amd64 → rooda-windows-amd64.exe
```

## Algorithm

### Build Process
```
1. Embed version metadata:
   go build -ldflags "-X main.Version=$(git describe --tags) \
                      -X main.CommitSHA=$(git rev-parse HEAD) \
                      -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

2. Embed default prompts:
   // In main package
   //go:embed prompts/*.md
   var defaultPrompts embed.FS

3. Build binary:
   # Using Makefile (recommended)
   make build
   
   # Or directly
   go build -o bin/rooda ./cmd/rooda

4. Cross-compile for each platform:
   for platform in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64; do
     GOOS=${platform%/*} GOARCH=${platform#*/} go build -o rooda-$platform
   done

5. Generate checksums:
   sha256sum rooda-* > checksums.txt

6. Package for distribution:
   - Homebrew: Create formula with download URL and SHA256
   - Direct download: Host binaries with install.sh script (includes checksum verification)
   - Go install: Tag release, push to GitHub
   - Include install.sh and checksums.txt in GitHub Release assets
```

### Installation Methods

**Homebrew (macOS):**
```bash
brew tap jomadu/rooda
brew install rooda
```

**Direct download (Linux/macOS/Windows):**
```bash
curl -fsSL https://github.com/jomadu/rooda/releases/latest/download/install.sh | sh
```

**Go toolchain:**
```bash
go install github.com/jomadu/rooda@latest
```

### Runtime Prompt Resolution
```
1. Check for custom prompts in workspace (./prompts/)
2. Check for custom prompts in global config (~/.config/rooda/prompts/)
3. Fall back to embedded defaults (extracted from embed.FS)
```

## Edge Cases

### Missing Embedded Prompts
- **Scenario:** Build succeeds but `go:embed` directive is malformed
- **Detection:** Unit test verifies all expected prompt files are accessible via embed.FS
- **Handling:** Build fails if any required prompt file is missing from embedded FS

### Platform Detection Failure
- **Scenario:** Install script runs on unsupported platform (e.g., BSD, Solaris)
- **Detection:** `uname -s` and `uname -m` don't match known patterns
- **Handling:** Script prints error with supported platforms and exits

### Windows Support
- **Scenario:** Windows user wants to install rooda
- **Options:**
  1. Direct download: `rooda-windows-amd64.exe` from GitHub Releases
  2. WSL: Use Linux install script
  3. Git Bash/MSYS: Install script detects Windows
- **Path conventions:** Windows uses `%USERPROFILE%\bin` instead of `/usr/local/bin`
- **Config directory:** Windows uses `%APPDATA%\rooda` instead of `~/.config/rooda` (see configuration.md)

### Version Mismatch
- **Scenario:** Binary reports version but doesn't match git tag
- **Detection:** CI verifies `rooda version` output matches `$GITHUB_REF_NAME`
- **Handling:** Release fails if version mismatch detected

### Homebrew Formula Outdated
- **Scenario:** New release published but Homebrew formula not updated
- **Detection:** GitHub Actions workflow updates formula automatically
- **Handling:** Workflow clones jomadu/homebrew-rooda, updates rooda.rb with new version and SHA256s, commits and pushes
- **Setup:** Requires HOMEBREW_TAP_TOKEN secret with repo + workflow scopes (see docs/HOMEBREW_SETUP.md)

### Large Binary Size
- **Scenario:** Binary exceeds 20MB due to embedded assets
- **Detection:** CI checks binary size after build
- **Handling:** Warning if > 20MB, failure if > 50MB

## Dependencies

### Build-time
- Go 1.21+ (for `go:embed` and modern stdlib)
- git (for version metadata)
- make (optional but recommended for unified build interface)
- Cross-compilation toolchain (built into Go)

### Runtime
- None (single binary, no external dependencies)

### Distribution
- GitHub Releases (for hosting binaries)
- Homebrew tap repository (for macOS package manager)
- install.sh script (for curl-based installation)

## Implementation Mapping

### Source Files
- `main.go` — Entry point, version flag, embedded prompts
- `internal/prompts/loader.go` — Prompt resolution (custom vs embedded)
- `scripts/build.sh` — Cross-compilation script
- `scripts/install.sh` — Platform detection and download script
- `jomadu/homebrew-rooda/rooda.rb` — Homebrew formula (separate repo, auto-updated by CI)
- `.github/workflows/release.yml` — CI pipeline for building and publishing

### Related Specs
- [configuration.md](configuration.md) — Prompt path resolution (workspace > global > embedded)
- [prompt-composition.md](prompt-composition.md) — How prompts are loaded and assembled
- [cli-interface.md](cli-interface.md) — `--version` flag implementation

## Examples

### Example 1: Build with Version Metadata
```bash
$ git describe --tags
v2.0.0

$ go build -ldflags "-X main.Version=v2.0.0 -X main.CommitSHA=$(git rev-parse HEAD)"

$ ./rooda version
rooda v2.0.0 (commit: a1b2c3d4e5f6, built: 2026-02-08T20:00:00Z)
```

**Verification:** Version string matches git tag and includes commit SHA.

### Example 2: Cross-Compile for Linux
```bash
$ GOOS=linux GOARCH=amd64 go build -o rooda-linux-amd64

$ file rooda-linux-amd64
rooda-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=..., not stripped

$ ls -lh rooda-linux-amd64
-rwxr-xr-x  1 user  staff   12M Feb  8 20:00 rooda-linux-amd64
```

**Verification:** Binary is ELF format, statically linked, reasonable size.

### Example 3: Install via Homebrew
```bash
$ brew tap jomadu/rooda
$ brew install rooda

$ which rooda
/opt/homebrew/bin/rooda

$ rooda version
rooda v2.0.0 (commit: a1b2c3d4e5f6, built: 2026-02-08T20:00:00Z)
```

**Verification:** Binary installed to standard Homebrew location, version correct.

### Example 4: Install via curl
```bash
$ curl -fsSL https://github.com/jomadu/rooda/releases/latest/download/install.sh | sh
Detecting platform... darwin/arm64
Downloading rooda v2.0.0...
Installing to /usr/local/bin/rooda...
Installation complete!

$ rooda version
rooda v2.0.0 (commit: a1b2c3d4e5f6, built: 2026-02-08T20:00:00Z)
```

**Verification:** Script detects platform, downloads correct binary, installs to PATH.

### Example 5: Embedded Prompts Fallback
```bash
$ rm -rf ./prompts  # No custom prompts
$ rm -rf ~/.config/rooda/prompts  # No global prompts

$ rooda run build --dry-run
# Observe: Plan, Specs, Implementation
...
# (Prompt assembled from embedded defaults)
```

**Verification:** Binary runs without external prompt files, uses embedded defaults.

### Example 6: CI/CD Installation
```yaml
# .github/workflows/test.yml
steps:
  - name: Install rooda
    run: |
      curl -fsSL https://github.com/jomadu/rooda/releases/latest/download/install.sh | sh
      rooda version
  
  - name: Run rooda procedure
    run: rooda run build --max-iterations 5
```

**Verification:** Installation succeeds in CI environment, binary executes.

## Notes

### Why Single Binary?

The v1 bash implementation required external dependencies (`yq`, AI CLI tools) and separate prompt files. This created friction:
- Users had to install `yq` separately
- Prompt files could get out of sync with script version
- Installation was multi-step (copy script, copy prompts, install deps)

A single Go binary with embedded prompts eliminates all external dependencies except the AI CLI tool (which is configurable and user-provided). This aligns with the "minimal setup friction" outcome.

### Why `go:embed` for Prompts?

Embedding default prompts ensures the binary is self-contained and version-locked. Users can still override with custom prompts (workspace or global config), but the binary always has a working set of defaults. This prevents "missing prompt file" errors and simplifies distribution.

### Why Cross-Compilation?

Go's built-in cross-compilation makes it trivial to support multiple platforms without separate build environments. A single `go build` command with `GOOS` and `GOARCH` produces binaries for macOS, Linux, and different architectures. This is essential for CI/CD environments where users may run rooda on various platforms.

### Why Homebrew + curl + go install?

Different users have different preferences:
- **Homebrew:** macOS users expect `brew install` for CLI tools
- **curl | sh:** Common pattern for Linux/CI environments (Docker, GitHub Actions)
- **go install:** Go developers prefer installing from source

Supporting all three maximizes adoption with minimal maintenance (Homebrew formula is auto-generated, install.sh is a simple script, `go install` is built into Go).

### Binary Size Considerations

Embedding prompts adds ~50KB (25 files × ~2KB each). Go binaries are statically linked, so base size is ~10-15MB. Total expected size: 12-18MB, well under the 20MB threshold. If size becomes an issue, consider:
- Strip debug symbols: `go build -ldflags "-s -w"`
- Compress with UPX (trade-off: slower startup)
- Lazy-load prompts from embedded FS (minimal impact)

### Version Embedding

Using `-ldflags` to inject version metadata at build time ensures `rooda version` always reports accurate information. This is critical for debugging ("what version are you running?") and for CI pipelines that verify version consistency.

### Installation Script Security

The `curl | sh` pattern is convenient but risky if the script is compromised. Mitigations:
- Host script in version-controlled repo (GitHub)
- Use HTTPS (fsSL flags: fail silently, show errors, follow redirects, use SSL)
- Verify binary SHA256 in script before installing
- Document alternative: download binary directly and verify manually

### Windows Support

Windows amd64 is supported via:
- `GOOS=windows GOARCH=amd64` build target produces `rooda-windows-amd64.exe`
- Install script detects Windows (Git Bash/MSYS) and handles `.exe` extension
- Config directory: `%APPDATA%\rooda` (e.g., `C:\Users\username\AppData\Roaming\rooda`)
- Binary install location: `%USERPROFILE%\bin` (user must add to PATH manually)

Windows users can also:
- Download `.exe` directly from GitHub Releases
- Use WSL (Windows Subsystem for Linux) with Linux binary
- Use Scoop package manager (future enhancement)

### Binary Naming Convention

**Release artifacts** (platform-specific):
- `rooda-darwin-arm64`
- `rooda-darwin-amd64`
- `rooda-linux-amd64`
- `rooda-linux-arm64`
- `rooda-windows-amd64.exe`

**Installed binary** (generic, in PATH):
- Unix/macOS: `rooda`
- Windows: `rooda.exe`

Install script renames platform-specific artifact to generic name during installation.

### Checksum Verification

Generate checksums during build:
```bash
sha256sum rooda-* > checksums.txt
```

Install script verifies before execution:
```bash
# Uses sha256sum (Linux) or shasum (macOS)
grep "$BINARY" checksums.txt | sha256sum -c -
```

Prevents MITM attacks and ensures binary integrity. Checksums.txt included in GitHub Release assets.
