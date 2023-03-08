package setup

import (
	"bytes"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type Config struct {
	Database        database.Config
	SystemDefaults  systemdefaults.SystemDefaults
	InternalAuthZ   authz.Config
	ExternalDomain  string
	ExternalPort    uint16
	ExternalSecure  bool
	Log             *logging.Config
	EncryptionKeys  *encryptionKeyConfig
	DefaultInstance command.InstanceSetup
	Machine         *id.Config
	Projections     projection.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
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

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}

type Steps struct {
	s1ProjectionTable    *ProjectionTable
	s2AssetsTable        *AssetTable
	FirstInstance        *FirstInstance
	s4EventstoreIndexes  *EventstoreIndexesNew
	s5LastFailed         *LastFailed
	s6OwnerRemoveColumns *OwnerRemoveColumns
	s7LogstoreTables     *LogstoreTables
	s8AuthTokens         *AuthTokenIndexes
	s9EventstoreIndexes2 *EventstoreIndexesNew
}

type encryptionKeyConfig struct {
	User *crypto.KeyConfig
	SMTP *crypto.KeyConfig
}

func MustNewSteps(v *viper.Viper) *Steps {
	v.AutomaticEnv()
	v.SetEnvPrefix("ZITADEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultSteps))
	logging.OnError(err).Fatal("unable to read setup steps")

	for _, file := range stepFiles {
		v.SetConfigFile(file)
		err := v.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read setup file")
	}

	steps := new(Steps)
	err = v.Unmarshal(steps,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
		)),
	)
	logging.OnError(err).Fatal("unable to read steps")
	return steps
}
