package model

import (
	"encoding/json"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	CustomTextKeyAggregateID = "aggregate_id"
	CustomTextKeyTemplate    = "template"
	CustomTextKeyLanguage    = "language"
	CustomTextKeyKey         = "key"
)

type CustomTextView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`

	Template string `json:"Template" gorm:"column:template;primary_key"`
	Language string `json:"Language" gorm:"column:language;primary_key"`
	Key      string `json:"Key" gorm:"column:key;primary_key"`
	Text     string `json:"Text" gorm:"column:text"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func CustomTextViewFromModel(template *model.CustomTextView) *CustomTextView {
	return &CustomTextView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		Language:     template.Language.String(),
		Template:     template.Template,
		Key:          template.Key,
		Text:         template.Text,
	}
}

func CustomTextViewsToDomain(texts []*CustomTextView) []*domain.CustomText {
	result := make([]*domain.CustomText, len(texts))
	for i, text := range texts {
		result[i] = CustomTextViewToDomain(text)
	}
	return result
}

func CustomTextViewToDomain(text *CustomTextView) *domain.CustomText {
	lang := language.Make(text.Language)
	return &domain.CustomText{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  text.AggregateID,
			Sequence:     text.Sequence,
			CreationDate: text.CreationDate,
			ChangeDate:   text.ChangeDate,
		},
		Template: text.Template,
		Language: lang,
		Key:      text.Key,
		Text:     text.Text,
	}
}

func (i *CustomTextView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	switch event.Type {
	case es_model.CustomTextSet, org_es_model.CustomTextSet:
		i.setRootData(event)
		err = i.SetData(event)
		if err != nil {
			return err
		}
		i.ChangeDate = event.CreationDate
	}
	return err
}

func (r *CustomTextView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *CustomTextView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-3n9fs").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5CVaR", "Could not unmarshal data")
	}
	return nil
}

func (r *CustomTextView) IsMessageTemplate() bool {
	return r.Template == domain.InitCodeMessageType ||
		r.Template == domain.PasswordResetMessageType ||
		r.Template == domain.VerifyEmailMessageType ||
		r.Template == domain.VerifyPhoneMessageType ||
		r.Template == domain.DomainClaimedMessageType
}

func CustomTextViewsToLoginDomain(aggregateID, lang string, texts []*CustomTextView) *domain.CustomLoginText {
	langTag := language.Make(lang)
	result := &domain.CustomLoginText{
		ObjectRoot: models.ObjectRoot{
			AggregateID: aggregateID,
		},
		Language: langTag,
	}
	for _, text := range texts {
		if text.CreationDate.Before(result.CreationDate) {
			result.CreationDate = text.CreationDate
		}
		if text.ChangeDate.After(result.ChangeDate) {
			result.ChangeDate = text.ChangeDate
		}
		if strings.HasPrefix(text.Key, domain.LoginKeySelectAccount) {
			selectAccountKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLogin) {
			loginKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPassword) {
			passwordKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyUsernameChange) {
			usernameChangeKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyUsernameChangeDone) {
			usernameChangeDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitPassword) {
			initPasswordKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitPasswordDone) {
			initPasswordDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyEmailVerification) {
			emailVerificationKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyEmailVerificationDone) {
			emailVerificationDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitializeUser) {
			initializeUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitUserDone) {
			initializeUserDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAPrompt) {
			initMFAPromptKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAOTP) {
			initMFAOTPKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAU2F) {
			initMFAU2FKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFADone) {
			initMFADoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyMFAProviders) {
			mfaProvidersKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAOTP) {
			verifyMFAOTPKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAU2F) {
			verifyMFAU2FKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordless) {
			passwordlessKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordChange) {
			passwordChangeKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordChangeDone) {
			passwordChangeDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordResetDone) {
			passwordResetDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationOption) {
			registrationOptionKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationUser) {
			registrationUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationOrg) {
			registrationOrgKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLinkingUserDone) {
			linkingUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyExternalNotFound) {
			externalUserNotFoundKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeySuccessLogin) {
			successLoginKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLogoutDone) {
			logoutDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyFooter) {
			footerKeyToDomain(text, result)
		}
	}
	return result
}

func selectAccountKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeySelectAccountTitle {
		result.SelectAccountScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescription {
		result.SelectAccountScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountTitleLinkingProcess {
		result.SelectAccountScreenText.TitleLinkingProcess = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescriptionLinkingProcess {
		result.SelectAccountScreenText.DescriptionLinkingProcess = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountOtherUser {
		result.SelectAccountScreenText.OtherUser = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateActive {
		result.SelectAccountScreenText.SessionStateActive = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateInactive {
		result.SelectAccountScreenText.SessionStateInactive = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountUserMustBeMemberOfOrg {
		result.SelectAccountScreenText.UserMustBeMemberOfOrg = text.Text
	}
}

func loginKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLoginTitle {
		result.LoginScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescription {
		result.LoginScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyLoginTitleLinkingProcess {
		result.LoginScreenText.TitleLinkingProcess = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescriptionLinkingProcess {
		result.LoginScreenText.DescriptionLinkingProcess = text.Text
	}
	if text.Key == domain.LoginKeyLoginNameLabel {
		result.LoginScreenText.LoginNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyLoginUsernamePlaceHolder {
		result.LoginScreenText.UsernamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginLoginnamePlaceHolder {
		result.LoginScreenText.LoginnamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginExternalUserDescription {
		result.LoginScreenText.ExternalUserDescription = text.Text
	}
	if text.Key == domain.LoginKeyLoginUserMustBeMemberOfOrg {
		result.LoginScreenText.UserMustBeMemberOfOrg = text.Text
	}
	if text.Key == domain.LoginKeyLoginRegisterButtonText {
		result.LoginScreenText.RegisterButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLoginNextButtonText {
		result.LoginScreenText.NextButtonText = text.Text
	}
}

func passwordKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		result.PasswordScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordDescription {
		result.PasswordScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordLabel {
		result.PasswordScreenText.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetLinkText {
		result.PasswordScreenText.ResetLinkText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordBackButtonText {
		result.PasswordScreenText.BackButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordNextButtonText {
		result.PasswordScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordMinLength {
		result.PasswordScreenText.MinLength = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasUppercase {
		result.PasswordScreenText.HasUppercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasLowercase {
		result.PasswordScreenText.HasLowercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasNumber {
		result.PasswordScreenText.HasNumber = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasSymbol {
		result.PasswordScreenText.HasSymbol = text.Text
	}
	if text.Key == domain.LoginKeyPasswordConfirmation {
		result.PasswordScreenText.Confirmation = text.Text
	}
}

func usernameChangeKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeTitle {
		result.UsernameChangeScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDescription {
		result.UsernameChangeScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeUsernameLabel {
		result.UsernameChangeScreenText.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeCancelButtonText {
		result.UsernameChangeScreenText.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeNextButtonText {
		result.UsernameChangeScreenText.NextButtonText = text.Text
	}
}

func usernameChangeDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeDoneTitle {
		result.UsernameChangeDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneDescription {
		result.UsernameChangeDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneNextButtonText {
		result.UsernameChangeDoneScreenText.NextButtonText = text.Text
	}
}

func initPasswordKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordTitle {
		result.InitPasswordScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDescription {
		result.InitPasswordScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordCodeLabel {
		result.InitPasswordScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordLabel {
		result.InitPasswordScreenText.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordConfirmLabel {
		result.InitPasswordScreenText.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNextButtonText {
		result.InitPasswordScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordResendButtonText {
		result.InitPasswordScreenText.ResendButtonText = text.Text
	}
}

func initPasswordDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordDoneTitle {
		result.InitPasswordDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneDescription {
		result.InitPasswordDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneNextButtonText {
		result.InitPasswordDoneScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneCancelButtonText {
		result.InitPasswordDoneScreenText.CancelButtonText = text.Text
	}
}

func emailVerificationKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationTitle {
		result.EmailVerificationScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDescription {
		result.EmailVerificationScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationCodeLabel {
		result.EmailVerificationScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationNextButtonText {
		result.EmailVerificationScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationResendButtonText {
		result.EmailVerificationScreenText.ResendButtonText = text.Text
	}
}

func emailVerificationDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationDoneTitle {
		result.EmailVerificationDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneDescription {
		result.EmailVerificationDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneNextButtonText {
		result.EmailVerificationDoneScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneCancelButtonText {
		result.EmailVerificationDoneScreenText.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneLoginButtonText {
		result.EmailVerificationDoneScreenText.LoginButtonText = text.Text
	}
}

func initializeUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitializeUserTitle {
		result.InitializeUserScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserDescription {
		result.InitializeUserScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserCodeLabel {
		result.InitializeUserScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordLabel {
		result.InitializeUserScreenText.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordConfirmLabel {
		result.InitializeUserScreenText.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserResendButtonText {
		result.InitializeUserScreenText.ResendButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNextButtonText {
		result.InitializeUserScreenText.NextButtonText = text.Text
	}
}

func initializeUserDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitUserDone {
		result.InitializeDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneDescription {
		result.InitializeDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneAbortButtonText {
		result.InitializeDoneScreenText.AbortButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneNextButtonText {
		result.InitializeDoneScreenText.NextButtonText = text.Text
	}
}

func initMFAPromptKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAPromptTitle {
		result.InitMFAPromptScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptDescription {
		result.InitMFAPromptScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptOTPOption {
		result.InitMFAPromptScreenText.OTPOption = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptU2FOption {
		result.InitMFAPromptScreenText.U2FOption = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptSkipButtonText {
		result.InitMFAPromptScreenText.SkipButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptNextButtonText {
		result.InitMFAPromptScreenText.NextButtonText = text.Text
	}
}

func initMFAOTPKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAOTPTitle {
		result.InitMFAOTPScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescription {
		result.InitMFAOTPScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescriptionOTP {
		result.InitMFAOTPScreenText.DescriptionOTP = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		result.InitMFAOTPScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		result.InitMFAOTPScreenText.SecretLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPNextButtonText {
		result.InitMFAOTPScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCancelButtonText {
		result.InitMFAOTPScreenText.CancelButtonText = text.Text
	}
}

func initMFAU2FKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAU2FTitle {
		result.InitMFAU2FScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FDescription {
		result.InitMFAU2FScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FTokenNameLabel {
		result.InitMFAU2FScreenText.TokenNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FRegisterTokenButtonText {
		result.InitMFAU2FScreenText.RegisterTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FNotSupported {
		result.InitMFAU2FScreenText.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FErrorRetry {
		result.InitMFAU2FScreenText.ErrorRetry = text.Text
	}
}

func initMFADoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFADoneTitle {
		result.InitMFADoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneDescription {
		result.InitMFADoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneAbortButtonText {
		result.InitMFADoneScreenText.AbortButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneNextButtonText {
		result.InitMFADoneScreenText.NextButtonText = text.Text
	}
}

func mfaProvidersKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyMFAProvidersChooseOther {
		result.MFAProvidersText.ChooseOther = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersOTP {
		result.MFAProvidersText.OTP = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersU2F {
		result.MFAProvidersText.U2F = text.Text
	}
}

func verifyMFAOTPKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAOTPTitle {
		result.VerifyMFAOTPScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPDescription {
		result.VerifyMFAOTPScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPCodeLabel {
		result.VerifyMFAOTPScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPNextButtonText {
		result.VerifyMFAOTPScreenText.NextButtonText = text.Text
	}
}

func verifyMFAU2FKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAU2FTitle {
		result.VerifyMFAU2FScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FDescription {
		result.VerifyMFAU2FScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FValidateTokenText {
		result.VerifyMFAU2FScreenText.ValidateTokenText = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FNotSupported {
		result.VerifyMFAU2FScreenText.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FErrorRetry {
		result.VerifyMFAU2FScreenText.ErrorRetry = text.Text
	}
}

func passwordlessKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessTitle {
		result.PasswordlessScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessDescription {
		result.PasswordlessScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessLoginWithPwButtonText {
		result.PasswordlessScreenText.LoginWithPwButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		result.PasswordlessScreenText.ValidateTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessNotSupported {
		result.PasswordlessScreenText.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessErrorRetry {
		result.PasswordlessScreenText.ErrorRetry = text.Text
	}
}

func passwordChangeKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeTitle {
		result.PasswordChangeScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDescription {
		result.PasswordChangeScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeOldPasswordLabel {
		result.PasswordChangeScreenText.OldPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordLabel {
		result.PasswordChangeScreenText.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordConfirmLabel {
		result.PasswordChangeScreenText.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeCancelButtonText {
		result.PasswordChangeScreenText.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNextButtonText {
		result.PasswordChangeScreenText.NextButtonText = text.Text
	}
}

func passwordChangeDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeDoneTitle {
		result.PasswordChangeDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneDescription {
		result.PasswordChangeDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneNextButtonText {
		result.PasswordChangeDoneScreenText.NextButtonText = text.Text
	}
}

func passwordResetDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordResetDoneTitle {
		result.PasswordResetDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneDescription {
		result.PasswordResetDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneNextButtonText {
		result.PasswordResetDoneScreenText.NextButtonText = text.Text
	}
}

func registrationOptionKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationOptionTitle {
		result.RegistrationOptionScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionDescription {
		result.RegistrationOptionScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionExternalLoginDescription {
		result.RegistrationOptionScreenText.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionUserNameButtonText {
		result.RegistrationOptionScreenText.UserNameButtonText = text.Text
	}
}

func registrationUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationUserTitle {
		result.RegistrationUserScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserDescription {
		result.RegistrationUserScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserFirstnameLabel {
		result.RegistrationUserScreenText.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLastnameLabel {
		result.RegistrationUserScreenText.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserEmailLabel {
		result.RegistrationUserScreenText.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserUsernameLabel {
		result.RegistrationUserScreenText.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLanguageLabel {
		result.RegistrationUserScreenText.LanguageLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserGenderLabel {
		result.RegistrationUserScreenText.GenderLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordLabel {
		result.RegistrationUserScreenText.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordConfirmLabel {
		result.RegistrationUserScreenText.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSAndPrivacyLabel {
		result.RegistrationUserScreenText.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSConfirm {
		result.RegistrationUserScreenText.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSLink {
		result.RegistrationUserScreenText.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSLinkText {
		result.RegistrationUserScreenText.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyConfirm {
		result.RegistrationUserScreenText.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyLink {
		result.RegistrationUserScreenText.PrivacyLink = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyLinkText {
		result.RegistrationUserScreenText.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserExternalLoginDescription {
		result.RegistrationUserScreenText.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserNextButtonText {
		result.RegistrationUserScreenText.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserBackButtonText {
		result.RegistrationUserScreenText.BackButtonText = text.Text
	}
}

func registrationOrgKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegisterOrgTitle {
		result.RegistrationOrgScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgDescription {
		result.RegistrationOrgScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgOrgNameLabel {
		result.RegistrationOrgScreenText.OrgNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgFirstnameLabel {
		result.RegistrationOrgScreenText.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgLastnameLabel {
		result.RegistrationOrgScreenText.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgUsernameLabel {
		result.RegistrationOrgScreenText.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgEmailLabel {
		result.RegistrationOrgScreenText.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordLabel {
		result.RegistrationOrgScreenText.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordConfirmLabel {
		result.RegistrationOrgScreenText.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSAndPrivacyLabel {
		result.RegistrationOrgScreenText.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSConfirm {
		result.RegistrationOrgScreenText.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSLink {
		result.RegistrationOrgScreenText.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSLinkText {
		result.RegistrationOrgScreenText.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyConfirm {
		result.RegistrationOrgScreenText.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyLink {
		result.RegistrationOrgScreenText.PrivacyLink = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyLinkText {
		result.RegistrationOrgScreenText.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgExternalLoginDescription {
		result.RegistrationOrgScreenText.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		result.RegistrationOrgScreenText.SaveButtonText = text.Text
	}
}

func linkingUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLinkingUserDoneTitle {
		result.LinkingUserDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneDescription {
		result.LinkingUserDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneCancelButtonText {
		result.LinkingUserDoneScreenText.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneNextButtonText {
		result.LinkingUserDoneScreenText.NextButtonText = text.Text
	}
}

func externalUserNotFoundKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyExternalNotFoundTitle {
		result.ExternalUserNotFoundScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundDescription {
		result.ExternalUserNotFoundScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundLinkButtonText {
		result.ExternalUserNotFoundScreenText.LinkButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundAutoRegisterButtonText {
		result.ExternalUserNotFoundScreenText.AutoRegisterButtonText = text.Text
	}
}

func successLoginKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeySuccessLoginTitle {
		result.SuccessLoginScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginAutoRedirectDescription {
		result.SuccessLoginScreenText.AutoRedirectDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginRedirectedDescription {
		result.SuccessLoginScreenText.RedirectedDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginNextButtonText {
		result.SuccessLoginScreenText.NextButtonText = text.Text
	}
}

func logoutDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLogoutDoneTitle {
		result.LogoutDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneDescription {
		result.LogoutDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneLoginButtonText {
		result.LogoutDoneScreenText.LoginButtonText = text.Text
	}
}

func footerKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyFooterTos {
		result.FooterText.TOS = text.Text
	}
	if text.Key == domain.LoginKeyFooterTosLink {
		result.FooterText.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyFooterPrivacy {
		result.FooterText.PrivacyPolicy = text.Text
	}
	if text.Key == domain.LoginKeyFooterPrivacyLink {
		result.FooterText.PrivacyPolicyLink = text.Text
	}
	if text.Key == domain.LoginKeyFooterHelp {
		result.FooterText.Help = text.Text
	}
	if text.Key == domain.LoginKeyFooterHelpLink {
		result.FooterText.HelpLink = text.Text
	}
}
