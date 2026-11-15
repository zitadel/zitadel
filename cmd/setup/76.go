package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed 76.sql
var addUniqueConstraintClaimants string

type UniqueConstraintClaimants struct {
	dbClient *database.DB
}

func (mig *UniqueConstraintClaimants) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addUniqueConstraintClaimants)
	return err
}

func (mig *UniqueConstraintClaimants) String() string {
	return "76_unique_constraint_claimants"
}
