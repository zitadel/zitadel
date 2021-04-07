package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	yearLayout            = "2006-01-02"
	defaultExpirationDate = "9999-01-01"
)

type AuthNKeyView struct {
	ID             string
	ObjectID       string
	ObjectType     ObjectType
	AuthIdentifier string
	Type           AuthNKeyType
	Sequence       uint64
	CreationDate   time.Time
	ExpirationDate time.Time
	PublicKey      []byte
	State          AuthNKeyState
}

type AuthNKey struct {
	models.ObjectRoot

	KeyID          string
	ObjectType     ObjectType
	Type           AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
}

type AuthNKeyType int32

const (
	AuthNKeyTypeNONE = iota
	AuthNKeyTypeJSON
)

type AuthNKeyState int32

const (
	AuthNKeyStateActive AuthNKeyState = iota
	AuthNKeyStateInactive
	AuthNKeyStateRemoved
)

type AuthNKeySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn AuthNKeySearchKey
	Asc           bool
	Queries       []*AuthNKeySearchQuery
}

type AuthNKeySearchKey int32

const (
	AuthNKeyKeyUnspecified AuthNKeySearchKey = iota
	AuthNKeyKeyID
	AuthNKeyObjectID
	AuthNKeyObjectType
)

type ObjectType int32

const (
	AuthNKeyObjectTypeUnspecified ObjectType = iota
	AuthNKeyObjectTypeUser
	AuthNKeyObjectTypeApplication
)

type AuthNKeySearchQuery struct {
	Key    AuthNKeySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type AuthNKeySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*AuthNKeyView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *AuthNKeySearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-f9ids", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func DefaultExpiration() (time.Time, error) {
	return time.Parse(yearLayout, defaultExpirationDate)
}
