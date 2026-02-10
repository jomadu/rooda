#!/usr/bin/env bash
set -euo pipefail

REPO="jomadu/rooda"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="rooda"

detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  
  case "$os" in
    linux) os="linux" ;;
    darwin) os="darwin" ;;
    mingw*|msys*|cygwin*) os="windows" ;;
    *) echo "Unsupported OS: $os" >&2; exit 1 ;;
  esac
  
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
  esac
  
  echo "${os}-${arch}"
}

get_latest_release() {
  local url="https://api.github.com/repos/${REPO}/releases/latest"
  if command -v curl >/dev/null 2>&1; then
    curl -sL "$url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "$url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  else
    echo "Error: curl or wget required" >&2
    exit 1
  fi
}

download_file() {
  local url="$1"
  local output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$output"
  else
    echo "Error: curl or wget required" >&2
    exit 1
  fi
}

verify_checksum() {
  local file="$1"
  local expected="$2"
  local actual
  
  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "$file" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "$file" | awk '{print $1}')"
  else
    echo "Warning: sha256sum/shasum not found, skipping checksum verification" >&2
    return 0
  fi
  
  if [ "$actual" != "$expected" ]; then
    echo "Error: Checksum mismatch" >&2
    echo "  Expected: $expected" >&2
    echo "  Actual:   $actual" >&2
    return 1
  fi
  return 0
}

main() {
  local platform version base_url binary_url checksums_url tmpdir binary_file checksums_file
  local binary_name expected_checksum
  
  echo "Detecting platform..."
  platform="$(detect_platform)"
  echo "Platform: $platform"
  
  echo "Fetching latest release..."
  version="$(get_latest_release)"
  if [ -z "$version" ]; then
    echo "Error: Could not determine latest release" >&2
    exit 1
  fi
  echo "Version: $version"
  
  base_url="https://github.com/${REPO}/releases/download/${version}"
  binary_name="${BINARY_NAME}-${platform}"
  [ "$platform" = "windows-amd64" ] && binary_name="${binary_name}.exe"
  
  binary_url="${base_url}/${binary_name}"
  checksums_url="${base_url}/checksums.txt"
  
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT
  
  binary_file="${tmpdir}/${binary_name}"
  checksums_file="${tmpdir}/checksums.txt"
  
  echo "Downloading ${binary_name}..."
  download_file "$binary_url" "$binary_file"
  
  echo "Downloading checksums..."
  download_file "$checksums_url" "$checksums_file"
  
  echo "Verifying checksum..."
  expected_checksum="$(grep "$binary_name" "$checksums_file" | awk '{print $1}')"
  if [ -z "$expected_checksum" ]; then
    echo "Error: Checksum not found for $binary_name" >&2
    exit 1
  fi
  
  if ! verify_checksum "$binary_file" "$expected_checksum"; then
    exit 1
  fi
  echo "Checksum verified"
  
  echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
  if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR" || {
      echo "Error: Cannot create $INSTALL_DIR (try with sudo)" >&2
      exit 1
    }
  fi
  
  if [ ! -w "$INSTALL_DIR" ]; then
    echo "Error: No write permission to $INSTALL_DIR (try with sudo)" >&2
    exit 1
  fi
  
  cp "$binary_file" "${INSTALL_DIR}/${BINARY_NAME}"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
  
  echo ""
  echo "✓ rooda ${version} installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
  echo ""
  echo "Run 'rooda --version' to verify installation"
  
  if ! command -v rooda >/dev/null 2>&1; then
    echo ""
    echo "Note: ${INSTALL_DIR} is not in your PATH"
    echo "Add it to your shell profile:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi
}

main "$@"
