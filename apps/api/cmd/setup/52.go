package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 52/alter.sql
	renameTableIfNotExisting string
	//go:embed 52/check.sql
	checkIfTableIsExisting string
)

type IDPTemplate6LDAP2 struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6LDAP2) Execute(ctx context.Context, _ eventstore.Event) error {
	var count int
	err := mig.dbClient.QueryRowContext(ctx,
		func(row *sql.Row) error {
			if err := row.Scan(&count); err != nil {
				return err
			}
			return row.Err()
		},
		checkIfTableIsExisting,
	)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, renameTableIfNotExisting)
	return err
}

func (mig *IDPTemplate6LDAP2) String() string {
	return "52_idp_templates6_ldap2"
}
