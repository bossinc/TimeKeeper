package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// treeNode is one node in the hierarchical window tree.
type treeNode struct {
	label    string
	timeMs   int64
	children []string // child keys
}

// treeData holds the flattened tree. Keys use tab as separator between levels.
type treeData struct {
	roots []string
	nodes map[string]*treeNode
}

// buildTreeData converts a flat WindowTime list into a 3-level tree.
// The rightmost " - " segment is the root (app name), working left toward the leaf.
func buildTreeData(windows []WindowTime) *treeData {
	td := &treeData{nodes: make(map[string]*treeNode)}
	for _, w := range windows {
		parts := strings.Split(w.Label, " - ")

		rootKey := parts[len(parts)-1]
		if td.nodes[rootKey] == nil {
			td.nodes[rootKey] = &treeNode{label: rootKey}
			td.roots = append(td.roots, rootKey)
		}
		if len(parts) == 1 {
			td.nodes[rootKey].timeMs += w.TimeMs
			continue
		}

		subLabel := parts[len(parts)-2]
		subKey := rootKey + "\t" + subLabel
		if td.nodes[subKey] == nil {
			td.nodes[subKey] = &treeNode{label: subLabel}
			td.nodes[rootKey].children = append(td.nodes[rootKey].children, subKey)
		}
		if len(parts) == 2 {
			td.nodes[subKey].timeMs += w.TimeMs
			continue
		}

		leafLabel := strings.Join(parts[:len(parts)-2], " - ")
		leafKey := subKey + "\t" + leafLabel
		if td.nodes[leafKey] == nil {
			td.nodes[leafKey] = &treeNode{label: leafLabel}
			td.nodes[subKey].children = append(td.nodes[subKey].children, leafKey)
		}
		td.nodes[leafKey].timeMs += w.TimeMs
	}
	// Roll up times to parents.
	for _, rk := range td.roots {
		root := td.nodes[rk]
		for _, sk := range root.children {
			sub := td.nodes[sk]
			for _, lk := range sub.children {
				sub.timeMs += td.nodes[lk].timeMs
			}
			root.timeMs += sub.timeMs
		}
	}
	return td
}

// AppState holds all runtime state for the application.
type AppState struct {
	sessions    []Session
	currentFile string
	tracker     *Tracker
	selectedIdx int
	td          *treeData

	win          *gtk.Window
	startBtn     *gtk.Button
	timeLabel    *gtk.Label
	projectLabel *gtk.Label
	startLabel   *gtk.Label
	endLabel     *gtk.Label
	notesView    *gtk.TextView
	reviewBtn    *gtk.Button

	sessionStore *gtk.ListStore
	windowStore  *gtk.TreeStore
}

func newMainWindow() *gtk.Window {
	st := &AppState{
		tracker:     NewTracker(),
		selectedIdx: -1,
		td:          &treeData{nodes: make(map[string]*treeNode)},
	}

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	st.win = win
	win.SetTitle("TimeKeeper")
	win.SetDefaultSize(960, 620)
	win.Connect("delete-event", func() bool {
		if st.tracker.IsRunning() {
			st.doStop()
		}
		return false
	})
	win.Connect("destroy", gtk.MainQuit)

	st.buildUI()
	return win
}

