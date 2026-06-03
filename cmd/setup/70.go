package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	addHistoryCountToPasswordComplexityPolicy string
)

type AddHistoryCountToPasswordComplexityPolicy struct {
	dbClient *database.DB
}

func (mig *AddHistoryCountToPasswordComplexityPolicy) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addHistoryCountToPasswordComplexityPolicy)
	return err
}

func (mig *AddHistoryCountToPasswordComplexityPolicy) String() string {
	return "70_add_history_count_to_password_complexity_policy"
}
