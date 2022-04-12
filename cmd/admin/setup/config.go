package setup

import (
	"bytes"

	"github.com/caos/logging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/hook"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/database"
)

type Config struct {
	Database       database.Config
	SystemDefaults systemdefaults.SystemDefaults
	InternalAuthZ  authz.Config
	ExternalPort   uint16
	ExternalDomain string
	ExternalSecure bool
	Log            *logging.Config
	EncryptionKeys *encryptionKeyConfig
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}

type Steps struct {
	s1ProjectionTable *ProjectionTable
	s2AssetsTable     *AssetTable
	S3DefaultInstance *DefaultInstance
}

type encryptionKeyConfig struct {
	User *crypto.KeyConfig
}

func MustNewSteps(v *viper.Viper) *Steps {
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultSteps))
	logging.OnError(err).Fatal("unable to read setup steps")

	steps := new(Steps)
	err = v.Unmarshal(steps,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)),
	)
	logging.OnError(err).Fatal("unable to read steps")
	return steps
}
