package repository

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type KeyRepository interface {
	GenerateSigningKeyPair(ctx context.Context, algorithm string) error
	GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, timer <-chan time.Time)
	GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error)
}
