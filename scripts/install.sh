#!/bin/bash
# Burnless one-liner installer
# Usage: curl -fsSL https://burnless.dev/install.sh | sh

set -e

REPO="burnless/burnless"
BINARY="burnless"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

LATEST=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
URL="https://github.com/$REPO/releases/download/$LATEST/${BINARY}-${OS}-${ARCH}.tar.gz"

echo "Installing burnless $LATEST for $OS/$ARCH..."
curl -fsSL "$URL" | tar -xz -C /tmp
chmod +x /tmp/$BINARY
sudo mv /tmp/$BINARY $INSTALL_DIR/$BINARY
echo "✓ Installed: $(burnless version)"
