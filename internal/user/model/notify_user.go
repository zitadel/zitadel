package model

import (
	"github.com/zitadel/zitadel/internal/domain"
)

type NotifyUserSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn NotifyUserSearchKey
	Asc           bool
	Queries       []*NotifyUserSearchQuery
}

type NotifyUserSearchKey int32

const (
	NotifyUserSearchKeyUnspecified NotifyUserSearchKey = iota
	NotifyUserSearchKeyUserID
	NotifyUserSearchKeyResourceOwner
	NotifyUserSearchKeyInstanceID
)

type NotifyUserSearchQuery struct {
	Key    NotifyUserSearchKey
	Method domain.SearchMethod
	Value  string
}
