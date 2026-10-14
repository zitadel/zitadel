package idp

import (
	"context"

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

// PromptParam allows to pass the prompt values for the authorization request to BeginAuth.
// If passed, it replaces the provider's default prompt (`select_account`), even if empty.
// It is currently only supported by OAuth and OIDC based providers.
type PromptParam []string

func (p PromptParam) setValue() {}

// AuthURLParam allows to pass an additional query parameter for the authorization request to BeginAuth.
// It takes precedence over parameters set by the provider itself (e.g. the default `prompt`).
// It is currently only supported by OAuth and OIDC based providers.
type AuthURLParam struct {
	Key   string
	Value string
}

func (p AuthURLParam) setValue() {}
