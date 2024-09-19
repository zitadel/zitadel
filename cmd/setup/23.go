package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

var (
	//go:embed 23.sql
	correctGlobalUniqueConstraints string
)

type CorrectGlobalUniqueConstraints struct {
	dbClient *database.DB
}

func (mig *CorrectGlobalUniqueConstraints) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, correctGlobalUniqueConstraints)
	return err
}

func (mig *CorrectGlobalUniqueConstraints) String() string {
	return "23_correct_global_unique_constraints"
}
