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

type PublicKeyV3WriteModel struct {
	eventstore.WriteModel
	UserID         string
	ExpirationDate time.Time
	PublicKey      []byte

	checkPermission domain.PermissionCheck
}

func (wm *PublicKeyV3WriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func NewPublicKeyV3WriteModel(resourceOwner, userID, id string, checkPermission domain.PermissionCheck) *PublicKeyV3WriteModel {
	return &PublicKeyV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID:          userID,
		checkPermission: checkPermission,
	}
}

func (wm *PublicKeyV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.PublicKeyCreatedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.UserID = e.UserID
			wm.PublicKey = e.PublicKey
			wm.ExpirationDate = e.ExpirationDate
		case *authenticator.PublicKeyDeletedEvent:
			wm.UserID = ""
			wm.PublicKey = nil
			wm.ExpirationDate = time.Time{}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PublicKeyV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(authenticator.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.PublicKeyCreatedType,
			authenticator.PublicKeyDeletedType,
		).Builder()
}

func (wm *PublicKeyV3WriteModel) checkPermissionWrite(ctx context.Context) error {
	if wm.UserID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, wm.UserID); err != nil {
		return err
	}
	return nil
}

func (wm *PublicKeyV3WriteModel) NewCreate(
	ctx context.Context,
	expirationDate time.Time,
	publicKey []byte,
) ([]eventstore.Command, error) {
	if err := wm.NotExists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewPublicKeyCreatedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			wm.UserID,
			expirationDate,
			publicKey,
		),
	}, nil
}

func (wm *PublicKeyV3WriteModel) NewDelete(ctx context.Context) ([]eventstore.Command, error) {
	if err := wm.Exists(); err != nil {
		return nil, err
	}
	if err := wm.checkPermissionWrite(ctx); err != nil {
		return nil, err
	}
	return []eventstore.Command{
		authenticator.NewPublicKeyDeletedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
		),
	}, nil
}

func (wm *PublicKeyV3WriteModel) Exists() error {
	if len(wm.PublicKey) == 0 {
		return zerrors.ThrowNotFound(nil, "COMMAND-CqNteIqtCt", "Errors.User.NotFound")
	}
	return nil
}

func (wm *PublicKeyV3WriteModel) NotExists() error {
	if err := wm.Exists(); err != nil {
		return nil
	}
	return zerrors.ThrowAlreadyExists(nil, "COMMAND-QkVpJv0DqA", "Errors.User.AlreadyExists")
}
