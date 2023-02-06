package execution

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	executionLogsTable     = "logstore.execution"
	executionTimestampCol  = "ts"
	executionTookCol       = "took_ms"
	executionMessageCol    = "message"
	executionLogLevelCol   = "loglevel"
	executionInstanceIdCol = "instance_id"
	executionProjectIdCol  = "project_id"
	executionActionIdCol   = "action_id"
	executionMetadataCol   = "metadata"
)

var _ logstore.UsageQuerier = (*databaseLogStorage)(nil)
var _ logstore.LogCleanupper = (*databaseLogStorage)(nil)

type databaseLogStorage struct {
	dbClient *sql.DB
}

func NewDatabaseLogStorage(dbClient *sql.DB) *databaseLogStorage {
	return &databaseLogStorage{dbClient: dbClient}
}

func (l *databaseLogStorage) QuotaUnit() quota.Unit {
	return quota.ActionsAllRunsSeconds
}

func (l *databaseLogStorage) Emit(ctx context.Context, bulk []logstore.LogRecord) error {
	builder := squirrel.Insert(executionLogsTable).
		Columns(
			executionTimestampCol,
			executionTookCol,
			executionMessageCol,
			executionLogLevelCol,
			executionInstanceIdCol,
			executionProjectIdCol,
			executionActionIdCol,
			executionMetadataCol,
		).
		PlaceholderFormat(squirrel.Dollar)

	for idx := range bulk {
		item := bulk[idx].(*Record)

		var took interface{}
		if item.TookMS > 0 {
			took = item.TookMS
		}

		builder = builder.Values(
			item.Timestamp,
			took,
			item.Message,
			item.LogLevel,
			item.InstanceID,
			item.ProjectID,
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

	logging.Debugf("successfully stored %d execution logs", rows)
	return nil
}

func (l *databaseLogStorage) QueryUsage(ctx context.Context, instanceId string, start time.Time) (uint64, error) {
	stmt, args, err := squirrel.Select(
		fmt.Sprintf("COALESCE(SUM(%s),0)", executionTookCol),
	).
		From(executionLogsTable).
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

	var milliSeconds uint64
	if err = l.dbClient.
		QueryRowContext(ctx, stmt, args...).
		Scan(&milliSeconds); err != nil {
		return 0, caos_errors.ThrowInternal(err, "EXEC-Ad8nP", "Errors.Access.ScanFailed")
	}

	return uint64(math.Ceil(float64(milliSeconds) / 1000)), nil
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
