package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// registrationMaxBodyBytes bounds the size of a dynamic client registration request body.
const registrationMaxBodyBytes = 100 * 1024

// dynamicClientRegistration handles POST requests to the OAuth 2.0 Dynamic Client
// Registration endpoint (RFC 7591). The route is always mounted; it only registers
// clients when the oidc_dynamic_client_registration feature is enabled for the instance,
// otherwise it behaves as if the endpoint did not exist.
func (s *Server) dynamicClientRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !authz.GetFeatures(ctx).OIDCDynamicClientRegistration {
		http.NotFound(w, r)
		return
	}

	resourceOwner, err := s.dynamicClientRegistrationResourceOwner(ctx, r)
	if err != nil {
		// Only missing/invalid credentials are a 401; unexpected errors (e.g. a failed
		// token lookup) must not be masked as invalid_token.
		if zerrors.IsUnauthenticated(err) || zerrors.IsPermissionDenied(err) {
			s.writeRegistrationUnauthorized(ctx, w)
			return
		}
		s.writeRegistrationServerError(ctx, w, err)
		return
	}

	var req clientRegistrationRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, registrationMaxBodyBytes)).Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			s.writeRegistrationJSON(ctx, w, http.StatusRequestEntityTooLarge, newRegistrationError(registrationErrorInvalidClientMetadata, "the request body is too large"))
			return
		}
		s.writeRegistrationError(ctx, w, newRegistrationError(registrationErrorInvalidClientMetadata, "the request body could not be parsed"))
		return
	}

	app, regErr := req.toOIDCApp()
	if regErr != nil {
		s.writeRegistrationError(ctx, w, regErr)
		return
	}

	projectID, err := s.ensureDCRProject(ctx, resourceOwner)
	if err != nil {
		s.writeRegistrationServerError(ctx, w, err)
		return
	}

	registered, err := s.command.AddDynamicOIDCClient(ctx, projectID, resourceOwner, app)
	if err != nil {
		s.writeRegistrationCommandError(ctx, w, err)
		return
	}

	s.writeRegistrationJSON(ctx, w, http.StatusCreated, newClientRegistrationResponse(registered, req.ClientName, time.Now().Unix()))
}

// dynamicClientRegistrationResourceOwner authorizes the registration and returns the
// organization the client is homed in. When an access token is presented it is used as the
// RFC 7591 §3 initial access token and the client is homed in the token's organization.
// Without a token, open registration must be enabled and the client is homed in the
// instance's default organization.
//
// Note on the trust model (open question on #9810): in token mode any valid access token is
// currently accepted and the registered client is scoped to that token's organization, so a
// caller can only ever create clients in its own organization. Whether registration should
// additionally require a dedicated permission or scope is left for maintainer input, as it
// is tied to the open vs. token-gated discussion. The whole endpoint stays behind the
// feature flag and the access interceptor's rate limiting.
func (s *Server) dynamicClientRegistrationResourceOwner(ctx context.Context, r *http.Request) (string, error) {
	if token := bearerToken(r); token != "" {
		accessToken, err := s.verifyAccessToken(ctx, token)
		if err != nil {
			return "", err
		}
		return accessToken.resourceOwner, nil
	}
	if !s.dynamicClientRegistrationConfig.AllowUnauthenticated {
		return "", zerrors.ThrowUnauthenticated(nil, "OIDC-Eich8", "Errors.Token.Invalid")
	}
	return authz.GetInstance(ctx).DefaultOrganisationID(), nil
}

// ensureDCRProject returns the dedicated project that holds dynamically registered clients
// for the organization, creating it on first use. The common case (the project already
// exists) is served from the projection; the creation and the concurrent-creation race are
// delegated to the command, which resolves them strongly consistently from the eventstore.
func (s *Server) ensureDCRProject(ctx context.Context, resourceOwner string) (string, error) {
	projectID, err := s.dcrProjectIDFromProjection(ctx, resourceOwner)
	if err != nil || projectID != "" {
		return projectID, err
	}
	return s.command.EnsureDCRProject(ctx, resourceOwner)
}

