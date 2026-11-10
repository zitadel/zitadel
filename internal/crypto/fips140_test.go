package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withFIPS140Enabled(t *testing.T) {
	t.Helper()
	prev := fips140Mode
	fips140Mode = func() bool { return true }
	t.Cleanup(func() { fips140Mode = prev })
}

func TestHashName_IsFIPSCompliant(t *testing.T) {
	tests := []struct {
		name HashName
		want bool
	}{
		{HashNamePBKDF2, true},
		{HashNameArgon2, false},
		{HashNameArgon2i, false},
		{HashNameArgon2id, false},
		{HashNameBcrypt, false},
		{HashNameMd5, false},
		{HashNameMd5Plain, false},
		{HashNameMd5Salted, false},
		{HashNamePHPass, false},
		{HashNameSha2, false},
		{HashNameScrypt, false},
		{HashNameDrupal7, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.name.IsFIPSCompliant())
		})
	}
}

func TestHashMode_IsFIPSCompliant(t *testing.T) {
	tests := []struct {
		mode HashMode
		want bool
	}{
		{HashModeSHA256, true},
		{HashModeSHA384, true},
		{HashModeSHA512, true},
		{HashModeSHA1, false},
		{HashModeSHA224, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mode.IsFIPSCompliant())
		})
	}
}

func TestNonFIPSVerifiersConfigured(t *testing.T) {
	got := nonFIPSVerifiersConfigured([]HashName{HashNamePBKDF2, HashNameBcrypt, HashNameMd5})
	assert.Equal(t, []HashName{HashNameBcrypt, HashNameMd5}, got)
	assert.Equal(t, []HashName{HashNameSha2}, nonFIPSVerifiersConfigured([]HashName{HashNamePBKDF2, HashNameSha2}))
	assert.Nil(t, nonFIPSVerifiersConfigured([]HashName{HashNamePBKDF2}))
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
	withFIPS140Enabled(t)

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
	withFIPS140Enabled(t)

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
	assert.Contains(t, err.Error(), "uncertified cryptographic state")
}
