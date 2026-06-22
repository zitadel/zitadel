package oidc

import (
	"encoding/json"
	"strings"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
)

// OAuth 2.0 Dynamic Client Registration error codes as defined in RFC 7591 §3.2.2.
const (
	registrationErrorInvalidRedirectURI    = "invalid_redirect_uri"
	registrationErrorInvalidClientMetadata = "invalid_client_metadata"
)

// registrationError is an OAuth 2.0 Dynamic Client Registration error response
// (RFC 7591 §3.2.2). It is always returned with HTTP status 400.
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
// interpreted. Free-form metadata (client_uri, logo_uri, contacts, …) is accepted but
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

	authMethod, regErr := registrationAuthMethodToDomain(req.TokenEndpointAuthMethod)
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
	applicationType, regErr := registrationApplicationTypeToDomain(req.ApplicationType)
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
// error. Grant/response combination problems are client metadata errors, everything else
// (redirect URI scheme and application type constraints) is reported as an invalid
// redirect URI. The internal compliance keys are not exposed to the client.
func registrationComplianceError(compliance *domain.Compliance) *registrationError {
	for _, problem := range compliance.Problems {
		switch problem {
		case "Application.OIDC.V1.GrantType", "Application.OIDC.V1.NotAllCombinationsAreAllowed":
			return newRegistrationError(registrationErrorInvalidClientMetadata, "the requested grant and response type combination is not supported")
		}
	}
	return newRegistrationError(registrationErrorInvalidRedirectURI, "one or more redirect_uris are invalid for the requested application type")
}

func registrationAuthMethodToDomain(method string) (domain.OIDCAuthMethodType, *registrationError) {
	switch method {
	case "", "client_secret_basic":
		// RFC 7591 §2: client_secret_basic is the default.
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
