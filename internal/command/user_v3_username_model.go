package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UsernameV3WriteModel struct {
	eventstore.WriteModel
	UserID        string
	Username      string
	IsOrgSpecific bool

	checkPermission domain.PermissionCheck
}

func (wm *UsernameV3WriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func NewUsernameV3WriteModel(resourceOwner, userID, id string, checkPermission domain.PermissionCheck) *UsernameV3WriteModel {
	return &UsernameV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID:          userID,
		checkPermission: checkPermission,
	}
}

func (wm *UsernameV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.UsernameCreatedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.UserID = e.UserID
			wm.Username = e.Username
			wm.IsOrgSpecific = e.IsOrgSpecific
		case *authenticator.UsernameDeletedEvent:
			wm.UserID = ""
			wm.Username = ""
			wm.IsOrgSpecific = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UsernameV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(authenticator.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.UsernameCreatedType,
			authenticator.UsernameDeletedType,
		).Builder()
}

func (wm *UsernameV3WriteModel) checkPermissionWrite(ctx context.Context) error {
	if wm.UserID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, wm.UserID); err != nil {
		return err
	}
	return nil
}

func (wm *UsernameV3WriteModel) NewCreate(
	ctx context.Context,
	isOrgSpecific bool,
	username string,
) ([]eventstore.Command, error) {
	if err := wm.NotExists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewUsernameCreatedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			wm.UserID,
			isOrgSpecific,
			username,
		),
	}, nil
}

func (wm *UsernameV3WriteModel) NewDelete(ctx context.Context) ([]eventstore.Command, error) {
	if err := wm.Exists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewUsernameDeletedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			wm.IsOrgSpecific,
			wm.Username,
		),
	}, nil
}

func (wm *UsernameV3WriteModel) Exists() error {
	if wm.Username == "" {
		return zerrors.ThrowNotFound(nil, "COMMAND-uEii8L6Awp", "Errors.User.NotFound")
	}
	return nil
}

func (wm *UsernameV3WriteModel) NotExists() error {
	if err := wm.Exists(); err != nil {
		return nil
	}
	return zerrors.ThrowAlreadyExists(nil, "COMMAND-rK7ZTzEEGU", "Errors.User.AlreadyExists")
}

func AuthenticatorAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          authenticator.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       authenticator.AggregateVersion,
	}
}
