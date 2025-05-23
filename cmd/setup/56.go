package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 56.sql
	sessionRecoveryCodeCheckedColumn string
)

type SessionRecoveryCodeCheckedColumn struct {
	dbClient *database.DB
}

func (mig *SessionRecoveryCodeCheckedColumn) Execute(ctx context.Context, e eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, sessionRecoveryCodeCheckedColumn)
	return err
}

func (mig *SessionRecoveryCodeCheckedColumn) String() string {
	return "56_session_recovery_code_checked_column"
}
