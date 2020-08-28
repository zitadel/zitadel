package eventstore

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type TokenVerifierRepo struct {
	TokenVerificationKey [32]byte
	IAMID                string
	IAMEvents            *iam_event.IAMEventstore
	ProjectEvents        *proj_event.ProjectEventstore
	View                 *view.View
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, clientID string) (userID string, agentID string, prefLang string, err error) {
	//TODO: use real key
	tokenID, err := crypto.DecryptAESString(tokenString, string(repo.TokenVerificationKey[:32]))
	if err != nil {
		return "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}
	token, err := repo.View.TokenByID(tokenID)
	if err != nil {
		return "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}

	for _, aud := range token.Audience {
		if clientID == aud {
			return token.UserID, token.UserAgentID, token.PreferredLanguage, nil
		}
	}
	return "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error) {
	app, err := repo.View.ApplicationByOIDCClientID(clientID)
	if err != nil {
		return "", nil, err
	}
	return app.ProjectID, app.OriginAllowList, nil
}

func (repo *TokenVerifierRepo) ExistsOrg(ctx context.Context, orgID string) error {
	_, err := repo.View.OrgByID(orgID)
	return err
}

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (string, error) {
	iam, err := repo.IAMEvents.IAMByID(ctx, repo.IAMID)
	if err != nil {
		return "", err
	}
	app, err := repo.View.ApplicationByProjecIDAndAppName(iam.IAMProjectID, appName)
	if err != nil {
		return "", err
	}
	return app.OIDCClientID, nil
}
