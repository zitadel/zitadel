package convert

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestChallengeOTPEmailGRPCToDomain(t *testing.T) {
	tests := []struct {
		name              string
		otpEmailChallenge *session_grpc.RequestChallenges_OTPEmail
		want              *domain.ChallengeTypeOTPEmail
		wantErr           error
	}{
		{
			name: "nil OTP email challenge",
		},
		{
			name: "otp email challenge - delivery type send code without url template",
			otpEmailChallenge: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{},
				},
			},
			want: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{},
				},
			},
		},
		{
			name: "otp email challenge - delivery type send code with url template",
			otpEmailChallenge: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{
						UrlTemplate: gu.Ptr("https://example.com/otp"),
					},
				},
			},
			want: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "https://example.com/otp",
					},
				},
			},
		},
		{
			name: "otp email challenge - delivery type return code",
			otpEmailChallenge: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_ReturnCode_{
					ReturnCode: &session_grpc.RequestChallenges_OTPEmail_ReturnCode{},
				},
			},
			want: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					ReturnCode: true,
				},
			},
		},
		{
			name:              "otp email challenge - delivery type not set",
			otpEmailChallenge: &session_grpc.RequestChallenges_OTPEmail{},
			want: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ChallengeOTPEmailGRPCToDomain(tt.otpEmailChallenge)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
