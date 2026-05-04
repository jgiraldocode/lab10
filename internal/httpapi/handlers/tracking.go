package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jgiraldo/focustrack/internal/app"
)

type TrackingHandler struct {
	trackerSvc *app.TrackerService
	segmentSvc *app.SegmentService
	exportSvc  *app.ExportService
}

func NewTrackingHandler(ts *app.TrackerService, ss *app.SegmentService, es *app.ExportService) *TrackingHandler {
	return &TrackingHandler{trackerSvc: ts, segmentSvc: ss, exportSvc: es}
}

func (h *TrackingHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	events, err := h.trackerSvc.GetEvents(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if events == nil {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *TrackingHandler) GetSegments(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	segments, err := h.segmentSvc.GetSegments(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, segments)
}

func (h *TrackingHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	timeline, err := h.segmentSvc.GetTimeline(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, timeline)
}

func (h *TrackingHandler) GetApps(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	apps, err := h.segmentSvc.GetAppSummaries(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, apps)
}

func (h *TrackingHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	data, err := h.exportSvc.ExportCSV(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=focustrack-"+date+".csv")
	w.Write(data)
}

func (h *TrackingHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	data, err := h.exportSvc.ExportJSON(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=focustrack-"+date+".json")
	w.Write(data)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
