package admin

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func SetInitCustomTextToDomain(msg *admin_pb.SetDefaultInitMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.InitCodeMessageType,
		Language:        langTag,
		Title:           msg.Title,
		PreHeader:       msg.PreHeader,
		Subject:         msg.Subject,
		Greeting:        msg.Greeting,
		Text:            msg.Text,
		ButtonText:      msg.ButtonText,
		FooterText:      msg.FooterText,
	}
}

func SetPasswordResetCustomTextToDomain(msg *admin_pb.SetDefaultPasswordResetMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.PasswordResetMessageType,
		Language:        langTag,
		Title:           msg.Title,
		PreHeader:       msg.PreHeader,
		Subject:         msg.Subject,
		Greeting:        msg.Greeting,
		Text:            msg.Text,
		ButtonText:      msg.ButtonText,
		FooterText:      msg.FooterText,
	}
}

func SetVerifyEmailCustomTextToDomain(msg *admin_pb.SetDefaultVerifyEmailMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.VerifyEmailMessageType,
		Language:        langTag,
		Title:           msg.Title,
		PreHeader:       msg.PreHeader,
		Subject:         msg.Subject,
		Greeting:        msg.Greeting,
		Text:            msg.Text,
		ButtonText:      msg.ButtonText,
		FooterText:      msg.FooterText,
	}
}

func SetVerifyPhoneCustomTextToDomain(msg *admin_pb.SetDefaultVerifyPhoneMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.VerifyPhoneMessageType,
		Language:        langTag,
		Title:           msg.Title,
		PreHeader:       msg.PreHeader,
		Subject:         msg.Subject,
		Greeting:        msg.Greeting,
		Text:            msg.Text,
		ButtonText:      msg.ButtonText,
		FooterText:      msg.FooterText,
	}
}

func SetDomainClaimedCustomTextToDomain(msg *admin_pb.SetDefaultDomainClaimedMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.DomainClaimedMessageType,
		Language:        langTag,
		Title:           msg.Title,
		PreHeader:       msg.PreHeader,
		Subject:         msg.Subject,
		Greeting:        msg.Greeting,
		Text:            msg.Text,
		ButtonText:      msg.ButtonText,
		FooterText:      msg.FooterText,
	}
}

