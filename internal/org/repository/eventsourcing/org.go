package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func OrgByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return OrgQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func OrgDomainUniqueQuery(domain string) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgDomainAggregate).
		AggregateIDFilter(domain).
		OrderDesc().
		SetLimit(1)
}

func OrgNameUniqueQuery(name string) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgNameAggregate).
		AggregateIDFilter(name).
		OrderDesc().
		SetLimit(1)
}

func OrgQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate).
		LatestSequenceFilter(latestSequence)
}

func OrgAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.OrgAggregate, model.OrgVersion, sequence, es_models.OverwriteResourceOwner(id))
}

func orgCreatedAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, org *model.Org, users func(context.Context, string) ([]*es_models.Aggregate, error)) (_ []*es_models.Aggregate, err error) {
	if org == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie7", "Errors.Internal")
	}

	agg, err := aggCreator.NewAggregate(ctx, org.AggregateID, model.OrgAggregate, model.OrgVersion, org.Sequence, es_models.OverwriteResourceOwner(org.AggregateID))
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.OrgAdded, org)
	if err != nil {
		return nil, err
	}
	aggregates := make([]*es_models.Aggregate, 0)
	aggregates, err = addDomainAggregateAndEvents(ctx, aggCreator, agg, aggregates, org, users)
	if err != nil {
		return nil, err
	}
	nameAggregate, err := reservedUniqueNameAggregate(ctx, aggCreator, org.AggregateID, org.Name)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, nameAggregate)
	return append(aggregates, agg), nil
}

func addDomainAggregateAndEvents(ctx context.Context, aggCreator *es_models.AggregateCreator, orgAggregate *es_models.Aggregate, aggregates []*es_models.Aggregate, org *model.Org, users func(context.Context, string) ([]*es_models.Aggregate, error)) ([]*es_models.Aggregate, error) {
	for _, domain := range org.Domains {
		orgAggregate, err := orgAggregate.AppendEvent(model.OrgDomainAdded, domain)
		if err != nil {
			return nil, err
		}
		if domain.Verified {
			domainAggregates, err := orgDomainVerified(ctx, aggCreator, orgAggregate, org, domain, users)
			if err != nil {
				return nil, err
			}
			aggregates = append(aggregates, domainAggregates...)
		}
		if domain.Primary {
			orgAggregate, err = orgAggregate.AppendEvent(model.OrgDomainPrimarySet, domain)
			if err != nil {
				return nil, err
			}
		}
	}
	return aggregates, nil
}

func OrgUpdateAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, updated *model.Org) ([]*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk83d", "Errors.Internal")
	}
	if updated == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
	}
	changes := existing.Changes(updated)
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-E0hc5", "Errors.NoChangesFound")
	}

	aggregates := make([]*es_models.Aggregate, 0, 3)

	if name, ok := changes["name"]; ok {
		nameAggregate, err := reservedUniqueNameAggregate(ctx, aggCreator, "", name.(string))
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, nameAggregate)
		nameReleasedAggregate, err := releasedUniqueNameAggregate(ctx, aggCreator, "", existing.Name)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, nameReleasedAggregate)
	}

	orgAggregate, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}

	orgAggregate, err = orgAggregate.AppendEvent(model.OrgChanged, changes)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, orgAggregate)

	return aggregates, nil
}

func orgDeactivateAggregate(aggCreator *es_models.AggregateCreator, org *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-R03z8", "Errors.Internal")
		}
		if org.State == int32(org_model.OrgStateInactive) {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-mcPH0", "Errors.Internal.AlreadyDeactivated")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.OrgDeactivated, nil)
	}
}

func orgReactivateAggregate(aggCreator *es_models.AggregateCreator, org *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-cTHLd", "Errors.Internal")
		}
		if org.State == int32(org_model.OrgStateActive) {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-pUSMs", "Errors.Org.AlreadyActive")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.OrgReactivated, nil)
	}
}

func reservedUniqueDomainAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, domain string) (*es_models.Aggregate, error) {
	aggregate, err := aggCreator.NewAggregate(ctx, domain, model.OrgDomainAggregate, model.OrgVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, domain, model.OrgDomainAggregate, model.OrgVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.OrgDomainReserved, nil)
	if err != nil {
		return nil, err
	}
	return aggregate.SetPrecondition(OrgDomainUniqueQuery(domain), isEventValidation(aggregate, model.OrgDomainReserved)), nil
}

func releasedUniqueDomainAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, domain string) (*es_models.Aggregate, error) {
	aggregate, err := aggCreator.NewAggregate(ctx, domain, model.OrgDomainAggregate, model.OrgVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, domain, model.OrgDomainAggregate, model.OrgVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.OrgDomainReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(OrgDomainUniqueQuery(domain), isEventValidation(aggregate, model.OrgDomainReleased)), nil
}

func reservedUniqueNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, name string) (aggregate *es_models.Aggregate, err error) {
	aggregate, err = aggCreator.NewAggregate(ctx, name, model.OrgNameAggregate, model.OrgVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, name, model.OrgNameAggregate, model.OrgVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.OrgNameReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(OrgNameUniqueQuery(name), isEventValidation(aggregate, model.OrgNameReserved)), nil
}

func releasedUniqueNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, name string) (aggregate *es_models.Aggregate, err error) {
	aggregate, err = aggCreator.NewAggregate(ctx, name, model.OrgNameAggregate, model.OrgVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, name, model.OrgNameAggregate, model.OrgVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.OrgNameReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(OrgNameUniqueQuery(name), isEventValidation(aggregate, model.OrgNameReleased)), nil
}

func OrgDomainAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if domain == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-OSid3", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgDomainAdded, domain)
	}
}

func OrgDomainValidationGeneratedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if domain == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-GD2gq", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgDomainVerificationAdded, domain)
	}
}

func OrgDomainValidationFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if domain == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-BHF52", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgDomainVerificationFailed, domain)
	}
}

func OrgDomainVerifiedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain, users func(context.Context, string) ([]*es_models.Aggregate, error)) ([]*es_models.Aggregate, error) {
	if domain == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-DHs7s", "Errors.Internal")
	}
	agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}

	aggregates, err := orgDomainVerified(ctx, aggCreator, agg, existing, domain, users)
	if err != nil {
		return nil, err
	}
	return append(aggregates, agg), nil
}

func orgDomainVerified(ctx context.Context, aggCreator *es_models.AggregateCreator, agg *es_models.Aggregate, existing *model.Org, domain *model.OrgDomain, users func(context.Context, string) ([]*es_models.Aggregate, error)) ([]*es_models.Aggregate, error) {
	agg, err := agg.AppendEvent(model.OrgDomainVerified, domain)
	if err != nil {
		return nil, err
	}
	domainAgregate, err := reservedUniqueDomainAggregate(ctx, aggCreator, existing.AggregateID, domain.Domain)
	if err != nil {
		return nil, err
	}
	aggregates := []*es_models.Aggregate{domainAgregate}
	if users != nil {
		userAggregates, err := users(ctx, domain.Domain)
		if err != nil {
			return nil, errors.ThrowPreconditionFailed(err, "EVENT-HBwsw", "Errors.Internal")
		}
		aggregates = append(aggregates, userAggregates...)
	}
	return aggregates, nil
}

func OrgDomainSetPrimaryAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if domain == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-PSw3j", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgDomainPrimarySet, domain)
	}
}

func OrgDomainRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, domain *model.OrgDomain) ([]*es_models.Aggregate, error) {
	if domain == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-si8dW", "Errors.Internal")
	}
	aggregates := make([]*es_models.Aggregate, 0, 2)
	agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.OrgDomainRemoved, domain)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, agg)
	domainAgregate, err := releasedUniqueDomainAggregate(ctx, aggCreator, existing.AggregateID, domain.Domain)
	if err != nil {
		return nil, err
	}
	return append(aggregates, domainAgregate), nil
}

func IDPConfigAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, idp *iam_es_model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-MSki8", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.IDPConfigAdded, idp)
		if idp.OIDCIDPConfig != nil {
			agg.AppendEvent(model.OIDCIDPConfigAdded, idp.OIDCIDPConfig)
		}
		return agg, nil
	}
}

func IDPConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, idp *iam_es_model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Akdi8", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, i := range existing.IDPs {
			if i.IDPConfigID == idp.IDPConfigID {
				changes = i.Changes(idp)
			}
		}
		agg.AppendEvent(model.IDPConfigChanged, changes)

		return agg, nil
	}
}

func IDPConfigRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.Org, idp *iam_es_model.IDPConfig, provider *iam_es_model.IDPProvider) (*es_models.Aggregate, error) {
	if idp == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlso9", "Errors.Internal")
	}
	agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.IDPConfigRemoved, &iam_es_model.IDPConfigID{IDPConfigID: idp.IDPConfigID})
	if err != nil {
		return nil, err
	}
	if provider != nil {
		return agg.AppendEvent(model.LoginPolicyIDPProviderCascadeRemoved, &iam_es_model.IDPProviderID{IDPConfigID: provider.IDPConfigID})
	}
	return agg, nil

}

func IDPConfigDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, idp *iam_es_model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-3sz7d", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.IDPConfigDeactivated, &iam_es_model.IDPConfigID{IDPConfigID: idp.IDPConfigID})

		return agg, nil
	}
}

func IDPConfigReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, idp *iam_es_model.IDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if idp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4jdiS", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.IDPConfigReactivated, &iam_es_model.IDPConfigID{IDPConfigID: idp.IDPConfigID})

		return agg, nil
	}
}

func OIDCIDPConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, config *iam_es_model.OIDCIDPConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if config == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-6Fjso", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, idp := range existing.IDPs {
			if idp.IDPConfigID == config.IDPConfigID {
				if idp.OIDCIDPConfig != nil {
					changes = idp.OIDCIDPConfig.Changes(config)
				}
			}
		}
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailedf(nil, "EVENT-Rks9u", "Errors.NoChangesFound")
		}
		agg.AppendEvent(model.OIDCIDPConfigChanged, changes)

		return agg, nil
	}
}

func isEventValidation(aggregate *es_models.Aggregate, eventType es_models.EventType) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == eventType {
			return errors.ThrowPreconditionFailedf(nil, "EVENT-eJQqe", "user is already %v", eventType)
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}
}
