package admin

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/grpc/text"
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

func SetPasswordlessRegistrationCustomTextToDomain(msg *admin_pb.SetDefaultPasswordlessRegistrationMessageTextRequest) *domain.CustomMessageText {
	langTag := language.Make(msg.Language)
	return &domain.CustomMessageText{
		MessageTextType: domain.PasswordlessRegistrationMessageType,
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

func SetLoginTextToDomain(req *admin_pb.SetCustomLoginTextsRequest) *domain.CustomLoginText {
	langTag := language.Make(req.Language)
	result := &domain.CustomLoginText{
		Language: langTag,
	}
	result.SelectAccount = text.SelectAccountScreenTextPbToDomain(req.SelectAccountText)
	result.Login = text.LoginScreenTextPbToDomain(req.LoginText)
	result.Password = text.PasswordScreenTextPbToDomain(req.PasswordText)
	result.UsernameChange = text.UsernameChangeScreenTextPbToDomain(req.UsernameChangeText)
	result.UsernameChangeDone = text.UsernameChangeDoneScreenTextPbToDomain(req.UsernameChangeDoneText)
	result.Password = text.PasswordScreenTextPbToDomain(req.PasswordText)
	result.InitPassword = text.InitPasswordScreenTextPbToDomain(req.InitPasswordText)
	result.InitPasswordDone = text.InitPasswordDoneScreenTextPbToDomain(req.InitPasswordDoneText)
	result.EmailVerification = text.EmailVerificationScreenTextPbToDomain(req.EmailVerificationText)
	result.EmailVerificationDone = text.EmailVerificationDoneScreenTextPbToDomain(req.EmailVerificationDoneText)
	result.InitUser = text.InitializeUserScreenTextPbToDomain(req.InitializeUserText)
	result.InitUserDone = text.InitializeDoneScreenTextPbToDomain(req.InitializeDoneText)
	result.InitMFAPrompt = text.InitMFAPromptScreenTextPbToDomain(req.InitMfaPromptText)
	result.InitMFAOTP = text.InitMFAOTPScreenTextPbToDomain(req.InitMfaOtpText)
	result.InitMFAU2F = text.InitMFAU2FScreenTextPbToDomain(req.InitMfaU2FText)
	result.InitMFADone = text.InitMFADoneScreenTextPbToDomain(req.InitMfaDoneText)
	result.MFAProvider = text.MFAProvidersTextPbToDomain(req.MfaProvidersText)
	result.VerifyMFAOTP = text.VerifyMFAOTPScreenTextPbToDomain(req.VerifyMfaOtpText)
	result.VerifyMFAU2F = text.VerifyMFAU2FScreenTextPbToDomain(req.VerifyMfaU2FText)
	result.Passwordless = text.PasswordlessScreenTextPbToDomain(req.PasswordlessText)
	result.PasswordlessRegistration = text.PasswordlessRegistrationScreenTextPbToDomain(req.PasswordlessRegistrationText)
	result.PasswordlessRegistrationDone = text.PasswordlessRegistrationDoneScreenTextPbToDomain(req.PasswordlessRegistrationDoneText)
	result.PasswordChange = text.PasswordChangeScreenTextPbToDomain(req.PasswordChangeText)
	result.PasswordChangeDone = text.PasswordChangeDoneScreenTextPbToDomain(req.PasswordChangeDoneText)
	result.PasswordResetDone = text.PasswordResetDoneScreenTextPbToDomain(req.PasswordResetDoneText)
	result.RegisterOption = text.RegistrationOptionScreenTextPbToDomain(req.RegistrationOptionText)
	result.RegistrationUser = text.RegistrationUserScreenTextPbToDomain(req.RegistrationUserText)
	result.RegistrationOrg = text.RegistrationOrgScreenTextPbToDomain(req.RegistrationOrgText)
	result.LinkingUsersDone = text.LinkingUserDoneScreenTextPbToDomain(req.LinkingUserDoneText)
	result.ExternalNotFoundOption = text.ExternalUserNotFoundScreenTextPbToDomain(req.ExternalUserNotFoundText)
	result.LoginSuccess = text.SuccessLoginScreenTextPbToDomain(req.SuccessLoginText)
	result.LogoutDone = text.LogoutDoneScreenTextPbToDomain(req.LogoutText)
	result.Footer = text.FooterTextPbToDomain(req.FooterText)
	return result
}
