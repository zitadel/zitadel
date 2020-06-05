package auth

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
)

const (
	mgmtName = "Management-API"
)

type TokenVerifier struct {
	mgmtID    string
	authZRepo *authz_repo.EsRepository
}

func Start(authZRepo *authz_repo.EsRepository) (v *TokenVerifier) {
	return &TokenVerifier{authZRepo: authZRepo}
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (string, string, string, error) {
	userID, clientID, agentID, err := v.authZRepo.VerifyAccessToken(ctx, token, mgmtName, v.mgmtID)
	if clientID != "" {
		v.mgmtID = clientID
	}
	return userID, clientID, agentID, err
}

func (v *TokenVerifier) ResolveGrant(ctx context.Context) (*auth.Grant, error) {
	return v.authZRepo.ResolveGrants(ctx)
}

func (v *TokenVerifier) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return v.authZRepo.ProjectIDByClientID(ctx, clientID)
}
