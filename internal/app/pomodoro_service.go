package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
	"github.com/jgiraldo/focustrack/internal/pomodoro"
)

type PomodoroRepository interface {
	CreateSession(ctx context.Context, s domain.PomodoroSession) (int64, error)
	UpdateSession(ctx context.Context, s domain.PomodoroSession) error
	GetActiveSession(ctx context.Context) (*domain.PomodoroSession, error)
	GetSessions(ctx context.Context) ([]domain.PomodoroSession, error)
	GetAllowlist(ctx context.Context, sessionID *int64) ([]domain.AllowlistEntry, error)
	AddToAllowlist(ctx context.Context, e domain.AllowlistEntry) (int64, error)
	RemoveFromAllowlist(ctx context.Context, id int64) error
	IsInAllowlist(ctx context.Context, sessionID int64, bundleID string) (bool, error)
	InsertViolation(ctx context.Context, v domain.FocusViolation) (int64, error)
	UpdateViolationAction(ctx context.Context, id int64, action domain.ViolationAction) error
}

type PomodoroService struct {
	repo         PomodoroRepository
	timer        *pomodoro.Timer
	mu           sync.RWMutex
	violation    *domain.FocusViolation
	lastViolTime time.Time
	cooldown     time.Duration
}

func NewPomodoroService(repo PomodoroRepository, cooldown time.Duration) *PomodoroService {
	return &PomodoroService{
		repo:     repo,
		timer:    pomodoro.NewTimer(),
		cooldown: cooldown,
	}
}

func (s *PomodoroService) Start(ctx context.Context, workMin, breakMin int) (*domain.PomodoroSession, error) {
	existing, err := s.repo.GetActiveSession(ctx)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrSessionActive
	}

	session := domain.PomodoroSession{
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  workMin,
		BreakMinutes: breakMin,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	}

	id, err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	session.ID = id

	s.timer.Start(&session, func(phase domain.Phase) {
		log.Printf("Pomodoro phase changed to: %s", phase)
		bgCtx := context.Background()
		sess := s.timer.Session()
		if sess != nil {
			s.repo.UpdateSession(bgCtx, *sess)
		}
	})

	return &session, nil
}

func (s *PomodoroService) Stop(ctx context.Context) error {
	s.timer.Stop()
	sess := s.timer.Session()
	if sess != nil {
		return s.repo.UpdateSession(ctx, *sess)
	}
	return nil
}

func (s *PomodoroService) GetState() domain.PomodoroState {
	state := s.timer.GetState()

	s.mu.RLock()
	state.Violation = s.violation
	s.mu.RUnlock()

	return state
}

func (s *PomodoroService) GetSessions(ctx context.Context) ([]domain.PomodoroSession, error) {
	return s.repo.GetSessions(ctx)
}

func (s *PomodoroService) CheckFocus(ctx context.Context, event domain.SampleEvent) {
	if !s.timer.IsActive() {
		return
	}

	sess := s.timer.Session()
	if sess == nil || sess.CurrentPhase != domain.PhaseWork {
		return
	}

	// Skip check for our own app
	if event.BundleID == "" {
		return
	}

	allowed, err := s.repo.IsInAllowlist(ctx, sess.ID, event.BundleID)
	if err != nil {
		log.Printf("Allowlist check error: %v", err)
		return
	}

	if allowed {
		s.mu.Lock()
		s.violation = nil
		s.mu.Unlock()
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if time.Since(s.lastViolTime) < s.cooldown {
		return
	}

	violation := domain.FocusViolation{
		SessionID: sess.ID,
		Timestamp: time.Now().UTC(),
		BundleID:  event.BundleID,
		AppName:   event.AppName,
	}

	id, err := s.repo.InsertViolation(ctx, violation)
	if err != nil {
		log.Printf("Insert violation error: %v", err)
		return
	}

	violation.ID = id
	s.violation = &violation
	s.lastViolTime = time.Now()
}

func (s *PomodoroService) DismissViolation(ctx context.Context, action domain.ViolationAction) error {
	s.mu.Lock()
	v := s.violation
	s.violation = nil
	s.mu.Unlock()

	if v != nil {
		return s.repo.UpdateViolationAction(ctx, v.ID, action)
	}
	return nil
}

func (s *PomodoroService) GetAllowlist(ctx context.Context, sessionID *int64) ([]domain.AllowlistEntry, error) {
	return s.repo.GetAllowlist(ctx, sessionID)
}

func (s *PomodoroService) AddToAllowlist(ctx context.Context, entry domain.AllowlistEntry) (int64, error) {
	return s.repo.AddToAllowlist(ctx, entry)
}

func (s *PomodoroService) RemoveFromAllowlist(ctx context.Context, id int64) error {
	return s.repo.RemoveFromAllowlist(ctx, id)
}
