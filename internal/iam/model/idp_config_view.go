package model

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

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
	ResourceOwner   string
	IDPProviderType IDPProviderType

	*IDPConfigOIDCView
	*IDPConfigAuthConnectorView
}

type IDPConfigOIDCView struct {
	OIDCClientID               string
	OIDCClientSecret           *crypto.CryptoValue
	OIDCIssuer                 string
	OIDCScopes                 []string
	OIDCIDPDisplayNameMapping  OIDCMappingField
	OIDCUsernameMapping        OIDCMappingField
	OAuthAuthorizationEndpoint string
	OAuthTokenEndpoint         string
}

type IDPConfigAuthConnectorView struct {
	AuthConnectorBaseURL     string
	AuthConnectorProviderID  string
	AuthConnectorMachineID   string
	AuthConnectorMachineName string
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
	IDPConfigSearchKeyAuthConnectorMachineID
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
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-Mv9sd", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *IDPConfigSearchRequest) AppendMyOrgQuery(orgID, iamID string) {
	r.Queries = append(r.Queries, &IDPConfigSearchQuery{Key: IDPConfigSearchKeyAggregateID, Method: domain.SearchMethodIsOneOf, Value: []string{orgID, iamID}})
}
