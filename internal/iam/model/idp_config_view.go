package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	Sequence        uint64
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
	IDPConfigSearchKeyOwnerRemoved
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
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IDPConfigSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-Mv9sd", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *IDPConfigSearchRequest) AppendMyOrgQuery(orgID, iamID string) {
	r.Queries = append(r.Queries, &IDPConfigSearchQuery{Key: IDPConfigSearchKeyAggregateID, Method: domain.SearchMethodIsOneOf, Value: []string{orgID, iamID}})
}
