package model

import "github.com/caos/zitadel/internal/model"

type GeneralSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn GeneralSearchKey
	Asc           bool
	Queries       []*GeneralSearchQuery
}

type GeneralSearchKey int32

const (
	GeneralSearchKeyUnspecified GeneralSearchKey = iota
)

type GeneralSearchQuery struct {
	Key    GeneralSearchKey
	Method model.SearchMethod
	Value  interface{}
}
