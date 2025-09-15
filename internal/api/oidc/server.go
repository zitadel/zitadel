package oidc

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Server struct {
	http.Handler
	*op.LegacyServer

	repo              repository.Repository
	query             *query.Queries
	command           *command.Commands
	accessTokenKeySet *oidcKeySet
	idTokenHintKeySet *oidcKeySet

	defaultLoginURL            string
	defaultLoginURLV2          string
	defaultLogoutURLV2         string
	defaultAccessTokenLifetime time.Duration
	defaultIdTokenLifetime     time.Duration
	jwksCacheControlMaxAge     time.Duration

	fallbackLogger            *slog.Logger
	hasher                    *crypto.Hasher
	signingKeyAlgorithm       string
	encAlg                    crypto.EncryptionAlgorithm
	targetEncryptionAlgorithm crypto.EncryptionAlgorithm
	opCrypto                  op.Crypto

	assetAPIPrefix func(ctx context.Context) string
}

func endpoints(endpointConfig *EndpointConfig) op.Endpoints {
	// some defaults. The new Server will disable endpoints that are nil.
	endpoints := op.Endpoints{
		Authorization:       op.NewEndpoint("/oauth/v2/authorize"),
		Token:               op.NewEndpoint("/oauth/v2/token"),
		Introspection:       op.NewEndpoint("/oauth/v2/introspect"),
		Userinfo:            op.NewEndpoint("/oidc/v1/userinfo"),
		Revocation:          op.NewEndpoint("/oauth/v2/revoke"),
		EndSession:          op.NewEndpoint("/oidc/v1/end_session"),
		JwksURI:             op.NewEndpoint("/oauth/v2/keys"),
		DeviceAuthorization: op.NewEndpoint("/oauth/v2/device_authorization"),
	}

	if endpointConfig == nil {
		return endpoints
	}
	if endpointConfig.Auth != nil {
		endpoints.Authorization = op.NewEndpointWithURL(endpointConfig.Auth.Path, endpointConfig.Auth.URL)
	}
	if endpointConfig.Token != nil {
		endpoints.Token = op.NewEndpointWithURL(endpointConfig.Token.Path, endpointConfig.Token.URL)
	}
	if endpointConfig.Introspection != nil {
		endpoints.Introspection = op.NewEndpointWithURL(endpointConfig.Introspection.Path, endpointConfig.Introspection.URL)
	}
	if endpointConfig.Userinfo != nil {
		endpoints.Userinfo = op.NewEndpointWithURL(endpointConfig.Userinfo.Path, endpointConfig.Userinfo.URL)
	}
	if endpointConfig.Revocation != nil {
		endpoints.Revocation = op.NewEndpointWithURL(endpointConfig.Revocation.Path, endpointConfig.Revocation.URL)
	}
	if endpointConfig.EndSession != nil {
		endpoints.EndSession = op.NewEndpointWithURL(endpointConfig.EndSession.Path, endpointConfig.EndSession.URL)
	}
	if endpointConfig.Keys != nil {
		endpoints.JwksURI = op.NewEndpointWithURL(endpointConfig.Keys.Path, endpointConfig.Keys.URL)
	}
	if endpointConfig.DeviceAuth != nil {
		endpoints.DeviceAuthorization = op.NewEndpointWithURL(endpointConfig.DeviceAuth.Path, endpointConfig.DeviceAuth.URL)
	}
	return endpoints
}

func (s *Server) getLogger(ctx context.Context) *slog.Logger {
	if logger, ok := logging.FromContext(ctx); ok {
		return logger
	}
	return s.fallbackLogger
}

func (s *Server) IssuerFromRequest(r *http.Request) string {
	return s.Provider().IssuerFromRequest(r)
}

func (s *Server) Health(ctx context.Context, r *op.Request[struct{}]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Health(ctx, r)
}

func (s *Server) Ready(ctx context.Context, r *op.Request[struct{}]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Ready(ctx, r)
}

func (s *Server) Discovery(ctx context.Context, r *op.Request[struct{}]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	restrictions, err := s.query.GetInstanceRestrictions(ctx)
	if err != nil {
		return nil, op.NewStatusError(oidc.ErrServerError().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("internal server error"), http.StatusInternalServerError)
	}
	allowedLanguages := restrictions.AllowedLanguages
	if len(allowedLanguages) == 0 {
		allowedLanguages = i18n.SupportedLanguages()
	}
	return op.NewResponse(s.createDiscoveryConfig(ctx, allowedLanguages)), nil
}

func (s *Server) VerifyAuthRequest(ctx context.Context, r *op.Request[oidc.AuthRequest]) (_ *op.ClientRequest[oidc.AuthRequest], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.VerifyAuthRequest(ctx, r)
}

