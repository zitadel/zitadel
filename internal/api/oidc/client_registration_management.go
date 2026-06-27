package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// getDynamicClientRegistration handles GET requests to the client configuration endpoint
// (RFC 7592 §2.1) and returns the current registration of the client. The read does not
// rotate the registration access token, so the response omits it; the client keeps the token
// it authenticated with.
func (s *Server) getDynamicClientRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, _, ok := s.authorizeClientManagement(w, r)
	if !ok {
		return
	}

	resp := clientRegistrationResponseFromClient(client)
	resp.RegistrationClientURI = s.registrationClientURI(ctx, client.ClientID)
	s.writeRegistrationJSON(ctx, w, http.StatusOK, resp)
}

// updateDynamicClientRegistration handles PUT requests to the client configuration endpoint
// (RFC 7592 §2.2). It replaces the client metadata with the submitted values and rotates the
// registration access token, returning the updated registration with the new token.
func (s *Server) updateDynamicClientRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, binding, ok := s.authorizeClientManagement(w, r)
	if !ok {
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
	app.AggregateID = client.ProjectID
	app.AppID = client.AppID

	updated, err := s.command.UpdateDynamicOIDCClient(ctx, app, binding.orgID)
	if err != nil {
		s.writeManagementCommandError(ctx, w, err)
		return
	}

	// The update response is a client information response (RFC 7592 §3); client_id_issued_at
	// reflects the original registration, which is not echoed back here.
	resp := newClientRegistrationResponse(updated, req.ClientName, 0)
	resp.RegistrationAccessToken, err = s.registrationAccessToken(client.ClientID, binding.orgID, updated.RegistrationAccessToken)
	if err != nil {
		s.writeRegistrationServerError(ctx, w, err)
		return
	}
	resp.RegistrationClientURI = s.registrationClientURI(ctx, client.ClientID)
	s.writeRegistrationJSON(ctx, w, http.StatusOK, resp)
}

// deleteDynamicClientRegistration handles DELETE requests to the client configuration
// endpoint (RFC 7592 §2.3). It removes the client and answers with 204 No Content.
func (s *Server) deleteDynamicClientRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, binding, ok := s.authorizeClientManagement(w, r)
	if !ok {
		return
	}

	if _, err := s.command.RemoveDynamicOIDCClient(ctx, client.ProjectID, client.AppID, binding.orgID); err != nil {
		s.writeManagementCommandError(ctx, w, err)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusNoContent)
}

// registrationAccessTokenBinding holds the values recovered from a registration access token.
type registrationAccessTokenBinding struct {
	clientID string
	orgID    string
	secret   string
}

// authorizeClientManagement performs the authorization shared by the RFC 7592 management
// endpoints. It requires the feature to be enabled, a registration access token that decrypts
// and is bound to the client_id in the path, an existing client, and a secret that matches the
// stored token hash. On any failure it writes the appropriate response (404 when the feature
// is off or the client does not exist, 401 for token problems) and returns ok=false.
func (s *Server) authorizeClientManagement(w http.ResponseWriter, r *http.Request) (*query.OIDCClient, *registrationAccessTokenBinding, bool) {
	ctx := r.Context()
	if !authz.GetFeatures(ctx).OIDCDynamicClientRegistration {
		http.NotFound(w, r)
		return nil, nil, false
	}

	clientID := chi.URLParam(r, "client_id")
	token := bearerToken(r)
	if token == "" {
		s.writeRegistrationUnauthorized(ctx, w)
		return nil, nil, false
	}
	binding, err := s.parseRegistrationAccessToken(token)
	if err != nil || binding.clientID != clientID {
		s.writeRegistrationUnauthorized(ctx, w)
		return nil, nil, false
	}

	client, err := s.query.ActiveOIDCClientByID(ctx, clientID, false)
	if err != nil {
		s.writeRegistrationNotFound(ctx, w)
		return nil, nil, false
	}
	if !s.verifyRegistrationAccessToken(ctx, client, binding.orgID, binding.secret) {
		s.writeRegistrationUnauthorized(ctx, w)
		return nil, nil, false
	}
	return client, binding, true
}

// verifyRegistrationAccessToken checks the presented registration access token secret against
// the client's stored token hash. The hash is served from the projection for an O(1) read on
// the common path; on a miss (most importantly a token that was just rotated and is not
// projected yet) it falls back to a strongly consistent read from the eventstore, so a rotated
// token is accepted immediately.
func (s *Server) verifyRegistrationAccessToken(ctx context.Context, client *query.OIDCClient, orgID, secret string) bool {
	if client.RegistrationTokenHash != "" {
		if _, err := s.hasher.Verify(client.RegistrationTokenHash, secret); err == nil {
			return true
		}
	}
	return s.command.VerifyDynamicClientRegistrationToken(ctx, client.ProjectID, client.AppID, orgID, secret) == nil
}

// registrationAccessToken builds the opaque registration access token (RFC 7592 §3) handed to
// a client on registration and rotation. It binds the plain secret to the client and its
// organization through authenticated encryption with the instance key (the same mechanism as
// refresh tokens), so the management endpoints can recover the binding before checking the
// secret against the stored hash.
func (s *Server) registrationAccessToken(clientID, orgID, secret string) (string, error) {
	return s.encAlg.EncryptToken(strings.Join([]string{clientID, orgID, secret}, ":"))
}

// parseRegistrationAccessToken decrypts a registration access token back into its binding.
func (s *Server) parseRegistrationAccessToken(token string) (*registrationAccessTokenBinding, error) {
	plain, err := s.encAlg.DecryptToken(token)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(plain, ":", 3)
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}
	return &registrationAccessTokenBinding{clientID: parts[0], orgID: parts[1], secret: parts[2]}, nil
}

// registrationClientURI builds the RFC 7592 §3 client configuration endpoint URL of a client.
func (s *Server) registrationClientURI(ctx context.Context, clientID string) string {
	return op.IssuerFromContext(ctx) + s.registrationEndpoint.Relative() + "/" + clientID
}

func (s *Server) writeRegistrationNotFound(ctx context.Context, w http.ResponseWriter) {
	s.writeRegistrationJSON(ctx, w, http.StatusNotFound, newRegistrationError("invalid_client", "the client does not exist or the registration access token is no longer valid"))
}

// writeManagementCommandError maps a command error from the update and delete endpoints. A
// client that no longer exists (a concurrent or retried delete, or a delete racing an update)
// is reported as 404, consistent with a read of a missing client; other errors fall back to
// the shared registration error mapping.
func (s *Server) writeManagementCommandError(ctx context.Context, w http.ResponseWriter, err error) {
	if zerrors.IsNotFound(err) {
		s.writeRegistrationNotFound(ctx, w)
		return
	}
	s.writeRegistrationCommandError(ctx, w, err)
}

// clientRegistrationResponseFromClient maps a persisted OIDC client to an RFC 7592 read
// response. The client secret is not returned because only its hash is stored; the original
// client name and issuance time are likewise not persisted in this version.
func clientRegistrationResponseFromClient(client *query.OIDCClient) *clientRegistrationResponse {
	return &clientRegistrationResponse{
		ClientID:                client.ClientID,
		RedirectURIs:            client.RedirectURIs,
		ResponseTypes:           responseTypesToRegistration(client.ResponseTypes),
		GrantTypes:              grantTypesToRegistration(client.GrantTypes),
		ApplicationType:         applicationTypeToRegistration(client.ApplicationType),
		TokenEndpointAuthMethod: authMethodToRegistration(client.AuthMethodType),
		PostLogoutRedirectURIs:  client.PostLogoutRedirectURIs,
	}
}
