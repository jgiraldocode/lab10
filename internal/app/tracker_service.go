package app

import (
	"context"
	"log"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
	"github.com/jgiraldo/focustrack/internal/tracker"
)

type SampleEventRepository interface {
	Insert(ctx context.Context, e domain.SampleEvent) (int64, error)
	GetByDate(ctx context.Context, date string) ([]domain.SampleEvent, error)
}

type TrackerService struct {
	tracker      tracker.Tracker
	eventRepo    SampleEventRepository
	pomodoroSvc  *PomodoroService
	pollInterval time.Duration
	lastEvent    *domain.SampleEvent
}

func NewTrackerService(t tracker.Tracker, repo SampleEventRepository, pomSvc *PomodoroService, interval time.Duration) *TrackerService {
	return &TrackerService{
		tracker:      t,
		eventRepo:    repo,
		pomodoroSvc:  pomSvc,
		pollInterval: interval,
	}
}

func (s *TrackerService) Run(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	log.Printf("Tracker started (interval: %s)", s.pollInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Tracker stopped")
			return
		case <-ticker.C:
			event, err := s.tracker.Poll(ctx)
			if err != nil {
				log.Printf("Poll error: %v", err)
				continue
			}

			if s.isDuplicate(event) {
				continue
			}

			if _, err := s.eventRepo.Insert(ctx, event); err != nil {
				log.Printf("Insert error: %v", err)
				continue
			}

			s.lastEvent = &event

			if s.pomodoroSvc != nil {
				s.pomodoroSvc.CheckFocus(ctx, event)
			}
		}
	}
}

func (s *TrackerService) isDuplicate(event domain.SampleEvent) bool {
	if s.lastEvent == nil {
		return false
	}
	return s.lastEvent.BundleID == event.BundleID &&
		s.lastEvent.TabURLHost == event.TabURLHost &&
		s.lastEvent.WindowTitle == event.WindowTitle
}

func (s *TrackerService) GetEvents(ctx context.Context, date string) ([]domain.SampleEvent, error) {
	return s.eventRepo.GetByDate(ctx, date)
}
