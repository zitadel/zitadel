package initialise

import (
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/id"
)

type Config struct {
	Database database.Config
	Machine  *id.Config
	Log      *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(database.DecodeHook),
	)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}
