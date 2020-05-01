package eventsourcing

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func (es *UserEventstore) generateInitUserCode(initCode *model.InitUserCode) error {
	initCodeCrypto, _, err := crypto.NewCode(es.InitializeUserCode)
	if err != nil {
		return err
	}
	initCode.Code = initCodeCrypto
	initCode.Expiry = es.InitializeUserCode.Expiry()
	return nil
}

func (es *UserEventstore) generatePhoneCode(phoneCode *model.PhoneCode) error {
	phoneCodeCrypto, _, err := crypto.NewCode(es.PhoneVerificationCode)
	if err != nil {
		return err
	}
	phoneCode.Code = phoneCodeCrypto
	phoneCode.Expiry = es.PhoneVerificationCode.Expiry()
	return nil
}

func (es *UserEventstore) generateEmailCode(emailCode *model.EmailCode) error {
	emailCodeCrypto, _, err := crypto.NewCode(es.EmailVerificationCode)
	if err != nil {
		return err
	}
	emailCode.Code = emailCodeCrypto
	emailCode.Expiry = es.EmailVerificationCode.Expiry()
	return nil
}