func (st *AppState) buildUI() {
	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	// ── Menu bar ──────────────────────────────────────────────────────────────
	mainBox.PackStart(st.buildMenuBar(), false, false, 0)

	// ── Top controls row 1 ────────────────────────────────────────────────────
	topBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	topBox.SetMarginStart(6)
	topBox.SetMarginEnd(6)
	topBox.SetMarginTop(4)
	topBox.SetMarginBottom(4)

	st.startBtn, _ = gtk.ButtonNewWithLabel("Start")
	st.startBtn.SetSizeRequest(80, -1)
	st.startBtn.Connect("clicked", func() { st.onStartStop() })
	topBox.PackStart(st.startBtn, false, false, 0)

	st.timeLabel, _ = gtk.LabelNew("0:00:00")
	topBox.PackStart(st.timeLabel, false, false, 8)

	sep1, _ := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	topBox.PackStart(sep1, false, false, 4)

	fileLbl, _ := gtk.LabelNew("File:")
	topBox.PackStart(fileLbl, false, false, 0)
	st.projectLabel, _ = gtk.LabelNew("")
	topBox.PackStart(st.projectLabel, false, false, 0)

	sep2, _ := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	topBox.PackStart(sep2, false, false, 4)

	startLbl, _ := gtk.LabelNew("Start:")
	topBox.PackStart(startLbl, false, false, 0)
	st.startLabel, _ = gtk.LabelNew("")
	topBox.PackStart(st.startLabel, false, false, 0)

	endLbl, _ := gtk.LabelNew("End:")
	topBox.PackStart(endLbl, false, false, 0)
	st.endLabel, _ = gtk.LabelNew("")
	topBox.PackStart(st.endLabel, false, false, 0)

	mainBox.PackStart(topBox, false, false, 0)

	// ── Top controls row 2 ────────────────────────────────────────────────────
	topBox2, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	topBox2.SetMarginStart(6)
	topBox2.SetMarginEnd(6)
	topBox2.SetMarginBottom(4)

	drawingCheck, _ := gtk.CheckButtonNewWithLabel("Drawing Mode (disable AFK)")
	drawingCheck.Connect("toggled", func() {
		st.tracker.DrawingMode = drawingCheck.GetActive()
	})
	topBox2.PackStart(drawingCheck, false, false, 0)

	st.reviewBtn, _ = gtk.ButtonNewWithLabel("Review")
	st.reviewBtn.Connect("clicked", func() { st.onReview() })
	topBox2.PackStart(st.reviewBtn, false, false, 0)

	mainBox.PackStart(topBox2, false, false, 0)

	// ── Split pane ────────────────────────────────────────────────────────────
	paned, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	paned.SetPosition(200)

	sessionScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	sessionScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	sessionView, sessionStore := st.buildSessionList()
	st.sessionStore = sessionStore
	sessionScroll.Add(sessionView)
	paned.Pack1(sessionScroll, false, false)

	windowScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	windowScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	windowView, windowStore := st.buildWindowTree()
	st.windowStore = windowStore
	windowScroll.Add(windowView)
	paned.Pack2(windowScroll, true, false)

	mainBox.PackStart(paned, true, true, 0)

	// ── Notes ─────────────────────────────────────────────────────────────────
	notesBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	notesBox.SetMarginStart(6)
	notesBox.SetMarginEnd(6)
	notesBox.SetMarginTop(4)
	notesBox.SetMarginBottom(6)

	notesLbl, _ := gtk.LabelNew("Notes:")
	notesLbl.SetHAlign(gtk.ALIGN_START)
	notesBox.PackStart(notesLbl, false, false, 0)

	notesScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	notesScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	notesScroll.SetSizeRequest(-1, 90)
	st.notesView, _ = gtk.TextViewNew()
	st.notesView.SetWrapMode(gtk.WRAP_WORD_CHAR)
	notesScroll.Add(st.notesView)
	notesBox.PackStart(notesScroll, false, false, 0)

	mainBox.PackStart(notesBox, false, false, 0)

	st.win.Add(mainBox)
}

func (st *AppState) buildMenuBar() *gtk.MenuBar {
	menuBar, _ := gtk.MenuBarNew()

	fileMenu, _ := gtk.MenuNew()
	newItem, _ := gtk.MenuItemNewWithLabel("New")
	newItem.Connect("activate", func() { st.onNew() })
	openItem, _ := gtk.MenuItemNewWithLabel("Open…")
	openItem.Connect("activate", func() { st.onOpen() })
	fileMenu.Append(newItem)
	fileMenu.Append(openItem)
	fileMenuItem, _ := gtk.MenuItemNewWithLabel("File")
	fileMenuItem.SetSubmenu(fileMenu)
	menuBar.Append(fileMenuItem)

	editMenu, _ := gtk.MenuNew()
	deleteItem, _ := gtk.MenuItemNewWithLabel("Delete Selection")
	deleteItem.Connect("activate", func() { st.onDeleteSelection() })
	editMenu.Append(deleteItem)
	editMenuItem, _ := gtk.MenuItemNewWithLabel("Edit")
	editMenuItem.SetSubmenu(editMenu)
	menuBar.Append(editMenuItem)

	return menuBar
}

