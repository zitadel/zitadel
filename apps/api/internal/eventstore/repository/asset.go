package repository

//Asset represents all information about a asset (img)
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
	AssetAdded AssetAction = iota
	AssetRemoved

	assetCount
)

func (f AssetAction) Valid() bool {
	return f >= 0 && f < assetCount
}
