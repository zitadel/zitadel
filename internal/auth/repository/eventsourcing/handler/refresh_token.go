package handler

import (
	"encoding/json"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
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

func (t *RefreshToken) Subscription() *v1.Subscription {
	return t.subscription
}

func (t *RefreshToken) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user.AggregateType, project.AggregateType}
}

func (t *RefreshToken) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := t.view.GetLatestRefreshTokenSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (t *RefreshToken) EventQuery() (*es_models.SearchQuery, error) {
	sequences, err := t.view.GetLatestRefreshTokenSequences()
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
		ExcludedInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (t *RefreshToken) Reduce(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case user.HumanRefreshTokenAddedType:
		token := new(view_model.RefreshTokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutRefreshToken(token, event)
	case user.HumanRefreshTokenRenewedType:
		e := new(user.HumanRefreshTokenRenewedEvent)
		if err := json.Unmarshal(event.Data, e); err != nil {
			logging.Log("EVEN-DBbn4").WithError(err).Error("could not unmarshal event data")
			return caos_errs.ThrowInternal(nil, "MODEL-BHn75", "could not unmarshal data")
		}
		token, err := t.view.RefreshTokenByID(e.TokenID, event.InstanceID)
		if err != nil {
			return err
		}
		err = token.AppendEvent(event)
		if err != nil {
			return err
		}
		return t.view.PutRefreshToken(token, event)
	case user.HumanRefreshTokenRemovedType:
		e := new(user.HumanRefreshTokenRemovedEvent)
		if err := json.Unmarshal(event.Data, e); err != nil {
			logging.Log("EVEN-BDbh3").WithError(err).Error("could not unmarshal event data")
			return caos_errs.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
		}
		return t.view.DeleteRefreshToken(e.TokenID, event.InstanceID, event)
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return t.view.DeleteUserRefreshTokens(event.AggregateID, event.InstanceID, event)
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
