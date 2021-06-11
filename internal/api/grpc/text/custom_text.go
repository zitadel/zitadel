package text

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	text_pb "github.com/caos/zitadel/pkg/grpc/text"
)

func ModelCustomMsgTextToPb(msg *model.MessageTextView) *text_pb.MessageCustomText {
	return &text_pb.MessageCustomText{
		Title:      msg.Title,
		PreHeader:  msg.PreHeader,
		Subject:    msg.Subject,
		Greeting:   msg.Greeting,
		Text:       msg.Text,
		ButtonText: msg.ButtonText,
		FooterText: msg.FooterText,
		Details: object.ToViewDetailsPb(
			msg.Sequence,
			msg.CreationDate,
			msg.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}

func SelectAccountScreenTextPbToDomain(text *text_pb.SelectAccountScreenText) domain.SelectAccountScreenText {
	if text == nil {
		return domain.SelectAccountScreenText{}
	}
	return domain.SelectAccountScreenText{
		Title:       text.Title,
		Description: text.Description,
		OtherUser:   text.OtherUser,
	}
}

func LoginScreenTextPbToDomain(text *text_pb.LoginScreenText) domain.LoginScreenText {
	if text == nil {
		return domain.LoginScreenText{}
	}
	return domain.LoginScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		LoginNameLabel:          text.NameLabel,
		RegisterButtonText:      text.RegisterButtonText,
		NextButtonText:          text.NextButtonText,
		ExternalUserDescription: text.ExternalUserDescription,
	}
}

func PasswordScreenTextPbToDomain(text *text_pb.PasswordScreenText) domain.PasswordScreenText {
	if text == nil {
		return domain.PasswordScreenText{}
	}
	return domain.PasswordScreenText{
		Title:          text.Title,
		Description:    text.Description,
		PasswordLabel:  text.PasswordLabel,
		ResetLinkText:  text.ResetLinkText,
		BackButtonText: text.BackButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func ResetPasswordScreenTextPbToDomain(text *text_pb.ResetPasswordScreenText) domain.ResetPasswordScreenText {
	if text == nil {
		return domain.ResetPasswordScreenText{}
	}
	return domain.ResetPasswordScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func InitializeUserScreenTextPbToDomain(text *text_pb.InitializeUserScreenText) domain.InitializeUserScreenText {
	if text == nil {
		return domain.InitializeUserScreenText{}
	}
	return domain.InitializeUserScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		CodeLabel:               text.CodeLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		ResendButtonText:        text.ResendButtonText,
		NextButtonText:          text.NextButtonText,
	}
}

func InitializeDoneScreenTextPbToDomain(text *text_pb.InitializeDoneScreenText) domain.InitializeDoneScreenText {
	if text == nil {
		return domain.InitializeDoneScreenText{}
	}
	return domain.InitializeDoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		AbortButtonText: text.AbortButtonText,
		NextButtonText:  text.NextButtonText,
	}
}

func InitMFAPromptScreenTextPbToDomain(text *text_pb.InitMFAPromptScreenText) domain.InitMFAPromptScreenText {
	if text == nil {
		return domain.InitMFAPromptScreenText{}
	}
	return domain.InitMFAPromptScreenText{
		Title:          text.Title,
		Description:    text.Description,
		OTPOption:      text.OtpOption,
		U2FOption:      text.U2FOption,
		SkipButtonText: text.SkipButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func InitMFAOTPScreenTextPbToDomain(text *text_pb.InitMFAOTPScreenText) domain.InitMFAOTPScreenText {
	if text == nil {
		return domain.InitMFAOTPScreenText{}
	}
	return domain.InitMFAOTPScreenText{
		Title:       text.Title,
		Description: text.Description,
		SecretLabel: text.SecretLabel,
		CodeLabel:   text.CodeLabel,
	}
}

func InitMFAU2FScreenTextPbToDomain(text *text_pb.InitMFAU2FScreenText) domain.InitMFAU2FScreenText {
	if text == nil {
		return domain.InitMFAU2FScreenText{}
	}
	return domain.InitMFAU2FScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		TokenNameLabel:          text.TokenNameLabel,
		RegisterTokenButtonText: text.RegisterTokenButtonText,
	}
}

func InitMFADoneScreenTextPbToDomain(text *text_pb.InitMFADoneScreenText) domain.InitMFADoneScreenText {
	if text == nil {
		return domain.InitMFADoneScreenText{}
	}
	return domain.InitMFADoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		AbortButtonText: text.AbortButtonText,
		NextButtonText:  text.NextButtonText,
	}
}

func VerifyMFAOTPScreenTextPbToDomain(text *text_pb.VerifyMFAOTPScreenText) domain.VerifyMFAOTPScreenText {
	if text == nil {
		return domain.VerifyMFAOTPScreenText{}
	}
	return domain.VerifyMFAOTPScreenText{
		Title:          text.Title,
		Description:    text.Description,
		CodeLabel:      text.CodeLabel,
		NextButtonText: text.NextButtonText,
	}
}

func VerifyMFAU2FScreenTextPbToDomain(text *text_pb.VerifyMFAU2FScreenText) domain.VerifyMFAU2FScreenText {
	if text == nil {
		return domain.VerifyMFAU2FScreenText{}
	}
	return domain.VerifyMFAU2FScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		ValidateTokenText:      text.ValidateTokenText,
		OtherOptionDescription: text.OtherOptionDescription,
		OptionOTPButtonText:    text.OptionOtpButtonText,
	}
}

func RegistrationOptionScreenTextPbToDomain(text *text_pb.RegistrationOptionScreenText) domain.RegistrationOptionScreenText {
	if text == nil {
		return domain.RegistrationOptionScreenText{}
	}
	return domain.RegistrationOptionScreenText{
		Title:                    text.Title,
		Description:              text.Description,
		UserNameButtonText:       text.UserNameButtonText,
		ExternalLoginDescription: text.ExternalLoginDescription,
	}
}

func RegistrationUserScreenTextPbToDomain(text *text_pb.RegistrationUserScreenText) domain.RegistrationUserScreenText {
	if text == nil {
		return domain.RegistrationUserScreenText{}
	}
	return domain.RegistrationUserScreenText{
		Title:                text.Title,
		Description:          text.Description,
		FirstnameLabel:       text.FirstnameLabel,
		LastnameLabel:        text.LastnameLabel,
		EmailLabel:           text.EmailLabel,
		LanguageLabel:        text.LanguageLabel,
		GenderLabel:          text.GenderLabel,
		PasswordLabel:        text.PasswordLabel,
		PasswordConfirmLabel: text.PasswordConfirmLabel,
		SaveButtonText:       text.SaveButtonText,
	}
}

func RegistrationOrgScreenTextPbToDomain(text *text_pb.RegistrationOrgScreenText) domain.RegistrationOrgScreenText {
	if text == nil {
		return domain.RegistrationOrgScreenText{}
	}
	return domain.RegistrationOrgScreenText{
		Title:                text.Title,
		Description:          text.Description,
		OrgNameLabel:         text.OrgnameLabel,
		FirstnameLabel:       text.FirstnameLabel,
		LastnameLabel:        text.LastnameLabel,
		EmailLabel:           text.EmailLabel,
		PasswordLabel:        text.PasswordLabel,
		PasswordConfirmLabel: text.PasswordConfirmLabel,
		SaveButtonText:       text.SaveButtonText,
	}
}

func PasswordlessScreenTextPbToDomain(text *text_pb.PasswordlessScreenText) domain.PasswordlessScreenText {
	if text == nil {
		return domain.PasswordlessScreenText{}
	}
	return domain.PasswordlessScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		LoginWithPwButtonText:   text.LoginWithPwButtonText,
		ValidateTokenButtonText: text.ValidateTokenButtonText,
	}
}

func SuccessLoginScreenTextPbToDomain(text *text_pb.SuccessLoginScreenText) domain.SuccessLoginScreenText {
	if text == nil {
		return domain.SuccessLoginScreenText{}
	}
	return domain.SuccessLoginScreenText{
		Title:                   text.Title,
		AutoRedirectDescription: text.AutoRedirectDescription,
		RedirectedDescription:   text.RedirectedDescription,
		NextButtonText:          text.NextButtonText,
	}
}
