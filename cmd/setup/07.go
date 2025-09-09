package setup

import (
	"context"
	_ "embed"
	"strings"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 07/logstore.sql
	createLogstoreSchema07 string
	//go:embed 07/access.sql
	createAccessLogsTable07 string
	//go:embed 07/execution.sql
	createExecutionLogsTable07 string
)

type LogstoreTables struct {
	dbClient *database.DB
	username string
}

func (mig *LogstoreTables) Execute(ctx context.Context, _ eventstore.Event) error {
	stmt := strings.ReplaceAll(createLogstoreSchema07, "%[1]s", mig.username) + createAccessLogsTable07 + createExecutionLogsTable07
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *LogstoreTables) String() string {
	return "07_logstore"
}
