package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/pquerna/otp/totp"
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
	Multifactors             global_model.Multifactors
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
	initCodeGen := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.InitializeUserCode, aesCrypto)
	emailVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.EmailVerificationCode, aesCrypto)
	phoneVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.PhoneVerificationCode, aesCrypto)
	passwordVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.PasswordVerificationCode, aesCrypto)
	aesOtpCrypto, err := crypto.NewAESCrypto(systemDefaults.Multifactors.OTP.VerificationKey)
	if err != nil {
		return nil, err
	}
	mfa := global_model.Multifactors{
		OTP: global_model.OTP{
			CryptoMFA: aesOtpCrypto,
			Issuer:    systemDefaults.Multifactors.OTP.Issuer,
		},
	}
	return &UserEventstore{
		Eventstore:               conf.Eventstore,
		userCache:                userCache,
		idGenerator:              idGenerator,
		InitializeUserCode:       initCodeGen,
		EmailVerificationCode:    emailVerificationCode,
		PhoneVerificationCode:    phoneVerificationCode,
		PasswordVerificationCode: passwordVerificationCode,
		Multifactors:             mfa,
	}, nil
}

func (es *UserEventstore) UserByID(ctx context.Context, id string) (*usr_model.User, error) {
	user := es.userCache.getUser(id)

	query, err := UserByIDQuery(user.AggregateID, user.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, user.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && user.Sequence == 0 {
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

	initCode := new(model.InitUserCode)
	if user.Email == nil || !user.IsEmailVerified || user.Password == nil || user.SecretString == "" {
		err := es.generateInitUserCode(initCode)
		if err != nil {
			return nil, err
		}
	}

	phoneCode := new(model.PhoneCode)
	if user.Phone != nil && !user.IsPhoneVerified {
		err := es.generatePhoneCode(phoneCode)
		if err != nil {
			return nil, err
		}
	}

	repoUser := model.UserFromModel(user)
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

	emailCode := new(model.EmailCode)
	if user.Email != nil && !user.IsEmailVerified {
		err := es.generateEmailCode(emailCode)
		if err != nil {
			return nil, err
		}
	}

	repoUser := model.UserFromModel(user)
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
	err = es.generateInitUserCode(initCode)
	if err != nil {
		return nil, err
	}

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

func (es *UserEventstore) UserPasswordByID(ctx context.Context, userID string) (*usr_model.Password, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Password != nil {
		return user.Password, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-d8e2", "password not found")
}

func (es *UserEventstore) CheckPassword(ctx context.Context, userID, password, authRequestID string) (*usr_model.User, error) {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing.Password == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s35Fa", "no password set")
	}
	if err := crypto.CompareHash(existing.Password.SecretCrypto, []byte(password), es.PasswordAlg); err == nil {
		return es.setPasswordCheckResult(ctx, existing, authRequestID, PasswordCheckSucceededAggregate)
	}
	user, err := es.setPasswordCheckResult(ctx, existing, authRequestID, PasswordCheckFailedAggregate)
	if err != nil {
		return nil, err
	}
	return user, caos_errs.ThrowInvalidArgument(nil, "EVENT-452ad", "invalid password")
}

func (es *UserEventstore) setPasswordCheckResult(ctx context.Context, user *usr_model.User, authRequestID string, check func(aggCreator *es_models.AggregateCreator, user *model.User, authRequestID string) es_sdk.AggregateFunc) (*usr_model.User, error) {
	repoUser := model.UserFromModel(user)

	agg := check(es.AggregateCreator(), repoUser, authRequestID)
	err := es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoUser)
	return model.UserToModel(repoUser), nil
}

func (es *UserEventstore) SetOneTimePassword(ctx context.Context, password *usr_model.Password) (*usr_model.Password, error) {
	return es.changedPassword(ctx, password, true)
}

func (es *UserEventstore) SetPassword(ctx context.Context, password *usr_model.Password) (*usr_model.Password, error) {
	return es.changedPassword(ctx, password, false)
}

func (es *UserEventstore) changedPassword(ctx context.Context, password *usr_model.Password, onetime bool) (*usr_model.Password, error) {
	if !password.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dosi3", "password invalid")
	}
	user, err := es.UserByID(ctx, password.AggregateID)
	if err != nil {
		return nil, err
	}

	secret, err := crypto.Hash([]byte(password.SecretString), es.PasswordAlg)
	if err != nil {
		return nil, err
	}
	repoPassword := &model.Password{Secret: secret, ChangeRequired: onetime}
	repoUser := model.UserFromModel(user)
	agg := PasswordChangeAggregate(es.AggregateCreator(), repoUser, repoPassword)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoUser)

	return model.PasswordToModel(repoUser.Password), nil
}

