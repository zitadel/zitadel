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

type IamEventstore struct {
	es_int.Eventstore
	iamCache     *IamCache
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
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

	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	return &IamEventstore{
		Eventstore:   conf.Eventstore,
		iamCache:     iamCache,
		idGenerator:  id.SonyFlakeGenerator,
		secretCrypto: aesCrypto,
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-89osr", "Errors.Iam.MemberInvalid")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.Iam.MemberAlreadyExisting")
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-s90pw", "Errors.Internal")
}

func (es *IamEventstore) ChangeIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s9ipe", "Errors.Iam.MemberInvalid")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s7ucs", "Errors.Iam.MemberNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoMember := model.IamMemberFromModel(member)

	projectAggregate := IamMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, projectAggregate)
	es.iamCache.cacheIam(repoIam)

	if _, m := model.GetIamMember(repoIam.Members, member.UserID); m != nil {
		return model.IamMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-29cws", "Errors.Internal")
}

func (es *IamEventstore) RemoveIamMember(ctx context.Context, member *iam_model.IamMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-0pors", "Errors.Iam.MemberInvalid")
	}
	existing, err := es.IamByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.Iam.MemberNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoMember := model.IamMemberFromModel(member)

	projectAggregate := IamMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, projectAggregate)
	es.iamCache.cacheIam(repoIam)
	return err
}

func (es *IamEventstore) GetIdpConfiguration(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IdpConfig, error) {
	existing, err := es.IamByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-Scj8s", "Errors.Iam.IdpNotExisting")
}

func (es *IamEventstore) AddIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	if idp == nil || !idp.IsValid(true) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Ms89d", "Errors.Iam.IdpInvalid")
	}
	existing, err := es.IamByID(ctx, idp.AggregateID)
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
	repoIam := model.IamFromModel(existing)
	repoIdp := model.IdpConfigFromModel(idp)

	addAggregate := IdpConfigurationAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, i := model.GetIdpConfig(repoIam.IDPs, idp.IDPConfigID); i != nil {
		return model.IdpConfigToModel(i), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Scj8s", "Errors.Internal")
}

func (es *IamEventstore) ChangeIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	if idp == nil || !idp.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Cms8o", "Errors.Iam.IdpInvalid")
	}
	existing, err := es.IamByID(ctx, idp.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Cmlos", "Errors.Iam.IdpNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoIdp := model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, i := model.GetIdpConfig(repoIam.IDPs, idp.IDPConfigID); i != nil {
		return model.IdpConfigToModel(i), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Xmlo0", "Errors.Internal")
}

func (es *IamEventstore) PrepareRemoveIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*model.Iam, *models.Aggregate, error) {
	if idp == nil || idp.IDPConfigID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Wz7sD", "Errors.Iam.IDMissing")
	}
	existing, err := es.IamByID(ctx, idp.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Smiu8", "Errors.Iam.IdpNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoIdp := model.IdpConfigFromModel(idp)
	provider := new(model.IdpProvider)
	if repoIam.DefaultLoginPolicy != nil {
		_, provider = model.GetIdpProvider(repoIam.DefaultLoginPolicy.IdpProviders, idp.IDPConfigID)
	}
	agg, err := IdpConfigurationRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoIam, repoIdp, provider)
	if err != nil {
		return nil, nil, err
	}
	return repoIam, agg, nil
}

func (es *IamEventstore) RemoveIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) error {
	repoIam, agg, err := es.PrepareRemoveIdpConfiguration(ctx, idp)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoIam.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.iamCache.cacheIam(repoIam)
	return nil
}

