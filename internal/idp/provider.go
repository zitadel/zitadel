package idp

import (
	"context"
	"math/rand"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
)

// Provider is the minimal implementation for a 3rd party authentication provider
type Provider interface {
	Name() string
	BeginAuth(ctx context.Context, state string, params ...Parameter) (Session, error)
	IsLinkingAllowed() bool
	IsCreationAllowed() bool
	IsAutoCreation() bool
	IsAutoUpdate() bool
}

// User contains the information of a federated user.
type User interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetDisplayName() string
	GetNickname() string
	GetPreferredUsername() string
	GetEmail() domain.EmailAddress
	IsEmailVerified() bool
	GetPhone() domain.PhoneNumber
	IsPhoneVerified() bool
	GetPreferredLanguage() language.Tag
	GetAvatarURL() string
	GetProfile() string
}

// Parameter allows to pass specific parameter to the BeginAuth function
type Parameter interface {
	setValue()
}

// UserAgentID allows to pass the user agent ID of the auth request to BeginAuth
type UserAgentID string

func (p UserAgentID) setValue() {}

// LoginHintParam allows to pass a login_hint to BeginAuth
type LoginHintParam string

func (p LoginHintParam) setValue() {}

func CodeVerifier() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~"

	nChar := rand.Intn(128-43) + 43

	b := make([]byte, nChar)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
