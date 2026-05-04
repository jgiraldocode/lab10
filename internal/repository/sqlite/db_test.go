package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

func TestOpenMemory(t *testing.T) {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer db.Close()

	// Verify tables exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sample_events'").Scan(&count)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if count != 1 {
		t.Errorf("expected sample_events table, count=%d", count)
	}
}

func TestSampleEventRepo_InsertAndGet(t *testing.T) {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer db.Close()

	repo := NewSampleEventRepo(db)
	ctx := context.Background()

	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	event := domain.SampleEvent{
		Timestamp:     now,
		BundleID:      "com.apple.Safari",
		AppName:       "Safari",
		WindowTitle:   "GitHub - Test",
		BrowserFamily: domain.BrowserSafari,
		TabTitle:      "GitHub",
		TabURLHost:    "github.com",
		Source:        domain.SourceAppleScript,
		Confidence:    domain.ConfidenceHigh,
	}

	id, err := repo.Insert(ctx, event)
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}

	events, err := repo.GetByDate(ctx, "2024-06-15")
	if err != nil {
		t.Fatalf("GetByDate: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.BundleID != "com.apple.Safari" {
		t.Errorf("expected com.apple.Safari, got %s", e.BundleID)
	}
	if e.TabURLHost != "github.com" {
		t.Errorf("expected github.com, got %s", e.TabURLHost)
	}
	if e.BrowserFamily != domain.BrowserSafari {
		t.Errorf("expected safari, got %s", e.BrowserFamily)
	}
}

func TestPomodoroRepo_SessionLifecycle(t *testing.T) {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer db.Close()

	repo := NewPomodoroRepo(db)
	ctx := context.Background()

	// No active session initially
	sess, err := repo.GetActiveSession(ctx)
	if err != nil {
		t.Fatalf("GetActiveSession: %v", err)
	}
	if sess != nil {
		t.Error("expected no active session")
	}

	// Create session
	session := domain.PomodoroSession{
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  25,
		BreakMinutes: 5,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	}

	id, err := repo.CreateSession(ctx, session)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	session.ID = id

	// Get active session
	active, err := repo.GetActiveSession(ctx)
	if err != nil {
		t.Fatalf("GetActiveSession: %v", err)
	}
	if active == nil {
		t.Fatal("expected active session")
	}
	if active.WorkMinutes != 25 {
		t.Errorf("expected 25 work minutes, got %d", active.WorkMinutes)
	}

	// Complete session
	now := time.Now()
	session.Status = domain.SessionCompleted
	session.EndedAt = &now
	if err := repo.UpdateSession(ctx, session); err != nil {
		t.Fatalf("UpdateSession: %v", err)
	}

	active, err = repo.GetActiveSession(ctx)
	if err != nil {
		t.Fatalf("GetActiveSession: %v", err)
	}
	if active != nil {
		t.Error("expected no active session after completion")
	}
}

func TestPomodoroRepo_Allowlist(t *testing.T) {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer db.Close()

	repo := NewPomodoroRepo(db)
	ctx := context.Background()

	// Create a session first
	sessID := int64(1)
	_, err = repo.CreateSession(ctx, domain.PomodoroSession{
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  25,
		BreakMinutes: 5,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	})
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	// Add default entry
	_, err = repo.AddToAllowlist(ctx, domain.AllowlistEntry{
		BundleID:  "com.apple.Terminal",
		AppName:   "Terminal",
		IsDefault: true,
	})
	if err != nil {
		t.Fatalf("AddToAllowlist: %v", err)
	}

	// Add session-specific entry
	_, err = repo.AddToAllowlist(ctx, domain.AllowlistEntry{
		BundleID:  "com.google.Chrome",
		AppName:   "Chrome",
		IsDefault: false,
		SessionID: &sessID,
	})
	if err != nil {
		t.Fatalf("AddToAllowlist: %v", err)
	}

	// Check allowlist
	allowed, err := repo.IsInAllowlist(ctx, sessID, "com.apple.Terminal")
	if err != nil {
		t.Fatalf("IsInAllowlist: %v", err)
	}
	if !allowed {
		t.Error("Terminal should be in allowlist (default)")
	}

	allowed, err = repo.IsInAllowlist(ctx, sessID, "com.google.Chrome")
	if err != nil {
		t.Fatalf("IsInAllowlist: %v", err)
	}
	if !allowed {
		t.Error("Chrome should be in allowlist (session)")
	}

	allowed, err = repo.IsInAllowlist(ctx, sessID, "com.spotify.client")
	if err != nil {
		t.Fatalf("IsInAllowlist: %v", err)
	}
	if allowed {
		t.Error("Spotify should NOT be in allowlist")
	}
}

func TestPomodoroRepo_Violations(t *testing.T) {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer db.Close()

	repo := NewPomodoroRepo(db)
	ctx := context.Background()

	// Create session
	sessID, err := repo.CreateSession(ctx, domain.PomodoroSession{
		StartedAt:    time.Now().UTC(),
		WorkMinutes:  25,
		BreakMinutes: 5,
		Status:       domain.SessionActive,
		CurrentPhase: domain.PhaseWork,
	})
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	// Insert violation
	vID, err := repo.InsertViolation(ctx, domain.FocusViolation{
		SessionID: sessID,
		Timestamp: time.Now().UTC(),
		BundleID:  "com.spotify.client",
		AppName:   "Spotify",
	})
	if err != nil {
		t.Fatalf("InsertViolation: %v", err)
	}

	// Update action
	err = repo.UpdateViolationAction(ctx, vID, domain.ViolationReturned)
	if err != nil {
		t.Fatalf("UpdateViolationAction: %v", err)
	}
}
