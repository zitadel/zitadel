package postgres

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/storage/database"
)

var _ database.Connector = (*Config)(nil)

type Config struct{ *pgxpool.Config }

// Connect implements [database.Connector].
func (c *Config) Connect(ctx context.Context) (database.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, c.Config)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &pgxPool{pool}, nil
}

func NameMatcher(name string) bool {
	return slices.Contains([]string{"postgres", "pg"}, strings.ToLower(name))
}

func DecodeConfig(_ string, config any) (database.Connector, error) {
	switch c := config.(type) {
	case string:
		config, err := pgxpool.ParseConfig(c)
		if err != nil {
			return nil, err
		}
		return &Config{config}, nil
	case map[string]any:
		return nil, errors.New("map configuration not implemented")
	}
	return nil, errors.New("invalid configuration")
}
