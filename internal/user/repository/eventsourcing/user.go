package eventsourcing

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
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

func MachineCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwner string, userLoginMustBeDomain bool) (_ []*es_models.Aggregate, err error) {
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

	agg, err = agg.AppendEvent(model.MachineAdded, user)
	if err != nil {
		return nil, err
	}

	userNameAggregate, err := reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, user.UserName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		userNameAggregate,
	}, nil
}

func HumanCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, initCode *model.InitUserCode, phoneCode *model.PhoneCode, resourceOwner string, userLoginMustBeDomain bool) (_ []*es_models.Aggregate, err error) {
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

	agg, err = agg.AppendEvent(model.HumanAdded, user)
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

	uniqueAggregates, err := getUniqueUserAggregates(ctx, aggCreator, user.UserName, user.EmailAddress, resourceOwner, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregates[0],
		uniqueAggregates[1],
	}, nil
}

func UserRegisterAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, resourceOwner string, initCode *model.InitUserCode, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	if user == nil || resourceOwner == "" || initCode == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "user, resourceowner, initcode must be set")
	}

	if user.Human == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ekuEA", "user must be type human")
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
	agg, err = agg.AppendEvent(model.HumanRegistered, user)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.InitializedHumanCodeAdded, initCode)
	if err != nil {
		return nil, err
	}
	uniqueAggregates, err := getUniqueUserAggregates(ctx, aggCreator, user.UserName, user.EmailAddress, resourceOwner, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregates[0],
		uniqueAggregates[1],
	}, nil
}

func getUniqueUserAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, userName, emailAddress, resourceOwner string, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	userNameAggregate, err := reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, userName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}

	emailAggregate, err := reservedUniqueEmailAggregate(ctx, aggCreator, resourceOwner, emailAddress)
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

func releasedUniqueUserNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, username string) (aggregate *es_models.Aggregate, err error) {
	aggregate, err = aggCreator.NewAggregate(ctx, username, model.UserUserNameAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, username, model.UserUserNameAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserUserNameReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserUserNameUniqueQuery(username), isEventValidation(aggregate, model.UserUserNameReleased)), nil
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

func UserInitCodeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, code *model.InitUserCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8i23", "code should not be nil")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedHumanCodeAdded, code)
	}
}

func UserInitCodeSentAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedHumanCodeSent, nil)
	}
}

func InitCodeVerifiedAggregate(aggCreator *es_models.AggregateCreator, user *model.User, password *model.Password) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		if user.Email != nil && !user.Email.IsEmailVerified {
			agg, err = agg.AppendEvent(model.HumanEmailVerified, nil)
			if err != nil {
				return nil, err
			}
		}
		if password != nil && password.Secret != nil {
			agg, err = agg.AppendEvent(model.HumanPasswordChanged, password)
			if err != nil {
				return nil, err
			}
		}
		return agg.AppendEvent(model.InitializedHumanCheckSucceeded, nil)
	}
}

func InitCodeCheckFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.InitializedHumanCheckFailed, nil)
	}
}

func SkipMfaAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaInitSkipped, nil)
	}
}

func PasswordChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, password *model.Password) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if password == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9832", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordChanged, password)
	}
}

func PasswordCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, user *model.User, check *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordCheckSucceeded, check)
	}
}
func PasswordCheckFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User, check *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordCheckFailed, check)
	}
}

func RequestSetPassword(aggCreator *es_models.AggregateCreator, user *model.User, request *model.PasswordCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if request == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8ei2", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordCodeAdded, request)
	}
}

func PasswordCodeSentAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordCodeSent, nil)
	}
}

func MachineChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, machine *model.Machine) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if machine == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		changes := user.Machine.Changes(machine)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0spow", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.MachineChanged, changes)
	}
}

func ProfileChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, profile *model.Profile) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if profile == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		changes := user.Profile.Changes(profile)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0spow", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.HumanProfileChanged, changes)
	}
}

func EmailChangeAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, email *model.Email, code *model.EmailCode) ([]*es_models.Aggregate, error) {
	if email == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "Errors.Internal")
	}
	if (!email.IsEmailVerified && code == nil) || (email.IsEmailVerified && code != nil) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-id934", "Errors.Internal")
	}
	changes := user.Email.Changes(email)
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-s90pw", "Errors.NoChangesFound")
	}
	aggregates := make([]*es_models.Aggregate, 0, 4)
	reserveEmailAggregate, err := reservedUniqueEmailAggregate(ctx, aggCreator, "", email.EmailAddress)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, reserveEmailAggregate)
	releaseEmailAggregate, err := releasedUniqueEmailAggregate(ctx, aggCreator, "", user.EmailAddress)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, releaseEmailAggregate)
	agg, err := UserAggregate(ctx, aggCreator, user)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.HumanEmailChanged, changes)
	if err != nil {
		return nil, err
	}
	if user.Email == nil {
		user.Email = new(model.Email)
	}
	if email.IsEmailVerified {
		agg, err = agg.AppendEvent(model.HumanEmailVerified, code)
		if err != nil {
			return nil, err
		}
	}
	if code != nil {
		agg, err = agg.AppendEvent(model.HumanEmailCodeAdded, code)
		if err != nil {
			return nil, err
		}
	}
	return append(aggregates, agg), nil
}

func EmailVerifiedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanEmailVerified, nil)
	}
}

func EmailVerificationFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanEmailVerificationFailed, nil)
	}
}

func EmailVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, code *model.EmailCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dki8s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanEmailCodeAdded, code)
	}
}

func EmailCodeSentAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanEmailCodeSent, nil)
	}
}

func PhoneChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, phone *model.Phone, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if phone == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkso3", "Errors.Internal")
		}
		if (!phone.IsPhoneVerified && code == nil) || (phone.IsPhoneVerified && code != nil) {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dksi8", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		if user.Phone == nil {
			user.Phone = new(model.Phone)
		}
		changes := user.Phone.Changes(phone)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sp0oc", "Errors.NoChangesFound")
		}
		agg, err = agg.AppendEvent(model.HumanPhoneChanged, changes)
		if err != nil {
			return nil, err
		}
		if phone.IsPhoneVerified {
			return agg.AppendEvent(model.HumanPhoneVerified, code)
		}
		if code != nil {
			return agg.AppendEvent(model.HumanPhoneCodeAdded, code)
		}
		return agg, nil
	}
}

func PhoneRemovedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPhoneRemoved, nil)
	}
}

func PhoneVerifiedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPhoneVerified, nil)
	}
}

func PhoneVerificationFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPhoneVerificationFailed, nil)
	}
}

func PhoneVerificationCodeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, code *model.PhoneCode) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dsue2", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPhoneCodeAdded, code)
	}
}

func PhoneCodeSentAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPhoneCodeSent, nil)
	}
}

func AddressChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, address *model.Address) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if address == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		if user.Address == nil {
			user.Address = new(model.Address)
		}
		changes := user.Address.Changes(address)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2tszw", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.HumanAddressChanged, changes)
	}
}

func MfaOTPAddAggregate(aggCreator *es_models.AggregateCreator, user *model.User, otp *model.OTP) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if otp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaOtpAdded, otp)
	}
}

func MfaOTPVerifyAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaOtpVerified, nil)
	}
}

func MfaOTPCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, user *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sd5DA", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaOtpCheckSucceeded, authReq)
	}
}

func MfaOTPCheckFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-64sd6", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaOtpCheckFailed, authReq)
	}
}

func MfaOTPRemoveAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMfaOtpRemoved, nil)
	}
}

func SignOutAggregates(aggCreator *es_models.AggregateCreator, users []*model.User, agentID string) func(ctx context.Context) ([]*es_models.Aggregate, error) {
	return func(ctx context.Context) ([]*es_models.Aggregate, error) {
		aggregates := make([]*es_models.Aggregate, len(users))
		for i, user := range users {
			agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, user.AggregateID)
			if err != nil {
				return nil, err
			}
			agg.AppendEvent(model.HumanSignedOut, map[string]interface{}{"userAgentID": agentID})
			aggregates[i] = agg
		}
		return aggregates, nil
	}
}

func DomainClaimedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, tempName string) ([]*es_models.Aggregate, error) {
	aggregates := make([]*es_models.Aggregate, 3)
	userAggregate, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, user.AggregateID)
	if err != nil {
		return nil, err
	}
	userAggregate, err = userAggregate.AppendEvent(model.DomainClaimed, map[string]interface{}{"userName": tempName})
	if err != nil {
		return nil, err
	}
	aggregates[0] = userAggregate
	releasedUniqueAggregate, err := releasedUniqueUserNameAggregate(ctx, aggCreator, user.ResourceOwner, user.UserName)
	if err != nil {
		return nil, err
	}
	aggregates[1] = releasedUniqueAggregate
	reservedUniqueAggregate, err := reservedUniqueUserNameAggregate(ctx, aggCreator, user.ResourceOwner, tempName, false)
	if err != nil {
		return nil, err
	}
	aggregates[2] = reservedUniqueAggregate
	return aggregates, nil
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
				domains = append(domains, domain)
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
						break
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
