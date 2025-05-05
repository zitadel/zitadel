package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 55.sql
	createSessions8UserIDIndex string
)

type CreateSessions8UserIDIndex struct {
	dbClient *database.DB
}

func (mig *CreateSessions8UserIDIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createSessions8UserIDIndex)
	return err
}

func (mig *CreateSessions8UserIDIndex) String() string {
	return "55_create_sessions8_user_id_index"
}
