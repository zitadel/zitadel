package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Version Version
}

func (c Config) Hooks() []viper.DecoderConfigOption {
	return []viper.DecoderConfigOption{
		viper.DecodeHook(decodeVersion),
	}
}

func (c Config) CurrentVersion() Version {
	return c.Version
}

var Path string

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if Path != "" {
		// Use config file from the flag.
		viper.SetConfigFile(Path)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".zitadel" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".zitadel")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("ZITADEL")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

type Version uint8

const (
	VersionUnknown Version = iota
	V3
)

func decodeVersion(from, to reflect.Value) (_ interface{}, err error) {
	if to.Type() != reflect.TypeOf(Version(0)) {
		return from.Interface(), nil
	}

	switch from.Interface().(string) {
	case "":
		return VersionUnknown, nil
	case "v3":
		return V3, nil

	}

	return nil, fmt.Errorf("unsupported version: %v", from.Interface())
}
