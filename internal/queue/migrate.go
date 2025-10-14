package queue

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/zitadel/zitadel/internal/database"
)

type Migrator struct {
	driver riverdriver.Driver[*sql.Tx]
}

func NewMigrator(client *database.DB) *Migrator {
	return &Migrator{
		driver: riverdatabasesql.New(client.DB),
	}
}

func (m *Migrator) Execute(ctx context.Context) error {
	err := m.driver.GetExecutor().Exec(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema)
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
