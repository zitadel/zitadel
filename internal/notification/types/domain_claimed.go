package types

import (
	"strings"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/templates"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(mailhtml string, i18n *i18n.Translator, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID})
	if err != nil {
		return err
	}
	var args = map[string]interface{}{
		"FirstName":    user.FirstName,
		"LastName":     user.LastName,
		"Username":     user.LastEmail,
		"TempUsername": username,
		"Domain":       strings.Split(user.LastEmail, "@")[1],
	}
	systemDefaults.Notifications.TemplateData.DomainClaimed.Translate(i18n, args, user.PreferredLanguage)
	data := &DomainClaimedData{TemplateData: systemDefaults.Notifications.TemplateData.DomainClaimed, URL: url}
	template, err := templates.GetParsedTemplate(mailhtml, data)
	if err != nil {
		return err
	}
	return generateEmail(user, systemDefaults.Notifications.TemplateData.DomainClaimed.Subject, template, systemDefaults.Notifications, true)
}
