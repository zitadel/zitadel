//go:build fips

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/crypto"
)

func TestLoadDefaultConfig_FIPSPasswordHasher(t *testing.T) {
	v := testViper(t)
	require.NoError(t, loadDefaultConfigInto(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.PasswordHasher")
	assert.Equal(t, crypto.HashNamePBKDF2, cfg.Hasher.Algorithm)
	assert.Equal(t, 290000, v.GetInt("SystemDefaults.PasswordHasher.Hasher.Rounds"))
	assert.Equal(t, "sha256", v.GetString("SystemDefaults.PasswordHasher.Hasher.Hash"))
	assert.Empty(t, cfg.Verifiers)
}

func TestLoadDefaultConfig_FIPSSecretHasher(t *testing.T) {
	v := testViper(t)
	require.NoError(t, loadDefaultConfigInto(v))
	cfg := loadAndUnmarshalHashConfig(t, v, "SystemDefaults.SecretHasher")
	assert.Equal(t, crypto.HashNamePBKDF2, cfg.Hasher.Algorithm)
	assert.Equal(t, 290000, v.GetInt("SystemDefaults.SecretHasher.Hasher.Rounds"))
	assert.Equal(t, "sha256", v.GetString("SystemDefaults.SecretHasher.Hasher.Hash"))
	assert.Empty(t, cfg.Verifiers)
}

func TestLoadDefaultConfig_FIPSVerifierKeysEmpty(t *testing.T) {
	v := testViper(t)
	require.NoError(t, loadDefaultConfigInto(v))
	assert.Empty(t, v.GetStringSlice("SystemDefaults.PasswordHasher.Verifiers"))
	assert.Empty(t, v.GetStringSlice("SystemDefaults.SecretHasher.Verifiers"))
}
