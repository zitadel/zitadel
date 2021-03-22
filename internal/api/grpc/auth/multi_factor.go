package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func (s *Server) ListMyAuthFactors(ctx context.Context, _ *auth_pb.ListMyAuthFactorsRequest) (*auth_pb.ListMyAuthFactorsResponse, error) {
	mfas, err := s.repo.MyUserMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyAuthFactorsResponse{
		Result: user_grpc.AuthFactorsToPb(mfas),
	}, nil
}

func (s *Server) AddMyAuthFactorOTP(ctx context.Context, _ *auth_pb.AddMyAuthFactorOTPRequest) (*auth_pb.AddMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	otp, err := s.command.AddHumanOTP(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorOTPResponse{
		Url:    otp.Url,
		Secret: otp.SecretString,
		Details: object.AddToDetailsPb(
			otp.Sequence,
			otp.ChangeDate,
			otp.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyAuthFactorOTP(ctx context.Context, req *auth_pb.VerifyMyAuthFactorOTPRequest) (*auth_pb.VerifyMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanCheckMFAOTPSetup(ctx, ctxData.UserID, req.Code, "", ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyAuthFactorOTPResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyAuthFactorOTP(ctx context.Context, _ *auth_pb.RemoveMyAuthFactorOTPRequest) (*auth_pb.RemoveMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanRemoveOTP(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAuthFactorOTPResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) AddMyAuthFactorU2F(ctx context.Context, _ *auth_pb.AddMyAuthFactorU2FRequest) (*auth_pb.AddMyAuthFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	u2f, err := s.command.HumanAddU2FSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorU2FResponse{
		Key: &user_pb.WebAuthNKey{
			Id:        u2f.WebAuthNTokenID,
			PublicKey: u2f.CredentialCreationData,
		},
		Details: object.AddToDetailsPb(
			u2f.Sequence,
			u2f.ChangeDate,
			u2f.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyAuthFactorU2F(ctx context.Context, req *auth_pb.VerifyMyAuthFactorU2FRequest) (*auth_pb.VerifyMyAuthFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanVerifyU2FSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Verification.TokenName, "", req.Verification.PublicKeyCredential)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyAuthFactorU2FResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyAuthFactorU2F(ctx context.Context, req *auth_pb.RemoveMyAuthFactorU2FRequest) (*auth_pb.RemoveMyAuthFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanRemoveU2F(ctx, ctxData.UserID, req.TokenId, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAuthFactorU2FResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
