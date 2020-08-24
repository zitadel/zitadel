package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func LoginPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smla8", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLoginPolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyAdded, policy)
	}
}

func LoginPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlco9", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.LoginPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.LoginPolicyChanged, changes)
	}
}

func LoginPolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-S8sio", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyRemoved, nil)
	}
}

func LoginPolicyIDPProviderAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, provider *iam_es_model.IDPProvider, iamID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate, iam_es_model.IamAggregate).
			AggregateIDsFilter(existing.AggregateID, iamID)

		validation := checkExistingLoginPolicyIDPProviderValidation(provider)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyIDPProviderAdded, provider)
	}
}

func LoginPolicyIDPProviderRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, provider *iam_es_model.IDPProviderID, cascade bool) (*es_models.Aggregate, error) {
	if provider == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
	}
	agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	eventType := model.LoginPolicyIDPProviderRemoved
	if cascade {
		eventType = model.LoginPolicyIDPProviderCascadeRemoved
	}
	return agg.AppendEvent(eventType, provider)
}

func checkExistingLoginPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, event := range events {
			switch event.Type {
			case model.LoginPolicyAdded:
				existing = true
			case model.LoginPolicyRemoved:
				existing = false
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Nsh8u", "Errors.Org.LoginPolicy.AlreadyExists")
		}
		return nil
	}
}

func checkExistingLoginPolicyIDPProviderValidation(idpProvider *iam_es_model.IDPProvider) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		idpConfigs := make([]*iam_es_model.IDPConfig, 0)
		idps := make([]*iam_es_model.IDPProvider, 0)
		for _, event := range events {
			switch event.Type {
			case model.IDPConfigAdded, iam_es_model.IDPConfigAdded:
				config := new(iam_es_model.IDPConfig)
				config.SetData(event)
				idpConfigs = append(idpConfigs, config)
				if event.AggregateType == model.OrgAggregate {
					config.Type = int32(iam_model.IDPProviderTypeOrg)
				} else {
					config.Type = int32(iam_model.IDPProviderTypeSystem)
				}
			case model.IDPConfigRemoved, iam_es_model.IDPConfigRemoved:
				config := new(iam_es_model.IDPConfig)
				config.SetData(event)
				for i, p := range idpConfigs {
					if p.IDPConfigID == config.IDPConfigID {
						idpConfigs[i] = idpConfigs[len(idpConfigs)-1]
						idpConfigs[len(idpConfigs)-1] = nil
						idpConfigs = idpConfigs[:len(idpConfigs)-1]
					}
				}
			case model.LoginPolicyIDPProviderAdded:
				idp := new(iam_es_model.IDPProvider)
				idp.SetData(event)
			case model.LoginPolicyIDPProviderRemoved:
				idp := new(iam_es_model.IDPProvider)
				idp.SetData(event)
				for i, p := range idps {
					if p.IDPConfigID == idp.IDPConfigID {
						idps[i] = idps[len(idps)-1]
						idps[len(idps)-1] = nil
						idps = idps[:len(idps)-1]
					}
				}
			}
		}
		exists := false
		for _, p := range idpConfigs {
			if p.IDPConfigID == idpProvider.IDPConfigID && p.Type == idpProvider.Type {
				exists = true
			}
		}
		if !exists {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Djlo9", "Errors.Iam.IdpNotExisting")
		}
		for _, p := range idps {
			if p.IDPConfigID == idpProvider.IDPConfigID {
				return errors.ThrowPreconditionFailed(nil, "EVENT-us5Zw", "Errors.Org.LoginPolicy.IdpProviderAlreadyExisting")
			}
		}
		return nil
	}
}
