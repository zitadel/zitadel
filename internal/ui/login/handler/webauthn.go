package handler

type webAuthNData struct {
	userData
	CredentialCreationData string
	SessionID              string
}

type webAuthNFormData struct {
	CredentialData string `schema:"credentialData"`
	SessionID      string `schema:"sessionID"`
	Resend         bool   `schema:"recreate"`
	Name           bool   `schema:"name"`
}
