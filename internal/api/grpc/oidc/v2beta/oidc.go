package oidc

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
)

func (s *Server) GetAuthRequest(ctx context.Context, req *oidc_pb.GetAuthRequestRequest) (*oidc_pb.GetAuthRequestResponse, error) {
	authRequest, err := s.query.AuthRequestByID(ctx, true, req.GetAuthRequestId(), true)
	if err != nil {
		logging.WithError(err).Error("query authRequest by ID")
		return nil, err
	}
	return &oidc_pb.GetAuthRequestResponse{
		AuthRequest: authRequestToPb(authRequest),
	}, nil
}

func authRequestToPb(a *query.AuthRequest) *oidc_pb.AuthRequest {
	pba := &oidc_pb.AuthRequest{
		Id:           a.ID,
		CreationDate: timestamppb.New(a.CreationDate),
		ClientId:     a.ClientID,
		Scope:        a.Scope,
		RedirectUri:  a.RedirectURI,
		Prompt:       promptsToPb(a.Prompt),
		UiLocales:    a.UiLocales,
		LoginHint:    a.LoginHint,
		HintUserId:   a.HintUserID,
	}
	if a.MaxAge != nil {
		pba.MaxAge = durationpb.New(*a.MaxAge)
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
	switch v := req.GetCallbackKind().(type) {
	case *oidc_pb.CreateCallbackRequest_Error:
		return s.failAuthRequest(ctx, req.GetAuthRequestId(), v.Error)
	case *oidc_pb.CreateCallbackRequest_Session:
		return s.linkSessionToAuthRequest(ctx, req.GetAuthRequestId(), v.Session)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "OIDCv2-zee7A", "verification oneOf %T in method CreateCallback not implemented", v)
	}
}

func (s *Server) failAuthRequest(ctx context.Context, authRequestID string, ae *oidc_pb.AuthorizationError) (*oidc_pb.CreateCallbackResponse, error) {
	details, aar, err := s.command.FailAuthRequest(ctx, authRequestID, errorReasonToDomain(ae.GetError()))
	if err != nil {
		return nil, err
	}
	authReq := &oidc.AuthRequestV2{CurrentAuthRequest: aar}
	callback, err := oidc.CreateErrorCallbackURL(authReq, errorReasonToOIDC(ae.GetError()), ae.GetErrorDescription(), ae.GetErrorUri(), s.op.Provider())
	if err != nil {
		return nil, err
	}
	return &oidc_pb.CreateCallbackResponse{
		Details:     object.DomainToDetailsPb(details),
		CallbackUrl: callback,
	}, nil
}

func (s *Server) linkSessionToAuthRequest(ctx context.Context, authRequestID string, session *oidc_pb.Session) (*oidc_pb.CreateCallbackResponse, error) {
	details, aar, err := s.command.LinkSessionToAuthRequest(ctx, authRequestID, session.GetSessionId(), session.GetSessionToken(), true)
	if err != nil {
		return nil, err
	}
	authReq := &oidc.AuthRequestV2{CurrentAuthRequest: aar}
	ctx = op.ContextWithIssuer(ctx, http.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), s.externalSecure))
	var callback string
	if aar.ResponseType == domain.OIDCResponseTypeCode {
		callback, err = oidc.CreateCodeCallbackURL(ctx, authReq, s.op.Provider())
	} else {
		callback, err = s.op.CreateTokenCallbackURL(ctx, authReq)
	}
	if err != nil {
		return nil, err
	}
	return &oidc_pb.CreateCallbackResponse{
		Details:     object.DomainToDetailsPb(details),
		CallbackUrl: callback,
	}, nil
}

func errorReasonToDomain(errorReason oidc_pb.ErrorReason) domain.OIDCErrorReason {
	switch errorReason {
	case oidc_pb.ErrorReason_ERROR_REASON_UNSPECIFIED:
		return domain.OIDCErrorReasonUnspecified
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST:
		return domain.OIDCErrorReasonInvalidRequest
	case oidc_pb.ErrorReason_ERROR_REASON_UNAUTHORIZED_CLIENT:
		return domain.OIDCErrorReasonUnauthorizedClient
	case oidc_pb.ErrorReason_ERROR_REASON_ACCESS_DENIED:
		return domain.OIDCErrorReasonAccessDenied
	case oidc_pb.ErrorReason_ERROR_REASON_UNSUPPORTED_RESPONSE_TYPE:
		return domain.OIDCErrorReasonUnsupportedResponseType
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_SCOPE:
		return domain.OIDCErrorReasonInvalidScope
	case oidc_pb.ErrorReason_ERROR_REASON_SERVER_ERROR:
		return domain.OIDCErrorReasonServerError
	case oidc_pb.ErrorReason_ERROR_REASON_TEMPORARY_UNAVAILABLE:
		return domain.OIDCErrorReasonTemporaryUnavailable
	case oidc_pb.ErrorReason_ERROR_REASON_INTERACTION_REQUIRED:
		return domain.OIDCErrorReasonInteractionRequired
	case oidc_pb.ErrorReason_ERROR_REASON_LOGIN_REQUIRED:
		return domain.OIDCErrorReasonLoginRequired
	case oidc_pb.ErrorReason_ERROR_REASON_ACCOUNT_SELECTION_REQUIRED:
		return domain.OIDCErrorReasonAccountSelectionRequired
	case oidc_pb.ErrorReason_ERROR_REASON_CONSENT_REQUIRED:
		return domain.OIDCErrorReasonConsentRequired
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_URI:
		return domain.OIDCErrorReasonInvalidRequestURI
	case oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_OBJECT:
		return domain.OIDCErrorReasonInvalidRequestObject
	case oidc_pb.ErrorReason_ERROR_REASON_REQUEST_NOT_SUPPORTED:
		return domain.OIDCErrorReasonRequestNotSupported
	case oidc_pb.ErrorReason_ERROR_REASON_REQUEST_URI_NOT_SUPPORTED:
		return domain.OIDCErrorReasonRequestURINotSupported
	case oidc_pb.ErrorReason_ERROR_REASON_REGISTRATION_NOT_SUPPORTED:
		return domain.OIDCErrorReasonRegistrationNotSupported
	default:
		return domain.OIDCErrorReasonUnspecified
	}
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
