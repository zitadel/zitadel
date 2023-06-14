package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterOTP(ctx context.Context, req *user.RegisterOTPRequest) (*user.RegisterOTPResponse, error) {
	return otpDetailsToPb(
		s.command.AddUserOTP(ctx, req.GetUserId(), authz.GetCtxData(ctx).ResourceOwner),
	)

}

func otpDetailsToPb(otp *domain.OTPv2, err error) (*user.RegisterOTPResponse, error) {
	if err != nil {
		return nil, err
	}
	return &user.RegisterOTPResponse{
		Details: object.DomainToDetailsPb(otp.ObjectDetails),
		Uri:     otp.URI,
		Secret:  otp.Secret,
	}, nil
}

func (s *Server) VerifyOTPRegistration(context.Context, *user.VerifyOTPRegistrationRequest) (*user.VerifyOTPRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyOTPRegistration not implemented")
}
