package mock

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageQuerier = (*InmemLogStorage)(nil)
var _ logstore.LogCleanupper = (*InmemLogStorage)(nil)

type InmemLogStorage struct {
	mux     sync.Mutex
	clock   clock.Clock
	emitted []*record
	bulks   []int
}

func NewInMemoryStorage(clock clock.Clock) *InmemLogStorage {
	return &InmemLogStorage{
		clock:   clock,
		emitted: make([]*record, 0),
		bulks:   make([]int, 0),
	}
}

func (l *InmemLogStorage) QuotaUnit() quota.Unit {
	return quota.Unimplemented
}

func (l *InmemLogStorage) Emit(_ context.Context, bulk []logstore.LogRecord) error {
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

func (l *InmemLogStorage) QueryUsage(_ context.Context, _ string, start time.Time) (uint64, error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	var count uint64
	for _, r := range l.emitted {
		if r.ts.After(start) {
			count++
		}
	}
	return count, nil
}

func (l *InmemLogStorage) Cleanup(_ context.Context, keep time.Duration) error {
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

func (l *InmemLogStorage) Bulks() []int {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.bulks
}

func (l *InmemLogStorage) Len() int {
	l.mux.Lock()
	defer l.mux.Unlock()

	return len(l.emitted)
}
