package ready

import (
	"log/slog"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/network"
)

type Config struct {
	Log  *logging.Config
	Port uint16
	TLS  network.TLS
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			hook.EnumHookFunc(internal_authz.MemberTypeString),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")

	config.Log.Formatter.Data = map[string]interface{}{
		"service": "zitadel",
		"version": build.Version(),
	}

	slog.SetDefault(config.Log.Slog())

	return config
}
