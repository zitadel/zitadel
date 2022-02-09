package config

import (
	"database/sql"

	"github.com/caos/zitadel/internal/errors"
)

var client *sql.DB

type Config struct{}

func Connect() (*sql.DB, error) {
	//TODO: viper read into Config

	return nil, errors.ThrowUnimplemented(nil, "CONFI-8bvVL", "connect is unimplemented")
}
