package start

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestReadConfigMaxFontSize(t *testing.T) {
	for _, tt := range []struct {
		name    string
		yaml    string
		env     string
		want    int64
		wantErr bool
	}{
		{name: "default", want: 524288},
		{name: "configured", yaml: "AssetStorage:\n  MaxFontSize: 33554432\n", want: 33554432},
		{name: "environment overrides yaml", yaml: "AssetStorage:\n  MaxFontSize: 1048576\n", env: "33554432", want: 33554432},
		{name: "zero", yaml: "AssetStorage:\n  MaxFontSize: 0\n", wantErr: true},
		{name: "negative", yaml: "AssetStorage:\n  MaxFontSize: -1\n", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			defaults, err := os.ReadFile("../defaults.yaml")
			require.NoError(t, err)
			v := viper.New()
			v.SetConfigType("yaml")
			require.NoError(t, v.ReadConfig(strings.NewReader(string(defaults))))
			if tt.yaml != "" {
				require.NoError(t, v.MergeConfig(strings.NewReader(tt.yaml)))
			}
			if tt.env != "" {
				t.Setenv("ZITADEL_ASSETSTORAGE_MAXFONTSIZE", tt.env)
				v.SetEnvPrefix("ZITADEL")
				v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
				v.AutomaticEnv()
			}
			config, err := readConfig(v)
			if tt.wantErr {
				require.ErrorContains(t, err, "AssetStorage.MaxFontSize")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, config.AssetStorage.MaxFontSize)
		})
	}
}

func TestReadConfigMaxFontSizeOmitted(t *testing.T) {
	config, err := readConfig(viper.New())
	require.NoError(t, err)
	require.Equal(t, int64(524288), config.AssetStorage.MaxFontSize)
}
