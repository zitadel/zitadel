package eventstore

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type TokenVerifierRepo struct {
	TokenVerificationKey crypto.EncryptionAlgorithm
	IAMID                string
	Eventstore           v1.Eventstore
	View                 *view.View
	Query                *query.Queries
}

func (repo *TokenVerifierRepo) tokenByID(ctx context.Context, tokenID, userID string) (*usr_model.TokenView, error) {
	token, viewErr := repo.View.TokenByID(tokenID, authz.GetInstance(ctx).InstanceID())
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
	}

	events, esErr := repo.getUserEvents(ctx, userID, token.InstanceID, token.Sequence)
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

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, verifierClientID, projectID string) (userID string, agentID string, clientID, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	tokenData, err := base64.RawURLEncoding.DecodeString(tokenString)
	if err != nil {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-ASdgg", "invalid token")
	}
	tokenIDSubject, err := repo.TokenVerificationKey.DecryptString(tokenData, repo.TokenVerificationKey.EncryptionKeyID())
	if err != nil {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}

	splittedToken := strings.Split(tokenIDSubject, ":")
	if len(splittedToken) != 2 {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-GDg3a", "invalid token")
	}
	token, err := repo.tokenByID(ctx, splittedToken[0], splittedToken[1])
	if err != nil {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}
	if token.IsPAT {
		return token.UserID, "", "", "", token.ResourceOwner, nil
	}
	for _, aud := range token.Audience {
		if verifierClientID == aud || projectID == aud {
			return token.UserID, token.UserAgentID, token.ApplicationID, token.PreferredLanguage, token.ResourceOwner, nil
		}
	}
	return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error) {
	app, err := repo.View.ApplicationByOIDCClientID(ctx, clientID)
	if err != nil {
		return "", nil, err
	}
	return app.ProjectID, app.OIDCConfig.AllowedOrigins, nil
}

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (clientID, projectID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := repo.View.ApplicationByProjecIDAndAppName(ctx, authz.GetInstance(ctx).ProjectID(), appName)
	if err != nil {
		return "", "", err
	}
	if app.OIDCConfig != nil {
		clientID = app.OIDCConfig.ClientID
	} else if app.APIConfig != nil {
		clientID = app.APIConfig.ClientID
	}
	return clientID, app.ProjectID, nil
}

func (r *TokenVerifierRepo) getUserEvents(ctx context.Context, userID, instanceID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, instanceID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}
