package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	addLockoutPolicyAutoUnlockColumn string
)

type LockoutPolicyAutoUnlockColumn struct {
	dbClient *database.DB
}

func (mig *LockoutPolicyAutoUnlockColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.DB.ExecContext(ctx, addLockoutPolicyAutoUnlockColumn)
	return err
}

func (mig *LockoutPolicyAutoUnlockColumn) String() string {
	return "70_lockout_policies3_add_auto_unlock_after_min_column"
}
