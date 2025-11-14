package database

import "context"

type Migrator interface {
	// Migrate executes migrations to setup the database.
	// The method can be called once per running Zitadel.
	Migrate(ctx context.Context) error
}
