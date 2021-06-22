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

func SetLoginTextToDomain(req *admin_pb.SetDefaultLoginTextsRequest) *domain.CustomLoginText {
	langTag := language.Make(req.Language)
	result := &domain.CustomLoginText{
		Language: langTag,
	}
	result.SelectAccountScreenText = text.SelectAccountScreenTextPbToDomain(req.SelectAccountText)
	result.LoginScreenText = text.LoginScreenTextPbToDomain(req.LoginText)
	result.PasswordScreenText = text.PasswordScreenTextPbToDomain(req.PasswordText)
	result.UsernameChangeScreenText = text.UsernameChangeScreenTextPbToDomain(req.UsernameChangeText)
	result.UsernameChangeDoneScreenText = text.UsernameChangeDoneScreenTextPbToDomain(req.UsernameChangeDoneText)
	result.PasswordScreenText = text.PasswordScreenTextPbToDomain(req.PasswordText)
	result.InitPasswordScreenText = text.InitPasswordScreenTextPbToDomain(req.InitPasswordText)
	result.InitPasswordDoneScreenText = text.InitPasswordDoneScreenTextPbToDomain(req.InitPasswordDoneText)
	result.EmailVerificationScreenText = text.EmailVerificationScreenTextPbToDomain(req.EmailVerificationText)
	result.EmailVerificationDoneScreenText = text.EmailVerificationDoneScreenTextPbToDomain(req.EmailVerificationDoneText)
	result.InitializeUserScreenText = text.InitializeUserScreenTextPbToDomain(req.InitializeUserText)
	result.InitializeDoneScreenText = text.InitializeDoneScreenTextPbToDomain(req.InitializeDoneText)
	result.InitMFAPromptScreenText = text.InitMFAPromptScreenTextPbToDomain(req.InitMfaPromptText)
	result.InitMFAOTPScreenText = text.InitMFAOTPScreenTextPbToDomain(req.InitMfaOtpText)
	result.InitMFAU2FScreenText = text.InitMFAU2FScreenTextPbToDomain(req.InitMfaU2FText)
	result.InitMFADoneScreenText = text.InitMFADoneScreenTextPbToDomain(req.InitMfaDoneText)
	result.MFAProvidersText = text.MFAProvidersTextPbToDomain(req.MfaProvidersText)
	result.VerifyMFAOTPScreenText = text.VerifyMFAOTPScreenTextPbToDomain(req.VerifyMfaOtpText)
	result.VerifyMFAU2FScreenText = text.VerifyMFAU2FScreenTextPbToDomain(req.VerifyMfaU2FText)
	result.PasswordlessScreenText = text.PasswordlessScreenTextPbToDomain(req.PasswordlessText)
	result.PasswordChangeScreenText = text.PasswordChangeScreenTextPbToDomain(req.PasswordChangeText)
	result.PasswordChangeDoneScreenText = text.PasswordChangeDoneScreenTextPbToDomain(req.PasswordChangeDoneText)
	result.PasswordResetDoneScreenText = text.PasswordResetDoneScreenTextPbToDomain(req.PasswordResetDoneText)
	result.RegistrationOptionScreenText = text.RegistrationOptionScreenTextPbToDomain(req.RegistrationOptionText)
	result.RegistrationUserScreenText = text.RegistrationUserScreenTextPbToDomain(req.RegistrationUserText)
	result.RegistrationOrgScreenText = text.RegistrationOrgScreenTextPbToDomain(req.RegistrationOrgText)
	result.LinkingUserDoneScreenText = text.LinkingUserDoneScreenTextPbToDomain(req.LinkingUserDoneText)
	result.ExternalUserNotFoundScreenText = text.ExternalUserNotFoundScreenTextPbToDomain(req.ExternalUserNotFoundText)
	result.SuccessLoginScreenText = text.SuccessLoginScreenTextPbToDomain(req.SuccessLoginText)
	result.LogoutDoneScreenText = text.LogoutDoneScreenTextPbToDomain(req.LogoutText)
	result.FooterText = text.FooterTextPbToDomain(req.FooterText)
	return result
}
