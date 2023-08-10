package management

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	text_grpc "github.com/zitadel/zitadel/internal/api/grpc/text"
	"github.com/zitadel/zitadel/internal/domain"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetCustomInitMessageText(ctx context.Context, req *mgmt_pb.GetCustomInitMessageTextRequest) (*mgmt_pb.GetCustomInitMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.InitCodeMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomInitMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultInitMessageText(ctx context.Context, req *mgmt_pb.GetDefaultInitMessageTextRequest) (*mgmt_pb.GetDefaultInitMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultInitMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomInitMessageText(ctx context.Context, req *mgmt_pb.SetCustomInitMessageTextRequest) (*mgmt_pb.SetCustomInitMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetInitCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomInitMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomInitMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomInitMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomInitMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.InitCodeMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomInitMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomPasswordResetMessageText(ctx context.Context, req *mgmt_pb.GetCustomPasswordResetMessageTextRequest) (*mgmt_pb.GetCustomPasswordResetMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordResetMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultPasswordResetMessageText(ctx context.Context, req *mgmt_pb.GetDefaultPasswordResetMessageTextRequest) (*mgmt_pb.GetDefaultPasswordResetMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomPasswordResetMessageText(ctx context.Context, req *mgmt_pb.SetCustomPasswordResetMessageTextRequest) (*mgmt_pb.SetCustomPasswordResetMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetPasswordResetCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomPasswordResetMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomPasswordResetMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomPasswordResetMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomPasswordResetMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordResetMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomPasswordResetMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomVerifyEmailMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifyEmailMessageTextRequest) (*mgmt_pb.GetCustomVerifyEmailMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifyEmailMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifyEmailMessageTextRequest) (*mgmt_pb.GetDefaultVerifyEmailMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomVerifyEmailMessageText(ctx context.Context, req *mgmt_pb.SetCustomVerifyEmailMessageTextRequest) (*mgmt_pb.SetCustomVerifyEmailMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetVerifyEmailCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomVerifyEmailMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifyEmailMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomVerifyEmailMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomVerifyEmailMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomVerifyEmailMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomVerifyPhoneMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifyPhoneMessageTextRequest) (*mgmt_pb.GetCustomVerifyPhoneMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyPhoneMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifyPhoneMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifyPhoneMessageTextRequest) (*mgmt_pb.GetDefaultVerifyPhoneMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomVerifyPhoneMessageText(ctx context.Context, req *mgmt_pb.SetCustomVerifyPhoneMessageTextRequest) (*mgmt_pb.SetCustomVerifyPhoneMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetVerifyPhoneCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomVerifyPhoneMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifyPhoneMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomVerifyPhoneMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomVerifyPhoneMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyPhoneMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomVerifyPhoneMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomVerifySMSOTPMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifySMSOTPMessageTextRequest) (*mgmt_pb.GetCustomVerifySMSOTPMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifySMSOTPMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifySMSOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifySMSOTPMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifySMSOTPMessageTextRequest) (*mgmt_pb.GetDefaultVerifySMSOTPMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.VerifySMSOTPMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifySMSOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomVerifySMSOTPMessageText(ctx context.Context, req *mgmt_pb.SetCustomVerifySMSOTPMessageTextRequest) (*mgmt_pb.SetCustomVerifySMSOTPMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetVerifySMSOTPCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomVerifySMSOTPMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifySMSOTPMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomVerifySMSOTPMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomVerifySMSOTPMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifySMSOTPMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomVerifySMSOTPMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomVerifyEmailOTPMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifyEmailOTPMessageTextRequest) (*mgmt_pb.GetCustomVerifyEmailOTPMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailOTPMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyEmailOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifyEmailOTPMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifyEmailOTPMessageTextRequest) (*mgmt_pb.GetDefaultVerifyEmailOTPMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.VerifyEmailOTPMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifyEmailOTPMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomVerifyEmailOTPMessageText(ctx context.Context, req *mgmt_pb.SetCustomVerifyEmailOTPMessageTextRequest) (*mgmt_pb.SetCustomVerifyEmailOTPMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetVerifyEmailOTPCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomVerifyEmailOTPMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomVerifyEmailOTPMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailOTPMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomVerifyEmailOTPMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.GetCustomDomainClaimedMessageTextRequest) (*mgmt_pb.GetCustomDomainClaimedMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.DomainClaimedMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.GetDefaultDomainClaimedMessageTextRequest) (*mgmt_pb.GetDefaultDomainClaimedMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomDomainClaimedMessageCustomText(ctx context.Context, req *mgmt_pb.SetCustomDomainClaimedMessageTextRequest) (*mgmt_pb.SetCustomDomainClaimedMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetDomainClaimedCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomDomainClaimedMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomDomainClaimedMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomDomainClaimedMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomDomainClaimedMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.DomainClaimedMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomDomainClaimedMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomPasswordChangeMessageText(ctx context.Context, req *mgmt_pb.GetCustomPasswordChangeMessageTextRequest) (*mgmt_pb.GetCustomPasswordChangeMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordChangeMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomPasswordChangeMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultPasswordChangeMessageText(ctx context.Context, req *mgmt_pb.GetDefaultPasswordChangeMessageTextRequest) (*mgmt_pb.GetDefaultPasswordChangeMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.PasswordChangeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordChangeMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomPasswordChangeMessageCustomText(ctx context.Context, req *mgmt_pb.SetCustomPasswordChangeMessageTextRequest) (*mgmt_pb.SetCustomPasswordChangeMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetPasswordChangeCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomPasswordChangeMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomPasswordChangeMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomPasswordChangeMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomPasswordChangeMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordChangeMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomPasswordChangeMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomPasswordlessRegistrationMessageText(ctx context.Context, req *mgmt_pb.GetCustomPasswordlessRegistrationMessageTextRequest) (*mgmt_pb.GetCustomPasswordlessRegistrationMessageTextResponse, error) {
	msg, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordlessRegistrationMessageType, req.Language, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomPasswordlessRegistrationMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultPasswordlessRegistrationMessageText(ctx context.Context, req *mgmt_pb.GetDefaultPasswordlessRegistrationMessageTextRequest) (*mgmt_pb.GetDefaultPasswordlessRegistrationMessageTextResponse, error) {
	msg, err := s.query.IAMMessageTextByTypeAndLanguage(ctx, domain.PasswordlessRegistrationMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordlessRegistrationMessageTextResponse{
		CustomText: text_grpc.ModelCustomMessageTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomPasswordlessRegistrationMessageCustomText(ctx context.Context, req *mgmt_pb.SetCustomPasswordlessRegistrationMessageTextRequest) (*mgmt_pb.SetCustomPasswordlessRegistrationMessageTextResponse, error) {
	result, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, SetPasswordlessRegistrationCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomPasswordlessRegistrationMessageTextResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomPasswordlessRegistrationMessageTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest) (*mgmt_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, error) {
	result, err := s.command.RemoveOrgMessageTexts(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordlessRegistrationMessageType, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetCustomLoginTexts(ctx context.Context, req *mgmt_pb.GetCustomLoginTextsRequest) (*mgmt_pb.GetCustomLoginTextsResponse, error) {
	msg, err := s.query.GetCustomLoginTexts(ctx, authz.GetCtxData(ctx).OrgID, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultLoginTexts(ctx context.Context, req *mgmt_pb.GetDefaultLoginTextsRequest) (*mgmt_pb.GetDefaultLoginTextsResponse, error) {
	msg, err := s.query.IAMLoginTexts(ctx, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomLoginText(ctx context.Context, req *mgmt_pb.SetCustomLoginTextsRequest) (*mgmt_pb.SetCustomLoginTextsResponse, error) {
	result, err := s.command.SetOrgLoginText(ctx, authz.GetCtxData(ctx).OrgID, SetLoginCustomTextToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetCustomLoginTextsResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomLoginTextToDefault(ctx context.Context, req *mgmt_pb.ResetCustomLoginTextsToDefaultRequest) (*mgmt_pb.ResetCustomLoginTextsToDefaultResponse, error) {
	result, err := s.command.RemoveOrgLoginTexts(ctx, authz.GetCtxData(ctx).OrgID, language.Make(req.Language))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetCustomLoginTextsToDefaultResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}
