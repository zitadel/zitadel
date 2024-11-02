package domain

const VerifiedTypeSuffix = "domain.verified"

type VerifiedPayload struct {
	Name string `json:"domain"`
}
