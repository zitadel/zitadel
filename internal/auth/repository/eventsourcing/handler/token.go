package handler

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	project_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

type Token struct {
	handler
	ProjectEvents *proj_event.ProjectEventstore
	subscription  *eventstore.Subscription
}

func newToken(
	handler handler,
	projectEvents *proj_event.ProjectEventstore,
) *Token {
	h := &Token{
		handler:       handler,
		ProjectEvents: projectEvents,
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

func (_ *Token) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user_es_model.UserAggregate, project_es_model.ProjectAggregate}
}

func (p *Token) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestTokenSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (t *Token) EventQuery() (*models.SearchQuery, error) {
	sequence, err := t.view.GetLatestTokenSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(user_es_model.UserAggregate, project_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (t *Token) Reduce(event *models.Event) (err error) {
	switch event.Type {
	case user_es_model.UserTokenAdded,
		project_es_model.TokenAdded:
		token := new(view_model.TokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutToken(token, event)
	case user_es_model.UserProfileChanged,
		user_es_model.HumanProfileChanged:
		user := new(view_model.UserView)
		user.AppendEvent(event)
		tokens, err := t.view.TokensByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, token := range tokens {
			token.PreferredLanguage = user.PreferredLanguage
		}
		return t.view.PutTokens(tokens, event)
	case user_es_model.SignedOut,
		user_es_model.HumanSignedOut:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		return t.view.DeleteSessionTokens(id, event.AggregateID, event)
	case user_es_model.UserLocked,
		user_es_model.UserDeactivated,
		user_es_model.UserRemoved:
		return t.view.DeleteUserTokens(event.AggregateID, event)
	case project_es_model.ApplicationDeactivated,
		project_es_model.ApplicationRemoved:
		application, err := applicationFromSession(event)
		if err != nil {
			return err
		}
		return t.view.DeleteApplicationTokens(event, application.AppID)
	case project_es_model.ProjectDeactivated,
		project_es_model.ProjectRemoved:
		project, err := t.ProjectEvents.ProjectByID(context.Background(), event.AggregateID)
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

func (t *Token) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-3jkl4", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, t.view.GetLatestTokenFailedEvent, t.view.ProcessedTokenFailedEvent, t.view.ProcessedTokenSequence, t.errorCountUntilSkip)
}

func agentIDFromSession(event *models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &session); err != nil {
		logging.Log("EVEN-s3bq9").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}

func applicationFromSession(event *models.Event) (*project_es_model.Application, error) {
	application := new(project_es_model.Application)
	if err := json.Unmarshal(event.Data, &application); err != nil {
		logging.Log("EVEN-GRE2q").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-Hrw1q", "could not unmarshal data")
	}
	return application, nil
}

func (t *Token) OnSuccess() error {
	return spooler.HandleSuccess(t.view.UpdateTokenSpoolerRunTimestamp)
}
