package execution

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/logstore"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	accessLogsTable          = "logstore.execution"
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
	// TODO: Implement
	return errors.New("not yet implemented")
}

func (l *databaseLogStorage) QueryUsage(ctx context.Context, instanceId string, start, end time.Time) (uint64, error) {
	// TODO: Implement
	return 0, errors.New("not yet implemented")
}

func (l *databaseLogStorage) Cleanup(ctx context.Context, keep time.Duration) error {
	// TODO: Implement
	return errors.New("not yet implemented")
}
