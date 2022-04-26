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
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type PasswordCodeData struct {
	templates.TemplateData
	FirstName string
	LastName  string
	URL       string
}

func SendPasswordCode(ctx context.Context, mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.PasswordCode, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), getTwilioConfig func(ctx context.Context) (*twilio.TwilioConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, assetsPrefix string, origin string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url := login.InitPasswordLink(origin, user.ID, codeString)
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	passwordResetData := &PasswordCodeData{
		TemplateData: GetTemplateData(translator, args, assetsPrefix, url, domain.PasswordResetMessageType, user.PreferredLanguage, colors),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, passwordResetData)
	if err != nil {
		return err
	}
	if code.NotificationType == int32(domain.NotificationTypeSms) {
		return generateSms(ctx, user, passwordResetData.Text, getTwilioConfig, getFileSystemProvider, getLogProvider, false)
	}
	return generateEmail(ctx, user, passwordResetData.Subject, template, smtpConfig, getFileSystemProvider, getLogProvider, true)

}
