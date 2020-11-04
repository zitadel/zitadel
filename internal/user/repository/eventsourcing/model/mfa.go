package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/duo-labs/webauthn/protocol"
)

type OTP struct {
	es_models.ObjectRoot

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  int32               `json:"-"`
}

type WebauthNToken struct {
	es_models.ObjectRoot

	SessionID                    string `json:"sessionID,omitempty"`
	CredentialCreationDataString string `json:"credentialCreationDataString,omitempty"`
	CredentialCreationData       *protocol.CredentialCreation
	Challenge                    string   `json:"challenge"`
	UserID                       []byte   `json:"user_id"`
	AllowedCredentialIDs         [][]byte `json:"allowed_credentials,omitempty"`
	UserVerification             string   `json:"userVerification"`
}

func OTPFromModel(otp *model.OTP) *OTP {
	return &OTP{
		ObjectRoot: otp.ObjectRoot,
		Secret:     otp.Secret,
		State:      int32(otp.State),
	}
}

func OTPToModel(otp *OTP) *model.OTP {
	return &model.OTP{
		ObjectRoot: otp.ObjectRoot,
		Secret:     otp.Secret,
		State:      model.MfaState(otp.State),
	}
}

func (u *Human) appendOTPAddedEvent(event *es_models.Event) error {
	u.OTP = &OTP{
		State: int32(model.MfaStateNotReady),
	}
	return u.OTP.setData(event)
}

func (u *Human) appendOTPVerifiedEvent() {
	u.OTP.State = int32(model.MfaStateReady)
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
