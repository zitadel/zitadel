package eventsourcing

import (
	"encoding/json"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstoreAge(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *PolicyEventstore {
	return &PolicyEventstore{
		Eventstore:  mockEs,
		policyCache: GetMockCacheAge(ctrl),
	}
}

func GetMockCacheAge(ctrl *gomock.Controller) *PolicyCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &PolicyCache{policyCache: mockCache}
}

func GetMockGetPasswordAgePolicyOK(ctrl *gomock.Controller) *PolicyEventstore {
	data, _ := json.Marshal(model.PasswordAgePolicy{Description: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordAgePolicyAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreAge(ctrl, mockEs)
}

func GetMockGetPasswordAgePolicyNoEvents(ctrl *gomock.Controller) *PolicyEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreAge(ctrl, mockEs)
}

func GetMockPasswordAgePolicy(ctrl *gomock.Controller) *PolicyEventstore {
	data, _ := json.Marshal(model.PasswordAgePolicy{Description: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordAgePolicyAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreAge(ctrl, mockEs)
}

func GetMockPasswordAgePolicyNoEvents(ctrl *gomock.Controller) *PolicyEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreAge(ctrl, mockEs)
}
