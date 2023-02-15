package setup

import (
	"context"
	"database/sql"
	"embed"
	"strings"
)

var (
	//go:embed 07/logstore.sql
	createLogstoreSchema07 string
	//go:embed 07/cockroach/access.sql
	//go:embed 07/postgres/access.sql
	createAccessLogsTable07 embed.FS
	//go:embed 07/cockroach/execution.sql
	//go:embed 07/postgres/execution.sql
	createExecutionLogsTable07 embed.FS
)

type LogstoreTables struct {
	dbClient *sql.DB
	username string
	dbType   string
}

func (mig *LogstoreTables) Execute(ctx context.Context) error {
	accessStmt, err := readStmt(createAccessLogsTable07, "07", mig.dbType, "access.sql")
	if err != nil {
		return err
	}
	executionStmt, err := readStmt(createExecutionLogsTable07, "07", mig.dbType, "execution.sql")
	if err != nil {
		return err
	}
	stmt := strings.ReplaceAll(createLogstoreSchema07, "%[1]s", mig.username) + accessStmt + executionStmt
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *LogstoreTables) String() string {
	return "07_logstore"
}
