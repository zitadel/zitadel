package mock

import (
	"context"
	"time"

	"github.com/thejerf/abtime"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageQuerier = (*inmemLogStorage)(nil)
var _ logstore.LogCleanupper = (*inmemLogStorage)(nil)

type inmemLogStorage struct {
	clock   *abtime.ManualTime
	emitted []*record
	bulks   []int
}

func NewInMemoryStorage(clock *abtime.ManualTime) *inmemLogStorage {
	return &inmemLogStorage{
		clock:   clock,
		emitted: make([]*record, 0),
		bulks:   make([]int, 0),
	}
}

func (l *inmemLogStorage) QuotaUnit() quota.Unit {
	return quota.Unimplemented
}

func (l *inmemLogStorage) Emit(_ context.Context, bulk []logstore.LogRecord) error {
	if len(bulk) == 0 {
		return nil
	}
	for idx := range bulk {
		l.emitted = append(l.emitted, bulk[idx].(*record))
	}
	l.bulks = append(l.bulks, len(bulk))
	return nil
}

func (l *inmemLogStorage) QueryUsage(_ context.Context, _ string, start, end time.Time) (uint64, error) {
	var count uint64
	for _, r := range l.emitted {
		if r.ts.After(start) && r.ts.Before(end) {
			count++
		}
	}
	return count, nil
}

func (l *inmemLogStorage) Cleanup(_ context.Context, keep time.Duration) error {
	clean := make([]*record, 0)
	from := l.clock.Now().Add(-keep)
	for _, r := range l.emitted {
		if r.ts.After(from) {
			clean = append(clean, r)
		}
	}
	l.emitted = clean
	return nil
}

func (l *inmemLogStorage) Bulks() []int {
	return l.bulks
}

func (l *inmemLogStorage) Len() int {
	return len(l.emitted)
}
