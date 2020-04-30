package eventsourcing

import (
	"encoding/json"
	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
	"github.com/sony/sonyflake"
	"time"
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
		hash := crypto.NewMockHashAlgorithm(ctrl)
		hash.EXPECT().Hash(gomock.Any()).Return(nil, nil)
		hash.EXPECT().Algorithm().Return("bcrypt")
		es.PasswordAlg = hash
	}
	return es
}

func GetMockCache(ctrl *gomock.Controller) *UserCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &UserCache{userCache: mockCache}
}

func GetSonyFlacke() *sonyflake.Sonyflake {
	return sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GetMockPwGenerator(ctrl *gomock.Controller) crypto.Generator {
	alg := crypto.CreateMockEncryptionAlg(ctrl)
	generator := crypto.NewMockGenerator(ctrl)
	generator.EXPECT().Length().Return(uint(10))
	generator.EXPECT().Runes().Return([]rune("abcdefghijklmnopqrstuvwxyz"))
	generator.EXPECT().Alg().AnyTimes().Return(alg)
	generator.EXPECT().Expiry().Return(time.Hour * 1)
	return generator
}

func GetMockUserByIDOK(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockUserByIDNoEvents(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithPWGenerator(ctrl *gomock.Controller, init, email, phone, password bool) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
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
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithPw(ctrl, mockEs, init, email, phone, password)
}

func GetMockManipulateUserWithOTPGen(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	es := GetMockedEventstore(ctrl, mockEs)
	hash := crypto.NewMockEncryptionAlgorithm(ctrl)
	hash.EXPECT().Algorithm().Return("aes")
	hash.EXPECT().Encrypt(gomock.Any()).Return(nil, nil)
	hash.EXPECT().EncryptionKeyID().Return("id")
	es.Multifactors = systemdefaults.Multifactors{OTP: systemdefaults.OTP{
		Issuer:    "Issuer",
		CryptoMFA: hash,
	}}
	return es
}

func GetMockManipulateInactiveUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 2, Type: usr_model.UserDeactivated},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateLockedUser(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	data, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserLocked},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithInitCode(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	code := model.InitUserCode{Expiry: time.Hour * 30}
	dataUser, _ := json.Marshal(user)
	dataCode, _ := json.Marshal(code)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.InitializedUserCodeCreated, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithEmailCode(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
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
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserEmailCodeAdded, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, true, false, false)
}
func GetMockManipulateUserVerifiedEmail(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
		Email: &model.Email{
			EmailAddress: "EmailAddress",
		},
	}
	dataUser, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserEmailVerified},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithPhoneCode(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
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
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserPhoneCodeAdded, Data: dataCode},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithPw(ctrl, mockEs, false, false, true, false)
}

func GetMockManipulateUserVerifiedPhone(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
		Phone: &model.Phone{
			PhoneNumber: "PhoneNumber",
		},
	}
	dataUser, _ := json.Marshal(user)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserPhoneVerified},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
func GetMockManipulateUserFull(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName:  "UserName",
			FirstName: "FirstName",
			LastName:  "LastName",
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
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserWithOTP(ctrl *gomock.Controller) *UserEventstore {
	user := model.User{
		Profile: &model.Profile{
			UserName: "UserName",
		},
	}
	otp := model.OTP{Secret: &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}}
	dataUser, _ := json.Marshal(user)
	dataOtp, _ := json.Marshal(otp)
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.UserAdded, Data: dataUser},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: usr_model.MfaOtpAdded, Data: dataOtp},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateUserNoEvents(ctrl *gomock.Controller) *UserEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}