func SetLoginTextToDomain(text *admin_pb.SetDefaultLoginTextsRequest) *domain.CustomLoginText {
	langTag := language.Make(text.Language)
	result := &domain.CustomLoginText{
		Language: langTag,
	}
	if text.SelectAccountText != nil {
		result.SelectAccountTitle = text.SelectAccountText.Title
		result.SelectAccountDescription = text.SelectAccountText.Description
		result.SelectAccountOtherUser = text.SelectAccountText.OtherUser
	}
	if text.LoginText != nil {
		result.LoginTitle = text.LoginText.Title
		result.LoginDescription = text.LoginText.Description
		result.LoginNameLabel = text.LoginText.NameLabel
		result.LoginRegisterButtonText = text.LoginText.RegisterButtonText
		result.LoginNextButtonText = text.LoginText.NextButtonText
		result.LoginExternalUserDescription = text.LoginText.ExternalUserDescription
	}
	if text.PasswordText != nil {
		result.PasswordTitle = text.PasswordText.Title
		result.PasswordDescription = text.PasswordText.Description
		result.PasswordLabel = text.PasswordText.PasswordLabel
		result.PasswordResetLinkText = text.PasswordText.ResetLinkText
		result.PasswordBackButtonText = text.PasswordText.BackButtonText
		result.PasswordNextButtonText = text.PasswordText.NextButtonText
	}
	if text.ResetPasswordText != nil {
		result.ResetPasswordTitle = text.ResetPasswordText.Title
		result.ResetPasswordDescription = text.ResetPasswordText.Description
		result.ResetPasswordNextButtonText = text.ResetPasswordText.NextButtonText
	}
	if text.InitializeUserText != nil {
		result.InitializeUserTitle = text.InitializeUserText.Title
		result.InitializeUserDescription = text.InitializeUserText.Description
		result.InitializeUserCodeLabel = text.InitializeUserText.CodeLabel
		result.InitializeUserNewPasswordLabel = text.InitializeUserText.NewPasswordLabel
		result.InitializeUserNewPasswordConfirmLabel = text.InitializeUserText.NewPasswordConfirmLabel
		result.InitializeUserResendButtonText = text.InitializeUserText.ResendButtonText
		result.InitializeUserNextButtonText = text.InitializeUserText.NextButtonText
	}
	if text.InitializeDoneText != nil {
		result.InitializeDoneTitle = text.InitializeDoneText.Title
		result.InitializeDoneDescription = text.InitializeDoneText.Description
		result.InitializeDoneAbortButtonText = text.InitializeDoneText.AbortButtonText
		result.InitializeDoneNextButtonText = text.InitializeDoneText.NextButtonText
	}
	if text.InitMfaPromptText != nil {
		result.InitMFAPromptTitle = text.InitMfaPromptText.Title
		result.InitMFAPromptDescription = text.InitMfaPromptText.Description
		result.InitMFAPromptOTPOption = text.InitMfaPromptText.OtpOption
		result.InitMFAPromptU2FOption = text.InitMfaPromptText.U2FOption
		result.InitMFAPromptSkipButtonText = text.InitMfaPromptText.SkipButtonText
		result.InitMFAPromptNextButtonText = text.InitMfaPromptText.NextButtonText
	}
	if text.InitMfaOtpText != nil {
		result.InitMFAOTPTitle = text.InitMfaOtpText.Title
		result.InitMFAOTPDescription = text.InitMfaOtpText.Description
		result.InitMFAOTPSecretLabel = text.InitMfaOtpText.SecretLabel
		result.InitMFAOTPCodeLabel = text.InitMfaOtpText.CodeLabel
	}
	if text.InitMfaU2FText != nil {
		result.InitMFAU2FTitle = text.InitMfaU2FText.Title
		result.InitMFAU2FDescription = text.InitMfaU2FText.Description
		result.InitMFAU2FTokenNameLabel = text.InitMfaU2FText.TokenNameLabel
		result.InitMFAU2FRegisterTokenButtonText = text.InitMfaU2FText.RegisterTokenButtonText
	}
	if text.InitMfaDoneText != nil {
		result.InitMFADoneTitle = text.InitMfaDoneText.Title
		result.InitMFADoneDescription = text.InitMfaDoneText.Description
		result.InitMFADoneAbortButtonText = text.InitMfaDoneText.AbortButtonText
		result.InitMFADoneNextButtonText = text.InitMfaDoneText.NextButtonText
	}
	if text.VerifyMfaOtpText != nil {
		result.VerifyMFAOTPTitle = text.VerifyMfaOtpText.Title
		result.VerifyMFAOTPDescription = text.VerifyMfaOtpText.Description
		result.VerifyMFAOTPCodeLabel = text.VerifyMfaOtpText.CodeLabel
		result.VerifyMFAOTPNextButtonText = text.VerifyMfaOtpText.NextButtonText
	}
	if text.VerifyMfaU2FText != nil {
		result.VerifyMFAU2FTitle = text.VerifyMfaU2FText.Title
		result.VerifyMFAU2FDescription = text.VerifyMfaU2FText.Description
		result.VerifyMFAU2FValidateTokenText = text.VerifyMfaU2FText.ValidateTokenText
		result.VerifyMFAU2FOtherOptionDescription = text.VerifyMfaU2FText.OtherOptionDescription
		result.VerifyMFAU2FOptionOTPButtonText = text.VerifyMfaU2FText.OptionOtpButtonText
	}
	if text.RegistrationOptionText != nil {
		result.RegistrationOptionTitle = text.RegistrationOptionText.Title
		result.RegistrationOptionDescription = text.RegistrationOptionText.Description
		result.RegistrationOptionUserNameButtonText = text.RegistrationOptionText.UserNameButtonText
		result.RegistrationOptionExternalLoginDescription = text.RegistrationOptionText.ExternalLoginDescription
	}
	if text.RegistrationUserText != nil {
		result.RegistrationUserTitle = text.RegistrationUserText.Title
		result.RegistrationUserDescription = text.RegistrationUserText.Description
		result.RegistrationUserFirstnameLabel = text.RegistrationUserText.FirstnameLabel
		result.RegistrationUserLastnameLabel = text.RegistrationUserText.LastnameLabel
		result.RegistrationUserEmailLabel = text.RegistrationUserText.EmailLabel
		result.RegistrationUserLanguageLabel = text.RegistrationUserText.LanguageLabel
		result.RegistrationUserGenderLabel = text.RegistrationUserText.GenderLabel
		result.RegistrationUserPasswordLabel = text.RegistrationUserText.PasswordLabel
		result.RegistrationUserPasswordConfirmLabel = text.RegistrationUserText.PasswordConfirmLabel
		result.RegistrationUserSaveButtonText = text.RegistrationUserText.SaveButtonText
	}
	if text.RegistrationOrgText != nil {
		result.RegisterOrgTitle = text.RegistrationOrgText.Title
		result.RegisterOrgDescription = text.RegistrationOrgText.Description
		result.RegisterOrgDescription = text.RegistrationOrgText.OrgnameLabel
		result.RegisterOrgFirstnameLabel = text.RegistrationOrgText.FirstnameLabel
		result.RegisterOrgLastnameLabel = text.RegistrationOrgText.LastnameLabel
		result.RegisterOrgUsernameLabel = text.RegistrationOrgText.UsernameLabel
		result.RegisterOrgEmailLabel = text.RegistrationOrgText.EmailLabel
		result.RegisterOrgPasswordLabel = text.RegistrationOrgText.PasswordLabel
		result.RegisterOrgPasswordConfirmLabel = text.RegistrationOrgText.PasswordConfirmLabel
		result.RegisterOrgSaveButtonText = text.RegistrationOrgText.SaveButtonText
	}
	if text.PasswordlessText != nil {
		result.PasswordlessTitle = text.PasswordlessText.Title
		result.PasswordlessDescription = text.PasswordlessText.Description
		result.PasswordlessLoginWithPwButtonText = text.PasswordlessText.LoginWithPwButtonText
		result.PasswordlessValidateTokenButtonText = text.PasswordlessText.ValidateTokenButtonText
	}
	if text.SuccessLoginText != nil {
		result.SuccessLoginTitle = text.SuccessLoginText.Title
		result.SuccessLoginAutoRedirectDescription = text.SuccessLoginText.AutoRedirectDescription
		result.SuccessLoginRedirectedDescription = text.SuccessLoginText.RedirectedDescription
		result.SuccessLoginNextButtonText = text.SuccessLoginText.NextButtonText
	}
	return result
}
