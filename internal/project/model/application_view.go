package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type ApplicationView struct {
	ID                   string
	ProjectID            string
	Name                 string
	CreationDate         time.Time
	ChangeDate           time.Time
	State                AppState
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool

	IsOIDC                     bool
	OIDCVersion                OIDCVersion
	OIDCClientID               string
	OIDCRedirectUris           []string
	OIDCResponseTypes          []OIDCResponseType
	OIDCGrantTypes             []OIDCGrantType
	OIDCApplicationType        OIDCApplicationType
	OIDCAuthMethodType         OIDCAuthMethodType
	OIDCPostLogoutRedirectUris []string
	NoneCompliant              bool
	ComplianceProblems         []string
	DevMode                    bool
	OriginAllowList            []string
	AdditionalOrigins          []string
	AccessTokenType            OIDCTokenType
	IDTokenRoleAssertion       bool
	AccessTokenRoleAssertion   bool
	IDTokenUserinfoAssertion   bool
	ClockSkew                  time.Duration

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
	Method domain.SearchMethod
	Value  interface{}
}

type ApplicationSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ApplicationView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ApplicationSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-3Mf8s", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
