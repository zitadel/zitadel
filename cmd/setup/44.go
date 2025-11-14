package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 44/*.sql
	replaceCurrentSequencesIndex embed.FS
)

type ReplaceCurrentSequencesIndex struct {
	dbClient *database.DB
}

func (mig *ReplaceCurrentSequencesIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(replaceCurrentSequencesIndex, "44")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *ReplaceCurrentSequencesIndex) String() string {
	return "44_replace_current_sequences_index"
}
