package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckIDPIntentGRPCToDomain(t *testing.T) {
	t.Parallel()
	tt := []struct {
		testName string
		input    *session_grpc.CheckIDPIntent
		expected *domain.CheckIDPIntentType
	}{
		{
			testName: "nil input",
			input:    nil,
			expected: nil,
		},
		{
			testName: "valid input",
			input: &session_grpc.CheckIDPIntent{
				IdpIntentId:    "intent-123",
				IdpIntentToken: "token-abc",
			},
			expected: &domain.CheckIDPIntentType{
				ID:    "intent-123",
				Token: "token-abc",
			},
		},
		{
			testName: "empty strings",
			input: &session_grpc.CheckIDPIntent{
				IdpIntentId:    "",
				IdpIntentToken: "",
			},
			expected: &domain.CheckIDPIntentType{
				ID:    "",
				Token: "",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			result := CheckIDPIntentGRPCToDomain(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
