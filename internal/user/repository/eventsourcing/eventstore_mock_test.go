package eventsourcing

import (
	"encoding/json"
	"time"

	"github.com/golang/mock/gomock"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *UserEventstore {
	return &UserEventstore{
		Eventstore:  mockEs,
		userCache:   GetMockCache(ctrl),
		idGenerator: GetSonyFlacke(),
	}
}

func GetMockedEventstoreWithPw(ctrl *gomock.Controller, mockEs *mock.MockEventstore, init, email, phone, password bool) *UserEventstore {
	es := &UserEventstore{
		Eventstore:  mockEs,
		userCache:   GetMockCache(ctrl),
		idGenerator: GetSonyFlacke(),
	}
	if init {
		es.InitializeUserCode = GetMockPwGenerator(ctrl)
	}
	if email {
		es.EmailVerificationCode = GetMockPwGenerator(ctrl)
	}
	if phone {
		es.PhoneVerificationCode = GetMockPwGenerator(ctrl)
	}
	if password {
		es.PasswordVerificationCode = GetMockPwGenerator(ctrl)
		es.PasswordAlg = crypto.CreateMockHashAlg(ctrl)
	}
	return es
}

func GetMockCache(ctrl *gomock.Controller) *UserCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &UserCache{userCache: mockCache}
}

func GetSonyFlacke() id.Generator {
	return id.SonyFlakeGenerator
}

func GetMockPwGenerator(ctrl *gomock.Controller) crypto.Generator {
	alg := crypto.CreateMockEncryptionAlg(ctrl)
	generator := crypto.NewMockGenerator(ctrl)
	generator.EXPECT().Length().Return(uint(10)).AnyTimes()
	generator.EXPECT().Runes().Return([]rune("abcdefghijklmnopqrstuvwxyz")).AnyTimes()
	generator.EXPECT().Alg().Return(alg).AnyTimes()
	generator.EXPECT().Expiry().Return(time.Hour * 1).AnyTimes()
	return generator
}

func GetMockUserByIDOK(ctrl *gomock.Controller, user model.User) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockUserByIDNoEvents(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		UserName: "UserName",
		Human: &model.Human{
			Profile: &model.Profile{
				DisplayName: "DisplayName",
			},
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithPWGenerator(ctrl *gomock.Controller, init, email, phone, password bool) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
		Email: &model.Email{
			EmailAddress: "EmailAddress",
		},
		Phone: &model.Phone{
			PhoneNumber: "PhoneNumber",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, init, email, phone, password)
}

func GetMockManipulateUserWithInitCodeGen(ctrl *gomock.Controller, user model.User) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, true, false, false, false)
}

func GetMockManipulateUserWithPasswordInitCodeGen(ctrl *gomock.Controller, user model.User) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, true, false, false, true)
}

func GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl *gomock.Controller, user model.User) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, true, false, true)
}

func GetMockManipulateUserWithEmailCodeGen(ctrl *gomock.Controller, user model.User, verified bool) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	if verified {
		email, _ := json.Marshal(model.Email{EmailAddress: "address"})
		events = append(events,
			&es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 2, Type: model.HumanEmailVerified},
			&es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 3, Type: model.HumanEmailChanged, Data: email},
		)
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, true, false, false)
}

func GetMockManipulateUserWithPhoneCodeGen(ctrl *gomock.Controller, user model.User) *UserEventstore {
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, false, true, false)
}

func GetMockManipulateUserWithPasswordCodeGen(ctrl *gomock.Controller, user model.User, verified bool) *UserEventstore {
	data, _ := json.Marshal(user)
	code, _ := json.Marshal(user.PasswordCode)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 2, Type: model.UserPasswordCodeAdded, Data: code},
	}
	if verified {
		events = append(events, &es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 3, Type: model.HumanEmailVerified})
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, false, false, true)
}

func GetMockManipulateUserWithOTPGen(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		UserName: "UserName",
		Human: &model.Human{
			Profile: &model.Profile{
				DisplayName: "DisplayName",
			},
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	es := GetMockedEventstore(ctrl, mockEs)
	hash := crypto.NewMockEncryptionAlgorithm(ctrl)
	hash.EXPECT().Algorithm().Return("aes").AnyTimes()
	hash.EXPECT().Encrypt(gomock.Any()).Return(nil, nil).AnyTimes()
	hash.EXPECT().EncryptionKeyID().Return("id").AnyTimes()
	es.Multifactors = global_model.Multifactors{OTP: global_model.OTP{
		Issuer:    "Issuer",
		CryptoMFA: hash,
	}}
	return es
}

func GetMockManipulateInactiveUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 2, Type: model.UserDeactivated},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateLockedUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: data},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserLocked},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithInitCode(ctrl *gomock.Controller, user model.User) *UserEventstore {
	code := model.InitUserCode{Code: &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}}
	dataUser, _ := json.Marshal(user)
	dataCode, _ := json.Marshal(code)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.InitializedUserCodeAdded, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, true, false, false, true)
}

