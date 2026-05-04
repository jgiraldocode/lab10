package app

import (
	"context"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type SegmentService struct {
	eventRepo SampleEventRepository
}

func NewSegmentService(repo SampleEventRepository) *SegmentService {
	return &SegmentService{eventRepo: repo}
}

func (s *SegmentService) GetSegments(ctx context.Context, date string) ([]domain.TimeSegment, error) {
	events, err := s.eventRepo.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}
	return MergeEventsToSegments(events), nil
}

func MergeEventsToSegments(events []domain.SampleEvent) []domain.TimeSegment {
	if len(events) == 0 {
		return nil
	}

	var segments []domain.TimeSegment
	current := segmentFromEvent(events[0])

	for i := 1; i < len(events); i++ {
		e := events[i]
		if e.BundleID == current.BundleID && e.TabURLHost == current.TabHost {
			current.EndTime = e.Timestamp
			current.Seconds = int(current.EndTime.Sub(current.StartTime).Seconds())
		} else {
			current.EndTime = e.Timestamp
			current.Seconds = int(current.EndTime.Sub(current.StartTime).Seconds())
			segments = append(segments, current)
			current = segmentFromEvent(e)
		}
	}

	if current.Seconds == 0 {
		current.Seconds = 1
	}
	segments = append(segments, current)

	return segments
}

func segmentFromEvent(e domain.SampleEvent) domain.TimeSegment {
	return domain.TimeSegment{
		StartTime: e.Timestamp,
		EndTime:   e.Timestamp,
		BundleID:  e.BundleID,
		AppName:   e.AppName,
		TabHost:   e.TabURLHost,
		TabTitle:  e.TabTitle,
		Seconds:   0,
	}
}

type AppSummary struct {
	BundleID     string `json:"bundle_id"`
	AppName      string `json:"app_name"`
	TotalSeconds int    `json:"total_seconds"`
}

func (s *SegmentService) GetAppSummaries(ctx context.Context, date string) ([]AppSummary, error) {
	segments, err := s.GetSegments(ctx, date)
	if err != nil {
		return nil, err
	}

	totals := make(map[string]*AppSummary)
	for _, seg := range segments {
		key := seg.BundleID
		if _, ok := totals[key]; !ok {
			totals[key] = &AppSummary{BundleID: seg.BundleID, AppName: seg.AppName}
		}
		totals[key].TotalSeconds += seg.Seconds
	}

	result := make([]AppSummary, 0, len(totals))
	for _, v := range totals {
		result = append(result, *v)
	}

	// Sort by total seconds descending
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].TotalSeconds > result[i].TotalSeconds {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

type TimelineEntry struct {
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	AppName   string  `json:"app_name"`
	BundleID  string  `json:"bundle_id"`
	TabHost   string  `json:"tab_host,omitempty"`
	Seconds   int     `json:"seconds"`
	Percent   float64 `json:"percent"`
}

func (s *SegmentService) GetTimeline(ctx context.Context, date string) ([]TimelineEntry, error) {
	segments, err := s.GetSegments(ctx, date)
	if err != nil {
		return nil, err
	}

	var totalSecs int
	for _, seg := range segments {
		totalSecs += seg.Seconds
	}

	entries := make([]TimelineEntry, len(segments))
	for i, seg := range segments {
		pct := 0.0
		if totalSecs > 0 {
			pct = float64(seg.Seconds) / float64(totalSecs) * 100
		}
		entries[i] = TimelineEntry{
			StartTime: seg.StartTime.Format(time.RFC3339),
			EndTime:   seg.EndTime.Format(time.RFC3339),
			AppName:   seg.AppName,
			BundleID:  seg.BundleID,
			TabHost:   seg.TabHost,
			Seconds:   seg.Seconds,
			Percent:   pct,
		}
	}

	return entries, nil
}
