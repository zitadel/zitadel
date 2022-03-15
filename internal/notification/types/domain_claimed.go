package types

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/notification/templates"
	"github.com/caos/zitadel/internal/query"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(ctx context.Context, mailhtml string, translator *i18n.Translator, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults, emailConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), colors *query.LabelPolicy, assetsPrefix string) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID})
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["TempUsername"] = username
	args["Domain"] = strings.Split(user.LastEmail, "@")[1]

	domainClaimedData := &DomainClaimedData{
		TemplateData: GetTemplateData(translator, args, assetsPrefix, url, domain.DomainClaimedMessageType, user.PreferredLanguage, colors),
		URL:          url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, domainClaimedData)
	if err != nil {
		return err
	}
	return generateEmail(ctx, user, domainClaimedData.Subject, template, systemDefaults.Notifications, emailConfig, getFileSystemProvider, getLogProvider, true)
}
