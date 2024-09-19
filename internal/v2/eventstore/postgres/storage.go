package postgres

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
)

var (
	_ eventstore.Pusher  = (*Storage)(nil)
	_ eventstore.Querier = (*Storage)(nil)
)

type Storage struct {
	client           *database.DB
	config           *Config
	pushPositionStmt string
}

type Config struct {
	MaxRetries uint32
}

func New(client *database.DB, config *Config) *Storage {
	initPushStmt(client.Type())
	return &Storage{
		client:           client,
		config:           config,
		pushPositionStmt: initPushStmt(client.Type()),
	}
}

func initPushStmt(typ string) string {
	switch typ {
	case "cockroach":
		return ", hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp()"
	case "postgres":
		return ", statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())"
	default:
		logging.WithFields("database_type", typ).Panic("position statement for type not implemented")
		return ""
	}
}

// Health implements eventstore.Pusher.
func (s *Storage) Health(ctx context.Context) error {
	return s.client.PingContext(ctx)
}
