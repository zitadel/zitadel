package avatar

const AvatarAddedTypeSuffix = ".avatar.added"

type AddedPayload struct {
	StoreKey string `json:"storeKey"`
}
