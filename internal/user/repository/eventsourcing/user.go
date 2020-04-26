package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func UserByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "id should be filled")
	}
	return UserQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(usr_model.UserAggregate).
		LatestSequenceFilter(latestSequence)
}

func UserAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User) (*es_models.Aggregate, error) {
	if user == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing user should not be nil")
	}
	return aggCreator.NewAggregate(ctx, user.AggregateID, usr_model.UserAggregate, model.UserVersion, user.Sequence)
}

func UserAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if user == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing user should not be nil")
	}

	return aggCreator.NewAggregate(ctx, user.AggregateID, usr_model.UserAggregate, model.UserVersion, user.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func UserCreateAggregate(aggCreator *es_models.AggregateCreator, user *model.User, initCode *model.InitUserCode, phoneCode *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if user == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "user should not be nil")
		}

		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}

		agg, err = agg.AppendEvent(usr_model.UserAdded, user)
		if err != nil {
			return nil, err
		}
		if user.Email != nil && user.EmailAddress != "" && user.IsEmailVerified {
			agg, err = agg.AppendEvent(usr_model.UserEmailVerified, nil)
			if err != nil {
				return nil, err
			}
		}
		if user.Phone != nil && user.PhoneNumber != "" && user.IsPhoneVerified {
			agg, err = agg.AppendEvent(usr_model.UserPhoneVerified, nil)
			if err != nil {
				return nil, err
			}
		}
		if user.Password != nil {
			agg, err = agg.AppendEvent(usr_model.UserPasswordSetRequested, user.Password)
			if err != nil {
				return nil, err
			}
		}
		if initCode != nil {
			agg, err = agg.AppendEvent(usr_model.InitializedUserCodeCreated, initCode)
			if err != nil {
				return nil, err
			}
		}
		if phoneCode != nil {
			agg, err = agg.AppendEvent(usr_model.UserPhoneCodeAdded, phoneCode)
			if err != nil {
				return nil, err
			}
		}
		return agg, err
	}
}

func UserRegisterAggregate(aggCreator *es_models.AggregateCreator, user *model.User, resourceOwner string, emailCode *model.EmailCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if user == nil || resourceOwner == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "user and resourceowner should not be nothing")
		}

		agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
		if err != nil {
			return nil, err
		}

		agg, err = agg.AppendEvent(usr_model.UserRegistered, user)
		if err != nil {
			return nil, err
		}

		if emailCode != nil {
			agg, err = agg.AppendEvent(usr_model.UserEmailCodeAdded, emailCode)
			if err != nil {
				return nil, err
			}
		}
		return agg, err
	}
}

func UserDeactivateAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, usr_model.UserDeactivated)
}

func UserReactivateAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, usr_model.UserReactivated)
}

func UserLockAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, usr_model.UserLocked)
}

func UserUnlockAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, user, usr_model.UserUnlocked)
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
		agg, err = agg.AppendEvent(usr_model.InitializedUserCodeCreated, code)
		if err != nil {
			return nil, err
		}
		return agg, err
	}
}

func SkipMfaAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg, err = agg.AppendEvent(usr_model.MfaInitSkipped, nil)
		if err != nil {
			return nil, err
		}
		return agg, err
	}
}

func PasswordChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, password *model.Password) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if password == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9832", "password should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg, err = agg.AppendEvent(usr_model.UserPasswordChanged, password)
		if err != nil {
			return nil, err
		}
		return agg, err
	}
}

func RequestSetPassword(aggCreator *es_models.AggregateCreator, existing *model.User, request *model.RequestPasswordSet) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if request == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8ei2", "password set request should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg, err = agg.AppendEvent(usr_model.UserPasswordSetRequested, request)
		if err != nil {
			return nil, err
		}
		return agg, err
	}
}

func ProfileChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, profile *model.Profile) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if profile == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "profile should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Profile.Changes(profile)
		return agg.AppendEvent(usr_model.UserProfileChanged, changes)
	}
}

func EmailChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, email *model.Email, code *model.EmailCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if email == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "email should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Email.Changes(email)
		agg, err = agg.AppendEvent(usr_model.UserEmailChanged, changes)
		if err != nil {
			return nil, err
		}
		if email.IsEmailVerified {
			return agg.AppendEvent(usr_model.UserEmailVerified, code)
		}
		if code != nil {
			return agg.AppendEvent(usr_model.UserEmailCodeAdded, code)
		}
		return agg, nil
	}
}
func EmailVerifiedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(usr_model.UserEmailVerified, nil)
	}
}

func EmailVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, code *model.EmailCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "code should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(usr_model.UserEmailCodeAdded, code)
	}
}

func PhoneChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, phone *model.Phone, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if phone == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkso3", "phone should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Phone == nil {
			existing.Phone = new(model.Phone)
		}
		changes := existing.Phone.Changes(phone)
		agg, err = agg.AppendEvent(usr_model.UserPhoneChanged, changes)
		if err != nil {
			return nil, err
		}
		if phone.IsPhoneVerified {
			return agg.AppendEvent(usr_model.UserPhoneVerified, code)
		}
		if code != nil {
			return agg.AppendEvent(usr_model.UserPhoneCodeAdded, code)
		}
		return agg, nil
	}
}
func PhoneVerifiedAggregate(aggCreator *es_models.AggregateCreator, existing *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(usr_model.UserPhoneVerified, nil)
	}
}

func PhoneVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dsue2", "code should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(usr_model.UserPhoneCodeAdded, code)
	}
}

func AddressChangeAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, address *model.Address) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if address == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "address should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Address == nil {
			existing.Address = new(model.Address)
		}
		changes := existing.Address.Changes(address)
		agg, err = agg.AppendEvent(usr_model.UserAddressChanged, changes)
		if err != nil {
			return nil, err
		}
		return agg, nil
	}
}

func MfaOTPAddAggregate(aggCreator *es_models.AggregateCreator, existing *model.User, address *model.Address) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if address == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "address should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		if existing.Address == nil {
			existing.Address = new(model.Address)
		}
		changes := existing.Address.Changes(address)
		agg, err = agg.AppendEvent(usr_model.UserAddressChanged, changes)
		if err != nil {
			return nil, err
		}
		return agg, nil
	}
}
