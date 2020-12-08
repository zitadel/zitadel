package webauthn

import (
	"bytes"
	"encoding/json"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type WebAuthN struct {
	webLogin   *webauthn.WebAuthn
	webConsole *webauthn.WebAuthn
}

func StartServer(sd systemdefaults.WebAuthN) (*WebAuthN, error) {
	webLogin, err := webauthn.New(&webauthn.Config{
		RPDisplayName: sd.DisplayName,
		RPID:          sd.ID,
		RPOrigin:      sd.OriginLogin,
	})
	if err != nil {
		return nil, err
	}
	webConsole, err := webauthn.New(&webauthn.Config{
		RPDisplayName: sd.DisplayName,
		RPID:          sd.ID,
		RPOrigin:      sd.OriginConsole,
	})
	if err != nil {
		return nil, err
	}
	return &WebAuthN{
		webLogin:   webLogin,
		webConsole: webConsole,
	}, err
}

type webUser struct {
	*usr_model.User
	credentials []webauthn.Credential
}

func (u *webUser) WebAuthnID() []byte {
	return []byte(u.AggregateID)
}

func (u *webUser) WebAuthnName() string {
	return u.PreferredLoginName
}

func (u *webUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *webUser) WebAuthnIcon() string {
	return ""
}

func (u *webUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (w *WebAuthN) BeginRegistration(user *usr_model.User, authType usr_model.AuthenticatorAttachment, userVerification usr_model.UserVerificationRequirement, isLoginUI bool, webAuthNs ...*usr_model.WebAuthNToken) (*usr_model.WebAuthNToken, error) {
	creds := WebAuthNsToCredentials(webAuthNs)
	existing := make([]protocol.CredentialDescriptor, len(creds))
	for i, cred := range creds {
		existing[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}
	credentialOptions, sessionData, err := w.web(isLoginUI).BeginRegistration(
		&webUser{
			User:        user,
			credentials: creds,
		},
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			UserVerification:        UserVerificationFromModel(userVerification),
			AuthenticatorAttachment: AuthenticatorAttachmentFromModel(authType),
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
		webauthn.WithExclusions(existing),
	)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-bM8sd", "Errors.User.WebAuthN.BeginRegisterFailed")
	}
	cred, err := json.Marshal(credentialOptions)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-D7cus", "Errors.User.WebAuthN.MarshalError")
	}
	return &usr_model.WebAuthNToken{
		Challenge:              sessionData.Challenge,
		CredentialCreationData: cred,
		AllowedCredentialIDs:   sessionData.AllowedCredentialIDs,
		UserVerification:       UserVerificationToModel(sessionData.UserVerification),
	}, nil
}

func (w *WebAuthN) FinishRegistration(user *usr_model.User, webAuthN *usr_model.WebAuthNToken, tokenName string, credData []byte, isLoginUI bool) (*usr_model.WebAuthNToken, error) {
	if webAuthN == nil {
		return nil, caos_errs.ThrowInternal(nil, "WEBAU-5M9so", "Errors.User.WebAuthN.NotFound")
	}
	credentialData, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credData))
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-sEr8c", "Errors.User.WebAuthN.ErrorOnParseCredential")
	}
	sessionData := WebAuthNToSessionData(webAuthN)
	credential, err := w.web(isLoginUI).CreateCredential(
		&webUser{
			User: user,
		},
		sessionData, credentialData)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-3Vb9s", "Errors.User.WebAuthN.CreateCredentialFailed")
	}

	webAuthN.KeyID = credential.ID
	webAuthN.PublicKey = credential.PublicKey
	webAuthN.AttestationType = credential.AttestationType
	webAuthN.AAGUID = credential.Authenticator.AAGUID
	webAuthN.SignCount = credential.Authenticator.SignCount
	webAuthN.WebAuthNTokenName = tokenName
	return webAuthN, nil
}

func (w *WebAuthN) BeginLogin(user *usr_model.User, userVerification usr_model.UserVerificationRequirement, isLoginUI bool, webAuthNs ...*usr_model.WebAuthNToken) (*usr_model.WebAuthNLogin, error) {
	assertion, sessionData, err := w.web(isLoginUI).BeginLogin(&webUser{
		User:        user,
		credentials: WebAuthNsToCredentials(webAuthNs),
	}, webauthn.WithUserVerification(UserVerificationFromModel(userVerification)))
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-4G8sw", "Errors.User.WebAuthN.BeginLoginFailed")
	}
	cred, err := json.Marshal(assertion)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "WEBAU-2M0s9", "Errors.User.WebAuthN.MarshalError")
	}
	return &usr_model.WebAuthNLogin{
		Challenge:               sessionData.Challenge,
		CredentialAssertionData: cred,
		AllowedCredentialIDs:    sessionData.AllowedCredentialIDs,
		UserVerification:        userVerification,
	}, nil
}

func (w *WebAuthN) FinishLogin(user *usr_model.User, webAuthN *usr_model.WebAuthNLogin, credData []byte, isLoginUI bool, webAuthNs ...*usr_model.WebAuthNToken) ([]byte, uint32, error) {
	assertionData, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(credData))
	if err != nil {
		return nil, 0, caos_errs.ThrowInternal(err, "WEBAU-ADgv4", "Errors.User.WebAuthN.ValidateLoginFailed")
	}
	webUser := &webUser{
		User:        user,
		credentials: WebAuthNsToCredentials(webAuthNs),
	}
	credential, err := w.web(isLoginUI).ValidateLogin(webUser, WebAuthNLoginToSessionData(webAuthN), assertionData)
	if err != nil {
		return nil, 0, caos_errs.ThrowInternal(err, "WEBAU-3M9si", "Errors.User.WebAuthN.ValidateLoginFailed")
	}

	if credential.Authenticator.CloneWarning {
		return credential.ID, credential.Authenticator.SignCount, caos_errs.ThrowInternal(err, "WEBAU-4M90s", "Errors.User.WebAuthN.CloneWarning")
	}
	return credential.ID, credential.Authenticator.SignCount, nil
}

func (w *WebAuthN) web(isLoginUI bool) *webauthn.WebAuthn {
	if isLoginUI {
		return w.webLogin
	}
	return w.webConsole
}
