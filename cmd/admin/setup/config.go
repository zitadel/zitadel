package setup

import (
	"bytes"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/hook"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/viper"
)

type Config struct {
	Database       database.Config
	SystemDefaults systemdefaults.SystemDefaults
	InternalAuthZ  authz.Config
	ExternalPort   uint16
	ExternalDomain string
	ExternalSecure bool
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	return config
}

type Steps struct {
	S1DefaultInstance *DefaultInstance
}

func MustNewSteps(v *viper.Viper) *Steps {
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultSteps))
	logging.OnError(err).Fatal("unable to read setup steps")

	steps := new(Steps)
	err = v.Unmarshal(steps,
		viper.DecodeHook(hook.Base64ToBytesHookFunc()),
		viper.DecodeHook(hook.TagToLanguageHookFunc()),
	)
	logging.OnError(err).Fatal("unable to read steps")
	return steps
}
