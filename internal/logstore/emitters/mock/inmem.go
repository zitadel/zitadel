package mock

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageQuerier = (*inmemLogStorage)(nil)
var _ logstore.LogCleanupper = (*inmemLogStorage)(nil)

type inmemLogStorage struct {
	mux     sync.Mutex
	clock   clock.Clock
	emitted []*record
	bulks   []int
}

func NewInMemoryStorage(clock clock.Clock) *inmemLogStorage {
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
	l.mux.Lock()
	defer l.mux.Unlock()
	for idx := range bulk {
		l.emitted = append(l.emitted, bulk[idx].(*record))
	}
	l.bulks = append(l.bulks, len(bulk))
	return nil
}

func (l *inmemLogStorage) QueryUsage(_ context.Context, _ string, start, end time.Time) (uint64, error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	var count uint64
	for _, r := range l.emitted {
		if r.ts.After(start) && r.ts.Before(end) {
			count++
		}
	}
	return count, nil
}

func (l *inmemLogStorage) Cleanup(_ context.Context, keep time.Duration) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	clean := make([]*record, 0)
	from := l.clock.Now().Add(-(keep + 1))
	for _, r := range l.emitted {
		if r.ts.After(from) {
			clean = append(clean, r)
		}
	}
	l.emitted = clean
	return nil
}

func (l *inmemLogStorage) Bulks() []int {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.bulks
}

func (l *inmemLogStorage) Len() int {
	l.mux.Lock()
	defer l.mux.Unlock()

	return len(l.emitted)
}
