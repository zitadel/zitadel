package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
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

var _ logstore.UsageQuerier = (*databaseLogStorage)(nil)
var _ logstore.LogCleanupper = (*databaseLogStorage)(nil)

type databaseLogStorage struct {
	dbClient *database.DB
}

func NewDatabaseLogStorage(dbClient *database.DB) *databaseLogStorage {
	return &databaseLogStorage{dbClient: dbClient}
}

func (l *databaseLogStorage) QuotaUnit() quota.Unit {
	return quota.ActionsAllRunsSeconds
}

func (l *databaseLogStorage) Emit(ctx context.Context, bulk []logstore.LogRecord) error {
	if len(bulk) == 0 {
		return nil
	}
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
		item := bulk[idx].(*Record)

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
	if err = l.dbClient.
		QueryRowContext(ctx, stmt, args...).
		Scan(&durationSeconds); err != nil {
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
