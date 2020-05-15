package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
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

func SendPasswordCodeCode(user *view_model.NotifyUser, code *es_model.PasswordCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.PasswordReset, &UrlData{UserID: user.ID, Code: codeString})
	if err != nil {
		return err
	}
	passwordCodeData := &PasswordCodeData{TemplateData: systemDefaults.Notifications.TemplateData.PasswordReset, FirstName: user.FirstName, LastName: user.LastName, URL: url}

	template, err := templates.GetParsedTemplate(passwordCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, template, systemDefaults.Notifications, false)
}
