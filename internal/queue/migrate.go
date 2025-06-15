package queue

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/zitadel/zitadel/internal/database"
)

type Migrator struct {
	driver riverdriver.Driver[pgx.Tx]
}

func NewMigrator(client *database.DB) *Migrator {
	return &Migrator{
		driver: riverpgxv5.New(client.Pool),
	}
}

func (m *Migrator) Execute(ctx context.Context) error {
	_, err := m.driver.GetExecutor().Exec(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema)
	if err != nil {
		return err
	}

	migrator, err := rivermigrate.New(m.driver, &rivermigrate.Config{Schema: schema})
	if err != nil {
		return err
	}
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err

}
