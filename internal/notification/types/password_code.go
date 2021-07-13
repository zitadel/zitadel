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

type PasswordCodeData struct {
	templates.TemplateData
	FirstName string
	LastName  string
	URL       string
}

func SendPasswordCode(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.PasswordCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *iam_model.LabelPolicyView, apiDomain string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.PasswordReset, &UrlData{UserID: user.ID, Code: codeString})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	passwordResetData := &PasswordCodeData{
		TemplateData: templates.GetTemplateData(translator, args, apiDomain, url, domain.PasswordResetMessageType, user.PreferredLanguage, colors),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, passwordResetData)
	if err != nil {
		return err
	}
	if code.NotificationType == int32(domain.NotificationTypeSms) {
		return generateSms(user, passwordResetData.Text, systemDefaults.Notifications, false)
	}
	return generateEmail(user, passwordResetData.Subject, template, systemDefaults.Notifications, true)

}
