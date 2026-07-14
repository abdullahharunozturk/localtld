#!/usr/bin/env bash
#
# localtld installer — when brew isn't used, clone the repo and link the bin onto PATH.
#   curl -fsSL https://localtld.sh | bash
#
# Transparency: this script only clones into ~/.local/share/localtld and creates a
# symlink at ~/.local/bin/localtld. Read the source first:
#   https://github.com/abdullahharunozturk/localtld/blob/master/install.sh

set -euo pipefail

REPO="https://github.com/abdullahharunozturk/localtld.git"
DEST="${XDG_DATA_HOME:-$HOME/.local/share}/localtld"
BIN_DIR="$HOME/.local/bin"

red()   { printf '\033[1;31m%s\033[0m\n' "$*"; }
green() { printf '\033[1;32m%s\033[0m\n' "$*"; }

[ "$(uname -s)" = "Darwin" ] || { red "Only macOS is supported for now."; exit 1; }
command -v git >/dev/null || { red "git is required."; exit 1; }

# When brew is available the formula is the real path (coming soon); for now, source install.
if [ -d "$DEST/.git" ]; then
  git -C "$DEST" pull --ff-only
else
  git clone --depth 1 "$REPO" "$DEST"
fi

mkdir -p "$BIN_DIR"
ln -sf "$DEST/bin/localtld" "$BIN_DIR/localtld"
chmod +x "$DEST/bin/localtld"

green "✓ localtld installed → $BIN_DIR/localtld"
case ":$PATH:" in
  *":$BIN_DIR:"*) : ;;
  *) printf '\nAdd it to your PATH:\n  export PATH="%s:$PATH"\n' "$BIN_DIR" ;;
esac
printf '\nNext step:\n  localtld setup\n'
