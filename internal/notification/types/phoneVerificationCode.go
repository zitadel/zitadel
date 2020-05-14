package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	templates.TemplateData
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
	_ = &PhoneVerificationCodeData{TemplateData: systemDefaults.Notifications.TemplateData.VerifyPhone, FirstName: user.FirstName, LastName: user.LastName, UserID: user.ID, Code: codeString}

	//TODO: generateSMS
	return nil
}
