package types

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	UserID string
}

func SendPhoneVerificationCode(translator *i18n.Translator, user *view_model.NotifyUser, code *es_model.PhoneCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	text := translator.Localize(fmt.Sprintf("%s.%s", domain.VerifyPhoneMessageType, domain.MessageText), args, user.PreferredLanguage)

	codeData := &PhoneVerificationCodeData{UserID: user.ID}
	template, err := templates.ParseTemplateText(text, codeData)
	if err != nil {
		return err
	}
	return generateSms(user, template, systemDefaults.Notifications, true)
}
