package tracker

import (
	"context"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type Tracker interface {
	Poll(ctx context.Context) (domain.SampleEvent, error)
}
