package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers/email"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type InitCodeEmailData struct {
	FirstName string
	LastName  string
	UserID    string
	Code      string
}

func SendUserInitCode(user *view_model.NotifyUser, code *es_model.InitUserCode, systemDefaults systemdefaults.SystemDefaults) error {
	template, err := templates.ParseTemplateFile("", systemDefaults.Notifications.TemplateData.InitCode)
	if err != nil {
		return err
	}
	_ = &InitCodeEmailData{FirstName: user.FirstName, LastName: user.LastName, UserID: user.ID}

	return generateEmail(user, template, systemDefaults.Notifications)
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

		return provider.HandleMessage(message)
	}
	return caos_errs.ThrowInternalf(nil, "NOTIF-s8ipw", "Could not send init message: userid: %v", user.ID)
}
