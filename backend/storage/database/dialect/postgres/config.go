package postgres

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/backend/cmd/configure/bla4"
	"github.com/zitadel/zitadel/backend/storage/database"
)

var (
	_    database.Connector = (*Config)(nil)
	Name                    = "postgres"
)

type Config struct {
	config *pgxpool.Config

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

// FinishAllowed implements [bla4.Iterator].
func (c *Config) FinishAllowed() bool {
	// Option can be skipped
	return len(c.configuredFields) < 2
}

// NextField implements [bla4.Iterator].
func (c *Config) NextField() string {
	if c.configuredFields == nil {
		c.configuredFields = []string{"Host", "Port", "Database", "MaxOpenConns", "MaxIdleConns", "MaxConnLifetime", "MaxConnIdleTime", "User", "Options"}
	}
	if len(c.configuredFields) == 0 {
		return ""
	}
	field := c.configuredFields[0]
	c.configuredFields = c.configuredFields[1:]
	return field
}

// Configure implements [bla4.Configurer].
func (c *Config) Configure() (value any, err error) {
	typeSelect := promptui.Select{
		Label: "Configure the database connection",
		Items: []string{"connection string", "fields"},
	}
	i, _, err := typeSelect.Run()
	if err != nil {
		return nil, err
	}
	if i > 0 {
		return nil, nil
	}

	if c.config == nil {
		c.config, _ = pgxpool.ParseConfig("host=localhost user=zitadel password= dbname=zitadel sslmode=disable")
	}

	prompt := promptui.Prompt{
		Label:     "Connection string",
		Default:   c.config.ConnString(),
		AllowEdit: c.config.ConnString() != "",
		Validate: func(input string) error {
			_, err := pgxpool.ParseConfig(input)
			return err
		},
	}

	return prompt.Run()
}

var _ bla4.Iterator = (*Config)(nil)

// Connect implements [database.Connector].
func (c *Config) Connect(ctx context.Context) (database.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, c.config)
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
		return &Config{config: config}, nil
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
			config: &pgxpool.Config{},
		}, nil
	}
	return nil, errors.New("invalid configuration")
}
