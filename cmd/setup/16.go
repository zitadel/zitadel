package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 16.sql
	uniqueConstraintLower string
)

type UniqueConstraintToLower struct {
	dbClient *database.DB
}

func (mig *UniqueConstraintToLower) Execute(ctx context.Context) error {
	res, err := mig.dbClient.ExecContext(ctx, uniqueConstraintLower)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	logging.WithFields("count", count).Info("unique constraints updated")
	return err
}

func (mig *UniqueConstraintToLower) String() string {
	return "16_unique_constraint_lower"
}
