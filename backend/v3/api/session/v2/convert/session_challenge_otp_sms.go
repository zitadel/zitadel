package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func ChallengeOTPSMSGRPCToDomain(otpSMSChallenge *session_grpc.RequestChallenges_OTPSMS) *domain.ChallengeTypeOTPSMS {
	if otpSMSChallenge == nil {
		return nil
	}
	return &domain.ChallengeTypeOTPSMS{
		ReturnCode: otpSMSChallenge.GetReturnCode(),
	}
}
