package idp

import (
	"context"

	"golang.org/x/text/language"
)

type Provider interface {
	Name() string
	BeginAuth(ctx context.Context, state string, params ...any) (Session, error)
	IsLinkingAllowed() bool
	IsCreationAllowed() bool
	IsAutoCreation() bool
	IsAutoUpdate() bool
}

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
