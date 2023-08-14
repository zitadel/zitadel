package domain

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	InitCodeMessageType                 = "InitCode"
	PasswordResetMessageType            = "PasswordReset"
	VerifyEmailMessageType              = "VerifyEmail"
	VerifyPhoneMessageType              = "VerifyPhone"
	VerifySMSOTPMessageType             = "VerifySMSOTP"
	VerifyEmailOTPMessageType           = "VerifyEmailOTP"
	DomainClaimedMessageType            = "DomainClaimed"
	PasswordlessRegistrationMessageType = "PasswordlessRegistration"
	PasswordChangeMessageType           = "PasswordChange"
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
	PasswordChange           CustomMessageText
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

func IsMessageTextType(textType string) bool {
	return textType == InitCodeMessageType ||
		textType == PasswordResetMessageType ||
		textType == VerifyEmailMessageType ||
		textType == VerifyPhoneMessageType ||
		textType == VerifySMSOTPMessageType ||
		textType == VerifyEmailOTPMessageType ||
		textType == DomainClaimedMessageType ||
		textType == PasswordlessRegistrationMessageType ||
		textType == PasswordChangeMessageType
}
