package app

import (
	"testing"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

func TestMergeEventsToSegments_Empty(t *testing.T) {
	result := MergeEventsToSegments(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestMergeEventsToSegments_SingleEvent(t *testing.T) {
	events := []domain.SampleEvent{
		{Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), BundleID: "com.app.test", AppName: "Test"},
	}

	result := MergeEventsToSegments(events)
	if len(result) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result))
	}
	if result[0].Seconds != 1 {
		t.Errorf("expected 1 second minimum, got %d", result[0].Seconds)
	}
	if result[0].BundleID != "com.app.test" {
		t.Errorf("expected com.app.test, got %s", result[0].BundleID)
	}
}

func TestMergeEventsToSegments_MergesSameApp(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []domain.SampleEvent{
		{Timestamp: base, BundleID: "com.app.a", AppName: "AppA"},
		{Timestamp: base.Add(3 * time.Second), BundleID: "com.app.a", AppName: "AppA"},
		{Timestamp: base.Add(6 * time.Second), BundleID: "com.app.a", AppName: "AppA"},
		{Timestamp: base.Add(9 * time.Second), BundleID: "com.app.a", AppName: "AppA"},
	}

	result := MergeEventsToSegments(events)
	if len(result) != 1 {
		t.Fatalf("expected 1 merged segment, got %d", len(result))
	}
	if result[0].Seconds != 9 {
		t.Errorf("expected 9 seconds, got %d", result[0].Seconds)
	}
}

func TestMergeEventsToSegments_DifferentApps(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []domain.SampleEvent{
		{Timestamp: base, BundleID: "com.app.a", AppName: "AppA"},
		{Timestamp: base.Add(3 * time.Second), BundleID: "com.app.a", AppName: "AppA"},
		{Timestamp: base.Add(6 * time.Second), BundleID: "com.app.b", AppName: "AppB"},
		{Timestamp: base.Add(9 * time.Second), BundleID: "com.app.b", AppName: "AppB"},
		{Timestamp: base.Add(12 * time.Second), BundleID: "com.app.a", AppName: "AppA"},
	}

	result := MergeEventsToSegments(events)
	if len(result) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(result))
	}

	if result[0].BundleID != "com.app.a" || result[0].Seconds != 6 {
		t.Errorf("segment 0: expected com.app.a 6s, got %s %ds", result[0].BundleID, result[0].Seconds)
	}
	if result[1].BundleID != "com.app.b" || result[1].Seconds != 6 {
		t.Errorf("segment 1: expected com.app.b 6s, got %s %ds", result[1].BundleID, result[1].Seconds)
	}
	if result[2].BundleID != "com.app.a" {
		t.Errorf("segment 2: expected com.app.a, got %s", result[2].BundleID)
	}
}

func TestMergeEventsToSegments_BrowserTabHost(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []domain.SampleEvent{
		{Timestamp: base, BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "github.com"},
		{Timestamp: base.Add(3 * time.Second), BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "github.com"},
		{Timestamp: base.Add(6 * time.Second), BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "google.com"},
		{Timestamp: base.Add(9 * time.Second), BundleID: "com.google.Chrome", AppName: "Chrome", TabURLHost: "google.com"},
	}

	result := MergeEventsToSegments(events)
	if len(result) != 2 {
		t.Fatalf("expected 2 segments (different tab hosts), got %d", len(result))
	}
	if result[0].TabHost != "github.com" {
		t.Errorf("expected github.com, got %s", result[0].TabHost)
	}
	if result[1].TabHost != "google.com" {
		t.Errorf("expected google.com, got %s", result[1].TabHost)
	}
}
