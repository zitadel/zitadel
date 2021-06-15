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

func CustomTextViewToModel(template *CustomTextView) *model.CustomTextView {
	lang := language.Make(template.Language)
	return &model.CustomTextView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		Template:     template.Template,
		Language:     lang,
		Key:          template.Key,
		Text:         template.Text,
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
		if strings.HasPrefix(text.Key, domain.LoginKeyResetPassword) {
			resetPasswordKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitializeUser) {
			initializeUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitializeDone) {
			initializeDoneKeyToDomain(text, result)
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
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAOTP) {
			verifyMFAOTPKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAU2F) {
			verifyMFAU2FKeyToDomain(text, result)
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
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordless) {
			passwordlessKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeySuccessLogin) {
			successLoginKeyToDomain(text, result)
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
	if text.Key == domain.LoginKeySelectAccountOtherUser {
		result.SelectAccountScreenText.OtherUser = text.Text
	}
}

func loginKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLoginTitle {
		result.LoginScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescription {
		result.LoginScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyLoginNameLabel {
		result.LoginScreenText.LoginNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyLoginExternalUserDescription {
		result.LoginScreenText.ExternalUserDescription = text.Text
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
}

func resetPasswordKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyResetPasswordTitle {
		result.ResetPasswordScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyResetPasswordDescription {
		result.ResetPasswordScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyResetPasswordNextButtonText {
		result.ResetPasswordScreenText.NextButtonText = text.Text
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
}

func initializeDoneKeyToDomain(text *CustomTextView, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitializeDone {
		result.InitializeDoneScreenText.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitializeDoneDescription {
		result.InitializeDoneScreenText.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitializeDoneAbortButtonText {
		result.InitializeDoneScreenText.AbortButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitializeDoneNextButtonText {
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
	if text.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		result.InitMFAOTPScreenText.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		result.InitMFAOTPScreenText.SecretLabel = text.Text
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
	if text.Key == domain.LoginKeyVerifyMFAU2FOtherOptionDescription {
		result.VerifyMFAU2FScreenText.OtherOptionDescription = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FOptionOTPButtonText {
		result.VerifyMFAU2FScreenText.OptionOTPButtonText = text.Text
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
	if text.Key == domain.LoginKeyRegistrationUserSaveButtonText {
		result.RegistrationUserScreenText.SaveButtonText = text.Text
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
	if text.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		result.RegistrationOrgScreenText.SaveButtonText = text.Text
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
