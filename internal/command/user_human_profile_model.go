package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
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
		ResourceOwner(wm.ResourceOwner).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanProfileChangedType,
			user.UserRemovedType)
}

func (wm *HumanProfileWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
) (*user.HumanProfileChangedEvent, bool, error) {
	changes := make([]user.ProfileChanges, 0)
	var err error

	if wm.FirstName != firstName {
		changes = append(changes, user.ChangeFirstName(firstName))
	}
	if wm.LastName != lastName {
		changes = append(changes, user.ChangeLastName(lastName))
	}
	if wm.NickName != nickName {
		changes = append(changes, user.ChangeNickName(nickName))
	}
	if wm.DisplayName != displayName {
		changes = append(changes, user.ChangeDisplayName(displayName))
	}
	if wm.PreferredLanguage != preferredLanguage {
		changes = append(changes, user.ChangePreferredLanguage(preferredLanguage))
	}
	if wm.Gender != gender {
		changes = append(changes, user.ChangeGender(gender))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := user.NewHumanProfileChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
