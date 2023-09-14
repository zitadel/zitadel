package access

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var _ logstore.UsageStorer[*record.AccessLog] = (*databaseLogStorage)(nil)

type databaseLogStorage struct {
	dbClient *database.DB
	commands *command.Commands
	queries  *query.Queries
}

func NewDatabaseLogStorage(dbClient *database.DB, commands *command.Commands, queries *query.Queries) *databaseLogStorage {
	return &databaseLogStorage{dbClient: dbClient, commands: commands, queries: queries}
}

func (l *databaseLogStorage) QuotaUnit() quota.Unit {
	return quota.RequestsAllAuthenticated
}

func (l *databaseLogStorage) Emit(ctx context.Context, bulk []*record.AccessLog) error {
	if len(bulk) == 0 {
		return nil
	}
	return l.incrementUsage(ctx, bulk)
}

func (l *databaseLogStorage) incrementUsage(ctx context.Context, bulk []*record.AccessLog) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	byInstance := make(map[string][]*record.AccessLog)
	for _, r := range bulk {
		if r.InstanceID != "" {
			byInstance[r.InstanceID] = append(byInstance[r.InstanceID], r)
		}
	}
	for instanceID, instanceBulk := range byInstance {
		q, getQuotaErr := l.queries.GetQuota(ctx, instanceID, quota.RequestsAllAuthenticated)
		if errors.Is(getQuotaErr, sql.ErrNoRows) {
			continue
		}
		err = errors.Join(err, getQuotaErr)
		if getQuotaErr != nil {
			continue
		}
		sum, incrementErr := l.incrementUsageFromAccessLogs(ctx, instanceID, q.CurrentPeriodStart, instanceBulk)
		err = errors.Join(err, incrementErr)
		if incrementErr != nil {
			continue
		}
		notifications, getNotificationErr := l.queries.GetDueQuotaNotifications(ctx, instanceID, quota.RequestsAllAuthenticated, q, q.CurrentPeriodStart, sum)
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

func (l *databaseLogStorage) incrementUsageFromAccessLogs(ctx context.Context, instanceID string, periodStart time.Time, records []*record.AccessLog) (sum uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var count uint64
	for _, r := range records {
		if r.IsAuthenticated() {
			count++
		}
	}
	return projection.QuotaProjection.IncrementUsage(ctx, quota.RequestsAllAuthenticated, instanceID, periodStart, count)
}
