package instance

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueSecretGeneratorType       = "secret_generator"
	secretGeneratorPrefix           = "secret.generator."
	SecretGeneratorAddedEventType   = instanceEventTypePrefix + secretGeneratorPrefix + "added"
	SecretGeneratorChangedEventType = instanceEventTypePrefix + secretGeneratorPrefix + "changed"
	SecretGeneratorRemovedEventType = instanceEventTypePrefix + secretGeneratorPrefix + "removed"
)

func NewAddSecretGeneratorTypeUniqueConstraint(generatorType domain.SecretGeneratorType) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueSecretGeneratorType,
		string(generatorType),
		"Errors.SecretGenerator.AlreadyExists")
}

func NewRemoveSecretGeneratorTypeUniqueConstraint(generatorType domain.SecretGeneratorType) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueSecretGeneratorType,
		string(generatorType))
}

type SecretGeneratorAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GeneratorType       domain.SecretGeneratorType `json:"generatorType"`
	Length              uint                       `json:"length,omitempty"`
	Expiry              time.Duration              `json:"expiry,omitempty"`
	IncludeLowerLetters bool                       `json:"includeLowerLetters,omitempty"`
	IncludeUpperLetters bool                       `json:"includeUpperLetters,omitempty"`
	IncludeDigits       bool                       `json:"includeDigits,omitempty"`
	IncludeSymbols      bool                       `json:"includeSymbols,omitempty"`
}

func NewSecretGeneratorAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorType domain.SecretGeneratorType,
	length uint,
	expiry time.Duration,
	includeLowerLetters,
	includeUpperLetters,
	includeDigits,
	includeSymbols bool,
) *SecretGeneratorAddedEvent {
	return &SecretGeneratorAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecretGeneratorAddedEventType,
		),
		GeneratorType:       generatorType,
		Length:              length,
		Expiry:              expiry,
		IncludeLowerLetters: includeLowerLetters,
		IncludeUpperLetters: includeUpperLetters,
		IncludeDigits:       includeDigits,
		IncludeSymbols:      includeSymbols,
	}
}

func (e *SecretGeneratorAddedEvent) Payload() any {
	return e
}

func (e *SecretGeneratorAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddSecretGeneratorTypeUniqueConstraint(e.GeneratorType)}
}

func SecretGeneratorAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	secretGeneratorAdded := &SecretGeneratorAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(secretGeneratorAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-en9f4", "unable to unmarshal secret generator added")
	}

	return secretGeneratorAdded, nil
}

type SecretGeneratorChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GeneratorType       domain.SecretGeneratorType `json:"generatorType"`
	Length              *uint                      `json:"length,omitempty"`
	Expiry              *time.Duration             `json:"expiry,omitempty"`
	IncludeLowerLetters *bool                      `json:"includeLowerLetters,omitempty"`
	IncludeUpperLetters *bool                      `json:"includeUpperLetters,omitempty"`
	IncludeDigits       *bool                      `json:"includeDigits,omitempty"`
	IncludeSymbols      *bool                      `json:"includeSymbols,omitempty"`
}

func (e *SecretGeneratorChangedEvent) Payload() any {
	return e
}

func (e *SecretGeneratorChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSecretGeneratorChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorType domain.SecretGeneratorType,
	changes []SecretGeneratorChanges,
) (*SecretGeneratorChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-j2jfw", "Errors.NoChangesFound")
	}
	changeEvent := &SecretGeneratorChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecretGeneratorChangedEventType,
		),
		GeneratorType: generatorType,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SecretGeneratorChanges func(event *SecretGeneratorChangedEvent)

func ChangeSecretGeneratorLength(length uint) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.Length = &length
	}
}

func ChangeSecretGeneratorExpiry(expiry time.Duration) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.Expiry = &expiry
	}
}

func ChangeSecretGeneratorIncludeLowerLetters(includeLowerLetters bool) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.IncludeLowerLetters = &includeLowerLetters
	}
}

func ChangeSecretGeneratorIncludeUpperLetters(includeUpperLetters bool) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.IncludeUpperLetters = &includeUpperLetters
	}
}

func ChangeSecretGeneratorIncludeDigits(includeDigits bool) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.IncludeDigits = &includeDigits
	}
}

func ChangeSecretGeneratorIncludeSymbols(includeSymbols bool) func(event *SecretGeneratorChangedEvent) {
	return func(e *SecretGeneratorChangedEvent) {
		e.IncludeSymbols = &includeSymbols
	}
}

func SecretGeneratorChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SecretGeneratorChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-2m09e", "unable to unmarshal secret generator changed")
	}

	return e, nil
}

type SecretGeneratorRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GeneratorType domain.SecretGeneratorType `json:"generatorType"`
}

func (e *SecretGeneratorRemovedEvent) Payload() any {
	return e
}

func (e *SecretGeneratorRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveSecretGeneratorTypeUniqueConstraint(e.GeneratorType)}
}

func NewSecretGeneratorRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorType domain.SecretGeneratorType,
) *SecretGeneratorRemovedEvent {
	return &SecretGeneratorRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecretGeneratorRemovedEventType,
		),
		GeneratorType: generatorType,
	}
}

func SecretGeneratorRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SecretGeneratorRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-m09ke", "unable to unmarshal secret generator removed")
	}

	return e, nil
}
