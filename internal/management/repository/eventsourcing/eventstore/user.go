package eventstore

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/domain"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	key_model "github.com/caos/zitadel/internal/key/model"
	key_view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserRepo struct {
	v1.Eventstore
	SearchLimit     uint64
	View            *view.View
	SystemDefaults  systemdefaults.SystemDefaults
	PrefixAvatarURL string
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		user = new(model.UserView)
	}

	events, esErr := repo.getUserEvents(ctx, id, user.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.User.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return model.UserToModel(user, repo.PrefixAvatarURL), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return model.UserToModel(user, repo.PrefixAvatarURL), nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4Fm9s", "Errors.User.NotFound")
	}
	return model.UserToModel(&userCopy, repo.PrefixAvatarURL), nil
}

func (repo *UserRepo) SearchUsers(ctx context.Context, request *usr_model.UserSearchRequest, ensureLimit bool) (*usr_model.UserSearchResponse, error) {
	if ensureLimit {
		err := request.EnsureLimit(repo.SearchLimit)

		if err != nil {
			return nil, err
		}
	}
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
		Result:      model.UsersToModel(users, repo.PrefixAvatarURL),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserRepo) UserIDsByDomain(ctx context.Context, domain string) ([]string, error) {
	return repo.View.UserIDsByDomain(domain)
}

func (repo *UserRepo) UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*usr_model.UserChanges, error) {
	changes, err := repo.getUserChanges(ctx, id, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		change.ModifierLoginName = change.ModifierID
		user, _ := repo.UserByID(ctx, change.ModifierID)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.HumanView != nil {
				change.ModifierName = user.HumanView.DisplayName
				change.ModifierAvatarURL = user.HumanView.AvatarURL
			}
			if user.MachineView != nil {
				change.ModifierName = user.MachineView.Name
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
	return model.UserToModel(user, repo.PrefixAvatarURL), nil
}

func (repo *UserRepo) IsUserUnique(ctx context.Context, userName, email string) (bool, error) {
	return repo.View.IsUserUnique(userName, email)
}

func (repo *UserRepo) GetMetadataByKey(ctx context.Context, userID, resourceOwner, key string) (*domain.Metadata, error) {
	data, err := repo.View.MetadataByKeyAndResourceOwner(userID, resourceOwner, key)
	if err != nil {
		return nil, err
	}
	return iam_model.MetadataViewToDomain(data), nil
}

func (repo *UserRepo) SearchMetadata(ctx context.Context, userID, resourceOwner string, req *domain.MetadataSearchRequest) (*domain.MetadataSearchResponse, error) {
	err := req.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestUserSequence()
	logging.Log("EVENT-m0ds3").OnError(sequenceErr).Warn("could not read latest user sequence")
	req.AppendAggregateIDQuery(userID)
	req.AppendResourceOwnerQuery(resourceOwner)
	metadata, count, err := repo.View.SearchMetadata(req)
	if err != nil {
		return nil, err
	}
	result := &domain.MetadataSearchResponse{
		Offset:      req.Offset,
		Limit:       req.Limit,
		TotalResult: count,
		Result:      iam_model.MetadataViewsToDomain(metadata),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserRepo) UserMFAs(ctx context.Context, userID string) ([]*usr_model.MultiFactor, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-xx0hV", "Errors.User.NotHuman")
	}
	mfas := make([]*usr_model.MultiFactor, 0)
	if user.OTPState != usr_model.MFAStateUnspecified {
		mfas = append(mfas, &usr_model.MultiFactor{Type: usr_model.MFATypeOTP, State: user.OTPState})
	}
	for _, u2f := range user.U2FTokens {
		mfas = append(mfas, &usr_model.MultiFactor{Type: usr_model.MFATypeU2F, State: u2f.State, Attribute: u2f.Name, ID: u2f.TokenID})
	}
	return mfas, nil
}

func (repo *UserRepo) GetPasswordless(ctx context.Context, userID string) ([]*usr_model.WebAuthNView, error) {
	user, err := repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9anf8", "Errors.User.NotHuman")
	}
	return user.HumanView.PasswordlessTokens, nil
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
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserRepo) ExternalIDPsByIDPConfigID(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error) {
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return model.ExternalIDPViewsToModel(externalIDPs), nil
}

func (repo *UserRepo) ExternalIDPsByIDPConfigIDAndResourceOwner(ctx context.Context, idpConfigID, resourceOwner string) ([]*usr_model.ExternalIDPView, error) {
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigIDAndResourceOwner(idpConfigID, resourceOwner)
	if err != nil {
		return nil, err
	}
	return model.ExternalIDPViewsToModel(externalIDPs), nil
}

func (repo *UserRepo) GetMachineKey(ctx context.Context, userID, keyID string) (*key_model.AuthNKeyView, error) {
	key, err := repo.View.AuthNKeyByIDs(userID, keyID)
	if err != nil {
		return nil, err
	}
	return key_view_model.AuthNKeyToModel(key), nil
}

func (repo *UserRepo) SearchMachineKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, seqErr := repo.View.GetLatestAuthNKeySequence()
	logging.Log("EVENT-Sk8fs").OnError(seqErr).Warn("could not read latest authn key sequence")
	keys, count, err := repo.View.SearchAuthNKeys(request)
	if err != nil {
		return nil, err
	}
	result := &key_model.AuthNKeySearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      key_view_model.AuthNKeysToModel(keys),
	}
	if seqErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
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

func (repo *UserRepo) SearchUserMemberships(ctx context.Context, request *usr_model.UserMembershipSearchRequest) (*usr_model.UserMembershipSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserRepo) UserMembershipsByUserID(ctx context.Context, userID string) ([]*usr_model.UserMembershipView, error) {
	memberships, err := repo.View.UserMembershipsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return model.UserMembershipsToModel(memberships), nil
}

func (r *UserRepo) getUserChanges(ctx context.Context, userID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*usr_model.UserChanges, error) {
	query := usr_view.ChangesQuery(userID, lastSequence, limit, sortAscending, retention)

	events, err := r.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-g9HCv").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-htuG9", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-6cAxe", "Errors.User.NoChanges")
	}

	result := make([]*usr_model.UserChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-8GTGS").OnError(err).Debug("unable to parse timestamp")
		change := &usr_model.UserChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierID: event.EditorUser,
			Sequence:   event.Sequence,
		}

		//TODO: now all types should be unmarshalled, e.g. password
		// if len(event.Data) != 0 {
		// 	user := new(model.User)
		// 	err := json.Unmarshal(event.Data, user)
		// 	logging.Log("EVENT-Rkg7X").OnError(err).Debug("unable to unmarshal data")
		// 	change.Data = user
		// }

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

func (r *UserRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
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
		request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyMemberType, Method: domain.SearchMethodNotEquals, Value: usr_model.MemberTypeIam})
	}
	if !orgPerm {
		request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyMemberType, Method: domain.SearchMethodNotEquals, Value: usr_model.MemberTypeOrganisation})
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
				result.Timestamp = sequence.LastSuccessfulSpoolerRun
			}
			return result
		}
	}
	request.Queries = append(request.Queries, &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyObjectID, Method: domain.SearchMethodIsOneOf, Value: ids})
	return nil
}
