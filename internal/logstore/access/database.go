package access

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
)

const (
	accessLogsTable          = "logstore.access"
	accessTimestampCol       = "ts"
	accessProtocolCol        = "protocol"
	accessRequestURLCol      = "request_url"
	accessResponseStatusCol  = "response_status"
	accessRequestHeadersCol  = "request_headers"
	accessResponseHeadersCol = "response_headers"
	accessInstanceIdCol      = "instance_id"
	accessProjectIdCol       = "project_id"
	accessRequestedDomainCol = "requested_domain"
	accessRequestedHostCol   = "requested_host"
	redacted                 = "[REDACTED]"
)

func newStorageBulkSink(dbClient *sql.DB) logstore.BulkSinkFunc {
	return func(ctx context.Context, bulk []any) error {
		return storeAccessLogs(ctx, dbClient, bulk)
	}
}

// Emit notification request events
// Notify centrally
func storeAccessLogs(ctx context.Context, dbClient *sql.DB, bulk []any) error {

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
		item := bulk[idx].(*logstore.AccessLogRecord)
		builder = builder.Values(
			item.Timestamp,
			item.Protocol,
			item.RequestURL,
			item.ResponseStatus,
			pruneHeaders(item.RequestHeaders),
			pruneHeaders(item.ResponseHeaders),
			item.InstanceID,
			item.ProjectID,
			item.RequestedDomain,
			item.RequestedHost,
		)
	}

	stmt, args, err := builder.ToSql()
	if err != nil {
		return caos_errors.ThrowInternal(err, "LOGCH-KOS7I", "Errors.Internal")
	}

	result, err := dbClient.ExecContext(ctx, stmt, args...)
	if err != nil {
		return caos_errors.ThrowInternal(err, "LOGCH-alnT9", "Errors.Access.StorageFailed")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return caos_errors.ThrowInternal(err, "LOGCH-7KIpL", "Errors.Internal")
	}

	logging.Debugf("successfully stored %d acccess logs", rows)
	return nil
}

func authenticatedInstanceRequests(ctx context.Context, dbClient *sql.DB, instanceId string) (uint64, error) {

	stmt, args, err := squirrel.Select(
		fmt.Sprintf("count(%s)", accessInstanceIdCol),
	).
		From(accessLogsTable + " AS OF SYSTEM TIME '-20s'").
		Where(squirrel.And{
			squirrel.Eq{accessInstanceIdCol: instanceId},
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{accessProtocolCol: logstore.HTTP},
					squirrel.Expr(fmt.Sprintf(`%s #>> '{%s,0}' = '[REDACTED]'`, accessRequestHeadersCol, zitadel_http.Authorization)),
					squirrel.NotEq{accessResponseStatusCol: http.StatusForbidden},
					squirrel.NotEq{accessResponseStatusCol: http.StatusInternalServerError},
					squirrel.NotEq{accessResponseStatusCol: http.StatusTooManyRequests},
				},
				squirrel.And{
					squirrel.Eq{accessProtocolCol: logstore.GRPC},
					squirrel.Expr(fmt.Sprintf(`%s #>> '{%s,0}' = '[REDACTED]'`, accessResponseHeadersCol, zitadel_http.Authorization)),
					squirrel.NotEq{accessResponseStatusCol: codes.ResourceExhausted},
				},
			},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, caos_errors.ThrowInternal(err, "QUOTA-V9Sde", "Errors.Internal")
	}

	var count uint64
	if err = dbClient.
		QueryRowContext(ctx, stmt, args...).
		Scan(&count); err != nil {
		return 0, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Access.ScanFailed")
	}

	return count, nil
}

func pruneHeaders(header http.Header) http.Header {
	clonedHeader := header.Clone()
	for key := range clonedHeader {
		if strings.ToLower(key) == strings.ToLower(zitadel_http.Authorization) {
			clonedHeader[key] = []string{redacted}
		}
	}
	return clonedHeader
}
