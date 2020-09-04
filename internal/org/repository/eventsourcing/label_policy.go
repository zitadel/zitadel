package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

// ToDo Michi
func LabelPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.LabelPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-TUWod", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLabelPolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LabelPolicyAdded, policy)
	}
}

func LabelPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.LabelPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-unRI2", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.LabelPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Tz130", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.LabelPolicyChanged, changes)
		//		return nil, nil
	}
}

func LabelPolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-v7E9b", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LabelPolicyRemoved, nil)
	}
}

func LabelPolicyIDPProviderAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, provider *iam_es_model.IDPProvider, iamID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-QGVxo", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
			AggregateIDsFilter(existing.AggregateID, iamID)

		validation := checkExistingLabelPolicyIDPProviderValidation(provider)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LabelPolicyIDPProviderAdded, provider)
	}
}

func LabelPolicyIDPProviderRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, provider *iam_es_model.IDPProviderID, cascade bool) (*es_models.Aggregate, error) {
	if provider == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-TgLUZ", "Errors.Internal")
	}
	agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	eventType := model.LabelPolicyIDPProviderRemoved
	if cascade {
		eventType = model.LabelPolicyIDPProviderCascadeRemoved
	}
	return agg.AppendEvent(eventType, provider)
}

func checkExistingLabelPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, event := range events {
			switch event.Type {
			case model.LabelPolicyAdded:
				existing = true
			case model.LabelPolicyRemoved:
				existing = false
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-g9mCI", "Errors.Org.LabelPolicy.AlreadyExists")
		}
		return nil
	}
}

func checkExistingLabelPolicyIDPProviderValidation(idpProvider *iam_es_model.IDPProvider) func(...*es_models.Event) error {
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
			case model.LabelPolicyIDPProviderAdded:
				idp := new(iam_es_model.IDPProvider)
				idp.SetData(event)
			case model.LabelPolicyIDPProviderRemoved:
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
			return errors.ThrowPreconditionFailed(nil, "EVENT-Er3po9", "Errors.IAM.IdpNotExisting")
		}
		for _, p := range idps {
			if p.IDPConfigID == idpProvider.IDPConfigID {
				return errors.ThrowPreconditionFailed(nil, "EVENT-VSWO5", "Errors.Org.LabelPolicy.IdpProviderAlreadyExisting")
			}
		}
		return nil
	}
}
