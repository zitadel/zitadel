package handler

type webAuthNData struct {
	userData
	CredentialCreationData string
}

type webAuthNFormData struct {
	CredentialData string `schema:"credentialData"`
	Name           string `schema:"name"`
	Recreate       bool   `schema:"recreate"`
}
