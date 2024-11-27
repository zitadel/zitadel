package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 39.sql
	deleteStaleOrgFields string
)

type DeleteStaleOrgFields struct {
	dbClient *database.DB
}

func (mig *DeleteStaleOrgFields) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, deleteStaleOrgFields)
	return err
}

func (mig *DeleteStaleOrgFields) String() string {
	return "39_delete_stale_org_fields"
}
