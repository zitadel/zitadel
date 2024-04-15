package user

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
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
	Country           string              `json:"country,omitempty"`
	Locality          string              `json:"locality,omitempty"`
	PostalCode        string              `json:"postalCode,omitempty"`
	Region            string              `json:"region,omitempty"`
	StreetAddress     string              `json:"streetAddress,omitempty"`

	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	Secret                 *crypto.CryptoValue `json:"secret,omitempty"` // legacy
	EncodedHash            string              `json:"encodedHash,omitempty"`
	PasswordChangeRequired bool                `json:"changeRequired,omitempty"`
}

type HumanRegisteredEvent humanRegisteredEvent
type humanRegisteredEvent = eventstore.Event[humanRegisteredPayload]

func HumanRegisteredEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanRegisteredEvent, error) {
	event, err := eventstore.EventFromStorage[humanRegisteredEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanRegisteredEvent)(event), nil
}
