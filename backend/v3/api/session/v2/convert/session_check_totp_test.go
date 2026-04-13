package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckTOTPGRPCToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    *session_grpc.CheckTOTP
		expected *domain.CheckTOTPType
	}{
		{name: "nil input returns nil"},
		{
			name:     "code is mapped",
			input:    &session_grpc.CheckTOTP{Code: "123456"},
			expected: &domain.CheckTOTPType{Code: "123456"},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := CheckTOTPGRPCToDomain(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
