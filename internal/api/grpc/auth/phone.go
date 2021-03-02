package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyPhone(ctx context.Context, _ *auth_pb.GetMyPhoneRequest) (*auth_pb.GetMyPhoneResponse, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPhoneResponse{
		Phone: user.ModelPhoneToPb(phone),
		Details: object.ToDetailsPb(
			phone.Sequence,
			phone.CreationDate,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) SetMyPhone(ctx context.Context, req *auth_pb.SetMyPhoneRequest) (*auth_pb.SetMyPhoneResponse, error) {
	phone, err := s.command.ChangeHumanPhone(ctx, UpdateMyPhoneToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyPhoneResponse{
		Details: object.ToDetailsPb(
			phone.Sequence,
			phone.CreationDate,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyPhone(ctx context.Context, req *auth_pb.VerifyMyPhoneRequest) (*auth_pb.VerifyMyPhoneResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	_, err := s.command.VerifyHumanPhone(ctx, ctxData.UserID, req.Code, ctxData.OrgID)
	if err != nil {
		return nil, err
	}

	//TODO: response from business
	return &auth_pb.VerifyMyPhoneResponse{
		//Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendMyPhoneVerification(ctx context.Context, _ *auth_pb.ResendMyPhoneVerificationRequest) (*auth_pb.ResendMyPhoneVerificationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.CreateHumanPhoneVerificationCode(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ResendMyPhoneVerificationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyPhone(ctx context.Context, _ *auth_pb.RemoveMyPhoneRequest) (*auth_pb.RemoveMyPhoneResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.RemoveHumanPhone(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyPhoneResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
