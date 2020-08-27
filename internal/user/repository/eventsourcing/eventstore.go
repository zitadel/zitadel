package eventsourcing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/id"
	org_model "github.com/caos/zitadel/internal/org/model"
	policy_model "github.com/caos/zitadel/internal/policy/model"

	"github.com/pquerna/otp/totp"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
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
)

type UserEventstore struct {
	es_int.Eventstore
	userCache                *UserCache
	idGenerator              id.Generator
	domain                   string
	PasswordAlg              crypto.HashAlgorithm
	InitializeUserCode       crypto.Generator
	EmailVerificationCode    crypto.Generator
	PhoneVerificationCode    crypto.Generator
	PasswordVerificationCode crypto.Generator
	Multifactors             global_model.Multifactors
	validateTOTP             func(string, string) bool
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
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	initCodeGen := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.InitializeUserCode, aesCrypto)
	emailVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.EmailVerificationCode, aesCrypto)
	phoneVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.PhoneVerificationCode, aesCrypto)
	passwordVerificationCode := crypto.NewEncryptionGenerator(systemDefaults.SecretGenerators.PasswordVerificationCode, aesCrypto)
	aesOtpCrypto, err := crypto.NewAESCrypto(systemDefaults.Multifactors.OTP.VerificationKey)
	passwordAlg := crypto.NewBCrypt(systemDefaults.SecretGenerators.PasswordSaltCost)

	return &UserEventstore{
		Eventstore:               conf.Eventstore,
		userCache:                userCache,
		idGenerator:              id.SonyFlakeGenerator,
		domain:                   systemDefaults.Domain,
		InitializeUserCode:       initCodeGen,
		EmailVerificationCode:    emailVerificationCode,
		PhoneVerificationCode:    phoneVerificationCode,
		PasswordVerificationCode: passwordVerificationCode,
		Multifactors: global_model.Multifactors{
			OTP: global_model.OTP{
				CryptoMFA: aesOtpCrypto,
				Issuer:    systemDefaults.Multifactors.OTP.Issuer,
			},
		},
		PasswordAlg:  passwordAlg,
		validateTOTP: totp.Validate,
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

func (es *UserEventstore) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := UserByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

func (es *UserEventstore) PrepareCreateUser(ctx context.Context, user *usr_model.User, pwPolicy *policy_model.PasswordComplexityPolicy, orgIAMPolicy *org_model.OrgIAMPolicy, resourceOwner string) (*model.User, []*es_models.Aggregate, error) {
	err := user.CheckOrgIAMPolicy(orgIAMPolicy)
	if err != nil {
		return nil, nil, err
	}
	user.SetNamesAsDisplayname()
	if !user.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.User.Invalid")
	}

	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	user.AggregateID = id

	err = user.HashPasswordIfExisting(pwPolicy, es.PasswordAlg, true)
	if err != nil {
		return nil, nil, err
	}
	err = user.GenerateInitCodeIfNeeded(es.InitializeUserCode)
	if err != nil {
		return nil, nil, err
	}
	err = user.GeneratePhoneCodeIfNeeded(es.PhoneVerificationCode)
	if err != nil {
		return nil, nil, err
	}

	repoUser := model.UserFromModel(user)
	repoInitCode := model.InitCodeFromModel(user.InitCode)
	repoPhoneCode := model.PhoneCodeFromModel(user.PhoneCode)

	createAggregates, err := UserCreateAggregate(ctx, es.AggregateCreator(), repoUser, repoInitCode, repoPhoneCode, resourceOwner, orgIAMPolicy.UserLoginMustBeDomain)

	return repoUser, createAggregates, err
}

func (es *UserEventstore) CreateUser(ctx context.Context, user *usr_model.User, pwPolicy *policy_model.PasswordComplexityPolicy, orgIAMPolicy *org_model.OrgIAMPolicy) (*usr_model.User, error) {
	repoUser, aggregates, err := es.PrepareCreateUser(ctx, user, pwPolicy, orgIAMPolicy, "")
	if err != nil {
		return nil, err
	}

	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoUser.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoUser)
	return model.UserToModel(repoUser), nil
}

