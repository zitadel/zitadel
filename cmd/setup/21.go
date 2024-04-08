package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 21.sql
	addBlockFieldToLimits string
)

type AddBlockFieldToLimits struct {
	dbClient *database.DB
}

func (mig *AddBlockFieldToLimits) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addBlockFieldToLimits)
	return err
}

func (mig *AddBlockFieldToLimits) String() string {
	return "21_add_block_field_to_limits"
}
