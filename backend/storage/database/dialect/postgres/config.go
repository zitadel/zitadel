package postgres

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/cmd/configure/bla"
	"github.com/zitadel/zitadel/backend/storage/database"
)

var (
	_    database.Connector = (*Config)(nil)
	Name                    = "postgres"

	Field = &configure.OneOf{
		Description: "Configuring postgres using one of the following options",
		SubFields: []configure.Updater{
			&configure.Field[string]{
				Description: "Connection string",
				Version:     semver.MustParse("v3"),
				Validate: func(s string) error {
					_, err := pgxpool.ParseConfig(s)
					return err
				},
			},
			&configure.Struct{
				Description: "Configuration for the connection",
				SubFields: []configure.Updater{
					&configure.Field[string]{
						FieldName:   "host",
						Value:       "localhost",
						Description: "The host to connect to",
						Version:     semver.MustParse("3"),
					},
					&configure.Field[uint32]{
						FieldName:   "port",
						Value:       5432,
						Description: "The port to connect to",
						Version:     semver.MustParse("3"),
					},
					&configure.Field[string]{
						FieldName:   "database",
						Value:       "zitadel",
						Description: "The database to connect to",
						Version:     semver.MustParse("3"),
					},
					&configure.Field[string]{
						FieldName:   "user",
						Description: "The user to connect as",
						Value:       "zitadel",
						Version:     semver.MustParse("3"),
					},
					&configure.Field[string]{
						FieldName:   "password",
						Description: "The password to connect with",
						Version:     semver.MustParse("3"),
						HideInput:   true,
					},
					&configure.OneOf{
						FieldName:   "sslMode",
						Description: "The SSL mode to use",
						SubFields: []configure.Updater{
							&configure.Constant[string]{
								Description: "Disable",
								Constant:    "disable",
								Version:     semver.MustParse("3"),
							},
							&configure.Constant[string]{
								Description: "Require",
								Constant:    "require",
								Version:     semver.MustParse("3"),
							},
							&configure.Constant[string]{
								Description: "Verify CA",
								Constant:    "verify-ca",
								Version:     semver.MustParse("3"),
							},
							&configure.Constant[string]{
								Description: "Verify Full",
								Constant:    "verify-full",
								Version:     semver.MustParse("3"),
							},
						},
					},
				},
			},
		},
	}
)

type Config struct{ pgxpool.Config }

// ConfigForIndex implements bla.OneOfField.
func (c Config) ConfigForIndex(i int) any {
	switch i {
	case 0:
		return new(string)
	case 1:
		return &c.Config
	}
	return nil
}

// Possibilities implements bla.OneOfField.
func (c Config) Possibilities() []string {
	return []string{"connection string", "fields"}
}

var _ bla.OneOfField = (*Config)(nil)

// Connect implements [database.Connector].
func (c *Config) Connect(ctx context.Context) (database.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, &c.Config)
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
		return &Config{Config: *config}, nil
	case map[string]any:
		return &Config{
			Config: pgxpool.Config{},
		}, nil
	}
	return nil, errors.New("invalid configuration")
}
