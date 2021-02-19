package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	"strings"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type TokenVerifierRepo struct {
	TokenVerificationKey [32]byte
	IAMID                string
	Eventstore           eventstore.Eventstore
	IAMEvents            *iam_event.IAMEventstore
	ProjectEvents        *proj_event.ProjectEventstore
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
	//TODO: use real key
	tokenIDSubject, err := crypto.DecryptAESString(tokenString, string(repo.TokenVerificationKey[:32]))
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

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	iam, err := repo.IAMEvents.IAMByID(ctx, repo.IAMID)
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
