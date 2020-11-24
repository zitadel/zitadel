package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func orgMemberAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, member *model.OrgMember, resourceOwner string) (agg *es_models.Aggregate, err error) {
	if member == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-c63Ap", "Errors.Internal")
	}

	if resourceOwner != "" {
		agg, err = aggCreator.NewAggregate(ctx, member.AggregateID, model.OrgAggregate, model.OrgVersion, member.Sequence, es_models.OverwriteResourceOwner(resourceOwner))
	} else {
		agg, err = aggCreator.NewAggregate(ctx, member.AggregateID, model.OrgAggregate, model.OrgVersion, member.Sequence)
	}
	if err != nil {
		return nil, err
	}

	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, usr_model.UserAggregate).
		AggregateIDsFilter(member.AggregateID, member.UserID)

	validation := addMemberValidation(agg, member)
	return agg.SetPrecondition(validationQuery, validation).AppendEvent(model.OrgMemberAdded, member)
}

func orgMemberChangedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, existingMember, member *model.OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil || org == nil || existingMember == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d34fs", "Errors.Internal")
		}

		changes := existingMember.Changes(member)
		if len(changes) == 0 {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-VLMGn", "Errors.NoChangesFound")
		}

		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgMemberChanged, changes)
	}
}

func orgMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, member *model.OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-vNPVX", "Errors.Internal")
		}

		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgMemberRemoved, member)
	}
}

func addMemberValidation(aggregate *es_models.Aggregate, member *model.OrgMember) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existsOrg := false
		existsUser := false
		isMember := false
		for _, event := range events {
			switch event.AggregateType {
			case usr_model.UserAggregate:
				existsUser = true
			case model.OrgAggregate:
				aggregate.PreviousSequence = event.Sequence
				existsOrg = true
				switch event.Type {
				case model.OrgMemberAdded, model.OrgMemberRemoved:
					manipulatedMember, err := model.OrgMemberFromEvent(new(model.OrgMember), event)
					if err != nil {
						return errors.ThrowInternal(err, "EVENT-Eg8St", "unable to validate object")
					}
					if manipulatedMember.UserID == member.UserID {
						isMember = event.Type == model.OrgMemberAdded
					}
				}
			}
		}
		if existsOrg && existsUser && !isMember {
			return nil
		}
		return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "conditions not met")
	}
}
