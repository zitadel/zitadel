package firebasescrypt

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap/verifier"
)

const (
	testPassword  = "user1password"
	testSalt      = "42xEC+ixf3L2lw=="
	testSaltSep   = "Bw=="
	testSignerKey = "jxspr8Ki0RYycVU8zykbdLGjFQ3McFUH0uiiTvC8pVMXAn210wjLNmdZJzxUECKbm0QsEmYUSDzZvpjeJ9WmXA=="
	testHashA     = "lSrfV15cpx95/sZS2W9c9Kp6i/LVgQNDNC/qzrCnh1SAyZvqmZqAjTdn3aoItz+VHjoZilo78198JAdRuid5lQ=="
	testHashB     = "NohHtJ2FmkEqzefZ0IHnOlqTzFN8n6vXXd0EW2OIGuUh4rcwvVh7XVRDVs5yOrueoudRPIoimS6jIwf4pMW3rg=="
	testHashC     = "EQ0tlNdKv5ZBDYN7ofOYlhhTEe5Tu9bueDmtDUhLPM+9dcS1u1ALgMcT6N1+U9GsDlVdd3FE2dmz37fU5Y8A9Q=="
	defaultParams = "ln=14,r=8"
)

func encode(params, salt, hash, saltSep, signerKey string) string {
	return Prefix + strings.Join([]string{params, salt, hash, saltSep, signerKey}, "$")
}

var (
	// Sample from the firebase/scrypt README.
	vectorA = encode(defaultParams, testSalt, testHashA, testSaltSep, testSignerKey)
	vectorB = encode(defaultParams, testSalt, testHashB, "", testSignerKey)
	vectorC = encode("ln=10,r=4", testSalt, testHashC, testSaltSep, testSignerKey)
	// Values from the Ory Kratos tests.
	vectorD = "$firebasescrypt$ln=14,r=8$sPtDhWcd1MfdAw==$xbSou7FOl6mChCyzpCPIQ7tku7nsQMTFtyOZSXXd7tjBa4NtimOx7v42Gv2SfzPQu1oxM2/k4SsbOu73wlKe1A==$Bw==$YE0dO4bwD4JnJafh6lZZfkp1MtKzuKAXQcDCJNJNyeCHairWHKENOkbh3dzwaCdizzOspwr/FITUVlnOAwPKyw=="
)

// Ory Kratos format of vector A.
var kratosVectorA = "$firescrypt$" + strings.Join([]string{"ln=14,r=8,p=1", testSalt, testHashA, testSaltSep, testSignerKey}, "$")

func urlSafe(s string) string {
	return strings.NewReplacer("+", "-", "/", "_").Replace(s)
}

func noPadding(s string) string {
	return strings.TrimRight(s, "=")
}

