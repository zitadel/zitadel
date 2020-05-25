package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	FirstName string
	LastName  string
	Code      string
	UserID    string
}

func SendPhoneVerificationCode(user *view_model.NotifyUser, code *es_model.PhoneCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	codeData := &PhoneVerificationCodeData{FirstName: user.FirstName, LastName: user.LastName, UserID: user.ID, Code: codeString}
	template, err := templates.ParseTemplateText(systemDefaults.Notifications.TemplateData.VerifyPhone.Text, codeData)
	if err != nil {
		return err
	}
	return generateSms(user, template, systemDefaults.Notifications, true)
}
