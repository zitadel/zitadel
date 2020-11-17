package handler

type webAuthNData struct {
	userData
	CredentialCreationData string
	SessionID              string
}

type webAuthNFormData struct {
	CredentialData string `schema:"credentialData"`
	SessionID      string `schema:"sessionID"`
	Name           string `schema:"name"`
	Resend         bool   `schema:"recreate"`
}
