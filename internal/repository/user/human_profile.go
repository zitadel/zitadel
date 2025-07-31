package user

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	profileEventPrefix      = humanEventPrefix + "profile."
	HumanProfileChangedType = profileEventPrefix + "changed"
)

type HumanProfileChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName         string         `json:"firstName,omitempty"`
	LastName          string         `json:"lastName,omitempty"`
	NickName          *string        `json:"nickName,omitempty"`
	DisplayName       *string        `json:"displayName,omitempty"`
	PreferredLanguage *language.Tag  `json:"preferredLanguage,omitempty"`
	Gender            *domain.Gender `json:"gender,omitempty"`
}

func (e *HumanProfileChangedEvent) Payload() any {
	return e
}

func (e *HumanProfileChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanProfileChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []ProfileChanges,
) (*HumanProfileChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "USER-33n8F", "Errors.NoChangesFound")
	}
	changeEvent := &HumanProfileChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanProfileChangedType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type ProfileChanges func(event *HumanProfileChangedEvent)

func ChangeFirstName(firstName string) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.FirstName = firstName
	}
}

func ChangeLastName(lastName string) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.LastName = lastName
	}
}

func ChangeNickName(nickName string) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.NickName = &nickName
	}
}

func ChangeDisplayName(displayName string) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.DisplayName = &displayName
	}
}

func ChangePreferredLanguage(language language.Tag) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.PreferredLanguage = &language
	}
}

func ChangeGender(gender domain.Gender) func(event *HumanProfileChangedEvent) {
	return func(e *HumanProfileChangedEvent) {
		e.Gender = &gender
	}
}

func HumanProfileChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	profileChanged := &HumanProfileChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(profileChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human profile changed")
	}

	return profileChanged, nil
}
