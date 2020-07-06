package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"strings"
)

func UserByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	return UserQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAggregate).
		LatestSequenceFilter(latestSequence)
}

func UserUserNameUniqueQuery(userName string) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserUserNameAggregate).
		AggregateIDFilter(userName).
		OrderDesc().
		SetLimit(1)
}

func UserEmailUniqueQuery(email string) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserEmailAggregate).
		AggregateIDFilter(email).
		OrderDesc().
		SetLimit(1)
}

func UserAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User) (*es_models.Aggregate, error) {
	if user == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, user.AggregateID, model.UserAggregate, model.UserVersion, user.Sequence)
}

func UserAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if user == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "Errors.Internal")
	}

	return aggCreator.NewAggregate(ctx, user.AggregateID, model.UserAggregate, model.UserVersion, user.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func UserCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, initCode *model.InitUserCode, phoneCode *model.PhoneCode, resourceOwner string, userLoginMustBeDomain bool) (_ []*es_models.Aggregate, err error) {
	if user == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "Errors.Internal")
	}

	var agg *es_models.Aggregate
	if resourceOwner != "" {
		agg, err = UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
	} else {
		agg, err = UserAggregate(ctx, aggCreator, user)
	}
	if err != nil {
		return nil, err
	}
	if !userLoginMustBeDomain {
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate).
			AggregateIDsFilter()

		validation := addUserNameValidation(user.UserName)
		agg.SetPrecondition(validationQuery, validation)
	}

	agg, err = agg.AppendEvent(model.UserAdded, user)
	if err != nil {
		return nil, err
	}
	if user.Email != nil && user.EmailAddress != "" && user.IsEmailVerified {
		agg, err = agg.AppendEvent(model.UserEmailVerified, nil)
		if err != nil {
			return nil, err
		}
	}
	if user.Phone != nil && user.PhoneNumber != "" && user.IsPhoneVerified {
		agg, err = agg.AppendEvent(model.UserPhoneVerified, nil)
		if err != nil {
			return nil, err
		}
	}
	if initCode != nil {
		agg, err = agg.AppendEvent(model.InitializedUserCodeAdded, initCode)
		if err != nil {
			return nil, err
		}
	}
	if phoneCode != nil {
		agg, err = agg.AppendEvent(model.UserPhoneCodeAdded, phoneCode)
		if err != nil {
			return nil, err
		}
	}
	uniqueAggregates, err := getUniqueUserAggregates(ctx, aggCreator, user, resourceOwner, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregates[0],
		uniqueAggregates[1],
	}, nil
}

func UserRegisterAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwner string, emailCode *model.EmailCode, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	if user == nil || resourceOwner == "" || emailCode == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "user, resourceowner, emailcode should not be nothing")
	}

	agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
	if err != nil {
		return nil, err
	}

	if !userLoginMustBeDomain {
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate).
			AggregateIDsFilter()

		validation := addUserNameValidation(user.UserName)
		agg.SetPrecondition(validationQuery, validation)
	}
	agg, err = agg.AppendEvent(model.UserRegistered, user)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.UserEmailCodeAdded, emailCode)
	if err != nil {
		return nil, err
	}
	uniqueAggregates, err := getUniqueUserAggregates(ctx, aggCreator, user, resourceOwner, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregates[0],
		uniqueAggregates[1],
	}, nil
}

func getUniqueUserAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwner string, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	userNameAggregate, err := reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, user.UserName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}

	emailAggregate, err := reservedUniqueEmailAggregate(ctx, aggCreator, resourceOwner, user.EmailAddress)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		userNameAggregate,
		emailAggregate,
	}, nil
}
func reservedUniqueUserNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, userName string, userLoginMustBeDomain bool) (*es_models.Aggregate, error) {
	uniqueUserName := userName
	if userLoginMustBeDomain {
		uniqueUserName = userName + resourceOwner
	}
	aggregate, err := aggCreator.NewAggregate(ctx, uniqueUserName, model.UserUserNameAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, uniqueUserName, model.UserUserNameAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserUserNameReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserUserNameUniqueQuery(uniqueUserName), isEventValidation(aggregate, model.UserUserNameReserved)), nil
}

func reservedUniqueEmailAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, email string) (aggregate *es_models.Aggregate, err error) {
	aggregate, err = aggCreator.NewAggregate(ctx, email, model.UserEmailAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, email, model.UserEmailAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserEmailReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserEmailUniqueQuery(email), isEventValidation(aggregate, model.UserEmailReserved)), nil
}

func releasedUniqueEmailAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, email string) (aggregate *es_models.Aggregate, err error) {
	aggregate, err = aggCreator.NewAggregate(ctx, email, model.UserEmailAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, email, model.UserEmailAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserEmailReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserEmailUniqueQuery(email), isEventValidation(aggregate, model.UserEmailReleased)), nil
}

func UserDeactivateAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, model.UserDeactivated)
}

func UserReactivateAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, model.UserReactivated)
}

func UserLockAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, model.UserLocked)
}

func UserUnlockAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, model.UserUnlocked)
}

func userStateAggregate(aggCreator *es_models.AggregateCreator, user *model.User, state es_models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(state, nil)
	}
}

func UserInitCodeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, code *model.InitUserCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8i23", "code should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedUserCodeAdded, code)
	}
}

func UserInitCodeSentAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedUserCodeSent, nil)
	}
}

func InitCodeVerifiedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, password *model.Password) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Email != nil && !existing.Email.IsEmailVerified {
			agg, err = agg.AppendEvent(model.UserEmailVerified, nil)
			if err != nil {
				return nil, err
			}
		}
		if password != nil {
			agg, err = agg.AppendEvent(model.UserPasswordChanged, password)
			if err != nil {
				return nil, err
			}
		}
		return agg.AppendEvent(model.InitializedUserCheckSucceeded, nil)
	}
}

func InitCodeCheckFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedUserCheckFailed, nil)
	}
}

func SkipMfaAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaInitSkipped, nil)
	}
}

func PasswordChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, password *model.Password) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if password == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9832", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPasswordChanged, password)
	}
}

func PasswordCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, check *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPasswordCheckSucceeded, check)
	}
}
func PasswordCheckFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, check *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPasswordCheckFailed, check)
	}
}

func RequestSetPassword(aggCreator *es_models.AggregateCreator, existing *model.User, request *model.PasswordCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if request == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8ei2", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPasswordCodeAdded, request)
	}
}

func PasswordCodeSentAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPasswordCodeSent, nil)
	}
}

func ProfileChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, profile *model.Profile) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if profile == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Profile.Changes(profile)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0spow", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.UserProfileChanged, changes)
	}
}

func EmailChangeAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.User, email *model.Email, code *model.EmailCode) ([]*es_models.Aggregate, error) {
	if email == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "Errors.Internal")
	}
	if (!email.IsEmailVerified && code == nil) || (email.IsEmailVerified && code != nil) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-id934", "Errors.Internal")
	}
	changes := existing.Email.Changes(email)
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-s90pw", "Errors.NoChangesFound")
	}
	aggregates := make([]*es_models.Aggregate, 0, 4)
	reserveEmailAggregate, err := reservedUniqueEmailAggregate(ctx, aggCreator, "", email.EmailAddress)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, reserveEmailAggregate)
	releaseEmailAggregate, err := releasedUniqueEmailAggregate(ctx, aggCreator, "", existing.EmailAddress)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, releaseEmailAggregate)
	agg, err := UserAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.UserEmailChanged, changes)
	if err != nil {
		return nil, err
	}
	if existing.Email == nil {
		existing.Email = new(model.Email)
	}
	if email.IsEmailVerified {
		agg, err = agg.AppendEvent(model.UserEmailVerified, code)
		if err != nil {
			return nil, err
		}
	}
	if code != nil {
		agg, err = agg.AppendEvent(model.UserEmailCodeAdded, code)
		if err != nil {
			return nil, err
		}
	}
	return append(aggregates, agg), nil
}

func EmailVerifiedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserEmailVerified, nil)
	}
}

func EmailVerificationFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserEmailVerificationFailed, nil)
	}
}

func EmailVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, code *model.EmailCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserEmailCodeAdded, code)
	}
}

func EmailCodeSentAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserEmailCodeSent, nil)
	}
}

func PhoneChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, phone *model.Phone, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if phone == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkso3", "Errors.Internal")
		}
		if (!phone.IsPhoneVerified && code == nil) || (phone.IsPhoneVerified && code != nil) {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dksi8", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Phone == nil {
			existing.Phone = new(model.Phone)
		}
		changes := existing.Phone.Changes(phone)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sp0oc", "Errors.NoChangesFound")
		}
		agg, err = agg.AppendEvent(model.UserPhoneChanged, changes)
		if err != nil {
			return nil, err
		}
		if phone.IsPhoneVerified {
			return agg.AppendEvent(model.UserPhoneVerified, code)
		}
		if code != nil {
			return agg.AppendEvent(model.UserPhoneCodeAdded, code)
		}
		return agg, nil
	}
}

func PhoneRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPhoneRemoved, nil)
	}
}

func PhoneVerifiedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPhoneVerified, nil)
	}
}

func PhoneVerificationFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPhoneVerificationFailed, nil)
	}
}

func PhoneVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dsue2", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPhoneCodeAdded, code)
	}
}

func PhoneCodeSentAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserPhoneCodeSent, nil)
	}
}

func AddressChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, address *model.Address) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if address == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Address == nil {
			existing.Address = new(model.Address)
		}
		changes := existing.Address.Changes(address)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2tszw", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.UserAddressChanged, changes)
	}
}

func MfaOTPAddAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, otp *model.OTP) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if otp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaOtpAdded, otp)
	}
}

func MfaOTPVerifyAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaOtpVerified, nil)
	}
}

func MfaOTPCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sd5DA", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaOtpCheckSucceeded, authReq)
	}
}

func MfaOTPCheckFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-64sd6", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaOtpCheckFailed, authReq)
	}
}

func MfaOTPRemoveAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MfaOtpRemoved, nil)
	}
}

func SignOutAggregates(aggCreator *es_models.AggregateCreator, existingUsers []*model.User, agentID string) func(ctx context.Context) ([]*es_models.Aggregate, error) {
	return func(ctx context.Context) ([]*es_models.Aggregate, error) {
		aggregates := make([]*es_models.Aggregate, len(existingUsers))
		for i, existing := range existingUsers {
			agg, err := UserAggregateOverwriteContext(ctx, aggCreator, existing, existing.ResourceOwner, existing.AggregateID)
			if err != nil {
				return nil, err
			}
			agg.AppendEvent(model.SignedOut, map[string]interface{}{"userAgentID": agentID})
			aggregates[i] = agg
		}
		return aggregates, nil
	}
}

func isEventValidation(aggregate *es_models.Aggregate, eventType es_models.EventType) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == eventType {
			return errors.ThrowPreconditionFailedf(nil, "EVENT-eJQqe", "user is already %v", eventType)
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}
}

func addUserNameValidation(userName string) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		domains := make([]*org_es_model.OrgDomain, 0)
		for _, event := range events {
			switch event.Type {
			case org_es_model.OrgDomainAdded:
				domain := new(org_es_model.OrgDomain)
				domain.SetData(event)
			case org_es_model.OrgDomainVerified:
				domain := new(org_es_model.OrgDomain)
				domain.SetData(event)
				for _, d := range domains {
					if d.Domain == domain.Domain {
						d.Verified = true
					}
				}
			case org_es_model.OrgDomainRemoved:
				domain := new(org_es_model.OrgDomain)
				domain.SetData(event)
				for i, d := range domains {
					if d.Domain == domain.Domain {
						domains[i] = domains[len(domains)-1]
						domains[len(domains)-1] = nil
						domains = domains[:len(domains)-1]
					}
				}
			}
		}
		split := strings.Split(userName, "@")
		if len(split) != 2 {
			return nil
		}
		for _, d := range domains {
			if d.Verified && d.Domain == split[1] {
				return errors.ThrowPreconditionFailed(nil, "EVENT-us5Zw", "Errors.User.DomainNotAllowedAsUsername")
			}
		}
		return nil
	}
}