func (es *UserEventstore) PrepareRegisterUser(ctx context.Context, user *usr_model.User, policy *policy_model.PasswordComplexityPolicy, orgIAMPolicy *org_model.OrgIAMPolicy, resourceOwner string) (*model.User, []*es_models.Aggregate, error) {
	err := user.CheckOrgIAMPolicy(orgIAMPolicy)
	if err != nil {
		return nil, nil, err
	}
	user.SetNamesAsDisplayname()
	if !user.IsValid() || user.Password == nil || user.SecretString == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.User.Invalid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	user.AggregateID = id

	err = user.HashPasswordIfExisting(policy, es.PasswordAlg, false)
	if err != nil {
		return nil, nil, err
	}
	err = user.GenerateInitCodeIfNeeded(es.InitializeUserCode)
	if err != nil {
		return nil, nil, err
	}

	repoUser := model.UserFromModel(user)
	repoInitCode := model.InitCodeFromModel(user.InitCode)

	aggregates, err := UserRegisterAggregate(ctx, es.AggregateCreator(), repoUser, resourceOwner, repoInitCode, orgIAMPolicy.UserLoginMustBeDomain)
	return repoUser, aggregates, err
}

func (es *UserEventstore) RegisterUser(ctx context.Context, user *usr_model.User, pwPolicy *policy_model.PasswordComplexityPolicy, orgIAMPolicy *org_model.OrgIAMPolicy, resourceOwner string) (*usr_model.User, error) {
	repoUser, createAggregates, err := es.PrepareRegisterUser(ctx, user, pwPolicy, orgIAMPolicy, resourceOwner)
	if err != nil {
		return nil, err
	}

	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoUser.AppendEvents, createAggregates...)
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "Errors.User.AlreadyInactive")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do94s", "Errors.User.NotInactive")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83s", "Errors.User.ShouldBeActiveOrInitial")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dks83", "Errors.User.NotLocked")
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

func (es *UserEventstore) UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*usr_model.UserChanges, error) {
	query := ChangesQuery(id, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-g9HCv").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-htuG9", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-6cAxe", "Errors.User.NoChanges")
	}

	result := make([]*usr_model.UserChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-8GTGS").OnError(err).Debug("unable to parse timestamp")
		change := &usr_model.UserChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			user := new(model.Profile)
			err := json.Unmarshal(event.Data, user)
			logging.Log("EVENT-Rkg7X").OnError(err).Debug("unable to unmarshal data")
			change.Data = user
		}

		result[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &usr_model.UserChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}, nil
}

func ChangesQuery(userID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAggregate)
	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(userID).
		SetLimit(limit)
	return query
}

func (es *UserEventstore) InitializeUserCodeByID(ctx context.Context, userID string) (*usr_model.InitUserCode, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d8diw", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.InitCode != nil {
		return user.InitCode, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-d8e2", "Erorrs.User.InitCodeNotFound")
}

