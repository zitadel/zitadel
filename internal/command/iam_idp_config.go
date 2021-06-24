package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) AddDefaultIDPConfig(ctx context.Context, config domain.IDPConfig) (string, *domain.ObjectDetails, error) {
	idpConfigID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	addedConfig := NewIAMIDPConfigWriteModel(idpConfigID)

	iamAgg := IAMAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.EventPusher{
		iam_repo.NewIDPConfigAddedEvent(
			ctx,
			iamAgg,
			idpConfigID,
			config.IDPConfigName(),
			config.IDPConfigType(),
			config.IDPConfigStylingType(),
		),
	}
	var configEventCreator configEventCreator
	switch conf := config.(type) {
	case *domain.OIDCIDPConfig:
		configEventCreator = c.addDefaultOIDCIDPConfig(conf)
	case *domain.AuthConnectorIDPConfig:
		configEventCreator = c.addDefaultAuthConnectorIDPConfig(conf)
	default:
		return "", nil, errors.ThrowInvalidArgument(nil, "IAM-eUpQU", "Errors.idp.config.notset")
	}
	configEvent, err := configEventCreator(ctx, iamAgg, idpConfigID)
	if err != nil {
		return "", nil, err
	}
	events = append(events, configEvent)
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return idpConfigID, writeModelToObjectDetails(&addedConfig.IDPConfigWriteModel.WriteModel), nil
}

type configEventCreator func(context.Context, *eventstore.Aggregate, string) (eventstore.EventPusher, error)

func (c *Commands) addDefaultOIDCIDPConfig(config *domain.OIDCIDPConfig) configEventCreator {
	return func(ctx context.Context, agg *eventstore.Aggregate, idpConfigID string) (eventstore.EventPusher, error) {
		clientSecret, err := crypto.Encrypt([]byte(config.ClientSecretString), c.idpConfigSecretCrypto)
		if err != nil {
			return nil, err
		}

		return iam_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			agg,
			config.ClientID,
			idpConfigID,
			config.Issuer,
			clientSecret,
			config.IDPDisplayNameMapping,
			config.UsernameMapping,
			config.Scopes...,
		), nil
	}
}

func (c *Commands) addDefaultAuthConnectorIDPConfig(config *domain.AuthConnectorIDPConfig) configEventCreator {
	return func(ctx context.Context, agg *eventstore.Aggregate, idpConfigID string) (eventstore.EventPusher, error) {
		return iam_repo.NewIDPAuthConnectorConfigAddedEvent(
			ctx,
			agg,
			idpConfigID,
			config.BaseURL,
			config.ProviderID,
			config.MachineID,
		), nil
	}
}

func (c *Commands) ChangeDefaultIDPConfig(ctx context.Context, config domain.IDPConfig) (*domain.ObjectDetails, error) {
	if config.ID() == "" {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-4m9gs", "Errors.IDMissing")
	}
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, config.ID())
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M9so", "Errors.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, iamAgg, config.ID(), config.IDPConfigName(), config.IDPConfigStylingType())
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9so", "Errors.IAM.IDPConfig.NotActive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewIDPConfigDeactivatedEvent(ctx, iamAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewIDPConfigReactivatedEvent(ctx, iamAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.ExternalIDP) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M0xy", "Errors.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.EventPusher{
		iam_repo.NewIDPConfigRemovedEvent(ctx, iamAgg, idpID, existingIDP.Name),
	}

	for _, idpProvider := range idpProviders {
		if idpProvider.AggregateID == domain.IAMID {
			userEvents := c.removeIDPProviderFromDefaultLoginPolicy(ctx, iamAgg, idpProvider, true, externalIDPs...)
			events = append(events, userEvents...)
		}
		orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
		orgEvents := c.removeIDPProviderFromLoginPolicy(ctx, orgAgg, idpID, true)
		events = append(events, orgEvents...)
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) getIAMIDPConfigByID(ctx context.Context, idpID string) (domain.IDPConfig, error) {
	config, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if !config.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M9so", "Errors.IDPConfig.NotExisting")
	}
	return writeModelToIDPConfig(&config.IDPConfigWriteModel), nil
}

func (c *Commands) iamIDPConfigWriteModelByID(ctx context.Context, idpID string) (policy *IAMIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMIDPConfigWriteModel(idpID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
