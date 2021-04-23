package eventstore

type Asset struct {
	// ID is to refer to the asset
	ID string
	//Asset is the actual image
	Asset []byte
	//Action defines if asset should be added or removed
	Action AssetAction
}

type AssetAction int32

const (
	AssetAdd AssetAction = iota
	AssetRemove

	assetActionCount
)

func NewAddAsset(
	id string,
	asset []byte) *Asset {
	return &Asset{
		ID:     id,
		Asset:  asset,
		Action: AssetAdd,
	}
}

func NewRemoveAsset(
	id string) *Asset {
	return &Asset{
		ID:     id,
		Action: AssetRemove,
	}
}
