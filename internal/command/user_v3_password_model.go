package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type PasswordV3WriteModel struct {
	eventstore.WriteModel
	UserID string

	EncodedHash    string
	ChangeRequired bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration
	GeneratorID      string
	VerificationID   string

	checkPermission domain.PermissionCheck
}

func (wm *PasswordV3WriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func NewPasswordV3WriteModel(resourceOwner, id string, checkPermission domain.PermissionCheck) *PasswordV3WriteModel {
	return &PasswordV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID:          id,
		checkPermission: checkPermission,
	}
}

func (wm *PasswordV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.PasswordCreatedEvent:
			wm.UserID = e.UserID
			wm.EncodedHash = e.EncodedHash
			wm.ChangeRequired = e.ChangeRequired
			wm.Code = nil
		case *authenticator.PasswordDeletedEvent:
			wm.UserID = ""
			wm.EncodedHash = ""
			wm.ChangeRequired = false
			wm.Code = nil
		case *authenticator.PasswordCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
			wm.GeneratorID = e.GeneratorID
		case *authenticator.PasswordCodeSentEvent:
			wm.GeneratorID = e.GeneratorInfo.GetID()
			wm.VerificationID = e.GeneratorInfo.GetVerificationID()
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(authenticator.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.PasswordCreatedType,
			authenticator.PasswordDeletedType,
			authenticator.PasswordCodeAddedType,
		).Builder()
}

func (wm *PasswordV3WriteModel) NewCreate(
	ctx context.Context,
	encodeHash string,
	changeRequired bool,
) ([]eventstore.Command, error) {
	return []eventstore.Command{
		authenticator.NewPasswordCreatedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			wm.UserID,
			encodeHash,
			changeRequired,
		),
	}, nil
}

func (wm *PasswordV3WriteModel) NewAddCode(
	ctx context.Context,
	notificationType domain.NotificationType,
	urlTemplate string,
	codeReturned bool,
	code func(context.Context, domain.NotificationType) (*EncryptedCode, string, error),
) (_ []eventstore.Command, plainCode string, err error) {
	crypt, generatorID, err := code(ctx, notificationType)
	if err != nil {
		return nil, "", err
	}

	events := []eventstore.Command{
		authenticator.NewPasswordCodeAddedEvent(ctx,
			AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()),
			crypt.CryptedCode(),
			crypt.CodeExpiry(),
			notificationType,
			urlTemplate,
			codeReturned,
			generatorID,
		),
	}
	if codeReturned {
		plainCode = crypt.Plain
	}
	return events, plainCode, nil
}

func (wm *PasswordV3WriteModel) NewDelete(ctx context.Context) ([]eventstore.Command, error) {
	if err := wm.Exists(); err != nil {
		return nil, err
	}
	return []eventstore.Command{authenticator.NewPasswordDeletedEvent(ctx, AuthenticatorAggregateFromWriteModel(wm.GetWriteModel()))}, nil
}

func (wm *PasswordV3WriteModel) Exists() error {
	if wm.EncodedHash == "" {
		return zerrors.ThrowNotFound(nil, "COMMAND-Joi3utDPIh", "Errors.User.Password.NotFound")
	}
	return nil
}
