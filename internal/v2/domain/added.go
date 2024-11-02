package domain

const AddedTypeSuffix = "domain.added"

type AddedPayload struct {
	Name string `json:"domain"`
}
