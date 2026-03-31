package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func ChallengeOTPEmailGRPCToDomain(otpEmailChallenge *session_grpc.RequestChallenges_OTPEmail) (*domain.ChallengeTypeOTPEmail, error) {
	if otpEmailChallenge == nil {
		return nil, nil
	}

	switch t := otpEmailChallenge.DeliveryType.(type) {
	case *session_grpc.RequestChallenges_OTPEmail_SendCode_:
		return &domain.ChallengeTypeOTPEmail{
			DeliveryType: domain.DeliveryType{
				SendCode: &domain.SendCode{
					URLTemplate: t.SendCode.GetUrlTemplate(),
				},
			},
		}, nil
	case *session_grpc.RequestChallenges_OTPEmail_ReturnCode_:
		return &domain.ChallengeTypeOTPEmail{
			DeliveryType: domain.DeliveryType{
				ReturnCode: true,
			},
		}, nil
	case nil:
		return &domain.ChallengeTypeOTPEmail{
			DeliveryType: domain.DeliveryType{},
		}, nil
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "SESSION-mfil3D", "delivery_type oneOf %T in OTPEmailChallenge not implemented", t)
	}
}
