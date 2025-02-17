package configure

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Unmarshaller interface {
	Hooks() []viper.DecoderConfigOption
}

func ReadConfigPreRun[C Unmarshaller](v *viper.Viper, config *C) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := v.Unmarshal(config, (*config).Hooks()...); err != nil {
			panic(err)
		}
	}
}

func ReadConfig[C Unmarshaller](v *viper.Viper) (*C, error) {
	var config C
	if err := v.Unmarshal(&config, config.Hooks()...); err != nil {
		return nil, err
	}
	return &config, nil
}
