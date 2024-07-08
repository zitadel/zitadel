package initialise

import (
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/sonyflake"
)

type Config struct {
	Database    database.Config
	IDGenerator id_generator.GeneratorType
	Machine     *sonyflake.Config
	Log         *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(database.DecodeHook),
	)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id_generator.SetGeneratorWithConfig(config.IDGenerator, config.Machine)

	return config
}
