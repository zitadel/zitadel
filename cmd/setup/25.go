package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 25.sql
	addLowerFieldsToVerifiedEmail string
)

type User11AddLowerFieldsToVerifiedEmail struct {
	dbClient *database.DB
}

func (mig *User11AddLowerFieldsToVerifiedEmail) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addLowerFieldsToVerifiedEmail)
	return err
}

func (mig *User11AddLowerFieldsToVerifiedEmail) String() string {
	return "25_user12_add_lower_fields_to_verified_email"
}
