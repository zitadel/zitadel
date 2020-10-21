package types

import (
	"net/http"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
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

func SendUserInitCode(dir http.FileSystem, i18n *i18n.Translator, user *view_model.NotifyUser, code *es_model.InitUserCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *iam_model.LabelPolicyView) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.InitCode, &UrlData{UserID: user.ID, Code: codeString, PasswordSet: user.PasswordSet})
	if err != nil {
		return err
	}
	var args = map[string]interface{}{
		"FirstName":          user.FirstName,
		"LastName":           user.LastName,
		"Code":               codeString,
		"PreferredLoginName": user.PreferredLoginName,
	}
	systemDefaults.Notifications.TemplateData.InitCode.Translate(i18n, args, user.PreferredLanguage)
	initCodeData := &InitCodeEmailData{TemplateData: systemDefaults.Notifications.TemplateData.InitCode, URL: url}

	// Set the color in initCodeData
	initCodeData.PrimaryColor = colors.PrimaryColor
	initCodeData.SecondaryColor = colors.SecondaryColor
	template, err := templates.GetParsedTemplate(dir, initCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, systemDefaults.Notifications.TemplateData.InitCode.Subject, template, systemDefaults.Notifications, true)
}
