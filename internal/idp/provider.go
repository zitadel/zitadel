package idp

import (
	"context"

	"golang.org/x/text/language"
)

// Provider is the minimal implementation for a 3rd party authentication provider
type Provider interface {
	Name() string
	BeginAuth(ctx context.Context, state string, params ...any) (Session, error)
	IsLinkingAllowed() bool
	IsCreationAllowed() bool
	IsAutoCreation() bool
	IsAutoUpdate() bool
}

// User contains the information of a federated user
// All the data from the provider can be found in the RawData field
type User struct {
	ID                string
	FirstName         string
	LastName          string
	DisplayName       string
	NickName          string
	PreferredUsername string
	Email             string
	IsEmailVerified   bool
	Phone             string
	IsPhoneVerified   bool
	PreferredLanguage language.Tag
	AvatarURL         string
	Profile           string
	RawData           any
}
