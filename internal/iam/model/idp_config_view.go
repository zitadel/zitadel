package model

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"

	"time"
)

type IDPConfigView struct {
	AggregateID     string
	IDPConfigID     string
	Name            string
	StylingType     IDPStylingType
	AutoRegister    bool
	State           IDPConfigState
	CreationDate    time.Time
	ChangeDate      time.Time
	IDPProviderType IDPProviderType

	IsOIDC                     bool
	OIDCClientID               string
	OIDCClientSecret           *crypto.CryptoValue
	OIDCIssuer                 string
	OIDCScopes                 []string
	OIDCIDPDisplayNameMapping  OIDCMappingField
	OIDCUsernameMapping        OIDCMappingField
	OAuthAuthorizationEndpoint string
	OAuthTokenEndpoint         string
	JWTEndpoint                string
	JWTIssuer                  string
	JWTKeysEndpoint            string
	JWTHeaderName              string
}

type IDPConfigSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IDPConfigSearchKey
	Asc           bool
	Queries       []*IDPConfigSearchQuery
}

type IDPConfigSearchKey int32

const (
	IDPConfigSearchKeyUnspecified IDPConfigSearchKey = iota
	IDPConfigSearchKeyName
	IDPConfigSearchKeyAggregateID
	IDPConfigSearchKeyIdpConfigID
	IDPConfigSearchKeyIdpProviderType
	IDPConfigSearchKeyInstanceID
)

type IDPConfigSearchQuery struct {
	Key    IDPConfigSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type IDPConfigSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IDPConfigView
	Timestamp   time.Time
}

func (r *IDPConfigSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-Mv9sd", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
