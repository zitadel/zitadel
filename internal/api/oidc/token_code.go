package oidc

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) CodeExchange(ctx context.Context, r *op.ClientRequest[oidc.AccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	client, ok := r.Client.(*Client)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "OIDC-Ae2ph", "Error.Internal")
	}

	plainCode, err := s.decryptCode(ctx, r.Data.Code)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "OIDC-ahLi2", "Errors.User.Code.Invalid")
	}

	var (
		session *command.OIDCSession
		state   string
	)
	if strings.HasPrefix(plainCode, command.IDPrefixV2) {
		session, state, err = s.command.CreateOIDCSessionFromCodeExchange(
			setContextUserSystem(ctx), plainCode, authRequestComplianceChecker(client, r.Data),
		)
	} else {
		session, state, err = s.codeExchangeV1(ctx, client, r.Data, plainCode)
	}
	if err != nil {
		return nil, err
	}
	return response(s.accessTokenResponseFromSession(ctx, client, session, state, client.client.ProjectID, client.client.ProjectRoleAssertion))
}

// codeExchangeV1 creates a v2 token from a v1 auth request.
func (s *Server) codeExchangeV1(ctx context.Context, client *Client, req *oidc.AccessTokenRequest, plainCode string) (session *command.OIDCSession, state string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authReq, err := s.getAuthRequestV1ByCode(ctx, plainCode)
	if err != nil {
		return nil, "", err
	}

	if challenge := authReq.GetCodeChallenge(); challenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
		if err = op.AuthorizeCodeChallenge(req.CodeVerifier, challenge); err != nil {
			return nil, "", err
		}
	}
	if req.RedirectURI != authReq.GetRedirectURI() {
		return nil, "", oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
	}
	userAgentID, _, userOrgID, authTime, authMethodsReferences, reason, actor := getInfoFromRequest(authReq)

	scope := authReq.GetScopes()
	session, err = s.command.CreateOIDCSession(ctx,
		authReq.GetSubject(),
		userOrgID,
		client.client.ClientID,
		scope,
		authReq.GetAudience(),
		AMRToAuthMethodTypes(authMethodsReferences),
		authTime,
		&domain.UserAgent{
			FingerprintID: &userAgentID,
		},
		reason,
		actor,
		slices.Contains(scope, oidc.ScopeOfflineAccess),
	)
	return session, authReq.GetState(), err
}

func (s *Server) getAuthRequestV1ByCode(ctx context.Context, plainCode string) (op.AuthRequest, error) {
	authReq, err := s.repo.AuthRequestByCode(ctx, plainCode)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(authReq)
}

func authRequestComplianceChecker(client *Client, req *oidc.AccessTokenRequest) command.AuthRequestComplianceChecker {
	return func(ctx context.Context, authReq *command.AuthRequestWriteModel) error {
		if authReq.CodeChallenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
			err := op.AuthorizeCodeChallenge(req.CodeVerifier, CodeChallengeToOIDC(authReq.CodeChallenge))
			if err != nil {
				return err
			}
		}
		if req.RedirectURI != authReq.RedirectURI {
			return oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
		}
		if err := authReq.CheckAuthenticated(); err != nil {
			return err
		}
		return nil
	}
}