func (s *Server) dcrProjectIDFromProjection(ctx context.Context, resourceOwner string) (string, error) {
	nameQuery, err := query.NewProjectNameSearchQuery(query.TextEquals, command.DCRProjectName)
	if err != nil {
		return "", err
	}
	ownerQuery, err := query.NewProjectResourceOwnerSearchQuery(resourceOwner)
	if err != nil {
		return "", err
	}
	// No permission check: access to the registration endpoint is the authorization
	// boundary for dynamic client registration.
	projects, err := s.query.SearchProjects(ctx, &query.ProjectSearchQueries{Queries: []query.SearchQuery{nameQuery, ownerQuery}}, nil)
	if err != nil {
		return "", err
	}
	if len(projects.Projects) == 0 {
		return "", nil
	}
	return projects.Projects[0].ID, nil
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	scheme, token, found := strings.Cut(auth, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}
	return strings.TrimSpace(token)
}

func (s *Server) writeRegistrationJSON(ctx context.Context, w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.getLogger(ctx).ErrorContext(ctx, "dynamic client registration: encode response", "err", err)
	}
}

func (s *Server) writeRegistrationError(ctx context.Context, w http.ResponseWriter, regErr *registrationError) {
	s.writeRegistrationJSON(ctx, w, http.StatusBadRequest, regErr)
}

func (s *Server) writeRegistrationUnauthorized(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
	s.writeRegistrationJSON(ctx, w, http.StatusUnauthorized, newRegistrationError("invalid_token", "a valid access token is required to register a client"))
}

func (s *Server) writeRegistrationServerError(ctx context.Context, w http.ResponseWriter, err error) {
	s.getLogger(ctx).ErrorContext(ctx, "dynamic client registration", "err", err)
	s.writeRegistrationJSON(ctx, w, http.StatusInternalServerError, newRegistrationError("server_error", "the client could not be registered"))
}

func (s *Server) writeRegistrationCommandError(ctx context.Context, w http.ResponseWriter, err error) {
	if zerrors.IsErrorInvalidArgument(err) {
		s.writeRegistrationError(ctx, w, newRegistrationError(registrationErrorInvalidClientMetadata, "the requested client metadata is not supported"))
		return
	}
	s.writeRegistrationServerError(ctx, w, err)
}

// OAuth 2.0 Dynamic Client Registration error codes as defined in RFC 7591 §3.2.2.
const (
	registrationErrorInvalidRedirectURI    = "invalid_redirect_uri"
	registrationErrorInvalidClientMetadata = "invalid_client_metadata"
)