func (st *AppState) buildSessionList() (*gtk.TreeView, *gtk.ListStore) {
	store, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	view, _ := gtk.TreeViewNewWithModel(store)
	view.SetHeadersVisible(false)

	renderer, _ := gtk.CellRendererTextNew()
	col, _ := gtk.TreeViewColumnNewWithAttribute("Session", renderer, "text", 0)
	col.SetExpand(true)
	view.AppendColumn(col)

	sel, _ := view.GetSelection()
	sel.Connect("changed", func() {
		_, iter, ok := sel.GetSelected()
		if !ok {
			st.selectedIdx = -1
			return
		}
		path, _ := store.GetPath(iter)
		if indices := path.GetIndices(); len(indices) > 0 {
			st.selectedIdx = indices[0]
		}
	})

	return view, store
}

func (st *AppState) buildWindowTree() (*gtk.TreeView, *gtk.TreeStore) {
	store, _ := gtk.TreeStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	view, _ := gtk.TreeViewNewWithModel(store)

	r1, _ := gtk.CellRendererTextNew()
	c1, _ := gtk.TreeViewColumnNewWithAttribute("Window", r1, "text", 0)
	c1.SetExpand(true)
	view.AppendColumn(c1)

	r2, _ := gtk.CellRendererTextNew()
	c2, _ := gtk.TreeViewColumnNewWithAttribute("Time", r2, "text", 1)
	c2.SetMinWidth(80)
	view.AppendColumn(c2)

	return view, store
}

func (st *AppState) refreshSessionList() {
	st.sessionStore.Clear()
	for _, s := range st.sessions {
		iter := st.sessionStore.Append()
		text := fmt.Sprintf("%s  %s", s.Start.Format("2006-01-02 15:04"), FormatTime(s.TotalMs()))
		st.sessionStore.SetValue(iter, 0, text)
	}
}

func (st *AppState) refreshWindowTree(td *treeData) {
	st.windowStore.Clear()
	for _, rootKey := range td.roots {
		root := td.nodes[rootKey]
		rootIter := st.windowStore.Append(nil)
		st.windowStore.SetValue(rootIter, 0, root.label)
		st.windowStore.SetValue(rootIter, 1, FormatTime(root.timeMs))
		for _, subKey := range root.children {
			sub := td.nodes[subKey]
			subIter := st.windowStore.Append(rootIter)
			st.windowStore.SetValue(subIter, 0, sub.label)
			st.windowStore.SetValue(subIter, 1, FormatTime(sub.timeMs))
			for _, leafKey := range sub.children {
				leaf := td.nodes[leafKey]
				leafIter := st.windowStore.Append(subIter)
				st.windowStore.SetValue(leafIter, 0, leaf.label)
				st.windowStore.SetValue(leafIter, 1, FormatTime(leaf.timeMs))
			}
		}
	}
}

func (st *AppState) onStartStop() {
	if st.tracker.IsRunning() {
		st.doStop()
	} else {
		st.doStart()
	}
}

func (st *AppState) doStart() {
	st.tracker.Reset()
	st.windowStore.Clear()
	st.setNotesText("")
	st.startLabel.SetText(time.Now().Format("15:04:05"))
	st.endLabel.SetText("")
	st.startBtn.SetLabel("Stop")
	st.reviewBtn.SetSensitive(false)

	st.tracker.Start(
		func() {
			snap := st.tracker.Snapshot()
			td := buildTreeData(snap)
			totalMs := st.tracker.TotalMs()
			glib.IdleAdd(func() bool {
				st.refreshWindowTree(td)
				st.timeLabel.SetText(FormatTime(totalMs))
				return false
			})
		},
		func() {
			glib.IdleAdd(func() bool {
				st.endLabel.SetText(time.Now().Format("15:04:05"))
				st.startBtn.SetLabel("Start")
				st.reviewBtn.SetSensitive(true)
				st.finishSession()
				st.showInfo("AFK Detected",
					"You were idle for 5 minutes.\n"+
						"5 minutes have been deducted and tracking has stopped.")
				return false
			})
		},
	)
}

func (st *AppState) doStop() {
	st.tracker.Stop()
	st.endLabel.SetText(time.Now().Format("15:04:05"))
	st.startBtn.SetLabel("Start")
	st.reviewBtn.SetSensitive(true)
	st.finishSession()
}

func (st *AppState) finishSession() {
	session := st.tracker.ToSession(st.getNotesText())
	st.sessions = append(st.sessions, session)
	st.refreshSessionList()

	if st.currentFile == "" {
		st.promptSaveAs()
	} else {
		if err := saveFile(st.currentFile, st.sessions); err != nil {
			st.showError(err.Error())
		}
	}
}

