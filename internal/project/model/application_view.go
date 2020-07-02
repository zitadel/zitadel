package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type ApplicationView struct {
	ID           string
	ProjectID    string
	Name         string
	CreationDate time.Time
	ChangeDate   time.Time
	State        AppState

	IsOIDC                     bool
	OIDCClientID               string
	OIDCRedirectUris           []string
	OIDCResponseTypes          []OIDCResponseType
	OIDCGrantTypes             []OIDCGrantType
	OIDCApplicationType        OIDCApplicationType
	OIDCAuthMethodType         OIDCAuthMethodType
	OIDCPostLogoutRedirectUris []string

	Sequence uint64
}

type ApplicationSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn AppSearchKey
	Asc           bool
	Queries       []*ApplicationSearchQuery
}

type AppSearchKey int32

const (
	AppSearchKeyUnspecified AppSearchKey = iota
	AppSearchKeyName
	AppSearchKeyOIDCClientID
	AppSearchKeyProjectID
	AppSearchKeyAppID
)

type ApplicationSearchQuery struct {
	Key    AppSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type ApplicationSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ApplicationView
}

func (r *ApplicationSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