func (s *Server) Authorize(ctx context.Context, r *op.ClientRequest[oidc.AuthRequest]) (_ *op.Redirect, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	// Use an own method to validate the id_token_hint, because in case of an error, we don't want to fail the request.
	// We just want to ignore the hint.
	userID, err := op.ValidateAuthReqIDTokenHint(ctx, r.Data.IDTokenHint, s.Provider().IDTokenHintVerifier(ctx))
	logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID()).
		OnError(err).Error("invalid id_token_hint")

	req, err := s.Provider().Storage().CreateAuthRequest(ctx, r.Data, userID)
	if err != nil {
		return op.TryErrorRedirect(ctx, r.Data, oidc.DefaultToServerError(err, "unable to save auth request"), s.Provider().Encoder(), s.Provider().Logger())
	}
	return op.NewRedirect(r.Client.LoginURL(req.GetID())), nil
}

func (s *Server) DeviceAuthorization(ctx context.Context, r *op.ClientRequest[oidc.DeviceAuthorizationRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.DeviceAuthorization(ctx, r)
}

func (s *Server) Revocation(ctx context.Context, r *op.ClientRequest[oidc.RevocationRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Revocation(ctx, r)
}

func (s *Server) EndSession(ctx context.Context, r *op.Request[oidc.EndSessionRequest]) (_ *op.Redirect, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.EndSession(ctx, r)
}

func (s *Server) createDiscoveryConfig(ctx context.Context, supportedUILocales oidc.Locales) *oidc.DiscoveryConfiguration {
	issuer := op.IssuerFromContext(ctx)
	backChannelLogoutSupported := authz.GetInstance(ctx).Features().EnableBackChannelLogout

	return &oidc.DiscoveryConfiguration{
		Issuer:                      issuer,
		AuthorizationEndpoint:       s.Endpoints().Authorization.Absolute(issuer),
		TokenEndpoint:               s.Endpoints().Token.Absolute(issuer),
		IntrospectionEndpoint:       s.Endpoints().Introspection.Absolute(issuer),
		UserinfoEndpoint:            s.Endpoints().Userinfo.Absolute(issuer),
		RevocationEndpoint:          s.Endpoints().Revocation.Absolute(issuer),
		EndSessionEndpoint:          s.Endpoints().EndSession.Absolute(issuer),
		JwksURI:                     s.Endpoints().JwksURI.Absolute(issuer),
		DeviceAuthorizationEndpoint: s.Endpoints().DeviceAuthorization.Absolute(issuer),
		ScopesSupported:             op.Scopes(s.Provider()),
		ResponseTypesSupported:      op.ResponseTypes(s.Provider()),
		ResponseModesSupported: []string{
			string(oidc.ResponseModeQuery),
			string(oidc.ResponseModeFragment),
			string(oidc.ResponseModeFormPost),
		},
		GrantTypesSupported:                                op.GrantTypes(s.Provider()),
		SubjectTypesSupported:                              op.SubjectTypes(s.Provider()),
		IDTokenSigningAlgValuesSupported:                   supportedSigningAlgs(),
		RequestObjectSigningAlgValuesSupported:             op.RequestObjectSigAlgorithms(s.Provider()),
		TokenEndpointAuthMethodsSupported:                  op.AuthMethodsTokenEndpoint(s.Provider()),
		TokenEndpointAuthSigningAlgValuesSupported:         op.TokenSigAlgorithms(s.Provider()),
		IntrospectionEndpointAuthSigningAlgValuesSupported: op.IntrospectionSigAlgorithms(s.Provider()),
		IntrospectionEndpointAuthMethodsSupported:          op.AuthMethodsIntrospectionEndpoint(s.Provider()),
		RevocationEndpointAuthSigningAlgValuesSupported:    op.RevocationSigAlgorithms(s.Provider()),
		RevocationEndpointAuthMethodsSupported:             op.AuthMethodsRevocationEndpoint(s.Provider()),
		ClaimsSupported:                                    op.SupportedClaims(s.Provider()),
		CodeChallengeMethodsSupported:                      op.CodeChallengeMethods(s.Provider()),
		UILocalesSupported:                                 supportedUILocales,
		RequestParameterSupported:                          s.Provider().RequestObjectSupported(),
		BackChannelLogoutSupported:                         backChannelLogoutSupported,
		BackChannelLogoutSessionSupported:                  backChannelLogoutSupported,
	}
}

func response(resp any, err error) (*op.Response, error) {
	if err != nil {
		return nil, err
	}
	return op.NewResponse(resp), nil
}
