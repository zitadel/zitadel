package repository

import (
	"context"
	"github.com/caos/zitadel/internal/view/model"
)

type AdministratorRepository interface {
	GetFailedEvents(context.Context) ([]*model.FailedEvent, error)
	RemoveFailedEvent(context.Context, *model.FailedEvent) error
	GetViews() ([]*model.View, error)
	GetSpoolerDiv(db, viewName string) int64
	ClearView(ctx context.Context, db, viewName string) error
}
