package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type UserLockView struct {
	ID                       string
	ChangeDate               time.Time
	ResourceOwner            string
	Sequence                 uint64
	State                    int32
	PasswordCheckFailedCount uint64
}

type UserLockSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserLockSearchKey
	Asc           bool
	Queries       []*UserLockSearchQuery
}

type UserLockSearchKey int32

const (
	UserLockSearchKeyUnspecified UserLockSearchKey = iota
	UserLockSearchKeyUserID
)

type UserLockSearchQuery struct {
	Key    UserLockSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type UserLockSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserLockView
	Sequence    uint64
	Timestamp   time.Time
}
