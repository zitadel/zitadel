package model

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type IDPConfigView struct {
	AggregateID     string
	IDPConfigID     string
	Name            string
	StylingType     IDPStylingType
	State           IDPConfigState
	CreationDate    time.Time
	ChangeDate      time.Time
	Sequence        uint64
	IDPProviderType IDPProviderType

	IsOIDC                    bool
	OIDCClientID              string
	OIDCClientSecret          *crypto.CryptoValue
	OIDCIssuer                string
	OIDCScopes                []string
	OIDCIDPDisplayNameMapping OIDCMappingField
	OIDCUsernameMapping       OIDCMappingField
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

func (r *IDPConfigSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *IDPConfigSearchRequest) AppendMyOrgQuery(orgID, iamID string) {
	r.Queries = append(r.Queries, &IDPConfigSearchQuery{Key: IDPConfigSearchKeyAggregateID, Method: domain.SearchMethodIsOneOf, Value: []string{orgID, iamID}})
}
