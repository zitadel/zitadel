package dialect

import (
	"context"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/dialect/postgres"
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

// Configure implements [configure.Configurer].
// func (c *Config) Configure() (any, error) {
// 	possibilities := make([]string, len(hooks))
// 	var cursor int
// 	for i, hook := range hooks {
// 		if _, ok := c.Dialects[hook.Name]; ok {
// 			cursor = i
// 		}
// 		possibilities[i] = hook.Name
// 	}

// 	prompt := promptui.Select{
// 		Label:     "Select a dialect",
// 		Items:     possibilities,
// 		CursorPos: cursor,
// 	}
// 	i, _, err := prompt.Run()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var config bla4.Configurer

// 	if dialect, ok := c.Dialects[hooks[i].Name]; ok {
// 		config, err = hooks[i].Decode(dialect)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		clear(c.Dialects)
// 		config = hooks[i].Constructor()
// 	}
// 	if c.Dialects == nil {
// 		c.Dialects = make(map[string]any)
// 	}
// 	c.Dialects[hooks[i].Name], err = config.Configure()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c, nil
// }

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
