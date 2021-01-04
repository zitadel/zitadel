package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if config.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-eUpQU", "Errors.idp.config.notset")
	}

	idpConfigID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	//TODO: check name unique on aggregate
	addedConfig := NewIAMIDPConfigWriteModel(config.AggregateID, idpConfigID)

	clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), r.idpConfigSecretCrypto)
	if err != nil {
		return nil, err
	}

	iamAgg := IAMAggregateFromWriteModel(&addedConfig.WriteModel)
	iamAgg.PushEvents(
		iam_repo.NewIDPConfigAddedEvent(
			ctx,
			idpConfigID,
			config.Name,
			domain.IDPConfigType(config.Type),
			domain.IDPConfigStylingType(config.StylingType),
		),
	)
	iamAgg.PushEvents(
		iam_repo.NewIDPOIDCConfigAddedEvent(
			ctx, config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			clientSecret,
			domain.OIDCMappingField(config.OIDCConfig.IDPDisplayNameMapping),
			domain.OIDCMappingField(config.OIDCConfig.UsernameMapping),
			config.OIDCConfig.Scopes...,
		),
	)
	err = r.eventstore.PushAggregate(ctx, addedConfig, iamAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(addedConfig), nil
}

func (r *CommandSide) ChangeDefaultIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, config.AggregateID, config.IDPConfigID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9so", "Errors.IAM.IDPConfig.NotExisting")
	}

	changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, config.IDPConfigID, config.Name, domain.IDPConfigStylingType(config.StylingType))
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingIDP, iamAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(existingIDP), nil
}

func (r *CommandSide) DeactivateDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, iamID, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9so", "Errors.IAM.IDPConfig.NotActive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIDPConfigDeactivatedEvent(ctx, idpID))

	err = r.eventstore.PushAggregate(ctx, existingIDP, iamAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(existingIDP), nil
}

func (r *CommandSide) ReactivateDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	existingIDP, err := r.iamIDPConfigWriteModelByID(ctx, iamID, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIDPConfigReactivatedEvent(ctx, idpID))

	err = r.eventstore.PushAggregate(ctx, existingIDP, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(existingIDP), nil
}

func (r *CommandSide) RemoveDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	writeModel, err := r.pushDefaultIDPWriteModel(ctx, iamID, idpID, func(a *iam.Aggregate, _ *IAMIDPConfigWriteModel) *iam.Aggregate {
		a.Aggregate = *a.PushEvents(iam_repo.NewIDPConfigRemovedEvent(ctx, idpID))
		return a
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) pushDefaultIDPWriteModel(ctx context.Context, iamID, idpID string, eventSetter func(*iam.Aggregate, *IAMIDPConfigWriteModel) *iam.Aggregate) (*IAMIDPConfigWriteModel, error) {
	writeModel := NewIAMIDPConfigWriteModel(iamID, idpID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := eventSetter(IAMAggregateFromWriteModel(&writeModel.WriteModel), writeModel)
	err = r.eventstore.PushAggregate(ctx, writeModel, aggregate)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}

func (r *CommandSide) iamIDPConfigWriteModelByID(ctx context.Context, iamID, idpID string) (policy *IAMIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMIDPConfigWriteModel(iamID, idpID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
