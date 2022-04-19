package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

var dbList = []string{"auth", "authz", "adminapi", "notification"}

type AdministratorRepo struct {
	View *view.View
}

func (repo *AdministratorRepo) GetFailedEvents(ctx context.Context) ([]*view_model.FailedEvent, error) {
	allFailedEvents := make([]*view_model.FailedEvent, 0)
	for _, db := range dbList {
		failedEvents, err := repo.View.AllFailedEvents(db)
		if err != nil {
			return nil, err
		}
		for _, failedEvent := range failedEvents {
			allFailedEvents = append(allFailedEvents, repository.FailedEventToModel(failedEvent))
		}
	}
	return allFailedEvents, nil
}

func (repo *AdministratorRepo) RemoveFailedEvent(ctx context.Context, failedEvent *view_model.FailedEvent) error {
	return repo.View.RemoveFailedEvent(failedEvent.Database, repository.FailedEventFromModel(failedEvent))
}

func (repo *AdministratorRepo) GetViews() ([]*view_model.View, error) {
	views := make([]*view_model.View, 0)
	for _, db := range dbList {
		sequences, err := repo.View.AllCurrentSequences(db)
		if err != nil {
			return nil, err
		}
		for _, sequence := range sequences {
			views = append(views, repository.CurrentSequenceToModel(sequence))
		}
	}
	return views, nil
}

func (repo *AdministratorRepo) ClearView(ctx context.Context, database, view string) error {
	return repo.View.ClearView(database, view)
}
