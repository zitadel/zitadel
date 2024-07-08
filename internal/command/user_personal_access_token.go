package command

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddPat struct {
	ExpirationDate time.Time
	Scopes         []string
}

type PersonalAccessToken struct {
	models.ObjectRoot

	ExpirationDate  time.Time
	Scopes          []string
	AllowedUserType domain.UserType

	TokenID string
	Token   string
}

func NewPersonalAccessToken(resourceOwner string, userID string, expirationDate time.Time, scopes []string, allowedUserType domain.UserType) *PersonalAccessToken {
	return &PersonalAccessToken{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		ExpirationDate:  expirationDate,
		Scopes:          scopes,
		AllowedUserType: allowedUserType,
	}
}

func (pat *PersonalAccessToken) content() error {
	if pat.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-xs0k2n", "Errors.ResourceOwnerMissing")
	}
	if pat.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-0pzb1", "Errors.User.UserIDMissing")
	}
	if pat.TokenID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-68xm2o", "Errors.IDMissing")
	}
	return nil
}

func (pat *PersonalAccessToken) valid() (err error) {
	if err := pat.content(); err != nil {
		return err
	}
	pat.ExpirationDate, err = domain.ValidateExpirationDate(pat.ExpirationDate)
	return err
}

func (pat *PersonalAccessToken) checkAggregate(ctx context.Context, filter preparation.FilterToQueryReducer) error {
	userWriteModel, err := userWriteModelByID(ctx, filter, pat.AggregateID, pat.ResourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(userWriteModel.UserState) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Dggw2", "Errors.User.NotFound")
	}
	if pat.AllowedUserType != domain.UserTypeUnspecified && userWriteModel.UserType != pat.AllowedUserType {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Df2f1", "Errors.User.WrongType")
	}
	return nil
}

func (c *Commands) AddPersonalAccessToken(ctx context.Context, pat *PersonalAccessToken) (_ *domain.ObjectDetails, err error) {
	if pat.TokenID == "" {
		pat.TokenID, err = id_generator.Next()
		if err != nil {
			return nil, err
		}
	}
	validation := prepareAddPersonalAccessToken(pat, c.keyAlgorithm)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareAddPersonalAccessToken(pat *PersonalAccessToken, algorithm crypto.EncryptionAlgorithm) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := pat.valid(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			if err := pat.checkAggregate(ctx, filter); err != nil {
				return nil, err
			}
			writeModel, err := getPersonalAccessTokenWriteModelByID(ctx, filter, pat.AggregateID, pat.TokenID, pat.ResourceOwner)
			if err != nil {
				return nil, err
			}

			pat.Token, err = createToken(algorithm, writeModel.TokenID, writeModel.AggregateID)
			if err != nil {
				return nil, err
			}

			return []eventstore.Command{
				user.NewPersonalAccessTokenAddedEvent(
					ctx,
					UserAggregateFromWriteModel(&writeModel.WriteModel),
					pat.TokenID,
					pat.ExpirationDate,
					pat.Scopes,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) RemovePersonalAccessToken(ctx context.Context, pat *PersonalAccessToken) (*domain.ObjectDetails, error) {
	validation := prepareRemovePersonalAccessToken(pat)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareRemovePersonalAccessToken(pat *PersonalAccessToken) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := pat.content(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			writeModel, err := getPersonalAccessTokenWriteModelByID(ctx, filter, pat.AggregateID, pat.TokenID, pat.ResourceOwner)
			if err != nil {
				return nil, err
			}
			if !writeModel.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.PAT.NotFound")
			}
			return []eventstore.Command{
				user.NewPersonalAccessTokenRemovedEvent(
					ctx,
					UserAggregateFromWriteModel(&writeModel.WriteModel),
					pat.TokenID,
				),
			}, nil
		}, nil
	}
}

func createToken(algorithm crypto.EncryptionAlgorithm, tokenID, userID string) (string, error) {
	encrypted, err := algorithm.Encrypt([]byte(tokenID + ":" + userID))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func getPersonalAccessTokenWriteModelByID(ctx context.Context, filter preparation.FilterToQueryReducer, userID, tokenID, resourceOwner string) (_ *PersonalAccessTokenWriteModel, err error) {
	writeModel := NewPersonalAccessTokenWriteModel(userID, tokenID, resourceOwner)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
