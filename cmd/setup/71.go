package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 71.sql
	addJWTAudience string
)

type JWTProvideAddAudienceColumn struct {
	dbClient *database.DB
}

func (mig *JWTProvideAddAudienceColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addJWTAudience)
	return err
}

func (mig *JWTProvideAddAudienceColumn) String() string {
	return "71_idp_templates6_jwt_add_audience"
}