func (es *UserEventstore) CreateInitializeUserCodeByID(ctx context.Context, userID string) (*usr_model.InitUserCode, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	initCode := new(usr_model.InitUserCode)
	err = initCode.GenerateInitUserCode(es.InitializeUserCode)
	if err != nil {
		return nil, err
	}

	repoUser := model.UserFromModel(user)
	repoInitCode := model.InitCodeFromModel(initCode)

	agg := UserInitCodeAggregate(es.AggregateCreator(), repoUser, repoInitCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.userCache.cacheUser(repoUser)
	return model.InitCodeToModel(repoUser.InitCode), nil
}

func (es *UserEventstore) InitCodeSent(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-0posw", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := UserInitCodeSentAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) VerifyInitCode(ctx context.Context, policy *policy_model.PasswordComplexityPolicy, userID, verificationCode, password string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9fd", "Errors.User.UserIDMissing")
	}
	if verificationCode == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9fd", "Errors.User.Code.Empty")
	}
	pw := &usr_model.Password{SecretString: password}
	err := pw.HashPasswordIfExisting(policy, es.PasswordAlg, false)
	if err != nil {
		return err
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.InitCode == nil {
		return caos_errs.ThrowNotFound(nil, "EVENT-spo9W", "Errors.User.Code.NotFound")
	}
	repoPassword := model.PasswordFromModel(pw)
	repoExisting := model.UserFromModel(existing)
	var updateAggregate func(ctx context.Context) (*es_models.Aggregate, error)
	if err := crypto.VerifyCode(existing.InitCode.CreationDate, existing.InitCode.Expiry, existing.InitCode.Code, verificationCode, es.InitializeUserCode); err != nil {
		updateAggregate = InitCodeCheckFailedAggregate(es.AggregateCreator(), repoExisting)
		es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
		return err
	} else {
		updateAggregate = InitCodeVerifiedAggregate(es.AggregateCreator(), repoExisting, repoPassword)
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) SkipMfaInit(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "Errors.User.UserIDMissing")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Password != nil {
		return user.Password, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-d8e2", "Errors.User.Password.NotFound")
}

func (es *UserEventstore) CheckPassword(ctx context.Context, userID, password string, authRequest *req_model.AuthRequest) error {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.Password == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-s35Fa", "Errors.User.Password.Empty")
	}
	if err := crypto.CompareHash(existing.Password.SecretCrypto, []byte(password), es.PasswordAlg); err == nil {
		return es.setPasswordCheckResult(ctx, existing, authRequest, PasswordCheckSucceededAggregate)
	}
	if err := es.setPasswordCheckResult(ctx, existing, authRequest, PasswordCheckFailedAggregate); err != nil {
		return err
	}
	return caos_errs.ThrowInvalidArgument(nil, "EVENT-452ad", "Errors.User.Password.Invalid")
}

func (es *UserEventstore) setPasswordCheckResult(ctx context.Context, user *usr_model.User, authRequest *req_model.AuthRequest, check func(*es_models.AggregateCreator, *model.User, *model.AuthRequest) es_sdk.AggregateFunc) error {
	repoUser := model.UserFromModel(user)
	repoAuthRequest := model.AuthRequestFromModel(authRequest)
	agg := check(es.AggregateCreator(), repoUser, repoAuthRequest)
	err := es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) SetOneTimePassword(ctx context.Context, policy *policy_model.PasswordComplexityPolicy, password *usr_model.Password) (*usr_model.Password, error) {
	user, err := es.UserByID(ctx, password.AggregateID)
	if err != nil {
		return nil, err
	}
	return es.changedPassword(ctx, user, policy, password.SecretString, true)
}

func (es *UserEventstore) SetPassword(ctx context.Context, policy *policy_model.PasswordComplexityPolicy, userID, code, password string) error {
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.PasswordCode == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-65sdr", "Errors.User.Code.NotFound")
	}
	if err := crypto.VerifyCode(user.PasswordCode.CreationDate, user.PasswordCode.Expiry, user.PasswordCode.Code, code, es.PasswordVerificationCode); err != nil {
		return err
	}
	_, err = es.changedPassword(ctx, user, policy, password, false)
	return err
}

func (es *UserEventstore) ChangePassword(ctx context.Context, policy *policy_model.PasswordComplexityPolicy, userID, old, new string) (*usr_model.Password, error) {
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Password == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Fds3s", "Errors.User.Password.Empty")
	}
	if err := crypto.CompareHash(user.Password.SecretCrypto, []byte(old), es.PasswordAlg); err != nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-s56a3", "Errors.User.Password.Invalid")
	}
	return es.changedPassword(ctx, user, policy, new, false)
}

