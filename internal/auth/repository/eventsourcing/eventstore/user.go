package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	key_model "github.com/caos/zitadel/internal/key/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	key_view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	SearchLimit    uint64
	Eventstore     eventstore.Eventstore
	UserEvents     *user_event.UserEventstore
	OrgEvents      *org_event.OrgEventstore
	View           *view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.UserEvents.Health(ctx)
}

func (repo *UserRepo) Register(ctx context.Context, user *model.User, orgMember *org_model.OrgMember, resourceOwner string) (*model.User, error) {
	return repo.registerUser(ctx, user, nil, orgMember, resourceOwner)
}

func (repo *UserRepo) RegisterExternalUser(ctx context.Context, user *model.User, externalIDP *model.ExternalIDP, orgMember *org_model.OrgMember, resourceOwner string) (*model.User, error) {
	return repo.registerUser(ctx, user, externalIDP, orgMember, resourceOwner)
}

func (repo *UserRepo) registerUser(ctx context.Context, registerUser *model.User, externalIDP *model.ExternalIDP, orgMember *org_model.OrgMember, resourceOwner string) (*model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(policyResourceOwner)
	if errors.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(policyResourceOwner)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	user, aggregates, err := repo.UserEvents.PrepareRegisterUser(ctx, registerUser, externalIDP, pwPolicyView, orgPolicyView, resourceOwner)
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
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-H2JIT", "Errors.User.NotHuman")
	}
	return user.GetProfile()
}

func (repo *UserRepo) ChangeMyProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error) {
	if err := checkIDs(ctx, profile.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) SearchMyExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, seqErr := repo.View.GetLatestExternalIDPSequence()
	logging.Log("EVENT-5Jsi8").OnError(seqErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest user sequence")
	request.AppendUserQuery(authz.GetCtxData(ctx).UserID)
	externalIDPS, count, err := repo.View.SearchExternalIDPs(request)
	if err != nil {
		return nil, err
	}
	result := &model.ExternalIDPSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      usr_view_model.ExternalIDPViewsToModel(externalIDPS),
	}
	if seqErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserRepo) AddMyExternalIDP(ctx context.Context, externalIDP *model.ExternalIDP) (*model.ExternalIDP, error) {
	if err := checkIDs(ctx, externalIDP.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.AddExternalIDP(ctx, externalIDP)
}

func (repo *UserRepo) RemoveMyExternalIDP(ctx context.Context, externalIDP *model.ExternalIDP) error {
	if err := checkIDs(ctx, externalIDP.ObjectRoot); err != nil {
		return err
	}
	return repo.UserEvents.RemoveExternalIDP(ctx, externalIDP)
}

func (repo *UserRepo) MyEmail(ctx context.Context) (*model.Email, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-oGRpc", "Errors.User.NotHuman")
	}
	return user.GetEmail()
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
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-DTWJb", "Errors.User.NotHuman")
	}
	return user.GetPhone()
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
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Ok9nI", "Errors.User.NotHuman")
	}
	return user.GetAddress()
}

func (repo *UserRepo) ChangeMyAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := checkIDs(ctx, address.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) ChangeMyPassword(ctx context.Context, old, new string) error {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(policy)
	_, err = repo.UserEvents.ChangePassword(ctx, pwPolicyView, authz.GetCtxData(ctx).UserID, old, new, "")
	return err
}

func (repo *UserRepo) ChangePassword(ctx context.Context, userID, old, new, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(policy)
	_, err = repo.UserEvents.ChangePassword(ctx, pwPolicyView, userID, old, new, userAgentID)
	return err
}

func (repo *UserRepo) MyUserMFAs(ctx context.Context) ([]*model.MultiFactor, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	mfas := make([]*model.MultiFactor, 0)
	if user.OTPState != model.MFAStateUnspecified {
		mfas = append(mfas, &model.MultiFactor{Type: model.MFATypeOTP, State: user.OTPState})
	}
	for _, u2f := range user.U2FTokens {
		mfas = append(mfas, &model.MultiFactor{Type: model.MFATypeU2F, State: u2f.State, Attribute: u2f.Name, ID: u2f.TokenID})
	}
	return mfas, nil
}

func (repo *UserRepo) AddMFAOTP(ctx context.Context, userID string) (*model.OTP, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-Fk93s").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddOTP(ctx, userID, accountName)
}

func (repo *UserRepo) AddMyMFAOTP(ctx context.Context) (*model.OTP, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		logging.Log("EVENT-Ml0sd").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddOTP(ctx, authz.GetCtxData(ctx).UserID, accountName)
}

func (repo *UserRepo) VerifyMFAOTPSetup(ctx context.Context, userID, code, userAgentID string) error {
	return repo.UserEvents.CheckMFAOTPSetup(ctx, userID, code, userAgentID)
}

func (repo *UserRepo) VerifyMyMFAOTPSetup(ctx context.Context, code string) error {
	return repo.UserEvents.CheckMFAOTPSetup(ctx, authz.GetCtxData(ctx).UserID, code, "")
}

