package types

import (
	"strings"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults, colors *query.LabelPolicy, apiDomain string) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID, OrgID: user.ResourceOwner})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["TempUsername"] = username
	args["Domain"] = strings.Split(user.LastEmail, "@")[1]

	domainClaimedData := &DomainClaimedData{
		TemplateData: GetTemplateData(translator, args, apiDomain, url, domain.DomainClaimedMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, domainClaimedData)
	if err != nil {
		return err
	}
	return generateEmail(user, domainClaimedData.Subject, template, systemDefaults.Notifications, true)
}
