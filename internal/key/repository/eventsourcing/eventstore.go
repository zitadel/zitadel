package eventsourcing

import (
	"context"
	"strconv"

	"github.com/sony/sonyflake"

	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

type KeyEventstore struct {
	es_int.Eventstore
	keyCache    *KeyCache
	idGenerator id.Generator
}

type KeyConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartKey(conf KeyConfig, systemDefaults sd.SystemDefaults) (*KeyEventstore, error) {
	keyCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &KeyEventstore{
		Eventstore:  conf.Eventstore,
		keyCache:    keyCache,
		idGenerator: id.SonyFlakeGenerator,
	}, nil
}

func (es *KeyEventstore) KeyByID(ctx context.Context, id string) (*key_model.Key, error) {
	project := es.projectCache.getKey(id)

	query, err := KeyByIDQuery(project.AggregateID, project.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, project.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && project.Sequence != 0) {
		return nil, err
	}
	es.projectCache.cacheKey(project)
	return model.KeyToModel(project), nil
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

	es.keyCache.cacheKey(repoKey)
	return model.KeyPairToModel(repoKey), nil
}
