package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type IAMEventstore struct {
	es_int.Eventstore
	iamCache     *IAMCache
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
}

type IAMConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartIAM(conf IAMConfig, systemDefaults sd.SystemDefaults) (*IAMEventstore, error) {
	iamCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}

	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	return &IAMEventstore{
		Eventstore:   conf.Eventstore,
		iamCache:     iamCache,
		idGenerator:  id.SonyFlakeGenerator,
		secretCrypto: aesCrypto,
	}, nil
}

func (es *IAMEventstore) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	iam := es.iamCache.getIAM(id)

	query, err := IAMByIDQuery(iam.AggregateID, iam.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	es.iamCache.cacheIAM(iam)
	return model.IAMToModel(iam), nil
}

func (es *IAMEventstore) IAMEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := IAMByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

//func (es *IAMEventstore) StartSetup(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
//	iam, err := es.IAMByID(ctx, iamID)
//	if err != nil && !caos_errs.IsNotFound(err) {
//		return nil, err
//	}
//
//	if iam != nil && (iam.SetUpStarted >= step || iam.SetUpStarted != iam.SetUpDone) {
//		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "Setup already started")
//	}
//
//	if iam == nil {
//		iam = &iam_model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: iamID}}
//	}
//	iam.SetUpStarted = step
//	repoIAM := model.IAMFromModel(iam)
//
//	createAggregate := IAMSetupStartedAggregate(es.AggregateCreator(), repoIAM)
//	err = es_sdk.Push(ctx, es.PushAggregates, repoIAM.AppendEvents, createAggregate)
//	if err != nil {
//		return nil, err
//	}
//
//	es.iamCache.cacheIAM(repoIAM)
//	return model.IAMToModel(repoIAM), nil
//}
//
//func (es *IAMEventstore) SetupDone(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
//	iam, err := es.IAMByID(ctx, iamID)
//	if err != nil {
//		return nil, err
//	}
//	iam.SetUpDone = step
//
//	repoIam := model.IAMFromModel(iam)
//	createAggregate := IAMSetupDoneAggregate(es.AggregateCreator(), repoIam)
//	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
//	if err != nil {
//		return nil, err
//	}
//
//	es.iamCache.cacheIAM(repoIam)
//	return model.IAMToModel(repoIam), nil
//}

func (es *IAMEventstore) PrepareSetupDone(ctx context.Context, iam *model.IAM, aggregate *models.Aggregate, step iam_model.Step) (*model.IAM, *models.Aggregate, func(ctx context.Context, aggregates ...*models.Aggregate) error, error) {
	iam.SetUpDone = model.Step(step)
	agg, err := IAMSetupDoneEvent(ctx, aggregate, iam)
	if err != nil {
		return nil, nil, nil, err
	}
	return iam, agg, es.PushAggregates, nil
}

func (es *IAMEventstore) SetGlobalOrg(ctx context.Context, iamID, globalOrg string) (*iam_model.IAM, error) {
	iam, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	repoIam := model.IAMFromModel(iam)
	createAggregate := IAMSetGlobalOrgAggregate(es.AggregateCreator(), repoIam, globalOrg)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIAM(repoIam)
	return model.IAMToModel(repoIam), nil
}

func (es *IAMEventstore) SetIAMProject(ctx context.Context, iamID, iamProjectID string) (*iam_model.IAM, error) {
	iam, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	repoIam := model.IAMFromModel(iam)
	createAggregate := IAMSetIamProjectAggregate(es.AggregateCreator(), repoIam, iamProjectID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIAM(repoIam)
	return model.IAMToModel(repoIam), nil
}

func (es *IAMEventstore) AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-89osr", "Errors.IAM.MemberInvalid")
	}
	existing, err := es.IAMByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.IAM.MemberAlreadyExisting")
	}
	repoIam := model.IAMFromModel(existing)
	repoMember := model.IAMMemberFromModel(member)

	addAggregate := IAMMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)

	if _, m := model.GetIAMMember(repoIam.Members, member.UserID); m != nil {
		return model.IAMMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-s90pw", "Errors.Internal")
}

