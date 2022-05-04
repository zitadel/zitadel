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

type PasswordCodeData struct {
	templates.TemplateData
	FirstName string
	LastName  string
	URL       string
}

func SendPasswordCode(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.PasswordCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm, colors *query.LabelPolicy, apiDomain string) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.PasswordReset, &UrlData{UserID: user.ID, Code: codeString, OrgID: user.ResourceOwner})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	passwordResetData := &PasswordCodeData{
		TemplateData: GetTemplateData(translator, args, apiDomain, url, domain.PasswordResetMessageType, user.PreferredLanguage, colors),
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
