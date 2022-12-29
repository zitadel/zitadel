package eventstore

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/command"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageReporter = (*esReporter)(nil)

type esReporter struct {
	dbClient *sql.DB
	commands *command.Commands
}

func NewEventstoreReporter(dbClient *sql.DB, commands *command.Commands) *esReporter {
	return &esReporter{
		dbClient: dbClient,
		commands: commands,
	}
}

func (e *esReporter) GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (*query.Quota, error) {
	return query.GetInstanceQuota(ctx, e.dbClient, instanceID, unit)
}
func (e *esReporter) Report(ctx context.Context, q *query.Quota, used uint64) (err error) {
	return e.commands.ReportUsage(ctx, q, used)
}