func (es *IAMEventstore) ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s9ipe", "Errors.IAM.MemberInvalid")
	}
	existing, err := es.IAMByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s7ucs", "Errors.IAM.MemberNotExisting")
	}
	repoIam := model.IAMFromModel(existing)
	repoMember := model.IAMMemberFromModel(member)

	projectAggregate := IAMMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, projectAggregate)
	es.iamCache.cacheIAM(repoIam)

	if _, m := model.GetIAMMember(repoIam.Members, member.UserID); m != nil {
		return model.IAMMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-29cws", "Errors.Internal")
}

func (es *IAMEventstore) RemoveIAMMember(ctx context.Context, member *iam_model.IAMMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-0pors", "Errors.IAM.MemberInvalid")
	}
	existing, err := es.IAMByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.IAM.MemberNotExisting")
	}
	repoIAM := model.IAMFromModel(existing)
	repoMember := model.IAMMemberFromModel(member)

	projectAggregate := IAMMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoIAM, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIAM.AppendEvents, projectAggregate)
	es.iamCache.cacheIAM(repoIAM)
	return err
}

func (es *IAMEventstore) GetIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, existingIDP := existing.GetIDP(idpConfigID); existingIDP != nil {
		return existingIDP, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-Scj8s", "Errors.IAM.IdpNotExisting")
}

func (es *IAMEventstore) AddIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if idp == nil || !idp.IsValid(true) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Ms89d", "Errors.IAM.IdpInvalid")
	}
	existing, err := es.IAMByID(ctx, idp.AggregateID)
	if err != nil {
		return nil, err
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	idp.IDPConfigID = id

	if idp.OIDCConfig != nil {
		idp.OIDCConfig.IDPConfigID = id
		err = idp.OIDCConfig.CryptSecret(es.secretCrypto)
		if err != nil {
			return nil, err
		}
	}
	repoIam := model.IAMFromModel(existing)
	repoIdp := model.IDPConfigFromModel(idp)

	addAggregate := IDPConfigAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, idpConfig := model.GetIDPConfig(repoIam.IDPs, idp.IDPConfigID); idpConfig != nil {
		return model.IDPConfigToModel(idpConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Scj8s", "Errors.Internal")
}

func (es *IAMEventstore) ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if idp == nil || !idp.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Cms8o", "Errors.IAM.IdpInvalid")
	}
	existing, err := es.IAMByID(ctx, idp.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, existingIDP := existing.GetIDP(idp.IDPConfigID); existingIDP == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Cmlos", "Errors.IAM.IdpNotExisting")
	}
	repoIam := model.IAMFromModel(existing)
	repoIdp := model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, idpConfig := model.GetIDPConfig(repoIam.IDPs, idp.IDPConfigID); idpConfig != nil {
		return model.IDPConfigToModel(idpConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Xmlo0", "Errors.Internal")
}

func (es *IAMEventstore) PrepareRemoveIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*model.IAM, *models.Aggregate, error) {
	if idp == nil || idp.IDPConfigID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Wz7sD", "Errors.IAM.IDMissing")
	}
	existing, err := es.IAMByID(ctx, idp.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, existingIDP := existing.GetIDP(idp.IDPConfigID); existingIDP == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Smiu8", "Errors.IAM.IdpNotExisting")
	}
	repoIam := model.IAMFromModel(existing)
	repoIdp := model.IDPConfigFromModel(idp)
	provider := new(model.IDPProvider)
	if repoIam.DefaultLoginPolicy != nil {
		_, provider = model.GetIDPProvider(repoIam.DefaultLoginPolicy.IDPProviders, idp.IDPConfigID)
	}
	agg, err := IDPConfigRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoIdp, provider)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, agg, nil
}

func (es *IAMEventstore) RemoveIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) error {
	repoIam, agg, err := es.PrepareRemoveIDPConfig(ctx, idp)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.iamCache.cacheIAM(repoIam)
	return nil
}

