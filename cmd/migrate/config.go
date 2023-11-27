package migrate

import (
	"bytes"
	_ "embed"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
)

type EventsConfig struct {
	Source      database.Config
	Destination database.Config

	Log     *logging.Config
	Machine *id.Config
}

var (
	//go:embed defaults.yaml
	defaultConfig []byte
	configPaths   []string
)

func mustNewEventsConfig(v *viper.Viper) *EventsConfig {
	config := new(EventsConfig)
	mustNewConfig(v, config)

	err := config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}

func mustNewProjectionsConfig(v *viper.Viper) *ProjectionsConfig {
	config := new(ProjectionsConfig)
	mustNewConfig(v, config)

	err := config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}

func mustNewConfig(v *viper.Viper, config any) {
	v.AutomaticEnv()
	v.SetEnvPrefix("ZITADEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read setup steps")

	for _, file := range configPaths {
		v.SetConfigFile(file)
		err := v.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}

	err = v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			database.DecodeHook,
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")
}
