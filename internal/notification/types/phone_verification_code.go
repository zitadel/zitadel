package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/notification/templates"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type PhoneVerificationCodeData struct {
	UserID string
}

func SendPhoneVerificationCode(text *iam_model.MessageTextView, user *view_model.NotifyUser, code *es_model.PhoneCode, systemDefaults systemdefaults.SystemDefaults, alg crypto.EncryptionAlgorithm) error {
	codeString, err := crypto.DecryptString(code.Code, alg)
	if err != nil {
		return err
	}
	var args = mapNotifyUserToArgs(user)
	args["Code"] = codeString

	text.Text, err = templates.ParseTemplateText(text.Text, args)

	codeData := &PhoneVerificationCodeData{UserID: user.ID}
	template, err := templates.ParseTemplateText(text.Text, codeData)
	if err != nil {
		return err
	}
	return generateSms(user, template, systemDefaults.Notifications, true)
}
