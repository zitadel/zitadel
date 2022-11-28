package access

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"

	"github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
)

const (
	accessLogsTable          = "logstore.access"
	accessColTimestamp       = "ts"
	accessColProtocol        = "protocol"
	accessColRequestURL      = "request_url"
	accessColResponseStatus  = "response_status"
	accessColRequestHeaders  = "request_headers"
	accessColResponseHeaders = "response_headers"
	redacted                 = "[REDACTED]"
)

func newStorageBulkSink(dbClient *sql.DB, storedHandler logstore.StoredAccessLogsReducer) logstore.BulkSinkFunc {
	return func(ctx context.Context, bulk []any) error {
		return storeAccessLogs(ctx, dbClient, bulk, storedHandler)
	}
}

func storeAccessLogs(ctx context.Context, dbClient *sql.DB, bulk []any, storedHandler logstore.StoredAccessLogsReducer) error {

	builder := squirrel.Insert(accessLogsTable).
		Columns(
			accessColTimestamp,
			accessColProtocol,
			accessColRequestURL,
			accessColResponseStatus,
			accessColRequestHeaders,
			accessColResponseHeaders,
		).
		PlaceholderFormat(squirrel.Dollar)

	accessLogs := make([]*logstore.AccessLogRecord, len(bulk))
	for idx := range bulk {
		item := bulk[idx].(*logstore.AccessLogRecord)
		builder = builder.Values(
			item.Timestamp,
			item.Protocol,
			item.RequestURL,
			item.ResponseStatus,
			pruneHeaders(item.RequestHeaders),
			pruneHeaders(item.ResponseHeaders),
		)
		accessLogs[idx] = item
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
	storedHandler.Reduce(ctx, accessLogs)
	return nil
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
