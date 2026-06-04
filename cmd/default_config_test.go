package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/crypto"
)

func testViper(t *testing.T) *viper.Viper {
	t.Helper()
	v := viper.New()
	v.SetConfigType("yaml")
	return v
}

func loadAndUnmarshalHashConfig(t *testing.T, v *viper.Viper, key string) crypto.HashConfig {
	t.Helper()
	var cfg crypto.HashConfig
	require.NoError(t, v.UnmarshalKey(key, &cfg))
	return cfg
}
