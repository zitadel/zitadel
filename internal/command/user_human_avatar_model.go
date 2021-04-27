package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanAvatarWriteModel struct {
	eventstore.WriteModel

	AssetID string
	Avatar  []byte

	UserState domain.UserState
}

func NewHumanAvatarWriteModel(userID, resourceOwner string) *HumanAvatarWriteModel {
	return &HumanAvatarWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanAvatarWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.UserState = domain.UserStateActive
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.HumanAvatarChangedEvent:
			wm.AssetID = e.AssetID
		case *user.HumanAvatarRemovedEvent:
			wm.AssetID = ""
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanAvatarWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.UserRemovedType,
			user.HumanAvatarChangedType,
			user.HumanAvatarRemovedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
