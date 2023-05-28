package eventstore

import (
	"time"

	"github.com/zitadel/zitadel/internal/database"
	z_sql "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

type Config struct {
	PushTimeout              time.Duration
	Client                   *database.DB
	AllowOrderByCreationDate bool

	pusher  Pusher
	querier Querier
}

func TestConfig(querier Querier, pusher Pusher) *Config {
	return &Config{pusher: pusher, querier: querier}
}

func Start(config *Config) (*Eventstore, error) {
	config.querier = z_sql.NewCRDB(config.Client, config.AllowOrderByCreationDate)
	config.pusher = eventstore.NewEventstore(config.Client)
	return NewEventstore(config), nil
}
