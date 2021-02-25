package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func (s *Server) ListMyMultiFactors(ctx context.Context, _ *auth_pb.ListMyMultiFactorsRequest) (*auth_pb.ListMyMultiFactorsResponse, error) {
	mfas, err := s.repo.MyUserMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMultiFactorsResponse{
		Result: user_grpc.MultiFactorsToPb(mfas),
	}, nil
}

func (s *Server) AddMyMultiFactorOTP(ctx context.Context, _ *auth_pb.AddMyMultiFactorOTPRequest) (*auth_pb.AddMyMultiFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	otp, err := s.command.AddHumanOTP(ctx, ctxData.UserID, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyMultiFactorOTPResponse{
		Url:    otp.Url,
		Secret: otp.SecretString,
		Details: object.ToDetailsPb(
			otp.Sequence,
			otp.CreationDate,
			otp.ChangeDate,
			otp.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyMultiFactorOTP(ctx context.Context, req *auth_pb.VerifyMyMultiFactorOTPRequest) (*auth_pb.VerifyMyMultiFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanCheckMFAOTPSetup(ctx, ctxData.UserID, req.Code, "", ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.VerifyMyMultiFactorOTPResponse{}, nil
}

func (s *Server) RemoveMyMultiFactorOTP(ctx context.Context, _ *auth_pb.RemoveMyMultiFactorOTPRequest) (*auth_pb.RemoveMyMultiFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanRemoveOTP(ctx, ctxData.UserID, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.RemoveMyMultiFactorOTPResponse{}, nil
}

func (s *Server) AddMyMultiFactorU2F(ctx context.Context, _ *auth_pb.AddMyMultiFactorU2FRequest) (*auth_pb.AddMyMultiFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	u2f, err := s.command.HumanAddU2FSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyMultiFactorU2FResponse{
		Key: &user_pb.WebAuthNKey{
			Id:        u2f.WebAuthNTokenID,
			PublicKey: u2f.CredentialCreationData,
		},
		Details: object.ToDetailsPb(
			u2f.Sequence,
			u2f.CreationDate,
			u2f.ChangeDate,
			u2f.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyMultiFactorU2F(ctx context.Context, req *auth_pb.VerifyMyMultiFactorU2FRequest) (*auth_pb.VerifyMyMultiFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanVerifyU2FSetup(ctx, ctxData.UserID, ctxData.OrgID, req.Verification.TokenName, "", req.Verification.PublicKeyCredential)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.VerifyMyMultiFactorU2FResponse{}, nil
}

func (s *Server) RemoveMyMultiFactorU2F(ctx context.Context, req *auth_pb.RemoveMyMultiFactorU2FRequest) (*auth_pb.RemoveMyMultiFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanRemovePasswordless(ctx, ctxData.UserID, req.TokenId, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.RemoveMyMultiFactorU2FResponse{}, nil
}
