package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) SetDefaultLoginText(ctx context.Context, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMailText, err := c.setDefaultLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMailText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMailText.WriteModel), nil
}

func (c *Commands) setDefaultLoginText(ctx context.Context, iamAgg *eventstore.Aggregate, text *domain.CustomLoginText) ([]eventstore.EventPusher, *IAMCustomLoginTextReadModel, error) {
	if !text.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomText.Invalid")
	}

	loginText, err := c.defaultLoginTextWriteModelByID(ctx, text.Language)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	if loginText.SelectAccountTitle != text.SelectAccountTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeySelectAccount, text.SelectAccountTitle, text.Language))
	}
	if loginText.SelectAccountDescription != text.SelectAccountDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, text.SelectAccountDescription, text.Language))
	}
	if loginText.SelectAccountOtherUser != text.SelectAccountOtherUser {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, text.SelectAccountOtherUser, text.Language))
	}
	if loginText.LoginTitle != text.LoginTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginTitle, text.LoginTitle, text.Language))
	}
	if loginText.LoginDescription != text.LoginDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginDescription, text.LoginDescription, text.Language))
	}
	if loginText.LoginNameLabel != text.LoginNameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, text.LoginNameLabel, text.Language))
	}
	if loginText.LoginRegisterButtonText != text.LoginRegisterButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, text.LoginRegisterButtonText, text.Language))
	}
	if loginText.LoginNextButtonText != text.LoginNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, text.LoginNextButtonText, text.Language))
	}
	if loginText.LoginExternalUserDescription != text.LoginExternalUserDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, text.LoginExternalUserDescription, text.Language))
	}
	if loginText.PasswordTitle != text.PasswordTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordTitle, text.PasswordTitle, text.Language))
	}
	if loginText.PasswordDescription != text.PasswordDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordDescription, text.PasswordDescription, text.Language))
	}
	if loginText.PasswordLabel != text.PasswordLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordLabel, text.PasswordLabel, text.Language))
	}
	if loginText.PasswordResetLinkText != text.PasswordResetLinkText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, text.PasswordResetLinkText, text.Language))
	}
	if loginText.PasswordBackButtonText != text.PasswordBackButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, text.PasswordBackButtonText, text.Language))
	}
	if loginText.PasswordNextButtonText != text.PasswordNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, text.PasswordNextButtonText, text.Language))
	}
	if loginText.PasswordResetTitle != text.PasswordResetTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordTitle, text.PasswordResetTitle, text.Language))
	}
	if loginText.PasswordResetDescription != text.PasswordResetDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordDescription, text.PasswordResetDescription, text.Language))
	}
	if loginText.PasswordResetNextButtonText != text.PasswordResetNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordNextButtonText, text.PasswordResetNextButtonText, text.Language))
	}
	if loginText.InitializeTitle != text.InitializeTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, text.InitializeTitle, text.Language))
	}
	if loginText.InitializeDescription != text.InitializeDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, text.InitializeDescription, text.Language))
	}
	if loginText.InitializeCodeLabel != text.InitializeCodeLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, text.InitializeCodeLabel, text.Language))
	}
	if loginText.InitializeNewPassword != text.InitializeNewPassword {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPassword, text.InitializeNewPassword, text.Language))
	}
	if loginText.InitializeNewPasswordConfirm != text.InitializeNewPasswordConfirm {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirm, text.InitializeNewPasswordConfirm, text.Language))
	}
	if loginText.InitializeResendButtonText != text.InitializeResendButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, text.InitializeResendButtonText, text.Language))
	}
	if loginText.InitializeNextButtonText != text.InitializeNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, text.InitializeNextButtonText, text.Language))
	}
	if loginText.InitializeDoneTitle != text.InitializeDoneTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneTitle, text.InitializeDoneTitle, text.Language))
	}
	if loginText.InitializeDoneDescription != text.InitializeDoneDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneDescription, text.InitializeDoneDescription, text.Language))
	}
	if loginText.InitializeDoneAbortButtonText != text.InitializeDoneAbortButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneAbortButtonText, text.InitializeDoneAbortButtonText, text.Language))
	}
	if loginText.InitializeDoneNextButtonText != text.InitializeDoneNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneNextButtonText, text.InitializeDoneNextButtonText, text.Language))
	}
	if loginText.InitMFATitle != text.InitMFATitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, text.InitMFATitle, text.Language))
	}
	if loginText.InitMFADescription != text.InitMFADescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, text.InitMFADescription, text.Language))
	}
	if loginText.InitMFAOTPOption != text.InitMFAOTPOption {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, text.InitMFAOTPOption, text.Language))
	}
	if loginText.InitMFAU2FOption != text.InitMFAU2FOption {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, text.InitMFAU2FOption, text.Language))
	}
	if loginText.InitMFASkipButtonText != text.InitMFASkipButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, text.InitMFASkipButtonText, text.Language))
	}
	if loginText.InitMFANextButtonText != text.InitMFANextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, text.InitMFANextButtonText, text.Language))
	}
	if loginText.InitMFAOTPTitle != text.InitMFAOTPTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, text.InitMFAOTPTitle, text.Language))
	}
	if loginText.InitMFAOTPDescription != text.InitMFAOTPDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, text.InitMFAOTPDescription, text.Language))
	}
	if loginText.InitMFAOTPSecretLabel != text.InitMFAOTPSecretLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, text.InitMFAOTPSecretLabel, text.Language))
	}
	if loginText.InitMFAOTPCodeLabel != text.InitMFAOTPCodeLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, text.InitMFAOTPCodeLabel, text.Language))
	}
	if loginText.InitMFAU2FTitle != text.InitMFAU2FTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, text.InitMFAU2FTitle, text.Language))
	}
	if loginText.InitMFAU2FDescription != text.InitMFAU2FDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, text.InitMFAU2FDescription, text.Language))
	}
	if loginText.InitMFAU2FTokenNameLabel != text.InitMFAU2FTokenNameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, text.InitMFAU2FTokenNameLabel, text.Language))
	}
	if loginText.InitMFAU2FRegisterTokenButtonText != text.InitMFAU2FRegisterTokenButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, text.InitMFAU2FRegisterTokenButtonText, text.Language))
	}
	if loginText.InitMFADoneTitle != text.InitMFADoneTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, text.InitMFADoneTitle, text.Language))
	}
	if loginText.InitMFADoneDescription != text.InitMFADoneDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, text.InitMFADoneDescription, text.Language))
	}
	if loginText.InitMFADoneAbortButtonText != text.InitMFADoneAbortButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneAbortButtonText, text.InitMFADoneAbortButtonText, text.Language))
	}
	if loginText.InitMFADoneNextButtonText != text.InitMFADoneNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, text.InitMFADoneNextButtonText, text.Language))
	}
	if loginText.VerifyMFAOTPTitle != text.VerifyMFAOTPTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, text.VerifyMFAOTPTitle, text.Language))
	}
	if loginText.VerifyMFAOTPDescription != text.VerifyMFAOTPDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, text.VerifyMFAOTPDescription, text.Language))
	}
	if loginText.VerifyMFAOTPCodeLabel != text.VerifyMFAOTPCodeLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, text.VerifyMFAOTPCodeLabel, text.Language))
	}
	if loginText.VerifyMFAOTPNextButtonText != text.VerifyMFAOTPNextButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, text.VerifyMFAOTPNextButtonText, text.Language))
	}
	if loginText.VerifyMFAU2FTitle != text.VerifyMFAU2FTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, text.VerifyMFAU2FTitle, text.Language))
	}
	if loginText.VerifyMFAU2FDescription != text.VerifyMFAU2FDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, text.VerifyMFAU2FDescription, text.Language))
	}
	if loginText.VerifyMFAU2FValidateTokenText != text.VerifyMFAU2FValidateTokenText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, text.VerifyMFAU2FValidateTokenText, text.Language))
	}
	if loginText.VerifyMFAU2FOtherOptionDescription != text.VerifyMFAU2FOtherOptionDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOtherOptionDescription, text.VerifyMFAU2FOtherOptionDescription, text.Language))
	}
	if loginText.VerifyMFAU2FOptionOTPButtonText != text.VerifyMFAU2FOptionOTPButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOptionOTPButtonText, text.VerifyMFAU2FOptionOTPButtonText, text.Language))
	}
	if loginText.RegistrationOptionTitle != text.RegistrationOptionTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, text.RegistrationOptionTitle, text.Language))
	}
	if loginText.RegistrationOptionDescription != text.RegistrationOptionDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, text.RegistrationOptionDescription, text.Language))
	}
	if loginText.RegistrationOptionUserNameButtonText != text.RegistrationOptionUserNameButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, text.RegistrationOptionUserNameButtonText, text.Language))
	}
	if loginText.RegistrationOptionExternalLoginDescription != text.RegistrationOptionExternalLoginDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, text.RegistrationOptionExternalLoginDescription, text.Language))
	}
	if loginText.RegistrationUserTitle != text.RegistrationUserTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, text.RegistrationUserTitle, text.Language))
	}
	if loginText.RegistrationUserDescription != text.RegistrationUserDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, text.RegistrationUserDescription, text.Language))
	}
	if loginText.RegistrationUserFirstnameLabel != text.RegistrationUserFirstnameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, text.RegistrationUserFirstnameLabel, text.Language))
	}
	if loginText.RegistrationUserLastnameLabel != text.RegistrationUserLastnameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, text.RegistrationUserLastnameLabel, text.Language))
	}
	if loginText.RegistrationUserEmailLabel != text.RegistrationUserEmailLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, text.RegistrationUserEmailLabel, text.Language))
	}
	if loginText.RegistrationUserLanguageLabel != text.RegistrationUserLanguageLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, text.RegistrationUserLanguageLabel, text.Language))
	}
	if loginText.RegistrationUserGenderLabel != text.RegistrationUserGenderLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, text.RegistrationUserGenderLabel, text.Language))
	}
	if loginText.RegistrationUserPasswordLabel != text.RegistrationUserPasswordLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, text.RegistrationUserPasswordLabel, text.Language))
	}
	if loginText.RegistrationUserPasswordConfirmLabel != text.RegistrationUserPasswordConfirmLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, text.RegistrationUserPasswordConfirmLabel, text.Language))
	}
	if loginText.RegistrationUserSaveButtonText != text.RegistrationUserSaveButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserSaveButtonText, text.RegistrationUserSaveButtonText, text.Language))
	}
	if loginText.RegisterOrgTitle != text.RegisterOrgTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, text.RegisterOrgTitle, text.Language))
	}
	if loginText.RegisterOrgDescription != text.RegisterOrgDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, text.RegisterOrgDescription, text.Language))
	}
	if loginText.RegisterOrgOrgNameLabel != text.RegisterOrgOrgNameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, text.RegisterOrgOrgNameLabel, text.Language))
	}
	if loginText.RegisterOrgFirstnameLabel != text.RegisterOrgFirstnameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, text.RegisterOrgFirstnameLabel, text.Language))
	}
	if loginText.RegisterOrgLastnameLabel != text.RegisterOrgLastnameLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, text.RegisterOrgLastnameLabel, text.Language))
	}
	if loginText.RegisterOrgEmailLabel != text.RegisterOrgEmailLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, text.RegisterOrgEmailLabel, text.Language))
	}
	if loginText.RegisterOrgPasswordLabel != text.RegisterOrgPasswordLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, text.RegisterOrgPasswordLabel, text.Language))
	}
	if loginText.RegisterOrgPasswordConfirmLabel != text.RegisterOrgPasswordConfirmLabel {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, text.RegisterOrgPasswordConfirmLabel, text.Language))
	}
	if loginText.RegisterOrgSaveButtonText != text.RegisterOrgSaveButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, text.RegisterOrgSaveButtonText, text.Language))
	}
	if loginText.PasswordlessTitle != text.PasswordlessTitle {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, text.PasswordlessTitle, text.Language))
	}
	if loginText.PasswordlessDescription != text.PasswordlessDescription {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, text.PasswordlessDescription, text.Language))
	}
	if loginText.PasswordlessLoginWithPwButtonText != text.PasswordlessLoginWithPwButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, text.PasswordlessLoginWithPwButtonText, text.Language))
	}
	if loginText.PasswordlessValidateTokenButtonText != text.PasswordlessValidateTokenButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, text.PasswordlessValidateTokenButtonText, text.Language))
	}
	return events, loginText, nil
}

func (c *Commands) defaultLoginTextWriteModelByID(ctx context.Context, lang language.Tag) (*IAMCustomLoginTextReadModel, error) {
	writeModel := NewIAMCustomLoginTextReadModel(lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
