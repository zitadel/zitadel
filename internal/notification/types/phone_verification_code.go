package types

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/templates"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	UserID string
}

func SendPhoneVerificationCode(ctx context.Context, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.PhoneCode, getTwilioConfig func(ctx context.Context) (*twilio.TwilioConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	text := translator.Localize(fmt.Sprintf("%s.%s", domain.VerifyPhoneMessageType, domain.MessageText), args, user.PreferredLanguage)

	codeData := &PhoneVerificationCodeData{UserID: user.ID}
	template, err := templates.ParseTemplateText(text, codeData)
	if err != nil {
		return err
	}
	return generateSms(ctx, user, template, getTwilioConfig, getFileSystemProvider, getLogProvider, true)
}
