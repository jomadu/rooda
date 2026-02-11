# Installation

## Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/rooda/main/scripts/install.sh | bash
```

This installs the latest release to `/usr/local/bin/rooda` (or `~/.local/bin/rooda` if `/usr/local/bin` is not writable).

## Homebrew (macOS/Linux)

```bash
brew tap jomadu/rooda
brew install rooda
```

## Direct Download

Download the latest release for your platform from [GitHub Releases](https://github.com/jomadu/rooda/releases):

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/jomadu/rooda/releases/latest/download/rooda-darwin-arm64 -o rooda
chmod +x rooda
sudo mv rooda /usr/local/bin/
```

**macOS (Intel):**
```bash
curl -L https://github.com/jomadu/rooda/releases/latest/download/rooda-darwin-amd64 -o rooda
chmod +x rooda
sudo mv rooda /usr/local/bin/
```

**Linux (x86_64):**
```bash
curl -L https://github.com/jomadu/rooda/releases/latest/download/rooda-linux-amd64 -o rooda
chmod +x rooda
sudo mv rooda /usr/local/bin/
```

**Linux (ARM64):**
```bash
curl -L https://github.com/jomadu/rooda/releases/latest/download/rooda-linux-arm64 -o rooda
chmod +x rooda
sudo mv rooda /usr/local/bin/
```

**Windows:**
```powershell
# Download from https://github.com/jomadu/rooda/releases/latest/download/rooda-windows-amd64.exe
# Add to PATH
```

## Build from Source

Requires Go >= 1.24.5:

```bash
git clone https://github.com/jomadu/rooda.git
cd rooda
go build -o bin/rooda ./cmd/rooda
sudo mv bin/rooda /usr/local/bin/
```

Or use the build script:

```bash
./scripts/build.sh
# Binaries created in bin/ for all platforms
```

## Verify Installation

```bash
rooda version
rooda list
```

## Next Steps

1. **Bootstrap your repository**: `rooda bootstrap --ai-cmd-alias kiro-cli`
2. **Configure AI command**: Set `--ai-cmd-alias` or configure in `~/.config/rooda/rooda-config.yml`
3. **Run a procedure**: `rooda audit-spec --ai-cmd-alias kiro-cli`

See [Configuration](configuration.md) for details on setting up AI commands and procedures.
