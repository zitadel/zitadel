package domain

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type MetaData struct {
	es_models.ObjectRoot

	State MetaDataState
	Key   string
	Value string
}

type MetaDataState int32

const (
	MetaDataStateUnspecified MetaDataState = iota
	MetaDataStateActive
	MetaDataStateRemoved
)

func (u *MetaData) IsValid() bool {
	return u.Key != "" && u.Value != ""
}

func (s MetaDataState) Exists() bool {
	return s != MetaDataStateUnspecified && s != MetaDataStateRemoved
}
