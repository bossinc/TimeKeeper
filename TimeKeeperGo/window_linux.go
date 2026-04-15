//go:build linux

package main

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
	gosync "sync"
	"sync/atomic"

	"github.com/godbus/dbus/v5"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/screensaver"
	"github.com/jezek/xgb/xproto"
)

// ── X11 ───────────────────────────────────────────────────────────────────

var (
	xOnce gosync.Once
	xConn *xgb.Conn
	xErr  error

	rootWindow          xproto.Window
	atomNetActiveWindow xproto.Atom
	atomNetWMName       xproto.Atom
	atomUTF8String      xproto.Atom
)

func initX() {
	xOnce.Do(func() {
		xConn, xErr = xgb.NewConn()
		if xErr != nil {
			return
		}
		setup := xproto.Setup(xConn)
		rootWindow = setup.Roots[0].Root
		if err := screensaver.Init(xConn); err != nil {
			xErr = err
			return
		}
		atomNetActiveWindow = internAtom("_NET_ACTIVE_WINDOW")
		atomNetWMName = internAtom("_NET_WM_NAME")
		atomUTF8String = internAtom("UTF8_STRING")
	})
}

func internAtom(name string) xproto.Atom {
	reply, err := xproto.InternAtom(xConn, false, uint16(len(name)), name).Reply()
	if err != nil {
		return xproto.AtomNone
	}
	return reply.Atom
}

// ── D-Bus (Wayland / GNOME) ───────────────────────────────────────────────

var (
	dbusOnce gosync.Once
	dbusConn *dbus.Conn
	dbusErr  error

	gnomeIntrospectOnce  gosync.Once
	gnomeIntrospectAvail bool

	// Set to 1 permanently after GetWindows returns AccessDenied (GNOME 41+).
	gnomeGetWindowsDenied atomic.Bool
)

func initDBus() {
	dbusOnce.Do(func() {
		dbusConn, dbusErr = dbus.SessionBus()
	})
}

func checkGnomeIntrospect() bool {
	gnomeIntrospectOnce.Do(func() {
		initDBus()
		if dbusErr != nil {
			return
		}
		obj := dbusConn.Object("org.gnome.Shell", "/org/gnome/Shell/Introspect")
		var xmlData string
		if err := obj.Call("org.freedesktop.DBus.Introspectable.Introspect", 0).Store(&xmlData); err != nil {
			return
		}
		gnomeIntrospectAvail = strings.Contains(xmlData, "org.gnome.Shell.Introspect")
	})
	return gnomeIntrospectAvail
}

// ── Active window title ───────────────────────────────────────────────────

func getActiveWindowTitle() string {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if title := gnomeActiveTitle(); title != "" {
			return title
		}
	}
	return x11ActiveTitle()
}

// gnomeActiveTitle tries the TimeKeeper extension first, then falls back to
// org.gnome.Shell.Introspect.GetWindows (restricted on GNOME 41+).
func gnomeActiveTitle() string {
	if title := extensionActiveTitle(); title != "" {
		return title
	}
	return introspectActiveTitle()
}

// extensionActiveTitle calls the window-watcher@timekeeper GNOME Shell extension.
func extensionActiveTitle() string {
	initDBus()
	if dbusErr != nil {
		return ""
	}
	obj := dbusConn.Object("io.timekeeper.WindowWatcher", "/io/timekeeper/WindowWatcher")
	var title string
	if err := obj.Call("io.timekeeper.WindowWatcher.GetActiveWindow", 0).Store(&title); err != nil {
		return ""
	}
	return title
}

// introspectActiveTitle uses org.gnome.Shell.Introspect.GetWindows as a fallback
// (works on GNOME ≤40 without unsafe-mode).
func introspectActiveTitle() string {
	if !checkGnomeIntrospect() || gnomeGetWindowsDenied.Load() {
		return ""
	}
	obj := dbusConn.Object("org.gnome.Shell", "/org/gnome/Shell/Introspect")
	// GetWindows returns a{ta{sv}}: map of window-id → property map.
	var windows map[uint64]map[string]dbus.Variant
	if err := obj.Call("org.gnome.Shell.Introspect.GetWindows", 0).Store(&windows); err != nil {
		var dbusErr dbus.Error
		if errors.As(err, &dbusErr) && dbusErr.Name == "org.freedesktop.DBus.Error.AccessDenied" {
			gnomeGetWindowsDenied.Store(true)
		}
		return ""
	}
	// The window with the highest focus-timestamp is the currently active one.
	var bestTitle string
	var bestTs uint32
	for _, props := range windows {
		tsVar, ok := props["focus-timestamp"]
		if !ok {
			continue
		}
		ts, ok := tsVar.Value().(uint32)
		if !ok || ts <= bestTs {
			continue
		}
		bestTs = ts
		if titleVar, ok := props["title"]; ok {
			if t, ok := titleVar.Value().(string); ok {
				bestTitle = t
			}
		}
	}
	return bestTitle
}

func x11ActiveTitle() string {
	initX()
	if xErr != nil {
		return ""
	}
	prop, err := xproto.GetProperty(xConn, false, rootWindow,
		atomNetActiveWindow, xproto.AtomWindow, 0, 1).Reply()
	if err != nil || len(prop.Value) < 4 {
		return ""
	}
	wid := xproto.Window(binary.LittleEndian.Uint32(prop.Value))

	name, err := xproto.GetProperty(xConn, false, wid,
		atomNetWMName, atomUTF8String, 0, 512).Reply()
	if err == nil && len(name.Value) > 0 {
		return strings.TrimRight(string(name.Value), "\x00")
	}
	name, err = xproto.GetProperty(xConn, false, wid,
		xproto.AtomWmName, xproto.AtomString, 0, 512).Reply()
	if err == nil && len(name.Value) > 0 {
		return strings.TrimRight(string(name.Value), "\x00")
	}
	return ""
}

// ── Idle time ─────────────────────────────────────────────────────────────

func getIdleTimeMs() int64 {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if ms := gnomeIdleTimeMs(); ms >= 0 {
			return ms
		}
	}
	return x11IdleTimeMs()
}

// gnomeIdleTimeMs queries GNOME Mutter's IdleMonitor for idle time in ms.
func gnomeIdleTimeMs() int64 {
	initDBus()
	if dbusErr != nil {
		return -1
	}
	obj := dbusConn.Object("org.gnome.Mutter.IdleMonitor",
		"/org/gnome/Mutter/IdleMonitor/Core")
	var idleMs uint64
	if err := obj.Call("org.gnome.Mutter.IdleMonitor.GetIdletime", 0).Store(&idleMs); err != nil {
		return -1
	}
	return int64(idleMs)
}

func x11IdleTimeMs() int64 {
	initX()
	if xErr != nil {
		return 0
	}
	reply, err := screensaver.QueryInfo(xConn, xproto.Drawable(rootWindow)).Reply()
	if err != nil {
		return 0
	}
	return int64(reply.MsSinceUserInput)
}