func GetMockManipulateUserWithEmailCode(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
		Email: &model.Email{
			EmailAddress: "EmailAddress",
		},
	}
	code := model.EmailCode{Code: &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}}
	dataUser, _ := json.Marshal(user)
	dataCode, _ := json.Marshal(code)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserEmailCodeAdded, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, true, false, false)
}
func GetMockManipulateUserVerifiedEmail(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
		Email: &model.Email{
			EmailAddress: "EmailAddress",
		},
	}
	dataUser, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserEmailVerified},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithPhoneCode(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
		Phone: &model.Phone{
			PhoneNumber: "PhoneNumber",
		},
	}
	code := model.PhoneCode{Code: &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}}
	dataUser, _ := json.Marshal(user)
	dataCode, _ := json.Marshal(code)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserPhoneCodeAdded, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, false, true, false)
}

func GetMockManipulateUserVerifiedPhone(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
		Phone: &model.Phone{
			PhoneNumber: "PhoneNumber",
		},
	}
	dataUser, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserPhoneVerified},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserFull(ctrl *gomock.Controller, verified bool) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
			FirstName:   "FirstName",
			LastName:    "LastName",
		},
		Password: &model.Password{
			Secret:         &crypto.CryptoValue{Algorithm: "bcrypt", KeyID: "KeyID"},
			ChangeRequired: true,
		},
		Email: &model.Email{
			EmailAddress: "EmailAddress",
		},
		Phone: &model.Phone{
			PhoneNumber: "PhoneNumber",
		},
		Address: &model.Address{
			Country: "Country",
		},
	}
	dataUser, _ := json.Marshal(user)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
	}
	if verified {
		events = append(events, &es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 2, Type: model.HumanEmailVerified})
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithOTP(ctrl *gomock.Controller, decrypt, verified bool) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
	}
	otp := model.OTP{
		Secret: &crypto.CryptoValue{
			CryptoType: crypto.TypeEncryption,
			Algorithm:  "enc",
			KeyID:      "id",
			Crypted:    []byte("code"),
		},
	}
	dataUser, _ := json.Marshal(user)
	dataOTP, _ := json.Marshal(otp)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.MFAOTPAdded, Data: dataOTP},
	}
	if verified {
		events = append(events, &es_models.Event{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.MFAOTPVerified})
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	es := GetMockedEventstore(ctrl, mockEs)
	if !decrypt {
		return es
	}
	enc := crypto.NewMockEncryptionAlgorithm(ctrl)
	enc.EXPECT().Algorithm().Return("enc").AnyTimes()
	enc.EXPECT().Encrypt(gomock.Any()).Return(nil, nil).AnyTimes()
	enc.EXPECT().EncryptionKeyID().Return("id").AnyTimes()
	enc.EXPECT().DecryptionKeyIDs().Return([]string{"id"}).AnyTimes()
	enc.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return("code", nil).AnyTimes()
	es.Multifactors = global_model.Multifactors{OTP: global_model.OTP{
		Issuer:    "Issuer",
		CryptoMFA: enc,
	}}
	return es
}

func GetMockManipulateUserWithExternalIDP(ctrl *gomock.Controller) *UserEventstore {
	user := model.Human{
		Profile: &model.Profile{
			DisplayName: "DisplayName",
		},
	}
	externalIDP := model.ExternalIDP{
		IDPConfigID: "IDPConfigID",
		UserID:      "UserID",
		DisplayName: "DisplayName",
	}
	dataUser, _ := json.Marshal(user)
	dataIDP, _ := json.Marshal(externalIDP)
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.UserAdded, Data: dataUser, ResourceOwner: "ResourceOwner"},
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, Type: model.HumanExternalIDPAdded, Data: dataIDP, ResourceOwner: "ResourceOwner"},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserNoEvents(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserNoEventsWithPw(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, false, false, true)
}

func GetMockedEventstoreComplexity(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *UserEventstore {
	return &UserEventstore{
		Eventstore: mockEs,
	}
}

func GetMockChangesUserOK(ctrl *gomock.Controller) *UserEventstore {
	user := model.Profile{
		FirstName:   "Hans",
		LastName:    "Muster",
		DisplayName: "DisplayName",
	}
	data, err := json.Marshal(user)
	if err != nil {

	}
	events := []*es_models.Event{
		{AggregateID: "AggregateID", AggregateVersion: "v1", Sequence: 1, AggregateType: model.UserAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesUserNoEvents(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}
