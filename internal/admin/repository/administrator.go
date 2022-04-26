package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/view/model"
)

type AdministratorRepository interface {
	GetFailedEvents(context.Context) ([]*model.FailedEvent, error)
	RemoveFailedEvent(context.Context, *model.FailedEvent) error
	GetViews() ([]*model.View, error)
	ClearView(ctx context.Context, db, viewName string) error
}
