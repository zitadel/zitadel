package oidc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gooidc "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type authRequestByCodeRepository struct {
	repository.Repository
	authRequest *domain.AuthRequest
}

func (r *authRequestByCodeRepository) AuthRequestByCode(context.Context, string) (*domain.AuthRequest, error) {
	return r.authRequest, nil
}

func TestServer_codeExchangeV1_PKCE(t *testing.T) {
	const (
		clientID    = "client"
		redirectURI = "https://client.example/callback"
	)
	server := &Server{
		repo: &authRequestByCodeRepository{
			authRequest: &domain.AuthRequest{
				ApplicationID: clientID,
				CallbackURI:   redirectURI,
				UserID:        "user",
				CodeUserID:    "user",
				PossibleSteps: []domain.NextStep{&domain.RedirectToCallbackStep{}},
				Request: domain.NewAuthRequestOIDC(
					[]string{gooidc.ScopeOpenID},
					domain.OIDCResponseTypeCode,
					domain.OIDCResponseModeQuery,
					"",
					nil,
				),
			},
		},
	}
	client := &Client{client: &query.OIDCClient{
		ClientID:       clientID,
		AuthMethodType: domain.OIDCAuthMethodTypeNone,
	}}
	// A different redirect URI provides a deterministic downstream error when
	// the missing PKCE verifier is not rejected first.
	request := &gooidc.AccessTokenRequest{RedirectURI: "https://client.example/other"}

	// A missing verifier must stop the exchange before session creation.
	_, err := server.codeExchangeV1(context.Background(), client, request, "code")
	var oidcErr *gooidc.Error
	require.ErrorAs(t, err, &oidcErr)
	assert.Equal(t, gooidc.InvalidRequest, oidcErr.ErrorType)
}

func TestCodeExchangeComplianceChecker_PKCE(t *testing.T) {
	const (
		clientID     = "client"
		redirectURI  = "https://client.example/callback"
		codeVerifier = "0123456789012345678901234567890123456789012"
	)
	validChallenge := &domain.OIDCCodeChallenge{
		Challenge: gooidc.NewSHACodeChallenge(codeVerifier),
		Method:    domain.CodeChallengeMethodS256,
	}
	tests := []struct {
		name          string
		authMethod    domain.OIDCAuthMethodType
		codeChallenge *domain.OIDCCodeChallenge
		codeVerifier  string
		wantErr       bool
	}{
		{
			name:       "public client without verifier or stored challenge",
			authMethod: domain.OIDCAuthMethodTypeNone,
			wantErr:    true,
		},
		{
			name:          "public client with valid S256 verifier",
			authMethod:    domain.OIDCAuthMethodTypeNone,
			codeChallenge: validChallenge,
			codeVerifier:  codeVerifier,
		},
		{
			name:          "public client with stored challenge and missing verifier",
			authMethod:    domain.OIDCAuthMethodTypeNone,
			codeChallenge: validChallenge,
			wantErr:       true,
		},
		{
			name:         "public client with verifier and missing stored challenge",
			authMethod:   domain.OIDCAuthMethodTypeNone,
			codeVerifier: codeVerifier,
			wantErr:      true,
		},
		{
			name:       "confidential client without PKCE",
			authMethod: domain.OIDCAuthMethodTypeBasic,
		},
		{
			name:          "confidential client with stored challenge and missing verifier",
			authMethod:    domain.OIDCAuthMethodTypeBasic,
			codeChallenge: validChallenge,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{client: &query.OIDCClient{
				ClientID:       clientID,
				AuthMethodType: tt.authMethod,
			}}
			request := &gooidc.AccessTokenRequest{
				RedirectURI:  redirectURI,
				CodeVerifier: tt.codeVerifier,
			}
			authRequest := &command.AuthRequestWriteModel{
				ClientID:         clientID,
				RedirectURI:      redirectURI,
				ResponseType:     domain.OIDCResponseTypeCode,
				CodeChallenge:    tt.codeChallenge,
				SessionID:        "session",
				AuthRequestState: domain.AuthRequestStateCodeExchanged,
			}

			err := codeExchangeComplianceChecker(client, request)(context.Background(), authRequest)
			if tt.wantErr {
				var oidcErr *gooidc.Error
				require.ErrorAs(t, err, &oidcErr)
				assert.Equal(t, gooidc.InvalidRequest, oidcErr.ErrorType)
				return
			}
			require.NoError(t, err)
		})
	}
}

// codeRepo returns a fixed auth request by code; any other repository call panics.
type codeRepo struct {
	repository.Repository
	authReq *domain.AuthRequest
}

func (r *codeRepo) AuthRequestByCode(context.Context, string) (*domain.AuthRequest, error) {
	return r.authReq, nil
}

func TestServer_codeExchangeV1_rejected(t *testing.T) {
	tests := []struct {
		name    string
		authReq *domain.AuthRequest
		wantMsg string
	}{
		{
			name: "user switched after issuance",
			authReq: &domain.AuthRequest{
				ApplicationID: "clientID",
				UserID:        "userB",
				CodeUserID:    "userA",
				PossibleSteps: []domain.NextStep{&domain.RedirectToCallbackStep{}},
				Request:       &domain.AuthRequestOIDC{},
			},
			wantMsg: "user changed after the code was issued",
		},
		{
			name: "user not authenticated",
			authReq: &domain.AuthRequest{
				ApplicationID: "clientID",
				UserID:        "userA",
				CodeUserID:    "userA",
				PossibleSteps: []domain.NextStep{&domain.LoginStep{}},
				Request:       &domain.AuthRequestOIDC{},
			},
			wantMsg: "user is not authenticated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{repo: &codeRepo{authReq: tt.authReq}}
			client := &Client{client: &query.OIDCClient{ClientID: "clientID"}}

			_, err := s.codeExchangeV1(context.Background(), client, &gooidc.AccessTokenRequest{}, "code")

			var oidcErr *gooidc.Error
			if assert.True(t, errors.As(err, &oidcErr), "got %v", err) {
				assert.Equal(t, gooidc.InvalidGrant, oidcErr.ErrorType)
				assert.Equal(t, tt.wantMsg, oidcErr.Description)
			}
		})
	}
}
