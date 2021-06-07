package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	LoginCustomText = "Login"
)

type CustomLoginText struct {
	models.ObjectRoot

	State           PolicyState
	Default         bool
	MessageTextType string
	Language        language.Tag

	SelectAccountTitle       string
	SelectAccountDescription string
	SelectAccountOtherUser   string

	LoginTitle                   string
	LoginDescription             string
	LoginNameLabel               string
	LoginRegisterButtonText      string
	LoginNextButtonText          string
	LoginExternalUserDescription string

	PasswordTitle          string
	PasswordDescription    string
	PasswordLabel          string
	PasswordBackButtonText string
	PasswordNextButtonText string

	VerifyMFATitle             string
	VerifyMFADescription       string
	VerifyMFAOTPCodeLabel      string
	VerifyMFAOTPNextButtonText string
	// Back button?

	RegistrationOptionTitle                    string
	RegistrationOptionDescription              string
	RegistrationOptionUserNameButtonText       string
	RegistrationOptionExternalLoginDescription string

	RegistrationTitle                string
	RegistrationDescription          string
	RegistrationFirstnameLabel       string
	RegistrationLastnameLabel        string
	RegistrationEmailLabel           string
	RegistrationLanguageLabel        string
	RegistrationGenderLabel          string
	RegistrationPasswordLabel        string
	RegistrationPasswordConfirmLabel string
	RegistrationSaveButtonText       string

	RegisterOrgTitle                string
	RegisterOrgDescription          string
	RegisterOrgOrgNameLabel         string
	RegisterOrgFirstnameLabel       string
	RegisterOrgLastnameLabel        string
	RegisterOrgUsernameLabel        string
	RegisterOrgEmailLabel           string
	RegisterOrgPasswordLabel        string
	RegisterOrgPasswordConfirmLabel string
	RegisterOrgSaveButtonText       string
}

func (m *CustomLoginText) IsValid() bool {
	return m.MessageTextType != "" && m.Language != language.Und
}
