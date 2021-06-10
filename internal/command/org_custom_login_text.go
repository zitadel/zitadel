package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) SetOrgLoginText(ctx context.Context, resourceOwner string, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-m29rF", "Errors.ResourceOwnerMissing")
	}
	iamAgg := org.NewAggregate(resourceOwner, resourceOwner)
	events, existingLoginText, err := c.setOrgLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingLoginText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingLoginText.WriteModel), nil
}

func (c *Commands) setOrgLoginText(ctx context.Context, orgAgg *eventstore.Aggregate, loginText *domain.CustomLoginText) ([]eventstore.EventPusher, *OrgCustomLoginTextReadModel, error) {
	if !loginText.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "ORG-PPo2w", "Errors.CustomText.Invalid")
	}

	existingLoginText, err := c.orgCustomLoginTextWriteModelByID(ctx, orgAgg.ID, loginText.Language)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	if existingLoginText.SelectAccountTitle != loginText.SelectAccountTitle {
		if loginText.SelectAccountTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountTitle, loginText.SelectAccountTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountTitle, loginText.Language))
		}
	}
	if existingLoginText.SelectAccountDescription != loginText.SelectAccountDescription {
		if loginText.SelectAccountDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, loginText.SelectAccountDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, loginText.Language))
		}
	}
	if existingLoginText.SelectAccountOtherUser != loginText.SelectAccountOtherUser {
		if loginText.SelectAccountOtherUser != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, loginText.SelectAccountOtherUser, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, loginText.Language))
		}
	}
	if existingLoginText.LoginTitle != loginText.LoginTitle {
		if loginText.LoginTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginTitle, loginText.LoginTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginTitle, loginText.Language))
		}
	}
	if existingLoginText.LoginDescription != loginText.LoginDescription {
		if loginText.LoginDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginDescription, loginText.LoginDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginDescription, loginText.Language))
		}
	}
	if existingLoginText.LoginNameLabel != loginText.LoginNameLabel {
		if loginText.LoginNameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, loginText.LoginNameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, loginText.Language))
		}
	}
	if existingLoginText.LoginRegisterButtonText != loginText.LoginRegisterButtonText {
		if loginText.LoginRegisterButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, loginText.LoginRegisterButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, loginText.Language))
		}
	}
	if existingLoginText.LoginNextButtonText != loginText.LoginNextButtonText {
		if loginText.LoginNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, loginText.LoginNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.LoginExternalUserDescription != loginText.LoginExternalUserDescription {
		if loginText.LoginExternalUserDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, loginText.LoginExternalUserDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, loginText.Language))
		}
	}
	if existingLoginText.PasswordTitle != loginText.PasswordTitle {
		if loginText.PasswordTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordTitle, loginText.PasswordTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordTitle, loginText.Language))
		}
	}
	if existingLoginText.PasswordDescription != loginText.PasswordDescription {
		if loginText.PasswordDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordDescription, loginText.PasswordDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordDescription, loginText.Language))
		}
	}
	if existingLoginText.PasswordLabel != loginText.PasswordLabel {
		if loginText.PasswordLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordLabel, loginText.PasswordLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordLabel, loginText.Language))
		}
	}
	if existingLoginText.PasswordResetLinkText != loginText.PasswordResetLinkText {
		if loginText.PasswordResetLinkText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, loginText.PasswordResetLinkText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, loginText.Language))
		}
	}
	if existingLoginText.PasswordBackButtonText != loginText.PasswordBackButtonText {
		if loginText.PasswordBackButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, loginText.PasswordBackButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, loginText.Language))
		}
	}
	if existingLoginText.PasswordNextButtonText != loginText.PasswordNextButtonText {
		if loginText.PasswordNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, loginText.PasswordNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.PasswordResetTitle != loginText.ResetPasswordTitle {
		if loginText.ResetPasswordTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordTitle, loginText.ResetPasswordTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordTitle, loginText.Language))
		}
	}
	if existingLoginText.PasswordResetDescription != loginText.ResetPasswordDescription {
		if loginText.ResetPasswordDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordDescription, loginText.ResetPasswordDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordDescription, loginText.Language))
		}
	}
	if existingLoginText.PasswordResetNextButtonText != loginText.ResetPasswordNextButtonText {
		if loginText.ResetPasswordNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordNextButtonText, loginText.ResetPasswordNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyResetPasswordNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitializeTitle != loginText.InitializeUserTitle {
		if loginText.InitializeUserTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, loginText.InitializeUserTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, loginText.Language))
		}
	}
	if existingLoginText.InitializeDescription != loginText.InitializeUserDescription {
		if loginText.InitializeUserDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, loginText.InitializeUserDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, loginText.Language))
		}
	}
	if existingLoginText.InitializeCodeLabel != loginText.InitializeUserCodeLabel {
		if loginText.InitializeUserCodeLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, loginText.InitializeUserCodeLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, loginText.Language))
		}
	}
	if existingLoginText.InitializeNewPassword != loginText.InitializeUserNewPasswordLabel {
		if loginText.InitializeUserNewPasswordLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPassword, loginText.InitializeUserNewPasswordLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPassword, loginText.Language))
		}
	}
	if existingLoginText.InitializeNewPasswordConfirm != loginText.InitializeUserNewPasswordConfirmLabel {
		if loginText.InitializeUserNewPasswordConfirmLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirm, loginText.InitializeUserNewPasswordConfirmLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirm, loginText.Language))
		}
	}
	if existingLoginText.InitializeResendButtonText != loginText.InitializeUserResendButtonText {
		if loginText.InitializeUserResendButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, loginText.InitializeUserResendButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitializeNextButtonText != loginText.InitializeUserNextButtonText {
		if loginText.InitializeUserNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, loginText.InitializeUserNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitializeDoneTitle != loginText.InitializeDoneTitle {
		if loginText.InitializeDoneTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneTitle, loginText.InitializeDoneTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneTitle, loginText.Language))
		}
	}
	if existingLoginText.InitializeDoneDescription != loginText.InitializeDoneDescription {
		if loginText.InitializeDoneDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneDescription, loginText.InitializeDoneDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneDescription, loginText.Language))
		}
	}
	if existingLoginText.InitializeDoneAbortButtonText != loginText.InitializeDoneAbortButtonText {
		if loginText.InitializeDoneAbortButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneAbortButtonText, loginText.InitializeDoneAbortButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneAbortButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitializeDoneNextButtonText != loginText.InitializeDoneNextButtonText {
		if loginText.InitializeDoneNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneNextButtonText, loginText.InitializeDoneNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitializeDoneNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitMFATitle != loginText.InitMFAPromptTitle {
		if loginText.InitMFAPromptTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, loginText.InitMFAPromptTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, loginText.Language))
		}
	}
	if existingLoginText.InitMFADescription != loginText.InitMFAPromptDescription {
		if loginText.InitMFAPromptDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, loginText.InitMFAPromptDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, loginText.Language))
		}
	}
	if existingLoginText.InitMFAOTPOption != loginText.InitMFAPromptOTPOption {
		if loginText.InitMFAPromptOTPOption != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, loginText.InitMFAPromptOTPOption, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, loginText.Language))
		}
	}
	if existingLoginText.InitMFAU2FOption != loginText.InitMFAPromptU2FOption {
		if loginText.InitMFAPromptU2FOption != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, loginText.InitMFAPromptU2FOption, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, loginText.Language))
		}
	}
	if existingLoginText.InitMFASkipButtonText != loginText.InitMFAPromptSkipButtonText {
		if loginText.InitMFAPromptSkipButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, loginText.InitMFAPromptSkipButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitMFANextButtonText != loginText.InitMFAPromptNextButtonText {
		if loginText.InitMFAPromptNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, loginText.InitMFAPromptNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitMFAOTPTitle != loginText.InitMFAOTPTitle {
		if loginText.InitMFAOTPTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, loginText.InitMFAOTPTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, loginText.Language))
		}
	}
	if existingLoginText.InitMFAOTPDescription != loginText.InitMFAOTPDescription {
		if loginText.InitMFAOTPDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, loginText.InitMFAOTPDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, loginText.Language))
		}
	}
	if existingLoginText.InitMFAOTPSecretLabel != loginText.InitMFAOTPSecretLabel {
		if loginText.InitMFAOTPSecretLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, loginText.InitMFAOTPSecretLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, loginText.Language))
		}
	}
	if existingLoginText.InitMFAOTPCodeLabel != loginText.InitMFAOTPCodeLabel {
		if loginText.InitMFAOTPCodeLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, loginText.InitMFAOTPCodeLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, loginText.Language))
		}
	}
	if existingLoginText.InitMFAU2FTitle != loginText.InitMFAU2FTitle {
		if loginText.InitMFAU2FTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, loginText.InitMFAU2FTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, loginText.Language))
		}
	}
	if existingLoginText.InitMFAU2FDescription != loginText.InitMFAU2FDescription {
		if loginText.InitMFAU2FDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, loginText.InitMFAU2FDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, loginText.Language))
		}
	}
	if existingLoginText.InitMFAU2FTokenNameLabel != loginText.InitMFAU2FTokenNameLabel {
		if loginText.InitMFAU2FTokenNameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, loginText.InitMFAU2FTokenNameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, loginText.Language))
		}
	}
	if existingLoginText.InitMFAU2FRegisterTokenButtonText != loginText.InitMFAU2FRegisterTokenButtonText {
		if loginText.InitMFAU2FRegisterTokenButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, loginText.InitMFAU2FRegisterTokenButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitMFADoneTitle != loginText.InitMFADoneTitle {
		if loginText.InitMFADoneTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, loginText.InitMFADoneTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, loginText.Language))
		}
	}
	if existingLoginText.InitMFADoneDescription != loginText.InitMFADoneDescription {
		if loginText.InitMFADoneDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, loginText.InitMFADoneDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, loginText.Language))
		}
	}
	if existingLoginText.InitMFADoneAbortButtonText != loginText.InitMFADoneAbortButtonText {
		if loginText.InitMFADoneAbortButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneAbortButtonText, loginText.InitMFADoneAbortButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneAbortButtonText, loginText.Language))
		}
	}
	if existingLoginText.InitMFADoneNextButtonText != loginText.InitMFADoneNextButtonText {
		if loginText.InitMFADoneNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, loginText.InitMFADoneNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAOTPTitle != loginText.VerifyMFAOTPTitle {
		if loginText.VerifyMFAOTPTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, loginText.VerifyMFAOTPTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAOTPDescription != loginText.VerifyMFAOTPDescription {
		if loginText.VerifyMFAOTPDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, loginText.VerifyMFAOTPDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAOTPCodeLabel != loginText.VerifyMFAOTPCodeLabel {
		if loginText.VerifyMFAOTPCodeLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, loginText.VerifyMFAOTPCodeLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAOTPNextButtonText != loginText.VerifyMFAOTPNextButtonText {
		if loginText.VerifyMFAOTPNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, loginText.VerifyMFAOTPNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAU2FTitle != loginText.VerifyMFAU2FTitle {
		if loginText.VerifyMFAU2FTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, loginText.VerifyMFAU2FTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAU2FDescription != loginText.VerifyMFAU2FDescription {
		if loginText.VerifyMFAU2FDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, loginText.VerifyMFAU2FDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAU2FValidateTokenText != loginText.VerifyMFAU2FValidateTokenText {
		if loginText.VerifyMFAU2FValidateTokenText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, loginText.VerifyMFAU2FValidateTokenText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAU2FOtherOptionDescription != loginText.VerifyMFAU2FOtherOptionDescription {
		if loginText.VerifyMFAU2FOtherOptionDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOtherOptionDescription, loginText.VerifyMFAU2FOtherOptionDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOtherOptionDescription, loginText.Language))
		}
	}
	if existingLoginText.VerifyMFAU2FOptionOTPButtonText != loginText.VerifyMFAU2FOptionOTPButtonText {
		if loginText.VerifyMFAU2FOptionOTPButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOptionOTPButtonText, loginText.VerifyMFAU2FOptionOTPButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FOptionOTPButtonText, loginText.Language))
		}
	}
	if existingLoginText.RegistrationOptionTitle != loginText.RegistrationOptionTitle {
		if loginText.RegistrationOptionTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, loginText.RegistrationOptionTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, loginText.Language))
		}
	}
	if existingLoginText.RegistrationOptionDescription != loginText.RegistrationOptionDescription {
		if loginText.RegistrationOptionDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, loginText.RegistrationOptionDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, loginText.Language))
		}
	}
	if existingLoginText.RegistrationOptionUserNameButtonText != loginText.RegistrationOptionUserNameButtonText {
		if loginText.RegistrationOptionUserNameButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, loginText.RegistrationOptionUserNameButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, loginText.Language))
		}
	}
	if existingLoginText.RegistrationOptionExternalLoginDescription != loginText.RegistrationOptionExternalLoginDescription {
		if loginText.RegistrationOptionExternalLoginDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, loginText.RegistrationOptionExternalLoginDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserTitle != loginText.RegistrationUserTitle {
		if loginText.RegistrationUserTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, loginText.RegistrationUserTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserDescription != loginText.RegistrationUserDescription {
		if loginText.RegistrationUserDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, loginText.RegistrationUserDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserFirstnameLabel != loginText.RegistrationUserFirstnameLabel {
		if loginText.RegistrationUserFirstnameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, loginText.RegistrationUserFirstnameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserLastnameLabel != loginText.RegistrationUserLastnameLabel {
		if loginText.RegistrationUserLastnameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, loginText.RegistrationUserLastnameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserEmailLabel != loginText.RegistrationUserEmailLabel {
		if loginText.RegistrationUserEmailLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, loginText.RegistrationUserEmailLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserLanguageLabel != loginText.RegistrationUserLanguageLabel {
		if loginText.RegistrationUserLanguageLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, loginText.RegistrationUserLanguageLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserGenderLabel != loginText.RegistrationUserGenderLabel {
		if loginText.RegistrationUserGenderLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, loginText.RegistrationUserGenderLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserPasswordLabel != loginText.RegistrationUserPasswordLabel {
		if loginText.RegistrationUserPasswordLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, loginText.RegistrationUserPasswordLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserPasswordConfirmLabel != loginText.RegistrationUserPasswordConfirmLabel {
		if loginText.RegistrationUserPasswordConfirmLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, loginText.RegistrationUserPasswordConfirmLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, loginText.Language))
		}
	}
	if existingLoginText.RegistrationUserSaveButtonText != loginText.RegistrationUserSaveButtonText {
		if loginText.RegistrationUserSaveButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserSaveButtonText, loginText.RegistrationUserSaveButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegistrationUserSaveButtonText, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgTitle != loginText.RegisterOrgTitle {
		if loginText.RegisterOrgTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, loginText.RegisterOrgTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgDescription != loginText.RegisterOrgDescription {
		if loginText.RegisterOrgDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, loginText.RegisterOrgDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgOrgNameLabel != loginText.RegisterOrgOrgNameLabel {
		if loginText.RegisterOrgOrgNameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, loginText.RegisterOrgOrgNameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgFirstnameLabel != loginText.RegisterOrgFirstnameLabel {
		if loginText.RegisterOrgFirstnameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, loginText.RegisterOrgFirstnameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgLastnameLabel != loginText.RegisterOrgLastnameLabel {
		if loginText.RegisterOrgLastnameLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, loginText.RegisterOrgLastnameLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgEmailLabel != loginText.RegisterOrgEmailLabel {
		if loginText.RegisterOrgEmailLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, loginText.RegisterOrgEmailLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgPasswordLabel != loginText.RegisterOrgPasswordLabel {
		if loginText.RegisterOrgPasswordLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, loginText.RegisterOrgPasswordLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgPasswordConfirmLabel != loginText.RegisterOrgPasswordConfirmLabel {
		if loginText.RegisterOrgPasswordConfirmLabel != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, loginText.RegisterOrgPasswordConfirmLabel, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, loginText.Language))
		}
	}
	if existingLoginText.RegisterOrgSaveButtonText != loginText.RegisterOrgSaveButtonText {
		if loginText.RegisterOrgSaveButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, loginText.RegisterOrgSaveButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, loginText.Language))
		}
	}
	if existingLoginText.PasswordlessTitle != loginText.PasswordlessTitle {
		if loginText.PasswordlessTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, loginText.PasswordlessTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, loginText.Language))
		}
	}
	if existingLoginText.PasswordlessDescription != loginText.PasswordlessDescription {
		if loginText.PasswordlessDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, loginText.PasswordlessDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, loginText.Language))
		}
	}
	if existingLoginText.PasswordlessLoginWithPwButtonText != loginText.PasswordlessLoginWithPwButtonText {
		if loginText.PasswordlessLoginWithPwButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, loginText.PasswordlessLoginWithPwButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, loginText.Language))
		}
	}
	if existingLoginText.PasswordlessValidateTokenButtonText != loginText.PasswordlessValidateTokenButtonText {
		if loginText.PasswordlessValidateTokenButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, loginText.PasswordlessValidateTokenButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, loginText.Language))
		}
	}
	if existingLoginText.SuccessLoginTitle != loginText.SuccessLoginTitle {
		if loginText.SuccessLoginTitle != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginTitle, loginText.SuccessLoginTitle, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginTitle, loginText.Language))
		}
	}
	if existingLoginText.SuccessLoginAutoRedirectDescription != loginText.SuccessLoginAutoRedirectDescription {
		if loginText.SuccessLoginAutoRedirectDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginAutoRedirectDescription, loginText.SuccessLoginAutoRedirectDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginAutoRedirectDescription, loginText.Language))
		}
	}
	if existingLoginText.SuccessLoginRedirectedDescription != loginText.SuccessLoginRedirectedDescription {
		if loginText.SuccessLoginRedirectedDescription != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginRedirectedDescription, loginText.SuccessLoginRedirectedDescription, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginRedirectedDescription, loginText.Language))
		}
	}
	if existingLoginText.SuccessLoginNextButtonText != loginText.SuccessLoginNextButtonText {
		if loginText.SuccessLoginNextButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginNextButtonText, loginText.SuccessLoginNextButtonText, loginText.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, domain.LoginCustomText, domain.LoginKeySuccessLoginNextButtonText, loginText.Language))
		}
	}
	return events, existingLoginText, nil
}

func (c *Commands) RemoveOrgLoginTexts(ctx context.Context, resourceOwner string, lang language.Tag) error {
	if resourceOwner == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-1B8dw", "Errors.ResourceOwnerMissing")
	}
	if lang == language.Und {
		return caos_errs.ThrowInvalidArgument(nil, "Org-5ZZmo", "Errors.CustomMailText.Invalid")
	}
	customText, err := c.orgCustomLoginTextWriteModelByID(ctx, resourceOwner, lang)
	if err != nil {
		return err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-9ru44", "Errors.CustomMailText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&customText.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, domain.LoginCustomText, lang))
	return err
}

func (c *Commands) orgCustomLoginTextWriteModelByID(ctx context.Context, orgID string, lang language.Tag) (*OrgCustomLoginTextReadModel, error) {
	writeModel := NewOrgCustomLoginTextReadModel(orgID, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
