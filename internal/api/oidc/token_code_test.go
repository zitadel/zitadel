package oidc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

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

			_, err := s.codeExchangeV1(context.Background(), client, &oidc.AccessTokenRequest{}, "code")

			var oidcErr *oidc.Error
			if assert.True(t, errors.As(err, &oidcErr), "got %v", err) {
				assert.Equal(t, oidc.InvalidGrant, oidcErr.ErrorType)
				assert.Equal(t, tt.wantMsg, oidcErr.Description)
			}
		})
	}
}