func TestVerifier_Verify(t *testing.T) {
	tests := []struct {
		name     string
		encoded  string
		password string
		want     verifier.Result
	}{
		{
			name:     "vector A",
			encoded:  vectorA,
			password: testPassword,
			want:     verifier.OK,
		},
		{
			name:     "vector B, empty salt separator",
			encoded:  vectorB,
			password: testPassword,
			want:     verifier.OK,
		},
		{
			name:     "vector C, non default params",
			encoded:  vectorC,
			password: testPassword,
			want:     verifier.OK,
		},
		{
			name:     "vector D, kratos",
			encoded:  vectorD,
			password: "8x4WjoDbSxJZdR",
			want:     verifier.OK,
		},
		{
			name:     "wrong password",
			encoded:  vectorA,
			password: "wrong",
			want:     verifier.Fail,
		},
		{
			name:     "wrong hash",
			encoded:  encode(defaultParams, testSalt, "m"+testHashA[1:], testSaltSep, testSignerKey),
			password: testPassword,
			want:     verifier.Fail,
		},
		{
			name:     "url safe base64",
			encoded:  encode(defaultParams, urlSafe(testSalt), urlSafe(testHashA), urlSafe(testSaltSep), urlSafe(testSignerKey)),
			password: testPassword,
			want:     verifier.OK,
		},
		{
			name:     "no padding",
			encoded:  encode(defaultParams, noPadding(testSalt), noPadding(testHashA), noPadding(testSaltSep), noPadding(testSignerKey)),
			password: testPassword,
			want:     verifier.OK,
		},
		{
			name:     "skip bcrypt",
			encoded:  "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
			password: testPassword,
			want:     verifier.Skip,
		},
		{
			name:     "skip kratos format",
			encoded:  kratosVectorA,
			password: testPassword,
			want:     verifier.Skip,
		},
	}
	v := NewVerifier(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Verify(tt.encoded, tt.password)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVerifier_Validate(t *testing.T) {
	tests := []struct {
		name      string
		opts      *ValidationOpts
		encoded   string
		want      verifier.Result
		wantParam string
	}{
		{
			name:    "vector A",
			encoded: vectorA,
			want:    verifier.OK,
		},
		{
			name: "vector C, wide opts",
			opts: &ValidationOpts{
				MinLN: 10,
				MaxLN: 14,
				MinR:  4,
				MaxR:  8,
			},
			encoded: vectorC,
			want:    verifier.OK,
		},
		{
			name:      "LN below min",
			encoded:   encode("ln=13,r=8", testSalt, testHashA, testSaltSep, testSignerKey),
			want:      verifier.Fail,
			wantParam: "LN",
		},
		{
			name:      "LN above max",
			encoded:   encode("ln=15,r=8", testSalt, testHashA, testSaltSep, testSignerKey),
			want:      verifier.Fail,
			wantParam: "LN",
		},
		{
			name:      "R below min",
			encoded:   encode("ln=14,r=7", testSalt, testHashA, testSaltSep, testSignerKey),
			want:      verifier.Fail,
			wantParam: "R",
		},
		{
			name:      "R above max",
			encoded:   encode("ln=14,r=9", testSalt, testHashA, testSaltSep, testSignerKey),
			want:      verifier.Fail,
			wantParam: "R",
		},
		{
			name:    "skip kratos format",
			encoded: kratosVectorA,
			want:    verifier.Skip,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVerifier(tt.opts).Validate(tt.encoded)
			assert.Equal(t, tt.want, got)
			if tt.wantParam == "" {
				require.NoError(t, err)
				return
			}
			var bounds *verifier.BoundsError
			require.True(t, errors.As(err, &bounds))
			assert.Equal(t, Identifier, bounds.Algorithm)
			assert.Equal(t, tt.wantParam, bounds.Param)
		})
	}
}

func TestVerifier_malformed(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{
			name:    "too few parts",
			encoded: Prefix + strings.Join([]string{defaultParams, testSalt, testHashA, testSaltSep}, "$"),
		},
		{
			name:    "too many parts",
			encoded: vectorA + "$foo",
		},
		{
			name:    "params scan error",
			encoded: encode("ln=foo,r=8", testSalt, testHashA, testSaltSep, testSignerKey),
		},
		{
			name:    "params with p",
			encoded: encode("ln=14,r=8,p=1", testSalt, testHashA, testSaltSep, testSignerKey),
		},
		{
			name:    "params not canonical",
			encoded: encode("ln=014,r=+8", testSalt, testHashA, testSaltSep, testSignerKey),
		},
		{
			name:    "negative ln",
			encoded: encode("ln=-1,r=8", testSalt, testHashA, testSaltSep, testSignerKey),
		},
		{
			name:    "salt error",
			encoded: encode(defaultParams, "!!!", testHashA, testSaltSep, testSignerKey),
		},
		{
			name:    "hash error",
			encoded: encode(defaultParams, testSalt, "!!!", testSaltSep, testSignerKey),
		},
		{
			name:    "salt separator error",
			encoded: encode(defaultParams, testSalt, testHashA, "!!!", testSignerKey),
		},
		{
			name:    "signer key error",
			encoded: encode(defaultParams, testSalt, testHashA, testSaltSep, "!!!"),
		},
		{
			name:    "empty hash",
			encoded: encode(defaultParams, testSalt, "", testSaltSep, testSignerKey),
		},
		{
			name:    "empty signer key",
			encoded: encode(defaultParams, testSalt, testHashA, testSaltSep, ""),
		},
		{
			name:    "hash and signer key length differ",
			encoded: encode(defaultParams, testSalt, "Bw==", testSaltSep, testSignerKey),
		},
	}
	v := NewVerifier(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Validate(tt.encoded)
			assert.Error(t, err)
			assert.Equal(t, verifier.Skip, got)

			got, err = v.Verify(tt.encoded, testPassword)
			assert.Error(t, err)
			assert.Equal(t, verifier.Skip, got)
		})
	}
}

func TestNewVerifier(t *testing.T) {
	tests := []struct {
		name string
		opts *ValidationOpts
		want *ValidationOpts
	}{
		{
			name: "nil",
			want: DefaultValidationOpts,
		},
		{
			name: "zero value",
			opts: &ValidationOpts{},
			want: DefaultValidationOpts,
		},
		{
			name: "partial",
			opts: &ValidationOpts{MinLN: 10, MaxR: 9},
			want: &ValidationOpts{
				MinLN: 10,
				MaxLN: DefaultMaxLN,
				MinR:  DefaultMinR,
				MaxR:  9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before ValidationOpts
			if tt.opts != nil {
				before = *tt.opts
			}
			got := NewVerifier(tt.opts)
			assert.Equal(t, tt.want, got.opts)
			if tt.opts != nil {
				assert.Equal(t, before, *tt.opts)
			}
		})
	}
}
