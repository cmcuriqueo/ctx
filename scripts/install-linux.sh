#!/usr/bin/env bash
# Installs ctx to ~/.local/bin and ensures it is in PATH.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_DIR="${HOME}/.local/bin"

if ! command -v gcc &>/dev/null; then
    echo "Error: gcc is required for Tree-sitter (CGO). Install build-essential or gcc." >&2
    exit 1
fi

mkdir -p "$INSTALL_DIR"

echo "Building ctx..."
cd "$REPO_ROOT"
CGO_ENABLED=1 go build -o ctx ./cmd/ctx

if [[ ! -f "$REPO_ROOT/ctx" ]]; then
    echo "Build failed: ctx binary not found" >&2
    exit 1
fi

echo "Installing ctx to $INSTALL_DIR..."
cp "$REPO_ROOT/ctx" "$INSTALL_DIR/ctx"
chmod +x "$INSTALL_DIR/ctx"

update_shell_rc() {
    local rc="$1"
    if [[ -f "$rc" ]]; then
        if ! grep -qF "export PATH=\"$INSTALL_DIR:\$PATH\"" "$rc" 2>/dev/null; then
            echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$rc"
            echo "Updated $rc"
        fi
    fi
}

update_shell_rc "$HOME/.bashrc"
update_shell_rc "$HOME/.zshrc"

echo ""
echo "Installation complete."
echo "Restart your terminal or run:"
echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
echo "Then run: ctx --help"
