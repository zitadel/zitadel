package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-eUpQU", "Errors.idp.config.notset")
	}

	idpConfigID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	addedConfig := NewIAMIDPConfigWriteModel(idpConfigID)

	clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), r.idpConfigSecretCrypto)
	if err != nil {
		return nil, err
	}

	iamAgg := IAMAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.EventPusher{
		iam_repo.NewIDPConfigAddedEvent(
			ctx,
			iamAgg,
			idpConfigID,
			config.Name,
			config.Type,
			config.StylingType,
		),
		iam_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			iamAgg,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			clientSecret,
			config.OIDCConfig.IDPDisplayNameMapping,
			config.OIDCConfig.UsernameMapping,
			config.OIDCConfig.Scopes...,
		),
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&addedConfig.IDPConfigWriteModel), nil
}

func (r *CommandSide) ChangeDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, config.IDPConfigID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M9so", "Errors.IAM.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, iamAgg, config.IDPConfigID, config.Name, config.StylingType)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&existingIDP.IDPConfigWriteModel), nil
}

func (r *CommandSide) DeactivateDefaultIDPConfig(ctx context.Context, idpID string) error {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9so", "Errors.IAM.IDPConfig.NotActive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, iam_repo.NewIDPConfigDeactivatedEvent(ctx, iamAgg, idpID))
	return err
}

func (r *CommandSide) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) error {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, iam_repo.NewIDPConfigReactivatedEvent(ctx, iamAgg, idpID))
	return err
}

func (r *CommandSide) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.ExternalIDP) error {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return caos_errs.ThrowNotFound(nil, "IAM-4M0xy", "Errors.IAM.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.EventPusher{
		iam_repo.NewIDPConfigRemovedEvent(ctx, iamAgg, idpID, existingIDP.Name),
	}

	for _, idpProvider := range idpProviders {
		if idpProvider.AggregateID == domain.IAMID {
			userEvents := r.removeIDPProviderFromDefaultLoginPolicy(ctx, iamAgg, idpProvider, true, externalIDPs...)
			events = append(events, userEvents...)
		}
		orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
		orgEvents := r.removeIDPProviderFromLoginPolicy(ctx, orgAgg, idpID, true)
		events = append(events, orgEvents...)
	}

	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) iamIDPConfigWriteModelByID(ctx context.Context, idpID string) (policy *IAMIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMIDPConfigWriteModel(idpID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
