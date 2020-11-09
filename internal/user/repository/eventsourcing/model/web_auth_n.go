package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

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
	WebauthNTokenID string `json:"webAuthNTokenId"`
	KeyID           []byte `json:"keyID"`
	PublicKey       []byte `json:"publicKey"`
	AttestationType string `json:"attestationType"`
	AAGUID          []byte `json:"aaguid"`
	SignCount       uint32 `json:"signCount"`
}

type WebAuthNTokenID struct {
	WebauthNTokenID string `json:"webAuthNTokenId"`
}

func GetWebauthn(webauthnTokens []*WebAuthNToken, id string) (int, *WebAuthNToken) {
	for i, webauthn := range webauthnTokens {
		if webauthn.WebauthNTokenID == id {
			return i, webauthn
		}
	}
	return -1, nil
}

func WebAuthNsToModel(u2fs []*WebAuthNToken) []*model.WebAuthNToken {
	convertedIDPs := make([]*model.WebAuthNToken, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNToModel(m)
	}
	return convertedIDPs
}

func WebAuthNsFromModel(u2fs []*model.WebAuthNToken) []*WebAuthNToken {
	convertedIDPs := make([]*WebAuthNToken, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNFromModel(m)
	}
	return convertedIDPs
}

func WebAuthNFromModel(webAuthN *model.WebAuthNToken) *WebAuthNToken {
	return &WebAuthNToken{
		ObjectRoot:      webAuthN.ObjectRoot,
		WebauthNTokenID: webAuthN.WebAuthNTokenID,
		Challenge:       webAuthN.Challenge,
		State:           int32(webAuthN.State),
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
	}
}

func WebAuthNToModel(webAuthN *WebAuthNToken) *model.WebAuthNToken {
	return &model.WebAuthNToken{
		ObjectRoot:      webAuthN.ObjectRoot,
		WebAuthNTokenID: webAuthN.WebauthNTokenID,
		Challenge:       webAuthN.Challenge,
		State:           model.MfaState(webAuthN.State),
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
	}
}

func WebAuthNVerifyFromModel(webAuthN *model.WebAuthNToken) *WebAuthNVerify {
	return &WebAuthNVerify{
		WebauthNTokenID: webAuthN.WebAuthNTokenID,
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
	}
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
	if _, token := GetWebauthn(u.U2FTokens, webauthn.WebauthNTokenID); token != nil {
		token.setData(event)
		token.State = int32(model.MfaStateReady)
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "MODEL-4hu9s", "Errors.Users.Mfa.U2F.NotExisting")
}

func (u *Human) appendU2FRemovedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	for i := len(u.U2FTokens) - 1; i >= 0; i-- {
		if u.U2FTokens[i].WebauthNTokenID == webauthn.WebauthNTokenID {
			copy(u.U2FTokens[i:], u.U2FTokens[i+1:])
			u.U2FTokens[len(u.U2FTokens)-1] = nil
			u.U2FTokens = u.U2FTokens[:len(u.U2FTokens)-1]
			return nil
		}
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
