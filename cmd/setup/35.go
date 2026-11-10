package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 35/*.sql
	addPositionToEsWmIndex embed.FS
)

type AddPositionToIndexEsWm struct {
	dbClient *database.DB
}

func (mig *AddPositionToIndexEsWm) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(addPositionToEsWmIndex, "35")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "execute statement", "file", stmt.file, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *AddPositionToIndexEsWm) String() string {
	return "35_add_position_to_index_es_wm"
}
