package database

import (
	"context"
)

// Connector abstracts the database driver.
type Connector interface {
	Connect(ctx context.Context) (Pool, error)
}
