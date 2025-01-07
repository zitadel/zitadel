package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 43.sql
	replaceCurrentSequencesIndex string
)

type ReplaceCurrentSequencesIndex struct {
	dbClient *database.DB
}

func (mig *ReplaceCurrentSequencesIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, replaceCurrentSequencesIndex)
	return err
}

func (mig *ReplaceCurrentSequencesIndex) String() string {
	return "43_replace_current_sequences_index"
}
