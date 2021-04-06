package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type IDPProviderView struct {
	AggregateID     string
	IDPConfigID     string
	IDPProviderType IDPProviderType
	Name            string
	StylingType     IDPStylingType
	IDPConfigType   IdpConfigType
	IDPState        IDPConfigState

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type IDPProviderSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IDPProviderSearchKey
	Asc           bool
	Queries       []*IDPProviderSearchQuery
}

type IDPProviderSearchKey int32

const (
	IDPProviderSearchKeyUnspecified IDPProviderSearchKey = iota
	IDPProviderSearchKeyAggregateID
	IDPProviderSearchKeyIdpConfigID
	IDPProviderSearchKeyState
)

type IDPProviderSearchQuery struct {
	Key    IDPProviderSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type IDPProviderSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IDPProviderView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IDPProviderSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-8fn7f", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *IDPProviderSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &IDPProviderSearchQuery{Key: IDPProviderSearchKeyAggregateID, Method: domain.SearchMethodEquals, Value: aggregateID})
}

func IdpProviderViewsToDomain(idpProviders []*IDPProviderView) []*domain.IDPProvider {
	providers := make([]*domain.IDPProvider, len(idpProviders))
	for i, provider := range idpProviders {
		p := &domain.IDPProvider{
			IDPConfigID:   provider.IDPConfigID,
			Type:          idpProviderTypeToDomain(provider.IDPProviderType),
			Name:          provider.Name,
			IDPConfigType: idpConfigTypeToDomain(provider.IDPConfigType),
			StylingType:   idpStylingTypeToDomain(provider.StylingType),
			IDPState:      idpStateToDomain(provider.IDPState),
		}
		providers[i] = p
	}
	return providers
}

func idpProviderTypeToDomain(idpType IDPProviderType) domain.IdentityProviderType {
	switch idpType {
	case IDPProviderTypeSystem:
		return domain.IdentityProviderTypeSystem
	case IDPProviderTypeOrg:
		return domain.IdentityProviderTypeOrg
	default:
		return domain.IdentityProviderTypeSystem
	}
}

func idpConfigTypeToDomain(idpType IdpConfigType) domain.IDPConfigType {
	switch idpType {
	case IDPConfigTypeOIDC:
		return domain.IDPConfigTypeOIDC
	case IDPConfigTypeSAML:
		return domain.IDPConfigTypeSAML
	default:
		return domain.IDPConfigTypeOIDC
	}
}

func idpStylingTypeToDomain(stylingType IDPStylingType) domain.IDPConfigStylingType {
	switch stylingType {
	case IDPStylingTypeGoogle:
		return domain.IDPConfigStylingTypeGoogle
	default:
		return domain.IDPConfigStylingTypeUnspecified
	}
}

func idpStateToDomain(state IDPConfigState) domain.IDPConfigState {
	switch state {
	case IDPConfigStateActive:
		return domain.IDPConfigStateActive
	case IDPConfigStateInactive:
		return domain.IDPConfigStateInactive
	case IDPConfigStateRemoved:
		return domain.IDPConfigStateRemoved
	default:
		return domain.IDPConfigStateActive
	}
}
