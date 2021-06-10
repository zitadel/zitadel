package management

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func SetInitCustomTextToDomain(msg *mgmt_pb.SetCustomInitMessageTextRequest) *domain.CustomMessageText {
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

func SetPasswordResetCustomTextToDomain(msg *mgmt_pb.SetCustomPasswordResetMessageTextRequest) *domain.CustomMessageText {
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

func SetVerifyEmailCustomTextToDomain(msg *mgmt_pb.SetCustomVerifyEmailMessageTextRequest) *domain.CustomMessageText {
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

func SetVerifyPhoneCustomTextToDomain(msg *mgmt_pb.SetCustomVerifyPhoneMessageTextRequest) *domain.CustomMessageText {
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

func SetDomainClaimedCustomTextToDomain(msg *mgmt_pb.SetCustomDomainClaimedMessageTextRequest) *domain.CustomMessageText {
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
