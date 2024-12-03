package saml

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zitadel/saml/pkg/provider"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
)

func (s *Server) GetAuthRequest(ctx context.Context, req *saml_pb.GetSAMLRequestRequest) (*saml_pb.GetSAMLRequestResponse, error) {
	authRequest, err := s.query.AuthRequestByID(ctx, true, req.GetSamlRequestId(), true)
	if err != nil {
		logging.WithError(err).Error("query samlRequest by ID")
		return nil, err
	}
	return &saml_pb.GetSAMLRequestResponse{
		SamlRequest: samlRequestToPb(authRequest),
	}, nil
}

func samlRequestToPb(a *query.AuthRequest) *saml_pb.SAMLRequest {
	return &saml_pb.SAMLRequest{
		Id:           a.ID,
		CreationDate: timestamppb.New(a.CreationDate),
	}
}

func (s *Server) CreateCallback(ctx context.Context, req *saml_pb.CreateCallbackRequest) (*saml_pb.CreateCallbackResponse, error) {
	switch v := req.GetCallbackKind().(type) {
	case *saml_pb.CreateCallbackRequest_Error:
		return s.failAuthRequest(ctx, req.GetSamlRequestId(), v.Error)
	case *saml_pb.CreateCallbackRequest_Session:
		return s.linkSessionToAuthRequest(ctx, req.GetSamlRequestId(), v.Session)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "OIDCv2-zee7A", "verification oneOf %T in method CreateCallback not implemented", v)
	}
}

func (s *Server) failAuthRequest(ctx context.Context, samlRequestID string, ae *saml_pb.AuthorizationError) (*saml_pb.CreateCallbackResponse, error) {
	details, aar, err := s.command.FailSAMLRequest(ctx, samlRequestID, errorReasonToDomain(ae.GetError()))
	if err != nil {
		return nil, err
	}
	authReq := &saml.AuthRequestV2{CurrentSAMLRequest: aar}
	callback, err := saml.CreateErrorCallbackURL(authReq, errorReasonToSAML(ae.GetError()), ae.GetErrorDescription(), ae.GetErrorUri())
	if err != nil {
		return nil, err
	}
	return &saml_pb.CreateCallbackResponse{
		Details:     object.DomainToDetailsPb(details),
		CallbackUrl: callback,
	}, nil
}

func (s *Server) linkSessionToAuthRequest(ctx context.Context, samlRequestID string, session *saml_pb.Session) (*saml_pb.CreateCallbackResponse, error) {
	details, aar, err := s.command.LinkSessionToAuthRequest(ctx, samlRequestID, session.GetSessionId(), session.GetSessionToken(), true)
	if err != nil {
		return nil, err
	}
	authReq := &oidc.AuthRequestV2{CurrentAuthRequest: aar}
	ctx = op.ContextWithIssuer(ctx, http.DomainContext(ctx).Origin())
	var callback string
	if aar.ResponseType == domain.OIDCResponseTypeCode {
		callback, err = oidc.CreateCodeCallbackURL(ctx, authReq, s.op.Provider())
	} else {
		callback, err = s.op.CreateTokenCallbackURL(ctx, authReq)
	}
	if err != nil {
		return nil, err
	}
	return &saml_pb.CreateCallbackResponse{
		Details:     object.DomainToDetailsPb(details),
		CallbackUrl: callback,
	}, nil
}

func errorReasonToDomain(errorReason saml_pb.ErrorReason) domain.SAMLErrorReason {
	switch errorReason {
	case saml_pb.ErrorReason_ERROR_REASON_UNSPECIFIED:
		return domain.SAMLErrorReasonUnspecified
	case saml_pb.ErrorReason_ERROR_REASON_VERSION_MISSMATCH:
		return domain.SAMLErrorReasonVersionMissmatch
	case saml_pb.ErrorReason_ERROR_REASON_AUTH_N_FAILED:
		return domain.SAMLErrorReasonAuthNFailed
	case saml_pb.ErrorReason_ERROR_REASON_INVALID_ATTR_NAME_OR_VALUE:
		return domain.SAMLErrorReasonInvalidAttrNameOrValue
	case saml_pb.ErrorReason_ERROR_REASON_INVALID_NAMEID_POLICY:
		return domain.SAMLErrorReasonInvalidNameIDPolicy
	case saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED:
		return domain.SAMLErrorReasonRequestDenied
	case saml_pb.ErrorReason_ERROR_REASON_REQUEST_UNSUPPORTED:
		return domain.SAMLErrorReasonRequestUnsupported
	case saml_pb.ErrorReason_ERROR_REASON_UNSUPPORTED_BINDING:
		return domain.SAMLErrorReasonUnsupportedBinding
	default:
		return domain.SAMLErrorReasonUnspecified
	}
}

func errorReasonToSAML(reason saml_pb.ErrorReason) string {
	switch reason {
	case saml_pb.ErrorReason_ERROR_REASON_UNSPECIFIED:
		return "unspecified error"
	case saml_pb.ErrorReason_ERROR_REASON_VERSION_MISSMATCH:
		return provider.StatusCodeVersionMissmatch
	case saml_pb.ErrorReason_ERROR_REASON_AUTH_N_FAILED:
		return provider.StatusCodeAuthNFailed
	case saml_pb.ErrorReason_ERROR_REASON_INVALID_ATTR_NAME_OR_VALUE:
		return provider.StatusCodeInvalidAttrNameOrValue
	case saml_pb.ErrorReason_ERROR_REASON_INVALID_NAMEID_POLICY:
		return provider.StatusCodeInvalidNameIDPolicy
	case saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED:
		return provider.StatusCodeRequestDenied
	case saml_pb.ErrorReason_ERROR_REASON_REQUEST_UNSUPPORTED:
		return provider.StatusCodeRequestUnsupported
	case saml_pb.ErrorReason_ERROR_REASON_UNSUPPORTED_BINDING:
		return provider.StatusCodeUnsupportedBinding
	default:
		return "unspecified error"
	}
}
