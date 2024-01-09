package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 20.sql
	addByUserIndexToSession string
)

type AddByUserIndexToSession struct {
	dbClient *database.DB
}

func (mig *AddByUserIndexToSession) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addByUserIndexToSession)
	return err
}

func (mig *AddByUserIndexToSession) String() string {
	return "20_add_by_user_index_on_session"
}
