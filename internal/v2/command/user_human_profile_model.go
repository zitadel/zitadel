package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanProfileWriteModel struct {
	eventstore.WriteModel

	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            domain.Gender

	UserState domain.UserState
}

func NewHumanProfileWriteModel(userID, resourceOwner string) *HumanProfileWriteModel {
	return &HumanProfileWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanProfileWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.FirstName = e.FirstName
			wm.LastName = e.LastName
			wm.NickName = e.NickName
			wm.DisplayName = e.DisplayName
			wm.PreferredLanguage = e.PreferredLanguage
			wm.Gender = e.Gender
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.FirstName = e.FirstName
			wm.LastName = e.LastName
			wm.NickName = e.NickName
			wm.DisplayName = e.DisplayName
			wm.PreferredLanguage = e.PreferredLanguage
			wm.Gender = e.Gender
			wm.UserState = domain.UserStateActive
		case *user.HumanProfileChangedEvent:
			if e.FirstName != "" {
				wm.FirstName = e.FirstName
			}
			if e.LastName != "" {
				wm.LastName = e.LastName
			}
			if e.NickName != nil {
				wm.NickName = *e.NickName
			}
			if e.DisplayName != nil {
				wm.DisplayName = *e.DisplayName
			}
			if e.PreferredLanguage != nil {
				wm.PreferredLanguage = *e.PreferredLanguage
			}
			if e.Gender != nil {
				wm.Gender = *e.Gender
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanProfileWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *HumanProfileWriteModel) NewChangedEvent(
	ctx context.Context,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
) (*user.HumanProfileChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanProfileChangedEvent(ctx, UserAggregateFromWriteModel(&wm.WriteModel))
	if wm.FirstName != firstName {
		hasChanged = true
		changedEvent.FirstName = firstName
	}
	if wm.LastName != lastName {
		hasChanged = true
		changedEvent.LastName = lastName
	}
	if wm.NickName != nickName {
		hasChanged = true
		changedEvent.NickName = &nickName
	}
	if wm.DisplayName != displayName {
		hasChanged = true
		changedEvent.DisplayName = &displayName
	}
	if wm.PreferredLanguage != preferredLanguage {
		hasChanged = true
		changedEvent.PreferredLanguage = &preferredLanguage
	}
	if gender.Valid() && wm.Gender != gender {
		hasChanged = true
		changedEvent.Gender = &gender
	}

	return changedEvent, hasChanged
}
