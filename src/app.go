package main

// App holds all business state and logic, with no UI dependencies.
type App struct {
	Sessions    []Session
	CurrentFile string
	Tracker     *Tracker
	SelectedIdx int
}

func NewApp() *App {
	return &App{
		Tracker:     NewTracker(),
		SelectedIdx: -1,
	}
}

// AddSession finalises the current tracking run and appends it to Sessions.
func (a *App) AddSession(notes string) {
	a.Sessions = append(a.Sessions, a.Tracker.ToSession(notes))
}

// DeleteSelected removes the selected session. Returns false if nothing is selected.
func (a *App) DeleteSelected() bool {
	if a.SelectedIdx < 0 || a.SelectedIdx >= len(a.Sessions) {
		return false
	}
	a.Sessions = append(a.Sessions[:a.SelectedIdx], a.Sessions[a.SelectedIdx+1:]...)
	a.SelectedIdx = -1
	return true
}

// SelectedSession returns the currently selected session, or false if none.
func (a *App) SelectedSession() (Session, bool) {
	if a.SelectedIdx < 0 || a.SelectedIdx >= len(a.Sessions) {
		return Session{}, false
	}
	return a.Sessions[a.SelectedIdx], true
}

// Load replaces the current sessions with those from path.
func (a *App) Load(path string) error {
	sessions, err := loadFile(path)
	if err != nil {
		return err
	}
	a.CurrentFile = path
	a.Sessions = sessions
	return nil
}

// Save writes sessions to CurrentFile. No-op if no file has been set.
func (a *App) Save() error {
	if a.CurrentFile == "" {
		return nil
	}
	return saveFile(a.CurrentFile, a.Sessions)
}

// Reset stops tracking and clears all state.
func (a *App) Reset() {
	if a.Tracker.IsRunning() {
		a.Tracker.Stop()
	}
	a.Sessions = nil
	a.CurrentFile = ""
	a.Tracker.Reset()
	a.SelectedIdx = -1
}
