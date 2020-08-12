package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IamByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0soe4", "Errors.Iam.IDMissing")
	}
	return IamQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IamQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate).
		LatestSequenceFilter(latestSequence)
}

func IamAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo04e", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence)
}

func IamAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "Errors.Internal")
	}

	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func IamSetupStartedAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamSetupStarted, nil)
	}
}

func IamSetupDoneAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.IamSetupDone, nil)
	}
}

func IamSetGlobalOrgAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam, globalOrg string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if globalOrg == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8siwa", "Errors.Iam.GlobalOrgMissing")
		}
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.GlobalOrgSet, &model.Iam{GlobalOrgID: globalOrg})
	}
}

func IamSetIamProjectAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam, projectID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if projectID == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sjuw3", "Errors.Iam.IamProjectIDMisisng")
		}
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamProjectSet, &model.Iam{IamProjectID: projectID})
	}
}

func IamMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9sope", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberAdded, member)
	}
}

func IamMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-38skf", "Errors.Internal")
		}

		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberChanged, member)
	}
}

func IamMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-90lsw", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberRemoved, member)
	}
}

func IdpConfigurationAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, idp *model.IdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-MSn7d", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg, err = agg.AppendEvent(model.IdpConfigAdded, idp)
		if err != nil {
			return nil, err
		}
		if idp.OIDCIDPConfig != nil {
			return agg.AppendEvent(model.OidcIdpConfigAdded, idp.OIDCIDPConfig)
		}
		return agg, nil
	}
}

func IdpConfigurationChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, idp *model.IdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Amc7s", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, i := range existing.IDPs {
			if i.IDPConfigID == idp.IDPConfigID {
				changes = i.Changes(idp)
			}
		}
		return agg.AppendEvent(model.IdpConfigChanged, changes)
	}
}

func IdpConfigurationRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, idp *model.IdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-se23g", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IdpConfigRemoved, &model.IdpConfigID{IdpConfigID: idp.IDPConfigID})
	}
}

func IdpConfigurationDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, idp *model.IdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slfi3", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IdpConfigDeactivated, &model.IdpConfigID{IdpConfigID: idp.IDPConfigID})
	}
}

func IdpConfigurationReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, idp *model.IdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IdpConfigReactivated, &model.IdpConfigID{IdpConfigID: idp.IDPConfigID})
	}
}

func OIDCIdpConfigurationChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, config *model.OidcIdpConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if config == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, idp := range existing.IDPs {
			if idp.IDPConfigID == config.IdpConfigID {
				if idp.OIDCIDPConfig != nil {
					changes = idp.OIDCIDPConfig.Changes(config)
				}
			}
		}
		return agg.AppendEvent(model.OidcIdpConfigChanged, changes)
	}
}

func LoginPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, policy *model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smla8", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyAdded, policy)
	}
}

func LoginPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, policy *model.LoginPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlco9", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
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

func LoginPolicyIdpProviderAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, provider *model.IdpProvider) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		//TODO: Check if idp provider existing
		return agg.AppendEvent(model.LoginPolicyIdpProviderAdded, provider)
	}
}

func LoginPolicyIdpProviderRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, provider *model.IdpProviderID) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if provider == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sml9d", "Errors.Internal")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.LoginPolicyIdpProviderRemoved, provider)
	}
}
