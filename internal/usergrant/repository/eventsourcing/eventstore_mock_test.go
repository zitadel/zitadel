package eventsourcing

import (
	"encoding/json"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *UserGrantEventStore {
	return &UserGrantEventStore{
		Eventstore:     mockEs,
		userGrantCache: GetMockCache(ctrl),
		idGenerator:    GetSonyFlacke(),
	}
}

func GetMockCache(ctrl *gomock.Controller) *UserGrantCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &UserGrantCache{userGrantCache: mockCache}
}

func GetSonyFlacke() id.Generator {
	return id.SonyFlakeGenerator
}

func GetMockUserGrantByIDOK(ctrl *gomock.Controller) *UserGrantEventStore {
	user := model.UserGrant{
		UserID:    "UserID",
		ProjectID: "ProjectID",
		RoleKeys:  []string{"Key"},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockUserGrantByIDRemoved(ctrl *gomock.Controller) *UserGrantEventStore {
	user := model.UserGrant{
		UserID:    "UserID",
		ProjectID: "ProjectID",
		RoleKeys:  []string{"Key"},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantRemoved},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockUserGrantByIDNoEvents(ctrl *gomock.Controller) *UserGrantEventStore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserGrant(ctrl *gomock.Controller) *UserGrantEventStore {
	user := model.UserGrant{
		UserID:    "UserID",
		ProjectID: "ProjectID",
		RoleKeys:  []string{"Key"},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserGrantInactive(ctrl *gomock.Controller) *UserGrantEventStore {
	user := model.UserGrant{
		UserID:    "UserID",
		ProjectID: "ProjectID",
		RoleKeys:  []string{"Key"},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.UserGrantDeactivated},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserGrantNoEvents(ctrl *gomock.Controller) *UserGrantEventStore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}
