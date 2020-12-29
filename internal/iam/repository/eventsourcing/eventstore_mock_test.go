package eventsourcing

import (
	"encoding/json"
	model2 "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/id"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *IAMEventstore {
	return &IAMEventstore{
		Eventstore:  mockEs,
		iamCache:    GetMockCache(ctrl),
		idGenerator: GetSonyFlacke(),
	}
}

func GetMockedEventstoreWithCrypto(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *IAMEventstore {
	return &IAMEventstore{
		Eventstore:   mockEs,
		iamCache:     GetMockCache(ctrl),
		idGenerator:  GetSonyFlacke(),
		secretCrypto: crypto.NewBCrypt(10),
	}
}
func GetMockCache(ctrl *gomock.Controller) *IAMCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &IAMCache{iamCache: mockCache}
}

func GetSonyFlacke() id.Generator {
	return id.SonyFlakeGenerator
}

func GetMockIAMByIDOK(ctrl *gomock.Controller) *IAMEventstore {
	data, _ := json.Marshal(model.IAM{GlobalOrgID: "GlobalOrgID"})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.GlobalOrgSet, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockIamByIDNoEvents(ctrl *gomock.Controller) *IAMEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAM(ctrl *gomock.Controller) *IAMEventstore {
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithCrypto(ctrl *gomock.Controller) *IAMEventstore {
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithCrypto(ctrl, mockEs)
}

func GetMockManipulateIAMWithMember(ctrl *gomock.Controller) *IAMEventstore {
	memberData, _ := json.Marshal(model.IAMMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithOIDCIdp(ctrl *gomock.Controller) *IAMEventstore {
	idpData, _ := json.Marshal(model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"})
	oidcData, _ := json.Marshal(model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IDPConfigAdded, Data: idpData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCIDPConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithLoginPolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.LoginPolicy{AllowRegister: true, AllowUsernamePassword: true, AllowExternalIdp: true})
	idpProviderData, _ := json.Marshal(model.IDPProvider{IDPConfigID: "IDPConfigID", Type: 1})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyAdded, Data: policyData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyIDPProviderAdded, Data: idpProviderData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithLoginPolicyWithMFAs(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.LoginPolicy{AllowRegister: true, AllowUsernamePassword: true, AllowExternalIdp: true})
	idpProviderData, _ := json.Marshal(model.IDPProvider{IDPConfigID: "IDPConfigID", Type: 1})
	secondFactor, _ := json.Marshal(model.MFA{MFAType: int32(model2.SecondFactorTypeOTP)})
	multiFactor, _ := json.Marshal(model.MFA{MFAType: int32(model2.MultiFactorTypeU2FWithPIN)})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyAdded, Data: policyData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyIDPProviderAdded, Data: idpProviderData},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicySecondFactorAdded, Data: secondFactor},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.LoginPolicyMultiFactorAdded, Data: multiFactor},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithPasswodComplexityPolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.PasswordComplexityPolicy{MinLength: 10})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordComplexityPolicyAdded, Data: policyData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithPasswordAgePolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordAgePolicyAdded, Data: policyData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithPasswordLockoutPolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.PasswordLockoutPolicy{MaxAttempts: 10, ShowLockOutFailures: true})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.PasswordLockoutPolicyAdded, Data: policyData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithOrgIAMPolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.OrgIAMPolicy{UserLoginMustBeDomain: true})
	events := []*es_models.Event{
		{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		{AggregateID: "AggregateID", Sequence: 1, Type: model.OrgIAMPolicyAdded, Data: policyData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMNotExisting(ctrl *gomock.Controller) *IAMEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithLabelPolicy(ctrl *gomock.Controller) *IAMEventstore {
	policyData, _ := json.Marshal(model.LabelPolicy{PrimaryColor: "000001", SecondaryColor: "FFFFF1"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.LabelPolicyAdded, Data: policyData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithMailTemplate(ctrl *gomock.Controller) *IAMEventstore {
	mailTemplate, _ := json.Marshal(model.MailTemplate{Template: []byte("<!doctype htm>")})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.MailTemplateAdded, Data: mailTemplate},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateIAMWithMailText(ctrl *gomock.Controller) *IAMEventstore {
	mailText, _ := json.Marshal(model.MailText{MailTextType: "Type", Language: "DE"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.IAMSetupStarted},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.MailTextAdded, Data: mailText},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
