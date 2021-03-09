package model

import (
	"encoding/json"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
)

type OTP struct {
	es_models.ObjectRoot

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  int32               `json:"-"`
}

type OTPVerified struct {
	UserAgentID string `json:"userAgentID,omitempty"`
}

func (u *Human) appendOTPAddedEvent(event *es_models.Event) error {
	u.OTP = &OTP{
		State: int32(model.MFAStateNotReady),
	}
	return u.OTP.setData(event)
}

func (u *Human) appendOTPVerifiedEvent() {
	u.OTP.State = int32(model.MFAStateReady)
}

func (u *Human) appendOTPRemovedEvent() {
	u.OTP = nil
}

func (o *OTP) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d9soe").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lo023", "could not unmarshal event")
	}
	return nil
}

func (o *OTPVerified) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-BF421").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-GB6hj", "could not unmarshal event")
	}
	return nil
}
