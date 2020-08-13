package eventsourcing

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/id"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *IamEventstore {
	return &IamEventstore{
		Eventstore:  mockEs,
		iamCache:    GetMockCache(ctrl),
		idGenerator: GetSonyFlacke(),
	}
}

func GetMockedEventstoreWithCrypto(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *IamEventstore {
	return &IamEventstore{
		Eventstore:   mockEs,
		iamCache:     GetMockCache(ctrl),
		idGenerator:  GetSonyFlacke(),
		secretCrypto: crypto.NewBCrypt(10),
	}
}
func GetMockCache(ctrl *gomock.Controller) *IamCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &IamCache{iamCache: mockCache}
}

func GetSonyFlacke() id.Generator {
	return id.SonyFlakeGenerator
}

func GetMockIamByIDOK(ctrl *gomock.Controller) *IamEventstore {
	data, _ := json.Marshal(model.Iam{GlobalOrgID: "GlobalOrgID"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.GlobalOrgSet, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockIamByIDNoEvents(ctrl *gomock.Controller) *IamEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIam(ctrl *gomock.Controller) *IamEventstore {
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIamWithCrypto(ctrl *gomock.Controller) *IamEventstore {
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithCrypto(ctrl, mockEs)
}

func GetMockManipulateIamWithMember(ctrl *gomock.Controller) *IamEventstore {
	memberData, _ := json.Marshal(model.IamMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIamWithOIDCIdp(ctrl *gomock.Controller) *IamEventstore {
	idpData, _ := json.Marshal(model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "Name"})
	oidcData, _ := json.Marshal(model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IdpConfigAdded, Data: idpData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OidcIdpConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIamWithLoginPolicy(ctrl *gomock.Controller) *IamEventstore {
	policyData, _ := json.Marshal(model.LoginPolicy{AllowRegister: true, AllowUsernamePassword: true, AllowExternalIdp: true})
	idpProviderData, _ := json.Marshal(model.IdpProvider{IdpConfigID: "IdpConfigID", Type: 1})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IamSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyAdded, Data: policyData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyIdpProviderAdded, Data: idpProviderData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIamNotExisting(ctrl *gomock.Controller) *IamEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
