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
	SuccessLoginRedirectedDescription   string
	SuccessLoginNextButtonText          string
}

func (wm *CustomLoginTextReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.CustomTextSetEvent:
			if e.Template != domain.LoginCustomText || wm.Language != e.Language {
				continue
			}
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
			if strings.HasPrefix(e.Key, domain.LoginKeyResetPassword) {
				wm.handleResetPasswordScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeUser) {
				wm.handleInitializeUserScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeDone) {
				wm.handleInitializeDoneScreenSetEvent(e)
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
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordless) {
				wm.handlePasswordlessScreenSetEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeySuccessLogin) {
				wm.handleSuccessLoginScreenSetEvent(e)
				continue
			}
			wm.State = domain.PolicyStateActive
		case *policy.CustomTextRemovedEvent:
			if e.Key != domain.LoginCustomText || wm.Language != e.Language {
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
			if strings.HasPrefix(e.Key, domain.LoginKeyResetPassword) {
				wm.handleResetPasswordScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeUser) {
				wm.handleInitializeUserScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeyInitializeDone) {
				wm.handleInitializeDoneScreenRemoveEvent(e)
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
			if strings.HasPrefix(e.Key, domain.LoginKeyPasswordless) {
				wm.handlePasswordlessScreenRemoveEvent(e)
				continue
			}
			if strings.HasPrefix(e.Key, domain.LoginKeySuccessLogin) {
				wm.handleSuccessLoginScreenRemoveEvent(e)
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
		wm.SelectAccountDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeySelectAccountOtherUser {
		wm.SelectAccountOtherUser = e.Text
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
	if e.Key == domain.LoginKeySelectAccountOtherUser {
		wm.SelectAccountOtherUser = ""
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
	if e.Key == domain.LoginKeyLoginNameLabel {
		wm.LoginNameLabel = e.Text
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
	if e.Key == domain.LoginKeyLoginNameLabel {
		wm.LoginNameLabel = ""
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
}

func (wm *CustomLoginTextReadModel) handleResetPasswordScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyResetPasswordTitle {
		wm.ResetPasswordTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyResetPasswordDescription {
		wm.ResetPasswordDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyResetPasswordNextButtonText {
		wm.ResetPasswordNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleResetPasswordScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyResetPasswordTitle {
		wm.ResetPasswordTitle = ""
		return
	}
	if e.Key == domain.LoginKeyResetPasswordDescription {
		wm.ResetPasswordDescription = ""
		return
	}
	if e.Key == domain.LoginKeyResetPasswordNextButtonText {
		wm.ResetPasswordNextButtonText = ""
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
	if e.Key == domain.LoginKeyInitializeUserNewPassword {
		wm.InitializeNewPassword = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordConfirm {
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
	if e.Key == domain.LoginKeyInitializeUserNewPassword {
		wm.InitializeNewPassword = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeUserNewPasswordConfirm {
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

func (wm *CustomLoginTextReadModel) handleInitializeDoneScreenSetEvent(e *policy.CustomTextSetEvent) {
	if e.Key == domain.LoginKeyInitializeDoneTitle {
		wm.InitializeDoneTitle = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneDescription {
		wm.InitializeDoneDescription = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneAbortButtonText {
		wm.InitializeDoneAbortButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneNextButtonText {
		wm.InitializeDoneNextButtonText = e.Text
		return
	}
}

func (wm *CustomLoginTextReadModel) handleInitializeDoneScreenRemoveEvent(e *policy.CustomTextRemovedEvent) {
	if e.Key == domain.LoginKeyInitializeDoneTitle {
		wm.InitializeDoneTitle = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneDescription {
		wm.InitializeDoneDescription = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneAbortButtonText {
		wm.InitializeDoneAbortButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitializeDoneNextButtonText {
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
	if e.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		wm.InitMFAOTPSecretLabel = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		wm.InitMFAOTPCodeLabel = e.Text
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
	if e.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		wm.InitMFAOTPSecretLabel = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		wm.InitMFAOTPCodeLabel = ""
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
	if e.Key == domain.LoginKeyInitMFADoneAbortButtonText {
		wm.InitMFADoneAbortButtonText = e.Text
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneAbortButtonText {
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
	if e.Key == domain.LoginKeyInitMFADoneAbortButtonText {
		wm.InitMFADoneAbortButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneAbortButtonText {
		wm.InitMFADoneAbortButtonText = ""
		return
	}
	if e.Key == domain.LoginKeyInitMFADoneNextButtonText {
		wm.InitMFADoneNextButtonText = ""
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
	if e.Key == domain.LoginKeyVerifyMFAU2FOtherOptionDescription {
		wm.VerifyMFAU2FOtherOptionDescription = e.Text
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FOptionOTPButtonText {
		wm.VerifyMFAU2FOptionOTPButtonText = e.Text
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
	if e.Key == domain.LoginKeyVerifyMFAU2FOtherOptionDescription {
		wm.VerifyMFAU2FOtherOptionDescription = ""
	}
	if e.Key == domain.LoginKeyVerifyMFAU2FOptionOTPButtonText {
		wm.VerifyMFAU2FOptionOTPButtonText = ""
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
	if e.Key == domain.LoginKeyRegistrationUserSaveButtonText {
		wm.RegistrationUserSaveButtonText = e.Text
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
	if e.Key == domain.LoginKeyRegistrationUserSaveButtonText {
		wm.RegistrationUserSaveButtonText = ""
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
	if e.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		wm.RegisterOrgSaveButtonText = ""
		return
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
