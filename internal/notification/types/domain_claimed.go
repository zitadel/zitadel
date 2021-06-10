package types

import (
	"html"
	"strings"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/notification/templates"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(mailhtml string, text *iam_model.MessageTextView, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults, colors *iam_model.LabelPolicyView, apiDomain string) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["TempUsername"] = username
	args["Domain"] = strings.Split(user.LastEmail, "@")[1]

	text.Greeting, err = templates.ParseTemplateText(text.Greeting, args)
	text.Text, err = templates.ParseTemplateText(text.Text, args)
	text.Text = html.UnescapeString(text.Text)

	emailCodeData := &DomainClaimedData{
		TemplateData: templates.GetTemplateData(apiDomain, url, text, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, text.Subject, template, systemDefaults.Notifications, true)
}
