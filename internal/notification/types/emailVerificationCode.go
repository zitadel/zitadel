package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type EmailVerificationCodeData struct {
	templates.TemplateData
	FirstName string
	LastName  string
	URL       string
}

func SendEmailVerificationCode(user *view_model.NotifyUser, code *es_model.EmailCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.VerifyEmail, &UrlData{UserID: user.ID, Code: codeString})
	if err != nil {
		return err
	}
	emailCodeData := &EmailVerificationCodeData{TemplateData: systemDefaults.Notifications.TemplateData.VerifyEmail, FirstName: user.FirstName, LastName: user.LastName, URL: url}

	template, err := templates.GetParsedTemplate(emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, template, systemDefaults.Notifications, true)
}