func (st *AppState) promptSaveAs() {
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Save Session", st.win, gtk.FILE_CHOOSER_ACTION_SAVE,
		"_Cancel", gtk.RESPONSE_CANCEL,
		"_Save", gtk.RESPONSE_ACCEPT,
	)
	dlg.SetCurrentName("session.json")
	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.json")
	filter.SetName("JSON Session Files (*.json)")
	dlg.AddFilter(filter)

	if dlg.Run() == gtk.RESPONSE_ACCEPT {
		path := dlg.GetFilename()
		if !strings.HasSuffix(path, ".json") {
			path += ".json"
		}
		st.currentFile = path
		st.projectLabel.SetText(filepath.Base(path))
		st.win.SetTitle(filepath.Base(path) + " - TimeKeeper")
		if err := saveFile(st.currentFile, st.sessions); err != nil {
			st.showError(err.Error())
		}
	}
	dlg.Destroy()
}

func (st *AppState) onNew() {
	if st.tracker.IsRunning() {
		st.tracker.Stop()
	}
	st.sessions = nil
	st.currentFile = ""
	st.tracker.Reset()
	st.timeLabel.SetText("0:00:00")
	st.startLabel.SetText("")
	st.endLabel.SetText("")
	st.projectLabel.SetText("")
	st.setNotesText("")
	st.startBtn.SetLabel("Start")
	st.reviewBtn.SetSensitive(true)
	st.sessionStore.Clear()
	st.windowStore.Clear()
	st.win.SetTitle("TimeKeeper")
}

func (st *AppState) onOpen() {
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Open Session", st.win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL,
		"_Open", gtk.RESPONSE_ACCEPT,
	)
	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.json")
	filter.SetName("JSON Session Files (*.json)")
	dlg.AddFilter(filter)

	if dlg.Run() == gtk.RESPONSE_ACCEPT {
		path := dlg.GetFilename()
		sessions, err := loadFile(path)
		if err != nil {
			dlg.Destroy()
			st.showError("Could not open file: " + err.Error())
			return
		}
		st.currentFile = path
		st.sessions = sessions
		st.projectLabel.SetText(filepath.Base(path))
		st.win.SetTitle(filepath.Base(path) + " - TimeKeeper")
		st.refreshSessionList()
		st.windowStore.Clear()
		st.timeLabel.SetText("0:00:00")
		st.startLabel.SetText("")
		st.endLabel.SetText("")
	}
	dlg.Destroy()
}

func (st *AppState) onReview() {
	if st.tracker.IsRunning() {
		return
	}
	if st.selectedIdx < 0 || st.selectedIdx >= len(st.sessions) {
		st.showInfo("No Selection", "Select a session from the list first.")
		return
	}
	s := st.sessions[st.selectedIdx]
	st.refreshWindowTree(buildTreeData(s.Windows))
	st.timeLabel.SetText(FormatTime(s.TotalMs()))
	st.startLabel.SetText(s.Start.Format("15:04:05"))
	st.endLabel.SetText(s.End.Format("15:04:05"))
	st.setNotesText(s.Notes)
}

func (st *AppState) onDeleteSelection() {
	if st.selectedIdx < 0 || st.selectedIdx >= len(st.sessions) {
		return
	}
	st.sessions = append(st.sessions[:st.selectedIdx], st.sessions[st.selectedIdx+1:]...)
	st.selectedIdx = -1
	st.refreshSessionList()
	if st.currentFile != "" {
		if err := saveFile(st.currentFile, st.sessions); err != nil {
			st.showError(err.Error())
		}
	}
	st.windowStore.Clear()
	st.timeLabel.SetText("0:00:00")
	st.startLabel.SetText("")
	st.endLabel.SetText("")
	st.setNotesText("")
}

func (st *AppState) getNotesText() string {
	buf, err := st.notesView.GetBuffer()
	if err != nil {
		return ""
	}
	start, end := buf.GetBounds()
	text, _ := buf.GetText(start, end, false)
	return text
}

func (st *AppState) setNotesText(s string) {
	if buf, err := st.notesView.GetBuffer(); err == nil {
		buf.SetText(s)
	}
}

func (st *AppState) showError(msg string) {
	dlg := gtk.MessageDialogNew(st.win, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func (st *AppState) showInfo(title, msg string) {
	dlg := gtk.MessageDialogNew(st.win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", title)
	dlg.FormatSecondaryText("%s", msg)
	dlg.Run()
	dlg.Destroy()
}
