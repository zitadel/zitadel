package types

import (
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type InitCodeEmailData struct {
	templates.TemplateData
	URL string
}

type UrlData struct {
	UserID      string
	Code        string
	PasswordSet bool
	OrgID       string
}

func SendUserInitCode(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.InitUserCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, apiDomain string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.InitCode, &UrlData{UserID: user.ID, Code: codeString, PasswordSet: user.PasswordSet, OrgID: user.ResourceOwner})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	initCodeData := &InitCodeEmailData{
		TemplateData: GetTemplateData(translator, args, apiDomain, url, domain.InitCodeMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, initCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, initCodeData.Subject, template, systemDefaults.Notifications, true)
}
