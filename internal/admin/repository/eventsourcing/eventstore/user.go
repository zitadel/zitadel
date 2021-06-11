package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/user/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	SearchLimit     uint64
	Eventstore      v1.Eventstore
	View            *view.View
	SystemDefaults  systemdefaults.SystemDefaults
	PrefixAvatarURL string
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
}

func (repo *UserRepo) SearchUsers(ctx context.Context, request *model.UserSearchRequest) (*model.UserSearchResponse, error) {
	sequence, sequenceErr := repo.View.GetLatestUserSequence()
	logging.Log("EVENT-Gdbgfw").OnError(sequenceErr).Warn("could not read latest user sequence")
	users, count, err := repo.View.SearchUsers(request)
	if err != nil {
		return nil, err
	}
	result := &model.UserSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      usr_view_model.UsersToModel(users, repo.PrefixAvatarURL),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}
