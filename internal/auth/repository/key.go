package repository

import (
	"context"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore/key"
	key_model "github.com/caos/zitadel/internal/key/model"
	"gopkg.in/square/go-jose.v2"
)

type KeyRepository interface {
	GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string, usage key_model.KeyUsage)
	GetCertificateAndKey(ctx context.Context, certAndKeyCh chan<- key.CertificateAndKey, algorithm string, usage key_model.KeyUsage)
	GetKeySet(ctx context.Context, usage key_model.KeyUsage) (*jose.JSONWebKeySet, error)
}
