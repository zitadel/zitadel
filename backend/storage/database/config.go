package database

import (
	"context"
)

type Connector interface {
	Connect(ctx context.Context) (Pool, error)
	// bla4.Configurer
}
