package database

import "context"

type Migrator interface {
	// Migrate executes migrations to setup the database.
	Migrate(ctx context.Context) error
}
