package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	text_grpc "github.com/caos/zitadel/internal/api/grpc/text"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetInitMessageCustomText(ctx context.Context, req *admin_pb.GetInitMessageCustomTextRequest) (*admin_pb.GetInitMessageCustomTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetInitMessageCustomTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateInitMessageCustomText(ctx context.Context, req *admin_pb.SetInitMessageCustomTextRequest) (*admin_pb.SetInitMessageCustomTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetInitCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetInitMessageCustomTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetPasswordResetMessageCustomText(ctx context.Context, req *admin_pb.GetPasswordResetMessageCustomTextRequest) (*admin_pb.GetPasswordResetMessageCustomTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetPasswordResetMessageCustomTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdatePasswordResetMessageCustomText(ctx context.Context, req *admin_pb.SetPasswordResetMessageCustomTextRequest) (*admin_pb.SetPasswordResetMessageCustomTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetPasswordResetCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetPasswordResetMessageCustomTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetVerifyEmailMessageCustomText(ctx context.Context, req *admin_pb.GetVerifyEmailMessageCustomTextRequest) (*admin_pb.GetVerifyEmailMessageCustomTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetVerifyEmailMessageCustomTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateVerifyEmailMessageCustomText(ctx context.Context, req *admin_pb.SetVerifyEmailMessageCustomTextRequest) (*admin_pb.SetVerifyEmailMessageCustomTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetVerifyEmailCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetVerifyEmailMessageCustomTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetVerifyPhoneMessageCustomText(ctx context.Context, req *admin_pb.GetVerifyPhoneMessageCustomTextRequest) (*admin_pb.GetVerifyPhoneMessageCustomTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetVerifyPhoneMessageCustomTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateVerifyPhoneMessageCustomText(ctx context.Context, req *admin_pb.SetVerifyPhoneMessageCustomTextRequest) (*admin_pb.SetVerifyPhoneMessageCustomTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetVerifyPhoneCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetVerifyPhoneMessageCustomTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDomainClaimedMessageCustomText(ctx context.Context, req *admin_pb.GetDomainClaimedMessageCustomTextRequest) (*admin_pb.GetDomainClaimedMessageCustomTextResponse, error) {
	msg, err := s.iam.GetDefaultMailText(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDomainClaimedMessageCustomTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) UpdateDomainClaimedMessageCustomText(ctx context.Context, req *admin_pb.SetDomainClaimedMessageCustomTextRequest) (*admin_pb.SetDomainClaimedMessageCustomTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, SetDomainClaimedCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDomainClaimedMessageCustomTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}
