#!/bin/sh
# MCP Hub installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Ricardo-M-L/mcphub/master/install.sh | sh

set -e

REPO="Ricardo-M-L/mcphub"
BINARY="mcphub"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin|linux) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

ASSET="${BINARY}-${OS}-${ARCH}"
echo "Downloading mcphub for ${OS}/${ARCH}..."

# Get latest release URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

# Download
TMPDIR=$(mktemp -d)
curl -fsSL "$DOWNLOAD_URL" -o "$TMPDIR/$BINARY"
chmod +x "$TMPDIR/$BINARY"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
else
    echo "Need sudo to install to $INSTALL_DIR"
    sudo mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
fi

rm -rf "$TMPDIR"

echo ""
echo "mcphub installed successfully!"
echo ""
echo "Get started:"
echo "  mcphub search filesystem"
echo "  mcphub install <server-name>"
echo ""
