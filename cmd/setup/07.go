package setup

import (
	"context"
	"database/sql"
	_ "embed"
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
	dbClient *sql.DB
}

func (mig *LogstoreTables) Execute(ctx context.Context) error {
	stmt := createLogstoreSchema07 + createAccessLogsTable07 + createExecutionLogsTable07
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *LogstoreTables) String() string {
	return "07_logstore"
}
