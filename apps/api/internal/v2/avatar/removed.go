package avatar

const AvatarRemovedTypeSuffix = ".avatar.removed"

type RemovedPayload struct {
	StoreKey string `json:"storeKey"`
}
