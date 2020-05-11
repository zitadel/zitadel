package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/model"
)

func OrgMemberAddedAggregate(aggCreator *es_models.AggregateCreator, member *OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-c63Ap", "member must not be nil")
		}

		aggregate, err := aggCreator.NewAggregate(ctx, member.AggregateID, model.OrgAggregate, orgVersion, member.Sequence)
		if err != nil {
			return nil, err
		}

		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter("org", "user").
			AggregateIDsFilter(member.AggregateID, member.UserID)

		validation := addMemberValidation(aggregate)

		return aggregate.SetPrecondition(validationQuery, validation).AppendEvent(model.OrgMemberAdded, member)
	}
}

func OrgMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existingMember *OrgMember, member *OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil || existingMember == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d34fs", "member should not be nil")
		}

		changes := existingMember.Changes(member)
		if len(changes) == 0 {
			return nil, errors.ThrowInvalidArgument(nil, "EVENT-VLMGn", "nothing changed")
		}

		agg, err := OrgAggregate(ctx, aggCreator, existingMember.AggregateID, existingMember.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgMemberChanged, changes)
	}
}

func OrgMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, member *OrgMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dieu7", "member should not be nil")
		}

		agg, err := OrgAggregate(ctx, aggCreator, member.AggregateID, member.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgMemberRemoved, member)
	}
}

func addMemberValidation(aggregate *es_models.Aggregate) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existsOrg := false
		existsUser := false
		for _, event := range events {
			switch event.AggregateType {
			case "user":
				existsUser = true
			case "org":
				aggregate.PreviousSequence = event.Sequence
				existsOrg = true
			}
		}
		if existsOrg && existsUser {
			return nil
		}
		return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "conditions not met")
	}
}
