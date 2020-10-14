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

type EmailVerificationCodeData struct {
	templates.TemplateData
	URL string
}

func SendEmailVerificationCode(dir http.FileSystem, i18n *i18n.Translator, user *view_model.NotifyUser, code *es_model.EmailCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *iam_model.LabelPolicyView) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.VerifyEmail, &UrlData{UserID: user.ID, Code: codeString})
	if err != nil {
		return err
	}
	var args = map[string]interface{}{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Code":      codeString,
	}
	systemDefaults.Notifications.TemplateData.VerifyEmail.Translate(i18n, args, user.PreferredLanguage)
	emailCodeData := &EmailVerificationCodeData{TemplateData: systemDefaults.Notifications.TemplateData.VerifyEmail, URL: url}

	// Set the color in initCodeData
	emailCodeData.PrimaryColor = colors.PrimaryColor
	emailCodeData.SecondaryColor = colors.SecondaryColor
	template, err := templates.GetParsedTemplate(dir, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, systemDefaults.Notifications.TemplateData.VerifyEmail.Subject, template, systemDefaults.Notifications, true)
}
