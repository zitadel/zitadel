package domain

const PrimarySetTypeSuffix = "domain.primary.set"

type PrimarySetPayload struct {
	Name string `json:"domain"`
}
