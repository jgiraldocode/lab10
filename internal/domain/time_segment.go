package domain

import "time"

type TimeSegment struct {
	ID                int64  `json:"id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	BundleID          string    `json:"bundle_id"`
	AppName           string    `json:"app_name"`
	TabHost           string    `json:"tab_host,omitempty"`
	TabTitle          string    `json:"tab_title,omitempty"`
	Seconds           int       `json:"seconds"`
	PomodoroSessionID *int64    `json:"pomodoro_session_id,omitempty"`
}
