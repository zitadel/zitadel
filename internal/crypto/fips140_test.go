package crypto

import (
	"crypto/fips140"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsNonFIPSHasherAlgorithm(t *testing.T) {
	tests := []struct {
		alg  HashName
		want bool
	}{
		{HashNameBcrypt, true},
		{HashNameScrypt, true},
		{HashNameArgon2i, true},
		{HashNameArgon2id, true},
		{HashNamePBKDF2, false},
		{HashNameSha2, false},
		{HashNameArgon2, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.alg), func(t *testing.T) {
			assert.Equal(t, tt.want, isNonFIPSHasherAlgorithm(tt.alg))
		})
	}
}

func TestIsNonFIPSVerifier(t *testing.T) {
	tests := []struct {
		name HashName
		want bool
	}{
		{HashNameArgon2, true},
		{HashNameBcrypt, true},
		{HashNameScrypt, true},
		{HashNameMd5, true},
		{HashNameMd5Plain, true},
		{HashNameMd5Salted, true},
		{HashNamePHPass, true},
		{HashNameDrupal7, true},
		{HashNamePBKDF2, false},
		{HashNameSha2, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			assert.Equal(t, tt.want, isNonFIPSVerifier(tt.name))
		})
	}
}

func TestIsNonFIPSPBKDF2HashMode(t *testing.T) {
	tests := []struct {
		mode HashMode
		want bool
	}{
		{HashModeSHA1, true},
		{HashModeSHA224, true},
		{HashModeSHA256, false},
		{HashModeSHA384, false},
		{HashModeSHA512, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			assert.Equal(t, tt.want, isNonFIPSPBKDF2HashMode(tt.mode))
		})
	}
}

func TestNonFIPSVerifiersConfigured(t *testing.T) {
	got := nonFIPSVerifiersConfigured([]HashName{HashNamePBKDF2, HashNameBcrypt, HashNameMd5})
	assert.Equal(t, []HashName{HashNameBcrypt, HashNameMd5}, got)
	assert.Nil(t, nonFIPSVerifiersConfigured([]HashName{HashNamePBKDF2, HashNameSha2}))
}

func TestValidateFIPSPBKDF2Hasher(t *testing.T) {
	tests := []struct {
		name    string
		hasher  HasherConfig
		wantErr string
	}{
		{
			name: "rounds below minimum",
			hasher: HasherConfig{
				Algorithm: HashNamePBKDF2,
				Params: map[string]any{
					"Rounds": 999,
					"Hash":   HashModeSHA256,
				},
			},
			wantErr: "iteration count 999 is below the FIPS minimum",
		},
		{
			name: "rounds at minimum",
			hasher: HasherConfig{
				Algorithm: HashNamePBKDF2,
				Params: map[string]any{
					"Rounds": 1000,
					"Hash":   HashModeSHA256,
				},
			},
		},
		{
			name: "sha1 hash mode",
			hasher: HasherConfig{
				Algorithm: HashNamePBKDF2,
				Params: map[string]any{
					"Rounds": 10000,
					"Hash":   HashModeSHA1,
				},
			},
			wantErr: "hash mode \"sha1\" is not FIPS 140-3 compliant",
		},
		{
			name: "sha224 hash mode",
			hasher: HasherConfig{
				Algorithm: HashNamePBKDF2,
				Params: map[string]any{
					"Rounds": 10000,
					"Hash":   HashModeSHA224,
				},
			},
			wantErr: "hash mode \"sha224\" is not FIPS 140-3 compliant",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFIPSPBKDF2Hasher(tt.hasher)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestHashConfig_validateFIPS140(t *testing.T) {
	if !fips140.Enabled() {
		t.Skip("FIPS mode not enabled; run with GODEBUG=fips140=on on a GOFIPS140 build")
	}

	t.Run("bcrypt hasher fails", func(t *testing.T) {
		cfg := &HashConfig{
			Hasher: HasherConfig{
				Algorithm: HashNameBcrypt,
				Params:    map[string]any{"Cost": 12},
			},
			Limits: HashLimitsConfig{
				Bcrypt: BcryptLimitsConfig{MinCost: 10, MaxCost: 16},
			},
		}
		err := cfg.validateFIPS140()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "uncertified cryptographic state")
		assert.Contains(t, err.Error(), "bcrypt")
	})

	t.Run("pbkdf2 compliant passes", func(t *testing.T) {
		cfg := &HashConfig{
			Hasher: HasherConfig{
				Algorithm: HashNamePBKDF2,
				Params: map[string]any{
					"Rounds": 290000,
					"Hash":   HashModeSHA256,
				},
			},
			Limits: HashLimitsConfig{
				PBKDF2: PBKDF2LimitsConfig{MinRounds: 1000, MaxRounds: 10000000},
			},
		}
		require.NoError(t, cfg.validateFIPS140())
	})
}

func TestHashConfig_NewHasher_FIPSBcryptFails(t *testing.T) {
	if !fips140.Enabled() {
		t.Skip("FIPS mode not enabled; run with GODEBUG=fips140=on on a GOFIPS140 build")
	}

	cfg := &HashConfig{
		Hasher: HasherConfig{
			Algorithm: HashNameBcrypt,
			Params:    map[string]any{"Cost": 12},
		},
		Limits: HashLimitsConfig{
			Bcrypt: BcryptLimitsConfig{MinCost: 10, MaxCost: 16},
		},
	}
	_, err := cfg.NewHasher()
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "uncertified cryptographic state"))
}
