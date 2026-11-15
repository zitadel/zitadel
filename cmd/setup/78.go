package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 78_owners_column.sql
	uniqueConstraintsOwnersColumn string
	//go:embed 78_owners_gin.sql
	uniqueConstraintsOwnersGIN string
)

type UniqueConstraintOwners struct {
	dbClient *database.DB
}

func (mig *UniqueConstraintOwners) Execute(ctx context.Context, _ eventstore.Event) error {
	logging.Info(ctx, "add unique constraint owners column", "migration", mig.String())
	if _, err := mig.dbClient.ExecContext(ctx, uniqueConstraintsOwnersColumn); err != nil {
		return err
	}
	logging.Info(ctx, "create unique constraint owners gin index concurrently", "migration", mig.String())
	_, err := mig.dbClient.ExecContext(ctx, uniqueConstraintsOwnersGIN)
	return err
}

func (mig *UniqueConstraintOwners) String() string {
	return "78_unique_constraint_owners"
}
