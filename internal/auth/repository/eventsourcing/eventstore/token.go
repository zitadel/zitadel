package eventstore

import (
	"context"
	"strings"
	"time"

	"github.com/caos/logging"

	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
)

type TokenRepo struct {
	UserEvents    *user_event.UserEventstore
	ProjectEvents *proj_event.ProjectEventstore
	View          *view.View
}

func (repo *TokenRepo) CreateToken(ctx context.Context, agentID, clientID, subject string, audience, scopes []string, lifetime time.Duration) (string, time.Time, error) {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, auth_req_model.ProjectIDScope) && strings.HasSuffix(scope, auth_req_model.AudSuffix) {
			audience = append(audience, strings.TrimSuffix(strings.TrimPrefix(scope, auth_req_model.ProjectIDScope), auth_req_model.AudSuffix))
		}
	}
	if !strings.Contains(subject, "@") { //TODO: improve, but how?!
		token, err := repo.createUserToken(ctx, agentID, clientID, subject, audience, scopes, lifetime)
		if err != nil {
			return "", time.Time{}, err
		}
		return token.TokenID, token.Expiration, nil
	}
	token, err := repo.createApplicationToken(ctx, subject, audience, scopes, lifetime)
	if err != nil {
		return "", time.Time{}, err
	}
	return token.TokenID, token.Expiration, nil
}

func (repo *TokenRepo) createUserToken(ctx context.Context, agentID, clientID, subject string, audience, scopes []string, lifetime time.Duration) (*usr_model.Token, error) {
	preferredLanguage := ""
	user, _ := repo.View.UserByID(subject)
	if user != nil {
		preferredLanguage = user.PreferredLanguage
	}
	now := time.Now().UTC()
	token := &usr_model.Token{
		ObjectRoot: models.ObjectRoot{
			AggregateID: subject,
		},
		UserAgentID:       agentID,
		ApplicationID:     clientID,
		Audience:          audience,
		Scopes:            scopes,
		Expiration:        now.Add(lifetime),
		PreferredLanguage: preferredLanguage,
	}
	return repo.UserEvents.TokenAdded(ctx, token)
}

func (repo *TokenRepo) createApplicationToken(ctx context.Context, subject string, audience, scopes []string, lifetime time.Duration) (*proj_model.Token, error) {
	app, err := repo.View.ApplicationByClientID(ctx, subject)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	token := &proj_model.Token{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectID,
		},
		ClientID:   subject,
		Audience:   audience,
		Scopes:     scopes,
		Expiration: now.Add(lifetime),
	}
	return repo.ProjectEvents.TokenAdded(ctx, token)
}

func (repo *TokenRepo) IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error) {
	token, err := repo.TokenByID(ctx, userID, tokenID)
	if err == nil {
		return token.Expiration.After(time.Now().UTC()), nil
	}
	if errors.IsNotFound(err) {
		return false, nil
	}
	return false, err
}

func (repo *TokenRepo) TokenByID(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error) {
	token, viewErr := repo.View.TokenByID(tokenID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
	}

	events, esErr := repo.UserEvents.UserEventsByID(ctx, userID, token.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
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
		return nil, errors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func AppendAudIfNotExisting(aud string, existingAud []string) []string {
	for _, a := range existingAud {
		if a == aud {
			return existingAud
		}
	}
	return append(existingAud, aud)
}
