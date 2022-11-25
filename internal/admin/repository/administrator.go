package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/view/model"
)

type AdministratorRepository interface {
	GetFailedEvents(ctx context.Context, instanceID string) ([]*model.FailedEvent, error)
	RemoveFailedEvent(context.Context, *model.FailedEvent) error
	GetViews(instanceID string) ([]*model.View, error)
	ClearView(ctx context.Context, db, viewName string) error
}
