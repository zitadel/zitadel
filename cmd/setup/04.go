package setup

import (
	"context"
	"database/sql"
	"embed"
)

var (
	//go:embed 04/cockroach/index.sql
	//go:embed 04/postgres/index.sql
	stmts embed.FS
)

type EventstoreIndexes struct {
	dbClient *sql.DB
	dbType   string
}

func (mig *EventstoreIndexes) Execute(ctx context.Context) error {
	stmt, err := readStmt(stmts, "04", mig.dbType, "index.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *EventstoreIndexes) String() string {
	return "04_eventstore_indexes"
}

func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(folder + "/" + typ + "/" + filename)
	return string(stmt), err
}
