package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	InitCodeMessageType                 = "InitCode"
	PasswordResetMessageType            = "PasswordReset"
	VerifyEmailMessageType              = "VerifyEmail"
	VerifyPhoneMessageType              = "VerifyPhone"
	DomainClaimedMessageType            = "DomainClaimed"
	PasswordlessRegistrationMessageType = "PasswordlessRegistration"
	MessageTitle                        = "Title"
	MessagePreHeader                    = "PreHeader"
	MessageSubject                      = "Subject"
	MessageGreeting                     = "Greeting"
	MessageText                         = "Text"
	MessageButtonText                   = "ButtonText"
	MessageFooterText                   = "Footer"
)

type MessageTexts struct {
	InitCode                 CustomMessageText
	PasswordReset            CustomMessageText
	VerifyEmail              CustomMessageText
	VerifyPhone              CustomMessageText
	DomainClaimed            CustomMessageText
	PasswordlessRegistration CustomMessageText
}

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

func (m *MessageTexts) GetMessageTextByType(msgType string) *CustomMessageText {
	switch msgType {
	case InitCodeMessageType:
		return &m.InitCode
	case PasswordResetMessageType:
		return &m.PasswordReset
	case VerifyEmailMessageType:
		return &m.VerifyEmail
	case VerifyPhoneMessageType:
		return &m.VerifyPhone
	case DomainClaimedMessageType:
		return &m.DomainClaimed
	case PasswordlessRegistrationMessageType:
		return &m.PasswordlessRegistration
	}
	return nil
}
