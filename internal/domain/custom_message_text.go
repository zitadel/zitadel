package domain

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (m *CustomMessageText) IsValid(supportedLanguages []language.Tag) error {
	if m.MessageTextType == "" {
		return zerrors.ThrowInvalidArgument(nil, "INSTANCE-kd9fs", "Errors.CustomMessageText.Invalid")
	}
	if err := LanguageIsDefined(m.Language); err != nil {
		return err
	}
	return LanguagesAreSupported(supportedLanguages, m.Language)
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
