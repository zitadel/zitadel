package model

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/model"
)

type KeyView struct {
	ID              string
	Private         bool
	Expiry          time.Time
	Algorithm       string
	Usage           string
	Key             *crypto.CryptoValue
	CurrentSequence uint64
}

type KeySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn KeySearchKey
	Asc           bool
	Queries       []*KeySearchQuery
}

type KeySearchKey int32

const (
	KEYSEARCHKEY_UNSPECIFIED KeySearchKey = iota
	KEYSEARCHKEY_PRIVATE
	KEYSEARCHKEY_EXPIRY
	KEYSEARCHKEY_USAGE
)

type KeySearchQuery struct {
	Key    KeySearchKey
	Method model.SearchMethod
	Value  interface{}
}

type KeySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*KeyView
}

func (r *KeySearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
