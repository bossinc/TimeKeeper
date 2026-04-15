package main

import (
	"sync"
	"time"
)

const (
	afkThreshold = 5 * time.Minute
	tickInterval = 500 * time.Millisecond
)

// Tracker polls the active window every tickInterval and accumulates time.
type Tracker struct {
	mu          sync.Mutex
	entries     map[string]int64 // label → ms
	order       []string         // insertion order
	totalMs     int64
	isRunning   bool
	DrawingMode bool
	StartTime   time.Time

	ticker *time.Ticker
	done   chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{entries: make(map[string]int64)}
}

// Start begins tracking. onTick is called after each tick; onAFK is called if
// idle time exceeds 5 minutes (tracker is stopped before calling onAFK).
func (t *Tracker) Start(onTick func(), onAFK func()) {
	t.mu.Lock()
	if t.isRunning {
		t.mu.Unlock()
		return
	}
	t.isRunning = true
	t.StartTime = time.Now()
	t.ticker = time.NewTicker(tickInterval)
	t.done = make(chan struct{})
	t.mu.Unlock()

	go func() {
		lastTick := time.Now()
		for {
			select {
			case now := <-t.ticker.C:
				elapsed := now.Sub(lastTick).Milliseconds()
				lastTick = now

				if !t.DrawingMode {
					if getIdleTimeMs() >= afkThreshold.Milliseconds() {
						t.mu.Lock()
						// Deduct 5 minutes from the last window and total.
						deduct := afkThreshold.Milliseconds()
						if len(t.order) > 0 {
							last := t.order[len(t.order)-1]
							t.entries[last] -= deduct
							if t.entries[last] < 0 {
								t.entries[last] = 0
							}
						}
						t.totalMs -= deduct
						if t.totalMs < 0 {
							t.totalMs = 0
						}
						t.mu.Unlock()
						t.Stop()
						onAFK()
						return
					}
				}

				label := getActiveWindowTitle()
				if label == "" {
					label = "(unknown)"
				}

				t.mu.Lock()
				if _, exists := t.entries[label]; !exists {
					t.order = append(t.order, label)
					t.entries[label] = 0
				}
				t.entries[label] += elapsed
				t.totalMs += elapsed
				t.mu.Unlock()

				onTick()

			case <-t.done:
				return
			}
		}
	}()
}

func (t *Tracker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isRunning {
		return
	}
	t.isRunning = false
	t.ticker.Stop()
	close(t.done)
}

func (t *Tracker) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.isRunning
}

func (t *Tracker) TotalMs() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.totalMs
}

// Snapshot returns a copy of the current window entries in insertion order.
func (t *Tracker) Snapshot() []WindowTime {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]WindowTime, 0, len(t.order))
	for _, label := range t.order {
		out = append(out, WindowTime{Label: label, TimeMs: t.entries[label]})
	}
	return out
}

func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]int64)
	t.order = nil
	t.totalMs = 0
}

func (t *Tracker) ToSession(notes string) Session {
	return Session{
		Start:   t.StartTime,
		End:     time.Now(),
		Notes:   notes,
		Windows: t.Snapshot(),
	}
}
