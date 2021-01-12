package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IAMByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0soe4", "Errors.IAM.IDMissing")
	}
	return IAMQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IAMQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(latestSequence)
}

func IAMAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.IAM) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo04e", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IAMAggregate, model.IAMVersion, iam.Sequence)
}

func IAMAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.IAM, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "Errors.Internal")
	}

	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IAMAggregate, model.IAMVersion, iam.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func IAMSetupStartedAggregate(aggCreator *es_models.AggregateCreator, iam *model.IAM) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IAMAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IAMSetupStarted, &struct{ Step model.Step }{Step: iam.SetUpStarted})
	}
}

func IAMSetupDoneAggregate(aggCreator *es_models.AggregateCreator, iam *model.IAM) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IAMAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.IAMSetupDone, &struct{ Step model.Step }{Step: iam.SetUpDone})
	}
}

func IAMSetupDoneEvent(ctx context.Context, agg *es_models.Aggregate, iam *model.IAM) (*es_models.Aggregate, error) {
	return agg.AppendEvent(model.IAMSetupDone, &struct{ Step model.Step }{Step: iam.SetUpDone})
}

func IAMSetGlobalOrgAggregate(aggCreator *es_models.AggregateCreator, iam *model.IAM, globalOrg string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if globalOrg == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8siwa", "Errors.IAM.GlobalOrgMissing")
		}
		agg, err := IAMAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.GlobalOrgSet, &model.IAM{GlobalOrgID: globalOrg})
	}
}

func IAMSetIamProjectAggregate(aggCreator *es_models.AggregateCreator, iam *model.IAM, projectID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if projectID == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sjuw3", "Errors.IAM.IAMProjectIDMissing")
		}
		agg, err := IAMAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IAMProjectSet, &model.IAM{IAMProjectID: projectID})
	}
}

func IAMMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existingIAM *model.IAM, member *model.IAMMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9sope", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existingIAM)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IAMMemberAdded, member)
	}
}

func IAMMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existingIAM *model.IAM, member *model.IAMMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-38skf", "Errors.Internal")
		}

		agg, err := IAMAggregate(ctx, aggCreator, existingIAM)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IAMMemberChanged, member)
	}
}

func IAMMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existingIAM *model.IAM, member *model.IAMMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-90lsw", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existingIAM)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IAMMemberRemoved, member)
	}
}

func IDPConfigAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, idp *model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-MSn7d", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg, err = agg.AppendEvent(model.IDPConfigAdded, idp)
		if err != nil {
			return nil, err
		}
		if idp.OIDCIDPConfig != nil {
			return agg.AppendEvent(model.OIDCIDPConfigAdded, idp.OIDCIDPConfig)
		}
		return agg, nil
	}
}

func IDPConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, idp *model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Amc7s", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, i := range existing.IDPs {
			if i.IDPConfigID == idp.IDPConfigID {
				changes = i.Changes(idp)
			}
		}
		return agg.AppendEvent(model.IDPConfigChanged, changes)
	}
}

func IDPConfigRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, idp *model.IDPConfig, provider *model.IDPProvider) (*es_models.Aggregate, error) {
	if idp == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-se23g", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.IDPConfigRemoved, &model.IDPConfigID{IDPConfigID: idp.IDPConfigID})
	if err != nil {
		return nil, err
	}
	if provider != nil {
		return agg.AppendEvent(model.LoginPolicyIDPProviderCascadeRemoved, &model.IDPConfigID{IDPConfigID: idp.IDPConfigID})
	}
	return agg, nil
}

func IDPConfigDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, idp *model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slfi3", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IDPConfigDeactivated, &model.IDPConfigID{IDPConfigID: idp.IDPConfigID})
	}
}

func IDPConfigReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, idp *model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IDPConfigReactivated, &model.IDPConfigID{IDPConfigID: idp.IDPConfigID})
	}
}

func OIDCIDPConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, config *model.OIDCIDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if config == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, idp := range existing.IDPs {
			if idp.IDPConfigID == config.IDPConfigID && idp.OIDCIDPConfig != nil {
				changes = idp.OIDCIDPConfig.Changes(config)
			}
		}
		if len(changes) <= 1 {
			return nil, errors.ThrowPreconditionFailedf(nil, "EVENT-Cml9s", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.OIDCIDPConfigChanged, changes)
	}
}
func LabelPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.LabelPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-e248Y", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.IAMAggregate).
			EventTypesFilter(model.LabelPolicyAdded).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLabelPolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LabelPolicyAdded, policy)
	}
}

func checkExistingLabelPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.LabelPolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-KyLIK", "Errors.IAM.LabelPolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func LabelPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.LabelPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-uP6HQ", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultLabelPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-hZE24", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.LabelPolicyChanged, changes)
	}
}

func LoginPolicyAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.LoginPolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smla8", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		EventTypesFilter(model.LoginPolicyAdded).
		AggregateIDFilter(existing.AggregateID)

	validation := checkExistingLoginPolicyValidation()
	agg.SetPrecondition(validationQuery, validation)
	return agg.AppendEvent(model.LoginPolicyAdded, policy)
}

func LoginPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlco9", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultLoginPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.LoginPolicyChanged, changes)
	}
}

func LoginPolicyIDPProviderAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, provider *model.IDPProvider) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.IAMAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLoginPolicyIDPProviderValidation(provider.IDPConfigID)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyIDPProviderAdded, provider)
	}
}

func LoginPolicyIDPProviderRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, provider *model.IDPProviderID) (*es_models.Aggregate, error) {
	if provider == nil || existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.LoginPolicyIDPProviderRemoved, provider)
}

func LoginPolicySecondFactorAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, mfa *model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4Gm9s", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.IAMAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLoginPolicySecondFactorValidation(mfa.MFAType)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicySecondFactorAdded, mfa)
	}
}

func LoginPolicySecondFactorRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, mfa *model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-5Bm9s", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicySecondFactorRemoved, mfa)
	}
}

func LoginPolicyMultiFactorAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, mfa *model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4Gm9s", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.IAMAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingLoginPolicyMultiFactorValidation(mfa.MFAType)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.LoginPolicyMultiFactorAdded, mfa)
	}
}

func LoginPolicyMultiFactorRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, mfa *model.MFA) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mfa == nil || existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-6Mso9", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyMultiFactorRemoved, mfa)
	}
}

func PasswordComplexityPolicyAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordComplexityPolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smla8", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		EventTypesFilter(model.PasswordComplexityPolicyAdded).
		AggregateIDFilter(existing.AggregateID)

	validation := checkExistingPasswordComplexityPolicyValidation()
	agg.SetPrecondition(validationQuery, validation)
	return agg.AppendEvent(model.PasswordComplexityPolicyAdded, policy)
}

func PasswordComplexityPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordComplexityPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlco9", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultPasswordComplexityPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.PasswordComplexityPolicyChanged, changes)
	}
}

func PasswordAgePolicyAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordAgePolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-T7sui", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		EventTypesFilter(model.PasswordAgePolicyAdded).
		AggregateIDFilter(existing.AggregateID)

	validation := checkExistingPasswordAgePolicyValidation()
	agg.SetPrecondition(validationQuery, validation)
	return agg.AppendEvent(model.PasswordAgePolicyAdded, policy)
}

func PasswordAgePolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-3Gs0o", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultPasswordAgePolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-3Wdos", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.PasswordAgePolicyChanged, changes)
	}
}

func PasswordLockoutPolicyAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordLockoutPolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-w5Tds", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		EventTypesFilter(model.PasswordLockoutPolicyAdded).
		AggregateIDFilter(existing.AggregateID)

	validation := checkExistingPasswordLockoutPolicyValidation()
	agg.SetPrecondition(validationQuery, validation)
	return agg.AppendEvent(model.PasswordLockoutPolicyAdded, policy)
}

func PasswordLockoutPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2D0fs", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultPasswordLockoutPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-7Hsk9", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.PasswordLockoutPolicyChanged, changes)
	}
}

func OrgIAMPolicyAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.OrgIAMPolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-w5Tds", "Errors.Internal")
	}
	agg, err := IAMAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		EventTypesFilter(model.OrgIAMPolicyAdded).
		AggregateIDFilter(existing.AggregateID)

	validation := checkExistingOrgIAMPolicyValidation()
	agg.SetPrecondition(validationQuery, validation)
	return agg.AppendEvent(model.OrgIAMPolicyAdded, policy)
}

func OrgIAMPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.IAM, policy *model.OrgIAMPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2D0fs", "Errors.Internal")
		}
		agg, err := IAMAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.DefaultOrgIAMPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-7Hsk9", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.OrgIAMPolicyChanged, changes)
	}
}

func checkExistingLoginPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.LoginPolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-Ski9d", "Errors.IAM.LoginPolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func checkExistingPasswordComplexityPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.PasswordComplexityPolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-Ski9d", "Errors.IAM.PasswordComplexityPolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func checkExistingPasswordAgePolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.PasswordAgePolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-Ski9d", "Errors.IAM.PasswordAgePolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func checkExistingPasswordLockoutPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.PasswordLockoutPolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-Ski9d", "Errors.IAM.PasswordLockoutPolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func checkExistingOrgIAMPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		for _, event := range events {
			switch event.Type {
			case model.OrgIAMPolicyAdded:
				return errors.ThrowPreconditionFailed(nil, "EVENT-bSm8f", "Errors.IAM.OrgIAMPolicy.AlreadyExists")
			}
		}
		return nil
	}
}

func checkExistingLoginPolicyIDPProviderValidation(idpConfigID string) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		idpConfigs := make([]*model.IDPConfig, 0)
		idps := make([]*model.IDPProvider, 0)
		for _, event := range events {
			switch event.Type {
			case model.IDPConfigAdded:
				config := new(model.IDPConfig)
				err := config.SetData(event)
				if err != nil {
					return err
				}
				idpConfigs = append(idpConfigs, config)
			case model.IDPConfigRemoved:
				config := new(model.IDPConfig)
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
				idp := new(model.IDPProvider)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				idps = append(idps, idp)
			case model.LoginPolicyIDPProviderRemoved:
				idp := new(model.IDPProvider)
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
			}
		}
		exists := false
		for _, p := range idpConfigs {
			if p.IDPConfigID == idpConfigID {
				exists = true
			}
		}
		if !exists {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Djlo9", "Errors.IAM.IdpNotExisting")
		}
		for _, p := range idps {
			if p.IDPConfigID == idpConfigID {
				return errors.ThrowPreconditionFailed(nil, "EVENT-us5Zw", "Errors.IAM.LoginPolicy.IdpProviderAlreadyExisting")
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
				idp := new(model.MFA)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				mfas = append(mfas, idp.MFAType)
			case model.LoginPolicySecondFactorRemoved:
				mfa := new(model.MFA)
				err := mfa.SetData(event)
				if err != nil {
					return err
				}
				for i := len(mfas) - 1; i >= 0; i-- {
					if mfas[i] == mfa.MFAType {
						mfas[i] = mfas[len(mfas)-1]
						mfas[len(mfas)-1] = 0
						mfas = mfas[:len(mfas)-1]
						break
					}
				}
			}
		}
		for _, m := range mfas {
			if m == mfaType {
				return errors.ThrowPreconditionFailed(nil, "EVENT-3vmHd", "Errors.IAM.LoginPolicy.MFA.AlreadyExisting")
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
				idp := new(model.MFA)
				err := idp.SetData(event)
				if err != nil {
					return err
				}
				mfas = append(mfas, idp.MFAType)
			case model.LoginPolicyMultiFactorRemoved:
				mfa := new(model.MFA)
				err := mfa.SetData(event)
				if err != nil {
					return err
				}
				for i := len(mfas) - 1; i >= 0; i-- {
					if mfas[i] == mfa.MFAType {
						mfas[i] = mfas[len(mfas)-1]
						mfas[len(mfas)-1] = 0
						mfas = mfas[:len(mfas)-1]
						break
					}
				}
			}
		}
		for _, m := range mfas {
			if m == mfaType {
				return errors.ThrowPreconditionFailed(nil, "EVENT-6Hsj89", "Errors.IAM.LoginPolicy.MFA.AlreadyExisting")
			}
		}
		return nil
	}
}
