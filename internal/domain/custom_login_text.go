package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	LoginCustomText = "Login"

	LoginKeyLogin                          = "Login."
	LoginKeyLoginTitle                     = LoginKeyLogin + "Title"
	LoginKeyLoginDescription               = LoginKeyLogin + "Description"
	LoginKeyLoginTitleLinkingProcess       = LoginKeyLogin + "TitleLinking"
	LoginKeyLoginDescriptionLinkingProcess = LoginKeyLogin + "DescriptionLinking"
	LoginKeyLoginNameLabel                 = LoginKeyLogin + "LoginNameLabel"
	LoginKeyLoginUsernamePlaceHolder       = LoginKeyLogin + "UsernamePlaceHolder"
	LoginKeyLoginLoginnamePlaceHolder      = LoginKeyLogin + "LoginnamePlaceHolder"
	LoginKeyLoginRegisterButtonText        = LoginKeyLogin + "RegisterButtonText"
	LoginKeyLoginNextButtonText            = LoginKeyLogin + "NextButtonText"
	LoginKeyLoginExternalUserDescription   = LoginKeyLogin + "ExternalUserDescription"
	LoginKeyLoginUserMustBeMemberOfOrg     = LoginKeyLogin + "MustBeMemberOfOrg"

	LoginKeySelectAccount                          = "SelectAccount."
	LoginKeySelectAccountTitle                     = LoginKeySelectAccount + "Title"
	LoginKeySelectAccountDescription               = LoginKeySelectAccount + "Description"
	LoginKeySelectAccountTitleLinkingProcess       = LoginKeySelectAccount + "TitleLinking"
	LoginKeySelectAccountDescriptionLinkingProcess = LoginKeySelectAccount + "DescriptionLinking"
	LoginKeySelectAccountOtherUser                 = LoginKeySelectAccount + "OtherUser"
	LoginKeySelectAccountSessionStateActive        = LoginKeySelectAccount + "SessionState0"
	LoginKeySelectAccountSessionStateInactive      = LoginKeySelectAccount + "SessionState1"
	LoginKeySelectAccountUserMustBeMemberOfOrg     = LoginKeySelectAccount + "MustBeMemberOfOrg"

	LoginKeyPassword               = "Password."
	LoginKeyPasswordTitle          = LoginKeyPassword + "Title"
	LoginKeyPasswordDescription    = LoginKeyPassword + "Description"
	LoginKeyPasswordLabel          = LoginKeyPassword + "PasswordLabel"
	LoginKeyPasswordMinLength      = LoginKeyPassword + "MinLength"
	LoginKeyPasswordHasUppercase   = LoginKeyPassword + "HasUppercase"
	LoginKeyPasswordHasLowercase   = LoginKeyPassword + "HasLowercase"
	LoginKeyPasswordHasNumber      = LoginKeyPassword + "HasNumber"
	LoginKeyPasswordHasSymbol      = LoginKeyPassword + "HasSymbol"
	LoginKeyPasswordResetLinkText  = LoginKeyPassword + "ResetLinkText"
	LoginKeyPasswordBackButtonText = LoginKeyPassword + "BackButtonText"
	LoginKeyPasswordNextButtonText = LoginKeyPassword + "NextButtonText"

	LoginKeyUsernameChange         = "UsernameChange."
	UsernameChangeTitle            = LoginKeyUsernameChange + "Title"
	UsernameChangeDescription      = LoginKeyUsernameChange + "Description"
	UsernameChangeUsernameLabel    = LoginKeyUsernameChange + "UsernameLabel"
	UsernameChangeCancelButtonText = LoginKeyUsernameChange + "CancelButtonText"
	UsernameChangeNextButtonText   = LoginKeyUsernameChange + "NextButtonText"

	LoginKeyUsernameChangeDone       = "UsernameChangeDone."
	UsernameChangeDoneTitle          = LoginKeyUsernameChangeDone + "Title"
	UsernameChangeDoneDescription    = LoginKeyUsernameChangeDone + "Description"
	UsernameChangeDoneNextButtonText = LoginKeyUsernameChangeDone + "NextButtonText"

	LoginKeyInitPassword                 = "InitPassword."
	LoginKeyInitPasswordTitle            = LoginKeyInitPassword + "Title"
	LoginKeyInitPasswordDescription      = LoginKeyInitPassword + "Description"
	LoginKeyInitPasswordCodeLabel        = LoginKeyInitPassword + "CodeLabel"
	LoginKeyInitPasswordNextButtonText   = LoginKeyInitPassword + "NextButtonText"
	LoginKeyInitPasswordResendButtonText = LoginKeyInitPassword + "ResendButtonText"

	LoginKeyInitPasswordDone                 = "InitPasswordDone."
	LoginKeyInitPasswordDoneTitle            = LoginKeyInitPasswordDone + "Title"
	LoginKeyInitPasswordDoneDescription      = LoginKeyInitPasswordDone + "Description"
	LoginKeyInitPasswordDoneNextButtonText   = LoginKeyInitPasswordDone + "NextButtonText"
	LoginKeyInitPasswordDoneCancelButtonText = LoginKeyInitPasswordDone + "CancelButtonText"

	LoginKeyEmailVerification                 = "EmailVerification."
	LoginKeyEmailVerificationTitle            = LoginKeyEmailVerification + "Title"
	LoginKeyEmailVerificationDescription      = LoginKeyEmailVerification + "Description"
	LoginKeyEmailVerificationCodeLabel        = LoginKeyEmailVerification + "CodeLabel"
	LoginKeyEmailVerificationNextButtonText   = LoginKeyEmailVerification + "NextButtonText"
	LoginKeyEmailVerificationResendButtonText = LoginKeyEmailVerification + "ResendButtonText"

	LoginKeyEmailVerificationDone                 = "EmailVerificationDone."
	LoginKeyEmailVerificationDoneTitle            = LoginKeyEmailVerificationDone + "Title"
	LoginKeyEmailVerificationDoneDescription      = LoginKeyEmailVerificationDone + "Description"
	LoginKeyEmailVerificationDoneNextButtonText   = LoginKeyEmailVerificationDone + "NextButtonText"
	LoginKeyEmailVerificationDoneCancelButtonText = LoginKeyEmailVerificationDone + "CancelButtonText"
	LoginKeyEmailVerificationDoneLoginButtonText  = LoginKeyEmailVerificationDone + "LoginButtonText"

	LoginKeyInitializeUser                        = "InitUser."
	LoginKeyInitializeUserTitle                   = LoginKeyInitializeUser + "Title"
	LoginKeyInitializeUserDescription             = LoginKeyInitializeUser + "Description"
	LoginKeyInitializeUserCodeLabel               = LoginKeyInitializeUser + "CodeLabel"
	LoginKeyInitializeUserNewPasswordLabel        = LoginKeyInitializeUser + "NewPasswordLabel"
	LoginKeyInitializeUserNewPasswordConfirmLabel = LoginKeyInitializeUser + "NewPasswordConfirm"
	LoginKeyInitializeUserResendButtonText        = LoginKeyInitializeUser + "ResendButtonText"
	LoginKeyInitializeUserNextButtonText          = LoginKeyInitializeUser + "NextButtonText"

	LoginKeyInitUserDone                = "InitUserDone."
	LoginKeyInitUserDoneTitle           = LoginKeyInitUserDone + "Title"
	LoginKeyInitUserDoneDescription     = LoginKeyInitUserDone + "Description"
	LoginKeyInitUserDoneAbortButtonText = LoginKeyInitUserDone + "AbortButtonText"
	LoginKeyInitUserDoneNextButtonText  = LoginKeyInitUserDone + "NextButtonText"

	LoginKeyInitMFAPrompt               = "InitMFAPrompt."
	LoginKeyInitMFAPromptTitle          = LoginKeyInitMFAPrompt + "Title"
	LoginKeyInitMFAPromptDescription    = LoginKeyInitMFAPrompt + "Description"
	LoginKeyInitMFAPromptOTPOption      = LoginKeyInitMFAPrompt + "Provider0"
	LoginKeyInitMFAPromptU2FOption      = LoginKeyInitMFAPrompt + "Provider1"
	LoginKeyInitMFAPromptSkipButtonText = LoginKeyInitMFAPrompt + "SkipButtonText"
	LoginKeyInitMFAPromptNextButtonText = LoginKeyInitMFAPrompt + "NextButtonText"

	LoginKeyInitMFAOTP               = "InitMFAOTP."
	LoginKeyInitMFAOTPTitle          = LoginKeyInitMFAOTP + "Title"
	LoginKeyInitMFAOTPDescription    = LoginKeyInitMFAOTP + "Description"
	LoginKeyInitMFAOTPDescriptionOTP = LoginKeyInitMFAOTP + "OTPDescription"
	LoginKeyInitMFAOTPSecretLabel    = LoginKeyInitMFAOTP + "SecretLabel"
	LoginKeyInitMFAOTPCodeLabel      = LoginKeyInitMFAOTP + "CodeLabel"
	LoginKeyInitMFANextButtonText    = LoginKeyInitMFAOTP + "NextButtonText"
	LoginKeyInitMFACancelButtonText  = LoginKeyInitMFAOTP + "CancelButtonText"

	LoginKeyInitMFAU2F                        = "InitMFAU2F."
	LoginKeyInitMFAU2FTitle                   = LoginKeyInitMFAU2F + "Title"
	LoginKeyInitMFAU2FDescription             = LoginKeyInitMFAU2F + "Description"
	LoginKeyInitMFAU2FTokenNameLabel          = LoginKeyInitMFAU2F + "TokenNameLabel"
	LoginKeyInitMFAU2FNotSupported            = LoginKeyInitMFAU2F + "NotSupported"
	LoginKeyInitMFAU2FRegisterTokenButtonText = LoginKeyInitMFAU2F + "RegisterTokenButtonText"
	LoginKeyInitMFAU2FErrorRetry              = LoginKeyInitMFAU2F + "Error.Retry"

	LoginKeyInitMFADone                = "InitMFADone."
	LoginKeyInitMFADoneTitle           = LoginKeyInitMFADone + "Title"
	LoginKeyInitMFADoneDescription     = LoginKeyInitMFADone + "Description"
	LoginKeyInitMFADoneAbortButtonText = LoginKeyInitMFADone + "AbortButtonText"
	LoginKeyInitMFADoneNextButtonText  = LoginKeyInitMFADone + "NextButtonText"

	LoginKeyMFAProviders            = "MFAProvider."
	LoginKeyMFAProvidersChooseOther = LoginKeyMFAProviders + "ChooseOther"
	LoginKeyMFAProvidersOTP         = LoginKeyMFAProviders + "Provider0"
	LoginKeyMFAProvidersU2F         = LoginKeyMFAProviders + "Provider1"

	LoginKeyVerifyMFAOTP               = "VerifyMFAOTP."
	LoginKeyVerifyMFAOTPTitle          = LoginKeyVerifyMFAOTP + "Title"
	LoginKeyVerifyMFAOTPDescription    = LoginKeyVerifyMFAOTP + "Description"
	LoginKeyVerifyMFAOTPCodeLabel      = LoginKeyVerifyMFAOTP + "CodeLabel"
	LoginKeyVerifyMFAOTPNextButtonText = LoginKeyVerifyMFAOTP + "NextButtonText"

	LoginKeyVerifyMFAU2F                  = "VerifyMFAU2F."
	LoginKeyVerifyMFAU2FTitle             = LoginKeyVerifyMFAU2F + "Title"
	LoginKeyVerifyMFAU2FDescription       = LoginKeyVerifyMFAU2F + "Description"
	LoginKeyVerifyMFAU2FNotSupported      = LoginKeyVerifyMFAU2F + "NotSupported"
	LoginKeyVerifyMFAU2FValidateTokenText = LoginKeyVerifyMFAU2F + "ValidateTokenButtonText"
	LoginKeyVerifyMFAU2FErrorRetry        = LoginKeyVerifyMFAU2F + "Error.Retry"

	LoginKeyPasswordless                        = "Passwordless."
	LoginKeyPasswordlessTitle                   = LoginKeyPasswordless + "Title"
	LoginKeyPasswordlessDescription             = LoginKeyPasswordless + "Description"
	LoginKeyPasswordlessLoginWithPwButtonText   = LoginKeyPasswordless + "LoginWithPwButtonText"
	LoginKeyPasswordlessValidateTokenButtonText = LoginKeyPasswordless + "ValidateTokenButtonText"
	LoginKeyPasswordlessNotSupported            = LoginKeyPasswordless + "NotSupported"
	LoginKeyPasswordlessErrorRetry              = LoginKeyPasswordless + "Error.Retry"

	LoginKeyPasswordChange                        = "PasswordChange."
	LoginKeyPasswordChangeTitle                   = LoginKeyPasswordChange + "Title"
	LoginKeyPasswordChangeDescription             = LoginKeyPasswordChange + "Description"
	LoginKeyPasswordChangeOldPasswordLabel        = LoginKeyPasswordChange + "OldPasswordLabel"
	LoginKeyPasswordChangeNewPasswordLabel        = LoginKeyPasswordChange + "NewPasswordLabel"
	LoginKeyPasswordChangeNewPasswordConfirmLabel = LoginKeyPasswordChange + "NewPasswordConfirmationLabel"
	LoginKeyPasswordChangeCancelButtonText        = LoginKeyPasswordChange + "CancelButtonText"
	LoginKeyPasswordChangeNextButtonText          = LoginKeyPasswordChange + "NextButtonText"

	LoginKeyPasswordChangeDone               = "PasswordChangeDone."
	LoginKeyPasswordChangeDoneTitle          = LoginKeyPasswordChangeDone + "Title"
	LoginKeyPasswordChangeDoneDescription    = LoginKeyPasswordChangeDone + "Description"
	LoginKeyPasswordChangeDoneNextButtonText = LoginKeyPasswordChangeDone + "NextButtonText"

	LoginKeyPasswordResetDone               = "PasswordResetDone."
	LoginKeyPasswordResetDoneTitle          = LoginKeyPasswordResetDone + "Title"
	LoginKeyPasswordResetDoneDescription    = LoginKeyPasswordResetDone + "Description"
	LoginKeyPasswordResetDoneNextButtonText = LoginKeyPasswordResetDone + "NextButtonText"

	LoginKeyRegistrationOption                         = "RegisterOption."
	LoginKeyRegistrationOptionTitle                    = LoginKeyRegistrationOption + "Title"
	LoginKeyRegistrationOptionDescription              = LoginKeyRegistrationOption + "Description"
	LoginKeyRegistrationOptionUserNameButtonText       = LoginKeyRegistrationOption + "RegisterUsernamePasswordButtonText"
	LoginKeyRegistrationOptionExternalLoginDescription = LoginKeyRegistrationOption + "ExternalLoginDescription"

	LoginKeyRegistrationUser                         = "RegistrationUser."
	LoginKeyRegistrationUserTitle                    = LoginKeyRegistrationUser + "Title"
	LoginKeyRegistrationUserDescription              = LoginKeyRegistrationUser + "Description"
	LoginKeyRegistrationUserFirstnameLabel           = LoginKeyRegistrationUser + "FirstnameLabel"
	LoginKeyRegistrationUserLastnameLabel            = LoginKeyRegistrationUser + "LastnameLabel"
	LoginKeyRegistrationUserEmailLabel               = LoginKeyRegistrationUser + "EmailLabel"
	LoginKeyRegistrationUserUsernameLabel            = LoginKeyRegistrationUser + "UsernameLabel"
	LoginKeyRegistrationUserLanguageLabel            = LoginKeyRegistrationUser + "LanguageLabel"
	LoginKeyRegistrationUserPasswordLabel            = LoginKeyRegistrationUser + "PasswordLabel"
	LoginKeyRegistrationUserPasswordConfirmLabel     = LoginKeyRegistrationUser + "PasswordConfirmLabel"
	LoginKeyRegistrationUserTOSAndPrivacy            = LoginKeyRegistrationUser + "TosAndPrivacyLabel"
	LoginKeyRegistrationUserTOSConfirm               = LoginKeyRegistrationUser + "TosConfirm"
	LoginKeyRegistrationUserTOSLink                  = LoginKeyRegistrationUser + "TosLink"
	LoginKeyRegistrationUserTOSLinkText              = LoginKeyRegistrationUser + "TosLinkText"
	LoginKeyRegistrationUserPrivacyConfirm           = LoginKeyRegistrationUser + "TosConfirmAnd"
	LoginKeyRegistrationUserPrivacyLink              = LoginKeyRegistrationUser + "PrivacyLink"
	LoginKeyRegistrationUserPrivacyLinkText          = LoginKeyRegistrationUser + "PrivacyLinkText"
	LoginKeyRegistrationUserExternalLoginDescription = LoginKeyRegistrationUser + "ExternalLogin"
	LoginKeyRegistrationUserNextButtonText           = LoginKeyRegistrationUser + "NextButtonText"
	LoginKeyRegistrationUserBackButtonText           = LoginKeyRegistrationUser + "BackButtonText"

	LoginKeyRegistrationOrg                     = "RegistrationOrg."
	LoginKeyRegisterOrgTitle                    = LoginKeyRegistrationOrg + "Title"
	LoginKeyRegisterOrgDescription              = LoginKeyRegistrationOrg + "Description"
	LoginKeyRegisterOrgOrgNameLabel             = LoginKeyRegistrationOrg + "OrgNameLabel"
	LoginKeyRegisterOrgFirstnameLabel           = LoginKeyRegistrationOrg + "FirstnameLabel"
	LoginKeyRegisterOrgLastnameLabel            = LoginKeyRegistrationOrg + "LastnameLabel"
	LoginKeyRegisterOrgUsernameLabel            = LoginKeyRegistrationOrg + "UsernameLabel"
	LoginKeyRegisterOrgEmailLabel               = LoginKeyRegistrationOrg + "EmailLabel"
	LoginKeyRegisterOrgPasswordLabel            = LoginKeyRegistrationOrg + "PasswordLabel"
	LoginKeyRegisterOrgPasswordConfirmLabel     = LoginKeyRegistrationOrg + "PasswordConfirmLabel"
	LoginKeyRegisterOrgTOSAndPrivacy            = LoginKeyRegistrationOrg + "TosAndPrivacyLabel"
	LoginKeyRegisterOrgTOSConfirm               = LoginKeyRegistrationOrg + "TosConfirm"
	LoginKeyRegisterOrgTOSLink                  = LoginKeyRegistrationOrg + "TosLink"
	LoginKeyRegisterOrgTOSLinkText              = LoginKeyRegistrationOrg + "TosLinkText"
	LoginKeyRegisterOrgPrivacyConfirm           = LoginKeyRegistrationOrg + "TosConfirmAnd"
	LoginKeyRegisterOrgPrivacyLink              = LoginKeyRegistrationOrg + "PrivacyLink"
	LoginKeyRegisterOrgPrivacyLinkText          = LoginKeyRegistrationOrg + "PrivacyLinkText"
	LoginKeyRegisterOrgExternalLoginDescription = LoginKeyRegistrationOrg + "ExternalLogin"
	LoginKeyRegisterOrgSaveButtonText           = LoginKeyRegistrationOrg + "SaveButtonText"

	LoginKeyLinkingUserDone                 = "LinkingUsersDone."
	LoginKeyLinkingUserDoneTitle            = LoginKeyLinkingUserDone + "Title"
	LoginKeyLinkingUserDoneDescription      = LoginKeyLinkingUserDone + "Description"
	LoginKeyLinkingUserDoneCancelButtonText = LoginKeyLinkingUserDone + "CancelButtonText"
	LoginKeyLinkingUserDoneNextButtonText   = LoginKeyLinkingUserDone + "NextButtonText"

	LoginKeyExternalNotFound                       = "ExternalNotFound."
	LoginKeyExternalNotFoundTitle                  = LoginKeyExternalNotFound + "Title"
	LoginKeyExternalNotFoundDescription            = LoginKeyExternalNotFound + "Description"
	LoginKeyExternalNotFoundLinkButtonText         = LoginKeyExternalNotFound + "LinkButtonText"
	LoginKeyExternalNotFoundAutoRegisterButtonText = LoginKeyExternalNotFound + "AutoRegisterButtonText"

	LoginKeySuccessLogin                        = "LoginSuccess."
	LoginKeySuccessLoginTitle                   = LoginKeySuccessLogin + "Title"
	LoginKeySuccessLoginAutoRedirectDescription = LoginKeySuccessLogin + "AutoRedirectDescription"
	LoginKeySuccessLoginRedirectedDescription   = LoginKeySuccessLogin + "RedirectedDescription"
	LoginKeySuccessLoginNextButtonText          = LoginKeySuccessLogin + "NextButtonText"

	LoginKeyLogoutDone                = "LogoutDone."
	LoginKeyLogoutDoneTitle           = LoginKeyLogoutDone + "Title"
	LoginKeyLogoutDoneDescription     = LoginKeyLogoutDone + "Description"
	LoginKeyLogoutDoneLoginButtonText = LoginKeyLogoutDone + "LoginButtonText"

	LoginKeyFooter            = "Footer."
	LoginKeyFooterTos         = LoginKeyFooter + "Tos"
	LoginKeyFooterTosLink     = LoginKeyFooter + "TosLink"
	LoginKeyFooterPrivacy     = LoginKeyFooter + "Privacy"
	LoginKeyFooterPrivacyLink = LoginKeyFooter + "PrivacyLink"
	LoginKeyFooterHelp        = LoginKeyFooter + "Help"
)

type CustomLoginText struct {
	models.ObjectRoot

	State    PolicyState
	Default  bool
	Language language.Tag

	SelectAccountScreenText      SelectAccountScreenText
	LoginScreenText              LoginScreenText
	PasswordScreenText           PasswordScreenText
	InitPasswordScreenText       InitPasswordScreenText
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
	Title                     string
	Description               string
	TitleLinkingProcess       string
	DescriptionLinkingProcess string
	OtherUser                 string
	SessionStateActive        string
	SessionStateInactive      string
	UserMustBeMemberOfOrg     string
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
	UserMustBeMemberOfOrg     string
}

type PasswordScreenText struct {
	Title          string
	Description    string
	PasswordLabel  string
	ResetLinkText  string
	BackButtonText string
	NextButtonText string
}

type InitPasswordScreenText struct {
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
