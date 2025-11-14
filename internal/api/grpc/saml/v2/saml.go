package saml

import (
	"context"

	"connectrpc.com/connect"
	"github.com/zitadel/logging"
	"github.com/zitadel/saml/pkg/provider"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
)

func (s *Server) GetSAMLRequest(ctx context.Context, req *connect.Request[saml_pb.GetSAMLRequestRequest]) (*connect.Response[saml_pb.GetSAMLRequestResponse], error) {
	authRequest, err := s.query.SamlRequestByID(ctx, true, req.Msg.GetSamlRequestId(), true)
	if err != nil {
		logging.WithError(err).Error("query samlRequest by ID")
		return nil, err
	}
	return connect.NewResponse(&saml_pb.GetSAMLRequestResponse{
		SamlRequest: samlRequestToPb(authRequest),
	}), nil
}

func samlRequestToPb(a *query.SamlRequest) *saml_pb.SAMLRequest {
	return &saml_pb.SAMLRequest{
		Id:           a.ID,
		CreationDate: timestamppb.New(a.CreationDate),
	}
}

func (s *Server) CreateResponse(ctx context.Context, req *connect.Request[saml_pb.CreateResponseRequest]) (*connect.Response[saml_pb.CreateResponseResponse], error) {
	switch v := req.Msg.GetResponseKind().(type) {
	case *saml_pb.CreateResponseRequest_Error:
		return s.failSAMLRequest(ctx, req.Msg.GetSamlRequestId(), v.Error)
	case *saml_pb.CreateResponseRequest_Session:
		return s.linkSessionToSAMLRequest(ctx, req.Msg.GetSamlRequestId(), v.Session)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "SAMLv2-0Tfak3fBS0", "verification oneOf %T in method CreateResponse not implemented", v)
	}
}

func (s *Server) failSAMLRequest(ctx context.Context, samlRequestID string, ae *saml_pb.AuthorizationError) (*connect.Response[saml_pb.CreateResponseResponse], error) {
	details, aar, err := s.command.FailSAMLRequest(ctx, samlRequestID, errorReasonToDomain(ae.GetError()))
	if err != nil {
		return nil, err
	}
	authReq := &saml.AuthRequestV2{CurrentSAMLRequest: aar}
	url, body, err := s.idp.CreateErrorResponse(authReq, errorReasonToDomain(ae.GetError()), ae.GetErrorDescription())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(createCallbackResponseFromBinding(details, url, body, authReq.RelayState)), nil
}

func (s *Server) checkPermission(ctx context.Context, issuer string, userID string) error {
	permission, err := s.query.CheckProjectPermissionByEntityID(ctx, issuer, userID)
	if err != nil {
		return err
	}
	if !permission.HasProjectChecked {
		return zerrors.ThrowPermissionDenied(nil, "SAML-foSyH49RvL", "Errors.User.ProjectRequired")
	}
	if !permission.ProjectRoleChecked {
		return zerrors.ThrowPermissionDenied(nil, "SAML-foSyH49RvL", "Errors.User.GrantRequired")
	}
	return nil
}

func (s *Server) linkSessionToSAMLRequest(ctx context.Context, samlRequestID string, session *saml_pb.Session) (*connect.Response[saml_pb.CreateResponseResponse], error) {
	details, aar, err := s.command.LinkSessionToSAMLRequest(ctx, samlRequestID, session.GetSessionId(), session.GetSessionToken(), true, s.checkPermission)
	if err != nil {
		return nil, err
	}
	authReq := &saml.AuthRequestV2{CurrentSAMLRequest: aar}
	responseIssuer := authReq.ResponseIssuer
	if responseIssuer == "" {
		responseIssuer = http_utils.DomainContext(ctx).Origin()
	}
	ctx = provider.ContextWithIssuer(ctx, responseIssuer)
	url, body, err := s.idp.CreateResponse(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(createCallbackResponseFromBinding(details, url, body, authReq.RelayState)), nil
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
