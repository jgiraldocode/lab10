package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	PollInterval time.Duration
	DBPath       string
	HTTPPort     int
	CooldownSecs int
	WorkMinutes  int
	BreakMinutes int
}

func Load() Config {
	c := Config{
		PollInterval: 3 * time.Second,
		DBPath:       "./data/tracker.db",
		HTTPPort:     8690,
		CooldownSecs: 10,
		WorkMinutes:  25,
		BreakMinutes: 5,
	}

	if v := os.Getenv("TRACKER_POLL_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.PollInterval = d
		}
	}
	if v := os.Getenv("TRACKER_DB_PATH"); v != "" {
		c.DBPath = v
	}
	if v := os.Getenv("TRACKER_HTTP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.HTTPPort = p
		}
	}
	if v := os.Getenv("TRACKER_COOLDOWN"); v != "" {
		if s, err := strconv.Atoi(v); err == nil {
			c.CooldownSecs = s
		}
	}

	return c
}
