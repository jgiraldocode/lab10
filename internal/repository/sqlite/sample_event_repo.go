package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type SampleEventRepo struct {
	db *sql.DB
}

func NewSampleEventRepo(db *sql.DB) *SampleEventRepo {
	return &SampleEventRepo{db: db}
}

func (r *SampleEventRepo) Insert(ctx context.Context, e domain.SampleEvent) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO sample_events (timestamp, bundle_id, app_name, window_title, browser_family, tab_title, tab_url_host, source, confidence)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.Timestamp.UTC().Format(time.RFC3339),
		e.BundleID, e.AppName, e.WindowTitle,
		string(e.BrowserFamily), e.TabTitle, e.TabURLHost,
		string(e.Source), string(e.Confidence),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SampleEventRepo) GetByDate(ctx context.Context, date string) ([]domain.SampleEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, timestamp, bundle_id, app_name, window_title, browser_family, tab_title, tab_url_host, source, confidence
		 FROM sample_events
		 WHERE date(timestamp) = ?
		 ORDER BY timestamp ASC`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.SampleEvent
	for rows.Next() {
		var e domain.SampleEvent
		var ts string
		var bf, src, conf string
		if err := rows.Scan(&e.ID, &ts, &e.BundleID, &e.AppName, &e.WindowTitle, &bf, &e.TabTitle, &e.TabURLHost, &src, &conf); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse(time.RFC3339, ts)
		e.BrowserFamily = domain.BrowserFamily(bf)
		e.Source = domain.Source(src)
		e.Confidence = domain.Confidence(conf)
		events = append(events, e)
	}
	return events, rows.Err()
}
