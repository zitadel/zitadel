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
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitle, existingText.SelectAccountTitle, text.SelectAccountScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescription, existingText.SelectAccountDescription, text.SelectAccountScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountTitleLinkingProcess, existingText.SelectAccountTitleLinkingProcess, text.SelectAccountScreenText.TitleLinkingProcess, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountDescriptionLinkingProcess, existingText.SelectAccountDescriptionLinkingProcess, text.SelectAccountScreenText.DescriptionLinkingProcess, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountOtherUser, existingText.SelectAccountOtherUser, text.SelectAccountScreenText.OtherUser, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateActive, existingText.SelectAccountSessionStateActive, text.SelectAccountScreenText.SessionStateActive, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountSessionStateInactive, existingText.SelectAccountSessionStateInactive, text.SelectAccountScreenText.SessionStateInactive, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, existingText.SelectAccountUserMustBeMemberOfOrg, text.SelectAccountScreenText.UserMustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLoginTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitle, existingText.LoginTitle, text.LoginScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescription, existingText.LoginDescription, text.LoginScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginTitleLinkingProcess, existingText.LoginTitleLinkingProcess, text.LoginScreenText.TitleLinkingProcess, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginDescriptionLinkingProcess, existingText.LoginDescriptionLinkingProcess, text.LoginScreenText.DescriptionLinkingProcess, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNameLabel, existingText.LoginNameLabel, text.LoginScreenText.LoginNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUsernamePlaceHolder, existingText.LoginUsernamePlaceholder, text.LoginScreenText.UsernamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginLoginnamePlaceHolder, existingText.LoginLoginnamePlaceholder, text.LoginScreenText.UsernamePlaceholder, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginRegisterButtonText, existingText.LoginRegisterButtonText, text.LoginScreenText.RegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginNextButtonText, existingText.LoginNextButtonText, text.LoginScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginExternalUserDescription, existingText.LoginExternalUserDescription, text.LoginScreenText.ExternalUserDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLoginUserMustBeMemberOfOrg, existingText.LoginUserMustBeMemberOfOrg, text.LoginScreenText.UserMustBeMemberOfOrg, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordTitle, existingText.PasswordTitle, text.PasswordScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordDescription, existingText.PasswordDescription, text.PasswordScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordLabel, existingText.PasswordLabel, text.PasswordScreenText.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetLinkText, existingText.PasswordResetLinkText, text.PasswordScreenText.ResetLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordBackButtonText, existingText.PasswordBackButtonText, text.PasswordScreenText.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordNextButtonText, existingText.PasswordNextButtonText, text.PasswordScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordMinLength, existingText.PasswordMinLength, text.PasswordScreenText.MinLength, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasUppercase, existingText.PasswordHasUppercase, text.PasswordScreenText.HasUppercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasLowercase, existingText.PasswordHasLowercase, text.PasswordScreenText.HasLowercase, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasNumber, existingText.PasswordHasNumber, text.PasswordScreenText.HasNumber, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordHasSymbol, existingText.PasswordHasSymbol, text.PasswordScreenText.HasSymbol, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordConfirmation, existingText.PasswordConfirmation, text.PasswordScreenText.Confirmation, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getUsernameChangeTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeTitle, existingText.UsernameChangeTitle, text.UsernameChangeScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDescription, existingText.UsernameChangeDescription, text.UsernameChangeScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeUsernameLabel, existingText.UsernameChangeUsernameLabel, text.UsernameChangeScreenText.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeCancelButtonText, existingText.UsernameChangeCancelButtonText, text.UsernameChangeScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeNextButtonText, existingText.UsernameChangeNextButtonText, text.UsernameChangeScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getUsernameChangeDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneTitle, existingText.UsernameChangeDoneTitle, text.UsernameChangeDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneDescription, existingText.UsernameChangeDoneDescription, text.UsernameChangeDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyUsernameChangeDoneNextButtonText, existingText.UsernameChangeDoneNextButtonText, text.UsernameChangeDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordInitTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordTitle, existingText.InitPasswordTitle, text.InitPasswordScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDescription, existingText.InitPasswordDescription, text.InitPasswordScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordCodeLabel, existingText.InitPasswordCodeLabel, text.InitPasswordScreenText.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordLabel, existingText.InitPasswordNewPasswordLabel, text.InitPasswordScreenText.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, existingText.InitPasswordNewPasswordConfirmLabel, text.InitPasswordScreenText.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNextButtonText, existingText.InitPasswordNextButtonText, text.InitPasswordScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordResendButtonText, existingText.InitPasswordResendButtonText, text.InitPasswordScreenText.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordInitDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneTitle, existingText.InitPasswordDoneTitle, text.InitPasswordDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneDescription, existingText.InitPasswordDoneDescription, text.InitPasswordDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneNextButtonText, existingText.InitPasswordDoneNextButtonText, text.InitPasswordDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDoneCancelButtonText, existingText.InitPasswordDoneCancelButtonText, text.InitPasswordDoneScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getEmailVerificationTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationTitle, existingText.EmailVerificationTitle, text.EmailVerificationScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDescription, existingText.EmailVerificationDescription, text.EmailVerificationScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationCodeLabel, existingText.EmailVerificationCodeLabel, text.EmailVerificationScreenText.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationNextButtonText, existingText.EmailVerificationNextButtonText, text.EmailVerificationScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationResendButtonText, existingText.EmailVerificationResendButtonText, text.EmailVerificationScreenText.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getEmailVerificationDoneTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneTitle, existingText.EmailVerificationDoneTitle, text.EmailVerificationDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneDescription, existingText.EmailVerificationDoneDescription, text.EmailVerificationDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneNextButtonText, existingText.EmailVerificationDoneNextButtonText, text.EmailVerificationDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneCancelButtonText, existingText.EmailVerificationDoneCancelButtonText, text.EmailVerificationDoneScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyEmailVerificationDoneLoginButtonText, existingText.EmailVerificationDoneLoginButtonText, text.EmailVerificationDoneScreenText.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserTitle, existingText.InitializeTitle, text.InitializeUserScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserDescription, existingText.InitializeDescription, text.InitializeUserScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserCodeLabel, existingText.InitializeCodeLabel, text.InitializeUserScreenText.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordLabel, existingText.InitializeNewPassword, text.InitializeUserScreenText.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, existingText.InitializeNewPasswordConfirm, text.InitializeUserScreenText.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserResendButtonText, existingText.InitializeResendButtonText, text.InitializeUserScreenText.ResendButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitializeUserNextButtonText, existingText.InitializeNextButtonText, text.InitializeUserScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitUserDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneTitle, existingText.InitializeDoneTitle, text.InitializeDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneDescription, existingText.InitializeDoneDescription, text.InitializeDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneAbortButtonText, existingText.InitializeDoneAbortButtonText, text.InitializeDoneScreenText.AbortButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitUserDoneNextButtonText, existingText.InitializeDoneNextButtonText, text.InitializeDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAPromptEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptTitle, existingText.InitMFAPromptTitle, text.InitMFAPromptScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptDescription, existingText.InitMFAPromptDescription, text.InitMFAPromptScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptOTPOption, existingText.InitMFAPromptOTPOption, text.InitMFAPromptScreenText.OTPOption, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptU2FOption, existingText.InitMFAPromptU2FOption, text.InitMFAPromptScreenText.U2FOption, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptSkipButtonText, existingText.InitMFAPromptSkipButtonText, text.InitMFAPromptScreenText.SkipButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAPromptNextButtonText, existingText.InitMFAPromptNextButtonText, text.InitMFAPromptScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPTitle, existingText.InitMFAOTPTitle, text.InitMFAOTPScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescription, existingText.InitMFAOTPDescription, text.InitMFAOTPScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPDescriptionOTP, existingText.InitMFAOTPDescriptionOTP, text.InitMFAOTPScreenText.DescriptionOTP, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPSecretLabel, existingText.InitMFAOTPSecretLabel, text.InitMFAOTPScreenText.SecretLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCodeLabel, existingText.InitMFAOTPCodeLabel, text.InitMFAOTPScreenText.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPNextButtonText, existingText.InitMFAOTPNextButtonText, text.InitMFAOTPScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCancelButtonText, existingText.InitMFAOTPCancelButtonText, text.InitMFAOTPScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTitle, existingText.InitMFAU2FTitle, text.InitMFAU2FScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FDescription, existingText.InitMFAU2FDescription, text.InitMFAU2FScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FTokenNameLabel, existingText.InitMFAU2FTokenNameLabel, text.InitMFAU2FScreenText.TokenNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, existingText.InitMFAU2FRegisterTokenButtonText, text.InitMFAU2FScreenText.RegisterTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FNotSupported, existingText.InitMFAU2FNotSupported, text.InitMFAU2FScreenText.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAU2FErrorRetry, existingText.InitMFAU2FErrorRetry, text.InitMFAU2FScreenText.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getInitMFADoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneTitle, existingText.InitMFADoneTitle, text.InitMFADoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneDescription, existingText.InitMFADoneDescription, text.InitMFADoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneAbortButtonText, existingText.InitMFADoneAbortButtonText, text.InitMFADoneScreenText.AbortButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFADoneNextButtonText, existingText.InitMFADoneNextButtonText, text.InitMFADoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getMFAProviderEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersChooseOther, existingText.MFAProvidersChooseOther, text.MFAProvidersText.ChooseOther, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersOTP, existingText.MFAProvidersOTP, text.MFAProvidersText.OTP, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyMFAProvidersU2F, existingText.MFAProvidersU2F, text.MFAProvidersText.U2F, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getVerifyMFAOTPEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPTitle, existingText.VerifyMFAOTPTitle, text.VerifyMFAOTPScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPDescription, existingText.VerifyMFAOTPDescription, text.VerifyMFAOTPScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPCodeLabel, existingText.VerifyMFAOTPCodeLabel, text.VerifyMFAOTPScreenText.CodeLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAOTPNextButtonText, existingText.VerifyMFAOTPNextButtonText, text.VerifyMFAOTPScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getVerifyMFAU2FEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FTitle, existingText.VerifyMFAU2FTitle, text.VerifyMFAU2FScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FDescription, existingText.VerifyMFAU2FDescription, text.VerifyMFAU2FScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FValidateTokenText, existingText.VerifyMFAU2FValidateTokenText, text.VerifyMFAU2FScreenText.ValidateTokenText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FNotSupported, existingText.VerifyMFAU2FNotSupported, text.VerifyMFAU2FScreenText.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FErrorRetry, existingText.VerifyMFAU2FErrorRetry, text.VerifyMFAU2FScreenText.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordlessEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessTitle, existingText.PasswordlessTitle, text.PasswordlessScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessDescription, existingText.PasswordlessDescription, text.PasswordlessScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessLoginWithPwButtonText, existingText.PasswordlessLoginWithPwButtonText, text.PasswordlessScreenText.LoginWithPwButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessValidateTokenButtonText, existingText.PasswordlessValidateTokenButtonText, text.PasswordlessScreenText.ValidateTokenButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessNotSupported, existingText.PasswordlessNotSupported, text.PasswordlessScreenText.NotSupported, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordlessErrorRetry, existingText.PasswordlessErrorRetry, text.PasswordlessScreenText.ErrorRetry, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordChangeEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeTitle, existingText.PasswordChangeTitle, text.PasswordChangeScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDescription, existingText.PasswordChangeDescription, text.PasswordChangeScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeOldPasswordLabel, existingText.PasswordChangeOldPasswordLabel, text.PasswordChangeScreenText.OldPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordLabel, existingText.PasswordChangeNewPasswordLabel, text.PasswordChangeScreenText.NewPasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, existingText.PasswordChangeNewPasswordConfirmLabel, text.PasswordChangeScreenText.NewPasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeCancelButtonText, existingText.PasswordChangeCancelButtonText, text.PasswordChangeScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeNextButtonText, existingText.PasswordChangeNextButtonText, text.PasswordChangeScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordChangeDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneTitle, existingText.PasswordChangeDoneTitle, text.PasswordChangeDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneDescription, existingText.PasswordChangeDoneDescription, text.PasswordChangeDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordChangeDoneNextButtonText, existingText.PasswordChangeDoneNextButtonText, text.PasswordChangeDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getPasswordResetDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneTitle, existingText.PasswordResetDoneTitle, text.PasswordResetDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneDescription, existingText.PasswordResetDoneDescription, text.PasswordResetDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyPasswordResetDoneNextButtonText, existingText.PasswordResetDoneNextButtonText, text.PasswordResetDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationOptionEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionTitle, existingText.RegistrationOptionTitle, text.RegistrationOptionScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionDescription, existingText.RegistrationOptionDescription, text.RegistrationOptionScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionUserNameButtonText, existingText.RegistrationOptionUserNameButtonText, text.RegistrationOptionScreenText.UserNameButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationOptionExternalLoginDescription, existingText.RegistrationOptionExternalLoginDescription, text.RegistrationOptionScreenText.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTitle, existingText.RegistrationUserTitle, text.RegistrationUserScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserDescription, existingText.RegistrationUserDescription, text.RegistrationUserScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserFirstnameLabel, existingText.RegistrationUserFirstnameLabel, text.RegistrationUserScreenText.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLastnameLabel, existingText.RegistrationUserLastnameLabel, text.RegistrationUserScreenText.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserEmailLabel, existingText.RegistrationUserEmailLabel, text.RegistrationUserScreenText.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserUsernameLabel, existingText.RegistrationUserUsernameLabel, text.RegistrationUserScreenText.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserLanguageLabel, existingText.RegistrationUserLanguageLabel, text.RegistrationUserScreenText.LanguageLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserGenderLabel, existingText.RegistrationUserGenderLabel, text.RegistrationUserScreenText.GenderLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordLabel, existingText.RegistrationUserPasswordLabel, text.RegistrationUserScreenText.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPasswordConfirmLabel, existingText.RegistrationUserPasswordConfirmLabel, text.RegistrationUserScreenText.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, existingText.RegistrationUserTOSAndPrivacyLabel, text.RegistrationUserScreenText.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSConfirm, existingText.RegistrationUserTOSConfirm, text.RegistrationUserScreenText.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSLink, existingText.RegistrationUserTOSLink, text.RegistrationUserScreenText.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserTOSLinkText, existingText.RegistrationUserTOSLinkText, text.RegistrationUserScreenText.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyConfirm, existingText.RegistrationUserPrivacyConfirm, text.RegistrationUserScreenText.PrivacyConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyLink, existingText.RegistrationUserPrivacyLink, text.RegistrationUserScreenText.PrivacyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserPrivacyLinkText, existingText.RegistrationUserPrivacyLinkText, text.RegistrationUserScreenText.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserExternalLoginDescription, existingText.RegistrationUserExternalLoginDescription, text.RegistrationUserScreenText.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserNextButtonText, existingText.RegistrationUserNextButtonText, text.RegistrationUserScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserBackButtonText, existingText.RegistrationUserBackButtonText, text.RegistrationUserScreenText.BackButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getRegistrationOrgEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTitle, existingText.RegisterOrgTitle, text.RegistrationOrgScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgDescription, existingText.RegisterOrgDescription, text.RegistrationOrgScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgOrgNameLabel, existingText.RegisterOrgOrgNameLabel, text.RegistrationOrgScreenText.OrgNameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgFirstnameLabel, existingText.RegisterOrgFirstnameLabel, text.RegistrationOrgScreenText.FirstnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgLastnameLabel, existingText.RegisterOrgLastnameLabel, text.RegistrationOrgScreenText.LastnameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgUsernameLabel, existingText.RegisterOrgUsernameLabel, text.RegistrationOrgScreenText.UsernameLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgEmailLabel, existingText.RegisterOrgEmailLabel, text.RegistrationOrgScreenText.EmailLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordLabel, existingText.RegisterOrgPasswordLabel, text.RegistrationOrgScreenText.PasswordLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPasswordConfirmLabel, existingText.RegisterOrgPasswordConfirmLabel, text.RegistrationOrgScreenText.PasswordConfirmLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, existingText.RegisterOrgTOSAndPrivacyLabel, text.RegistrationOrgScreenText.TOSAndPrivacyLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSConfirm, existingText.RegisterOrgTOSConfirm, text.RegistrationOrgScreenText.TOSConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSLink, existingText.RegisterOrgTOSLink, text.RegistrationOrgScreenText.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgTOSLinkText, existingText.RegisterOrgTOSLinkText, text.RegistrationOrgScreenText.TOSLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyConfirm, existingText.RegisterOrgPrivacyConfirm, text.RegistrationOrgScreenText.PrivacyConfirm, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyLink, existingText.RegisterOrgPrivacyLink, text.RegistrationOrgScreenText.PrivacyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgPrivacyLinkText, existingText.RegisterOrgPrivacyLinkText, text.RegistrationOrgScreenText.PrivacyLinkText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgExternalLoginDescription, existingText.RegisterOrgExternalLoginDescription, text.RegistrationOrgScreenText.ExternalLoginDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgSaveButtonText, existingText.RegisterOrgSaveButtonText, text.RegistrationOrgScreenText.SaveButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLinkingUserEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneTitle, existingText.LinkingUserDoneTitle, text.LinkingUserDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneDescription, existingText.LinkingUserDoneDescription, text.LinkingUserDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneCancelButtonText, existingText.LinkingUserDoneCancelButtonText, text.LinkingUserDoneScreenText.CancelButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLinkingUserDoneNextButtonText, existingText.LinkingUserDoneNextButtonText, text.LinkingUserDoneScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getExternalUserNotFoundEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundTitle, existingText.ExternalUserNotFoundTitle, text.ExternalUserNotFoundScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundDescription, existingText.ExternalUserNotFoundDescription, text.ExternalUserNotFoundScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundLinkButtonText, existingText.ExternalUserNotFoundLinkButtonText, text.ExternalUserNotFoundScreenText.LinkButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, existingText.ExternalUserNotFoundAutoRegisterButtonText, text.ExternalUserNotFoundScreenText.AutoRegisterButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getSuccessLoginEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginTitle, existingText.SuccessLoginTitle, text.SuccessLoginScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginAutoRedirectDescription, existingText.SuccessLoginAutoRedirectDescription, text.SuccessLoginScreenText.AutoRedirectDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginRedirectedDescription, existingText.SuccessLoginRedirectedDescription, text.SuccessLoginScreenText.RedirectedDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeySuccessLoginNextButtonText, existingText.SuccessLoginNextButtonText, text.SuccessLoginScreenText.NextButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getLogoutDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneTitle, existingText.LogoutDoneTitle, text.LogoutDoneScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneDescription, existingText.LogoutDoneDescription, text.LogoutDoneScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyLogoutDoneLoginButtonText, existingText.LogoutDoneLoginButtonText, text.LogoutDoneScreenText.LoginButtonText, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	return events
}

func (c *Commands) getFooterTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterTos, existingText.FooterTOS, text.FooterText.TOS, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterTosLink, existingText.FooterTOSLink, text.FooterText.TOSLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterPrivacy, existingText.FooterPrivacyPolicy, text.FooterText.PrivacyPolicy, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterPrivacyLink, existingText.FooterPrivacyPolicyLink, text.FooterText.PrivacyPolicyLink, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelp, existingText.FooterHelp, text.FooterText.Help, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyFooterHelpLink, existingText.FooterHelpLink, text.FooterText.HelpLink, text.Language, defaultText)
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
