package gosql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/storage/database"
)

var (
	_ database.Connector = (*Config)(nil)

	Name  = "gosql"
	Field = &configure.Field[string]{
		Description: "Connection string",
		Version:     semver.MustParse("v3"),
		Validate: func(s string) error {
			_, err := pgxpool.ParseConfig(s)
			return err
		},
	}
)

type Config struct {
	db *sql.DB
}

// Connect implements [database.Connector].
func (c *Config) Connect(ctx context.Context) (database.Pool, error) {
	if err := c.db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &sqlPool{c.db}, nil
}

func NameMatcher(name string) bool {
	name = strings.ToLower(name)
	for _, driver := range sql.Drivers() {
		if driver == name {
			return true
		}
	}
	return false
}

func DecodeConfig(name string, config any) (database.Connector, error) {
	switch c := config.(type) {
	case string:
		db, err := sql.Open(name, c)
		if err != nil {
			return nil, err
		}
		return &Config{db}, nil
	case map[string]any:
		return nil, errors.New("map configuration not implemented")
	}
	return nil, errors.New("invalid configuration")
}
