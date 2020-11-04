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

type webUser struct {
	*usr_model.User
	credentials []webauthn.Credential
}

func (u *webUser) WebAuthnID() []byte {
	return []byte(u.AggregateID)
}

func (u *webUser) WebAuthnName() string {
	return u.UserName
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

func (w *WebAuthN) BeginRegistration(user *usr_model.User, authType protocol.AuthenticatorAttachment, userVerification protocol.UserVerificationRequirement, creds ...webauthn.Credential) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	//residentKeyRequirement := false
	existing := make([]protocol.CredentialDescriptor, len(creds))
	for i, cred := range creds {
		existing[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}
	credentialOptions, sessionData, err := w.web.BeginRegistration(&webUser{
		User:        user,
		credentials: creds,
	},
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			//RequireResidentKey: &residentKeyRequirement,
			UserVerification: userVerification,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
		webauthn.WithExclusions(existing),
	)
	if err != nil {
		return nil, nil, err
	}
	return credentialOptions, sessionData, nil
}

func (w *WebAuthN) FinishRegistration(user *usr_model.User, sessionData webauthn.SessionData, credentialData *protocol.ParsedCredentialCreationData) (*webauthn.Credential, error) {
	//data, err := json.Marshal(credentialData)
	//parsedCredentialData, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(data))
	credential, err := w.web.CreateCredential(
		&webUser{
			User: user,
		},
		sessionData, credentialData)
	if err != nil {
		//return nil, err
	}
	return credential, nil
}

func (w *WebAuthN) BeginLogin(user *usr_model.User, userVerification protocol.UserVerificationRequirement, creds ...webauthn.Credential) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	assertion, sessionData, err := w.web.BeginLogin(&webUser{
		User:        user,
		credentials: creds,
	}) //webauthn.WithUserVerification(userVerification),

	if err != nil {
		return nil, nil, err
	}
	return assertion, sessionData, nil
}

func (w *WebAuthN) FinishLogin(user *usr_model.User, sessionData webauthn.SessionData, assertionData *protocol.ParsedCredentialAssertionData, creds ...webauthn.Credential) error {
	webUser := &webUser{
		User:        user,
		credentials: creds,
	}
	credential, err := w.web.ValidateLogin(webUser, sessionData, assertionData)
	if err != nil {
		return err
	}

	if credential.Authenticator.CloneWarning {
		return nil //ErrCredentialCloned
	}
	for _, cred := range webUser.WebAuthnCredentials() {
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
