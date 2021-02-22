package repository

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type KeyRepository interface {
	GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, timer <-chan time.Time)
	GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error)
}
