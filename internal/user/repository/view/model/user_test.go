package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"testing"
)

func mockUserData(user *es_model.User) []byte {
	data, _ := json.Marshal(user)
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

func getFullUser() *es_model.User {
	return &es_model.User{
		Profile: &es_model.Profile{
			UserName:  "UserName",
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
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullUser())},
				user:  &UserView{},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INITIAL)},
		},
		{
			name: "append change user profile event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserProfileChanged, ResourceOwner: "OrgID", Data: mockProfileData(&es_model.Profile{FirstName: "FirstNameChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INITIAL)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstNameChanged", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INITIAL)},
		},
		{
			name: "append change user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailChanged, ResourceOwner: "OrgID", Data: mockEmailData(&es_model.Email{EmailAddress: "EmailChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "EmailChanged", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append verify user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INITIAL)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append change user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneChanged, ResourceOwner: "OrgID", Data: mockPhoneData(&es_model.Phone{PhoneNumber: "PhoneChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "PhoneChanged", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append verify user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", IsPhoneVerified: true, Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append change user address event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAddressChanged, ResourceOwner: "OrgID", Data: mockAddressData(&es_model.Address{Country: "CountryChanged"})},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", IsEmailVerified: true, Phone: "Phone", Country: "CountryChanged", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append user deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserDeactivated, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INACTIVE)},
		},
		{
			name: "append user reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserReactivated, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_INACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append user lock event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserLocked, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_LOCKED)},
		},
		{
			name: "append user unlock event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserUnlocked, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_LOCKED)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
		},
		{
			name: "append user add otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpAdded, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE), OTPState: int32(model.MFASTATE_NOTREADY)},
		},
		{
			name: "append user verify otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpVerified, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE), OTPState: int32(model.MFASTATE_NOTREADY)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE), OTPState: int32(model.MFASTATE_READY)},
		},
		{
			name: "append user remove otp event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.MfaOtpRemoved, ResourceOwner: "OrgID"},
				user:  &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE), OTPState: int32(model.MFASTATE_READY)},
			},
			result: &UserView{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", Email: "Email", Phone: "Phone", Country: "Country", State: int32(model.USERSTATE_ACTIVE), OTPState: int32(model.MFASTATE_UNSPECIFIED)},
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
		})
	}
}
