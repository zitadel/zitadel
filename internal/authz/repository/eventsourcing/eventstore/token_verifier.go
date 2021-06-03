package eventstore

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	features_view_model "github.com/caos/zitadel/internal/features/repository/view/model"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type TokenVerifierRepo struct {
	TokenVerificationKey crypto.EncryptionAlgorithm
	IAMID                string
	Eventstore           v1.Eventstore
	View                 *view.View
}

func (repo *TokenVerifierRepo) TokenByID(ctx context.Context, tokenID, userID string) (*usr_model.TokenView, error) {
	token, viewErr := repo.View.TokenByID(tokenID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
	}

	events, esErr := repo.getUserEvents(ctx, userID, token.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-5Nm9s").WithError(viewErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return model.TokenViewToModel(token), nil
	}
	viewToken := *token
	for _, event := range events {
		err := token.AppendEventIfMyToken(event)
		if err != nil {
			return model.TokenViewToModel(&viewToken), nil
		}
	}
	if !token.Expiration.After(time.Now().UTC()) || token.Deactivated {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, clientID string) (userID string, agentID string, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	tokenData, err := base64.RawURLEncoding.DecodeString(tokenString)
	if err != nil {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-ASdgg", "invalid token")
	}
	tokenIDSubject, err := repo.TokenVerificationKey.DecryptString(tokenData, repo.TokenVerificationKey.EncryptionKeyID())
	if err != nil {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}

	splittedToken := strings.Split(tokenIDSubject, ":")
	if len(splittedToken) != 2 {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-GDg3a", "invalid token")
	}
	token, err := repo.TokenByID(ctx, splittedToken[0], splittedToken[1])
	if err != nil {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}

	projectID, _, err := repo.ProjectIDAndOriginsByClientID(ctx, clientID)
	if err != nil {
		return "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-5M9so", "invalid token")
	}
	for _, aud := range token.Audience {
		if clientID == aud || projectID == aud {
			return token.UserID, token.UserAgentID, token.PreferredLanguage, token.ResourceOwner, nil
		}
	}
	return "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
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

func (repo *TokenVerifierRepo) CheckOrgFeatures(ctx context.Context, orgID string, requiredFeatures ...string) error {
	features, err := repo.View.FeaturesByAggregateID(orgID)
	if caos_errs.IsNotFound(err) {
		return repo.checkDefaultFeatures(ctx, requiredFeatures...)
	}
	if err != nil {
		return err
	}
	return checkFeatures(features, requiredFeatures...)
}

func checkFeatures(features *features_view_model.FeaturesView, requiredFeatures ...string) error {
	for _, requiredFeature := range requiredFeatures {
		if strings.HasPrefix(requiredFeature, domain.FeatureLoginPolicy) {
			if err := checkLoginPolicyFeatures(features, requiredFeature); err != nil {
				return err
			}
		}
		if requiredFeature == domain.FeaturePasswordComplexityPolicy && !features.PasswordComplexityPolicy {
			return MissingFeatureErr(requiredFeature)
		}
		if requiredFeature == domain.FeatureLabelPolicy && !features.PasswordComplexityPolicy {
			return MissingFeatureErr(requiredFeature)
		}
	}
	return nil
}

func checkLoginPolicyFeatures(features *features_view_model.FeaturesView, requiredFeature string) error {
	switch requiredFeature {
	case domain.FeatureLoginPolicyFactors:
		if !features.LoginPolicyFactors {
			return MissingFeatureErr(requiredFeature)
		}
	case domain.FeatureLoginPolicyIDP:
		if !features.LoginPolicyIDP {
			return MissingFeatureErr(requiredFeature)
		}
	case domain.FeatureLoginPolicyPasswordless:
		if !features.LoginPolicyPasswordless {
			return MissingFeatureErr(requiredFeature)
		}
	case domain.FeatureLoginPolicyRegistration:
		if !features.LoginPolicyRegistration {
			return MissingFeatureErr(requiredFeature)
		}
	case domain.FeatureLoginPolicyUsernameLogin:
		if !features.LoginPolicyUsernameLogin {
			return MissingFeatureErr(requiredFeature)
		}
	case domain.FeatureLoginPolicyPasswordReset:
		if !features.LoginPolicyPasswordReset {
			return MissingFeatureErr(requiredFeature)
		}
	default:
		if !features.LoginPolicyFactors && !features.LoginPolicyIDP && !features.LoginPolicyPasswordless && !features.LoginPolicyRegistration && !features.LoginPolicyUsernameLogin {
			return MissingFeatureErr(requiredFeature)
		}
	}
	return nil
}

func MissingFeatureErr(feature string) error {
	return caos_errs.ThrowPermissionDeniedf(nil, "AUTH-Dvgsf", "missing feature %v", feature)
}

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	iam, err := repo.getIAMByID(ctx)
	if err != nil {
		return "", err
	}
	app, err := repo.View.ApplicationByProjecIDAndAppName(ctx, iam.IAMProjectID, appName)
	if err != nil {
		return "", err
	}
	return app.OIDCClientID, nil
}

func (r *TokenVerifierRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}

func (u *TokenVerifierRepo) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &iam_es_model.IAM{
		ObjectRoot: models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (repo *TokenVerifierRepo) checkDefaultFeatures(ctx context.Context, requiredFeatures ...string) error {
	features, viewErr := repo.View.FeaturesByAggregateID(domain.IAMID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		features = new(features_view_model.FeaturesView)
	}
	events, esErr := repo.getIAMEvents(ctx, features.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return checkFeatures(features, requiredFeatures...)
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return esErr
	}
	featuresCopy := *features
	for _, event := range events {
		if err := featuresCopy.AppendEvent(event); err != nil {
			return checkFeatures(features, requiredFeatures...)
		}
	}
	return checkFeatures(&featuresCopy, requiredFeatures...)
}

func (repo *TokenVerifierRepo) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}
