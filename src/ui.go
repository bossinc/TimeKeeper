package main

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// AppState holds all UI widget references. Business state lives in app.
type AppState struct {
	app *App

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
		app: NewApp(),
	}

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	st.win = win
	win.SetTitle("TimeKeeper")
	win.SetDefaultSize(960, 620)
	win.Connect("delete-event", func() bool {
		if st.app.Tracker.IsRunning() {
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
		st.app.Tracker.DrawingMode = drawingCheck.GetActive()
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
			st.app.SelectedIdx = -1
			return
		}
		path, _ := store.GetPath(iter)
		if indices := path.GetIndices(); len(indices) > 0 {
			st.app.SelectedIdx = indices[0]
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
	for _, s := range st.app.Sessions {
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
