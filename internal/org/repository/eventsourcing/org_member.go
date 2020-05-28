package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func orgMemberAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, member *OrgMember) (*es_models.Aggregate, error) {
	if member == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-c63Ap", "member must not be nil")
	}

	aggregate, err := aggCreator.NewAggregate(ctx, member.AggregateID, org_model.OrgAggregate, orgVersion, member.Sequence)
	if err != nil {
		return nil, err
	}

	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(org_model.OrgAggregate, usr_model.UserAggregate).
		AggregateIDsFilter(member.AggregateID, member.UserID)

	validation := addMemberValidation(aggregate, member)

	return aggregate.SetPrecondition(validationQuery, validation).AppendEvent(org_model.OrgMemberAdded, member)
}

func orgMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existingMember *OrgMember, member *OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil || existingMember == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d34fs", "member must not be nil")
		}

		changes := existingMember.Changes(member)
		if len(changes) == 0 {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-VLMGn", "nothing changed")
		}

		agg, err := OrgAggregate(ctx, aggCreator, existingMember.AggregateID, existingMember.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(org_model.OrgMemberChanged, changes)
	}
}

func orgMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, member *OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dieu7", "member must not be nil")
		}

		agg, err := OrgAggregate(ctx, aggCreator, member.AggregateID, member.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(org_model.OrgMemberRemoved, member)
	}
}

func addMemberValidation(aggregate *es_models.Aggregate, member *OrgMember) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existsOrg := false
		existsUser := false
		isMember := false
		for _, event := range events {
			switch event.AggregateType {
			case usr_model.UserAggregate:
				existsUser = true
			case org_model.OrgAggregate:
				aggregate.PreviousSequence = event.Sequence
				existsOrg = true
				switch event.Type {
				case org_model.OrgMemberAdded, org_model.OrgMemberRemoved:
					manipulatedMember, err := OrgMemberFromEvent(new(OrgMember), event)
					if err != nil {
						return errors.ThrowInternal(err, "EVENT-Eg8St", "unable to validate object")
					}
					if manipulatedMember.UserID == member.UserID {
						isMember = event.Type == org_model.OrgMemberAdded
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
