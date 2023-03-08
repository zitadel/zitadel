package eventstore

import (
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	z_sql "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
)

type Config struct {
	PushTimeout time.Duration
	Client      *database.DB

	repo repository.Repository
}

func TestConfig(repo repository.Repository) *Config {
	return &Config{repo: repo}
}

func Start(config *Config) (*Eventstore, error) {
	config.repo = z_sql.NewCRDB(config.Client)
	return NewEventstore(config), nil
}
