package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers/email"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func SendUserInitCode(user *view_model.NotifyUser, code *es_model.InitUserCode, config systemdefaults.Notifications) error {
	template, err := templates.ParseTemplateFile("", config.TemplateData.InitCode)
	if err != nil {
		return err
	}
	return generateEmail(user, template, config)
}

func generateEmail(user *view_model.NotifyUser, content string, config systemdefaults.Notifications) error {
	provider, err := email.InitEmailProvider(&config.Providers.Email)
	if err != nil {
		return err
	}
	message := &email.EmailMessage{
		SenderEmail: config.Providers.Email.From,
		Recipients:  []string{user.LastEmail},
		Subject:     config.TemplateData.InitCode.Subject,
		Content:     content,
	}
	if provider.CanHandleMessage(message) {
		provider.HandleMessage(message)
	}
	return caos_errs.ThrowInternalf(nil, "NOTIF-s8ipw", "Could not send init message: userid: %v", user.ID)
}
