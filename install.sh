#!/bin/sh
# MCP Hub installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Ricardo-M-L/mcphub/master/install.sh | sh

set -e

REPO="Ricardo-M-L/mcphub"
BINARY="mcphub"
MCP_BINARY="mcphub-mcp"

# Prefer ~/.local/bin (no sudo needed), fallback to /usr/local/bin
if [ -d "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
else
    mkdir -p "$HOME/.local/bin"
    INSTALL_DIR="$HOME/.local/bin"
fi

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

echo "Installing mcphub for ${OS}/${ARCH}..."

# Download CLI
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}-${OS}-${ARCH}"
TMPDIR=$(mktemp -d)
curl -fsSL "$DOWNLOAD_URL" -o "$TMPDIR/$BINARY"
chmod +x "$TMPDIR/$BINARY"
mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"

# Download MCP server binary
MCP_URL="https://github.com/${REPO}/releases/latest/download/${MCP_BINARY}-${OS}-${ARCH}"
if curl -fsSL "$MCP_URL" -o "$TMPDIR/$MCP_BINARY" 2>/dev/null; then
    chmod +x "$TMPDIR/$MCP_BINARY"
    mv "$TMPDIR/$MCP_BINARY" "$INSTALL_DIR/$MCP_BINARY"
fi

rm -rf "$TMPDIR"

# Check if INSTALL_DIR is in PATH
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        echo ""
        echo "Add this to your shell profile (~/.zshrc or ~/.bashrc):"
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        ;;
esac

echo ""
echo "mcphub installed successfully!"
echo ""
echo "Get started:"
echo "  mcphub search filesystem"
echo "  mcphub install <server-name>"
echo ""
echo "Add to Claude Code:"
echo "  claude mcp add mcphub mcphub-mcp"
echo ""
