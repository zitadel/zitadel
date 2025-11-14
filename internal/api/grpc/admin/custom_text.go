package admin

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	text_grpc "github.com/zitadel/zitadel/internal/api/grpc/text"
	"github.com/zitadel/zitadel/internal/domain"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultInitMessageText(ctx context.Context, req *admin_pb.GetDefaultInitMessageTextRequest) (*admin_pb.GetDefaultInitMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultInitMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomInitMessageText(ctx context.Context, req *admin_pb.GetCustomInitMessageTextRequest) (*admin_pb.GetCustomInitMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.InitCodeMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomInitMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultInitMessageText(ctx context.Context, req *admin_pb.SetDefaultInitMessageTextRequest) (*admin_pb.SetDefaultInitMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetInitCustomTextToDomain(req))
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

func (s *Server) ResetCustomInitMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomInitMessageTextToDefaultRequest) (*admin_pb.ResetCustomInitMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.InitCodeMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomInitMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultPasswordResetMessageText(ctx context.Context, req *admin_pb.GetDefaultPasswordResetMessageTextRequest) (*admin_pb.GetDefaultPasswordResetMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomPasswordResetMessageText(ctx context.Context, req *admin_pb.GetCustomPasswordResetMessageTextRequest) (*admin_pb.GetCustomPasswordResetMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.PasswordResetMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultPasswordResetMessageText(ctx context.Context, req *admin_pb.SetDefaultPasswordResetMessageTextRequest) (*admin_pb.SetDefaultPasswordResetMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetPasswordResetCustomTextToDomain(req))
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

func (s *Server) ResetCustomPasswordResetMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomPasswordResetMessageTextToDefaultRequest) (*admin_pb.ResetCustomPasswordResetMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.PasswordResetMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomPasswordResetMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultVerifyEmailMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifyEmailMessageTextRequest) (*admin_pb.GetDefaultVerifyEmailMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifyEmailMessageText(ctx context.Context, req *admin_pb.GetCustomVerifyEmailMessageTextRequest) (*admin_pb.GetCustomVerifyEmailMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.VerifyEmailMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifyEmailMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifyEmailMessageTextRequest) (*admin_pb.SetDefaultVerifyEmailMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetVerifyEmailCustomTextToDomain(req))
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

func (s *Server) ResetCustomVerifyEmailMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomVerifyEmailMessageTextToDefaultRequest) (*admin_pb.ResetCustomVerifyEmailMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.VerifyEmailMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomVerifyEmailMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultVerifyPhoneMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.GetDefaultVerifyPhoneMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifyPhoneMessageText(ctx context.Context, req *admin_pb.GetCustomVerifyPhoneMessageTextRequest) (*admin_pb.GetCustomVerifyPhoneMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.VerifyPhoneMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifyPhoneMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifyPhoneMessageTextRequest) (*admin_pb.SetDefaultVerifyPhoneMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetVerifyPhoneCustomTextToDomain(req))
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

func (s *Server) ResetCustomVerifyPhoneMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomVerifyPhoneMessageTextToDefaultRequest) (*admin_pb.ResetCustomVerifyPhoneMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.VerifyPhoneMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomVerifyPhoneMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultVerifySMSOTPMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifySMSOTPMessageTextRequest) (*admin_pb.GetDefaultVerifySMSOTPMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.VerifySMSOTPMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifySMSOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifySMSOTPMessageText(ctx context.Context, req *admin_pb.GetCustomVerifySMSOTPMessageTextRequest) (*admin_pb.GetCustomVerifySMSOTPMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.VerifySMSOTPMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifySMSOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifySMSOTPMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifySMSOTPMessageTextRequest) (*admin_pb.SetDefaultVerifySMSOTPMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetVerifySMSOTPCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultVerifySMSOTPMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifySMSOTPMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomVerifySMSOTPMessageTextToDefaultRequest) (*admin_pb.ResetCustomVerifySMSOTPMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.VerifySMSOTPMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomVerifySMSOTPMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultVerifyEmailOTPMessageText(ctx context.Context, req *admin_pb.GetDefaultVerifyEmailOTPMessageTextRequest) (*admin_pb.GetDefaultVerifyEmailOTPMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.VerifyEmailOTPMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultVerifyEmailOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomVerifyEmailOTPMessageText(ctx context.Context, req *admin_pb.GetCustomVerifyEmailOTPMessageTextRequest) (*admin_pb.GetCustomVerifyEmailOTPMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.VerifyEmailOTPMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomVerifyEmailOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultVerifyEmailOTPMessageText(ctx context.Context, req *admin_pb.SetDefaultVerifyEmailOTPMessageTextRequest) (*admin_pb.SetDefaultVerifyEmailOTPMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetVerifyEmailOTPCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultVerifyEmailOTPMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifyEmailOTPMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultRequest) (*admin_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.VerifyEmailOTPMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultDomainClaimedMessageText(ctx context.Context, req *admin_pb.GetDefaultDomainClaimedMessageTextRequest) (*admin_pb.GetDefaultDomainClaimedMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomDomainClaimedMessageText(ctx context.Context, req *admin_pb.GetCustomDomainClaimedMessageTextRequest) (*admin_pb.GetCustomDomainClaimedMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.DomainClaimedMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultDomainClaimedMessageText(ctx context.Context, req *admin_pb.SetDefaultDomainClaimedMessageTextRequest) (*admin_pb.SetDefaultDomainClaimedMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetDomainClaimedCustomTextToDomain(req))
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

func (s *Server) ResetCustomDomainClaimedMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomDomainClaimedMessageTextToDefaultRequest) (*admin_pb.ResetCustomDomainClaimedMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.DomainClaimedMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomDomainClaimedMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultPasswordChangeMessageText(ctx context.Context, req *admin_pb.GetDefaultPasswordChangeMessageTextRequest) (*admin_pb.GetDefaultPasswordChangeMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.PasswordChangeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordChangeMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomPasswordChangeMessageText(ctx context.Context, req *admin_pb.GetCustomPasswordChangeMessageTextRequest) (*admin_pb.GetCustomPasswordChangeMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.PasswordChangeMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomPasswordChangeMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultPasswordChangeMessageText(ctx context.Context, req *admin_pb.SetDefaultPasswordChangeMessageTextRequest) (*admin_pb.SetDefaultPasswordChangeMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetPasswordChangeCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultPasswordChangeMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomPasswordChangeMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomPasswordChangeMessageTextToDefaultRequest) (*admin_pb.ResetCustomPasswordChangeMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.PasswordChangeMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomPasswordChangeMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultInviteUserMessageText(ctx context.Context, req *admin_pb.GetDefaultInviteUserMessageTextRequest) (*admin_pb.GetDefaultInviteUserMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.InviteUserMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultInviteUserMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomInviteUserMessageText(ctx context.Context, req *admin_pb.GetCustomInviteUserMessageTextRequest) (*admin_pb.GetCustomInviteUserMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.InviteUserMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomInviteUserMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultInviteUserMessageText(ctx context.Context, req *admin_pb.SetDefaultInviteUserMessageTextRequest) (*admin_pb.SetDefaultInviteUserMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetInviteUserCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultInviteUserMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomInviteUserMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomInviteUserMessageTextToDefaultRequest) (*admin_pb.ResetCustomInviteUserMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.InviteUserMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomInviteUserMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultPasswordlessRegistrationMessageText(ctx context.Context, req *admin_pb.GetDefaultPasswordlessRegistrationMessageTextRequest) (*admin_pb.GetDefaultPasswordlessRegistrationMessageTextResponse, error) {
	msg, err := s.query.DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx, domain.PasswordlessRegistrationMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordlessRegistrationMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetCustomPasswordlessRegistrationMessageText(ctx context.Context, req *admin_pb.GetCustomPasswordlessRegistrationMessageTextRequest) (*admin_pb.GetCustomPasswordlessRegistrationMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetInstance(ctx).InstanceID(), domain.PasswordlessRegistrationMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomPasswordlessRegistrationMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetDefaultPasswordlessRegistrationMessageText(ctx context.Context, req *admin_pb.SetDefaultPasswordlessRegistrationMessageTextRequest) (*admin_pb.SetDefaultPasswordlessRegistrationMessageTextResponse, error) {
	result, err := s.command.SetDefaultMessageText(ctx, authz.GetInstance(ctx).InstanceID(), SetPasswordlessRegistrationCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultPasswordlessRegistrationMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomPasswordlessRegistrationMessageTextToDefault(ctx context.Context, req *admin_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest) (*admin_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveInstanceMessageTexts(ctx, domain.PasswordlessRegistrationMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetDefaultLoginTexts(ctx context.Context, req *admin_pb.GetDefaultLoginTextsRequest) (*admin_pb.GetDefaultLoginTextsResponse, error) {
	msg, err := s.query.GetDefaultLoginTexts(ctx, req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}
func (s *Server) GetCustomLoginTexts(ctx context.Context, req *admin_pb.GetCustomLoginTextsRequest) (*admin_pb.GetCustomLoginTextsResponse, error) {
	msg, err := s.query.GetCustomLoginTexts(ctx, authz.GetInstance(ctx).InstanceID(), req.Language)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomLoginText(ctx context.Context, req *admin_pb.SetCustomLoginTextsRequest) (*admin_pb.SetCustomLoginTextsResponse, error) {
	result, err := s.command.SetCustomInstanceLoginText(ctx, SetLoginTextToDomain(req))
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
	result, err := s.command.RemoveCustomInstanceLoginTexts(ctx, language.Make(req.Language))
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
