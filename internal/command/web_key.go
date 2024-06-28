package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/webkey"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) GenerateWebKey(ctx context.Context, conf crypto.WebKeyConfig) (*domain.ObjectDetails, error) {
	_, activeID, err := c.getAllWebKeys(ctx)
	if err != nil {
		return nil, err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	encryptedPrivate, public, err := crypto.GenerateEncryptedWebKey(keyID, c.keyAlgorithm, conf)
	if err != nil {
		return nil, err
	}

	aggregate := webkey.NewAggregate(keyID, authz.GetInstance(ctx).InstanceID())
	addedCmd, err := webkey.NewAddedEvent(ctx, aggregate, encryptedPrivate, public, conf)
	if err != nil {
		return nil, err
	}
	commands := []eventstore.Command{addedCmd}

	// make sure the first key gets activated by default
	if activeID == "" {
		commands = append(commands, webkey.NewActivatedEvent(ctx, aggregate))
	}

	events, err := c.eventstore.Push(ctx, commands...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) ActivateWebKey(ctx context.Context, keyID string) (_ *domain.ObjectDetails, err error) {
	keys, activeID, err := c.getAllWebKeys(ctx)
	if err != nil {
		return nil, err
	}
	nextActive, ok := keys[keyID]
	if !ok {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-teiG3", "Errors.WebKey.NotFound")
	}

	commands := make([]eventstore.Command, 0, 2)
	commands = append(commands, webkey.NewActivatedEvent(ctx,
		webkey.AggregateFromWriteModel(ctx, &nextActive.WriteModel),
	))
	if activeID != "" {
		commands = append(commands, webkey.NewDeactivatedEvent(ctx,
			webkey.AggregateFromWriteModel(ctx, &keys[activeID].WriteModel),
		))
	}

	events, err := c.eventstore.Push(ctx, commands...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) getAllWebKeys(ctx context.Context) (_ map[string]*WebKeyWriteModel, activeID string, err error) {
	models := newWebKeyWriteModels(authz.GetInstance(ctx).InstanceID())
	if err = c.eventstore.FilterToQueryReducer(ctx, models); err != nil {
		return nil, "", err
	}
	return models.keys, models.activeID, nil
}

func (c *Commands) RemoveWebKey(ctx context.Context, keyID string) (_ *domain.ObjectDetails, err error) {
	model := NewWebKeyWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	if err = c.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	if model.State == domain.WebKeyStateUndefined {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ooCa7", "Errors.WebKey.NotFound")
	}
	if model.State == domain.WebKeyStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Chai1", "Errors.WebKey.Active")
	}
	events, err := c.eventstore.Push(ctx, webkey.NewRemovedEvent(ctx,
		webkey.AggregateFromWriteModel(ctx, &model.WriteModel),
	))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
