package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type OTP struct {
	es_models.ObjectRoot

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  int32               `json:"-"`
}

func OTPFromModel(otp *model.OTP) *OTP {
	return &OTP{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  otp.ObjectRoot.AggregateID,
			Sequence:     otp.Sequence,
			ChangeDate:   otp.ChangeDate,
			CreationDate: otp.CreationDate,
		},
		Secret: otp.Secret,
		State:  int32(otp.State),
	}
}

func OTPToModel(otp *OTP) *model.OTP {
	return &model.OTP{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  otp.ObjectRoot.AggregateID,
			Sequence:     otp.Sequence,
			ChangeDate:   otp.ChangeDate,
			CreationDate: otp.CreationDate,
		},
		Secret: otp.Secret,
		State:  model.MfaState(otp.State),
	}
}
