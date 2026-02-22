package domain

import (
	"time"

	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/crypto"
)

type WebKeyState int

const (
	WebKeyStateUnspecified WebKeyState = iota
	WebKeyStateInitial
	WebKeyStateActive
	WebKeyStateInactive
	WebKeyStateRemoved
)

// WebKey represents a cached web key for signing or verification
type WebKey struct {
	KeyID        string
	InstanceID   string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     int64
	State        WebKeyState
	PrivateKey   *crypto.CryptoValue
	PublicKey    *jose.JSONWebKey
	Config       crypto.WebKeyConfig
}

// WebKeys represents a collection of web keys for caching
type WebKeys struct {
	InstanceID string
	Keys       []*WebKey
}
