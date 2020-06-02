package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"time"
)

type TokenVerifierRepo struct {
	TokenVerificationKey string
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
	tokenID, err := crypto.DecryptAESString(tokenString, repo.TokenVerificationKey)
	if err != nil {
		return "", "", "", caos_errs.ThrowPermissionDenied(nil, "APP-8EF0zZ", "invalid token")
	}
	token, err := repo.View.TokenByID(tokenID)
	if err != nil {
		return "", "", "", caos_errs.ThrowPermissionDenied(err, "APP-BxUSiL", "invalid token")
	}
	if token.Expiration.Before(time.Now().UTC()) {
		return "", "", "", caos_errs.ThrowPermissionDenied(err, "APP-BJEyZ1", "token expired")
	}
	for _, aud := range token.Scopes {
		if clientID == aud {
			return token.UserID, clientID, token.UserAgentID, nil
		}
	}
	return "", "", "", caos_errs.ThrowPermissionDenied(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) verifierClientID(ctx context.Context, appName, appClientID string) (string, error) {
	if appClientID != "" {
		return appClientID, nil
	}
	iam, err := repo.IamEvents.IamByID(ctx, repo.IamID)
	if err != nil {
		return "", err
	}
	apps, _, err := repo.View.SearchApplications(&model.ApplicationSearchRequest{
		Queries: []*model.ApplicationSearchQuery{
			&model.ApplicationSearchQuery{Key: model.APPLICATIONSEARCHKEY_PROJECT_ID, Method: global_model.SEARCHMETHOD_EQUALS, Value: iam.IamProjectID},
			&model.ApplicationSearchQuery{Key: model.APPLICATIONSEARCHKEY_NAME, Method: global_model.SEARCHMETHOD_EQUALS, Value: appName},
		},
		Limit: 1,
	},
	)
	if err != nil {
		return "", err
	}
	if len(apps) != 1 {
		return "", caos_errs.ThrowNotFound(nil, "APP-ZAQlLQ", "client not found")
	}
	return apps[0].OIDCClientID, nil
}
