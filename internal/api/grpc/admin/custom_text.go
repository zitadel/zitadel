package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	text_grpc "github.com/caos/zitadel/internal/api/grpc/text"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetInitMessageCustomText(ctx context.Context, req *admin_pb.GetDefaultInitMessageTextRequest) (*admin_pb.GetDefaultInitMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultInitMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateInitMessageCustomText(ctx context.Context, req *admin_pb.SetDefaultInitMessageTextRequest) (*admin_pb.SetDefaultInitMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetInitCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultInitMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetPasswordResetMessageCustomText(ctx context.Context, req *admin_pb.GetDefaultPasswordResetMessageTextRequest) (*admin_pb.GetDefaultPasswordResetMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdatePasswordResetMessageCustomText(ctx context.Context, req *admin_pb.SetDefaultPasswordResetMessageTextRequest) (*admin_pb.SetDefaultPasswordResetMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetPasswordResetCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultPasswordResetMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetVerifyEmailMessageCustomText(ctx context.Context, req *admin_pb.GetDefaultVerifyEmailMessageTextRequest) (*admin_pb.GetDefaultVerifyEmailMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateVerifyEmailMessageCustomText(ctx context.Context, req *admin_pb.SetDefaultVerifyEmailMessageTextRequest) (*admin_pb.SetDefaultVerifyEmailMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetVerifyEmailCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultVerifyEmailMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetVerifyPhoneMessageCustomText(ctx context.Context, req *admin_pb.GetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.GetDefaultVerifyPhoneMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateVerifyPhoneMessageCustomText(ctx context.Context, req *admin_pb.SetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.SetDefaultVerifyPhoneMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetVerifyPhoneCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultVerifyPhoneMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDomainClaimedMessageCustomText(ctx context.Context, req *admin_pb.GetDefaultDomainClaimedMessageTextRequest) (*admin_pb.GetDefaultDomainClaimedMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateDomainClaimedMessageCustomText(ctx context.Context, req *admin_pb.SetDefaultDomainClaimedMessageTextRequest) (*admin_pb.SetDefaultDomainClaimedMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetDomainClaimedCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultDomainClaimedMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}
