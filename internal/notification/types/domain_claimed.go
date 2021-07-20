package types

import (
	"strings"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/i18n"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/notification/templates"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults, colors *iam_model.LabelPolicyView, apiDomain string) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["TempUsername"] = username
	args["Domain"] = strings.Split(user.LastEmail, "@")[1]

	domainClaimedData := &DomainClaimedData{
		TemplateData: templates.GetTemplateData(translator, args, apiDomain, url, domain.DomainClaimedMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, domainClaimedData)
	if err != nil {
		return err
	}
	return generateEmail(user, domainClaimedData.Subject, template, systemDefaults.Notifications, true)
}
