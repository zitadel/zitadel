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
		{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: repo_model.OrgAggregate, Data: data},
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
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: data},
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
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IDPConfigAdded, Data: idpData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCIDPConfigAdded, Data: oidcData},
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
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyAdded, Data: loginPolicy},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyIDPProviderAdded, Data: idpData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithPasswordComplexityPolicy(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	passwordComplexityPolicy, _ := json.Marshal(iam_es_model.PasswordComplexityPolicy{
		MinLength:    10,
		HasLowercase: true,
		HasUppercase: true,
		HasSymbol:    true,
		HasNumber:    true,
	})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordComplexityPolicyAdded, Data: passwordComplexityPolicy},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithPasswordLockoutPolicy(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	passwordLockoutPolicy, _ := json.Marshal(iam_es_model.PasswordLockoutPolicy{
		MaxAttempts:         10,
		ShowLockOutFailures: true,
	})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordLockoutPolicyAdded, Data: passwordLockoutPolicy},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithPasswordAgePolicy(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	passwordAgePolicy, _ := json.Marshal(iam_es_model.PasswordAgePolicy{
		MaxAgeDays:     10,
		ExpireWarnDays: 10,
	})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordAgePolicyAdded, Data: passwordAgePolicy},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithLabelPolicy(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	labelPolicy, _ := json.Marshal(iam_es_model.LabelPolicy{PrimaryColor: "000001", SecondaryColor: "FFFFF1"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LabelPolicyAdded, Data: labelPolicy},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithMailTemplate(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	mailTemplate, _ := json.Marshal(iam_es_model.MailTemplate{Template: []byte("<!doctype htm>")})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.MailTemplateAdded, Data: mailTemplate},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgWithMailText(ctrl *gomock.Controller) *OrgEventstore {
	orgData, _ := json.Marshal(model.Org{Name: "MusterOrg"})
	mailText, _ := json.Marshal(iam_es_model.MailText{MailTextType: "Type", Language: "DE"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgAdded, Data: orgData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.MailTextAdded, Data: mailText},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
