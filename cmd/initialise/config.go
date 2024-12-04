package initialise

import (
	"log/slog"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
)

type Config struct {
	Database database.Config
	Machine  *id.Config
	Log      *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			database.DecodeHook,
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	logging.OnError(err).Fatal("unable to read config")

	config.Log.Formatter.Data = map[string]interface{}{
		"service": "zitadel",
		"version": build.Version(),
	}

	slog.SetDefault(config.Log.Slog())

	id.Configure(config.Machine)

	return config
}
