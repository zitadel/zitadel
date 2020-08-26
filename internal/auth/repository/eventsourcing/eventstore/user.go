package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	policy_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	"github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	Eventstore   eventstore.Eventstore
	UserEvents   *user_event.UserEventstore
	OrgEvents    *org_event.OrgEventstore
	PolicyEvents *policy_event.PolicyEventstore
	View         *view.View
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.UserEvents.Health(ctx)
}

func (repo *UserRepo) Register(ctx context.Context, registerUser *model.User, orgMember *org_model.OrgMember, resourceOwner string) (*model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := repo.OrgEvents.GetOrgIAMPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	user, aggregates, err := repo.UserEvents.PrepareRegisterUser(ctx, registerUser, pwPolicy, orgPolicy, resourceOwner)
	if err != nil {
		return nil, err
	}
	if orgMember != nil {
		orgMember.UserID = user.AggregateID
		_, memberAggregate, err := repo.OrgEvents.PrepareAddOrgMember(ctx, orgMember, policyResourceOwner)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, memberAggregate)
	}

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, user.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}
	return usr_model.UserToModel(user), nil
}

func (repo *UserRepo) MyUser(ctx context.Context) (*model.UserView, error) {
	return repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) MyProfile(ctx context.Context) (*model.Profile, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return user.GetProfile(), nil
}

func (repo *UserRepo) ChangeMyProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error) {
	if err := checkIDs(ctx, profile.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) MyEmail(ctx context.Context) (*model.Email, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return user.GetEmail(), nil
}

func (repo *UserRepo) ChangeMyEmail(ctx context.Context, email *model.Email) (*model.Email, error) {
	if err := checkIDs(ctx, email.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeEmail(ctx, email)
}

func (repo *UserRepo) VerifyEmail(ctx context.Context, userID, code string) error {
	return repo.UserEvents.VerifyEmail(ctx, userID, code)
}

func (repo *UserRepo) VerifyMyEmail(ctx context.Context, code string) error {
	return repo.UserEvents.VerifyEmail(ctx, authz.GetCtxData(ctx).UserID, code)
}

func (repo *UserRepo) ResendEmailVerificationMail(ctx context.Context, userID string) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, userID)
}

func (repo *UserRepo) ResendMyEmailVerificationMail(ctx context.Context) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) MyPhone(ctx context.Context) (*model.Phone, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return user.GetPhone(), nil
}

func (repo *UserRepo) ChangeMyPhone(ctx context.Context, phone *model.Phone) (*model.Phone, error) {
	if err := checkIDs(ctx, phone.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangePhone(ctx, phone)
}

func (repo *UserRepo) RemoveMyPhone(ctx context.Context) error {
	return repo.UserEvents.RemovePhone(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) VerifyMyPhone(ctx context.Context, code string) error {
	return repo.UserEvents.VerifyPhone(ctx, authz.GetCtxData(ctx).UserID, code)
}

func (repo *UserRepo) ResendMyPhoneVerificationCode(ctx context.Context) error {
	return repo.UserEvents.CreatePhoneVerificationCode(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) MyAddress(ctx context.Context) (*model.Address, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return user.GetAddress(), nil
}

func (repo *UserRepo) ChangeMyAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := checkIDs(ctx, address.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) ChangeMyPassword(ctx context.Context, old, new string) error {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return err
	}
	_, err = repo.UserEvents.ChangePassword(ctx, policy, authz.GetCtxData(ctx).UserID, old, new)
	return err
}

func (repo *UserRepo) ChangePassword(ctx context.Context, userID, old, new string) error {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return err
	}
	_, err = repo.UserEvents.ChangePassword(ctx, policy, userID, old, new)
	return err
}

func (repo *UserRepo) MyUserMfas(ctx context.Context) ([]*model.MultiFactor, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	if user.OTPState == model.MfaStateUnspecified {
		return []*model.MultiFactor{}, nil
	}
	return []*model.MultiFactor{{Type: model.MfaTypeOTP, State: user.OTPState}}, nil
}

func (repo *UserRepo) AddMfaOTP(ctx context.Context, userID string) (*model.OTP, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-Fk93s").OnError(err).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddOTP(ctx, userID, accountName)
}

func (repo *UserRepo) AddMyMfaOTP(ctx context.Context) (*model.OTP, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		logging.Log("EVENT-Ml0sd").OnError(err).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddOTP(ctx, authz.GetCtxData(ctx).UserID, accountName)
}

func (repo *UserRepo) VerifyMfaOTPSetup(ctx context.Context, userID, code string) error {
	return repo.UserEvents.CheckMfaOTPSetup(ctx, userID, code)
}

func (repo *UserRepo) VerifyMyMfaOTPSetup(ctx context.Context, code string) error {
	return repo.UserEvents.CheckMfaOTPSetup(ctx, authz.GetCtxData(ctx).UserID, code)
}

func (repo *UserRepo) RemoveMyMfaOTP(ctx context.Context) error {
	return repo.UserEvents.RemoveOTP(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ResendInitVerificationMail(ctx context.Context, userID string) error {
	_, err := repo.UserEvents.CreateInitializeUserCodeByID(ctx, userID)
	return err
}

func (repo *UserRepo) VerifyInitCode(ctx context.Context, userID, code, password string) error {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return err
	}
	return repo.UserEvents.VerifyInitCode(ctx, policy, userID, code, password)
}

func (repo *UserRepo) SkipMfaInit(ctx context.Context, userID string) error {
	return repo.UserEvents.SkipMfaInit(ctx, userID)
}

func (repo *UserRepo) RequestPasswordReset(ctx context.Context, loginname string) error {
	user, err := repo.View.UserByLoginName(loginname)
	if err != nil {
		return err
	}
	return repo.UserEvents.RequestSetPassword(ctx, user.ID, model.NotificationTypeEmail)
}

func (repo *UserRepo) SetPassword(ctx context.Context, userID, code, password string) error {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return err
	}
	return repo.UserEvents.SetPassword(ctx, policy, userID, code, password)
}

func (repo *UserRepo) SignOut(ctx context.Context, agentID string) error {
	userSessions, err := repo.View.UserSessionsByAgentID(agentID)
	if err != nil {
		return err
	}
	userIDs := make([]string, len(userSessions))
	for i, session := range userSessions {
		userIDs[i] = session.UserID
	}
	return repo.UserEvents.SignOut(ctx, agentID, userIDs)
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*model.UserView, error) {
	user, err := repo.View.UserByID(id)
	if err != nil {
		return nil, err
	}
	events, err := repo.UserEvents.UserEventsByID(ctx, id, user.Sequence)
	if err != nil {
		logging.Log("EVENT-PSoc3").WithError(err).Debug("error retrieving new events")
		return usr_view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_view_model.UserToModel(user), nil
		}
	}
	return usr_view_model.UserToModel(&userCopy), nil
}

func (repo *UserRepo) MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error) {
	changes, err := repo.UserEvents.UserChanges(ctx, authz.GetCtxData(ctx).UserID, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierId)
		if user != nil {
			change.ModifierName = user.DisplayName
		}
	}
	return changes, nil
}

func checkIDs(ctx context.Context, obj es_models.ObjectRoot) error {
	if obj.AggregateID != authz.GetCtxData(ctx).UserID {
		return errors.ThrowPermissionDenied(nil, "EVENT-kFi9w", "object does not belong to user")
	}
	return nil
}
