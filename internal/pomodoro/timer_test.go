package pomodoro

import (
	"testing"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

func TestTimer_StartsAndReportsState(t *testing.T) {
	timer := NewTimer()

	session := &domain.PomodoroSession{
		ID:           1,
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  1,
		BreakMinutes: 1,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	}

	timer.Start(session, func(phase domain.Phase) {})

	if !timer.IsActive() {
		t.Error("timer should be active")
	}

	state := timer.GetState()
	if !state.Active {
		t.Error("state should be active")
	}
	if state.Phase != domain.PhaseWork {
		t.Errorf("expected work phase, got %s", state.Phase)
	}
	if state.TotalSeconds != 60 {
		t.Errorf("expected 60 total seconds, got %d", state.TotalSeconds)
	}
	if state.RemainingSeconds <= 0 || state.RemainingSeconds > 60 {
		t.Errorf("unexpected remaining seconds: %d", state.RemainingSeconds)
	}

	timer.Stop()

	if timer.IsActive() {
		t.Error("timer should be stopped")
	}

	state = timer.GetState()
	if state.Active {
		t.Error("state should be inactive after stop")
	}
}

func TestTimer_StopSetsCancelled(t *testing.T) {
	timer := NewTimer()

	session := &domain.PomodoroSession{
		ID:           1,
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  25,
		BreakMinutes: 5,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	}

	timer.Start(session, func(phase domain.Phase) {})
	timer.Stop()

	sess := timer.Session()
	if sess.Status != domain.SessionCancelled {
		t.Errorf("expected cancelled, got %s", sess.Status)
	}
	if sess.EndedAt == nil {
		t.Error("expected EndedAt to be set")
	}
}
