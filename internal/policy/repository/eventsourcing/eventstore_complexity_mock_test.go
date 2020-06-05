package eventsourcing

import (
	"encoding/json"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstoreComplexity(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *PolicyEventstore {
	return &PolicyEventstore{
		Eventstore:  mockEs,
		policyCache: GetMockCacheComplexity(ctrl),
		idGenerator: GetSonyFlacke(),
	}
}

func GetMockCacheComplexity(ctrl *gomock.Controller) *PolicyCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &PolicyCache{policyCache: mockCache}
}

func GetMockGetPasswordComplexityPolicyOK(ctrl *gomock.Controller) *PolicyEventstore {
	data, _ := json.Marshal(model.PasswordComplexityPolicy{Description: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordComplexityPolicyAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockGetPasswordComplexityPolicyNoEvents(ctrl *gomock.Controller) *PolicyEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockPasswordComplexityPolicy(ctrl *gomock.Controller) *PolicyEventstore {
	data, _ := json.Marshal(model.PasswordComplexityPolicy{Description: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordComplexityPolicyAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockPasswordComplexityPolicyNoEvents(ctrl *gomock.Controller) *PolicyEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}
