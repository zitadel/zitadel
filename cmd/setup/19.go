package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 19.sql
	addCurrentSequencesIndex string
)

type AddCurrentSequencesIndex struct {
	dbClient *database.DB
}

func (mig *AddCurrentSequencesIndex) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addCurrentSequencesIndex)
	return err
}

func (mig *AddCurrentSequencesIndex) String() string {
	return "19_add_current_sequences_index"
}
