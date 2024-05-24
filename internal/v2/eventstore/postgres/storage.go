package postgres

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	_ eventstore.Pusher  = (*Storage)(nil)
	_ eventstore.Querier = (*Storage)(nil)

	pushPositionStmt string
)

type Storage struct {
	client *database.DB
	config *Config
}

type Config struct {
	MaxRetries uint32
}

func New(client *database.DB, config *Config) *Storage {
	initPushStmt(client.Type())
	return &Storage{
		client: client,
		config: config,
	}
}

func initPushStmt(typ string) {
	switch typ {
	case "cockroach":
		pushPositionStmt = ", hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp()"
	case "postgres":
		pushPositionStmt = ", statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())"
	default:
		logging.WithFields("database_type", typ).Panic("position statement for type not implemented")
	}
}

// Health implements eventstore.Pusher.
func (s *Storage) Health(ctx context.Context) error {
	return s.client.PingContext(ctx)
}
