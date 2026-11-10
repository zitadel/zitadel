package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestChallengePasskeyGRPCToDomain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		challengePasskey *session_grpc.RequestChallenges_WebAuthN
		want             *domain.ChallengeTypePasskey
	}{
		{
			name: "no challenge passkey",
		},
		{
			name: "challenge passkey - user verification required",
			challengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				Domain:                      "example.com",
				UserVerificationRequirement: session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
			},
			want: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
			},
		},
		{
			name: "challenge passkey - user verification preferred",
			challengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				Domain:                      "example.com",
				UserVerificationRequirement: session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED,
			},
			want: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
		},
		{
			name: "challenge passkey - user verification unspecified",
			challengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				Domain:                      "example.com",
				UserVerificationRequirement: session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_UNSPECIFIED,
			},
			want: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementUnspecified,
			},
		},
		{
			name: "challenge passkey - user verification discouraged",
			challengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				Domain:                      "example.com",
				UserVerificationRequirement: session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_DISCOURAGED,
			},
			want: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementDiscouraged,
			},
		},
		{
			name: "challenge passkey - user verification with no value",
			challengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				Domain: "example.com",
			},
			want: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementUnspecified,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, ChallengePasskeyGRPCToDomain(tt.challengePasskey))
		})
	}
}
