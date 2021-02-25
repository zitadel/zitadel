package eventsourcing

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
)

type KeyEventstore struct {
	es_int.Eventstore
	keySize            int
	keyAlgorithm       crypto.EncryptionAlgorithm
	privateKeyLifetime time.Duration
	publicKeyLifetime  time.Duration
	idGenerator        id.Generator
}

type KeyConfig struct {
	Size                     int
	PrivateKeyLifetime       types.Duration
	PublicKeyLifetime        types.Duration
	EncryptionConfig         *crypto.KeyConfig
	SigningKeyRotationCheck  types.Duration
	SigningKeyGracefulPeriod types.Duration
}

func StartKey(eventstore es_int.Eventstore, config KeyConfig, keyAlgorithm crypto.EncryptionAlgorithm, generator id.Generator) (*KeyEventstore, error) {
	return &KeyEventstore{
		Eventstore:         eventstore,
		keySize:            config.Size,
		keyAlgorithm:       keyAlgorithm,
		privateKeyLifetime: config.PrivateKeyLifetime.Duration,
		publicKeyLifetime:  config.PublicKeyLifetime.Duration,
		idGenerator:        generator,
	}, nil
}

func (es *KeyEventstore) GenerateKeyPair(ctx context.Context, usage key_model.KeyUsage, algorithm string) (*key_model.KeyPair, error) {
	privateKey, publicKey, err := crypto.GenerateEncryptedKeyPair(es.keySize, es.keyAlgorithm)
	if err != nil {
		return nil, err
	}
	privateKeyExp := time.Now().UTC().Add(es.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(es.publicKeyLifetime)
	return es.CreateKeyPair(ctx, &key_model.KeyPair{
		ObjectRoot: models.ObjectRoot{},
		Usage:      usage,
		Algorithm:  algorithm,
		PrivateKey: &key_model.Key{
			Key:    privateKey,
			Expiry: privateKeyExp,
		},
		PublicKey: &key_model.Key{
			Key:    publicKey,
			Expiry: publicKeyExp,
		},
	})
}

func (es *KeyEventstore) CreateKeyPair(ctx context.Context, pair *key_model.KeyPair) (*key_model.KeyPair, error) {
	if !pair.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-G34ga", "Name is required")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	pair.AggregateID = id
	repoKey := model.KeyPairFromModel(pair)

	createAggregate := KeyPairCreateAggregate(es.AggregateCreator(), repoKey)
	err = es_sdk.Push(ctx, es.PushAggregates, repoKey.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}
	return model.KeyPairToModel(repoKey), nil
}

func (es *KeyEventstore) LatestKeyEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	return es.FilterEvents(ctx, KeyPairQuery(sequence))
}
