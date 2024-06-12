package user

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type humanRegisteredPayload struct {
	Username          string              `json:"userName"`
	FirstName         string              `json:"firstName,omitempty"`
	LastName          string              `json:"lastName,omitempty"`
	NickName          string              `json:"nickName,omitempty"`
	DisplayName       string              `json:"displayName,omitempty"`
	PreferredLanguage language.Tag        `json:"preferredLanguage,omitempty"`
	Gender            domain.Gender       `json:"gender,omitempty"`
	EmailAddress      domain.EmailAddress `json:"email,omitempty"`
	PhoneNumber       domain.PhoneNumber  `json:"phone,omitempty"`

	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	Secret                 *crypto.CryptoValue `json:"secret,omitempty"` // legacy
	EncodedHash            string              `json:"encodedHash,omitempty"`
	PasswordChangeRequired bool                `json:"changeRequired,omitempty"`
}

type HumanRegisteredEvent eventstore.Event[humanRegisteredPayload]

const HumanRegisteredType = humanPrefix + ".selfregistered"

var _ eventstore.TypeChecker = (*HumanRegisteredEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanRegisteredEvent) ActionType() string {
	return HumanRegisteredType
}

func HumanRegisteredEventFromStorage(event *eventstore.StorageEvent) (e *HumanRegisteredEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-8HvGi", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanRegisteredPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanRegisteredEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