func (es *IAMEventstore) DeactivateIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	if idpID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Fbs8o", "Errors.IAM.IDMissing")
	}
	existing, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IDPConfig{IDPConfigID: idpID}
	if _, existingIDP := existing.GetIDP(idp.IDPConfigID); existingIDP == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Mci32", "Errors.IAM.IdpNotExisting")
	}
	repoIam := model.IAMFromModel(existing)
	repoIdp := model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, idpConfig := model.GetIDPConfig(repoIam.IDPs, idp.IDPConfigID); idpConfig != nil {
		return model.IDPConfigToModel(idpConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Xnc8d", "Errors.Internal")
}

func (es *IAMEventstore) ReactivateIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	if idpID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Wkjsf", "Errors.IAM.IDMissing")
	}
	iam, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IDPConfig{IDPConfigID: idpID}
	if _, existingIDP := iam.GetIDP(idp.IDPConfigID); existingIDP == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Sjc78", "Errors.IAM.IdpNotExisting")
	}
	repoIam := model.IAMFromModel(iam)
	repoIdp := model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigReactivatedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, idpConfig := model.GetIDPConfig(repoIam.IDPs, idp.IDPConfigID); idpConfig != nil {
		return model.IDPConfigToModel(idpConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Snd4f", "Errors.Internal")
}

func (es *IAMEventstore) ChangeIDPOIDCConfig(ctx context.Context, config *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	if config == nil || !config.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-*5ki8", "Errors.IAM.OIDCConfigInvalid")
	}
	iam, err := es.IAMByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var idp *iam_model.IDPConfig
	if _, idp = iam.GetIDP(config.IDPConfigID); idp == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-pso0s", "Errors.IAM.IdpNoExisting")
	}
	if idp.Type != iam_model.IDPConfigTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Fms8w", "Errors.IAM.IdpIsNotOIDC")
	}
	if config.ClientSecretString != "" {
		err = config.CryptSecret(es.secretCrypto)
		if err != nil {
			return nil, err
		}
	} else {
		config.ClientSecret = nil
	}
	repoIam := model.IAMFromModel(iam)
	repoConfig := model.OIDCIDPConfigFromModel(config)

	iamAggregate := OIDCIDPConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, idpConfig := model.GetIDPConfig(repoIam.IDPs, idp.IDPConfigID); idpConfig != nil {
		return model.OIDCIDPConfigToModel(idpConfig.OIDCIDPConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Sldk8", "Errors.Internal")
}

func (es *IAMEventstore) PrepareAddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-VwlDv", "Errors.IAM.LabelPolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	labelPolicy := model.LabelPolicyFromModel(policy)

	addAggregate := LabelPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, labelPolicy)
	aggregate, err := addAggregate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, aggregate, nil
}

func (es *IAMEventstore) AddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-aAPWI", "Errors.IAM.LabelPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLabelPolicy := model.LabelPolicyFromModel(policy)

	addAggregate := LabelPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLabelPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LabelPolicyToModel(repoIam.DefaultLabelPolicy), nil
}

func (es *IAMEventstore) ChangeLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-vRqjg", "Errors.IAM.LabelPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLabelPolicy := model.LabelPolicyFromModel(policy)

	addAggregate := LabelPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLabelPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LabelPolicyToModel(repoIam.DefaultLabelPolicy), nil
}

func (es *IAMEventstore) PrepareAddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || !policy.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-3mP0s", "Errors.IAM.LoginPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	addAggregate, err := LoginPolicyAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	repoIam, addAggregate, err := es.PrepareAddLoginPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IAMEventstore) PrepareChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || !policy.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-3M0so", "Errors.IAM.LoginPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	changeAgg, err := LoginPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)(ctx)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, changeAgg, nil
}

func (es *IAMEventstore) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	repoIam, changeAggregate, err := es.PrepareChangeLoginPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, changeAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IAMEventstore) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	if provider == nil || !provider.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-bMS8i", "Errors.IdpProviderInvalid")
	}
	iam, err := es.IAMByID(ctx, provider.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := iam.DefaultLoginPolicy.GetIdpProvider(provider.IDPConfigID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.IAM.LoginPolicy.IdpProviderAlreadyExisting")
	}
	repoIam := model.IAMFromModel(iam)
	repoProvider := model.IDPProviderFromModel(provider)

	addAggregate := LoginPolicyIDPProviderAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoProvider)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	if _, m := model.GetIDPProvider(repoIam.DefaultLoginPolicy.IDPProviders, provider.IDPConfigID); m != nil {
		return model.IDPProviderToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Slf9s", "Errors.Internal")
}

