package repository

import (
	"context"

	"gopkg.in/square/go-jose.v2"
)

type KeyRepository interface {
	GenerateSigningKeyPair(ctx context.Context, algorithm string) error
	GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string)
	GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error)
}
