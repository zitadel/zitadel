package options

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/zitadel/logging"
)

func InitViper() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZITADEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return fmt.Errorf("unable to read default config: %w", err)
	}
	return nil
}

func MergeToViper(configFiles ...string) {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}
}
