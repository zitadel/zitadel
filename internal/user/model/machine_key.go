package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type MachineKeyView struct {
	ID             string
	UserID         string
	Type           MachineKeyType
	Sequence       uint64
	CreationDate   time.Time
	ExpirationDate time.Time
	PublicKey      []byte
}

type MachineKeySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn MachineKeySearchKey
	Asc           bool
	Queries       []*MachineKeySearchQuery
}

type MachineKeySearchKey int32

const (
	MachineKeyKeyUnspecified MachineKeySearchKey = iota
	MachineKeyKeyID
	MachineKeyKeyUserID
)

type MachineKeySearchQuery struct {
	Key    MachineKeySearchKey
	Method model.SearchMethod
	Value  interface{}
}

type MachineKeySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*MachineKeyView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *MachineKeySearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
