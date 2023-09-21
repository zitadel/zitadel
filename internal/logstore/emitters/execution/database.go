package execution

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageStorer[*record.ExecutionLog] = (*databaseLogStorage)(nil)

type databaseLogStorage struct {
	dbClient *database.DB
	commands *command.Commands
	queries  *query.Queries
}

func NewDatabaseLogStorage(dbClient *database.DB, commands *command.Commands, queries *query.Queries) *databaseLogStorage {
	return &databaseLogStorage{dbClient: dbClient, commands: commands, queries: queries}
}

func (l *databaseLogStorage) QuotaUnit() quota.Unit {
	return quota.ActionsAllRunsSeconds
}

func (l *databaseLogStorage) Emit(ctx context.Context, bulk []*record.ExecutionLog) error {
	if len(bulk) == 0 {
		return nil
	}
	return l.incrementUsage(ctx, bulk)
}

func (l *databaseLogStorage) incrementUsage(ctx context.Context, bulk []*record.ExecutionLog) (err error) {
	byInstance := make(map[string][]*record.ExecutionLog)
	for _, r := range bulk {
		if r.InstanceID != "" {
			byInstance[r.InstanceID] = append(byInstance[r.InstanceID], r)
		}
	}
	for instanceID, instanceBulk := range byInstance {
		q, getQuotaErr := l.queries.GetQuota(ctx, instanceID, quota.ActionsAllRunsSeconds)
		if errors.Is(getQuotaErr, sql.ErrNoRows) {
			continue
		}
		err = errors.Join(err, getQuotaErr)
		if getQuotaErr != nil {
			continue
		}
		sum, incrementErr := l.incrementUsageFromExecutionLogs(ctx, instanceID, q.CurrentPeriodStart, instanceBulk)
		err = errors.Join(err, incrementErr)
		if incrementErr != nil {
			continue
		}

		notifications, getNotificationErr := l.queries.GetDueQuotaNotifications(ctx, instanceID, quota.ActionsAllRunsSeconds, q, q.CurrentPeriodStart, sum)
		err = errors.Join(err, getNotificationErr)
		if getNotificationErr != nil || len(notifications) == 0 {
			continue
		}
		ctx = authz.WithInstanceID(ctx, instanceID)
		reportErr := l.commands.ReportQuotaUsage(ctx, notifications)
		err = errors.Join(err, reportErr)
		if reportErr != nil {
			continue
		}
	}
	return err
}

func (l *databaseLogStorage) incrementUsageFromExecutionLogs(ctx context.Context, instanceID string, periodStart time.Time, records []*record.ExecutionLog) (sum uint64, err error) {
	var total time.Duration
	for _, r := range records {
		total += r.Took
	}
	return projection.QuotaProjection.IncrementUsage(ctx, quota.ActionsAllRunsSeconds, instanceID, periodStart, uint64(math.Floor(total.Seconds())))
}