func (es *UserEventstore) RequestSetPassword(ctx context.Context, userID string, notifyType usr_model.NotificationType) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	passwordCode := new(model.PasswordCode)
	err = es.generatePasswordCode(passwordCode, notifyType)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := RequestSetPassword(es.AggregateCreator(), repoUser, passwordCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) ProfileByID(ctx context.Context, userID string) (*usr_model.Profile, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Profile != nil {
		return user.Profile, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-dk23f", "profile not found")
}

func (es *UserEventstore) ChangeProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	if !profile.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d82i3", "profile is invalid")
	}
	existing, err := es.UserByID(ctx, profile.AggregateID)
	if err != nil {
		return nil, err
	}
	repoExisting := model.UserFromModel(existing)
	repoNew := model.ProfileFromModel(profile)

	updateAggregate := ProfileChangeAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.ProfileToModel(repoExisting.Profile), nil
}

func (es *UserEventstore) EmailByID(ctx context.Context, userID string) (*usr_model.Email, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Email != nil {
		return user.Email, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-dki89", "email not found")
}

func (es *UserEventstore) ChangeEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	if !email.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-lco09", "email is invalid")
	}
	existing, err := es.UserByID(ctx, email.AggregateID)
	if err != nil {
		return nil, err
	}

	var emailCode *model.EmailCode
	if !email.IsEmailVerified {
		emailCode = new(model.EmailCode)
		err := es.generateEmailCode(emailCode)
		if err != nil {
			return nil, err
		}
	}

	repoExisting := model.UserFromModel(existing)
	repoNew := model.EmailFromModel(email)
	updateAggregate := EmailChangeAggregate(es.AggregateCreator(), repoExisting, repoNew, emailCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.EmailToModel(repoExisting.Email), nil
}

func (es *UserEventstore) VerifyEmail(ctx context.Context, userID, verificationCode string) error {
	if userID == "" || verificationCode == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9fd", "userId or Code empty")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.EmailCode == nil {
		return caos_errs.ThrowNotFound(nil, "EVENT-lso9w", "code not found")
	}
	if err := crypto.VerifyCode(existing.EmailCode.CreationDate, existing.EmailCode.Expiry, existing.EmailCode.Code, verificationCode, es.EmailVerificationCode); err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	updateAggregate := EmailVerifiedAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lco09", "userID missing")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.Email == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-pdo9s", "no email existing")
	}
	if existing.IsEmailVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-pdo9s", "email already verified")
	}

	emailCode := new(model.EmailCode)
	err = es.generateEmailCode(emailCode)
	if err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	updateAggregate := EmailVerificationCodeAggregate(es.AggregateCreator(), repoExisting, emailCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) PhoneByID(ctx context.Context, userID string) (*usr_model.Phone, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9se", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Phone != nil {
		return user.Phone, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-pos9e", "phone not found")
}

