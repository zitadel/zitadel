package database

import (
	"context"

	"github.com/zitadel/zitadel/backend/cmd/configure/bla4"
)

type Connector interface {
	Connect(ctx context.Context) (Pool, error)
	bla4.Configurer
}
