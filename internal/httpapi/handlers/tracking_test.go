package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jgiraldo/focustrack/internal/app"
	"github.com/jgiraldo/focustrack/internal/domain"
	"github.com/jgiraldo/focustrack/internal/repository/sqlite"
)

func setupTestServices(t *testing.T) (*app.TrackerService, *app.SegmentService, *app.ExportService) {
	t.Helper()
	db, err := sqlite.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	eventRepo := sqlite.NewSampleEventRepo(db)
	segmentSvc := app.NewSegmentService(eventRepo)
	exportSvc := app.NewExportService(segmentSvc)

	// Use a dummy tracker that doesn't poll
	trackerSvc := app.NewTrackerService(nil, eventRepo, nil, time.Hour)

	// Insert test data
	ctx := context.Background()
	now := time.Now().UTC()
	events := []domain.SampleEvent{
		{Timestamp: now.Add(-10 * time.Second), BundleID: "com.app.test", AppName: "TestApp", Source: domain.SourceAppleScript, Confidence: domain.ConfidenceHigh},
		{Timestamp: now.Add(-7 * time.Second), BundleID: "com.app.test", AppName: "TestApp", Source: domain.SourceAppleScript, Confidence: domain.ConfidenceHigh},
		{Timestamp: now.Add(-4 * time.Second), BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "github.com", Source: domain.SourceAppleScript, Confidence: domain.ConfidenceHigh},
		{Timestamp: now, BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "github.com", Source: domain.SourceAppleScript, Confidence: domain.ConfidenceHigh},
	}
	for _, e := range events {
		if _, err := eventRepo.Insert(ctx, e); err != nil {
			t.Fatalf("Insert: %v", err)
		}
	}

	return trackerSvc, segmentSvc, exportSvc
}

func TestGetEvents(t *testing.T) {
	trackerSvc, segmentSvc, exportSvc := setupTestServices(t)
	h := NewTrackingHandler(trackerSvc, segmentSvc, exportSvc)

	today := time.Now().UTC().Format("2006-01-02")
	req := httptest.NewRequest(http.MethodGet, "/api/events?date="+today, nil)
	w := httptest.NewRecorder()

	h.GetEvents(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var events []domain.SampleEvent
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 4 {
		t.Errorf("expected 4 events, got %d", len(events))
	}
}

func TestGetSegments(t *testing.T) {
	trackerSvc, segmentSvc, exportSvc := setupTestServices(t)
	h := NewTrackingHandler(trackerSvc, segmentSvc, exportSvc)

	today := time.Now().UTC().Format("2006-01-02")
	req := httptest.NewRequest(http.MethodGet, "/api/segments?date="+today, nil)
	w := httptest.NewRecorder()

	h.GetSegments(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var segments []domain.TimeSegment
	if err := json.NewDecoder(w.Body).Decode(&segments); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(segments) != 2 {
		t.Errorf("expected 2 segments (TestApp + Chrome), got %d", len(segments))
	}
}

func TestGetApps(t *testing.T) {
	trackerSvc, segmentSvc, exportSvc := setupTestServices(t)
	h := NewTrackingHandler(trackerSvc, segmentSvc, exportSvc)

	today := time.Now().UTC().Format("2006-01-02")
	req := httptest.NewRequest(http.MethodGet, "/api/apps?date="+today, nil)
	w := httptest.NewRecorder()

	h.GetApps(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var apps []app.AppSummary
	if err := json.NewDecoder(w.Body).Decode(&apps); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(apps) != 2 {
		t.Errorf("expected 2 apps, got %d", len(apps))
	}
}

func TestExportCSV(t *testing.T) {
	trackerSvc, segmentSvc, exportSvc := setupTestServices(t)
	h := NewTrackingHandler(trackerSvc, segmentSvc, exportSvc)

	today := time.Now().UTC().Format("2006-01-02")
	req := httptest.NewRequest(http.MethodGet, "/api/export/csv?date="+today, nil)
	w := httptest.NewRecorder()

	h.ExportCSV(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected text/csv, got %s", ct)
	}
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty CSV body")
	}
}
