package command

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddInstanceMember(a *instance.Aggregate, userID string, roles ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTA-SDSfs", "Errors.Invalid.Argument")
		}
		// TODO: check roles
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				if exists, err := ExistsUser(ctx, filter, userID, ""); err != nil || !exists {
					return nil, errors.ThrowNotFound(err, "INSTA-GSXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsInstanceMember(ctx, filter, a.ID, userID); err != nil || isMember {
					return nil, errors.ThrowAlreadyExists(err, "INSTA-pFDwe", "Errors.Instance.Member.AlreadyExists")
				}
				return []eventstore.Command{instance.NewMemberAddedEvent(ctx, &a.Aggregate, userID, roles...)}, nil
			},
			nil
	}
}

func IsInstanceMember(ctx context.Context, filter preparation.FilterToQueryReducer, instanceID, userID string) (isMember bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderAsc().
		AddQuery().
		AggregateIDs(instanceID).
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.MemberAddedEventType,
			instance.MemberRemovedEventType,
			instance.MemberCascadeRemovedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch e := event.(type) {
		case *instance.MemberAddedEvent:
			if e.UserID == userID {
				isMember = true
			}
		case *instance.MemberRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		case *instance.MemberCascadeRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		}
	}

	return isMember, nil
}
