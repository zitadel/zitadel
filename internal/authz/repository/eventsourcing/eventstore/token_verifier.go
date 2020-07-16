package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"time"
)

type TokenVerifierRepo struct {
	TokenVerificationKey [32]byte
	IamID                string
	IamEvents            *iam_event.IamEventstore
	ProjectEvents        *proj_event.ProjectEventstore
	View                 *view.View
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, clientID string) (userID string, agentID string, err error) {
	//TODO: use real key
	tokenID, err := crypto.DecryptAESString(tokenString, string(repo.TokenVerificationKey[:32]))
	if err != nil {
		return "", "", caos_errs.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}
	token, err := repo.View.TokenByID(tokenID)
	if err != nil {
		return "", "", caos_errs.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", caos_errs.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}

	for _, aud := range token.Audience {
		if clientID == aud {
			return token.UserID, token.UserAgentID, nil
		}
	}
	return "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) ProjectIDByClientID(ctx context.Context, clientID string) (projectID string, err error) {
	app, err := repo.View.ApplicationByOIDCClientID(clientID)
	if err != nil {
		return "", err
	}
	return app.ProjectID, nil
}

func (repo *TokenVerifierRepo) ExistsOrg(ctx context.Context, orgID string) error {
	_, err := repo.View.OrgByID(orgID)
	return err
}

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (string, error) {
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
