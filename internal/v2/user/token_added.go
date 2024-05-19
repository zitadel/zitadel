package user

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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

type TokenAddedEvent eventstore.Event[tokenAddedPayload]

const TokenAddedType = AggregateType + ".token.added"

var _ eventstore.TypeChecker = (*TokenAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *TokenAddedEvent) ActionType() string {
	return TokenAddedType
}

func TokenAddedEventFromStorage(event *eventstore.StorageEvent) (e *TokenAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-jeeON", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[tokenAddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &TokenAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
