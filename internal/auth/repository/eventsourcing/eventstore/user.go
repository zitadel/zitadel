package eventstore

import (
	"context"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	key_model "github.com/caos/zitadel/internal/key/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	key_view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	SearchLimit    uint64
	Eventstore     eventstore.Eventstore
	OrgEvents      *org_event.OrgEventstore
	View           *view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
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

func (repo *UserRepo) GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNView, error) {
	user, err := repo.UserByID(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "USER-9kF98", "Errors.User.NotHuman")
	}
	return user.HumanView.PasswordlessTokens, nil
}

func (repo *UserRepo) UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(userSessions))
	for i, session := range userSessions {
		userIDs[i] = session.UserID
	}
	return userIDs, nil
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*model.UserView, error) {
	user, err := repo.View.UserByID(id)
	if err != nil {
		return nil, err
	}
	events, err := repo.getUserEvents(ctx, id, user.Sequence)
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

func (repo *UserRepo) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	return repo.getUserEvents(ctx, id, sequence)
}

func (repo *UserRepo) UserByLoginName(ctx context.Context, loginname string) (*model.UserView, error) {
	user, err := repo.View.UserByLoginName(loginname)
	if err != nil {
		return nil, err
	}
	events, err := repo.getUserEvents(ctx, user.ID, user.Sequence)
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
	changes, err := repo.getUserChanges(ctx, authz.GetCtxData(ctx).UserID, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierID
		user, _ := repo.UserByID(ctx, change.ModifierID)
		if user != nil {
			if user.HumanView != nil {
				change.ModifierName = user.DisplayName
			}
			if user.MachineView != nil {
				change.ModifierName = user.MachineView.Name
			}
		}
	}
	return changes, nil
}

func (repo *UserRepo) MachineKeyByID(ctx context.Context, keyID string) (*key_model.AuthNKeyView, error) {
	key, err := repo.View.AuthNKeyByID(keyID)
	if err != nil {
		return nil, err
	}
	return key_view_model.AuthNKeyToModel(key), nil
}

func (r *UserRepo) getUserChanges(ctx context.Context, userID string, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error) {
	query := usr_view.ChangesQuery(userID, lastSequence, limit, sortAscending)

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
