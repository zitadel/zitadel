package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) RegisterPasskey(ctx context.Context, req *connect.Request[user.RegisterPasskeyRequest]) (resp *connect.Response[user.RegisterPasskeyResponse], err error) {
	var (
		authenticator = passkeyAuthenticatorToDomain(req.Msg.GetAuthenticator())
	)
	if code := req.Msg.GetCode(); code != nil {
		return passkeyRegistrationDetailsToPb(
			s.command.RegisterUserPasskeyWithCode(ctx, req.Msg.GetUserId(), "", authenticator, code.Id, code.Code, req.Msg.GetDomain(), s.userCodeAlg),
		)
	}
	return passkeyRegistrationDetailsToPb(
		s.command.RegisterUserPasskey(ctx, req.Msg.GetUserId(), "", req.Msg.GetDomain(), authenticator),
	)
}

func passkeyAuthenticatorToDomain(pa user.PasskeyAuthenticator) domain.AuthenticatorAttachment {
	switch pa {
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_UNSPECIFIED:
		return domain.AuthenticatorAttachmentUnspecified
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM:
		return domain.AuthenticatorAttachmentPlattform
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_CROSS_PLATFORM:
		return domain.AuthenticatorAttachmentCrossPlattform
	default:
		return domain.AuthenticatorAttachmentUnspecified
	}
}

func webAuthNRegistrationDetailsToPb(details *domain.WebAuthNRegistrationDetails, err error) (*object_pb.Details, *structpb.Struct, error) {
	if err != nil {
		return nil, nil, err
	}
	options := new(structpb.Struct)
	if err := options.UnmarshalJSON(details.PublicKeyCredentialCreationOptions); err != nil {
		return nil, nil, zerrors.ThrowInternal(err, "USERv2-Dohr6", "Errors.Internal")
	}
	return object.DomainToDetailsPb(details.ObjectDetails), options, nil
}

func passkeyRegistrationDetailsToPb(details *domain.WebAuthNRegistrationDetails, err error) (*connect.Response[user.RegisterPasskeyResponse], error) {
	objectDetails, options, err := webAuthNRegistrationDetailsToPb(details, err)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RegisterPasskeyResponse{
		Details:                            objectDetails,
		PasskeyId:                          details.ID,
		PublicKeyCredentialCreationOptions: options,
	}), nil
}

func (s *Server) VerifyPasskeyRegistration(ctx context.Context, req *connect.Request[user.VerifyPasskeyRegistrationRequest]) (*connect.Response[user.VerifyPasskeyRegistrationResponse], error) {
	pkc, err := req.Msg.GetPublicKeyCredential().MarshalJSON()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USERv2-Pha2o", "Errors.Internal")
	}
	objectDetails, err := s.command.HumanHumanPasswordlessSetup(ctx, req.Msg.GetUserId(), "", req.Msg.GetPasskeyName(), "", pkc)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyPasskeyRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}), nil
}

func (s *Server) CreatePasskeyRegistrationLink(ctx context.Context, req *connect.Request[user.CreatePasskeyRegistrationLinkRequest]) (resp *connect.Response[user.CreatePasskeyRegistrationLinkResponse], err error) {
	switch medium := req.Msg.Medium.(type) {
	case nil:
		return passkeyDetailsToPb(
			s.command.AddUserPasskeyCode(ctx, req.Msg.GetUserId(), "", s.userCodeAlg),
		)
	case *user.CreatePasskeyRegistrationLinkRequest_SendLink:
		return passkeyDetailsToPb(
			s.command.AddUserPasskeyCodeURLTemplate(ctx, req.Msg.GetUserId(), "", s.userCodeAlg, medium.SendLink.GetUrlTemplate()),
		)
	case *user.CreatePasskeyRegistrationLinkRequest_ReturnCode:
		return passkeyCodeDetailsToPb(
			s.command.AddUserPasskeyCodeReturn(ctx, req.Msg.GetUserId(), "", s.userCodeAlg),
		)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "USERv2-gaD8y", "verification oneOf %T in method CreatePasskeyRegistrationLink not implemented", medium)
	}
}

func passkeyDetailsToPb(details *domain.ObjectDetails, err error) (*connect.Response[user.CreatePasskeyRegistrationLinkResponse], error) {
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.CreatePasskeyRegistrationLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func passkeyCodeDetailsToPb(details *domain.PasskeyCodeDetails, err error) (*connect.Response[user.CreatePasskeyRegistrationLinkResponse], error) {
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.CreatePasskeyRegistrationLinkResponse{
		Details: object.DomainToDetailsPb(details.ObjectDetails),
		Code: &user.PasskeyRegistrationCode{
			Id:   details.CodeID,
			Code: details.Code,
		},
	}), nil
}

func (s *Server) RemovePasskey(ctx context.Context, req *connect.Request[user.RemovePasskeyRequest]) (*connect.Response[user.RemovePasskeyResponse], error) {
	objectDetails, err := s.command.HumanRemovePasswordless(ctx, req.Msg.GetUserId(), req.Msg.GetPasskeyId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemovePasskeyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}), nil
}

func (s *Server) ListPasskeys(ctx context.Context, req *connect.Request[user.ListPasskeysRequest]) (*connect.Response[user.ListPasskeysResponse], error) {
	query := new(query.UserAuthMethodSearchQueries)
	err := query.AppendUserIDQuery(req.Msg.UserId)
	if err != nil {
		return nil, err
	}
	err = query.AppendAuthMethodQuery(domain.UserAuthMethodTypePasswordless)
	if err != nil {
		return nil, err
	}
	err = query.AppendStateQuery(domain.MFAStateReady)
	if err != nil {
		return nil, err
	}
	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ListPasskeysResponse{
		Details: object.ToListDetails(authMethods.SearchResponse),
		Result:  authMethodsToPasskeyPb(authMethods),
	}), nil
}

func authMethodsToPasskeyPb(methods *query.AuthMethods) []*user.Passkey {
	t := make([]*user.Passkey, len(methods.AuthMethods))
	for i, token := range methods.AuthMethods {
		t[i] = authMethodToPasskeyPb(token)
	}
	return t
}

func authMethodToPasskeyPb(token *query.AuthMethod) *user.Passkey {
	return &user.Passkey{
		Id:    token.TokenID,
		State: mfaStateToPb(token.State),
		Name:  token.Name,
	}
}

func mfaStateToPb(state domain.MFAState) user.AuthFactorState {
	switch state {
	case domain.MFAStateNotReady:
		return user.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY
	case domain.MFAStateReady:
		return user.AuthFactorState_AUTH_FACTOR_STATE_READY
	case domain.MFAStateUnspecified, domain.MFAStateRemoved:
		// Handle all remaining cases so the linter succeeds
		return user.AuthFactorState_AUTH_FACTOR_STATE_UNSPECIFIED
	default:
		return user.AuthFactorState_AUTH_FACTOR_STATE_UNSPECIFIED
	}
}
