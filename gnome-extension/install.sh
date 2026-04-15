#!/usr/bin/env bash
set -e

UUID="window-watcher@timekeeper"
DEST="$HOME/.local/share/gnome-shell/extensions/$UUID"

mkdir -p "$DEST"
cp -r "$UUID/." "$DEST/"

echo "Installed to $DEST"
echo "Enable with: gnome-extensions enable $UUID"
echo "Then log out and back in, or run: gnome-shell --replace &  (X11 only)"
