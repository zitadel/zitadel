package login

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

const (
	// encKeyTTL is how long an ephemeral server private key is retained.
	// Auth requests typically expire in minutes; 10 minutes provides headroom
	// without accumulating stale keys.
	encKeyTTL = 10 * time.Minute
)

// encKeyEntry holds an ephemeral server-side private key and its expiry time.
type encKeyEntry struct {
	key     *ecdh.PrivateKey
	expires time.Time
}

// passwordEncKeyStore holds per-auth-request ephemeral ECDH private keys.
// Keys are single-use: retrieved and deleted in one atomic step.
// A background goroutine periodically evicts expired entries.
type passwordEncKeyStore struct {
	mu      sync.Mutex
	entries map[string]*encKeyEntry
	done    chan struct{}
}

func newPasswordEncKeyStore() *passwordEncKeyStore {
	s := &passwordEncKeyStore{
		entries: make(map[string]*encKeyEntry),
		done:    make(chan struct{}),
	}
	go s.evictLoop()
	return s
}

// generate creates a new P-256 ephemeral keypair for the given authRequestID,
// stores the private key, and returns the public key for embedding in the page.
// Calling generate twice for the same ID replaces the previous key.
func (s *passwordEncKeyStore) generate(authRequestID string) (*ecdh.PublicKey, error) {
	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ephemeral ECDH key: %w", err)
	}
	s.mu.Lock()
	s.entries[authRequestID] = &encKeyEntry{key: privKey, expires: time.Now().Add(encKeyTTL)}
	s.mu.Unlock()
	return privKey.PublicKey(), nil
}

// retrieve pops the private key for the given authRequestID.
// It returns an error if the key is absent or expired; in both cases the entry
// is removed so there is no second-attempt window.
func (s *passwordEncKeyStore) retrieve(authRequestID string) (*ecdh.PrivateKey, error) {
	s.mu.Lock()
	entry, ok := s.entries[authRequestID]
	delete(s.entries, authRequestID)
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no ephemeral key for auth request %q", authRequestID)
	}
	if time.Now().After(entry.expires) {
		return nil, fmt.Errorf("ephemeral key for auth request %q has expired", authRequestID)
	}
	return entry.key, nil
}

func (s *passwordEncKeyStore) evictLoop() {
	ticker := time.NewTicker(encKeyTTL)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.mu.Lock()
			for id, e := range s.entries {
				if now.After(e.expires) {
					delete(s.entries, id)
				}
			}
			s.mu.Unlock()
		case <-s.done:
			return
		}
	}
}

func (s *passwordEncKeyStore) close() { close(s.done) }
