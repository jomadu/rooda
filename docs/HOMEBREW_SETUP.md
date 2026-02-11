# Homebrew Tap Setup Guide

This guide walks through setting up automated Homebrew formula updates for rooda releases.

## Prerequisites

- GitHub account: `jomadu`
- Main repo: `jomadu/rooda`
- Homebrew tap repo: `jomadu/homebrew-rooda` (to be created)

## Step 1: Create Homebrew Tap Repository

```bash
# Create new repo on GitHub
gh repo create jomadu/homebrew-rooda --public --description "Homebrew tap for rooda"

# Clone locally
git clone https://github.com/jomadu/homebrew-rooda.git
cd homebrew-rooda
```

## Step 2: Create Initial Formula

Create `rooda.rb`:

```ruby
class Rooda < Formula
  desc "OODA loop orchestrator for AI coding agents"
  homepage "https://github.com/jomadu/rooda"
  version "0.1.0"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/jomadu/rooda/releases/download/v0.1.0/rooda-darwin-arm64"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    else
      url "https://github.com/jomadu/rooda/releases/download/v0.1.0/rooda-darwin-amd64"
      sha256 "PLACEHOLDER_AMD64_SHA256"
    end
  end

  def install
    bin.install "rooda-darwin-#{Hardware::CPU.arch}" => "rooda"
  end

  test do
    system "#{bin}/rooda", "--version"
  end
end
```

Commit and push:

```bash
git add rooda.rb
git commit -m "Initial rooda formula"
git push origin main
```

## Step 3: Generate GitHub Personal Access Token

1. Go to https://github.com/settings/tokens/new
2. Token name: `rooda-homebrew-automation`
3. Expiration: No expiration (or 1 year with calendar reminder)
4. Scopes:
   - ✅ `repo` (full control)
   - ✅ `workflow` (update GitHub Actions)
