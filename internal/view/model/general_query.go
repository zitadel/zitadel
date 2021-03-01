package model

import (
	"github.com/caos/zitadel/internal/domain"
)

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
	Method domain.SearchMethod
	Value  interface{}
}
