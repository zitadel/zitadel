package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 65.sql
	userMetadata5Index string
)

type FixUserMetadata5Index struct {
	dbClient *database.DB
}

func (mig *FixUserMetadata5Index) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, userMetadata5Index)
	return err
}

func (mig *FixUserMetadata5Index) String() string {
	return "65_fix_user_metadata5_index"
}
