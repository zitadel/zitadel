package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

func (c *Commands) GetCurrentQuotaPeriod(ctx context.Context, instanceID string, unit quota.Unit) (*quota.AddedEvent, time.Time, error) {
	wm, err := c.getQuotaWriteModel(ctx, instanceID, instanceID, unit)
	if err != nil || !wm.active {
		return nil, time.Time{}, err
	}

	return wm.config, pushPeriodStart(wm.config.From, wm.config.ResetInterval, time.Now()), nil
}

func pushPeriodStart(from time.Time, interval time.Duration, now time.Time) time.Time {
	next := from.Add(interval)
	if next.After(now) {
		return from
	}
	return pushPeriodStart(next, interval, now)
}
