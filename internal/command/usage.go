package command

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

func (c *Commands) IncrementUsageFromAccessLogs(ctx context.Context, instanceID string, periodStart time.Time, records []*record.AccessLog) (sum uint64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("incrementing access relevant usage failed for at least one quota period: %w", err)
		}
	}()
	var count uint64
	for _, r := range records {
		if r.IsAuthenticated() {
			count++
		}
	}
	return projection.QuotaProjection.IncrementUsage(ctx, quota.RequestsAllAuthenticated, instanceID, periodStart, count)
}

func (c *Commands) IncrementUsageFromExecutionLogs(ctx context.Context, instanceID string, periodStart time.Time, records []*record.ExecutionLog) (sum uint64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("incrementing access relevant usage failed for at least one quota period: %w", err)
		}
	}()
	var total time.Duration
	for _, r := range records {
		total += r.Took
	}
	return projection.QuotaProjection.IncrementUsage(ctx, quota.ActionsAllRunsSeconds, instanceID, periodStart, uint64(math.Floor(total.Seconds())))
}
