package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WindowTime is a single tracked window entry with accumulated time.
type WindowTime struct {
	Label     string      `json:"label"`
	TimeMs    int64       `json:"time_ms"`
	Switches  []time.Time `json:"switches,omitempty"`
}

// Session is one tracking session (start → stop).
type Session struct {
	Start   time.Time    `json:"start"`
	End     time.Time    `json:"end"`
	Notes   string       `json:"notes"`
	Windows []WindowTime `json:"windows"`
}

func (s *Session) TotalMs() int64 {
	var total int64
	for _, w := range s.Windows {
		total += w.TimeMs
	}
	return total
}

// FormatTime converts milliseconds to H:MM:SS.
func FormatTime(ms int64) string {
	s := ms / 1000
	h := s / 3600
	m := (s % 3600) / 60
	sec := s % 60
	return fmt.Sprintf("%d:%02d:%02d", h, m, sec)
}

func saveFile(path string, sessions []Session) error {
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadFile(path string) ([]Session, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sessions []Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("could not parse file: %w", err)
	}
	return sessions, nil
}