func (es *IAMEventstore) PrepareRemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*model.IAM, *models.Aggregate, error) {
	if provider == nil || !provider.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Esi8c", "Errors.IdpProviderInvalid")
	}
	iam, err := es.IAMByID(ctx, provider.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, m := iam.DefaultLoginPolicy.GetIdpProvider(provider.IDPConfigID); m == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.IAM.LoginPolicy.IdpProviderNotExisting")
	}
	repoIam := model.IAMFromModel(iam)
	removeAgg, err := LoginPolicyIDPProviderRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, &model.IDPProviderID{provider.IDPConfigID})
	if err != nil {
		return nil, nil, err
	}
	return repoIam, removeAgg, nil
}

func (es *IAMEventstore) RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error {
	repoIam, removeAgg, err := es.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, removeAgg)
	if err != nil {
		return err
	}
	es.iamCache.cacheIAM(repoIam)
	return nil
}

func (es *IAMEventstore) AddSecondFactorToLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	repoIAM, addAggregate, err := es.PrepareAddSecondFactorToLoginPolicy(ctx, aggregateID, mfa)
	if err != nil {
		return 0, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIAM.AppendEvents, addAggregate)
	if err != nil {
		return 0, err
	}
	es.iamCache.cacheIAM(repoIAM)
	if _, m := model.GetMFA(repoIAM.DefaultLoginPolicy.SecondFactors, int32(mfa)); m != 0 {
		return iam_model.SecondFactorType(m), nil
	}
	return 0, caos_errs.ThrowInternal(nil, "EVENT-5N9so", "Errors.Internal")
}

func (es *IAMEventstore) PrepareAddSecondFactorToLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.SecondFactorType) (*model.IAM, *models.Aggregate, error) {
	if mfa == iam_model.SecondFactorTypeUnspecified {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-1M8Js", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	iam, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, m := iam.DefaultLoginPolicy.GetSecondFactor(mfa); m != 0 {
		return nil, nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-4Rk09", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}
	repoIAM := model.IAMFromModel(iam)
	repoMFA := model.SecondFactorFromModel(mfa)

	addAggregate := LoginPolicySecondFactorAddedAggregate(es.Eventstore.AggregateCreator(), repoIAM, repoMFA)
	aggregate, err := addAggregate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return repoIAM, aggregate, nil
}

func (es *IAMEventstore) RemoveSecondFactorFromLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.SecondFactorType) error {
	if mfa == iam_model.SecondFactorTypeUnspecified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-4gJ9s", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	iam, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return err
	}
	if _, m := iam.DefaultLoginPolicy.GetSecondFactor(mfa); m == 0 {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-gBm9s", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	repoIam := model.IAMFromModel(iam)
	repoMFA := model.SecondFactorFromModel(mfa)

	removeAgg := LoginPolicySecondFactorRemovedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMFA)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, removeAgg)
	if err != nil {
		return err
	}
	es.iamCache.cacheIAM(repoIam)
	return nil
}

func (es *IAMEventstore) PrepareAddMultiFactorToLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.MultiFactorType) (*model.IAM, *models.Aggregate, error) {
	if mfa == iam_model.MultiFactorTypeUnspecified {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-2Dh7J", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	iam, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, m := iam.DefaultLoginPolicy.GetMultiFactor(mfa); m != 0 {
		return nil, nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-4Rk09", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}
	repoIam := model.IAMFromModel(iam)
	repoMFA := model.MultiFactorFromModel(mfa)

	addAggregate, err := LoginPolicyMultiFactorAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMFA)(ctx)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddMultiFactorToLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	repoIAM, addAggregate, err := es.PrepareAddMultiFactorToLoginPolicy(ctx, aggregateID, mfa)
	if err != nil {
		return 0, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIAM.AppendEvents, addAggregate)
	if err != nil {
		return 0, err
	}
	es.iamCache.cacheIAM(repoIAM)
	if _, m := model.GetMFA(repoIAM.DefaultLoginPolicy.MultiFactors, int32(mfa)); m != 0 {
		return iam_model.MultiFactorType(m), nil
	}
	return 0, caos_errs.ThrowInternal(nil, "EVENT-5N9so", "Errors.Internal")
}

