package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type PATV3WriteModel struct {
	eventstore.WriteModel
	UserID         string
	ExpirationDate time.Time
	Scopes         []string

	checkPermission domain.PermissionCheck
}

func (wm *PATV3WriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func NewPATV3WriteModel(resourceOwner, userID, id string, checkPermission domain.PermissionCheck) *PATV3WriteModel {
	return &PATV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID:          userID,
		checkPermission: checkPermission,
	}
}

func (wm *PATV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.PATCreatedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.UserID = e.UserID
			wm.Scopes = e.Scopes
			wm.ExpirationDate = e.ExpirationDate
		case *authenticator.PATDeletedEvent:
			wm.UserID = ""
			wm.Scopes = nil
			wm.ExpirationDate = time.Time{}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PATV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(authenticator.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.PATCreatedType,
			authenticator.PATDeletedType,
		).Builder()
}

func (wm *PATV3WriteModel) checkPermissionWrite(ctx context.Context) error {
	if wm.UserID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, wm.UserID); err != nil {
		return err
	}
	return nil
}

func (wm *PATV3WriteModel) NewCreate(
	ctx context.Context,
	expirationDate time.Time,
	scopes []string,
) ([]eventstore.Command, error) {
	if err := wm.NotExists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewPATCreatedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			wm.UserID,
			expirationDate,
			scopes,
		),
	}, nil
}

func (wm *PATV3WriteModel) NewDelete(ctx context.Context) ([]eventstore.Command, error) {
	if err := wm.Exists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewPATDeletedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
		),
	}, nil
}

func (wm *PATV3WriteModel) Exists() error {
	if len(wm.Scopes) == 0 {
		return zerrors.ThrowNotFound(nil, "COMMAND-ur4kxtxIhW", "Errors.User.NotFound")
	}
	return nil
}

func (wm *PATV3WriteModel) NotExists() error {
	if err := wm.Exists(); err != nil {
		return nil
	}
	return zerrors.ThrowAlreadyExists(nil, "COMMAND-iBM2bOhvYH", "Errors.User.AlreadyExists")
}
