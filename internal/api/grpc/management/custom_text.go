package management

import (
	"context"

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
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
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

func (s *Server) GetCustomPasswordResetMessageText(ctx context.Context, req *mgmt_pb.GetCustomPasswordResetMessageTextRequest) (*mgmt_pb.GetCustomPasswordResetMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.PasswordResetMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomPasswordResetMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
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

func (s *Server) GetCustomVerifyEmailMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifyEmailMessageTextRequest) (*mgmt_pb.GetCustomVerifyEmailMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyEmailMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyEmailMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
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

func (s *Server) GetCustomVerifyPhoneMessageText(ctx context.Context, req *mgmt_pb.GetCustomVerifyPhoneMessageTextRequest) (*mgmt_pb.GetCustomVerifyPhoneMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.VerifyPhoneMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomVerifyPhoneMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
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

func (s *Server) GetCustomDomainClaimedMessageText(ctx context.Context, req *mgmt_pb.GetCustomDomainClaimedMessageTextRequest) (*mgmt_pb.GetCustomDomainClaimedMessageTextResponse, error) {
	msg, err := s.org.GetMessageText(ctx, authz.GetCtxData(ctx).OrgID, domain.DomainClaimedMessageType, req.Language)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetCustomDomainClaimedMessageTextResponse{
		CustomText: text_grpc.ModelCustomMsgTextToPb(msg),
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