func (es *IAMEventstore) RemoveMultiFactorFromLoginPolicy(ctx context.Context, aggregateID string, mfa iam_model.MultiFactorType) error {
	if mfa == iam_model.MultiFactorTypeUnspecified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-4gJ9s", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	iam, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return err
	}
	if _, m := iam.DefaultLoginPolicy.GetMultiFactor(mfa); m == 0 {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-gBm9s", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	repoIam := model.IAMFromModel(iam)
	repoMFA := model.MultiFactorFromModel(mfa)

	removeAgg := LoginPolicyMultiFactorRemovedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMFA)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, removeAgg)
	if err != nil {
		return err
	}
	es.iamCache.cacheIAM(repoIam)
	return nil
}

func (es *IAMEventstore) PrepareAddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Ks8Fs", "Errors.IAM.PasswordComplexityPolicy.Empty")
	}
	if err := policy.IsValid(); err != nil {
		return nil, nil, err
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordComplexityPolicy := model.PasswordComplexityPolicyFromModel(policy)

	addAggregate, err := PasswordComplexityPolicyAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoPasswordComplexityPolicy)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	repoIam, addAggregate, err := es.PrepareAddPasswordComplexityPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordComplexityPolicyToModel(repoIam.DefaultPasswordComplexityPolicy), nil
}

func (es *IAMEventstore) ChangePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	if policy == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.IAM.PasswordComplexityPolicy.Empty")
	}
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordComplexityPolicy := model.PasswordComplexityPolicyFromModel(policy)

	addAggregate := PasswordComplexityPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoPasswordComplexityPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordComplexityPolicyToModel(repoIam.DefaultPasswordComplexityPolicy), nil
}

func (es *IAMEventstore) PrepareAddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-2dGt6", "Errors.IAM.PasswordAgePolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordAgePolicy := model.PasswordAgePolicyFromModel(policy)

	addAggregate, err := PasswordAgePolicyAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoPasswordAgePolicy)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	repoIam, addAggregate, err := es.PrepareAddPasswordAgePolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordAgePolicyToModel(repoIam.DefaultPasswordAgePolicy), nil
}

func (es *IAMEventstore) ChangePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-2Fgt6", "Errors.IAM.PasswordAgePolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordAgePolicy := model.PasswordAgePolicyFromModel(policy)

	addAggregate := PasswordAgePolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoPasswordAgePolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordAgePolicyToModel(repoIam.DefaultPasswordAgePolicy), nil
}

func (es *IAMEventstore) PrepareAddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-3R56z", "Errors.IAM.PasswordLockoutPolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordLockoutPolicy := model.PasswordLockoutPolicyFromModel(policy)

	addAggregate, err := PasswordLockoutPolicyAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoPasswordLockoutPolicy)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	repoIam, addAggregate, err := es.PrepareAddPasswordLockoutPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordLockoutPolicyToModel(repoIam.DefaultPasswordLockoutPolicy), nil
}

func (es *IAMEventstore) ChangePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-6Zsj9", "Errors.IAM.PasswordLockoutPolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoPasswordLockoutPolicy := model.PasswordLockoutPolicyFromModel(policy)

	addAggregate := PasswordLockoutPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoPasswordLockoutPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.PasswordLockoutPolicyToModel(repoIam.DefaultPasswordLockoutPolicy), nil
}

func (es *IAMEventstore) GetOrgIAMPolicy(ctx context.Context, iamID string) (*iam_model.OrgIAMPolicy, error) {
	existingIAM, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	if existingIAM.DefaultOrgIAMPolicy == nil {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2Fj8s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return existingIAM.DefaultOrgIAMPolicy, nil
}

func (es *IAMEventstore) PrepareAddOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*model.IAM, *models.Aggregate, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-3R56z", "Errors.IAM.OrgIAMPolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoOrgIAMPolicy := model.OrgIAMPolicyFromModel(policy)

	addAggregate, err := OrgIAMPolicyAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoOrgIAMPolicy)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, addAggregate, nil
}

func (es *IAMEventstore) AddOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	repoIam, addAggregate, err := es.PrepareAddOrgIAMPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.OrgIAMPolicyToModel(repoIam.DefaultOrgIAMPolicy), nil
}

func (es *IAMEventstore) ChangeOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-6Zsj9", "Errors.IAM.OrgIAMPolicy.Empty")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoOrgIAMPolicy := model.OrgIAMPolicyFromModel(policy)

	addAggregate := OrgIAMPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoOrgIAMPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.OrgIAMPolicyToModel(repoIam.DefaultOrgIAMPolicy), nil
}
