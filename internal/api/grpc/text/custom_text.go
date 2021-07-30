package text

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	text_pb "github.com/caos/zitadel/pkg/grpc/text"
)

func DomainCustomMsgTextToPb(msg *domain.CustomMessageText) *text_pb.MessageCustomText {
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
		SelectAccountText:                SelectAccountScreenToPb(text.SelectAccount),
		LoginText:                        LoginScreenTextToPb(text.Login),
		PasswordText:                     PasswordScreenTextToPb(text.Password),
		UsernameChangeText:               UsernameChangeScreenTextToPb(text.UsernameChange),
		UsernameChangeDoneText:           UsernameChangeDoneScreenTextToPb(text.UsernameChangeDone),
		InitPasswordText:                 InitPasswordScreenTextToPb(text.InitPassword),
		InitPasswordDoneText:             InitPasswordDoneScreenTextToPb(text.InitPasswordDone),
		EmailVerificationText:            EmailVerificationScreenTextToPb(text.EmailVerification),
		EmailVerificationDoneText:        EmailVerificationDoneScreenTextToPb(text.EmailVerificationDone),
		InitializeUserText:               InitializeUserScreenTextToPb(text.InitUser),
		InitializeDoneText:               InitializeUserDoneScreenTextToPb(text.InitUserDone),
		InitMfaPromptText:                InitMFAPromptScreenTextToPb(text.InitMFAPrompt),
		InitMfaOtpText:                   InitMFAOTPScreenTextToPb(text.InitMFAOTP),
		InitMfaU2FText:                   InitMFAU2FScreenTextToPb(text.InitMFAU2F),
		InitMfaDoneText:                  InitMFADoneScreenTextToPb(text.InitMFADone),
		MfaProvidersText:                 MFAProvidersTextToPb(text.MFAProvider),
		VerifyMfaOtpText:                 VerifyMFAOTPScreenTextToPb(text.VerifyMFAOTP),
		VerifyMfaU2FText:                 VerifyMFAU2FScreenTextToPb(text.VerifyMFAU2F),
		PasswordlessText:                 PasswordlessScreenTextToPb(text.Passwordless),
		PasswordlessRegistrationText:     PasswordlessRegistrationScreenTextToPb(text.PasswordlessRegistration),
		PasswordlessRegistrationDoneText: PasswordlessRegistrationDoneScreenTextToPb(text.PasswordlessRegistrationDone),
		PasswordChangeText:               PasswordChangeScreenTextToPb(text.PasswordChange),
		PasswordChangeDoneText:           PasswordChangeDoneScreenTextToPb(text.PasswordChangeDone),
		PasswordResetDoneText:            PasswordResetDoneScreenTextToPb(text.PasswordResetDone),
		RegistrationOptionText:           RegistrationOptionScreenTextToPb(text.RegisterOption),
		RegistrationUserText:             RegistrationUserScreenTextToPb(text.RegistrationUser),
		RegistrationOrgText:              RegistrationOrgScreenTextToPb(text.RegistrationOrg),
		LinkingUserDoneText:              LinkingUserDoneScreenTextToPb(text.LinkingUsersDone),
		ExternalUserNotFoundText:         ExternalUserNotFoundScreenTextToPb(text.ExternalNotFoundOption),
		SuccessLoginText:                 SuccessLoginScreenTextToPb(text.LoginSuccess),
		LogoutText:                       LogoutDoneScreenTextToPb(text.LogoutDone),
		FooterText:                       FooterTextToPb(text.Footer),
	}
}

func SelectAccountScreenToPb(text domain.SelectAccountScreenText) *text_pb.SelectAccountScreenText {
	return &text_pb.SelectAccountScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.TitleLinking,
		DescriptionLinkingProcess: text.DescriptionLinking,
		OtherUser:                 text.OtherUser,
		SessionStateActive:        text.SessionState0,
		SessionStateInactive:      text.SessionState1,
		UserMustBeMemberOfOrg:     text.MustBeMemberOfOrg,
	}
}

