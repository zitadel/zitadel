package eventstore

import (
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
	z_sql "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
)

type Config struct {
	PushTimeout time.Duration
	Client      *sql.DB

	repo repository.Repository
}

func Start(config *Config) (*Eventstore, error) {
	config.repo = z_sql.NewCRDB(config.Client)
	return NewEventstore(config), nil
}
