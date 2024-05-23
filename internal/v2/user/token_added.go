package user

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type tokenAddedPayload struct {
	TokenID               string             `json:"tokenId,omitempty"`
	ApplicationID         string             `json:"applicationId,omitempty"`
	UserAgentID           string             `json:"userAgentId,omitempty"`
	RefreshTokenID        string             `json:"refreshTokenID,omitempty"`
	Audience              []string           `json:"audience,omitempty"`
	Scopes                []string           `json:"scopes,omitempty"`
	AuthMethodsReferences []string           `json:"authMethodsReferences,omitempty"`
	AuthTime              time.Time          `json:"authTime,omitempty"`
	Expiration            time.Time          `json:"expiration,omitempty"`
	PreferredLanguage     string             `json:"preferredLanguage,omitempty"`
	Reason                domain.TokenReason `json:"reason,omitempty"`
	Actor                 *domain.TokenActor `json:"actor,omitempty"`
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
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-0YSt4", "Errors.Invalid.Event.Type")
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
