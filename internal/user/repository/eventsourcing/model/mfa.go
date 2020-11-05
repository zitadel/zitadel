package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/duo-labs/webauthn/webauthn"
)

type OTP struct {
	es_models.ObjectRoot

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  int32               `json:"-"`
}

type WebAuthNToken struct {
	es_models.ObjectRoot

	WebauthNTokenID string `json:"webAuthNTokenId"`
	Challenge       string `json:"challenge"`
	State           int32  `json:"-"`

	KeyID           []byte `json:"keyID"`
	PublicKey       []byte `json:"publicKey"`
	AttestationType string `json:"attestationType"`
	AAGUID          []byte `json:"aaguid"`
	SignCount       uint32 `json:"signCount"`
}

type WebAuthNVerify struct {
	es_models.ObjectRoot

	WebauthNTokenID string `json:"webAuthNTokenId"`
	State           int32  `json:"-"`

	KeyID           []byte `json:"keyID"`
	PublicKey       []byte `json:"publicKey"`
	AttestationType string `json:"attestationType"`
	AAGUID          []byte `json:"aaguid"`
	SignCount       uint32 `json:"signCount"`
}

func GetWebauthn(webauthnTokens []*WebAuthNToken, id string) (int, *WebAuthNToken) {
	for i, webauthn := range webauthnTokens {
		if webauthn.WebauthNTokenID == id {
			return i, webauthn
		}
	}
	return -1, nil
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

func WebAuthNsToModel(u2fs []*WebAuthNToken) []*model.WebauthNToken {
	convertedIDPs := make([]*model.WebauthNToken, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNToModel(m)
	}
	return convertedIDPs
}

func WebAuthNsFromModel(u2fs []*model.WebauthNToken) []*WebAuthNToken {
	convertedIDPs := make([]*WebAuthNToken, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNFromModel(m)
	}
	return convertedIDPs
}

func WebAuthNFromModel(webAuthN *model.WebauthNToken) *WebAuthNToken {
	return &WebAuthNToken{
		ObjectRoot:      webAuthN.ObjectRoot,
		WebauthNTokenID: webAuthN.SessionID,
		Challenge:       webAuthN.SessionData.Challenge,
		State:           int32(webAuthN.State),
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
	}
}

func WebAuthNToModel(webAuthN *WebAuthNToken) *model.WebauthNToken {
	return &model.WebauthNToken{
		ObjectRoot: webAuthN.ObjectRoot,
		SessionID:  webAuthN.WebauthNTokenID,
		SessionData: &webauthn.SessionData{
			UserID:    []byte(webAuthN.AggregateID),
			Challenge: webAuthN.Challenge,
		},
		State:           model.MfaState(webAuthN.State),
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
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

func (u *Human) appendU2FAddedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	webauthn.ObjectRoot.CreationDate = event.CreationDate
	webauthn.State = int32(model.MfaStateNotReady)
	for i, token := range u.U2FTokens {
		if token.State == int32(model.MfaStateNotReady) {
			u.U2FTokens[i] = webauthn
			return nil
		}
	}
	u.U2FTokens = append(u.U2FTokens, webauthn)
	return nil
}

func (u *Human) appendU2FVerifiedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	if i, token := GetWebauthn(u.U2FTokens, webauthn.WebauthNTokenID); token != nil {
		u.U2FTokens[i] = u.U2FTokens[len(u.U2FTokens)-1]
		return nil
	}
	webauthn.State = int32(model.MfaStateNotReady)
	u.U2FTokens = append(u.U2FTokens, webauthn)
	return nil
}

func (u *Human) appendU2FRemovedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	if i, token := GetWebauthn(u.U2FTokens, webauthn.WebauthNTokenID); token != nil {
		err = u.U2FTokens[i].setData(event)
		if err != nil {
			return err
		}
		webauthn.State = int32(model.MfaStateReady)
		return nil
	}
	return nil
}

func (w *WebAuthNToken) setData(event *es_models.Event) error {
	w.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, w); err != nil {
		logging.Log("EVEN-4M9is").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lo023", "could not unmarshal event")
	}
	return nil
}
