package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/i18n"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type InitCodeEmailData struct {
	templates.TemplateData
	URL string
}

type UrlData struct {
	UserID      string
	Code        string
	PasswordSet bool
}

func SendUserInitCode(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.InitUserCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *iam_model.LabelPolicyView, apiDomain string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.InitCode, &UrlData{UserID: user.ID, Code: codeString, PasswordSet: user.PasswordSet})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	initCodeData := &InitCodeEmailData{
		TemplateData: templates.GetTemplateData(translator, args, apiDomain, url, domain.InitCodeMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, initCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, initCodeData.Subject, template, systemDefaults.Notifications, true)
}
