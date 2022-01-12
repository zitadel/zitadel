package eventstore

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/query"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type ProjectRepo struct {
	v1.Eventstore
	SearchLimit     uint64
	View            *view.View
	Roles           []string
	IAMID           string
	PrefixAvatarURL string
	Query           *query.Queries
}

func (repo *ProjectRepo) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ProjectChanges, error) {
	changes, err := repo.getProjectChanges(ctx, id, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		change.ModifierLoginName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
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

func (repo *ProjectRepo) ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ApplicationChanges, error) {
	changes, err := repo.getApplicationChanges(ctx, projectID, appID, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		change.ModifierLoginName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
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

func (repo *ProjectRepo) GetProjectMemberRoles(ctx context.Context) ([]string, error) {
	iam, err := repo.Query.IAMByID(ctx, repo.IAMID)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	global := authz.GetCtxData(ctx).OrgID == iam.GlobalOrgID
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT") && !strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			if global && !strings.HasSuffix(roleMap, "GLOBAL") {
				continue
			}
			roles = append(roles, roleMap)
		}
	}
	return roles, nil
}

func (repo *ProjectRepo) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *ProjectRepo) userByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		user = new(usr_es_model.UserView)
	}
	events, esErr := repo.getUserEvents(ctx, id, user.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4n8Fs", "Errors.User.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2m0Fs", "Errors.User.NotFound")
	}
	return usr_es_model.UserToModel(&userCopy, repo.PrefixAvatarURL), nil
}

func (r *ProjectRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}

func (repo *ProjectRepo) getProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ProjectChanges, error) {
	query := proj_view.ChangesQuery(id, lastSequence, limit, sortAscending, retention)

	events, err := repo.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, caos_errs.ThrowInternal(err, "EVENT-328b1", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*proj_model.ProjectChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ProjectChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			var data interface{}
			if strings.Contains(change.EventType, "application") {
				data = new(proj_model.Application)
			} else {
				data = new(proj_model.Project)
			}
			err = json.Unmarshal(event.Data, data)
			logging.Log("EVENT-NCkpN").OnError(err).Debug("unable to unmarshal data")
			change.Data = data
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ProjectChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func (repo *ProjectRepo) getProjectEvents(ctx context.Context, id string, sequence uint64) ([]*models.Event, error) {
	query, err := proj_view.ProjectByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}

func (repo *ProjectRepo) getApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ApplicationChanges, error) {
	query := proj_view.ChangesQuery(projectID, lastSequence, limit, sortAscending, retention)

	events, err := repo.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, caos_errs.ThrowInternal(err, "EVENT-sw6Ku", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-9IHLP", "Errors.Changes.NotFound")
	}

	result := make([]*proj_model.ApplicationChange, 0)
	for _, event := range events {
		if !strings.Contains(event.Type.String(), "application") || event.Data == nil {
			continue
		}

		app := new(proj_model.Application)
		err := json.Unmarshal(event.Data, app)
		logging.Log("EVENT-GIiKD").OnError(err).Debug("unable to unmarshal data")
		if app.AppID != appID {
			continue
		}

		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-MJzeN").OnError(err).Debug("unable to parse timestamp")

		result = append(result, &proj_model.ApplicationChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
			Data:       app,
		})
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ApplicationChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}, nil
}
