package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrSessionActive = errors.New("a pomodoro session is already active")
	ErrNoSession     = errors.New("no active pomodoro session")
)
