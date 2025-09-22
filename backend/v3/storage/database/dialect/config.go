package dialect

import (
	"context"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
)

type Hook struct {
	Match       func(string) bool
	Decode      func(config any) (database.Connector, error)
	Name        string
	Constructor func() database.Connector
}

var hooks = []Hook{
	{
		Match:       postgres.NameMatcher,
		Decode:      postgres.DecodeConfig,
		Name:        postgres.Name,
		Constructor: func() database.Connector { return new(postgres.Config) },
	},
	// {
	// 	Match:       gosql.NameMatcher,
	// 	Decode:      gosql.DecodeConfig,
	// 	Name:        gosql.Name,
	// 	Constructor: func() database.Connector { return new(gosql.Config) },
	// },
}

type Config struct {
	Dialects map[string]any `mapstructure:",remain" yaml:",inline"`

	connector database.Connector
}

func (c Config) Connect(ctx context.Context) (database.Pool, error) {
	if len(c.Dialects) != 1 {
		return nil, errors.New("exactly one dialect must be configured")
	}

	return c.connector.Connect(ctx)
}

// Hooks implements [configure.Unmarshaller].
func (c Config) Hooks() []viper.DecoderConfigOption {
	return []viper.DecoderConfigOption{
		viper.DecodeHook(decodeHook),
	}
}

func decodeHook(from, to reflect.Value) (_ any, err error) {
	if to.Type() != reflect.TypeOf(Config{}) {
		return from.Interface(), nil
	}

	config := new(Config)
	if err = mapstructure.Decode(from.Interface(), config); err != nil {
		return nil, err
	}

	if err = config.decodeDialect(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) decodeDialect() error {
	for _, hook := range hooks {
		for name, config := range c.Dialects {
			if !hook.Match(name) {
				continue
			}

			connector, err := hook.Decode(config)
			if err != nil {
				return err
			}

			c.connector = connector
			return nil
		}
	}
	return errors.New("no dialect found")
}
