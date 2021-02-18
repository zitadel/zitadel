package model

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
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
	Method model.SearchMethod
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

func (r *AuthNKeySearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func DefaultExpiration() (time.Time, error) {
	return time.Parse(yearLayout, defaultExpirationDate)
}
