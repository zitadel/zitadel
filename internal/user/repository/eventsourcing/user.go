package eventsourcing

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
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

func UserExternalIDPUniqueQuery(externalIDPUserID string) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserExternalIDPAggregate).
		AggregateIDFilter(externalIDPUserID).
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
	if user == nil || user.Machine == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "Errors.Internal")
	}

	var agg *es_models.Aggregate
	if resourceOwner != "" {
		agg, err = UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
	} else {
		resourceOwner = authz.GetCtxData(ctx).OrgID
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
	if user == nil || user.Human == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "Errors.Internal")
	}

	var agg *es_models.Aggregate
	if resourceOwner != "" {
		agg, err = UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
	} else {
		resourceOwner = authz.GetCtxData(ctx).OrgID
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
	return append(uniqueAggregates, agg), nil
}

func UserRegisterAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, externalIDP *model.ExternalIDP, resourceOwner string, initCode *model.InitUserCode, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	if user == nil || resourceOwner == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-duxk2", "user, resourceowner, initcode must be set")
	}

	if user.Human == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ekuEA", "user must be type human")
	}

	aggregates := make([]*es_models.Aggregate, 0)
	agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, resourceOwner, user.AggregateID)
	if err != nil {
		return nil, err
	}

	agg, err = agg.AppendEvent(model.HumanRegistered, user)
	if err != nil {
		return nil, err
	}
	if initCode != nil {
		agg, err = agg.AppendEvent(model.InitializedHumanCodeAdded, initCode)
		if err != nil {
			return nil, err
		}
	}
	if user.Email != nil && user.EmailAddress != "" && user.IsEmailVerified {
		agg, err = agg.AppendEvent(model.HumanEmailVerified, nil)
		if err != nil {
			return nil, err
		}
	}

	if externalIDP != nil {
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate, iam_es_model.IAMAggregate).
			AggregateIDsFilter()

		if !userLoginMustBeDomain {
			validation := addUserNameAndIDPConfigExistingValidation(user.UserName, externalIDP)
			agg.SetPrecondition(validationQuery, validation)
		} else {
			validation := addIDPConfigExistingValidation(externalIDP)
			agg.SetPrecondition(validationQuery, validation)
		}

		agg, err = agg.AppendEvent(model.HumanExternalIDPAdded, externalIDP)
		uniqueExternalIDPAggregate, err := reservedUniqueExternalIDPAggregate(ctx, aggCreator, resourceOwner, externalIDP)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, uniqueExternalIDPAggregate)
	} else if !userLoginMustBeDomain {
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate).
			AggregateIDsFilter()
		validation := addUserNameValidation(user.UserName)
		agg.SetPrecondition(validationQuery, validation)
	}

	uniqueAggregates, err := getUniqueUserAggregates(ctx, aggCreator, user.UserName, user.EmailAddress, resourceOwner, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, uniqueAggregates...)
	return append(aggregates, agg), nil
}

func getUniqueUserAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, userName, emailAddress, resourceOwner string, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	userNameAggregate, err := reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, userName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}

	return []*es_models.Aggregate{
		userNameAggregate,
	}, nil
}

func reservedUniqueUserNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, userName string, userLoginMustBeDomain bool) (*es_models.Aggregate, error) {
	if userLoginMustBeDomain && resourceOwner == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-wWfgH", "Errors.Internal")
	}

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

func releasedUniqueUserNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, userName string, userLoginMustBeDomain bool) (aggregate *es_models.Aggregate, err error) {
	if userLoginMustBeDomain && resourceOwner == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sK9Kg", "Errors.Internal")
	}

	uniqueUserName := userName
	if userLoginMustBeDomain {
		uniqueUserName = userName + resourceOwner
	}

	aggregate, err = aggCreator.NewAggregate(ctx, uniqueUserName, model.UserUserNameAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, uniqueUserName, model.UserUserNameAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserUserNameReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserUserNameUniqueQuery(uniqueUserName), isEventValidation(aggregate, model.UserUserNameReleased)), nil
}

func changeUniqueUserNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner, oldUsername, username string, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	aggregates := make([]*es_models.Aggregate, 2)
	var err error
	aggregates[0], err = releasedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, oldUsername, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	aggregates[1], err = reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, username, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return aggregates, nil
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

func UserRemoveAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	agg, err := UserAggregate(ctx, aggCreator, user)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.UserRemoved, nil)
	if err != nil {
		return nil, err
	}
	uniqueAgg, err := releasedUniqueUserNameAggregate(ctx, aggCreator, user.ResourceOwner, user.UserName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAgg,
	}, nil
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
		return agg.AppendEvent(model.HumanMFAInitSkipped, nil)
	}
}

func PasswordChangeAggregate(aggCreator *es_models.AggregateCreator, user *model.User, passwordChange *model.PasswordChange) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if passwordChange == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9832", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanPasswordChanged, passwordChange)
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

func ResendInitialPasswordAggregate(aggCreator *es_models.AggregateCreator, user *model.User, code *usr_model.InitUserCode, email string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if code == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dfs3q", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		if email != "" && user.Email != nil && email != user.Email.EmailAddress {
			agg, err = agg.AppendEvent(model.HumanEmailChanged, map[string]interface{}{"email": email})
			if err != nil {
				return nil, err
			}
		}
		return agg.AppendEvent(model.InitializedHumanCodeAdded, code)
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

func ExternalLoginCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, user *model.User, check *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, user.AggregateID)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanExternalLoginCheckSucceeded, check)
	}
}

func TokenAddedAggregate(aggCreator *es_models.AggregateCreator, user *model.User, token *model.Token) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, user.AggregateID)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserTokenAdded, token)
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

func EmailChangeAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, email *model.Email, code *model.EmailCode) (*es_models.Aggregate, error) {
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
	return agg, nil
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

func MFAOTPAddAggregate(aggCreator *es_models.AggregateCreator, user *model.User, otp *model.OTP) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if otp == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dkx9s", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMFAOTPAdded, otp)
	}
}

func MFAOTPVerifyAggregate(aggCreator *es_models.AggregateCreator, user *model.User, userAgentID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMFAOTPVerified, &model.OTPVerified{UserAgentID: userAgentID})
	}
}

func MFAOTPCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, user *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sd5DA", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMFAOTPCheckSucceeded, authReq)
	}
}

func MFAOTPCheckFailedAggregate(aggCreator *es_models.AggregateCreator, user *model.User, authReq *model.AuthRequest) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if authReq == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-64sd6", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMFAOTPCheckFailed, authReq)
	}
}

func MFAOTPRemoveAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.HumanMFAOTPRemoved, nil)
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
	aggregates, err := changeUniqueUserNameAggregate(ctx, aggCreator, user.ResourceOwner, user.UserName, tempName, false)
	if err != nil {
		return nil, err
	}
	userAggregate, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, user.AggregateID)
	if err != nil {
		return nil, err
	}
	userAggregate, err = userAggregate.AppendEvent(model.DomainClaimed, map[string]interface{}{"userName": tempName})
	if err != nil {
		return nil, err
	}
	return append(aggregates, userAggregate), nil
}

func DomainClaimedSentAggregate(aggCreator *es_models.AggregateCreator, user *model.User) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAggregate(ctx, aggCreator, user)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.DomainClaimedSent, nil)
	}
}

func ExternalIDPAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, externalIDPs ...*model.ExternalIDP) ([]*es_models.Aggregate, error) {
	if externalIDPs == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Di9os", "Errors.Internal")
	}
	aggregates := make([]*es_models.Aggregate, 0)
	agg, err := UserAggregate(ctx, aggCreator, user)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(org_es_model.OrgAggregate, iam_es_model.IAMAggregate).
		AggregateIDsFilter()

	validation := addIDPConfigExistingValidation(externalIDPs...)
	agg.SetPrecondition(validationQuery, validation)
	for _, externalIDP := range externalIDPs {
		agg, err = agg.AppendEvent(model.HumanExternalIDPAdded, externalIDP)
		uniqueAggregate, err := reservedUniqueExternalIDPAggregate(ctx, aggCreator, "", externalIDP)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, uniqueAggregate)
	}
	return append(aggregates, agg), nil
}

func ExternalIDPRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, externalIDP *model.ExternalIDP, cascade bool) ([]*es_models.Aggregate, error) {
	if externalIDP == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mlo0s", "Errors.Internal")
	}

	aggregates := make([]*es_models.Aggregate, 0)
	agg, err := UserAggregateOverwriteContext(ctx, aggCreator, user, user.ResourceOwner, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	if cascade {
		agg, err = agg.AppendEvent(model.HumanExternalIDPCascadeRemoved, externalIDP)
	} else {
		agg, err = agg.AppendEvent(model.HumanExternalIDPRemoved, externalIDP)
	}
	uniqueReleasedAggregate, err := releasedUniqueExternalIDPAggregate(ctx, aggCreator, externalIDP, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, uniqueReleasedAggregate)
	return append(aggregates, agg), nil
}

func reservedUniqueExternalIDPAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, resourceOwner string, externalIDP *model.ExternalIDP) (*es_models.Aggregate, error) {
	uniqueExternlIDP := externalIDP.IDPConfigID + externalIDP.UserID
	aggregate, err := aggCreator.NewAggregate(ctx, uniqueExternlIDP, model.UserExternalIDPAggregate, model.UserVersion, 0)
	if resourceOwner != "" {
		aggregate, err = aggCreator.NewAggregate(ctx, uniqueExternlIDP, model.UserExternalIDPAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwner))
	}
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.HumanExternalIDPReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserExternalIDPUniqueQuery(uniqueExternlIDP), isEventValidation(aggregate, model.HumanExternalIDPReserved)), nil
}

func releasedUniqueExternalIDPAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, externalIDP *model.ExternalIDP, resourceOwnerID string) (aggregate *es_models.Aggregate, err error) {
	uniqueExternlIDP := externalIDP.IDPConfigID + externalIDP.UserID
	aggregate, err = aggCreator.NewAggregate(ctx, uniqueExternlIDP, model.UserExternalIDPAggregate, model.UserVersion, 0, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(authz.GetCtxData(ctx).UserID))

	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.HumanExternalIDPReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserExternalIDPUniqueQuery(uniqueExternlIDP), isEventValidation(aggregate, model.HumanExternalIDPReleased)), nil
}

func UsernameChangedAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, user *model.User, oldUsername string, userLoginMustBeDomain bool) ([]*es_models.Aggregate, error) {
	aggregates, err := changeUniqueUserNameAggregate(ctx, aggCreator, user.ResourceOwner, oldUsername, user.UserName, userLoginMustBeDomain)
	if err != nil {
		return nil, err
	}
	userAggregate, err := UserAggregate(ctx, aggCreator, user)
	if err != nil {
		return nil, err
	}
	userAggregate, err = userAggregate.AppendEvent(model.UserUserNameChanged, map[string]interface{}{"userName": user.UserName})
	if err != nil {
		return nil, err
	}
	if !userLoginMustBeDomain {
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate).
			AggregateIDsFilter()

		validation := addUserNameValidation(user.UserName)
		userAggregate.SetPrecondition(validationQuery, validation)
	}
	return append(aggregates, userAggregate), nil
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
			domains = handleDomainEvents(domains, event)
		}
		return handleCheckDomainAllowedAsUsername(domains, userName)
	}
}

func handleDomainEvents(domains []*org_es_model.OrgDomain, event *es_models.Event) []*org_es_model.OrgDomain {
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
	return domains
}

func handleCheckDomainAllowedAsUsername(domains []*org_es_model.OrgDomain, userName string) error {
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

func addIDPConfigExistingValidation(externalIDPs ...*model.ExternalIDP) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		iamLoginPolicy := new(iam_es_model.LoginPolicy)
		orgPolicyExisting := false
		orgLoginPolicy := new(iam_es_model.LoginPolicy)
		for _, event := range events {
			switch event.AggregateType {
			case org_es_model.OrgAggregate:
				orgPolicyExisting = handleOrgLoginPolicy(event, orgPolicyExisting, orgLoginPolicy)
			case iam_es_model.IAMAggregate:
				handleIAMLoginPolicy(event, iamLoginPolicy)
			}
		}
		return handleIDPConfigExisting(iamLoginPolicy, orgLoginPolicy, orgPolicyExisting, externalIDPs...)
	}
}

func handleIDPConfigExisting(iamLoginPolicy, orgLoginPolicy *iam_es_model.LoginPolicy, orgPolicyExisting bool, externalIDPs ...*model.ExternalIDP) error {
	if orgPolicyExisting {
		if !orgLoginPolicy.AllowExternalIdp {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Wmi9s", "Errors.User.ExternalIDP.NotAllowed")
		}
		for _, externalIDP := range externalIDPs {
			existing := false
			for _, provider := range orgLoginPolicy.IDPProviders {
				if provider.IDPConfigID == externalIDP.IDPConfigID {
					existing = true
					break
				}
			}
			if !existing {
				return errors.ThrowPreconditionFailed(nil, "EVENT-Ms9it", "Errors.User.ExternalIDP.IDPConfigNotExisting")
			}
		}
	}
	if !iamLoginPolicy.AllowExternalIdp {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Ns7uf", "Errors.User.ExternalIDP.NotAllowed")
	}
	for _, externalIDP := range externalIDPs {
		existing := false
		for _, provider := range iamLoginPolicy.IDPProviders {
			if provider.IDPConfigID == externalIDP.IDPConfigID {
				existing = true
				break
			}
		}
		if !existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Ms9it", "Errors.User.ExternalIDP.IDPConfigNotExisting")
		}
	}
	return nil
}

