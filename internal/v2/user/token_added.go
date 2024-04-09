package user

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type tokenAddedPayload struct {
	TokenID           string       `json:"tokenId"`
	ApplicationID     string       `json:"applicationId"`
	UserAgentID       string       `json:"userAgentId"`
	RefreshTokenID    string       `json:"refreshTokenID,omitempty"`
	Audience          []string     `json:"audience"`
	Scopes            []string     `json:"scopes"`
	Expiration        time.Time    `json:"expiration"`
	PreferredLanguage language.Tag `json:"preferredLanguage"`
}

type TokenAddedEvent tokenAddedEvent
type tokenAddedEvent = eventstore.Event[tokenAddedPayload]

func TokenAddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*TokenAddedEvent, error) {
	event, err := eventstore.EventFromStorage[tokenAddedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*TokenAddedEvent)(event), nil
}
