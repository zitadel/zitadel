//go:build !fips

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/crypto"
)

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