func LoginScreenTextToPb(text domain.LoginScreenText) *text_pb.LoginScreenText {
	return &text_pb.LoginScreenText{
		Title:                     text.Title,
		Description:               text.Description,
		TitleLinkingProcess:       text.TitleLinking,
		DescriptionLinkingProcess: text.DescriptionLinking,
		LoginNameLabel:            text.LoginNameLabel,
		UserNamePlaceholder:       text.UsernamePlaceholder,
		LoginNamePlaceholder:      text.LoginnamePlaceholder,
		RegisterButtonText:        text.RegisterButtonText,
		NextButtonText:            text.NextButtonText,
		ExternalUserDescription:   text.ExternalUserDescription,
		UserMustBeMemberOfOrg:     text.MustBeMemberOfOrg,
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
		MinLength:      text.MinLength,
		HasUppercase:   text.HasUppercase,
		HasLowercase:   text.HasLowercase,
		HasNumber:      text.HasNumber,
		HasSymbol:      text.HasSymbol,
		Confirmation:   text.Confirmation,
	}
}

func UsernameChangeScreenTextToPb(text domain.UsernameChangeScreenText) *text_pb.UsernameChangeScreenText {
	return &text_pb.UsernameChangeScreenText{
		Title:            text.Title,
		Description:      text.Description,
		UsernameLabel:    text.UsernameLabel,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func UsernameChangeDoneScreenTextToPb(text domain.UsernameChangeDoneScreenText) *text_pb.UsernameChangeDoneScreenText {
	return &text_pb.UsernameChangeDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func InitPasswordScreenTextToPb(text domain.InitPasswordScreenText) *text_pb.InitPasswordScreenText {
	return &text_pb.InitPasswordScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		CodeLabel:               text.CodeLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		NextButtonText:          text.NextButtonText,
		ResendButtonText:        text.ResendButtonText,
	}
}

func InitPasswordDoneScreenTextToPb(text domain.InitPasswordDoneScreenText) *text_pb.InitPasswordDoneScreenText {
	return &text_pb.InitPasswordDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
	}
}

func EmailVerificationScreenTextToPb(text domain.EmailVerificationScreenText) *text_pb.EmailVerificationScreenText {
	return &text_pb.EmailVerificationScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CodeLabel:        text.CodeLabel,
		NextButtonText:   text.NextButtonText,
		ResendButtonText: text.ResendButtonText,
	}
}

