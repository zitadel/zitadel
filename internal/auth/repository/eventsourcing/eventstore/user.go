package eventstore

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/user/model"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
)

type UserRepo struct {
	SearchLimit     uint64
	Eventstore      v1.Eventstore
	View            *view.View
	Query           *query.Queries
	SystemDefaults  systemdefaults.SystemDefaults
	PrefixAvatarURL string
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
}

func (repo *UserRepo) UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(userSessions))
	for _, session := range userSessions {
		if session.State == int32(domain.UserSessionStateActive) {
			userIDs = append(userIDs, session.UserID)
		}
	}
	return userIDs, nil
}

func (repo *UserRepo) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*models.Event, error) {
	return repo.getUserEvents(ctx, id, sequence)
}

func (repo *UserRepo) MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error) {
	changes, err := repo.getUserChanges(ctx, authz.GetCtxData(ctx).UserID, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		change.ModifierLoginName = change.ModifierID
		user, _ := repo.Query.GetUserByID(ctx, change.ModifierID, false)
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

func (r *UserRepo) getUserChanges(ctx context.Context, userID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error) {
	query := usr_view.ChangesQuery(userID, lastSequence, limit, sortAscending, retention)

	events, err := r.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-g9HCv").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-htuG9", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-6cAxe", "Errors.User.NoChanges")
	}

	result := make([]*model.UserChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-8GTGS").OnError(err).Debug("unable to parse timestamp")
		change := &model.UserChange{
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

	return &model.UserChanges{
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