func (es *UserEventstore) changedPassword(ctx context.Context, user *usr_model.User, policy *policy_model.PasswordComplexityPolicy, password string, onetime bool) (*usr_model.Password, error) {
	pw := &usr_model.Password{SecretString: password}
	err := pw.HashPasswordIfExisting(policy, es.PasswordAlg, onetime)
	if err != nil {
		return nil, err
	}
	repoPassword := model.PasswordFromModel(pw)
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
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dic8s", "Errors.User.UserIDMissing")
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

func (es *UserEventstore) PasswordCodeSent(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-s09ow", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := PasswordCodeSentAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) ProfileByID(ctx context.Context, userID string) (*usr_model.Profile, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Profile != nil {
		return user.Profile, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-dk23f", "Errors.User.ProfileNotFound")
}

func (es *UserEventstore) ChangeProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	profile.SetNamesAsDisplayname()
	if !profile.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d82i3", "Errors.User.ProfileInvalid")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di834", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Email != nil {
		return user.Email, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-dki89", "Errors.User.EmailNotFound")
}

func (es *UserEventstore) ChangeEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	if !email.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-lco09", "Errors.User.EmailInvalid")
	}
	existing, err := es.UserByID(ctx, email.AggregateID)
	if err != nil {
		return nil, err
	}

	emailCode, err := email.GenerateEmailCodeIfNeeded(es.EmailVerificationCode)
	if err != nil {
		return nil, err
	}

	repoExisting := model.UserFromModel(existing)
	repoNew := model.EmailFromModel(email)
	repoEmailCode := model.EmailCodeFromModel(emailCode)

	updateAggregate, err := EmailChangeAggregate(ctx, es.AggregateCreator(), repoExisting, repoNew, repoEmailCode)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.EmailToModel(repoExisting.Email), nil
}

func (es *UserEventstore) VerifyEmail(ctx context.Context, userID, verificationCode string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9fd", "Errors.User.UserIDMissing")
	}
	if verificationCode == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-skDws", "Errors.User.Code.Empty")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.EmailCode == nil {
		return caos_errs.ThrowNotFound(nil, "EVENT-lso9w", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCode(existing.EmailCode.CreationDate, existing.EmailCode.Expiry, existing.EmailCode.Code, verificationCode, es.EmailVerificationCode)
	if err == nil {
		return es.setEmailVerifyResult(ctx, existing, EmailVerifiedAggregate)
	}
	if err := es.setEmailVerifyResult(ctx, existing, EmailVerificationFailedAggregate); err != nil {
		return err
	}
	return caos_errs.ThrowInvalidArgument(err, "EVENT-dtGaa", "Errors.User.Code.Invalid")
}

