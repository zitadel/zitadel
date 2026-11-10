package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestChallengeOTPSMSGRPCToDomain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		challenge *session_grpc.RequestChallenges_OTPSMS
		want      *domain.ChallengeTypeOTPSMS
	}{
		{
			name: "nil OTP SMS challenge",
		},
		{
			name: "OTP SMS challenge true",
			challenge: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			want: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
		},
		{
			name: "OTP SMS challenge false",
			challenge: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: false,
			},
			want: &domain.ChallengeTypeOTPSMS{
				ReturnCode: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ChallengeOTPSMSGRPCToDomain(tt.challenge)
			assert.Equal(t, tt.want, got)
		})
	}
}
