package oidc

import (
	"context"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Server struct {
	http.Handler
	*op.LegacyServer
	signingKeyAlgorithm string
}

func endpoints(endpointConfig *EndpointConfig) op.Endpoints {
	// some defaults. The new Server will disable enpoints that are nil.
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
	defer func() { span.EndWithError(err) }()

	return op.NewResponse(s.createDiscoveryConfig(ctx)), nil
}

func (s *Server) Keys(ctx context.Context, r *op.Request[struct{}]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Keys(ctx, r)
}

func (s *Server) VerifyAuthRequest(ctx context.Context, r *op.Request[oidc.AuthRequest]) (_ *op.ClientRequest[oidc.AuthRequest], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.VerifyAuthRequest(ctx, r)
}

func (s *Server) Authorize(ctx context.Context, r *op.ClientRequest[oidc.AuthRequest]) (_ *op.Redirect, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Authorize(ctx, r)
}

func (s *Server) DeviceAuthorization(ctx context.Context, r *op.ClientRequest[oidc.DeviceAuthorizationRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.DeviceAuthorization(ctx, r)
}

func (s *Server) VerifyClient(ctx context.Context, r *op.Request[op.ClientCredentials]) (_ op.Client, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.VerifyClient(ctx, r)
}

func (s *Server) CodeExchange(ctx context.Context, r *op.ClientRequest[oidc.AccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.CodeExchange(ctx, r)
}

func (s *Server) RefreshToken(ctx context.Context, r *op.ClientRequest[oidc.RefreshTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.RefreshToken(ctx, r)
}

func (s *Server) JWTProfile(ctx context.Context, r *op.Request[oidc.JWTProfileGrantRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.JWTProfile(ctx, r)
}

func (s *Server) TokenExchange(ctx context.Context, r *op.ClientRequest[oidc.TokenExchangeRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.TokenExchange(ctx, r)
}

func (s *Server) ClientCredentialsExchange(ctx context.Context, r *op.ClientRequest[oidc.ClientCredentialsRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.ClientCredentialsExchange(ctx, r)
}

func (s *Server) DeviceToken(ctx context.Context, r *op.ClientRequest[oidc.DeviceAccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.DeviceToken(ctx, r)
}

func (s *Server) Introspect(ctx context.Context, r *op.Request[op.IntrospectionRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.Introspect(ctx, r)
}

func (s *Server) UserInfo(ctx context.Context, r *op.Request[oidc.UserInfoRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return s.LegacyServer.UserInfo(ctx, r)
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

func (s *Server) createDiscoveryConfig(ctx context.Context) *oidc.DiscoveryConfiguration {
	issuer := op.IssuerFromContext(ctx)
	return &oidc.DiscoveryConfiguration{
		Issuer:                                     issuer,
		AuthorizationEndpoint:                      s.Endpoints().Authorization.Absolute(issuer),
		TokenEndpoint:                              s.Endpoints().Token.Absolute(issuer),
		IntrospectionEndpoint:                      s.Endpoints().Introspection.Absolute(issuer),
		UserinfoEndpoint:                           s.Endpoints().Userinfo.Absolute(issuer),
		RevocationEndpoint:                         s.Endpoints().Revocation.Absolute(issuer),
		EndSessionEndpoint:                         s.Endpoints().EndSession.Absolute(issuer),
		JwksURI:                                    s.Endpoints().JwksURI.Absolute(issuer),
		DeviceAuthorizationEndpoint:                s.Endpoints().DeviceAuthorization.Absolute(issuer),
		ScopesSupported:                            op.Scopes(s.Provider()),
		ResponseTypesSupported:                     op.ResponseTypes(s.Provider()),
		GrantTypesSupported:                        op.GrantTypes(s.Provider()),
		SubjectTypesSupported:                      op.SubjectTypes(s.Provider()),
		IDTokenSigningAlgValuesSupported:           []string{s.signingKeyAlgorithm},
		RequestObjectSigningAlgValuesSupported:     op.RequestObjectSigAlgorithms(s.Provider()),
		TokenEndpointAuthMethodsSupported:          op.AuthMethodsTokenEndpoint(s.Provider()),
		TokenEndpointAuthSigningAlgValuesSupported: op.TokenSigAlgorithms(s.Provider()),
		IntrospectionEndpointAuthSigningAlgValuesSupported: op.IntrospectionSigAlgorithms(s.Provider()),
		IntrospectionEndpointAuthMethodsSupported:          op.AuthMethodsIntrospectionEndpoint(s.Provider()),
		RevocationEndpointAuthSigningAlgValuesSupported:    op.RevocationSigAlgorithms(s.Provider()),
		RevocationEndpointAuthMethodsSupported:             op.AuthMethodsRevocationEndpoint(s.Provider()),
		ClaimsSupported:                                    op.SupportedClaims(s.Provider()),
		CodeChallengeMethodsSupported:                      op.CodeChallengeMethods(s.Provider()),
		UILocalesSupported:                                 s.Provider().SupportedUILocales(),
		RequestParameterSupported:                          s.Provider().RequestObjectSupported(),
	}
}
