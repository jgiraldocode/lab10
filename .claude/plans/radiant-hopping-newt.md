# Plan: macOS App Tracking + Pomodoro Application

## Context

The user wants a complete macOS productivity app that tracks which applications are in the foreground (with browser tab detail via AppleScript), stores usage data in SQLite, provides a Pomodoro timer with an app allowlist and focus violation warnings, and presents everything through a Vue.js dashboard. No code exists yet -- this is greenfield.

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (stdlib HTTP, `modernc.org/sqlite`, `os/exec` for AppleScript) |
| Persistence | SQLite with WAL mode |
| Frontend | Vue 3 + Pinia + Tailwind CSS v4 + Vite |

## Project Structure

```
cmd/api/main.go
internal/
  config/config.go
  domain/
    sample_event.go, time_segment.go, pomodoro.go, errors.go
  app/
    tracker_service.go, segment_service.go, pomodoro_service.go, export_service.go
  repository/sqlite/
    db.go, migrations/001_create_tables.sql
    sample_event_repo.go, time_segment_repo.go, pomodoro_repo.go
  tracker/
    tracker.go, frontmost.go, browser.go, browser_safari.go, browser_chrome.go
  pomodoro/timer.go
  httpapi/
    router.go, middleware.go
    handlers/ (tracking.go, pomodoro.go, allowlist.go, export.go)
frontend/
  src/
    stores/ (tracking.ts, pomodoro.ts)
    views/ (DashboardView, PomodoroView, StatsView, ExportView)
    components/ui/ (Button, Card, Badge, Input, Modal, Progress)
    components/tracking/ (DailyTimeline, AppUsageCard, SegmentList)
    components/pomodoro/ (PomodoroTimer, AllowlistManager, FocusWarning, SessionHistory)
    components/stats/ (StatsOverview, FilterBar)
    lib/api.ts
    composables/ (usePolling.ts)
    router/index.ts
```

## Implementation Phases

### Phase 1: Backend Core (F0-F1)
1. Init Go module, create `cmd/api/main.go` with config, DB, router wiring
2. Domain types: `SampleEvent`, `TimeSegment` with enums (Source, Confidence, BrowserFamily)
3. SQLite setup: `db.go` with `modernc.org/sqlite`, embedded migration `001_create_tables.sql`
4. Tracker: `frontmost.go` using AppleScript via `osascript` to get app name + bundle ID + window title
5. `TrackerService`: polling goroutine at configurable interval (default 3s)
6. `SegmentService`: on-demand aggregation (merge consecutive same-app events)
7. `ExportService`: CSV/JSON formatters
8. Repos: `sample_event_repo.go`, `time_segment_repo.go`
9. HTTP handlers: `GET /api/events`, `GET /api/segments`, `GET /api/timeline`, `GET /api/export/*`

### Phase 2: Browser Tab Integration (F2-F3)
1. `browser_safari.go`: AppleScript for Safari active tab URL + title
2. `browser_chrome.go`: AppleScript for Chrome/Brave/Edge (shared Chromium template)
3. `browser.go`: registry mapping bundle IDs to browser families and script app names
4. Integrate into `MacTracker.Poll()`: detect browser, run tab script, extract host via `net/url`
5. Graceful degradation: catch errors, fall back to window title with `confidence=low`

### Phase 3: Pomodoro Backend (F4)
1. Domain: `PomodoroSession`, `AllowlistEntry`, `FocusViolation` types
2. `pomodoro/timer.go`: server-side countdown goroutine with phase transitions
3. `PomodoroService`: session lifecycle, allowlist checking, violation detection with 10s cooldown
4. `pomodoro_repo.go`: CRUD for sessions, allowlist, violations
5. HTTP handlers: `POST /api/pomodoro/start`, `GET /api/pomodoro/state`, `POST /api/pomodoro/stop`, allowlist CRUD, violation recording
6. Integration: `TrackerService` checks Pomodoro state on each poll, flags violations

### Phase 4: Frontend (all phases)
1. Scaffold Vue project with Vite + Pinia + Tailwind + Router
2. UI primitives: Button, Card, Badge, Input, Modal, Progress
3. `lib/api.ts`: fetch wrapper for all backend endpoints
4. `stores/tracking.ts`: segments, timeline, daily stats
5. `stores/pomodoro.ts`: timer state polling (1s when active), allowlist, violations
6. `DashboardView`: DailyTimeline (horizontal stacked bars by app), top apps summary
7. `PomodoroView`: PomodoroTimer (circular countdown), AllowlistManager, SessionHistory
8. `FocusWarning.vue`: modal overlay on violation -- "Return to focus" / "Add to allowlist"
9. `StatsView`: StatsOverview cards, FilterBar (date range, app filter, Pomodoro-only)
10. `ExportView`: date picker + format buttons + download

### Phase 5: Polish (F5)
1. Dark mode across all components
2. Loading/error states
3. Privacy: strip query strings from URLs before storage
4. Configurable cooldown, poll interval via UI

## Key Technical Decisions

- **AppleScript via `osascript`** (not CGO): simpler build, ~100ms per call, acceptable at 3s intervals
- **`modernc.org/sqlite`** (pure Go): no CGO dependency
- **Polling (not SSE)** for Pomodoro state: 1 req/s to localhost is negligible
- **On-demand segment aggregation**: merge events at query time, ~28K events/day max
- **SQLite WAL mode**: concurrent reads from HTTP + single writer from tracker goroutine

## SQLite Schema

Tables: `sample_events`, `time_segments`, `pomodoro_sessions`, `allowlist_entries`, `focus_violations` with appropriate indexes on timestamps and foreign keys.

## API Endpoints

| Group | Endpoints |
|-------|----------|
| Tracking | `GET /api/events`, `GET /api/segments`, `GET /api/timeline`, `GET /api/apps` |
| Pomodoro | `POST /api/pomodoro/start`, `GET /api/pomodoro/state`, `POST /api/pomodoro/stop`, `GET /api/pomodoro/sessions` |
| Allowlist | `GET/POST/DELETE /api/allowlist` |
| Export | `GET /api/export/csv`, `GET /api/export/json` |
| Config | `GET/PATCH /api/config` |

## Testing Strategy

- **Go unit tests**: domain types, segment merge algorithm (table-driven), Pomodoro allowlist/violation logic, AppleScript output parsing
- **Go integration tests**: SQLite repos with `:memory:` DB, HTTP handlers with `httptest`
- **Frontend tests (Vitest)**: Pinia stores with mocked API, FocusWarning component behavior

## Verification

1. `make test` -- all Go + frontend tests pass
2. `make dev` -- start backend + Vite dev server
3. Open http://localhost:5173 -- dashboard shows real-time app tracking
4. Start Pomodoro, switch to non-allowed app -- FocusWarning appears
5. Export CSV for today -- file downloads with correct data
