package handler

import (
	"encoding/json"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	project_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

type RefreshToken struct {
	handler
	subscription *v1.Subscription
}

func newRefreshToken(
	handler handler,
) *RefreshToken {
	h := &RefreshToken{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (t *RefreshToken) subscribe() {
	t.subscription = t.es.Subscribe(t.AggregateTypes()...)
	go func() {
		for event := range t.subscription.Events {
			query.ReduceEvent(t, event)
		}
	}()
}

func (t *RefreshToken) ViewModel() string {
	return refreshTokenTable
}

func (t *RefreshToken) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user_es_model.UserAggregate, project_es_model.ProjectAggregate}
}

func (t *RefreshToken) CurrentSequence() (uint64, error) {
	sequence, err := t.view.GetLatestRefreshTokenSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (t *RefreshToken) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := t.view.GetLatestRefreshTokenSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(user_es_model.UserAggregate, project_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (t *RefreshToken) Reduce(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case user_repo.HumanRefreshTokenAddedType:
		token := new(view_model.RefreshTokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutRefreshToken(token, event)
	case user_repo.HumanRefreshTokenRenewedType:
		e := new(user_repo.HumanRefreshTokenRenewedEvent)
		if err := json.Unmarshal(event.Data, e); err != nil {
			logging.Log("EVEN-DBbn4").WithError(err).Error("could not unmarshal event data")
			return caos_errs.ThrowInternal(nil, "MODEL-BHn75", "could not unmarshal data")
		}
		token, err := t.view.RefreshTokenByID(e.TokenID)
		if err != nil {
			return err
		}
		err = token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutRefreshToken(token, event)
	case user_repo.HumanRefreshTokenRemovedType:
		e := new(user_repo.HumanRefreshTokenRemovedEvent)
		if err := json.Unmarshal(event.Data, e); err != nil {
			logging.Log("EVEN-BDbh3").WithError(err).Error("could not unmarshal event data")
			return caos_errs.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
		}
		return t.view.DeleteRefreshToken(e.TokenID, event)
	//case user_es_model.UserProfileChanged,
	//	user_es_model.HumanProfileChanged:
	//	user := new(view_model.UserView)
	//	user.AppendEvent(event)
	//	tokens, err := t.view.RefreshTokensByUserID(event.AggregateID)
	//	if err != nil {
	//		return err
	//	}
	//	for _, token := range tokens {
	//		token.PreferredLanguage = user.PreferredLanguage
	//	}
	//	return t.view.PutRefreshTokens(tokens, event)
	//case user_es_model.SignedOut,
	//	user_es_model.HumanSignedOut:
	//	id, err := agentIDFromSession(event)
	//	if err != nil {
	//		return err
	//	}
	//	return t.view.DeleteSessionRefreshTokens(id, event.AggregateID, event)
	//case user_es_model.UserLocked,
	//	user_es_model.UserDeactivated,
	//	user_es_model.UserRemoved:
	//	return t.view.DeleteUserRefreshTokens(event.AggregateID, event)
	//case project_es_model.ApplicationDeactivated,
	//	project_es_model.ApplicationRemoved:
	//	application, err := applicationFromSession(event)
	//	if err != nil {
	//		return err
	//	}
	//	return t.view.DeleteApplicationRefreshTokens(event, application.AppID)
	//case project_es_model.ProjectDeactivated,
	//	project_es_model.ProjectRemoved:
	//	project, err := t.getProjectByID(context.Background(), event.AggregateID)
	//	if err != nil {
	//		return err
	//	}
	//	applicationsIDs := make([]string, 0, len(project.Applications))
	//	for _, app := range project.Applications {
	//		applicationsIDs = append(applicationsIDs, app.AppID)
	//	}
	//	return t.view.DeleteApplicationRefreshTokens(event, applicationsIDs...)
	default:
		return t.view.ProcessedRefreshTokenSequence(event)
	}
}

func (t *RefreshToken) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3jkl4", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, t.view.GetLatestTokenFailedEvent, t.view.ProcessedTokenFailedEvent, t.view.ProcessedTokenSequence, t.errorCountUntilSkip)
}

func (t *RefreshToken) OnSuccess() error {
	return spooler.HandleSuccess(t.view.UpdateTokenSpoolerRunTimestamp)
}
