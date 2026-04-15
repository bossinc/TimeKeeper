TimeKeeper
==========

Tells you how you are spending time on your computer.

## Build

```
go build ./src
```

## Install GNOME Extension

```
cd gnome-extension
./install.sh
gnome-extensions enable window-watcher@timekeeper
```

Log out and back in to activate (or on X11: `gnome-shell --replace &`).
