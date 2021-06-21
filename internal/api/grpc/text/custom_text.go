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

func CustomLoginTextToPb(text *domain.CustomLoginText) *text_pb.LoginCustomText {
	return &text_pb.LoginCustomText{
		Details: object.ToViewDetailsPb(
			text.Sequence,
			text.CreationDate,
			text.ChangeDate,
			text.AggregateID,
		),
		SelectAccountText:      SelectAccountScreenToPb(text.SelectAccountScreenText),
		LoginText:              LoginScreenTextToPb(text.LoginScreenText),
		PasswordText:           PasswordScreenTextToPb(text.PasswordScreenText),
		ResetPasswordText:      ResetPasswordScreenTextToPb(text.InitPasswordScreenText),
		InitializeUserText:     InitializeUserScreenTextToPb(text.InitializeUserScreenText),
		InitializeDoneText:     InitializeDoneScreenTextToPb(text.InitializeDoneScreenText),
		InitMfaPromptText:      InitMFAPromptScreenTextToPb(text.InitMFAPromptScreenText),
		InitMfaOtpText:         InitMFAOTPScreenTextToPb(text.InitMFAOTPScreenText),
		InitMfaU2FText:         InitMFAU2FScreenTextToPb(text.InitMFAU2FScreenText),
		InitMfaDoneText:        InitMFADoneScreenTextToPb(text.InitMFADoneScreenText),
		VerifyMfaOtpText:       VerifyMFAOTPScreenTextToPb(text.VerifyMFAOTPScreenText),
		VerifyMfaU2FText:       VerifyMFAU2FScreenTextToPb(text.VerifyMFAU2FScreenText),
		RegistrationOptionText: RegistrationOptionScreenTextToPb(text.RegistrationOptionScreenText),
		RegistrationUserText:   RegistrationUserScreenTextToPb(text.RegistrationUserScreenText),
		RegistrationOrgText:    RegistrationOrgScreenTextToPb(text.RegistrationOrgScreenText),
		PasswordlessText:       PasswordlessScreenTextToPb(text.PasswordlessScreenText),
		SuccessLoginText:       SuccessLoginScreenTextToPb(text.SuccessLoginScreenText),
	}
}

func SelectAccountScreenToPb(text domain.SelectAccountScreenText) *text_pb.SelectAccountScreenText {
	return &text_pb.SelectAccountScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.Title,
		DescriptionLinkingProcess: text.Description,
		OtherUser:                 text.OtherUser,
		SesstionStateActive:       text.SessionStateActive,
		SesstionStateInactive:     text.SessionStateInactive,
		UserMustBeMemberOfOrg:     text.UserMustBeMemberOfOrg,
	}
}

func LoginScreenTextToPb(text domain.LoginScreenText) *text_pb.LoginScreenText {
	return &text_pb.LoginScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.TitleLinkingProcess,
		DescriptionLinkingProcess: text.DescriptionLinkingProcess,
		LoginNameLabel:            text.LoginNameLabel,
		RegisterButtonText:        text.RegisterButtonText,
		NextButtonText:            text.NextButtonText,
		ExternalUserDescription:   text.ExternalUserDescription,
		UserMustBeMemberOfOrg:     text.UserMustBeMemberOfOrg,
	}
}

func PasswordScreenTextToPb(text domain.PasswordScreenText) *text_pb.PasswordScreenText {
	return &text_pb.PasswordScreenText{
		Title:          text.Title,
		Description:    text.Description,
		PasswordLabel:  text.PasswordLabel,
		ResetLinkText:  text.ResetLinkText,
		BackButtonText: text.BackButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func ResetPasswordScreenTextToPb(text domain.InitPasswordScreenText) *text_pb.ResetPasswordScreenText {
	return &text_pb.ResetPasswordScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func InitializeUserScreenTextToPb(text domain.InitializeUserScreenText) *text_pb.InitializeUserScreenText {
	return &text_pb.InitializeUserScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		CodeLabel:               text.CodeLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		ResendButtonText:        text.ResendButtonText,
		NextButtonText:          text.NextButtonText,
	}
}

func InitializeDoneScreenTextToPb(text domain.InitializeDoneScreenText) *text_pb.InitializeDoneScreenText {
	return &text_pb.InitializeDoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		AbortButtonText: text.AbortButtonText,
		NextButtonText:  text.NextButtonText,
	}
}

func InitMFAPromptScreenTextToPb(text domain.InitMFAPromptScreenText) *text_pb.InitMFAPromptScreenText {
	return &text_pb.InitMFAPromptScreenText{
		Title:          text.Title,
		Description:    text.Description,
		OtpOption:      text.OTPOption,
		U2FOption:      text.U2FOption,
		SkipButtonText: text.SkipButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func InitMFAOTPScreenTextToPb(text domain.InitMFAOTPScreenText) *text_pb.InitMFAOTPScreenText {
	return &text_pb.InitMFAOTPScreenText{
		Title:       text.Title,
		Description: text.Description,
		SecretLabel: text.SecretLabel,
		CodeLabel:   text.CodeLabel,
	}
}

func InitMFAU2FScreenTextToPb(text domain.InitMFAU2FScreenText) *text_pb.InitMFAU2FScreenText {
	return &text_pb.InitMFAU2FScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		TokenNameLabel:          text.TokenNameLabel,
		RegisterTokenButtonText: text.RegisterTokenButtonText,
	}
}

func InitMFADoneScreenTextToPb(text domain.InitMFADoneScreenText) *text_pb.InitMFADoneScreenText {
	return &text_pb.InitMFADoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		AbortButtonText: text.AbortButtonText,
		NextButtonText:  text.NextButtonText,
	}
}

func VerifyMFAOTPScreenTextToPb(text domain.VerifyMFAOTPScreenText) *text_pb.VerifyMFAOTPScreenText {
	return &text_pb.VerifyMFAOTPScreenText{
		Title:          text.Title,
		Description:    text.Description,
		CodeLabel:      text.CodeLabel,
		NextButtonText: text.NextButtonText,
	}
}

func VerifyMFAU2FScreenTextToPb(text domain.VerifyMFAU2FScreenText) *text_pb.VerifyMFAU2FScreenText {
	return &text_pb.VerifyMFAU2FScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		ValidateTokenText:      text.ValidateTokenText,
		OtherOptionDescription: text.OtherOptionDescription,
		OptionOtpButtonText:    text.OptionOTPButtonText,
	}
}

