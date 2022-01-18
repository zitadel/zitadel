package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	v1.Eventstore
	SearchLimit     uint64
	View            *view.View
	Query           *query.Queries
	SystemDefaults  systemdefaults.SystemDefaults
	PrefixAvatarURL string
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		user = new(model.UserView)
	}

	events, esErr := repo.getUserEvents(ctx, id, user.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.User.NotFound")
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
		return nil, errors.ThrowNotFound(nil, "EVENT-4Fm9s", "Errors.User.NotFound")
	}
	return model.UserToModel(&userCopy, repo.PrefixAvatarURL), nil
}

func (repo *UserRepo) UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*usr_model.UserChanges, error) {
	changes, err := repo.getUserChanges(ctx, id, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		change.ModifierLoginName = change.ModifierID
		user, _ := repo.Query.GetUserByID(ctx, change.ModifierID)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.Human != nil {
				change.ModifierName = user.Human.DisplayName
				change.ModifierAvatarURL = domain.AvatarURL(repo.PrefixAvatarURL, user.ResourceOwner, user.Human.AvatarKey)
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
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
