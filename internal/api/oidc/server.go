package oidc

import (
	"context"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/exp/slog"
)

type Server struct {
	http.Handler
	*op.LegacyServer
	storage *OPStorage
}

func newServer(provider op.OpenIDProvider, storage *OPStorage, endpointConfig *EndpointConfig, logger *slog.Logger) *Server {
	server := &Server{
		LegacyServer: op.NewLegacyServer(provider, endpoints(endpointConfig)),
	}
	server.Handler = op.RegisterLegacyServer(server, op.WithFallbackLogger(logger))
	return server
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

func (s *Server) AuthCallbackURL() func(context.Context, string) string {
	return op.AuthCallbackURL(s.Provider())
}

// CodeExchange is an example how we can override / implement the handler methods.
func (s *Server) CodeExchange(ctx context.Context, req *op.ClientRequest[oidc.AccessTokenRequest]) (*op.Response, error) {
	return s.LegacyServer.CodeExchange(ctx, req)
}
