package mock

import (
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*record)(nil)

func NewRecord(clock clock.Clock) *record {
	return &record{ts: clock.Now()}
}

type record struct {
	ts       time.Time
	redacted bool
}

func (r record) Normalize() logstore.LogRecord {
	r.redacted = true
	return &r
}
