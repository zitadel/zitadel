package queue

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	new_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/database"
)

type Migrator struct {
	client *database.DB
}

func NewMigrator(client *database.DB) *Migrator {
	return &Migrator{
		client: client,
	}
}

func (m *Migrator) Execute(ctx context.Context) error {
	_, err := m.client.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema)
	if err != nil {
		return err
	}

	var mig migrator
	switch pool := m.client.DB.(type) {
	case *new_sql.Pool:
		mig, err = m.sqlMigrator(pool.DB)
	case *postgres.Pool:
		mig, err = m.pgxMigrator(pool.Pool)
	}
	if err != nil {
		return err
	}

	_, err = mig.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err

}

func (m *Migrator) pgxMigrator(pool *pgxpool.Pool) (migrator, error) {
	return rivermigrate.New(riverpgxv5.New(pool), &rivermigrate.Config{Schema: schema})
}

func (m *Migrator) sqlMigrator(pool *sql.DB) (migrator, error) {
	return rivermigrate.New(riverdatabasesql.New(pool), &rivermigrate.Config{Schema: schema})
}

type migrator interface {
	Migrate(ctx context.Context, direction rivermigrate.Direction, opts *rivermigrate.MigrateOpts) (*rivermigrate.MigrateResult, error)
}
