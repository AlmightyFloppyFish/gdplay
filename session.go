package gdplay

import (
	"time"
)

// audioAction is a Signal for a manage channel
type audioAction int

// signals for manage channel
const (
	audioStop = iota
	audioPause
	audioResume
)

// Wait untill playing is stopped or completed
func (s *AudioSession) Wait() {
	for {
		time.Sleep(300 * time.Millisecond)
		if s.IsPlaying {
			continue
		}
		return
	}
}

// Stop playing and optionally disconnect
func (s *AudioSession) Stop(disconnect bool) {
	s.manage <- audioStop
	if disconnect {
		s.Vc.Disconnect()
		s.Vc.Close()
	}
}

// Pause playing
func (s *AudioSession) Pause() {
	s.manage <- audioPause
}

// Resume playing
func (s *AudioSession) Resume() {
	s.manage <- audioResume
}
