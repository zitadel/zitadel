package command

import (
	"strings"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type CustomLoginTextReadModel struct {
	eventstore.WriteModel

	Language language.Tag
	State    domain.PolicyState

	SelectAccountTitle                     string
	SelectAccountDescription               string
	SelectAccountTitleLinkingProcess       string
	SelectAccountDescriptionLinkingProcess string
	SelectAccountOtherUser                 string
	SelectAccountSessionStateActive        string
	SelectAccountSessionStateInactive      string
	SelectAccountUserMustBeMemberOfOrg     string

	LoginTitle                     string
	LoginDescription               string
	LoginTitleLinkingProcess       string
	LoginDescriptionLinkingProcess string
	LoginNameLabel                 string
	LoginUsernamePlaceholder       string
	LoginLoginnamePlaceholder      string
	LoginRegisterButtonText        string
	LoginNextButtonText            string
	LoginExternalUserDescription   string
	LoginUserMustBeMemberOfOrg     string

	PasswordTitle          string
	PasswordDescription    string
	PasswordLabel          string
	PasswordResetLinkText  string
	PasswordBackButtonText string
	PasswordNextButtonText string
	PasswordMinLength      string
	PasswordHasUppercase   string
	PasswordHasLowercase   string
	PasswordHasNumber      string
	PasswordHasSymbol      string
	PasswordConfirmation   string

	UsernameChangeTitle            string
	UsernameChangeDescription      string
	UsernameChangeUsernameLabel    string
	UsernameChangeCancelButtonText string
	UsernameChangeNextButtonText   string

	UsernameChangeDoneTitle          string
	UsernameChangeDoneDescription    string
	UsernameChangeDoneNextButtonText string

	InitPasswordTitle                   string
	InitPasswordDescription             string
	InitPasswordCodeLabel               string
	InitPasswordNewPasswordLabel        string
	InitPasswordNewPasswordConfirmLabel string
	InitPasswordNextButtonText          string
	InitPasswordResendButtonText        string

	InitPasswordDoneTitle            string
	InitPasswordDoneDescription      string
	InitPasswordDoneNextButtonText   string
	InitPasswordDoneCancelButtonText string

	EmailVerificationTitle            string
	EmailVerificationDescription      string
	EmailVerificationCodeLabel        string
	EmailVerificationNextButtonText   string
	EmailVerificationResendButtonText string

	EmailVerificationDoneTitle            string
	EmailVerificationDoneDescription      string
	EmailVerificationDoneNextButtonText   string
	EmailVerificationDoneCancelButtonText string
	EmailVerificationDoneLoginButtonText  string

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

	InitMFAPromptTitle          string
	InitMFAPromptDescription    string
	InitMFAPromptOTPOption      string
	InitMFAPromptU2FOption      string
	InitMFAPromptSkipButtonText string
	InitMFAPromptNextButtonText string

	InitMFAOTPTitle            string
	InitMFAOTPDescription      string
	InitMFAOTPDescriptionOTP   string
	InitMFAOTPSecretLabel      string
	InitMFAOTPCodeLabel        string
	InitMFAOTPNextButtonText   string
	InitMFAOTPCancelButtonText string

	InitMFAU2FTitle                   string
	InitMFAU2FDescription             string
	InitMFAU2FTokenNameLabel          string
	InitMFAU2FRegisterTokenButtonText string
	InitMFAU2FNotSupported            string
	InitMFAU2FErrorRetry              string

	InitMFADoneTitle           string
	InitMFADoneDescription     string
	InitMFADoneAbortButtonText string
	InitMFADoneNextButtonText  string

	MFAProvidersChooseOther string
	MFAProvidersOTP         string
	MFAProvidersU2F         string

	VerifyMFAOTPTitle          string
	VerifyMFAOTPDescription    string
	VerifyMFAOTPCodeLabel      string
	VerifyMFAOTPNextButtonText string

	VerifyMFAU2FTitle             string
	VerifyMFAU2FDescription       string
	VerifyMFAU2FValidateTokenText string
	VerifyMFAU2FNotSupported      string
	VerifyMFAU2FErrorRetry        string

	PasswordlessTitle                   string
	PasswordlessDescription             string
	PasswordlessLoginWithPwButtonText   string
	PasswordlessValidateTokenButtonText string
	PasswordlessNotSupported            string
	PasswordlessErrorRetry              string

	PasswordChangeTitle                   string
	PasswordChangeDescription             string
	PasswordChangeOldPasswordLabel        string
	PasswordChangeNewPasswordLabel        string
	PasswordChangeNewPasswordConfirmLabel string
	PasswordChangeCancelButtonText        string
	PasswordChangeNextButtonText          string

	PasswordChangeDoneTitle          string
	PasswordChangeDoneDescription    string
	PasswordChangeDoneNextButtonText string

	PasswordResetDoneTitle          string
	PasswordResetDoneDescription    string
	PasswordResetDoneNextButtonText string

	RegistrationOptionTitle                    string
	RegistrationOptionDescription              string
	RegistrationOptionUserNameButtonText       string
	RegistrationOptionExternalLoginDescription string

	RegistrationUserTitle                  string
	RegistrationUserDescription            string
	RegistrationUserDescriptionOrgRegister string
	RegistrationUserFirstnameLabel         string
	RegistrationUserLastnameLabel          string
	RegistrationUserEmailLabel             string
	RegistrationUserUsernameLabel          string
	RegistrationUserLanguageLabel          string
	RegistrationUserGenderLabel            string
	RegistrationUserPasswordLabel          string
	RegistrationUserPasswordConfirmLabel   string
	RegistrationUserTOSAndPrivacyLabel     string
	RegistrationUserTOSConfirm             string
	RegistrationUserTOSLink                string
	RegistrationUserTOSLinkText            string
	RegistrationUserTOSConfirmAnd          string
	RegistrationUserPrivacyLink            string
	RegistrationUserPrivacyLinkText        string
	RegistrationUserNextButtonText         string
	RegistrationUserBackButtonText         string

	RegisterOrgTitle                string
	RegisterOrgDescription          string
	RegisterOrgOrgNameLabel         string
	RegisterOrgFirstnameLabel       string
	RegisterOrgLastnameLabel        string
	RegisterOrgUsernameLabel        string
	RegisterOrgEmailLabel           string
	RegisterOrgPasswordLabel        string
	RegisterOrgPasswordConfirmLabel string
	RegisterOrgTOSAndPrivacyLabel   string
	RegisterOrgTOSConfirm           string
	RegisterOrgTOSLinkText          string
	RegisterOrgTOSConfirmAnd        string
	RegisterOrgPrivacyLinkText      string
	RegisterOrgSaveButtonText       string

	LinkingUserDoneTitle            string
	LinkingUserDoneDescription      string
	LinkingUserDoneCancelButtonText string
	LinkingUserDoneNextButtonText   string

	ExternalUserNotFoundTitle                  string
	ExternalUserNotFoundDescription            string
	ExternalUserNotFoundLinkButtonText         string
	ExternalUserNotFoundAutoRegisterButtonText string

	SuccessLoginTitle                   string
	SuccessLoginAutoRedirectDescription string
	SuccessLoginRedirectedDescription   string
	SuccessLoginNextButtonText          string

	LogoutDoneTitle           string
	LogoutDoneDescription     string
	LogoutDoneLoginButtonText string

	FooterTOS           string
	FooterPrivacyPolicy string
	FooterHelp          string
	FooterHelpLink      string
}

func (wm *CustomLoginTextReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.CustomTextSetEvent:
			if e.Template != domain.LoginCustomText || wm.Language != e.Language {
				continue
			}
			wm.State = domain.PolicyStateActive
			if strings.HasPrefix(e.Key, domain.LoginKeySelectAccount) {
				wm.handleSelectAccountScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLogin) {
				wm.handleLoginScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPassword) {
				wm.handlePasswordScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyUsernameChange) {
				wm.handleUsernameChangeScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyUsernameChangeDone) {
				wm.handleUsernameChangeDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitPassword) {
				wm.handleInitPasswordScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitPasswordDone) {
				wm.handleInitPasswordDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyEmailVerification) {
				wm.handleEmailVerificationScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyEmailVerificationDone) {
				wm.handleEmailVerificationDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeUser) {
				wm.handleInitializeUserScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitUserDone) {
				wm.handleInitializeUserDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAPrompt) {
				wm.handleInitMFAPromptScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAOTP) {
				wm.handleInitMFAOTPScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAU2F) {
				wm.handleInitMFAU2FScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyMFAProviders) {
				wm.handleMFAProvidersTextSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFADone) {
				wm.handleInitMFADoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyVerifyMFAOTP) {
				wm.handleVerifyMFAOTPScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyVerifyMFAU2F) {
				wm.handleVerifyMFAU2FScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordless) {
				wm.handlePasswordlessScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordChange) {
				wm.handlePasswordChangeScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordChangeDone) {
				wm.handlePasswordChangeDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordResetDone) {
				wm.handlePasswordResetDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationOption) {
				wm.handleRegistrationOptionScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationUser) {
				wm.handleRegistrationUserScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationOrg) {
				wm.handleRegistrationOrgScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLinkingUserDone) {
				wm.handleLinkingUserDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyExternalNotFound) {
				wm.handleExternalUserNotFoundScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeySuccessLogin) {
				wm.handleSuccessLoginScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLogoutDone) {
				wm.handleLogoutDoneScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyFooter) {
				wm.handleFooterTextSetEvent(e)
				continue
			}
		case *policy.CustomTextRemovedEvent:
			if e.Template != domain.LoginCustomText || wm.Language != e.Language {
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeySelectAccount) {
				wm.handleSelectAccountScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLogin) {
				wm.handleLoginScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPassword) {
				wm.handlePasswordScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyUsernameChange) {
				wm.handleUsernameChangeRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyUsernameChangeDone) {
				wm.handleUsernameChangeDoneRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitPassword) {
				wm.handleInitPasswordScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitPasswordDone) {
				wm.handleInitPasswordDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyEmailVerification) {
				wm.handleEmailVerificationScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyEmailVerificationDone) {
				wm.handleEmailVerificationDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeUser) {
				wm.handleInitializeUserScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitUserDone) {
				wm.handleInitializeDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyMFAProviders) {
				wm.handleMFAProvidersTextRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAPrompt) {
				wm.handleInitMFAPromptScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAOTP) {
				wm.handleInitMFAOTPScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFAU2F) {
				wm.handleInitMFAU2FScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitMFADone) {
				wm.handleInitMFADoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyVerifyMFAOTP) {
				wm.handleVerifyMFAOTPScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyVerifyMFAU2F) {
				wm.handleVerifyMFAU2FScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordless) {
				wm.handlePasswordlessScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordChange) {
				wm.handlePasswordChangeScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordChangeDone) {
				wm.handlePasswordChangeDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordResetDone) {
				wm.handlePasswordResetDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationOption) {
				wm.handleRegistrationOptionScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationUser) {
				wm.handleRegistrationUserScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyRegistrationOrg) {
				wm.handleRegistrationOrgScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLinkingUserDone) {
				wm.handleLinkingUserDoneRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyExternalNotFound) {
				wm.handleExternalUserNotFoundScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeySuccessLogin) {
				wm.handleSuccessLoginScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyLogoutDone) {
				wm.handleLogoutDoneScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyFooter) {
				wm.handleFooterTextRemoveEvent(e)
				continue
			}
		case *policy.CustomTextTemplateRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *CustomLoginTextReadModel) handleSelectAccountScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeySelectAccountTitle {
		wm.SelectAccountTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountDescription {
		wm.SelectAccountDescriptionLinkingProcess = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountTitleLinkingProcess {
		wm.SelectAccountTitleLinkingProcess = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountDescriptionLinkingProcess {
		wm.SelectAccountDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountOtherUser {
		wm.SelectAccountOtherUser = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountSessionStateActive {
		wm.SelectAccountSessionStateActive = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountSessionStateInactive {
		wm.SelectAccountSessionStateInactive = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountUserMustBeMemberOfOrg {
		wm.SelectAccountUserMustBeMemberOfOrg = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleSelectAccountScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeySelectAccountTitle {
		wm.SelectAccountTitle = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountDescription {
		wm.SelectAccountDescription = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountTitleLinkingProcess {
		wm.SelectAccountTitleLinkingProcess = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountDescriptionLinkingProcess {
		wm.SelectAccountDescription = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountOtherUser {
		wm.SelectAccountOtherUser = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountSessionStateActive {
		wm.SelectAccountSessionStateActive = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountSessionStateInactive {
		wm.SelectAccountSessionStateInactive = ""
		return
	}
	if e.Key == domain.LoginKeySelectAccountUserMustBeMemberOfOrg {
		wm.SelectAccountUserMustBeMemberOfOrg = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLoginScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyLoginTitle {
		wm.LoginTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginDescription {
		wm.LoginDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginTitleLinkingProcess {
		wm.LoginTitleLinkingProcess = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginDescriptionLinkingProcess {
		wm.LoginDescriptionLinkingProcess = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginNameLabel {
		wm.LoginNameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginUsernamePlaceHolder {
		wm.LoginUsernamePlaceholder = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginLoginnamePlaceHolder {
		wm.LoginLoginnamePlaceholder = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginRegisterButtonText {
		wm.LoginRegisterButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginNextButtonText {
		wm.LoginNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginExternalUserDescription {
		wm.LoginExternalUserDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyLoginUserMustBeMemberOfOrg {
		wm.LoginUserMustBeMemberOfOrg = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLoginScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyLoginTitle {
		wm.LoginTitle = ""
		return
	}
	if e.Key == domain.LoginKeyLoginDescription {
		wm.LoginDescription = ""
		return
	}
	if e.Key == domain.LoginKeyLoginTitleLinkingProcess {
		wm.LoginTitleLinkingProcess = ""
		return
	}
	if e.Key == domain.LoginKeyLoginDescriptionLinkingProcess {
		wm.LoginDescriptionLinkingProcess = ""
		return
	}
	if e.Key == domain.LoginKeyLoginNameLabel {
		wm.LoginNameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyLoginUsernamePlaceHolder {
		wm.LoginUsernamePlaceholder = ""
		return
	}
	if e.Key == domain.LoginKeyLoginLoginnamePlaceHolder {
		wm.LoginLoginnamePlaceholder = ""
		return
	}
	if e.Key == domain.LoginKeyLoginRegisterButtonText {
		wm.LoginRegisterButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyLoginNextButtonText {
		wm.LoginNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyLoginExternalUserDescription {
		wm.LoginExternalUserDescription = ""
		return
	}
	if e.Key == domain.LoginKeyLoginUserMustBeMemberOfOrg {
		wm.LoginUserMustBeMemberOfOrg = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyPasswordTitle {
		wm.PasswordTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordDescription {
		wm.PasswordDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordLabel {
		wm.PasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordResetLinkText {
		wm.PasswordResetLinkText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordBackButtonText {
		wm.PasswordBackButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordNextButtonText {
		wm.PasswordNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordMinLength {
		wm.PasswordMinLength = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordHasUppercase {
		wm.PasswordHasUppercase = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordHasLowercase {
		wm.PasswordHasLowercase = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordHasNumber {
		wm.PasswordHasNumber = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordHasSymbol {
		wm.PasswordHasSymbol = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordConfirmation {
		wm.PasswordConfirmation = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyPasswordTitle {
		wm.PasswordTitle = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordDescription {
		wm.PasswordDescription = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordLabel {
		wm.PasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordResetLinkText {
		wm.PasswordResetLinkText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordBackButtonText {
		wm.PasswordBackButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordNextButtonText {
		wm.PasswordNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordMinLength {
		wm.PasswordMinLength = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordHasUppercase {
		wm.PasswordHasUppercase = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordHasLowercase {
		wm.PasswordHasLowercase = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordHasNumber {
		wm.PasswordHasNumber = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordHasSymbol {
		wm.PasswordHasSymbol = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordConfirmation {
		wm.PasswordConfirmation = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleUsernameChangeScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyUsernameChangeTitle {
		wm.UsernameChangeTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDescription {
		wm.UsernameChangeDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeUsernameLabel {
		wm.UsernameChangeUsernameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeCancelButtonText {
		wm.UsernameChangeCancelButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeNextButtonText {
		wm.UsernameChangeNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleUsernameChangeRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyUsernameChangeTitle {
		wm.UsernameChangeTitle = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDescription {
		wm.UsernameChangeDescription = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeUsernameLabel {
		wm.UsernameChangeUsernameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeCancelButtonText {
		wm.UsernameChangeCancelButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeNextButtonText {
		wm.UsernameChangeNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleUsernameChangeDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyUsernameChangeDoneTitle {
		wm.UsernameChangeDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDoneDescription {
		wm.UsernameChangeDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDoneNextButtonText {
		wm.UsernameChangeDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleUsernameChangeDoneRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyUsernameChangeDoneTitle {
		wm.UsernameChangeDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDoneDescription {
		wm.UsernameChangeDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyUsernameChangeDoneNextButtonText {
		wm.UsernameChangeDoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitPasswordScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitPasswordTitle {
		wm.InitPasswordTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDescription {
		wm.InitPasswordDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordCodeLabel {
		wm.InitPasswordCodeLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNewPasswordLabel {
		wm.InitPasswordNewPasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNewPasswordConfirmLabel {
		wm.InitPasswordNewPasswordConfirmLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNextButtonText {
		wm.InitPasswordNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordResendButtonText {
		wm.InitPasswordResendButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitPasswordScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitPasswordTitle {
		wm.InitPasswordTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDescription {
		wm.InitPasswordDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordCodeLabel {
		wm.InitPasswordCodeLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNewPasswordLabel {
		wm.InitPasswordNewPasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNewPasswordConfirmLabel {
		wm.InitPasswordNewPasswordConfirmLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordNextButtonText {
		wm.InitPasswordNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordResendButtonText {
		wm.InitPasswordResendButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitPasswordDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitPasswordDoneTitle {
		wm.InitPasswordDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneDescription {
		wm.InitPasswordDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneNextButtonText {
		wm.InitPasswordDoneNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneCancelButtonText {
		wm.InitPasswordDoneCancelButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitPasswordDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitPasswordDoneTitle {
		wm.InitPasswordDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneDescription {
		wm.InitPasswordDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneNextButtonText {
		wm.InitPasswordDoneNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitPasswordDoneCancelButtonText {
		wm.InitPasswordDoneCancelButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleEmailVerificationScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyEmailVerificationTitle {
		wm.EmailVerificationTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDescription {
		wm.EmailVerificationDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationCodeLabel {
		wm.EmailVerificationCodeLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationNextButtonText {
		wm.EmailVerificationNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationResendButtonText {
		wm.EmailVerificationResendButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleEmailVerificationScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyEmailVerificationTitle {
		wm.EmailVerificationTitle = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDescription {
		wm.EmailVerificationDescription = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationCodeLabel {
		wm.EmailVerificationCodeLabel = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationNextButtonText {
		wm.EmailVerificationNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationResendButtonText {
		wm.EmailVerificationResendButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleEmailVerificationDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyEmailVerificationDoneTitle {
		wm.EmailVerificationDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneDescription {
		wm.EmailVerificationDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneNextButtonText {
		wm.EmailVerificationDoneNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneCancelButtonText {
		wm.EmailVerificationDoneCancelButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneLoginButtonText {
		wm.EmailVerificationDoneLoginButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleEmailVerificationDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyEmailVerificationDoneTitle {
		wm.EmailVerificationDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneDescription {
		wm.EmailVerificationDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneNextButtonText {
		wm.EmailVerificationDoneNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneCancelButtonText {
		wm.EmailVerificationDoneCancelButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyEmailVerificationDoneLoginButtonText {
		wm.EmailVerificationDoneLoginButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitializeUserScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitializeUserTitle {
		wm.InitializeTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserDescription {
		wm.InitializeDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserCodeLabel {
		wm.InitializeCodeLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordLabel {
		wm.InitializeNewPassword = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordConfirmLabel {
		wm.InitializeNewPasswordConfirm = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserResendButtonText {
		wm.InitializeResendButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNextButtonText {
		wm.InitializeNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitializeUserScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitializeUserTitle {
		wm.InitializeTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserDescription {
		wm.InitializeDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserCodeLabel {
		wm.InitializeCodeLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordLabel {
		wm.InitializeNewPassword = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordConfirmLabel {
		wm.InitializeNewPasswordConfirm = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserResendButtonText {
		wm.InitializeResendButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNextButtonText {
		wm.InitializeNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitializeUserDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitUserDoneTitle {
		wm.InitializeDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneDescription {
		wm.InitializeDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneCancelButtonText {
		wm.InitializeDoneAbortButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneNextButtonText {
		wm.InitializeDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitializeDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitUserDoneTitle {
		wm.InitializeDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneDescription {
		wm.InitializeDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneCancelButtonText {
		wm.InitializeDoneAbortButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitUserDoneNextButtonText {
		wm.InitializeDoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAPromptScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitMFAPromptTitle {
		wm.InitMFAPromptTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptDescription {
		wm.InitMFAPromptDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptOTPOption {
		wm.InitMFAPromptOTPOption = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptU2FOption {
		wm.InitMFAPromptU2FOption = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptSkipButtonText {
		wm.InitMFAPromptSkipButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptNextButtonText {
		wm.InitMFAPromptNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAPromptScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitMFAPromptTitle {
		wm.InitMFAPromptTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptDescription {
		wm.InitMFAPromptDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptOTPOption {
		wm.InitMFAPromptOTPOption = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptU2FOption {
		wm.InitMFAPromptU2FOption = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptSkipButtonText {
		wm.InitMFAPromptSkipButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAPromptNextButtonText {
		wm.InitMFAPromptNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAOTPScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitMFAOTPTitle {
		wm.InitMFAOTPTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPDescription {
		wm.InitMFAOTPDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPDescriptionOTP {
		wm.InitMFAOTPDescriptionOTP = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		wm.InitMFAOTPSecretLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		wm.InitMFAOTPCodeLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPNextButtonText {
		wm.InitMFAOTPNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCancelButtonText {
		wm.InitMFAOTPCancelButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAOTPScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitMFAOTPTitle {
		wm.InitMFAOTPTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPDescription {
		wm.InitMFAOTPDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPDescriptionOTP {
		wm.InitMFAOTPDescriptionOTP = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		wm.InitMFAOTPSecretLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		wm.InitMFAOTPCodeLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPNextButtonText {
		wm.InitMFAOTPNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCancelButtonText {
		wm.InitMFAOTPCancelButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAU2FScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitMFAU2FTitle {
		wm.InitMFAU2FTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FDescription {
		wm.InitMFAU2FDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FTokenNameLabel {
		wm.InitMFAU2FTokenNameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FRegisterTokenButtonText {
		wm.InitMFAU2FRegisterTokenButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FNotSupported {
		wm.InitMFAU2FNotSupported = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FErrorRetry {
		wm.InitMFAU2FErrorRetry = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFAU2FScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitMFAU2FTitle {
		wm.InitMFAU2FTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FDescription {
		wm.InitMFAU2FDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FTokenNameLabel {
		wm.InitMFAU2FTokenNameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FRegisterTokenButtonText {
		wm.InitMFAU2FRegisterTokenButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FNotSupported {
		wm.InitMFAU2FNotSupported = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAU2FErrorRetry {
		wm.InitMFAU2FErrorRetry = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFADoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitMFADoneTitle {
		wm.InitMFADoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneDescription {
		wm.InitMFADoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneCancelButtonText {
		wm.InitMFADoneAbortButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneNextButtonText {
		wm.InitMFADoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitMFADoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitMFADoneTitle {
		wm.InitMFADoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneDescription {
		wm.InitMFADoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneCancelButtonText {
		wm.InitMFADoneAbortButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneNextButtonText {
		wm.InitMFADoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleMFAProvidersTextSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyMFAProvidersChooseOther {
		wm.MFAProvidersChooseOther = e.Text
		return
	}
	if e.Key == domain.LoginKeyMFAProvidersOTP {
		wm.MFAProvidersOTP = e.Text
		return
	}
	if e.Key == domain.LoginKeyMFAProvidersU2F {
		wm.MFAProvidersU2F = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleMFAProvidersTextRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyMFAProvidersChooseOther {
		wm.MFAProvidersChooseOther = ""
		return
	}
	if e.Key == domain.LoginKeyMFAProvidersOTP {
		wm.MFAProvidersOTP = ""
		return
	}
	if e.Key == domain.LoginKeyMFAProvidersU2F {
		wm.MFAProvidersU2F = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleVerifyMFAOTPScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyVerifyMFAOTPTitle {
		wm.VerifyMFAOTPTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPDescription {
		wm.VerifyMFAOTPDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPCodeLabel {
		wm.VerifyMFAOTPCodeLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPNextButtonText {
		wm.VerifyMFAOTPNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleVerifyMFAOTPScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyVerifyMFAOTPTitle {
		wm.VerifyMFAOTPTitle = ""
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPDescription {
		wm.VerifyMFAOTPDescription = ""
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPCodeLabel {
		wm.VerifyMFAOTPCodeLabel = ""
		return
	}
	if e.Key == domain.LoginKeyVerifyMFAOTPNextButtonText {
		wm.VerifyMFAOTPNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleVerifyMFAU2FScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyVerifyMFAU2FTitle {
		wm.VerifyMFAU2FTitle = e.Text
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FDescription {
		wm.VerifyMFAU2FDescription = e.Text
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FValidateTokenText {
		wm.VerifyMFAU2FValidateTokenText = e.Text
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FNotSupported {
		wm.VerifyMFAU2FNotSupported = e.Text
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FErrorRetry {
		wm.VerifyMFAU2FErrorRetry = e.Text
	}
}

func (wm *CustomLoginTextReadModel) handleVerifyMFAU2FScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyVerifyMFAU2FTitle {
		wm.VerifyMFAU2FTitle = ""
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FDescription {
		wm.VerifyMFAU2FDescription = ""
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FValidateTokenText {
		wm.VerifyMFAU2FValidateTokenText = ""
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FNotSupported {
		wm.VerifyMFAU2FNotSupported = ""
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FErrorRetry {
		wm.VerifyMFAU2FErrorRetry = ""
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordlessScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyPasswordlessTitle {
		wm.PasswordlessTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordlessDescription {
		wm.PasswordlessDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordlessLoginWithPwButtonText {
		wm.PasswordlessLoginWithPwButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		wm.PasswordlessValidateTokenButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordlessNotSupported {
		wm.PasswordlessNotSupported = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordlessErrorRetry {
		wm.PasswordlessErrorRetry = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordlessScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyPasswordlessTitle {
		wm.PasswordlessTitle = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordlessDescription {
		wm.PasswordlessDescription = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordlessLoginWithPwButtonText {
		wm.PasswordlessLoginWithPwButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		wm.PasswordlessValidateTokenButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordlessNotSupported {
		wm.PasswordlessNotSupported = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordlessErrorRetry {
		wm.PasswordlessErrorRetry = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordChangeScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyPasswordChangeTitle {
		wm.PasswordChangeTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDescription {
		wm.PasswordChangeDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeOldPasswordLabel {
		wm.PasswordChangeOldPasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNewPasswordLabel {
		wm.PasswordChangeNewPasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNewPasswordConfirmLabel {
		wm.PasswordChangeNewPasswordConfirmLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeCancelButtonText {
		wm.PasswordChangeCancelButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNextButtonText {
		wm.PasswordChangeNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordChangeScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyPasswordChangeTitle {
		wm.PasswordChangeTitle = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDescription {
		wm.PasswordChangeDescription = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeOldPasswordLabel {
		wm.PasswordChangeOldPasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNewPasswordLabel {
		wm.PasswordChangeNewPasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNewPasswordConfirmLabel {
		wm.PasswordChangeNewPasswordConfirmLabel = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeCancelButtonText {
		wm.PasswordChangeCancelButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeNextButtonText {
		wm.PasswordChangeNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordChangeDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyPasswordChangeDoneTitle {
		wm.PasswordChangeDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDoneDescription {
		wm.PasswordChangeDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDoneNextButtonText {
		wm.PasswordChangeDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordChangeDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyPasswordChangeDoneTitle {
		wm.PasswordChangeDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDoneDescription {
		wm.PasswordChangeDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordChangeDoneNextButtonText {
		wm.PasswordChangeDoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordResetDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyPasswordResetDoneTitle {
		wm.PasswordResetDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordResetDoneDescription {
		wm.PasswordResetDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyPasswordResetDoneNextButtonText {
		wm.PasswordResetDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handlePasswordResetDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyPasswordResetDoneTitle {
		wm.PasswordResetDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordResetDoneDescription {
		wm.PasswordResetDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyPasswordResetDoneNextButtonText {
		wm.PasswordResetDoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationOptionScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyRegistrationOptionTitle {
		wm.RegistrationOptionTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionDescription {
		wm.RegistrationOptionDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionUserNameButtonText {
		wm.RegistrationOptionUserNameButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionExternalLoginDescription {
		wm.RegistrationOptionExternalLoginDescription = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationOptionScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyRegistrationOptionTitle {
		wm.RegistrationOptionTitle = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionDescription {
		wm.RegistrationOptionDescription = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionUserNameButtonText {
		wm.RegistrationOptionUserNameButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationOptionExternalLoginDescription {
		wm.RegistrationOptionExternalLoginDescription = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationUserScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyRegistrationUserTitle {
		wm.RegistrationUserTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserDescription {
		wm.RegistrationUserDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserDescriptionOrgRegister {
		wm.RegistrationUserDescriptionOrgRegister = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserFirstnameLabel {
		wm.RegistrationUserFirstnameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserLastnameLabel {
		wm.RegistrationUserLastnameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserEmailLabel {
		wm.RegistrationUserEmailLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserUsernameLabel {
		wm.RegistrationUserUsernameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserLanguageLabel {
		wm.RegistrationUserLanguageLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserGenderLabel {
		wm.RegistrationUserGenderLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPasswordLabel {
		wm.RegistrationUserPasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPasswordConfirmLabel {
		wm.RegistrationUserPasswordConfirmLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSAndPrivacyLabel {
		wm.RegistrationUserTOSAndPrivacyLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSConfirm {
		wm.RegistrationUserTOSConfirm = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSLinkText {
		wm.RegistrationUserTOSLinkText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSConfirmAnd {
		wm.RegistrationUserTOSConfirmAnd = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPrivacyLinkText {
		wm.RegistrationUserPrivacyLinkText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserNextButtonText {
		wm.RegistrationUserNextButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserBackButtonText {
		wm.RegistrationUserBackButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationUserScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyRegistrationUserTitle {
		wm.RegistrationUserTitle = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserDescription {
		wm.RegistrationUserDescription = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserDescriptionOrgRegister {
		wm.RegistrationUserDescriptionOrgRegister = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserFirstnameLabel {
		wm.RegistrationUserFirstnameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserLastnameLabel {
		wm.RegistrationUserLastnameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserEmailLabel {
		wm.RegistrationUserEmailLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserUsernameLabel {
		wm.RegistrationUserUsernameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserLanguageLabel {
		wm.RegistrationUserLanguageLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserGenderLabel {
		wm.RegistrationUserGenderLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPasswordLabel {
		wm.RegistrationUserPasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPasswordConfirmLabel {
		wm.RegistrationUserPasswordConfirmLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSAndPrivacyLabel {
		wm.RegistrationUserTOSAndPrivacyLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSConfirm {
		wm.RegistrationUserTOSConfirm = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSLinkText {
		wm.RegistrationUserTOSLinkText = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserTOSConfirmAnd {
		wm.RegistrationUserTOSConfirmAnd = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserPrivacyLinkText {
		wm.RegistrationUserPrivacyLinkText = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserNextButtonText {
		wm.RegistrationUserNextButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyRegistrationUserBackButtonText {
		wm.RegistrationUserBackButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationOrgScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyRegisterOrgTitle {
		wm.RegisterOrgTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgDescription {
		wm.RegisterOrgDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgOrgNameLabel {
		wm.RegisterOrgOrgNameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgFirstnameLabel {
		wm.RegisterOrgFirstnameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgLastnameLabel {
		wm.RegisterOrgLastnameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgUsernameLabel {
		wm.RegisterOrgUsernameLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgEmailLabel {
		wm.RegisterOrgEmailLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPasswordLabel {
		wm.RegisterOrgPasswordLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPasswordConfirmLabel {
		wm.RegisterOrgPasswordConfirmLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSAndPrivacyLabel {
		wm.RegisterOrgTOSAndPrivacyLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSConfirm {
		wm.RegisterOrgTOSConfirm = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSLinkText {
		wm.RegisterOrgTOSLinkText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTosConfirmAnd {
		wm.RegisterOrgTOSConfirmAnd = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPrivacyLinkText {
		wm.RegisterOrgPrivacyLinkText = e.Text
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		wm.RegisterOrgSaveButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleRegistrationOrgScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyRegisterOrgTitle {
		wm.RegisterOrgTitle = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgDescription {
		wm.RegisterOrgDescription = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgOrgNameLabel {
		wm.RegisterOrgOrgNameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgFirstnameLabel {
		wm.RegisterOrgFirstnameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgLastnameLabel {
		wm.RegisterOrgLastnameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgUsernameLabel {
		wm.RegisterOrgUsernameLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgEmailLabel {
		wm.RegisterOrgEmailLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPasswordLabel {
		wm.RegisterOrgPasswordLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPasswordConfirmLabel {
		wm.RegisterOrgPasswordConfirmLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSAndPrivacyLabel {
		wm.RegisterOrgTOSAndPrivacyLabel = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSConfirm {
		wm.RegisterOrgTOSConfirm = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTOSLinkText {
		wm.RegisterOrgTOSLinkText = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgTosConfirmAnd {
		wm.RegisterOrgTOSConfirmAnd = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgPrivacyLinkText {
		wm.RegisterOrgPrivacyLinkText = ""
		return
	}
	if e.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		wm.RegisterOrgSaveButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLinkingUserDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyLinkingUserDoneTitle {
		wm.LinkingUserDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneDescription {
		wm.LinkingUserDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneCancelButtonText {
		wm.LinkingUserDoneCancelButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneNextButtonText {
		wm.LinkingUserDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLinkingUserDoneRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyLinkingUserDoneTitle {
		wm.LinkingUserDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneDescription {
		wm.LinkingUserDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneCancelButtonText {
		wm.LinkingUserDoneCancelButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyLinkingUserDoneNextButtonText {
		wm.LinkingUserDoneNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleExternalUserNotFoundScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyExternalNotFoundTitle {
		wm.ExternalUserNotFoundTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundDescription {
		wm.ExternalUserNotFoundDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundLinkButtonText {
		wm.ExternalUserNotFoundLinkButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundAutoRegisterButtonText {
		wm.ExternalUserNotFoundAutoRegisterButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleExternalUserNotFoundScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyExternalNotFoundTitle {
		wm.ExternalUserNotFoundTitle = ""
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundDescription {
		wm.ExternalUserNotFoundDescription = ""
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundLinkButtonText {
		wm.ExternalUserNotFoundLinkButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyExternalNotFoundAutoRegisterButtonText {
		wm.ExternalUserNotFoundAutoRegisterButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleSuccessLoginScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeySuccessLoginTitle {
		wm.SuccessLoginTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeySuccessLoginAutoRedirectDescription {
		wm.SuccessLoginAutoRedirectDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeySuccessLoginRedirectedDescription {
		wm.SuccessLoginRedirectedDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeySuccessLoginNextButtonText {
		wm.SuccessLoginNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleSuccessLoginScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeySuccessLoginTitle {
		wm.SuccessLoginTitle = ""
		return
	}
	if e.Key == domain.LoginKeySuccessLoginAutoRedirectDescription {
		wm.SuccessLoginAutoRedirectDescription = ""
		return
	}
	if e.Key == domain.LoginKeySuccessLoginRedirectedDescription {
		wm.SuccessLoginRedirectedDescription = ""
		return
	}
	if e.Key == domain.LoginKeySuccessLoginNextButtonText {
		wm.SuccessLoginNextButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLogoutDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyLogoutDoneTitle {
		wm.LogoutDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyLogoutDoneDescription {
		wm.LogoutDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyLogoutDoneLoginButtonText {
		wm.LogoutDoneLoginButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleLogoutDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyLogoutDoneTitle {
		wm.LogoutDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyLogoutDoneDescription {
		wm.LogoutDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyLogoutDoneLoginButtonText {
		wm.LogoutDoneLoginButtonText = ""
		return
	}
}

func (wm *CustomLoginTextReadModel) handleFooterTextSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyFooterTOS {
		wm.FooterTOS = e.Text
		return
	}
	if e.Key == domain.LoginKeyFooterPrivacyPolicy {
		wm.FooterPrivacyPolicy = e.Text
		return
	}
	if e.Key == domain.LoginKeyFooterHelp {
		wm.FooterHelp = e.Text
		return
	}
	if e.Key == domain.LoginKeyFooterHelpLink {
		wm.FooterHelpLink = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleFooterTextRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyFooterTOS {
		wm.FooterTOS = ""
		return
	}
	if e.Key == domain.LoginKeyFooterPrivacyPolicy {
		wm.FooterPrivacyPolicy = ""
		return
	}
	if e.Key == domain.LoginKeyFooterHelp {
		wm.FooterHelp = ""
		return
	}
	if e.Key == domain.LoginKeyFooterHelpLink {
		wm.FooterHelpLink = ""
		return
	}
}