// registrationError is an OAuth 2.0 Dynamic Client Registration error response
// (RFC 7591 §3.2.2). Invalid client metadata is returned with HTTP status 400; the same
// shape is reused for the 401 (invalid_token) and 500 (server_error) responses.
type registrationError struct {
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func (e *registrationError) Error() string {
	if e.ErrorDescription == "" {
		return e.ErrorType
	}
	return e.ErrorType + ": " + e.ErrorDescription
}

func newRegistrationError(errorType, description string) *registrationError {
	return &registrationError{ErrorType: errorType, ErrorDescription: description}
}

// clientRegistrationRequest holds the client metadata of an OAuth 2.0 Dynamic Client
// Registration request (RFC 7591 §2, extended by OpenID Connect Dynamic Client
// Registration 1.0 §2). Only the members relevant for a ZITADEL OIDC application are
// interpreted. Free-form metadata (client_uri, logo_uri, contacts, etc.) is accepted but
// not persisted in this version. Members that imply unsupported behaviour (jwks/jwks_uri,
// private_key_jwt) are rejected explicitly.
type clientRegistrationRequest struct {
	RedirectURIs            []string        `json:"redirect_uris,omitempty"`
	ResponseTypes           []string        `json:"response_types,omitempty"`
	GrantTypes              []string        `json:"grant_types,omitempty"`
	ApplicationType         string          `json:"application_type,omitempty"`
	ClientName              string          `json:"client_name,omitempty"`
	TokenEndpointAuthMethod string          `json:"token_endpoint_auth_method,omitempty"`
	PostLogoutRedirectURIs  []string        `json:"post_logout_redirect_uris,omitempty"`
	JWKsURI                 string          `json:"jwks_uri,omitempty"`
	JWKs                    json.RawMessage `json:"jwks,omitempty"`
}

// clientRegistrationResponse is the client information response of a successful
// registration (RFC 7591 §3.2.1).
type clientRegistrationResponse struct {
	ClientID              string `json:"client_id"`
	ClientSecret          string `json:"client_secret,omitempty"`
	ClientIDIssuedAt      int64  `json:"client_id_issued_at,omitempty"`
	ClientSecretExpiresAt *int64 `json:"client_secret_expires_at,omitempty"`

	RedirectURIs            []string `json:"redirect_uris,omitempty"`
	ResponseTypes           []string `json:"response_types,omitempty"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	ApplicationType         string   `json:"application_type,omitempty"`
	ClientName              string   `json:"client_name,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	PostLogoutRedirectURIs  []string `json:"post_logout_redirect_uris,omitempty"`
}

// toOIDCApp maps the registration request to a domain OIDC application, applying the RFC
// 7591 defaults and rejecting unsupported metadata. The returned application is not yet
// persisted; the application name and the client credentials are assigned by the command.
func (req *clientRegistrationRequest) toOIDCApp() (*domain.OIDCApp, *registrationError) {
	if req.JWKsURI != "" || len(req.JWKs) > 0 {
		return nil, newRegistrationError(registrationErrorInvalidClientMetadata, "jwks and jwks_uri are not supported")
	}

	applicationType, regErr := registrationApplicationTypeToDomain(req.ApplicationType)
	if regErr != nil {
		return nil, regErr
	}
	authMethod, regErr := registrationAuthMethodToDomain(req.TokenEndpointAuthMethod, applicationType)
	if regErr != nil {
		return nil, regErr
	}
	grantTypes, regErr := registrationGrantTypesToDomain(req.GrantTypes)
	if regErr != nil {
		return nil, regErr
	}
	responseTypes, regErr := registrationResponseTypesToDomain(req.ResponseTypes)
	if regErr != nil {
		return nil, regErr
	}

	redirectURIs := trimSpaceSlice(req.RedirectURIs)
	if len(redirectURIs) == 0 {
		return nil, newRegistrationError(registrationErrorInvalidRedirectURI, "at least one redirect_uri is required")
	}

	app := &domain.OIDCApp{
		AppName:                strings.TrimSpace(req.ClientName),
		RedirectUris:           redirectURIs,
		ResponseTypes:          responseTypes,
		GrantTypes:             grantTypes,
		ApplicationType:        gu.Ptr(applicationType),
		AuthMethodType:         gu.Ptr(authMethod),
		PostLogoutRedirectUris: trimSpaceSlice(req.PostLogoutRedirectURIs),
		OIDCVersion:            gu.Ptr(domain.OIDCVersionV1),
		AccessTokenType:        gu.Ptr(domain.OIDCTokenTypeBearer),
	}

	if !app.IsValid() {
		return nil, newRegistrationError(registrationErrorInvalidClientMetadata, "the requested grant and response type combination is not supported")
	}
	if compliance := domain.GetOIDCV1Compliance(app.ApplicationType, app.GrantTypes, app.AuthMethodType, app.RedirectUris); compliance.NoneCompliant {
		return nil, registrationComplianceError(compliance)
	}
	return app, nil
}

// registrationComplianceError maps a domain compliance failure to the matching RFC 7591
// error. Auth-method and grant/response combination problems are client metadata errors;
// redirect URI scheme problems are reported as an invalid redirect URI. The internal
// compliance keys (see internal/domain/application_oidc.go) are not exposed to the client.
func registrationComplianceError(compliance *domain.Compliance) *registrationError {
	for _, problem := range compliance.Problems {
		if strings.Contains(problem, "AuthMethodType") || strings.Contains(problem, "GrantType") || strings.Contains(problem, "Combinations") {
			return newRegistrationError(registrationErrorInvalidClientMetadata, "the requested client metadata is not supported")
		}
	}
	return newRegistrationError(registrationErrorInvalidRedirectURI, "one or more redirect_uris are invalid for the requested application type")
}

func registrationAuthMethodToDomain(method string, applicationType domain.OIDCApplicationType) (domain.OIDCAuthMethodType, *registrationError) {
	switch method {
	case "":
		// Default the auth method to the application type. Native and user-agent
		// applications must be public (none); web applications default to
		// client_secret_basic per RFC 7591 §2.
		if applicationType == domain.OIDCApplicationTypeNative || applicationType == domain.OIDCApplicationTypeUserAgent {
			return domain.OIDCAuthMethodTypeNone, nil
		}
		return domain.OIDCAuthMethodTypeBasic, nil
	case "client_secret_basic":
		return domain.OIDCAuthMethodTypeBasic, nil
	case "client_secret_post":
		return domain.OIDCAuthMethodTypePost, nil
	case "none":
		return domain.OIDCAuthMethodTypeNone, nil
	default:
		return 0, newRegistrationError(registrationErrorInvalidClientMetadata, "token_endpoint_auth_method "+method+" is not supported")
	}
}

func registrationGrantTypesToDomain(grantTypes []string) ([]domain.OIDCGrantType, *registrationError) {
	if len(grantTypes) == 0 {
		// RFC 7591 §2: authorization_code is the default.
		return []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode}, nil
	}
	mapped := make([]domain.OIDCGrantType, 0, len(grantTypes))
	for _, grantType := range grantTypes {
		switch grantType {
		case "authorization_code":
			mapped = append(mapped, domain.OIDCGrantTypeAuthorizationCode)
		case "refresh_token":
			mapped = append(mapped, domain.OIDCGrantTypeRefreshToken)
		case "implicit":
			mapped = append(mapped, domain.OIDCGrantTypeImplicit)
		case "urn:ietf:params:oauth:grant-type:device_code":
			mapped = append(mapped, domain.OIDCGrantTypeDeviceCode)
		case "urn:ietf:params:oauth:grant-type:token-exchange":
			mapped = append(mapped, domain.OIDCGrantTypeTokenExchange)
		default:
			return nil, newRegistrationError(registrationErrorInvalidClientMetadata, "grant_type "+grantType+" is not supported")
		}
	}
	return mapped, nil
}

