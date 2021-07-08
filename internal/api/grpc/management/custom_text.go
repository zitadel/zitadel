package management

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	text_grpc "github.com/caos/zitadel/internal/api/grpc/text"
	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetCustomInitMessageText(ctx context.Context, req *mgmt_pb.GetCustomInitMessageTextRequest) (*mgmt_pb.GetCustomInitMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomInitMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultInitMessageText(ctx context.Context, req *mgmt_pb.GetDefaultInitMessageTextRequest) (*mgmt_pb.GetDefaultInitMessageTextResponse, error) {
	msg, err := s.org.GetDefaultMessageText(ctx, domain.InitCodeMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultInitMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
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
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomPasswordResetMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultPasswordResetMessageText(ctx context.Context, req *mgmt_pb.GetDefaultPasswordResetMessageTextRequest) (*mgmt_pb.GetDefaultPasswordResetMessageTextResponse, error) {
	msg, err := s.org.GetDefaultMessageText(ctx, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordResetMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
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
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyEmailMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifyEmailMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifyEmailMessageTextRequest) (*mgmt_pb.GetDefaultVerifyEmailMessageTextResponse, error) {
	msg, err := s.org.GetDefaultMessageText(ctx, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifyEmailMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
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
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultVerifyPhoneMessageText(ctx context.Context, req *mgmt_pb.GetDefaultVerifyPhoneMessageTextRequest) (*mgmt_pb.GetDefaultVerifyPhoneMessageTextResponse, error) {
	msg, err := s.org.GetDefaultMessageText(ctx, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
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

func (s *Server) GetCustomDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.GetCustomDomainClaimedMessageTextRequest) (*mgmt_pb.GetCustomDomainClaimedMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomDomainClaimedMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.GetDefaultDomainClaimedMessageTextRequest) (*mgmt_pb.GetDefaultDomainClaimedMessageTextResponse, error) {
	msg, err := s.org.GetDefaultMessageText(ctx, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultDomainClaimedMessageTextResponse{
		CustomText: text_grpc.DomainCustomMsgTextToPb(msg),
	}, nil
}

func (s *Server) SetCustomDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.SetCustomDomainClaimedMessageTextRequest) (*mgmt_pb.SetCustomDomainClaimedMessageTextResponse, error) {
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

func (s *Server) GetCustomLoginTexts(ctx context.Context, req *mgmt_pb.GetCustomLoginTextsRequest) (*mgmt_pb.GetCustomLoginTextsResponse, error) {
	msg, err := s.org.GetLoginTexts(ctx, authz.GetCtxData(ctx).OrgID, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomLoginTextsResponse{
		CustomText: text_grpc.CustomLoginTextToPb(msg),
	}, nil
}

func (s *Server) GetDefaultLoginTexts(ctx context.Context, req *mgmt_pb.GetDefaultLoginTextsRequest) (*mgmt_pb.GetDefaultLoginTextsResponse, error) {
	msg, err := s.org.GetDefaultLoginTexts(ctx, req.Language)
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
