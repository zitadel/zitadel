package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyPhone(ctx context.Context, _ *auth_pb.GetMyPhoneRequest) (*auth_pb.GetMyPhoneResponse, error) {
	phone, err := s.query.GetHumanPhone(ctx, authz.GetCtxData(ctx).UserID, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPhoneResponse{
		Phone: user.ModelPhoneToPb(phone),
		Details: object.ToViewDetailsPb(
			phone.Sequence,
			phone.CreationDate,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) SetMyPhone(ctx context.Context, req *auth_pb.SetMyPhoneRequest) (*auth_pb.SetMyPhoneResponse, error) {
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	phone, err := s.command.ChangeHumanPhone(ctx, UpdateMyPhoneToDomain(ctx, req), authz.GetCtxData(ctx).ResourceOwner, phoneCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyPhoneResponse{
		Details: object.ChangeToDetailsPb(
			phone.Sequence,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyPhone(ctx context.Context, req *auth_pb.VerifyMyPhoneRequest) (*auth_pb.VerifyMyPhoneResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.VerifyHumanPhone(ctx, ctxData.UserID, req.Code, ctxData.ResourceOwner, phoneCodeGenerator)
	if err != nil {
		return nil, err
	}

	return &auth_pb.VerifyMyPhoneResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendMyPhoneVerification(ctx context.Context, _ *auth_pb.ResendMyPhoneVerificationRequest) (*auth_pb.ResendMyPhoneVerificationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.CreateHumanPhoneVerificationCode(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ResendMyPhoneVerificationResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyPhone(ctx context.Context, _ *auth_pb.RemoveMyPhoneRequest) (*auth_pb.RemoveMyPhoneResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.RemoveHumanPhone(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyPhoneResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