func registrationResponseTypesToDomain(responseTypes []string) ([]domain.OIDCResponseType, *registrationError) {
	if len(responseTypes) == 0 {
		// RFC 7591 §2: code is the default.
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}, nil
	}
	mapped := make([]domain.OIDCResponseType, 0, len(responseTypes))
	for _, responseType := range responseTypes {
		switch normalizeResponseType(responseType) {
		case "code":
			mapped = append(mapped, domain.OIDCResponseTypeCode)
		case "id_token":
			mapped = append(mapped, domain.OIDCResponseTypeIDToken)
		case "id_token token":
			mapped = append(mapped, domain.OIDCResponseTypeIDTokenToken)
		default:
			return nil, newRegistrationError(registrationErrorInvalidClientMetadata, "response_type "+responseType+" is not supported")
		}
	}
	return mapped, nil
}

// normalizeResponseType sorts the space separated values of a response_type so that the
// order of "token" and "id_token" does not matter.
func normalizeResponseType(responseType string) string {
	fields := strings.Fields(responseType)
	hasToken, hasIDToken, hasCode := false, false, false
	for _, field := range fields {
		switch field {
		case "token":
			hasToken = true
		case "id_token":
			hasIDToken = true
		case "code":
			hasCode = true
		}
	}
	switch {
	case hasCode && !hasToken && !hasIDToken:
		return "code"
	case hasIDToken && hasToken && !hasCode:
		return "id_token token"
	case hasIDToken && !hasToken && !hasCode:
		return "id_token"
	default:
		return strings.Join(fields, " ")
	}
}

