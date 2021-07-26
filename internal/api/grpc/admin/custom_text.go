package admin

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/grpc/object"
	text_grpc "github.com/caos/zitadel/internal/api/grpc/text"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultInitMessageText(ctx context.Context, req *admin_pb.GetDefaultInitMessageTextRequest) (*admin_pb.GetDefaultInitMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMessageText(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultInitMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomInitMessageText(ctx context.Context, req *admin_pb.GetCustomInitMessageTextRequest) (*admin_pb.GetCustomInitMessageTextResponse, error) {
	msg, err := s.iam.GetCustomMessageText(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomInitMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultInitMessageText(ctx context.Context, req *admin_pb.SetDefaultInitMessageTextRequest) (*admin_pb.SetDefaultInitMessageTextResponse, error) {
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

func (s *Server) GetDefaultPasswordResetMessageText(ctx context.Context, req *admin_pb.GetDefaultPasswordResetMessageTextRequest) (*admin_pb.GetDefaultPasswordResetMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMessageText(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordResetMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomPasswordResetMessageText(ctx context.Context, req *admin_pb.GetCustomPasswordResetMessageTextRequest) (*admin_pb.GetCustomPasswordResetMessageTextResponse, error) {
	msg, err := s.iam.GetCustomMessageText(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomPasswordResetMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultPasswordResetMessageText(ctx context.Context, req *admin_pb.SetDefaultPasswordResetMessageTextRequest) (*admin_pb.SetDefaultPasswordResetMessageTextResponse, error) {
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

func (s *Server) GetDefaultVerifyEmailMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifyEmailMessageTextRequest) (*admin_pb.GetDefaultVerifyEmailMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMessageText(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyEmailMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifyEmailMessageText(ctx context.Context, req *admin_pb.GetCustomVerifyEmailMessageTextRequest) (*admin_pb.GetCustomVerifyEmailMessageTextResponse, error) {
	msg, err := s.iam.GetCustomMessageText(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifyEmailMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifyEmailMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifyEmailMessageTextRequest) (*admin_pb.SetDefaultVerifyEmailMessageTextResponse, error) {
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

func (s *Server) GetDefaultVerifyPhoneMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.GetDefaultVerifyPhoneMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMessageText(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifyPhoneMessageText(ctx context.Context, req *admin_pb.GetCustomVerifyPhoneMessageTextRequest) (*admin_pb.GetCustomVerifyPhoneMessageTextResponse, error) {
	msg, err := s.iam.GetCustomMessageText(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifyPhoneMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.SetDefaultVerifyPhoneMessageTextResponse, error) {
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

func (s *Server) GetDefaultDomainClaimedMessageText(ctx context.Context, req *admin_pb.GetDefaultDomainClaimedMessageTextRequest) (*admin_pb.GetDefaultDomainClaimedMessageTextResponse, error) {
	msg, err := s.iam.GetDefaultMessageText(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultDomainClaimedMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomDomainClaimedMessageText(ctx context.Context, req *admin_pb.GetCustomDomainClaimedMessageTextRequest) (*admin_pb.GetCustomDomainClaimedMessageTextResponse, error) {
	msg, err := s.iam.GetCustomMessageText(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomDomainClaimedMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultDomainClaimedMessageText(ctx context.Context, req *admin_pb.SetDefaultDomainClaimedMessageTextRequest) (*admin_pb.SetDefaultDomainClaimedMessageTextResponse, error) {
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

func (s *Server) GetDefaultLoginTexts(ctx context.Context, req *admin_pb.GetDefaultLoginTextsRequest) (*admin_pb.GetDefaultLoginTextsResponse, error) {
	msg, err := s.iam.GetDefaultLoginTexts(ctx, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}
func (s *Server) GetCustomLoginTexts(ctx context.Context, req *admin_pb.GetCustomLoginTextsRequest) (*admin_pb.GetCustomLoginTextsResponse, error) {
	msg, err := s.iam.GetCustomLoginTexts(ctx, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomLoginText(ctx context.Context, req *admin_pb.SetCustomLoginTextsRequest) (*admin_pb.SetCustomLoginTextsResponse, error) {
	result, err := s.command.SetCustomIAMLoginText(ctx, SetLoginTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetCustomLoginTextsResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomLoginTextToDefault(ctx context.Context, req *admin_pb.ResetCustomLoginTextsToDefaultRequest) (*admin_pb.ResetCustomLoginTextsToDefaultResponse, error) {
	result, err := s.command.RemoveCustomIAMLoginTexts(ctx, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomLoginTextsToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}
