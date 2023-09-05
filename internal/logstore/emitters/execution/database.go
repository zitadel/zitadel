package execution

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	executionLogsTable     = "logstore.execution"
	executionTimestampCol  = "log_date"
	executionTookCol       = "took"
	executionMessageCol    = "message"
	executionLogLevelCol   = "loglevel"
	executionInstanceIdCol = "instance_id"
	executionActionIdCol   = "action_id"
	executionMetadataCol   = "metadata"
)

var _ logstore.UsageStorer[*record.ExecutionLog] = (*databaseLogStorage)(nil)
var _ logstore.LogCleanupper[*record.ExecutionLog] = (*databaseLogStorage)(nil)

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
	incrementErr := l.incrementUsage(ctx, bulk)
	storeErr := l.store(ctx, bulk)
	joinedErr := errors.Join(incrementErr, storeErr)
	if joinedErr != nil {
		joinedErr = fmt.Errorf("storing execution logs and/or incrementing quota usage failed: %w", joinedErr)
	}
	return joinedErr
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
		sum, incrementErr := l.commands.IncrementUsageFromExecutionLogs(ctx, instanceID, q.CurrentPeriodStart, instanceBulk)
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

func (l *databaseLogStorage) store(ctx context.Context, bulk []*record.ExecutionLog) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("storing execution logs failed: %w", err)
		}
	}()
	builder := squirrel.Insert(executionLogsTable).
		Columns(
			executionTimestampCol,
			executionTookCol,
			executionMessageCol,
			executionLogLevelCol,
			executionInstanceIdCol,
			executionActionIdCol,
			executionMetadataCol,
		).
		PlaceholderFormat(squirrel.Dollar)

	for idx := range bulk {
		item := bulk[idx]

		var took interface{}
		if item.Took > 0 {
			took = item.Took
		}

		builder = builder.Values(
			item.LogDate,
			took,
			item.Message,
			item.LogLevel,
			item.InstanceID,
			item.ActionID,
			item.Metadata,
		)
	}

	stmt, args, err := builder.ToSql()
	if err != nil {
		return caos_errors.ThrowInternal(err, "EXEC-KOS7I", "Errors.Internal")
	}

	result, err := l.dbClient.ExecContext(ctx, stmt, args...)
	if err != nil {
		return caos_errors.ThrowInternal(err, "EXEC-0j6i5", "Errors.Access.StorageFailed")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return caos_errors.ThrowInternal(err, "EXEC-MGchJ", "Errors.Internal")
	}

	logging.WithFields("rows", rows).Debug("successfully stored execution logs")
	return nil
}

func (l *databaseLogStorage) QueryUsage(ctx context.Context, instanceId string, start time.Time) (uint64, error) {
	stmt, args, err := squirrel.Select(
		fmt.Sprintf("COALESCE(SUM(%s)::INT,0)", executionTookCol),
	).
		From(executionLogsTable + l.dbClient.Timetravel(call.Took(ctx))).
		Where(squirrel.And{
			squirrel.Eq{executionInstanceIdCol: instanceId},
			squirrel.GtOrEq{executionTimestampCol: start},
			squirrel.NotEq{executionTookCol: nil},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, caos_errors.ThrowInternal(err, "EXEC-DXtzg", "Errors.Internal")
	}

	var durationSeconds uint64
	err = l.dbClient.
		QueryRowContext(ctx,
			func(row *sql.Row) error {
				return row.Scan(&durationSeconds)
			},
			stmt, args...,
		)
	if err != nil {
		return 0, caos_errors.ThrowInternal(err, "EXEC-Ad8nP", "Errors.Logstore.Execution.ScanFailed")
	}
	return durationSeconds, nil
}

func (l *databaseLogStorage) Cleanup(ctx context.Context, keep time.Duration) error {
	stmt, args, err := squirrel.Delete(executionLogsTable).
		Where(squirrel.LtOrEq{executionTimestampCol: time.Now().Add(-keep)}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return caos_errors.ThrowInternal(err, "EXEC-Bja8V", "Errors.Internal")
	}

	execCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = l.dbClient.ExecContext(execCtx, stmt, args...)
	return err
}