5. Click "Generate token"
6. **Copy the token immediately** (you won't see it again)

## Step 4: Add Token to Main Repo Secrets

```bash
# In jomadu/rooda repo
gh secret set HOMEBREW_TAP_TOKEN
# Paste the token when prompted
```

## Step 5: Create Release Workflow

In `jomadu/rooda`, create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build binaries
        run: |
          VERSION=${GITHUB_REF#refs/tags/}
          COMMIT_SHA=$(git rev-parse HEAD)
          BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
          LDFLAGS="-X main.Version=$VERSION -X main.CommitSHA=$COMMIT_SHA -X main.BuildDate=$BUILD_DATE"
          
          # macOS arm64
          GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o rooda-darwin-arm64
          
          # macOS amd64
          GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o rooda-darwin-amd64
          
          # Linux amd64
          GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o rooda-linux-amd64
          
          # Linux arm64
          GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o rooda-linux-arm64
          
          # Windows amd64
          GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o rooda-windows-amd64.exe
      
      - name: Generate checksums
        run: |
          sha256sum rooda-* > checksums.txt
      
      - name: Create release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            rooda-darwin-arm64
            rooda-darwin-amd64
            rooda-linux-amd64
            rooda-linux-arm64
            rooda-windows-amd64.exe
            checksums.txt
            scripts/install.sh
      
      - name: Update Homebrew formula
        env:
          GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
        run: |
          VERSION=${GITHUB_REF#refs/tags/}
          
          # Calculate SHA256 for macOS binaries
          ARM64_SHA=$(sha256sum rooda-darwin-arm64 | cut -d' ' -f1)
          AMD64_SHA=$(sha256sum rooda-darwin-amd64 | cut -d' ' -f1)
          
          # Clone tap repo
          git clone https://x-access-token:${GITHUB_TOKEN}@github.com/jomadu/homebrew-rooda.git
          cd homebrew-rooda
          
          # Update formula
          cat > rooda.rb <<EOF
          class Rooda < Formula
            desc "OODA loop orchestrator for AI coding agents"
            homepage "https://github.com/jomadu/rooda"
            version "$VERSION"
            
            on_macos do
              if Hardware::CPU.arm?
                url "https://github.com/jomadu/rooda/releases/download/$VERSION/rooda-darwin-arm64"
                sha256 "$ARM64_SHA"
              else
                url "https://github.com/jomadu/rooda/releases/download/$VERSION/rooda-darwin-amd64"
                sha256 "$AMD64_SHA"
              end
            end

            def install
              bin.install "rooda-darwin-#{Hardware::CPU.arch}" => "rooda"
            end

            test do
              system "#{bin}/rooda", "--version"
            end
          end
          EOF
          
          # Commit and push
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add rooda.rb
          git commit -m "Update rooda to $VERSION"
          git push
```

## Step 6: Create Install Script

Create `scripts/install.sh` in `jomadu/rooda`:

```bash
#!/bin/sh
set -e

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine binary name
if [ "$OS" = "windows" ]; then
  BINARY="rooda-${OS}-${ARCH}.exe"
  INSTALL_NAME="rooda.exe"
else
  BINARY="rooda-${OS}-${ARCH}"
  INSTALL_NAME="rooda"
fi

# Get latest release
LATEST_RELEASE=$(curl -s https://api.github.com/repos/jomadu/rooda/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
  echo "Failed to fetch latest release"
  exit 1
fi

echo "Installing rooda $LATEST_RELEASE for $OS/$ARCH..."

# Download binary
DOWNLOAD_URL="https://github.com/jomadu/rooda/releases/download/${LATEST_RELEASE}/${BINARY}"
TMP_FILE="/tmp/${BINARY}"

curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"

# Download checksums
curl -fsSL "https://github.com/jomadu/rooda/releases/download/${LATEST_RELEASE}/checksums.txt" -o /tmp/checksums.txt

# Verify checksum
cd /tmp
if command -v sha256sum >/dev/null 2>&1; then
  grep "$BINARY" checksums.txt | sha256sum -c -
elif command -v shasum >/dev/null 2>&1; then
  grep "$BINARY" checksums.txt | shasum -a 256 -c -
else
  echo "Warning: No checksum tool found, skipping verification"
fi

# Install
chmod +x "$TMP_FILE"

if [ "$OS" = "windows" ]; then
  INSTALL_DIR="$HOME/bin"
else
  INSTALL_DIR="/usr/local/bin"
fi

mkdir -p "$INSTALL_DIR"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "$INSTALL_DIR/$INSTALL_NAME"
else
  echo "Installing to $INSTALL_DIR requires sudo..."
  sudo mv "$TMP_FILE" "$INSTALL_DIR/$INSTALL_NAME"
fi

echo "✓ rooda installed to $INSTALL_DIR/$INSTALL_NAME"
echo ""
echo "Run 'rooda version' to verify installation"
```

Make it executable:

```bash
chmod +x scripts/install.sh
git add scripts/install.sh
git commit -m "Add install script with Windows support"
git push
```

## Step 7: Test the Setup

### Test Release Process

```bash
# In jomadu/rooda
git tag v0.1.0
git push origin v0.1.0
```

Watch the GitHub Actions workflow:
- Go to https://github.com/jomadu/rooda/actions
- Verify "Release" workflow runs successfully
- Check that binaries are uploaded to release
- Verify homebrew-rooda repo was updated

### Test Homebrew Installation

```bash
brew tap jomadu/rooda
brew install rooda
rooda version
```

### Test Direct Installation

```bash
# macOS/Linux
curl -fsSL https://github.com/jomadu/rooda/releases/latest/download/install.sh | sh

# Windows (PowerShell)
# Manual download for now - see Windows section below
```

## Windows Installation Notes

Windows users have three options:

1. **Direct download** (recommended):
   ```powershell
   # Download from releases page
   Invoke-WebRequest -Uri "https://github.com/jomadu/rooda/releases/latest/download/rooda-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\rooda.exe"
   ```

2. **WSL** (Linux subsystem):
   ```bash
   curl -fsSL https://github.com/jomadu/rooda/releases/latest/download/install.sh | sh
   ```

3. **Scoop** (future enhancement):
   Create a Scoop manifest for Windows package management.

## Troubleshooting

### Formula update fails
- Verify `HOMEBREW_TAP_TOKEN` secret is set correctly
- Check token has `repo` and `workflow` scopes
- Ensure homebrew-rooda repo exists and is accessible

### Checksum verification fails
- Ensure checksums.txt is generated correctly
- Verify binary wasn't corrupted during upload
- Check that sha256sum/shasum is available on target system

### Binary not executable
- Ensure `chmod +x` is run in install script
- Check file permissions after download

## Maintenance

### Rotating the GitHub Token
1. Generate new token (same scopes)
2. Update `HOMEBREW_TAP_TOKEN` secret
3. Delete old token

### Adding New Platforms
1. Add build step to release.yml
2. Update install.sh platform detection
3. Update checksums generation
4. Test on target platform
