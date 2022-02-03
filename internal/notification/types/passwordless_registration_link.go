package types

import (
	"context"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/notification/templates"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/user"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type PasswordlessRegistrationLinkData struct {
	templates.TemplateData
	URL string
}

func SendPasswordlessRegistrationLink(ctx context.Context, mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *user.HumanPasswordlessInitCodeRequestedEvent, systemDefaults systemdefaults.SystemDefaults, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, apiDomain string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url := domain.PasswordlessInitCodeLink(systemDefaults.Notifications.Endpoints.PasswordlessRegistration, user.ID, user.ResourceOwner, code.ID, codeString)
	var args = mapNotifyUserToArgs(user)

	emailCodeData := &PasswordlessRegistrationLinkData{
		TemplateData: GetTemplateData(translator, args, apiDomain, url, domain.PasswordlessRegistrationMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}

	template, err := templates.GetParsedTemplate(mailhtml, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(ctx, user, emailCodeData.Subject, template, systemDefaults.Notifications, smtpConfig, true)
}
