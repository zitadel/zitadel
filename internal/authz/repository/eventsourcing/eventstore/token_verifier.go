package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type TokenVerifierRepo struct {
	TokenVerificationKey [32]byte
	IamID                string
	IamEvents            *iam_event.IamEventstore
	ProjectEvents        *proj_event.ProjectEventstore
	View                 *view.View
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, appName, appID string) (userID string, clientID string, agentID string, err error) {
	clientID, err = repo.verifierClientID(ctx, appName, appID)
	if err != nil {
		return "", "", "", caos_errs.ThrowPermissionDenied(nil, "APP-ptTIF2", "invalid token")
	}
	//TODO: use real key
	tokenID, err := crypto.DecryptAESString(tokenString, string(repo.TokenVerificationKey[:32]))
	if err != nil {
		return "", "", "", caos_errs.ThrowPermissionDenied(nil, "APP-8EF0zZ", "invalid token")
	}
	token, err := repo.View.TokenByID(tokenID)
	if err != nil {
		return "", "", "", caos_errs.ThrowPermissionDenied(err, "APP-BxUSiL", "invalid token")
	}
	valid, err := repo.View.IsTokenValid(tokenID)
	if err != nil {
		return "", "", "", err
	}
	if !valid {
		return "", "", "", caos_errs.ThrowPermissionDenied(err, "APP-k9KS0", "invalid token")
	}

	for _, aud := range token.Audience {
		if clientID == aud {
			return token.UserID, clientID, token.UserAgentID, nil
		}
	}
	return "", "", "", caos_errs.ThrowPermissionDenied(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) ProjectIDByClientID(ctx context.Context, clientID string) (projectID string, err error) {
	app, err := repo.View.ApplicationByOIDCClientID(clientID)
	if err != nil {
		return "", err
	}
	return app.ID, nil
}

func (repo *TokenVerifierRepo) verifierClientID(ctx context.Context, appName, appClientID string) (string, error) {
	if appClientID != "" {
		return appClientID, nil
	}
	iam, err := repo.IamEvents.IamByID(ctx, repo.IamID)
	if err != nil {
		return "", err
	}
	app, err := repo.View.ApplicationByProjecIDAndAppName(iam.IamProjectID, appName)
	if err != nil {
		return "", err
	}
	return app.OIDCClientID, nil
}
