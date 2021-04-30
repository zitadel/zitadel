package command

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanRefreshTokenWriteModel struct {
	eventstore.WriteModel

	TokenID      string
	RefreshToken string
	//IsRefreshTokenVerified bool
	//
	//Code             *crypto.CryptoValue
	//CodeCreationDate time.Time
	//CodeExpiry       time.Duration

	UserState      domain.UserState
	IdleExpiration time.Time
	Expiration     time.Time
}

func NewHumanRefreshTokenWriteModel(userID, resourceOwner, tokenID string) *HumanRefreshTokenWriteModel {
	return &HumanRefreshTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		TokenID: tokenID,
	}
}

func (wm *HumanRefreshTokenWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanRefreshTokenAddedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanRefreshTokenRenewedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanRefreshTokenRemovedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *HumanRefreshTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanRefreshTokenAddedEvent:
			wm.TokenID = e.TokenID
			wm.RefreshToken = e.TokenID
			wm.IdleExpiration = e.CreationDate().Add(e.IdleExpiration)
			wm.Expiration = e.CreationDate().Add(e.Expiration)
			wm.UserState = domain.UserStateActive
		case *user.HumanRefreshTokenRenewedEvent:
			if wm.UserState == domain.UserStateActive {
				wm.RefreshToken = e.RefreshToken
			}
			wm.RefreshToken = e.RefreshToken
			wm.IdleExpiration = e.CreationDate().Add(e.IdleExpiration)
		case *user.HumanRefreshTokenRemovedEvent:
			wm.UserState = domain.UserStateDeleted
			//case *user.HumanInitializedCheckSucceededEvent:
			//	wm.UserState = domain.UserStateActive
			//case *user.HumanRefreshTokenChangedEvent:
			//	wm.RefreshToken = e.RefreshTokenAddress
			//	wm.IsRefreshTokenVerified = false
			//	wm.Code = nil
			//case *user.HumanRefreshTokenCodeAddedEvent:
			//	wm.Code = e.Code
			//	wm.CodeCreationDate = e.CreationDate()
			//	wm.CodeExpiry = e.Expiry
			//case *user.HumanRefreshTokenVerifiedEvent:
			//	wm.IsRefreshTokenVerified = true
			//	wm.Code = nil
			//case *user.UserRemovedEvent:
			//	wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanRefreshTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.UserV1AddedType,
			//user.HumanAddedType,
			//user.UserV1RegisteredType,
			//user.HumanRegisteredType,
			//user.UserV1InitialCodeAddedType,
			//user.HumanInitialCodeAddedType,
			//user.UserV1InitializedCheckSucceededType,
			//user.HumanInitializedCheckSucceededType,
			//user.UserV1RefreshTokenChangedType,
			//user.HumanRefreshTokenChangedType,
			//user.UserV1RefreshTokenCodeAddedType,
			//user.HumanRefreshTokenCodeAddedType,
			//user.UserV1RefreshTokenVerifiedType,
			//user.HumanRefreshTokenVerifiedType,
			user.HumanRefreshTokenAddedType,
			user.HumanRefreshTokenRenewedType,
			user.HumanRefreshTokenRemovedType,
			user.UserRemovedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
