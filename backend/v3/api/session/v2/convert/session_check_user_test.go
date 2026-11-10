package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckUserGRPCToQueryOpt(t *testing.T) {
	t.Parallel()
	tt := []struct {
		testName                string
		inputCheckUser          *session_grpc.CheckUser
		expectedError           error
		expectedDomainCheckUser *domain.CheckUserType
	}{
		{
			testName: "when input checkUser is nil should return nil, nil",
		},
		{
			testName:       "when input checkUser query is not matching type should return invalid argument error",
			inputCheckUser: &session_grpc.CheckUser{},
			expectedError:  zerrors.ThrowInvalidArgumentf(nil, "CONV-7B2m0b", "user search %T not implemented", nil),
		},
		{
			testName: "when input checkUser query is UserID should set user id on domain CheckUserType",
			inputCheckUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{UserId: "user-123"},
			},
			expectedDomainCheckUser: &domain.CheckUserType{UserID: "user-123"},
		},
		{
			testName: "when input checkUser query is LoginName should set user id on domain CheckUserType",
			inputCheckUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_LoginName{LoginName: "user@example.com"},
			},
			expectedDomainCheckUser: &domain.CheckUserType{LoginName: "user@example.com"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Test
			res, err := CheckUserGRPCToDomain(tc.inputCheckUser)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedDomainCheckUser, res)
		})
	}
}
