#!/usr/bin/env bash
# Goslide One-Liner Installer for Linux & macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/joincdream/goslide/main/install.sh | bash

set -e

REPO="joincdream/goslide"
INSTALL_DIR="${GOSLIDE_INSTALL_DIR:-$HOME/.local/bin}"

# 1. Detect Operating System
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux*)  OS="linux" ;;
  darwin*) OS="darwin" ;;
  msys*|cygwin*|mingw*)
    echo "❌ Windows environment detected. Please download and extract goslide_windows_amd64.zip manually from https://github.com/${REPO}/releases"
    exit 1
    ;;
  *)
    echo "❌ Unsupported operating system: $OS"
    exit 1
    ;;
esac

# 2. Detect CPU Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "❌ Unsupported CPU architecture: $ARCH"
    exit 1
    ;;
esac

# 3. Resolve Target Version
if [ -z "$VERSION" ]; then
  echo "🔍 Resolving latest release for ${OS}/${ARCH}..."
  LATEST_JSON="$(curl -fsSL -H "User-Agent: Goslide-Installer" "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null || true)"
  VERSION="$(echo "$LATEST_JSON" | grep -m1 '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/' || true)"
fi

if [ -z "$VERSION" ]; then
  echo "⚠️  Could not resolve latest release via GitHub API, falling back to default v0.1.1-test"
  VERSION="v0.1.1-test"
fi

echo "📦 Installing Goslide ${VERSION} (${OS}/${ARCH})..."

# 4. Download and Extract Archive
ARTIFACT="goslide_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARTIFACT}"

TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'goslide-install')"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

echo "⬇️  Downloading ${DOWNLOAD_URL}..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARTIFACT"; then
  echo "❌ Failed to download $DOWNLOAD_URL"
  exit 1
fi

echo "📂 Extracting archive..."
tar -xzf "$TMP_DIR/$ARTIFACT" -C "$TMP_DIR"

# Locate binary whether it is at archive root or inside a subfolder
BINARY_PATH="$(find "$TMP_DIR" -type f -name "goslide" 2>/dev/null | head -n 1)"

if [ -z "$BINARY_PATH" ] || [ ! -f "$BINARY_PATH" ]; then
  echo "❌ Binary 'goslide' not found inside archive"
  exit 1
fi

# 5. Install to Target Directory
mkdir -p "$INSTALL_DIR"
cp "$BINARY_PATH" "$INSTALL_DIR/goslide"
chmod 755 "$INSTALL_DIR/goslide"

# 6. macOS Gatekeeper Quarantine Removal
if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine "$INSTALL_DIR/goslide" 2>/dev/null || true
fi

echo ""
echo "🎉 Successfully installed Goslide to: ${INSTALL_DIR}/goslide"
echo ""

# 7. Check if INSTALL_DIR is in PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*)
    echo "✨ ${INSTALL_DIR} is already in your PATH."
    ;;
  *)
    echo "⚠️  ${INSTALL_DIR} is NOT in your current PATH!"
    echo "👉 Please add the following line to your shell configuration profile (~/.bashrc, ~/.zshrc):"
    echo ""
    echo "   export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo ""
    ;;
esac

echo "🚀 Quick Start:"
echo "   goslide demo"
echo "   goslide serve demo/demo.md"
echo "   goslide update --check"
echo ""
