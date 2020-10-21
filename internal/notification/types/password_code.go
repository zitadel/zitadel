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

type PasswordCodeData struct {
	templates.TemplateData
	FirstName string
	LastName  string
	URL       string
}

func SendPasswordCode(dir http.FileSystem, i18n *i18n.Translator, user *view_model.NotifyUser, code *es_model.PasswordCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *iam_model.LabelPolicyView) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.PasswordReset, &UrlData{UserID: user.ID, Code: codeString})
	if err != nil {
		return err
	}
	var args = map[string]interface{}{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Code":      codeString,
	}
	systemDefaults.Notifications.TemplateData.PasswordReset.Translate(i18n, args, user.PreferredLanguage)
	passwordCodeData := &PasswordCodeData{TemplateData: systemDefaults.Notifications.TemplateData.PasswordReset, FirstName: user.FirstName, LastName: user.LastName, URL: url}

	// Set the color in initCodeData
	passwordCodeData.PrimaryColor = colors.PrimaryColor
	passwordCodeData.SecondaryColor = colors.SecondaryColor
	template, err := templates.GetParsedTemplate(dir, passwordCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, systemDefaults.Notifications.TemplateData.PasswordReset.Subject, template, systemDefaults.Notifications, false)
}
