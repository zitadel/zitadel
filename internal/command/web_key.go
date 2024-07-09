package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/webkey"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type WebKeyDetails struct {
	KeyID         string
	ObjectDetails *domain.ObjectDetails
}

// GenerateWebKey creates one web key pair for the instance.
// If the instance does not have an active key, the new key is activated.
func (c *Commands) GenerateWebKey(ctx context.Context, conf crypto.WebKeyConfig) (*WebKeyDetails, error) {
	_, activeID, err := c.getAllWebKeys(ctx)
	if err != nil {
		return nil, err
	}
	addedCmd, aggregate, err := c.generateWebKeyCommand(ctx, authz.GetInstance(ctx).InstanceID(), conf)
	if err != nil {
		return nil, err
	}
	commands := []eventstore.Command{addedCmd}
	if activeID == "" {
		commands = append(commands, webkey.NewActivatedEvent(ctx, aggregate))
	}
	events, err := c.eventstore.Push(ctx, commands...)
	if err != nil {
		return nil, err
	}
	return &WebKeyDetails{
		KeyID:         aggregate.ID,
		ObjectDetails: pushedEventsToObjectDetails(events),
	}, nil
}

// GenerateInitialWebKeys creates 2 web key pairs for the instance.
// The first key is activated for signing use.
// If the instance already has keys, this is noop.
func (c *Commands) GenerateInitialWebKeys(ctx context.Context, conf crypto.WebKeyConfig) error {
	keys, _, err := c.getAllWebKeys(ctx)
	if err != nil {
		return err
	}
	if len(keys) != 0 {
		return nil
	}
	commands, err := c.generateInitialWebKeysCommands(ctx, authz.GetInstance(ctx).InstanceID(), conf)
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, commands...)
	return err
}

func (c *Commands) generateInitialWebKeysCommands(ctx context.Context, instanceID string, conf crypto.WebKeyConfig) ([]eventstore.Command, error) {
	commands := make([]eventstore.Command, 0, 3)
	for i := 0; i < 2; i++ {
		addedCmd, aggregate, err := c.generateWebKeyCommand(ctx, instanceID, conf)
		if err != nil {
			return nil, err
		}
		commands = append(commands, addedCmd)
		if i == 0 {
			commands = append(commands, webkey.NewActivatedEvent(ctx, aggregate))
		}
	}
	return commands, nil
}

func (c *Commands) generateWebKeyCommand(ctx context.Context, instanceID string, conf crypto.WebKeyConfig) (eventstore.Command, *eventstore.Aggregate, error) {
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	encryptedPrivate, public, err := crypto.GenerateEncryptedWebKey(keyID, c.keyAlgorithm, conf)
	if err != nil {
		return nil, nil, err
	}
	aggregate := webkey.NewAggregate(keyID, instanceID)
	addedCmd, err := webkey.NewAddedEvent(ctx, aggregate, encryptedPrivate, public, conf)
	if err != nil {
		return nil, nil, err
	}
	return addedCmd, aggregate, nil
}

// ActivateWebKey activates the key identified by keyID.
// Any previously activated key on the current instance is deactivated.
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

// getAllWebKeys searches for all web keys on the instance and returns a map of key IDs.
// activeID is the id of the currently active key.
func (c *Commands) getAllWebKeys(ctx context.Context) (_ map[string]*WebKeyWriteModel, activeID string, err error) {
	models := newWebKeyWriteModels(authz.GetInstance(ctx).InstanceID())
	if err = c.eventstore.FilterToQueryReducer(ctx, models); err != nil {
		return nil, "", err
	}
	return models.keys, models.activeID, nil
}

func (c *Commands) DeleteWebKey(ctx context.Context, keyID string) (_ *domain.ObjectDetails, err error) {
	model := NewWebKeyWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	if err = c.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	if model.State == domain.WebKeyStateUnspecified {
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

func (c *Commands) prepareGenerateInitialWebKeys(instanceID string, conf crypto.WebKeyConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return c.generateInitialWebKeysCommands(ctx, instanceID, conf)
		}, nil
	}
}