func registrationApplicationTypeToDomain(applicationType string) (domain.OIDCApplicationType, *registrationError) {
	switch applicationType {
	case "", "web":
		// OpenID Connect Dynamic Client Registration 1.0 §2: web is the default.
		return domain.OIDCApplicationTypeWeb, nil
	case "native":
		return domain.OIDCApplicationTypeNative, nil
	default:
		return 0, newRegistrationError(registrationErrorInvalidClientMetadata, "application_type "+applicationType+" is not supported")
	}
}

// newClientRegistrationResponse builds the RFC 7591 §3.2.1 response from the registered
// application. The client name is echoed from the request, as the human-readable name is
// not persisted as-is in this version.
func newClientRegistrationResponse(app *domain.OIDCApp, clientName string, issuedAt int64) *clientRegistrationResponse {
	resp := &clientRegistrationResponse{
		ClientID:                app.ClientID,
		ClientSecret:            app.ClientSecretString,
		ClientIDIssuedAt:        issuedAt,
		RedirectURIs:            app.RedirectUris,
		ResponseTypes:           responseTypesToRegistration(app.ResponseTypes),
		GrantTypes:              grantTypesToRegistration(app.GrantTypes),
		ApplicationType:         applicationTypeToRegistration(gu.Value(app.ApplicationType)),
		ClientName:              strings.TrimSpace(clientName),
		TokenEndpointAuthMethod: authMethodToRegistration(gu.Value(app.AuthMethodType)),
		PostLogoutRedirectURIs:  app.PostLogoutRedirectUris,
	}
	if app.ClientSecretString != "" {
		// RFC 7591 §3.2.1: client_secret_expires_at is REQUIRED if client_secret is
		// issued. 0 means the secret never expires.
		resp.ClientSecretExpiresAt = gu.Ptr(int64(0))
	}
	return resp
}

func authMethodToRegistration(authMethod domain.OIDCAuthMethodType) string {
	switch authMethod {
	case domain.OIDCAuthMethodTypeBasic:
		return "client_secret_basic"
	case domain.OIDCAuthMethodTypePost:
		return "client_secret_post"
	case domain.OIDCAuthMethodTypeNone:
		return "none"
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		return "private_key_jwt"
	default:
		return ""
	}
}

func grantTypesToRegistration(grantTypes []domain.OIDCGrantType) []string {
	mapped := make([]string, 0, len(grantTypes))
	for _, grantType := range grantTypes {
		switch grantType {
		case domain.OIDCGrantTypeAuthorizationCode:
			mapped = append(mapped, "authorization_code")
		case domain.OIDCGrantTypeRefreshToken:
			mapped = append(mapped, "refresh_token")
		case domain.OIDCGrantTypeImplicit:
			mapped = append(mapped, "implicit")
		case domain.OIDCGrantTypeDeviceCode:
			mapped = append(mapped, "urn:ietf:params:oauth:grant-type:device_code")
		case domain.OIDCGrantTypeTokenExchange:
			mapped = append(mapped, "urn:ietf:params:oauth:grant-type:token-exchange")
		}
	}
	return mapped
}

func responseTypesToRegistration(responseTypes []domain.OIDCResponseType) []string {
	mapped := make([]string, 0, len(responseTypes))
	for _, responseType := range responseTypes {
		switch responseType {
		case domain.OIDCResponseTypeCode:
			mapped = append(mapped, "code")
		case domain.OIDCResponseTypeIDToken:
			mapped = append(mapped, "id_token")
		case domain.OIDCResponseTypeIDTokenToken:
			mapped = append(mapped, "id_token token")
		case domain.OIDCResponseTypeUnspecified:
			// not exposed through dynamic client registration
		}
	}
	return mapped
}

func applicationTypeToRegistration(applicationType domain.OIDCApplicationType) string {
	switch applicationType {
	case domain.OIDCApplicationTypeWeb:
		return "web"
	case domain.OIDCApplicationTypeNative:
		return "native"
	case domain.OIDCApplicationTypeUserAgent:
		return "user_agent"
	default:
		return ""
	}
}

func trimSpaceSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	trimmed := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			trimmed = append(trimmed, value)
		}
	}
	if len(trimmed) == 0 {
		return nil
	}
	return trimmed
}