func (es *UserEventstore) setEmailVerifyResult(ctx context.Context, existing *usr_model.User, check func(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc) error {
	repoExisting := model.UserFromModel(existing)
	err := es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, check(es.AggregateCreator(), repoExisting))
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-lco09", "Errors.User.UserIDMissing")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.Email == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-pdo9s", "Errors.User.EmailNotFound")
	}
	if existing.IsEmailVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-pdo9s", "Errors.User.EmailAlreadyVerified")
	}

	emailCode := new(usr_model.EmailCode)
	err = emailCode.GenerateEmailCode(es.EmailVerificationCode)
	if err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	repoEmailCode := model.EmailCodeFromModel(emailCode)
	updateAggregate := EmailVerificationCodeAggregate(es.AggregateCreator(), repoExisting, repoEmailCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) EmailVerificationCodeSent(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-spo0w", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := EmailCodeSentAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) PhoneByID(ctx context.Context, userID string) (*usr_model.Phone, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9se", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Phone != nil {
		return user.Phone, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-pos9e", "Errors.User.PhoneNotFound")
}

func (es *UserEventstore) ChangePhone(ctx context.Context, phone *usr_model.Phone) (*usr_model.Phone, error) {
	if !phone.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9s4", "Errors.User.PhoneInvalid")
	}
	existing, err := es.UserByID(ctx, phone.AggregateID)
	if err != nil {
		return nil, err
	}

	phoneCode, err := phone.GeneratePhoneCodeIfNeeded(es.PhoneVerificationCode)
	if err != nil {
		return nil, err
	}

	repoExisting := model.UserFromModel(existing)
	repoNew := model.PhoneFromModel(phone)
	repoPhoneCode := model.PhoneCodeFromModel(phoneCode)

	updateAggregate := PhoneChangeAggregate(es.AggregateCreator(), repoExisting, repoNew, repoPhoneCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.userCache.cacheUser(repoExisting)
	return model.PhoneToModel(repoExisting.Phone), nil
}

func (es *UserEventstore) VerifyPhone(ctx context.Context, userID, verificationCode string) error {
	if userID == "" || verificationCode == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dsi8s", "Errors.User.UserIDMissing")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.PhoneCode == nil {
		return caos_errs.ThrowNotFound(nil, "EVENT-slp0s", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCode(existing.PhoneCode.CreationDate, existing.PhoneCode.Expiry, existing.PhoneCode.Code, verificationCode, es.PhoneVerificationCode)
	if err == nil {
		return es.setPhoneVerifyResult(ctx, existing, PhoneVerifiedAggregate)
	}
	if err := es.setPhoneVerifyResult(ctx, existing, PhoneVerificationFailedAggregate); err != nil {
		return err
	}
	return caos_errs.ThrowInvalidArgument(err, "EVENT-dsf4G", "Errors.User.Code.Invalid")
}

func (es *UserEventstore) setPhoneVerifyResult(ctx context.Context, existing *usr_model.User, check func(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc) error {
	repoExisting := model.UserFromModel(existing)
	err := es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, check(es.AggregateCreator(), repoExisting))
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) CreatePhoneVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-do9sw", "Errors.User.UserIDMissing")
	}
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if existing.Phone == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp9fs", "Errors.User.PhoneNotFound")
	}
	if existing.IsPhoneVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sleis", "Errors.User.PhoneAlreadyVerified")
	}

	phoneCode := new(usr_model.PhoneCode)
	err = phoneCode.GeneratePhoneCode(es.PhoneVerificationCode)
	if err != nil {
		return err
	}

	repoExisting := model.UserFromModel(existing)
	repoPhoneCode := model.PhoneCodeFromModel(phoneCode)
	updateAggregate := PhoneVerificationCodeAggregate(es.AggregateCreator(), repoExisting, repoPhoneCode)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) PhoneVerificationCodeSent(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp0wa", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := PhoneCodeSentAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) RemovePhone(ctx context.Context, userID string) error {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	repoExisting := model.UserFromModel(existing)
	removeAggregate := PhoneRemovedAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, removeAggregate)
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoExisting)
	return nil
}