func (es *UserEventstore) ChangePhone(ctx context.Context, phone *usr_model.Phone) (*usr_model.Phone, error) {
	if !phone.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9s4", "phone is invalid")
	}
	existing, err := es.UserByID(ctx, phone.AggregateID)
	if err != nil {
		return nil, err
	}

	var phoneCode *model.PhoneCode
	if !phone.IsPhoneVerified {
		phoneCode = new(model.PhoneCode)
		err := es.generatePhoneCode(phoneCode)
		if err != nil {
			return nil, err
		}
	}

	repoExisting := model.UserFromModel(existing)
	repoNew := model.PhoneFromModel(phone)
	updateAggregate := PhoneChangeAggregate(es.AggregateCreator(), repoExisting, repoNew, phoneCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.PhoneToModel(repoExisting.Phone), nil
}

func (es *UserEventstore) VerifyPhone(ctx context.Context, userID, verificationCode string) error {
	if userID == "" || verificationCode == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dsi8s", "userId or Code empty")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.PhoneCode == nil {
		return caos_errs.ThrowNotFound(nil, "EVENT-slp0s", "code not found")
	}
	if err := crypto.VerifyCode(existing.PhoneCode.CreationDate, existing.PhoneCode.Expiry, existing.PhoneCode.Code, verificationCode, es.PhoneVerificationCode); err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	updateAggregate := PhoneVerifiedAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) CreatePhoneVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9sw", "userID missing")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.Phone == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp9fs", "no phone existing")
	}
	if existing.IsPhoneVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sleis", "phone already verified")
	}

	phoneCode := new(model.PhoneCode)
	err = es.generatePhoneCode(phoneCode)
	if err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	updateAggregate := PhoneVerificationCodeAggregate(es.AggregateCreator(), repoExisting, phoneCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) AddressByID(ctx context.Context, userID string) (*usr_model.Address, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di8ws", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Address != nil {
		return user.Address, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-so9wa", "address not found")
}

func (es *UserEventstore) ChangeAddress(ctx context.Context, address *usr_model.Address) (*usr_model.Address, error) {
	existing, err := es.UserByID(ctx, address.AggregateID)
	if err != nil {
		return nil, err
	}
	repoExisting := model.UserFromModel(existing)
	repoNew := model.AddressFromModel(address)

	updateAggregate := AddressChangeAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.AddressToModel(repoExisting.Address), nil
}

func (es *UserEventstore) OTPByID(ctx context.Context, userID string) (*usr_model.OTP, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9se", "userID missing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.OTP != nil {
		return user.OTP, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-dps09", "otp not found")
}

func (es *UserEventstore) AddOTP(ctx context.Context, userID string) (*usr_model.OTP, error) {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing.OTP != nil && existing.OTP.State == usr_model.MFASTATE_READY {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-do9se", "user has already configured otp")
	}
	key, err := totp.Generate(totp.GenerateOpts{Issuer: es.Multifactors.OTP.Issuer, AccountName: userID})
	if err != nil {
		return nil, err
	}
	encrypted, err := crypto.Encrypt([]byte(key.Secret()), es.Multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	repoOtp := &model.OTP{Secret: encrypted}
	repoExisting := model.UserFromModel(existing)
	updateAggregate := MfaOTPAddAggregate(es.AggregateCreator(), repoExisting, repoOtp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	otp := model.OTPToModel(repoExisting.OTP)
	otp.Url = key.URL()
	otp.SecretString = key.Secret()
	return otp, nil
}

func (es *UserEventstore) RemoveOTP(ctx context.Context, userID string) error {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.OTP == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp0de", "no otp existing")
	}
	repoExisting := model.UserFromModel(existing)
	updateAggregate := MfaOTPRemoveAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) CheckMfaOTP(ctx context.Context, userID, code string) error {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.OTP == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp0de", "no otp existing")
	}
	decrypt, err := crypto.DecryptString(existing.OTP.Secret, es.Multifactors.OTP.CryptoMFA)
	if err != nil {
		return err
	}

	valid := totp.Validate(code, decrypt)
	if !valid {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-8isk2", "Invalid code")
	}
	return nil
}