func (es *IamEventstore) DeactivateIdpConfiguration(ctx context.Context, iamID, idpID string) (*iam_model.IdpConfig, error) {
	if idpID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Fbs8o", "Errors.Iam.IDMissing")
	}
	existing, err := es.IamByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IdpConfig{IDPConfigID: idpID}
	if _, app := existing.GetIDP(idp.IDPConfigID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Mci32", "Errors.Iam.IdpNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoIdp := model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, i := model.GetIdpConfig(repoIam.IDPs, idp.IDPConfigID); i != nil {
		return model.IdpConfigToModel(i), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Xnc8d", "Errors.Internal")
}

func (es *IamEventstore) ReactivateIdpConfiguration(ctx context.Context, iamID, idpID string) (*iam_model.IdpConfig, error) {
	if idpID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Wkjsf", "Errors.Iam.IDMissing")
	}
	existing, err := es.IamByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IdpConfig{IDPConfigID: idpID}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Sjc78", "Errors.Iam.IdpNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoIdp := model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationReactivatedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, i := model.GetIdpConfig(repoIam.IDPs, idp.IDPConfigID); i != nil {
		return model.IdpConfigToModel(i), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Snd4f", "Errors.Internal")
}

func (es *IamEventstore) ChangeIdpOidcConfiguration(ctx context.Context, config *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error) {
	if config == nil || !config.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-*5ki8", "Errors.Iam.OIDCConfigInvalid")
	}
	existing, err := es.IamByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var idp *iam_model.IdpConfig
	if _, idp = existing.GetIDP(config.IDPConfigID); idp == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-pso0s", "Errors.Iam.IdpNoExisting")
	}
	if idp.Type != iam_model.IDPConfigTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Fms8w", "Errors.Iam.IdpIsNotOIDC")
	}
	if config.ClientSecretString != "" {
		err = idp.OIDCConfig.CryptSecret(es.secretCrypto)
	} else {
		config.ClientSecret = nil
	}
	repoIam := model.IamFromModel(existing)
	repoConfig := model.OidcIdpConfigFromModel(config)

	iamAggregate := OIDCIdpConfigurationChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, a := model.GetIdpConfig(repoIam.IDPs, idp.IDPConfigID); a != nil {
		return model.OidcIdpConfigToModel(a.OIDCIDPConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Sldk8", "Errors.Internal")
}

func (es *IamEventstore) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.Iam.LoginPolicyInvalid")
	}
	existing, err := es.IamByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IamFromModel(existing)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IamEventstore) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.Iam.LoginPolicyInvalid")
	}
	existing, err := es.IamByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	repoIam := model.IamFromModel(existing)
	repoLoginPolicy := model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	return model.LoginPolicyToModel(repoIam.DefaultLoginPolicy), nil
}

func (es *IamEventstore) AddIdpProviderToLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) (*iam_model.IdpProvider, error) {
	if provider == nil || !provider.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.IdpProviderInvalid")
	}
	existing, err := es.IamByID(ctx, provider.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.DefaultLoginPolicy.GetIdpProvider(provider.IdpConfigID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.Iam.LoginPolicy.IdpProviderAlreadyExisting")
	}
	repoIam := model.IamFromModel(existing)
	repoProvider := model.IdpProviderFromModel(provider)

	addAggregate := LoginPolicyIdpProviderAddedAggregate(es.Eventstore.AggregateCreator(), repoIam, repoProvider)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.iamCache.cacheIam(repoIam)
	if _, m := model.GetIdpProvider(repoIam.DefaultLoginPolicy.IdpProviders, provider.IdpConfigID); m != nil {
		return model.IdpProviderToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-Slf9s", "Errors.Internal")
}

func (es *IamEventstore) RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) error {
	if provider == nil || !provider.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Esi8c", "Errors.IdpProviderInvalid")
	}
	existing, err := es.IamByID(ctx, provider.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.DefaultLoginPolicy.GetIdpProvider(provider.IdpConfigID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.Iam.LoginPolicy.IdpProviderNotExisting")
	}
	repoIam := model.IamFromModel(existing)
	addAggregate := LoginPolicyIdpProviderRemovedAggregate(es.Eventstore.AggregateCreator(), repoIam, &model.IdpProviderID{provider.IdpConfigID})
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, addAggregate)
	if err != nil {
		return err
	}
	es.iamCache.cacheIam(repoIam)
	return nil
}
