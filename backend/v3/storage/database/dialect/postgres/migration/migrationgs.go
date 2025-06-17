package migration

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

var migrations []*migrate.Migration

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	// we need to ensure that the schema exists before we can run the migration
	// because creating the migrations table already required the schema
	_, err := conn.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS zitadel")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewMigrator(ctx, conn, "zitadel.migrations")
	if err != nil {
		return err
	}
	migrator.Migrations = migrations
	return migrator.Migrate(ctx)
}

func registerSQLMigration(sequence int32, up, down string) {
	migrations = append(migrations, &migrate.Migration{
		Sequence: sequence,
		UpSQL:    up,
		DownSQL:  down,
	})
}