func RegistrationOptionScreenTextToPb(text domain.RegistrationOptionScreenText) *text_pb.RegistrationOptionScreenText {
	return &text_pb.RegistrationOptionScreenText{
		Title:                    text.Title,
		Description:              text.Description,
		UserNameButtonText:       text.UserNameButtonText,
		ExternalLoginDescription: text.ExternalLoginDescription,
	}
}

func RegistrationUserScreenTextToPb(text domain.RegistrationUserScreenText) *text_pb.RegistrationUserScreenText {
	return &text_pb.RegistrationUserScreenText{
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

func RegistrationOrgScreenTextToPb(text domain.RegistrationOrgScreenText) *text_pb.RegistrationOrgScreenText {
	return &text_pb.RegistrationOrgScreenText{
		Title:                text.Title,
		Description:          text.Description,
		OrgnameLabel:         text.OrgNameLabel,
		FirstnameLabel:       text.FirstnameLabel,
		LastnameLabel:        text.LastnameLabel,
		EmailLabel:           text.EmailLabel,
		PasswordLabel:        text.PasswordLabel,
		PasswordConfirmLabel: text.PasswordConfirmLabel,
		SaveButtonText:       text.SaveButtonText,
	}
}

func PasswordlessScreenTextToPb(text domain.PasswordlessScreenText) *text_pb.PasswordlessScreenText {
	return &text_pb.PasswordlessScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		LoginWithPwButtonText:   text.LoginWithPwButtonText,
		ValidateTokenButtonText: text.ValidateTokenButtonText,
	}
}

func SuccessLoginScreenTextToPb(text domain.SuccessLoginScreenText) *text_pb.SuccessLoginScreenText {
	return &text_pb.SuccessLoginScreenText{
		Title:                   text.Title,
		AutoRedirectDescription: text.AutoRedirectDescription,
		RedirectedDescription:   text.RedirectedDescription,
		NextButtonText:          text.NextButtonText,
	}
}

func SelectAccountScreenTextPbToDomain(text *text_pb.SelectAccountScreenText) domain.SelectAccountScreenText {
	if text == nil {
		return domain.SelectAccountScreenText{}
	}
	return domain.SelectAccountScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.TitleLinkingProcess,
		DescriptionLinkingProcess: text.DescriptionLinkingProcess,
		OtherUser:                 text.OtherUser,
		SessionStateActive:        text.SesstionStateActive,
		SessionStateInactive:      text.SesstionStateInactive,
		UserMustBeMemberOfOrg:     text.UserMustBeMemberOfOrg,
	}
}

func LoginScreenTextPbToDomain(text *text_pb.LoginScreenText) domain.LoginScreenText {
	if text == nil {
		return domain.LoginScreenText{}
	}
	return domain.LoginScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.TitleLinkingProcess,
		DescriptionLinkingProcess: text.DescriptionLinkingProcess,
		LoginNameLabel:            text.LoginNameLabel,
		RegisterButtonText:        text.RegisterButtonText,
		NextButtonText:            text.NextButtonText,
		ExternalUserDescription:   text.ExternalUserDescription,
		UserMustBeMemberOfOrg:     text.UserMustBeMemberOfOrg,
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

func ResetPasswordScreenTextPbToDomain(text *text_pb.ResetPasswordScreenText) domain.InitPasswordScreenText {
	if text == nil {
		return domain.InitPasswordScreenText{}
	}
	return domain.InitPasswordScreenText{
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
