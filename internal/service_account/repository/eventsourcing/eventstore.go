package eventsourcing

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/service_account/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	svcacc_model "github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccountEventstore struct {
	es_int.Eventstore
	// userCache                *ServiceAccountCache
	idGenerator  id.Generator
	keySize      int
	keyAlgorithm crypto.EncryptionAlgorithm
}

type ServiceAccountConfig struct {
	es_int.Eventstore
}

func StartServiceAccount(conf ServiceAccountConfig) (*ServiceAccountEventstore, error) {
	// userCache, err := StartCache(conf.Cache)
	// if err != nil {
	// 	return nil, err
	// }

	return &ServiceAccountEventstore{
		Eventstore: conf.Eventstore,
	}, nil
}

func (es *ServiceAccountEventstore) ServiceAccountByID(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)

	query, err := ServiceAccountByIDQuery(id, 0)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, serviceAccount.AppendEvents, query)
	if err != nil && errors.IsNotFound(err) && serviceAccount.Sequence == 0 {
		return nil, err
	}
	return model.ServiceAccountToModel(serviceAccount), nil
}

func (es *ServiceAccountEventstore) ServiceAccountEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := ServiceAccountByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

func (es *ServiceAccountEventstore) CreateServiceAccount(ctx context.Context, account *svcacc_model.ServiceAccount) (*svcacc_model.ServiceAccount, error) {

	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	account.AggregateID = id

	serviceAccount := model.ServiceAccountFromModel(account)

	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-JHkFp", "Errors.Serviceaccount.Unimplemented")
}

// func (es *ServiceAccountEventstore) GenerateKeyPair(ctx context.Context, usage key_model.KeyUsage, algorithm string) (*key_model.KeyPair, error) {
// 	privateKey, publicKey, err := crypto.GenerateEncryptedKeyPair(es.keySize, es.keyAlgorithm)
// 	if err != nil {
// 		return nil, err
// 	}
// 	privateKeyExp := time.Now().UTC().Add(es.privateKeyLifetime)
// 	publicKeyExp := time.Now().UTC().Add(es.publicKeyLifetime)
// 	return es.CreateKeyPair(ctx, &key_model.KeyPair{
// 		ObjectRoot: models.ObjectRoot{},
// 		Usage:      usage,
// 		Algorithm:  algorithm,
// 		PrivateKey: &key_model.Key{
// 			Key:    privateKey,
// 			Expiry: privateKeyExp,
// 		},
// 		PublicKey: &key_model.Key{
// 			Key:    publicKey,
// 			Expiry: publicKeyExp,
// 		},
// 	})
// }

func (es *ServiceAccountEventstore) UpdateServiceAccount(ctx context.Context, account *svcacc_model.ServiceAccount) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := model.ServiceAccountFromModel(account)
	//TODO: update logic
	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-EVEMY", "Errors.Serviceaccount.Unimplemented")
}

func (es *ServiceAccountEventstore) DeactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)
	//TODO: deactivate logic
	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-kRhs6", "Errors.Serviceaccount.Unimplemented")
}

func (es *ServiceAccountEventstore) ReactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)
	//TODO: reactivate logic
	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-ff2em", "Errors.Serviceaccount.Unimplemented")
}

func (es *ServiceAccountEventstore) LockServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)
	//TODO: lock logic
	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-DIzFN", "Errors.Serviceaccount.Unimplemented")
}

func (es *ServiceAccountEventstore) UnlockServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)
	//TODO: unlock logic
	return model.ServiceAccountToModel(serviceAccount), errors.ThrowUnimplemented(nil, "EVENT-iIBKd", "Errors.Serviceaccount.Unimplemented")
}

func (es *ServiceAccountEventstore) DeleteServiceAccount(ctx context.Context, id string) error {
	//TODO: delete logic
	return errors.ThrowUnimplemented(nil, "EVENT-jg2ZX", "Errors.Serviceaccount.Unimplemented")
}
