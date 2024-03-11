package user

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	_ eventstore.Command = (*TokenAdded)(nil)
)

type TokenAdded struct {
	TokenID           string       `json:"tokenId"`
	ApplicationID     string       `json:"applicationId"`
	UserAgentID       string       `json:"userAgentId"`
	RefreshTokenID    string       `json:"refreshTokenID,omitempty"`
	Audience          []string     `json:"audience"`
	Scopes            []string     `json:"scopes"`
	Expiration        time.Time    `json:"expiration"`
	PreferredLanguage language.Tag `json:"preferredLanguage"`

	creator string
}

// Creator implements [eventstore.Command].
func (t *TokenAdded) Creator() string {
	return t.creator
}

// Payload implements [eventstore.Command].
func (t *TokenAdded) Payload() any {
	return t
}

// Revision implements [eventstore.Command].
func (t *TokenAdded) Revision() uint16 {
	return 1
}

// Type implements [eventstore.Command].
func (t *TokenAdded) Type() string {
	return "user.token.added"
}

// UniqueConstraints implements [eventstore.Command].
func (t *TokenAdded) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
