package authncache

import (
	"github.com/zitadel/zitadel/internal/cachekey"
)


type CachedPublicKey struct {
	Algorithm  string
	Use        string
	KeyID      string
	InstanceID string
	UserType   string
	Key        any
	Expiry     int64
}


func (c *CachedPublicKey) Keys(index cachekey.AuthnKeyIndex) []string {
	switch index {
	case cachekey.InstanceID:
		return []string{c.InstanceID}
	case cachekey.UserType:
		return []string{c.UserType}
	case cachekey.KeyID:
		return []string{c.KeyID}
	default:
		return nil
	}
}
