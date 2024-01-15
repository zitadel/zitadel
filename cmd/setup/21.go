package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 18.sql
	addLimitFieldsToInstances string
)

type AddLimitFieldsToInstances struct {
	dbClient *database.DB
}

func (mig *AddLimitFieldsToInstances) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addLimitFieldsToInstances)
	return err
}

func (mig *AddLimitFieldsToInstances) String() string {
	return "21_add_limit_fields_to_instances"
}