func (repo *UserRepo) RemoveMyMFAOTP(ctx context.Context) error {
	return repo.UserEvents.RemoveOTP(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) AddMFAU2F(ctx context.Context, userID string) (*model.WebAuthNToken, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-DAqe1").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddU2F(ctx, userID, accountName, true)
}

func (repo *UserRepo) AddMyMFAU2F(ctx context.Context) (*model.WebAuthNToken, error) {
	userID := authz.GetCtxData(ctx).UserID
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-Ghwl1").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddU2F(ctx, userID, accountName, false)
}

func (repo *UserRepo) VerifyMFAU2FSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error {
	return repo.UserEvents.VerifyU2FSetup(ctx, userID, tokenName, userAgentID, credentialData)
}

func (repo *UserRepo) VerifyMyMFAU2FSetup(ctx context.Context, tokenName string, credentialData []byte) error {
	return repo.UserEvents.VerifyU2FSetup(ctx, authz.GetCtxData(ctx).UserID, tokenName, "", credentialData)
}

func (repo *UserRepo) RemoveMFAU2F(ctx context.Context, userID, webAuthNTokenID string) error {
	return repo.UserEvents.RemoveU2FToken(ctx, userID, webAuthNTokenID)
}

func (repo *UserRepo) RemoveMyMFAU2F(ctx context.Context, webAuthNTokenID string) error {
	return repo.UserEvents.RemoveU2FToken(ctx, authz.GetCtxData(ctx).UserID, webAuthNTokenID)
}

func (repo *UserRepo) GetPasswordless(ctx context.Context, userID string) ([]*model.WebAuthNToken, error) {
	return repo.UserEvents.GetPasswordless(ctx, userID)
}

func (repo *UserRepo) AddPasswordless(ctx context.Context, userID string) (*model.WebAuthNToken, error) {
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-Vj2k1").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddPasswordless(ctx, userID, accountName, true)
}

func (repo *UserRepo) GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNToken, error) {
	return repo.UserEvents.GetPasswordless(ctx, authz.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) AddMyPasswordless(ctx context.Context) (*model.WebAuthNToken, error) {
	userID := authz.GetCtxData(ctx).UserID
	accountName := ""
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		logging.Log("EVENT-AEq21").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get user for loginname")
	} else {
		accountName = user.PreferredLoginName
	}
	return repo.UserEvents.AddPasswordless(ctx, authz.GetCtxData(ctx).UserID, accountName, false)
}

func (repo *UserRepo) VerifyPasswordlessSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error {
	return repo.UserEvents.VerifyPasswordlessSetup(ctx, userID, tokenName, userAgentID, credentialData)
}

func (repo *UserRepo) VerifyMyPasswordlessSetup(ctx context.Context, tokenName string, credentialData []byte) error {
	return repo.UserEvents.VerifyPasswordlessSetup(ctx, authz.GetCtxData(ctx).UserID, tokenName, "", credentialData)
}

func (repo *UserRepo) RemovePasswordless(ctx context.Context, userID, webAuthNTokenID string) error {
	return repo.UserEvents.RemovePasswordlessToken(ctx, userID, webAuthNTokenID)
}

func (repo *UserRepo) RemoveMyPasswordless(ctx context.Context, webAuthNTokenID string) error {
	return repo.UserEvents.RemovePasswordlessToken(ctx, authz.GetCtxData(ctx).UserID, webAuthNTokenID)
}

func (repo *UserRepo) ChangeMyUsername(ctx context.Context, username string) error {
	ctxData := authz.GetCtxData(ctx)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(ctxData.OrgID)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.ChangeUsername(ctx, ctxData.UserID, username, orgPolicyView)
}
func (repo *UserRepo) ResendInitVerificationMail(ctx context.Context, userID string) error {
	_, err := repo.UserEvents.CreateInitializeUserCodeByID(ctx, userID)
	return err
}

func (repo *UserRepo) VerifyInitCode(ctx context.Context, userID, code, password string) error {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(policy)
	return repo.UserEvents.VerifyInitCode(ctx, pwPolicyView, userID, code, password)
}

func (repo *UserRepo) SkipMFAInit(ctx context.Context, userID string) error {
	return repo.UserEvents.SkipMFAInit(ctx, userID)
}

func (repo *UserRepo) RequestPasswordReset(ctx context.Context, loginname string) error {
	user, err := repo.View.UserByLoginName(loginname)
	if err != nil {
		return err
	}
	return repo.UserEvents.RequestSetPassword(ctx, user.ID, model.NotificationTypeEmail)
}

func (repo *UserRepo) SetPassword(ctx context.Context, userID, code, password, userAgentID string) error {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(policy)
	return repo.UserEvents.SetPassword(ctx, pwPolicyView, userID, code, password, userAgentID)
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
		logging.Log("EVENT-PSoc3").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return usr_view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_view_model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "EVENT-vZ8us", "Errors.User.NotFound")
	}
	return usr_view_model.UserToModel(&userCopy), nil
}

func (repo *UserRepo) MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error) {
	changes, err := repo.UserEvents.UserChanges(ctx, authz.GetCtxData(ctx).UserID, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierID)
		if user != nil {
			if user.Human != nil {
				change.ModifierName = user.DisplayName
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *UserRepo) ChangeUsername(ctx context.Context, userID, username string) error {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(policyResourceOwner)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.ChangeUsername(ctx, userID, username, orgPolicyView)
}

func checkIDs(ctx context.Context, obj es_models.ObjectRoot) error {
	if obj.AggregateID != authz.GetCtxData(ctx).UserID {
		return errors.ThrowPermissionDenied(nil, "EVENT-kFi9w", "object does not belong to user")
	}
	return nil
}

func (repo *UserRepo) MachineKeyByID(ctx context.Context, keyID string) (*key_model.AuthNKeyView, error) {
	key, err := repo.View.AuthNKeyByID(keyID)
	if err != nil {
		return nil, err
	}
	return key_view_model.AuthNKeyToModel(key), nil
}
