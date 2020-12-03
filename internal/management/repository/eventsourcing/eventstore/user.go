package eventstore

import (
	"context"

	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserRepo struct {
	es_int.Eventstore
	SearchLimit     uint64
	UserEvents      *usr_event.UserEventstore
	OrgEvents       *org_event.OrgEventstore
	UserGrantEvents *usr_grant_event.UserGrantEventStore
	View            *view.View
	SystemDefaults  systemdefaults.SystemDefaults
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		user = new(model.UserView)
	}
	events, esErr := repo.UserEvents.UserEventsByID(ctx, id, user.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.User.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4Fm9s", "Errors.User.NotFound")
	}
	return model.UserToModel(&userCopy), nil
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.CreateUser(ctx, user, pwPolicyView, orgPolicyView)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(policyResourceOwner)
	if err != nil && caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.RegisterUser(ctx, user, pwPolicyView, orgPolicyView, resourceOwner)
}

func (repo *UserRepo) DeactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.DeactivateUser(ctx, id)
}

func (repo *UserRepo) ReactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.ReactivateUser(ctx, id)
}

func (repo *UserRepo) LockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.LockUser(ctx, id)
}

func (repo *UserRepo) UnlockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.UnlockUser(ctx, id)
}

func (repo *UserRepo) RemoveUser(ctx context.Context, id string) error {
	aggregates := make([]*es_models.Aggregate, 0)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	user, agg, err := repo.UserEvents.PrepareRemoveUser(ctx, id, orgPolicyView)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg...)

	// remove user_grants
	usergrants, err := repo.View.UserGrantsByUserID(id)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		_, aggs, err := repo.UserGrantEvents.PrepareRemoveUserGrant(ctx, grant.ID, true)
		if err != nil {
			return err
		}
		for _, agg := range aggs {
			aggregates = append(aggregates, agg)
		}
	}

	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, user.AppendEvents, aggregates...)
}

func (repo *UserRepo) SearchUsers(ctx context.Context, request *usr_model.UserSearchRequest) (*usr_model.UserSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestUserSequence()
	logging.Log("EVENT-Lcn7d").OnError(sequenceErr).Warn("could not read latest user sequence")
	users, count, err := repo.View.SearchUsers(request)
	if err != nil {
		return nil, err
	}
	result := &usr_model.UserSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.UsersToModel(users),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *UserRepo) UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*usr_model.UserChanges, error) {
	changes, err := repo.UserEvents.UserChanges(ctx, id, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		// user, _ := repo.View.UserByID(change.ModifierID)
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierID)
		if user != nil {
			if user.Human != nil {
				change.ModifierName = user.Human.DisplayName
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *UserRepo) GetUserByLoginNameGlobal(ctx context.Context, loginName string) (*usr_model.UserView, error) {
	user, err := repo.View.GetGlobalUserByLoginName(loginName)
	if err != nil {
		return nil, err
	}
	return model.UserToModel(user), nil
}

func (repo *UserRepo) IsUserUnique(ctx context.Context, userName, email string) (bool, error) {
	return repo.View.IsUserUnique(userName, email)
}

func (repo *UserRepo) UserMfas(ctx context.Context, userID string) ([]*usr_model.MultiFactor, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-xx0hV", "Errors.User.NotHuman")
	}
	if user.OTPState == usr_model.MfaStateUnspecified {
		return []*usr_model.MultiFactor{}, nil
	}
	return []*usr_model.MultiFactor{{Type: usr_model.MfaTypeOTP, State: user.OTPState}}, nil
}

func (repo *UserRepo) RemoveOTP(ctx context.Context, userID string) error {
	return repo.UserEvents.RemoveOTP(ctx, userID)
}

func (repo *UserRepo) SetOneTimePassword(ctx context.Context, password *usr_model.Password) (*usr_model.Password, error) {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && caos_errs.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(policy)
	return repo.UserEvents.SetOneTimePassword(ctx, pwPolicyView, password)
}

func (repo *UserRepo) RequestSetPassword(ctx context.Context, id string, notifyType usr_model.NotificationType) error {
	return repo.UserEvents.RequestSetPassword(ctx, id, notifyType)
}

func (repo *UserRepo) ResendInitialMail(ctx context.Context, userID, email string) error {
	return repo.UserEvents.ResendInitialMail(ctx, userID, email)
}

func (repo *UserRepo) ProfileByID(ctx context.Context, userID string) (*usr_model.Profile, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-gDFC2", "Errors.User.NotHuman")
	}
	return user.GetProfile()
}

func (repo *UserRepo) SearchExternalIDPs(ctx context.Context, request *usr_model.ExternalIDPSearchRequest) (*usr_model.ExternalIDPSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, seqErr := repo.View.GetLatestExternalIDPSequence()
	logging.Log("EVENT-Qs7uf").OnError(seqErr).Warn("could not read latest external idp sequence")
	externalIDPS, count, err := repo.View.SearchExternalIDPs(request)
	if err != nil {
		return nil, err
	}
	result := &usr_model.ExternalIDPSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ExternalIDPViewsToModel(externalIDPS),
	}
	if seqErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *UserRepo) RemoveExternalIDP(ctx context.Context, externalIDP *usr_model.ExternalIDP) error {
	return repo.UserEvents.RemoveExternalIDP(ctx, externalIDP)
}

func (repo *UserRepo) ChangeMachine(ctx context.Context, machine *usr_model.Machine) (*usr_model.Machine, error) {
	return repo.UserEvents.ChangeMachine(ctx, machine)
}

