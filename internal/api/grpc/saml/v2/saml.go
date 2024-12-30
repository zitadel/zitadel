package saml

import (
	"context"

	"github.com/zitadel/logging"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
)

func (s *Server) GetAuthRequest(ctx context.Context, req *saml_pb.GetSAMLRequestRequest) (*saml_pb.GetSAMLRequestResponse, error) {
	authRequest, err := s.query.SamlRequestByID(ctx, true, req.GetSamlRequestId(), true)
	if err != nil {
		logging.WithError(err).Error("query samlRequest by ID")
		return nil, err
	}
	return &saml_pb.GetSAMLRequestResponse{
		SamlRequest: samlRequestToPb(authRequest),
	}, nil
}

func samlRequestToPb(a *query.SamlRequest) *saml_pb.SAMLRequest {
	return &saml_pb.SAMLRequest{
		Id:           a.ID,
		CreationDate: timestamppb.New(a.CreationDate),
	}
}

func (s *Server) CreateResponse(ctx context.Context, req *saml_pb.CreateResponseRequest) (*saml_pb.CreateResponseResponse, error) {
	switch v := req.GetResponseKind().(type) {
	case *saml_pb.CreateResponseRequest_Error:
		return s.failSAMLRequest(ctx, req.GetSamlRequestId(), v.Error)
	case *saml_pb.CreateResponseRequest_Session:
		return s.linkSessionToSAMLRequest(ctx, req.GetSamlRequestId(), v.Session)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "SAMLv2-0Tfak3fBS0", "verification oneOf %T in method CreateResponse not implemented", v)
	}
}

func (s *Server) failSAMLRequest(ctx context.Context, samlRequestID string, ae *saml_pb.AuthorizationError) (*saml_pb.CreateResponseResponse, error) {
	details, aar, err := s.command.FailSAMLRequest(ctx, samlRequestID, errorReasonToDomain(ae.GetError()))
	if err != nil {
		return nil, err
	}
	authReq := &saml.AuthRequestV2{CurrentSAMLRequest: aar}
	url, body, err := s.idp.CreateErrorResponse(authReq, errorReasonToDomain(ae.GetError()), ae.GetErrorDescription())
	if err != nil {
		return nil, err
	}
	return createCallbackResponseFromBinding(details, url, body, authReq.RelayState), nil
}

func (s *Server) linkSessionToSAMLRequest(ctx context.Context, samlRequestID string, session *saml_pb.Session) (*saml_pb.CreateResponseResponse, error) {
	details, aar, err := s.command.LinkSessionToSAMLRequest(ctx, samlRequestID, session.GetSessionId(), session.GetSessionToken(), true)
	if err != nil {
		return nil, err
	}
	authReq := &saml.AuthRequestV2{CurrentSAMLRequest: aar}
	url, body, err := s.idp.CreateResponse(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return createCallbackResponseFromBinding(details, url, body, authReq.RelayState), nil
}

func createCallbackResponseFromBinding(details *domain.ObjectDetails, url string, body string, relayState string) *saml_pb.CreateResponseResponse {
	resp := &saml_pb.CreateResponseResponse{
		Details: object.DomainToDetailsPb(details),
		Url:     url,
	}

	if body != "" {
		resp.Binding = &saml_pb.CreateResponseResponse_Post{
			Post: &saml_pb.PostResponse{
				RelayState:   relayState,
				SamlResponse: body,
			},
		}
	} else {
		resp.Binding = &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}}
	}
	return resp
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
