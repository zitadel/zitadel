//go:build fips

package cmd

import (
	"bytes"
	_ "embed"

	"github.com/spf13/viper"
)

//go:embed defaults_fips.yaml
var defaultFipsConfig []byte

func mergeFipsDefaultConfig(v *viper.Viper) error {
	return v.MergeConfig(bytes.NewBuffer(defaultFipsConfig))
}
