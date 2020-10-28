package webauthn

import (
	"bytes"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"

	usr_model "github.com/caos/zitadel/internal/user/model"
)

type WebAuthN struct {
	web *webauthn.WebAuthn
}

func StartServer(displayName, id, origin string) (*WebAuthN, error) {
	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: displayName,
		RPID:          id,
		RPOrigin:      origin,
		Debug:         true,
	})
	if err != nil {
		return nil, err
	}
	return &WebAuthN{
		web: web,
	}, err
}

type user struct {
	id          string
	username    string
	displayName string
	credentials []webauthn.Credential
}

func (u *user) WebAuthnID() []byte {
	return []byte(u.id)
}

func (u *user) WebAuthnName() string {
	return u.username
}

func (u *user) WebAuthnDisplayName() string {
	return u.displayName
}

func (u *user) WebAuthnIcon() string {
	return ""
}

func (u *user) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (w *WebAuthN) BeginRegistration(view *usr_model.UserView, authType protocol.AuthenticatorAttachment, userVerification protocol.UserVerificationRequirement, creds ...webauthn.Credential) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	residentKeyRequirement := false
	user := &user{
		id:          view.ID,
		username:    view.UserName,
		displayName: view.DisplayName,
		credentials: creds,
	}
	existing := make([]protocol.CredentialDescriptor, len(creds))
	for i, cred := range creds {
		existing[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}
	credentialOptions, sessionData, err := w.web.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			RequireResidentKey: &residentKeyRequirement,
			UserVerification:   userVerification,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
		webauthn.WithExclusions(existing),
	)
	if err != nil {
		return nil, nil, err
	}
	return credentialOptions, sessionData, nil
}

func (w *WebAuthN) FinishRegistration(view *usr_model.UserView, sessionData webauthn.SessionData, credentialData *protocol.ParsedCredentialCreationData) (*webauthn.Credential, error) {
	user := &user{
		id:          view.ID,
		username:    view.UserName,
		displayName: view.DisplayName,
		credentials: nil,
	}
	credential, err := w.web.CreateCredential(user, sessionData, credentialData)
	if err != nil {
		//return nil, err
	}
	return credential, nil
}

func (w *WebAuthN) BeginLogin(view *usr_model.UserView, userVerification protocol.UserVerificationRequirement, creds ...webauthn.Credential) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user := &user{
		id:          view.ID,
		username:    view.UserName,
		displayName: view.DisplayName,
		credentials: creds,
	}
	assertion, sessionData, err := w.web.BeginLogin(user)//webauthn.WithUserVerification(userVerification),

	if err != nil {
		return nil, nil, err
	}
	return assertion, sessionData, nil
}

func (w *WebAuthN) FinishLogin(view *usr_model.UserView, sessionData webauthn.SessionData, assertionData *protocol.ParsedCredentialAssertionData, creds ...webauthn.Credential) error {
	user := &user{
		id:          view.ID,
		username:    view.UserName,
		displayName: view.DisplayName,
		credentials: creds,
	}
	credential, err := w.web.ValidateLogin(user, sessionData, assertionData)
	if err != nil {
		return err
	}

	if credential.Authenticator.CloneWarning {
		return nil //ErrCredentialCloned
	}
	for _, cred := range user.WebAuthnCredentials() {
		if bytes.Equal(cred.ID, credential.ID) {

		}
	}

	//w.storage.UpdateSignCount(credential.AuthenticatorID, credential.Authenticator.SignCount)
	return nil
}

//let options = JSON.parse(atob(document.getElementsByName('credentialCreationData')[0].value));
//options.publicKey.challenge = base64js.toByteArray(options.publicKey.challenge);
//options.publicKey.user.id = atob(options.publicKey.user.id);
//navigator.credentials.get({publicKey: options.publicKey})
//.then(function (credential) {
//console.log(credential);
//verifyAssertion(credential);
//}).catch(function (err) {
//console.log(err.name);
//alert(err.message);
//});
