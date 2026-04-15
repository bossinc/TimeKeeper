package main

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (st *AppState) onStartStop() {
	if st.app.Tracker.IsRunning() {
		st.doStop()
	} else {
		st.doStart()
	}
}

func (st *AppState) doStart() {
	st.app.Tracker.Reset()
	st.windowStore.Clear()
	st.setNotesText("")
	st.startLabel.SetText(time.Now().Format("15:04:05"))
	st.endLabel.SetText("")
	st.startBtn.SetLabel("Stop")
	st.reviewBtn.SetSensitive(false)

	st.app.Tracker.Start(
		func() {
			snap := st.app.Tracker.Snapshot()
			td := buildTreeData(snap)
			totalMs := st.app.Tracker.TotalMs()
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
	st.app.Tracker.Stop()
	st.endLabel.SetText(time.Now().Format("15:04:05"))
	st.startBtn.SetLabel("Start")
	st.reviewBtn.SetSensitive(true)
	st.finishSession()
}

func (st *AppState) finishSession() {
	st.app.AddSession(st.getNotesText())
	st.refreshSessionList()

	if st.app.CurrentFile == "" {
		st.promptSaveAs()
	} else if err := st.app.Save(); err != nil {
		st.showError(err.Error())
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
		st.app.CurrentFile = path
		st.projectLabel.SetText(filepath.Base(path))
		st.win.SetTitle(filepath.Base(path) + " - TimeKeeper")
		if err := st.app.Save(); err != nil {
			st.showError(err.Error())
		}
	}
	dlg.Destroy()
}

func (st *AppState) onNew() {
	st.app.Reset()
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
		if err := st.app.Load(path); err != nil {
			dlg.Destroy()
			st.showError("Could not open file: " + err.Error())
			return
		}
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
	if st.app.Tracker.IsRunning() {
		return
	}
	s, ok := st.app.SelectedSession()
	if !ok {
		st.showInfo("No Selection", "Select a session from the list first.")
		return
	}
	st.refreshWindowTree(buildTreeData(s.Windows))
	st.timeLabel.SetText(FormatTime(s.TotalMs()))
	st.startLabel.SetText(s.Start.Format("15:04:05"))
	st.endLabel.SetText(s.End.Format("15:04:05"))
	st.setNotesText(s.Notes)
}

func (st *AppState) onDeleteSelection() {
	if !st.app.DeleteSelected() {
		return
	}
	st.refreshSessionList()
	if err := st.app.Save(); err != nil {
		st.showError(err.Error())
	}
	st.windowStore.Clear()
	st.timeLabel.SetText("0:00:00")
	st.startLabel.SetText("")
	st.endLabel.SetText("")
	st.setNotesText("")
}
