package handler

import (
	"context"
	"encoding/json"

	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

type RefreshToken struct {
	handler
	subscription *v1.Subscription
}

func newRefreshToken(
	ctx context.Context,
	handler handler,
) *RefreshToken {
	h := &RefreshToken{
		handler: handler,
	}

	h.subscribe(ctx)

	return h
}

func (t *RefreshToken) subscribe(ctx context.Context) {
	t.subscription = t.es.Subscribe(t.AggregateTypes()...)
	go func() {
		for event := range t.subscription.Events {
			query.ReduceEvent(ctx, t, event)
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
	return []es_models.AggregateType{user.AggregateType, project.AggregateType, instance.AggregateType}
}

func (t *RefreshToken) CurrentSequence(ctx context.Context, instanceID string) (uint64, error) {
	sequence, err := t.view.GetLatestRefreshTokenSequence(ctx, instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (t *RefreshToken) EventQuery(ctx context.Context, instanceIDs []string) (*es_models.SearchQuery, error) {
	sequences, err := t.view.GetLatestRefreshTokenSequences(ctx, instanceIDs)
	if err != nil {
		return nil, err
	}
	return newSearchQuery(sequences, t.AggregateTypes(), instanceIDs), nil
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
			logging.WithError(err).Error("could not unmarshal event data")
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
			logging.WithError(err).Error("could not unmarshal event data")
			return caos_errs.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
		}
		return t.view.DeleteRefreshToken(e.TokenID, event.InstanceID, event)
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return t.view.DeleteUserRefreshTokens(event.AggregateID, event.InstanceID, event)
	case instance.InstanceRemovedEventType:
		return t.view.DeleteInstanceRefreshTokens(event)
	case org.OrgRemovedEventType:
		return t.view.DeleteOrgRefreshTokens(event)
	default:
		return t.view.ProcessedRefreshTokenSequence(event)
	}
}

func (t *RefreshToken) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, t.view.GetLatestRefreshTokenFailedEvent, t.view.ProcessedRefreshTokenFailedEvent, t.view.ProcessedRefreshTokenSequence, t.errorCountUntilSkip)
}

func (t *RefreshToken) OnSuccess(instanceIDs []string) error {
	return spooler.HandleSuccess(t.view.UpdateRefreshTokenSpoolerRunTimestamp, instanceIDs)
}
