package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type PomodoroRepo struct {
	db *sql.DB
}

func NewPomodoroRepo(db *sql.DB) *PomodoroRepo {
	return &PomodoroRepo{db: db}
}

func (r *PomodoroRepo) CreateSession(ctx context.Context, s domain.PomodoroSession) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO pomodoro_sessions (started_at, work_minutes, break_minutes, status, current_phase)
		 VALUES (?, ?, ?, ?, ?)`,
		s.StartedAt.UTC().Format(time.RFC3339),
		s.WorkMinutes, s.BreakMinutes,
		string(s.Status), string(s.CurrentPhase),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PomodoroRepo) UpdateSession(ctx context.Context, s domain.PomodoroSession) error {
	var endedAt *string
	if s.EndedAt != nil {
		v := s.EndedAt.UTC().Format(time.RFC3339)
		endedAt = &v
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE pomodoro_sessions SET ended_at=?, status=?, current_phase=? WHERE id=?`,
		endedAt, string(s.Status), string(s.CurrentPhase), s.ID,
	)
	return err
}

func (r *PomodoroRepo) GetActiveSession(ctx context.Context) (*domain.PomodoroSession, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, started_at, ended_at, work_minutes, break_minutes, status, current_phase
		 FROM pomodoro_sessions WHERE status = 'active' LIMIT 1`)
	return r.scanSession(row)
}

func (r *PomodoroRepo) GetSessions(ctx context.Context) ([]domain.PomodoroSession, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, started_at, ended_at, work_minutes, break_minutes, status, current_phase
		 FROM pomodoro_sessions ORDER BY started_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.PomodoroSession
	for rows.Next() {
		var s domain.PomodoroSession
		var startedAt string
		var endedAt sql.NullString
		var status, phase string
		if err := rows.Scan(&s.ID, &startedAt, &endedAt, &s.WorkMinutes, &s.BreakMinutes, &status, &phase); err != nil {
			return nil, err
		}
		s.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
		if endedAt.Valid {
			t, _ := time.Parse(time.RFC3339, endedAt.String)
			s.EndedAt = &t
		}
		s.Status = domain.SessionStatus(status)
		s.CurrentPhase = domain.Phase(phase)
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *PomodoroRepo) scanSession(row *sql.Row) (*domain.PomodoroSession, error) {
	var s domain.PomodoroSession
	var startedAt string
	var endedAt sql.NullString
	var status, phase string
	err := row.Scan(&s.ID, &startedAt, &endedAt, &s.WorkMinutes, &s.BreakMinutes, &status, &phase)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
	if endedAt.Valid {
		t, _ := time.Parse(time.RFC3339, endedAt.String)
		s.EndedAt = &t
	}
	s.Status = domain.SessionStatus(status)
	s.CurrentPhase = domain.Phase(phase)
	return &s, nil
}

// Allowlist

func (r *PomodoroRepo) GetAllowlist(ctx context.Context, sessionID *int64) ([]domain.AllowlistEntry, error) {
	query := `SELECT id, bundle_id, app_name, is_default, session_id FROM allowlist_entries WHERE is_default = 1`
	args := []any{}
	if sessionID != nil {
		query += ` OR session_id = ?`
		args = append(args, *sessionID)
	}
	query += ` ORDER BY app_name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.AllowlistEntry
	for rows.Next() {
		var e domain.AllowlistEntry
		var isDefault int
		var sessID sql.NullInt64
		if err := rows.Scan(&e.ID, &e.BundleID, &e.AppName, &isDefault, &sessID); err != nil {
			return nil, err
		}
		e.IsDefault = isDefault == 1
		if sessID.Valid {
			e.SessionID = &sessID.Int64
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *PomodoroRepo) AddToAllowlist(ctx context.Context, e domain.AllowlistEntry) (int64, error) {
	isDefault := 0
	if e.IsDefault {
		isDefault = 1
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO allowlist_entries (bundle_id, app_name, is_default, session_id)
		 VALUES (?, ?, ?, ?)`,
		e.BundleID, e.AppName, isDefault, e.SessionID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PomodoroRepo) RemoveFromAllowlist(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM allowlist_entries WHERE id = ?`, id)
	return err
}

func (r *PomodoroRepo) IsInAllowlist(ctx context.Context, sessionID int64, bundleID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM allowlist_entries WHERE bundle_id = ? AND (is_default = 1 OR session_id = ?)`,
		bundleID, sessionID,
	).Scan(&count)
	return count > 0, err
}

// Violations

func (r *PomodoroRepo) InsertViolation(ctx context.Context, v domain.FocusViolation) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO focus_violations (session_id, timestamp, bundle_id, app_name, action)
		 VALUES (?, ?, ?, ?, ?)`,
		v.SessionID, v.Timestamp.UTC().Format(time.RFC3339),
		v.BundleID, v.AppName, string(v.Action),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PomodoroRepo) UpdateViolationAction(ctx context.Context, id int64, action domain.ViolationAction) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE focus_violations SET action = ? WHERE id = ?`,
		string(action), id,
	)
	return err
}
