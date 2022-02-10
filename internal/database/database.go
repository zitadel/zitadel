package database

import (
	"database/sql"
)

func Connect(config Config) (*sql.DB, error) {
	client, err := sql.Open("postgres", config.String())
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(int(config.MaxOpenConns))
	client.SetConnMaxLifetime(config.MaxConnLifetime.Duration)
	client.SetConnMaxIdleTime(config.MaxConnIdleTime.Duration)

	return client, nil
}
