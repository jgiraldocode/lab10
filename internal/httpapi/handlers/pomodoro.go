package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jgiraldo/focustrack/internal/app"
	"github.com/jgiraldo/focustrack/internal/domain"
)

type PomodoroHandler struct {
	svc *app.PomodoroService
}

func NewPomodoroHandler(svc *app.PomodoroService) *PomodoroHandler {
	return &PomodoroHandler{svc: svc}
}

type startRequest struct {
	WorkMinutes  int `json:"work_minutes"`
	BreakMinutes int `json:"break_minutes"`
}

func (h *PomodoroHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WorkMinutes <= 0 {
		req.WorkMinutes = 25
	}
	if req.BreakMinutes <= 0 {
		req.BreakMinutes = 5
	}

	session, err := h.svc.Start(r.Context(), req.WorkMinutes, req.BreakMinutes)
	if err != nil {
		if err == domain.ErrSessionActive {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

func (h *PomodoroHandler) Stop(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Stop(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func (h *PomodoroHandler) GetState(w http.ResponseWriter, r *http.Request) {
	state := h.svc.GetState()
	writeJSON(w, http.StatusOK, state)
}

func (h *PomodoroHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.svc.GetSessions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

type violationActionRequest struct {
	Action string `json:"action"`
}

func (h *PomodoroHandler) HandleViolation(w http.ResponseWriter, r *http.Request) {
	var req violationActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	action := domain.ViolationAction(req.Action)
	if action != domain.ViolationReturned && action != domain.ViolationAddedToAllowlist && action != domain.ViolationDismissed {
		writeError(w, http.StatusBadRequest, "invalid action")
		return
	}

	if err := h.svc.DismissViolation(r.Context(), action); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
