package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/model"
)

func AddOrgMemberAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existingOrg *Org, addedMember *OrgMember) (*es_models.Aggregate, error) {
	if existingOrg == nil || addedMember == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-c63Ap", "members must not be nil")
	}

	aggregate, err := aggCreator.NewAggregate(ctx, existingOrg.AggregateID, model.OrgAggregate, orgVersion, existingOrg.Sequence)
	if err != nil {
		return nil, err
	}
	//TODO: add user exists precondition

	return aggregate.AppendEvent(model.OrgMemberAdded, addedMember)
}

func ChangeOrgMemberAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existingMember, updatedMember *OrgMember) (*es_models.Aggregate, error) {
	if existingMember == nil || updatedMember == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-nnxKg", "members must not be nil")
	}

	changes := existingMember.Changes(updatedMember)
	aggregate, err := aggCreator.NewAggregate(ctx, existingMember.AggregateID, model.OrgAggregate, orgVersion, existingMember.Sequence)
	if err != nil {
		return nil, err
	}

	return aggregate.AppendEvent(model.OrgMemberChanged, changes)
}

func RemoveOrgMemberAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existingMember *OrgMember) (*es_models.Aggregate, error) {
	if existingMember == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-knnVI", "member must not be nil")
	}

	aggregate, err := aggCreator.NewAggregate(ctx, existingMember.AggregateID, model.OrgAggregate, orgVersion, existingMember.Sequence)
	if err != nil {
		return nil, err
	}

	return aggregate.AppendEvent(model.OrgMemberRemoved, nil)
}
