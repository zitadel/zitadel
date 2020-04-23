package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/sony/sonyflake"
	"strconv"
)

type UserEventstore struct {
	es_int.Eventstore
	userCache                *UserCache
	idGenerator              *sonyflake.Sonyflake
	PasswordAlg              crypto.HashAlgorithm
	InitializeUserCode       crypto.Generator
	EmailVerificationCode    crypto.Generator
	PhoneVerificationCode    crypto.Generator
	PasswordVerificationCode crypto.Generator
}

type UserConfig struct {
	es_int.Eventstore
	Cache            *config.CacheConfig
	PasswordSaltCost int
}

func StartUser(conf UserConfig, systemDefaults sd.SystemDefaults) (*UserEventstore, error) {
	userCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	idGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	initCodeGen := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, aesCrypto)
	emailVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, aesCrypto)
	phoneVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, aesCrypto)
	passwordVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, aesCrypto)

	return &UserEventstore{
		Eventstore:               conf.Eventstore,
		userCache:                userCache,
		idGenerator:              idGenerator,
		InitializeUserCode:       initCodeGen,
		EmailVerificationCode:    emailVerificationCode,
		PhoneVerificationCode:    phoneVerificationCode,
		PasswordVerificationCode: passwordVerificationCode,
	}, nil
}

func (es *UserEventstore) UserByID(ctx context.Context, id string) (*usr_model.User, error) {
	user := es.userCache.getUser(id)

	query, err := UserByIDQuery(user.AggregateID, user.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, user.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && user.Sequence != 0) {
		return nil, err
	}
	es.userCache.cacheUser(user)
	return model.UserToModel(user), nil
}

func (es *UserEventstore) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	if user.Profile != nil && user.UserName == "" && user.Email != nil {
		user.UserName = user.EmailAddress
	}
	if !user.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	//TODO: Check Uniqueness
	id, err := es.idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	user.AggregateID = strconv.FormatUint(id, 10)
	if user.Password != nil && user.SecretString != "" {
		secret, err := crypto.Hash([]byte(user.SecretString), es.PasswordAlg)
		if err != nil {
			return nil, err
		}
		user.Password.SecretCrypto = secret
		user.Password.ChangeRequired = true
	}

	repoUser := model.UserFromModel(user)
	initCode := new(model.InitUserCode)
	if user.Email == nil || !user.IsEmailVerified || user.Password == nil || user.SecretString == "" {
		initCodeCrypto, _, err := crypto.NewCode(es.InitializeUserCode)
		if err != nil {
			return nil, err
		}
		initCode.Code = initCodeCrypto
		initCode.Expiry = es.InitializeUserCode.Expiry()
	}

	phoneCode := new(model.PhoneCode)
	if user.Phone != nil && !user.IsPhoneVerified {
		phoneCodeCrypto, _, err := crypto.NewCode(es.InitializeUserCode)
		if err != nil {
			return nil, err
		}
		phoneCode.Code = phoneCodeCrypto
		phoneCode.Expiry = es.PhoneVerificationCode.Expiry()
	}

	createAggregate := UserCreateAggregate(es.AggregateCreator(), repoUser, initCode, phoneCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoUser)
	return model.UserToModel(repoUser), nil
}

func (es *UserEventstore) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	if user.Profile != nil && user.UserName == "" && user.Email != nil {
		user.UserName = user.EmailAddress
	}
	if !user.IsValid() || user.Password == nil || user.SecretString == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "user is invalid")
	}
	//TODO: Check Uniqueness
	id, err := es.idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	user.AggregateID = strconv.FormatUint(id, 10)

	secret, err := crypto.Hash([]byte(user.SecretString), es.PasswordAlg)
	if err != nil {
		return nil, err
	}
	user.Password = &usr_model.Password{SecretCrypto: secret, ChangeRequired: false}

	repoUser := model.UserFromModel(user)

	emailCode := new(model.EmailCode)
	if user.Email != nil && !user.IsEmailVerified {
		emailCodeCrypto, _, err := crypto.NewCode(es.EmailVerificationCode)
		if err != nil {
			return nil, err
		}
		emailCode.Code = emailCodeCrypto
		emailCode.Expiry = es.EmailVerificationCode.Expiry()
	}

	createAggregate := UserRegisterAggregate(es.AggregateCreator(), repoUser, resourceOwner, emailCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoUser)
	return model.UserToModel(repoUser), nil
}

func (es *UserEventstore) DeactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	existing, err := es.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.IsInactive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "cant deactivate inactive user")
	}

	repoExisting := model.UserFromModel(existing)
	aggregate := UserDeactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoExisting)
	return model.UserToModel(repoExisting), nil
}

func (es *UserEventstore) ReactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	existing, err := es.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existing.IsInactive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do94s", "user must be inactive")
	}

	repoExisting := model.UserFromModel(existing)
	aggregate := UserReactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoExisting)
	return model.UserToModel(repoExisting), nil
}

func (es *UserEventstore) LockUser(ctx context.Context, id string) (*usr_model.User, error) {
	existing, err := es.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existing.IsActive() && !existing.IsInitial() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83s", "user must be active or initial")
	}

	repoExisting := model.UserFromModel(existing)
	aggregate := UserLockAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoExisting)
	return model.UserToModel(repoExisting), nil
}

func (es *UserEventstore) UnlockUser(ctx context.Context, id string) (*usr_model.User, error) {
	existing, err := es.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existing.IsLocked() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dks83", "user must be locked")
	}

	repoExisting := model.UserFromModel(existing)
	aggregate := UserUnlockAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoExisting)
	return model.UserToModel(repoExisting), nil
}

func (es *UserEventstore) InitializeUserCodeByID(ctx context.Context, userID string) (*usr_model.InitUserCode, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d8diw", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.InitCode != nil {
		return user.InitCode, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-d8e2", "init code not found")
}

func (es *UserEventstore) CreateInitializeUserCodeByID(ctx context.Context, userID string) (*usr_model.InitUserCode, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	initCode := new(model.InitUserCode)
	initCodeCrypto, _, err := crypto.NewCode(es.InitializeUserCode)
	if err != nil {
		return nil, err
	}
	initCode.Code = initCodeCrypto
	initCode.Expiry = es.InitializeUserCode.Expiry()

	repoUser := model.UserFromModel(user)
	agg := UserInitCodeAggregate(es.AggregateCreator(), repoUser, initCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoUser)
	return model.InitCodeToModel(repoUser.InitCode), nil
}

func (es *UserEventstore) SkipMfaInit(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := SkipMfaAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}
