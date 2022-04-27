package types

import (
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type PasswordlessRegistrationLinkData struct {
	templates.TemplateData
	URL string
}

func SendPasswordlessRegistrationLink(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *user.HumanPasswordlessInitCodeRequestedEvent, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, apiDomain string) error {
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
	return generateEmail(user, emailCodeData.Subject, template, systemDefaults.Notifications, true)
}
