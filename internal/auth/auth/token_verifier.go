package auth

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
)

type TokenVerifier struct {
}

func Start() (v *TokenVerifier) {
	return new(TokenVerifier)
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (string, string, string, error) {
	return "", "", "", nil
}

func (v *TokenVerifier) ResolveGrants(ctx context.Context, userID, orgID string) ([]*auth.Grant, error) {
	return nil, nil
}

func (v *TokenVerifier) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return "", nil
}
