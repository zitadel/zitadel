package configure

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Unmarshaller interface {
	Hooks() []viper.DecoderConfigOption
}

func ReadConfigPreRun[C Unmarshaller](v *viper.Viper, config C) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := v.Unmarshal(config, config.Hooks()...); err != nil {
			panic(err)
		}
	}
}

func ReadConfig[C Unmarshaller](v *viper.Viper) (config C, err error) {
	if err := v.Unmarshal(&config, config.Hooks()...); err != nil {
		return config, err
	}
	return config, nil
}
