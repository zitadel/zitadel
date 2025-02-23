package dialect

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/cmd/config"
	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/dialect/gosql"
	"github.com/zitadel/zitadel/backend/storage/database/dialect/postgres"
)

type Hook struct {
	Match  func(string) bool
	Decode func(name string, config any) (database.Connector, error)
	Name   string
	Field  configure.Updater
}

var hooks = make([]Hook, 0)

func init() {
	hooks = append(hooks,
		Hook{
			Match:  postgres.NameMatcher,
			Decode: postgres.DecodeConfig,
			Name:   postgres.Name,
			Field:  postgres.Field,
		},
		Hook{
			Match:  gosql.NameMatcher,
			Decode: gosql.DecodeConfig,
			Name:   gosql.Name,
			Field:  gosql.Field,
		},
	)
}

type Config struct {
	Dialects map[string]any `mapstructure:",remain"`

	connector database.Connector
}

// Fields implements [configure.StructUpdater].
func (c *Config) Fields() []configure.Updater {
	dialects := configure.OneOf{
		Description: "The database dialect Zitadel connects to",
		SubFields:   []configure.Updater{},
	}
	for _, hook := range hooks {
		if hook.Field == nil {
			panic("hook must configure its config fields")
		}
		dialects.SubFields = append(dialects.SubFields, &configure.Struct{
			FieldName:   hook.Name,
			Description: fmt.Sprintf("Configuration for %s", hook.Name),
			SubFields:   []configure.Updater{hook.Field},
		})
	}

	return []configure.Updater{
		dialects,
	}
}

// Name implements [configure.StructUpdater].
func (c *Config) Name() string {
	return "database"
}

func (c Config) Connect(ctx context.Context) (database.Pool, error) {
	if len(c.Dialects) != 1 {
		return nil, errors.New("Exactly one dialect must be configured")
	}

	return c.connector.Connect(ctx)
}

// Hooks implements [configure.Unmarshaller].
func (c Config) Hooks() []viper.DecoderConfigOption {
	return []viper.DecoderConfigOption{
		viper.DecodeHook(decodeHook),
	}
}

// var _ configure.StructUpdater = (*Config)(nil)

func (c Config) Configure(v *viper.Viper, currentVersion config.Version) Config {
	return c
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
