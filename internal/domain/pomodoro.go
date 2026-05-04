package domain

import "time"

type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionCompleted SessionStatus = "completed"
	SessionCancelled SessionStatus = "cancelled"
)

type Phase string

const (
	PhaseWork  Phase = "work"
	PhaseBreak Phase = "break"
)

type ViolationAction string

const (
	ViolationReturned         ViolationAction = "returned"
	ViolationAddedToAllowlist ViolationAction = "added_to_allowlist"
	ViolationDismissed        ViolationAction = "dismissed"
)

type PomodoroSession struct {
	ID           int64         `json:"id"`
	StartedAt    time.Time     `json:"started_at"`
	EndedAt      *time.Time    `json:"ended_at,omitempty"`
	WorkMinutes  int           `json:"work_minutes"`
	BreakMinutes int           `json:"break_minutes"`
	Status       SessionStatus `json:"status"`
	CurrentPhase Phase         `json:"current_phase"`
}

type AllowlistEntry struct {
	ID        int64  `json:"id"`
	BundleID  string `json:"bundle_id"`
	AppName   string `json:"app_name"`
	IsDefault bool   `json:"is_default"`
	SessionID *int64 `json:"session_id,omitempty"`
}

type FocusViolation struct {
	ID        int64           `json:"id"`
	SessionID int64           `json:"session_id"`
	Timestamp time.Time       `json:"timestamp"`
	BundleID  string          `json:"bundle_id"`
	AppName   string          `json:"app_name"`
	Action    ViolationAction `json:"action"`
}

type PomodoroState struct {
	Active           bool             `json:"active"`
	SessionID        int64            `json:"session_id,omitempty"`
	Phase            Phase            `json:"phase,omitempty"`
	RemainingSeconds int              `json:"remaining_seconds"`
	TotalSeconds     int              `json:"total_seconds"`
	Violation        *FocusViolation  `json:"violation,omitempty"`
}
