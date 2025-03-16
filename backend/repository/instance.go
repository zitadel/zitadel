package repository

import "github.com/zitadel/zitadel/backend/storage/cache"

type Instance struct {
	ID   string
	Name string
}

type InstanceIndex uint8

var InstanceIndices = []InstanceIndex{
	InstanceByID,
	InstanceByDomain,
}

const (
	InstanceByID InstanceIndex = iota
	InstanceByDomain
)

var _ cache.Entry[InstanceIndex, string] = (*Instance)(nil)

// Keys implements [cache.Entry].
func (i *Instance) Keys(index InstanceIndex) (key []string) {
	switch index {
	case InstanceByID:
		return []string{i.ID}
	case InstanceByDomain:
		return []string{i.Name}
	}
	return nil
}

type ListRequest struct {
	Limit uint16
}
