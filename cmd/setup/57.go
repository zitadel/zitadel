package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 57.sql
	addSessionRecoveryCodeCheckedAt string
)

type SessionRecoveryCodeCheckedAt struct {
	dbClient *database.DB
}

func (mig *SessionRecoveryCodeCheckedAt) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addSessionRecoveryCodeCheckedAt)
	return err
}

func (mig *SessionRecoveryCodeCheckedAt) String() string {
	return "57_session_recovery_code_checked_at"
}
