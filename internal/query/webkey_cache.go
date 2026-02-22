package query

import (
	jose "github.com/go-jose/go-jose/v4"
)

// webkeyActiveSigningKeyCacheIndex represents the index for active signing key cache
type webkeyActiveSigningKeyCacheIndex int

const (
	webkeyActiveSigningKeyCacheInstanceIndex webkeyActiveSigningKeyCacheIndex = iota
)

// webkeyActiveSigningKeyCacheEntry represents a cached active signing key
type webkeyActiveSigningKeyCacheEntry struct {
	instanceID string
	webKey     *jose.JSONWebKey
}

func (e *webkeyActiveSigningKeyCacheEntry) Keys(index webkeyActiveSigningKeyCacheIndex) []string {
	switch index {
	case webkeyActiveSigningKeyCacheInstanceIndex:
		return []string{e.instanceID}
	default:
		return nil
	}
}

// webkeyPublicKeysCacheIndex represents the index for public keys cache
type webkeyPublicKeysCacheIndex int

const (
	webkeyPublicKeysCacheInstanceIndex webkeyPublicKeysCacheIndex = iota
)

// webkeyPublicKeysCacheEntry represents cached public keys
type webkeyPublicKeysCacheEntry struct {
	instanceID string
	keySet     *jose.JSONWebKeySet
}

func (e *webkeyPublicKeysCacheEntry) Keys(index webkeyPublicKeysCacheIndex) []string {
	switch index {
	case webkeyPublicKeysCacheInstanceIndex:
		return []string{e.instanceID}
	default:
		return nil
	}
}

func webkeyActiveSigningKeyCacheIndices() []webkeyActiveSigningKeyCacheIndex {
	return []webkeyActiveSigningKeyCacheIndex{webkeyActiveSigningKeyCacheInstanceIndex}
}

func webkeyPublicKeysCacheIndices() []webkeyPublicKeysCacheIndex {
	return []webkeyPublicKeysCacheIndex{webkeyPublicKeysCacheInstanceIndex}
}
