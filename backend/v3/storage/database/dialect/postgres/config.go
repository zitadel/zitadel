package postgres

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var (
	_    database.Connector = (*Config)(nil)
	Name                    = "postgres"
)

type Config struct {
	*pgxpool.Config

	// Host               string
	// Port               int32
	// Database           string
	// MaxOpenConns       uint32
	// MaxIdleConns       uint32
	// MaxConnLifetime    time.Duration
	// MaxConnIdleTime    time.Duration
	// User               User
	// // Additional options to be appended as options=<Options>
	// // The value will be taken as is. Multiple options are space separated.
	// Options string

	configuredFields []string
}

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

func DecodeConfig(input any) (database.Connector, error) {
	switch c := input.(type) {
	case string:
		config, err := pgxpool.ParseConfig(c)
		if err != nil {
			return nil, err
		}
		return &Config{Config: config}, nil
	case map[string]any:
		connector := new(Config)
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
			WeaklyTypedInput: true,
			Result:           connector,
		})
		if err != nil {
			return nil, err
		}
		if err = decoder.Decode(c); err != nil {
			return nil, err
		}
		return &Config{
			Config: &pgxpool.Config{},
		}, nil
	}
	return nil, errors.New("invalid configuration")
}
