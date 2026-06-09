package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

func loadBaseDefaultConfig(t *testing.T, v *viper.Viper) {
	t.Helper()
	require.NoError(t, v.ReadConfig(bytes.NewBuffer(defaultConfig)))
}

func TestLoadDefaultConfig_NonFIPSPasswordHasher(t *testing.T) {
	v := testViper(t)
	require.NoError(t, loadDefaultConfigInto(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.PasswordHasher")
	assert.Equal(t, crypto.HashNameBcrypt, cfg.Hasher.Algorithm)
}

func TestLoadDefaultConfig_NonFIPSSecretHasher(t *testing.T) {
	v := testViper(t)
	require.NoError(t, loadDefaultConfigInto(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.SecretHasher")
	assert.Equal(t, crypto.HashNameBcrypt, cfg.Hasher.Algorithm)
}

func TestLoadDefaultConfig_FIPSPasswordHasher(t *testing.T) {
	v := testViper(t)
	loadBaseDefaultConfig(t, v)
	require.NoError(t, applyFipsDefaultOverlay(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.PasswordHasher")
	assert.Equal(t, crypto.HashNamePBKDF2, cfg.Hasher.Algorithm)
	assert.Equal(t, 290000, v.GetInt("SystemDefaults.PasswordHasher.Hasher.Rounds"))
	assert.Equal(t, "sha256", v.GetString("SystemDefaults.PasswordHasher.Hasher.Hash"))
	assert.Empty(t, cfg.Verifiers)
}

func TestLoadDefaultConfig_FIPSSecretHasher(t *testing.T) {
	v := testViper(t)
	loadBaseDefaultConfig(t, v)
	require.NoError(t, applyFipsDefaultOverlay(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.SecretHasher")
	assert.Equal(t, crypto.HashNamePBKDF2, cfg.Hasher.Algorithm)
	assert.Equal(t, 290000, v.GetInt("SystemDefaults.SecretHasher.Hasher.Rounds"))
	assert.Equal(t, "sha256", v.GetString("SystemDefaults.SecretHasher.Hasher.Hash"))
	assert.Empty(t, cfg.Verifiers)
}

func TestLoadDefaultConfig_FIPSVerifierKeysEmpty(t *testing.T) {
	v := testViper(t)
	loadBaseDefaultConfig(t, v)
	require.NoError(t, applyFipsDefaultOverlay(v))
	assert.Empty(t, v.GetStringSlice("SystemDefaults.PasswordHasher.Verifiers"))
	assert.Empty(t, v.GetStringSlice("SystemDefaults.SecretHasher.Verifiers"))
}
