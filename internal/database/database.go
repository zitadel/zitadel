package database

import (
	"database/sql"

	"github.com/caos/zitadel/internal/errors"
)

func Connect(config Config) (*sql.DB, error) {
	client, err := sql.Open("postgres", config.String())
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(int(config.MaxOpenConns))
	client.SetConnMaxLifetime(config.MaxConnLifetime)
	client.SetConnMaxIdleTime(config.MaxConnIdleTime)

	if err := client.Ping(); err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "DATAB-0pIWD", "Errors.Database.Connection.Failed")
	}

	return client, nil
}
