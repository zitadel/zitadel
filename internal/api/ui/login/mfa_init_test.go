package login

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
)

// Test_isMFAPromptStep guards the fix for the unauthenticated MFA-enrollment / phone-overwrite
// issue: MFA init/enrollment handlers must only proceed when the auth request has legitimately
// reached the MFA prompt step (i.e. after the first factor was checked). A request that only
// submitted a loginname carries a first-factor step, not an *domain.MFAPromptStep.
func Test_isMFAPromptStep(t *testing.T) {
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
			name: "existing factor - MFA verification step, not enrollment",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{
					&domain.MFAVerificationStep{MFAProviders: []domain.MFAType{domain.MFATypeTOTP}},
				},
			},
			want: false,
		},
		{
			name: "legitimate enrollment - MFA prompt step",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{
					&domain.MFAPromptStep{Required: true, MFAProviders: []domain.MFAType{domain.MFATypeTOTP}},
				},
			},
			want: true,
		},
		{
			name: "MFA prompt step not first - only the current step counts",
			authReq: &domain.AuthRequest{
				PossibleSteps: []domain.NextStep{
					&domain.PasswordStep{},
					&domain.MFAPromptStep{Required: true, MFAProviders: []domain.MFAType{domain.MFATypeTOTP}},
				},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, isMFAPromptStep(tc.authReq))
		})
	}
}