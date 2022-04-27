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
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type PasswordlessRegistrationLinkData struct {
	templates.TemplateData
	URL string
}

func SendPasswordlessRegistrationLink(ctx context.Context, mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *user.HumanPasswordlessInitCodeRequestedEvent, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, assetsPrefix string, origin string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url := domain.PasswordlessInitCodeLink(origin+login.HandlerPrefix+login.EndpointPasswordlessRegistration, user.ID, user.ResourceOwner, code.ID, codeString)
	var args = mapNotifyUserToArgs(user)

	emailCodeData := &PasswordlessRegistrationLinkData{
		TemplateData: GetTemplateData(translator, args, assetsPrefix, url, domain.PasswordlessRegistrationMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}

	template, err := templates.GetParsedTemplate(mailhtml, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(ctx, user, emailCodeData.Subject, template, smtpConfig, getFileSystemProvider, getLogProvider, true)
}
