package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type EmailVerificationCodeData struct {
	templates.TemplateData
	URL string
}

func SendEmailVerificationCode(ctx context.Context, mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.EmailCode, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, assetsPrefix string, origin string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url := login.MailVerificationLink(origin, user.ID, codeString)
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	emailCodeData := &EmailVerificationCodeData{
		TemplateData: GetTemplateData(translator, args, assetsPrefix, url, domain.VerifyEmailMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}

	template, err := templates.GetParsedTemplate(mailhtml, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(ctx, user, emailCodeData.Subject, template, smtpConfig, getFileSystemProvider, getLogProvider, true)
}
