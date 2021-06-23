package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) getAllLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	events = append(events, c.getSelectLoginTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getLoginTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getUsernameChangeTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getUsernameChangeDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordInitTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordInitDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getEmailVerificationTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getEmailVerificationDoneTextEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitUserDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitMFAPromptEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitMFAOTPEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitMFAU2FEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getInitMFADoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getMFAProviderEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getVerifyMFAOTPEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getVerifyMFAU2FEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordlessEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordChangeEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordChangeDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getPasswordResetDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getRegistrationOptionEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getRegistrationUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getRegistrationOrgEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getLinkingUserEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getExternalUserNotFoundEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getSuccessLoginEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getLogoutDoneEvents(ctx, agg, existingText, text, defaultText)...)
	events = append(events, c.getFooterTextEvents(ctx, agg, existingText, text, defaultText)...)
	return events
}

func (c *Commands) getSelectLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitle, existingText.SelectAccountTitle, text.SelectAccount.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescription, existingText.SelectAccountDescription, text.SelectAccount.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitleLinkingProcess, existingText.SelectAccountTitleLinkingProcess, text.SelectAccount.TitleLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescriptionLinkingProcess, existingText.SelectAccountDescriptionLinkingProcess, text.SelectAccount.DescriptionLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountOtherUser, existingText.SelectAccountOtherUser, text.SelectAccount.OtherUser, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateActive, existingText.SelectAccountSessionStateActive, text.SelectAccount.SessionState0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateInactive, existingText.SelectAccountSessionStateInactive, text.SelectAccount.SessionState1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, existingText.SelectAccountUserMustBeMemberOfOrg, text.SelectAccount.MustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitle, existingText.LoginTitle, text.Login.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescription, existingText.LoginDescription, text.Login.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitleLinkingProcess, existingText.LoginTitleLinkingProcess, text.Login.TitleLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescriptionLinkingProcess, existingText.LoginDescriptionLinkingProcess, text.Login.DescriptionLinking, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNameLabel, existingText.LoginNameLabel, text.Login.LoginNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUsernamePlaceHolder, existingText.LoginUsernamePlaceholder, text.Login.UsernamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginLoginnamePlaceHolder, existingText.LoginLoginnamePlaceholder, text.Login.UsernamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginRegisterButtonText, existingText.LoginRegisterButtonText, text.Login.RegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNextButtonText, existingText.LoginNextButtonText, text.Login.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginExternalUserDescription, existingText.LoginExternalUserDescription, text.Login.ExternalLogin, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUserMustBeMemberOfOrg, existingText.LoginUserMustBeMemberOfOrg, text.Login.MustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordTitle, existingText.PasswordTitle, text.Password.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordDescription, existingText.PasswordDescription, text.Password.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordLabel, existingText.PasswordLabel, text.Password.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetLinkText, existingText.PasswordResetLinkText, text.Password.ResetLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordBackButtonText, existingText.PasswordBackButtonText, text.Password.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordNextButtonText, existingText.PasswordNextButtonText, text.Password.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordMinLength, existingText.PasswordMinLength, text.Password.MinLength, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasUppercase, existingText.PasswordHasUppercase, text.Password.HasUppercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasLowercase, existingText.PasswordHasLowercase, text.Password.HasLowercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasNumber, existingText.PasswordHasNumber, text.Password.HasNumber, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasSymbol, existingText.PasswordHasSymbol, text.Password.HasSymbol, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordConfirmation, existingText.PasswordConfirmation, text.Password.Confirmation, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getUsernameChangeTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeTitle, existingText.UsernameChangeTitle, text.UsernameChange.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDescription, existingText.UsernameChangeDescription, text.UsernameChange.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeUsernameLabel, existingText.UsernameChangeUsernameLabel, text.UsernameChange.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeCancelButtonText, existingText.UsernameChangeCancelButtonText, text.UsernameChange.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeNextButtonText, existingText.UsernameChangeNextButtonText, text.UsernameChange.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getUsernameChangeDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneTitle, existingText.UsernameChangeDoneTitle, text.UsernameChangeDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneDescription, existingText.UsernameChangeDoneDescription, text.UsernameChangeDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneNextButtonText, existingText.UsernameChangeDoneNextButtonText, text.UsernameChangeDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordInitTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordTitle, existingText.InitPasswordTitle, text.InitPassword.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDescription, existingText.InitPasswordDescription, text.InitPassword.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordCodeLabel, existingText.InitPasswordCodeLabel, text.InitPassword.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordLabel, existingText.InitPasswordNewPasswordLabel, text.InitPassword.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, existingText.InitPasswordNewPasswordConfirmLabel, text.InitPassword.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNextButtonText, existingText.InitPasswordNextButtonText, text.InitPassword.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordResendButtonText, existingText.InitPasswordResendButtonText, text.InitPassword.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordInitDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneTitle, existingText.InitPasswordDoneTitle, text.InitPasswordDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneDescription, existingText.InitPasswordDoneDescription, text.InitPasswordDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneNextButtonText, existingText.InitPasswordDoneNextButtonText, text.InitPasswordDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneCancelButtonText, existingText.InitPasswordDoneCancelButtonText, text.InitPasswordDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getEmailVerificationTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationTitle, existingText.EmailVerificationTitle, text.EmailVerification.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDescription, existingText.EmailVerificationDescription, text.EmailVerification.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationCodeLabel, existingText.EmailVerificationCodeLabel, text.EmailVerification.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationNextButtonText, existingText.EmailVerificationNextButtonText, text.EmailVerification.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationResendButtonText, existingText.EmailVerificationResendButtonText, text.EmailVerification.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getEmailVerificationDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneTitle, existingText.EmailVerificationDoneTitle, text.EmailVerificationDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneDescription, existingText.EmailVerificationDoneDescription, text.EmailVerificationDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneNextButtonText, existingText.EmailVerificationDoneNextButtonText, text.EmailVerificationDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneCancelButtonText, existingText.EmailVerificationDoneCancelButtonText, text.EmailVerificationDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneLoginButtonText, existingText.EmailVerificationDoneLoginButtonText, text.EmailVerificationDone.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserTitle, existingText.InitializeTitle, text.InitUser.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserDescription, existingText.InitializeDescription, text.InitUser.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserCodeLabel, existingText.InitializeCodeLabel, text.InitUser.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordLabel, existingText.InitializeNewPassword, text.InitUser.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, existingText.InitializeNewPasswordConfirm, text.InitUser.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserResendButtonText, existingText.InitializeResendButtonText, text.InitUser.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNextButtonText, existingText.InitializeNextButtonText, text.InitUser.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitUserDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneTitle, existingText.InitializeDoneTitle, text.InitUserDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneDescription, existingText.InitializeDoneDescription, text.InitUserDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneAbortButtonText, existingText.InitializeDoneAbortButtonText, text.InitUserDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneNextButtonText, existingText.InitializeDoneNextButtonText, text.InitUserDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAPromptEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptTitle, existingText.InitMFAPromptTitle, text.InitMFAPrompt.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptDescription, existingText.InitMFAPromptDescription, text.InitMFAPrompt.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptOTPOption, existingText.InitMFAPromptOTPOption, text.InitMFAPrompt.Provider0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptU2FOption, existingText.InitMFAPromptU2FOption, text.InitMFAPrompt.Provider1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptSkipButtonText, existingText.InitMFAPromptSkipButtonText, text.InitMFAPrompt.SkipButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptNextButtonText, existingText.InitMFAPromptNextButtonText, text.InitMFAPrompt.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPTitle, existingText.InitMFAOTPTitle, text.InitMFAOTP.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescription, existingText.InitMFAOTPDescription, text.InitMFAOTP.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescriptionOTP, existingText.InitMFAOTPDescriptionOTP, text.InitMFAOTP.OTPDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPSecretLabel, existingText.InitMFAOTPSecretLabel, text.InitMFAOTP.SecretLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCodeLabel, existingText.InitMFAOTPCodeLabel, text.InitMFAOTP.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPNextButtonText, existingText.InitMFAOTPNextButtonText, text.InitMFAOTP.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCancelButtonText, existingText.InitMFAOTPCancelButtonText, text.InitMFAOTP.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTitle, existingText.InitMFAU2FTitle, text.InitMFAU2F.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FDescription, existingText.InitMFAU2FDescription, text.InitMFAU2F.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTokenNameLabel, existingText.InitMFAU2FTokenNameLabel, text.InitMFAU2F.TokenNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, existingText.InitMFAU2FRegisterTokenButtonText, text.InitMFAU2F.RegisterTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FNotSupported, existingText.InitMFAU2FNotSupported, text.InitMFAU2F.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FErrorRetry, existingText.InitMFAU2FErrorRetry, text.InitMFAU2F.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFADoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneTitle, existingText.InitMFADoneTitle, text.InitMFADone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneDescription, existingText.InitMFADoneDescription, text.InitMFADone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneAbortButtonText, existingText.InitMFADoneAbortButtonText, text.InitMFADone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneNextButtonText, existingText.InitMFADoneNextButtonText, text.InitMFADone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getMFAProviderEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersChooseOther, existingText.MFAProvidersChooseOther, text.MFAProvider.ChooseOther, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersOTP, existingText.MFAProvidersOTP, text.MFAProvider.Provider0, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersU2F, existingText.MFAProvidersU2F, text.MFAProvider.Provider1, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getVerifyMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPTitle, existingText.VerifyMFAOTPTitle, text.VerifyMFAOTP.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPDescription, existingText.VerifyMFAOTPDescription, text.VerifyMFAOTP.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPCodeLabel, existingText.VerifyMFAOTPCodeLabel, text.VerifyMFAOTP.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPNextButtonText, existingText.VerifyMFAOTPNextButtonText, text.VerifyMFAOTP.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getVerifyMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FTitle, existingText.VerifyMFAU2FTitle, text.VerifyMFAU2F.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FDescription, existingText.VerifyMFAU2FDescription, text.VerifyMFAU2F.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FValidateTokenText, existingText.VerifyMFAU2FValidateTokenText, text.VerifyMFAU2F.ValidateTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FNotSupported, existingText.VerifyMFAU2FNotSupported, text.VerifyMFAU2F.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FErrorRetry, existingText.VerifyMFAU2FErrorRetry, text.VerifyMFAU2F.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordlessEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessTitle, existingText.PasswordlessTitle, text.Passwordless.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessDescription, existingText.PasswordlessDescription, text.Passwordless.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessLoginWithPwButtonText, existingText.PasswordlessLoginWithPwButtonText, text.Passwordless.LoginWithPwButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessValidateTokenButtonText, existingText.PasswordlessValidateTokenButtonText, text.Passwordless.ValidateTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessNotSupported, existingText.PasswordlessNotSupported, text.Passwordless.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessErrorRetry, existingText.PasswordlessErrorRetry, text.Passwordless.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordChangeEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeTitle, existingText.PasswordChangeTitle, text.PasswordChange.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDescription, existingText.PasswordChangeDescription, text.PasswordChange.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeOldPasswordLabel, existingText.PasswordChangeOldPasswordLabel, text.PasswordChange.OldPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordLabel, existingText.PasswordChangeNewPasswordLabel, text.PasswordChange.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, existingText.PasswordChangeNewPasswordConfirmLabel, text.PasswordChange.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeCancelButtonText, existingText.PasswordChangeCancelButtonText, text.PasswordChange.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNextButtonText, existingText.PasswordChangeNextButtonText, text.PasswordChange.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordChangeDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneTitle, existingText.PasswordChangeDoneTitle, text.PasswordChangeDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneDescription, existingText.PasswordChangeDoneDescription, text.PasswordChangeDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneNextButtonText, existingText.PasswordChangeDoneNextButtonText, text.PasswordChangeDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordResetDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneTitle, existingText.PasswordResetDoneTitle, text.PasswordResetDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneDescription, existingText.PasswordResetDoneDescription, text.PasswordResetDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneNextButtonText, existingText.PasswordResetDoneNextButtonText, text.PasswordResetDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationOptionEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionTitle, existingText.RegistrationOptionTitle, text.RegisterOption.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionDescription, existingText.RegistrationOptionDescription, text.RegisterOption.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionUserNameButtonText, existingText.RegistrationOptionUserNameButtonText, text.RegisterOption.RegisterUsernamePasswordButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionExternalLoginDescription, existingText.RegistrationOptionExternalLoginDescription, text.RegisterOption.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTitle, existingText.RegistrationUserTitle, text.RegistrationUser.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserDescription, existingText.RegistrationUserDescription, text.RegistrationUser.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserFirstnameLabel, existingText.RegistrationUserFirstnameLabel, text.RegistrationUser.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLastnameLabel, existingText.RegistrationUserLastnameLabel, text.RegistrationUser.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserEmailLabel, existingText.RegistrationUserEmailLabel, text.RegistrationUser.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserUsernameLabel, existingText.RegistrationUserUsernameLabel, text.RegistrationUser.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLanguageLabel, existingText.RegistrationUserLanguageLabel, text.RegistrationUser.LanguageLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserGenderLabel, existingText.RegistrationUserGenderLabel, text.RegistrationUser.GenderLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordLabel, existingText.RegistrationUserPasswordLabel, text.RegistrationUser.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordConfirmLabel, existingText.RegistrationUserPasswordConfirmLabel, text.RegistrationUser.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, existingText.RegistrationUserTOSAndPrivacyLabel, text.RegistrationUser.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSConfirm, existingText.RegistrationUserTOSConfirm, text.RegistrationUser.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSLink, existingText.RegistrationUserTOSLink, text.RegistrationUser.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSLinkText, existingText.RegistrationUserTOSLinkText, text.RegistrationUser.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyConfirm, existingText.RegistrationUserPrivacyConfirm, text.RegistrationUser.PrivacyConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyLink, existingText.RegistrationUserPrivacyLink, text.RegistrationUser.PrivacyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyLinkText, existingText.RegistrationUserPrivacyLinkText, text.RegistrationUser.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserExternalLoginDescription, existingText.RegistrationUserExternalLoginDescription, text.RegistrationUser.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserNextButtonText, existingText.RegistrationUserNextButtonText, text.RegistrationUser.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserBackButtonText, existingText.RegistrationUserBackButtonText, text.RegistrationUser.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationOrgEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTitle, existingText.RegisterOrgTitle, text.RegistrationOrg.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgDescription, existingText.RegisterOrgDescription, text.RegistrationOrg.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgOrgNameLabel, existingText.RegisterOrgOrgNameLabel, text.RegistrationOrg.OrgNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgFirstnameLabel, existingText.RegisterOrgFirstnameLabel, text.RegistrationOrg.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgLastnameLabel, existingText.RegisterOrgLastnameLabel, text.RegistrationOrg.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgUsernameLabel, existingText.RegisterOrgUsernameLabel, text.RegistrationOrg.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgEmailLabel, existingText.RegisterOrgEmailLabel, text.RegistrationOrg.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordLabel, existingText.RegisterOrgPasswordLabel, text.RegistrationOrg.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordConfirmLabel, existingText.RegisterOrgPasswordConfirmLabel, text.RegistrationOrg.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, existingText.RegisterOrgTOSAndPrivacyLabel, text.RegistrationOrg.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSConfirm, existingText.RegisterOrgTOSConfirm, text.RegistrationOrg.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSLink, existingText.RegisterOrgTOSLink, text.RegistrationOrg.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSLinkText, existingText.RegisterOrgTOSLinkText, text.RegistrationOrg.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyConfirm, existingText.RegisterOrgPrivacyConfirm, text.RegistrationOrg.PrivacyConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyLink, existingText.RegisterOrgPrivacyLink, text.RegistrationOrg.PrivacyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyLinkText, existingText.RegisterOrgPrivacyLinkText, text.RegistrationOrg.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgExternalLoginDescription, existingText.RegisterOrgExternalLoginDescription, text.RegistrationOrg.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgSaveButtonText, existingText.RegisterOrgSaveButtonText, text.RegistrationOrg.SaveButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLinkingUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneTitle, existingText.LinkingUserDoneTitle, text.LinkingUsersDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneDescription, existingText.LinkingUserDoneDescription, text.LinkingUsersDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneCancelButtonText, existingText.LinkingUserDoneCancelButtonText, text.LinkingUsersDone.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneNextButtonText, existingText.LinkingUserDoneNextButtonText, text.LinkingUsersDone.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getExternalUserNotFoundEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundTitle, existingText.ExternalUserNotFoundTitle, text.ExternalNotFoundOption.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundDescription, existingText.ExternalUserNotFoundDescription, text.ExternalNotFoundOption.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundLinkButtonText, existingText.ExternalUserNotFoundLinkButtonText, text.ExternalNotFoundOption.LinkButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, existingText.ExternalUserNotFoundAutoRegisterButtonText, text.ExternalNotFoundOption.AutoRegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getSuccessLoginEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginTitle, existingText.SuccessLoginTitle, text.LoginSuccess.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginAutoRedirectDescription, existingText.SuccessLoginAutoRedirectDescription, text.LoginSuccess.AutoRedirectDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginRedirectedDescription, existingText.SuccessLoginRedirectedDescription, text.LoginSuccess.RedirectedDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginNextButtonText, existingText.SuccessLoginNextButtonText, text.LoginSuccess.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLogoutDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneTitle, existingText.LogoutDoneTitle, text.LogoutDone.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneDescription, existingText.LogoutDoneDescription, text.LogoutDone.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneLoginButtonText, existingText.LogoutDoneLoginButtonText, text.LogoutDone.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getFooterTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterTos, existingText.FooterTOS, text.Footer.TOS, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterTosLink, existingText.FooterTOSLink, text.Footer.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterPrivacy, existingText.FooterPrivacyPolicy, text.Footer.PrivacyPolicy, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterPrivacyLink, existingText.FooterPrivacyPolicyLink, text.Footer.PrivacyPolicyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelp, existingText.FooterHelp, text.Footer.Help, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelpLink, existingText.FooterHelpLink, text.Footer.HelpLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getCustomLoginTextEvent(ctx context.Context, agg *eventstore.Aggregate, textKey, existingText, newText string, lang language.Tag, defaultText bool) eventstore.EventPusher {
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
