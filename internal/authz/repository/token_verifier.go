package repository

import (
	"context"
)

type TokenVerifierRepository interface {
	VerifyAccessToken(ctx context.Context, tokenString, verifierClientID, projectID string) (userID string, agentID string, clientID, prefLang, resourceOwner string, err error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	VerifierClientID(ctx context.Context, appName string) (clientID, projectID string, err error)
}
