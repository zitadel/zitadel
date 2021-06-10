package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	LoginCustomText = "Login"

	LoginKeySelectAccount            = "SelectAccount."
	LoginKeySelectAccountTitle       = LoginKeySelectAccount + "Title"
	LoginKeySelectAccountDescription = LoginKeySelectAccount + "Description"
	LoginKeySelectAccountOtherUser   = LoginKeySelectAccount + "OtherUser"

	LoginKeyLogin                        = "Login."
	LoginKeyLoginTitle                   = LoginKeyLogin + "Title"
	LoginKeyLoginDescription             = LoginKeyLogin + "Description"
	LoginKeyLoginNameLabel               = LoginKeyLogin + "NameLabel"
	LoginKeyLoginRegisterButtonText      = LoginKeyLogin + "RegisterButtonText"
	LoginKeyLoginNextButtonText          = LoginKeyLogin + "NextButtonText"
	LoginKeyLoginExternalUserDescription = LoginKeyLogin + "ExternalUserDescription"

	LoginKeyPassword               = "Password."
	LoginKeyPasswordTitle          = LoginKeyPassword + "Title"
	LoginKeyPasswordDescription    = LoginKeyPassword + "Description"
	LoginKeyPasswordLabel          = LoginKeyPassword + "Label"
	LoginKeyPasswordResetLinkText  = LoginKeyPassword + "ResetLinkText"
	LoginKeyPasswordBackButtonText = LoginKeyPassword + "BackButtonText"
	LoginKeyPasswordNextButtonText = LoginKeyPassword + "NextButtonText"

	LoginKeyResetPassword               = "ResetPassword."
	LoginKeyResetPasswordTitle          = LoginKeyResetPassword + "Title"
	LoginKeyResetPasswordDescription    = LoginKeyResetPassword + "Description"
	LoginKeyResetPasswordNextButtonText = LoginKeyResetPassword + "NextButtonText"

	LoginKeyInitializeUser                   = "InitializeUser."
	LoginKeyInitializeUserTitle              = LoginKeyInitializeUser + "Title"
	LoginKeyInitializeUserDescription        = LoginKeyInitializeUser + "Description"
	LoginKeyInitializeUserCodeLabel          = LoginKeyInitializeUser + "CodeLabel"
	LoginKeyInitializeUserNewPassword        = LoginKeyInitializeUser + "NewPassword"
	LoginKeyInitializeUserNewPasswordConfirm = LoginKeyInitializeUser + "NewPasswordConfirm"
	LoginKeyInitializeUserResendButtonText   = LoginKeyInitializeUser + "ResendButtonText"
	LoginKeyInitializeUserNextButtonText     = LoginKeyInitializeUser + "NextButtonText"

	LoginKeyInitializeDone                = "InitializeDone."
	LoginKeyInitializeDoneTitle           = LoginKeyInitializeDone + "Title"
	LoginKeyInitializeDoneDescription     = LoginKeyInitializeDone + "Description"
	LoginKeyInitializeDoneAbortButtonText = LoginKeyInitializeDone + "AbortButtonText"
	LoginKeyInitializeDoneNextButtonText  = LoginKeyInitializeDone + "NextButtonText"

	LoginKeyInitMFAPrompt               = "InitMFAPrompt."
	LoginKeyInitMFAPromptTitle          = LoginKeyInitMFAPrompt + "Title"
	LoginKeyInitMFAPromptDescription    = LoginKeyInitMFAPrompt + "Description"
	LoginKeyInitMFAPromptOTPOption      = LoginKeyInitMFAPrompt + "OTPOption"
	LoginKeyInitMFAPromptU2FOption      = LoginKeyInitMFAPrompt + "U2FOption"
	LoginKeyInitMFAPromptSkipButtonText = LoginKeyInitMFAPrompt + "SkipButtonText"
	LoginKeyInitMFAPromptNextButtonText = LoginKeyInitMFAPrompt + "NextButtonText"

	LoginKeyInitMFAOTP            = "InitMFAOTP."
	LoginKeyInitMFAOTPTitle       = LoginKeyInitMFAOTP + "Title"
	LoginKeyInitMFAOTPDescription = LoginKeyInitMFAOTP + "Description"
	LoginKeyInitMFAOTPSecretLabel = LoginKeyInitMFAOTP + "SecretLabel"
	LoginKeyInitMFAOTPCodeLabel   = LoginKeyInitMFAOTP + "CodeLabel"

	LoginKeyInitMFAU2F                        = "InitMFAU2F."
	LoginKeyInitMFAU2FTitle                   = LoginKeyInitMFAU2F + "Title"
	LoginKeyInitMFAU2FDescription             = LoginKeyInitMFAU2F + "Description"
	LoginKeyInitMFAU2FTokenNameLabel          = LoginKeyInitMFAU2F + "TokenNameLabel"
	LoginKeyInitMFAU2FRegisterTokenButtonText = LoginKeyInitMFAU2F + "RegisterTokenButtonText"

	LoginKeyInitMFADone                = "InitMFADone."
	LoginKeyInitMFADoneTitle           = LoginKeyInitMFADone + "Title"
	LoginKeyInitMFADoneDescription     = LoginKeyInitMFADone + "Description"
	LoginKeyInitMFADoneAbortButtonText = LoginKeyInitMFADone + "AbortButtonText"
	LoginKeyInitMFADoneNextButtonText  = LoginKeyInitMFADone + "NextButtonText"

	LoginKeyVerifyMFAOTP               = "VerifyMFAOTP."
	LoginKeyVerifyMFAOTPTitle          = LoginKeyVerifyMFAOTP + "Title"
	LoginKeyVerifyMFAOTPDescription    = LoginKeyVerifyMFAOTP + "Description"
	LoginKeyVerifyMFAOTPCodeLabel      = LoginKeyVerifyMFAOTP + "CodeLabel"
	LoginKeyVerifyMFAOTPNextButtonText = LoginKeyVerifyMFAOTP + "NextButtonText"

	LoginKeyVerifyMFAU2F                       = "VerifyMFAU2F."
	LoginKeyVerifyMFAU2FTitle                  = LoginKeyVerifyMFAU2F + "Title"
	LoginKeyVerifyMFAU2FDescription            = LoginKeyVerifyMFAU2F + "Description"
	LoginKeyVerifyMFAU2FValidateTokenText      = LoginKeyVerifyMFAU2F + "ValidateTokenText"
	LoginKeyVerifyMFAU2FOtherOptionDescription = LoginKeyVerifyMFAU2F + "OtherOptionDescription"
	LoginKeyVerifyMFAU2FOptionOTPButtonText    = LoginKeyVerifyMFAU2F + "OptionOTPButtonText"

	LoginKeyRegistrationOption                         = "RegistrationOption."
	LoginKeyRegistrationOptionTitle                    = LoginKeyRegistrationOption + "Title"
	LoginKeyRegistrationOptionDescription              = LoginKeyRegistrationOption + "Description"
	LoginKeyRegistrationOptionUserNameButtonText       = LoginKeyRegistrationOption + "UserNameButtonText"
	LoginKeyRegistrationOptionExternalLoginDescription = LoginKeyRegistrationOption + "ExternalLoginDescription"

	LoginKeyRegistrationUser                     = "RegistrationUser."
	LoginKeyRegistrationUserTitle                = LoginKeyRegistrationUser + "Title"
	LoginKeyRegistrationUserDescription          = LoginKeyRegistrationUser + "Description"
	LoginKeyRegistrationUserFirstnameLabel       = LoginKeyRegistrationUser + "FirstnameLabel"
	LoginKeyRegistrationUserLastnameLabel        = LoginKeyRegistrationUser + "LastnameLabel"
	LoginKeyRegistrationUserEmailLabel           = LoginKeyRegistrationUser + "EmailLabel"
	LoginKeyRegistrationUserLanguageLabel        = LoginKeyRegistrationUser + "LanguageLabel"
	LoginKeyRegistrationUserGenderLabel          = LoginKeyRegistrationUser + "GenderLabel"
	LoginKeyRegistrationUserPasswordLabel        = LoginKeyRegistrationUser + "PasswordLabel"
	LoginKeyRegistrationUserPasswordConfirmLabel = LoginKeyRegistrationUser + "PasswordConfirmLabel"
	LoginKeyRegistrationUserSaveButtonText       = LoginKeyRegistrationUser + "SaveButtonText"

	LoginKeyRegistrationOrg                 = "RegistrationOrg."
	LoginKeyRegisterOrgTitle                = LoginKeyRegistrationOrg + "Title"
	LoginKeyRegisterOrgDescription          = LoginKeyRegistrationOrg + "Description"
	LoginKeyRegisterOrgOrgNameLabel         = LoginKeyRegistrationOrg + "OrgNameLabel"
	LoginKeyRegisterOrgFirstnameLabel       = LoginKeyRegistrationOrg + "FirstnameLabel"
	LoginKeyRegisterOrgLastnameLabel        = LoginKeyRegistrationOrg + "LastnameLabel"
	LoginKeyRegisterOrgUsernameLabel        = LoginKeyRegistrationOrg + "UsernameLabel"
	LoginKeyRegisterOrgEmailLabel           = LoginKeyRegistrationOrg + "EmailLabel"
	LoginKeyRegisterOrgPasswordLabel        = LoginKeyRegistrationOrg + "PasswordLabel"
	LoginKeyRegisterOrgPasswordConfirmLabel = LoginKeyRegistrationOrg + "PasswordConfirmLabel"
	LoginKeyRegisterOrgSaveButtonText       = LoginKeyRegistrationOrg + "SaveButtonText"

	LoginKeyPasswordless                        = "Passwordless."
	LoginKeyPasswordlessTitle                   = LoginKeyPasswordless + "Title"
	LoginKeyPasswordlessDescription             = LoginKeyPasswordless + "Description"
	LoginKeyPasswordlessLoginWithPwButtonText   = LoginKeyPasswordless + "LoginWithPwButtonText"
	LoginKeyPasswordlessValidateTokenButtonText = LoginKeyPasswordless + "ValidateTokenButtonText"
)

