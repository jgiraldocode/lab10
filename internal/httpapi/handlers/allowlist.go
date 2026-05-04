package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jgiraldo/focustrack/internal/app"
	"github.com/jgiraldo/focustrack/internal/domain"
)

type AllowlistHandler struct {
	svc *app.PomodoroService
}

func NewAllowlistHandler(svc *app.PomodoroService) *AllowlistHandler {
	return &AllowlistHandler{svc: svc}
}

func (h *AllowlistHandler) GetAllowlist(w http.ResponseWriter, r *http.Request) {
	var sessionID *int64
	if v := r.URL.Query().Get("session_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			sessionID = &id
		}
	}

	entries, err := h.svc.GetAllowlist(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entries == nil {
		entries = []domain.AllowlistEntry{}
	}
	writeJSON(w, http.StatusOK, entries)
}

type addAllowlistRequest struct {
	BundleID  string `json:"bundle_id"`
	AppName   string `json:"app_name"`
	IsDefault bool   `json:"is_default"`
	SessionID *int64 `json:"session_id,omitempty"`
}

func (h *AllowlistHandler) AddToAllowlist(w http.ResponseWriter, r *http.Request) {
	var req addAllowlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.BundleID == "" {
		writeError(w, http.StatusBadRequest, "bundle_id is required")
		return
	}

	entry := domain.AllowlistEntry{
		BundleID:  req.BundleID,
		AppName:   req.AppName,
		IsDefault: req.IsDefault,
		SessionID: req.SessionID,
	}

	id, err := h.svc.AddToAllowlist(r.Context(), entry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	entry.ID = id
	writeJSON(w, http.StatusCreated, entry)
}

func (h *AllowlistHandler) RemoveFromAllowlist(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.RemoveFromAllowlist(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
