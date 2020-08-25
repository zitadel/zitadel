package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func mockUserData(user *es_model.Human) []byte {
	data, _ := json.Marshal(user)
	return data
}

func mockPasswordData(password *es_model.Password) []byte {
	data, _ := json.Marshal(password)
	return data
}

func mockProfileData(profile *es_model.Profile) []byte {
	data, _ := json.Marshal(profile)
	return data
}

func mockEmailData(email *es_model.Email) []byte {
	data, _ := json.Marshal(email)
	return data
}

func mockPhoneData(phone *es_model.Phone) []byte {
	data, _ := json.Marshal(phone)
	return data
}

func mockAddressData(address *es_model.Address) []byte {
	data, _ := json.Marshal(address)
	return data
}

func getFullUser(password *es_model.Password) *es_model.Human {
	return &es_model.Human{
		Profile: &es_model.Profile{
			FirstName: "FirstName",
			LastName:  "LastName",
		},
		Email: &es_model.Email{
			EmailAddress: "Email",
		},
		Phone: &es_model.Phone{
			PhoneNumber: "Phone",
		},
		Address: &es_model.Address{
			Country: "Country",
		},
		Password: password,
	}
}

func TestUserAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		user  *UserView
	}
	tests := []struct {
		name   string
		args   args
		result *UserView
	}{
		{
			name: "append added user event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullUser(nil))},
				user:  &UserView{},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInitial)},
		},
		{
			name: "append added user with password event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullUser(&es_model.Password{Secret: &crypto.CryptoValue{}}))},
				user:  &UserView{},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", PasswordSet: true}, State: int32(model.UserStateInitial)},
		},
		{
			name: "append added user with password but change required event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullUser(&es_model.Password{ChangeRequired: true, Secret: &crypto.CryptoValue{}}))},
				user:  &UserView{},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", PasswordSet: true, PasswordChangeRequired: true}, State: int32(model.UserStateInitial)},
		},
		{
			name: "append password change event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPasswordChanged, ResourceOwner: "OrgID", Data: mockPasswordData(&es_model.Password{Secret: &crypto.CryptoValue{}})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", PasswordSet: true}, State: int32(model.UserStateActive)},
		},
		{
			name: "append password change with change required event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPasswordChanged, ResourceOwner: "OrgID", Data: mockPasswordData(&es_model.Password{ChangeRequired: true, Secret: &crypto.CryptoValue{}})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", PasswordSet: true, PasswordChangeRequired: true}, State: int32(model.UserStateActive)},
		},
		{
			name: "append change user profile event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserProfileChanged, ResourceOwner: "OrgID", Data: mockProfileData(&es_model.Profile{FirstName: "FirstNameChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInitial)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstNameChanged", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInitial)},
		},
		{
			name: "append change user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailChanged, ResourceOwner: "OrgID", Data: mockEmailData(&es_model.Email{EmailAddress: "EmailChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "EmailChanged", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append verify user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInitial)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append change user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneChanged, ResourceOwner: "OrgID", Data: mockPhoneData(&es_model.Phone{PhoneNumber: "PhoneChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "PhoneChanged", Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append verify user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", IsPhoneVerified: true, Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append change user address event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAddressChanged, ResourceOwner: "OrgID", Data: mockAddressData(&es_model.Address{Country: "CountryChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "CountryChanged"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append user deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserDeactivated, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInactive)},
		},
		{
			name: "append user reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserReactivated, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateInactive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append user lock event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserLocked, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateLocked)},
		},
		{
			name: "append user unlock event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserUnlocked, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateLocked)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
		},
		{
			name: "append user add otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpAdded, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", OTPState: int32(model.MfaStateNotReady)}, State: int32(model.UserStateActive)},
		},
		{
			name: "append user verify otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", OTPState: int32(model.MfaStateNotReady)}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", OTPState: int32(model.MfaStateReady)}, State: int32(model.UserStateActive)},
		},
		{
			name: "append user remove otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpRemoved, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", OTPState: int32(model.MfaStateReady)}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", OTPState: int32(model.MfaStateUnspecified)}, State: int32(model.UserStateActive)},
		},
		{
			name: "append mfa init skipped event",
			args: args{
				event: &es_models.Event{Sequence: 1, CreationDate: time.Now().UTC(), Type: es_model.MfaInitSkipped, AggregateID: "AggregateID", ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country"}, State: int32(model.UserStateActive)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", HumanView: &HumanView{FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", MfaInitSkipped: time.Now().UTC()}, State: int32(model.UserStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.AppendEvent(tt.args.event)
			if tt.args.user.ID != tt.result.ID {
				t.Errorf("got wrong result ID: expected: %v, actual: %v ", tt.result.ID, tt.args.user.ID)
			}
			if tt.args.user.FirstName != tt.result.FirstName {
				t.Errorf("got wrong result FirstName: expected: %v, actual: %v ", tt.result.FirstName, tt.args.user.FirstName)
			}
			if tt.args.user.LastName != tt.result.LastName {
				t.Errorf("got wrong result FirstName: expected: %v, actual: %v ", tt.result.FirstName, tt.args.user.FirstName)
			}
			if tt.args.user.ResourceOwner != tt.result.ResourceOwner {
				t.Errorf("got wrong result ResourceOwner: expected: %v, actual: %v ", tt.result.ResourceOwner, tt.args.user.ResourceOwner)
			}
			if tt.args.user.Email != tt.result.Email {
				t.Errorf("got wrong result email: expected: %v, actual: %v ", tt.result.Email, tt.args.user.Email)
			}
			if tt.args.user.IsEmailVerified != tt.result.IsEmailVerified {
				t.Errorf("got wrong result IsEmailVerified: expected: %v, actual: %v ", tt.result.IsEmailVerified, tt.args.user.IsEmailVerified)
			}
			if tt.args.user.Phone != tt.result.Phone {
				t.Errorf("got wrong result Phone: expected: %v, actual: %v ", tt.result.Phone, tt.args.user.Phone)
			}
			if tt.args.user.IsPhoneVerified != tt.result.IsPhoneVerified {
				t.Errorf("got wrong result IsPhoneVerified: expected: %v, actual: %v ", tt.result.IsPhoneVerified, tt.args.user.IsPhoneVerified)
			}
			if tt.args.user.Country != tt.result.Country {
				t.Errorf("got wrong result Country: expected: %v, actual: %v ", tt.result.Country, tt.args.user.Country)
			}
			if tt.args.user.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, tt.args.user.State)
			}
			if tt.args.user.OTPState != tt.result.OTPState {
				t.Errorf("got wrong result OTPState: expected: %v, actual: %v ", tt.result.OTPState, tt.args.user.OTPState)
			}
			if tt.args.user.MfaInitSkipped.Round(1*time.Second) != tt.result.MfaInitSkipped.Round(1*time.Second) {
				t.Errorf("got wrong result MfaInitSkipped: expected: %v, actual: %v ", tt.result.MfaInitSkipped.Round(1*time.Second), tt.args.user.MfaInitSkipped.Round(1*time.Second))
			}
			if tt.args.user.PasswordSet != tt.result.PasswordSet {
				t.Errorf("got wrong result PasswordSet: expected: %v, actual: %v ", tt.result.PasswordSet, tt.args.user.PasswordSet)
			}
			if tt.args.user.PasswordChangeRequired != tt.result.PasswordChangeRequired {
				t.Errorf("got wrong result PasswordChangeRequired: expected: %v, actual: %v ", tt.result.PasswordChangeRequired, tt.args.user.PasswordChangeRequired)
			}
		})
	}
}
