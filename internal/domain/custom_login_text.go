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
	PasswordResetLinkText  string
	PasswordBackButtonText string
	PasswordNextButtonText string

	PasswordResetTitle          string
	PasswordResetDescription    string
	PasswordResetNextButtonText string

	InitializeTitle              string
	InitializeDescription        string
	InitializeCodeLabel          string
	InitializeNewPassword        string
	InitializeNewPasswordConfirm string
	InitializeResendButtonText   string
	InitializeNextButtonText     string

	InitializeDoneTitle           string
	InitializeDoneDescription     string
	InitializeDoneAbortButtonText string
	InitializeDoneNextButtonText  string

	InitMFATitle          string
	InitMFADescription    string
	InitMFAOTPOption      string
	InitMFAU2FOption      string
	InitMFASkipButtonText string
	InitMFANextButtonText string

	InitMFAOTPTitle       string
	InitMFAOTPDescription string
	InitMFAOTPSecretLabel string
	InitMFAOTPCodeLabel   string

	InitMFAU2FTitle                   string
	InitMFAU2FDescription             string
	InitMFAU2FTokenNameLabel          string
	InitMFAU2FRegisterTokenButtonText string

	InitMFADoneTitle           string
	InitMFADoneDescription     string
	InitMFADoneAbortButtonText string
	InitMFADoneNextButtonText  string

	VerifyMFAOTPTitle          string
	VerifyMFAOTPDescription    string
	VerifyMFAOTPCodeLabel      string
	VerifyMFAOTPNextButtonText string

	VerifyMFAU2FTitle                  string
	VerifyMFAU2FDescription            string
	VerifyMFAU2FValidateTokenText      string
	VerifyMFAU2FOtherOptionDescription string
	VerifyMFAU2FOptionOTPButtonText    string

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

	PasswordlessTitle                   string
	PasswordlessDescription             string
	PasswordlessLoginWithPwButtonText   string
	PasswordlessValidateTokenButtonText string
}

func (m *CustomLoginText) IsValid() bool {
	return m.MessageTextType != "" && m.Language != language.Und
}
