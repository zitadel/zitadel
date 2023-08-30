package access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/database"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	accessLogsTable          = "logstore.access"
	accessTimestampCol       = "log_date"
	accessProtocolCol        = "protocol"
	accessRequestURLCol      = "request_url"
	accessResponseStatusCol  = "response_status"
	accessRequestHeadersCol  = "request_headers"
	accessResponseHeadersCol = "response_headers"
	accessInstanceIdCol      = "instance_id"
	accessProjectIdCol       = "project_id"
	accessRequestedDomainCol = "requested_domain"
	accessRequestedHostCol   = "requested_host"
)

var _ logstore.UsageStorer[*record.AccessLog] = (*databaseLogStorage)(nil)
var _ logstore.LogCleanupper[*record.AccessLog] = (*databaseLogStorage)(nil)

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
	incrementErr := l.incrementUsage(ctx, bulk)
	storeErr := l.store(ctx, bulk)
	joinedErr := errors.Join(incrementErr, storeErr)
	if joinedErr != nil {
		joinedErr = fmt.Errorf("storing access logs and/or incrementing quota usage failed: %w", joinedErr)
	}
	return joinedErr
}

func (l *databaseLogStorage) incrementUsage(ctx context.Context, bulk []*record.AccessLog) (err error) {
	byInstance := make(map[string][]*record.AccessLog)
	for _, r := range bulk {
		if r.InstanceID != "" {
			byInstance[r.InstanceID] = append(byInstance[r.InstanceID], r)
		}
	}
	for instanceID, instanceBulk := range byInstance {
		q, getQuotaErr := l.queries.GetQuota(ctx, instanceID, quota.RequestsAllAuthenticated)
		if errors.Is(getQuotaErr, sql.ErrNoRows) {
			getQuotaErr = nil
			continue
		}
		err = errors.Join(err, getQuotaErr)
		if getQuotaErr != nil {
			continue
		}
		incrementErr := l.commands.IncrementUsageFromAccessLogs(ctx, instanceID, q.CurrentPeriodStart, instanceBulk)
		err = errors.Join(err, incrementErr)
		if incrementErr != nil {
			continue
		}
	}
	return err
}

func (l *databaseLogStorage) store(ctx context.Context, bulk []*record.AccessLog) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("storing access logs failed: %w", err)
		}
	}()
	builder := squirrel.Insert(accessLogsTable).
		Columns(
			accessTimestampCol,
			accessProtocolCol,
			accessRequestURLCol,
			accessResponseStatusCol,
			accessRequestHeadersCol,
			accessResponseHeadersCol,
			accessInstanceIdCol,
			accessProjectIdCol,
			accessRequestedDomainCol,
			accessRequestedHostCol,
		).
		PlaceholderFormat(squirrel.Dollar)

	for idx := range bulk {
		item := bulk[idx]
		builder = builder.Values(
			item.LogDate,
			item.Protocol,
			item.RequestURL,
			item.ResponseStatus,
			item.RequestHeaders,
			item.ResponseHeaders,
			item.InstanceID,
			item.ProjectID,
			item.RequestedDomain,
			item.RequestedHost,
		)
	}

	stmt, args, err := builder.ToSql()
	if err != nil {
		return caos_errors.ThrowInternal(err, "ACCESS-KOS7I", "Errors.Internal")
	}

	result, err := l.dbClient.ExecContext(ctx, stmt, args...)
	if err != nil {
		return caos_errors.ThrowInternal(err, "ACCESS-alnT9", "Errors.Access.StorageFailed")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return caos_errors.ThrowInternal(err, "ACCESS-7KIpL", "Errors.Internal")
	}

	logging.WithFields("rows", rows).Debug("successfully stored access logs")
	return nil
}

func (l *databaseLogStorage) Cleanup(ctx context.Context, keep time.Duration) error {
	stmt, args, err := squirrel.Delete(accessLogsTable).
		Where(squirrel.LtOrEq{accessTimestampCol: time.Now().Add(-keep)}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return caos_errors.ThrowInternal(err, "ACCESS-2oTh6", "Errors.Internal")
	}

	execCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = l.dbClient.ExecContext(execCtx, stmt, args...)
	return err
}
