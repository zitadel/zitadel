package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 12/12_add_otp_columns.sql
	addOTPColumns string
)

type AddOTPColumns struct {
	dbClient *database.DB
}

func (mig *AddOTPColumns) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addOTPColumns)
	return err
}

func (mig *AddOTPColumns) String() string {
	return "12_auth_users_otp_columns"
}
