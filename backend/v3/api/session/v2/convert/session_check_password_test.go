package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckPasswordGRPCToDomain(t *testing.T) {
	t.Parallel()
	tt := []struct {
		testName      string
		inputCheckPsw *session_grpc.CheckPassword
		expected      *domain.CheckPasswordType
	}{
		{testName: "when nil input should return nil output"},
		{
			testName:      "when not nil input should return valid output",
			inputCheckPsw: &session_grpc.CheckPassword{Password: "pw"},
			expected:      &domain.CheckPasswordType{Password: "pw"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Test
			res := CheckPasswordGRPCToDomain(tc.inputCheckPsw)

			// Verify
			assert.Equal(t, tc.expected, res)
		})
	}
}
