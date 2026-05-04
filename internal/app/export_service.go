package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"
)

type ExportService struct {
	segmentSvc *SegmentService
}

func NewExportService(segSvc *SegmentService) *ExportService {
	return &ExportService{segmentSvc: segSvc}
}

func (s *ExportService) ExportCSV(ctx context.Context, date string) ([]byte, error) {
	segments, err := s.segmentSvc.GetSegments(ctx, date)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{"start_time", "end_time", "bundle_id", "app_name", "tab_host", "tab_title", "seconds"})

	for _, seg := range segments {
		w.Write([]string{
			seg.StartTime.Format(time.RFC3339),
			seg.EndTime.Format(time.RFC3339),
			seg.BundleID,
			seg.AppName,
			seg.TabHost,
			seg.TabTitle,
			fmt.Sprintf("%d", seg.Seconds),
		})
	}

	w.Flush()
	return buf.Bytes(), w.Error()
}

func (s *ExportService) ExportJSON(ctx context.Context, date string) ([]byte, error) {
	segments, err := s.segmentSvc.GetSegments(ctx, date)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(segments, "", "  ")
}
