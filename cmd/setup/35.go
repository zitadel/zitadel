package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 35_1.sql
	createEsWmTemp string
	//go:embed 35_2.sql
	dropAndRenameEsWm string
)

type AddPositionToIndexEsWm struct {
	dbClient *database.DB
}

func (mig *AddPositionToIndexEsWm) Execute(ctx context.Context, _ eventstore.Event) error {
	if _, err := mig.dbClient.ExecContext(ctx, createEsWmTemp); err != nil {
		return err
	}
	_, err := mig.dbClient.ExecContext(ctx, dropAndRenameEsWm)
	return err
}

func (mig *AddPositionToIndexEsWm) String() string {
	return "35_add_position_to_index_es_wm"
}
