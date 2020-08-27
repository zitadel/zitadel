package eventsourcing

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/id"

	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	repo_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *OrgEventstore {
	return &OrgEventstore{
		Eventstore:  mockEs,
		idGenerator: GetSonyFlake(),
	}
}

func GetMockedEventstoreWithCrypto(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *OrgEventstore {
	return &OrgEventstore{
		Eventstore:   mockEs,
		idGenerator:  GetSonyFlake(),
		secretCrypto: crypto.NewBCrypt(10),
	}
}

func GetSonyFlake() id.Generator {
	return id.SonyFlakeGenerator
}
func GetMockChangesOrgOK(ctrl *gomock.Controller) *OrgEventstore {
	org := model.Org{
		Name: "MusterOrg",
	}
	data, err := json.Marshal(org)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: repo_model.OrgAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgNoEvents(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithCrypto(ctrl *gomock.Controller) *OrgEventstore {
	org := model.Org{
		Name: "MusterOrg",
	}
	data, _ := json.Marshal(org)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithCrypto(ctrl, mockEs)
}

func GetMockChangesOrgWithOIDCIdp(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	idpData, _ := json.Marshal(iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"})
	oidcData, _ := json.Marshal(iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IDPConfigAdded, Data: idpData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCIDPConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithLoginPolicy(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	loginPolicy, _ := json.Marshal(iam_es_model.LoginPolicy{AllowRegister: true, AllowExternalIdp: true, AllowUsernamePassword: true})
	idpData, _ := json.Marshal(iam_es_model.IDPProvider{IDPConfigID: "IDPConfigID", Type: int32(iam_model.IDPProviderTypeSystem)})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyAdded, Data: loginPolicy},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyIDPProviderAdded, Data: idpData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