func addUserNameAndIDPConfigExistingValidation(userName string, externalIDP *model.ExternalIDP) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		domains := make([]*org_es_model.OrgDomain, 0)
		iamLoginPolicy := new(iam_es_model.LoginPolicy)
		orgPolicyExisting := false
		orgLoginPolicy := new(iam_es_model.LoginPolicy)

		for _, event := range events {
			domains = handleDomainEvents(domains, event)

			switch event.AggregateType {
			case org_es_model.OrgAggregate:
				orgPolicyExisting = handleOrgLoginPolicy(event, orgPolicyExisting, orgLoginPolicy)
			case iam_es_model.IAMAggregate:
				handleIAMLoginPolicy(event, iamLoginPolicy)
			}
		}
		err := handleCheckDomainAllowedAsUsername(domains, userName)
		if err != nil {
			return err
		}
		return handleIDPConfigExisting(iamLoginPolicy, orgLoginPolicy, orgPolicyExisting, externalIDP)
	}
}

func handleOrgLoginPolicy(event *es_models.Event, orgPolicyExisting bool, orgLoginPolicy *iam_es_model.LoginPolicy) bool {
	switch event.Type {
	case org_es_model.LoginPolicyAdded:
		orgPolicyExisting = true
		orgLoginPolicy.SetData(event)
	case org_es_model.LoginPolicyChanged:
		orgLoginPolicy.SetData(event)
	case org_es_model.LoginPolicyRemoved:
		orgPolicyExisting = false
	case org_es_model.LoginPolicyIDPProviderAdded:
		idp := new(iam_es_model.IDPProvider)
		idp.SetData(event)
		orgLoginPolicy.IDPProviders = append(orgLoginPolicy.IDPProviders, idp)
	case org_es_model.LoginPolicyIDPProviderRemoved, org_es_model.LoginPolicyIDPProviderCascadeRemoved:
		idp := new(iam_es_model.IDPProvider)
		idp.SetData(event)
		for i, provider := range orgLoginPolicy.IDPProviders {
			if provider.IDPConfigID == idp.IDPConfigID {
				orgLoginPolicy.IDPProviders[i] = orgLoginPolicy.IDPProviders[len(orgLoginPolicy.IDPProviders)-1]
				orgLoginPolicy.IDPProviders[len(orgLoginPolicy.IDPProviders)-1] = nil
				orgLoginPolicy.IDPProviders = orgLoginPolicy.IDPProviders[:len(orgLoginPolicy.IDPProviders)-1]
				break
			}
		}
	}
	return orgPolicyExisting
}

func handleIAMLoginPolicy(event *es_models.Event, iamLoginPolicy *iam_es_model.LoginPolicy) {
	switch event.Type {
	case iam_es_model.LoginPolicyAdded, iam_es_model.LoginPolicyChanged:
		iamLoginPolicy.SetData(event)
	case iam_es_model.LoginPolicyIDPProviderAdded:
		idp := new(iam_es_model.IDPProvider)
		idp.SetData(event)
		iamLoginPolicy.IDPProviders = append(iamLoginPolicy.IDPProviders, idp)
	case iam_es_model.LoginPolicyIDPProviderRemoved, iam_es_model.LoginPolicyIDPProviderCascadeRemoved:
		idp := new(iam_es_model.IDPProvider)
		idp.SetData(event)
		for i, provider := range iamLoginPolicy.IDPProviders {
			if provider.IDPConfigID == idp.IDPConfigID {
				iamLoginPolicy.IDPProviders[i] = iamLoginPolicy.IDPProviders[len(iamLoginPolicy.IDPProviders)-1]
				iamLoginPolicy.IDPProviders[len(iamLoginPolicy.IDPProviders)-1] = nil
				iamLoginPolicy.IDPProviders = iamLoginPolicy.IDPProviders[:len(iamLoginPolicy.IDPProviders)-1]
				break
			}
		}
	}
}