type CustomLoginText struct {
	models.ObjectRoot

	State    PolicyState
	Default  bool
	Language language.Tag

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

	ResetPasswordTitle          string
	ResetPasswordDescription    string
	ResetPasswordNextButtonText string

	InitializeUserTitle                   string
	InitializeUserDescription             string
	InitializeUserCodeLabel               string
	InitializeUserNewPasswordLabel        string
	InitializeUserNewPasswordConfirmLabel string
	InitializeUserResendButtonText        string
	InitializeUserNextButtonText          string

	InitializeDoneTitle           string
	InitializeDoneDescription     string
	InitializeDoneAbortButtonText string
	InitializeDoneNextButtonText  string

	InitMFAPromptTitle          string
	InitMFAPromptDescription    string
	InitMFAPromptOTPOption      string
	InitMFAPromptU2FOption      string
	InitMFAPromptSkipButtonText string
	InitMFAPromptNextButtonText string

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

	RegistrationUserTitle                string
	RegistrationUserDescription          string
	RegistrationUserFirstnameLabel       string
	RegistrationUserLastnameLabel        string
	RegistrationUserEmailLabel           string
	RegistrationUserLanguageLabel        string
	RegistrationUserGenderLabel          string
	RegistrationUserPasswordLabel        string
	RegistrationUserPasswordConfirmLabel string
	RegistrationUserSaveButtonText       string

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

	SuccessLoginTitle                   string
	SuccessLoginAutoRedirectDescription string
	SuccessLoginRedirectDescription     string
	SuccessLoginNextButtonText          string
}

func (m *CustomLoginText) IsValid() bool {
	return m.Language != language.Und
}