func (repo *UserRepo) GetMachineKey(ctx context.Context, userID, keyID string) (*usr_model.MachineKeyView, error) {
	key, err := repo.View.MachineKeyByIDs(userID, keyID)
	if err != nil {
		return nil, err
	}
	return model.MachineKeyToModel(key), nil
}

func (repo *UserRepo) SearchMachineKeys(ctx context.Context, request *usr_model.MachineKeySearchRequest) (*usr_model.MachineKeySearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, seqErr := repo.View.GetLatestMachineKeySequence()
	logging.Log("EVENT-Sk8fs").OnError(seqErr).Warn("could not read latest user sequence")
	keys, count, err := repo.View.SearchMachineKeys(request)
	if err != nil {
		return nil, err
	}
	result := &usr_model.MachineKeySearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.MachineKeysToModel(keys),
	}
	if seqErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *UserRepo) AddMachineKey(ctx context.Context, key *usr_model.MachineKey) (*usr_model.MachineKey, error) {
	return repo.UserEvents.AddMachineKey(ctx, key)
}

func (repo *UserRepo) RemoveMachineKey(ctx context.Context, userID, keyID string) error {
	return repo.UserEvents.RemoveMachineKey(ctx, userID, keyID)
}

func (repo *UserRepo) ChangeProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) ChangeUsername(ctx context.Context, userID, userName string) error {
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.ChangeUsername(ctx, userID, userName, orgPolicyView)
}

func (repo *UserRepo) EmailByID(ctx context.Context, userID string) (*usr_model.Email, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-pt7HY", "Errors.User.NotHuman")
	}
	return user.GetEmail()
}

func (repo *UserRepo) ChangeEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	return repo.UserEvents.ChangeEmail(ctx, email)
}

func (repo *UserRepo) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, userID)
}

func (repo *UserRepo) PhoneByID(ctx context.Context, userID string) (*usr_model.Phone, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-hliQl", "Errors.User.NotHuman")
	}
	return user.GetPhone()
}

func (repo *UserRepo) ChangePhone(ctx context.Context, email *usr_model.Phone) (*usr_model.Phone, error) {
	return repo.UserEvents.ChangePhone(ctx, email)
}

func (repo *UserRepo) RemovePhone(ctx context.Context, userID string) error {
	return repo.UserEvents.RemovePhone(ctx, userID)
}

func (repo *UserRepo) CreatePhoneVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreatePhoneVerificationCode(ctx, userID)
}

func (repo *UserRepo) AddressByID(ctx context.Context, userID string) (*usr_model.Address, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-LQh4I", "Errors.User.NotHuman")
	}
	return user.GetAddress()
}

func (repo *UserRepo) ChangeAddress(ctx context.Context, address *usr_model.Address) (*usr_model.Address, error) {
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) SearchUserMemberships(ctx context.Context, request *usr_model.UserMembershipSearchRequest) (*usr_model.UserMembershipSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestUserMembershipSequence()
	logging.Log("EVENT-Dn7sf").OnError(sequenceErr).Warn("could not read latest user sequence")

	result := handleSearchUserMembershipsPermissions(ctx, request, sequence)
	if result != nil {
		return result, nil
	}

	memberships, count, err := repo.View.SearchUserMemberships(request)
	if err != nil {
		return nil, err
	}
	result = &usr_model.UserMembershipSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.UserMembershipsToModel(memberships),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func handleSearchUserMembershipsPermissions(ctx context.Context, request *usr_model.UserMembershipSearchRequest, sequence *repository.CurrentSequence) *usr_model.UserMembershipSearchResponse {
	permissions := authz.GetAllPermissionsFromCtx(ctx)
	iamPerm := authz.HasGlobalExplicitPermission(permissions, iamMemberReadPerm)
	orgPerm := authz.HasGlobalExplicitPermission(permissions, orgMemberReadPerm)
	projectPerm := authz.HasGlobalExplicitPermission(permissions, projectMemberReadPerm)
	projectGrantPerm := authz.HasGlobalExplicitPermission(permissions, projectGrantMemberReadPerm)
	if iamPerm && orgPerm && projectPerm && projectGrantPerm {
		return nil
	}
	if !iamPerm {
		request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyMemberType, Method: global_model.SearchMethodNotEquals, Value: usr_model.MemberTypeIam})
	}
	if !orgPerm {
		request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyMemberType, Method: global_model.SearchMethodNotEquals, Value: usr_model.MemberTypeOrganisation})
	}

	ids := authz.GetExplicitPermissionCtxIDs(permissions, projectMemberReadPerm)
	ids = append(ids, authz.GetExplicitPermissionCtxIDs(permissions, projectGrantMemberReadPerm)...)
	if _, q := request.GetSearchQuery(usr_model.UserMembershipSearchKeyObjectID); q != nil {
		containsID := false
		for _, id := range ids {
			if id == q.Value {
				containsID = true
				break
			}
		}
		if !containsID {
			result := &usr_model.UserMembershipSearchResponse{
				Offset:      request.Offset,
				Limit:       request.Limit,
				TotalResult: uint64(0),
				Result:      []*usr_model.UserMembershipView{},
			}
			if sequence != nil {
				result.Sequence = sequence.CurrentSequence
				result.Timestamp = sequence.CurrentTimestamp
			}
			return result
		}
	}
	request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyObjectID, Method: global_model.SearchMethodIsOneOf, Value: ids})
	return nil
}
