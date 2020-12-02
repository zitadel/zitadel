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

	KeyID           []byte `json:"keyId"`
	PublicKey       []byte `json:"publicKey"`
	AttestationType string `json:"attestationType"`
	AAGUID          []byte `json:"aaguid"`
	SignCount       uint32 `json:"signCount"`
}

type WebAuthNVerify struct {
	WebAuthNTokenID   string `json:"webAuthNTokenId"`
	KeyID             []byte `json:"keyId"`
	PublicKey         []byte `json:"publicKey"`
	AttestationType   string `json:"attestationType"`
	AAGUID            []byte `json:"aaguid"`
	SignCount         uint32 `json:"signCount"`
	WebAuthNTokenName string `json:"webAuthNTokenName"`
}

type WebAuthNSignCount struct {
	WebauthNTokenID string `json:"webAuthNTokenId"`
	SignCount       uint32 `json:"signCount"`
}

type WebAuthNTokenID struct {
	WebauthNTokenID string `json:"webAuthNTokenId"`
}

type WebAuthNLogin struct {
	es_models.ObjectRoot

	WebauthNTokenID string `json:"webAuthNTokenId"`
	Challenge       string `json:"challenge"`
	*AuthRequest
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
		State:           model.MFAState(webAuthN.State),
		KeyID:           webAuthN.KeyID,
		PublicKey:       webAuthN.PublicKey,
		AAGUID:          webAuthN.AAGUID,
		SignCount:       webAuthN.SignCount,
		AttestationType: webAuthN.AttestationType,
	}
}

func WebAuthNVerifyFromModel(webAuthN *model.WebAuthNToken) *WebAuthNVerify {
	return &WebAuthNVerify{
		WebAuthNTokenID:   webAuthN.WebAuthNTokenID,
		KeyID:             webAuthN.KeyID,
		PublicKey:         webAuthN.PublicKey,
		AAGUID:            webAuthN.AAGUID,
		SignCount:         webAuthN.SignCount,
		AttestationType:   webAuthN.AttestationType,
		WebAuthNTokenName: webAuthN.WebAuthNTokenName,
	}
}

func WebAuthNLoginsToModel(u2fs []*WebAuthNLogin) []*model.WebAuthNLogin {
	convertedIDPs := make([]*model.WebAuthNLogin, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNLoginToModel(m)
	}
	return convertedIDPs
}

func WebAuthNLoginsFromModel(u2fs []*model.WebAuthNLogin) []*WebAuthNLogin {
	convertedIDPs := make([]*WebAuthNLogin, len(u2fs))
	for i, m := range u2fs {
		convertedIDPs[i] = WebAuthNLoginFromModel(m)
	}
	return convertedIDPs
}

func WebAuthNLoginFromModel(webAuthN *model.WebAuthNLogin) *WebAuthNLogin {
	return &WebAuthNLogin{
		ObjectRoot:  webAuthN.ObjectRoot,
		Challenge:   webAuthN.Challenge,
		AuthRequest: AuthRequestFromModel(webAuthN.AuthRequest),
	}
}

func WebAuthNLoginToModel(webAuthN *WebAuthNLogin) *model.WebAuthNLogin {
	return &model.WebAuthNLogin{
		ObjectRoot:  webAuthN.ObjectRoot,
		Challenge:   webAuthN.Challenge,
		AuthRequest: AuthRequestToModel(webAuthN.AuthRequest),
	}
}

func (u *Human) appendU2FAddedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	webauthn.ObjectRoot.CreationDate = event.CreationDate
	webauthn.State = int32(model.MFAStateNotReady)
	for i, token := range u.U2FTokens {
		if token.State == int32(model.MFAStateNotReady) {
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
		err := token.setData(event)
		if err != nil {
			return err
		}
		token.State = int32(model.MFAStateReady)
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "MODEL-4hu9s", "Errors.Users.MFA.U2F.NotExisting")
}

func (u *Human) appendU2FChangeSignCountEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	if _, token := GetWebauthn(u.U2FTokens, webauthn.WebauthNTokenID); token != nil {
		token.setData(event)
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "MODEL-5Ms8h", "Errors.Users.MFA.U2F.NotExisting")
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

func (u *Human) appendPasswordlessAddedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	webauthn.ObjectRoot.CreationDate = event.CreationDate
	webauthn.State = int32(model.MFAStateNotReady)
	for i, token := range u.PasswordlessTokens {
		if token.State == int32(model.MFAStateNotReady) {
			u.PasswordlessTokens[i] = webauthn
			return nil
		}
	}
	u.PasswordlessTokens = append(u.PasswordlessTokens, webauthn)
	return nil
}

func (u *Human) appendPasswordlessVerifiedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	if _, token := GetWebauthn(u.PasswordlessTokens, webauthn.WebauthNTokenID); token != nil {
		err := token.setData(event)
		if err != nil {
			return err
		}
		token.State = int32(model.MFAStateReady)
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "MODEL-mKns8", "Errors.Users.MFA.Passwordless.NotExisting")
}

func (u *Human) appendPasswordlessChangeSignCountEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	if _, token := GetWebauthn(u.PasswordlessTokens, webauthn.WebauthNTokenID); token != nil {
		err := token.setData(event)
		if err != nil {
			return err
		}
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "MODEL-2Mv9s", "Errors.Users.MFA.Passwordless.NotExisting")
}

func (u *Human) appendPasswordlessRemovedEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNToken)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	for i := len(u.PasswordlessTokens) - 1; i >= 0; i-- {
		if u.PasswordlessTokens[i].WebauthNTokenID == webauthn.WebauthNTokenID {
			copy(u.PasswordlessTokens[i:], u.PasswordlessTokens[i+1:])
			u.PasswordlessTokens[len(u.PasswordlessTokens)-1] = nil
			u.PasswordlessTokens = u.PasswordlessTokens[:len(u.PasswordlessTokens)-1]
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

func (u *Human) appendU2FLoginEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNLogin)
	webauthn.ObjectRoot.AppendEvent(event)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	webauthn.ObjectRoot.CreationDate = event.CreationDate
	for i, token := range u.U2FLogins {
		if token.AuthRequest.ID == webauthn.AuthRequest.ID {
			u.U2FLogins[i] = webauthn
			return nil
		}
	}
	u.U2FLogins = append(u.U2FLogins, webauthn)
	return nil
}

func (u *Human) appendPasswordlessLoginEvent(event *es_models.Event) error {
	webauthn := new(WebAuthNLogin)
	webauthn.ObjectRoot.AppendEvent(event)
	err := webauthn.setData(event)
	if err != nil {
		return err
	}
	webauthn.ObjectRoot.CreationDate = event.CreationDate
	for i, token := range u.PasswordlessLogins {
		if token.AuthRequest.ID == webauthn.AuthRequest.ID {
			u.PasswordlessLogins[i] = webauthn
			return nil
		}
	}
	u.PasswordlessLogins = append(u.PasswordlessLogins, webauthn)
	return nil
}

func (w *WebAuthNLogin) setData(event *es_models.Event) error {
	w.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, w); err != nil {
		logging.Log("EVEN-hmSlo").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lo023", "could not unmarshal event")
	}
	return nil
}