func EmailVerificationDoneScreenTextToPb(text domain.EmailVerificationDoneScreenText) *text_pb.EmailVerificationDoneScreenText {
	return &text_pb.EmailVerificationDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
		LoginButtonText:  text.LoginButtonText,
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

func InitializeUserDoneScreenTextToPb(text domain.InitializeUserDoneScreenText) *text_pb.InitializeUserDoneScreenText {
	return &text_pb.InitializeUserDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func InitMFAPromptScreenTextToPb(text domain.InitMFAPromptScreenText) *text_pb.InitMFAPromptScreenText {
	return &text_pb.InitMFAPromptScreenText{
		Title:          text.Title,
		Description:    text.Description,
		OtpOption:      text.Provider0,
		U2FOption:      text.Provider1,
		SkipButtonText: text.SkipButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func InitMFAOTPScreenTextToPb(text domain.InitMFAOTPScreenText) *text_pb.InitMFAOTPScreenText {
	return &text_pb.InitMFAOTPScreenText{
		Title:            text.Title,
		Description:      text.Description,
		DescriptionOtp:   text.OTPDescription,
		SecretLabel:      text.SecretLabel,
		CodeLabel:        text.CodeLabel,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
	}
}

func InitMFAU2FScreenTextToPb(text domain.InitMFAU2FScreenText) *text_pb.InitMFAU2FScreenText {
	return &text_pb.InitMFAU2FScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		TokenNameLabel:          text.TokenNameLabel,
		RegisterTokenButtonText: text.RegisterTokenButtonText,
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func InitMFADoneScreenTextToPb(text domain.InitMFADoneScreenText) *text_pb.InitMFADoneScreenText {
	return &text_pb.InitMFADoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func MFAProvidersTextToPb(text domain.MFAProvidersText) *text_pb.MFAProvidersText {
	return &text_pb.MFAProvidersText{
		ChooseOther: text.ChooseOther,
		Otp:         text.Provider0,
		U2F:         text.Provider1,
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
		Title:             text.Title,
		Description:       text.Description,
		ValidateTokenText: text.ValidateTokenButtonText,
		NotSupported:      text.NotSupported,
		ErrorRetry:        text.ErrorRetry,
	}
}

func PasswordlessScreenTextToPb(text domain.PasswordlessScreenText) *text_pb.PasswordlessScreenText {
	return &text_pb.PasswordlessScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		LoginWithPwButtonText:   text.LoginWithPwButtonText,
		ValidateTokenButtonText: text.ValidateTokenButtonText,
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func PasswordlessRegistrationScreenTextToPb(text domain.PasswordlessRegistrationScreenText) *text_pb.PasswordlessRegistrationScreenText {
	return &text_pb.PasswordlessRegistrationScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		RegisterTokenButtonText: text.RegisterTokenButtonText,
		TokenNameLabel:          text.TokenNameLabel,
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func PasswordlessRegistrationDoneScreenTextToPb(text domain.PasswordlessRegistrationDoneScreenText) *text_pb.PasswordlessRegistrationDoneScreenText {
	return &text_pb.PasswordlessRegistrationDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func PasswordChangeScreenTextToPb(text domain.PasswordChangeScreenText) *text_pb.PasswordChangeScreenText {
	return &text_pb.PasswordChangeScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		OldPasswordLabel:        text.OldPasswordLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		CancelButtonText:        text.CancelButtonText,
		NextButtonText:          text.NextButtonText,
	}
}

func PasswordChangeDoneScreenTextToPb(text domain.PasswordChangeDoneScreenText) *text_pb.PasswordChangeDoneScreenText {
	return &text_pb.PasswordChangeDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func PasswordResetDoneScreenTextToPb(text domain.PasswordResetDoneScreenText) *text_pb.PasswordResetDoneScreenText {
	return &text_pb.PasswordResetDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func RegistrationOptionScreenTextToPb(text domain.RegistrationOptionScreenText) *text_pb.RegistrationOptionScreenText {
	return &text_pb.RegistrationOptionScreenText{
		Title:                    text.Title,
		Description:              text.Description,
		UserNameButtonText:       text.RegisterUsernamePasswordButtonText,
		ExternalLoginDescription: text.ExternalLoginDescription,
	}
}

func RegistrationUserScreenTextToPb(text domain.RegistrationUserScreenText) *text_pb.RegistrationUserScreenText {
	return &text_pb.RegistrationUserScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		DescriptionOrgRegister: text.DescriptionOrgRegister,
		FirstnameLabel:         text.FirstnameLabel,
		LastnameLabel:          text.LastnameLabel,
		EmailLabel:             text.EmailLabel,
		UsernameLabel:          text.UsernameLabel,
		LanguageLabel:          text.LanguageLabel,
		GenderLabel:            text.GenderLabel,
		PasswordLabel:          text.PasswordLabel,
		PasswordConfirmLabel:   text.PasswordConfirmLabel,
		TosAndPrivacyLabel:     text.TOSAndPrivacyLabel,
		TosConfirm:             text.TOSConfirm,
		TosLinkText:            text.TOSLinkText,
		TosConfirmAnd:          text.TOSConfirmAnd,
		PrivacyLinkText:        text.PrivacyLinkText,
		NextButtonText:         text.NextButtonText,
		BackButtonText:         text.BackButtonText,
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
		UsernameLabel:        text.UsernameLabel,
		PasswordLabel:        text.PasswordLabel,
		PasswordConfirmLabel: text.PasswordConfirmLabel,
		TosAndPrivacyLabel:   text.TOSAndPrivacyLabel,
		TosConfirm:           text.TOSConfirm,
		TosLinkText:          text.TOSLinkText,
		TosConfirmAnd:        text.TOSConfirmAnd,
		PrivacyLinkText:      text.PrivacyLinkText,
		SaveButtonText:       text.SaveButtonText,
	}
}

func LinkingUserDoneScreenTextToPb(text domain.LinkingUserDoneScreenText) *text_pb.LinkingUserDoneScreenText {
	return &text_pb.LinkingUserDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func ExternalUserNotFoundScreenTextToPb(text domain.ExternalUserNotFoundScreenText) *text_pb.ExternalUserNotFoundScreenText {
	return &text_pb.ExternalUserNotFoundScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		LinkButtonText:         text.LinkButtonText,
		AutoRegisterButtonText: text.AutoRegisterButtonText,
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

func LogoutDoneScreenTextToPb(text domain.LogoutDoneScreenText) *text_pb.LogoutDoneScreenText {
	return &text_pb.LogoutDoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		LoginButtonText: text.LoginButtonText,
	}
}

func FooterTextToPb(text domain.FooterText) *text_pb.FooterText {
	return &text_pb.FooterText{
		Tos:           text.TOS,
		PrivacyPolicy: text.PrivacyPolicy,
		Help:          text.Help,
		HelpLink:      text.HelpLink,
	}
}

func SelectAccountScreenTextPbToDomain(text *text_pb.SelectAccountScreenText) domain.SelectAccountScreenText {
	if text == nil {
		return domain.SelectAccountScreenText{}
	}
	return domain.SelectAccountScreenText{
		Title:              text.Title,
		Description:        text.Description,
		TitleLinking:       text.TitleLinkingProcess,
		DescriptionLinking: text.DescriptionLinkingProcess,
		OtherUser:          text.OtherUser,
		SessionState0:      text.SessionStateActive,
		SessionState1:      text.SessionStateInactive,
		MustBeMemberOfOrg:  text.UserMustBeMemberOfOrg,
	}
}

func LoginScreenTextPbToDomain(text *text_pb.LoginScreenText) domain.LoginScreenText {
	if text == nil {
		return domain.LoginScreenText{}
	}
	return domain.LoginScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		TitleLinking:            text.TitleLinkingProcess,
		DescriptionLinking:      text.DescriptionLinkingProcess,
		LoginNameLabel:          text.LoginNameLabel,
		UsernamePlaceholder:     text.UserNamePlaceholder,
		LoginnamePlaceholder:    text.LoginNamePlaceholder,
		RegisterButtonText:      text.RegisterButtonText,
		NextButtonText:          text.NextButtonText,
		ExternalUserDescription: text.ExternalUserDescription,
		MustBeMemberOfOrg:       text.UserMustBeMemberOfOrg,
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
		MinLength:      text.MinLength,
		HasUppercase:   text.HasUppercase,
		HasLowercase:   text.HasLowercase,
		HasNumber:      text.HasNumber,
		HasSymbol:      text.HasSymbol,
		Confirmation:   text.Confirmation,
	}
}

func UsernameChangeScreenTextPbToDomain(text *text_pb.UsernameChangeScreenText) domain.UsernameChangeScreenText {
	if text == nil {
		return domain.UsernameChangeScreenText{}
	}
	return domain.UsernameChangeScreenText{
		Title:            text.Title,
		Description:      text.Description,
		UsernameLabel:    text.UsernameLabel,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func UsernameChangeDoneScreenTextPbToDomain(text *text_pb.UsernameChangeDoneScreenText) domain.UsernameChangeDoneScreenText {
	if text == nil {
		return domain.UsernameChangeDoneScreenText{}
	}
	return domain.UsernameChangeDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func InitPasswordScreenTextPbToDomain(text *text_pb.InitPasswordScreenText) domain.InitPasswordScreenText {
	if text == nil {
		return domain.InitPasswordScreenText{}
	}
	return domain.InitPasswordScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		CodeLabel:               text.CodeLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		NextButtonText:          text.NextButtonText,
		ResendButtonText:        text.ResendButtonText,
	}
}

func InitPasswordDoneScreenTextPbToDomain(text *text_pb.InitPasswordDoneScreenText) domain.InitPasswordDoneScreenText {
	if text == nil {
		return domain.InitPasswordDoneScreenText{}
	}
	return domain.InitPasswordDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
	}
}

func EmailVerificationScreenTextPbToDomain(text *text_pb.EmailVerificationScreenText) domain.EmailVerificationScreenText {
	if text == nil {
		return domain.EmailVerificationScreenText{}
	}
	return domain.EmailVerificationScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CodeLabel:        text.CodeLabel,
		NextButtonText:   text.NextButtonText,
		ResendButtonText: text.ResendButtonText,
	}
}

func EmailVerificationDoneScreenTextPbToDomain(text *text_pb.EmailVerificationDoneScreenText) domain.EmailVerificationDoneScreenText {
	if text == nil {
		return domain.EmailVerificationDoneScreenText{}
	}
	return domain.EmailVerificationDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
		LoginButtonText:  text.LoginButtonText,
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

func InitializeDoneScreenTextPbToDomain(text *text_pb.InitializeUserDoneScreenText) domain.InitializeUserDoneScreenText {
	if text == nil {
		return domain.InitializeUserDoneScreenText{}
	}
	return domain.InitializeUserDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func InitMFAPromptScreenTextPbToDomain(text *text_pb.InitMFAPromptScreenText) domain.InitMFAPromptScreenText {
	if text == nil {
		return domain.InitMFAPromptScreenText{}
	}
	return domain.InitMFAPromptScreenText{
		Title:          text.Title,
		Description:    text.Description,
		Provider0:      text.OtpOption,
		Provider1:      text.U2FOption,
		SkipButtonText: text.SkipButtonText,
		NextButtonText: text.NextButtonText,
	}
}

func InitMFAOTPScreenTextPbToDomain(text *text_pb.InitMFAOTPScreenText) domain.InitMFAOTPScreenText {
	if text == nil {
		return domain.InitMFAOTPScreenText{}
	}
	return domain.InitMFAOTPScreenText{
		Title:            text.Title,
		Description:      text.Description,
		OTPDescription:   text.DescriptionOtp,
		SecretLabel:      text.SecretLabel,
		CodeLabel:        text.CodeLabel,
		NextButtonText:   text.NextButtonText,
		CancelButtonText: text.CancelButtonText,
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
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func InitMFADoneScreenTextPbToDomain(text *text_pb.InitMFADoneScreenText) domain.InitMFADoneScreenText {
	if text == nil {
		return domain.InitMFADoneScreenText{}
	}
	return domain.InitMFADoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func MFAProvidersTextPbToDomain(text *text_pb.MFAProvidersText) domain.MFAProvidersText {
	if text == nil {
		return domain.MFAProvidersText{}
	}
	return domain.MFAProvidersText{
		ChooseOther: text.ChooseOther,
		Provider0:   text.Otp,
		Provider1:   text.U2F,
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
		Title:                   text.Title,
		Description:             text.Description,
		ValidateTokenButtonText: text.ValidateTokenText,
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
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
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func PasswordlessPromptScreenTextPbToDomain(text *text_pb.PasswordlessPromptScreenText) domain.PasswordlessPromptScreenText {
	if text == nil {
		return domain.PasswordlessPromptScreenText{}
	}
	return domain.PasswordlessPromptScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		DescriptionInit:        text.DescriptionInit,
		PasswordlessButtonText: text.PasswordlessButtonText,
		NextButtonText:         text.NextButtonText,
		SkipButtonText:         text.SkipButtonText,
	}
}

func PasswordlessRegistrationScreenTextPbToDomain(text *text_pb.PasswordlessRegistrationScreenText) domain.PasswordlessRegistrationScreenText {
	if text == nil {
		return domain.PasswordlessRegistrationScreenText{}
	}
	return domain.PasswordlessRegistrationScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		RegisterTokenButtonText: text.RegisterTokenButtonText,
		TokenNameLabel:          text.TokenNameLabel,
		NotSupported:            text.NotSupported,
		ErrorRetry:              text.ErrorRetry,
	}
}

func PasswordlessRegistrationDoneScreenTextPbToDomain(text *text_pb.PasswordlessRegistrationDoneScreenText) domain.PasswordlessRegistrationDoneScreenText {
	if text == nil {
		return domain.PasswordlessRegistrationDoneScreenText{}
	}
	return domain.PasswordlessRegistrationDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func PasswordChangeScreenTextPbToDomain(text *text_pb.PasswordChangeScreenText) domain.PasswordChangeScreenText {
	if text == nil {
		return domain.PasswordChangeScreenText{}
	}
	return domain.PasswordChangeScreenText{
		Title:                   text.Title,
		Description:             text.Description,
		OldPasswordLabel:        text.OldPasswordLabel,
		NewPasswordLabel:        text.NewPasswordLabel,
		NewPasswordConfirmLabel: text.NewPasswordConfirmLabel,
		CancelButtonText:        text.CancelButtonText,
		NextButtonText:          text.NextButtonText,
	}
}

func PasswordChangeDoneScreenTextPbToDomain(text *text_pb.PasswordChangeDoneScreenText) domain.PasswordChangeDoneScreenText {
	if text == nil {
		return domain.PasswordChangeDoneScreenText{}
	}
	return domain.PasswordChangeDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func PasswordResetDoneScreenTextPbToDomain(text *text_pb.PasswordResetDoneScreenText) domain.PasswordResetDoneScreenText {
	if text == nil {
		return domain.PasswordResetDoneScreenText{}
	}
	return domain.PasswordResetDoneScreenText{
		Title:          text.Title,
		Description:    text.Description,
		NextButtonText: text.NextButtonText,
	}
}

func RegistrationOptionScreenTextPbToDomain(text *text_pb.RegistrationOptionScreenText) domain.RegistrationOptionScreenText {
	if text == nil {
		return domain.RegistrationOptionScreenText{}
	}
	return domain.RegistrationOptionScreenText{
		Title:                              text.Title,
		Description:                        text.Description,
		RegisterUsernamePasswordButtonText: text.UserNameButtonText,
		ExternalLoginDescription:           text.ExternalLoginDescription,
	}
}

func RegistrationUserScreenTextPbToDomain(text *text_pb.RegistrationUserScreenText) domain.RegistrationUserScreenText {
	if text == nil {
		return domain.RegistrationUserScreenText{}
	}
	return domain.RegistrationUserScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		DescriptionOrgRegister: text.DescriptionOrgRegister,
		FirstnameLabel:         text.FirstnameLabel,
		LastnameLabel:          text.LastnameLabel,
		EmailLabel:             text.EmailLabel,
		UsernameLabel:          text.UsernameLabel,
		LanguageLabel:          text.LanguageLabel,
		GenderLabel:            text.GenderLabel,
		PasswordLabel:          text.PasswordLabel,
		PasswordConfirmLabel:   text.PasswordConfirmLabel,
		TOSAndPrivacyLabel:     text.TosAndPrivacyLabel,
		TOSConfirm:             text.TosConfirm,
		TOSLinkText:            text.TosLinkText,
		TOSConfirmAnd:          text.TosConfirmAnd,
		PrivacyLinkText:        text.PrivacyLinkText,
		NextButtonText:         text.NextButtonText,
		BackButtonText:         text.BackButtonText,
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
		UsernameLabel:        text.UsernameLabel,
		PasswordLabel:        text.PasswordLabel,
		PasswordConfirmLabel: text.PasswordConfirmLabel,
		TOSAndPrivacyLabel:   text.TosAndPrivacyLabel,
		TOSConfirm:           text.TosConfirm,
		TOSLinkText:          text.TosLinkText,
		TOSConfirmAnd:        text.TosConfirmAnd,
		PrivacyLinkText:      text.PrivacyLinkText,
		SaveButtonText:       text.SaveButtonText,
	}
}

func LinkingUserDoneScreenTextPbToDomain(text *text_pb.LinkingUserDoneScreenText) domain.LinkingUserDoneScreenText {
	if text == nil {
		return domain.LinkingUserDoneScreenText{}
	}
	return domain.LinkingUserDoneScreenText{
		Title:            text.Title,
		Description:      text.Description,
		CancelButtonText: text.CancelButtonText,
		NextButtonText:   text.NextButtonText,
	}
}

func ExternalUserNotFoundScreenTextPbToDomain(text *text_pb.ExternalUserNotFoundScreenText) domain.ExternalUserNotFoundScreenText {
	if text == nil {
		return domain.ExternalUserNotFoundScreenText{}
	}
	return domain.ExternalUserNotFoundScreenText{
		Title:                  text.Title,
		Description:            text.Description,
		LinkButtonText:         text.LinkButtonText,
		AutoRegisterButtonText: text.AutoRegisterButtonText,
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

func LogoutDoneScreenTextPbToDomain(text *text_pb.LogoutDoneScreenText) domain.LogoutDoneScreenText {
	if text == nil {
		return domain.LogoutDoneScreenText{}
	}
	return domain.LogoutDoneScreenText{
		Title:           text.Title,
		Description:     text.Description,
		LoginButtonText: text.LoginButtonText,
	}
}

func FooterTextPbToDomain(text *text_pb.FooterText) domain.FooterText {
	if text == nil {
		return domain.FooterText{}
	}
	return domain.FooterText{
		TOS:           text.Tos,
		PrivacyPolicy: text.PrivacyPolicy,
		Help:          text.Help,
		HelpLink:      text.HelpLink,
	}
}
