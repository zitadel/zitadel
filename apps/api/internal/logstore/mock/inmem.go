package mock

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageStorer[*Record] = (*InmemLogStorage)(nil)
var _ logstore.LogCleanupper[*Record] = (*InmemLogStorage)(nil)
var _ logstore.Queries = (*InmemLogStorage)(nil)

type InmemLogStorage struct {
	mux     sync.Mutex
	clock   clock.Clock
	emitted []*Record
	bulks   []int
	quota   *query.Quota
}

func NewInMemoryStorage(clock clock.Clock, quota *query.Quota) *InmemLogStorage {
	return &InmemLogStorage{
		clock:   clock,
		emitted: make([]*Record, 0),
		bulks:   make([]int, 0),
		quota:   quota,
	}
}

func (l *InmemLogStorage) QuotaUnit() quota.Unit {
	return quota.Unimplemented
}

func (l *InmemLogStorage) Emit(_ context.Context, bulk []*Record) error {
	if len(bulk) == 0 {
		return nil
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	l.emitted = append(l.emitted, bulk...)
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

	clean := make([]*Record, 0)
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

func (l *InmemLogStorage) GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (qu *query.Quota, err error) {
	return l.quota, nil
}

func (l *InmemLogStorage) GetQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit, periodStart time.Time) (usage uint64, err error) {
	return uint64(l.Len()), nil
}

func (l *InmemLogStorage) GetRemainingQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit) (remaining *uint64, err error) {
	if !l.quota.Limit {
		return nil, nil
	}
	var r uint64
	used := uint64(l.Len())
	if used > l.quota.Amount {
		return &r, nil
	}
	r = l.quota.Amount - used
	return &r, nil
}

func (l *InmemLogStorage) GetDueQuotaNotifications(ctx context.Context, instanceID string, unit quota.Unit, qu *query.Quota, periodStart time.Time, usedAbs uint64) (dueNotifications []*quota.NotificationDueEvent, err error) {
	return nil, nil
}

func (l *InmemLogStorage) ReportQuotaUsage(ctx context.Context, dueNotifications []*quota.NotificationDueEvent) error {
	return nil
}
