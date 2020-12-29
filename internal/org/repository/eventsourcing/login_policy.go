package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func LoginPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, policy *iam_es_model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smla8", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(org.AggregateID)

		validation := checkExistingLoginPolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyAdded, policy)
	}
}

func LoginPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, policy *iam_es_model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlco9", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		changes := org.LoginPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.LoginPolicyChanged, changes)
	}
}

func LoginPolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-S8sio", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyRemoved, nil)
	}
}

func LoginPolicyIDPProviderAddedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, provider *iam_es_model.IDPProvider, iamID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
			AggregateIDsFilter(org.AggregateID, iamID)

		validation := checkExistingLoginPolicyIDPProviderValidation(provider)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyIDPProviderAdded, provider)
	}
}

func LoginPolicyIDPProviderRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, org *model.Org, provider *iam_es_model.IDPProviderID, cascade bool) (*es_models.Aggregate, error) {
	if provider == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
	}
	agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
	if err != nil {
		return nil, err
	}
	eventType := model.LoginPolicyIDPProviderRemoved
	if cascade {
		eventType = model.LoginPolicyIDPProviderCascadeRemoved
	}
	return agg.AppendEvent(eventType, provider)
}

func LoginPolicySecondFactorAddedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, mfa *iam_es_model.MFA, iamID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-5Gk9s", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDsFilter(org.AggregateID)

		validation := checkExistingLoginPolicySecondFactorValidation(mfa.MFAType)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicySecondFactorAdded, mfa)
	}
}

func LoginPolicySecondFactorRemovedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, mfa *iam_es_model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicySecondFactorRemoved, mfa)
	}
}

func LoginPolicyMultiFactorAddedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, mfa *iam_es_model.MFA, iamID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4Bm9s", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDsFilter(org.AggregateID)

		validation := checkExistingLoginPolicyMultiFactorValidation(mfa.MFAType)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyMultiFactorAdded, mfa)
	}
}

func LoginPolicyMultiFactorRemovedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, mfa *iam_es_model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-6Nm9s", "Errors.Internal")
		}

		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyMultiFactorRemoved, mfa)
	}
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
				err := config.SetData(event)
				if err != nil {
					return err
				}
				idpConfigs = append(idpConfigs, config)
				if event.AggregateType == model.OrgAggregate {
					config.Type = int32(iam_model.IDPProviderTypeOrg)
				} else {
					config.Type = int32(iam_model.IDPProviderTypeSystem)
				}
			case model.IDPConfigRemoved, iam_es_model.IDPConfigRemoved:
				config := new(iam_es_model.IDPConfig)
				err := config.SetData(event)
				if err != nil {
					return err
				}
				for i := len(idpConfigs) - 1; i >= 0; i-- {
					if idpConfigs[i].IDPConfigID == config.IDPConfigID {
						idpConfigs[i] = idpConfigs[len(idpConfigs)-1]
						idpConfigs[len(idpConfigs)-1] = nil
						idpConfigs = idpConfigs[:len(idpConfigs)-1]
						break
					}
				}
			case model.LoginPolicyIDPProviderAdded:
				idp := new(iam_es_model.IDPProvider)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				idps = append(idps, idp)
			case model.LoginPolicyIDPProviderRemoved, model.LoginPolicyIDPProviderCascadeRemoved:
				idp := new(iam_es_model.IDPProvider)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				for i := len(idps) - 1; i >= 0; i-- {
					if idps[i].IDPConfigID == idp.IDPConfigID {
						idps[i] = idps[len(idps)-1]
						idps[len(idps)-1] = nil
						idps = idps[:len(idps)-1]
						break
					}
				}
			case model.LoginPolicyRemoved:
				idps = make([]*iam_es_model.IDPProvider, 0)
			}
		}
		exists := false
		for _, p := range idpConfigs {
			if p.IDPConfigID == idpProvider.IDPConfigID && p.Type == idpProvider.Type {
				exists = true
			}
		}
		if !exists {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Djlo9", "Errors.IAM.IdpNotExisting")
		}
		for _, p := range idps {
			if p.IDPConfigID == idpProvider.IDPConfigID {
				return errors.ThrowPreconditionFailed(nil, "EVENT-us5Zw", "Errors.Org.LoginPolicy.IdpProviderAlreadyExisting")
			}
		}
		return nil
	}
}

func checkExistingLoginPolicySecondFactorValidation(mfaType int32) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		mfas := make([]int32, 0)
		for _, event := range events {
			switch event.Type {
			case model.LoginPolicySecondFactorAdded:
				mfa := new(iam_es_model.MFA)
				err := mfa.SetData(event)
				if err != nil {
					return err
				}
				mfas = append(mfas, mfa.MFAType)
			case model.LoginPolicySecondFactorRemoved:
				idp := new(iam_es_model.IDPProvider)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				for i := len(mfas) - 1; i >= 0; i-- {
					if mfas[i] == mfaType {
						mfas[i] = mfas[len(mfas)-1]
						mfas[len(mfas)-1] = 0
						mfas = mfas[:len(mfas)-1]
						break
					}
				}
			case model.LoginPolicyRemoved:
				mfas = make([]int32, 0)
			}
		}

		for _, m := range mfas {
			if m == mfaType {
				return errors.ThrowPreconditionFailed(nil, "EVENT-4Bo0sw", "Errors.Org.LoginPolicy.MFA.AlreadyExisting")
			}
		}
		return nil
	}
}

func checkExistingLoginPolicyMultiFactorValidation(mfaType int32) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		mfas := make([]int32, 0)
		for _, event := range events {
			switch event.Type {
			case model.LoginPolicyMultiFactorAdded:
				mfa := new(iam_es_model.MFA)
				err := mfa.SetData(event)
				if err != nil {
					return err
				}
				mfas = append(mfas, mfa.MFAType)
			case model.LoginPolicyMultiFactorRemoved:
				idp := new(iam_es_model.IDPProvider)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				for i := len(mfas) - 1; i >= 0; i-- {
					if mfas[i] == mfaType {
						mfas[i] = mfas[len(mfas)-1]
						mfas[len(mfas)-1] = 0
						mfas = mfas[:len(mfas)-1]
						break
					}
				}
			case model.LoginPolicyRemoved:
				mfas = make([]int32, 0)
			}
		}

		for _, m := range mfas {
			if m == mfaType {
				return errors.ThrowPreconditionFailed(nil, "EVENT-4Bo0sw", "Errors.Org.LoginPolicy.MFA.AlreadyExisting")
			}
		}
		return nil
	}
}
