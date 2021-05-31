package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	UserID string
}

func SendPhoneVerificationCode(i18n *i18n.Translator, user *view_model.NotifyUser, code *es_model.PhoneCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	systemDefaults.Notifications.TemplateData.VerifyPhone.Translate(i18n, args, user.PreferredLanguage)
	codeData := &PhoneVerificationCodeData{UserID: user.ID}
	template, err := templates.ParseTemplateText(systemDefaults.Notifications.TemplateData.VerifyPhone.Text, codeData)
	if err != nil {
		return err
	}
	return generateSms(user, template, systemDefaults.Notifications, true)
}
