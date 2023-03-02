package access

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"google.golang.org/grpc/codes"

	"github.com/zitadel/zitadel/internal/api/call"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
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

var _ logstore.UsageQuerier = (*databaseLogStorage)(nil)
var _ logstore.LogCleanupper = (*databaseLogStorage)(nil)

type databaseLogStorage struct {
	dbClient *database.DB
}

func NewDatabaseLogStorage(dbClient *database.DB) *databaseLogStorage {
	return &databaseLogStorage{dbClient: dbClient}
}

func (l *databaseLogStorage) QuotaUnit() quota.Unit {
	return quota.RequestsAllAuthenticated
}

func (l *databaseLogStorage) Emit(ctx context.Context, bulk []logstore.LogRecord) error {
	if len(bulk) == 0 {
		return nil
	}
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
		item := bulk[idx].(*Record)
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

func (l *databaseLogStorage) QueryUsage(ctx context.Context, instanceId string, start time.Time) (uint64, error) {
	stmt, args, err := squirrel.Select(
		fmt.Sprintf("count(%s)", accessInstanceIdCol),
	).
		From(accessLogsTable + l.dbClient.Timetravel(call.Took(ctx))).
		Where(squirrel.And{
			squirrel.Eq{accessInstanceIdCol: instanceId},
			squirrel.GtOrEq{accessTimestampCol: start},
			squirrel.Expr(fmt.Sprintf(`%s #>> '{%s,0}' = '[REDACTED]'`, accessRequestHeadersCol, strings.ToLower(zitadel_http.Authorization))),
			squirrel.NotLike{accessRequestURLCol: "%/zitadel.system.v1.SystemService/%"},
			squirrel.NotLike{accessRequestURLCol: "%/system/v1/%"},
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{accessProtocolCol: HTTP},
					squirrel.NotEq{accessResponseStatusCol: http.StatusForbidden},
					squirrel.NotEq{accessResponseStatusCol: http.StatusInternalServerError},
					squirrel.NotEq{accessResponseStatusCol: http.StatusTooManyRequests},
				},
				squirrel.And{
					squirrel.Eq{accessProtocolCol: GRPC},
					squirrel.NotEq{accessResponseStatusCol: codes.PermissionDenied},
					squirrel.NotEq{accessResponseStatusCol: codes.Internal},
					squirrel.NotEq{accessResponseStatusCol: codes.ResourceExhausted},
				},
			},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, caos_errors.ThrowInternal(err, "ACCESS-V9Sde", "Errors.Internal")
	}

	var count uint64
	if err = l.dbClient.
		QueryRowContext(ctx, stmt, args...).
		Scan(&count); err != nil {
		return 0, caos_errors.ThrowInternal(err, "ACCESS-pBPrM", "Errors.Logstore.Access.ScanFailed")
	}

	return count, nil
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
