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

	Template string `json:"template" gorm:"column:template;primary_key"`
	Language string `json:"language" gorm:"column:language;primary_key"`
	Key      string `json:"key" gorm:"column:key;primary_key"`
	Text     string `json:"text" gorm:"column:text"`

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

func CustomTextViewsToMessageDomain(aggregateID, lang string, texts []*CustomTextView) *domain.CustomMessageText {
	langTag := language.Make(lang)
	result := &domain.CustomMessageText{
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
		if text.Key == domain.MessageTitle {
			result.Title = text.Text
		}
		if text.Key == domain.MessagePreHeader {
			result.PreHeader = text.Text
		}
		if text.Key == domain.MessageSubject {
			result.Subject = text.Text
		}
		if text.Key == domain.MessageGreeting {
			result.Greeting = text.Text
		}
		if text.Key == domain.MessageText {
			result.Text = text.Text
		}
		if text.Key == domain.MessageButtonText {
			result.ButtonText = text.Text
		}
		if text.Key == domain.MessageFooterText {
			result.FooterText = text.Text
		}
	}
	return result
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
		result.SelectAccount.Title = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescription {
		result.SelectAccount.Description = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountTitleLinkingProcess {
		result.SelectAccount.TitleLinking = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescriptionLinkingProcess {
		result.SelectAccount.DescriptionLinking = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountOtherUser {
		result.SelectAccount.OtherUser = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateActive {
		result.SelectAccount.SessionState0 = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateInactive {
		result.SelectAccount.SessionState1 = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountUserMustBeMemberOfOrg {
		result.SelectAccount.MustBeMemberOfOrg = text.Text
	}
}

func loginKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLoginTitle {
		result.Login.Title = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescription {
		result.Login.Description = text.Text
	}
	if text.Key == domain.LoginKeyLoginTitleLinkingProcess {
		result.Login.TitleLinking = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescriptionLinkingProcess {
		result.Login.DescriptionLinking = text.Text
	}
	if text.Key == domain.LoginKeyLoginNameLabel {
		result.Login.LoginNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyLoginUsernamePlaceHolder {
		result.Login.UsernamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginLoginnamePlaceHolder {
		result.Login.LoginnamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginExternalUserDescription {
		result.Login.ExternalUserDescription = text.Text
	}
	if text.Key == domain.LoginKeyLoginUserMustBeMemberOfOrg {
		result.Login.MustBeMemberOfOrg = text.Text
	}
	if text.Key == domain.LoginKeyLoginRegisterButtonText {
		result.Login.RegisterButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLoginNextButtonText {
		result.Login.NextButtonText = text.Text
	}
}

func passwordKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		result.Password.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordDescription {
		result.Password.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordLabel {
		result.Password.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetLinkText {
		result.Password.ResetLinkText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordBackButtonText {
		result.Password.BackButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordNextButtonText {
		result.Password.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordMinLength {
		result.Password.MinLength = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasUppercase {
		result.Password.HasUppercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasLowercase {
		result.Password.HasLowercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasNumber {
		result.Password.HasNumber = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasSymbol {
		result.Password.HasSymbol = text.Text
	}
	if text.Key == domain.LoginKeyPasswordConfirmation {
		result.Password.Confirmation = text.Text
	}
}

func usernameChangeKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeTitle {
		result.UsernameChange.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDescription {
		result.UsernameChange.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeUsernameLabel {
		result.UsernameChange.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeCancelButtonText {
		result.UsernameChange.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeNextButtonText {
		result.UsernameChange.NextButtonText = text.Text
	}
}

func usernameChangeDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeDoneTitle {
		result.UsernameChangeDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneDescription {
		result.UsernameChangeDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneNextButtonText {
		result.UsernameChangeDone.NextButtonText = text.Text
	}
}

func initPasswordKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordTitle {
		result.InitPassword.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDescription {
		result.InitPassword.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordCodeLabel {
		result.InitPassword.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordLabel {
		result.InitPassword.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordConfirmLabel {
		result.InitPassword.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNextButtonText {
		result.InitPassword.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordResendButtonText {
		result.InitPassword.ResendButtonText = text.Text
	}
}

func initPasswordDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordDoneTitle {
		result.InitPasswordDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneDescription {
		result.InitPasswordDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneNextButtonText {
		result.InitPasswordDone.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneCancelButtonText {
		result.InitPasswordDone.CancelButtonText = text.Text
	}
}

func emailVerificationKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationTitle {
		result.EmailVerification.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDescription {
		result.EmailVerification.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationCodeLabel {
		result.EmailVerification.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationNextButtonText {
		result.EmailVerification.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationResendButtonText {
		result.EmailVerification.ResendButtonText = text.Text
	}
}

func emailVerificationDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationDoneTitle {
		result.EmailVerificationDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneDescription {
		result.EmailVerificationDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneNextButtonText {
		result.EmailVerificationDone.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneCancelButtonText {
		result.EmailVerificationDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneLoginButtonText {
		result.EmailVerificationDone.LoginButtonText = text.Text
	}
}

func initializeUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitializeUserTitle {
		result.InitUser.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserDescription {
		result.InitUser.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserCodeLabel {
		result.InitUser.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordLabel {
		result.InitUser.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordConfirmLabel {
		result.InitUser.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserResendButtonText {
		result.InitUser.ResendButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNextButtonText {
		result.InitUser.NextButtonText = text.Text
	}
}

func initializeUserDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitUserDone {
		result.InitUserDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneDescription {
		result.InitUserDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneCancelButtonText {
		result.InitUserDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneNextButtonText {
		result.InitUserDone.NextButtonText = text.Text
	}
}

func initMFAPromptKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAPromptTitle {
		result.InitMFAPrompt.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptDescription {
		result.InitMFAPrompt.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptOTPOption {
		result.InitMFAPrompt.Provider0 = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptU2FOption {
		result.InitMFAPrompt.Provider1 = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptSkipButtonText {
		result.InitMFAPrompt.SkipButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptNextButtonText {
		result.InitMFAPrompt.NextButtonText = text.Text
	}
}

func initMFAOTPKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAOTPTitle {
		result.InitMFAOTP.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescription {
		result.InitMFAOTP.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescriptionOTP {
		result.InitMFAOTP.OTPDescription = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		result.InitMFAOTP.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		result.InitMFAOTP.SecretLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPNextButtonText {
		result.InitMFAOTP.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCancelButtonText {
		result.InitMFAOTP.CancelButtonText = text.Text
	}
}

func initMFAU2FKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAU2FTitle {
		result.InitMFAU2F.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FDescription {
		result.InitMFAU2F.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FTokenNameLabel {
		result.InitMFAU2F.TokenNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FRegisterTokenButtonText {
		result.InitMFAU2F.RegisterTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FNotSupported {
		result.InitMFAU2F.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FErrorRetry {
		result.InitMFAU2F.ErrorRetry = text.Text
	}
}

func initMFADoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFADoneTitle {
		result.InitMFADone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneDescription {
		result.InitMFADone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneCancelButtonText {
		result.InitMFADone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneNextButtonText {
		result.InitMFADone.NextButtonText = text.Text
	}
}

func mfaProvidersKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyMFAProvidersChooseOther {
		result.MFAProvider.ChooseOther = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersOTP {
		result.MFAProvider.Provider0 = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersU2F {
		result.MFAProvider.Provider1 = text.Text
	}
}

func verifyMFAOTPKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAOTPTitle {
		result.VerifyMFAOTP.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPDescription {
		result.VerifyMFAOTP.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPCodeLabel {
		result.VerifyMFAOTP.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPNextButtonText {
		result.VerifyMFAOTP.NextButtonText = text.Text
	}
}

func verifyMFAU2FKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAU2FTitle {
		result.VerifyMFAU2F.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FDescription {
		result.VerifyMFAU2F.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FValidateTokenText {
		result.VerifyMFAU2F.ValidateTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FNotSupported {
		result.VerifyMFAU2F.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FErrorRetry {
		result.VerifyMFAU2F.ErrorRetry = text.Text
	}
}

func passwordlessKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessTitle {
		result.Passwordless.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessDescription {
		result.Passwordless.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessLoginWithPwButtonText {
		result.Passwordless.LoginWithPwButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		result.Passwordless.ValidateTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessNotSupported {
		result.Passwordless.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessErrorRetry {
		result.Passwordless.ErrorRetry = text.Text
	}
}

func passwordChangeKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeTitle {
		result.PasswordChange.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDescription {
		result.PasswordChange.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeOldPasswordLabel {
		result.PasswordChange.OldPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordLabel {
		result.PasswordChange.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordConfirmLabel {
		result.PasswordChange.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeCancelButtonText {
		result.PasswordChange.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNextButtonText {
		result.PasswordChange.NextButtonText = text.Text
	}
}

func passwordChangeDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeDoneTitle {
		result.PasswordChangeDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneDescription {
		result.PasswordChangeDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneNextButtonText {
		result.PasswordChangeDone.NextButtonText = text.Text
	}
}

func passwordResetDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordResetDoneTitle {
		result.PasswordResetDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneDescription {
		result.PasswordResetDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneNextButtonText {
		result.PasswordResetDone.NextButtonText = text.Text
	}
}

func registrationOptionKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationOptionTitle {
		result.RegisterOption.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionDescription {
		result.RegisterOption.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionExternalLoginDescription {
		result.RegisterOption.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionUserNameButtonText {
		result.RegisterOption.RegisterUsernamePasswordButtonText = text.Text
	}
}

func registrationUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationUserTitle {
		result.RegistrationUser.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserDescription {
		result.RegistrationUser.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserDescriptionOrgRegister {
		result.RegistrationUser.DescriptionOrgRegister = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserFirstnameLabel {
		result.RegistrationUser.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLastnameLabel {
		result.RegistrationUser.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserEmailLabel {
		result.RegistrationUser.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserUsernameLabel {
		result.RegistrationUser.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLanguageLabel {
		result.RegistrationUser.LanguageLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserGenderLabel {
		result.RegistrationUser.GenderLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordLabel {
		result.RegistrationUser.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordConfirmLabel {
		result.RegistrationUser.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSAndPrivacyLabel {
		result.RegistrationUser.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSConfirm {
		result.RegistrationUser.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSLink {
		result.RegistrationUser.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSLinkText {
		result.RegistrationUser.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyConfirm {
		result.RegistrationUser.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyLink {
		result.RegistrationUser.PrivacyLink = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyLinkText {
		result.RegistrationUser.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserExternalLoginDescription {
		result.RegistrationUser.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserNextButtonText {
		result.RegistrationUser.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserBackButtonText {
		result.RegistrationUser.BackButtonText = text.Text
	}
}

func registrationOrgKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegisterOrgTitle {
		result.RegistrationOrg.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgDescription {
		result.RegistrationOrg.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgOrgNameLabel {
		result.RegistrationOrg.OrgNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgFirstnameLabel {
		result.RegistrationOrg.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgLastnameLabel {
		result.RegistrationOrg.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgUsernameLabel {
		result.RegistrationOrg.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgEmailLabel {
		result.RegistrationOrg.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordLabel {
		result.RegistrationOrg.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordConfirmLabel {
		result.RegistrationOrg.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSAndPrivacyLabel {
		result.RegistrationOrg.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSConfirm {
		result.RegistrationOrg.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSLink {
		result.RegistrationOrg.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSLinkText {
		result.RegistrationOrg.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyConfirm {
		result.RegistrationOrg.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyLink {
		result.RegistrationOrg.PrivacyLink = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyLinkText {
		result.RegistrationOrg.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgExternalLoginDescription {
		result.RegistrationOrg.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		result.RegistrationOrg.SaveButtonText = text.Text
	}
}

func linkingUserKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLinkingUserDoneTitle {
		result.LinkingUsersDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneDescription {
		result.LinkingUsersDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneCancelButtonText {
		result.LinkingUsersDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneNextButtonText {
		result.LinkingUsersDone.NextButtonText = text.Text
	}
}

func externalUserNotFoundKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyExternalNotFoundTitle {
		result.ExternalNotFoundOption.Title = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundDescription {
		result.ExternalNotFoundOption.Description = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundLinkButtonText {
		result.ExternalNotFoundOption.LinkButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundAutoRegisterButtonText {
		result.ExternalNotFoundOption.AutoRegisterButtonText = text.Text
	}
}

func successLoginKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeySuccessLoginTitle {
		result.LoginSuccess.Title = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginAutoRedirectDescription {
		result.LoginSuccess.AutoRedirectDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginRedirectedDescription {
		result.LoginSuccess.RedirectedDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginNextButtonText {
		result.LoginSuccess.NextButtonText = text.Text
	}
}

func logoutDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLogoutDoneTitle {
		result.LogoutDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneDescription {
		result.LogoutDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneLoginButtonText {
		result.LogoutDone.LoginButtonText = text.Text
	}
}

func footerKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyFooterTOS {
		result.Footer.TOS = text.Text
	}
	if text.Key == domain.LoginKeyFooterTOSLink {
		result.Footer.TOSLink = text.Text
	}
	if text.Key == domain.LoginKeyFooterPrivacy {
		result.Footer.PrivacyPolicy = text.Text
	}
	if text.Key == domain.LoginKeyFooterPrivacyLink {
		result.Footer.PrivacyPolicyLink = text.Text
	}
	if text.Key == domain.LoginKeyFooterHelp {
		result.Footer.Help = text.Text
	}
	if text.Key == domain.LoginKeyFooterHelpLink {
		result.Footer.HelpLink = text.Text
	}
}
