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
	SortingColumn ApplicationSearchKey
	Asc           bool
	Queries       []*ApplicationSearchQuery
}

type ApplicationSearchKey int32

const (
	APPLICATIONSEARCHKEY_UNSPECIFIED ApplicationSearchKey = iota
	APPLICATIONSEARCHKEY_NAME
	APPLICATIONSEARCHKEY_OIDC_CLIENT_ID
	APPLICATIONSEARCHKEY_PROJECT_ID
	APPLICATIONSEARCHKEY_APP_ID
)

type ApplicationSearchQuery struct {
	Key    ApplicationSearchKey
	Method model.SearchMethod
	Value  string
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
