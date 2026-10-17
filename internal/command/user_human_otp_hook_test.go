package command

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/action/otp"
	"github.com/zitadel/zitadel/internal/crypto"
)

func TestPublicGeneratorConfigFrom(t *testing.T) {
	got := publicGeneratorConfigFrom(&crypto.GeneratorConfig{
		Length:              8,
		Expiry:              time.Hour,
		IncludeLowerLetters: true,
		IncludeDigits:       true,
	})
	want := &otp.PublicGeneratorConfig{
		Length:              8,
		Expiry:              otp.Duration(time.Hour),
		IncludeLowerLetters: true,
		IncludeDigits:       true,
	}
	assert.Equal(t, want, got)

	assert.Nil(t, publicGeneratorConfigFrom(nil))
}

func TestApplyGenerationOverrides(t *testing.T) {
	base := &crypto.GeneratorConfig{
		Length:              8,
		Expiry:              time.Hour,
		IncludeLowerLetters: true,
		IncludeUpperLetters: true,
		IncludeDigits:       true,
		IncludeSymbols:      true,
	}
	tests := []struct {
		name   string
		gen    *otp.GenerationOverrides
		expiry *otp.Duration
		want   *crypto.GeneratorConfig
	}{
		{
			name: "no overrides returns clone",
			want: base,
		},
		{
			name: "length only",
			gen:  &otp.GenerationOverrides{Length: gu.Ptr(uint32(4))},
			want: &crypto.GeneratorConfig{
				Length:              4,
				Expiry:              time.Hour,
				IncludeLowerLetters: true,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      true,
			},
		},
		{
			name: "digits only voice SMS",
			gen: &otp.GenerationOverrides{
				Length:              gu.Ptr(uint32(4)),
				IncludeDigits:       gu.Ptr(true),
				IncludeLowerLetters: gu.Ptr(false),
				IncludeUpperLetters: gu.Ptr(false),
				IncludeSymbols:      gu.Ptr(false),
			},
			want: &crypto.GeneratorConfig{
				Length:              4,
				Expiry:              time.Hour,
				IncludeLowerLetters: false,
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
		},
		{
			name:   "expiry only",
			expiry: gu.Ptr(otp.Duration(5 * time.Minute)),
			want: &crypto.GeneratorConfig{
				Length:              8,
				Expiry:              5 * time.Minute,
				IncludeLowerLetters: true,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      true,
			},
		},
		{
			name:   "gen and expiry together",
			gen:    &otp.GenerationOverrides{Length: gu.Ptr(uint32(6))},
			expiry: gu.Ptr(otp.Duration(2 * time.Minute)),
			want: &crypto.GeneratorConfig{
				Length:              6,
				Expiry:              2 * time.Minute,
				IncludeLowerLetters: true,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyGenerationOverrides(base, tt.gen, tt.expiry)
			assert.Equal(t, tt.want, got)
			// Ensure base is not mutated.
			assert.Equal(t, uint(8), base.Length)
			assert.Equal(t, time.Hour, base.Expiry)
		})
	}
}

func TestOverrideExpiry(t *testing.T) {
	assert.Equal(t, time.Hour, overrideExpiry(time.Hour, nil))
	assert.Equal(t, 5*time.Minute, overrideExpiry(time.Hour, gu.Ptr(otp.Duration(5*time.Minute))))
}

func TestEncryptOverriddenOTPCode(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	code, err := encryptOverriddenOTPCode("A7F2B9", alg, 2*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "A7F2B9", code.Plain)
	assert.Equal(t, 2*time.Minute, code.Expiry)
	assert.NotNil(t, code.Crypted)
}

func TestValidateGenerationConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *crypto.GeneratorConfig
		wantErr bool
	}{
		{"length zero is tolerated", &crypto.GeneratorConfig{Length: 0}, false},
		{"digits only is valid", &crypto.GeneratorConfig{Length: 6, IncludeDigits: true}, false},
		{"lower+upper is valid", &crypto.GeneratorConfig{Length: 6, IncludeLowerLetters: true, IncludeUpperLetters: true}, false},
		{"length with no classes is rejected", &crypto.GeneratorConfig{Length: 6}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGenerationConfig(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
