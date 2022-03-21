package org

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/command/v2/user"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func AddMemberCommand(a *org.Aggregate, userID string, roles ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument")
		}
		// TODO: check roles
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				if exists, err := user.ExistsUser(ctx, filter, userID, a.ID); err != nil || !exists {
					return nil, errors.ThrowNotFound(err, "ORG-GoXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsMember(ctx, filter, a.ID, userID); err != nil || isMember {
					return nil, errors.ThrowNotFound(err, "ORG-poWwe", "Errors.Org.Member.AlreadyExists")
				}
				return []eventstore.Command{org.NewMemberAddedEvent(ctx, &a.Aggregate, userID, roles...)}, nil
			},
			nil
	}
}

func IsMember(ctx context.Context, filter preparation.FilterToQueryReducer, orgID, userID string) (isMember bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateIDs(orgID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.MemberAddedEventType,
			org.MemberRemovedEventType,
			org.MemberCascadeRemovedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		id := struct {
			ID string `json:"userId"`
		}{}
		switch event.(type) {
		case *org.MemberAddedEvent:
			err = json.Unmarshal(event.DataAsBytes(), &id)
			if err != nil {
				return false, errors.ThrowInternal(err, "ORG-WcoNL", "Errors.Internal")
			}
			if id.ID == userID {
				isMember = true
			}
		case *org.MemberRemovedEvent, *org.MemberCascadeRemovedEvent:
			err = json.Unmarshal(event.DataAsBytes(), &id)
			if err != nil {
				return false, errors.ThrowInternal(err, "ORG-hJkpr", "Errors.Internal")
			}
			if id.ID == userID {
				isMember = false
			}
		}
	}

	return isMember, nil
}
