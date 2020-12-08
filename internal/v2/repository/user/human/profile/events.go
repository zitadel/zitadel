package profile

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/user/human"
	"golang.org/x/text/language"
)

const (
	profileEventPrefix      = eventstore.EventType("user.human.profile.")
	HumanProfileChangedType = profileEventPrefix + "changed"
)

type HumanProfileChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName         string       `json:"firstName,omitempty"`
	LastName          string       `json:"lastName,omitempty"`
	NickName          string       `json:"nickName,omitempty"`
	DisplayName       string       `json:"displayName,omitempty"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty"`
	Gender            human.Gender `json:"gender,omitempty"`
}

func (e *HumanProfileChangedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanProfileChangedEvent) Data() interface{} {
	return e
}

func NewHumanProfileChangedEvent(
	ctx context.Context,
	current *HumanProfileWriteModel,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender human.Gender,
) *HumanProfileChangedEvent {
	e := &HumanProfileChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanProfileChangedType,
		),
	}
	if current.FirstName != firstName {
		e.FirstName = firstName
	}
	if current.LastName != lastName {
		e.LastName = lastName
	}
	if current.NickName != nickName {
		e.NickName = nickName
	}
	if current.DisplayName != displayName {
		e.DisplayName = displayName
	}
	if current.PreferredLanguage != preferredLanguage {
		e.PreferredLanguage = preferredLanguage
	}
	if current.Gender != gender {
		e.Gender = gender
	}
	return e
}

func HumanProfileChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	profileChanged := &HumanProfileChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, profileChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human profile changed")
	}

	return profileChanged, nil
}
