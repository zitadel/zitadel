package model

import (
	"github.com/caos/zitadel/internal/crypto"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/model"
	"time"
)

type IdpConfigView struct {
	ResourceOwner string
	IdpConfigID   string
	Name          string
	LogoSrc       string
	State         iam_model.IdpConfigState
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64

	IsOidc           bool
	OidcClientID     string
	OidcClientSecret crypto.CryptoValue
	OidcIssuer       string
	OidcScopes       []string
}

type IdpConfigSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IdpConfigSearchKey
	Asc           bool
	Queries       []*IdpConfigSearchQuery
}

type IdpConfigSearchKey int32

const (
	IdpConfigSearchKeyUnspecified IdpConfigSearchKey = iota
	IdpConfigSearchKeyName
	IdpConfigSearchKeyIamID
	IdpConfigSearchKeyIdpConfigID
)

type IdpConfigSearchQuery struct {
	Key    IdpConfigSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type IdpConfigSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IdpConfigView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IdpConfigSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
