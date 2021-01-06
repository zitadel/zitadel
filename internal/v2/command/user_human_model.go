package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"golang.org/x/text/language"
)

type HumanWriteModel struct {
	eventstore.WriteModel

	UserName string

	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            domain.Gender

	Email           string
	IsEmailVerified bool

	Phone           string
	IsPhoneVerified bool

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string

	Secret               *crypto.CryptoValue
	SecretChangeRequired bool

	UserState domain.UserState
}

func NewHumanWriteModel(userID string) *HumanWriteModel {
	return &HumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent,
			*user.HumanRegisteredEvent,
			*user.HumanProfileChangedEvent,
			*user.HumanEmailChangedEvent,
			*user.HumanEmailVerifiedEvent,
			*user.HumanPhoneChangedEvent,
			*user.HumanPhoneVerifiedEvent,
			*user.HumanAddressChangedEvent,
			*user.HumanPasswordChangedEvent,
			*user.UserDeactivatedEvent,
			*user.UserReactivatedEvent,
			*user.UserLockedEvent,
			*user.UserUnlockedEvent,
			*user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

//TODO: Compute State? initial/active
func (wm *HumanWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.reduceHumanAddedEvent(e)
		case *user.HumanRegisteredEvent:
			wm.reduceHumanRegisteredEvent(e)
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.HumanProfileChangedEvent:
			wm.reduceHumanProfileChangedEvent(e)
		case *user.HumanEmailChangedEvent:
			wm.reduceHumanEmailChangedEvent(e)
		case *user.HumanEmailVerifiedEvent:
			wm.reduceHumanEmailVerifiedEvent()
		case *user.HumanPhoneChangedEvent:
			wm.reduceHumanPhoneChangedEvent(e)
		case *user.HumanPhoneVerifiedEvent:
			wm.reduceHumanPhoneVerifiedEvent()
		case *user.HumanPasswordChangedEvent:
			wm.reduceHumanPasswordChangedEvent(e)
		case *user.UserLockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateLocked
			}
		case *user.UserUnlockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserDeactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateInactive
			}
		case *user.UserReactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *HumanWriteModel) reduceHumanAddedEvent(e *user.HumanAddedEvent) {
	wm.UserName = e.UserName
	wm.FirstName = e.FirstName
	wm.LastName = e.LastName
	wm.NickName = e.NickName
	wm.DisplayName = e.DisplayName
	wm.PreferredLanguage = e.PreferredLanguage
	wm.Gender = e.Gender
	wm.Email = e.EmailAddress
	wm.Phone = e.PhoneNumber
	wm.Country = e.Country
	wm.Locality = e.Locality
	wm.PostalCode = e.PostalCode
	wm.Region = e.Region
	wm.StreetAddress = e.StreetAddress
	wm.Secret = e.Secret
	wm.SecretChangeRequired = e.ChangeRequired
	wm.UserState = domain.UserStateInitial
}

func (wm *HumanWriteModel) reduceHumanRegisteredEvent(e *user.HumanRegisteredEvent) {
	wm.UserName = e.UserName
	wm.FirstName = e.FirstName
	wm.LastName = e.LastName
	wm.NickName = e.NickName
	wm.DisplayName = e.DisplayName
	wm.PreferredLanguage = e.PreferredLanguage
	wm.Gender = e.Gender
	wm.Email = e.EmailAddress
	wm.Phone = e.PhoneNumber
	wm.Country = e.Country
	wm.Locality = e.Locality
	wm.PostalCode = e.PostalCode
	wm.Region = e.Region
	wm.StreetAddress = e.StreetAddress
	wm.Secret = e.Secret
	wm.SecretChangeRequired = e.ChangeRequired
	wm.UserState = domain.UserStateInitial
}

func (wm *HumanWriteModel) reduceHumanProfileChangedEvent(e *user.HumanProfileChangedEvent) {
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
}

func (wm *HumanWriteModel) reduceHumanEmailChangedEvent(e *user.HumanEmailChangedEvent) {
	wm.Email = e.EmailAddress
	wm.IsEmailVerified = false
}

func (wm *HumanWriteModel) reduceHumanEmailVerifiedEvent() {
	wm.IsEmailVerified = true
}

func (wm *HumanWriteModel) reduceHumanPhoneChangedEvent(e *user.HumanPhoneChangedEvent) {
	wm.Phone = e.PhoneNumber
	wm.IsPhoneVerified = false
}

func (wm *HumanWriteModel) reduceHumanPhoneVerifiedEvent() {
	wm.IsPhoneVerified = true
}

func (wm *HumanWriteModel) reduceHumanAddressChangedEvent(e *user.HumanAddressChangedEvent) {
	if e.Country != nil {
		wm.Country = *e.Country
	}
	if e.Locality != nil {
		wm.Locality = *e.Locality
	}
	if e.PostalCode != nil {
		wm.PostalCode = *e.PostalCode
	}
	if e.Region != nil {
		wm.Region = *e.Region
	}
	if e.StreetAddress != nil {
		wm.StreetAddress = *e.StreetAddress
	}
}

func (wm *HumanWriteModel) reduceHumanPasswordChangedEvent(e *user.HumanPasswordChangedEvent) {
	wm.Secret = e.Secret
	wm.SecretChangeRequired = e.ChangeRequired
}
