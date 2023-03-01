package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

type EventstoreIndexesNew struct {
	dbClient *database.DB
	name     string
	step     string
	fileName string
	stmts    embed.FS
}

func (mig *EventstoreIndexesNew) Execute(ctx context.Context) error {
	stmt, err := readStmt(mig.stmts, mig.step, mig.dbClient.Type(), mig.fileName)
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *EventstoreIndexesNew) String() string {
	return mig.name
}
