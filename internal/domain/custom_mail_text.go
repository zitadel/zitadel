package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	InitCodeMessageType      = "InitCode"
	PasswordResetMessageType = "PasswordReset"
	VerifyEmailMessageType   = "VerifyEmail"
	VerifyPhoneMessageType   = "VerifyPhone"
	DomainClaimedMessageType = "DomainClaimed"
	MessageTitle             = "Title"
	MessagePreHeader         = "PreHeader"
	MessageSubject           = "Subject"
	MessageGreeting          = "Greeting"
	MessageText              = "Text"
	MessageButtonText        = "ButtonText"
	MessageFooterText        = "FooterText"
)

type CustomMessageText struct {
	models.ObjectRoot

	State           PolicyState
	Default         bool
	MessageTextType string
	Language        language.Tag
	Title           string
	PreHeader       string
	Subject         string
	Greeting        string
	Text            string
	ButtonText      string
	FooterText      string
}

func (m *CustomMessageText) IsValid() bool {
	return m.MessageTextType != "" && m.Language != language.Und
}
