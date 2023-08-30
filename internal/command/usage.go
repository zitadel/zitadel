package command

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

func (c *Commands) IncrementUsageFromAccessLogs(ctx context.Context, instanceID string, authenticatedRequestsPeriodStart time.Time, records []*record.AccessLog) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("incrementing access relevant usage failed for at least one quota period: %w", err)
		}
	}()
	count := uint64(0)
	for _, r := range records {
		if r.IsAuthenticated() {
			count++
		}
	}
	return projection.QuotaProjection.IncrementUsage(ctx, quota.RequestsAllAuthenticated, instanceID, authenticatedRequestsPeriodStart, count)
}
