CREATE TABLE IF NOT EXISTS sample_events (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp      TEXT    NOT NULL,
    bundle_id      TEXT    NOT NULL DEFAULT '',
    app_name       TEXT    NOT NULL DEFAULT '',
    window_title   TEXT    NOT NULL DEFAULT '',
    browser_family TEXT    NOT NULL DEFAULT '',
    tab_title      TEXT    NOT NULL DEFAULT '',
    tab_url_host   TEXT    NOT NULL DEFAULT '',
    source         TEXT    NOT NULL,
    confidence     TEXT    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON sample_events(timestamp);

CREATE TABLE IF NOT EXISTS pomodoro_sessions (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at     TEXT    NOT NULL,
    ended_at       TEXT,
    work_minutes   INTEGER NOT NULL DEFAULT 25,
    break_minutes  INTEGER NOT NULL DEFAULT 5,
    status         TEXT    NOT NULL DEFAULT 'active',
    current_phase  TEXT    NOT NULL DEFAULT 'work'
);

CREATE TABLE IF NOT EXISTS time_segments (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time          TEXT    NOT NULL,
    end_time            TEXT    NOT NULL,
    bundle_id           TEXT    NOT NULL DEFAULT '',
    app_name            TEXT    NOT NULL DEFAULT '',
    tab_host            TEXT    NOT NULL DEFAULT '',
    tab_title           TEXT    NOT NULL DEFAULT '',
    seconds             INTEGER NOT NULL,
    pomodoro_session_id INTEGER,
    FOREIGN KEY (pomodoro_session_id) REFERENCES pomodoro_sessions(id)
);
CREATE INDEX IF NOT EXISTS idx_segments_start ON time_segments(start_time);

CREATE TABLE IF NOT EXISTS allowlist_entries (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    bundle_id  TEXT    NOT NULL,
    app_name   TEXT    NOT NULL,
    is_default INTEGER NOT NULL DEFAULT 0,
    session_id INTEGER,
    FOREIGN KEY (session_id) REFERENCES pomodoro_sessions(id)
);

CREATE TABLE IF NOT EXISTS focus_violations (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    timestamp  TEXT    NOT NULL,
    bundle_id  TEXT    NOT NULL,
    app_name   TEXT    NOT NULL DEFAULT '',
    action     TEXT    NOT NULL DEFAULT '',
    FOREIGN KEY (session_id) REFERENCES pomodoro_sessions(id)
);
