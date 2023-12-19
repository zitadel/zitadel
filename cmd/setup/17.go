package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 17.sql
	addOffsetField string
)

type AddOffsetToCurrentStates struct {
	dbClient *database.DB
}

func (mig *AddOffsetToCurrentStates) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addOffsetField)
	return err
}

func (mig *AddOffsetToCurrentStates) String() string {
	return "17_add_offset_col_to_current_states"
}
