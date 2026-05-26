package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 71.sql
	addLockoutPolicyShowRemainingLockoutTimeColumn string
)

type LockoutPolicyShowRemainingLockoutTimeColumn struct {
	dbClient *database.DB
}

func (mig *LockoutPolicyShowRemainingLockoutTimeColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.DB.ExecContext(ctx, addLockoutPolicyShowRemainingLockoutTimeColumn)
	return err
}

func (mig *LockoutPolicyShowRemainingLockoutTimeColumn) String() string {
	return "71_lockout_policies3_add_show_remaining_lockout_time_column"
}
