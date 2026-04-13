package convert

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_accessTokenTypeToPb(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		input    domain.OIDCTokenType
		expected user.AccessTokenType
	}{
		{
			name:     "Bearer token type",
			input:    domain.OIDCTokenTypeBearer,
			expected: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		},
		{
			name:     "JWT token type",
			input:    domain.OIDCTokenTypeJWT,
			expected: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
		},
		{
			name:     "Unknown token type returns Bearer",
			input:    domain.OIDCTokenType(2),
			expected: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := accessTokenTypeToPb(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}

func Test_machineToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    *query.Machine
		expected *user.MachineUser
	}{
		{
			name: "All fields set, Bearer token, HasSecret true",
			input: &query.Machine{
				Name:            "machine1",
				Description:     "desc",
				EncodedSecret:   "secret",
				AccessTokenType: domain.OIDCTokenTypeBearer,
			},
			expected: &user.MachineUser{
				Name:            "machine1",
				Description:     "desc",
				HasSecret:       true,
				AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
			},
		},
		{
			name: "No secret, JWT token",
			input: &query.Machine{
				Name:            "machine2",
				Description:     "desc2",
				EncodedSecret:   "",
				AccessTokenType: domain.OIDCTokenTypeJWT,
			},
			expected: &user.MachineUser{
				Name:            "machine2",
				Description:     "desc2",
				HasSecret:       false,
				AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
			},
		},
		{
			name: "Unknown token type, HasSecret false",
			input: &query.Machine{
				Name:            "machine3",
				Description:     "desc3",
				EncodedSecret:   "",
				AccessTokenType: domain.OIDCTokenType(99),
			},
			expected: &user.MachineUser{
				Name:            "machine3",
				Description:     "desc3",
				HasSecret:       false,
				AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := machineToPb(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}
