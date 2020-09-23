package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/id"
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

func (es *IAMEventstore) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
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

func (es *IAMEventstore) StartSetup(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
	iam, err := es.IAMByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}

	if iam != nil && (iam.SetUpStarted >= step || iam.SetUpStarted != iam.SetUpDone) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "Setup already started")
	}

	repoIam := &model.IAM{ObjectRoot: iam.ObjectRoot, SetUpStarted: model.Step(step)}
	createAggregate := IAMSetupStartedAggregate(es.AggregateCreator(), repoIam)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIAM(repoIam)
	return model.IAMToModel(repoIam), nil
}

func (es *IAMEventstore) SetupDone(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
	iam, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	iam.SetUpDone = step

	repoIam := model.IAMFromModel(iam)
	createAggregate := IAMSetupDoneAggregate(es.AggregateCreator(), repoIam)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIAM(repoIam)
	return model.IAMToModel(repoIam), nil
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

func (es *IAMEventstore) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.IAM.LoginPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IAMEventstore) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.IAM.LoginPolicyInvalid")
	}
	iam, err := es.IAMByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IAMFromModel(iam)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIAM(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IAMEventstore) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	if provider == nil || !provider.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.IdpProviderInvalid")
	}
	iam, err := es.IAMByID(ctx, provider.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := iam.DefaultLoginPolicy.GetIdpProvider(provider.IdpConfigID); m != nil {
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
	if _, m := model.GetIDPProvider(repoIam.DefaultLoginPolicy.IDPProviders, provider.IdpConfigID); m != nil {
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
	if _, m := iam.DefaultLoginPolicy.GetIdpProvider(provider.IdpConfigID); m == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.IAM.LoginPolicy.IdpProviderNotExisting")
	}
	repoIam := model.IAMFromModel(iam)
	removeAgg, err := LoginPolicyIDPProviderRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, &model.IDPProviderID{provider.IdpConfigID})
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
