package handler

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	proj_model "github.com/caos/zitadel/internal/project/model"
	project_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

type Token struct {
	handler
	subscription *v1.Subscription
}

func newToken(
	handler handler,
) *Token {
	h := &Token{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (t *Token) subscribe() {
	t.subscription = t.es.Subscribe(t.AggregateTypes()...)
	go func() {
		for event := range t.subscription.Events {
			query.ReduceEvent(t, event)
		}
	}()
}

func (t *Token) ViewModel() string {
	return tokenTable
}

func (t *Token) Subscription() *v1.Subscription {
	return t.subscription
}

func (_ *Token) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user.AggregateType, project.AggregateType}
}

func (p *Token) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := p.view.GetLatestTokenSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (t *Token) EventQuery() (*es_models.SearchQuery, error) {
	sequences, err := t.view.GetLatestTokenSequences()
	if err != nil {
		return nil, err
	}
	query := es_models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(t.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(t.AggregateTypes()...).
		LatestSequenceFilter(0).
		IgnoredInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (t *Token) Reduce(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case user.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		token := new(view_model.TokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutToken(token, event)
	case user.UserV1ProfileChangedType,
		user.HumanProfileChangedType:
		user := new(view_model.UserView)
		user.AppendEvent(event)
		tokens, err := t.view.TokensByUserID(event.AggregateID, event.InstanceID)
		if err != nil {
			return err
		}
		for _, token := range tokens {
			token.PreferredLanguage = user.PreferredLanguage
		}
		return t.view.PutTokens(tokens, event)
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		return t.view.DeleteSessionTokens(id, event.AggregateID, event.InstanceID, event)
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return t.view.DeleteUserTokens(event.AggregateID, event.InstanceID, event)
	case user_repo.UserTokenRemovedType,
		user_repo.PersonalAccessTokenRemovedType:
		id, err := tokenIDFromRemovedEvent(event)
		if err != nil {
			return err
		}
		return t.view.DeleteToken(id, event.InstanceID, event)
	case user_repo.HumanRefreshTokenRemovedType:
		id, err := refreshTokenIDFromRemovedEvent(event)
		if err != nil {
			return err
		}
		return t.view.DeleteTokensFromRefreshToken(id, event.InstanceID, event)
	case project.ApplicationDeactivatedType,
		project.ApplicationRemovedType:
		application, err := applicationFromSession(event)
		if err != nil {
			return err
		}
		return t.view.DeleteApplicationTokens(event, application.AppID)
	case project.ProjectDeactivatedType,
		project.ProjectRemovedType:
		project, err := t.getProjectByID(context.Background(), event.AggregateID, event.InstanceID)
		if err != nil {
			return err
		}
		applicationsIDs := make([]string, 0, len(project.Applications))
		for _, app := range project.Applications {
			applicationsIDs = append(applicationsIDs, app.AppID)
		}
		return t.view.DeleteApplicationTokens(event, applicationsIDs...)
	default:
		return t.view.ProcessedTokenSequence(event)
	}
}

func (t *Token) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3jkl4", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, t.view.GetLatestTokenFailedEvent, t.view.ProcessedTokenFailedEvent, t.view.ProcessedTokenSequence, t.errorCountUntilSkip)
}

func agentIDFromSession(event *es_models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &session); err != nil {
		logging.Log("EVEN-s3bq9").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}

func applicationFromSession(event *es_models.Event) (*project_es_model.Application, error) {
	application := new(project_es_model.Application)
	if err := json.Unmarshal(event.Data, &application); err != nil {
		logging.Log("EVEN-GRE2q").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-Hrw1q", "could not unmarshal data")
	}
	return application, nil
}

func tokenIDFromRemovedEvent(event *es_models.Event) (string, error) {
	removed := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &removed); err != nil {
		logging.Log("EVEN-Sdff3").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-Sff32", "could not unmarshal data")
	}
	return removed["tokenId"].(string), nil
}

func refreshTokenIDFromRemovedEvent(event *es_models.Event) (string, error) {
	removed := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &removed); err != nil {
		logging.Log("EVEN-Ff23g").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-Dfb3w", "could not unmarshal data")
	}
	return removed["tokenId"].(string), nil
}

func (t *Token) OnSuccess() error {
	return spooler.HandleSuccess(t.view.UpdateTokenSpoolerRunTimestamp)
}

func (t *Token) getProjectByID(ctx context.Context, projID, instanceID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, instanceID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &project_es_model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, t.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Dsdw2", "Errors.Project.NotFound")
	}

	return project_es_model.ProjectToModel(esProject), nil
}
