package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

func (q *Queries) GetCurrentQuotaPeriod(ctx context.Context, instanceID string, unit quota.Unit) (*quota.AddedEvent, time.Time, error) {
	rm, err := q.getQuotaReadModel(ctx, instanceID, instanceID, unit)
	if err != nil || !rm.active {
		return nil, time.Time{}, err
	}

	return rm.config, pushPeriodStart(rm.config.From, rm.config.ResetInterval, time.Now()), nil
}

func pushPeriodStart(from time.Time, interval time.Duration, now time.Time) time.Time {
	next := from.Add(interval)
	if next.After(now) {
		return from
	}
	return pushPeriodStart(next, interval, now)
}

func (q *Queries) getQuotaReadModel(ctx context.Context, instanceId, resourceOwner string, unit quota.Unit) (*quotaReadModel, error) {
	rm := newQuotaReadModel(instanceId, resourceOwner, unit)
	return rm, q.eventstore.FilterToQueryReducer(ctx, rm)
}
