package dialect

import (
	"context"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/dialect/gosql"
	"github.com/zitadel/zitadel/backend/storage/database/dialect/postgres"
)

type Hook struct {
	Match  func(string) bool
	Decode func(name string, config any) (database.Connector, error)
}

var hooks = make([]Hook, 0)

func init() {
	hooks = append(hooks,
		Hook{
			Match:  postgres.NameMatcher,
			Decode: postgres.DecodeConfig,
		},
		Hook{
			Match:  gosql.NameMatcher,
			Decode: gosql.DecodeConfig,
		},
	)
}

type Config struct {
	Dialects map[string]any `mapstructure:",remain"`

	connector database.Connector
}

// Hooks implements [configure.Unmarshaller].
func (c Config) Hooks() []viper.DecoderConfigOption {
	return []viper.DecoderConfigOption{
		viper.DecodeHook(decodeHook),
	}
}

func (c Config) Connect(ctx context.Context) (database.Pool, error) {
	return c.connector.Connect(ctx)
}

func (c *Config) decodeDialect() error {
	for _, hook := range hooks {
		for name, config := range c.Dialects {
			if !hook.Match(name) {
				continue
			}

			connector, err := hook.Decode(name, config)
			if err != nil {
				return err
			}

			c.connector = connector
			return nil
		}
	}
	return errors.New("no dialect found")
}

func decodeHook(from, to reflect.Value) (_ interface{}, err error) {
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
