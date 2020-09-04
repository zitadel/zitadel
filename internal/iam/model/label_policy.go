package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

// ToDo Michi
type LabelPolicy struct {
	models.ObjectRoot

	State          PolicyState
	Default        bool
	PrimaryColor   string
	SecundaryColor string
}

// type IDPProvider struct {
// 	models.ObjectRoot
// 	Type        IDPProviderType
// 	IdpConfigID string
// }

// type PolicyState int32

//const (
//	PolicyStateActive PolicyState = iota
//	PolicyStateRemoved
//)

// type IDPProviderType int32

// const (
// 	IDPProviderTypeSystem IDPProviderType = iota
// 	IDPProviderTypeOrg
// )

func (p *LabelPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}

// func (p *IDPProvider) IsValid() bool {
// 	return p.ObjectRoot.AggregateID != "" && p.IdpConfigID != ""
// }

// func (p *LabelPolicy) GetIdpProvider(id string) (int, *IDPProvider) {
// 	for i, m := range p.IDPProviders {
// 		if m.IdpConfigID == id {
// 			return i, m
// 		}
// 	}
// 	return -1, nil
// }
