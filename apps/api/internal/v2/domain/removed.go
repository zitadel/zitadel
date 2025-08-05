package domain

const RemovedTypeSuffix = "domain.removed"

type RemovedPayload struct {
	Name string `json:"domain"`
}
