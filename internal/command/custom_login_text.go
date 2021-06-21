package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

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
	return events
}

func (c *Commands) getPasswordResetTextEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, 0)
	event := c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordTitle, existingText.InitPasswordTitle, text.InitPasswordScreenText.Title, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordDescription, existingText.InitPasswordDescription, text.InitPasswordScreenText.Description, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitPasswordNextButtonText, existingText.InitPasswordNextButtonText, text.InitPasswordScreenText.NextButtonText, text.Language, defaultText)
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

func (c *Commands) getInitDoneEvents(ctx context.Context, agg *eventstore.Aggregate, existingText *CustomLoginTextReadModel, text *domain.CustomLoginText, defaultText bool) []eventstore.EventPusher {
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
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPSecretLabel, existingText.InitMFAOTPSecretLabel, text.InitMFAOTPScreenText.SecretLabel, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyInitMFAOTPCodeLabel, existingText.InitMFAOTPCodeLabel, text.InitMFAOTPScreenText.CodeLabel, text.Language, defaultText)
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
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FOtherOptionDescription, existingText.VerifyMFAU2FOtherOptionDescription, text.VerifyMFAU2FScreenText.OtherOptionDescription, text.Language, defaultText)
	if event != nil {
		events = append(events, event)
	}
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyVerifyMFAU2FOptionOTPButtonText, existingText.VerifyMFAU2FOptionOTPButtonText, text.VerifyMFAU2FScreenText.OptionOTPButtonText, text.Language, defaultText)
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
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegistrationUserSaveButtonText, existingText.RegistrationUserSaveButtonText, text.RegistrationUserScreenText.SaveButtonText, text.Language, defaultText)
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
	event = c.getCustomLoginTextEvent(ctx, agg, domain.LoginKeyRegisterOrgSaveButtonText, existingText.RegisterOrgSaveButtonText, text.RegistrationOrgScreenText.SaveButtonText, text.Language, defaultText)
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
