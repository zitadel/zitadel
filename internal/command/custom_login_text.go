package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) createAllLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	events = append(events, c.createSelectLoginTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createLoginTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createUsernameChangeTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createUsernameChangeDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordInitTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordInitDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createEmailVerificationTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createEmailVerificationDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitUserDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitMFAPromptEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitMFAOTPEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitMFAU2FEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createInitMFADoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.crateMFAProviderEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createVerifyMFAOTPEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createVerifyMFAU2FEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordlessEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordlessRegistrationEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordlessRegistrationDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordChangeEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordChangeDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createPasswordResetDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createRegistrationOptionEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createRegistrationUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createRegistrationOrgEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createLinkingUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createExternalUserNotFoundEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createSuccessLoginEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createLogoutDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.createFooterTextEvents(ctx, agg, existingText, text, defaultText)...)
	return events
}

func (c *Commands) createSelectLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitle, existingText.SelectAccountTitle, text.SelectAccount.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescription, existingText.SelectAccountDescription, text.SelectAccount.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitleLinkingProcess, existingText.SelectAccountTitleLinkingProcess, text.SelectAccount.TitleLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescriptionLinkingProcess, existingText.SelectAccountDescriptionLinkingProcess, text.SelectAccount.DescriptionLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountOtherUser, existingText.SelectAccountOtherUser, text.SelectAccount.OtherUser, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateActive, existingText.SelectAccountSessionStateActive, text.SelectAccount.SessionState0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateInactive, existingText.SelectAccountSessionStateInactive, text.SelectAccount.SessionState1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, existingText.SelectAccountUserMustBeMemberOfOrg, text.SelectAccount.MustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitle, existingText.LoginTitle, text.Login.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescription, existingText.LoginDescription, text.Login.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitleLinkingProcess, existingText.LoginTitleLinkingProcess, text.Login.TitleLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescriptionLinkingProcess, existingText.LoginDescriptionLinkingProcess, text.Login.DescriptionLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNameLabel, existingText.LoginNameLabel, text.Login.LoginNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUsernamePlaceHolder, existingText.LoginUsernamePlaceholder, text.Login.UsernamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginLoginnamePlaceHolder, existingText.LoginLoginnamePlaceholder, text.Login.LoginnamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginRegisterButtonText, existingText.LoginRegisterButtonText, text.Login.RegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNextButtonText, existingText.LoginNextButtonText, text.Login.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginExternalUserDescription, existingText.LoginExternalUserDescription, text.Login.ExternalUserDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUserMustBeMemberOfOrg, existingText.LoginUserMustBeMemberOfOrg, text.Login.MustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordTitle, existingText.PasswordTitle, text.Password.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordDescription, existingText.PasswordDescription, text.Password.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordLabel, existingText.PasswordLabel, text.Password.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetLinkText, existingText.PasswordResetLinkText, text.Password.ResetLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordBackButtonText, existingText.PasswordBackButtonText, text.Password.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordNextButtonText, existingText.PasswordNextButtonText, text.Password.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordMinLength, existingText.PasswordMinLength, text.Password.MinLength, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasUppercase, existingText.PasswordHasUppercase, text.Password.HasUppercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasLowercase, existingText.PasswordHasLowercase, text.Password.HasLowercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasNumber, existingText.PasswordHasNumber, text.Password.HasNumber, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasSymbol, existingText.PasswordHasSymbol, text.Password.HasSymbol, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordConfirmation, existingText.PasswordConfirmation, text.Password.Confirmation, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createUsernameChangeTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeTitle, existingText.UsernameChangeTitle, text.UsernameChange.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDescription, existingText.UsernameChangeDescription, text.UsernameChange.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeUsernameLabel, existingText.UsernameChangeUsernameLabel, text.UsernameChange.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeCancelButtonText, existingText.UsernameChangeCancelButtonText, text.UsernameChange.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeNextButtonText, existingText.UsernameChangeNextButtonText, text.UsernameChange.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createUsernameChangeDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneTitle, existingText.UsernameChangeDoneTitle, text.UsernameChangeDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneDescription, existingText.UsernameChangeDoneDescription, text.UsernameChangeDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneNextButtonText, existingText.UsernameChangeDoneNextButtonText, text.UsernameChangeDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordInitTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordTitle, existingText.InitPasswordTitle, text.InitPassword.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDescription, existingText.InitPasswordDescription, text.InitPassword.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordCodeLabel, existingText.InitPasswordCodeLabel, text.InitPassword.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordLabel, existingText.InitPasswordNewPasswordLabel, text.InitPassword.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, existingText.InitPasswordNewPasswordConfirmLabel, text.InitPassword.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNextButtonText, existingText.InitPasswordNextButtonText, text.InitPassword.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordResendButtonText, existingText.InitPasswordResendButtonText, text.InitPassword.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordInitDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneTitle, existingText.InitPasswordDoneTitle, text.InitPasswordDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneDescription, existingText.InitPasswordDoneDescription, text.InitPasswordDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneNextButtonText, existingText.InitPasswordDoneNextButtonText, text.InitPasswordDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneCancelButtonText, existingText.InitPasswordDoneCancelButtonText, text.InitPasswordDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createEmailVerificationTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationTitle, existingText.EmailVerificationTitle, text.EmailVerification.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDescription, existingText.EmailVerificationDescription, text.EmailVerification.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationCodeLabel, existingText.EmailVerificationCodeLabel, text.EmailVerification.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationNextButtonText, existingText.EmailVerificationNextButtonText, text.EmailVerification.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationResendButtonText, existingText.EmailVerificationResendButtonText, text.EmailVerification.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createEmailVerificationDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneTitle, existingText.EmailVerificationDoneTitle, text.EmailVerificationDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneDescription, existingText.EmailVerificationDoneDescription, text.EmailVerificationDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneNextButtonText, existingText.EmailVerificationDoneNextButtonText, text.EmailVerificationDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneCancelButtonText, existingText.EmailVerificationDoneCancelButtonText, text.EmailVerificationDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneLoginButtonText, existingText.EmailVerificationDoneLoginButtonText, text.EmailVerificationDone.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserTitle, existingText.InitializeTitle, text.InitUser.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserDescription, existingText.InitializeDescription, text.InitUser.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserCodeLabel, existingText.InitializeCodeLabel, text.InitUser.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordLabel, existingText.InitializeNewPassword, text.InitUser.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, existingText.InitializeNewPasswordConfirm, text.InitUser.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserResendButtonText, existingText.InitializeResendButtonText, text.InitUser.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNextButtonText, existingText.InitializeNextButtonText, text.InitUser.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitUserDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneTitle, existingText.InitializeDoneTitle, text.InitUserDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneDescription, existingText.InitializeDoneDescription, text.InitUserDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneCancelButtonText, existingText.InitializeDoneAbortButtonText, text.InitUserDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneNextButtonText, existingText.InitializeDoneNextButtonText, text.InitUserDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitMFAPromptEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptTitle, existingText.InitMFAPromptTitle, text.InitMFAPrompt.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptDescription, existingText.InitMFAPromptDescription, text.InitMFAPrompt.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptOTPOption, existingText.InitMFAPromptOTPOption, text.InitMFAPrompt.Provider0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptU2FOption, existingText.InitMFAPromptU2FOption, text.InitMFAPrompt.Provider1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptSkipButtonText, existingText.InitMFAPromptSkipButtonText, text.InitMFAPrompt.SkipButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptNextButtonText, existingText.InitMFAPromptNextButtonText, text.InitMFAPrompt.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPTitle, existingText.InitMFAOTPTitle, text.InitMFAOTP.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescription, existingText.InitMFAOTPDescription, text.InitMFAOTP.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescriptionOTP, existingText.InitMFAOTPDescriptionOTP, text.InitMFAOTP.OTPDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPSecretLabel, existingText.InitMFAOTPSecretLabel, text.InitMFAOTP.SecretLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCodeLabel, existingText.InitMFAOTPCodeLabel, text.InitMFAOTP.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPNextButtonText, existingText.InitMFAOTPNextButtonText, text.InitMFAOTP.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCancelButtonText, existingText.InitMFAOTPCancelButtonText, text.InitMFAOTP.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTitle, existingText.InitMFAU2FTitle, text.InitMFAU2F.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FDescription, existingText.InitMFAU2FDescription, text.InitMFAU2F.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTokenNameLabel, existingText.InitMFAU2FTokenNameLabel, text.InitMFAU2F.TokenNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, existingText.InitMFAU2FRegisterTokenButtonText, text.InitMFAU2F.RegisterTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FNotSupported, existingText.InitMFAU2FNotSupported, text.InitMFAU2F.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FErrorRetry, existingText.InitMFAU2FErrorRetry, text.InitMFAU2F.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createInitMFADoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneTitle, existingText.InitMFADoneTitle, text.InitMFADone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneDescription, existingText.InitMFADoneDescription, text.InitMFADone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneCancelButtonText, existingText.InitMFADoneAbortButtonText, text.InitMFADone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneNextButtonText, existingText.InitMFADoneNextButtonText, text.InitMFADone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) crateMFAProviderEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersChooseOther, existingText.MFAProvidersChooseOther, text.MFAProvider.ChooseOther, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersOTP, existingText.MFAProvidersOTP, text.MFAProvider.Provider0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersU2F, existingText.MFAProvidersU2F, text.MFAProvider.Provider1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createVerifyMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPTitle, existingText.VerifyMFAOTPTitle, text.VerifyMFAOTP.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPDescription, existingText.VerifyMFAOTPDescription, text.VerifyMFAOTP.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPCodeLabel, existingText.VerifyMFAOTPCodeLabel, text.VerifyMFAOTP.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPNextButtonText, existingText.VerifyMFAOTPNextButtonText, text.VerifyMFAOTP.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createVerifyMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FTitle, existingText.VerifyMFAU2FTitle, text.VerifyMFAU2F.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FDescription, existingText.VerifyMFAU2FDescription, text.VerifyMFAU2F.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FValidateTokenText, existingText.VerifyMFAU2FValidateTokenText, text.VerifyMFAU2F.ValidateTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FNotSupported, existingText.VerifyMFAU2FNotSupported, text.VerifyMFAU2F.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FErrorRetry, existingText.VerifyMFAU2FErrorRetry, text.VerifyMFAU2F.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordlessEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessTitle, existingText.PasswordlessTitle, text.Passwordless.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessDescription, existingText.PasswordlessDescription, text.Passwordless.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessLoginWithPwButtonText, existingText.PasswordlessLoginWithPwButtonText, text.Passwordless.LoginWithPwButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessValidateTokenButtonText, existingText.PasswordlessValidateTokenButtonText, text.Passwordless.ValidateTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessNotSupported, existingText.PasswordlessNotSupported, text.Passwordless.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessErrorRetry, existingText.PasswordlessErrorRetry, text.Passwordless.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordlessRegistrationEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationTitle, existingText.PasswordlessRegistrationTitle, text.PasswordlessRegistration.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationDescription, existingText.PasswordlessRegistrationDescription, text.PasswordlessRegistration.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationRegisterTokenButtonText, existingText.PasswordlessRegistrationRegisterTokenButtonText, text.PasswordlessRegistration.RegisterTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationTokenNameLabel, existingText.PasswordlessRegistrationTokenNameLabel, text.PasswordlessRegistration.TokenNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationNotSupported, existingText.PasswordlessRegistrationNotSupported, text.PasswordlessRegistration.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationErrorRetry, existingText.PasswordlessRegistrationErrorRetry, text.PasswordlessRegistration.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordlessRegistrationDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationDoneTitle, existingText.PasswordlessRegistrationDoneTitle, text.PasswordlessRegistrationDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationDoneDescription, existingText.PasswordlessRegistrationDoneDescription, text.PasswordlessRegistrationDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessRegistrationDoneNextButtonText, existingText.PasswordlessRegistrationDoneNextButtonText, text.PasswordlessRegistrationDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordChangeEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeTitle, existingText.PasswordChangeTitle, text.PasswordChange.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDescription, existingText.PasswordChangeDescription, text.PasswordChange.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeOldPasswordLabel, existingText.PasswordChangeOldPasswordLabel, text.PasswordChange.OldPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordLabel, existingText.PasswordChangeNewPasswordLabel, text.PasswordChange.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, existingText.PasswordChangeNewPasswordConfirmLabel, text.PasswordChange.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeCancelButtonText, existingText.PasswordChangeCancelButtonText, text.PasswordChange.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNextButtonText, existingText.PasswordChangeNextButtonText, text.PasswordChange.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordChangeDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneTitle, existingText.PasswordChangeDoneTitle, text.PasswordChangeDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneDescription, existingText.PasswordChangeDoneDescription, text.PasswordChangeDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneNextButtonText, existingText.PasswordChangeDoneNextButtonText, text.PasswordChangeDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createPasswordResetDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneTitle, existingText.PasswordResetDoneTitle, text.PasswordResetDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneDescription, existingText.PasswordResetDoneDescription, text.PasswordResetDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneNextButtonText, existingText.PasswordResetDoneNextButtonText, text.PasswordResetDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createRegistrationOptionEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionTitle, existingText.RegistrationOptionTitle, text.RegisterOption.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionDescription, existingText.RegistrationOptionDescription, text.RegisterOption.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionUserNameButtonText, existingText.RegistrationOptionUserNameButtonText, text.RegisterOption.RegisterUsernamePasswordButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionExternalLoginDescription, existingText.RegistrationOptionExternalLoginDescription, text.RegisterOption.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createRegistrationUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTitle, existingText.RegistrationUserTitle, text.RegistrationUser.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserDescription, existingText.RegistrationUserDescription, text.RegistrationUser.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserDescriptionOrgRegister, existingText.RegistrationUserDescriptionOrgRegister, text.RegistrationUser.DescriptionOrgRegister, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserFirstnameLabel, existingText.RegistrationUserFirstnameLabel, text.RegistrationUser.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLastnameLabel, existingText.RegistrationUserLastnameLabel, text.RegistrationUser.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserEmailLabel, existingText.RegistrationUserEmailLabel, text.RegistrationUser.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserUsernameLabel, existingText.RegistrationUserUsernameLabel, text.RegistrationUser.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLanguageLabel, existingText.RegistrationUserLanguageLabel, text.RegistrationUser.LanguageLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserGenderLabel, existingText.RegistrationUserGenderLabel, text.RegistrationUser.GenderLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordLabel, existingText.RegistrationUserPasswordLabel, text.RegistrationUser.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordConfirmLabel, existingText.RegistrationUserPasswordConfirmLabel, text.RegistrationUser.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, existingText.RegistrationUserTOSAndPrivacyLabel, text.RegistrationUser.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSConfirm, existingText.RegistrationUserTOSConfirm, text.RegistrationUser.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSLinkText, existingText.RegistrationUserTOSLinkText, text.RegistrationUser.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSConfirmAnd, existingText.RegistrationUserTOSConfirmAnd, text.RegistrationUser.TOSConfirmAnd, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyLinkText, existingText.RegistrationUserPrivacyLinkText, text.RegistrationUser.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserNextButtonText, existingText.RegistrationUserNextButtonText, text.RegistrationUser.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserBackButtonText, existingText.RegistrationUserBackButtonText, text.RegistrationUser.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createRegistrationOrgEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTitle, existingText.RegisterOrgTitle, text.RegistrationOrg.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgDescription, existingText.RegisterOrgDescription, text.RegistrationOrg.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgOrgNameLabel, existingText.RegisterOrgOrgNameLabel, text.RegistrationOrg.OrgNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgFirstnameLabel, existingText.RegisterOrgFirstnameLabel, text.RegistrationOrg.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgLastnameLabel, existingText.RegisterOrgLastnameLabel, text.RegistrationOrg.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgUsernameLabel, existingText.RegisterOrgUsernameLabel, text.RegistrationOrg.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgEmailLabel, existingText.RegisterOrgEmailLabel, text.RegistrationOrg.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordLabel, existingText.RegisterOrgPasswordLabel, text.RegistrationOrg.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordConfirmLabel, existingText.RegisterOrgPasswordConfirmLabel, text.RegistrationOrg.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, existingText.RegisterOrgTOSAndPrivacyLabel, text.RegistrationOrg.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSConfirm, existingText.RegisterOrgTOSConfirm, text.RegistrationOrg.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSLinkText, existingText.RegisterOrgTOSLinkText, text.RegistrationOrg.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTosConfirmAnd, existingText.RegisterOrgTOSConfirmAnd, text.RegistrationOrg.TOSConfirmAnd, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyLinkText, existingText.RegisterOrgPrivacyLinkText, text.RegistrationOrg.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgSaveButtonText, existingText.RegisterOrgSaveButtonText, text.RegistrationOrg.SaveButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createLinkingUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneTitle, existingText.LinkingUserDoneTitle, text.LinkingUsersDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneDescription, existingText.LinkingUserDoneDescription, text.LinkingUsersDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneCancelButtonText, existingText.LinkingUserDoneCancelButtonText, text.LinkingUsersDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneNextButtonText, existingText.LinkingUserDoneNextButtonText, text.LinkingUsersDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createExternalUserNotFoundEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundTitle, existingText.ExternalUserNotFoundTitle, text.ExternalNotFoundOption.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundDescription, existingText.ExternalUserNotFoundDescription, text.ExternalNotFoundOption.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundLinkButtonText, existingText.ExternalUserNotFoundLinkButtonText, text.ExternalNotFoundOption.LinkButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, existingText.ExternalUserNotFoundAutoRegisterButtonText, text.ExternalNotFoundOption.AutoRegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createSuccessLoginEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginTitle, existingText.SuccessLoginTitle, text.LoginSuccess.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginAutoRedirectDescription, existingText.SuccessLoginAutoRedirectDescription, text.LoginSuccess.AutoRedirectDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginRedirectedDescription, existingText.SuccessLoginRedirectedDescription, text.LoginSuccess.RedirectedDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginNextButtonText, existingText.SuccessLoginNextButtonText, text.LoginSuccess.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createLogoutDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneTitle, existingText.LogoutDoneTitle, text.LogoutDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneDescription, existingText.LogoutDoneDescription, text.LogoutDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneLoginButtonText, existingText.LogoutDoneLoginButtonText, text.LogoutDone.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createFooterTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterTOS, existingText.FooterTOS, text.Footer.TOS, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterPrivacyPolicy, existingText.FooterPrivacyPolicy, text.Footer.PrivacyPolicy, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelp, existingText.FooterHelp, text.Footer.Help, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.createCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelpLink, existingText.FooterHelpLink, text.Footer.HelpLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) createCustomLoginTextEvent(ctx context.Context, agg *eventstore.Aggregate, textKey, existingText, newText string, lang language.Tag, defaultText bool) eventstore.EventPusher {
	if existingText == newText {
		return nil
	}
	if defaultText {
		if newText != "" {
			return iam.NewCustomTextSetEvent(ctx, agg, domain.LoginCustomText, textKey, newText, lang)
		}
		return iam.NewCustomTextRemovedEvent(ctx, agg, domain.LoginCustomText, textKey, lang)
	}
	if newText != "" {
		return org.NewCustomTextSetEvent(ctx, agg, domain.LoginCustomText, textKey, newText, lang)
	}
	return org.NewCustomTextRemovedEvent(ctx, agg, domain.LoginCustomText, textKey, lang)
}