func (es *UserEventstore) AddressByID(ctx context.Context, userID string) (*usr_model.Address, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di8ws", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Address != nil {
		return user.Address, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-so9wa", "Errors.User.AddressNotFound")
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

func (es *UserEventstore) AddOTP(ctx context.Context, userID, accountName string) (*usr_model.OTP, error) {
	existing, err := es.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing.IsOTPReady() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-do9se", "Errors.User.Mfa.Otp.AlreadyReady")
	}
	if accountName == "" {
		accountName = existing.UserName
		if existing.Email != nil {
			accountName = existing.EmailAddress
		}
	}
	key, err := totp.Generate(totp.GenerateOpts{Issuer: es.Multifactors.OTP.Issuer, AccountName: accountName})
	if err != nil {
		return nil, err
	}
	encryptedSecret, err := crypto.Encrypt([]byte(key.Secret()), es.Multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	repoOtp := &model.OTP{Secret: encryptedSecret}
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
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sp0de", "Errors.User.Mfa.Otp.NotExisting")
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

func (es *UserEventstore) CheckMfaOTPSetup(ctx context.Context, userID, code string) error {
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.OTP == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sd5NJ", "Errors.Users.Mfa.Otp.NotExisting")
	}
	if user.IsOTPReady() {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sd5NJ", "Errors.Users.Mfa.Otp.AlreadyReady")
	}
	if err := es.verifyMfaOTP(user.OTP, code); err != nil {
		return err
	}
	repoUser := model.UserFromModel(user)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, MfaOTPVerifyAggregate(es.AggregateCreator(), repoUser))
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) CheckMfaOTP(ctx context.Context, userID, code string, authRequest *req_model.AuthRequest) error {
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsOTPReady() {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-sd5NJ", "Errors.User.Mfa.Otp.NotReady")
	}

	repoUser := model.UserFromModel(user)
	repoAuthReq := model.AuthRequestFromModel(authRequest)
	var aggregate func(*es_models.AggregateCreator, *model.User, *model.AuthRequest) es_sdk.AggregateFunc
	var checkErr error
	if checkErr = es.verifyMfaOTP(user.OTP, code); checkErr != nil {
		aggregate = MfaOTPCheckFailedAggregate
	} else {
		aggregate = MfaOTPCheckSucceededAggregate
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, aggregate(es.AggregateCreator(), repoUser, repoAuthReq))
	if checkErr != nil {
		return checkErr
	}
	if err != nil {
		return err
	}

	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) verifyMfaOTP(otp *usr_model.OTP, code string) error {
	decrypt, err := crypto.DecryptString(otp.Secret, es.Multifactors.OTP.CryptoMFA)
	if err != nil {
		return err
	}

	valid := es.validateTOTP(code, decrypt)
	if !valid {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.Mfa.Otp.InvalidCode")
	}
	return nil
}

func (es *UserEventstore) SignOut(ctx context.Context, agentID string, userIDs []string) error {
	users := make([]*model.User, len(userIDs))
	for i, id := range userIDs {
		user, err := es.UserByID(ctx, id)
		if err != nil {
			return err
		}
		users[i] = model.UserFromModel(user)
	}

	aggFunc := SignOutAggregates(es.AggregateCreator(), users, agentID)
	aggregates, err := aggFunc(ctx)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, nil, aggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (es *UserEventstore) PrepareDomainClaimed(ctx context.Context, userIDs []string) ([]*es_models.Aggregate, error) {
	aggregates := make([]*es_models.Aggregate, 0)
	for _, userID := range userIDs {
		user, err := es.UserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		repoUser := model.UserFromModel(user)
		name, err := es.generateTemporaryLoginName()
		if err != nil {
			return nil, err
		}
		userAgg, err := DomainClaimedAggregate(ctx, es.AggregateCreator(), repoUser, name)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, userAgg...)
	}
	return aggregates, nil
}

func (es *UserEventstore) DomainClaimedSent(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-0posw", "Errors.User.UserIDMissing")
	}
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	repoUser := model.UserFromModel(user)
	agg := DomainClaimedSentAggregate(es.AggregateCreator(), repoUser)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUser.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) ChangeUsername(ctx context.Context, userID, username string, orgIamPolicy *org_model.OrgIAMPolicy) error {
	user, err := es.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	oldUsername := user.UserName
	user.UserName = username
	if err := user.CheckOrgIAMPolicy(orgIamPolicy); err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D23s2", "Errors.Users.Mfa.Otp.NotExisting")
	}
	repoUser := model.UserFromModel(user)
	aggregates, err := UsernameChangedAggregates(ctx, es.AggregateCreator(), repoUser, oldUsername, orgIamPolicy.UserLoginMustBeDomain)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoUser.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	es.userCache.cacheUser(repoUser)
	return nil
}

func (es *UserEventstore) generateTemporaryLoginName() (string, error) {
	id, err := es.idGenerator.Next()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@temporary.%s", id, es.domain), nil
}
