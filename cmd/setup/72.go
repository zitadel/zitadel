package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 72.sql
	addLockoutPolicyShowAbsoluteLockoutTimeColumn string
)

type LockoutPolicyShowAbsoluteLockoutTimeColumn struct {
	dbClient *database.DB
}

func (mig *LockoutPolicyShowAbsoluteLockoutTimeColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.DB.ExecContext(ctx, addLockoutPolicyShowAbsoluteLockoutTimeColumn)
	return err
}

func (mig *LockoutPolicyShowAbsoluteLockoutTimeColumn) String() string {
	return "72_lockout_policies3_add_show_absolute_lockout_time_column"
}
