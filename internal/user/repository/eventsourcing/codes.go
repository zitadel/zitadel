package eventsourcing

import (
	"github.com/caos/zitadel/internal/crypto"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func (es *UserEventstore) generatePasswordCode(passwordCode *model.PasswordCode, notifyType usr_model.NotificationType) error {
	passwordCodeCrypto, _, err := crypto.NewCode(es.PasswordVerificationCode)
	if err != nil {
		return err
	}
	passwordCode.Code = passwordCodeCrypto
	passwordCode.Expiry = es.PasswordVerificationCode.Expiry()
	passwordCode.NotificationType = int32(notifyType)
	return nil
}
