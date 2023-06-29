package oidc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
)

func (s *Server) GetAuthRequest(ctx context.Context, req *oidc_pb.GetAuthRequestRequest) (*oidc_pb.GetAuthRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthRequest not implemented")
}

func authRequestToPb(a *domain.AuthRequest) *oidc_pb.AuthRequest {
	pba := &oidc_pb.AuthRequest{
		Id:           a.ID,
		CreationDate: timestamppb.New(a.CreationDate),
		ClientId:     a.ApplicationID,
		RedirectUri:  a.CallbackURI,
		Prompt:       promptsToPb(a.Prompt),
		UiLocales:    a.UiLocales,
	}
	if oidcAuthReq, ok := a.Request.(*domain.AuthRequestOIDC); ok { // should we throw an error here if the type != OIDC?
		pba.Scope = oidcAuthReq.Scopes
	}
	if a.LoginHint != "" {
		pba.LoginHint = &a.LoginHint
	}
	if a.MaxAuthAge != nil {
		pba.MaxAge = durationpb.New(*a.MaxAuthAge)
	}
	if a.UserID != "" {
		pba.HintUserId = &a.UserID
	}
	return pba
}

func promptsToPb(promps []domain.Prompt) []oidc_pb.Prompt {
	out := make([]oidc_pb.Prompt, len(promps))
	for i, p := range promps {
		out[i] = promptToPb(p)
	}
	return out
}

func promptToPb(p domain.Prompt) oidc_pb.Prompt {
	switch p {
	case domain.PromptUnspecified:
		return oidc_pb.Prompt_PROMPT_UNSPECIFIED
	case domain.PromptNone:
		return oidc_pb.Prompt_PROMPT_NONE
	case domain.PromptLogin:
		return oidc_pb.Prompt_PROMPT_LOGIN
	case domain.PromptConsent:
		return oidc_pb.Prompt_PROMPT_CONSENT
	case domain.PromptSelectAccount:
		return oidc_pb.Prompt_PROMPT_SELECT_ACCOUNT
	case domain.PromptCreate:
		return oidc_pb.Prompt_PROMPT_CREATE
	default:
		return oidc_pb.Prompt_PROMPT_UNSPECIFIED
	}
}

func (s *Server) CreateCallback(ctx context.Context, req *oidc_pb.CreateCallbackRequest) (*oidc_pb.CreateCallbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCallback not implemented")
}

func errorReasonToOIDC(reason oidc_pb.ErrorReason) string {
	switch reason {
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST:
		return "invalid_request"
	case oidc_pb.ErrorReason_ERROR_REASON_UNAUTHORIZED_CLIENT:
		return "unauthorized_client"
	case oidc_pb.ErrorReason_ERROR_REASON_ACCESS_DENIED:
		return "access_denied"
	case oidc_pb.ErrorReason_ERROR_REASON_UNSUPPORTED_RESPONSE_TYPE:
		return "unsupported_response_type"
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_SCOPE:
		return "invalid_scope"
	case oidc_pb.ErrorReason_ERROR_REASON_TEMPORARY_UNAVAILABLE:
		return "temporarily_unavailable"
	case oidc_pb.ErrorReason_ERROR_REASON_INTERACTION_REQUIRED:
		return "interaction_required"
	case oidc_pb.ErrorReason_ERROR_REASON_LOGIN_REQUIRED:
		return "login_required"
	case oidc_pb.ErrorReason_ERROR_REASON_ACCOUNT_SELECTION_REQUIRED:
		return "account_selection_required"
	case oidc_pb.ErrorReason_ERROR_REASON_CONSENT_REQUIRED:
		return "consent_required"
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_URI:
		return "invalid_request_uri"
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_OBJECT:
		return "invalid_request_object"
	case oidc_pb.ErrorReason_ERROR_REASON_REQUEST_NOT_SUPPORTED:
		return "request_not_supported"
	case oidc_pb.ErrorReason_ERROR_REASON_REQUEST_URI_NOT_SUPPORTED:
		return "request_uri_not_supported"
	case oidc_pb.ErrorReason_ERROR_REASON_REGISTRATION_NOT_SUPPORTED:
		return "registration_not_supported"
	case oidc_pb.ErrorReason_ERROR_REASON_UNSPECIFIED, oidc_pb.ErrorReason_ERROR_REASON_SERVER_ERROR:
		fallthrough
	default:
		return "server_error"
	}
}
