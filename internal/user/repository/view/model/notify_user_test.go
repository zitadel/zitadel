package model

import (
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func TestNotifyUserAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		user  *NotifyUser
	}
	tests := []struct {
		name   string
		args   args
		result *NotifyUser
	}{
		{
			name: "append added user event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullHuman(nil))},
				user:  &NotifyUser{},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
		},
		{
			name: "append added human event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.HumanAdded, ResourceOwner: "OrgID", Data: mockUserData(getFullHuman(nil))},
				user:  &NotifyUser{},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
		},
		{
			name: "append change user profile event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserProfileChanged, ResourceOwner: "OrgID", Data: mockProfileData(&es_model.Profile{FirstName: "FirstNameChanged"})},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstNameChanged", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
		},
		{
			name: "append change user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailChanged, ResourceOwner: "OrgID", Data: mockEmailData(&es_model.Email{EmailAddress: "EmailChanged"})},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "EmailChanged", LastPhone: "Phone"},
		},
		{
			name: "append change user email event, existing email",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailChanged, ResourceOwner: "OrgID", Data: mockEmailData(&es_model.Email{EmailAddress: "EmailChanged"})},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", VerifiedEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "EmailChanged", VerifiedEmail: "Email", LastPhone: "Phone"},
		},
		{
			name: "append verify user email event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserEmailVerified, ResourceOwner: "OrgID"},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", VerifiedEmail: "Email", LastPhone: "Phone"},
		},
		{
			name: "append change user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneChanged, ResourceOwner: "OrgID", Data: mockPhoneData(&es_model.Phone{PhoneNumber: "PhoneChanged"})},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "PhoneChanged"},
		},
		{
			name: "append change user phone event, existing phone",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneChanged, ResourceOwner: "OrgID", Data: mockPhoneData(&es_model.Phone{PhoneNumber: "PhoneChanged"})},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone", VerifiedPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "PhoneChanged", VerifiedPhone: "Phone"},
		},
		{
			name: "append verify user phone event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserPhoneVerified, ResourceOwner: "OrgID"},
				user:  &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone"},
			},
			result: &NotifyUser{ID: "AggregateID", ResourceOwner: "OrgID", UserName: "UserName", FirstName: "FirstName", LastName: "LastName", LastEmail: "Email", LastPhone: "Phone", VerifiedPhone: "Phone"},
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
			if tt.args.user.LastEmail != tt.result.LastEmail {
				t.Errorf("got wrong result LastEmail: expected: %v, actual: %v ", tt.result.LastEmail, tt.args.user.LastEmail)
			}
			if tt.args.user.VerifiedEmail != tt.result.VerifiedEmail {
				t.Errorf("got wrong result VerifiedEmail: expected: %v, actual: %v ", tt.result.VerifiedEmail, tt.args.user.VerifiedEmail)
			}
			if tt.args.user.LastPhone != tt.result.LastPhone {
				t.Errorf("got wrong result LastPhone: expected: %v, actual: %v ", tt.result.LastPhone, tt.args.user.LastPhone)
			}
			if tt.args.user.VerifiedPhone != tt.result.VerifiedPhone {
				t.Errorf("got wrong result VerifiedPhone: expected: %v, actual: %v ", tt.result.VerifiedPhone, tt.args.user.VerifiedPhone)
			}
		})
	}
}
