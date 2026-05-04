package httpapi

import (
	"net/http"

	"github.com/jgiraldo/focustrack/internal/httpapi/handlers"
)

func NewRouter(
	trackingH *handlers.TrackingHandler,
	pomodoroH *handlers.PomodoroHandler,
	allowlistH *handlers.AllowlistHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Tracking
	mux.HandleFunc("GET /api/events", trackingH.GetEvents)
	mux.HandleFunc("GET /api/segments", trackingH.GetSegments)
	mux.HandleFunc("GET /api/timeline", trackingH.GetTimeline)
	mux.HandleFunc("GET /api/apps", trackingH.GetApps)

	// Export
	mux.HandleFunc("GET /api/export/csv", trackingH.ExportCSV)
	mux.HandleFunc("GET /api/export/json", trackingH.ExportJSON)

	// Pomodoro
	mux.HandleFunc("POST /api/pomodoro/start", pomodoroH.Start)
	mux.HandleFunc("POST /api/pomodoro/stop", pomodoroH.Stop)
	mux.HandleFunc("GET /api/pomodoro/state", pomodoroH.GetState)
	mux.HandleFunc("GET /api/pomodoro/sessions", pomodoroH.GetSessions)
	mux.HandleFunc("POST /api/pomodoro/violations", pomodoroH.HandleViolation)

	// Allowlist
	mux.HandleFunc("GET /api/allowlist", allowlistH.GetAllowlist)
	mux.HandleFunc("POST /api/allowlist", allowlistH.AddToAllowlist)
	mux.HandleFunc("DELETE /api/allowlist/{id}", allowlistH.RemoveFromAllowlist)

	return CORSMiddleware(mux)
}
