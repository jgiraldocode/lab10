package pomodoro

import (
	"sync"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type Timer struct {
	mu        sync.RWMutex
	session   *domain.PomodoroSession
	phaseEnd  time.Time
	running   bool
	stopCh    chan struct{}
	onPhaseEnd func(domain.Phase)
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Start(session *domain.PomodoroSession, onPhaseEnd func(domain.Phase)) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.session = session
	t.onPhaseEnd = onPhaseEnd
	t.running = true
	t.stopCh = make(chan struct{})
	t.phaseEnd = time.Now().Add(time.Duration(session.WorkMinutes) * time.Minute)

	go t.tick()
}

func (t *Timer) tick() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			t.mu.Lock()
			if !t.running || t.session == nil {
				t.mu.Unlock()
				return
			}

			if time.Now().After(t.phaseEnd) {
				oldPhase := t.session.CurrentPhase
				if oldPhase == domain.PhaseWork {
					t.session.CurrentPhase = domain.PhaseBreak
					t.phaseEnd = time.Now().Add(time.Duration(t.session.BreakMinutes) * time.Minute)
				} else {
					now := time.Now()
					t.session.Status = domain.SessionCompleted
					t.session.EndedAt = &now
					t.running = false
				}
				cb := t.onPhaseEnd
				newPhase := t.session.CurrentPhase
				t.mu.Unlock()
				if cb != nil {
					cb(newPhase)
				}
				continue
			}
			t.mu.Unlock()
		}
	}
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		close(t.stopCh)
		t.running = false
		now := time.Now()
		if t.session != nil {
			t.session.Status = domain.SessionCancelled
			t.session.EndedAt = &now
		}
	}
}

func (t *Timer) GetState() domain.PomodoroState {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.running || t.session == nil {
		return domain.PomodoroState{Active: false}
	}

	remaining := time.Until(t.phaseEnd)
	if remaining < 0 {
		remaining = 0
	}

	var totalSecs int
	if t.session.CurrentPhase == domain.PhaseWork {
		totalSecs = t.session.WorkMinutes * 60
	} else {
		totalSecs = t.session.BreakMinutes * 60
	}

	return domain.PomodoroState{
		Active:           true,
		SessionID:        t.session.ID,
		Phase:            t.session.CurrentPhase,
		RemainingSeconds: int(remaining.Seconds()),
		TotalSeconds:     totalSecs,
	}
}

func (t *Timer) IsActive() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

func (t *Timer) Session() *domain.PomodoroSession {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.session
}
