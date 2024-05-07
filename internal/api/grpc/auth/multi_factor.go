package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func (s *Server) ListMyAuthFactors(ctx context.Context, _ *auth_pb.ListMyAuthFactorsRequest) (*auth_pb.ListMyAuthFactorsResponse, error) {
	query := new(query.UserAuthMethodSearchQueries)
	err := query.AppendUserIDQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	err = query.AppendAuthMethodsQuery(domain.UserAuthMethodTypeU2F, domain.UserAuthMethodTypeTOTP, domain.UserAuthMethodTypeOTPSMS, domain.UserAuthMethodTypeOTPEmail)
	if err != nil {
		return nil, err
	}
	err = query.AppendStateQuery(domain.MFAStateReady)
	if err != nil {
		return nil, err
	}
	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyAuthFactorsResponse{
		Result: user_grpc.AuthMethodsToPb(authMethods),
	}, nil
}

func (s *Server) AddMyAuthFactorOTP(ctx context.Context, _ *auth_pb.AddMyAuthFactorOTPRequest) (*auth_pb.AddMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	otp, err := s.command.AddHumanTOTP(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorOTPResponse{
		Url:    otp.URI,
		Secret: otp.Secret,
		Details: object.AddToDetailsPb(
			otp.Sequence,
			otp.EventDate,
			otp.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyAuthFactorOTP(ctx context.Context, req *auth_pb.VerifyMyAuthFactorOTPRequest) (*auth_pb.VerifyMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanCheckMFATOTPSetup(ctx, ctxData.UserID, req.Code, "", ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyAuthFactorOTPResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyAuthFactorOTP(ctx context.Context, _ *auth_pb.RemoveMyAuthFactorOTPRequest) (*auth_pb.RemoveMyAuthFactorOTPResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanRemoveTOTP(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAuthFactorOTPResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) AddMyAuthFactorOTPSMS(ctx context.Context, _ *auth_pb.AddMyAuthFactorOTPSMSRequest) (*auth_pb.AddMyAuthFactorOTPSMSResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	details, err := s.command.AddHumanOTPSMS(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorOTPSMSResponse{
		Details: object.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemoveMyAuthFactorOTPSMS(ctx context.Context, _ *auth_pb.RemoveMyAuthFactorOTPSMSRequest) (*auth_pb.RemoveMyAuthFactorOTPSMSResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	details, err := s.command.RemoveHumanOTPSMS(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAuthFactorOTPSMSResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddMyAuthFactorOTPEmail(ctx context.Context, _ *auth_pb.AddMyAuthFactorOTPEmailRequest) (*auth_pb.AddMyAuthFactorOTPEmailResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	details, err := s.command.AddHumanOTPEmail(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorOTPEmailResponse{
		Details: object.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemoveMyAuthFactorOTPEmail(ctx context.Context, _ *auth_pb.RemoveMyAuthFactorOTPEmailRequest) (*auth_pb.RemoveMyAuthFactorOTPEmailResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	details, err := s.command.RemoveHumanOTPEmail(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAuthFactorOTPEmailResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddMyAuthFactorU2F(ctx context.Context, _ *auth_pb.AddMyAuthFactorU2FRequest) (*auth_pb.AddMyAuthFactorU2FResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	u2f, err := s.command.HumanAddU2FSetup(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyAuthFactorU2FResponse{
		Key: &user_pb.WebAuthNKey{
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
