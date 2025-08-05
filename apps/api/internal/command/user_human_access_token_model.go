package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserAccessTokenWriteModel struct {
	eventstore.WriteModel

	TokenID           string
	ApplicationID     string
	UserAgentID       string
	Audience          []string
	Scopes            []string
	Expiration        time.Time
	PreferredLanguage string
	Reason            domain.TokenReason
	Actor             *domain.TokenActor

	UserState domain.UserState
}

func NewUserAccessTokenWriteModel(userID, resourceOwner, tokenID string) *UserAccessTokenWriteModel {
	return &UserAccessTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		TokenID: tokenID,
	}
}

func (wm *UserAccessTokenWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.UserTokenAddedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserTokenRemovedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanSignedOutEvent:
			if wm.UserAgentID != e.UserAgentID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserLockedEvent,
			*user.UserDeactivatedEvent,
			*user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *UserAccessTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.UserTokenAddedEvent:
			wm.TokenID = e.TokenID
			wm.ApplicationID = e.ApplicationID
			wm.UserAgentID = e.UserAgentID
			wm.Audience = e.Audience
			wm.Scopes = e.Scopes
			wm.Expiration = e.Expiration
			wm.PreferredLanguage = e.PreferredLanguage
			wm.UserState = domain.UserStateActive
			wm.Reason = e.Reason
			wm.Actor = e.Actor
			if e.Expiration.Before(time.Now()) {
				wm.UserState = domain.UserStateDeleted
			}
		case *user.UserTokenRemovedEvent,
			*user.HumanSignedOutEvent,
			*user.UserLockedEvent,
			*user.UserDeactivatedEvent,
			*user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserAccessTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.UserTokenAddedType,
			user.UserTokenRemovedType,
			user.HumanSignedOutType,
			user.UserLockedType,
			user.UserDeactivatedType,
			user.UserRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
