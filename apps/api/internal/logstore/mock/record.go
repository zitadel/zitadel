package mock

import (
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord[*Record] = (*Record)(nil)

func NewRecord(clock clock.Clock) *Record {
	return &Record{ts: clock.Now()}
}

type Record struct {
	ts       time.Time
	redacted bool
}

func (r Record) Normalize() *Record {
	r.redacted = true
	return &r
}
