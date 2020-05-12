package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type IamEventstore struct {
	es_int.Eventstore
	iamCache *IamCache
}

type IamConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartIam(conf IamConfig, systemDefaults sd.SystemDefaults) (*IamEventstore, error) {
	iamCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}

	return &IamEventstore{
		Eventstore: conf.Eventstore,
		iamCache:   iamCache,
	}, nil
}

func (es *IamEventstore) IamByID(ctx context.Context, id string) (*iam_model.Iam, error) {
	iam := es.iamCache.getIam(id)

	query, err := IamByIDQuery(iam.AggregateID, iam.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	es.iamCache.cacheIam(iam)
	return model.IamToModel(iam), nil
}

func (es *IamEventstore) StartSetup(ctx context.Context, iamID string) (*iam_model.Iam, error) {
	iam, err := es.IamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}

	if iam != nil && iam.SetUpStarted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "Setup already started")
	}

	repoIam := &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: iamID}}
	createAggregate := IamSetupStartedAggregate(es.AggregateCreator(), repoIam)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIam(repoIam)
	return model.IamToModel(repoIam), nil
}

func (es *IamEventstore) SetupDone(ctx context.Context, iamID string) (*iam_model.Iam, error) {
	iam, err := es.IamByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	repoIam := model.IamFromModel(iam)
	createAggregate := IamSetupDoneAggregate(es.AggregateCreator(), repoIam)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIam(repoIam)
	return model.IamToModel(repoIam), nil
}

func (es *IamEventstore) SetGlobalOrg(ctx context.Context, iamID, globalOrg string) (*iam_model.Iam, error) {
	iam, err := es.IamByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	repoIam := model.IamFromModel(iam)
	createAggregate := IamSetGlobalOrgAggregate(es.AggregateCreator(), repoIam, globalOrg)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIam(repoIam)
	return model.IamToModel(repoIam), nil
}

func (es *IamEventstore) SetIamProject(ctx context.Context, iamID, iamProjectID string) (*iam_model.Iam, error) {
	iam, err := es.IamByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	repoIam := model.IamFromModel(iam)
	createAggregate := IamSetIamProjectAggregate(es.AggregateCreator(), repoIam, iamProjectID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIam(repoIam)
	return model.IamToModel(repoIam), nil
}

func (es *IamEventstore) AddIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-89osr", "UserID and Roles are required")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "User is already member of this Iam")
	}
	repoIam := model.IamFromModel(existing)
	repoMember := model.IamMemberFromModel(member)

	addAggregate := IamMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)

	if _, m := model.GetIamMember(repoIam.Members, member.UserID); m != nil {
		return model.IamMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-s90pw", "Could not find member in list")
}

func (es *IamEventstore) ChangeIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s9ipe", "UserID and Roles are required")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s7ucs", "User is not member of this project")
	}
	repoIam := model.IamFromModel(existing)
	repoMember := model.IamMemberFromModel(member)

	projectAggregate := IamMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, projectAggregate)
	es.iamCache.cacheIam(repoIam)

	if _, m := model.GetIamMember(repoIam.Members, member.UserID); m != nil {
		return model.IamMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-29cws", "Could not find member in list")
}

func (es *IamEventstore) RemoveIamMember(ctx context.Context, member *iam_model.IamMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-0pors", "UserID and Roles are required")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "User is not member of this project")
	}
	repoIam := model.IamFromModel(existing)
	repoMember := model.IamMemberFromModel(member)

	projectAggregate := IamMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, projectAggregate)
	es.iamCache.cacheIam(repoIam)
	return err
}
