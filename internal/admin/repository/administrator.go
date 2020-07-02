package repository

import (
	"context"
	"github.com/caos/zitadel/internal/view/model"
)

type AdministratorRepository interface {
	GetFailedEvents(context.Context) ([]*model.FailedEvent, error)
	RemoveFailedEvent(context.Context, *model.FailedEvent) error
	GetViews(context.Context) ([]*model.View, error)
	ClearView(ctx context.Context, db, view string) error
}
