package login

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
)

// Test_isChangeUsernameStep guards the fix for the unauthenticated account-rename issue:
// handleChangeUsername must only proceed when the auth request has legitimately reached the
// change-username step (which nextSteps appends only after the first factor and any required MFA have
// succeeded). A request that merely submitted a loginname carries a first-factor step, so an
// unauthenticated caller must not be able to rename the bound account.
func Test_isChangeUsernameStep(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		authReq *domain.AuthRequest
		want    bool
	}{
		{
			name:    "nil auth request",
			authReq: nil,
			want:    false,
		},
		{
			name:    "no possible steps",
			authReq: &domain.AuthRequest{},
			want:    false,
		},
		{
			name: "loginname only - first factor (password) step",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{&domain.PasswordStep{}},
			},
			want: false,
		},
		{
			name: "not yet authenticated - MFA verification step",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{
					&domain.MFAVerificationStep{MFAProviders: []domain.MFAType{domain.MFATypeTOTP}},
				},
			},
			want: false,
		},
		{
			name: "legitimate change-username step",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{&domain.ChangeUsernameStep{}},
			},
			want: true,
		},
		{
			name: "change-username step not first - only the current step counts",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{
					&domain.ChangePasswordStep{},
					&domain.ChangeUsernameStep{},
				},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, isChangeUsernameStep(tc.authReq))
		})
	}
}
