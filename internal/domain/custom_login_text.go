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

	LoginKeyLogin                          = "Login."
	LoginKeyLoginTitle                     = LoginKeyLogin + "Title"
	LoginKeyLoginDescription               = LoginKeyLogin + "Description"
	LoginKeyLoginTitleLinkingProcess       = LoginKeyLogin + "TitleLinking"
	LoginKeyLoginDescriptionLinkingProcess = LoginKeyLogin + "DescriptionLinking"
	LoginKeyLoginNameLabel                 = LoginKeyLogin + "LoginNameLabel"
	LoginKeyLoginRegisterButtonText        = LoginKeyLogin + "RegisterButtonText"
	LoginKeyLoginNextButtonText            = LoginKeyLogin + "NextButtonText"
	LoginKeyLoginExternalUserDescription   = LoginKeyLogin + "ExternalUserDescription"
	LoginKeyUserMustBeMemberOfOrg          = LoginKeyLogin + "MustBeMemberOfOrg"

	LoginKeyPassword               = "Password."
	LoginKeyPasswordTitle          = LoginKeyPassword + "Title"
	LoginKeyPasswordDescription    = LoginKeyPassword + "Description"
	LoginKeyPasswordLabel          = LoginKeyPassword + "PasswordLabel"
	LoginKeyPasswordResetLinkText  = LoginKeyPassword + "ResetLinkText"
	LoginKeyPasswordBackButtonText = LoginKeyPassword + "BackButtonText"
	LoginKeyPasswordNextButtonText = LoginKeyPassword + "NextButtonText"

	LoginKeyResetPassword               = "ResetPassword."
	LoginKeyResetPasswordTitle          = LoginKeyResetPassword + "Title"
	LoginKeyResetPasswordDescription    = LoginKeyResetPassword + "Description"
	LoginKeyResetPasswordNextButtonText = LoginKeyResetPassword + "NextButtonText"

	LoginKeyInitializeUser                        = "InitializeUser."
	LoginKeyInitializeUserTitle                   = LoginKeyInitializeUser + "Title"
	LoginKeyInitializeUserDescription             = LoginKeyInitializeUser + "Description"
	LoginKeyInitializeUserCodeLabel               = LoginKeyInitializeUser + "CodeLabel"
	LoginKeyInitializeUserNewPasswordLabel        = LoginKeyInitializeUser + "NewPasswordLabel"
	LoginKeyInitializeUserNewPasswordConfirmLabel = LoginKeyInitializeUser + "NewPasswordConfirm"
	LoginKeyInitializeUserResendButtonText        = LoginKeyInitializeUser + "ResendButtonText"
	LoginKeyInitializeUserNextButtonText          = LoginKeyInitializeUser + "NextButtonText"

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

	LoginKeySuccessLogin                        = "SuccessLogin."
	LoginKeySuccessLoginTitle                   = LoginKeySuccessLogin + "Title"
	LoginKeySuccessLoginAutoRedirectDescription = LoginKeySuccessLogin + "AutoRedirectDescription"
	LoginKeySuccessLoginRedirectedDescription   = LoginKeySuccessLogin + "RedirectedDescription"
	LoginKeySuccessLoginNextButtonText          = LoginKeySuccessLogin + "NextButtonText"
)

type CustomLoginText struct {
	models.ObjectRoot

	State    PolicyState
	Default  bool
	Language language.Tag

	SelectAccountScreenText      SelectAccountScreenText
	LoginScreenText              LoginScreenText
	PasswordScreenText           PasswordScreenText
	ResetPasswordScreenText      ResetPasswordScreenText
	InitializeUserScreenText     InitializeUserScreenText
	InitializeDoneScreenText     InitializeDoneScreenText
	InitMFAPromptScreenText      InitMFAPromptScreenText
	InitMFAOTPScreenText         InitMFAOTPScreenText
	InitMFAU2FScreenText         InitMFAU2FScreenText
	InitMFADoneScreenText        InitMFADoneScreenText
	VerifyMFAOTPScreenText       VerifyMFAOTPScreenText
	VerifyMFAU2FScreenText       VerifyMFAU2FScreenText
	RegistrationOptionScreenText RegistrationOptionScreenText
	RegistrationUserScreenText   RegistrationUserScreenText
	RegistrationOrgScreenText    RegistrationOrgScreenText
	PasswordlessScreenText       PasswordlessScreenText
	SuccessLoginScreenText       SuccessLoginScreenText
}

func (m *CustomLoginText) IsValid() bool {
	return m.Language != language.Und
}

type SelectAccountScreenText struct {
	Title       string
	Description string
	OtherUser   string
}

type LoginScreenText struct {
	Title                     string
	Description               string
	TitleLinkingProcess       string
	DescriptionLinkingProcess string
	LoginNameLabel            string
	RegisterButtonText        string
	NextButtonText            string
	ExternalUserDescription   string
	UserMustBeOrgMemberOfOrg  string
}

type PasswordScreenText struct {
	Title          string
	Description    string
	PasswordLabel  string
	ResetLinkText  string
	BackButtonText string
	NextButtonText string
}

type ResetPasswordScreenText struct {
	Title          string
	Description    string
	NextButtonText string
}

type InitializeUserScreenText struct {
	Title                   string
	Description             string
	CodeLabel               string
	NewPasswordLabel        string
	NewPasswordConfirmLabel string
	ResendButtonText        string
	NextButtonText          string
}

type InitializeDoneScreenText struct {
	Title           string
	Description     string
	AbortButtonText string
	NextButtonText  string
}

type InitMFAPromptScreenText struct {
	Title          string
	Description    string
	OTPOption      string
	U2FOption      string
	SkipButtonText string
	NextButtonText string
}

type InitMFAOTPScreenText struct {
	Title       string
	Description string
	SecretLabel string
	CodeLabel   string
}

type InitMFAU2FScreenText struct {
	Title                   string
	Description             string
	TokenNameLabel          string
	RegisterTokenButtonText string
}

type InitMFADoneScreenText struct {
	Title           string
	Description     string
	AbortButtonText string
	NextButtonText  string
}

type VerifyMFAOTPScreenText struct {
	Title          string
	Description    string
	CodeLabel      string
	NextButtonText string
}

type VerifyMFAU2FScreenText struct {
	Title                  string
	Description            string
	ValidateTokenText      string
	OtherOptionDescription string
	OptionOTPButtonText    string
}

type RegistrationOptionScreenText struct {
	Title                    string
	Description              string
	UserNameButtonText       string
	ExternalLoginDescription string
}

type RegistrationUserScreenText struct {
	Title                string
	Description          string
	FirstnameLabel       string
	LastnameLabel        string
	EmailLabel           string
	LanguageLabel        string
	GenderLabel          string
	PasswordLabel        string
	PasswordConfirmLabel string
	SaveButtonText       string
}

type RegistrationOrgScreenText struct {
	Title                string
	Description          string
	OrgNameLabel         string
	FirstnameLabel       string
	LastnameLabel        string
	UsernameLabel        string
	EmailLabel           string
	PasswordLabel        string
	PasswordConfirmLabel string
	SaveButtonText       string
}

type PasswordlessScreenText struct {
	Title                   string
	Description             string
	LoginWithPwButtonText   string
	ValidateTokenButtonText string
}

type SuccessLoginScreenText struct {
	Title                   string
	AutoRedirectDescription string
	RedirectedDescription   string
	NextButtonText          string
}
